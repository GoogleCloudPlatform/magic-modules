/*
* Copyright 2024 Google LLC. All Rights Reserved.
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
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"text/template"
	"time"

	membership "magician/github"

	"github.com/google/go-github/v61/github"
	"github.com/spf13/cobra"

	"golang.org/x/exp/slices"

	_ "embed"
)

var (
	// used for flags
	dryRun bool

	//go:embed templates/SCHEDULED_PR_WAITING_FOR_CONTRIBUTOR.md.tmpl
	waitingForContributorTemplate string

	//go:embed templates/SCHEDULED_PR_WAITING_FOR_MERGE.md.tmpl
	waitingForMergeTemplate string

	//go:embed templates/SCHEDULED_PR_WAITING_FOR_REVIEW.md.tmpl
	waitingForReviewTemplate string
)

type reminderCommentData struct {
	User          string
	SinceDays     int
	CoreReviewers []string
}

// scheduledPrReminders sends automated PR notifications and closes stale PRs
var scheduledPrReminders = &cobra.Command{
	Use:   "scheduled-pr-reminders [--dry-run]",
	Short: "Sends automated PR notifications and closes stale PRs",
	Long:  "Sends automated PR notifications and closes stale PRs",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN environment variable")
		}
		gh := github.NewClient(nil).WithAuthToken(githubToken)
		return execScheduledPrReminders(gh)
	},
}

func execScheduledPrReminders(gh *github.Client) error {
	ctx := context.Background()
	opt := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allPulls []*github.PullRequest
	for {
		pulls, resp, err := gh.PullRequests.List(
			ctx,
			"GoogleCloudPlatform",
			"magic-modules",
			opt,
		)
		if err != nil {
			return err
		}
		allPulls = append(allPulls, pulls...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	for index, pr := range allPulls {
		// Skip drafts
		if *pr.Draft {
			fmt.Printf(
				"%d/%d: PR %d: Skipping draft pr\n",
				index+1,
				len(allPulls),
				*pr.Number,
			)
			continue
		}
		var allEvents []*github.IssueEvent
		eventsOpt := &github.ListOptions{PerPage: 100}
		for {
			events, resp, err := gh.Issues.ListIssueEvents(
				ctx,
				"GoogleCloudPlatform",
				"magic-modules",
				*pr.Number,
				eventsOpt,
			)
			if err != nil {
				return err
			}
			allEvents = append(allEvents, events...)
			if resp.NextPage == 0 {
				break
			}
			eventsOpt.Page = resp.NextPage
		}

		var allReviews []*github.PullRequestReview
		reviewsOpt := &github.ListOptions{PerPage: 100}
		for {
			reviews, resp, err := gh.PullRequests.ListReviews(
				ctx,
				"GoogleCloudPlatform",
				"magic-modules",
				*pr.Number,
				reviewsOpt,
			)
			if err != nil {
				return err
			}
			allReviews = append(allReviews, reviews...)
			if resp.NextPage == 0 {
				break
			}
			reviewsOpt.Page = resp.NextPage
		}
		state, since, err := notificationState(pr, allEvents, allReviews)
		if err != nil {
			fmt.Printf(
				"%d/%d: PR %d: error computing notification state: %s\n",
				index+1,
				len(allPulls),
				*pr.Number,
				err,
			)
			continue
		}
		fmt.Printf(
			"%d/%d: PR %d: %s since %v\n",
			index+1,
			len(allPulls),
			*pr.Number,
			state,
			since,
		)
		sinceDays := businessDaysDiff(since, time.Now())
		if shouldNotify(pr, state, sinceDays) {
			comment, err := formatReminderComment(pr, state, sinceDays)
			if err != nil {
				fmt.Printf(
					"%d/%d: PR %d: error rendering comment: %s\n",
					index+1,
					len(allPulls),
					*pr.Number,
					err,
				)
				continue
			}
			if dryRun {
				fmt.Printf("DRY RUN: Would post comment: %s\n", comment)
			} else {
				_, _, err := gh.Issues.CreateComment(
					ctx,
					"GoogleCloudPlatform",
					"magic-modules",
					*pr.Number,
					&github.IssueComment{
						Body: github.String(comment),
					},
				)
				if err != nil {
					return fmt.Errorf("error posting comment to PR %d: %w", *pr.Number, err)
				}
			}
		}

		if shouldClose(pr, state, sinceDays) {
			if dryRun {
				fmt.Printf("DRY RUN: Would close PR %d\n", *pr.Number)
			} else {
				_, _, err := gh.Issues.Edit(
					ctx,
					"GoogleCloudPlatform",
					"magic-modules",
					*pr.Number,
					&github.IssueRequest{
						State: github.String("closed"),
					},
				)
				if err != nil {
					return fmt.Errorf("error closing PR %d: %w", *pr.Number, err)
				}
			}
		}
	}
	return nil
}

type pullRequestReviewState int64

const (
	waitingForReviewerAssignment pullRequestReviewState = iota
	waitingForReview
	waitingForMerge
	waitingForContributor
)

func (s pullRequestReviewState) String() string {
	switch s {
	case waitingForReviewerAssignment:
		return "Waiting for reviewer assignment"
	case waitingForReview:
		return "Waiting for review"
	case waitingForMerge:
		return "Waiting for merge"
	case waitingForContributor:
		return "Waiting for contributor"
	default:
		return fmt.Sprintf("%d", s)
	}
}

// Returns the current state and the time that state was entered. This requires reconciling
// several data sources, since GitHub doesn't return all types of data in all sources.
// The basic algorithm is:
// - find the most recent request for review from a core contributor
//   - if there are none, the state is waitingForReviewerAssignment
//
// - check for any reviews from core reviewers since that review request.
//   - if there are none, the state is waitingForReview and the time is the
//     review request time
//   - if any are change requests, the state is waitingForContributor and the time
//     is the earliest change request
//   - if any are approvals, the state is waitingForMerge and the time is the
//     earliest approval
//   - otherwise there are reviews and all are comment reviews; the state is
//     waitingForContributor and the time is the earliest review time
//
// We don't specially handle cases where the contributor has "acted" because it would be
// significant additional effort, and this case is already handled by re-requesting review
// automatically based on contributor actions.
func notificationState(pr *github.PullRequest, issueEvents []*github.IssueEvent, reviews []*github.PullRequestReview) (pullRequestReviewState, time.Time, error) {
	slices.SortFunc(issueEvents, func(a, b *github.IssueEvent) int {
		if a.CreatedAt.Before(*b.CreatedAt.GetTime()) {
			return 1
		}
		if a.CreatedAt.After(*b.CreatedAt.GetTime()) {
			return -1
		}
		return 0
	})
	slices.SortFunc(reviews, func(a, b *github.PullRequestReview) int {
		if a.SubmittedAt.Before(*b.SubmittedAt.GetTime()) {
			return 1
		}
		if a.SubmittedAt.After(*b.SubmittedAt.GetTime()) {
			return -1
		}
		return 0
	})

	var latestReviewRequest *github.IssueEvent
	removedRequests := map[string]struct{}{}
	for _, event := range issueEvents {
		if *event.Event == "review_request_removed" && event.RequestedReviewer != nil {
			removedRequests[*event.RequestedReviewer.Login] = struct{}{}
			continue
		}
		if *event.Event != "review_requested" {
			continue
		}
		// Ignore review requests for users who no longer exist.
		if event.RequestedReviewer == nil {
			continue
		}
		// Ignore review requests that were later removed.
		if _, ok := removedRequests[*event.RequestedReviewer.Login]; ok {
			continue
		}
		if membership.IsCoreReviewer(*event.RequestedReviewer.Login) {
			latestReviewRequest = event
			break
		}
	}

	if latestReviewRequest == nil {
		return waitingForReviewerAssignment, *pr.CreatedAt.GetTime(), nil
	}

	var earliestApproved *github.PullRequestReview
	var earliestChangesRequested *github.PullRequestReview
	var earliestCommented *github.PullRequestReview

	ignoreBy := map[string]struct{}{}
	for _, review := range reviews {
		if review.SubmittedAt.Before(*latestReviewRequest.CreatedAt.GetTime()) {
			break
		}
		// Ignore reviews by deleted accounts
		if review.User == nil {
			continue
		}
		if !membership.IsCoreReviewer(*review.User.Login) {
			continue
		}
		reviewer := *review.User.Login

		// ignore any reviews by reviewers who had a later approval
		if _, ok := ignoreBy[reviewer]; ok {
			continue
		}
		switch *review.State {
		case "DISMISSED":
			// ignore dismissed reviews
			continue
		case "APPROVED":
			earliestApproved = review
			// ignore all earlier reviews from this reviewer
			ignoreBy[reviewer] = struct{}{}
		case "CHANGES_REQUESTED":
			earliestChangesRequested = review
			// ignore all earlier reviews from this reviewer
			ignoreBy[reviewer] = struct{}{}
		case "COMMENTED":
			earliestCommented = review
		}
	}

	if earliestChangesRequested != nil {
		return waitingForContributor, *earliestChangesRequested.SubmittedAt.GetTime(), nil
	}
	if earliestApproved != nil {
		return waitingForMerge, *earliestApproved.SubmittedAt.GetTime(), nil
	}
	if earliestCommented != nil {
		return waitingForContributor, *earliestCommented.SubmittedAt.GetTime(), nil
	}
	return waitingForReview, *latestReviewRequest.CreatedAt.GetTime(), nil
}

// Calculates the number of PDT days between from and to (by calendar date, not # of hours).
func businessDaysDiff(from, to time.Time) int {
	if to.Before(from) {
		from, to = to, from
	}
	pdtLoc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
	minFrom := time.Date(2024, 5, 1, 0, 0, 0, 0, pdtLoc)
	if from.Before(minFrom) {
		from = minFrom
	}
	from = from.In(pdtLoc)
	to = to.In(pdtLoc)
	// Timezone-aware truncation to day start
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	to = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location())

	// Calculate offset from previous Monday & truncate dates
	fromOffset := int(from.Weekday()) - 1
	if fromOffset == -1 {
		fromOffset = 6
	}
	toOffset := int(to.Weekday()) - 1
	if toOffset == -1 {
		toOffset = 6
	}

	// Calculate # of business days for full weeks
	from = from.AddDate(0, 0, -fromOffset)
	to = to.AddDate(0, 0, -toOffset)
	daysDiff := int(math.Floor(to.Sub(from).Hours()/24/7)) * 5

	// Adjust days based on weekdays from the offsets. For "from", count weekends as
	// 5 days (that is, they "become" the following Monday). For "to", count weekends
	// as 4 days (that is, they "become" the previous Friday).
	fromOffset = int(math.Min(float64(fromOffset), 5))
	toOffset = int(math.Min(float64(toOffset), 4))
	daysDiff += toOffset - fromOffset

	// Special case: daysDiff may be < 0 if from & to are a saturday and sunday
	// from the same weekend. Count this as 0.
	if daysDiff < 0 {
		daysDiff = 0
	}

	return daysDiff
}

func shouldNotify(pr *github.PullRequest, state pullRequestReviewState, sinceDays int) bool {
	labels := map[string]struct{}{}
	for _, label := range pr.Labels {
		labels[*label.Name] = struct{}{}
	}
	switch state {
	case waitingForMerge:
		if _, ok := labels["disable-review-reminders"]; ok {
			return false
		}
		return sinceDays > 0 && sinceDays%5 == 0
	case waitingForContributor:
		if _, ok := labels["disable-automatic-closure"]; ok {
			return false
		}
		return slices.Contains([]int{10, 20, 28, 30}, sinceDays)
	case waitingForReview:
		if _, ok := labels["disable-review-reminders"]; ok {
			return false
		}
		return sinceDays == 3 || (sinceDays > 0 && sinceDays%5 == 0)
	}
	return false
}

func formatReminderComment(pullRequest *github.PullRequest, state pullRequestReviewState, sinceDays int) (string, error) {
	embeddedTemplate := ""
	switch state {
	case waitingForMerge:
		embeddedTemplate = waitingForMergeTemplate
	case waitingForContributor:
		embeddedTemplate = waitingForContributorTemplate
	case waitingForReview:
		embeddedTemplate = waitingForReviewTemplate
	default:
		return "", fmt.Errorf("state does not have corresponding template: %s", state.String())
	}
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"weekdaysToWeeks": func(a int) int {
			return a / 5
		},
	}).Parse(embeddedTemplate)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse template for %s: %s", state.String(), err))
	}

	coreReviewers := []string{}
	for _, reviewer := range pullRequest.RequestedReviewers {
		if membership.IsCoreReviewer(*reviewer.Login) {
			coreReviewers = append(coreReviewers, *reviewer.Login)
		}
	}

	data := reminderCommentData{
		User:          *pullRequest.User.Login,
		SinceDays:     sinceDays,
		CoreReviewers: coreReviewers,
	}

	sb := new(strings.Builder)
	err = tmpl.Execute(sb, data)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func shouldClose(pr *github.PullRequest, state pullRequestReviewState, sinceDays int) bool {
	for _, label := range pr.Labels {
		if *label.Name == "disable-automatic-closure" {
			return false
		}
	}
	return state == waitingForContributor && sinceDays >= 30
}

func init() {
	rootCmd.AddCommand(scheduledPrReminders)
	scheduledPrReminders.Flags().BoolVar(&dryRun, "dry-run", false, "Only log write actions instead of updating PRs")
}
