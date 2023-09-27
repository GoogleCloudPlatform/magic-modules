package cmd

import "magician/github"

type mockGithub struct {
	author            string
	userType          github.UserType
	firstReviewer     string
	previousReviewers []string
	calledMethods     map[string]bool
}

func (m *mockGithub) GetPullRequestAuthor(string) (string, error) {
	m.calledMethods["GetPullRequestAuthor"] = true
	return m.author, nil
}

func (m *mockGithub) GetUserType(string) github.UserType {
	m.calledMethods["GetUserType"] = true
	return m.userType
}

func (m *mockGithub) GetPullRequestRequestedReviewer(string) (string, error) {
	m.calledMethods["GetPullRequestRequestedReviewer"] = true
	return m.firstReviewer, nil
}

func (m *mockGithub) GetPullRequestPreviousAssignedReviewers(string) ([]string, error) {
	m.calledMethods["GetPullRequestPreviousAssignedReviewers"] = true
	return m.previousReviewers, nil
}

func (m *mockGithub) RequestPullRequestReviewer(prNumber string, reviewer string) error {
	m.calledMethods["RequestPullRequestReviewer"] = true
	return nil
}

func (m *mockGithub) PostComment(prNumber string, comment string) error {
	m.calledMethods["PostComment"] = true
	return nil
}

func (m *mockGithub) AddLabel(prNumber string, label string) error {
	m.calledMethods["AddLabel"] = true
	return nil
}

func (m *mockGithub) RemoveLabel(prNumber string, label string) error {
	m.calledMethods["RemoveLabel"] = true
	return nil
}

func (m *mockGithub) PostBuildStatus(prNumber string, title string, state string, targetUrl string, commitSha string) error {
	m.calledMethods["PostBuildStatus"] = true
	return nil
}

func (m *mockGithub) CreateWorkflowDispatchEvent(workflowFileName string, inputs map[string]any) error {
	m.calledMethods["CreateWorkflowDispatchEvent"] = true
	return nil
}
