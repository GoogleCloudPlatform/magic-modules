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
	"magician/utility"

	gh "github.com/google/go-github/v68/github"
)

const (
	defaultOwner = "GoogleCloudPlatform"
	defaultRepo  = "magic-modules"
)

// GetPullRequest fetches a single pull request
func (c *Client) GetPullRequest(prNumber string) (*gh.PullRequest, error) {
	num := utility.ParseInt(prNumber)
	pr, _, err := c.gh.PullRequests.Get(c.ctx, defaultOwner, defaultRepo, num)
	return pr, err
}

// GetPullRequests fetches multiple pull requests
func (c *Client) GetPullRequests(state, base, sort, direction string) ([]*gh.PullRequest, error) {
	opts := &gh.PullRequestListOptions{
		State:     state,
		Base:      base,
		Sort:      sort,
		Direction: direction,
	}

	prs, _, err := c.gh.PullRequests.List(c.ctx, defaultOwner, defaultRepo, opts)
	return prs, err
}

// GetPullRequestRequestedReviewers gets requested reviewers for a PR
func (c *Client) GetPullRequestRequestedReviewers(prNumber string) ([]*gh.User, error) {
	num := utility.ParseInt(prNumber)
	reviewers, _, err := c.gh.PullRequests.ListReviewers(c.ctx, defaultOwner, defaultRepo, num, nil)
	if err != nil {
		return nil, err
	}

	return reviewers.Users, nil
}

// GetPullRequestPreviousReviewers gets previous reviewers for a PR
func (c *Client) GetPullRequestPreviousReviewers(prNumber string) ([]*gh.User, error) {
	num := utility.ParseInt(prNumber)
	reviews, _, err := c.gh.PullRequests.ListReviews(c.ctx, defaultOwner, defaultRepo, num, nil)
	if err != nil {
		return nil, err
	}

	// Use a map to deduplicate reviewers
	reviewerMap := make(map[string]*gh.User)

	for _, review := range reviews {
		if review.User != nil && review.User.Login != nil {
			login := review.User.GetLogin()
			reviewerMap[login] = review.User
		}
	}

	// Convert map to slice
	reviewers := make([]*gh.User, 0, len(reviewerMap))
	for _, user := range reviewerMap {
		reviewers = append(reviewers, user)
	}

	return reviewers, nil
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

// GetPullRequestComments gets comments on a PR
func (c *Client) GetPullRequestComments(prNumber string) ([]*gh.IssueComment, error) {
	num := utility.ParseInt(prNumber)
	comments, _, err := c.gh.Issues.ListComments(c.ctx, defaultOwner, defaultRepo, num, nil)
	return comments, err
}

// GetTeamMembers gets members of a team
func (c *Client) GetTeamMembers(organization, team string) ([]*gh.User, error) {
	members, _, err := c.gh.Teams.ListTeamMembersBySlug(c.ctx, organization, team, nil)
	return members, err
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
