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

type testIssue struct {
	testLists            []string
	issueNumber          string
	ErrorMessage         string
	DebugLogLink         string
	GaFailureRate        string
	GaFailureRateLabel   testFailureRateLabel
	BetaFailureRate      string
	BetaFailureRateLabel testFailureRateLabel
}

// manageTestFailureTicketCmd represents the manageTestFailureTicket command
var manageTestFailureTicketCmd = &cobra.Command{
	Use:   "manage-test-failure-ticket",
	Short: "Manages GitHub test failure tickets",
	Long: `This command manages the GitHub test failure tickets. 
 
	 It performs the following operations:
	 1. Lists out GitHub issues with test-failure and forward/review labels.
	 2. Removes forward/review labels from these issues.
	 3. Closes 100% test ticket if it starts to pass
 
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

	////////////////////////////////////
	// Comment out for debug
	///////////////////////////////////
	/*
		Remove review labels to forward test failure tickets
		for _, issue := range issues {
			_, err := gh.Issues.RemoveLabelForIssue(ctx, GithubOwner, GithubRepo, issue.GetNumber(), "forward/review")
			if err != nil {
				return err
			}
		}
	*/
	////////////////////////////////////

	// Get today's test status
	gaTestInfoList, err := getTestInfoList(provider.GA, now, gcs)
	if err != nil {
		return fmt.Errorf("error getting test info list: %w", err)
	}
	betaTestInfoList, err := getTestInfoList(provider.Beta, now, gcs)
	if err != nil {
		return fmt.Errorf("error getting test info list: %w", err)
	}
	gaTestFailures := testFailureSet(gaTestInfoList)
	betaTestFailures := testFailureSet(betaTestInfoList)
	gaTestSuccess := testSuccessSet(gaTestInfoList)
	betaTestSuccess := testSuccessSet(betaTestInfoList)

	fmt.Println("gaTestFailures: ", gaTestFailures)
	fmt.Println("betaTestFailures: ", betaTestFailures)

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

	var closeTicketsList []int
	for _, issue := range issues {
		tests, err := testNamesFromIssue(issue) // TODO: refactor testNamesFromIssues
		if err != nil {
			return err
		}
		if len(tests) == 0 {
			fmt.Println("No tests found for issue ", issue.GetNumber())
			continue
		}
		shouldClose := true
		for _, test := range tests {
			if _, ok := gaTestFailures[test]; ok {
				shouldClose = false
				break
			}
			if _, ok := betaTestFailures[test]; ok {
				shouldClose = false
				break
			}
			_, foundGaTest := gaTestSuccess[test]
			_, foundBetaTest := betaTestSuccess[test]
			if !foundGaTest && !foundBetaTest {
				fmt.Printf("test %s not found in either GA or Beta, might be skipped\n", test)
				shouldClose = false
				break
			}
		}
		if shouldClose {
			closeTicketsList = append(closeTicketsList, issue.GetNumber())
		}
	}

	comment := "All failing tests listed in this ticket passed in last night's integration run. Closing the ticket."
	for _, ticketNumber := range closeTicketsList {
		fmt.Println(ticketNumber)
		////////////////////////////////////
		// Comment out for debug
		///////////////////////////////////
		/*
			issueComment := &github.IssueComment{
				Body: github.String(comment),
			}
			_, _, err = gh.Issues.CreateComment(ctx, GithubOwner, GithubRepo, ticketNumber, issueComment)
			issueRquest := &github.IssueRequest{
				State: github.String("closed"),
			}
			_, _, err = gh.Issues.Edit(ctx, GithubOwner, GithubRepo, ticketNumber, issueRquest)
		*/
		////////////////////////////////////
	}

	return nil
}

func init() {
	rootCmd.AddCommand(manageTestFailureTicketCmd)
}

func testFailureSet(testInfoList []TestInfo) map[string]struct{} {
	testFailures := make(map[string]struct{})
	for _, testInfo := range testInfoList {
		if testInfo.Status == "FAILURE" {
			testFailures[testInfo.Name] = struct{}{}
		}
	}
	return testFailures
}

func testSuccessSet(testInfoList []TestInfo) map[string]struct{} {
	testFailures := make(map[string]struct{})
	for _, testInfo := range testInfoList {
		if testInfo.Status == "SUCCESS" {
			testFailures[testInfo.Name] = struct{}{}
		}
	}
	return testFailures
}
