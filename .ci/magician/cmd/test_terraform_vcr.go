package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"magician/exec"
	"magician/github"
	"magician/provider"
	"magician/source"
	"magician/vcr"

	_ "embed"
)

var (
	//go:embed templates/vcr/post_replay.tmpl
	postReplayTmplText string
	//go:embed templates/vcr/record_replay.tmpl
	recordReplayTmplText string
)

var ttvRequiredEnvironmentVariables = [...]string{
	"GOCACHE",
	"GOPATH",
	"GOOGLE_BILLING_ACCOUNT",
	"GOOGLE_CUST_ID",
	"GOOGLE_IDENTITY_USER",
	"GOOGLE_MASTER_BILLING_ACCOUNT",
	"GOOGLE_ORG",
	"GOOGLE_ORG_2",
	"GOOGLE_ORG_DOMAIN",
	"GOOGLE_PROJECT",
	"GOOGLE_PROJECT_NUMBER",
	"GOOGLE_REGION",
	"GOOGLE_SERVICE_ACCOUNT",
	"GOOGLE_PUBLIC_AVERTISED_PREFIX_DESCRIPTION",
	"GOOGLE_ZONE",
	"HOME",
	"PATH",
	"SA_KEY",
	"USER",
}

var ttvOptionalEnvironmentVariables = [...]string{
	"GOOGLE_CHRONICLE_INSTANCE_ID",
	"GOOGLE_VMWAREENGINE_PROJECT",
}

type postReplay struct {
	RunFullVCR       bool
	AffectedServices []string
	NotRunBetaTests  []string
	NotRunGATests    []string
	ReplayingResult  vcr.Result
	ReplayingErr     error
	LogBucket        string
	Version          string
	Head             string
	BuildID          string
}

type recordReplay struct {
	RecordingResult               vcr.Result
	ReplayingAfterRecordingResult vcr.Result
	HasTerminatedTests            bool
	RecordingErr                  error
	AllRecordingPassed            bool
	LogBucket                     string
	Version                       string
	Head                          string
	BuildID                       string
	LogBaseUrl                    string
	BrowseLogBaseUrl              string
}

var testTerraformVCRCmd = &cobra.Command{
	Use:   "test-terraform-vcr",
	Short: "Run vcr tests for affected packages",
	Long: `This command runs on new pull requests to replay VCR cassettes and re-record failing cassettes.

It expects the following arguments:
	1. PR number
	2. SHA of the latest magic-modules commit
	3. Build ID
	4. Project ID where Cloud Builds are located
	5. Build step number
	
The following environment variables are required:
` + listTTVRequiredEnvironmentVariables(),
	RunE: func(cmd *cobra.Command, args []string) error {
		env := make(map[string]string)
		for _, ev := range ttvRequiredEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				return fmt.Errorf("did not provide %s environment variable", ev)
			}
			env[ev] = val
		}
		for _, ev := range ttvOptionalEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if ok {
				env[ev] = val
			} else {
				fmt.Printf("🟡 Did not provide %s environment variable\n", ev)
			}
		}

		for _, tokenName := range []string{"GITHUB_TOKEN_DOWNSTREAMS", "GITHUB_TOKEN_MAGIC_MODULES"} {
			val, ok := lookupGithubTokenOrFallback(tokenName)
			if !ok {
				return fmt.Errorf("did not provide %s or GITHUB_TOKEN environment variable", tokenName)
			}
			env[tokenName] = val
		}

		baseBranch := os.Getenv("BASE_BRANCH")
		if baseBranch == "" {
			baseBranch = "main"
		}

		gh := github.NewClient(env["GITHUB_TOKEN_MAGIC_MODULES"])
		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating a runner: %w", err)
		}
		ctlr := source.NewController(env["GOPATH"], "modular-magician", env["GITHUB_TOKEN_DOWNSTREAMS"], rnr)

		vt, err := vcr.NewTester(env, "ci-vcr-cassettes", "ci-vcr-logs", rnr)
		if err != nil {
			return fmt.Errorf("error creating VCR tester: %w", err)
		}

		if len(args) != 5 {
			return fmt.Errorf("wrong number of arguments %d, expected 5", len(args))
		}

		return execTestTerraformVCR(args[0], args[1], args[2], args[3], args[4], baseBranch, gh, rnr, ctlr, vt)
	},
}

func listTTVRequiredEnvironmentVariables() string {
	var result string
	for i, ev := range ttvRequiredEnvironmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execTestTerraformVCR(prNumber, mmCommitSha, buildID, projectID, buildStep, baseBranch string, gh GithubClient, rnr ExecRunner, ctlr *source.Controller, vt *vcr.Tester) error {
	newBranch := "auto-pr-" + prNumber
	oldBranch := newBranch + "-old"

	tpgRepo := &source.Repo{
		Name:   "terraform-provider-google",
		Owner:  "modular-magician",
		Branch: newBranch,
	}
	tpgbRepo := &source.Repo{
		Name:   "terraform-provider-google-beta",
		Owner:  "modular-magician",
		Branch: newBranch,
	}
	// Initialize repos
	for _, repo := range []*source.Repo{tpgRepo, tpgbRepo} {
		ctlr.SetPath(repo)
		if err := ctlr.Clone(repo); err != nil {
			return fmt.Errorf("error cloning repo: %w", err)
		}
		if err := ctlr.Fetch(repo, oldBranch); err != nil {
			return fmt.Errorf("failed to fetch old branch: %w", err)
		}
		changedFiles, err := ctlr.DiffNameOnly(repo, oldBranch, newBranch)
		if err != nil {
			return fmt.Errorf("failed to compute name-only diff: %w", err)
		}
		repo.ChangedFiles = changedFiles
		repo.UnifiedZeroDiff, err = ctlr.DiffUnifiedZero(repo, oldBranch, newBranch)
		if err != nil {
			return fmt.Errorf("failed to compute unified=0 diff: %w", err)
		}
	}

	vt.SetRepoPath(provider.Beta, tpgbRepo.Path)

	if err := rnr.PushDir(tpgbRepo.Path); err != nil {
		return fmt.Errorf("error changing to tpgbRepo dir: %w", err)
	}

	services, runFullVCR := modifiedPackages(tpgbRepo.ChangedFiles, provider.Beta)
	if len(services) == 0 && !runFullVCR {
		fmt.Println("Skipping tests: No go files or test fixtures changed")
		return nil
	}
	fmt.Println("Running tests: Go files or test fixtures changed")

	if err := vt.FetchCassettes(provider.Beta, baseBranch, newBranch); err != nil {
		return fmt.Errorf("error fetching cassettes: %w", err)
	}

	buildStatusTargetURL := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s;step=%s?project=%s", buildID, buildStep, projectID)
	if err := gh.PostBuildStatus(prNumber, "VCR-test", "pending", buildStatusTargetURL, mmCommitSha); err != nil {
		return fmt.Errorf("error posting pending status: %w", err)
	}

	replayingResult, testDirs, replayingErr := runReplaying(runFullVCR, provider.Beta, services, vt)
	testState := "success"
	if replayingErr != nil {
		testState = "failure"
	}

	if err := vt.UploadLogs(vcr.UploadLogsOptions{
		Head:    newBranch,
		BuildID: buildID,
		Mode:    vcr.Replaying,
		Version: provider.Beta,
	}); err != nil {
		return fmt.Errorf("error uploading replaying logs: %w", err)
	}

	if hasPanics, err := handlePanics(prNumber, buildID, buildStatusTargetURL, mmCommitSha, replayingResult, vcr.Replaying, gh); err != nil {
		return fmt.Errorf("error handling panics: %w", err)
	} else if hasPanics {
		return nil
	}

	var servicesArr []string
	for s := range services {
		servicesArr = append(servicesArr, s)
	}

	notRunBeta, notRunGa := notRunTests(tpgRepo.UnifiedZeroDiff, tpgbRepo.UnifiedZeroDiff, replayingResult)
	postReplayData := postReplay{
		RunFullVCR:       runFullVCR,
		AffectedServices: sort.StringSlice(servicesArr),
		NotRunBetaTests:  notRunBeta,
		NotRunGATests:    notRunGa,
		ReplayingResult:  subtestResult(replayingResult),
		ReplayingErr:     replayingErr,
		LogBucket:        "ci-vcr-logs",
		Version:          provider.Beta.String(),
		Head:             newBranch,
		BuildID:          buildID,
	}

	comment, err := formatPostReplay(postReplayData)
	if err != nil {
		return fmt.Errorf("error formatting post replay comment: %w", err)
	}
	if err := gh.PostComment(prNumber, comment); err != nil {
		return fmt.Errorf("error posting comment: %w", err)
	}
	if len(replayingResult.FailedTests) > 0 {
		recordingResult, recordingErr := vt.RunParallel(vcr.RunOptions{
			Mode:     vcr.Recording,
			Version:  provider.Beta,
			TestDirs: testDirs,
			Tests:    replayingResult.FailedTests,
		})
		if recordingErr != nil {
			testState = "failure"
		} else {
			testState = "success"
		}

		if err := vt.UploadCassettes(newBranch, provider.Beta); err != nil {
			return fmt.Errorf("error uploading cassettes: %w", err)
		}

		if err := vt.UploadLogs(vcr.UploadLogsOptions{
			Head:     newBranch,
			BuildID:  buildID,
			Parallel: true,
			Mode:     vcr.Recording,
			Version:  provider.Beta,
		}); err != nil {
			return fmt.Errorf("error uploading recording logs: %w", err)
		}

		if hasPanics, err := handlePanics(prNumber, buildID, buildStatusTargetURL, mmCommitSha, recordingResult, vcr.Recording, gh); err != nil {
			return fmt.Errorf("error handling panics: %w", err)
		} else if hasPanics {
			return nil
		}

		replayingAfterRecordingResult := vcr.Result{}
		var replayingAfterRecordingErr error
		if len(recordingResult.PassedTests) > 0 {
			replayingAfterRecordingResult, replayingAfterRecordingErr = vt.RunParallel(vcr.RunOptions{
				Mode:     vcr.Replaying,
				Version:  provider.Beta,
				TestDirs: testDirs,
				Tests:    recordingResult.PassedTests,
			})
			if replayingAfterRecordingErr != nil {
				testState = "failure"
			}

			if err := vt.UploadLogs(vcr.UploadLogsOptions{
				Head:           newBranch,
				BuildID:        buildID,
				AfterRecording: true,
				Parallel:       true,
				Mode:           vcr.Replaying,
				Version:        provider.Beta,
			}); err != nil {
				return fmt.Errorf("error uploading recording logs: %w", err)
			}

		}

		hasTerminatedTests := (len(recordingResult.PassedTests) + len(recordingResult.FailedTests)) < len(replayingResult.FailedTests)
		allRecordingPassed := len(recordingResult.FailedTests) == 0 && !hasTerminatedTests && recordingErr == nil

		recordReplayData := recordReplay{
			RecordingResult:               subtestResult(recordingResult),
			ReplayingAfterRecordingResult: subtestResult(replayingAfterRecordingResult),
			RecordingErr:                  recordingErr,
			HasTerminatedTests:            hasTerminatedTests,
			AllRecordingPassed:            allRecordingPassed,
			LogBucket:                     "ci-vcr-logs",
			Version:                       provider.Beta.String(),
			Head:                          newBranch,
			BuildID:                       buildID,
		}
		recordReplayComment, err := formatRecordReplay(recordReplayData)
		if err != nil {
			return fmt.Errorf("error formatting record replay comment: %w", err)
		}
		if err := gh.PostComment(prNumber, recordReplayComment); err != nil {
			return fmt.Errorf("error posting comment: %w", err)
		}
	}

	if err := gh.PostBuildStatus(prNumber, "VCR-test", testState, buildStatusTargetURL, mmCommitSha); err != nil {
		return fmt.Errorf("error posting build status: %w", err)
	}
	return nil
}

var addedTestsRegexp = regexp.MustCompile(`(?m)^\+func (TestAcc\w+)\(t \*testing.T\) {`)

func notRunTests(gaDiff, betaDiff string, result vcr.Result) ([]string, []string) {
	fmt.Println("Checking for new acceptance tests that were not run")
	addedGaTests := addedTestsRegexp.FindAllStringSubmatch(gaDiff, -1)
	addedBetaTests := addedTestsRegexp.FindAllStringSubmatch(betaDiff, -1)

	if len(addedGaTests) == 0 && len(addedBetaTests) == 0 {
		return []string{}, []string{}
	}

	// Consider tests "run" only if they passed or failed.
	runTests := map[string]struct{}{}
	for _, t := range result.PassedTests {
		runTests[t] = struct{}{}
	}
	for _, t := range result.FailedTests {
		runTests[t] = struct{}{}
	}

	notRunBeta := []string{}
	for _, t := range addedBetaTests {
		if _, ok := runTests[t[1]]; !ok {
			notRunBeta = append(notRunBeta, t[1])
		}
	}
	// Always count GA-only tests because we never run GA tests
	notRunGa := []string{}
	addedBetaTestsMap := map[string]struct{}{}
	for _, t := range addedBetaTests {
		addedBetaTestsMap[t[1]] = struct{}{}
	}
	for _, t := range addedGaTests {
		if _, ok := addedBetaTestsMap[t[1]]; !ok {
			notRunGa = append(notRunGa, t[1])
		}
	}

	sort.Strings(notRunBeta)
	sort.Strings(notRunGa)
	return notRunBeta, notRunGa
}

func subtestResult(original vcr.Result) vcr.Result {
	return vcr.Result{
		PassedTests:  excludeCompoundTests(original.PassedTests, original.PassedSubtests),
		FailedTests:  excludeCompoundTests(original.FailedTests, original.FailedSubtests),
		SkippedTests: excludeCompoundTests(original.SkippedTests, original.SkippedSubtests),
		Panics:       original.Panics,
	}
}

// Returns the name of the compound test that the given subtest belongs to.
func compoundTest(subtest string) string {
	parts := strings.Split(subtest, "__")
	if len(parts) != 2 {
		return subtest
	}
	return parts[0]
}

// Returns subtests and tests that are not compound tests.
func excludeCompoundTests(allTests, subtests []string) []string {
	res := make([]string, 0, len(allTests)+len(subtests))
	compoundTests := make(map[string]struct{}, len(subtests))
	for _, subtest := range subtests {
		if compound := compoundTest(subtest); compound != subtest {
			compoundTests[compound] = struct{}{}
			res = append(res, subtest)
		}
	}
	for _, test := range allTests {
		if _, ok := compoundTests[test]; !ok {
			res = append(res, test)
		}
	}
	sort.Strings(res)
	return res
}

func modifiedPackages(changedFiles []string, version provider.Version) (map[string]struct{}, bool) {
	var goFiles []string
	for _, line := range changedFiles {
		if strings.HasSuffix(line, ".go") || strings.Contains(line, "test-fixtures") || strings.HasSuffix(line, "go.mod") || strings.HasSuffix(line, "go.sum") {
			goFiles = append(goFiles, line)
		}
	}
	services := make(map[string]struct{})
	runFullVCR := false
	for _, file := range goFiles {
		if strings.HasPrefix(file, version.ProviderName()+"/services/") {
			fileParts := strings.Split(file, "/")
			services[fileParts[2]] = struct{}{}
		} else if file == version.ProviderName()+"/provider/provider_mmv1_resources.go" || file == version.ProviderName()+"/provider/provider_dcl_resources.go" {
			fmt.Println("ignore changes in ", file)
		} else {
			fmt.Println("run full tests ", file)
			runFullVCR = true
			break
		}
	}
	return services, runFullVCR
}

func runReplaying(runFullVCR bool, version provider.Version, services map[string]struct{}, vt *vcr.Tester) (vcr.Result, []string, error) {
	result := vcr.Result{}
	var testDirs []string
	var replayingErr error
	if runFullVCR {
		fmt.Println("runReplaying: full VCR tests")
		result, replayingErr = vt.Run(vcr.RunOptions{
			Mode:    vcr.Replaying,
			Version: version,
		})
	} else if len(services) > 0 {
		fmt.Printf("runReplaying: %d specific services: %v\n", len(services), services)
		for service := range services {
			servicePath := "./" + filepath.Join(version.ProviderName(), "services", service)
			testDirs = append(testDirs, servicePath)
			fmt.Println("run VCR tests in ", service)
			serviceResult, serviceReplayingErr := vt.Run(vcr.RunOptions{
				Mode:     vcr.Replaying,
				Version:  version,
				TestDirs: []string{servicePath},
			})
			if serviceReplayingErr != nil {
				replayingErr = serviceReplayingErr
			}
			result.PassedTests = append(result.PassedTests, serviceResult.PassedTests...)
			result.SkippedTests = append(result.SkippedTests, serviceResult.SkippedTests...)
			result.FailedTests = append(result.FailedTests, serviceResult.FailedTests...)
			result.Panics = append(result.Panics, serviceResult.Panics...)
		}
	} else {
		fmt.Println("runReplaying: no impacted services")
	}

	return result, testDirs, replayingErr
}

func handlePanics(prNumber, buildID, buildStatusTargetURL, mmCommitSha string, result vcr.Result, mode vcr.Mode, gh GithubClient) (bool, error) {
	if len(result.Panics) > 0 {
		comment := color("red", fmt.Sprintf("The provider crashed while running the VCR tests in %s mode\n", mode.Upper()))
		comment += fmt.Sprintf(`Please fix it to complete your PR.
View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-%s/artifacts/%s/build-log/%s_test.log)`, prNumber, buildID, mode.Lower())
		if err := gh.PostComment(prNumber, comment); err != nil {
			return true, fmt.Errorf("error posting comment: %v", err)
		}
		if err := gh.PostBuildStatus(prNumber, "VCR-test", "failure", buildStatusTargetURL, mmCommitSha); err != nil {
			return true, fmt.Errorf("error posting failure status: %v", err)
		}
		return true, nil
	}
	return false, nil
}

func init() {
	rootCmd.AddCommand(testTerraformVCRCmd)
}

func formatComment(fileName string, tmplText string, data any) (string, error) {
	funcs := template.FuncMap{
		"join":         strings.Join,
		"add":          func(i, j int) int { return i + j },
		"color":        color,
		"compoundTest": compoundTest,
	}
	tmpl, err := template.New(fileName).Funcs(funcs).Parse(tmplText)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse %s: %s", fileName, err))
	}
	sb := new(strings.Builder)
	err = tmpl.Execute(sb, data)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(sb.String()), nil
}

func formatPostReplay(data postReplay) (string, error) {
	return formatComment("post_replay.tmpl", postReplayTmplText, data)
}

func formatRecordReplay(data recordReplay) (string, error) {
	logBasePath := fmt.Sprintf("%s/%s/refs/heads/%s/artifacts/%s", data.LogBucket, data.Version, data.Head, data.BuildID)
	if data.BuildID == "" {
		logBasePath = fmt.Sprintf("%s/%s/refs/heads/%s", data.LogBucket, data.Version, data.Head)
	}
	data.LogBaseUrl = fmt.Sprintf("https://storage.cloud.google.com/%s", logBasePath)
	data.BrowseLogBaseUrl = fmt.Sprintf("https://console.cloud.google.com/storage/browser/%s", logBasePath)
	return formatComment("record_replay.tmpl", recordReplayTmplText, data)
}
