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
	"magician/cloudstorage"
	"magician/provider"
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
	 3. Closes 100% test ticket if it starts to pass for 3 days
 
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

		gcs := cloudstorage.NewClient()

		now := time.Now()

		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			return fmt.Errorf("Error loading location: %s", err)
		}
		date := now.In(loc)

		return execManageTestFailureTicket(date, gh, gcs)
	},
}

func listMTFTRequiredEnvironmentVariables() string {
	var result string
	for i, ev := range mtftRequiredEnvironmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execManageTestFailureTicket(now time.Time, gh *github.Client, gcs CloudstorageClient) error {
	ctx := context.Background()

	// Get test tickets with "forward/review" labels
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

	// Get Test status for past 3 days
	gaTestFailuresMap := make(map[string][]bool)
	betaTestFailuresMap := make(map[string][]bool)

	lastNDaysTestNonSuccessMap(provider.GA, 3, now, gcs, gaTestFailuresMap)
	lastNDaysTestNonSuccessMap(provider.Beta, 3, now, gcs, betaTestFailuresMap)

	// Get 100% failing test tickets
	opts = &github.IssueListByRepoOptions{
		State:       "open",
		Labels:      []string{"test-failure-100"},
		ListOptions: github.ListOptions{PerPage: 100},
	}
	issues, err = ListIssuesWithOpts(ctx, gh, opts)
	if err != nil {
		return err
	}

	var shouldCloseTickets []int // store GH issue number

	for _, issue := range issues {
		// Get failing test names
		tests, err := testNamesFromIssue(issue)
		if err != nil {
			return err
		}
		if len(tests) == 0 {
			fmt.Println("No tests found for issue ", issue.GetNumber())
			continue
		}

		if shouldCloseTestTicket(tests, gaTestFailuresMap, betaTestFailuresMap) {
			shouldCloseTickets = append(shouldCloseTickets, issue.GetNumber())
		}
	}

	comment := "All failing tests listed in this ticket have passed in the last three consecutive nightly runs. Closing the ticket."
	for _, ticketNumber := range shouldCloseTickets {
		fmt.Println("Closing ticket ", ticketNumber)
		issueComment := &github.IssueComment{
			Body: github.String(comment),
		}
		_, _, err = gh.Issues.CreateComment(ctx, GithubOwner, GithubRepo, ticketNumber, issueComment)
		if err != nil {
			return fmt.Errorf("error posting comment to issue %d: %w", ticketNumber, err)
		}
		issueRquest := &github.IssueRequest{
			State: github.String("closed"),
		}
		_, _, err = gh.Issues.Edit(ctx, GithubOwner, GithubRepo, ticketNumber, issueRquest)
		if err != nil {
			return fmt.Errorf("error closing issue %d: %w", ticketNumber, err)
		}
	}

	return nil
}

func lastNDaysTestNonSuccessMap(pVersion provider.Version, n int, now time.Time, gcs CloudstorageClient, testFailuresMap map[string][]bool) error {
	for i := 0; i < n; i++ {
		date := now.AddDate(0, 0, -i)
		testInfoList, err := getTestInfoList(pVersion, date, gcs)
		if err != nil {
			return fmt.Errorf("error getting test info list: %w", err)
		}
		for _, testInfo := range testInfoList {
			testName := testInfo.Name
			if _, ok := testFailuresMap[testName]; !ok {
				testFailuresMap[testName] = make([]bool, n)
			}
			testFailuresMap[testName][i] = (testInfo.Status == "FAILURE" || testInfo.Status == "UNKNOWN") // failed or skipped
		}
	}
	return nil
}

func shouldCloseTestTicket(tests []string, gaTestFailuresMap, betaTestFailuresMap map[string][]bool) bool {
	for _, test := range tests {
		gaFailures, foundGaTest := gaTestFailuresMap[test]
		betaFailures, foundBetaTest := betaTestFailuresMap[test]

		if !foundGaTest && !foundBetaTest {
			fmt.Printf("test %s not found in either GA or Beta, might be skipped\n", test)
			return false
		}

		if foundGaTest {
			for _, fail := range gaFailures {
				if fail {
					return false
				}
			}
		}
		if foundBetaTest {
			for _, fail := range betaFailures {
				if fail {
					return false
				}
			}
		}
	}
	return true
}

func init() {
	rootCmd.AddCommand(manageTestFailureTicketCmd)
}
