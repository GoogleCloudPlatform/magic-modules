package cmd

import (
	_ "embed"
	"errors"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/labeler"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

var enrolledTeamsYaml = []byte(`
service/google-x:
  resources:
  - google_x_resource`)

func TestAddLabelsCmdRun(t *testing.T) {
	cases := map[string]struct {
		args           []string
		oldResourceMap map[string]*schema.Resource
		newResourceMap map[string]*schema.Resource
		githubIssue    *labeler.Issue
		updateErrors   bool
		expectedLabels []string
		expectError    bool
	}{
		"empty resource map": {
			args: []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{},
			newResourceMap: map[string]*schema.Resource{},
			githubIssue: &labeler.Issue{
				Number:      12345,
				Body:        "Unused",
				Labels:      []labeler.Label{},
				PullRequest: map[string]any{},
			},
			expectedLabels: nil,
		},
		"resource changed that doesn't match mapping": {
			args: []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{
				"google_y_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_y_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Required: true},
					},
				},
			},
			githubIssue: &labeler.Issue{
				Number:      12345,
				Body:        "Unused",
				Labels:      []labeler.Label{},
				PullRequest: map[string]any{},
			},
			expectedLabels: nil,
		},
		"resource matches mapping but isn't changed": {
			args: []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			githubIssue: &labeler.Issue{
				Number:      12345,
				Body:        "Unused",
				Labels:      []labeler.Label{},
				PullRequest: map[string]any{},
			},
			expectedLabels: nil,
		},
		"resource changed that matches mapping": {
			args: []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Required: true},
					},
				},
			},
			githubIssue: &labeler.Issue{
				Number:      12345,
				Body:        "Unused",
				Labels:      []labeler.Label{},
				PullRequest: map[string]any{},
			},
			expectedLabels: []string{"service/google-x"},
		},
		"service labels are deduped": {
			args: []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
				"google_x_resource2": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Required: true},
					},
				},
				"google_x_resource2": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Required: true},
					},
				},
			},
			githubIssue: &labeler.Issue{
				Number:      12345,
				Body:        "Unused",
				Labels:      []labeler.Label{},
				PullRequest: map[string]any{},
			},
			expectedLabels: []string{"service/google-x"},
		},
		"existing labels are preserved": {
			args: []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Required: true},
					},
				},
			},
			githubIssue: &labeler.Issue{
				Number:      12345,
				Body:        "Unused",
				Labels:      []labeler.Label{{Name: "override-breaking-change"}},
				PullRequest: map[string]any{},
			},
			expectedLabels: []string{"override-breaking-change", "service/google-x"},
		},
		"existing service label prevents new service labels": {
			args: []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Required: true},
					},
				},
			},
			githubIssue: &labeler.Issue{
				Number:      12345,
				Body:        "Unused",
				Labels:      []labeler.Label{{Name: "service/google-z"}},
				PullRequest: map[string]any{},
			},
			// nil indicates that the issue won't be updated at all (preserving existing labels)
			expectedLabels: nil,
		},
		"error fetching issue": {
			args: []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Required: true},
					},
				},
			},
			githubIssue: nil,
			expectError: true,
		},
		"error parsing PR id": {
			args: []string{"foobar"},
			oldResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Required: true},
					},
				},
			},
			githubIssue: &labeler.Issue{
				Number:      12345,
				Body:        "Unused",
				Labels:      []labeler.Label{{Name: "service/google-z"}},
				PullRequest: map[string]any{},
			},
			expectError: true,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			var gotLabels []string
			o := addLabelsOptions{
				computeSchemaDiff: func() diff.SchemaDiff {
					return diff.ComputeSchemaDiff(tc.oldResourceMap, tc.newResourceMap)
				},
				enrolledTeamsYaml: enrolledTeamsYaml,
				getIssue: func(repository string, id uint64) (labeler.Issue, error) {
					if tc.githubIssue != nil {
						return *tc.githubIssue, nil
					}
					var issue labeler.Issue
					return issue, errors.New("Error getting issue")
				},
				updateIssues: func(repository string, issueUpdates []labeler.IssueUpdate, dryRun bool) {
					gotLabels = issueUpdates[0].Labels
				},
			}

			err := o.run([]string{"1"})
			if err != nil {
				if tc.expectError {
					return
				}
				t.Errorf("Error running command: %s", err)
			}

			if tc.expectedLabels == nil {
				if gotLabels != nil {
					t.Errorf("Expected updateIssues to not run. Got %v as new labels", gotLabels)
				}
			}

			less := func(a, b string) bool { return a < b }
			if (len(tc.expectedLabels) > 0 || len(gotLabels) > 0) && !cmp.Equal(tc.expectedLabels, gotLabels, cmpopts.SortSlices(less)) {
				t.Errorf("Unexpected final labels. Want %v, got %v", tc.expectedLabels, gotLabels)
			}
		})
	}
}
