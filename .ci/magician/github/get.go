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
	"strconv"
	"time"

	gh "github.com/google/go-github/v68/github"
)

const (
	defaultOwner = "GoogleCloudPlatform"
	defaultRepo  = "magic-modules"
)

// Types for external interface compatibility
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

// GetPullRequest fetches a single pull request
func (c *Client) GetPullRequest(prNumber string) (PullRequest, error) {
	num, err := strconv.Atoi(prNumber)
	if err != nil {
		return PullRequest{}, err
	}

	pr, _, err := c.gh.PullRequests.Get(c.ctx, defaultOwner, defaultRepo, num)
	if err != nil {
		return PullRequest{}, err
	}

	return convertGHPullRequest(pr), nil
}

// GetPullRequests fetches multiple pull requests
func (c *Client) GetPullRequests(state, base, sort, direction string) ([]PullRequest, error) {
	opts := &gh.PullRequestListOptions{
		State:     state,
		Base:      base,
		Sort:      sort,
		Direction: direction,
	}

	prs, _, err := c.gh.PullRequests.List(c.ctx, defaultOwner, defaultRepo, opts)
	if err != nil {
		return nil, err
	}

	result := make([]PullRequest, len(prs))
	for i, pr := range prs {
		result[i] = convertGHPullRequest(pr)
	}

	return result, nil
}

// GetPullRequestRequestedReviewers gets requested reviewers for a PR
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

// GetPullRequestPreviousReviewers gets previous reviewers for a PR
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

// GetCommitMessage gets a commit message
func (c *Client) GetCommitMessage(owner, repo, sha string) (string, error) {
	commit, _, err := c.gh.Repositories.GetCommit(c.ctx, owner, repo, sha, nil)
	if err != nil {
		return "", err
	}

	if commit.Commit != nil && commit.Commit.Message != nil {
		return *commit.Commit.Message, nil
	}

	return "", fmt.Errorf("no commit message found")
}

// GetPullRequestComments gets all comments on a PR, handling pagination
func (c *Client) GetPullRequestComments(prNumber string) ([]PullRequestComment, error) {
	num, err := strconv.Atoi(prNumber)
	if err != nil {
		return nil, err
	}

	var allComments []*gh.IssueComment
	opts := &gh.IssueListCommentsOptions{
		ListOptions: gh.ListOptions{
			PerPage: 100,
		},
	}

	for {
		comments, resp, err := c.gh.Issues.ListComments(c.ctx, defaultOwner, defaultRepo, num, opts)
		if err != nil {
			return nil, err
		}

		allComments = append(allComments, comments...)

		if resp.NextPage == 0 {
			break // No more pages
		}

		// Set up for the next page
		opts.Page = resp.NextPage
	}

	return convertGHComments(allComments), nil
}

// GetTeamMembers gets all members of a team, handling pagination
func (c *Client) GetTeamMembers(organization, team string) ([]User, error) {
	var allMembers []*gh.User
	opts := &gh.TeamListTeamMembersOptions{
		ListOptions: gh.ListOptions{
			PerPage: 100,
		},
	}

	for {
		members, resp, err := c.gh.Teams.ListTeamMembersBySlug(c.ctx, organization, team, opts)
		if err != nil {
			return nil, err
		}

		allMembers = append(allMembers, members...)

		if resp.NextPage == 0 {
			break // No more pages
		}

		// Set up for the next page
		opts.Page = resp.NextPage
	}

	return convertGHUsers(allMembers), nil
}

// IsOrgMember checks if a user is a member of an organization
func (c *Client) IsOrgMember(username, org string) bool {
	isMember, _, err := c.gh.Organizations.IsMember(c.ctx, org, username)
	if err != nil {
		return false
	}

	return isMember
}

// IsTeamMember checks if a user is a member of a team
func (c *Client) IsTeamMember(organization, teamSlug, username string) bool {
	membership, _, err := c.gh.Teams.GetTeamMembershipBySlug(c.ctx, organization, teamSlug, username)
	if err != nil {
		return false
	}

	return membership != nil && membership.State != nil && *membership.State == "active"
}
