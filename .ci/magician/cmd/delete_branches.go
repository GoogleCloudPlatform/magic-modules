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
	"magician/github"
	"magician/source"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var deleteBranchesCmd = &cobra.Command{
	Use:   "delete-branches",
	Short: "Delete the auto pr branches after pushing the given commit",
	Long: `This command deletes auto pr branches after pushing the given commit SHA to downstreams.

	It expects the following parameters:
	1. COMMIT_SHA
	2. BASE_BRANCH

	It also expects the following environment variables:
	1. GITHUB_TOKEN_CLASSIC
	2. GOPATH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		baseBranch := args[0]
		sha := args[1]

		githubToken, ok := os.LookupEnv("GITHUB_TOKEN_CLASSIC")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN_CLASSIC environment variable")
		}

		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating Runner: %s", err)
		}

		gh := github.NewClient(githubToken)
		if err != nil {
			return fmt.Errorf("error creating GitHub client: %s", err)
		}

		return execDeleteBranchesCmd(baseBranch, sha, githubToken, rnr, gh)
	},
}

func execDeleteBranchesCmd(baseBranch, sha, githubToken string, runner ExecRunner, gh GithubClient) error {
	prNumber, err := fetchPRNumber(sha, baseBranch, runner, gh)

	if err != nil {
		return err
	}

	err = deleteBranches(prNumber, githubToken, runner)

	return err
}

func fetchPRNumber(sha, baseBranch string, runner ExecRunner, gh GithubClient) (string, error) {
	message, err := gh.GetCommitMessage("hashicorp", "terraform-provider-google-beta", sha)
	if err != nil {
		return "", fmt.Errorf("error getting commit message: %s", err)
	}

	messageLines := strings.Split(message, "\n")

	messageParts := strings.Split(messageLines[0], " ")

	prNumber := strings.Trim(messageParts[len(messageParts)-1], "()#\n")

	_, err = strconv.ParseInt(prNumber, 10, 64)
	if err != nil {
		return "", fmt.Errorf("error parsing PR number: %s", err)
	}

	return prNumber, nil
}

var repoList = []string{
	"terraform-provider-google",
	"terraform-provider-google-beta",
	"terraform-google-conversion",
	"tf-oics",
}

func deleteBranches(prNumber, githubToken string, runner source.Runner) error {
	for _, repo := range repoList {
		for _, branch := range []string{
			fmt.Sprintf(":auto-pr-%s", prNumber),
			fmt.Sprintf(":auto-pr-%s-old", prNumber),
		} {
			_, err := runner.Run("git", []string{
				"push",
				fmt.Sprintf("https://modular-magician:%s@github.com/modular-magician/%s", githubToken, repo),
				branch,
			}, nil)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(deleteBranchesCmd)
}
