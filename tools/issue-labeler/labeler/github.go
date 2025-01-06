package labeler

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/google/go-github/github"
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

func ensureLabelWithColor(client *github.Client, owner, repo, labelName, color string) error {
	existingLabel, err := getLabel(client, owner, repo, labelName)
	desiredColor := strings.ToUpper(color)
	ctx := context.Background()

	if err != nil {
		return fmt.Errorf("failed to check for existing label %s: %w", labelName, err)
	} else if existingLabel != nil && strings.ToUpper(existingLabel.GetColor()) != desiredColor {
		existingLabel.Color = &desiredColor
		_, _, err = client.Issues.EditLabel(ctx, owner, repo, labelName, existingLabel)
		if err != nil {
			return fmt.Errorf("failed to update label %s color: %w", labelName, err)
		}
		glog.Infof("Updated label %q color from %q to %q", labelName, existingLabel.GetColor(), color)
	} else if existingLabel == nil {
		_, _, err = client.Issues.CreateLabel(ctx, owner, repo, &github.Label{
			Name:  &labelName,
			Color: &color,
		})
		if err != nil {
			return fmt.Errorf("failed to create label %s: %w", labelName, err)
		}
		glog.Infof("Created new label %q with color %q", labelName, color)

	} else {
		glog.Infof("Label %q already exists with correct color", labelName)
	}

	return nil
}

// getLabel attempts to get an existing label
func getLabel(client *github.Client, owner, repo, labelName string) (*github.Label, error) {
	ctx := context.Background()
	label, resp, err := client.Issues.GetLabel(ctx, owner, repo, labelName)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get label: %w", err)
	}
	return label, nil
}
