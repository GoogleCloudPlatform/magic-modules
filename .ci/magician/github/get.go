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

type User struct {
	Login string `json:"login"`
}

type Label struct {
	Name string `json:"name"`
}

type PullRequest struct {
	User struct {
		Login string `json:"login"`
	} `json:"user"`
	Labels []Label `json:"labels"`
}

func (gh *Client) GetPullRequest(prNumber string) (PullRequest, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s", prNumber)

	var pullRequest PullRequest

	_, err := utils.RequestCall(url, "GET", gh.token, &pullRequest, nil)
	if err != nil {
		return pullRequest, err
	}

	return pullRequest, nil
}

func (gh *Client) GetPullRequestRequestedReviewer(prNumber string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/requested_reviewers", prNumber)

	var requestedReviewers struct {
		Users []struct {
			Login string `json:"login"`
		} `json:"users"`
	}

	_, err := utils.RequestCall(url, "GET", gh.token, &requestedReviewers, nil)
	if err != nil {
		return "", err
	}

	if requestedReviewers.Users == nil || len(requestedReviewers.Users) == 0 {
		return "", nil
	}

	return requestedReviewers.Users[0].Login, nil
}

func (gh *Client) GetPullRequestPreviousAssignedReviewers(prNumber string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/reviews", prNumber)

	var reviews []struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	}

	_, err := utils.RequestCall(url, "GET", gh.token, &reviews, nil)
	if err != nil {
		return nil, err
	}

	previousAssignedReviewers := map[string]struct{}{}
	for _, review := range reviews {
		previousAssignedReviewers[review.User.Login] = struct{}{}
	}

	result := []string{}
	for key := range previousAssignedReviewers {
		result = append(result, key)
	}

	return result, nil
}
