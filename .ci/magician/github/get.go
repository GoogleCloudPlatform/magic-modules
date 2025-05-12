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
	utils "magician/utility"
	"time"
)

type User struct {
	Login string `json:"login"`
}

type Label struct {
	Name string `json:"name"`
}

type PullRequest struct {
	HTMLUrl        string  `json:"html_url"`
	Number         int     `json:"number"`
	Title          string  `json:"title"`
	User           User    `json:"user"`
	Body           string  `json:"body"`
	Labels         []Label `json:"labels"`
	MergeCommitSha string  `json:"merge_commit_sha"`
	Merged         bool    `json:"merged"`
}

type PullRequestComment struct {
	User      User      `json:"user"`
	Body      string    `json:"body"`
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func (gh *Client) GetPullRequest(prNumber string) (PullRequest, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s", prNumber)

	var pullRequest PullRequest

	err := utils.RequestCallWithRetry(url, "GET", gh.token, &pullRequest, nil)

	return pullRequest, err
}

func (gh *Client) GetPullRequests(state, base, sort, direction string) ([]PullRequest, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls?state=%s&base=%s&sort=%s&direction=%s", state, base, sort, direction)

	var pullRequests []PullRequest

	err := utils.RequestCallWithRetry(url, "GET", gh.token, &pullRequests, nil)

	return pullRequests, err
}

func (gh *Client) GetPullRequestRequestedReviewers(prNumber string) ([]User, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/requested_reviewers", prNumber)

	var requestedReviewers struct {
		Users []User `json:"users"`
	}

	err := utils.RequestCallWithRetry(url, "GET", gh.token, &requestedReviewers, nil)
	if err != nil {
		return nil, err
	}

	return requestedReviewers.Users, nil
}

func (gh *Client) GetPullRequestPreviousReviewers(prNumber string) ([]User, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/reviews", prNumber)

	var reviews []struct {
		User User `json:"user"`
	}

	err := utils.RequestCallWithRetry(url, "GET", gh.token, &reviews, nil)
	if err != nil {
		return nil, err
	}

	previousAssignedReviewers := map[string]User{}
	for _, review := range reviews {
		previousAssignedReviewers[review.User.Login] = review.User
	}

	result := []User{}
	for _, user := range previousAssignedReviewers {
		result = append(result, user)
	}

	return result, nil
}

func (gh *Client) GetCommitMessage(owner, repo, sha string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, sha)

	var commit struct {
		Commit struct {
			Message string `json:"message"`
		} `json:"commit"`
	}

	err := utils.RequestCall(url, "GET", gh.token, &commit, nil)
	if err != nil {
		return "", err
	}

	return commit.Commit.Message, nil
}

func (gh *Client) GetPullRequestComments(prNumber string) ([]PullRequestComment, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/comments", prNumber)

	var comments []PullRequestComment
	err := utils.RequestCallWithRetry(url, "GET", gh.token, &comments, nil)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (gh *Client) GetTeamMembers(organization, team string) ([]User, error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/teams/%s/members", organization, team)

	var members []User
	err := utils.RequestCallWithRetry(url, "GET", gh.token, &members, nil)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (gh *Client) IsOrgMember(author, org string) bool {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/members/%s", org, author)
	err := utils.RequestCallWithRetry(url, "GET", gh.token, nil, nil)
	return err == nil
}

func (gh *Client) IsTeamMember(organization, teamSlug, username string) bool {
	type TeamMembership struct {
		URL   string `json:"url"`
		Role  string `json:"role"`
		State string `json:"state"`
	}

	url := fmt.Sprintf("https://api.github.com/orgs/%s/teams/%s/memberships/%s", organization, teamSlug, username)
	var membership TeamMembership
	err := utils.RequestCallWithRetry(url, "GET", gh.token, &membership, nil)

	if err != nil {
		return false
	}

	// User is considered a member if state is "active"
	return membership.State == "active"
}
