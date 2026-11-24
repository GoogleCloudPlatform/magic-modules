package cmd

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestFetchPRNumber(t *testing.T) {
	sb := newSandbox(t)
	sb.RequireAllowlist()
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
		commitMessage: "Add `additional_group_keys` attribute to `google_cloud_identity_group` resource (#9217)\n\n* Add `additional_group_keys` attribute to `google_cloud_identity_group` resource\n\n* Update acceptance test to check for attribute\n\n* Fix test check\n\n* Add `output: true` to nested properties in output field\n\nSigned-off-by: Modular Magician <magic-modules@google.com>",
	}

	prNumber, err := fetchPRNumber("49d3741f9d4d810a0a4768363bb8498afa21c688", "main", sb.Runner, gh)
	if err != nil {
		t.Errorf("error fetching PR number: %s", err)
	}

	if prNumber != "9217" {
		t.Errorf("PR number is %s, expected 9217", prNumber)
	}

	// Verify GitHub API was called
	if calls, ok := gh.calledMethods["GetCommitMessage"]; !ok || len(calls) == 0 {
		t.Errorf("Expected GetCommitMessage to be called")
	} else {
		args := calls[0]
		if len(args) != 3 {
			t.Errorf("Expected GetCommitMessage to be called with 3 arguments, got %d", len(args))
		} else {
			if args[0] != "GoogleCloudPlatform" {
				t.Errorf("Expected owner to be 'GoogleCloudPlatform', got '%s'", args[0])
			}
			if args[1] != "magic-modules" {
				t.Errorf("Expected repo to be 'magic-modules', got '%s'", args[1])
			}
			if args[2] != "49d3741f9d4d810a0a4768363bb8498afa21c688" {
				t.Errorf("Expected SHA to be '49d3741f9d4d810a0a4768363bb8498afa21c688', got '%s'", args[2])
			}
		}
	}
}

func TestDeleteBranchWithRetry(t *testing.T) {
	oldBackoffs := pushRetryBackoffs
	pushRetryBackoffs = []time.Duration{0, 0, 0}
	defer func() { pushRetryBackoffs = oldBackoffs }()

	t.Run("succeeds first try", func(t *testing.T) {
		rnr := &flakyPushRunner{MockRunner: NewMockRunner(), failures: 0}
		if err := deleteBranchWithRetry(rnr, "https://example.com/repo", ":auto-pr-19408"); err != nil {
			t.Fatalf("deleteBranchWithRetry() unexpected error: %v", err)
		}
		if rnr.pushCalls != 1 {
			t.Errorf("got %d push calls, want 1", rnr.pushCalls)
		}
	})

	t.Run("succeeds after transient failures", func(t *testing.T) {
		rnr := &flakyPushRunner{MockRunner: NewMockRunner(), failures: 2}
		if err := deleteBranchWithRetry(rnr, "https://example.com/repo", ":auto-pr-19408"); err != nil {
			t.Fatalf("deleteBranchWithRetry() unexpected error: %v", err)
		}
		if rnr.pushCalls != 3 {
			t.Errorf("got %d push calls, want 3", rnr.pushCalls)
		}
	})

	t.Run("ignores remote ref does not exist without retrying", func(t *testing.T) {
		rnr := &missingBranchPushRunner{MockRunner: NewMockRunner()}
		if err := deleteBranchWithRetry(rnr, "https://example.com/repo", ":auto-pr-19408"); err != nil {
			t.Fatalf("deleteBranchWithRetry() unexpected error: %v", err)
		}
		if rnr.pushCalls != 1 {
			t.Errorf("got %d push calls, want 1", rnr.pushCalls)
		}
	})

	t.Run("fails after exhausting retries", func(t *testing.T) {
		rnr := &flakyPushRunner{MockRunner: NewMockRunner(), failures: 10}
		err := deleteBranchWithRetry(rnr, "https://example.com/repo", ":auto-pr-19408")
		if err == nil || !strings.Contains(err.Error(), "commit_refs") {
			t.Fatalf("expected commit_refs error, got: %v", err)
		}
		if rnr.pushCalls != 4 {
			t.Errorf("got %d push calls, want 4", rnr.pushCalls)
		}
	})
}

type missingBranchPushRunner struct {
	MockRunner
	pushCalls int
}

func (r *missingBranchPushRunner) Run(name string, args []string, env map[string]string) (string, error) {
	if name == "git" && len(args) > 0 && args[0] == "push" {
		r.pushCalls++
		return "", errors.New("error: unable to delete 'auto-pr-19408': remote ref does not exist")
	}
	return r.MockRunner.Run(name, args, env)
}
