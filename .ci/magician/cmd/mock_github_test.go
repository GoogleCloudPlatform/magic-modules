package cmd

import "magician/github"

type mockGithub struct {
	author            string
	userType          github.UserType
	firstReviewer     string
	previousReviewers []string
	calledMethods     map[string][]any
}

func (m *mockGithub) GetPullRequestAuthor(prNumber string) (string, error) {
	m.calledMethods["GetPullRequestAuthor"] = []any{prNumber}
	return m.author, nil
}

func (m *mockGithub) GetUserType(user string) github.UserType {
	m.calledMethods["GetUserType"] = []any{user}
	return m.userType
}

func (m *mockGithub) GetPullRequestRequestedReviewer(prNumber string) (string, error) {
	m.calledMethods["GetPullRequestRequestedReviewer"] = []any{prNumber}
	return m.firstReviewer, nil
}

func (m *mockGithub) GetPullRequestPreviousAssignedReviewers(prNumber string) ([]string, error) {
	m.calledMethods["GetPullRequestPreviousAssignedReviewers"] = []any{prNumber}
	return m.previousReviewers, nil
}

func (m *mockGithub) RequestPullRequestReviewer(prNumber string, reviewer string) error {
	m.calledMethods["RequestPullRequestReviewer"] = []any{prNumber, reviewer}
	return nil
}

func (m *mockGithub) PostComment(prNumber string, comment string) error {
	m.calledMethods["PostComment"] = []any{prNumber, comment}
	return nil
}

func (m *mockGithub) AddLabel(prNumber string, label string) error {
	m.calledMethods["AddLabel"] = []any{prNumber, label}
	return nil
}

func (m *mockGithub) RemoveLabel(prNumber string, label string) error {
	m.calledMethods["RemoveLabel"] = []any{prNumber, label}
	return nil
}

func (m *mockGithub) PostBuildStatus(prNumber string, title string, state string, targetUrl string, commitSha string) error {
	m.calledMethods["PostBuildStatus"] = []any{prNumber, title, state, targetUrl, commitSha}
	return nil
}

func (m *mockGithub) CreateWorkflowDispatchEvent(workflowFileName string, inputs map[string]any) error {
	m.calledMethods["CreateWorkflowDispatchEvent"] = []any{workflowFileName, inputs}
	return nil
}
