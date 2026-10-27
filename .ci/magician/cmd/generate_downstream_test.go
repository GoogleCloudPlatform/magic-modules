/*
* Copyright 2025 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package cmd

import (
	"errors"
	"magician/github"
	"strings"
	"testing"
	"time"
)

// flakyPushRunner fails the first failures calls to `git push`, then succeeds.
type flakyPushRunner struct {
	MockRunner
	failures  int
	pushCalls int
}

func (r *flakyPushRunner) Run(name string, args []string, env map[string]string) (string, error) {
	if name == "git" && len(args) > 0 && args[0] == "push" {
		r.pushCalls++
		if r.pushCalls <= r.failures {
			return "", errors.New("remote: fatal error in commit_refs")
		}
		return "", nil
	}
	return r.MockRunner.Run(name, args, env)
}

func TestPushWithRetry(t *testing.T) {
	oldBackoffs := pushRetryBackoffs
	pushRetryBackoffs = []time.Duration{0, 0, 0}
	defer func() { pushRetryBackoffs = oldBackoffs }()

	cases := map[string]struct {
		failures      int
		wantPushCalls int
		wantErr       bool
	}{
		"succeeds first try": {
			failures:      0,
			wantPushCalls: 1,
		},
		"succeeds after transient failures": {
			failures:      2,
			wantPushCalls: 3,
		},
		"succeeds on last attempt": {
			failures:      3,
			wantPushCalls: 4,
		},
		"fails after exhausting retries": {
			failures:      10,
			wantPushCalls: 4,
			wantErr:       true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			rnr := &flakyPushRunner{MockRunner: NewMockRunner(), failures: tc.failures}
			err := pushWithRetry(rnr, "https://example.com/repo", "auto-pr-123456")
			if (err != nil) != tc.wantErr {
				t.Fatalf("pushWithRetry() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "commit_refs") {
				t.Errorf("pushWithRetry() error = %v, want it to wrap the underlying push error", err)
			}
			if rnr.pushCalls != tc.wantPushCalls {
				t.Errorf("pushWithRetry() made %d push calls, want %d", rnr.pushCalls, tc.wantPushCalls)
			}
		})
	}
}

func TestSyncBranchPushWithRetry(t *testing.T) {
	oldBackoffs := pushRetryBackoffs
	pushRetryBackoffs = []time.Duration{0, 0, 0}
	defer func() { pushRetryBackoffs = oldBackoffs }()

	cases := map[string]struct {
		failures      int
		wantPushCalls int
		wantErr       bool
	}{
		"succeeds first try": {
			failures:      0,
			wantPushCalls: 1,
		},
		"succeeds after transient failures": {
			failures:      2,
			wantPushCalls: 3,
		},
		"fails after exhausting retries": {
			failures:      10,
			wantPushCalls: 4,
			wantErr:       true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			rnr := &flakyPushRunner{MockRunner: NewMockRunner(), failures: tc.failures}
			err := syncBranchPushWithRetry(rnr, "https://example.com/repo", "abc1234:tpg-sync")
			if (err != nil) != tc.wantErr {
				t.Fatalf("syncBranchPushWithRetry() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "commit_refs") {
				t.Errorf("syncBranchPushWithRetry() error = %v, want it to wrap the underlying push error", err)
			}
			if rnr.pushCalls != tc.wantPushCalls {
				t.Errorf("syncBranchPushWithRetry() made %d push calls, want %d", rnr.pushCalls, tc.wantPushCalls)
			}
		})
	}
}

func TestGetPullRequest(t *testing.T) {
	t.Run("resolves PR by commit message PR number even when outside GetPullRequests page", func(t *testing.T) {
		gh := &mockGithub{
			commitMessage: "compute: fix something (#19518)\n\nDetails",
			pullRequests: []github.PullRequest{
				{
					Number:         19518,
					MergeCommitSha: "2f7c8800d8e83831a117ea98d981a67f064ed039",
				},
			},
			// Default pullRequest returned by GetPullRequests has a different SHA.
			pullRequest: github.PullRequest{
				Number:         19999,
				MergeCommitSha: "other-sha",
			},
			calledMethods: make(map[string][][]any),
		}

		pr, err := getPullRequest("main", "2f7c8800d8e83831a117ea98d981a67f064ed039", gh)
		if err != nil {
			t.Fatalf("getPullRequest() unexpected error: %v", err)
		}
		if pr.Number != 19518 {
			t.Errorf("getPullRequest() got PR #%d, want #19518", pr.Number)
		}
		if _, called := gh.calledMethods["GetPullRequests"]; called {
			t.Errorf("expected GetPullRequests not to be called when GetPullRequest succeeds")
		}
	})

	t.Run("falls back to GetPullRequests when commit message has no PR number", func(t *testing.T) {
		gh := &mockGithub{
			commitMessage: "manual commit without pr number",
			pullRequest: github.PullRequest{
				Number:         123,
				MergeCommitSha: "target-sha",
			},
			calledMethods: make(map[string][][]any),
		}

		pr, err := getPullRequest("main", "target-sha", gh)
		if err != nil {
			t.Fatalf("getPullRequest() unexpected error: %v", err)
		}
		if pr.Number != 123 {
			t.Errorf("getPullRequest() got PR #%d, want #123", pr.Number)
		}
		if _, called := gh.calledMethods["GetPullRequests"]; !called {
			t.Errorf("expected GetPullRequests fallback to be called")
		}
	})
}
