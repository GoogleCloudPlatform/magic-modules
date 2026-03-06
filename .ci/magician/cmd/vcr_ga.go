package cmd

import (
	"fmt"
	"magician/exec"
	"magician/provider"
	"magician/source"
	"magician/vcr"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	_ "embed"
)

var vcuRequiredEnvironmentVariables = [...]string{
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
	"GITHUB_TOKEN_CLASSIC",
}

var vcuOptionalEnvironmentVariables = [...]string{
	"GOOGLE_CHRONICLE_INSTANCE_ID",
	"GOOGLE_VMWAREENGINE_PROJECT",
}

var (
	//go:embed templates/vcr/vcr_cassettes_update_replaying.tmpl
	replayingTmplText string
	//go:embed templates/vcr/vcr_cassettes_update_recording.tmpl
	recordingTmplText string
)

type vcrCassetteUpdateReplayingResult struct {
	ReplayingResult    vcr.Result
	ReplayingErr       error
	AllReplayingPassed bool
}

type vcrCassetteUpdateRecordingResult struct {
	RecordingResult    vcr.Result
	HasTerminatedTests bool
	RecordingErr       error
	AllRecordingPassed bool
}

var vcrCassetteUpdateCmd = &cobra.Command{
	Use:   "vcr-ga-gce",
	Short: "Update VCR cassettes",
	Long: `This command is triggered in .ci/gcb-vcr-ga-gce.yml to update vcr ga gce cassettes.
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		env := make(map[string]string)
		for _, ev := range vcuRequiredEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				return fmt.Errorf("did not provide %s environment variable", ev)
			}
			env[ev] = val
		}
		for _, ev := range vcuOptionalEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if ok {
				env[ev] = val
			} else {
				fmt.Printf("🟡 Did not provide %s environment variable\n", ev)
			}
		}

		buildID := args[0]

		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating Runner: %w", err)
		}
		ctlr := source.NewController(env["GOPATH"], "hashicorp", env["GITHUB_TOKEN_CLASSIC"], rnr)

		vt, err := vcr.NewTester(env, "ci-vcr-cassettes", "", rnr, true)
		if err != nil {
			return fmt.Errorf("error creating VCR tester: %w", err)
		}

		today := time.Now().Format("2006-01-02")
		return execVCRCassetteUpdate(buildID, today, rnr, ctlr, vt)
	},
}

func execVCRCassetteUpdate(buildID, today string, rnr ExecRunner, ctlr *source.Controller, vt *vcr.Tester) error {
	if err := vt.FetchCassettes(provider.GA, "main", ""); err != nil {
		return fmt.Errorf("error fetching cassettes: %w", err)
	}

	bucketPrefix := fmt.Sprintf("gs://vcr-nightly/ga/%s/%s", today, buildID)

	providerRepo := &source.Repo{
		Name: provider.GA.RepoName(),
	}
	ctlr.SetPath(providerRepo)
	if err := ctlr.Clone(providerRepo); err != nil {
		return fmt.Errorf("error cloning provider: %w", err)
	}
	vt.SetRepoPath(provider.GA, providerRepo.Path)

	var testDirs []string
	service := "apikeys"
	servicePath := "./" + filepath.Join(provider.GA.ProviderName(), "services", service)
	testDirs = append(testDirs, servicePath)

	fmt.Println("running tests in REPLAYING mode now")
	replayingResult, replayingErr := vt.Run(vcr.RunOptions{
		Mode:     vcr.Replaying,
		Version:  provider.GA,
		TestDirs: testDirs,
	})

	// upload replay build and test logs
	buildLogPath := filepath.Join(rnr.GetCWD(), "testlogs", fmt.Sprintf("%s_test.log", vcr.Replaying.Lower()))
	if _, err := uploadLogsToGCS(buildLogPath, bucketPrefix+"/logs/replaying/", rnr); err != nil {
		fmt.Printf("Warning: error uploading replaying test log: %s\n", err)
	}

	testLogPath := vt.LogPath(vcr.Replaying, provider.GA)
	if _, err := uploadLogsToGCS(filepath.Join(testLogPath, "*"), bucketPrefix+"/logs/build-log/", rnr); err != nil {
		fmt.Printf("Warning: error uploading replaying build log: %s\n", err)
	}

	replayingData := vcrCassetteUpdateReplayingResult{
		ReplayingResult:    replayingResult,
		ReplayingErr:       replayingErr,
		AllReplayingPassed: len(replayingResult.FailedTests) == 0 && replayingErr == nil,
	}
	comment, err := formatVCRCassettesUpdateReplaying(replayingData)
	if err != nil {
		return fmt.Errorf("error formatting replaying result: %w", err)
	}
	fmt.Println(comment)

	if len(replayingResult.Panics) != 0 {
		return fmt.Errorf("provider crashed while running the VCR tests in REPLAYING mode: %v", replayingResult.Panics)
	}

	if len(replayingResult.FailedTests) != 0 {
		fmt.Println("running tests in RECORDING mode now")

		recordingResult, recordingErr := vt.RunParallel(vcr.RunOptions{
			Mode:     vcr.Recording,
			Version:  provider.GA,
			TestDirs: testDirs,
			Tests:    replayingResult.FailedTests,
		})

		// upload build and test logs first to preserve debugging logs in case
		// uploading cassettes failed because recording not work
		buildLogPath := filepath.Join(rnr.GetCWD(), "testlogs", fmt.Sprintf("%s_test.log", vcr.Recording.Lower()))
		if _, err := uploadLogsToGCS(buildLogPath, bucketPrefix+"/logs/recording/", rnr); err != nil {
			fmt.Printf("Warning: error uploading recording test log: %s\n", err)
		}

		testLogPath := vt.LogPath(vcr.Recording, provider.GA)
		if _, err := uploadLogsToGCS(filepath.Join(testLogPath, "*"), bucketPrefix+"/logs/build-log/", rnr); err != nil {
			fmt.Printf("Warning: error uploading recording build log: %s\n", err)
		}

		if len(recordingResult.PassedTests) > 0 {
			cassettesPath := vt.CassettePath(provider.GA)
			if _, err := uploadCassettesToGCS(cassettesPath+"/*", "gs://ci-vcr-cassettes/ga/fixtures/", rnr); err != nil {
				// There could be cases that the tests do not generate any cassettes.
				fmt.Printf("Warning: error uploading cassettes: %s\n", err)
			}
		} else {
			fmt.Println("No tests passed in recording mode, not uploading cassettes.")
		}

		hasTerminatedTests := (len(recordingResult.PassedTests) + len(recordingResult.FailedTests)) < len(replayingResult.FailedTests)
		allRecordingPassed := len(recordingResult.FailedTests) == 0 && !hasTerminatedTests && recordingErr == nil

		recordingData := vcrCassetteUpdateRecordingResult{
			RecordingResult:    recordingResult,
			RecordingErr:       recordingErr,
			AllRecordingPassed: allRecordingPassed,
		}
		comment, err := formatVCRCassettesUpdateRecording(recordingData)
		if err != nil {
			return fmt.Errorf("error formatting recording result: %w", err)
		}
		fmt.Println(comment)

		if len(recordingResult.Panics) != 0 {
			return fmt.Errorf("provider crashed while running the VCR tests in RECORDING mode: %v", recordingResult.Panics)
		}
	}
	return nil
}

func uploadLogsToGCS(src, dest string, rnr ExecRunner) (string, error) {
	return uploadToGCS(src, dest, []string{"-h", "Content-Type:text/plain", "-q", "cp", "-r"}, rnr)
}

func uploadCassettesToGCS(src, dest string, rnr ExecRunner) (string, error) {
	return uploadToGCS(src, dest, []string{"-m", "-q", "cp"}, rnr)
}

func uploadToGCS(src, dest string, opts []string, rnr ExecRunner) (string, error) {
	fmt.Printf("uploading from %s to %s\n", src, dest)
	args := append(opts, src, dest)
	fmt.Println("gsutil", args)
	return rnr.Run("gsutil", args, nil)
}

func formatVCRCassettesUpdateReplaying(data vcrCassetteUpdateReplayingResult) (string, error) {
	return formatComment("vcr_cassette_update_replayinging.tmpl", replayingTmplText, data)
}

func formatVCRCassettesUpdateRecording(data vcrCassetteUpdateRecordingResult) (string, error) {
	return formatComment("vcr_cassette_update_recording.tmpl", recordingTmplText, data)
}

func init() {
	rootCmd.AddCommand(vcrCassetteUpdateCmd)
}
