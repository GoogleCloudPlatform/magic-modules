package cmd

import (
	"fmt"
	"magician/github"
	"os"

	"github.com/spf13/cobra"
)

var testTGCCmd = &cobra.Command{
	Use:   "test-tgc",
	Short: "Run tgc unit tests via workflow dispatch",
	Long: `This command runs tgc unit tests via workflow dispatch

	The following PR details are expected as environment variables:
        1. COMMIT_SHA
        2. PR_NUMBER
	`,
	Run: func(cmd *cobra.Command, args []string) {
		commit := os.Getenv("COMMIT_SHA")
		pr := os.Getenv("PR_NUMBER")

		gh := github.NewGithubService()

		execTestTGC(commit, pr, gh)
	},
}

func execTestTGC(commit, pr string, gh ttGithub) {
	if err := gh.CreateWorkflowDispatchEvent("test-tgc.yml", map[string]any{
		"owner":  "modular-magician",
		"repo":   "terraform-google-conversion",
		"branch": "auto-pr-" + pr,
		"sha":    commit,
	}); err != nil {
		fmt.Printf("Error creating workflow dispatch event: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(testTGCCmd)
}
