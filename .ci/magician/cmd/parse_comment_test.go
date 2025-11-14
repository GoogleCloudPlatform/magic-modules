package cmd

import (
	"magician/github"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseReassignReviewerCommand tests various formats of the reassign-reviewer command
func TestParseReassignReviewerCommand(t *testing.T) {
	tests := []struct {
		name         string
		commandLine  string
		wantMatch    bool
		wantReviewer string
	}{
		{
			name:         "reassign-reviewer with hyphen and username",
			commandLine:  "reassign-reviewer user123",
			wantMatch:    true,
			wantReviewer: "user123",
		},
		{
			name:         "reassign reviewer with space and username",
			commandLine:  "reassign reviewer user456",
			wantMatch:    true,
			wantReviewer: "user456",
		},
		{
			name:         "assign-reviewer with hyphen",
			commandLine:  "assign-reviewer newuser",
			wantMatch:    true,
			wantReviewer: "newuser",
		},
		{
			name:         "assign reviewer with space",
			commandLine:  "assign reviewer someone",
			wantMatch:    true,
			wantReviewer: "someone",
		},
		{
			name:         "reassign-review without er suffix",
			commandLine:  "reassign-review john-doe",
			wantMatch:    true,
			wantReviewer: "john-doe",
		},
		{
			name:         "assign review without er suffix",
			commandLine:  "assign review jane_doe",
			wantMatch:    true,
			wantReviewer: "jane_doe",
		},
		{
			name:         "with @ prefix on username",
			commandLine:  "reassign-reviewer @github-user",
			wantMatch:    true,
			wantReviewer: "github-user", // @ is stripped by regex
		},
		{
			name:         "no username specified",
			commandLine:  "reassign-reviewer",
			wantMatch:    true,
			wantReviewer: "",
		},
		{
			name:         "no space between command and username",
			commandLine:  "reassign-revieweruser123",
			wantMatch:    true,
			wantReviewer: "user123",
		},
		{
			name:         "no space between command and @username",
			commandLine:  "reassign-reviewer@user123",
			wantMatch:    true,
			wantReviewer: "user123",
		},
		{
			name:         "extra spaces before username",
			commandLine:  "reassign-reviewer   user789",
			wantMatch:    true,
			wantReviewer: "user789",
		},
		{
			name:         "assignreviewer without any separator",
			commandLine:  "assignreviewer user999",
			wantMatch:    true,
			wantReviewer: "user999",
		},
		{
			name:         "assignrevieweruser999 all together",
			commandLine:  "assignrevieweruser999",
			wantMatch:    true,
			wantReviewer: "user999",
		},
		{
			name:         "assignreviewer@user999 all together with @",
			commandLine:  "assignreviewer@user999",
			wantMatch:    true,
			wantReviewer: "user999",
		},
		{
			name:         "username with underscore and hyphen",
			commandLine:  "reassign-reviewer test_user-123",
			wantMatch:    true,
			wantReviewer: "test_user-123",
		},
		{
			name:         "invalid characters in username (space)",
			commandLine:  "reassign-reviewer user name",
			wantMatch:    true,
			wantReviewer: "user", // Only captures up to the space
		},
		{
			name:         "invalid characters in username (dot)",
			commandLine:  "reassign-reviewer user.name",
			wantMatch:    true,
			wantReviewer: "user", // Only captures up to the dot
		},
		{
			name:         "unrelated command",
			commandLine:  "cherry-pick branch-name",
			wantMatch:    false,
			wantReviewer: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := reassignReviewerRegex.FindStringSubmatch(tt.commandLine)
			gotMatch := matches != nil

			if gotMatch != tt.wantMatch {
				t.Errorf("reassignReviewerRegex.FindStringSubmatch(%q) match = %v, want %v",
					tt.commandLine, gotMatch, tt.wantMatch)
				return
			}

			if gotMatch {
				gotReviewer := ""
				if len(matches) > 1 {
					gotReviewer = matches[1]
					// No cleanup needed - regex handles everything
				}

				if gotReviewer != tt.wantReviewer {
					t.Errorf("extracted reviewer = %q, want %q", gotReviewer, tt.wantReviewer)
				}
			}
		})
	}
}

// TestMagicianInvocationRegex tests the extraction of command line from comments
func TestMagicianInvocationRegex(t *testing.T) {
	tests := []struct {
		name        string
		comment     string
		wantFound   bool
		wantCommand string
	}{
		{
			name:        "simple command",
			comment:     "@modular-magician reassign-reviewer user1",
			wantFound:   true,
			wantCommand: "reassign-reviewer user1",
		},
		{
			name:        "command in middle of comment",
			comment:     "LGTM!\n@modular-magician assign reviewer @john\nGreat work!",
			wantFound:   true,
			wantCommand: "assign reviewer @john",
		},
		{
			name:        "multiple commands (only first is processed)",
			comment:     "@modular-magician reassign-reviewer alice\n@modular-magician assign-reviewer bob",
			wantFound:   true,
			wantCommand: "reassign-reviewer alice",
		},
		{
			name:        "no command",
			comment:     "Just a regular comment without any magician invocation",
			wantFound:   false,
			wantCommand: "",
		},
		{
			name:        "command with extra spaces",
			comment:     "@modular-magician    reassign  reviewer   user2",
			wantFound:   true,
			wantCommand: "reassign  reviewer   user2",
		},
		{
			name:        "command with text after",
			comment:     "@modular-magician reassign-reviewer newuser please",
			wantFound:   true,
			wantCommand: "reassign-reviewer newuser please",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := magicianInvocationRegex.FindStringSubmatch(tt.comment)

			gotFound := match != nil
			if gotFound != tt.wantFound {
				t.Errorf("found command = %v, want %v", gotFound, tt.wantFound)
				return
			}

			if gotFound {
				if len(match) < 2 {
					t.Errorf("match has insufficient captures")
					return
				}

				gotCommand := match[1]
				// Note: In the actual implementation, this gets trimmed
				gotCommand = strings.TrimSpace(gotCommand)
				wantCommand := strings.TrimSpace(tt.wantCommand)

				if gotCommand != wantCommand {
					t.Errorf("command = %q, want %q", gotCommand, wantCommand)
				}
			}
		})
	}
}

func TestExecParseComment(t *testing.T) {
	availableReviewers := github.AvailableReviewers(nil)
	if len(availableReviewers) < 3 {
		t.Fatalf("not enough available reviewers (%v) to run TestExecParseComment (need at least 3)", availableReviewers)
	}

	cases := map[string]struct {
		comment                 string
		existingComments        []github.PullRequestComment
		expectSpecificReviewers []string
		expectRemovedReviewers  []string
		expectNoAction          bool
		expectCommentUpdate     bool
		expectCommentCreate     bool
	}{
		// "reassign-reviewer with hyphen and specific user": {
		// 	comment: "LGTM! @modular-magician reassign-reviewer alice",
		// 	existingComments: []github.PullRequestComment{
		// 		{
		// 			Body: github.FormatReviewerComment("bob"),
		// 			ID:   1234,
		// 		},
		// 	},
		// 	expectSpecificReviewers: []string{"alice"},
		// 	expectRemovedReviewers:  []string{"bob"},
		// 	expectCommentUpdate:     true,
		// },
		// "reassign to random reviewer (no username specified)": {
		// 	comment: "@modular-magician reassign-reviewer",
		// 	existingComments: []github.PullRequestComment{
		// 		{
		// 			Body: github.FormatReviewerComment("george"),
		// 			ID:   3456,
		// 		},
		// 	},
		// 	expectRemovedReviewers: []string{"george"},
		// 	expectCommentUpdate:    true,
		// 	// Can't check specific reviewer since it's random
		// },
		// "multiple @modular-magician invocations (only first processed)": {
		// 	comment: "@modular-magician reassign-reviewer larry\n@modular-magician reassign-reviewer mary",
		// 	existingComments: []github.PullRequestComment{
		// 		{
		// 			Body: github.FormatReviewerComment("nancy"),
		// 			ID:   1111,
		// 		},
		// 	},
		// 	expectSpecificReviewers: []string{"larry"}, // Only larry, not mary
		// 	expectRemovedReviewers:  []string{"nancy"},
		// 	expectCommentUpdate:     true,
		// },
		"no @modular-magician invocation": {
			comment:        "Just a regular comment without magician",
			expectNoAction: true,
		},
		"@modular-magician with no command": {
			comment:        "@modular-magician",
			expectNoAction: true,
		},
		"@modular-magician with unrecognized command": {
			comment:        "@modular-magician cherry-pick branch-xyz",
			expectNoAction: true,
		},
		// 		"command in middle of multi-line comment": {
		// 			comment: `This looks good to me.
		// LGTM!

		// @modular-magician reassign-reviewer rachel

		// Thanks for the great work!`,
		// 	existingComments: []github.PullRequestComment{
		// 		{
		// 			Body: github.FormatReviewerComment("steve"),
		// 			ID:   3333,
		// 		},
		// 	},
		// 	expectSpecificReviewers: []string{"rachel"},
		// 	expectRemovedReviewers:  []string{"steve"},
		// 	expectCommentUpdate:     true,
		// },
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			gh := &mockGithub{
				pullRequest: github.PullRequest{
					User: github.User{Login: "author"},
				},
				calledMethods:       make(map[string][][]any),
				pullRequestComments: tc.existingComments,
			}

			err := execParseComment("1", tc.comment, gh)
			if err != nil {
				t.Fatalf("execParseComment failed: %v", err)
			}

			// Check if no action was expected
			if tc.expectNoAction {
				assert.Empty(t, gh.calledMethods["RequestPullRequestReviewers"], "Expected no RequestPullRequestReviewers calls")
				assert.Empty(t, gh.calledMethods["PostComment"], "Expected no PostComment calls")
				assert.Empty(t, gh.calledMethods["UpdateComment"], "Expected no UpdateComment calls")
				return
			}

			// Check reviewer assignment
			if !tc.expectNoAction {
				assert.Len(t, gh.calledMethods["RequestPullRequestReviewers"], 1, "Expected RequestPullRequestReviewers called exactly once")
			}

			var assignedReviewers []string
			for _, args := range gh.calledMethods["RequestPullRequestReviewers"] {
				assignedReviewers = append(assignedReviewers, args[1].([]string)...)
			}

			var removedReviewers []string
			for _, args := range gh.calledMethods["RemovePullRequestReviewers"] {
				removedReviewers = append(removedReviewers, args[1].([]string)...)
			}

			// Check specific reviewers if expected
			if tc.expectSpecificReviewers != nil {
				for _, reviewer := range assignedReviewers {
					assert.Contains(t, tc.expectSpecificReviewers, reviewer)
				}
			}

			// Check removed reviewers if expected
			if tc.expectRemovedReviewers != nil {
				for _, reviewer := range removedReviewers {
					assert.Contains(t, tc.expectRemovedReviewers, reviewer)
				}
			}

			// Check comment creation/update
			if tc.expectCommentCreate {
				assert.Len(t, gh.calledMethods["PostComment"], 1, "Expected PostComment called once")
				assert.Empty(t, gh.calledMethods["UpdateComment"], "Expected no UpdateComment calls")
			}

			if tc.expectCommentUpdate {
				assert.Len(t, gh.calledMethods["UpdateComment"], 1, "Expected UpdateComment called once")
				assert.Empty(t, gh.calledMethods["PostComment"], "Expected no PostComment calls")
			}
		})
	}
}
