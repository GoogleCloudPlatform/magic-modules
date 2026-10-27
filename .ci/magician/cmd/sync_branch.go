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
	"time"

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

	remote := fmt.Sprintf("https://modular-magician:%s@github.com/GoogleCloudPlatform/magic-modules", githubToken)
	refspec := fmt.Sprintf("%s:%s", sha, syncBranch)
	return syncBranchPushWithRetry(runner, remote, refspec)
}

func syncBranchPushWithRetry(runner source.Runner, remote, refspec string) error {
	var err error
	for attempt := 0; ; attempt++ {
		if _, err = runner.Run("git", []string{"push", remote, refspec}, nil); err == nil {
			return nil
		}
		if attempt >= len(pushRetryBackoffs) {
			break
		}
		fmt.Printf("Push of %s failed (attempt %d of %d), retrying in %s: %v\n", refspec, attempt+1, len(pushRetryBackoffs)+1, pushRetryBackoffs[attempt], err)
		time.Sleep(pushRetryBackoffs[attempt])
	}
	return fmt.Errorf("push of %s failed after %d attempts: %w", refspec, len(pushRetryBackoffs)+1, err)
}

func init() {
	rootCmd.AddCommand(syncBranchCmd)
}
