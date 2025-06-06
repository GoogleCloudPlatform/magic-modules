/*
* Copyright 2023 Google LLC. All Rights Reserved.
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
	"magician/cloudbuild"
	"magician/github"

	"github.com/spf13/cobra"
)

// membershipCheckerCmd represents the membershipChecker command
var membershipCheckerCmd = &cobra.Command{
	Use:   "membership-checker",
	Short: "Assigns reviewers and manages pull request processing based on the author's trust level.",
	Long: `This command conducts a series of validations and actions based on the details and authorship of a provided pull request.

	The command expects the following pull request details as arguments:
	1. PR Number
	2. Commit SHA

	It then performs the following operations:
	1. Extracts and displays the pull request details.
	2. Fetches the author of the pull request and determines their contribution type.
	3. For trusted authors (Core Contributors and Googlers):
			a. Automatically approves the community-checker run.
	4. For external or untrusted contributors:
			a. Adds the 'awaiting-approval' label.
			b. Posts a link prompting approval for the build.
	`,
	// This can change to cobra.ExactArgs(2) after at least a 2-week soak
	Args: cobra.RangeArgs(2, 6),
	RunE: func(cmd *cobra.Command, args []string) error {
		prNumber := args[0]
		fmt.Println("PR Number: ", prNumber)

		commitSha := args[1]
		fmt.Println("Commit SHA: ", commitSha)

		githubToken, ok := lookupGithubTokenOrFallback("GITHUB_TOKEN_MAGIC_MODULES")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN_MAGIC_MODULES or GITHUB_TOKEN environment variables")
		}
		gh := github.NewClient(githubToken)
		cb := cloudbuild.NewClient()
		return execMembershipChecker(prNumber, commitSha, gh, cb)
	},
}

func execMembershipChecker(prNumber, commitSha string, gh GithubClient, cb CloudbuildClient) error {
	pullRequest, err := gh.GetPullRequest(prNumber)
	if err != nil {
		return err
	}

	author := pullRequest.User.Login
	authorUserType := gh.GetUserType(author)
	trusted := authorUserType == github.CoreContributorUserType || authorUserType == github.GooglerUserType

	// 1. auto approve community-checker run for trusted users
	// 2. add awaiting-approval label to external contributor PRs
	if trusted {
		err = cb.ApproveDownstreamGenAndTest(prNumber, commitSha)
		if err != nil {
			return err
		}
	} else {
		gh.AddLabels(prNumber, []string{"awaiting-approval"})
	}
	return nil
}

func init() {
	rootCmd.AddCommand(membershipCheckerCmd)
}
