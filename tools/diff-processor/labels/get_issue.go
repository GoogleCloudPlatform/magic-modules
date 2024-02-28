package labels

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	labeler "github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/labeler"
)

func GetIssue(repository string, id uint64) (labeler.Issue, error) {
	var issue labeler.Issue
	client := &http.Client{}
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%d", repository, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return issue, fmt.Errorf("Error creating request: %w", err)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN_MAGIC_MODULES"))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	resp, err := client.Do(req)
	if err != nil {
		return issue, fmt.Errorf("Error getting issue: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return issue, fmt.Errorf("Error reading response body: %w", err)
	}

	err = json.Unmarshal(body, &issue)
	if err != nil {
		var errorResponse labeler.ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return issue, fmt.Errorf("Error unmarshalling response body: %w", err)
		}
		return issue, fmt.Errorf("Error from API: %s", errorResponse.Message)
	}

	return issue, nil
}
