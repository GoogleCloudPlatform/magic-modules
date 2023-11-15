package cmd

import (
	"fmt"
	"magician/vcr"
	"os"

	"github.com/spf13/cobra"
)

// TODO(trodge): Move this into magician/github along with repo cloning
const githubUsername = "modular-magician"

var checkCassettesCmd = &cobra.Command{
	Use:   "check-cassettes",
	Short: "Run VCR tests on downstream main branch",
	Long: `This command runs after downstream changes are merged and runs the most recent
	VCR cassettes using the newly built beta provider.

	The following environment variables are expected:
	1. GOPATH
	2. SA_KEY
	3. GITHUB_TOKEN

	It prints a list of tests that failed in replaying mode along with all test output.`,
	Run: func(cmd *cobra.Command, args []string) {
		goPath, ok := os.LookupEnv("GOPATH")
		if !ok {
			fmt.Println("Did not provide GOPATH environment variable")
		}

		saKey := os.Getenv("SA_KEY")

		githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			fmt.Println("Did not provide GITHUB_TOKEN environment variable")
			os.Exit(1)
		}

		t, err := vcr.NewTester(goPath, saKey)
		if err != nil {
			fmt.Println("Error creating VCR tester: ", err)
		}
		execCheckCassettes(t, goPath, githubToken)
	},
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
