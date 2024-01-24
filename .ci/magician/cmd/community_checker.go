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
	"os"

	"github.com/spf13/cobra"
)

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

		gh := github.NewClient()
		cb := cloudbuild.NewClient()
		execCommunityChecker(prNumber, commitSha, branchName, headRepoUrl, headBranch, baseBranch, gh, cb)
	},
}

func execCommunityChecker(prNumber, commitSha, branchName, headRepoUrl, headBranch, baseBranch string, gh GithubClient, cb CloudbuildClient) {
	substitutions := map[string]string{
		"BRANCH_NAME":    branchName,
		"_PR_NUMBER":     prNumber,
		"_HEAD_REPO_URL": headRepoUrl,
		"_HEAD_BRANCH":   headBranch,
		"_BASE_BRANCH":   baseBranch,
	}

	pullRequest, err := gh.GetPullRequest(prNumber)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	author := pullRequest.User.Login
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
