package labeler

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

func newGitHubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// Helper functions
func splitRepository(repository string) (string, string, error) {
	var owner, repo string
	or := strings.Split(repository, "/")
	if len(or) != 2 {
		return "", "", fmt.Errorf("unexpected repository format %s", repository)
	}

	owner = or[0]
	repo = or[1]
	return owner, repo, nil
}

// ListLabels returns all labels for a repository
func listLabels(repository string) ([]*github.Label, error) {
	client := newGitHubClient()
	owner, repo, err := splitRepository(repository)
	if err != nil {
		return nil, fmt.Errorf("invalid repository format: %w", err)
	}

	ctx := context.Background()
	opts := &github.ListOptions{
		PerPage: 100,
	}

	var allLabels []*github.Label
	for {
		labels, resp, err := client.Issues.ListLabels(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list labels: %w", err)
		}
		allLabels = append(allLabels, labels...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allLabels, nil
}
