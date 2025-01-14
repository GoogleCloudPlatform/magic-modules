package labeler

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-github/v61/github"
)

func testIssueBodyWithResources(resources []string) *string {
	return github.String(fmt.Sprintf(`
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

	cases := map[string]struct {
		issues               []*github.Issue
		regexpLabels         []RegexpLabel
		expectedIssueUpdates []IssueUpdate
	}{
		"no issues -> no updates": {
			issues:               []*github.Issue{},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		"exempt labels -> no updates": {
			issues: []*github.Issue{
				{
					Number: github.Int(1),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels: []*github.Label{{Name: github.String("service/terraform")}},
				},
				{
					Number: github.Int(2),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels: []*github.Label{{Name: github.String("forward/exempt")}},
				},
			},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		"add resource & review labels": {
			issues: []*github.Issue{
				{
					Number: github.Int(1),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
				},
				{
					Number: github.Int(2),
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
		"don't update issues if all service labels are already present": {
			issues: []*github.Issue{
				{
					Number: github.Int(1),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels: []*github.Label{{Name: github.String("service/service1")}},
				},
				{
					Number: github.Int(2),
					Body:   testIssueBodyWithResources([]string{"google_service2_resource1"}),
					Labels: []*github.Label{{Name: github.String("service/service2-subteam1")}},
				},
			},
			regexpLabels:         defaultRegexpLabels,
			expectedIssueUpdates: []IssueUpdate{},
		},
		"add missing service labels": {
			issues: []*github.Issue{
				{
					Number: github.Int(1),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels: []*github.Label{{Name: github.String("service/service2-subteam1")}},
				},
				{
					Number: github.Int(2),
					Body:   testIssueBodyWithResources([]string{"google_service2_resource2"}),
					Labels: []*github.Label{{Name: github.String("service/service1")}},
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
			issues: []*github.Issue{
				{
					Number: github.Int(1),
					Body:   testIssueBodyWithResources([]string{"google_service1_resource1"}),
					Labels: []*github.Label{{Name: github.String("service/service2-subteam1")}, {Name: github.String("forward/linked")}},
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
			if !issueUpdatesEqual(issueUpdates, tc.expectedIssueUpdates) {
				t.Errorf("Expected %v, got %v", tc.expectedIssueUpdates, issueUpdates)
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
