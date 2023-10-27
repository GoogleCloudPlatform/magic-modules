package cmd

import (
	"fmt"
	"magician/github"
	"os"

	"github.com/spf13/cobra"
)

type ttGithub interface {
	CreateWorkflowDispatchEvent(string, map[string]any) error
}

var testTPGCmd = &cobra.Command{
	Use:   "test-tpg",
	Short: "Run provider unit tests via workflow dispatch",
	Long: `This command runs provider unit tests via workflow dispatch

	The following PR details are expected as environment variables:
        1. VERSION (beta or ga)
        2. COMMIT_SHA
        3. PR_NUMBER
	`,
	Run: func(cmd *cobra.Command, args []string) {
		version := os.Getenv("VERSION")
		commit := os.Getenv("COMMIT_SHA")
		pr := os.Getenv("PR_NUMBER")

		gh := github.NewGithubService()

		execTestTPG(version, commit, pr, gh)
	},
}

func execTestTPG(version, commit, pr string, gh ttGithub) {
	var repo string
	if version == "ga" {
		repo = "terraform-provider-google"
	} else if version == "beta" {
		repo = "terraform-provider-google-beta"
	} else {
		fmt.Println("invalid version specified")
		os.Exit(1)
	}

	if err := gh.CreateWorkflowDispatchEvent("test-tpg.yml", map[string]any{
		"owner":  "modular-magician",
		"repo":   repo,
		"branch": "auto-pr-" + pr,
		"sha":    commit,
	}); err != nil {
		fmt.Printf("Error creating workflow dispatch event: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(testTPGCmd)
}
