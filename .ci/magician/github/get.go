package github

import (
	"fmt"
	utils "magician/utility"
)

func (gh *github) GetPullRequestAuthor(prNumber string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s", prNumber)

	var pullRequest struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	}

	_, err := utils.RequestCall(url, "GET", gh.token, &pullRequest, nil)
	if err != nil {
		return "", err
	}

	return pullRequest.User.Login, nil
}

func (gh *github) GetPullRequestRequestedReviewer(prNumber string) (string, error) {
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

func (gh *github) GetPullRequestPreviousAssignedReviewers(prNumber string) ([]string, error) {
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
	for key, _ := range previousAssignedReviewers {
		result = append(result, key)
	}

	return result, nil
}
