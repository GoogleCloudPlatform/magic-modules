/*
* Copyright 2024 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package cmd

import (
	"fmt"
	"magician/github"
	"os"

	"github.com/spf13/cobra"
)

// requestReviewerCmd represents the requestReviewer command
var requestReviewerCmd = &cobra.Command{
	Use:   "request-reviewer",
	Short: "Assigns and re-requests reviewers",
	Long: `This command automatically requests (or re-requests) core contributor reviews for a PR based on whether the user is a core contributor.

	The command expects the following pull request details as arguments:
	1. PR Number

	It then performs the following operations:
	1. Determines the author of the pull request
	2. If the author is not a core contributor:
			a. Identifies the initially requested reviewer and those who previously reviewed this PR.
			b. Determines and requests reviewers based on the above.
			c. As appropriate, posts a welcome comment on the PR.
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prNumber := args[0]
		fmt.Println("PR Number: ", prNumber)
		githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN environment variable")
		}
		gh := github.NewClient(githubToken)
		return execRequestReviewer(prNumber, gh)
	},
}

func execRequestReviewer(prNumber string, gh GithubClient) error {
	pullRequest, err := gh.GetPullRequest(prNumber)
	if err != nil {
		return err
	}

	author := pullRequest.User.Login
	if !github.IsCoreContributor(author) {
		fmt.Println("Not core contributor - assigning reviewer")

		requestedReviewers, err := gh.GetPullRequestRequestedReviewers(prNumber)
		if err != nil {
			return err
		}

		previousReviewers, err := gh.GetPullRequestPreviousReviewers(prNumber)
		if err != nil {
			return err
		}

		reviewersToRequest, newPrimaryReviewer := github.ChooseCoreReviewers(requestedReviewers, previousReviewers)

		if len(reviewersToRequest) > 0 {
			err = gh.RequestPullRequestReviewers(prNumber, reviewersToRequest)
			if err != nil {
				return err
			}
		}

		if newPrimaryReviewer != "" {
			comment := github.FormatReviewerComment(newPrimaryReviewer)
			err = gh.PostComment(prNumber, comment)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(requestReviewerCmd)
}
