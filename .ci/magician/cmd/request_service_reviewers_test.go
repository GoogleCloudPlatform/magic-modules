package cmd

import (
	"github.com/stretchr/testify/assert"
	"magician/github"
	"testing"
)

var enrolledTeamsYaml = []byte(`
service/google-x:
    team: google-x
    resources:
    - google_x_resource
service/google-y:
    team: google-y
    resources:
    - google_y_resource
service/google-z:
    resources:
    - google_z_resource
`)

func TestExecRequestServiceReviewersMembershipChecker(t *testing.T) {
	cases := map[string]struct {
		pullRequest             github.PullRequest
		requestedReviewers      []string
		previousReviewers       []string
		teamMembers             map[string][]string
		expectSpecificReviewers []string
	}{
		"no service labels means no service team reviewers": {
			pullRequest: github.PullRequest{
				User: github.User{Login: "googler_author"},
			},
			expectSpecificReviewers: []string{},
		},
		"unregistered service labels will not trigger review": {
			pullRequest: github.PullRequest{
				User:   github.User{Login: "googler_author"},
				Labels: []github.Label{{Name: "service/google-a"}},
			},
			expectSpecificReviewers: []string{},
		},
		"no configured team means no reviewers": {
			pullRequest: github.PullRequest{
				User:   github.User{Login: "googler_author"},
				Labels: []github.Label{{Name: "service/google-z"}},
			},
			expectSpecificReviewers: []string{},
		},
		"no previous reviewers means all reviews will be requested": {
			pullRequest: github.PullRequest{
				User:   github.User{Login: "googler_author"},
				Labels: []github.Label{{Name: "service/google-x"}},
			},
			teamMembers:             map[string][]string{"google-x": []string{"googler_team_member"}},
			expectSpecificReviewers: []string{"googler_team_member"},
		},
		"previous reviewers will be re-requested": {
			pullRequest: github.PullRequest{
				User:   github.User{Login: "googler_author"},
				Labels: []github.Label{{Name: "service/google-x"}},
			},
			previousReviewers:       []string{"googler_team_member"},
			teamMembers:             map[string][]string{"google-x": []string{"googler_team_member", "googler_team_member_2", "googler_team_member_3", "googler_team_member_4", "googler_team_member_5"}},
			expectSpecificReviewers: []string{"googler_team_member"},
		},
		"active reviewers will not be re-requested": {
			pullRequest: github.PullRequest{
				User:   github.User{Login: "googler_author"},
				Labels: []github.Label{{Name: "service/google-x"}},
			},
			requestedReviewers:      []string{"googler_team_member"},
			teamMembers:             map[string][]string{"google-x": []string{"googler_team_member"}},
			expectSpecificReviewers: []string{},
		},
		"authors will not be requested on their own PRs": {
			pullRequest: github.PullRequest{
				User:   github.User{Login: "googler_team_member"},
				Labels: []github.Label{{Name: "service/google-x"}},
			},
			teamMembers:             map[string][]string{"google-x": []string{"googler_team_member"}},
			expectSpecificReviewers: []string{},
		},
		"other team members be requested even if the author is excluded": {
			pullRequest: github.PullRequest{
				User:   github.User{Login: "googler_team_member"},
				Labels: []github.Label{{Name: "service/google-x"}},
			},
			teamMembers:             map[string][]string{"google-x": []string{"googler_team_member", "googler_team_member_2"}},
			expectSpecificReviewers: []string{"googler_team_member_2"},
		},
		"multiple teams can be requested at once": {
			pullRequest: github.PullRequest{
				User:   github.User{Login: "googler_author"},
				Labels: []github.Label{{Name: "service/google-x"}, {Name: "service/google-y"}, {Name: "service/google-z"}},
			},
			teamMembers:             map[string][]string{"google-x": []string{"googler_team_member"}, "google-y": []string{"googler_y_team_member"}},
			expectSpecificReviewers: []string{"googler_team_member", "googler_y_team_member"},
		},
		">3 service teams will not be requested": {
			pullRequest: github.PullRequest{
				User:   github.User{Login: "googler_author"},
				Labels: []github.Label{{Name: "service/google-x"}, {Name: "service/google-y"}, {Name: "service/google-z"}, {Name: "service/google-a"}},
			},
			teamMembers:             map[string][]string{"google-x": []string{"googler_team_member"}, "google-y": []string{"googler_y_team_member"}},
			expectSpecificReviewers: []string{},
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			requestedReviewers := []github.User{}
			for _, login := range tc.requestedReviewers {
				requestedReviewers = append(requestedReviewers, github.User{Login: login})
			}
			previousReviewers := []github.User{}
			for _, login := range tc.previousReviewers {
				previousReviewers = append(previousReviewers, github.User{Login: login})
			}
			teamMembers := map[string][]github.User{}
			for team, logins := range tc.teamMembers {
				teamMembers[team] = []github.User{}
				for _, login := range logins {
					teamMembers[team] = append(teamMembers[team], github.User{Login: login})
				}
			}
			gh := &mockGithub{
				pullRequest:        tc.pullRequest,
				requestedReviewers: requestedReviewers,
				previousReviewers:  previousReviewers,
				teamMembers:        teamMembers,
				calledMethods:      make(map[string][][]any),
			}

			execRequestServiceReviewers("1", gh, enrolledTeamsYaml)

			actualReviewers := []string{}
			for _, args := range gh.calledMethods["RequestPullRequestReviewer"] {
				actualReviewers = append(actualReviewers, args[1].(string))
			}

			if tc.expectSpecificReviewers != nil {
				assert.ElementsMatch(t, tc.expectSpecificReviewers, actualReviewers)
			}
		})
	}
}
