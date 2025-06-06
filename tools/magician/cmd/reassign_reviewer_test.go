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
	"magician/github"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecReassignReviewer(t *testing.T) {
	availableReviewers := github.AvailableReviewers(nil)
	if len(availableReviewers) < 3 {
		t.Fatalf("not enough available reviewers (%v) to run TestExecRequestReviewer (need at least 3)", availableReviewers)
	}
	cases := map[string]struct {
		newPrimaryReviewer      string
		comments                []github.PullRequestComment
		expectSpecificReviewers []string
		expectRemovedReviewers  []string
	}{
		"reassign from no current reviewer to random reviewer": {
			comments: []github.PullRequestComment{},
		},
		"reassign from no current reviewer to alice": {
			newPrimaryReviewer:      "alice",
			comments:                []github.PullRequestComment{},
			expectSpecificReviewers: []string{"alice"},
		},
		"reassign from bob to random reviewer": {
			comments: []github.PullRequestComment{
				{
					Body: github.FormatReviewerComment("bob"),
					ID:   1234,
				},
			},
			expectRemovedReviewers: []string{"bob"},
		},
		"reassign from bob to alice": {
			newPrimaryReviewer: "alice",
			comments: []github.PullRequestComment{
				{
					Body: github.FormatReviewerComment("bob"),
					ID:   1234,
				},
			},
			expectSpecificReviewers: []string{"alice"},
			expectRemovedReviewers:  []string{"bob"},
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			gh := &mockGithub{
				pullRequest: github.PullRequest{
					User: github.User{Login: "author"},
				},
				calledMethods:       make(map[string][][]any),
				pullRequestComments: tc.comments,
			}

			err := execReassignReviewer("1", tc.newPrimaryReviewer, gh)
			if err != nil {
				t.Fatalf("execReassignReviewer failed: %v", err)
			}
			if len(gh.calledMethods["RequestPullRequestReviewers"]) != 1 {
				t.Errorf("Expected RequestPullRequestReviewers called 1 time, got %v", len(gh.calledMethods["RequestPullRequestReviewers"]))
			}

			var assignedReviewers []string
			for _, args := range gh.calledMethods["RequestPullRequestReviewers"] {
				assignedReviewers = append(assignedReviewers, args[1].([]string)...)
			}
			var removedReviewers []string
			for _, args := range gh.calledMethods["RemovePullRequestReviewers"] {
				removedReviewers = append(removedReviewers, args[1].([]string)...)
			}

			if tc.expectSpecificReviewers != nil {
				for _, reviewer := range assignedReviewers {
					assert.Contains(t, tc.expectSpecificReviewers, reviewer)
				}
			}
			if tc.expectRemovedReviewers != nil {
				for _, reviewer := range removedReviewers {
					assert.Contains(t, tc.expectRemovedReviewers, reviewer)
				}
			}
		})
	}
}
