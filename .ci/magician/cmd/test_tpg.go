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
	Run: func(cmd *cobra.Command, args []string) {
		version := os.Getenv("VERSION")
		commit := os.Getenv("COMMIT_SHA")
		pr := os.Getenv("PR_NUMBER")

		gh := github.NewClient()

		execTestTPG(version, commit, pr, gh)
	},
}

func execTestTPG(version, commit, pr string, gh ttGithub) {
	var repo string
	if version == "ga" {
		repo = "terraform-provider-google"
	} else if version == "beta" {
		repo = "terraform-provider-google-beta"
	} else {
		fmt.Println("invalid version specified")
		os.Exit(1)
	}

	if err := gh.CreateWorkflowDispatchEvent("test-tpg.yml", map[string]any{
		"owner":  "modular-magician",
		"repo":   repo,
		"branch": "auto-pr-" + pr,
		"sha":    commit,
	}); err != nil {
		fmt.Printf("Error creating workflow dispatch event: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(testTPGCmd)
}
