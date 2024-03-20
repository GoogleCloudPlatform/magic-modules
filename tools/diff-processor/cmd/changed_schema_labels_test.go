package cmd

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var enrolledTeamsYaml = []byte(`
service/google-x:
  resources:
  - google_x_resource
service/google-z:
  resources:
  - google_z_resource`)

func TestChangedSchemaLabelsCmdRun(t *testing.T) {
	cases := map[string]struct {
		args           []string
		oldResourceMap map[string]*schema.Resource
		newResourceMap map[string]*schema.Resource
		expectedLabels []string
		expectError    bool
	}{
		"empty resource map": {
			args:           []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{},
			newResourceMap: map[string]*schema.Resource{},
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
			expectedLabels: []string{"service/google-x"},
		},
		"resources changed that match multiple mappings": {
			args: []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
				"google_z_resource": {
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
				"google_z_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Required: true},
					},
				},
			},
			expectedLabels: []string{"service/google-x", "service/google-z"},
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			o := changedSchemaLabelsOptions{
				computeSchemaDiff: func() diff.SchemaDiff {
					return diff.ComputeSchemaDiff(tc.oldResourceMap, tc.newResourceMap)
				},
				enrolledTeamsYaml: enrolledTeamsYaml,
				stdout:            &buf,
			}

			err := o.run()
			if err != nil {
				if tc.expectError {
					return
				}
				t.Errorf("Error running command: %s", err)
			}

			out := make([]byte, buf.Len())
			buf.Read(out)
			var gotLabels []string
			if err = json.Unmarshal(out, &gotLabels); err != nil {
				t.Fatalf("Unable to unmarshal labels (%q): %s", out, err)
			}

			less := func(a, b string) bool { return a < b }
			if (len(tc.expectedLabels) > 0 || len(gotLabels) > 0) && !cmp.Equal(tc.expectedLabels, gotLabels, cmpopts.SortSlices(less)) {
				t.Errorf("Unexpected final labels. Want %q, got %q", tc.expectedLabels, gotLabels)
			}
		})
	}
}
