package cmd

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"magician/provider"
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
			packages: map[string]struct{}{"servicename": struct{}{}},
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
				"serviceone":   struct{}{},
				"servicetwo":   struct{}{},
				"servicethree": struct{}{},
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
				"#### Tests analytics",
				"Total tests: 6",
				"Passed tests: 3",
				"Skipped tests: 2",
				"Affected tests: 1",
				"",
				"<details>",
				"<summary>Click here to see the affected service packages</summary>",
				"<blockquote>",
				"",
				"None",
				"",
				"</blockquote>",
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
				"#### Tests analytics",
				"Total tests: 6",
				"Passed tests: 3",
				"Skipped tests: 2",
				"Affected tests: 1",
				"",
				"<details>",
				"<summary>Click here to see the affected service packages</summary>",
				"<blockquote>",
				"",
				"<ul>",
				"<li>svc-a</li>",
				"<li>svc-b</li>",
				"",
				"</ul>",
				"",
				"</blockquote>",
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
				"#### Tests analytics",
				"Total tests: 6",
				"Passed tests: 3",
				"Skipped tests: 2",
				"Affected tests: 1",
				"",
				"<details>",
				"<summary>Click here to see the affected service packages</summary>",
				"<blockquote>",
				"",
				"All service packages are affected",
				"",
				"</blockquote>",
				"</details>",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatPostReplay(tc.data)
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
				"#### Non-exercised tests",
				"",
				color("red", "Tests were added that are skipped in VCR:"),
				"- beta-1",
				"- beta-2",
			},
		},
		{
			name: "with not run ga tests",
			data: postReplay{
				NotRunGATests: []string{"ga-1", "ga-2"},
			},
			wantContains: []string{
				"#### Non-exercised tests",
				"",
				"",
				"",
				color("red", "Tests were added that are GA-only additions and require manual runs:"),
				"- ga-1",
				"- ga-2",
			},
		},
		{
			name: "with not run ga tests and not run beta tests",
			data: postReplay{
				NotRunGATests:   []string{"ga-1", "ga-2"},
				NotRunBetaTests: []string{"beta-1", "beta-2"},
			},
			wantContains: []string{
				"#### Non-exercised tests",
				"",
				color("red", "Tests were added that are skipped in VCR:"),
				"- beta-1",
				"- beta-2",
				"",
				"",
				"",
				color("red", "Tests were added that are GA-only additions and require manual runs:"),
				"- ga-1",
				"- ga-2",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatPostReplay(tc.data)
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
	tests := []struct {
		name         string
		data         postReplay
		wantContains []string
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
				"<details>",
				"<summary>Found 2 affected test(s) by replaying old test recordings. Starting RECORDING based on the most recent commit. Click here to see the affected tests",
				"</summary>",
				"<blockquote>",
				"<ul>",
				"<li>a</li>",
				"<li>b</li>",
				"", // Empty line
				"</ul>",
				"</blockquote>",
				"</details>",
				"",
				"[Get to know how VCR tests work](https://googlecloudplatform.github.io/magic-modules/develop/test/test/)",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatPostReplay(tc.data)
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
				color("red", "Errors occurred during REPLAYING mode. Please fix them to complete your PR."),
				"View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/replaying_test.log)",
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
				color("green", "All tests passed!"),
				"View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/replaying_test.log)",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatPostReplay(tc.data)
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
	tests := []struct {
		name         string
		data         recordReplay
		wantContains []string
	}{
		{
			name: "ReplayingAfterRecordingResult has failed tests",
			data: recordReplay{
				RecordingResult: vcr.Result{
					PassedTests: []string{"a", "b", "c"},
					FailedTests: []string{"d", "e"},
				},
				ReplayingAfterRecordingResult: vcr.Result{
					PassedTests: []string{"a"},
					FailedTests: []string{"b", "c"},
				},
				HasTerminatedTests: true,
				RecordingErr:       fmt.Errorf("some error"),
				BuildID:            "build-123",
				LogBucket:          "ci-vcr-logs",
				Version:            provider.Beta.String(),
				Head:               "auto-pr-123",
			},
			wantContains: []string{
				color("green", "Tests passed during RECORDING mode:"),
				"`a` [[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/a.log)]",
				"`b` [[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/b.log)]",
				"`c` [[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/c.log)]",
				color("red", "Tests failed when rerunning REPLAYING mode:"),
				"`b` [[Error message](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/replaying_build_after_recording/b_replaying_test.log)] [[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/replaying_after_recording/b.log)]",
				"`c` [[Error message](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/replaying_build_after_recording/c_replaying_test.log)] [[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/replaying_after_recording/c.log)]",
				"Tests failed due to non-determinism or randomness when the VCR replayed the response after the HTTP request was made.",
				"Please fix these to complete your PR. If you believe these test failures to be incorrect or unrelated to your change, or if you have any questions, please raise the concern with your reviewer.",
				color("red", "Tests failed during RECORDING mode:"),
				"`d` [[Error message](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/recording_build/d_recording_test.log)] [[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/d.log)]",
				"`e` [[Error message](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/recording_build/e_recording_test.log)] [[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/e.log)]",
				color("red", "Several tests terminated during RECORDING mode."),
				"Errors occurred during RECORDING mode. Please fix them to complete your PR.",
				"[build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/recording_test.log)",
				"[debug log](https://console.cloud.google.com/storage/browser/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording)",
			},
		},
		{
			name: "ReplayingAfterRecordingResult does not have failed tests",
			data: recordReplay{
				RecordingResult: vcr.Result{
					PassedTests: []string{"a", "b", "c"},
				},
				ReplayingAfterRecordingResult: vcr.Result{
					PassedTests: []string{"a", "b", "c"},
				},
				AllRecordingPassed: true,
				BuildID:            "build-123",
				Head:               "auto-pr-123",
				Version:            provider.Beta.String(),
				LogBucket:          "ci-vcr-logs",
			},
			wantContains: []string{
				color("green", "Tests passed during RECORDING mode:"),
				"`a`",
				"[[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/a.log)]",
				"`b`",
				"[[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/b.log)]",
				"`c`",
				"[[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording/c.log)]",
				color("green", "No issues found for passed tests after REPLAYING rerun."),
				color("green", "All tests passed!"),
				"[build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/build-log/recording_test.log)",
				"[debug log](https://console.cloud.google.com/storage/browser/ci-vcr-logs/beta/refs/heads/auto-pr-123/artifacts/build-123/recording)",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatRecordReplay(tc.data)
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			for _, wc := range tc.wantContains {
				if !strings.Contains(got, wc) {
					t.Errorf("formatRecordReplay() return value:\n%s\n\ndoes not contain %q", got, wc)
				}
			}
		})
	}
}
