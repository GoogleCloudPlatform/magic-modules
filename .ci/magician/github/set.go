package github

import (
	"fmt"
	utils "magician/utility"
	"net/http"
)

func PostBuildStatus(prNumber, title, state, target_url, commitSha string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/%s", commitSha)

	postBody := map[string]string{
		"context":    title,
		"state":      state,
		"target_url": target_url,
	}

	_, err := utils.RequestCall(url, "POST", github_token, nil, postBody)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully posted build status to pull request %s\n", prNumber)

	return nil
}

func PostComment(prNumber, comment string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/comments", prNumber)

	body := map[string]string{
		"body": comment,
	}

	reqStatusCode, err := utils.RequestCall(url, "POST", github_token, nil, body)
	if err != nil {
		return err
	}

	if reqStatusCode != http.StatusCreated {
		return fmt.Errorf("error posting comment for PR %s", prNumber)
	}

	fmt.Printf("Successfully posted comment to pull request %s\n", prNumber)

	return nil
}

func RequestPullRequestReviewer(prNumber, assignee string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/requested_reviewers", prNumber)

	body := map[string][]string{
		"reviewers":      {assignee},
		"team_reviewers": {},
	}

	reqStatusCode, err := utils.RequestCall(url, "POST", github_token, nil, body)
	if err != nil {
		return err
	}

	if reqStatusCode != http.StatusCreated {
		return fmt.Errorf("error adding reviewer for PR %s", prNumber)
	}

	fmt.Printf("Successfully added reviewer %s to pull request %s\n", assignee, prNumber)

	return nil
}

func AddLabel(prNumber, label string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/labels", prNumber)

	body := map[string][]string{
		"labels": {label},
	}
	_, err := utils.RequestCall(url, "POST", github_token, nil, body)

	if err != nil {
		return fmt.Errorf("failed to add %s label: %s", label, err)
	}

	return nil

}

func RemoveLabel(prNumber, label string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/labels/%s", prNumber, label)
	_, err := utils.RequestCall(url, "DELETE", github_token, nil, nil)

	if err != nil {
		return fmt.Errorf("failed to remove %s label: %s", label, err)
	}

	return nil
}
