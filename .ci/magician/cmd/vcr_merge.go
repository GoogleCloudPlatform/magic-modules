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
	Run: func(cmd *cobra.Command, args []string) {
		reference := args[0]
		fmt.Println("Reference commit SHA: ", reference)
		githubToken, ok := os.LookupEnv("GITHUB_TOKEN_CLASSIC")
		if !ok {
			fmt.Println("Did not provide GITHUB_TOKEN_CLASSIC environment variable")
			os.Exit(1)
		}
		rnr, err := exec.NewRunner()
		if err != nil {
			fmt.Println("Error creating Runner: ", err)
			os.Exit(1)
		}

		baseBranch := os.Getenv("BASE_BRANCH")
		if baseBranch == "" {
			baseBranch = "main"
		}
		gh := github.NewClient(githubToken)
		execVCRMerge(gh, reference, baseBranch, rnr)
	},
}

func execVCRMerge(gh GithubClient, sha string, baseBranch string, runner source.Runner) {
	arr, err := gh.GetPullRequests("closed", baseBranch, "updated", "desc")
	if err != nil {
		fmt.Println("Error getting pull requests: ", err)
		os.Exit(1)
	}
	pr := findPRBySHA(sha, arr)
	if pr == nil {
		fmt.Printf("Not finding any matching PR with commit SHA %s\n", sha)
		return
	}

	mergeCassettes(false, baseBranch, pr.Number, runner)
	mergeCassettes(true, baseBranch, pr.Number, runner)
}

func mergeCassettes(isBeta bool, baseBranch string, prNumber int, runner source.Runner) {
	prefix := "gs://ci-vcr-cassettes"
	if isBeta {
		prefix = "gs://ci-vcr-cassettes/beta"
	}
	prPath := fmt.Sprintf("refs/heads/auto-pr-%d", prNumber)
	branchPath := ""
	if baseBranch != "main" {
		branchPath = "/refs/branches/" + baseBranch
	}

	lsPath := fmt.Sprintf("%s/%s/fixtures/", prefix, prPath)
	err := listCassettes(lsPath, runner)
	if err != nil {
		fmt.Println(err)
		return
	}

	cpSrcPath := fmt.Sprintf("%s/%s/fixtures/*", prefix, prPath)
	cpDestPath := fmt.Sprintf("%s%s/fixtures/", prefix, branchPath)
	cpCassettes(cpSrcPath, cpDestPath, runner)

	rmDestPath := fmt.Sprintf("%s/%s/", prefix, prPath)
	rmCassettes(rmDestPath, runner)
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
