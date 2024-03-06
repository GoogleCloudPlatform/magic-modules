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
	"github.com/stretchr/testify/assert"
	"magician/github"
	"testing"
)

func TestExecRequestReviewer(t *testing.T) {
	availableReviewers := github.AvailableReviewers()
	cases := map[string]struct {
		pullRequest             github.PullRequest
		requestedReviewers      []string
		previousReviewers       []string
		teamMembers             map[string][]string
		expectSpecificReviewers []string
		expectReviewersFromList []string
	}{
		"core contributor author doesn't get a new reviewer, re-request, or comment with no previous reviewers": {
			pullRequest: github.PullRequest{
				User: github.User{Login: availableReviewers[0]},
			},
			expectSpecificReviewers: []string{},
		},
		"core contributor author doesn't get a new reviewer, re-request, or comment with previous reviewers": {
			pullRequest: github.PullRequest{
				User: github.User{Login: availableReviewers[0]},
			},
			previousReviewers:       []string{availableReviewers[1]},
			expectSpecificReviewers: []string{},
		},
		"non-core-contributor author gets a new reviewer with no previous reviewers": {
			pullRequest: github.PullRequest{
				User: github.User{Login: "author"},
			},
			expectReviewersFromList: availableReviewers,
		},
		"non-core-contributor author doesn't get a new reviewer (but does get re-request) with previous reviewers": {
			pullRequest: github.PullRequest{
				User: github.User{Login: "author"},
			},
			previousReviewers:       []string{availableReviewers[1], "author2", availableReviewers[2]},
			expectSpecificReviewers: []string{availableReviewers[1], availableReviewers[2]},
		},
		"non-core-contributor author doesn't get a new reviewer or a re-request with already-requested reviewers": {
			pullRequest: github.PullRequest{
				User: github.User{Login: "author"},
			},
			requestedReviewers:      []string{availableReviewers[1], "author2", availableReviewers[2]},
			expectSpecificReviewers: []string{},
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			requestedReviewers := []github.User{}
			for _, login := range tc.requestedReviewers {
				requestedReviewers = append(requestedReviewers, github.User{Login: login})
			}
			previousReviewers := []github.User{}
			for _, login := range tc.previousReviewers {
				previousReviewers = append(previousReviewers, github.User{Login: login})
			}
			gh := &mockGithub{
				pullRequest:        tc.pullRequest,
				requestedReviewers: requestedReviewers,
				previousReviewers:  previousReviewers,
				calledMethods:      make(map[string][][]any),
			}

			execRequestReviewer("1", gh)

			actualReviewers := []string{}
			for _, args := range gh.calledMethods["RequestPullRequestReviewer"] {
				actualReviewers = append(actualReviewers, args[1].(string))
			}

			if tc.expectSpecificReviewers != nil {
				assert.ElementsMatch(t, tc.expectSpecificReviewers, actualReviewers)
			}
			if tc.expectReviewersFromList != nil {
				for _, reviewer := range actualReviewers {
					assert.Contains(t, tc.expectReviewersFromList, reviewer)
				}
			}
		})
	}
}
