package github

import (
	"fmt"
	"os"
)

// GithubService represents the service for GitHub interactions.
type github struct {
	token string
}

type GithubService interface {
	GetPullRequestAuthor(prNumber string) (string, error)
	GetPullRequestRequestedReviewer(prNumber string) (string, error)
	GetPullRequestPreviousAssignedReviewers(prNumber string) ([]string, error)
	GetPullRequestLabelIDs(prNumber string) (map[int]struct{}, error)
	GetUserType(user string) UserType
	PostBuildStatus(prNumber, title, state, targetURL, commitSha string) error
	PostComment(prNumber, comment string) error
	RequestPullRequestReviewer(prNumber, assignee string) error
	AddLabel(prNumber, label string) error
	RemoveLabel(prNumber, label string) error
	CreateWorkflowDispatchEvent(workflowFileName string, inputs map[string]any) error
}

func NewGithubService() GithubService {
	githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		fmt.Println("Did not provide GITHUB_TOKEN environment variable")
		os.Exit(1)
	}

	return &github{token: githubToken}
}
