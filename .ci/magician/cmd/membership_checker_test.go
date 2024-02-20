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

func TestExecMembershipChecker_CoreContributorFlow(t *testing.T) {
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

	execMembershipChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	if _, ok := gh.calledMethods["RequestPullRequestReviewer"]; ok {
		t.Fatal("Incorrectly requested review for core contributor")
	}

	method := "TriggerMMPresubmitRuns"
	expected := [][]any{{"sha1", map[string]string{"BRANCH_NAME": "branch1", "_BASE_BRANCH": "base1", "_HEAD_BRANCH": "head1", "_HEAD_REPO_URL": "url1", "_PR_NUMBER": "pr1"}}}
	if calls, ok := cb.calledMethods[method]; !ok {
		t.Fatal("Presubmit runs not triggered for core author")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}

	method = "ApproveCommunityChecker"
	expected = [][]any{{"pr1", "sha1"}}
	if calls, ok := cb.calledMethods[method]; !ok {
		t.Fatal("Community checker not approved for core author")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}

}

func TestExecMembershipChecker_GooglerFlow(t *testing.T) {
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

	execMembershipChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	method := "TriggerMMPresubmitRuns"
	expected := [][]any{{"sha1", map[string]string{"BRANCH_NAME": "branch1", "_BASE_BRANCH": "base1", "_HEAD_BRANCH": "head1", "_HEAD_REPO_URL": "url1", "_PR_NUMBER": "pr1"}}}
	if calls, ok := cb.calledMethods[method]; !ok {
		t.Fatal("Presubmit runs not triggered for googler")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}

	method = "ApproveCommunityChecker"
	expected = [][]any{{"pr1", "sha1"}}
	if calls, ok := cb.calledMethods[method]; !ok {
		t.Fatal("Community checker not approved for googler")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}
}

func TestExecMembershipChecker_AmbiguousUserFlow(t *testing.T) {
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

	execMembershipChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	method := "AddLabel"
	expected := [][]any{{"pr1", "awaiting-approval"}}
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("Label wasn't posted to pull request")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}

	method = "GetAwaitingApprovalBuildLink"
	expected = [][]any{{"pr1", "sha1"}}
	if calls, ok := cb.calledMethods[method]; !ok {
		t.Fatal("Awaiting approval build link wasn't gotten from pull request")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}

	if _, ok := gh.calledMethods["ApproveCommunityChecker"]; ok {
		t.Fatal("Incorrectly approved community checker for ambiguous user")
	}

	if _, ok := gh.calledMethods["TriggerMMPresubmitRuns"]; ok {
		t.Fatal("Incorrectly triggered presubmit runs for ambiguous user")
	}
}

func TestExecMembershipChecker_CommentForNewPrimaryReviewer(t *testing.T) {
	gh := &mockGithub{
		pullRequest: github.PullRequest{
			User: github.User{
				Login: "googler_author",
			},
		},
		userType:           github.GooglerUserType,
		calledMethods:      make(map[string][][]any),
		requestedReviewers: []github.User{},
		previousReviewers:  []github.User{github.User{Login: "reviewer3"}},
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string][][]any),
	}

	execMembershipChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)
}
