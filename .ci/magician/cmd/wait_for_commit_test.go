package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type runResult struct {
	out string
	err error
}
type waitForCommitMockRunner struct {
	mockRunner
	runResults map[string][]runResult
}

func (mr *waitForCommitMockRunner) Run(name string, args []string, env map[string]string) (string, error) {
	mr.calledMethods["Run"] = append(mr.calledMethods["Run"], ParameterList{mr.cwd, name, args, env})
	cmd := fmt.Sprintf("%s %s %v %s", mr.cwd, name, args, sortedEnvString(env))
	if result, ok := mr.runResults[cmd]; ok {
		if result == nil {
			return "", fmt.Errorf("no results")
		}
		headRet := result[0]
		result = result[1:]
		mr.runResults[cmd] = result
		return headRet.out, headRet.err
	}
	return "", fmt.Errorf("unknown command %s\n", cmd)
}

func TestExecWaitForCommit(t *testing.T) {
	testCases := []struct {
		name          string
		baseBranch    string
		calledMethods []string
		runResults    map[string][]runResult
	}{
		{
			name:       "base branch is main",
			baseBranch: "main",
			calledMethods: []string{
				"git merge-base --is-ancestor sha origin/sync-branch",
				"git rev-parse --short origin/sync-branch",
				"git rev-parse --short sha~",
				"git fetch origin sync-branch",
				"git rev-parse --short origin/sync-branch",
				"git rev-parse --short sha~",
			},
			runResults: map[string][]runResult{
				"cwd git [merge-base --is-ancestor sha origin/sync-branch] map[]": {
					{
						out: "",
						err: fmt.Errorf("exit error 1"),
					},
				},
				"cwd git [rev-parse --short origin/sync-branch] map[]": {
					{
						out: "sha-x",
					},
					{
						out: "sha-z",
					},
				},
				"cwd git [rev-parse --short sha~] map[]": {
					{
						out: "sha-y",
					},
					{
						out: "sha-z",
					},
				},
				"cwd git [fetch origin sync-branch] map[]": {
					{
						out: "",
					},
				},
			},
		},
		{
			name:       "base branch is not main",
			baseBranch: "feature-branch",
			calledMethods: []string{
				"git merge-base --is-ancestor sha origin/sync-branch-feature-branch",
				"git rev-parse --short origin/sync-branch-feature-branch",
				"git rev-parse --short sha~",
				"git fetch origin sync-branch-feature-branch",
				"git rev-parse --short origin/sync-branch-feature-branch",
				"git rev-parse --short sha~",
			},
			runResults: map[string][]runResult{
				"cwd git [merge-base --is-ancestor sha origin/sync-branch-feature-branch] map[]": {
					{
						out: "",
						err: fmt.Errorf("exit error 1"),
					},
				},
				"cwd git [rev-parse --short origin/sync-branch-feature-branch] map[]": {
					{
						out: "sha-x",
					},
					{
						out: "sha-z",
					},
				},
				"cwd git [rev-parse --short sha~] map[]": {
					{
						out: "sha-y",
					},
					{
						out: "sha-z",
					},
				},
				"cwd git [fetch origin sync-branch-feature-branch] map[]": {
					{
						out: "",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			origWaitFunc := waitFunc
			defer func() {
				waitFunc = origWaitFunc
			}()
			waitFunc = func() {}

			runner := &waitForCommitMockRunner{
				mockRunner: mockRunner{
					cwd:           "cwd",
					calledMethods: make(map[string][]ParameterList),
				},
				runResults: test.runResults,
			}

			err := execWaitForCommit("sync-branch", test.baseBranch, "sha", runner)
			if err != nil {
				t.Fatalf("execWaitForCommit = %s, want = nil", err)
			}

			got, ok := runner.Calls("Run")
			if !ok {
				t.Fatalf("execWaitForCommit() got no calls")
			}
			var want []ParameterList
			for _, cmd := range test.calledMethods {
				words := strings.Split(cmd, " ")
				if len(words) > 0 {
					want = append(want, []any{"cwd", words[0], words[1:], map[string]string(nil)})
				}
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Fatalf("execWaitForCommit() executed commands diff = %s\n want = %+v, got = %+v", diff, want, got)
			}
		})
	}
}
