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
	//go:embed test_terraform_vcr_test_analytics.tmpl
	testsAnalyticsTmplText string
	//go:embed test_terraform_vcr_non_exercised_tests.tmpl
	nonExercisedTestsTmplText string
	//go:embed test_terraform_vcr_with_replay_failed_tests.tmpl
	withReplayFailedTestsTmplText string
	//go:embed test_terraform_vcr_without_replay_failed_tests.tmpl
	withoutReplayFailedTestsTmplText string
	//go:embed test_terraform_vcr_record_replay.tmpl
	recordReplayTmplText string
)

var ttvEnvironmentVariables = [...]string{
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

type analytics struct {
	ReplayingResult  *vcr.Result
	RunFullVCR       bool
	AffectedServices []string
}

type nonExercisedTests struct {
	NotRunBetaTests []string
	NotRunGATests   []string
}

type withReplayFailedTests struct {
	ReplayingResult *vcr.Result
}

type withoutReplayFailedTests struct {
	ReplayingErr error
	PRNumber     string
	BuildID      string
}

type recordReplay struct {
	RecordingResult               *vcr.Result
	ReplayingAfterRecordingResult *vcr.Result
	HasTerminatedTests            bool
	RecordingErr                  error
	AllRecordingPassed            bool
	PRNumber                      string
	BuildID                       string
}

var testTerraformVCRCmd = &cobra.Command{
	Use:   "test-terraform-vcr",
	Short: "Run vcr tests for affected packages",
	Long:  `This command runs on new pull requests to replay VCR cassettes and re-record failing cassettes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env := make(map[string]string, len(ttvEnvironmentVariables))
		for _, ev := range ttvEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				return fmt.Errorf("did not provide %s environment variable", ev)
			}
			env[ev] = val
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

		vt, err := vcr.NewTester(env, rnr)
		if err != nil {
			return fmt.Errorf("error creating VCR tester: %w", err)
		}

		if len(args) != 5 {
			return fmt.Errorf("wrong number of arguments %d, expected 5", len(args))
		}

		return execTestTerraformVCR(args[0], args[1], args[2], args[3], args[4], baseBranch, gh, rnr, ctlr, vt)
	},
}

func execTestTerraformVCR(prNumber, mmCommitSha, buildID, projectID, buildStep, baseBranch string, gh GithubClient, rnr exec.ExecRunner, ctlr *source.Controller, vt *vcr.Tester) error {
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

	services, runFullVCR := modifiedPackages(tpgbRepo.ChangedFiles)
	if len(services) == 0 && !runFullVCR {
		fmt.Println("Skipping tests: No go files or test fixtures changed")
		return nil
	}
	fmt.Println("Running tests: Go files or test fixtures changed")

	if err := vt.FetchCassettes(provider.Beta, baseBranch, prNumber); err != nil {
		return fmt.Errorf("error fetching cassettes: %w", err)
	}

	buildStatusTargetURL := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s;step=%s?project=%s", buildID, buildStep, projectID)
	if err := gh.PostBuildStatus(prNumber, "VCR-test", "pending", buildStatusTargetURL, mmCommitSha); err != nil {
		return fmt.Errorf("error posting pending status: %w", err)
	}

	replayingResult, testDirs, replayingErr := runReplaying(runFullVCR, services, vt)
	testState := "success"
	if replayingErr != nil {
		testState = "failure"
	}

	if err := vt.UploadLogs("ci-vcr-logs", prNumber, buildID, false, false, vcr.Replaying, provider.Beta); err != nil {
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
	analyticsData := analytics{
		ReplayingResult:  replayingResult,
		RunFullVCR:       runFullVCR,
		AffectedServices: sort.StringSlice(servicesArr),
	}
	testsAnalyticsComment, err := formatTestsAnalytics(analyticsData)
	if err != nil {
		return fmt.Errorf("error formatting test_analytics comment: %w", err)
	}

	notRunBeta, notRunGa := notRunTests(tpgRepo.UnifiedZeroDiff, tpgbRepo.UnifiedZeroDiff, replayingResult)

	nonExercisedTestsData := nonExercisedTests{
		NotRunBetaTests: notRunBeta,
		NotRunGATests:   notRunGa,
	}
	nonExercisedTestsComment, err := formatNonExercisedTests(nonExercisedTestsData)
	if err != nil {
		return fmt.Errorf("error formatting non exercised tests comment: %w", err)
	}

	if len(replayingResult.FailedTests) > 0 {
		withReplayFailedTestsData := withReplayFailedTests{
			ReplayingResult: replayingResult,
		}
		withReplayFailedTestsComment, err := formatWithReplayFailedTests(withReplayFailedTestsData)
		if err != nil {
			return fmt.Errorf("error formatting action taken comment: %w", err)
		}

		comment := strings.Join([]string{testsAnalyticsComment, nonExercisedTestsComment, withReplayFailedTestsComment}, "\n")
		if err := gh.PostComment(prNumber, comment); err != nil {
			return fmt.Errorf("error posting comment: %w", err)
		}

		recordingResult, recordingErr := vt.RunParallel(vcr.Recording, provider.Beta, testDirs, replayingResult.FailedTests)
		if recordingErr != nil {
			testState = "failure"
		} else {
			testState = "success"
		}

		if err := vt.UploadCassettes("ci-vcr-cassettes", prNumber, provider.Beta); err != nil {
			return fmt.Errorf("error uploading cassettes: %w", err)
		}

		if err := vt.UploadLogs("ci-vcr-logs", prNumber, buildID, true, false, vcr.Recording, provider.Beta); err != nil {
			return fmt.Errorf("error uploading recording logs: %w", err)
		}

		if hasPanics, err := handlePanics(prNumber, buildID, buildStatusTargetURL, mmCommitSha, recordingResult, vcr.Recording, gh); err != nil {
			return fmt.Errorf("error handling panics: %w", err)
		} else if hasPanics {
			return nil
		}

		var replayingAfterRecordingResult *vcr.Result
		var replayingAfterRecordingErr error
		if len(recordingResult.PassedTests) > 0 {
			replayingAfterRecordingResult, replayingAfterRecordingErr = vt.RunParallel(vcr.Replaying, provider.Beta, testDirs, recordingResult.PassedTests)
			if replayingAfterRecordingErr != nil {
				testState = "failure"
			}

			if err := vt.UploadLogs("ci-vcr-logs", prNumber, buildID, true, true, vcr.Replaying, provider.Beta); err != nil {
				return fmt.Errorf("error uploading recording logs: %w", err)
			}
		}

		hasTerminatedTests := (len(recordingResult.PassedTests) + len(recordingResult.FailedTests)) < len(replayingResult.FailedTests)
		allRecordingPassed := len(recordingResult.FailedTests) == 0 && !hasTerminatedTests && recordingErr == nil

		recordReplayData := recordReplay{
			RecordingResult:               recordingResult,
			ReplayingAfterRecordingResult: replayingAfterRecordingResult,
			RecordingErr:                  recordingErr,
			HasTerminatedTests:            hasTerminatedTests,
			AllRecordingPassed:            allRecordingPassed,
			PRNumber:                      prNumber,
			BuildID:                       buildID,
		}
		recordReplayComment, err := formatRecordReplay(recordReplayData)
		if err != nil {
			return fmt.Errorf("error formatting record replay comment: %w", err)
		}
		if err := gh.PostComment(prNumber, recordReplayComment); err != nil {
			return fmt.Errorf("error posting comment: %w", err)
		}

	} else { //  len(replayingResult.FailedTests) == 0
		withoutReplayFailedTestsData := withoutReplayFailedTests{
			ReplayingErr: replayingErr,
			PRNumber:     prNumber,
			BuildID:      buildID,
		}
		withoutReplayFailedTestsComment, err := formatWithoutReplayFailedTests(withoutReplayFailedTestsData)
		if err != nil {
			return fmt.Errorf("error formatting action taken comment: %w", err)
		}

		comment := strings.Join([]string{testsAnalyticsComment, nonExercisedTestsComment, withoutReplayFailedTestsComment}, "\n")
		if err := gh.PostComment(prNumber, comment); err != nil {
			return fmt.Errorf("error posting comment: %w", err)
		}
	}

	if err := gh.PostBuildStatus(prNumber, "VCR-test", testState, buildStatusTargetURL, mmCommitSha); err != nil {
		return fmt.Errorf("error posting build status: %w", err)
	}
	return nil
}

var addedTestsRegexp = regexp.MustCompile(`(?m)^\+func (Test\w+)\(t \*testing.T\) {`)

func notRunTests(gaDiff, betaDiff string, result *vcr.Result) ([]string, []string) {
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

func modifiedPackages(changedFiles []string) (map[string]struct{}, bool) {
	var goFiles []string
	for _, line := range changedFiles {
		if strings.HasSuffix(line, ".go") || strings.Contains(line, "test-fixtures") || strings.HasSuffix(line, "go.mod") || strings.HasSuffix(line, "go.sum") {
			goFiles = append(goFiles, line)
		}
	}
	services := make(map[string]struct{})
	runFullVCR := false
	for _, file := range goFiles {
		if strings.HasPrefix(file, "google-beta/services/") {
			fileParts := strings.Split(file, "/")
			services[fileParts[2]] = struct{}{}
		} else if file == "google-beta/provider/provider_mmv1_resources.go" || file == "google-beta/provider/provider_dcl_resources.go" {
			fmt.Println("ignore changes in ", file)
		} else {
			fmt.Println("run full tests ", file)
			runFullVCR = true
			break
		}
	}
	return services, runFullVCR
}

func runReplaying(runFullVCR bool, services map[string]struct{}, vt *vcr.Tester) (*vcr.Result, []string, error) {
	var result *vcr.Result
	var testDirs []string
	var replayingErr error
	if runFullVCR {
		fmt.Println("run full VCR tests")
		result, replayingErr = vt.Run(vcr.Replaying, provider.Beta, nil)
	} else if len(services) > 0 {
		result = &vcr.Result{}
		for service := range services {
			servicePath := "./" + filepath.Join("google-beta", "services", service)
			testDirs = append(testDirs, servicePath)
			fmt.Println("run VCR tests in ", service)
			serviceResult, serviceReplayingErr := vt.Run(vcr.Replaying, provider.Beta, []string{servicePath})
			if serviceReplayingErr != nil {
				replayingErr = serviceReplayingErr
			}
			result.PassedTests = append(result.PassedTests, serviceResult.PassedTests...)
			result.SkippedTests = append(result.SkippedTests, serviceResult.SkippedTests...)
			result.FailedTests = append(result.FailedTests, serviceResult.FailedTests...)
			result.Panics = append(result.Panics, serviceResult.Panics...)
		}
	}

	return result, testDirs, replayingErr
}

func handlePanics(prNumber, buildID, buildStatusTargetURL, mmCommitSha string, result *vcr.Result, mode vcr.Mode, gh GithubClient) (bool, error) {
	if len(result.Panics) > 0 {
		comment := fmt.Sprintf(`$\textcolor{red}{\textsf{The provider crashed while running the VCR tests in %s mode}}$
$\textcolor{red}{\textsf{Please fix it to complete your PR}}$
View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-%s/artifacts/%s/build-log/%s_test.log)`, mode.Upper(), prNumber, buildID, mode.Lower())
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
		"join": strings.Join,
		"add":  func(i, j int) int { return i + j },
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

func formatTestsAnalytics(data analytics) (string, error) {
	return formatComment("test_terraform_vcr_test_analytics.tmpl", testsAnalyticsTmplText, data)
}

func formatNonExercisedTests(data nonExercisedTests) (string, error) {
	return formatComment("test_terraform_vcr_recording_mode_results.tmpl", nonExercisedTestsTmplText, data)
}

func formatWithReplayFailedTests(data withReplayFailedTests) (string, error) {
	return formatComment("test_terraform_vcr_with_replay_failed_tests.tmpl", withReplayFailedTestsTmplText, data)
}

func formatWithoutReplayFailedTests(data withoutReplayFailedTests) (string, error) {
	return formatComment("test_terraform_vcr_without_replay_failed_tests.tmpl", withoutReplayFailedTestsTmplText, data)
}

func formatRecordReplay(data recordReplay) (string, error) {
	return formatComment("test_terraform_vcr_record_replay.tmpl", recordReplayTmplText, data)
}
