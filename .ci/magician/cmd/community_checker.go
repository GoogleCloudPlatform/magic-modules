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
	1. Trigger cloud presubmits with specific substitutions for the PR.
	2. Remove the 'awaiting-approval' label from the PR.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		githubToken, ok := lookupGithubTokenOrFallback("GITHUB_TOKEN_MAGIC_MODULES")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN_MAGIC_MODULES or GITHUB_TOKEN environment variables")
		}
		gh := github.NewClient(githubToken)
		cb := cloudbuild.NewClient()
		return execCommunityChecker(prNumber, commitSha, branchName, headRepoUrl, headBranch, baseBranch, gh, cb)
	},
}

func execCommunityChecker(prNumber, commitSha, branchName, headRepoUrl, headBranch, baseBranch string, gh GithubClient, cb CloudbuildClient) error {
	substitutions := map[string]string{
		"BRANCH_NAME":    branchName,
		"_PR_NUMBER":     prNumber,
		"_HEAD_REPO_URL": headRepoUrl,
		"_HEAD_BRANCH":   headBranch,
		"_BASE_BRANCH":   baseBranch,
	}

	// trigger presubmit builds - community-checker requires approval
	// (explicitly or via membership-checker)
	err := cb.TriggerMMPresubmitRuns(commitSha, substitutions)
	if err != nil {
		return err
	}

	// in community-checker job:
	// remove awaiting-approval label from external contributor PRs
	gh.RemoveLabel(prNumber, "awaiting-approval")
	return nil
}

func init() {
	rootCmd.AddCommand(communityApprovalCmd)
}
