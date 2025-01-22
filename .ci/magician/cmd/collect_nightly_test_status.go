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
	NIGHTLY_DATA_BUCKET = "nightly-test-data"
)

var cntsRequiredEnvironmentVariables = [...]string{
	"TEAMCITY_TOKEN",
}

type TestInfo struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	Service      string `json:"service"`
	ErrorMessage string `json:"error_message"`
	LogLink      string `json"log_link`
}

// collectNightlyTestStatusCmd represents the collectNightlyTestStatus command
var collectNightlyTestStatusCmd = &cobra.Command{
	Use:   "collect-nightly-test-status",
	Short: "Collects and stores nightly test status",
	Long: `This command collects nightly test status, stores data in json files and upload the files to GCS.


	The command expects the following argument(s):
	1. Custom test date in YYYY-MM-DD format

	It then performs the following operations:
	1. Collect nightly test status of the execution day or a given day if Test date is provided
	2. Stores data in json files
	3. Upload the files to GCS

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
		if customDate != "" {
			parsedDate, err := time.Parse("2006-01-02", customDate) // input format YYYY-MM-DD
			// Set the time to 7pm PT
			date = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 18, 0, 0, 0, loc)
			fmt.Println("parsedDate.Location(): ", parsedDate.Location())
			if err != nil {
				return fmt.Errorf("invalid input time format: %w", err)
			}
		}
		fmt.Println("now: ", now)
		fmt.Println("date: ", date)

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
		serviceTestResults, err := tc.GetTestResults(build)
		if err != nil {
			return fmt.Errorf("failed to get test results: %v", err)
		}
		if len(serviceTestResults.TestResults) == 0 {
			continue
		}
		// Get service package name
		serviceName, err := convertServiceName(build.BuildTypeId)
		if err != nil {
			return fmt.Errorf("failed to convert test service name for %s: %v", build.BuildTypeId, err)
		}

		for _, testResult := range serviceTestResults.TestResults {
			var errorMessage string
			// compute test debug log gcs link
			logLink := fmt.Sprintf("https://storage.cloud.google.com/teamcity-logs/nightly/%s/%s/%s/debug-%s-%s-%s-%s.txt", pVersion.TeamCityNightlyProjectName(), date, build.Number, pVersion.ProviderName(), build.Number, strconv.Itoa(build.Id), testResult.Name)
			// Get concise error message
			if testResult.Status == "FAILURE" {
				errorMessage = convertErrorMessage(testResult.ErrorMessage)
			}
			testInfoList = append(testInfoList, TestInfo{
				Name:         testResult.Name,
				Status:       testResult.Status,
				Service:      serviceName,
				ErrorMessage: errorMessage,
				LogLink:      logLink,
			})
		}
	}

	fmt.Println("Write test status")
	testStatusFileName := fmt.Sprintf("%s-%s.json", date, pVersion.String())
	err = utils.WriteToJson(testInfoList, testStatusFileName)
	if err != nil {
		return err
	}

	objectName := pVersion.String() + "/" + testStatusFileName
	err = gcs.WriteToGCSBucket(NIGHTLY_DATA_BUCKET, objectName, testStatusFileName)
	if err != nil {
		return err
	}

	return nil
}

func convertServiceName(servicePath string) (string, error) {
	idx := strings.LastIndex(servicePath, "_")

	if idx != -1 {
		return strings.ToLower(servicePath[idx+1:]), nil
	}
	return "", fmt.Errorf("wrong service path format for %s", servicePath)
}

func convertErrorMessage(rawErrorMessage string) string {

	startMarker := "------- Stdout: -------"
	endMarker := "------- Stderr: -------"
	startIndex := strings.Index(rawErrorMessage, startMarker)
	endIndex := strings.Index(rawErrorMessage, endMarker)

	if startIndex != -1 {
		startIndex += len(startMarker)
	}

	if endIndex == -1 {
		endIndex = len(rawErrorMessage)
	}

	return strings.TrimSpace(rawErrorMessage[startIndex:endIndex])
}

func init() {
	rootCmd.AddCommand(collectNightlyTestStatusCmd)
}
