package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"

	"github.com/golang/glog"
)

type ErrorResponse struct {
	Message string
}

type Issue struct {
	Number      int
	Body        string
	Labels      []Label
	PullRequest map[string]any `json:"pull_request"`
}

type Label struct {
	Name string
}

type IssueUpdate struct {
	Number    int
	Labels    []string
	OldLabels []string
}

type IssueUpdateBody struct {
	Labels []string `json:"labels"`
}

func getIssues(since string) []Issue {
	client := &http.Client{}
	done := false
	page := 1
	var issues []Issue
	for !done {
		url := fmt.Sprintf("https://api.github.com/repos/hashicorp/terraform-provider-google/issues?since=%s&per_page=100&page=%d", since, page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			glog.Exitf("Error creating request: %v", err)
		}
		req.Header.Add("Accept", "application/vnd.github+json")
		req.Header.Add("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
		req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
		resp, err := client.Do(req)
		if err != nil {
			glog.Exitf("Error listing issues: %v", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			glog.Exitf("Error reading response body: %v", err)
		}
		var newIssues []Issue
		json.Unmarshal(body, &newIssues)
		if len(newIssues) == 0 {
			var err ErrorResponse
			json.Unmarshal(body, &err)
			if err.Message == "Bad credentials" {
				glog.Exitf("Error from API: Bad credentials")
			}
			glog.Infof("API returned message: %s", err.Message)
			done = true
		} else {
			issues = append(issues, newIssues...)
			page++
		}
	}
	return issues
}

func computeIssueUpdates(issues []Issue, regexpLabels []regexpLabel) []IssueUpdate {
	var issueUpdates []IssueUpdate

	for _, issue := range issues {
		if len(issue.PullRequest) > 0 {
			continue
		}

		desired := make(map[string]struct{})
		for _, existing := range issue.Labels {
			desired[existing.Name] = struct{}{}
		}

		_, terraform := desired["service/terraform"]
		_, linked := desired["forward/linked"]
		_, exempt := desired["forward/exempt"]
		if terraform || exempt {
			continue
		}

		var issueUpdate IssueUpdate
		for label := range desired {
			issueUpdate.OldLabels = append(issueUpdate.OldLabels, label)
		}

		affectedResources := extractAffectedResources(issue.Body)
		for _, needed := range computeLabels(affectedResources, regexpLabels) {
			desired[needed] = struct{}{}
		}

		if len(desired) > len(issueUpdate.OldLabels) {
			if !linked {
				issueUpdate.Labels = append(issueUpdate.Labels, "forward/review")
			}
			for label := range desired {
				issueUpdate.Labels = append(issueUpdate.Labels, label)
			}
			sort.Strings(issueUpdate.Labels)

			issueUpdate.Number = issue.Number

			issueUpdates = append(issueUpdates, issueUpdate)
		}
	}

	return issueUpdates
}

func updateIssues(issueUpdates []IssueUpdate, dryRun bool) {
	client := &http.Client{}
	for _, issueUpdate := range issueUpdates {
		url := fmt.Sprintf("https://api.github.com/repos/hashicorp/terraform-provider-google/issues/%d", issueUpdate.Number)
		updateBody := IssueUpdateBody{Labels: issueUpdate.Labels}
		body, err := json.Marshal(updateBody)
		if err != nil {
			glog.Errorf("Error marshalling json: %v", err)
			continue
		}
		buf := bytes.NewReader(body)
		req, err := http.NewRequest("PATCH", url, buf)
		req.Header.Add("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
		req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
		if err != nil {
			glog.Errorf("Error creating request: %v", err)
			continue
		}
		fmt.Printf("Existing labels: %v\n", issueUpdate.OldLabels)
		fmt.Printf("New labels: %v\n", issueUpdate.Labels)
		fmt.Printf("%s %s (https://github.com/hashicorp/terraform-provider-google/issues/%d)\n", req.Method, req.URL, issueUpdate.Number)
		b, err := json.MarshalIndent(updateBody, "", "  ")
		if err != nil {
			glog.Errorf("Error marshalling json: %v", err)
			continue
		}
		fmt.Println(string(b))
		if !dryRun {
			resp, err := client.Do(req)
			if err != nil {
				glog.Errorf("Error updating issue: %v", err)
				continue
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				glog.Errorf("Error reading response body: %v", err)
				continue
			}
			var errResp ErrorResponse
			json.Unmarshal(body, &errResp)
			glog.Infof("API returned message: %s", errResp.Message)
		}
	}
}
