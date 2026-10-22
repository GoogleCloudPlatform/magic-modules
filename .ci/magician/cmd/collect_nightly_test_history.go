/*
* Copyright 2026 Google LLC. All Rights Reserved.
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
	"errors"
	"fmt"
	"magician/cloudstorage"
	"magician/provider"
	"magician/teamcity"
	utils "magician/utility"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/spf13/cobra"
)

// Number of days of nightly history used to determine whether a test is flakey
const nightlyHistoryDays = 30

// NightlyTestHistory aggregates a test's outcomes across a window of nightly runs.
type NightlyTestHistory struct {
	Service         string `json:"service"`
	Passes          int    `json:"passes"`
	Failures        int    `json:"failures"`
	Skips           int    `json:"skips"`
	LastStatus      string `json:"last_status"`
	LastFailureDate string `json:"last_failure_date,omitempty"`
	// TestNameId is TeamCity's stable cross-build identifier, used to link to the test's
	// nightly history page.
	TestNameId string `json:"test_name_id,omitempty"`
}

// NightlyTestHistoryReport is the rolling history consumed by PR CI runs to detect flakey tests.
type NightlyTestHistoryReport struct {
	ProviderVersion string                         `json:"provider_version"`
	EndDate         string                         `json:"end_date"`
	DaysBack        int                            `json:"days_back"`
	DaysFound       int                            `json:"days_found"`
	Tests           map[string]*NightlyTestHistory `json:"tests"`
}

var collectNightlyTestHistoryCmd = &cobra.Command{
	Use:   "collect-nightly-test-history",
	Short: "Aggregates recent nightly test status into a rolling history",
	Long: `This command aggregates the daily nightly test status files stored in GCS into a rolling
	history of test outcomes, and uploads it to GCS so PR CI runs can detect flakey tests without
	querying TeamCity. Days without a daily status file in GCS (e.g. on the first run) are
	backfilled from TeamCity.

	The command expects the following argument(s):
	1. End date in YYYY-MM-DD format. default: ""(current date in America/Los_Angeles)

	The following environment variables are required:
` + listCNTSRequiredEnvironmentVariables(),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token, ok := os.LookupEnv("TEAMCITY_TOKEN")
		if !ok {
			return fmt.Errorf("did not provide TEAMCITY_TOKEN environment variable")
		}

		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			return fmt.Errorf("Error loading location: %s", err)
		}
		date := time.Now().In(loc).Format("2006-01-02")
		if args[0] != "" {
			if _, err := time.Parse("2006-01-02", args[0]); err != nil {
				return fmt.Errorf("invalid input time format: %w", err)
			}
			date = args[0]
		}

		tc := teamcity.NewClient(token)
		gcs := cloudstorage.NewClient()
		for _, pVersion := range []provider.Version{provider.GA, provider.Beta} {
			if err := createNightlyTestHistory(pVersion, tc, gcs, loc, date, nightlyHistoryDays); err != nil {
				return fmt.Errorf("Error creating %s nightly test history: %w", pVersion.String(), err)
			}
		}
		return nil
	},
}

// GCS directory holding files produced solely for the nightly test history.
const nightlyTestHistoryDir = "nightly-test-history"

// nightlyTestHistoryObjectName is the fixed GCS object holding the latest rolling history.
func nightlyTestHistoryObjectName(pVersion provider.Version) string {
	return fmt.Sprintf("%s/%s/nightly-test-history.json", nightlyTestHistoryDir, pVersion.String())
}

// createNightlyTestHistory builds a rolling history of test outcomes from the daily
// test-metadata files in GCS ending on date, backfilling missing days from TeamCity.
// A single build isn't enough to judge flakiness since it may itself be flakey.
func createNightlyTestHistory(pVersion provider.Version, tc TeamcityClient, gcs CloudstorageClient, loc *time.Location, date string, daysBack int) error {
	endDate, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return fmt.Errorf("failed to parse date %q: %w", date, err)
	}

	report := NightlyTestHistoryReport{
		ProviderVersion: strings.ToUpper(pVersion.String()),
		EndDate:         date,
		DaysBack:        daysBack,
		Tests:           make(map[string]*NightlyTestHistory),
	}

	// Iterate oldest to newest so LastStatus reflects the most recent run.
	for i := daysBack - 1; i >= 0; i-- {
		dayTime := endDate.AddDate(0, 0, -i)
		day := dayTime.Format("2006-01-02")
		dayResults, found, err := getDailyTestStatus(pVersion, tc, gcs, dayTime)
		if err != nil {
			return fmt.Errorf("failed to get %s test status for %s: %w", pVersion.String(), day, err)
		}
		if !found {
			fmt.Printf("Skipping %s history for %s: no finished nightly run\n", strings.ToUpper(pVersion.String()), day)
			continue
		}
		report.DaysFound++

		for _, t := range dayResults {
			h, ok := report.Tests[t.Name]
			if !ok {
				h = &NightlyTestHistory{Service: t.Service}
				report.Tests[t.Name] = h
			}
			switch t.Status {
			case "SUCCESS":
				h.Passes++
			case "FAILURE":
				h.Failures++
				h.LastFailureDate = day
			case "UNKNOWN":
				h.Skips++
			}
			h.LastStatus = t.Status
			// Days are processed oldest to newest, so this keeps the most recently seen id.
			if t.TestNameId != "" {
				h.TestNameId = t.TestNameId
			}
		}
	}

	if report.DaysFound == 0 {
		return fmt.Errorf("no daily test status files found for the %d days ending %s", daysBack, date)
	}

	historyFileName := fmt.Sprintf("nightly-test-history-%s.json", pVersion.String())
	if err := utils.WriteToJson(report, historyFileName); err != nil {
		return err
	}
	defer os.Remove(historyFileName)
	return gcs.WriteToGCSBucket(nightlyDataBucket, nightlyTestHistoryObjectName(pVersion), historyFileName)
}

// getDailyTestStatus returns the nightly test status for the given day, reading the daily
// file from GCS. If it doesn't exist, the file is backfilled from TeamCity via createTestReport
// (which also uploads it to GCS). found is false when no finished nightly run exists for the day.
func getDailyTestStatus(pVersion provider.Version, tc TeamcityClient, gcs CloudstorageClient, day time.Time) ([]TestInfo, bool, error) {
	date := day.Format("2006-01-02")
	fileName := fmt.Sprintf("%s-%s.json", date, pVersion.String())
	objectName := fmt.Sprintf("test-metadata/%s/%s", pVersion.String(), fileName)
	defer os.Remove(fileName)

	err := gcs.DownloadFile(nightlyDataBucket, objectName, fileName)
	if errors.Is(err, storage.ErrObjectNotExist) {
		os.Remove(fileName) // DownloadFile creates an empty local file before reading
		fmt.Printf("Backfilling %s test status for %s from TeamCity\n", strings.ToUpper(pVersion.String()), date)
		// Match the collect-nightly-test-status window: 7pm PT the previous day to 7pm PT on the day.
		finishCut := time.Date(day.Year(), day.Month(), day.Day(), 19, 0, 0, 0, day.Location())
		startCut := finishCut.AddDate(0, 0, -1)
		if err := createTestReport(pVersion, tc, gcs, startCut.Format(time.RFC3339), finishCut.Format(time.RFC3339), date); err != nil {
			return nil, false, err
		}
		// createTestReport only writes the file when the nightly run has finished.
		if _, statErr := os.Stat(fileName); errors.Is(statErr, os.ErrNotExist) {
			return nil, false, nil
		}
	} else if err != nil {
		return nil, false, err
	}

	var results []TestInfo
	if err := utils.ReadFromJson(&results, fileName); err != nil {
		return nil, false, fmt.Errorf("failed to read %s: %w", fileName, err)
	}
	return results, true, nil
}

func init() {
	rootCmd.AddCommand(collectNightlyTestHistoryCmd)
}
