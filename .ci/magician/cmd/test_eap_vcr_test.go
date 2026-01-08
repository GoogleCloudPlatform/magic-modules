package cmd

import (
	"fmt"
	"strings"
	"testing"

	"magician/provider"
	"magician/vcr"
)

func TestAnalyticsCommentEAP(t *testing.T) {
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
				"Affected service packages",
				"None",
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
				"Affected service packages",
				"`svc-a`",
				"`svc-b`",
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
				"Affected service packages",
				"All service packages are affected",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatPostReplayEAP(tc.data)
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			for _, wc := range tc.wantContains {
				if !strings.Contains(got, wc) {
					t.Errorf("formatPostReplayEAP() returned %q, which does not contain %q", got, wc)
				}
			}
		})
	}
}

func TestNonExercisedTestsCommentEAP(t *testing.T) {
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
			got, err := formatPostReplayEAP(tc.data)
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			for _, wc := range tc.wantContains {
				if !strings.Contains(got, wc) {
					t.Errorf("formatPostReplayEAP() returned %q, which does not contain %q", got, wc)
				}
			}
		})
	}
}

func TestWithReplayFailedTestsEAP(t *testing.T) {
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
				"Found 2 affected test(s) by replaying old test recordings. Starting RECORDING based on the most recent commit. Affected tests",
				"`a`",
				"`b`",
				"[Get to know how VCR tests work](https://googlecloudplatform.github.io/magic-modules/develop/test/test/)",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatPostReplayEAP(tc.data)
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			for _, wc := range tc.wantContains {
				if !strings.Contains(got, wc) {
					t.Errorf("formatPostReplayEAP() returned %q, which does not contain %q", got, wc)
				}
			}
		})
	}
}

func TestWithoutReplayFailedTestsEAP(t *testing.T) {
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
			got, err := formatPostReplayEAP(tc.data)
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			for _, wc := range tc.wantContains {
				if !strings.Contains(got, wc) {
					t.Errorf("formatPostReplayEAP() returned %q, which does not contain %q", got, wc)
				}
			}
		})
	}
}
