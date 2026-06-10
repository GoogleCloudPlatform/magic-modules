package cmd

import (
	"strings"
	"testing"
)

func TestExecWaitForCommit(t *testing.T) {
	testCases := []struct {
		name       string
		baseBranch string
		syncBranch string
	}{
		{
			name:       "base branch is main",
			baseBranch: "main",
			syncBranch: "sync-branch",
		},
		{
			name:       "base branch is not main",
			baseBranch: "feature-branch",
			syncBranch: "sync-branch-feature-branch",
		},
		{
			name:       "already in sync branch",
			baseBranch: "main",
			syncBranch: "sync-branch",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sb := newSandbox(t)

			originDir := t.TempDir()
			sb.Runner.MustRun("git", []string{"init", "--bare", "-b", "main", originDir}, nil)
			sb.Runner.MustRun("git", []string{"remote", "add", "origin", originDir}, nil)

			sb.Runner.MustRun("git", []string{"commit", "--allow-empty", "-m", "Commit A"}, nil)
			shaA := strings.TrimSpace(sb.Runner.MustRun("git", []string{"rev-parse", "HEAD"}, nil))

			sb.Runner.MustRun("git", []string{"commit", "--allow-empty", "-m", "Commit B"}, nil)
			shaB := strings.TrimSpace(sb.Runner.MustRun("git", []string{"rev-parse", "HEAD"}, nil))

			sb.Runner.MustRun("git", []string{"commit", "--allow-empty", "-m", "Commit C"}, nil)
			shaC := strings.TrimSpace(sb.Runner.MustRun("git", []string{"rev-parse", "HEAD"}, nil))

			if tc.name == "already in sync branch" {
				sb.Runner.MustRun("git", []string{"push", "origin", shaC + ":refs/heads/" + tc.syncBranch}, nil)
			} else {
				sb.Runner.MustRun("git", []string{"push", "origin", shaA + ":refs/heads/" + tc.syncBranch}, nil)
			}

			sb.Runner.MustRun("git", []string{"fetch", "origin"}, nil)

			origWaitFunc := waitFunc
			defer func() { waitFunc = origWaitFunc }()

			waitFunc = func() {
				if tc.name == "already in sync branch" {
					t.Fatalf("waitFunc was called, but test should have exited early!")
				}
				sb.Runner.MustRun("git", []string{"push", "origin", shaB + ":refs/heads/" + tc.syncBranch}, nil)
			}

			err := execWaitForCommit("sync-branch", tc.baseBranch, shaC, sb.Runner)
			if err != nil {
				t.Fatalf("execWaitForCommit() failed: %v", err)
			}
		})
	}
}
