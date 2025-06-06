package github

import (
	gh "github.com/google/go-github/v68/github"
)

// Convert from GitHub types to our types
func convertGHPullRequest(pr *gh.PullRequest) PullRequest {
	if pr == nil {
		return PullRequest{}
	}

	var labels []Label
	if pr.Labels != nil {
		for _, l := range pr.Labels {
			if l.Name != nil {
				labels = append(labels, Label{Name: *l.Name})
			}
		}
	}

	return PullRequest{
		HTMLUrl:        pr.GetHTMLURL(),
		Number:         pr.GetNumber(),
		Title:          pr.GetTitle(),
		User:           User{Login: pr.GetUser().GetLogin()},
		Body:           pr.GetBody(),
		Labels:         labels,
		MergeCommitSha: pr.GetMergeCommitSHA(),
		Merged:         pr.GetMerged(),
	}
}

func convertGHUser(user *gh.User) User {
	if user == nil {
		return User{}
	}
	return User{
		Login: user.GetLogin(),
	}
}

func convertGHUsers(users []*gh.User) []User {
	result := make([]User, len(users))
	for i, u := range users {
		result[i] = convertGHUser(u)
	}
	return result
}

func convertGHComment(comment *gh.IssueComment) PullRequestComment {
	if comment == nil {
		return PullRequestComment{}
	}

	return PullRequestComment{
		User:      convertGHUser(comment.User),
		Body:      comment.GetBody(),
		ID:        int(comment.GetID()),
		CreatedAt: comment.GetCreatedAt().Time,
	}
}

func convertGHComments(comments []*gh.IssueComment) []PullRequestComment {
	result := make([]PullRequestComment, len(comments))
	for i, c := range comments {
		result[i] = convertGHComment(c)
	}
	return result
}
