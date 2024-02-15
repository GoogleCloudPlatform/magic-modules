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

// membershipCheckerCmd represents the membershipChecker command
var membershipCheckerCmd = &cobra.Command{
	Use:   "membership-checker",
	Short: "Assigns reviewers and manages pull request processing based on the author's trust level.",
	Long: `This command conducts a series of validations and actions based on the details and authorship of a provided pull request.

	The command expects the following pull request details as arguments:
	1. PR Number
	2. Commit SHA
	3. Branch Name
	4. Head Repo URL
	5. Head Branch
	6. Base Branch

	It then performs the following operations:
	1. Extracts and displays the pull request details.
	2. Fetches the author of the pull request and determines their contribution type.
	3. If the author is not a core contributor:
			a. Identifies the initially requested reviewer and those who previously reviewed this PR.
			b. Determines and requests reviewers based on the above.
			c. Posts comments tailored to the contribution type, the trust level of the contributor, and the primary reviewer.
	4. For trusted authors (Core Contributors and Googlers):
			a. Triggers generate-diffs using the provided PR details.
			b. Automatically approves the community-checker run.
	5. For external or untrusted contributors:
			a. Adds the 'awaiting-approval' label.
			b. Posts a link prompting approval for the build.
	`,
	Args: cobra.ExactArgs(6),
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
		execMembershipChecker(prNumber, commitSha, branchName, headRepoUrl, headBranch, baseBranch, gh, cb)
	},
}

func execMembershipChecker(prNumber, commitSha, branchName, headRepoUrl, headBranch, baseBranch string, gh GithubClient, cb CloudbuildClient) {
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
	rootCmd.AddCommand(membershipCheckerCmd)
}
