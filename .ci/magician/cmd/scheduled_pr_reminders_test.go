package cmd

import (
	"testing"
	"time"

	membership "magician/github"

	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
)

func TestNotificationState(t *testing.T) {
	firstCoreReviewer := membership.AvailableReviewers()[0]
	secondCoreReviewer := membership.AvailableReviewers()[1]
	cases := map[string]struct {
		pullRequest *github.PullRequest
		issueEvents []*github.IssueEvent
		reviews     []*github.PullRequestReview
		expectState pullRequestReviewState
		expectSince time.Time
	}{
		// expectState: waitingForReviewerAssignment
		"no review requests, and no reviews": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			expectState: waitingForReviewerAssignment,
			expectSince: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		"request for non-core reviewer, and no reviews": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String("reviewer")},
				},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String("reviewer")},
				},
			},
			expectState: waitingForReviewerAssignment,
			expectSince: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		"request for team reviewer, and no reviews": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				RequestedTeams: []*github.Team{
					&github.Team{Name: github.String("terraform-team")},
				},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:         github.String("review_requested"),
					CreatedAt:     &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedTeam: &github.Team{Name: github.String("terraform-team")},
				},
			},
			expectState: waitingForReviewerAssignment,
			expectSince: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		"request for team reviewer which was later removed, and no reviews": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				RequestedTeams: []*github.Team{
					&github.Team{Name: github.String("terraform-team")},
				},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_request_removed"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			expectState: waitingForReviewerAssignment,
			expectSince: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},

		// expectState: waitingForReview
		"no reviews": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			expectState: waitingForReview,
			expectSince: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		"review requested, removed, and rerequested, with no reviews": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_request_removed"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			expectState: waitingForReview,
			expectSince: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
		},
		"no reviews since latest review request": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(secondCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			reviews: []*github.PullRequestReview{
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(firstCoreReviewer)},
					State:       github.String("APPROVED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(firstCoreReviewer)},
					State:       github.String("CHANGES_REQUESTED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(firstCoreReviewer)},
					State:       github.String("COMMENTED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(firstCoreReviewer)},
					State:       github.String("DISMISSED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectState: waitingForReview,
			expectSince: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		"ignore reviews from deleted accounts": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			reviews: []*github.PullRequestReview{
				&github.PullRequestReview{
					User:        nil,
					State:       github.String("APPROVED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectState: waitingForReview,
			expectSince: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},

		// waitingForContributor
		"change request followed by comment review from same reviewer": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(secondCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			reviews: []*github.PullRequestReview{
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("APPROVED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("CHANGES_REQUESTED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("COMMENTED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 3, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectState: waitingForContributor,
			expectSince: time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
		},
		"approved review with a change request review": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(secondCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			reviews: []*github.PullRequestReview{
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(firstCoreReviewer)},
					State:       github.String("APPROVED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("CHANGES_REQUESTED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectState: waitingForContributor,
			expectSince: time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
		},
		"approved review followed by change request review from same user": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(secondCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			reviews: []*github.PullRequestReview{
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("APPROVED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("CHANGES_REQUESTED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectState: waitingForContributor,
			expectSince: time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
		},

		// expectState: waitingForMerge
		"approved review on its own": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(secondCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			reviews: []*github.PullRequestReview{
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("APPROVED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 3, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("COMMENTED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("CHANGES_REQUESTED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectState: waitingForMerge,
			expectSince: time.Date(2024, 2, 3, 0, 0, 0, 0, time.UTC),
		},
		"approved review with a comment review": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(secondCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			reviews: []*github.PullRequestReview{
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(firstCoreReviewer)},
					State:       github.String("APPROVED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("COMMENTED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectState: waitingForMerge,
			expectSince: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		"approved review followed by comment review from same user": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(secondCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			reviews: []*github.PullRequestReview{
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("APPROVED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("COMMENTED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectState: waitingForMerge,
			expectSince: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		"approved review followed by dismissed review from same user": {
			pullRequest: &github.PullRequest{
				User:      &github.User{Login: github.String("author")},
				CreatedAt: &github.Timestamp{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			issueEvents: []*github.IssueEvent{
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(secondCoreReviewer)},
				},
				&github.IssueEvent{
					Event:             github.String("review_requested"),
					CreatedAt:         &github.Timestamp{time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
					RequestedReviewer: &github.User{Login: github.String(firstCoreReviewer)},
				},
			},
			reviews: []*github.PullRequestReview{
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(firstCoreReviewer)},
					State:       github.String("APPROVED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(firstCoreReviewer)},
					State:       github.String("DISMISSED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)},
				},
				&github.PullRequestReview{
					User:        &github.User{Login: github.String(secondCoreReviewer)},
					State:       github.String("COMMENTED"),
					SubmittedAt: &github.Timestamp{time.Date(2024, 2, 3, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectState: waitingForMerge,
			expectSince: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			state, since, err := notificationState(
				tc.pullRequest,
				tc.issueEvents,
				tc.reviews,
			)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectState.String(), state.String())
			assert.Equal(t, tc.expectSince, since)
		})
	}
}

func TestBusinessDaysDiff(t *testing.T) {
	pdtLoc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
	cases := map[string]struct {
		from, to   time.Time
		expectDays int
	}{
		"same time": {
			from:       time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			expectDays: 0,
		},
		"same day": {
			from:       time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 6, 1, 0, 0, 0, pdtLoc),
			expectDays: 0,
		},
		"next day, earlier": {
			from:       time.Date(2024, 5, 6, 5, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 7, 0, 0, 0, 0, pdtLoc),
			expectDays: 1,
		},
		"next day, later": {
			from:       time.Date(2024, 5, 6, 5, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 7, 6, 0, 0, 0, pdtLoc),
			expectDays: 1,
		},
		"next week, earlier": {
			from:       time.Date(2024, 5, 6, 5, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 13, 3, 0, 0, 0, pdtLoc),
			expectDays: 5,
		},
		"next month": {
			from:       time.Date(2024, 5, 6, 23, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 6, 6, 1, 0, 0, 0, pdtLoc),
			expectDays: 23,
		},
		"previous day, earlier": {
			from:       time.Date(2024, 5, 7, 5, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			expectDays: 1,
		},
		"previous day, later": {
			from:       time.Date(2024, 5, 7, 5, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 6, 10, 0, 0, 0, pdtLoc),
			expectDays: 1,
		},
		"earlier than minFrom": {
			from:       time.Date(2022, 1, 1, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 3, 11, 1, 0, 0, pdtLoc),
			expectDays: 2,
		},
		"UTC times": {
			from:       time.Date(2024, 5, 7, 0, 0, 0, 0, time.UTC),
			to:         time.Date(2024, 5, 7, 23, 0, 0, 0, time.UTC),
			expectDays: 1,
		},
		"saturday to sunday": {
			from:       time.Date(2024, 5, 4, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 5, 0, 0, 0, 0, pdtLoc),
			expectDays: 0,
		},
		"saturday to monday": {
			from:       time.Date(2024, 5, 4, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			expectDays: 0,
		},
		"sunday to monday": {
			from:       time.Date(2024, 5, 5, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			expectDays: 0,
		},
		"friday to monday": {
			from:       time.Date(2024, 5, 3, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			expectDays: 1,
		},
		"monday to friday": {
			from:       time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 10, 0, 0, 0, 0, pdtLoc),
			expectDays: 4,
		},
		"monday to saturday": {
			from:       time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 11, 0, 0, 0, 0, pdtLoc),
			expectDays: 4,
		},
		"monday to sunday": {
			from:       time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 12, 0, 0, 0, 0, pdtLoc),
			expectDays: 4,
		},
		"monday to friday plus 1 week": {
			from:       time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 17, 0, 0, 0, 0, pdtLoc),
			expectDays: 9,
		},
		"monday to saturday plus 1 week": {
			from:       time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 18, 0, 0, 0, 0, pdtLoc),
			expectDays: 9,
		},
		"monday to sunday plus 1 week": {
			from:       time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 19, 0, 0, 0, 0, pdtLoc),
			expectDays: 9,
		},
		"monday plus 1 week": {
			from:       time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 13, 0, 0, 0, 0, pdtLoc),
			expectDays: 5,
		},
		"monday plus 2 weeks": {
			from:       time.Date(2024, 5, 6, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 20, 0, 0, 0, 0, pdtLoc),
			expectDays: 10,
		},
		"sunday plus 1 week": {
			from:       time.Date(2024, 5, 5, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 12, 0, 0, 0, 0, pdtLoc),
			expectDays: 4,
		},
		"sunday plus 2 weeks": {
			from:       time.Date(2024, 5, 5, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 19, 0, 0, 0, 0, pdtLoc),
			expectDays: 9,
		},
		"saturday plus 1 weeks": {
			from:       time.Date(2024, 5, 4, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 11, 0, 0, 0, 0, pdtLoc),
			expectDays: 4,
		},
		"saturday plus 2 weeks": {
			from:       time.Date(2024, 5, 4, 0, 0, 0, 0, pdtLoc),
			to:         time.Date(2024, 5, 18, 0, 0, 0, 0, pdtLoc),
			expectDays: 9,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			days := businessDaysDiff(tc.from, tc.to)
			assert.Equal(t, tc.expectDays, days)
		})
	}
}

func TestShouldNotify(t *testing.T) {
	cases := map[string]struct {
		pullRequest *github.PullRequest
		state       pullRequestReviewState
		sinceDays   int
		want        bool
	}{
		// waitingForMerge
		"waitingForMerge first day": {
			pullRequest: &github.PullRequest{},
			state:       waitingForMerge,
			sinceDays:   0,
			want:        false,
		},
		"waitingForMerge too early": {
			pullRequest: &github.PullRequest{},
			state:       waitingForMerge,
			sinceDays:   4,
			want:        false,
		},
		"waitingForMerge first week": {
			pullRequest: &github.PullRequest{},
			state:       waitingForMerge,
			sinceDays:   5,
			want:        true,
		},
		"waitingForMerge not on a week": {
			pullRequest: &github.PullRequest{},
			state:       waitingForMerge,
			sinceDays:   6,
			want:        false,
		},
		"waitingForMerge after many weeks": {
			pullRequest: &github.PullRequest{},
			state:       waitingForMerge,
			sinceDays:   5 * 57,
			want:        true,
		},
		"waitingForMerge skip with label": {
			pullRequest: &github.PullRequest{
				Labels: []*github.Label{{Name: github.String("disable-review-reminders")}},
			},
			state:     waitingForMerge,
			sinceDays: 5,
			want:      false,
		},
		"waitingForMerge ignore disable-automatic-closure": {
			pullRequest: &github.PullRequest{
				Labels: []*github.Label{{Name: github.String("disable-automatic-closure")}},
			},
			state:     waitingForMerge,
			sinceDays: 10,
			want:      true,
		},

		// waitingForReview
		"waitingForReview first day": {
			pullRequest: &github.PullRequest{},
			state:       waitingForReview,
			sinceDays:   0,
			want:        false,
		},
		"waitingForReview too early": {
			pullRequest: &github.PullRequest{},
			state:       waitingForReview,
			sinceDays:   1,
			want:        false,
		},
		"waitingForReview three days": {
			pullRequest: &github.PullRequest{},
			state:       waitingForReview,
			sinceDays:   3,
			want:        true,
		},
		"waitingForReview first week": {
			pullRequest: &github.PullRequest{},
			state:       waitingForReview,
			sinceDays:   5,
			want:        true,
		},
		"waitingForReview not on a week": {
			pullRequest: &github.PullRequest{},
			state:       waitingForReview,
			sinceDays:   6,
			want:        false,
		},
		"waitingForReview after many weeks": {
			pullRequest: &github.PullRequest{},
			state:       waitingForReview,
			sinceDays:   5 * 57,
			want:        true,
		},
		"waitingForReview skip with label": {
			pullRequest: &github.PullRequest{
				Labels: []*github.Label{{Name: github.String("disable-review-reminders")}},
			},
			state:     waitingForReview,
			sinceDays: 5,
			want:      false,
		},
		"waitingForReview ignore disable-automatic-closure": {
			pullRequest: &github.PullRequest{
				Labels: []*github.Label{{Name: github.String("disable-automatic-closure")}},
			},
			state:     waitingForReview,
			sinceDays: 10,
			want:      true,
		},

		// waitingForContributor
		"waitingForContributor first day": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   0,
			want:        false,
		},
		"waitingForContributor too early": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   4,
			want:        false,
		},
		"waitingForContributor two weeks": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   10,
			want:        true,
		},
		"waitingForContributor four weeks": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   20,
			want:        true,
		},
		"waitingForContributor 28 weekdays": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   28,
			want:        true,
		},
		"waitingForContributor six weeks": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   30,
			want:        true,
		},
		"waitingForContributor other sinceDays": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   9,
			want:        false,
		},
		"waitingForContributor skip with label": {
			pullRequest: &github.PullRequest{
				Labels: []*github.Label{{Name: github.String("disable-automatic-closure")}},
			},
			state:     waitingForContributor,
			sinceDays: 10,
			want:      false,
		},
		"waitingForContributor ignore disable-review-reminders": {
			pullRequest: &github.PullRequest{
				Labels: []*github.Label{{Name: github.String("disable-review-reminders")}},
			},
			state:     waitingForContributor,
			sinceDays: 10,
			want:      true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			got := shouldNotify(tc.pullRequest, tc.state, tc.sinceDays)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFormatReminderComment(t *testing.T) {
	firstCoreReviewer := membership.AvailableReviewers()[0]
	secondCoreReviewer := membership.AvailableReviewers()[1]
	cases := map[string]struct {
		pullRequest        *github.PullRequest
		state              pullRequestReviewState
		sinceDays          int
		expectedStrings    []string
		notExpectedStrings []string
	}{
		// waitingForMerge
		"waitingForMerge one week": {
			pullRequest: &github.PullRequest{
				User: &github.User{Login: github.String("pr-author")},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String(firstCoreReviewer)},
					&github.User{Login: github.String(secondCoreReviewer)},
					&github.User{Login: github.String("other-reviewer")},
				},
			},
			state: waitingForMerge,
			sinceDays: 5,
			expectedStrings: []string{
				"waiting for merge for 1 week",
				"disable-review-reminders",
				"@" + firstCoreReviewer,
				"@" + secondCoreReviewer,
			},
			notExpectedStrings: []string{
				"@pr-author",
				"@other-reviewer",
			},
		},
		"waitingForMerge two weeks": {
			pullRequest: &github.PullRequest{
				User: &github.User{Login: github.String("pr-author")},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String(firstCoreReviewer)},
					&github.User{Login: github.String(secondCoreReviewer)},
					&github.User{Login: github.String("other-reviewer")},
				},
			},
			state:     waitingForMerge,
			sinceDays: 5 * 2,
			expectedStrings: []string{
				"waiting for merge for 2 weeks",
				"disable-review-reminders",
				"@" + firstCoreReviewer,
				"@" + secondCoreReviewer,
			},
			notExpectedStrings: []string{
				"@pr-author",
				"@other-reviewer",
			},
		},
		"waitingForMerge many weeks": {
			pullRequest: &github.PullRequest{
				User: &github.User{Login: github.String("pr-author")},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String(firstCoreReviewer)},
					&github.User{Login: github.String(secondCoreReviewer)},
					&github.User{Login: github.String("other-reviewer")},
				},
			},
			state:     waitingForMerge,
			sinceDays: 5 * 57,
			expectedStrings: []string{
				"waiting for merge for 57 weeks",
				"disable-review-reminders",
				"@" + firstCoreReviewer,
				"@" + secondCoreReviewer,
			},
			notExpectedStrings: []string{
				"@pr-author",
				"@other-reviewer",
			},
		},

		// waitingForReview
		"waitingForReview three days": {
			pullRequest: &github.PullRequest{
				User: &github.User{Login: github.String("pr-author")},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String(firstCoreReviewer)},
					&github.User{Login: github.String(secondCoreReviewer)},
					&github.User{Login: github.String("other-reviewer")},
				},
			},
			state:     waitingForReview,
			sinceDays: 3,
			expectedStrings: []string{
				"waiting for review for 3 weekdays",
				"disable-review-reminders",
				"@" + firstCoreReviewer,
				"@" + secondCoreReviewer,
			},
			notExpectedStrings: []string{
				"@GoogleCloudPlatform/terraform-team",
				"@pr-author",
				"@other-reviewer",
			},
		},
		"waitingForReview one week": {
			pullRequest: &github.PullRequest{
				User: &github.User{Login: github.String("pr-author")},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String(firstCoreReviewer)},
					&github.User{Login: github.String(secondCoreReviewer)},
					&github.User{Login: github.String("other-reviewer")},
				},
			},
			state:     waitingForReview,
			sinceDays: 5,
			expectedStrings: []string{
				"@GoogleCloudPlatform/terraform-team",
				"waiting for review for 1 week",
				"disable-review-reminders",
				"@" + firstCoreReviewer,
				"@" + secondCoreReviewer,
			},
			notExpectedStrings: []string{
				"@pr-author",
				"@other-reviewer",
			},
		},
		"waitingForReview two weeks": {
			pullRequest: &github.PullRequest{
				User: &github.User{Login: github.String("pr-author")},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String(firstCoreReviewer)},
					&github.User{Login: github.String(secondCoreReviewer)},
					&github.User{Login: github.String("other-reviewer")},
				},
			},
			state:     waitingForReview,
			sinceDays: 10,
			expectedStrings: []string{
				"@GoogleCloudPlatform/terraform-team",
				"waiting for review for 2 weeks",
				"disable-review-reminders",
				"@" + firstCoreReviewer,
				"@" + secondCoreReviewer,
			},
			notExpectedStrings: []string{
				"@pr-author",
				"@other-reviewer",
			},
		},

		// waitingForContributor
		"waitingForContributor two weeks": {
			pullRequest: &github.PullRequest{
				User: &github.User{Login: github.String("pr-author")},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String(firstCoreReviewer)},
					&github.User{Login: github.String(secondCoreReviewer)},
					&github.User{Login: github.String("other-reviewer")},
				},
			},
			state:     waitingForContributor,
			sinceDays: 10,
			expectedStrings: []string{
				"@pr-author",
				"If no action is taken, this PR will be closed in 28 days",
				"disable-automatic-closure",
			},
			notExpectedStrings: []string{
				"@" + firstCoreReviewer,
				"@" + secondCoreReviewer,
				"@other-reviewer",
			},
		},
		"waitingForContributor four weeks": {
			pullRequest: &github.PullRequest{
				User: &github.User{Login: github.String("pr-author")},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String(firstCoreReviewer)},
					&github.User{Login: github.String(secondCoreReviewer)},
					&github.User{Login: github.String("other-reviewer")},
				},
			},
			state:     waitingForContributor,
			sinceDays: 20,
			expectedStrings: []string{
				"@pr-author",
				"If no action is taken, this PR will be closed in 14 days",
				"disable-automatic-closure",
			},
			notExpectedStrings: []string{
				"@" + firstCoreReviewer,
				"@" + secondCoreReviewer,
				"@other-reviewer",
			},
		},
		"waitingForContributor 28 days": {
			pullRequest: &github.PullRequest{
				User: &github.User{Login: github.String("pr-author")},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String(firstCoreReviewer)},
					&github.User{Login: github.String(secondCoreReviewer)},
					&github.User{Login: github.String("other-reviewer")},
				},
			},
			state:     waitingForContributor,
			sinceDays: 28,
			expectedStrings: []string{
				"@pr-author",
				"If no action is taken, this PR will be closed in 2 weekdays",
				"disable-automatic-closure",
			},
			notExpectedStrings: []string{
				"@" + firstCoreReviewer,
				"@" + secondCoreReviewer,
				"@other-reviewer",
			},
		},
		"waitingForContributor six weeks": {
			pullRequest: &github.PullRequest{
				User: &github.User{Login: github.String("pr-author")},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String(firstCoreReviewer)},
					&github.User{Login: github.String(secondCoreReviewer)},
					&github.User{Login: github.String("other-reviewer")},
				},
			},
			state:           waitingForContributor,
			sinceDays:       30,
			expectedStrings: []string{"@pr-author", "PR is being closed due to inactivity"},
			notExpectedStrings: []string{
				"If no action is taken, this PR will be closed",
				"disable-automatic-closure",
				"@" + firstCoreReviewer,
				"@" + secondCoreReviewer,
				"@other-reviewer",
			},
		},
		"waitingForContributor seven weeks": {
			pullRequest: &github.PullRequest{
				User: &github.User{Login: github.String("pr-author")},
				RequestedReviewers: []*github.User{
					&github.User{Login: github.String(firstCoreReviewer)},
					&github.User{Login: github.String(secondCoreReviewer)},
					&github.User{Login: github.String("other-reviewer")},
				},
			},
			state:           waitingForContributor,
			sinceDays:       35,
			expectedStrings: []string{"@pr-author", "PR is being closed due to inactivity"},
			notExpectedStrings: []string{
				"If no action is taken, this PR will be closed",
				"disable-automatic-closure",
				"@" + firstCoreReviewer,
				"@" + secondCoreReviewer,
				"@other-reviewer",
			},
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			comment, err := formatReminderComment(tc.pullRequest, tc.state, tc.sinceDays)
			assert.Nil(t, err)

			for _, s := range tc.expectedStrings {
				assert.Contains(t, comment, s)
			}

			for _, s := range tc.notExpectedStrings {
				assert.NotContains(t, comment, s)
			}
		})
	}
}

func TestShouldClose(t *testing.T) {
	cases := map[string]struct {
		pullRequest *github.PullRequest
		state       pullRequestReviewState
		sinceDays   int
		want        bool
	}{
		// waitingForContributor
		"waitingForContributor first day": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   0,
			want:        false,
		},
		"waitingForContributor two weeks": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   10,
			want:        false,
		},
		"waitingForContributor four weeks": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   20,
			want:        false,
		},
		"waitingForContributor six weeks": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   30,
			want:        true,
		},
		"waitingForContributor seven weeks": {
			pullRequest: &github.PullRequest{},
			state:       waitingForContributor,
			sinceDays:   35,
			want:        true,
		},
		"waitingForMerge six weeks": {
			pullRequest: &github.PullRequest{},
			state:       waitingForMerge,
			sinceDays:   30,
			want:        false,
		},
		"waitingForReview six weeks": {
			pullRequest: &github.PullRequest{},
			state:       waitingForReview,
			sinceDays:   30,
			want:        false,
		},
		"waitingForContributor skip with label": {
			pullRequest: &github.PullRequest{
				Labels: []*github.Label{{Name: github.String("disable-automatic-closure")}},
			},
			state:     waitingForContributor,
			sinceDays: 30,
			want:      false,
		},
		"waitingForContributor ignore disable-review-reminders": {
			pullRequest: &github.PullRequest{
				Labels: []*github.Label{{Name: github.String("disable-review-reminders")}},
			},
			state:     waitingForContributor,
			sinceDays: 30,
			want:      true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			got := shouldClose(tc.pullRequest, tc.state, tc.sinceDays)
			assert.Equal(t, tc.want, got)
		})
	}
}
