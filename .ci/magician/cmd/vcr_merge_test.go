package cmd

import (
	"magician/github"
	"os"
	"strings"
	"testing"
)

func TestExecVCRMerge(t *testing.T) {
	testCases := []struct {
		name            string
		baseBranch      string
		commitSha       string
		lsReturnedError bool
	}{
		{
			name:       "base branch is main",
			baseBranch: "main",
			commitSha:  "sha",
		},
		{
			name:       "base branch is not main",
			baseBranch: "test-branch",
			commitSha:  "sha",
		},
		{
			name:       "pr not found",
			baseBranch: "main",
			commitSha:  "random-sha",
		},
		{
			name:            "ls returns error",
			commitSha:       "sha",
			lsReturnedError: true,
			baseBranch:      "main",
		},
	}

	githubClient := &mockGithub{
		pullRequest: github.PullRequest{
			Number:         123,
			MergeCommitSha: "sha",
		},
		calledMethods: make(map[string][][]any),
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sb := newSandbox(t)

			// Intercepts the hardcoded `gcloud storage` commands and translates gs:// URLs into local directory operations.
			fakeGcloud := `#!/bin/bash
				if [ "$1" = "storage" ]; then
					if [ "$2" = "ls" ]; then
						TARGET=$(echo $3 | sed 's|gs://|gs/|')
						ls "$TARGET" > /dev/null 2>&1
						exit $?
					elif [ "$2" = "cp" ]; then
						SRC=$(echo $3 | sed 's|gs://|gs/|' | sed 's|/\*$||')
						DEST=$(echo $4 | sed 's|gs://|gs/|')
						mkdir -p "$DEST"
						cp -r "$SRC"/* "$DEST"
					elif [ "$2" = "rm" ]; then
						TARGET=$(echo $4 | sed 's|gs://|gs/|')
						rm -r "$TARGET"
					fi
				fi`
			sb.Runner.WriteFile("gcloud", fakeGcloud)
			sb.Runner.MustRun("chmod", []string{"+x", "gcloud"}, nil)

			if !tc.lsReturnedError && tc.name != "pr not found" {
				sb.Runner.MustRun("mkdir", []string{"-p", "gs/ci-vcr-cassettes/refs/heads/auto-pr-123/fixtures"}, nil)
				sb.Runner.MustRun("mkdir", []string{"-p", "gs/ci-vcr-cassettes/beta/refs/heads/auto-pr-123/fixtures"}, nil)
				sb.Runner.WriteFile("gs/ci-vcr-cassettes/refs/heads/auto-pr-123/fixtures/dummy1.txt", "data")
				sb.Runner.WriteFile("gs/ci-vcr-cassettes/beta/refs/heads/auto-pr-123/fixtures/dummy2.txt", "data")
			}

			err := execVCRMerge(githubClient, tc.commitSha, tc.baseBranch, sb.Runner)
			if err != nil {
				t.Fatalf("execVCRMerge() failed: %v", err)
			}

			if !tc.lsReturnedError && tc.name != "pr not found" {
				destBranchPath := ""
				if tc.baseBranch != "main" {
					destBranchPath = "/refs/branches/" + tc.baseBranch
				}

				if _, err := os.Stat(sb.Dir + "/gs/ci-vcr-cassettes" + destBranchPath + "/fixtures/dummy1.txt"); os.IsNotExist(err) {
					t.Fatalf("Expected file to be copied to /gs/ci-vcr-cassettes%s/fixtures/dummy1.txt", destBranchPath)
				}
				if _, err := os.Stat(sb.Dir + "/gs/ci-vcr-cassettes/beta" + destBranchPath + "/fixtures/dummy2.txt"); os.IsNotExist(err) {
					t.Fatalf("Expected file to be copied to /gs/ci-vcr-cassettes/beta%s/fixtures/dummy2.txt", destBranchPath)
				}

				if _, err := os.Stat(sb.Dir + "/gs/ci-vcr-cassettes/refs/heads/auto-pr-123/"); !os.IsNotExist(err) {
					t.Fatalf("Expected source directory /gs/ci-vcr-cassettes/refs/heads/auto-pr-123/ to be deleted")
				}
				if _, err := os.Stat(sb.Dir + "/gs/ci-vcr-cassettes/beta/refs/heads/auto-pr-123/"); !os.IsNotExist(err) {
					t.Fatalf("Expected source directory /gs/ci-vcr-cassettes/beta/refs/heads/auto-pr-123/ to be deleted")
				}
			}
		})
	}
}

func TestVCRMergeRunE(t *testing.T) {
	testCases := []struct {
		name    string
		envVars map[string]string
		errMsg  string
	}{
		{
			name: "GITHUB_TOKEN_CLASSIC is missing in env var",
			envVars: map[string]string{
				"BASE_BRANCH":          "main",
				"GITHUB_TOKEN_CLASSIC": "",
			},
		},
		{
			name: "BASE_BRANCH env var is not set",
			envVars: map[string]string{
				"GITHUB_TOKEN_CLASSIC": "123",
			},
		},
		{
			name: "BASE_BRANCH env var is empty",
			envVars: map[string]string{
				"GITHUB_TOKEN_CLASSIC": "123",
				"BASE_BRANCH":          "",
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			origEnvVars := make(map[string]string)
			for k, v := range test.envVars {
				if origVal, ok := os.LookupEnv(k); ok {
					origEnvVars[k] = origVal
				}
				os.Setenv(k, v)
			}
			defer func() {
				for k := range test.envVars {
					if origVal, ok := origEnvVars[k]; ok {
						os.Setenv(k, origVal)
					} else {
						os.Unsetenv(k)
					}
				}
			}()
			err := vcrMergeCmd.RunE(nil, []string{"sha"})
			if err == nil || !strings.Contains(err.Error(), test.errMsg) {
				t.Fatalf("vcrMergeCmd got %s, want error message with %q", err, test.errMsg)
			}
		})
	}
}
