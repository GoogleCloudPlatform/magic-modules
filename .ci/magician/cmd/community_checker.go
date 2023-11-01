/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"magician/cloudbuild"
	"magician/github"
	"os"

	"github.com/spf13/cobra"
)

type ccGithub interface {
	GetPullRequestAuthor(prNumber string) (string, error)
	GetUserType(user string) github.UserType
	RemoveLabel(prNumber string, label string) error
	PostBuildStatus(prNumber string, title string, state string, targetUrl string, commitSha string) error
}

type ccCloudbuild interface {
	TriggerMMPresubmitRuns(commitSha string, substitutions map[string]string) error
}

// communityApprovalCmd represents the communityApproval command
var communityApprovalCmd = &cobra.Command{
	Use:   "community-checker",
	Short: "Run presubmit generate diffs for untrusted users and remove awaiting-approval label",
	Long: `This command processes pull requests and performs various validations and actions based on the PR's metadata and author.

	The following PR details are expected as arguments:
	1. PR Number
	2. Commit SHA
	3. Branch Name
	4. Head Repo URL
	5. Head Branch
	6. Base Branch

	The command performs the following steps:
	1. Retrieve and print the provided pull request details.
	2. Get the author of the pull request and determine their user type.
	3. If the author is not a trusted user (neither a Core Contributor nor a Googler):
			a. Trigger cloud builds with specific substitutions for the PR.
	4. For all pull requests, the 'awaiting-approval' label is removed.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		prNumber := args[0]
		fmt.Println("PR Number: ", prNumber)

		commitSha := args[1]
		fmt.Println("Commit SHA: ", commitSha)

		branchName := args[2]
		fmt.Println("Branch Name: ", branchName)

		headRepoUrl := args[3]
		fmt.Println("Head Repo URL: ", headRepoUrl)

		headBranch := args[4]
		fmt.Println("Head Branch: ", headBranch)

		baseBranch := args[5]
		fmt.Println("Base Branch: ", baseBranch)

		gh := github.NewGithubService()
		cb := cloudbuild.NewCloudBuildService()
		execCommunityChecker(prNumber, commitSha, branchName, headRepoUrl, headBranch, baseBranch, gh, cb)
	},
}

func execCommunityChecker(prNumber, commitSha, branchName, headRepoUrl, headBranch, baseBranch string, gh ccGithub, cb ccCloudbuild) {
	substitutions := map[string]string{
		"BRANCH_NAME":    branchName,
		"_PR_NUMBER":     prNumber,
		"_HEAD_REPO_URL": headRepoUrl,
		"_HEAD_BRANCH":   headBranch,
		"_BASE_BRANCH":   baseBranch,
	}

	author, err := gh.GetPullRequestAuthor(prNumber)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	authorUserType := gh.GetUserType(author)
	trusted := authorUserType == github.CoreContributorUserType || authorUserType == github.GooglerUserType

	// only triggers build for untrusted users (because trusted users will be handled by membership-checker)
	if !trusted {
		err = cb.TriggerMMPresubmitRuns(commitSha, substitutions)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// in community-checker job:
	// remove awaiting-approval label from external contributor PRs
	gh.RemoveLabel(prNumber, "awaiting-approval")
}

func init() {
	rootCmd.AddCommand(communityApprovalCmd)
}
