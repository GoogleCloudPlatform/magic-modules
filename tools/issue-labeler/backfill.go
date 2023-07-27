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

type Issue struct {
	Number int
	Body   string
	Labels []Label
}

type Label struct {
	Name string
}

type IssueUpdate struct {
	Labels []string `json:"labels"`
}

func backfill(since string, enrolledTeams map[string][]string, dryRun bool) {
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
			done = true
		} else {
			issues = append(issues, newIssues...)
			page++
		}
	}
	for _, issue := range issues {
		desired := make(map[string]struct{})
		for _, existing := range issue.Labels {
			desired[existing.Name] = struct{}{}
		}
		oldLength := len(desired)

		for _, needed := range labels(issue.Body, enrolledTeams, false) {
			desired[needed] = struct{}{}
		}

		if len(desired) > oldLength {
			desiredSlice := []string{"forward/review"}
			for label := range desired {
				desiredSlice = append(desiredSlice, label)
			}
			sort.Strings(desiredSlice)
			url := fmt.Sprintf("https://api.github.com/repos/hashicorp/terraform-provider-google/issues/%d", issue.Number)
			update := IssueUpdate{Labels: desiredSlice}
			body, err := json.Marshal(update)
			if err != nil {
				glog.Exitf("Error marshalling json: %v", err)
			}
			buf := bytes.NewReader(body)
			req, err := http.NewRequest("PATCH", url, buf)
			if err != nil {
				glog.Exitf("Error creating request: %v", err)
			}
			if dryRun {
				fmt.Printf("%s %s\n", req.Method, req.URL)
				b, err := json.MarshalIndent(update, "", "  ")
				if err != nil {
					glog.Exitf("Error marshalling json: %v", err)
				}
				fmt.Println(string(b))
			} else {
				_, err = client.Do(req)
				if err != nil {
					glog.Exitf("Error updating issue: %v", err)
				}
			}
		}
	}
}
