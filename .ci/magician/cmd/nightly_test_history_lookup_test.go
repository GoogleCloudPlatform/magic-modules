package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"magician/provider"
	"magician/vcr"

	"github.com/stretchr/testify/assert"
)

type fakeGCS struct {
	objects map[string][]byte
}

func (f *fakeGCS) WriteToGCSBucket(bucket, object, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	f.objects[bucket+"/"+object] = data
	return nil
}

func (f *fakeGCS) DownloadFile(bucket, object, filePath string) error {
	data, ok := f.objects[bucket+"/"+object]
	if !ok {
		return errors.New("object not found")
	}
	return os.WriteFile(filePath, data, 0644)
}

func TestClassifyNightlyStatus(t *testing.T) {
	const end = "2026-09-30"
	history := map[string]*NightlyTestHistory{
		"TestAccFailing":           {Failures: 10, LastStatus: "FAILURE", LastFailureDate: end},
		"TestAccRecentlyFailing":   {Passes: 25, Failures: 3, LastStatus: "FAILURE", LastFailureDate: end},
		"TestAccFlaky":             {Passes: 20, Failures: 4, LastStatus: "SUCCESS", LastFailureDate: "2026-09-29"},
		"TestAccFlakyLastFailed":   {Passes: 20, Failures: 1, LastStatus: "FAILURE", LastFailureDate: end},
		"TestAccNewFailure":        {Failures: 1, LastStatus: "FAILURE", LastFailureDate: end},
		"TestAccPassing":           {Passes: 30, LastStatus: "SUCCESS"},
		"TestAccSkipped":           {Skips: 30, LastStatus: "UNKNOWN"},
		"TestAccParent":            {Passes: 20, Failures: 5, LastStatus: "SUCCESS", LastFailureDate: "2026-09-29"},
		"TestAccWithSub/sub_fails": {Failures: 30, LastStatus: "FAILURE", LastFailureDate: end},

		// Fixed in main partway through the window: many stale failures, but green ever since.
		"TestAccFixed": {Passes: 10, Failures: 20, LastStatus: "SUCCESS", LastFailureDate: "2026-09-27"},
		// Only two green nights so far, which is too short a streak to call it fixed.
		"TestAccJustFixed": {Passes: 2, Failures: 20, LastStatus: "SUCCESS", LastFailureDate: "2026-09-28"},
		// A single stale failure in a long window is not a useful flakiness signal.
		"TestAccRareStale": {Passes: 29, Failures: 1, LastStatus: "SUCCESS", LastFailureDate: "2026-09-29"},
		// Without a last_failure_date the recently-fixed check cannot run, so fall back to rate.
		"TestAccNoFailureDate": {Passes: 10, Failures: 20, LastStatus: "SUCCESS"},
	}
	cases := map[string]string{
		"TestAccFailing":            NightlyStatusFailing,
		"TestAccRecentlyFailing":    NightlyStatusFailing,
		"TestAccFlaky":              NightlyStatusFlaky,
		"TestAccFlakyLastFailed":    NightlyStatusFlaky,
		"TestAccNewFailure":         NightlyStatusFailing,
		"TestAccPassing":            NightlyStatusPassing,
		"TestAccSkipped":            NightlyStatusNotFound,
		"TestAccMissing":            NightlyStatusNotFound,
		"TestAccParent__sub":        NightlyStatusFlaky,
		"TestAccWithSub__sub_fails": NightlyStatusFailing,
		"TestAccFixed":              NightlyStatusRecentlyFixed,
		"TestAccJustFixed":          NightlyStatusFlaky,
		"TestAccRareStale":          NightlyStatusPassing,
		"TestAccNoFailureDate":      NightlyStatusFlaky,
	}
	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, want, classifyNightlyStatus(lookupNightlyHistory(name, history), end))
		})
	}
}

// A test fixed in main mid-window must not stay labelled flakey forever just because the window
// still contains its pre-fix failures.
func TestClassifyNightlyStatusRecentlyFixedStreak(t *testing.T) {
	const end = "2026-09-30"
	for nights, want := range map[int]string{
		0: NightlyStatusFailing, // still failing on the most recent night
		1: NightlyStatusFlaky,
		2: NightlyStatusFlaky,
		3: NightlyStatusRecentlyFixed,
		9: NightlyStatusRecentlyFixed,
	} {
		lastFailure := time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -nights)
		h := &NightlyTestHistory{
			Passes:          nights,
			Failures:        20,
			LastStatus:      "SUCCESS",
			LastFailureDate: lastFailure.Format("2006-01-02"),
		}
		if nights == 0 {
			h.LastStatus = "FAILURE"
		}
		assert.Equal(t, want, classifyNightlyStatus(h, end), "%d nights since last failure", nights)
	}

	// A missing history end date disables the check rather than misreporting.
	assert.Equal(t, NightlyStatusFlaky, classifyNightlyStatus(
		&NightlyTestHistory{Passes: 10, Failures: 20, LastStatus: "SUCCESS", LastFailureDate: "2026-09-01"}, ""))
}

func TestLoadNightlyTestHistory(t *testing.T) {
	report := NightlyTestHistoryReport{
		Tests: map[string]*NightlyTestHistory{"TestAccA": {Passes: 1, LastStatus: "SUCCESS"}},
	}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	gcs := &fakeGCS{objects: map[string][]byte{
		nightlyDataBucket + "/" + nightlyTestHistoryObjectName(provider.Beta): data,
	}}

	got, err := loadNightlyTestHistory(provider.Beta, gcs)
	assert.NoError(t, err)
	assert.Equal(t, report.Tests, got.Tests)

	_, err = loadNightlyTestHistory(provider.GA, gcs)
	assert.Error(t, err)

	got, err = loadNightlyTestHistory(provider.Beta, nil)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

func TestBuildVCRTestRowsNightlyStatus(t *testing.T) {
	replaying := vcr.Result{FailedTests: []string{"TestAccA", "TestAccB", "TestAccC"}}
	recording := vcr.Result{PassedTests: []string{"TestAccA"}, FailedTests: []string{"TestAccB", "TestAccC"}}
	replayingAfter := vcr.Result{PassedTests: []string{"TestAccA"}}
	history := &NightlyTestHistoryReport{
		EndDate: "2026-09-30",
		Tests: map[string]*NightlyTestHistory{
			"TestAccA": {Failures: 30, LastStatus: "FAILURE", LastFailureDate: "2026-09-30"},
			"TestAccB": {Failures: 30, LastStatus: "FAILURE", LastFailureDate: "2026-09-30"},
		},
	}

	rows := buildVCRTestRows(replaying, recording, replayingAfter, "https://logs", history)
	got := map[string]string{}
	for _, r := range rows {
		got[r.DisplayName] = r.NightlyStatus
	}
	assert.Equal(t, map[string]string{
		"TestAccA": "", // only tests failing in recording are annotated
		"TestAccB": NightlyStatusFailing,
		"TestAccC": NightlyStatusNotFound,
	}, got)

	for _, r := range buildVCRTestRows(replaying, recording, replayingAfter, "https://logs", nil) {
		assert.Empty(t, r.NightlyStatus)
	}
}

func TestNightlyFailureRate(t *testing.T) {
	history := map[string]*NightlyTestHistory{
		"TestAccAlways":  {Failures: 30},
		"TestAccNever":   {Passes: 30},
		"TestAccFlaky":   {Passes: 18, Failures: 12},
		"TestAccRounded": {Passes: 2, Failures: 1},
		// Skipped runs neither pass nor fail, so they are excluded from the rate.
		"TestAccSkips":     {Passes: 9, Failures: 1, Skips: 20},
		"TestAccOnlySkips": {Skips: 30},
		"TestAccNoRuns":    {},
		"TestAccParent":    {Passes: 27, Failures: 3},
		// A last failure date is appended so a high rate is not misread as the current state.
		"TestAccDated": {Passes: 10, Failures: 20, LastFailureDate: "2026-09-27"},
		// A clean test has no failure to date, so no suffix.
		"TestAccCleanDated": {Passes: 30, LastFailureDate: "2026-09-27"},
	}

	cases := map[string]string{
		"TestAccAlways":       "30/30 nightly runs failed (100%)",
		"TestAccNever":        "0/30 nightly runs failed (0%)",
		"TestAccFlaky":        "12/30 nightly runs failed (40%)",
		"TestAccRounded":      "1/3 nightly runs failed (33%)",
		"TestAccSkips":        "1/10 nightly runs failed (10%)",
		"TestAccOnlySkips":    "",
		"TestAccNoRuns":       "",
		"TestAccMissing":      "",
		"TestAccParent__sub1": "3/30 nightly runs failed (10%)",
		"TestAccDated":        "20/30 nightly runs failed (67%), last failed 2026-09-27",
		"TestAccCleanDated":   "0/30 nightly runs failed (0%)",
	}
	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, want, nightlyFailureRate(lookupNightlyHistory(name, history)))
		})
	}
}

func TestRecordReplayNightlyColumn(t *testing.T) {
	data := recordReplay{
		TestRows: []VCRTestTableRow{
			{DisplayName: "TestAcc_a", RecordingStatus: "Failed", ReplayingAfterRecordingStatus: "-", NightlyStatus: NightlyStatusFailing, NightlyFailureRate: "27/30 nightly runs failed (90%)"},
			{DisplayName: "TestAcc_b", RecordingStatus: "Failed", ReplayingAfterRecordingStatus: "-", NightlyStatus: NightlyStatusPassing, NightlyFailureRate: "0/30 nightly runs failed (0%)"},
			{DisplayName: "TestAcc_c", RecordingStatus: "Failed", ReplayingAfterRecordingStatus: "-", NightlyStatus: NightlyStatusRecentlyFixed, NightlyFailureRate: "20/30 nightly runs failed (67%), last failed 2026-09-27"},
		},
		RecordingResult:      vcr.Result{FailedTests: []string{"TestAcc_a", "TestAcc_b", "TestAcc_c"}},
		HasNightlyHistory:    true,
		NightlyKnownFailures: 1,
		Version:              provider.Beta.String(),
		Head:                 "auto-pr-123",
		LogBucket:            "ci-vcr-logs",
	}
	got, err := formatRecordReplay(data, new(strings.Builder))
	assert.NoError(t, err)
	assert.Contains(t, got, "| Recording Mode | Replaying Rerun | Nightly | Test Name |")
	assert.Contains(t, got, "| ❌ | - | 🔴 Failing in nightly<br>27/30 nightly runs failed (90%) | TestAcc_a |")
	assert.Contains(t, got, "| ❌ | - | 🟢 Passing in nightly<br>0/30 nightly runs failed (0%) | TestAcc_b |")
	// The date keeps the high rate from contradicting the "recently fixed" label.
	assert.Contains(t, got, "| ❌ | - | 🟢 Recently fixed in nightly<br>20/30 nightly runs failed (67%), last failed 2026-09-27 | TestAcc_c |")
	assert.Contains(t, got, "**Known Nightly Failures**: 1 of the tests")

	data.HasNightlyHistory = false
	data.NightlyKnownFailures = 0
	got, err = formatRecordReplay(data, new(strings.Builder))
	assert.NoError(t, err)
	assert.Contains(t, got, "| Recording Mode | Replaying Rerun | Test Name |")
	assert.NotContains(t, got, "Nightly")
}
