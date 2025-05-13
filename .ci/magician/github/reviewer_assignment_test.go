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
package github

import (
	"fmt"
	"strings"
	"testing"
	"time"

	gh "github.com/google/go-github/v68/github"
	"golang.org/x/exp/slices"
)

func TestChooseCoreReviewers(t *testing.T) {
	if len(AvailableReviewers(nil)) < 2 {
		t.Fatalf("not enough available reviewers (%v) to test (need at least 2)", AvailableReviewers(nil))
	}
	firstCoreReviewer := AvailableReviewers(nil)[0]
	secondCoreReviewer := AvailableReviewers(nil)[1]

	cases := map[string]struct {
		RequestedReviewers                               []*gh.User
		PreviousReviewers                                []*gh.User
		ExpectReviewersFromList, ExpectSpecificReviewers []string
		ExpectPrimaryReviewer                            bool
	}{
		"no previous review requests assigns new reviewer from team": {
			RequestedReviewers:      []*gh.User{},
			PreviousReviewers:       []*gh.User{},
			ExpectReviewersFromList: AvailableReviewers(nil),
			ExpectPrimaryReviewer:   true,
		},
		"requested reviewer from team means that primary reviewer was already selected": {
			RequestedReviewers:    []*gh.User{ghUser(firstCoreReviewer)},
			PreviousReviewers:     []*gh.User{},
			ExpectPrimaryReviewer: false,
		},
		"requested off-team reviewer does not mean that primary reviewer was already selected": {
			RequestedReviewers:    []*gh.User{ghUser("foobar")},
			PreviousReviewers:     []*gh.User{},
			ExpectPrimaryReviewer: true,
		},
		"previously involved team member reviewers should have review requested and mean that primary reviewer was already selected": {
			RequestedReviewers:      []*gh.User{},
			PreviousReviewers:       []*gh.User{ghUser(firstCoreReviewer)},
			ExpectSpecificReviewers: []string{firstCoreReviewer},
			ExpectPrimaryReviewer:   false,
		},
		"previously involved reviewers that are not team members are ignored": {
			RequestedReviewers:      []*gh.User{},
			PreviousReviewers:       []*gh.User{ghUser("foobar")},
			ExpectReviewersFromList: AvailableReviewers(nil),
			ExpectPrimaryReviewer:   true,
		},
		"only previously involved team member reviewers will have review requested": {
			RequestedReviewers:      []*gh.User{},
			PreviousReviewers:       []*gh.User{ghUser(firstCoreReviewer), ghUser("foobar"), ghUser(secondCoreReviewer)},
			ExpectSpecificReviewers: []string{firstCoreReviewer, secondCoreReviewer},
			ExpectPrimaryReviewer:   false,
		},
		"primary reviewer will not have review requested even if other team members previously reviewed": {
			RequestedReviewers:      []*gh.User{ghUser(secondCoreReviewer)},
			PreviousReviewers:       []*gh.User{ghUser(firstCoreReviewer)},
			ExpectSpecificReviewers: []string{firstCoreReviewer},
			ExpectPrimaryReviewer:   false,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			reviewers, primaryReviewer := ChooseCoreReviewers(tc.RequestedReviewers, tc.PreviousReviewers)
			if tc.ExpectPrimaryReviewer && primaryReviewer == "" {
				t.Error("wanted primary reviewer to be returned; got none")
			}
			if !tc.ExpectPrimaryReviewer && primaryReviewer != "" {
				t.Errorf("wanted no primary reviewer; got %s", primaryReviewer)
			}
			if len(tc.ExpectReviewersFromList) > 0 {
				for _, reviewer := range reviewers {
					if !slices.Contains(tc.ExpectReviewersFromList, reviewer) {
						t.Errorf("wanted reviewer %s to be in list %v but they were not", reviewer, tc.ExpectReviewersFromList)
					}
				}
			}
			if len(tc.ExpectSpecificReviewers) > 0 {
				if !slices.Equal(reviewers, tc.ExpectSpecificReviewers) {
					t.Errorf("wanted reviewers to be %v; instead got %v", tc.ExpectSpecificReviewers, reviewers)
				}
			}
		})
	}
}

// Helper function to create a github User struct
func ghUser(login string) *gh.User {
	return &gh.User{Login: gh.Ptr(login)}
}

func TestFormatReviewerComment(t *testing.T) {
	cases := map[string]struct {
		Reviewer       string
		AuthorUserType UserType
		Trusted        bool
	}{
		"community contributor": {
			Reviewer:       "foobar",
			AuthorUserType: CommunityUserType,
			Trusted:        false,
		},
		"googler": {
			Reviewer:       "foobar",
			AuthorUserType: GooglerUserType,
			Trusted:        true,
		},
		"core contributor": {
			Reviewer:       "foobar",
			AuthorUserType: CoreContributorUserType,
			Trusted:        true,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			comment := FormatReviewerComment(tc.Reviewer)
			t.Log(comment)
			if !strings.Contains(comment, fmt.Sprintf("@%s", tc.Reviewer)) {
				t.Errorf("wanted comment to contain @%s; does not.", tc.Reviewer)
			}
			if !strings.Contains(comment, "Tests will require approval") {
				t.Errorf("wanted comment to say tests will require approval; does not")
			}
		})
	}
}

func TestFindReviewerComment(t *testing.T) {
	now := time.Now()

	cases := map[string]struct {
		Comments        []*gh.IssueComment
		ExpectReviewer  string
		ExpectCommentID int64
	}{
		"no reviewer comment": {
			Comments: []*gh.IssueComment{
				{
					Body: gh.Ptr("this is not a reviewer comment"),
				},
			},
			ExpectReviewer:  "",
			ExpectCommentID: 0,
		},
		"reviewer comment": {
			Comments: []*gh.IssueComment{
				{
					Body: gh.Ptr(FormatReviewerComment("trodge")),
					ID:   gh.Int64(1234),
				},
			},
			ExpectReviewer:  "trodge",
			ExpectCommentID: 1234,
		},
		"multiple reviewer comments": {
			Comments: []*gh.IssueComment{
				{
					Body:      gh.Ptr(FormatReviewerComment("trodge")),
					ID:        gh.Int64(1234),
					CreatedAt: &gh.Timestamp{Time: now.Add(-48 * time.Hour)}, // 2 days ago
				},
				{
					Body:      gh.Ptr(FormatReviewerComment("c2thorn")),
					ID:        gh.Int64(5678),
					CreatedAt: &gh.Timestamp{Time: now.Add(-24 * time.Hour)}, // 1 day ago
				},
				{
					Body:      gh.Ptr(FormatReviewerComment("melinath")),
					ID:        gh.Int64(91011),
					CreatedAt: &gh.Timestamp{Time: now.Add(-36 * time.Hour)}, // 1.5 days ago
				},
			},
			ExpectReviewer:  "c2thorn",
			ExpectCommentID: 5678,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			comment, reviewer := FindReviewerComment(tc.Comments)
			if reviewer != tc.ExpectReviewer {
				t.Errorf("wanted reviewer to be %s; got %s", tc.ExpectReviewer, reviewer)
			}
			if (comment == nil && tc.ExpectCommentID != 0) ||
				(comment != nil && *comment.ID != tc.ExpectCommentID) {
				var actualID int64
				if comment != nil {
					actualID = *comment.ID
				}
				t.Errorf("wanted comment ID to be %d; got %d", tc.ExpectCommentID, actualID)
			}
		})
	}
}
