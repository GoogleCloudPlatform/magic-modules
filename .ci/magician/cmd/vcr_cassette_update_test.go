package cmd

import (
	"fmt"
	"magician/source"
	"magician/vcr"
	"os"
	"strings"
	"testing"

	_ "embed"

	"github.com/google/go-cmp/cmp"
)

func TestFormatVCRCassettesUpdateReplaying(t *testing.T) {
	tests := []struct {
		name string
		data vcrCassetteUpdateReplayingResult
		want string
	}{
		{
			name: "replay error",
			data: vcrCassetteUpdateReplayingResult{
				ReplayingErr: fmt.Errorf("some error"),
				ReplayingResult: vcr.Result{
					PassedTests:  []string{"a", "b"},
					FailedTests:  []string{"c", "d"},
					SkippedTests: []string{"e"},
				},
				AllReplayingPassed: false,
			},
			want: strings.Join(
				[]string{
					"#################################",
					"Tests Analytics",
					"#################################",
					"",
					"Total tests: 5",
					"Passed tests: 2",
					"Skipped tests: 1",
					"Affected tests: 2",
					"",
					"Affected tests list:",
					"- c",
					"- d",
					"",
					"",
					"#################################",
					"",
					"#################################",
					"Errors occurred during REPLAYING mode.",
					"#################################",
				},
				"\n",
			),
		},
		{
			name: "replay success",
			data: vcrCassetteUpdateReplayingResult{
				ReplayingResult: vcr.Result{
					PassedTests:  []string{"a", "b"},
					SkippedTests: []string{"e"},
				},
				AllReplayingPassed: true,
			},
			want: strings.Join(
				[]string{
					"#################################",
					"Tests Analytics",
					"#################################",
					"",
					"Total tests: 3",
					"Passed tests: 2",
					"Skipped tests: 1",
					"Affected tests: 0",
					"",
					"#################################",
					"",
					"",
					"#################################",
					"All tests passed in REPLAYING mode.",
					"#################################",
				},
				"\n",
			),
		},
		{
			name: "replay failure without error",
			data: vcrCassetteUpdateReplayingResult{
				ReplayingResult: vcr.Result{
					PassedTests:  []string{"a", "b"},
					FailedTests:  []string{"c", "d"},
					SkippedTests: []string{"e"},
				},
				AllReplayingPassed: false,
			},
			want: strings.Join(
				[]string{
					"#################################",
					"Tests Analytics",
					"#################################",
					"",
					"Total tests: 5",
					"Passed tests: 2",
					"Skipped tests: 1",
					"Affected tests: 2",
					"",
					"Affected tests list:",
					"- c",
					"- d",
					"",
					"",
					"#################################",
				},
				"\n",
			),
		},
		{
			name: "replay panic",
			data: vcrCassetteUpdateReplayingResult{
				ReplayingResult: vcr.Result{
					PassedTests:  []string{"a", "b"},
					FailedTests:  []string{"c", "d"},
					SkippedTests: []string{"e"},
					Panics:       []string{"f", "g"},
				},
				AllReplayingPassed: false,
			},
			want: strings.Join(
				[]string{
					"The provider crashed while running the VCR tests in REPLAYING mode",
				},
				"\n",
			),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatVCRCassettesUpdateReplaying(tc.data)
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("formatVCRCassettesUpdateReplaying() returned unexpected difference (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFormatVCRCassettesUpdateRecording(t *testing.T) {
	tests := []struct {
		name string
		data vcrCassetteUpdateRecordingResult
		want string
	}{
		{
			name: "record error",
			data: vcrCassetteUpdateRecordingResult{
				RecordingErr: fmt.Errorf("some error"),
				RecordingResult: vcr.Result{
					PassedTests: []string{"a", "b"},
					FailedTests: []string{"c", "d"},
				},
				AllRecordingPassed: false,
			},
			want: strings.Join(
				[]string{
					"#################################",
					"RECORDING Tests Report",
					"#################################",
					"",
					"",
					"Tests passed during RECORDING mode:",
					"- a",
					"- b",
					"",
					"Tests failed during RECORDING mode:",
					"- c",
					"- d",
					"",
					"",
					"",
					"",
					"Errors occurred during RECORDING mode: some error.",
				},
				"\n",
			),
		},
		{
			name: "record success",
			data: vcrCassetteUpdateRecordingResult{
				RecordingResult: vcr.Result{
					PassedTests: []string{"a", "b"},
				},
				AllRecordingPassed: true,
			},
			want: strings.Join(
				[]string{
					"#################################",
					"RECORDING Tests Report",
					"#################################",
					"",
					"",
					"Tests passed during RECORDING mode:",
					"- a",
					"- b",
					"",
					"",
					"",
					"",
					"",
					"",
					"All tests passed!",
				},
				"\n",
			),
		},
		{
			name: "record failed without error",
			data: vcrCassetteUpdateRecordingResult{
				RecordingResult: vcr.Result{
					PassedTests: []string{"a", "b"},
					FailedTests: []string{"c", "d"},
				},
				AllRecordingPassed: false,
			},
			want: strings.Join(
				[]string{
					"#################################",
					"RECORDING Tests Report",
					"#################################",
					"",
					"",
					"Tests passed during RECORDING mode:",
					"- a",
					"- b",
					"",
					"Tests failed during RECORDING mode:",
					"- c",
					"- d",
				},
				"\n",
			),
		},
		{
			name: "record panic",
			data: vcrCassetteUpdateRecordingResult{
				RecordingResult: vcr.Result{
					PassedTests: []string{"a", "b"},
					FailedTests: []string{"c", "d"},
					Panics:      []string{"e"},
				},
				AllRecordingPassed: false,
			},
			want: strings.Join(
				[]string{
					"#################################",
					"The provider crashed while running the VCR tests in RECORDING mode",
					"#################################",
				},
				"\n",
			),
		},
		{
			name: "has terminated test",
			data: vcrCassetteUpdateRecordingResult{
				RecordingResult: vcr.Result{
					PassedTests: []string{"a", "b"},
				},
				HasTerminatedTests: true,
				AllRecordingPassed: false,
			},
			want: strings.Join(
				[]string{
					"#################################",
					"RECORDING Tests Report",
					"#################################",
					"",
					"",
					"Tests passed during RECORDING mode:",
					"- a",
					"- b",
					"",
					"",
					"Several tests got terminated during RECORDING mode",
				},
				"\n",
			),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatVCRCassettesUpdateRecording(tc.data)
			if err != nil {
				t.Fatalf("Failed to format comment: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("formatVCRCassettesUpdateRecording() returned unexpected difference (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExecVCRCassetteUpdate(t *testing.T) {
	tests := []struct {
		name                  string
		simulateReplayFailure bool
	}{
		{
			name:                  "replay passed",
			simulateReplayFailure: false,
		},
		{
			name:                  "replay failed then record",
			simulateReplayFailure: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sb := newSandbox(t)
			sb.RequireAllowlist()

			// Intercepts gs:// URLs and translates to local copy operations.
			fakeGcloud := `#!/bin/bash
				if [ "$1" = "storage" ] && [ "$2" = "cp" ]; then
					ARGS=()
					for arg in "$@"; do
						if [[ ! "$arg" == -* ]] && [[ ! "$arg" == "storage" ]] && [[ ! "$arg" == "cp" ]]; then
							ARGS+=("$arg")
						fi
					done
					length=${#ARGS[@]}
					DEST=${ARGS[$length-1]}
					DEST=$(echo "$DEST" | sed 's|gs://|gs/|')
					mkdir -p "$DEST"
					for (( i=0; i<length-1; i++ )); do
						SRC=${ARGS[$i]}
						SRC=$(echo "$SRC" | sed 's|gs://|gs/|' | sed 's|/\*$||')
						cp -r "$SRC"/* "$DEST" 2>/dev/null || cp -r "$SRC" "$DEST" 2>/dev/null
					done
				fi`
			sb.Runner.WriteFile("gcloud", fakeGcloud)
			sb.Runner.MustRun("chmod", []string{"+x", "gcloud"}, nil)

			// Safely mocks git clone by making the destination directory.
			fakeGit := `#!/bin/bash
if [ "$1" = "clone" ]; then
  mkdir -p "$3"
  exit 0
fi
exec /usr/bin/git "$@"`
			sb.Runner.WriteFile("git", fakeGit)
			sb.Runner.MustRun("chmod", []string{"+x", "git"}, nil)

			// Simulates go list and integration tests natively via environment variables.
			fakeGo := `#!/bin/bash
				if [ "$1" = "list" ]; then
					echo "github.com/hashicorp/terraform-provider-google-beta/google-beta"
					exit 0
				elif [ "$1" = "test" ]; then
					if [ "$VCR_MODE" = "REPLAYING" ] && [ "$SIMULATE_REPLAY_FAILURE" = "true" ]; then
						echo "--- FAIL: TestAccContainerNodePool_defaultDriverInstallation (590.29s)"
						exit 1
					elif [ "$VCR_MODE" = "RECORDING" ]; then
						echo "--- PASS: TestAccContainerNodePool_defaultDriverInstallation (590.29s)"
						# Write a dummy cassette so we can verify it gets uploaded!
						echo "data" > "$VCR_PATH/new-cassette.txt"
						exit 0
					fi
					exit 0
				fi`
			sb.Runner.WriteFile("go", fakeGo)
			sb.Runner.MustRun("chmod", []string{"+x", "go"}, nil)

			testEnv := map[string]string{
				"SA_KEY": "sa_key",
			}
			if tc.simulateReplayFailure {
				testEnv["SIMULATE_REPLAY_FAILURE"] = "true"
			} else {
				testEnv["SIMULATE_REPLAY_FAILURE"] = "false"
			}

			// Setup dummy cassettes for the test to copy
			sb.Runner.MustRun("mkdir", []string{"-p", "gs/ci-vcr-cassettes/beta/fixtures"}, nil)
			sb.Runner.WriteFile("gs/ci-vcr-cassettes/beta/fixtures/dummy-cassette.txt", "data")

			ctlr := source.NewController("gopath", "hashicorp", "token", sb.Runner)
			vt, err := vcr.NewTester(testEnv, "ci-vcr-cassettes", "", sb.Runner.(vcr.ExecRunner), false)
			if err != nil {
				t.Fatalf("Failed to create new tester: %v", err)
			}

			err = execVCRCassetteUpdate("buildID", "2024-07-08", sb.Runner, ctlr, vt)
			if err != nil {
				t.Fatalf("execVCRCassetteUpdate returned error: %v", err)
			}

			if tc.simulateReplayFailure {
				if _, err := os.Stat(sb.Dir + "/gs/ci-vcr-cassettes/beta/fixtures/new-cassette.txt"); os.IsNotExist(err) {
					t.Fatalf("Expected newly recorded cassettes to be copied back to gs/ci-vcr-cassettes/beta/fixtures/")
				}
			} else {
				if _, err := os.Stat(sb.Dir + "/gs/vcr-nightly/beta/2024-07-08/buildID/logs/replaying/replaying_test.log"); os.IsNotExist(err) {
					t.Fatalf("Expected replaying test log to be uploaded to GCS")
				}
			}
		})
	}
}

func TestExecVCRCassetteUpdate_BuildFailure(t *testing.T) {
	sb := newSandbox(t)
	sb.RequireAllowlist()

	// Simulates a build failure during integration testing.
	fakeGo := `#!/bin/bash
		if [ "$1" = "list" ]; then
			echo "github.com/hashicorp/terraform-provider-google-beta/google-beta/services/corebilling"
			exit 0
		elif [ "$1" = "test" ]; then
			echo "FAIL	github.com/hashicorp/terraform-provider-google-beta/google-beta/services/corebilling [build failed]"
			exit 1
		fi`
	sb.Runner.WriteFile("go", fakeGo)
	sb.Runner.MustRun("chmod", []string{"+x", "go"}, nil)

	fakeGit := `#!/bin/bash
		if [ "$1" = "clone" ]; then
			mkdir -p "$3"
			exit 0
		fi
		exec /usr/bin/git "$@"`
	sb.Runner.WriteFile("git", fakeGit)
	sb.Runner.MustRun("chmod", []string{"+x", "git"}, nil)

	fakeGcloud := `#!/bin/bash
		if [ "$1" = "storage" ] && [ "$2" = "cp" ]; then
			exit 0
		fi`
	sb.Runner.WriteFile("gcloud", fakeGcloud)
	sb.Runner.MustRun("chmod", []string{"+x", "gcloud"}, nil)

	ctlr := source.NewController("gopath", "hashicorp", "token", sb.Runner)
	vt, err := vcr.NewTester(map[string]string{
		"SA_KEY": "sa_key",
	}, "ci-vcr-cassettes", "", sb.Runner.(vcr.ExecRunner), false)
	if err != nil {
		t.Fatalf("Failed to create new tester: %v", err)
	}

	err = execVCRCassetteUpdate("buildID", "2024-07-08", sb.Runner, ctlr, vt)
	if err == nil {
		t.Fatalf("execVCRCassetteUpdate expected to return error on build failure, got nil")
	}

	if !strings.Contains(err.Error(), "provider failed to build during VCR tests in REPLAYING mode") {
		t.Errorf("Unexpected error message: %v", err)
	}
}
