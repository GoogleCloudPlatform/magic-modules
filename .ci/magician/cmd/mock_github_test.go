/*
* Copyright 2023 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package cmd

import (
	"errors"

	"magician/github"

	ghi "github.com/google/go-github/v68/github"
)

type mockGithub struct {
	pullRequest         *ghi.PullRequest
	userType            github.UserType
	requestedReviewers  []*ghi.User
	previousReviewers   []*ghi.User
	pullRequestComments []*ghi.IssueComment
	teamMembers         map[string][]*ghi.User
	calledMethods       map[string][][]any
	commitMessage       string
}

func (m *mockGithub) GetPullRequest(prNumber string) (*ghi.PullRequest, error) {
	m.calledMethods["GetPullRequest"] = append(m.calledMethods["GetPullRequest"], []any{prNumber})
	return m.pullRequest, nil
}

func (m *mockGithub) GetPullRequests(state, base, sort, direction string) ([]*ghi.PullRequest, error) {
	m.calledMethods["GetPullRequests"] = append(m.calledMethods["GetPullRequests"], []any{state, base, sort, direction})
	return []*ghi.PullRequest{m.pullRequest}, nil
}

func (m *mockGithub) GetUserType(user string) github.UserType {
	m.calledMethods["GetUserType"] = append(m.calledMethods["GetUserType"], []any{user})
	return m.userType
}

func (m *mockGithub) GetPullRequestRequestedReviewers(prNumber string) ([]*ghi.User, error) {
	m.calledMethods["GetPullRequestRequestedReviewers"] = append(m.calledMethods["GetPullRequestRequestedReviewers"], []any{prNumber})
	return m.requestedReviewers, nil
}

func (m *mockGithub) GetPullRequestPreviousReviewers(prNumber string) ([]*ghi.User, error) {
	m.calledMethods["GetPullRequestPreviousReviewers"] = append(m.calledMethods["GetPullRequestPreviousReviewers"], []any{prNumber})
	return m.previousReviewers, nil
}

func (m *mockGithub) GetPullRequestComments(prNumber string) ([]*ghi.IssueComment, error) {
	m.calledMethods["GetPullRequestComments"] = append(m.calledMethods["GetPullRequestComments"], []any{prNumber})
	return m.pullRequestComments, nil
}

func (m *mockGithub) GetCommitMessage(owner, repo, sha string) (string, error) {
	m.calledMethods["GetCommitMessage"] = append(m.calledMethods["GetCommitMessage"], []any{owner, repo, sha})
	return m.commitMessage, nil
}

func (m *mockGithub) GetTeamMembers(organization, team string) ([]*ghi.User, error) {
	m.calledMethods["GetTeamMembers"] = append(m.calledMethods["GetTeamMembers"], []any{organization, team})
	if team == "" {
		return nil, errors.New("No team members set")
	}
	return m.teamMembers[team], nil
}

func (m *mockGithub) RequestPullRequestReviewers(prNumber string, reviewers []string) error {
	m.calledMethods["RequestPullRequestReviewers"] = append(m.calledMethods["RequestPullRequestReviewers"], []any{prNumber, reviewers})
	return nil
}

func (m *mockGithub) RemovePullRequestReviewers(prNumber string, reviewers []string) error {
	m.calledMethods["RemovePullRequestReviewers"] = append(m.calledMethods["RemovePullRequestReviewers"], []any{prNumber, reviewers})
	return nil
}

func (m *mockGithub) PostComment(prNumber string, comment string) error {
	m.calledMethods["PostComment"] = append(m.calledMethods["PostComment"], []any{prNumber, comment})
	return nil
}

func (m *mockGithub) UpdateComment(prNumber, comment string, id int) error {
	m.calledMethods["UpdateComment"] = append(m.calledMethods["UpdateComment"], []any{prNumber, comment, id})
	return nil
}

func (m *mockGithub) AddLabels(prNumber string, labels []string) error {
	m.calledMethods["AddLabels"] = append(m.calledMethods["AddLabels"], []any{prNumber, labels})
	return nil
}

func (m *mockGithub) RemoveLabel(prNumber string, label string) error {
	m.calledMethods["RemoveLabel"] = append(m.calledMethods["RemoveLabel"], []any{prNumber, label})
	return nil
}

func (m *mockGithub) PostBuildStatus(prNumber string, title string, state string, targetUrl string, commitSha string) error {
	m.calledMethods["PostBuildStatus"] = append(m.calledMethods["PostBuildStatus"], []any{prNumber, title, state, targetUrl, commitSha})
	return nil
}

func (m *mockGithub) CreateWorkflowDispatchEvent(workflowFileName string, inputs map[string]any) error {
	m.calledMethods["CreateWorkflowDispatchEvent"] = append(m.calledMethods["CreateWorkflowDispatchEvent"], []any{workflowFileName, inputs})
	return nil
}

func (m *mockGithub) MergePullRequest(owner, repo, prNumber, commitSha string) error {
	m.calledMethods["MergePullRequest"] = append(m.calledMethods["MergePullRequest"], []any{owner, repo, prNumber, commitSha})
	return nil
}
