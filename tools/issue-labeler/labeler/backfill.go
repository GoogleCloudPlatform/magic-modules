package labeler

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/golang/glog"
	"github.com/google/go-github/v61/github"
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

	sinceTime, err := time.Parse(time.RFC3339, since)
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

	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repo, opt)
		if err != nil {
			return nil, fmt.Errorf("listing issues: %w", err)
		}

		// Convert github.Issue to our Issue type
		for _, issue := range issues {
			labels := make([]Label, len(issue.Labels))
			for i, l := range issue.Labels {
				labels[i] = Label{Name: *l.Name}
			}

			allIssues = append(allIssues, issue)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allIssues, nil
}

// ComputeIssueUpdates remains the same as it doesn't interact with GitHub API
func ComputeIssueUpdates(issues []*github.Issue, regexpLabels []RegexpLabel) []IssueUpdate {
	var issueUpdates []IssueUpdate

	for _, issue := range issues {
		if !issue.IsPullRequest() {
			continue
		}

		desired := make(map[string]struct{})
		for _, existing := range issue.Labels {
			desired[*existing.Name] = struct{}{}
		}

		_, terraform := desired["service/terraform"]
		_, linked := desired["forward/linked"]
		_, exempt := desired["forward/exempt"]
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

		affectedResources := ExtractAffectedResources(*issue.Body)
		for _, needed := range ComputeLabels(affectedResources, regexpLabels) {
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

			issueUpdate.Number = *issue.Number
			issueUpdates = append(issueUpdates, issueUpdate)
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
