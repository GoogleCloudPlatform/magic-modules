package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ghi "github.com/google/go-github/v68/github"
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
		pullRequest             *ghi.PullRequest
		requestedReviewers      []string
		previousReviewers       []string
		teamMembers             map[string][]string
		expectSpecificReviewers []string
	}{
		"no service labels means no service team reviewers": {
			pullRequest: &ghi.PullRequest{
				User: &ghi.User{Login: ghi.String("googler_author")},
			},
			expectSpecificReviewers: []string{},
		},
		"unregistered service labels will not trigger review": {
			pullRequest: &ghi.PullRequest{
				User:   &ghi.User{Login: ghi.String("googler_author")},
				Labels: []*ghi.Label{{Name: ghi.String("service/google-a")}},
			},
			expectSpecificReviewers: []string{},
		},
		"no configured team means no reviewers": {
			pullRequest: &ghi.PullRequest{
				User:   &ghi.User{Login: ghi.String("googler_author")},
				Labels: []*ghi.Label{{Name: ghi.String("service/google-z")}},
			},
			expectSpecificReviewers: []string{},
		},
		"no previous reviewers means all reviews will be requested": {
			pullRequest: &ghi.PullRequest{
				User:   &ghi.User{Login: ghi.String("googler_author")},
				Labels: []*ghi.Label{{Name: ghi.String("service/google-x")}},
			},
			teamMembers:             map[string][]string{"google-x": {"googler_team_member"}},
			expectSpecificReviewers: []string{"googler_team_member"},
		},
		"previous reviewers will be re-requested": {
			pullRequest: &ghi.PullRequest{
				User:   &ghi.User{Login: ghi.String("googler_author")},
				Labels: []*ghi.Label{{Name: ghi.String("service/google-x")}},
			},
			previousReviewers:       []string{"googler_team_member"},
			teamMembers:             map[string][]string{"google-x": {"googler_team_member", "googler_team_member_2", "googler_team_member_3", "googler_team_member_4", "googler_team_member_5"}},
			expectSpecificReviewers: []string{"googler_team_member"},
		},
		"active reviewers will not be re-requested": {
			pullRequest: &ghi.PullRequest{
				User:   &ghi.User{Login: ghi.String("googler_author")},
				Labels: []*ghi.Label{{Name: ghi.String("service/google-x")}},
			},
			requestedReviewers:      []string{"googler_team_member"},
			teamMembers:             map[string][]string{"google-x": {"googler_team_member"}},
			expectSpecificReviewers: []string{},
		},
		"authors will not be requested on their own PRs": {
			pullRequest: &ghi.PullRequest{
				User:   &ghi.User{Login: ghi.String("googler_team_member")},
				Labels: []*ghi.Label{{Name: ghi.String("service/google-x")}},
			},
			teamMembers:             map[string][]string{"google-x": {"googler_team_member"}},
			expectSpecificReviewers: []string{},
		},
		"authors will not be requested on their own PRs even if they left comments on it previously": {
			pullRequest: &ghi.PullRequest{
				User:   &ghi.User{Login: ghi.String("googler_team_member")},
				Labels: []*ghi.Label{{Name: ghi.String("service/google-x")}},
			},
			teamMembers:             map[string][]string{"google-x": {"googler_team_member"}},
			previousReviewers:       []string{"googler_team_member"},
			expectSpecificReviewers: []string{},
		},
		"other team members be requested even if the author is excluded": {
			pullRequest: &ghi.PullRequest{
				User:   &ghi.User{Login: ghi.String("googler_team_member")},
				Labels: []*ghi.Label{{Name: ghi.String("service/google-x")}},
			},
			teamMembers:             map[string][]string{"google-x": {"googler_team_member", "googler_team_member_2"}},
			expectSpecificReviewers: []string{"googler_team_member_2"},
		},
		"multiple teams can be requested at once": {
			pullRequest: &ghi.PullRequest{
				User:   &ghi.User{Login: ghi.String("googler_author")},
				Labels: []*ghi.Label{{Name: ghi.String("service/google-x")}, {Name: ghi.String("service/google-y")}, {Name: ghi.String("service/google-z")}},
			},
			teamMembers:             map[string][]string{"google-x": {"googler_team_member"}, "google-y": {"googler_y_team_member"}},
			expectSpecificReviewers: []string{"googler_team_member", "googler_y_team_member"},
		},
		">3 service teams will not be requested": {
			pullRequest: &ghi.PullRequest{
				User:   &ghi.User{Login: ghi.String("googler_author")},
				Labels: []*ghi.Label{{Name: ghi.String("service/google-x")}, {Name: ghi.String("service/google-y")}, {Name: ghi.String("service/google-z")}, {Name: ghi.String("service/google-a")}},
			},
			teamMembers:             map[string][]string{"google-x": {"googler_team_member"}, "google-y": {"googler_y_team_member"}},
			expectSpecificReviewers: []string{},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			requestedReviewers := []*ghi.User{}
			for _, login := range tc.requestedReviewers {
				requestedReviewers = append(requestedReviewers, &ghi.User{Login: ghi.String(login)})
			}

			previousReviewers := []*ghi.User{}
			for _, login := range tc.previousReviewers {
				previousReviewers = append(previousReviewers, &ghi.User{Login: ghi.String(login)})
			}

			teamMembers := map[string][]*ghi.User{}
			for team, logins := range tc.teamMembers {
				teamMembers[team] = []*ghi.User{}
				for _, login := range logins {
					teamMembers[team] = append(teamMembers[team], &ghi.User{Login: ghi.String(login)})
				}
			}

			mockGH := &mockGithub{
				pullRequest:        tc.pullRequest,
				requestedReviewers: requestedReviewers,
				previousReviewers:  previousReviewers,
				teamMembers:        teamMembers,
				calledMethods:      make(map[string][][]any),
			}

			execRequestServiceReviewers("1", mockGH, enrolledTeamsYaml)

			actualReviewers := []string{}
			for _, args := range mockGH.calledMethods["RequestPullRequestReviewers"] {
				actualReviewers = append(actualReviewers, args[1].([]string)...)
			}

			if tc.expectSpecificReviewers != nil {
				assert.ElementsMatch(t, tc.expectSpecificReviewers, actualReviewers)
			}
		})
	}
}
