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
	utils "magician/utility"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v68/github"
	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/labeler"

	_ "embed"
)

var (
	//go:embed templates/TEST_FAILURE_ISSUE.md.tmpl
	testFailureIssueTemplate string
)

const (
	GithubOwner = "hashicorp"
	GithubRepo  = "terraform-provider-google"
	TotalDays   = 7
)

var ctftRequiredEnvironmentVariables = [...]string{
	"GITHUB_TOKEN",
}

type testFailureRateLabel int64

const (
	testFailureNone testFailureRateLabel = iota // 0% failure rate
	testFailure0                                // 0% < failure rate < 10%
	testFailure10                               // 10% <= failure rate < 50%
	testFailure50                               // 50% <= failure rate < 100%
	testFailure100                              // 100% failure rate
)

func (s testFailureRateLabel) String() string {
	switch s {
	case testFailure0:
		return "test-failure-0"
	case testFailure10:
		return "test-failure-10"
	case testFailure50:
		return "test-failure-50"
	case testFailure100:
		return "test-failure-100"
	default:
		return fmt.Sprintf("%d", s)
	}
}

type testFailure struct {
	TestName          string
	AffectedResource  string
	DebugLogLinks     map[provider.Version]string
	ErrorMessageLinks map[provider.Version]string
	FailureRates      map[provider.Version]string
	FailureRateLabels map[provider.Version]testFailureRateLabel
}

// createTestFailureTicketCmd represents the createTestFailureTicket command
var createTestFailureTicketCmd = &cobra.Command{
	Use:   "create-test-failure-ticket",
	Short: "Creates GitHub test failure tickets",
	Long: `This command creates GitHub test failure tickets based on nighlty test status.

	  It then performs the following operations:
	  1. Calculates test failure rate for last 7 days for all tests.
	  2. Identifies test that 
	  		a. failed 100% in last 3 days, or 
			b. failed 50%+ in last 7 days
	  3. Retrieves existing active and recently closed(within 24 hours) test failure tickets
	  4. Creates new tickets for identified failing tests detected in step 3 that don't already have a corresponding ticket.
  
	  The following environment variables are required:
  ` + listCTFTRequiredEnvironmentVariables(),
	RunE: func(cmd *cobra.Command, args []string) error {
		env := make(map[string]string)
		for _, ev := range ctftRequiredEnvironmentVariables {
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

		return execCreateTestFailureTicket(date, gh, gcs)
	},
}

func listCTFTRequiredEnvironmentVariables() string {
	var result string
	for i, ev := range ctftRequiredEnvironmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execCreateTestFailureTicket(now time.Time, gh *github.Client, gcs CloudstorageClient) error {
	ctx := context.Background()

	gaTestFailuresMap := make(map[string][]bool)
	betaTestFailuresMap := make(map[string][]bool)
	testFailuresToday := make(map[string]*testFailure)

	// Get last N-day test status map and today's failing test list
	lastNDaysTestFailureMap(provider.GA, TotalDays, now, gcs, gaTestFailuresMap, testFailuresToday)
	lastNDaysTestFailureMap(provider.Beta, TotalDays, now, gcs, betaTestFailuresMap, testFailuresToday)

	// Calculate failure rate
	for tName, tFailure := range testFailuresToday {
		if _, ok := gaTestFailuresMap[tName]; ok {
			gaRate, gaRateLabel := testFailureRate(gaTestFailuresMap[tName])
			tFailure.FailureRates[provider.GA] = gaRate
			tFailure.FailureRateLabels[provider.GA] = gaRateLabel
		}
		if _, ok := betaTestFailuresMap[tName]; ok {
			betaRate, betaRateLabel := testFailureRate(betaTestFailuresMap[tName])
			tFailure.FailureRates[provider.Beta] = betaRate
			tFailure.FailureRateLabels[provider.Beta] = betaRateLabel
		}
	}

	// Get existing GitHub test failure issues
	existTestNames, err := failingTestNamesFromActiveIssues(ctx, gh)
	if err != nil {
		return fmt.Errorf("error getting active test failure issues: %w", err)
	}

	// Get Github test failue issues closed (fixed) today
	closedTestNames, err := failingTestNamesFromClosedIssuesToday(ctx, gh, now)
	if err != nil {
		return fmt.Errorf("error getting today's closed test failure issues: %w", err)
	}

	// Create tickets
	for _, testFailure := range testFailuresToday {
		if shouldCreateTicket(*testFailure, existTestNames, closedTestNames) {
			err := createTicket(ctx, gh, testFailure)
			if err != nil {
				return fmt.Errorf("error creating test failure ticket: %w", err)
			}
		}
	}
	return nil
}

func lastNDaysTestFailureMap(pVersion provider.Version, n int, now time.Time, gcs CloudstorageClient, testFailuresMap map[string][]bool, testFailuresToday map[string]*testFailure) error {
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
			testFailuresMap[testName][i] = testInfo.Status == "FAILURE"

			if i == 0 && testInfo.Status == "FAILURE" {
				if _, ok := testFailuresToday[testName]; !ok {
					testFailuresToday[testName] = &testFailure{
						TestName:          testName,
						AffectedResource:  convertTestNameToResource(testName),
						ErrorMessageLinks: map[provider.Version]string{provider.GA: "", provider.Beta: ""},
						DebugLogLinks:     map[provider.Version]string{provider.GA: "", provider.Beta: ""},
						FailureRates:      map[provider.Version]string{provider.GA: "N/A", provider.Beta: "N/A"},
						FailureRateLabels: map[provider.Version]testFailureRateLabel{provider.GA: testFailure0, provider.Beta: testFailure0},
					}
				}
				// store error message
				d := date.Format("2006-01-02")
				fileName := fmt.Sprintf("%s-%s-%s.txt", testName, pVersion, d)
				errorMessageLink, err := storeErrorMessage(pVersion, gcs, testInfo.ErrorMessage, fileName, d)
				if err != nil {
					return err
				}
				testFailuresToday[testName].ErrorMessageLinks[pVersion] = errorMessageLink
				testFailuresToday[testName].DebugLogLinks[pVersion] = testInfo.LogLink
			}
		}
	}
	return nil
}

func testFailureRate(testFailures []bool) (string, testFailureRateLabel) {
	if testFailures == nil || len(testFailures) == 0 {
		return "N/A", testFailure0
	}
	n := len(testFailures)
	failCount := 0
	last3DaysFailed := true

	for i, fail := range testFailures {
		if fail {
			failCount++
		} else if i < 3 {
			last3DaysFailed = false
		}
	}

	// test passed consistently for last n days
	if failCount == 0 {
		return "0%", testFailureNone
	}

	// test failed consistently for last 3 days
	if last3DaysFailed {
		return "100%", testFailure100
	}

	failRate := float64(failCount) / float64(n)
	failRateStr := fmt.Sprintf("%.0f%%", failRate*100)

	switch {
	case failRate >= 0.5:
		return failRateStr, testFailure50
	case failRate >= 0.1:
		return failRateStr, testFailure10
	default:
		return failRateStr, testFailure0
	}
}

func getTestInfoList(pVersion provider.Version, date time.Time, gcs CloudstorageClient) ([]TestInfo, error) {
	lookupDate := date.Format("2006-01-02")

	testStatusFileName := fmt.Sprintf("%s-%s.json", lookupDate, pVersion.String())
	objectName := fmt.Sprintf("test-metadata/%s/%s", pVersion.String(), testStatusFileName)

	var testInfoList []TestInfo
	err := gcs.DownloadFile(nightlyDataBucket, objectName, testStatusFileName)
	if err != nil {
		return testInfoList, err
	}

	err = utils.ReadFromJson(&testInfoList, testStatusFileName)
	if err != nil {
		return testInfoList, err
	}
	return testInfoList, nil
}

func shouldCreateTicket(testfailure testFailure, existTestNames []string, todayClosedTestNames []string) bool {
	if testfailure.FailureRateLabels[provider.GA] == testFailureNone && testfailure.FailureRateLabels[provider.Beta] == testFailureNone {
		return false
	}
	for _, t := range existTestNames {
		if t == testfailure.TestName {
			return false
		}
	}
	for _, t := range todayClosedTestNames {
		if t == testfailure.TestName {
			return false
		}
	}

	if testfailure.FailureRateLabels[provider.GA] >= testFailure50 || testfailure.FailureRateLabels[provider.Beta] >= testFailure50 {
		return true
	}

	return false
}

func convertTestNameToResource(testName string) string {
	resourceName := strings.TrimPrefix(testName, "TestAcc")

	parts := strings.Split(resourceName, "_")
	if len(parts) > 0 {
		resourceName = parts[0]
	}

	// Handle datasources
	if strings.HasPrefix(resourceName, "DataSource") {
		resourceName = strings.TrimPrefix(resourceName, "DataSource")
	} else {
		resourceName = "Google" + resourceName
	}

	// Handle IAM resources
	re := regexp.MustCompile(`(IamMember|IamBinding|IamPolicy)Generated$`)
	resourceName = re.ReplaceAllString(resourceName, "${1}")

	// Convert camel case to snake case
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	resourceName = matchFirstCap.ReplaceAllString(resourceName, "${1}_${2}")
	resourceName = matchAllCap.ReplaceAllString(resourceName, "${1}_${2}")

	resourceName = strings.ToLower(resourceName)

	if val, ok := resourceNameConverter[resourceName]; ok {
		return val
	}

	return resourceName
}

func failingTestNamesFromActiveIssues(ctx context.Context, gh *github.Client) ([]string, error) {

	opts := &github.IssueListByRepoOptions{
		State:       "open",
		Labels:      []string{"test-failure"},
		ListOptions: github.ListOptions{PerPage: 100},
	}
	issues, err := ListIssuesWithOpts(ctx, gh, opts)
	if err != nil {
		return nil, err
	}
	tests, err := testNamesFromIssues(issues)
	if err != nil {
		return nil, err
	}

	return tests, nil

}

func failingTestNamesFromClosedIssuesToday(ctx context.Context, gh *github.Client, date time.Time) ([]string, error) {
	lastday := date.AddDate(0, 0, -1)
	opts := &github.IssueListByRepoOptions{
		State:       "closed",
		Labels:      []string{"test-failure"},
		Since:       lastday,
		ListOptions: github.ListOptions{PerPage: 100},
	}
	issues, err := ListIssuesWithOpts(ctx, gh, opts)
	if err != nil {
		return nil, err
	}
	tests, err := testNamesFromIssues(issues)
	if err != nil {
		return nil, err
	}

	return tests, nil
}

func ListIssuesWithOpts(ctx context.Context, gh *github.Client, opts *github.IssueListByRepoOptions) ([]*github.Issue, error) {

	var allIssues []*github.Issue
	for {
		issues, resp, err := gh.Issues.ListByRepo(ctx, GithubOwner, GithubRepo, opts)
		if err != nil {
			return nil, fmt.Errorf("error listing issues: %w", err)
		}

		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allIssues, nil
}

func createTicket(ctx context.Context, gh *github.Client, testFailure *testFailure) error {
	issueTitle := fmt.Sprintf("Failing test(s): %s", testFailure.TestName)
	issueBody, err := formatIssueBody(*testFailure)
	if err != nil {
		return fmt.Errorf("error formatting issue body: %w", err)
	}

	failureRatelabel := testFailure.FailureRateLabels[provider.GA].String()

	if testFailure.FailureRateLabels[provider.Beta] > testFailure.FailureRateLabels[provider.GA] {
		failureRatelabel = testFailure.FailureRateLabels[provider.Beta].String()
	}

	ticketLabels := []string{
		"size/xs",
		"test-failure",
		failureRatelabel,
	}

	// Apply service labels to forward test failure ticket automatically
	regexpLabels, err := labeler.BuildRegexLabels(labeler.EnrolledTeamsYaml)
	if err != nil {
		return fmt.Errorf("error building regex labels: %w", err)
	}

	labels := labeler.ComputeLabels([]string{testFailure.AffectedResource}, regexpLabels)
	ticketLabels = append(ticketLabels, labels...)

	issueRquest := &github.IssueRequest{
		Title:  github.String(issueTitle),
		Body:   github.String(issueBody),
		Labels: &ticketLabels,
		// Milestone: Near-Term Goals
		// https://github.com/hashicorp/terraform-provider-google/milestone/11
		Milestone: github.Int(11),
	}

	_, _, err = gh.Issues.Create(ctx, GithubOwner, GithubRepo, issueRquest)
	if err != nil {
		return fmt.Errorf("error creating issue: %w", err)
	}
	return nil
}

func formatIssueBody(testFailure testFailure) (string, error) {
	tmpl, err := template.New("issue").Parse(testFailureIssueTemplate)

	sb := new(strings.Builder)
	err = tmpl.Execute(sb, testFailure)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func testNamesFromIssues(issues []*github.Issue) ([]string, error) {
	var testNames []string
	for _, issue := range issues {
		tns, err := testNamesFromIssue(issue)
		if err != nil {
			return testNames, err
		}
		testNames = append(testNames, tns...)
	}
	return testNames, nil
}

func testNamesFromIssue(issue *github.Issue) ([]string, error) {
	var testNames []string
	if issue.IsPullRequest() {
		return testNames, nil
	}

	affectedTests := strings.ReplaceAll(issue.GetBody(), "<!-- List all impacted tests for searchability. The title of the issue can instead list one or more groups of tests, or describe the overall root cause. -->", "")
	impactTestRegexp := regexp.MustCompile(`Impacted tests:?[\r?\n]+((?:-? ?TestAcc[^\r\n]*\r?\n)*)`)
	matches := impactTestRegexp.FindStringSubmatch(affectedTests)

	if len(matches) > 1 {
		tests := strings.Split(matches[1], "\r\n")

		for _, test := range tests {
			subtests := strings.Split(test, "\n")
			for _, subtest := range subtests {
				if strings.HasPrefix(subtest, "- ") {
					subtest = strings.TrimSpace(subtest[2:])
					subtestParts := strings.Fields(subtest)
					subtest = subtestParts[0]
					testNames = append(testNames, subtest)
				} else {
					singleTestRegexp := regexp.MustCompile(`TestAcc[^\r\n]*`)
					if singleTestRegexp.MatchString(subtest) {
						testNames = append(testNames, subtest)
					}
				}
			}
		}

	}
	return testNames, nil
}

func storeErrorMessage(pVersion provider.Version, gcs CloudstorageClient, errorMessage, fileName, date string) (string, error) {
	// write error message to file
	data := []byte(errorMessage)
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write error message to file %s : %w", fileName, err)
	}

	// upload file to GCS
	objectName := fmt.Sprintf("test-errors/%s/%s/%s", pVersion.String(), date, fileName)
	err = gcs.WriteToGCSBucket(nightlyDataBucket, objectName, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to upload error message file %s to GCS bucket: %w", objectName, err)
	}

	// compute object view path
	link := fmt.Sprintf("https://storage.cloud.google.com/%s/%s", nightlyDataBucket, objectName)
	return link, nil
}

func init() {
	rootCmd.AddCommand(createTestFailureTicketCmd)
}

var (
	// TODO: add all mismatch resource names
	resourceNameConverter = map[string]string{
		"google_iam3_projects_policy_binding":        "google_iam_projects_policy_binding",
		"google_iam3_organizations_policy_binding":   "google_iam_organizations_policy_binding",
		"google_cloud_backup_dr_data_source":         "google_backup_dr_data_source",
		"google_cloud_backup_dr_backup":              "google_backup_dr_backup",
		"google_security_posture_posture_deployment": "google_securityposture_posture_deployment",
	}
)
