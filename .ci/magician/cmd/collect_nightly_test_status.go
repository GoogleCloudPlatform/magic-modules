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
	"os"
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
	ErrorMessage    string    `json:"error_message"`
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

		now := time.Now()

		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			return fmt.Errorf("Error loading location: %s", err)
		}
		date := now.In(loc)
		customDate := args[0]
		// check if a specific date is provided
		if customDate != "" {
			parsedDate, err := time.Parse("2006-01-02", customDate) // input format YYYY-MM-DD
			// Set the time to 7pm PT
			date = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 19, 0, 0, 0, loc)
			if err != nil {
				return fmt.Errorf("invalid input time format: %w", err)
			}
		}

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
	// Get all service test builds
	builds, err := tc.GetBuilds(pVersion.TeamCityNightlyProjectName(), formattedFinishCut, formattedStartCut)
	if err != nil {
		return err
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
			// Get test debug log gcs link
			logLink := fmt.Sprintf("https://storage.cloud.google.com/teamcity-logs/nightly/%s/%s/%s/debug-%s-%s-%s-%s.txt", pVersion.TeamCityNightlyProjectName(), date, build.Number, pVersion.ProviderName(), build.Number, strconv.Itoa(build.Id), testResult.Name)
			// Get concise error message for failed and skipped tests
			// Skipped tests have a status of "UNKNOWN" on TC
			if testResult.Status == "FAILURE" || testResult.Status == "UNKNOWN" {
				errorMessage = convertErrorMessage(testResult.ErrorMessage)
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
				ErrorMessage:    errorMessage,
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

func init() {
	rootCmd.AddCommand(collectNightlyTestStatusCmd)
}
