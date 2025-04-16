package cmd

import (
	"magician/github"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExecVCRMerge(t *testing.T) {
	testCases := []struct {
		name            string
		baseBranch      string
		commitSha       string
		lsReturnedError bool
		calledMethods   []string
	}{
		{
			name:       "base branch is main",
			baseBranch: "main",
			commitSha:  "sha",
			calledMethods: []string{
				"gsutil ls gs://ci-vcr-cassettes/refs/heads/auto-pr-123/fixtures/",
				"gsutil -m cp gs://ci-vcr-cassettes/refs/heads/auto-pr-123/fixtures/* gs://ci-vcr-cassettes/fixtures/",
				"gsutil -m rm -r gs://ci-vcr-cassettes/refs/heads/auto-pr-123/",
				"gsutil ls gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-123/fixtures/",
				"gsutil -m cp gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-123/fixtures/* gs://ci-vcr-cassettes/beta/fixtures/",
				"gsutil -m rm -r gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-123/",
			},
		},
		{
			name:       "base branch is not main",
			baseBranch: "test-branch",
			commitSha:  "sha",
			calledMethods: []string{
				"gsutil ls gs://ci-vcr-cassettes/refs/heads/auto-pr-123/fixtures/",
				"gsutil -m cp gs://ci-vcr-cassettes/refs/heads/auto-pr-123/fixtures/* gs://ci-vcr-cassettes/refs/branches/test-branch/fixtures/",
				"gsutil -m rm -r gs://ci-vcr-cassettes/refs/heads/auto-pr-123/",
				"gsutil ls gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-123/fixtures/",
				"gsutil -m cp gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-123/fixtures/* gs://ci-vcr-cassettes/beta/refs/branches/test-branch/fixtures/",
				"gsutil -m rm -r gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-123/",
			},
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
			calledMethods: []string{
				"gsutil ls gs://ci-vcr-cassettes/refs/heads/auto-pr-123/fixtures/",
				"gsutil ls gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-123/fixtures/",
			},
		},
	}

	githubClient := &mockGithub{
		pullRequest: github.PullRequest{
			Number:         123,
			MergeCommitSha: "sha",
		},
		calledMethods: make(map[string][][]any),
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			runner := &mockRunner{
				cwd:           "cwd",
				calledMethods: make(map[string][]ParameterList),
				cmdResults:    make(map[string]string),
			}
			if test.lsReturnedError {
				runner.notifyError = true
			}
			err := execVCRMerge(githubClient, test.commitSha, test.baseBranch, runner)
			if err != nil {
				t.Fatalf("execVCRMerge = %s, want = nil", err)
			}

			got, ok := runner.Calls("Run")
			if !ok && test.calledMethods != nil {
				t.Fatalf("execVCRMerge() expect %d calls, got none", len(test.calledMethods))
			}
			var want []ParameterList
			for _, cmd := range test.calledMethods {
				words := strings.Split(cmd, " ")
				if len(words) > 0 {
					want = append(want, []any{"cwd", words[0], words[1:], map[string]string(nil)})
				}
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Fatalf("execVCRMerge() executed commands diff = %s\n want = %+v, got = %+v", diff, want, got)
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
				"BASE_BRANCH": "main",
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
