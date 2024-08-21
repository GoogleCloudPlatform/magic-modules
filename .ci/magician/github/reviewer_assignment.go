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
func ChooseCoreReviewers(requestedReviewers, previousReviewers []User) (reviewersToRequest []string, newPrimaryReviewer string) {
	hasPrimaryReviewer := false
	newPrimaryReviewer = ""

	for _, reviewer := range requestedReviewers {
		if IsCoreReviewer(reviewer.Login) {
			hasPrimaryReviewer = true
			break
		}
	}

	for _, reviewer := range previousReviewers {
		if IsCoreReviewer(reviewer.Login) {
			hasPrimaryReviewer = true
			reviewersToRequest = append(reviewersToRequest, reviewer.Login)
		}
	}

	if !hasPrimaryReviewer {
		newPrimaryReviewer = GetRandomReviewer()
		reviewersToRequest = append(reviewersToRequest, newPrimaryReviewer)
	}

	return reviewersToRequest, newPrimaryReviewer
}

func FormatReviewerComment(newPrimaryReviewer string) string {
	tmpl, err := template.New("REVIEWER_ASSIGNMENT_COMMENT.md").Parse(reviewerAssignmentComment)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse REVIEWER_ASSIGNMENT_COMMENT.md: %s", err))
	}
	sb := new(strings.Builder)
	tmpl.Execute(sb, map[string]any{
		"reviewer": newPrimaryReviewer,
	})
	return sb.String()
}
