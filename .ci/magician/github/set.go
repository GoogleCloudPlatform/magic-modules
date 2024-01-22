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
)

func (gh *Client) PostBuildStatus(prNumber, title, state, targetURL, commitSha string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/%s", commitSha)

	postBody := map[string]string{
		"context":    title,
		"state":      state,
		"target_url": targetURL,
	}

	err := utils.RequestCall(url, "POST", gh.token, nil, postBody)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully posted build status to pull request %s\n", prNumber)

	return nil
}

func (gh *Client) PostComment(prNumber, comment string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/comments", prNumber)

	body := map[string]string{
		"body": comment,
	}

	err := utils.RequestCall(url, "POST", gh.token, nil, body)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully posted comment to pull request %s\n", prNumber)

	return nil
}

func (gh *Client) RequestPullRequestReviewer(prNumber, assignee string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/requested_reviewers", prNumber)

	body := map[string][]string{
		"reviewers":      {assignee},
		"team_reviewers": {},
	}

	err := utils.RequestCall(url, "POST", gh.token, nil, body)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully added reviewer %s to pull request %s\n", assignee, prNumber)

	return nil
}

func (gh *Client) AddLabel(prNumber, label string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/labels", prNumber)

	body := map[string][]string{
		"labels": {label},
	}
	err := utils.RequestCall(url, "POST", gh.token, nil, body)

	if err != nil {
		return fmt.Errorf("failed to add %s label: %s", label, err)
	}

	return nil

}

func (gh *Client) RemoveLabel(prNumber, label string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/labels/%s", prNumber, label)
	err := utils.RequestCall(url, "DELETE", gh.token, nil, nil)

	if err != nil {
		return fmt.Errorf("failed to remove %s label: %s", label, err)
	}

	return nil
}

func (gh *Client) CreateWorkflowDispatchEvent(workflowFileName string, inputs map[string]any) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/actions/workflows/%s/dispatches", workflowFileName)
	err := utils.RequestCall(url, "POST", gh.token, nil, map[string]any{
		"ref":    "main",
		"inputs": inputs,
	})

	if err != nil {
		return fmt.Errorf("failed to create workflow dispatch event: %s", err)
	}

	fmt.Printf("Successfully created workflow dispatch event for %s with inputs %v\n", workflowFileName, inputs)

	return nil
}

func (gh *Client) MergePullRequest(owner, repo, prNumber string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%s/merge", owner, repo, prNumber)
	err := utils.RequestCall(url, "PUT", gh.token, nil, map[string]any{
		"merge_method": "squash",
	})

	if err != nil {
		return fmt.Errorf("failed to merge pull request: %s", err)
	}

	fmt.Printf("Successfully merged pull request %s\n", prNumber)

	return nil
}
