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

		goPath, ok := os.LookupEnv("GOPATH")
		if !ok {
			return fmt.Errorf("did not provide GOPATH environment variable")
		}

		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating Runner: %s", err)
		}

		ctlr := source.NewController(goPath, "modular-magician", githubToken, rnr)

		return execDeleteBranchesCmd(baseBranch, sha, githubToken, rnr, ctlr)
	},
}

func execDeleteBranchesCmd(baseBranch, sha, githubToken string, runner source.Runner, controller *source.Controller) error {
	prNumber, err := fetchPRNumber(sha, baseBranch, runner, controller)

	if err != nil {
		return err
	}

	err = deleteBranches(prNumber, githubToken, runner)

	return err
}

func fetchPRNumber(sha, baseBranch string, runner source.Runner, controller *source.Controller) (string, error) {
	repo := &source.Repo{
		Name:   "terraform-provider-google-beta",
		Branch: baseBranch,
	}
	controller.SetPath(repo)

	if err := controller.Clone(repo); err != nil {
		return "", err
	}

	controller.SetPath(repo)

	if err := runner.PushDir(repo.Path); err != nil {
		return "", err
	}

	message, err := runner.Run("git", []string{
		"show",
		"-s",
		"--format=%s",
		sha,
	}, nil)
	if err != nil {
		return "", fmt.Errorf("error getting commit message: %s", err)
	}

	messageParts := strings.Split(message, " ")

	prNumber := messageParts[len(messageParts)-1]

	return strings.Trim(prNumber, "()#\n"), nil
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
