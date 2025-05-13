/*
* Copyright 2023 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
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

	ghi "github.com/google/go-github/v68/github"
)

func TestExecMembershipChecker_CoreContributorFlow(t *testing.T) {
	mockGH := &mockGithub{
		pullRequest: &ghi.PullRequest{
			User: &ghi.User{
				Login: ghi.Ptr("core_author"),
			},
		},
		userType:      github.CoreContributorUserType,
		calledMethods: make(map[string][][]any),
	}

	cb := &mockCloudBuild{
		calledMethods: make(map[string][][]any),
	}

	execMembershipChecker("pr1", "sha1", mockGH, cb)

	method := "ApproveDownstreamGenAndTest"
	expected := [][]any{{"pr1", "sha1"}}

	if calls, ok := cb.calledMethods[method]; !ok {
		t.Fatal("Community checker not approved for core author")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}
}

func TestExecMembershipChecker_GooglerFlow(t *testing.T) {
	randomReviewer := github.GetRandomReviewer(nil)

	mockGH := &mockGithub{
		pullRequest: &ghi.PullRequest{
			User: &ghi.User{
				Login: ghi.Ptr("googler_author"),
			},
		},
		userType:      github.GooglerUserType,
		calledMethods: make(map[string][][]any),
		requestedReviewers: []*ghi.User{
			{Login: ghi.Ptr("reviewer1")},
		},
		previousReviewers: []*ghi.User{
			{Login: ghi.Ptr(randomReviewer)},
			{Login: ghi.Ptr("reviewer3")},
		},
	}

	cb := &mockCloudBuild{
		calledMethods: make(map[string][][]any),
	}

	execMembershipChecker("pr1", "sha1", mockGH, cb)

	method := "ApproveDownstreamGenAndTest"
	expected := [][]any{{"pr1", "sha1"}}

	if calls, ok := cb.calledMethods[method]; !ok {
		t.Fatal("Community checker not approved for googler")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}
}

func TestExecMembershipChecker_AmbiguousUserFlow(t *testing.T) {
	randomReviewer := github.GetRandomReviewer(nil)

	mockGH := &mockGithub{
		pullRequest: &ghi.PullRequest{
			User: &ghi.User{
				Login: ghi.Ptr("ambiguous_author"),
			},
		},
		userType:      github.CommunityUserType,
		calledMethods: make(map[string][][]any),
		requestedReviewers: []*ghi.User{
			{Login: ghi.Ptr(randomReviewer)},
		},
		previousReviewers: []*ghi.User{
			{Login: ghi.Ptr(randomReviewer)},
			{Login: ghi.Ptr("reviewer3")},
		},
	}

	cb := &mockCloudBuild{
		calledMethods: make(map[string][][]any),
	}

	execMembershipChecker("pr1", "sha1", mockGH, cb)

	method := "AddLabels"
	expected := [][]any{{"pr1", []string{"awaiting-approval"}}}

	if calls, ok := mockGH.calledMethods[method]; !ok {
		t.Fatal("Label wasn't posted to pull request")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}

	if _, ok := mockGH.calledMethods["ApproveDownstreamGenAndTest"]; ok {
		t.Fatal("Incorrectly approved community checker for ambiguous user")
	}
}

func TestExecMembershipChecker_CommentForNewPrimaryReviewer(t *testing.T) {
	mockGH := &mockGithub{
		pullRequest: &ghi.PullRequest{
			User: &ghi.User{
				Login: ghi.Ptr("googler_author"),
			},
		},
		userType:           github.GooglerUserType,
		calledMethods:      make(map[string][][]any),
		requestedReviewers: []*ghi.User{},
		previousReviewers: []*ghi.User{
			{Login: ghi.Ptr("reviewer3")},
		},
	}

	cb := &mockCloudBuild{
		calledMethods: make(map[string][][]any),
	}

	execMembershipChecker("pr1", "sha1", mockGH, cb)
}
