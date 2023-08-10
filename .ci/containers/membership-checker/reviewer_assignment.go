package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	_ "embed"
)

var (
	//go:embed REVIEWER_ASSIGNMENT_COMMENT.md
	reviewerAssignmentComment string
)

// Returns a list of users to request review from, as well as a new primary reviewer if this is the first run.
func chooseReviewers(firstRequestedReviewer string, previouslyInvolvedReviewers []string) (reviewersToRequest []string, newPrimaryReviewer string) {
	hasPrimaryReviewer := false
	newPrimaryReviewer = ""

	if firstRequestedReviewer != "" {
		hasPrimaryReviewer = true
	}

	if previouslyInvolvedReviewers != nil {
		for _, reviewer := range previouslyInvolvedReviewers {
			if isTeamReviewer(reviewer) {
				hasPrimaryReviewer = true
				reviewersToRequest = append(reviewersToRequest, reviewer)
			}
		}
	}

	if !hasPrimaryReviewer {
		newPrimaryReviewer = getRandomReviewer()
		reviewersToRequest = append(reviewersToRequest, newPrimaryReviewer)
	}

	return reviewersToRequest, newPrimaryReviewer
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

func formatReviewerComment(newPrimaryReviewer string, authorUserType userType, trusted bool) string {
	tmpl, err := template.New("REVIEWER_ASSIGNMENT_COMMENT.md").Parse(reviewerAssignmentComment)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse REVIEWER_ASSIGNMENT_COMMENT.md: %s", err))
	}
	sb := new(strings.Builder)
	tmpl.Execute(sb, map[string]interface{}{
		"reviewer":       newPrimaryReviewer,
		"authorUserType": authorUserType.String(),
		"trusted":        trusted,
	})
	return sb.String()
}

func postComment(prNumber, comment, GITHUB_TOKEN string, authorUserType userType) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/comments", prNumber)

	body := map[string]string{
		"body": comment,
	}

	reqStatusCode, err := requestCall(url, "POST", GITHUB_TOKEN, nil, body)
	if err != nil {
		return err
	}

	if reqStatusCode != http.StatusCreated {
		return fmt.Errorf("Error posting comment for PR %s", prNumber)
	}

	fmt.Printf("Successfully posted comment to pull request %s\n", prNumber)

	return nil
}
