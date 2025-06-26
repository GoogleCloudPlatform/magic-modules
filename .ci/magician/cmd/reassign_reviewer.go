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
	"errors"
	"fmt"
	"magician/github"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// reassignReviewerCmd represents the reassignReviewer command
var reassignReviewerCmd = &cobra.Command{
	Use:   "reassign-reviewer PR_NUMBER [REVIEWER]",
	Short: "Reassigns primary reviewer to the given reviewer or a random reviewer if none given",
	Long: `This command reassigns reviewers when invoked via a comment on a pull request.

	The command expects the following PR details as arguments:
	1. PR_NUMBER
	2. COMMENT_AUTHOR
	3. REVIEWER (optional)


	It then performs the following operations:
	1. Updates the reviewer comment to reflect the new primary reviewer.
	2. Requests a review from the new primary reviewer.
	`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		prNumber := args[0]
		fmt.Println("PR Number: ", prNumber)

		githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN environment variable")
		}
		gh := github.NewClient(githubToken)

		author := args[1]
		if gh.GetUserType(author) != github.CoreContributorUserType {
			return fmt.Errorf("comment author is not a core contributor")
		}

		var newPrimaryReviewer string
		if len(args) > 2 {
			newPrimaryReviewer = strings.TrimPrefix(args[2], "@")
		}
		return execReassignReviewer(prNumber, newPrimaryReviewer, gh)
	},
}

func execReassignReviewer(prNumber, newPrimaryReviewer string, gh GithubClient) error {
	pullRequest, err := gh.GetPullRequest(prNumber)
	if err != nil {
		return err
	}
	comments, err := gh.GetPullRequestComments(prNumber)
	if err != nil {
		return err
	}

	reviewerComment, currentReviewer := github.FindReviewerComment(comments)
	if newPrimaryReviewer == "" {
		newPrimaryReviewer = github.GetRandomReviewer([]string{currentReviewer, pullRequest.User.Login})
	}

	if newPrimaryReviewer == "" {
		return errors.New("no primary reviewer found")
	}
	if newPrimaryReviewer == currentReviewer {
		return fmt.Errorf("primary reviewer is already %s", newPrimaryReviewer)
	}

	fmt.Println("New primary reviewer is ", newPrimaryReviewer)
	comment := github.FormatReviewerComment(newPrimaryReviewer)

	if currentReviewer == "" {
		fmt.Println("No reviewer comment found, creating one")
		err := gh.PostComment(prNumber, comment)
		if err != nil {
			return err
		}
	} else {
		if err := gh.RemovePullRequestReviewers(prNumber, []string{currentReviewer}); err != nil {
			fmt.Printf("Failed to remove reviewer %s from pull request: %s\n", currentReviewer, err)
		}
		fmt.Println("Updating reviewer comment")
		err := gh.UpdateComment(prNumber, comment, reviewerComment.ID)
		if err != nil {
			return err
		}
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
