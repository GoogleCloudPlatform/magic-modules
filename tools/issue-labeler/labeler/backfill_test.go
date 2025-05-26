package labeler

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-github/v68/github"
)

func testIssueBodyWithResources(resources []string) *string {
	return github.Ptr(fmt.Sprintf(`
### New or Affected Resource(s):

%s

#
`, strings.Join(resources, "\n")))
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

	cases := []struct {
		name, description    string
		issues               []*github.Issue
		regexpLabels         []RegexpLabel
		expectedIssueUpdates []IssueUpdate
	}{
		{
			name:                 "no issues",
			description:          "no issues means no updates",
			issues:               []*github.Issue{},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		{
			name:        "nil body",
			description: "gracefully handle a nil issue body",
			issues: []*github.Issue{
				{
					Number: github.Ptr(1),
				},
			},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		{
			name:        "nil number",
			description: "gracefully handle a nil issue number",
			issues: []*github.Issue{
				{
					Body: testIssueBodyWithResources([]string{"google_service1_resource1"}),
				},
			},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		{
			name: "no listed resources",
			issues: []*github.Issue{
				{
					Number: github.Ptr(1),
					Body:   github.Ptr("Body with unusual structure"),
				},
			},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		{
			name:        "service/terraform",
			description: "issues with service/terraform shouldn't get new labels",
			issues: []*github.Issue{
				{
					Number: github.Ptr(1),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels: []*github.Label{{Name: github.Ptr("service/terraform")}},
				},
				{
					Number: github.Ptr(2),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels: []*github.Label{{Name: github.Ptr("forward/exempt")}},
				},
			},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		{
			name:        "add resource & review labels",
			description: "issues with affected resources should normally get new labels added",
			issues: []*github.Issue{
				{
					Number: github.Ptr(1),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
				},
				{
					Number: github.Ptr(2),
					Body:   testIssueBodyWithResources([]string{"google_service2_resource1"}),
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
		{
			name:        "labels already correct",
			description: "don't update issues if all expected service labels are already present",
			issues: []*github.Issue{
				{
					Number: github.Ptr(1),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels: []*github.Label{{Name: github.Ptr("service/service1")}},
				},
				{
					Number: github.Ptr(2),
					Body:   testIssueBodyWithResources([]string{"google_service2_resource1"}),
					Labels: []*github.Label{{Name: github.Ptr("service/service2-subteam1")}},
				},
			},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		{
			name:        "missing labels",
			description: "add missing service labels",
			issues: []*github.Issue{
				{
					Number: github.Ptr(1),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels: []*github.Label{{Name: github.Ptr("service/service2-subteam1")}},
				},
				{
					Number: github.Ptr(2),
					Body:   testIssueBodyWithResources([]string{"google_service2_resource2"}),
					Labels: []*github.Label{{Name: github.Ptr("service/service1")}},
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
		{
			name:        "forward/linked",
			description: "don't add missing service labels if already linked",
			issues: []*github.Issue{
				{
					Number: github.Ptr(1),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels: []*github.Label{{Name: github.Ptr("service/service2-subteam1")}, {Name: github.Ptr("forward/linked")}},
				},
			},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		{
			name:        "test failure",
			description: "add service labels if missed but don't add forward/review label for test failure ticket",
			issues: []*github.Issue{
				{
					Number: github.Ptr(1),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels: []*github.Label{{Name: github.Ptr("test-failure")}, {Name: github.Ptr("test-failure-100")}},
				},
				{
					Number: github.Ptr(2),
					Body:   testIssueBodyWithResources([]string{"google_service2_resource1"}),
					Labels: []*github.Label{{Name: github.Ptr("test-failure")}, {Name: github.Ptr("test-failure-50")}, {Name: github.Ptr("service/service2-subteam1")}},
				},
			},
			regexpLabels: defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{
				{
					Number:    1,
					Labels:    []string{"service/service1", "test-failure", "test-failure-100"},
					OldLabels: []string{"test-failure", "test-failure-100"},
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			issueUpdates := ComputeIssueUpdates(tc.issues, tc.regexpLabels)
			if !issueUpdatesEqual(issueUpdates, tc.expectedIssueUpdates) {
				t.Errorf("ComputeIssueUpdates(%s) expected %v, got %v", tc.name, tc.expectedIssueUpdates, issueUpdates)
			}
		})
	}
}

func TestSplitRepository(t *testing.T) {
	tests := []struct {
		name       string
		repository string
		wantOwner  string
		wantRepo   string
		wantErr    bool
	}{
		{
			name:       "valid repository",
			repository: "owner/repo",
			wantOwner:  "owner",
			wantRepo:   "repo",
			wantErr:    false,
		},
		{
			name:       "invalid repository",
			repository: "invalid-format",
			wantOwner:  "",
			wantRepo:   "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := splitRepository(tt.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if owner != tt.wantOwner {
				t.Errorf("splitRepository() owner = %v, want %v", owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("splitRepository() repo = %v, want %v", repo, tt.wantRepo)
			}
		})
	}
}

// Helper function to compare issue updates while handling nil/empty slice equality
func issueUpdatesEqual(a, b []IssueUpdate) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return reflect.DeepEqual(a, b)
}
