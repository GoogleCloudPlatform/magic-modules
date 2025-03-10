/*
* Copyright 2025 Google LLC. All Rights Reserved.
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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v68/github"
	"github.com/spf13/cobra"

	_ "embed"
)

var mtftRequiredEnvironmentVariables = [...]string{
	"GITHUB_TOKEN",
}

// manageTestFailureTicketCmd represents the manageTestFailureTicket command
var manageTestFailureTicketCmd = &cobra.Command{
	Use:   "manage-test-failure-ticket",
	Short: "Manages GitHub test failure tickets",
	Long: `This command manages the GitHub test failure tickets. 
 
	 It performs the following operations:
	 1. Lists out GitHub issues with test-failure and forward/review labels.
	 2. Removes forward/review labels from these issues.
 
	 The following environment variables are required:
 ` + listMTFTRequiredEnvironmentVariables(),
	RunE: func(cmd *cobra.Command, args []string) error {
		env := make(map[string]string)
		for _, ev := range mtftRequiredEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				return fmt.Errorf("did not provide %s environment variable", ev)
			}
			env[ev] = val
		}

		gh := github.NewClient(nil).WithAuthToken(env["GITHUB_TOKEN"])

		now := time.Now()

		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			return fmt.Errorf("Error loading location: %s", err)
		}
		date := now.In(loc)

		return execManageTestFailureTicket(date, gh)
	},
}

func listMTFTRequiredEnvironmentVariables() string {
	var result string
	for i, ev := range mtftRequiredEnvironmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execManageTestFailureTicket(now time.Time, gh *github.Client) error {
	ctx := context.Background()
	opts := &github.IssueListByRepoOptions{
		State:       "open",
		Labels:      []string{"test-failure", "forward/review"},
		ListOptions: github.ListOptions{PerPage: 100},
	}
	issues, err := ListIssuesWithOpts(ctx, gh, opts)
	if err != nil {
		return err
	}

	// Remove review labels to forward test failure tickets
	for _, issue := range issues {
		_, err := gh.Issues.RemoveLabelForIssue(ctx, GithubOwner, GithubRepo, issue.GetNumber(), "forward/review")
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(manageTestFailureTicketCmd)
}
