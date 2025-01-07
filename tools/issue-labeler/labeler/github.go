package labeler

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v61/github"
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
