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
	"magician/github"

	"github.com/spf13/cobra"
)

// removeLabelCmd represents the remove-awaiting-approval-label command
var removeLabelCmd = &cobra.Command{
	Use:   "remove-awaiting-approval-label",
	Short: "remove awaiting-approval label",
	Long: `This command processes pull requests and performs various validations and actions based on the PR's metadata and author.

	 The following PR details are expected as arguments:
	 1. PR Number

	 The command performs the following steps:
	 1. Remove the 'awaiting-approval' label from the PR.
	 `,
	RunE: func(cmd *cobra.Command, args []string) error {
		prNumber := args[0]
		labelName := args[1]
		fmt.Println("PR Number: ", prNumber)

		githubToken, ok := lookupGithubTokenOrFallback("GITHUB_TOKEN_MAGIC_MODULES")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN_MAGIC_MODULES or GITHUB_TOKEN environment variables")
		}
		gh := github.NewClient(githubToken)

		execRemoveLabel(prNumber, gh, labelName)
		return nil
	},
}

func execRemoveLabel(prNumber string, gh GithubClient, labelName string) {
	gh.RemoveLabel(prNumber, labelName)
}

func init() {
	rootCmd.AddCommand(removeLabelCmd)
}
