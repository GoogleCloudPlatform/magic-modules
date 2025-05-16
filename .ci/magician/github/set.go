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
	"strconv"
	"strings"
	"time"

	utils "magician/utility"

	gh "github.com/google/go-github/v68/github"
)

// PostBuildStatus creates a commit status for a specific SHA
func (c *Client) PostBuildStatus(prNumber, title, state, targetURL, commitSha string) error {
	repoStatus := &gh.RepoStatus{
		Context:   gh.Ptr(title),
		State:     gh.Ptr(state),
		TargetURL: gh.Ptr(targetURL),
	}

	_, _, err := c.gh.Repositories.CreateStatus(c.ctx, defaultOwner, defaultRepo, commitSha, repoStatus)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully posted build status to pull request %s\n", prNumber)
	return nil
}

// PostComment adds a comment to a pull request
func (c *Client) PostComment(prNumber, comment string) error {
	num, err := strconv.Atoi(prNumber)
	if err != nil {
		return err
	}

	issueComment := &gh.IssueComment{
		Body: gh.Ptr(comment),
	}

	_, _, err = c.gh.Issues.CreateComment(c.ctx, defaultOwner, defaultRepo, num, issueComment)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully posted comment to pull request %s\n", prNumber)
	return nil
}

// UpdateComment updates an existing comment
func (c *Client) UpdateComment(prNumber, comment string, id int) error {
	issueComment := &gh.IssueComment{
		Body: gh.Ptr(comment),
	}

	_, _, err := c.gh.Issues.EditComment(c.ctx, defaultOwner, defaultRepo, int64(id), issueComment)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully updated comment %d in pull request %s\n", id, prNumber)
	return nil
}

func (gh *Client) RequestPullRequestReviewers(prNumber string, reviewers []string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/requested_reviewers", prNumber)

	body := map[string][]string{
		"reviewers":      reviewers,
		"team_reviewers": {},
	}

	err := utils.RequestCallWithRetry(url, "POST", gh.token, nil, body)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully added reviewers %v to pull request %s\n", reviewers, prNumber)

	return nil
}

func (gh *Client) RemovePullRequestReviewers(prNumber string, reviewers []string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/requested_reviewers", prNumber)

	body := map[string][]string{
		"reviewers":      reviewers,
		"team_reviewers": {},
	}

	err := utils.RequestCall(url, "DELETE", gh.token, nil, body)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully removed reviewers %v to pull request %s\n", reviewers, prNumber)

	return nil
}

// AddLabels adds labels to an issue or pull request
func (c *Client) AddLabels(prNumber string, labels []string) error {
	num, err := strconv.Atoi(prNumber)
	if err != nil {
		return err
	}

	_, _, err = c.gh.Issues.AddLabelsToIssue(c.ctx, defaultOwner, defaultRepo, num, labels)
	if err != nil {
		return fmt.Errorf("failed to add %q labels: %s", labels, err)
	}

	return nil
}

// RemoveLabel removes a label from an issue or pull request
func (c *Client) RemoveLabel(prNumber, label string) error {
	num, err := strconv.Atoi(prNumber)
	if err != nil {
		return err
	}

	_, err = c.gh.Issues.RemoveLabelForIssue(c.ctx, defaultOwner, defaultRepo, num, label)
	if err != nil {
		return fmt.Errorf("failed to remove %s label: %s", label, err)
	}

	return nil
}

// CreateWorkflowDispatchEvent triggers a workflow run
func (c *Client) CreateWorkflowDispatchEvent(workflowFileName string, inputs map[string]any) error {
	stringInputs := make(map[string]interface{})
	for k, v := range inputs {
		stringInputs[k] = v
	}

	event := gh.CreateWorkflowDispatchEventRequest{
		Ref:    "main",
		Inputs: stringInputs,
	}

	_, err := c.gh.Actions.CreateWorkflowDispatchEventByFileName(c.ctx, defaultOwner, defaultRepo, workflowFileName, event)
	if err != nil {
		return fmt.Errorf("failed to create workflow dispatch event: %s", err)
	}

	fmt.Printf("Successfully created workflow dispatch event for %s with inputs %v\n", workflowFileName, inputs)
	return nil
}

func (gh *Client) MergePullRequest(owner, repo, prNumber, commitSha string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%s/merge", owner, repo, prNumber)

	err := utils.RequestCallWithRetry(url, "PUT", gh.token, nil, map[string]any{
		"merge_method": "squash",
		"sha":          commitSha,
	})

	if err != nil {
		// Check if the error is "Merge already in progress" (405)
		if strings.Contains(err.Error(), "Merge already in progress") {
			fmt.Printf("Pull request %s is already being merged\n", prNumber)
			// This status does not indicate that the Pull Request was merged
			// Try again after 20s
			time.Sleep(20 * time.Second)
			return gh.MergePullRequest(owner, repo, prNumber, commitSha)
		}
		// Check if the PR is already merged (returns 405 Pull Request is not mergeable)
		if strings.Contains(err.Error(), "Pull Request is not mergeable") {
			fmt.Printf("Pull request %s is not mergeable; checking if it was already merged\n", prNumber)
			pr, err := gh.GetPullRequest(prNumber)
			if err != nil {
				return fmt.Errorf("failed to check if PR was already merged: %w", err)
			}
			if pr.Merged {
				fmt.Printf("Pull request %s was already merged\n", prNumber)
				return nil
			}
			fmt.Printf("Pull request %s wasn't already merged\n", prNumber)
		}
		return fmt.Errorf("failed to merge pull request: %w", err)
	}

	fmt.Printf("Successfully merged pull request %s\n", prNumber)
	return nil
}
