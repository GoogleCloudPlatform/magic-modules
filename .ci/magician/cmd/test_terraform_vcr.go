package cmd

import (
	"fmt"
	"magician/exec"
	"magician/github"
	"magician/provider"
	"magician/source"
	"magician/vcr"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var ttvEnvironmentVariables = [...]string{
	"GITHUB_TOKEN",
	"GOCACHE",
	"GOPATH",
	"GOOGLE_BILLING_ACCOUNT",
	"GOOGLE_CUST_ID",
	"GOOGLE_FIRESTORE_PROJECT",
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
	"GOOGLE_TPU_V2_VM_RUNTIME_VERSION",
	"GOOGLE_ZONE",
	"HOME",
	"PATH",
	"SA_KEY",
	"USER",
}

var testTerraformVCRCmd = &cobra.Command{
	Use:   "test-terraform-vcr",
	Short: "Run vcr tests for affected packages",
	Long:  `This command runs on new pull requests to replay VCR cassettes and re-record failing cassettes.`,
	Run: func(cmd *cobra.Command, args []string) {
		env := make(map[string]string, len(ttvEnvironmentVariables))
		for _, ev := range ttvEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				fmt.Printf("Did not provide %s environment variable\n", ev)
				os.Exit(1)
			}
			env[ev] = val
		}

		baseBranch := os.Getenv("BASE_BRANCH")
		if baseBranch == "" {
			baseBranch = "main"
		}

		gh := github.NewClient()
		rnr, err := exec.NewRunner()
		if err != nil {
			fmt.Println("Error creating a runner: ", err)
			os.Exit(1)
		}
		ctlr := source.NewController(env["GOPATH"], "modular-magician", env["GITHUB_TOKEN"], rnr)

		vt, err := vcr.NewTester(env, rnr)
		if err != nil {
			fmt.Println("Error creating VCR tester: ", err)
		}

		if len(args) != 5 {
			fmt.Printf("Wrong number of arguments %d, expected 5\n", len(args))
			os.Exit(1)
		}

		execTestTerraformVCR(args[0], args[1], args[2], args[3], args[4], baseBranch, gh, rnr, ctlr, vt)
	},
}

func execTestTerraformVCR(prNumber, mmCommitSha, buildID, projectID, buildStep, baseBranch string, gh GithubClient, rnr ExecRunner, ctlr *source.Controller, vt *vcr.Tester) {
	newBranch := "auto-pr-" + prNumber
	oldBranch := newBranch + "-old"
	repo := &source.Repo{
		Name:   "terraform-provider-google-beta",
		Owner:  "modular-magician",
		Branch: newBranch,
	}
	ctlr.SetPath(repo)
	if err := ctlr.Clone(repo); err != nil {
		fmt.Println("Error cloning repo: ", err)
		os.Exit(1)
	}

	vt.SetRepoPath(provider.Beta, repo.Path)

	if err := rnr.PushDir(repo.Path); err != nil {
		fmt.Println("Error changing to repo dir: ", err)
		os.Exit(1)
	}

	services, runFullVCR, err := modifiedPackages(newBranch, oldBranch, rnr)
	if err != nil {
		fmt.Println("Error getting modified packages: ", err)
		os.Exit(1)
	}

	if len(services) == 0 && !runFullVCR {
		fmt.Println("Skipping tests: No go files changed")
		os.Exit(0)
	}
	fmt.Println("Running tests: Go files changed")

	if err := vt.FetchCassettes(provider.Beta, baseBranch, prNumber); err != nil {
		fmt.Println("Error fetching cassettes: ", err)
		os.Exit(1)
	}

	buildStatusTargetURL := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s;step=%s?project=%s", buildID, buildStep, projectID)
	if err := gh.PostBuildStatus(prNumber, "VCR-test", "pending", buildStatusTargetURL, mmCommitSha); err != nil {
		fmt.Println("Error posting pending status: ", err)
		os.Exit(1)
	}

	replayingResult, affectedServicesComment, testDirs, replayingErr := runReplaying(runFullVCR, services, vt)
	testState := "success"
	if replayingErr != nil {
		testState = "failure"
	}

	if err := vt.UploadLogs("ci-vcr-logs", prNumber, buildID, false, false, vcr.Replaying, provider.Beta); err != nil {
		fmt.Println("Error uploading replaying logs: ", err)
		os.Exit(1)
	}

	if hasPanics, err := handlePanics(prNumber, buildID, buildStatusTargetURL, mmCommitSha, replayingResult, vcr.Replaying, gh); err != nil {
		fmt.Println("Error handling panics: ", err)
		os.Exit(1)
	} else if hasPanics {
		os.Exit(0)
	}

	failedTestsPattern := strings.Join(replayingResult.FailedTests, "|")

	comment := `#### Tests analytics
Total tests: ` + fmt.Sprintf("`%d`", len(replayingResult.PassedTests)+len(replayingResult.SkippedTests)+len(replayingResult.FailedTests)) + `
Passed tests: ` + fmt.Sprintf("`%d`", len(replayingResult.PassedTests)) + `
Skipped tests: ` + fmt.Sprintf("`%d`", len(replayingResult.SkippedTests)) + `
Affected tests: ` + fmt.Sprintf("`%d`", len(replayingResult.FailedTests)) + `

<details><summary>Click here to see the affected service packages</summary><blockquote>` + affectedServicesComment + `</blockquote></details>

`
	if len(replayingResult.FailedTests) > 0 {
		comment += fmt.Sprintf(`#### Action taken
<details> <summary>Found %d affected test(s) by replaying old test recordings. Starting RECORDING based on the most recent commit. Click here to see the affected tests</summary><blockquote>%s </blockquote></details>

[Get to know how VCR tests work](https://googlecloudplatform.github.io/magic-modules/docs/getting-started/contributing/#general-contributing-steps)`, len(replayingResult.FailedTests), failedTestsPattern)

		if err := gh.PostComment(prNumber, comment); err != nil {
			fmt.Println("Error posting comment: ", err)
			os.Exit(1)
		}

		recordingResult, recordingErr := vt.RunParallel(vcr.Recording, provider.Beta, testDirs, replayingResult.FailedTests)
		if recordingErr != nil {
			testState = "failure"
		} else {
			testState = "success"
		}

		if err := vt.UploadCassettes("ci-vcr-cassettes", prNumber, provider.Beta); err != nil {
			fmt.Println("Error uploading cassettes: ", err)
			os.Exit(1)
		}

		if err := vt.UploadLogs("ci-vcr-logs", prNumber, buildID, true, false, vcr.Recording, provider.Beta); err != nil {
			fmt.Println("Error uploading recording logs: ", err)
			os.Exit(1)
		}

		if hasPanics, err := handlePanics(prNumber, buildID, buildStatusTargetURL, mmCommitSha, recordingResult, vcr.Recording, gh); err != nil {
			fmt.Println("Error handling panics: ", err)
			os.Exit(1)
		} else if hasPanics {
			os.Exit(0)
		}

		comment = ""
		if len(recordingResult.PassedTests) > 0 {
			comment += "$\\textcolor{green}{\\textsf{Tests passed during RECORDING mode:}}$\n"
			for _, passedTest := range recordingResult.PassedTests {
				comment += fmt.Sprintf("`%s`[[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-%s/artifacts/%s/recording/%s.log)]\n", passedTest, prNumber, buildID, passedTest)
			}
			comment += "\n\n"

			replayingAfterRecordingResult, replayingAfterRecordingErr := vt.RunParallel(vcr.Replaying, provider.Beta, testDirs, recordingResult.PassedTests)
			if replayingAfterRecordingErr != nil {
				testState = "failure"
			}

			if err := vt.UploadLogs("ci-vcr-logs", prNumber, buildID, true, true, vcr.Replaying, provider.Beta); err != nil {
				fmt.Println("Error uploading recording logs: ", err)
				os.Exit(1)
			}

			if len(replayingAfterRecordingResult.FailedTests) > 0 {
				comment += "$\\textcolor{red}{\\textsf{Tests failed when rerunning REPLAYING mode:}}$\n"
				for _, failedTest := range replayingAfterRecordingResult.FailedTests {
					comment += fmt.Sprintf("`%s`[[Error message](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-%s/artifacts/%s/build-log/replaying_build_after_recording/%s_replaying_test.log)] [[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-%s/artifacts/%s/replaying_after_recording/%s.log)]\n", failedTest, prNumber, buildID, failedTest, prNumber, buildID, failedTest)
				}
				comment += "\n\n"
				comment += `Tests failed due to non-determinism or randomness when the VCR replayed the response after the HTTP request was made.

Please fix these to complete your PR. If you believe these test failures to be incorrect or unrelated to your change, or if you have any questions, please raise the concern with your reviewer.
`
			} else {
				comment += "$\\textcolor{green}{\\textsf{No issues found for passed tests after REPLAYING rerun.}}$\n"
			}
			comment += "\n---\n"

		}

		if len(recordingResult.FailedTests) > 0 {
			comment += "$\\textcolor{red}{\\textsf{Tests failed during RECORDING mode:}}$\n"
			for _, failedTest := range recordingResult.FailedTests {
				comment += fmt.Sprintf("`%s`[[Error message](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-%s/artifacts/%s/build-log/recording_build/%s_recording_test.log)] [[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-%s/artifacts/%s/recording/%s.log)]\n", failedTest, prNumber, buildID, failedTest, prNumber, buildID, failedTest)
			}
			comment += "\n\n"
			if len(recordingResult.PassedTests)+len(recordingResult.FailedTests) < len(replayingResult.FailedTests) {
				comment += "$\\textcolor{red}{\\textsf{Several tests got terminated during RECORDING mode.}}$\n"
			}
			comment += "$\\textcolor{red}{\\textsf{Please fix these to complete your PR.}}$\n"
		} else {
			if len(recordingResult.PassedTests)+len(recordingResult.FailedTests) < len(replayingResult.FailedTests) {
				comment += "$\\textcolor{red}{\\textsf{Several tests got terminated during RECORDING mode.}}$\n"
			} else if recordingErr != nil {
				// Check for any uncaught errors in RECORDING mode.
				comment += "$\\textcolor{red}{\\textsf{Errors occurred during RECORDING mode. Please fix them to complete your PR.}}$\n"
			} else {
				comment += "$\\textcolor{green}{\\textsf{All tests passed!}}$\n"
			}
		}

		comment += fmt.Sprintf("View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-%s/artifacts/%s/build-log/recording_test.log) or the [debug log](https://console.cloud.google.com/storage/browser/ci-vcr-logs/beta/refs/heads/auto-pr-%s/artifacts/%s/recording) for each test", prNumber, buildID, prNumber, buildID)

		if err := gh.PostComment(prNumber, comment); err != nil {
			fmt.Println("Error posting comment: ", err)
			os.Exit(1)
		}
	} else {
		if replayingErr != nil {
			// Check for any uncaught errors in REPLAYING mode.
			comment += "$\\textcolor{red}{\\textsf{Errors occurred during RECORDING mode. Please fix them to complete your PR.}}$\n"
		} else {
			comment += "$\\textcolor{green}{\\textsf{All tests passed!}}$\n"
		}
		comment += fmt.Sprintf("View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-%s/artifacts/%s/build-log/replaying_test.log)", prNumber, buildID)

		if err := gh.PostComment(prNumber, comment); err != nil {
			fmt.Println("Error posting comment: ", err)
			os.Exit(1)
		}
	}

	if err := gh.PostBuildStatus(prNumber, "VCR-test", testState, buildStatusTargetURL, mmCommitSha); err != nil {
		fmt.Println("Error posting build status: ", err)
		os.Exit(1)
	}
}

func modifiedPackages(newBranch, oldBranch string, rnr ExecRunner) (map[string]struct{}, bool, error) {
	fmt.Println("Checking for modified go files")

	if _, err := rnr.Run("git", []string{"fetch", "origin", fmt.Sprintf("%s:%s", oldBranch, oldBranch), "--depth", "1"}, nil); err != nil {

	}
	diffs, err := rnr.Run("git", []string{"diff", newBranch, oldBranch, "--name-only"}, nil)
	if err != nil {
		return nil, false, err
	}
	var goFiles []string
	for _, line := range strings.Split(diffs, "\n") {
		if strings.HasSuffix(line, ".go") || line == "go.mod" || line == "go.sum" {
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
	return services, runFullVCR, nil
}

func runReplaying(runFullVCR bool, services map[string]struct{}, vt *vcr.Tester) (*vcr.Result, string, []string, error) {
	var result *vcr.Result
	affectedServicesComment := "None"
	var testDirs []string
	var replayingErr error
	if runFullVCR {
		fmt.Println("run full VCR tests")
		affectedServicesComment = "all service packages are affected"
		result, replayingErr = vt.Run(vcr.Replaying, provider.Beta, nil)
	} else if len(services) > 0 {
		affectedServicesComment = "<ul>"
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
			affectedServicesComment += fmt.Sprintf("<li>%s</li>", service)
		}
		affectedServicesComment += "</ul>"
	}

	return result, affectedServicesComment, testDirs, replayingErr
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
