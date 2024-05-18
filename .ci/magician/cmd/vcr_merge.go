package cmd

import (
	"fmt"
	"magician/exec"
	"magician/github"
	"magician/source"
	"os"

	"github.com/spf13/cobra"
)

var vcrMergeCmd = &cobra.Command{
	Use:   "vcr-merge",
	Short: "Merge VCR cassettes",
	Long: `This command is triggered in .ci/gcb-push-downstream.yml to merge vcr cassettes.

	The command expects the following as arguments:
	1. Reference commit SHA

	It then performs the following operations:
	1. Get the latest closed PR matching the reference commit SHA.
	2. Run gsutil to list, copy, and remove the vcr cassettes fixtures.
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reference := args[0]
		fmt.Println("Reference commit SHA: ", reference)

		githubToken, ok := os.LookupEnv("GITHUB_TOKEN_CLASSIC")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN_CLASSIC environment variable")
		}

		baseBranch := os.Getenv("BASE_BRANCH")
		if baseBranch == "" {
			return fmt.Errorf("environment variable BASE_BRANCH is empty")
		}

		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating Runner: %w", err)
		}

		gh := github.NewClient(githubToken)
		return execVCRMerge(gh, reference, baseBranch, rnr)
	},
}

func execVCRMerge(gh GithubClient, sha string, baseBranch string, runner source.Runner) error {
	arr, err := gh.GetPullRequests("closed", baseBranch, "updated", "desc")
	if err != nil {
		return fmt.Errorf("error getting pull requests: %w", err)
	}
	pr := findPRBySHA(sha, arr)
	if pr == nil {
		fmt.Printf("Not finding any matching PR with commit SHA %s\n", sha)
		return nil
	}

	mergeCassettes("gs://ci-vcr-cassettes", baseBranch, fmt.Sprintf("refs/heads/auto-pr-%d", pr.Number), runner)
	mergeCassettes("gs://ci-vcr-cassettes/beta", baseBranch, fmt.Sprintf("refs/heads/auto-pr-%d", pr.Number), runner)
	return nil
}

func mergeCassettes(basePath, baseBranch, prPath string, runner source.Runner) {
	branchPath := ""
	if baseBranch != "main" {
		branchPath = "/refs/branches/" + baseBranch
	}

	if err := listCassettes(
		fmt.Sprintf("%s/%s/fixtures/", basePath, prPath),
		runner,
	); err != nil {
		fmt.Println(err)
		return
	}

	cpCassettes(
		fmt.Sprintf("%s/%s/fixtures/*", basePath, prPath),
		fmt.Sprintf("%s%s/fixtures/", basePath, branchPath),
		runner,
	)

	rmCassettes(fmt.Sprintf("%s/%s/", basePath, prPath), runner)
}

func listCassettes(path string, runner source.Runner) error {
	lsArgs := []string{
		"ls",
		path,
	}
	fmt.Println("Running command: ", "gsutil", lsArgs)
	ret, err := runner.Run("gsutil", lsArgs, nil)
	if err != nil {
		return err
	}
	fmt.Println(ret)
	return nil
}

func cpCassettes(src, dest string, runner source.Runner) {
	cpArgs := []string{
		"-m",
		"cp",
		src,
		dest,
	}
	fmt.Println("Running command: ", "gsutil", cpArgs)
	if _, err := runner.Run("gsutil", cpArgs, nil); err != nil {
		fmt.Println("Error in copy: ", err)
	}
}

func rmCassettes(dest string, runner source.Runner) {
	rmArgs := []string{
		"-m",
		"rm",
		"-r",
		dest,
	}
	fmt.Println("Running command: ", "gsutil", rmArgs)
	if _, err := runner.Run("gsutil", rmArgs, nil); err != nil {
		fmt.Println("Error in remove: ", err)
	}
}

func findPRBySHA(sha string, arr []github.PullRequest) *github.PullRequest {
	for _, pr := range arr {
		if pr.MergeCommitSha == sha {
			return &pr
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(vcrMergeCmd)
}
