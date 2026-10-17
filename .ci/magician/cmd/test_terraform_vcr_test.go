package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"magician/github"
	"magician/provider"
	"magician/source"
	"magician/vcr"
)

func TestModifiedPackagesFromDiffs(t *testing.T) {
	for _, tc := range []struct {
		name     string
		diffs    []string
		packages map[string]struct{}
		all      bool
	}{
		{
			name:     "one-package",
			diffs:    []string{"google-beta/services/servicename/resource.go"},
			packages: map[string]struct{}{"servicename": {}},
			all:      false,
		},
		{
			name: "multiple-packages",
			diffs: []string{
				"google-beta/services/serviceone/resource.go",
				"google-beta/services/servicetwo/test-fixtures/fixture.txt",
				"google-beta/services/servicethree/resource_test.go",
			},
			packages: map[string]struct{}{
				"serviceone":   {},
				"servicetwo":   {},
				"servicethree": {},
			},
			all: false,
		},
		{
			name:     "all-packages",
			diffs:    []string{"google-beta/provider/provider.go"},
			packages: map[string]struct{}{},
			all:      true,
		},
		{
			name:     "all-packages-go-mod",
			diffs:    []string{"scripts/go.mod"},
			packages: map[string]struct{}{},
			all:      true,
		},
		{
			name:     "all-packages-go-sum",
			diffs:    []string{"go.sum"},
			packages: map[string]struct{}{},
			all:      true,
		},
		{
			name:     "no-packages",
			diffs:    []string{"website/docs/d/notebooks_runtime_iam_policy.html.markdown"},
			packages: map[string]struct{}{},
			all:      false,
		},
	} {
		if packages, all := modifiedPackages(tc.diffs, provider.Beta); !reflect.DeepEqual(packages, tc.packages) {
			t.Errorf("Unexpected packages found for test %s: %v, expected %v", tc.name, packages, tc.packages)
		} else if all != tc.all {
			t.Errorf("Unexpected value for all packages for test %s: %v, expected %v", tc.name, all, tc.all)
		}
	}
}

func TestNotRunTests(t *testing.T) {
	cases := map[string]struct {
		gaDiff, betaDiff string
		result           vcr.Result
		wantNotRunBeta   []string
		wantNotRunGa     []string
	}{
		"no diff": {
			gaDiff:   "",
			betaDiff: "",
			result: vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"no added tests": {
			gaDiff:   "+// some change",
			betaDiff: "+// some change",
			result: vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"test added and passed": {
			gaDiff:   "+func TestAccTwo(t *testing.T) {",
			betaDiff: "+func TestAccTwo(t *testing.T) {",
			result: vcr.Result{
				PassedTests: []string{"TestAccTwo"},
				FailedTests: []string{},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"multiple tests added and passed": {
			gaDiff: `+func TestAccTwo(t *testing.T) {
+func TestAccThree(t *testing.T) {`,
			betaDiff: `+func TestAccTwo(t *testing.T) {
+func TestAccThree(t *testing.T) {`,
			result: vcr.Result{
				PassedTests: []string{"TestAccTwo", "TestAccThree"},
				FailedTests: []string{},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"test added and failed": {
			gaDiff:   "+func TestAccTwo(t *testing.T) {",
			betaDiff: "+func TestAccTwo(t *testing.T) {",
			result: vcr.Result{
				PassedTests: []string{},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"tests removed and run": {
			gaDiff:   "-func TestAccOne(t *testing.T) {",
			betaDiff: "-func TestAccTwo(t *testing.T) {",
			result: vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"test added and not run": {
			gaDiff:   "+func TestAccThree(t *testing.T) {",
			betaDiff: "+func TestAccFour(t *testing.T) {",
			result: vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{"TestAccFour"},
			wantNotRunGa:   []string{"TestAccThree"},
		},
		"multiple tests added and not run": {
			gaDiff: `+func TestAccTwo(t *testing.T) {
+func TestAccThree(t *testing.T) {`,
			betaDiff: `+func TestAccTwo(t *testing.T) {
+func TestAccThree(t *testing.T) {`,
			result: vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccFour"},
			},
			wantNotRunBeta: []string{"TestAccThree", "TestAccTwo"},
			wantNotRunGa:   []string{},
		},
		"tests removed and not run": {
			gaDiff:   "-func TestAccThree(t *testing.T) {",
			betaDiff: "-func TestAccFour(t *testing.T) {",
			result: vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"tests added but commented out": {
			gaDiff:   "+//func TestAccThree(t *testing.T) {",
			betaDiff: "+//func TestAccFour(t *testing.T) {",
			result: vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{},
		},
		"multiline diffs": {
			gaDiff: `diff --git a/google/services/alloydb/resource_alloydb_backup_generated_test.go b/google/services/alloydb/resource_alloydb_backup_generated_test.go
+func TestAccAlloydbBackup_alloydbBackupFullTestNewExample(t *testing.T) {
+func TestAccCloudRunService_cloudRunServiceMulticontainerExample(t *testing.T) {`,
			betaDiff: `diff --git a/google-beta/services/alloydb/resource_alloydb_backup_generated_test.go b/google-beta/services/alloydb/resource_alloydb_backup_generated_test.go
+func TestAccAlloydbBackup_alloydbBackupFullTestNewExample(t *testing.T) {`,
			result: vcr.Result{
				PassedTests: []string{},
				FailedTests: []string{},
			},
			wantNotRunBeta: []string{"TestAccAlloydbBackup_alloydbBackupFullTestNewExample"},
			wantNotRunGa:   []string{"TestAccCloudRunService_cloudRunServiceMulticontainerExample"},
		},
		"always count GA-only added tests": {
			gaDiff:   "+func TestAccOne(t *testing.T) {",
			betaDiff: "",
			result: vcr.Result{
				PassedTests: []string{"TestAccOne"},
				FailedTests: []string{"TestAccTwo"},
			},
			wantNotRunBeta: []string{},
			wantNotRunGa:   []string{"TestAccOne"},
		},
	}
	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			notRunBeta, notRunGa := notRunTests(tc.gaDiff, tc.betaDiff, tc.result)
			assert.Equal(t, tc.wantNotRunBeta, notRunBeta)
			assert.Equal(t, tc.wantNotRunGa, notRunGa)
		})
	}
}

func TestAnalyticsComment(t *testing.T) {
	tests := []struct {
		name         string
		data         postReplay
		wantContains []string
	}{
		{
			name: "run full vcr is false and no affected services",
			data: postReplay{
				ReplayingResult: vcr.Result{
					PassedTests:  []string{"a", "b", "c"},
					SkippedTests: []string{"d", "e"},
					FailedTests:  []string{"f"},
				},
				RunFullVCR:       false,
				AffectedServices: []string{},
			},
			wantContains: []string{
				"#### Analytics",
				"| Total Tests | Passed | Skipped | Affected |",
				"| 6 | 3 | 2 | 1 |",
				"<details>",
				"<summary><b>Affected Service Packages</b></summary>",
				"* None",
				"</details>",
			},
		},
		{
			name: "run full vcr is false and has affected services",
			data: postReplay{
				ReplayingResult: vcr.Result{
					PassedTests:  []string{"a", "b", "c"},
					SkippedTests: []string{"d", "e"},
					FailedTests:  []string{"f"},
				},
				RunFullVCR:       false,
				AffectedServices: []string{"svc-a", "svc-b"},
			},
			wantContains: []string{
				"#### Analytics",
				"| Total Tests | Passed | Skipped | Affected |",
				"| 6 | 3 | 2 | 1 |",
				"<details>",
				"<summary><b>Affected Service Packages</b></summary>",
				"* svc-a",
				"* svc-b",
				"</details>",
			},
		},
		{
			name: "run full vcr is true",
			data: postReplay{
				ReplayingResult: vcr.Result{
					PassedTests:  []string{"a", "b", "c"},
					SkippedTests: []string{"d", "e"},
					FailedTests:  []string{"f"},
				},
				RunFullVCR:       true,
				AffectedServices: []string{},
			},
			wantContains: []string{
				"#### Analytics",
				"| Total Tests | Passed | Skipped | Affected |",
				"| 6 | 3 | 2 | 1 |",
				"<details>",
				"<summary><b>Affected Service Packages</b></summary>",
				"* All service packages are affected",
				"</details>",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatPostReplay(tc.data, new(strings.Builder))
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			for _, wc := range tc.wantContains {
				if !strings.Contains(got, wc) {
					t.Errorf("formatPostReplay() returned %q, which does not contain %q", got, wc)
				}
			}
		})
	}
}

func TestNonExercisedTestsComment(t *testing.T) {
	tests := []struct {
		name         string
		data         postReplay
		wantContains []string
	}{
		{
			name: "with not run beta tests",
			data: postReplay{
				NotRunBetaTests: []string{"beta-1", "beta-2"},
			},
			wantContains: []string{
				"> [!IMPORTANT]",
				"**Manual Verification Required**",
				"VCR cannot automatically execute the following tests added in this PR. Please verify them manually:",
				"> 🔴 beta-1",
				"> 🔴 beta-2",
				"> [!CAUTION]",
				"**Issues requiring attention before PR completion**",
			},
		},
		{
			name: "with not run ga tests",
			data: postReplay{
				NotRunGATests: []string{"ga-1", "ga-2"},
			},
			wantContains: []string{
				"> [!IMPORTANT]",
				"**Manual Verification Required (GA-only additions)**",
				"The following tests are GA-only additions and cannot be run by VCR in Beta mode. Please verify them manually:",
				"> 🔴 ga-1",
				"> 🔴 ga-2",
				"> [!CAUTION]",
				"**Issues requiring attention before PR completion**",
			},
		},
		{
			name: "with not run ga tests and not run beta tests",
			data: postReplay{
				NotRunGATests:   []string{"ga-1", "ga-2"},
				NotRunBetaTests: []string{"beta-1", "beta-2"},
			},
			wantContains: []string{
				"> [!IMPORTANT]",
				"**Manual Verification Required**",
				"VCR cannot automatically execute the following tests added in this PR. Please verify them manually:",
				"> 🔴 beta-1",
				"> 🔴 beta-2",
				"**Manual Verification Required (GA-only additions)**",
				"The following tests are GA-only additions and cannot be run by VCR in Beta mode. Please verify them manually:",
				"> 🔴 ga-1",
				"> 🔴 ga-2",
				"> [!CAUTION]",
				"**Issues requiring attention before PR completion**",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatPostReplay(tc.data, new(strings.Builder))
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			for _, wc := range tc.wantContains {
				if !strings.Contains(got, wc) {
					t.Errorf("formatPostReplay() returned %q, which does not contain %q", got, wc)
				}
			}
		})
	}
}

func TestWithReplayFailedTests(t *testing.T) {
	manyFailedTests := make([]string, 101)
	for i := range manyFailedTests {
		manyFailedTests[i] = "TestAccAlloydbBackup_alloydbBackupFullTestNewExample"
	}
	tests := []struct {
		name               string
		data               postReplay
		wantContains       []string
		wantNotContains    []string
		wantStdoutContains []string
	}{
		{
			name: "with failed tests",
			data: postReplay{
				ReplayingResult: vcr.Result{
					FailedTests: []string{"a", "b"},
				},
			},
			wantContains: []string{
				"#### Action taken",
				"Found 2 affected test(s) by replaying old test recordings. Starting RECORDING based on the most recent commit.",
				"<details>",
				"<summary>Click here to see the affected tests</summary>",
				"* a",
				"* b",
				"</details>",
				"[Learn how VCR tests work](https://googlecloudplatform.github.io/magic-modules/develop/test/test/)",
			},
		},
		{
			name: "Too many failed tests",
			data: postReplay{
				ReplayingResult: vcr.Result{
					FailedTests: manyFailedTests,
				},
				RunFullVCR: true,
			},
			wantContains: []string{
				"More than 100 tests failed in replaying; this is too many to display on GitHub.",
			},
			wantNotContains: []string{manyFailedTests[0]},
			wantStdoutContains: []string{
				"Failed replaying tests",
				manyFailedTests[0],
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sb := new(strings.Builder)
			got, err := formatPostReplay(tc.data, sb)
			gotStdout := sb.String()
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			for _, wc := range tc.wantContains {
				if !strings.Contains(got, wc) {
					t.Errorf("formatPostReplay() returned %q, which does not contain %q", got, wc)
				}
			}
			for _, s := range tc.wantNotContains {
				assert.NotContains(t, got, s)
			}
			for _, s := range tc.wantStdoutContains {
				assert.Contains(t, gotStdout, s)
			}
		})
	}
}

func TestWithoutReplayFailedTests(t *testing.T) {
	tests := []struct {
		name         string
		data         postReplay
		wantContains []string
	}{
		{
			name: "with replay error",
			data: postReplay{
				ReplayingErr: fmt.Errorf("some error"),
				BuildID:      "build-123",
				Head:         "auto-pr-123",
				LogBucket:    "ci-vcr-logs",
				Version:      provider.Beta.String(),
			},
			wantContains: []string{
				"> [!CAUTION]",
				"🔴 Errors occurred during REPLAYING mode. Please check the build log for details.",
				"View the [replaying VCR build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/replaying_test.log)",
			},
		},
		{
			name: "without replay error",
			data: postReplay{
				BuildID:   "build-123",
				Head:      "auto-pr-123",
				LogBucket: "ci-vcr-logs",
				Version:   provider.Beta.String(),
			},
			wantContains: []string{
				"🟢 **All tests passed in Replaying mode! No Recording was needed.**",
				"View the [replaying VCR build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/replaying_test.log)",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatPostReplay(tc.data, new(strings.Builder))
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			for _, wc := range tc.wantContains {
				if !strings.Contains(got, wc) {
					t.Errorf("formatPostReplay() returned %q, which does not contain %q", got, wc)
				}
			}
		})
	}
}

func TestRecordReplay(t *testing.T) {
	manyRows := make([]VCRTestTableRow, 101)
	for i := range manyRows {
		manyRows[i] = VCRTestTableRow{
			DisplayName:                   "TestAcc_a",
			RecordingStatus:               "Passed",
			RecordingLogUrl:               "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/a.log",
			ReplayingAfterRecordingStatus: "Passed",
		}
	}
	tests := []struct {
		name               string
		data               recordReplay
		wantContains       []string
		wantNotContains    []string
		wantStdoutContains []string
	}{
		{
			name: "ReplayingAfterRecordingResult has failed tests",
			data: recordReplay{
				TestRows: []VCRTestTableRow{
					{
						DisplayName:                   "TestAcc_a",
						RecordingStatus:               "Passed",
						RecordingLogUrl:               "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/a.log",
						ReplayingAfterRecordingStatus: "Passed",
					},
					{
						DisplayName:                     "TestAcc_b",
						RecordingStatus:                 "Passed",
						RecordingLogUrl:                 "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/b.log",
						ReplayingAfterRecordingStatus:   "Failed",
						ReplayingAfterRecordingErrorUrl: "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/replaying_build_after_recording/b_replaying_test.log",
						ReplayingAfterRecordingLogUrl:   "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/replaying_after_recording/b.log",
					},
					{
						DisplayName:                     "TestAcc_c",
						RecordingStatus:                 "Passed",
						RecordingLogUrl:                 "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/c.log",
						ReplayingAfterRecordingStatus:   "Failed",
						ReplayingAfterRecordingErrorUrl: "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/replaying_build_after_recording/c_replaying_test.log",
						ReplayingAfterRecordingLogUrl:   "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/replaying_after_recording/c.log",
					},
					{
						DisplayName:                   "TestAcc_d",
						RecordingStatus:               "Failed",
						RecordingErrorUrl:             "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/recording_build/d_recording_test.log",
						RecordingLogUrl:               "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/d.log",
						ReplayingAfterRecordingStatus: "-",
					},
					{
						DisplayName:                   "TestAcc_e",
						RecordingStatus:               "Failed",
						RecordingErrorUrl:             "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/recording_build/e_recording_test.log",
						RecordingLogUrl:               "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/e.log",
						ReplayingAfterRecordingStatus: "-",
					},
				},
				HasTerminatedTests: true,
				RecordingErr:       fmt.Errorf("some error"),
				BuildID:            "build-123",
				LogBucket:          "ci-vcr-logs",
				Version:            provider.Beta.String(),
				Head:               "auto-pr-123",
			},
			wantContains: []string{
				"| Recording Mode | Findings | Test Name |",
				"| ✅&nbsp;[Log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/a.log) | - | TestAcc_a |",
				"| ✅&nbsp;[Log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/b.log) | `❌ Replay rerun failed`&nbsp;[Error](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/replaying_build_after_recording/b_replaying_test.log)&nbsp;·&nbsp;[Log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/replaying_after_recording/b.log) | TestAcc_b |",
				"| ✅&nbsp;[Log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/c.log) | `❌ Replay rerun failed`&nbsp;[Error](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/replaying_build_after_recording/c_replaying_test.log)&nbsp;·&nbsp;[Log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/replaying_after_recording/c.log) | TestAcc_c |",
				"| ❌&nbsp;[Error](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/recording_build/d_recording_test.log)&nbsp;·&nbsp;[Log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/d.log) | - | TestAcc_d |",
				"| ❌&nbsp;[Error](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/recording_build/e_recording_test.log)&nbsp;·&nbsp;[Log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/e.log) | - | TestAcc_e |",

				"Please address these issues to complete your PR",
				"View the [recording VCR build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/recording_test.log) or the [debug logs folder](https://console.cloud.google.com/storage/browser/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording) for detailed results.",
			},
		},
		{
			name: "ReplayingAfterRecordingResult does not have failed tests",
			data: recordReplay{
				TestRows: []VCRTestTableRow{
					{
						DisplayName:                   "TestAcc_a",
						RecordingStatus:               "Passed",
						RecordingLogUrl:               "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/a.log",
						ReplayingAfterRecordingStatus: "Passed",
					},
					{
						DisplayName:                   "TestAcc_b",
						RecordingStatus:               "Passed",
						RecordingLogUrl:               "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/b.log",
						ReplayingAfterRecordingStatus: "Passed",
					},
					{
						DisplayName:                   "TestAcc_c",
						RecordingStatus:               "Passed",
						RecordingLogUrl:               "https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/c.log",
						ReplayingAfterRecordingStatus: "Passed",
					},
				},
				AllRecordingPassed: true,
				BuildID:            "build-123",
				Head:               "auto-pr-123",
				Version:            provider.Beta.String(),
				LogBucket:          "ci-vcr-logs",
			},
			wantContains: []string{
				"| Recording Mode | Findings | Test Name |",
				"| ✅&nbsp;[Log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/a.log) | - | TestAcc_a |",
				"| ✅&nbsp;[Log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/b.log) | - | TestAcc_b |",
				"| ✅&nbsp;[Log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/c.log) | - | TestAcc_c |",
				"🟢 **All tests passed!**",
				"View the [recording VCR build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/recording_test.log) or the [debug logs folder](https://console.cloud.google.com/storage/browser/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording) for detailed results.",
			},
		},
		{
			name: "Too many rows",
			data: recordReplay{
				TestRows:           manyRows,
				AllRecordingPassed: true,
				BuildID:            "build-123",
				Head:               "auto-pr-123",
				Version:            provider.Beta.String(),
				LogBucket:          "ci-vcr-logs",
			},
			wantContains: []string{
				"More than 100 tests ran in recording mode; this is too many to display on GitHub.",
			},
			wantNotContains: []string{
				manyRows[0].DisplayName,
			},
			wantStdoutContains: []string{
				"| Recording Mode | Findings | Test Name |",
				"| ✅&nbsp;[Log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/a.log) | - | TestAcc_a |",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sb := new(strings.Builder)
			got, err := formatRecordReplay(tc.data, sb)
			gotStdout := sb.String()
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			for _, s := range tc.wantContains {
				assert.Contains(t, got, s)
			}
			for _, s := range tc.wantNotContains {
				assert.NotContains(t, got, s)
			}
			for _, s := range tc.wantStdoutContains {
				assert.Contains(t, gotStdout, s)
			}
		})
	}
}

func TestHandleBuildFailures(t *testing.T) {
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
	}
	result := vcr.Result{
		BuildFailures: []string{"package1", "package2"},
	}

	sb := newSandbox(t)
	sb.RequireAllowlist()
	handled, err := handleBuildFailures("123", "build-456", "http://target", "sha789", sb.Dir, result, vcr.Replaying, gh, sb.Runner)

	assert.NoError(t, err)
	assert.True(t, handled)

	assert.Len(t, gh.calledMethods["PostComment"], 1)
	assert.Equal(t, "123", gh.calledMethods["PostComment"][0][0])
	comment := gh.calledMethods["PostComment"][0][1].(string)
	assert.Contains(t, comment, "**Step 1: Replaying Mode**")
	assert.Contains(t, comment, "package1")
	assert.Contains(t, comment, "package2")

	assert.Len(t, gh.calledMethods["PostBuildStatus"], 1)
	assert.Equal(t, "123", gh.calledMethods["PostBuildStatus"][0][0])
	assert.Equal(t, "VCR-test", gh.calledMethods["PostBuildStatus"][0][1])
	assert.Equal(t, "failure", gh.calledMethods["PostBuildStatus"][0][2])
}

func TestHandleBuildFailures_Recording(t *testing.T) {
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
	}
	result := vcr.Result{
		BuildFailures: []string{"package1", "package2"},
	}

	sb := newSandbox(t)
	sb.RequireAllowlist()
	handled, err := handleBuildFailures("123", "build-456", "http://target", "sha789", sb.Dir, result, vcr.Recording, gh, sb.Runner)

	assert.NoError(t, err)
	assert.True(t, handled)

	assert.Len(t, gh.calledMethods["PostComment"], 1)
	comment := gh.calledMethods["PostComment"][0][1].(string)
	assert.Contains(t, comment, "**Step 2: Recording Mode**")
	assert.Contains(t, comment, "package1")
	assert.Contains(t, comment, "package2")
}

func TestHandleBuildFailures_NoFailures(t *testing.T) {
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
	}
	result := vcr.Result{}

	sb := newSandbox(t)
	sb.RequireAllowlist()
	handled, err := handleBuildFailures("123", "build-456", "http://target", "sha789", sb.Dir, result, vcr.Replaying, gh, sb.Runner)

	assert.NoError(t, err)
	assert.False(t, handled)
	assert.Len(t, gh.calledMethods["PostComment"], 0)
	assert.Len(t, gh.calledMethods["PostBuildStatus"], 0)
}

func TestAppendVCRResultToDiffComment_NotExists(t *testing.T) {
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
		pullRequest: github.PullRequest{
			User: github.User{Login: "author1"},
		},
		requestedReviewers: []github.User{
			{Login: "reviewer1"},
			{Login: "reviewer2"},
		},
		pullRequestComments: []github.PullRequestComment{
			{
				ID:   456,
				Body: "Some other comment",
			},
		},
	}

	sb := newSandbox(t)
	sb.RequireAllowlist()
	err := appendVCRResultToDiffComment("123", "VCR Results", sb.Dir, gh, sb.Runner)

	assert.NoError(t, err)
	assert.Len(t, gh.calledMethods["PostComment"], 1)
	assert.Equal(t, "123", gh.calledMethods["PostComment"][0][0])
	assert.Equal(t, "VCR Results", gh.calledMethods["PostComment"][0][1])
}
func TestAppendVCRResultToDiffComment_UseFileID(t *testing.T) {
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
		pullRequest: github.PullRequest{
			User: github.User{Login: "author1"},
		},
		requestedReviewers: []github.User{
			{Login: "reviewer1"},
			{Login: "reviewer2"},
		},
		pullRequestComments: []github.PullRequestComment{
			{
				ID:   456,
				Body: "Some comment",
			},
		},
	}
	sb := newSandbox(t)
	sb.RequireAllowlist()

	sb.Runner.WriteFile("diff_comment_id.txt", "456")

	err := appendVCRResultToDiffComment("123", "VCR Results", sb.Dir, gh, sb.Runner)

	assert.NoError(t, err)
	assert.Len(t, gh.calledMethods["UpdateComment"], 1)
	assert.Equal(t, "123", gh.calledMethods["UpdateComment"][0][0])
	assert.Contains(t, gh.calledMethods["UpdateComment"][0][1].(string), "VCR Results")
	assert.Equal(t, 456, gh.calledMethods["UpdateComment"][0][2])
}

func TestIsCommitOutdated(t *testing.T) {
	cases := map[string]struct {
		mmCommitSha       string
		headSha           string
		getPullRequestErr error
		wantOutdated      bool
		wantHeadSha       string
	}{
		"superseded commit": {
			mmCommitSha:  "old-sha-123",
			headSha:      "new-sha-456",
			wantOutdated: true,
			wantHeadSha:  "new-sha-456",
		},
		"latest commit": {
			mmCommitSha: "old-sha-123",
			headSha:     "old-sha-123",
		},
		"empty head sha": {
			mmCommitSha: "old-sha-123",
		},
		"empty commit sha": {
			headSha: "new-sha-456",
		},
		"pull request lookup error fails open": {
			mmCommitSha:       "old-sha-123",
			headSha:           "new-sha-456",
			getPullRequestErr: errors.New("rate limited"),
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			gh := &mockGithub{
				calledMethods:     make(map[string][][]any),
				pullRequest:       github.PullRequest{HeadSha: tc.headSha},
				getPullRequestErr: tc.getPullRequestErr,
			}

			outdated, headSha := isCommitOutdated("123", tc.mmCommitSha, gh)

			assert.Equal(t, tc.wantOutdated, outdated)
			assert.Equal(t, tc.wantHeadSha, headSha)
			assert.Equal(t, [][]any{{"123"}}, gh.calledMethods["GetPullRequest"])
		})
	}
}

// TestExecTestTerraformVCR_OutdatedCommit pins where execTestTerraformVCR checks for a newer
// commit: at entry (before cloning/replaying), and again before recording (after posting the Step 1 report).
func TestExecTestTerraformVCR_OutdatedCommit(t *testing.T) {
	cases := map[string]struct {
		// Head SHAs returned by successive pull request lookups.
		headShas           []string
		diffCommentBody    string
		tpgbNewSha         string
		tpgbOldSha         string
		updateCommentErr   error
		postBuildStatusErr error
		wantStatuses       []string
		wantReplayed       bool
		wantUpdateCalls    int
		wantComment        string
		wantErr            bool
	}{
		"superseded before replaying": {
			headShas:        []string{"new-sha"},
			wantStatuses:    []string{"failure"},
			wantUpdateCalls: 1,
			wantComment:     "## Test report\n\n**Step 1: Replaying Mode**\n\nSkipped REPLAYING mode because a newer commit (new-sha) was pushed to the PR.",
		},
		"superseded before replaying without diffs in comment": {
			headShas:        []string{"new-sha"},
			diffCommentBody: "Your PR hasn't generated any diffs, but I'll let you know if a future commit does.",
		},
		"superseded before replaying without downstream commits": {
			headShas:   []string{"new-sha"},
			tpgbNewSha: "same-downstream-sha",
			tpgbOldSha: "same-downstream-sha",
		},
		"superseded before recording": {
			headShas:        []string{"mm-sha", "new-sha"},
			wantStatuses:    []string{"pending", "failure"},
			wantReplayed:    true,
			wantUpdateCalls: 2,
			wantComment:     "---\n\n**Step 2: Recording Mode**\n\nSkipped RECORDING mode because a newer commit (new-sha) was pushed to the PR.",
		},
		"comment error before replaying still posts status": {
			headShas:         []string{"new-sha"},
			updateCommentErr: errors.New("comment too long"),
			wantStatuses:     []string{"failure"},
			wantUpdateCalls:  1,
			wantComment:      "Skipped REPLAYING mode",
		},
		"comment error before recording still posts status": {
			headShas:         []string{"mm-sha", "new-sha"},
			updateCommentErr: errors.New("comment too long"),
			wantStatuses:     []string{"pending", "failure"},
			wantReplayed:     true,
			wantUpdateCalls:  2,
			wantComment:      "Skipped RECORDING mode",
		},
		"status error before replaying is returned": {
			headShas:           []string{"new-sha"},
			postBuildStatusErr: errors.New("forbidden"),
			wantStatuses:       []string{"failure"},
			wantUpdateCalls:    1,
			wantComment:        "Skipped REPLAYING mode",
			wantErr:            true,
		},
		"status error before recording is returned": {
			headShas:           []string{"mm-sha", "new-sha"},
			postBuildStatusErr: errors.New("forbidden"),
			wantStatuses:       []string{"pending", "failure"},
			wantReplayed:       true,
			wantUpdateCalls:    2,
			wantComment:        "Skipped RECORDING mode",
			wantErr:            true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			sb := newSandbox(t)
			sb.RequireAllowlist()
			writeScript := func(name, script string) {
				if err := os.WriteFile(filepath.Join(sb.Dir, name), []byte(script), 0755); err != nil {
					t.Fatalf("Failed to write fake %s: %v", name, err)
				}
			}
			// Cloning creates the destination directory, and the diff reports a change to one service.
			writeScript("git", `#!/bin/bash
case "$1" in
  clone) mkdir -p "${@: -1}" ;;
  diff) [[ " $* " == *" --name-only "* ]] && echo "google-beta/services/compute/resource_compute_instance.go" ;;
esac
exit 0`)
			writeScript("gcloud", "#!/bin/bash\nexit 0")
			// Records which VCR modes ran. Replaying fails, so recording would run next.
			writeScript("go", fmt.Sprintf(`#!/bin/bash
if [ "$1" = "test" ]; then
  echo > "%s/go_test_$VCR_MODE"
  if [ "$VCR_MODE" = "REPLAYING" ]; then
    echo "--- FAIL: TestAccComputeInstance_basic (1.00s)"
    exit 1
  fi
fi
exit 0`, sb.Dir))
			sb.Runner.WriteFile("diff_comment_id.txt", "456")
			if tc.tpgbNewSha != "" {
				sb.Runner.WriteFile("commitSHA_modular-magician_terraform-provider-google-beta.txt", tc.tpgbNewSha)
			}
			if tc.tpgbOldSha != "" {
				sb.Runner.WriteFile("commitSHA_modular-magician_terraform-provider-google-beta-old.txt", tc.tpgbOldSha)
			}

			commentBody := tc.diffCommentBody
			if commentBody == "" {
				commentBody = "Diff report"
			}
			var pullRequests []github.PullRequest
			for _, sha := range tc.headShas {
				pullRequests = append(pullRequests, github.PullRequest{HeadSha: sha})
			}
			gh := &mockGithub{
				calledMethods:       make(map[string][][]any),
				pullRequests:        pullRequests,
				pullRequestComments: []github.PullRequestComment{{ID: 456, Body: commentBody}},
				updateCommentErr:    tc.updateCommentErr,
				postBuildStatusErr:  tc.postBuildStatusErr,
			}
			ctlr := source.NewController(filepath.Join(sb.Dir, "gopath"), "modular-magician", "token", sb.Runner)
			vt, err := vcr.NewTester(map[string]string{"SA_KEY": "sa_key"}, "ci-vcr-cassettes", "ci-vcr-logs", sb.Runner.(vcr.ExecRunner), false)
			if err != nil {
				t.Fatalf("Failed to create VCR tester: %v", err)
			}

			err = execTestTerraformVCR("123", "mm-sha", "build-1", "project-1", "22", "main", sb.Dir, gh, sb.Runner, ctlr, vt, &fakeGCS{objects: map[string][]byte{}})
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Len(t, gh.calledMethods["GetPullRequest"], len(tc.headShas))

			var statuses []string
			for _, call := range gh.calledMethods["PostBuildStatus"] {
				assert.Equal(t, "mm-sha", call[4])
				statuses = append(statuses, call[2].(string))
			}
			assert.Equal(t, tc.wantStatuses, statuses)

			_, err = os.Stat(filepath.Join(sb.Dir, "gopath"))
			assert.Equal(t, tc.wantReplayed, err == nil, "whether repos were cloned")
			_, err = os.Stat(filepath.Join(sb.Dir, "go_test_REPLAYING"))
			assert.Equal(t, tc.wantReplayed, err == nil, "whether replaying ran")
			_, err = os.Stat(filepath.Join(sb.Dir, "go_test_RECORDING"))
			assert.True(t, os.IsNotExist(err), "recording should not run")

			if assert.Len(t, gh.calledMethods["UpdateComment"], tc.wantUpdateCalls) && tc.wantUpdateCalls > 0 {
				lastComment := gh.calledMethods["UpdateComment"][tc.wantUpdateCalls-1][1].(string)
				assert.Contains(t, lastComment, tc.wantComment)
				if tc.wantReplayed {
					assert.Contains(t, lastComment, "**Step 1: Replaying Mode**")
				}
			}
			assert.Empty(t, gh.calledMethods["PostComment"])
		})
	}
}
