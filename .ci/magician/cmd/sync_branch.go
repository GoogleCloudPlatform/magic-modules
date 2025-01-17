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
	"magician/exec"
	"magician/source"
	"os"

	"github.com/spf13/cobra"
)

var syncBranchCmd = &cobra.Command{
	Use:   "sync-branch",
	Short: "Push the given commit to the given sync branch",
	Long: `This command updates the given sync branch with the given commit SHA.

	It expects the following parameters:
	1. SYNC_BRANCH_PREFIX
	2. BASE_BRANCH
	3. SHA

	It also expects the following environment variables:
	1. GITHUB_TOKEN_CLASSIC`,
	RunE: func(cmd *cobra.Command, args []string) error {
		syncBranchPrefix := args[0]
		baseBranch := args[1]
		sha := args[2]

		githubToken, ok := os.LookupEnv("GITHUB_TOKEN_CLASSIC")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN_CLASSIC environment variable")
		}

		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating Runner: %s", err)
		}
		return execSyncBranchCmd(syncBranchPrefix, baseBranch, sha, githubToken, rnr)
	},
}

func execSyncBranchCmd(syncBranchPrefix, baseBranch, sha, githubToken string, runner source.Runner) error {
	syncBranch := getSyncBranch(syncBranchPrefix, baseBranch)
	fmt.Println("SYNC_BRANCH: ", syncBranch)

	if syncBranchHasCommit(sha, syncBranch, runner) {
		fmt.Printf("Commit %s already in sync branch %s, skipping sync\n", sha, syncBranch)
		return nil
	}

	_, err := runner.Run("git", []string{"push", fmt.Sprintf("https://modular-magician:%s@github.com/GoogleCloudPlatform/magic-modules", githubToken), fmt.Sprintf("%s:%s", sha, syncBranch)}, nil)
	return err
}

func init() {
	rootCmd.AddCommand(syncBranchCmd)
}
