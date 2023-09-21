package cmd

import (
	"magician/github"
	"testing"
)

func TestExecMembershipChecker_CoreContributorFlow(t *testing.T) {
	gh := &mockGithub{
		author:        "core_author",
		userType:      github.CoreContributorUserType,
		calledMethods: make(map[string]bool),
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string]bool),
	}

	execMembershipChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	if gh.calledMethods["RequestPullRequestReviewer"] {
		t.Fatal("Incorrectly requested review for core contributor")
	}

	if !cb.calledMethods["TriggerMMPresubmitRuns"] {
		t.Fatal("presubmit runs not triggered for core author")
	}

	if !cb.calledMethods["ApproveCommunityChecker"] {
		t.Fatal("community checker not approved for core author")
	}

}

func TestExecMembershipChecker_GooglerFlow(t *testing.T) {
	gh := &mockGithub{
		author:            "googler_author",
		userType:          github.GooglerUserType,
		calledMethods:     make(map[string]bool),
		firstReviewer:     "reviewer1",
		previousReviewers: []string{github.GetRandomReviewer(), "reviewer3"},
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string]bool),
	}

	execMembershipChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	if !gh.calledMethods["RequestPullRequestReviewer"] {
		t.Fatal("Review wasn't requested for googler")
	}

	if !cb.calledMethods["TriggerMMPresubmitRuns"] {
		t.Fatal("Presubmit runs not triggered for googler")
	}

	if !cb.calledMethods["ApproveCommunityChecker"] {
		t.Fatal("Community checker not approved for googler")
	}
}

func TestExecMembershipChecker_AmbiguousUserFlow(t *testing.T) {
	gh := &mockGithub{
		author:            "ambiguous_author",
		userType:          github.CommunityUserType,
		calledMethods:     make(map[string]bool),
		firstReviewer:     github.GetRandomReviewer(),
		previousReviewers: []string{github.GetRandomReviewer(), "reviewer3"},
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string]bool),
	}

	execMembershipChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	if !gh.calledMethods["RequestPullRequestReviewer"] {
		t.Fatal("Review wasn't requested for ambiguous user")
	}

	if !gh.calledMethods["AddLabel"] || !cb.calledMethods["GetAwaitingApprovalBuildLink"] {
		t.Fatal("Label wasn't posted to pull request")
	}

	if cb.calledMethods["ApproveCommunityChecker"] {
		t.Fatal("Incorrectly approved community checker for ambiguous user")
	}

	if cb.calledMethods["TriggerMMPresubmitRuns"] {
		t.Fatal("Incorrectly triggered presubmit runs for ambiguous user")
	}
}

func TestExecMembershipChecker_CommentForNewPrimaryReviewer(t *testing.T) {
	gh := &mockGithub{
		author:            "googler_author",
		userType:          github.GooglerUserType,
		calledMethods:     make(map[string]bool),
		firstReviewer:     "",
		previousReviewers: []string{"reviewer3"},
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string]bool),
	}

	execMembershipChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	if !gh.calledMethods["PostComment"] {
		t.Fatal("Review wasn't requested for googler")
	}

	if !gh.calledMethods["PostComment"] {
		t.Fatal("Comment wasn't posted stating user status")
	}
}
