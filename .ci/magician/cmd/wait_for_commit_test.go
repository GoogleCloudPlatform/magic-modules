package cmd

import (
	"strings"
	"testing"
)

func TestExecWaitForCommit(t *testing.T) {
	testCases := []struct {
		name          string
		baseBranch    string
		syncBranch    string
		setupRemote   func(sb *sandbox, shas []string, targetBranch string)
		waitFunc      func(sb *sandbox, shas []string, targetBranch string, calls int)
		expectedCalls int
		expectErr     bool
	}{
		{
			name:       "already in sync branch",
			baseBranch: "main",
			syncBranch: "sync-branch",
			setupRemote: func(sb *sandbox, shas []string, targetBranch string) {
				sb.Runner.MustRun("git", []string{"push", "origin", shas[3] + ":refs/heads/" + targetBranch}, nil)
			},
			waitFunc: func(sb *sandbox, shas []string, targetBranch string, calls int) {
			},
			expectedCalls: 0,
			expectErr:     false,
		},
		{
			name:       "needs update once",
			baseBranch: "main",
			syncBranch: "sync-branch",
			setupRemote: func(sb *sandbox, shas []string, targetBranch string) {
				sb.Runner.MustRun("git", []string{"push", "origin", shas[1] + ":refs/heads/" + targetBranch}, nil)
			},
			waitFunc: func(sb *sandbox, shas []string, targetBranch string, calls int) {
				sb.Runner.MustRun("git", []string{"push", "origin", shas[2] + ":refs/heads/" + targetBranch}, nil)
			},
			expectedCalls: 1,
			expectErr:     false,
		},
		{
			name:       "base branch is not main",
			baseBranch: "feature-branch",
			syncBranch: "sync-branch",
			setupRemote: func(sb *sandbox, shas []string, targetBranch string) {
				sb.Runner.MustRun("git", []string{"push", "origin", shas[1] + ":refs/heads/" + targetBranch}, nil)
			},
			waitFunc: func(sb *sandbox, shas []string, targetBranch string, calls int) {
				sb.Runner.MustRun("git", []string{"push", "origin", shas[2] + ":refs/heads/" + targetBranch}, nil)
			},
			expectedCalls: 1,
			expectErr:     false,
		},
		{
			name:       "needs update multiple times",
			baseBranch: "main",
			syncBranch: "sync-branch",
			setupRemote: func(sb *sandbox, shas []string, targetBranch string) {
				sb.Runner.MustRun("git", []string{"push", "origin", shas[0] + ":refs/heads/" + targetBranch}, nil)
			},
			waitFunc: func(sb *sandbox, shas []string, targetBranch string, calls int) {
				if calls == 0 {
					sb.Runner.MustRun("git", []string{"push", "origin", shas[1] + ":refs/heads/" + targetBranch}, nil)
				} else if calls == 1 {
					sb.Runner.MustRun("git", []string{"push", "origin", shas[2] + ":refs/heads/" + targetBranch}, nil)
				}
			},
			expectedCalls: 2,
			expectErr:     false,
		},
		{
			name:       "error case bad sync branch",
			baseBranch: "main",
			syncBranch: "bad-branch",
			setupRemote: func(sb *sandbox, shas []string, targetBranch string) {
			},
			waitFunc:      func(sb *sandbox, shas []string, targetBranch string, calls int) {},
			expectedCalls: 0,
			expectErr:     true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sb := newSandbox(t)
			sb.RequireAllowlist()
			sb.AllowPassthrough("git")
			originDir := t.TempDir()
			sb.Runner.MustRun("git", []string{"init", "--bare", "-b", "main", originDir}, nil)
			sb.Runner.MustRun("git", []string{"remote", "add", "origin", originDir}, nil)

			var shas []string
			for _, msg := range []string{"Commit A", "Commit B", "Commit C", "Commit D"} {
				sb.Runner.MustRun("git", []string{"commit", "--allow-empty", "-m", msg}, nil)
				shas = append(shas, strings.TrimSpace(sb.Runner.MustRun("git", []string{"rev-parse", "HEAD"}, nil)))
			}

			targetBranch := getSyncBranch(tc.syncBranch, tc.baseBranch)
			tc.setupRemote(sb, shas, targetBranch)
			sb.Runner.MustRun("git", []string{"fetch", "origin"}, nil)

			calls := 0
			testWaitFunc := func() {
				tc.waitFunc(sb, shas, targetBranch, calls)
				calls++
			}

			err := execWaitForCommit(tc.syncBranch, tc.baseBranch, shas[3], sb.Runner, testWaitFunc)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("execWaitForCommit() failed: %v", err)
				}
			}

			if calls != tc.expectedCalls {
				t.Fatalf("expected waitFunc to be called %d times, got %d", tc.expectedCalls, calls)
			}
		})
	}
}
