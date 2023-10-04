package labels

import (
	"encoding/json"
	"fmt"

	issueLabeler "github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler"
)

func getIssue(id string) (issueLabeler.Issue, err) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.github.com/repos/hashicorp/terraform-provider-google/issues?since=%s&per_page=100&page=%d", since, page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error getting issue: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %w", err)
	}

	var issue issueLabeler.Issue
	err = json.Unmarshal(body, &issue)
	if err != nil {
		var errorResponse issueLabeler.ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling response body: %w", err)
		}
		return nil, fmt.Errorf("Error from API: %s", errorResponse.Message)
	}

	return issue, nil
}