package github

import (
	"fmt"
	"strings"
	"text/template"

	_ "embed"
)

var (
	//go:embed REVIEWER_ASSIGNMENT_COMMENT.md
	reviewerAssignmentComment string
)

// Returns a list of users to request review from, as well as a new primary reviewer if this is the first run.
func ChooseReviewers(firstRequestedReviewer string, previouslyInvolvedReviewers []string) (reviewersToRequest []string, newPrimaryReviewer string) {
	hasPrimaryReviewer := false
	newPrimaryReviewer = ""

	if firstRequestedReviewer != "" {
		hasPrimaryReviewer = true
	}

	for _, reviewer := range previouslyInvolvedReviewers {
		if IsTeamReviewer(reviewer) {
			hasPrimaryReviewer = true
			reviewersToRequest = append(reviewersToRequest, reviewer)
		}
	}

	if !hasPrimaryReviewer {
		newPrimaryReviewer = GetRandomReviewer()
		reviewersToRequest = append(reviewersToRequest, newPrimaryReviewer)
	}

	return reviewersToRequest, newPrimaryReviewer
}

func FormatReviewerComment(newPrimaryReviewer string, authorUserType UserType, trusted bool) string {
	tmpl, err := template.New("REVIEWER_ASSIGNMENT_COMMENT.md").Parse(reviewerAssignmentComment)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse REVIEWER_ASSIGNMENT_COMMENT.md: %s", err))
	}
	sb := new(strings.Builder)
	tmpl.Execute(sb, map[string]any{
		"reviewer":       newPrimaryReviewer,
		"authorUserType": authorUserType.String(),
		"trusted":        trusted,
	})
	return sb.String()
}
