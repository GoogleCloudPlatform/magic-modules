package cmd

import (
	_ "embed"
	"fmt"
	"magician/exec"
	"magician/provider"
	"magician/vcr"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var tevEnvironmentVariables = [...]string{
	"GEN_PATH",
	"GOCACHE",
	"GOPATH",
	"GOOGLE_REGION",
	"GOOGLE_ZONE",
	"ORG_ID",
	"GOOGLE_PROJECT",
	"GOOGLE_BILLING_ACCOUNT",
	"GOOGLE_ORG",
	"GOOGLE_ORG_DOMAIN",
	"GOOGLE_PROJECT_NUMBER",
	"GOOGLE_USE_DEFAULT_CREDENTIALS",
	"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT",
	"KOKORO_ARTIFACTS_DIR",
	"HOME",
	"MODIFIED_FILE_PATH",
	"PATH",
	"USER",
}

var testEAPVCRCmd = &cobra.Command{
	Use:   "test-eap-vcr",
	Short: "Run vcr tests for affected packages in EAP",
	Long: `This command runs on new change lists to replay VCR cassettes and re-record failing cassettes.

It expects the following arguments:
	1. Change number


The following environment variables are required:
` + listTEVEnvironmentVariables(),
	RunE: func(cmd *cobra.Command, args []string) error {
		env := make(map[string]string, len(tevEnvironmentVariables))
		for _, ev := range tevEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				return fmt.Errorf("did not provide %s environment variable", ev)
			}
			env[ev] = val
		}
		rnr, err := exec.NewRunner()
		if err != nil {
			return err
		}
		vt, err := vcr.NewTester(env, "ci-vcr-cassettes", "ci-vcr-logs", rnr)
		if err != nil {
			return err
		}
		return execTestEAPVCR(args[0], env["GEN_PATH"], env["KOKORO_ARTIFACTS_DIR"], env["MODIFIED_FILE_PATH"], rnr, vt)
	},
}

func listTEVEnvironmentVariables() string {
	var result string
	for i, ev := range tevEnvironmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execTestEAPVCR(changeNumber, genPath, kokoroArtifactsDir, modifiedFilePath string, rnr ExecRunner, vt *vcr.Tester) error {
	vt.SetRepoPath(provider.Private, genPath)
	if err := rnr.PushDir(genPath); err != nil {
		return fmt.Errorf("error changing to gen path: %w", err)
	}

	changedFiles, err := rnr.Run("git", []string{"diff", "--name-only"}, nil)
	if err != nil {
		return fmt.Errorf("error diffing gen path: %w", err)
	}

	services, runFullVCR := modifiedPackages(strings.Split(changedFiles, "\n"), provider.Private)
	if len(services) == 0 && !runFullVCR {
		fmt.Println("Skipping tests: No go files or test fixtures changed")
		return nil
	}
	fmt.Println("Running tests: Go files or test fixtures changed")

	head := "auto-cl-" + changeNumber
	if err := vt.FetchCassettes(provider.Private, "main", head); err != nil {
		return fmt.Errorf("error fetching cassettes: %w", err)
	}
	replayingResult, testDirs, replayingErr := runReplaying(runFullVCR, provider.Private, services, vt)
	if err := vt.UploadLogs(vcr.UploadLogsOptions{
		Head:    head,
		Mode:    vcr.Replaying,
		Version: provider.Private,
	}); err != nil {
		return fmt.Errorf("error uploading replaying logs: %w", err)
	}

	if hasPanics, err := handleEAPVCRPanics(head, kokoroArtifactsDir, modifiedFilePath, replayingResult, vcr.Replaying, rnr); err != nil {
		return fmt.Errorf("error handling panics: %w", err)
	} else if hasPanics {
		return nil
	}

	var servicesArr []string
	for s := range services {
		if _, ok := allowedPrivateServices[s]; ok {
			servicesArr = append(servicesArr, s)
		}
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
	if len(replayingResult.FailedTests) > 0 {
		withReplayFailedTestsData := withReplayFailedTests{
			ReplayingResult: replayingResult,
		}

		withReplayFailedTestsComment, err := formatWithReplayFailedTests(withReplayFailedTestsData)
		if err != nil {
			return fmt.Errorf("error formatting action taken comment: %w", err)
		}
		comment := strings.Join([]string{testsAnalyticsComment, withReplayFailedTestsComment}, "\n")
		if err := postGerritComment(kokoroArtifactsDir, modifiedFilePath, comment, rnr); err != nil {
			return fmt.Errorf("error posting comment: %w", err)
		}

		recordingResult, recordingErr := vt.RunParallel(vcr.RunOptions{
			Mode:     vcr.Recording,
			Version:  provider.Beta,
			TestDirs: testDirs,
			Tests:    replayingResult.FailedTests,
		})

		if err := vt.UploadCassettes(head, provider.Private); err != nil {
			return fmt.Errorf("error uploading cassettes: %w", err)
		}

		if hasPanics, err := handleEAPVCRPanics(head, kokoroArtifactsDir, modifiedFilePath, recordingResult, vcr.Recording, rnr); err != nil {
			return fmt.Errorf("error handling panics: %w", err)
		} else if hasPanics {
			return nil
		}

		replayingAfterRecordingResult := vcr.Result{}
		if len(recordingResult.PassedTests) > 0 {
			replayingAfterRecordingResult, _ = vt.RunParallel(vcr.RunOptions{
				Mode:     vcr.Replaying,
				Version:  provider.Private,
				TestDirs: testDirs,
				Tests:    recordingResult.PassedTests,
			})
			if err := vt.UploadLogs(vcr.UploadLogsOptions{
				Head:           head,
				Parallel:       true,
				AfterRecording: true,
				Mode:           vcr.Recording,
				Version:        provider.Private,
			}); err != nil {
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
		}
		recordReplayComment, err := formatRecordReplay(recordReplayData)
		if err != nil {
			return fmt.Errorf("error formatting record replay comment: %w", err)
		}
		if err := postGerritComment(kokoroArtifactsDir, modifiedFilePath, recordReplayComment, rnr); err != nil {
			return fmt.Errorf("error posting comment: %w", err)
		}
	} else { //  len(replayingResult.FailedTests) == 0
		withoutReplayFailedTestsData := withoutReplayFailedTests{
			ReplayingErr: replayingErr,
		}
		withoutReplayFailedTestsComment, err := formatWithoutReplayFailedTests(withoutReplayFailedTestsData)
		if err != nil {
			return fmt.Errorf("error formatting action taken comment: %w", err)
		}
		comment := strings.Join([]string{testsAnalyticsComment, withoutReplayFailedTestsComment}, "\n")
		if err := postGerritComment(kokoroArtifactsDir, modifiedFilePath, comment, rnr); err != nil {
			return fmt.Errorf("error posting comment: %w", err)
		}
	}
	return nil
}

var allowedPrivateServices = map[string]struct{}{
	"accesscontextmanager":      {},
	"chronicle":                 {},
	"cloudbuildv2":              {},
	"contactcenteraiplatform":   {},
	"gkeonprem":                 {},
	"kms":                       {},
	"netapp":                    {},
	"oracledatabase":            {},
	"remotebuildexecutionadmin": {},
	"spanner":                   {},
	"storageinsights":           {},
	"vmwareengine":              {},
	"bigquery":                  {},
	"cloudbuild":                {},
	"compute":                   {},
	"dataprocgdc":               {},
	"healthcare":                {},
	"managedkafka":              {},
	"networksecurity":           {},
	"osconfig":                  {},
	"saasmanagement":            {},
	"stackdriver":               {},
	"tpuv2":                     {},
	"workstations":              {},
}

func handleEAPVCRPanics(head, kokoroArtifactsDir, modifiedFilePath string, result vcr.Result, mode vcr.Mode, rnr ExecRunner) (bool, error) {
	if len(result.Panics) > 0 {
		comment := fmt.Sprintf(`The provider crashed while running the VCR tests in %s mode.
Please fix it to complete your CL
View the [build log](https://storage.cloud.google.com/ci-vcr-logs/%s/refs/heads/%s/build-log/%s_test.log)`,
			provider.Private.String(), mode.Upper(), head, mode.Lower())
		if err := postGerritComment(kokoroArtifactsDir, modifiedFilePath, comment, rnr); err != nil {
			return true, fmt.Errorf("error posting comment: %v", err)
		}
		return true, nil
	}
	return false, nil
}

func postGerritComment(kokoroArtifactsDir, modifiedFilePath, comment string, rnr ExecRunner) error {
	return rnr.AppendFile(filepath.Join(kokoroArtifactsDir, "gerrit_comments.json"), fmt.Sprintf("\n{path: \"%s\", message: \"%s\"}", modifiedFilePath, comment))
}

func init() {
	rootCmd.AddCommand(testEAPVCRCmd)
}
