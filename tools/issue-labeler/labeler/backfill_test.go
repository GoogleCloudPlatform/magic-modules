package labeler

import (
	"fmt"

	"reflect"
	"regexp"
	"strings"
	"testing"
)

func testIssueBodyWithResources(resources []string) string {
	return fmt.Sprintf(`
### New or Affected Resource(s):

%s

#
`, strings.Join(resources, "\n"))
}

func TestComputeIssueUpdates(t *testing.T) {
	defaultRegexpLabels := []RegexpLabel{
		{
			Regexp: regexp.MustCompile("google_service1_.*"),
			Label:  "service/service1",
		},
		{
			Regexp: regexp.MustCompile("google_service2_resource1"),
			Label:  "service/service2-subteam1",
		},
		{
			Regexp: regexp.MustCompile("google_service2_resource2"),
			Label:  "service/service2-subteam2",
		},
	}
	cases := map[string]struct {
		issues               []Issue
		regexpLabels         []RegexpLabel
		expectedIssueUpdates []IssueUpdate
	}{
		"no issues -> no updates": {
			issues:               []Issue{},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		"exempt labels -> no updates": {
			issues: []Issue{
				{
					Number:      1,
					Body:        testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels:      []Label{{Name: "service/terraform"}},
					PullRequest: map[string]any{},
				},
				{
					Number:      2,
					Body:        testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels:      []Label{{Name: "forward/exempt"}},
					PullRequest: map[string]any{},
				},
			},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		"add resource & review labels": {
			issues: []Issue{
				{
					Number:      1,
					Body:        testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels:      []Label{},
					PullRequest: map[string]any{},
				},
				{
					Number:      2,
					Body:        testIssueBodyWithResources([]string{"google_service2_resource1"}),
					Labels:      []Label{},
					PullRequest: map[string]any{},
				},
			},
			regexpLabels: defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{
				{
					Number: 1,
					Labels: []string{"forward/review", "service/service1"},
				},
				{
					Number: 2,
					Labels: []string{"forward/review", "service/service2-subteam1"},
				},
			},
		},
		"don't update issues if all service labels are already present": {
			issues: []Issue{
				{
					Number:      1,
					Body:        testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels:      []Label{{Name: "service/service1"}},
					PullRequest: map[string]any{},
				},
				{
					Number:      2,
					Body:        testIssueBodyWithResources([]string{"google_service2_resource1"}),
					Labels:      []Label{{Name: "service/service2-subteam1"}},
					PullRequest: map[string]any{},
				},
			},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		"add missing service labels": {
			issues: []Issue{
				{
					Number:      1,
					Body:        testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels:      []Label{{Name: "service/service2-subteam1"}},
					PullRequest: map[string]any{},
				},
				{
					Number:      2,
					Body:        testIssueBodyWithResources([]string{"google_service2_resource2"}),
					Labels:      []Label{{Name: "service/service1"}},
					PullRequest: map[string]any{},
				},
			},
			regexpLabels: defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{
				{
					Number:    1,
					Labels:    []string{"forward/review", "service/service1", "service/service2-subteam1"},
					OldLabels: []string{"service/service2-subteam1"},
				},
				{
					Number:    2,
					Labels:    []string{"forward/review", "service/service1", "service/service2-subteam2"},
					OldLabels: []string{"service/service1"},
				},
			},
		},
		"don't add missing service labels if already linked": {
			issues: []Issue{
				{
					Number:      1,
					Body:        testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels:      []Label{{Name: "service/service2-subteam1"}, {Name: "forward/linked"}},
					PullRequest: map[string]any{},
				},
			},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			issueUpdates := ComputeIssueUpdates(tc.issues, tc.regexpLabels)
			// reflect.DeepEqual treats nil & empty slices as not equal so ignore diffs if both slices are empty.
			if (len(issueUpdates) > 0 || len(tc.expectedIssueUpdates) > 0) && !reflect.DeepEqual(issueUpdates, tc.expectedIssueUpdates) {
				t.Errorf("Expected %v, got %v", tc.expectedIssueUpdates, issueUpdates)
			}
		})
	}
}
