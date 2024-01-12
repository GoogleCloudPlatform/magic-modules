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

var environmentVariables = [...]string{
	"COMMIT_SHA",
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
	"PATH",
	"SA_KEY",
}

var checkCassettesCmd = &cobra.Command{
	Use:   "check-cassettes",
	Short: "Run VCR tests on downstream main branch",
	Long: `This command runs after downstream changes are merged and runs the most recent
	VCR cassettes using the newly built beta provider.

	The following environment variables are expected:
` + listEnvironmentVariables() + `

	It prints a list of tests that failed in replaying mode along with all test output.`,
	Run: func(cmd *cobra.Command, args []string) {
		env := make(map[string]string, len(environmentVariables))
		for _, ev := range environmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				fmt.Printf("Did not provide %s environment variable\n", ev)
				os.Exit(1)
			}
			env[ev] = val
		}

		rnr, err := exec.NewRunner()
		if err != nil {
			fmt.Println("Error creating Runner: ", err)
			os.Exit(1)
		}

		ctlr := source.NewController(env["GOPATH"], "modular-magician", env["GITHUB_TOKEN"], rnr)

		t, err := vcr.NewTester(env, rnr)
		if err != nil {
			fmt.Println("Error creating VCR tester: ", err)
			os.Exit(1)
		}
		execCheckCassettes(env["COMMIT_SHA"], t, ctlr)
	},
}

func listEnvironmentVariables() string {
	var result string
	for i, ev := range environmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execCheckCassettes(commit string, t vcr.Tester, ctlr *source.Controller) {
	if err := t.FetchCassettes(provider.Beta); err != nil {
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
	t.SetRepoPath(provider.Beta, providerRepo.Path)

	result, err := t.Run(vcr.Replaying, provider.Beta)
	if err != nil {
		fmt.Println("Error running VCR: ", err)
		os.Exit(1)
	}
	fmt.Println(len(result.FailedTests), " failed tests: ", result.FailedTests)
	// TODO(trodge) report these failures to bigquery
	fmt.Println(len(result.PassedTests), " passed tests: ", result.PassedTests)
	fmt.Println(len(result.SkippedTests), " skipped tests: ", result.SkippedTests)

	if err := t.Cleanup(); err != nil {
		fmt.Println("Error cleaning up vcr tester: ", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(checkCassettesCmd)
}
