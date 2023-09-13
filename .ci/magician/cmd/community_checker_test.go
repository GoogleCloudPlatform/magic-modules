package cmd

import (
	"magician/github"
	"testing"
)

func TestExecCommunityChecker_CoreContributorFlow(t *testing.T) {
	gh := &mockGithub{
		author:        "core_author",
		userType:      github.CoreContributorUserType,
		calledMethods: make(map[string]bool),
	}
	cb := &mockCloudBuild{
		calledMethods: make(map[string]bool),
	}

	execCommunityChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	if cb.calledMethods["TriggerMMPresubmitRuns"] {
		t.Fatal("presubmit runs redundantly triggered for core contributor")
	}

	if !gh.calledMethods["RemoveLabel"] {
		t.Fatal("awaiting-approval label not removed for PR ")
	}

}

func TestExecCommunityChecker_GooglerFlow(t *testing.T) {
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

	execCommunityChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	if cb.calledMethods["TriggerMMPresubmitRuns"] {
		t.Fatal("presubmit runs redundantly triggered for googler")
	}

	if !gh.calledMethods["RemoveLabel"] {
		t.Fatal("awaiting-approval label not removed for PR ")
	}
}

func TestExecCommunityChecker_AmbiguousUserFlow(t *testing.T) {
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

	execCommunityChecker("pr1", "sha1", "branch1", "url1", "head1", "base1", gh, cb)

	if !cb.calledMethods["TriggerMMPresubmitRuns"] {
		t.Fatal("presubmit runs not triggered for ambiguous user")
	}

	if !gh.calledMethods["RemoveLabel"] {
		t.Fatal("awaiting-approval label not removed for PR ")
	}
}
