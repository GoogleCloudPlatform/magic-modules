package github

import (
	"fmt"
	utils "magician/utility"
	"net/http"
)

func (gh *github) PostBuildStatus(prNumber, title, state, targetURL, commitSha string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/%s", commitSha)

	postBody := map[string]string{
		"context":    title,
		"state":      state,
		"target_url": targetURL,
	}

	_, err := utils.RequestCall(url, "POST", gh.token, nil, postBody)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully posted build status to pull request %s\n", prNumber)

	return nil
}

func (gh *github) PostComment(prNumber, comment string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/comments", prNumber)

	body := map[string]string{
		"body": comment,
	}

	reqStatusCode, err := utils.RequestCall(url, "POST", gh.token, nil, body)
	if err != nil {
		return err
	}

	if reqStatusCode != http.StatusCreated {
		return fmt.Errorf("error posting comment for PR %s", prNumber)
	}

	fmt.Printf("Successfully posted comment to pull request %s\n", prNumber)

	return nil
}

func (gh *github) RequestPullRequestReviewer(prNumber, assignee string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/%s/requested_reviewers", prNumber)

	body := map[string][]string{
		"reviewers":      {assignee},
		"team_reviewers": {},
	}

	reqStatusCode, err := utils.RequestCall(url, "POST", gh.token, nil, body)
	if err != nil {
		return err
	}

	if reqStatusCode != http.StatusCreated {
		return fmt.Errorf("error adding reviewer for PR %s", prNumber)
	}

	fmt.Printf("Successfully added reviewer %s to pull request %s\n", assignee, prNumber)

	return nil
}

func (gh *github) AddLabel(prNumber, label string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/labels", prNumber)

	body := map[string][]string{
		"labels": {label},
	}
	_, err := utils.RequestCall(url, "POST", gh.token, nil, body)

	if err != nil {
		return fmt.Errorf("failed to add %s label: %s", label, err)
	}

	return nil

}

func (gh *github) RemoveLabel(prNumber, label string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/labels/%s", prNumber, label)
	_, err := utils.RequestCall(url, "DELETE", gh.token, nil, nil)

	if err != nil {
		return fmt.Errorf("failed to remove %s label: %s", label, err)
	}

	return nil
}

func (gh *github) CreateWorkflowDispatchEvent(workflowFileName string, inputs map[string]any) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/actions/workflows/%s/dispatches", workflowFileName)
	resp, err := utils.RequestCall(url, "POST", gh.token, nil, map[string]any{
		"ref":    "main",
		"inputs": inputs,
	})

	if resp != 200 && resp != 204 {
		return fmt.Errorf("server returned %d creating workflow dispatch event", resp)
	}

	if err != nil {
		return fmt.Errorf("failed to create workflow dispatch event: %s", err)
	}

	fmt.Printf("Successfully created workflow dispatch event for %s with inputs %v", workflowFileName, inputs)

	return nil
}
