package cmd

import (
	"magician/github"
	"reflect"
	"regexp"
	"testing"
)

func TestExecMembershipChecker_CoreContributorFlow(t *testing.T) {
	gh := &mockGithub{
		author:        "core_author",
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
		author:            "googler_author",
		userType:          github.GooglerUserType,
		calledMethods:     make(map[string][][]any),
		firstReviewer:     "reviewer1",
		previousReviewers: []string{github.GetRandomReviewer(), "reviewer3"},
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string][][]any),
	}

	execMembershipChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	method := "RequestPullRequestReviewer"
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("Review wasn't requested for googler")
	} else if len(calls) != 1 {
		t.Fatalf("Wrong number of calls for %s, got %d, expected 1", method, len(calls))
	} else if params := calls[0]; len(params) != 2 {
		t.Fatalf("Wrong number of params for %s, got %d, expected 2", method, len(params))
	} else if param := params[0]; param != "pr1" {
		t.Fatalf("Wrong first param for %s, got %v, expected pr1", method, param)
	} else if param := params[1]; !github.IsTeamReviewer(param.(string)) {
		t.Fatalf("Wrong second param for %s, got %v, expected a team reviewer", method, param)
	}

	method = "TriggerMMPresubmitRuns"
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
		author:            "ambiguous_author",
		userType:          github.CommunityUserType,
		calledMethods:     make(map[string][][]any),
		firstReviewer:     github.GetRandomReviewer(),
		previousReviewers: []string{github.GetRandomReviewer(), "reviewer3"},
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string][][]any),
	}

	execMembershipChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	method := "RequestPullRequestReviewer"
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("Review wasn't requested for ambiguous user")
	} else if len(calls) != 1 {
		t.Fatalf("Wrong number of calls for %s, got %d, expected 1", method, len(calls))
	} else if params := calls[0]; len(params) != 2 {
		t.Fatalf("Wrong number of params for %s, got %d, expected 2", method, len(params))
	} else if param := params[0]; param != "pr1" {
		t.Fatalf("Wrong first param for %s, got %v, expected pr1", method, param)
	} else if param := params[1]; !github.IsTeamReviewer(param.(string)) {
		t.Fatalf("Wrong second param for %s, got %v, expected a team reviewer", method, param)
	}

	method = "AddLabel"
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
		author:            "googler_author",
		userType:          github.GooglerUserType,
		calledMethods:     make(map[string][][]any),
		firstReviewer:     "",
		previousReviewers: []string{"reviewer3"},
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string][][]any),
	}

	execMembershipChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	method := "RequestPullRequestReviewer"
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("Review wasn't requested for googler")
	} else if len(calls) != 1 {
		t.Fatalf("Wrong number of calls for %s, got %d, expected 1", method, len(calls))
	} else if params := calls[0]; len(params) != 2 {
		t.Fatalf("Wrong number of params for %s, got %d, expected 2", method, len(params))
	} else if param := params[0]; param != "pr1" {
		t.Fatalf("Wrong first param for %s, got %v, expected pr1", method, param)
	} else if param := params[1]; !github.IsTeamReviewer(param.(string)) {
		t.Fatalf("Wrong second param for %s, got %v, expected a team reviewer", method, param)
	}

	method = "PostComment"
	reviewerExp := regexp.MustCompile(`@(.*?),`)
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("Comment wasn't posted stating user status")
	} else if len(calls) != 1 {
		t.Fatalf("Wrong number of calls for %s, got %d, expected 1", method, len(calls))
	} else if params := calls[0]; len(params) != 2 {
		t.Fatalf("Wrong number of params for %s, got %d, expected 2", method, len(params))
	} else if param := params[0]; param != "pr1" {
		t.Fatalf("Wrong first param for %s, got %v, expected pr1", method, param)
	} else if param, ok := params[1].(string); !ok {
		t.Fatalf("Got non-string second param for %s", method)
	} else if submatches := reviewerExp.FindStringSubmatch(param); len(submatches) != 2 || !github.IsTeamReviewer(submatches[1]) {
		t.Fatalf("%s called without a team reviewer (found %v) in the comment: %s", method, submatches, param)
	}
}
