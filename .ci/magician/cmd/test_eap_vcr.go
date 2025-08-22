package cmd

import (
	_ "embed"
	"encoding/json"
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

var (
	//go:embed templates/vcr/post_replay_eap.tmpl
	postReplayEAPTmplText string
)

var tevRequiredEnvironmentVariables = [...]string{
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
	"PATH",
	"USER",
}

var tevOptionalEnvironmentVariables = [...]string{
	"GOOGLE_CHRONICLE_INSTANCE_ID",
	"GOOGLE_CUST_ID",
	"GOOGLE_IDENTITY_USER",
	"GOOGLE_MASTER_BILLING_ACCOUNT",
	"GOOGLE_ORG_2",
	"GOOGLE_PUBLIC_AVERTISED_PREFIX_DESCRIPTION",
	"GOOGLE_SERVICE_ACCOUNT",
	"GOOGLE_VMWAREENGINE_PROJECT",
}

// GerritComment is a single inline comment for a Gerrit CL.
// See go/kokoro-gob-scm#gerrit-inline-comments.
type GerritComment struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// GerritCommenter is used to add comments to a Gerrit CL.
type GerritCommenter struct {
	gerritCommentsFilename string
	rnr                    ExecRunner
	comments               []GerritComment
}

func NewGerritCommenter(gerritCommentsFilename string, rnr ExecRunner) *GerritCommenter {
	return &GerritCommenter{
		gerritCommentsFilename: gerritCommentsFilename,
		rnr:                    rnr,
	}
}

// Add adds a comment to the gerrit_comments_file json file. If a path is not
// specified, the comment is added at the patchset level, just like other
// kokoro messages.
func (g *GerritCommenter) Add(c GerritComment) error {
	if c.Path == "" {
		c.Path = "/PATCHSET_LEVEL"
	}
	g.comments = append(g.comments, c)
	b, err := json.Marshal(g.comments)
	if err != nil {
		return err
	}
	return g.rnr.WriteFile(g.gerritCommentsFilename, string(b))
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
		env := make(map[string]string, len(tevRequiredEnvironmentVariables))
		for _, ev := range tevRequiredEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				return fmt.Errorf("did not provide %s environment variable", ev)
			}
			env[ev] = val
		}
		for _, ev := range tevOptionalEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if ok {
				env[ev] = val
			} else {
				fmt.Printf("🟡 Did not provide %s environment variable\n", ev)
			}
		}

		rnr, err := exec.NewRunner()
		if err != nil {
			return err
		}
		vt, err := vcr.NewTester(env, "ci-vcr-cassettes", "ci-vcr-logs", rnr)
		if err != nil {
			return err
		}

		if len(args) != 1 {
			return fmt.Errorf("wrong number of arguments %d, expected 1", len(args))
		}

		return execTestEAPVCR(args[0], env["GEN_PATH"], env["KOKORO_ARTIFACTS_DIR"], rnr, vt)
	},
}

func listTEVEnvironmentVariables() string {
	var result string
	for i, ev := range tevRequiredEnvironmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execTestEAPVCR(changeNumber, genPath, kokoroArtifactsDir string, rnr ExecRunner, vt *vcr.Tester) error {
	vt.SetRepoPath(provider.Private, genPath)
	if err := rnr.PushDir(genPath); err != nil {
		return fmt.Errorf("error changing to gen path: %w", err)
	}

	changedFiles, err := rnr.ReadFile("diff.log")
	if err != nil {
		return fmt.Errorf("error reading diff log: %w", err)
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

	// Comments for VCR must go in the gerrit_comments_acctest.json json file.
	commenter := NewGerritCommenter(filepath.Join(kokoroArtifactsDir, "gerrit_comments_acctest.json"), rnr)

	if hasPanics, err := handleEAPVCRPanics(head, replayingResult, vcr.Replaying, commenter); err != nil {
		return fmt.Errorf("error handling panics: %w", err)
	} else if hasPanics {
		return nil
	}

	var servicesArr []string
	for s := range services {
		servicesArr = append(servicesArr, s)
	}
	postReplayData := postReplay{
		RunFullVCR:       runFullVCR,
		AffectedServices: sort.StringSlice(servicesArr),
		ReplayingResult:  replayingResult,
		ReplayingErr:     replayingErr,
		LogBucket:        "ci-vcr-logs",
		Version:          provider.Private.String(),
		Head:             head,
	}
	comment, err := formatPostReplayEAP(postReplayData)
	if err != nil {
		return fmt.Errorf("error formatting post replay comment: %w", err)
	}
	c := GerritComment{
		Message: comment,
	}
	if err := commenter.Add(c); err != nil {
		return fmt.Errorf("error adding comment: %w", err)
	}
	if len(replayingResult.FailedTests) > 0 {
		recordingResult, recordingErr := vt.RunParallel(vcr.RunOptions{
			Mode:     vcr.Recording,
			Version:  provider.Private,
			TestDirs: testDirs,
			Tests:    replayingResult.FailedTests,
		})

		if recordingErr != nil {
			fmt.Println("error during recording:", recordingErr)
		}

		if err := vt.UploadCassettes(head, provider.Private); err != nil {
			return fmt.Errorf("error uploading cassettes: %w", err)
		}

		if err := vt.UploadLogs(vcr.UploadLogsOptions{
			Head:     head,
			Parallel: true,
			Mode:     vcr.Recording,
			Version:  provider.Private,
		}); err != nil {
			return fmt.Errorf("error uploading recording logs: %w", err)
		}

		if hasPanics, err := handleEAPVCRPanics(head, recordingResult, vcr.Recording, commenter); err != nil {
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
				Mode:           vcr.Replaying,
				Version:        provider.Private,
			}); err != nil {
				return fmt.Errorf("error uploading replaying after recording logs: %w", err)
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
			LogBucket:                     "ci-vcr-logs",
			Version:                       provider.Private.String(),
			Head:                          head,
		}
		recordReplayComment, err := formatRecordReplay(recordReplayData)
		if err != nil {
			return fmt.Errorf("error formatting record replay comment: %w", err)
		}
		c = GerritComment{
			Message: recordReplayComment,
		}
		if err := commenter.Add(c); err != nil {
			return fmt.Errorf("error adding comment: %w", err)
		}
	}
	return nil
}

func handleEAPVCRPanics(head string, result vcr.Result, mode vcr.Mode, commenter *GerritCommenter) (bool, error) {
	if len(result.Panics) > 0 {
		c := GerritComment{
			Message: fmt.Sprintf(`The provider crashed while running the VCR tests in %s mode.
Please fix it to complete your CL
View the [build log](https://storage.cloud.google.com/ci-vcr-logs/%s/refs/heads/%s/build-log/%s_test.log)`,
				provider.Private.String(), mode.Upper(), head, mode.Lower()),
		}
		return true, commenter.Add(c)
	}
	return false, nil
}

func init() {
	rootCmd.AddCommand(testEAPVCRCmd)
}

func formatPostReplayEAP(data postReplay) (string, error) {
	return formatComment("post_replay_eap.tmpl", postReplayEAPTmplText, data)
}
