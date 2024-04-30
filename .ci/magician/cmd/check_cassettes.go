package cmd

import (
	"fmt"
	"magician/exec"
	"magician/provider"
	"magician/source"
	"magician/vcr"
	"os"

	"github.com/spf13/cobra"
)

var ccEnvironmentVariables = [...]string{
	"COMMIT_SHA",
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
	"GOOGLE_ZONE",
	"PATH",
	"SA_KEY",
}

var checkCassettesCmd = &cobra.Command{
	Use:   "check-cassettes",
	Short: "Run VCR tests on downstream main branch",
	Long: `This command runs after downstream changes are merged and runs the most recent
	VCR cassettes using the newly built beta provider.

	The following environment variables are expected:
` + listCCEnvironmentVariables() + `

	It prints a list of tests that failed in replaying mode along with all test output.`,
	Run: func(cmd *cobra.Command, args []string) {
		env := make(map[string]string, len(ccEnvironmentVariables))
		for _, ev := range ccEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				fmt.Printf("Did not provide %s environment variable\n", ev)
				os.Exit(1)
			}
			env[ev] = val
		}

		githubToken, ok := lookupGithubTokenOrFallback("GITHUB_TOKEN_DOWNSTREAMS")
		if !ok {
			fmt.Println("Did not provide GITHUB_TOKEN_DOWNSTREAMS or GITHUB_TOKEN environment variables")
			os.Exit(1)
		}

		rnr, err := exec.NewRunner()
		if err != nil {
			fmt.Println("Error creating Runner: ", err)
			os.Exit(1)
		}

		ctlr := source.NewController(env["GOPATH"], "modular-magician", githubToken, rnr)

		vt, err := vcr.NewTester(env, rnr)
		if err != nil {
			fmt.Println("Error creating VCR tester: ", err)
			os.Exit(1)
		}
		execCheckCassettes(env["COMMIT_SHA"], vt, ctlr)
	},
}

func lookupGithubTokenOrFallback(tokenName string) (string, bool) {
	val, ok := os.LookupEnv(tokenName)
	if !ok {
		return os.LookupEnv("GITHUB_TOKEN")
	}
	return val, ok
}

func listCCEnvironmentVariables() string {
	var result string
	for i, ev := range ccEnvironmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execCheckCassettes(commit string, vt *vcr.Tester, ctlr *source.Controller) {
	if err := vt.FetchCassettes(provider.Beta, "main", ""); err != nil {
		fmt.Println("Error fetching cassettes: ", err)
		os.Exit(1)
	}

	providerRepo := &source.Repo{
		Name:   provider.Beta.RepoName(),
		Branch: "downstream-pr-" + commit,
	}
	ctlr.SetPath(providerRepo)
	if err := ctlr.Clone(providerRepo); err != nil {
		fmt.Println("Error cloning provider: ", err)
		os.Exit(1)
	}
	vt.SetRepoPath(provider.Beta, providerRepo.Path)

	result, err := vt.Run(vcr.Replaying, provider.Beta, nil)
	if err != nil {
		fmt.Println("Error running VCR: ", err)
	}
	if err := vt.UploadLogs("vcr-check-cassettes", "", "", false, false, vcr.Replaying, provider.Beta); err != nil {
		fmt.Println("Error uploading logs: ", err)
		os.Exit(1)
	}
	fmt.Println(len(result.FailedTests), " failed tests: ", result.FailedTests)
	// TODO(trodge) report these failures to bigquery
	fmt.Println(len(result.PassedTests), " passed tests: ", result.PassedTests)
	fmt.Println(len(result.SkippedTests), " skipped tests: ", result.SkippedTests)

	if err := vt.Cleanup(); err != nil {
		fmt.Println("Error cleaning up vcr tester: ", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(checkCassettesCmd)
}
