/*
* Copyright 2023 Google LLC. All Rights Reserved.
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
	"magician/github"
	"reflect"
	"testing"
)

func TestExecCommunityChecker_CoreContributorFlow(t *testing.T) {
	gh := &mockGithub{
		pullRequest: github.PullRequest{
			User: github.User{
				Login: "core_author",
			},
		},
		userType:      github.CoreContributorUserType,
		calledMethods: make(map[string][][]any),
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string][][]any),
	}

	execCommunityChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	if _, ok := cb.calledMethods["TriggerMMPresubmitRuns"]; ok {
		t.Fatal("Presubmit runs redundantly triggered for core contributor")
	}

	method := "RemoveLabel"
	expected := [][]any{{"pr1", "awaiting-approval"}}
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("awaiting-approval label not removed for PR ")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}

}

func TestExecCommunityChecker_GooglerFlow(t *testing.T) {
	gh := &mockGithub{
		pullRequest: github.PullRequest{
			User: github.User{
				Login: "googler_author",
			},
		},
		userType:           github.GooglerUserType,
		calledMethods:      make(map[string][][]any),
		requestedReviewers: []github.User{github.User{Login: "reviewer1"}},
		previousReviewers:  []github.User{github.User{Login: github.GetRandomReviewer()}, github.User{Login: "reviewer3"}},
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string][][]any),
	}

	execCommunityChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	if _, ok := cb.calledMethods["TriggerMMPresubmitRuns"]; ok {
		t.Fatal("Presubmit runs redundantly triggered for googler")
	}

	method := "RemoveLabel"
	expected := [][]any{{"pr1", "awaiting-approval"}}
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("awaiting-approval label not removed for PR ")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}
}

func TestExecCommunityChecker_AmbiguousUserFlow(t *testing.T) {
	gh := &mockGithub{
		pullRequest: github.PullRequest{
			User: github.User{
				Login: "ambiguous_author",
			},
		},
		userType:           github.CommunityUserType,
		calledMethods:      make(map[string][][]any),
		requestedReviewers: []github.User{github.User{Login: github.GetRandomReviewer()}},
		previousReviewers:  []github.User{github.User{Login: github.GetRandomReviewer()}, github.User{Login: "reviewer3"}},
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string][][]any),
	}

	execCommunityChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	method := "TriggerMMPresubmitRuns"
	expected := [][]any{{"sha1", map[string]string{"BRANCH_NAME": "branch1", "_BASE_BRANCH": "base1", "_HEAD_BRANCH": "head1", "_HEAD_REPO_URL": "url1", "_PR_NUMBER": "pr1"}}}
	if calls, ok := cb.calledMethods[method]; !ok {
		t.Fatal("Presubmit runs not triggered for ambiguous user")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}

	method = "RemoveLabel"
	expected = [][]any{{"pr1", "awaiting-approval"}}
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("awaiting-approval label not removed for PR ")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}
}
