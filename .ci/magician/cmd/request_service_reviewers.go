/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"magician/github"
	"os"

	"github.com/spf13/cobra"
)

// requestServiceReviewersCmd represents the requestServiceReviewers command
var requestServiceReviewersCmd = &cobra.Command{
	Use:   "request-service-reviewers PR_NUMBER",
	Short: "Assigns reviewers based on the PR's service labels.",
	Long: `This command requests (or re-requests) review based on the PR's service labels.

	If a PR has more than 3 service labels, the command will not do anything.
	`,
	Args: cobra.ExactArgs(6),
	Run: func(cmd *cobra.Command, args []string) {
		prNumber := args[0]
		fmt.Println("PR Number: ", prNumber)

		gh := github.NewGithubService()
		execRequestServiceReviewers(prNumber, gh)
	},
}

func execRequestServiceReviewers(prNumber string, gh github.GithubService) {
	pullRequest, err := gh.GetPullRequest(prNumber)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if authorUserType != github.CoreContributorUserType {
		fmt.Println("Not core contributor - assigning reviewer")

		previouslyInvolvedReviewers, err := gh.GetPullRequestPreviousAssignedReviewers(prNumber)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		reviewersToRequest, newPrimaryReviewer := github.ChooseCoreReviewers(firstRequestedReviewer, previouslyInvolvedReviewers)

		for _, reviewer := range reviewersToRequest {
			err = gh.RequestPullRequestReviewer(prNumber, reviewer)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if newPrimaryReviewer != "" {
			comment := github.FormatReviewerComment(newPrimaryReviewer, authorUserType, trusted)
			err = gh.PostComment(prNumber, comment)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	// auto_run(contributor-membership-checker) will be run on every commit or /gcbrun:
	// only triggers builds for trusted users
	if trusted {
		err = cb.TriggerMMPresubmitRuns(commitSha, substitutions)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// in contributor-membership-checker job:
	// 1. auto approve community-checker run for trusted users
	// 2. add awaiting-approval label to external contributor PRs
	if trusted {
		err = cb.ApproveCommunityChecker(prNumber, commitSha)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		gh.AddLabel(prNumber, "awaiting-approval")
		targetUrl, err := cb.GetAwaitingApprovalBuildLink(prNumber, commitSha)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = gh.PostBuildStatus(prNumber, "Approve Build", "success", targetUrl, commitSha)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func init() {
	rootCmd.AddCommand(requestServiceReviewersCmd)
}
