package main

import (
	"fmt"
	"net/http"
	"strings"
)

func requestReviewer(author, prNumber, GITHUB_TOKEN string) error {
	if isTeamMember(author) {
		fmt.Println("author is a team member, not assigning")
		return nil
	}

	firstRequestedReviewer, err := getPullRequestRequestedReviewer(prNumber, GITHUB_TOKEN)
	if err != nil {
		return err
	}

	previouslyInvolvedReviewers, err := getPullRequestPreviousAssignedReviewers(prNumber, GITHUB_TOKEN)
	if err != nil {
		return err
	}

	foundTeamReviewer := false

	if firstRequestedReviewer != "" {
		foundTeamReviewer = true
	}

	if previouslyInvolvedReviewers != nil {
		for _, reviewer := range previouslyInvolvedReviewers {
			if isTeamReviewer(reviewer) {
				foundTeamReviewer = true
				err = requestPullRequestReviewer(prNumber, reviewer, GITHUB_TOKEN)
				if err != nil {
					return err
				}
			}
		}
	}

	if !foundTeamReviewer {
		err = requestRandomReviewer(prNumber, GITHUB_TOKEN)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPullRequestAuthor(prNumber, GITHUB_TOKEN string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s", prNumber)

	var pullRequest struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	}

	_, err := requestCall(url, "GET", GITHUB_TOKEN, &pullRequest, nil)
	if err != nil {
		return "", err
	}

	return pullRequest.User.Login, nil
}

func getPullRequestRequestedReviewer(prNumber, GITHUB_TOKEN string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/requested_reviewers", prNumber)

	var requestedReviewers struct {
		Users []struct {
			Login string `json:"login"`
		} `json:"users"`
	}

	_, err := requestCall(url, "GET", GITHUB_TOKEN, &requestedReviewers, nil)
	if err != nil {
		return "", err
	}

	if requestedReviewers.Users == nil || len(requestedReviewers.Users) == 0 {
		return "", nil
	}

	return requestedReviewers.Users[0].Login, nil
}

func getPullRequestPreviousAssignedReviewers(prNumber, GITHUB_TOKEN string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/reviews", prNumber)

	var reviews []struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	}

	_, err := requestCall(url, "GET", GITHUB_TOKEN, &reviews, nil)
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

func requestPullRequestReviewer(prNumber, assignee, GITHUB_TOKEN string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/requested_reviewers", prNumber)

	body := map[string][]string{
		"reviewers":      []string{assignee},
		"team_reviewers": []string{},
	}

	reqStatusCode, err := requestCall(url, "POST", GITHUB_TOKEN, nil, body)
	if err != nil {
		return err
	}

	if reqStatusCode != http.StatusCreated {
		return fmt.Errorf("Error adding reviewer for PR %s", prNumber)
	}

	fmt.Printf("Successfully added reviewer %s to pull request %s\n", assignee, prNumber)

	return nil
}

func requestRandomReviewer(prNumber, GITHUB_TOKEN string) error {
	assignee := getRandomReviewer()
	err := requestPullRequestReviewer(prNumber, assignee, GITHUB_TOKEN)
	if err != nil {
		return err
	}
	err = postComment(prNumber, assignee, GITHUB_TOKEN)
	if err != nil {
		return err
	}
	return nil

}

func postComment(prNumber, reviewer, GITHUB_TOKEN string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/comments", prNumber)
	comment, err := readFile(".ci/containers/membership-checker/REVIEWER_ASSIGNMENT_COMMENT.md")
	if err != nil {
		return err
	}

	comment = strings.Replace(comment, "{{reviewer}}", reviewer, 1)

	body := map[string]string{
		"body": comment,
	}

	reqStatusCode, err := requestCall(url, "POST", GITHUB_TOKEN, nil, body)
	if err != nil {
		return err
	}

	if reqStatusCode != http.StatusCreated {
		return fmt.Errorf("Error posting reviewer assignment comment for PR %s", prNumber)
	}

	fmt.Printf("Successfully posted reviewer assignment comment to pull request %s\n", prNumber)

	return nil
}
