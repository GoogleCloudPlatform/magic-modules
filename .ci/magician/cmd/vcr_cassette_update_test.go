package cmd

import (
	"container/list"
	"fmt"
	"magician/source"
	"magician/vcr"
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
					"Errors occurred during REPLAYING mode: some error.",
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
		name          string
		cmdResults    map[string]string
		expectedCalls map[string][]ParameterList
	}{
		{
			name:       "replay passed",
			cmdResults: make(map[string]string),
			expectedCalls: map[string][]ParameterList{
				"Run": {
					{"/mock/dir/magic-modules/.ci/magician", "gsutil", []string{"-m", "-q", "cp", "gs://ci-vcr-cassettes/beta/fixtures/*", "/mock/dir/magic-modules/.ci/magician/cassettes/beta"}, map[string]string(nil)},
					{"/mock/dir/magic-modules/.ci/magician", "gsutil", []string{"-m", "-q", "cp", "/mock/dir/magic-modules/.ci/magician/cassettes/beta/*", "gs://vcr-nightly/beta/2024-07-08/buildID/main_cassettes_backup/fixtures/"}, map[string]string(nil)},
					{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "https://hashicorp:token@github.com/hashicorp/terraform-provider-google-beta", "gopath/src/github.com/hashicorp/terraform-provider-google-beta"}, map[string]string(nil)},
					{"gopath/src/github.com/hashicorp/terraform-provider-google-beta", "go", []string{"list", "./..."}, map[string]string(nil)},
					{"gopath/src/github.com/hashicorp/terraform-provider-google-beta", "go", []string{"test", "", "-parallel", "32", "-v", "-run=TestAcc", "-timeout", "240m", "-ldflags=-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc", "-vet=off"}, map[string]string{
						"ACCTEST_PARALLELISM":            "32",
						"GOOGLE_APPLICATION_CREDENTIALS": "/mock/dir/magic-modules/.ci/magician/sa_key.json",
						"GOOGLE_CREDENTIALS":             "sa_key",
						"GOOGLE_TEST_DIRECTORY":          "",
						"SA_KEY":                         "sa_key",
						"TF_ACC":                         "1",
						"TF_LOG":                         "DEBUG",
						"TF_LOG_PATH_MASK":               "/mock/dir/magic-modules/.ci/magician/testlogs/replaying/beta/%s.log",
						"TF_LOG_SDK_FRAMEWORK":           "INFO",
						"TF_SCHEMA_PANIC_ON_ERROR":       "1",
						"VCR_MODE":                       "REPLAYING",
						"VCR_PATH":                       "/mock/dir/magic-modules/.ci/magician/cassettes/beta",
					}},
					{"/mock/dir/magic-modules/.ci/magician", "gsutil", []string{"-h", "Content-Type:text/plain", "-q", "cp", "-r", "/mock/dir/magic-modules/.ci/magician/testlogs/replaying_test.log", "gs://vcr-nightly/beta/2024-07-08/buildID/logs/replaying/"}, map[string]string(nil)},
					{"/mock/dir/magic-modules/.ci/magician", "gsutil", []string{"-h", "Content-Type:text/plain", "-q", "cp", "-r", "/mock/dir/magic-modules/.ci/magician/testlogs/replaying/beta/*", "gs://vcr-nightly/beta/2024-07-08/buildID/logs/build-log/"}, map[string]string(nil)},
				},
			},
		},
		{
			name: "replay failed then record",
			cmdResults: map[string]string{
				"gopath/src/github.com/hashicorp/terraform-provider-google-beta go [test  -parallel 32 -v -run=TestAcc -timeout 240m -ldflags=-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc -vet=off] map[ACCTEST_PARALLELISM:32 GOOGLE_APPLICATION_CREDENTIALS:/mock/dir/magic-modules/.ci/magician/sa_key.json GOOGLE_CREDENTIALS:sa_key GOOGLE_TEST_DIRECTORY: SA_KEY:sa_key TF_ACC:1 TF_LOG:DEBUG TF_LOG_PATH_MASK:/mock/dir/magic-modules/.ci/magician/testlogs/replaying/beta/%s.log TF_LOG_SDK_FRAMEWORK:INFO TF_SCHEMA_PANIC_ON_ERROR:1 VCR_MODE:REPLAYING VCR_PATH:/mock/dir/magic-modules/.ci/magician/cassettes/beta]":                                           "--- FAIL: TestAccContainerNodePool_defaultDriverInstallation (590.29s)",
				"gopath/src/github.com/hashicorp/terraform-provider-google-beta go [test  -parallel 1 -v -run=TestAccContainerNodePool_defaultDriverInstallation$ -timeout 240m -ldflags=-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc -vet=off] map[ACCTEST_PARALLELISM:1 GOOGLE_APPLICATION_CREDENTIALS:/mock/dir/magic-modules/.ci/magician/sa_key.json GOOGLE_CREDENTIALS:sa_key GOOGLE_TEST_DIRECTORY: SA_KEY:sa_key TF_ACC:1 TF_LOG:DEBUG TF_LOG_PATH_MASK:/mock/dir/magic-modules/.ci/magician/testlogs/recording/beta/%s.log TF_LOG_SDK_FRAMEWORK:INFO TF_SCHEMA_PANIC_ON_ERROR:1 VCR_MODE:RECORDING VCR_PATH:/mock/dir/magic-modules/.ci/magician/cassettes/beta]": "--- PASS: TestAccContainerNodePool_defaultDriverInstallation (590.29s)",
			},
			expectedCalls: map[string][]ParameterList{
				"Run": {
					// replay
					{"/mock/dir/magic-modules/.ci/magician", "gsutil", []string{"-m", "-q", "cp", "gs://ci-vcr-cassettes/beta/fixtures/*", "/mock/dir/magic-modules/.ci/magician/cassettes/beta"}, map[string]string(nil)},
					{"/mock/dir/magic-modules/.ci/magician", "gsutil", []string{"-m", "-q", "cp", "/mock/dir/magic-modules/.ci/magician/cassettes/beta/*", "gs://vcr-nightly/beta/2024-07-08/buildID/main_cassettes_backup/fixtures/"}, map[string]string(nil)},
					{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "https://hashicorp:token@github.com/hashicorp/terraform-provider-google-beta", "gopath/src/github.com/hashicorp/terraform-provider-google-beta"}, map[string]string(nil)},
					{"gopath/src/github.com/hashicorp/terraform-provider-google-beta", "go", []string{"list", "./..."}, map[string]string(nil)},
					{"gopath/src/github.com/hashicorp/terraform-provider-google-beta", "go", []string{"test", "", "-parallel", "32", "-v", "-run=TestAcc", "-timeout", "240m", "-ldflags=-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc", "-vet=off"}, map[string]string{
						"ACCTEST_PARALLELISM":            "32",
						"GOOGLE_APPLICATION_CREDENTIALS": "/mock/dir/magic-modules/.ci/magician/sa_key.json",
						"GOOGLE_CREDENTIALS":             "sa_key",
						"GOOGLE_TEST_DIRECTORY":          "",
						"SA_KEY":                         "sa_key",
						"TF_ACC":                         "1",
						"TF_LOG":                         "DEBUG",
						"TF_LOG_PATH_MASK":               "/mock/dir/magic-modules/.ci/magician/testlogs/replaying/beta/%s.log",
						"TF_LOG_SDK_FRAMEWORK":           "INFO",
						"TF_SCHEMA_PANIC_ON_ERROR":       "1",
						"VCR_MODE":                       "REPLAYING",
						"VCR_PATH":                       "/mock/dir/magic-modules/.ci/magician/cassettes/beta",
					}},
					{"/mock/dir/magic-modules/.ci/magician", "gsutil", []string{"-h", "Content-Type:text/plain", "-q", "cp", "-r", "/mock/dir/magic-modules/.ci/magician/testlogs/replaying_test.log", "gs://vcr-nightly/beta/2024-07-08/buildID/logs/replaying/"}, map[string]string(nil)},
					{"/mock/dir/magic-modules/.ci/magician", "gsutil", []string{"-h", "Content-Type:text/plain", "-q", "cp", "-r", "/mock/dir/magic-modules/.ci/magician/testlogs/replaying/beta/*", "gs://vcr-nightly/beta/2024-07-08/buildID/logs/build-log/"}, map[string]string(nil)},
					// record
					{"gopath/src/github.com/hashicorp/terraform-provider-google-beta", "go", []string{"list", "./..."}, map[string]string(nil)},
					{"gopath/src/github.com/hashicorp/terraform-provider-google-beta", "go", []string{"test", "", "-parallel", "1", "-v", "-run=TestAccContainerNodePool_defaultDriverInstallation$", "-timeout", "240m", "-ldflags=-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc", "-vet=off"}, map[string]string{
						"ACCTEST_PARALLELISM":            "1",
						"GOOGLE_APPLICATION_CREDENTIALS": "/mock/dir/magic-modules/.ci/magician/sa_key.json",
						"GOOGLE_CREDENTIALS":             "sa_key",
						"GOOGLE_TEST_DIRECTORY":          "",
						"SA_KEY":                         "sa_key",
						"TF_ACC":                         "1",
						"TF_LOG":                         "DEBUG",
						"TF_LOG_PATH_MASK":               "/mock/dir/magic-modules/.ci/magician/testlogs/recording/beta/%s.log",
						"TF_LOG_SDK_FRAMEWORK":           "INFO",
						"TF_SCHEMA_PANIC_ON_ERROR":       "1",
						"VCR_MODE":                       "RECORDING",
						"VCR_PATH":                       "/mock/dir/magic-modules/.ci/magician/cassettes/beta",
					}},
					{"/mock/dir/magic-modules/.ci/magician", "gsutil", []string{"-h", "Content-Type:text/plain", "-q", "cp", "-r", "/mock/dir/magic-modules/.ci/magician/testlogs/recording_test.log", "gs://vcr-nightly/beta/2024-07-08/buildID/logs/recording/"}, map[string]string(nil)},
					{"/mock/dir/magic-modules/.ci/magician", "gsutil", []string{"-h", "Content-Type:text/plain", "-q", "cp", "-r", "/mock/dir/magic-modules/.ci/magician/testlogs/recording/beta/*", "gs://vcr-nightly/beta/2024-07-08/buildID/logs/build-log/"}, map[string]string(nil)},
					{"/mock/dir/magic-modules/.ci/magician", "gsutil", []string{"-m", "-q", "cp", "/mock/dir/magic-modules/.ci/magician/cassettes/beta/*", "gs://ci-vcr-cassettes/beta/fixtures/"}, map[string]string(nil)},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rnr := &mockRunner{
				calledMethods: make(map[string][]ParameterList),
				cwd:           "/mock/dir/magic-modules/.ci/magician",
				dirStack:      list.New(),
				cmdResults:    tc.cmdResults,
			}

			ctlr := source.NewController("gopath", "hashicorp", "token", rnr)
			vt, err := vcr.NewTester(map[string]string{
				"SA_KEY": "sa_key",
			}, "ci-vcr-cassettes", "", rnr)
			if err != nil {
				t.Fatalf("Failed to create new tester: %v", err)
			}

			err = execVCRCassetteUpdate("buildID", "2024-07-08", rnr, ctlr, vt)
			if err != nil {
				t.Fatalf("execVCRCassetteUpdate returned error: %v", err)
			}

			for method, expectedCalls := range tc.expectedCalls {
				if actualCalls, ok := rnr.Calls(method); !ok {
					t.Fatalf("Found no calls for %s", method)
				} else if len(actualCalls) != len(expectedCalls) {
					t.Fatalf("Unexpected number of calls for %s, got %d, expected %d", method, len(actualCalls), len(expectedCalls))
				} else {
					for i, actualParams := range actualCalls {
						if expectedParams := expectedCalls[i]; cmp.Diff(expectedParams, actualParams) != "" {
							t.Errorf("Wrong params for call %d to %s, got %v, expected %v, diff = %s", i, method, actualParams, expectedParams, cmp.Diff(expectedParams, actualParams))
						}
					}
				}
			}
		})
	}
}
