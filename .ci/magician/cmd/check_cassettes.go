package cmd

import (
	"fmt"
	"magician/vcr"
	"os"

	"github.com/spf13/cobra"
)

// TODO(trodge): Move this into magician/github along with repo cloning
const githubUsername = "modular-magician"

var environmentVariables = [...]string{
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
	"GOOGLE_ZONE",
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

		t, err := vcr.NewTester(env)
		if err != nil {
			fmt.Println("Error creating VCR tester: ", err)
			os.Exit(1)
		}
		execCheckCassettes(t, env["GOPATH"], env["GITHUB_TOKEN"])
	},
}

func listEnvironmentVariables() string {
	var result string
	for i, ev := range environmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execCheckCassettes(t vcr.Tester, goPath, githubToken string) {
	if err := t.FetchCassettes(vcr.Beta); err != nil {
		fmt.Println("Error fetching cassettes: ", err)
		os.Exit(1)
	}

	if err := t.CloneProvider(goPath, githubUsername, githubToken, vcr.Beta); err != nil {
		fmt.Println("Error cloning provider: ", err)
		os.Exit(1)
	}

	result, err := t.Run(vcr.Replaying, vcr.Beta)
	if err != nil {
		fmt.Println("Error running VCR: ", err)
		os.Exit(1)
	}
	fmt.Println("Failing tests: ", result.FailedTests)
	// TODO(trodge) report these failures to bigquery
	fmt.Println("Passing tests: ", result.PassedTests)
	fmt.Println("Skipping tests: ", result.SkippedTests)

	if err := t.Cleanup(); err != nil {
		fmt.Println("Error cleaning up vcr tester: ", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(checkCassettesCmd)
}
