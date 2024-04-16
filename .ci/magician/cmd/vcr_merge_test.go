package cmd

import (
	"magician/github"
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
			execVCRMerge(githubClient, test.commitSha, test.baseBranch, runner)
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
