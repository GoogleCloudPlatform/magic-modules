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
	"os"

	"github.com/spf13/cobra"
)

var testTGCCmd = &cobra.Command{
	Use:   "test-tgc",
	Short: "Run tgc unit tests via workflow dispatch",
	Long: `This command runs tgc unit tests via workflow dispatch

	The following PR details are expected as environment variables:
        1. COMMIT_SHA
        2. PR_NUMBER
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		commit := os.Getenv("COMMIT_SHA")
		pr := os.Getenv("PR_NUMBER")

		githubToken, ok := lookupGithubTokenOrFallback("GITHUB_TOKEN_MAGIC_MODULES")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN_MAGIC_MODULES or GITHUB_TOKEN environment variables")
		}
		gh := github.NewClient(githubToken)

		return execTestTGC(commit, pr, gh)
	},
}

func execTestTGC(commit, pr string, gh ttGithub) error {
	contentTPGB, err := os.ReadFile("/workspace/upstreamCommitSHA-terraform-provider-google-beta.txt")
	if err != nil {
		fmt.Println("Error:", err)
	}

	contentTGC, err := os.ReadFile("/workspace/upstreamCommitSHA-terraform-google-conversion.txt")
	if err != nil {
		fmt.Println("Error:", err)
	}

	commitShaOrBranchUpstreamTPGB := string(contentTPGB)
	commitShaOrBranchUpstreamTGC := string(contentTGC)

	if commitShaOrBranchUpstreamTPGB == "" {
		commitShaOrBranchUpstreamTPGB = "auto-pr-" + pr
	}

	if commitShaOrBranchUpstreamTGC == "" {
		commitShaOrBranchUpstreamTGC = "auto-pr-" + pr
	}

	fmt.Println("commitShaOrBranchUpstreamTPGB: ", commitShaOrBranchUpstreamTPGB)
	fmt.Println("commitShaOrBranchUpstreamTGC: ", commitShaOrBranchUpstreamTGC)

	if err := gh.CreateWorkflowDispatchEvent("test-tgc.yml", map[string]any{
		"owner":      "modular-magician",
		"repo":       "terraform-google-conversion",
		"tpgbbranch": commitShaOrBranchUpstreamTPGB,
		"tgcbranch":  commitShaOrBranchUpstreamTGC,
		"sha":        commit,
	}); err != nil {
		return fmt.Errorf("error creating workflow dispatch event: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(testTGCCmd)
}
