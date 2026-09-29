package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

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
	history := map[string]*NightlyTestHistory{
		"TestAccFailing":           {Failures: 10, LastStatus: "FAILURE"},
		"TestAccRecentlyFailing":   {Passes: 25, Failures: 3, LastStatus: "FAILURE"},
		"TestAccFlaky":             {Passes: 20, Failures: 2, LastStatus: "SUCCESS"},
		"TestAccFlakyLastFailed":   {Passes: 20, Failures: 1, LastStatus: "FAILURE"},
		"TestAccNewFailure":        {Failures: 1, LastStatus: "FAILURE"},
		"TestAccPassing":           {Passes: 30, LastStatus: "SUCCESS"},
		"TestAccSkipped":           {Skips: 30, LastStatus: "UNKNOWN"},
		"TestAccParent":            {Passes: 30, Failures: 5, LastStatus: "SUCCESS"},
		"TestAccWithSub/sub_fails": {Failures: 30, LastStatus: "FAILURE"},
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
	}
	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, want, classifyNightlyStatus(lookupNightlyHistory(name, history)))
		})
	}
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
	assert.Equal(t, report.Tests, got)

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
	history := map[string]*NightlyTestHistory{
		"TestAccA": {Failures: 30, LastStatus: "FAILURE"},
		"TestAccB": {Failures: 30, LastStatus: "FAILURE"},
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
		},
		RecordingResult:      vcr.Result{FailedTests: []string{"TestAcc_a", "TestAcc_b"}},
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
	assert.Contains(t, got, "**Known Nightly Failures**: 1 of the tests")

	data.HasNightlyHistory = false
	data.NightlyKnownFailures = 0
	got, err = formatRecordReplay(data, new(strings.Builder))
	assert.NoError(t, err)
	assert.Contains(t, got, "| Recording Mode | Replaying Rerun | Test Name |")
	assert.NotContains(t, got, "Nightly")
}
