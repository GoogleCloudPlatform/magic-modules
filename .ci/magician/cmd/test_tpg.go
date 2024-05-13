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

type ttGithub interface {
	CreateWorkflowDispatchEvent(string, map[string]any) error
}

var testTPGCmd = &cobra.Command{
	Use:   "test-tpg",
	Short: "Run provider unit tests via workflow dispatch",
	Long: `This command runs provider unit tests via workflow dispatch

	The following PR details are expected as environment variables:
        1. VERSION (beta or ga)
        2. COMMIT_SHA
        3. PR_NUMBER
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		version := os.Getenv("VERSION")
		commit := os.Getenv("COMMIT_SHA")
		pr := os.Getenv("PR_NUMBER")

		githubToken, ok := lookupGithubTokenOrFallback("GITHUB_TOKEN_MAGIC_MODULES")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN_MAGIC_MODULES or GITHUB_TOKEN environment variables")
		}
		gh := github.NewClient(githubToken)

		return execTestTPG(version, commit, pr, gh)
	},
}

func execTestTPG(version, commit, pr string, gh ttGithub) error {
	var repo string
	var content []byte
	var err error
	if version == "ga" {
		repo = "terraform-provider-google"
		content, err = os.ReadFile("/workspace/upstreamCommitSHA-terraform-provider-google.txt")
		if err != nil {
			fmt.Println("Error:", err)
		}
	} else if version == "beta" {
		repo = "terraform-provider-google-beta"
		content, err = os.ReadFile("/workspace/upstreamCommitSHA-terraform-provider-google-beta.txt")
		if err != nil {
			fmt.Println("Error:", err)
		}
	} else {
		return fmt.Errorf("invalid version specified")
	}

	commitShaOrBranchUpstream := string(content)

	if commitShaOrBranchUpstream == ""{
		commitShaOrBranchUpstream = "auto-pr-" + pr
	}

	fmt.Println("commitShaOrBranchUpstream: ", commitShaOrBranchUpstream)

	if err := gh.CreateWorkflowDispatchEvent("test-tpg.yml", map[string]any{
		"owner":  "modular-magician",
		"repo":   repo,
		"branch": commitShaOrBranchUpstream,
		"sha":    commit,
	}); err != nil {
		return fmt.Errorf("error creating workflow dispatch event: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(testTPGCmd)
}
