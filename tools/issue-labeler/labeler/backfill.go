package labeler

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/google/go-github/v68/github"
)

type Label struct {
	Name string
}

type IssueUpdate struct {
	Number    int
	Labels    []string
	OldLabels []string
}

func GetIssues(repository, since string) ([]*github.Issue, error) {
	client := newGitHubClient()
	owner, repo, err := splitRepository(repository)
	if err != nil {
		return nil, fmt.Errorf("invalid repository format: %w", err)
	}

	sinceTime, err := time.Parse("2006-01-02", since) // input format YYYY-MM-DD
	if err != nil {
		return nil, fmt.Errorf("invalid since time format: %w", err)
	}

	opt := &github.IssueListByRepoOptions{
		Since:     sinceTime,
		State:     "all",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allIssues []*github.Issue
	ctx := context.Background()

	issues, resp, err := client.Issues.ListByRepo(ctx, owner, repo, opt)
	if err != nil {
		return nil, fmt.Errorf("listing issues: %w", err)
	}
	allIssues = append(allIssues, issues...)

	for {
		// use link headers instead of page parameter based pagination as
		// it is not supported for large datasets

		next := parseNextLink(resp.Response)
		if next == "" {
			break
		}

		req, err := client.NewRequest("GET", next, nil)
		if err != nil {
			return allIssues, err
		}
		req.Header.Set("Accept", "application/vnd.github.raw+json")

		var issues []*github.Issue
		resp, err = client.Do(ctx, req, &issues)
		if err != nil {
			return allIssues, err
		}

		allIssues = append(allIssues, issues...)
	}

	return allIssues, nil
}

// parseNextLink finds the next page for a GitHub API request by parsing the previous response's Link header.
// https://docs.github.com/en/rest/using-the-rest-api/using-pagination-in-the-rest-api?apiVersion=2022-11-28#using-link-headers
func parseNextLink(resp *http.Response) string {
	var next string
	for _, hdr := range resp.Header.Values("Link") {
		links := strings.Split(hdr, ",")
		for _, link := range links {
			pair := strings.Split(strings.TrimSpace(link), ";")
			if len(pair) == 2 {
				if strings.TrimSpace(pair[0]) == `rel="next"` {
					next = strings.Trim(pair[1], "<> ")
				} else if strings.TrimSpace(pair[1]) == `rel="next"` {
					next = strings.Trim(pair[0], "<> ")
				}
				if next != "" {
					break
				}
			}
		}
	}
	return next
}

// ComputeIssueUpdates remains the same as it doesn't interact with GitHub API
func ComputeIssueUpdates(issues []*github.Issue, regexpLabels []RegexpLabel) []IssueUpdate {
	var issueUpdates []IssueUpdate

	for _, issue := range issues {
		// Skip pull requests
		if issue.IsPullRequest() {
			continue
		}

		desired := make(map[string]struct{})
		for _, existing := range issue.Labels {
			desired[*existing.Name] = struct{}{}
		}

		_, terraform := desired["service/terraform"]
		_, linked := desired["forward/linked"]
		_, exempt := desired["forward/exempt"]
		_, testfailure := desired["test-failure"]
		if terraform || exempt {
			continue
		}

		// Decision was made to no longer add new service labels to linked tickets, because it is
		// more difficult to know which teams have received those tickets and which haven't.
		// Forwarding a ticket to a different service team should involve removing the old service
		// label and `linked` label.
		if linked {
			continue
		}

		var issueUpdate IssueUpdate
		for label := range desired {
			issueUpdate.OldLabels = append(issueUpdate.OldLabels, label)
		}
		sort.Strings(issueUpdate.OldLabels)

		affectedResources := ExtractAffectedResources(issue.GetBody())
		for _, needed := range ComputeLabels(affectedResources, regexpLabels) {
			desired[needed] = struct{}{}
		}

		if len(desired) > len(issueUpdate.OldLabels) {
			// Forwarding test failure ticket directly
			if !linked && !testfailure {
				issueUpdate.Labels = append(issueUpdate.Labels, "forward/review")
			}
			for label := range desired {
				issueUpdate.Labels = append(issueUpdate.Labels, label)
			}
			sort.Strings(issueUpdate.Labels)

			issueUpdate.Number = issue.GetNumber()
			if issueUpdate.Number > 0 {
				issueUpdates = append(issueUpdates, issueUpdate)
			}
		}
	}

	return issueUpdates
}

func UpdateIssues(repository string, issueUpdates []IssueUpdate, dryRun bool) error {
	client := newGitHubClient()
	owner, repo, err := splitRepository(repository)
	if err != nil {
		return fmt.Errorf("invalid repository format: %w", err)
	}

	ctx := context.Background()
	failed := 0

	for _, update := range issueUpdates {
		fmt.Printf("Existing labels: %v\n", update.OldLabels)
		fmt.Printf("New labels: %v\n", update.Labels)
		fmt.Printf("Updating issue: https://github.com/%s/issues/%d\n", repository, update.Number)
		if dryRun {
			continue
		}
		_, _, err := client.Issues.Edit(ctx, owner, repo, int(update.Number), &github.IssueRequest{
			Labels: &update.Labels,
		})

		if err != nil {
			glog.Errorf("Error updating issue %d: %v", update.Number, err)
			failed++
			continue
		}

		fmt.Printf("GitHub Issue %s %d updated successfully\n", repository, update.Number)
	}

	if failed > 0 {
		return fmt.Errorf("failed to update %d / %d issues", failed, len(issueUpdates))
	}
	return nil
}
