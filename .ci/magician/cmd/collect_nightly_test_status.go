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
	"fmt"
	"magician/cloudstorage"
	"magician/provider"
	"magician/teamcity"
	utils "magician/utility"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	nightlyDataBucket = "nightly-test-data"
	tcTimeFormat      = "20060102T150405Z0700"
)

var cntsRequiredEnvironmentVariables = [...]string{
	"TEAMCITY_TOKEN",
}

type TestInfo struct {
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	Service         string    `json:"service"`
	Resource        string    `json:"resource"`
	CommitSha       string    `json:"commit_sha"`
	ErrorMessage    string    `json:"error_message"`
	ErrorType       string    `json:"error_type"`
	LogLink         string    `json:"log_link"`
	ProviderVersion string    `json:"provider_version"`
	QueuedDate      time.Time `json:"queued_date"`
	StartDate       time.Time `json:"start_date"`
	FinishDate      time.Time `json:"finish_date"`
	Duration        int       `json:"duration"`
}

// collectNightlyTestStatusCmd represents the collectNightlyTestStatus command
var collectNightlyTestStatusCmd = &cobra.Command{
	Use:   "collect-nightly-test-status",
	Short: "Collects and stores nightly test status",
	Long: `This command collects nightly test status, stores the data in JSON files and upload the files to GCS.


	The command expects the following argument(s):
	1. Custom test date in YYYY-MM-DD format. default: ""(current time when the job is executed)

	It then performs the following operations:
	1. Collects nightly test status of the execution day or the specified test date (if provided)
	2. Stores the collected data in JSON files
	3. Uploads the JSON files to GCS

	The following environment variables are required:
` + listCNTSRequiredEnvironmentVariables(),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env := make(map[string]string)
		for _, ev := range cntsRequiredEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				return fmt.Errorf("did not provide %s environment variable", ev)
			}
			env[ev] = val
		}

		tc := teamcity.NewClient(env["TEAMCITY_TOKEN"])
		gcs := cloudstorage.NewClient()

		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			return fmt.Errorf("Error loading location: %s", err)
		}

		now := time.Now().In(loc)
		year, month, day := now.Date()

		customDate := args[0]
		// check if a specific date is provided
		if customDate != "" {
			parsedDate, err := time.Parse("2006-01-02", customDate) // input format YYYY-MM-DD
			if err != nil {
				return fmt.Errorf("invalid input time format: %w", err)
			}
			year, month, day = parsedDate.Date()
		}

		// Set the time to 7pm PT
		date := time.Date(year, month, day, 19, 0, 0, 0, loc)

		return execCollectNightlyTestStatus(date, tc, gcs)
	},
}

func listCNTSRequiredEnvironmentVariables() string {
	var result string
	for i, ev := range cntsRequiredEnvironmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execCollectNightlyTestStatus(now time.Time, tc TeamcityClient, gcs CloudstorageClient) error {
	lastday := now.AddDate(0, 0, -1)
	formattedStartCut := lastday.Format(time.RFC3339)
	formattedFinishCut := now.Format(time.RFC3339)
	date := now.Format("2006-01-02")

	err := createTestReport(provider.GA, tc, gcs, formattedStartCut, formattedFinishCut, date)
	if err != nil {
		return fmt.Errorf("Error getting GA nightly test status: %w", err)
	}

	err = createTestReport(provider.Beta, tc, gcs, formattedStartCut, formattedFinishCut, date)
	if err != nil {
		return fmt.Errorf("Error getting Beta nightly test status: %w", err)
	}

	return nil
}

func createTestReport(pVersion provider.Version, tc TeamcityClient, gcs CloudstorageClient, formattedStartCut, formattedFinishCut, date string) error {

	baseLocator := fmt.Sprintf("count:500,project:%s,branch:refs/heads/nightly-test,queuedDate:(date:%s,condition:before),queuedDate:(date:%s,condition:after)", pVersion.TeamCityNightlyProjectName(), formattedFinishCut, formattedStartCut)
	fields := "build(id,buildTypeId,buildConfName,webUrl,number,queuedDate,startDate,finishDate)"
	params := url.Values{}

	// Check Queued Builds
	params.Set("locator", fmt.Sprintf("%s,state:queued", baseLocator))
	queuedBuilds, err := tc.GetBuilds(params)
	if err != nil {
		return fmt.Errorf("failed to get queued builds: %w", err)
	}
	if len(queuedBuilds.Builds) > 0 {
		fmt.Printf("%s Test unfinished: there are still %d builds queued.\n", strings.ToUpper(pVersion.String()), len(queuedBuilds.Builds))
		return nil
	}

	// Check Running Builds
	params.Set("locator", fmt.Sprintf("%s,state:running,tag:cron-trigger", baseLocator))
	params.Set("fields", fields)
	runningBuilds, err := tc.GetBuilds(params)
	if err != nil {
		return fmt.Errorf("failed to get running builds: %w", err)
	}
	if len(runningBuilds.Builds) > 0 {
		fmt.Printf("%s Test unfinished: there are still %d builds running.\n", strings.ToUpper(pVersion.String()), len(runningBuilds.Builds))
		return nil
	}

	// Get all service test builds
	params.Set("locator", fmt.Sprintf("%s,state:finished,tag:cron-trigger", baseLocator))
	params.Set("fields", fields)
	builds, err := tc.GetBuilds(params)
	if err != nil {
		return fmt.Errorf("failed to get finished builds: %w", err)
	}

	var testInfoList []TestInfo
	for _, build := range builds.Builds {
		// Get service package name
		serviceName, err := convertServiceName(build.BuildTypeId)
		if err != nil {
			return fmt.Errorf("failed to convert test service name for %s: %v", build.BuildTypeId, err)
		}
		// Skip sweeper package
		if serviceName == "sweeper" {
			continue
		}

		// Get test results
		serviceTestResults, err := tc.GetTestResults(build)
		if err != nil {
			return fmt.Errorf("failed to get test results: %v", err)
		}
		if len(serviceTestResults.TestResults) == 0 {
			fmt.Printf("Service %s has no tests\n", serviceName)
			continue
		}

		for _, testResult := range serviceTestResults.TestResults {
			var errorMessage string
			var errorType string
			// Get test debug log gcs link
			logLink := fmt.Sprintf("https://storage.cloud.google.com/teamcity-logs/nightly/%s/%s/%s/debug-%s-%s-%s-%s.txt", pVersion.TeamCityNightlyProjectName(), date, build.Number, pVersion.ProviderName(), build.Number, strconv.Itoa(build.Id), testResult.Name)
			// Get concise error message for failed and skipped tests
			// Skipped tests have a status of "UNKNOWN" on TC
			if testResult.Status == "FAILURE" || testResult.Status == "UNKNOWN" {
				errorMessage = convertErrorMessage(testResult.ErrorMessage)
				errorType = categorizeError(errorMessage)
			}

			queuedTime, err := time.Parse(tcTimeFormat, build.QueuedDate)
			if err != nil {
				return fmt.Errorf("failed to parse QueuedDate: %v", err)
			}
			startTime, err := time.Parse(tcTimeFormat, build.StartDate)
			if err != nil {
				return fmt.Errorf("failed to parse StartDate: %v", err)
			}
			finishTime, err := time.Parse(tcTimeFormat, build.FinishDate)
			if err != nil {
				return fmt.Errorf("failed to parse FinishDate: %v", err)
			}

			testInfoList = append(testInfoList, TestInfo{
				Name:            testResult.Name,
				Status:          testResult.Status,
				Service:         serviceName,
				Resource:        convertTestNameToResource(testResult.Name),
				CommitSha:       build.Number,
				ErrorMessage:    errorMessage,
				ErrorType:       errorType,
				LogLink:         logLink,
				ProviderVersion: strings.ToUpper(pVersion.String()),
				Duration:        testResult.Duration,
				QueuedDate:      queuedTime,
				StartDate:       startTime,
				FinishDate:      finishTime,
			})
		}
	}

	// Write test status data to a JSON file
	fmt.Println("Write test status")
	testStatusFileName := fmt.Sprintf("%s-%s.json", date, pVersion.String())
	err = utils.WriteToJson(testInfoList, testStatusFileName)
	if err != nil {
		return err
	}

	// Upload test status data file to gcs bucket
	objectName := fmt.Sprintf("test-metadata/%s/%s", pVersion.String(), testStatusFileName)
	err = gcs.WriteToGCSBucket(nightlyDataBucket, objectName, testStatusFileName)
	if err != nil {
		return err
	}

	return nil
}

// convertServiceName extracts service package name from teamcity build type id
// input: TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS_GOOGLE_PACKAGE_SECRETMANAGER
// output: secretmanager
func convertServiceName(servicePath string) (string, error) {
	idx := strings.LastIndex(servicePath, "_")

	if idx != -1 {
		return strings.ToLower(servicePath[idx+1:]), nil
	}
	return "", fmt.Errorf("wrong service path format for %s", servicePath)
}

// convertErrorMessage returns concise error message
func convertErrorMessage(rawErrorMessage string) string {

	startMarker := "------- Stdout: -------"
	endMarker := "------- Stderr: -------"
	startIndex := strings.Index(rawErrorMessage, startMarker)
	endIndex := strings.Index(rawErrorMessage, endMarker)

	if startIndex != -1 {
		startIndex += len(startMarker)
	} else {
		startIndex = 0
	}

	if endIndex == -1 {
		endIndex = len(rawErrorMessage)
	}

	return strings.TrimSpace(rawErrorMessage[startIndex:endIndex])
}

var (
	reSubnetNotReady   = regexp.MustCompile(`The resource '[^']+/subnetworks/[^']+' is not ready`)
	reApiEnv           = regexp.MustCompile(`has not been used in project (ci-test-project-188019|1067888929963|ci-test-project-nightly-ga|594424405950|ci-test-project-nightly-beta|653407317329|tf-vcr-private|808590572184) before or it is disabled`)
	reAttrSet          = regexp.MustCompile(`Attribute '[^']+' expected to be set`)
	reQuotaLimit       = regexp.MustCompile(`Quota limit '[^']+' has been exceeded`)
	reGoogleApi4xx     = regexp.MustCompile(`googleapi: Error 4\d\d`)
	reGoogleApi5xx     = regexp.MustCompile(`googleapi: Error 5\d\d`)
	reGoogleApiGeneric = regexp.MustCompile(`googleapi: Error`)
)

func categorizeError(errMsg string) string {
	if strings.Contains(errMsg, "Error code 13") {
		return "Error code 13"
	}
	if strings.Contains(errMsg, "Precondition check failed") {
		return "Precondition check failed"
	}

	// Diff Category
	if strings.Contains(errMsg, "After applying this test step, the plan was not empty") ||
		strings.Contains(errMsg, "After applying this test step and performing a `terraform refresh`") ||
		strings.Contains(errMsg, "Expected a non-empty plan, but got an empty plan") ||
		strings.Contains(errMsg, "error: Check failed") {
		return "Diff"
	}

	if strings.Contains(errMsg, "timeout while waiting for state") {
		return "Operation timeout"
	}

	// Regex: Subnetwork not ready
	if reSubnetNotReady.MatchString(errMsg) {
		return "Subnetwork not ready"
	}

	// ImportStateVerify Category
	if strings.Contains(errMsg, "ImportStateVerify attributes not equivalent") ||
		strings.Contains(errMsg, "Cannot import non-existent remote object") ||
		strings.Contains(errMsg, "Error: Unexpected Import Identifier") {
		return "ImportStateVerify"
	}

	// Deprecated (Case-insensitive check)
	if strings.Contains(strings.ToLower(errMsg), "deprecated") {
		return "Deprecated"
	}

	if strings.Contains(errMsg, "Provider produced inconsistent result after apply") &&
		strings.Contains(errMsg, "Root object was present, but now absent") {
		return "Root object was present, but now absent"
	}

	if strings.Contains(errMsg, "Provider produced inconsistent final plan") {
		return "Provider produced inconsistent final plan"
	}

	// API Enablement
	if reApiEnv.MatchString(errMsg) {
		return "API enablement (Test environment)"
	}
	if strings.Contains(errMsg, "has not been used in project") && strings.Contains(errMsg, "before or it is disabled") {
		return "API enablement (Created project)"
	}

	if strings.Contains(errMsg, "does not have required permissions") {
		return "Permissions"
	}
	if strings.Contains(errMsg, "bootstrap_iam_test_utils.go") {
		return "Bootstrapping"
	}

	// Bad Config Category
	if strings.Contains(errMsg, "Inconsistent dependency lock file") ||
		strings.Contains(errMsg, "Invalid resource type") ||
		strings.Contains(errMsg, "Blocks of type") && strings.Contains(errMsg, "are not expected here") ||
		strings.Contains(errMsg, "Conflicting configuration arguments") ||
		reAttrSet.MatchString(errMsg) {
		return "Bad config"
	}

	// Quota Category
	if strings.Contains(errMsg, "Quota exhausted") ||
		strings.Contains(errMsg, "Quota exceeded") ||
		strings.Contains(errMsg, "You do not have quota") ||
		reQuotaLimit.MatchString(errMsg) {
		return "Quota"
	}

	if strings.Contains(errMsg, "does not have enough resources available") {
		return "Resource availability"
	}

	// API Create/Read/Update/Delete
	if strings.Contains(errMsg, "Error: Error waiting to create") ||
		strings.Contains(errMsg, "Error: Error waiting for Create") ||
		strings.Contains(errMsg, "Error: Error waiting for creating") ||
		strings.Contains(errMsg, "Error: Error creating") ||
		strings.Contains(errMsg, "was created in the error state") ||
		strings.Contains(errMsg, "Error: Error changing instance status after creation:") {
		return "API Create"
	}

	if strings.Contains(errMsg, "Error: Error reading") {
		return "API Read"
	}

	if strings.Contains(errMsg, "Error setting IAM policy") ||
		strings.Contains(errMsg, "Error applying IAM policy") {
		return "API IAM"
	}

	if strings.Contains(errMsg, "Error: Error waiting for Updating") ||
		strings.Contains(errMsg, "Error: Error updating") {
		return "API Update"
	}

	if strings.Contains(errMsg, "Error: Error waiting for Deleting") ||
		strings.Contains(errMsg, "Error running post-test destroy") {
		return "API Delete"
	}

	// Google API Errors (Order matters: check specific codes before generic)
	if reGoogleApi4xx.MatchString(errMsg) {
		return "API (4xx)"
	}
	if reGoogleApi5xx.MatchString(errMsg) {
		return "API (5xx)"
	}
	if reGoogleApiGeneric.MatchString(errMsg) ||
		strings.Contains(errMsg, "Error: Error when reading or editing") ||
		strings.Contains(errMsg, "Error: Error waiting") ||
		strings.Contains(errMsg, "unable to queue the operation") ||
		strings.Contains(errMsg, "Error waiting for Switching runtime") {
		return "API (Other)"
	}

	return "Other"
}

func init() {
	rootCmd.AddCommand(collectNightlyTestStatusCmd)
}
