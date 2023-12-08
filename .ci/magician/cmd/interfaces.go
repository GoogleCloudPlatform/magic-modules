package cmd

import (
	"magician/github"
)

type GithubClient interface {
	GetPullRequest(prNumber string) (github.PullRequest, error)
	GetPullRequestRequestedReviewer(prNumber string) (string, error)
	GetPullRequestPreviousAssignedReviewers(prNumber string) ([]string, error)
	GetUserType(user string) github.UserType
	PostBuildStatus(prNumber, title, state, targetURL, commitSha string) error
	PostComment(prNumber, comment string) error
	RequestPullRequestReviewer(prNumber, assignee string) error
	AddLabel(prNumber, label string) error
	RemoveLabel(prNumber, label string) error
	CreateWorkflowDispatchEvent(workflowFileName string, inputs map[string]any) error
}

type CloudbuildClient interface {
	ApproveCommunityChecker(prNumber, commitSha string) error
	GetAwaitingApprovalBuildLink(prNumber, commitSha string) (string, error)
	TriggerMMPresubmitRuns(commitSha string, substitutions map[string]string) error
}

type ExecRunner interface {
	GetCWD() string
	Copy(src, dest string) error
	RemoveAll(path string) error
	PushDir(path string) error
	PopDir() error
	Run(name string, args, env []string) (string, error)
	MustRun(name string, args, env []string) string
}
