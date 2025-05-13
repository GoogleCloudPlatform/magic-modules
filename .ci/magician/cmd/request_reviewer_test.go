/*
* Copyright 2024 Google LLC. All Rights Reserved.
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
	"testing"

	ghi "github.com/google/go-github/v68/github"
	"github.com/stretchr/testify/assert"
)

func TestExecRequestReviewer(t *testing.T) {
	availableReviewers := github.AvailableReviewers(nil)
	if len(availableReviewers) < 3 {
		t.Fatalf("not enough available reviewers (%v) to run TestExecRequestReviewer (need at least 3)", availableReviewers)
	}

	cases := map[string]struct {
		pullRequest             *ghi.PullRequest
		requestedReviewers      []string
		previousReviewers       []string
		teamMembers             map[string][]string
		expectSpecificReviewers []string
		expectReviewersFromList []string
	}{
		"core contributor author doesn't get a new reviewer, re-request, or comment with no previous reviewers": {
			pullRequest: &ghi.PullRequest{
				User: &ghi.User{Login: ghi.String(availableReviewers[0])},
			},
			expectSpecificReviewers: []string{},
		},
		"core contributor author doesn't get a new reviewer, re-request, or comment with previous reviewers": {
			pullRequest: &ghi.PullRequest{
				User: &ghi.User{Login: ghi.String(availableReviewers[0])},
			},
			previousReviewers:       []string{availableReviewers[1]},
			expectSpecificReviewers: []string{},
		},
		"non-core-contributor author gets a new reviewer with no previous reviewers": {
			pullRequest: &ghi.PullRequest{
				User: &ghi.User{Login: ghi.String("author")},
			},
			expectReviewersFromList: availableReviewers,
		},
		"non-core-contributor author doesn't get a new reviewer (but does get re-request) with previous reviewers": {
			pullRequest: &ghi.PullRequest{
				User: &ghi.User{Login: ghi.String("author")},
			},
			previousReviewers:       []string{availableReviewers[1], "author2", availableReviewers[2]},
			expectSpecificReviewers: []string{availableReviewers[1], availableReviewers[2]},
		},
		"non-core-contributor author doesn't get a new reviewer or a re-request with already-requested reviewers": {
			pullRequest: &ghi.PullRequest{
				User: &ghi.User{Login: ghi.String("author")},
			},
			requestedReviewers:      []string{availableReviewers[1], "author2", availableReviewers[2]},
			expectSpecificReviewers: []string{},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			requestedReviewers := []*ghi.User{}
			for _, login := range tc.requestedReviewers {
				requestedReviewers = append(requestedReviewers, &ghi.User{Login: ghi.String(login)})
			}

			previousReviewers := []*ghi.User{}
			for _, login := range tc.previousReviewers {
				previousReviewers = append(previousReviewers, &ghi.User{Login: ghi.String(login)})
			}

			mockGH := &mockGithub{
				pullRequest:        tc.pullRequest,
				requestedReviewers: requestedReviewers,
				previousReviewers:  previousReviewers,
				calledMethods:      make(map[string][][]any),
			}

			execRequestReviewer("1", mockGH)

			actualReviewers := []string{}
			for _, args := range mockGH.calledMethods["RequestPullRequestReviewers"] {
				actualReviewers = append(actualReviewers, args[1].([]string)...)
			}

			if tc.expectSpecificReviewers != nil {
				assert.ElementsMatch(t, tc.expectSpecificReviewers, actualReviewers)
				if len(tc.expectSpecificReviewers) == 0 {
					assert.Len(t, mockGH.calledMethods["RequestPullRequestReviewers"], 0)
				}
			}

			if tc.expectReviewersFromList != nil {
				for _, reviewer := range actualReviewers {
					assert.Contains(t, tc.expectReviewersFromList, reviewer)
				}
			}
		})
	}
}
