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
	"magician/provider"
	"magician/source"
	"reflect"
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

func TestRunMake(t *testing.T) {
	cases := map[string]struct {
		repo     *source.Repo
		command  string
		wantRuns []ParameterList
	}{
		"terraform runs make provider without separate clean-provider": {
			repo: &source.Repo{
				Title:   "terraform",
				Path:    "/mock/dir/tpg",
				Version: provider.GA,
			},
			command: "downstream",
			wantRuns: []ParameterList{
				{"/mock/dir/magic-modules/.ci/magician", "make", []string{"provider", "OUTPUT_PATH=/mock/dir/tpg", "VERSION=ga"}, map[string]string(nil)},
			},
		},
		"terraform-google-conversion downstream updates go.mod without make build": {
			repo: &source.Repo{
				Title:  "terraform-google-conversion",
				Path:   "/mock/dir/tgc",
				Branch: "main",
			},
			command: "downstream",
			wantRuns: []ParameterList{
				{"/mock/dir/magic-modules/.ci/magician", "make", []string{"clean-tgc", "OUTPUT_PATH=/mock/dir/tgc"}, map[string]string(nil)},
				{"/mock/dir/magic-modules/.ci/magician", "make", []string{"tgc", "OUTPUT_PATH=/mock/dir/tgc"}, map[string]string(nil)},
				{"/mock/dir/tgc", "go", []string{"get", "-d", "github.com/hashicorp/terraform-provider-google-beta@main"}, map[string]string(nil)},
				{"/mock/dir/tgc", "go", []string{"mod", "tidy"}, map[string]string(nil)},
			},
		},
		"terraform-google-conversion head only cleans and generates tgc": {
			repo: &source.Repo{
				Title:  "terraform-google-conversion",
				Path:   "/mock/dir/tgc",
				Branch: "main",
			},
			command: "head",
			wantRuns: []ParameterList{
				{"/mock/dir/magic-modules/.ci/magician", "make", []string{"clean-tgc", "OUTPUT_PATH=/mock/dir/tgc"}, map[string]string(nil)},
				{"/mock/dir/magic-modules/.ci/magician", "make", []string{"tgc", "OUTPUT_PATH=/mock/dir/tgc"}, map[string]string(nil)},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			rnr := NewMockRunner()
			if err := runMake(tc.repo, tc.command, rnr); err != nil {
				t.Fatalf("runMake() unexpected error: %v", err)
			}
			gotRuns, _ := rnr.Calls("Run")
			if !reflect.DeepEqual(gotRuns, tc.wantRuns) {
				t.Errorf("runMake() Run calls = %v, want %v", gotRuns, tc.wantRuns)
			}
		})
	}
}

func TestDownstreamBaseBranch(t *testing.T) {
	cases := []struct {
		name       string
		baseBranch string
		repo       string
		want       string
	}{
		{
			name:       "main branch for terraform",
			baseBranch: "main",
			repo:       "terraform",
			want:       "main",
		},
		{
			name:       "main branch for tf-oics maps to master",
			baseBranch: "main",
			repo:       "tf-oics",
			want:       "master",
		},
		{
			name:       "stacked PR branch for terraform normalizes to main",
			baseBranch: "stack-pr-1",
			repo:       "terraform",
			want:       "main",
		},
		{
			name:       "stacked PR branch for tgc normalizes to main",
			baseBranch: "user/feature-stack-1",
			repo:       "terraform-google-conversion",
			want:       "main",
		},
		{
			name:       "stacked PR branch for tf-oics normalizes to master",
			baseBranch: "stack-pr-1",
			repo:       "tf-oics",
			want:       "master",
		},
		{
			name:       "feature release branch is preserved for terraform",
			baseBranch: "FEATURE-BRANCH-major-release-7.0.0",
			repo:       "terraform",
			want:       "FEATURE-BRANCH-major-release-7.0.0",
		},
		{
			name:       "feature release branch is preserved for tf-oics",
			baseBranch: "FEATURE-BRANCH-major-release-7.0.0",
			repo:       "tf-oics",
			want:       "FEATURE-BRANCH-major-release-7.0.0",
		},
		{
			name:       "feature branch with unsafe characters normalizes to main",
			baseBranch: "FEATURE-BRANCH-foo; evil",
			repo:       "terraform",
			want:       "main",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := downstreamBaseBranch(tc.baseBranch, tc.repo); got != tc.want {
				t.Errorf("downstreamBaseBranch(%q, %q) = %q, want %q", tc.baseBranch, tc.repo, got, tc.want)
			}
		})
	}
}

func TestCloneRepo_StackedBaseBranch(t *testing.T) {
	rnr := NewMockRunner()
	ctlr := source.NewController("/go", "modular-magician", "secret-token", rnr)
	mmRepo := &source.Repo{
		Name: "magic-modules",
		Path: "/workspace/mm-terraform-beta-base",
	}

	downstreamRepo, scratchRepo, commitMessage, err := cloneRepo(mmRepo, "stack-pr-1", "terraform", "beta", "base", "12345", rnr, ctlr)
	if err != nil {
		t.Fatalf("cloneRepo() unexpected error: %v", err)
	}
	if downstreamRepo.Branch != "main" {
		t.Errorf("downstreamRepo.Branch = %q, want \"main\"", downstreamRepo.Branch)
	}
	if scratchRepo.Branch != "auto-pr-12345-old" {
		t.Errorf("scratchRepo.Branch = %q, want \"auto-pr-12345-old\"", scratchRepo.Branch)
	}
	if commitMessage != "Old generated code for MM PR 12345." {
		t.Errorf("commitMessage = %q, want \"Old generated code for MM PR 12345.\"", commitMessage)
	}

	calls, ok := rnr.Calls("Run")
	if !ok {
		t.Fatal("expected Run calls on runner")
	}
	wantCalls := []ParameterList{
		{
			"/mock/dir/magic-modules/.ci/magician",
			"git",
			[]string{"clone", "-b", "main", "https://modular-magician:secret-token@github.com/hashicorp/terraform-provider-google-beta", "/go/src/github.com/hashicorp/terraform-provider-google-beta"},
			map[string]string(nil),
		},
		{
			"/workspace/mm-terraform-beta-base",
			"git",
			[]string{"checkout", "HEAD~"},
			map[string]string(nil),
		},
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Errorf("Run calls = %v, want %v", calls, wantCalls)
	}
}

func TestExecGenerateDownstream_RejectsNonReleaseBranchOnDownstreamPush(t *testing.T) {
	rnr := NewMockRunner()
	ctlr := source.NewController("/go", "modular-magician", "secret-token", rnr)
	gh := &mockGithub{}

	err := execGenerateDownstream("stack-pr-1", "downstream", "terraform", "beta", "abc1234", gh, rnr, ctlr)
	if err == nil {
		t.Fatal("execGenerateDownstream() with stacked baseBranch and command=\"downstream\" succeeded, want error")
	}
	if !strings.Contains(err.Error(), "unsupported base branch") {
		t.Errorf("execGenerateDownstream() error = %v, want 'unsupported base branch'", err)
	}
}
