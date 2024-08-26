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
	"magician/github"
)

type GithubClient interface {
	GetPullRequest(prNumber string) (github.PullRequest, error)
	GetPullRequests(state, base, sort, direction string) ([]github.PullRequest, error)
	GetPullRequestRequestedReviewers(prNumber string) ([]github.User, error)
	GetPullRequestPreviousReviewers(prNumber string) ([]github.User, error)
	GetUserType(user string) github.UserType
	GetTeamMembers(organization, team string) ([]github.User, error)
	MergePullRequest(owner, repo, prNumber, commitSha string) error
	PostBuildStatus(prNumber, title, state, targetURL, commitSha string) error
	PostComment(prNumber, comment string) error
	RequestPullRequestReviewers(prNumber string, reviewers []string) error
	AddLabels(prNumber string, labels []string) error
	RemoveLabel(prNumber, label string) error
	CreateWorkflowDispatchEvent(workflowFileName string, inputs map[string]any) error
}

type CloudbuildClient interface {
	ApproveCommunityChecker(prNumber, commitSha string) error
	TriggerMMPresubmitRuns(commitSha string, substitutions map[string]string) error
}

type ExecRunner interface {
	GetCWD() string
	Copy(src, dest string) error
	Mkdir(path string) error
	RemoveAll(path string) error
	PushDir(path string) error
	PopDir() error
	WriteFile(name, data string) error
	Run(name string, args []string, env map[string]string) (string, error)
	MustRun(name string, args []string, env map[string]string) string
}
