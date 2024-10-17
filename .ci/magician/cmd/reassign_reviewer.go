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

	"github.com/spf13/cobra"
)

// reassignReviewerCmd represents the reassignReviewer command
var reassignReviewerCmd = &cobra.Command{
	Use:   "reassign-reviewer PR_NUMBER",
	Short: "Reassigns primary reviewer to the given reviewer or a random reviewer if none given",
	Long: `This command reassigns reviewers when invoked via a comment on a pull request.

	The command expects the following PR details as arguments:
	1. PR_NUMBER
	2. REVIEWER (optional)


	It then performs the following operations:
	1. Updates the reviewer comment to reflect the new primary reviewer.
	`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prNumber := args[0]
		fmt.Println("PR Number: ", prNumber)

		githubToken, ok := lookupGithubTokenOrFallback("GITHUB_TOKEN_MAGIC_MODULES")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN_MAGIC_MODULES or GITHUB_TOKEN environment variable")
		}
		gh := github.NewClient(githubToken)
		var newPrimaryReviewer string
		if len(args) > 1 {
			newPrimaryReviewer = args[1]
		}
		return execReassignReviewer(prNumber, newPrimaryReviewer, gh)
	},
}

func execReassignReviewer(prNumber, newPrimaryReviewer string, gh GithubClient) error {
	comments, err := gh.GetPullRequestCommentsByUser(prNumber, "modular-magician")
	if err != nil {
		return err
	}

	reviewerComment, currentReviewer := github.FindReviewerComment(comments)

	if currentReviewer == "" {
		return fmt.Errorf("no reviewer comment found in PR %s", prNumber)
	}

	if newPrimaryReviewer == "" {
		newPrimaryReviewer = github.GetNewRandomReviewer(currentReviewer)
	}

	if currentReviewer == newPrimaryReviewer {
		return fmt.Errorf("primary reviewer is already %s", newPrimaryReviewer)
	}

	err = gh.UpdateComment(prNumber, github.FormatReviewerComment(newPrimaryReviewer), reviewerComment.ID)
	if err != nil {
		return err
	}

	err = gh.RequestPullRequestReviewers(prNumber, []string{newPrimaryReviewer})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(reassignReviewerCmd)
}
