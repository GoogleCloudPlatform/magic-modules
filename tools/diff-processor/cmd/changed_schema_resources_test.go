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

func TestChangedSchemaResourcesCmdRun(t *testing.T) {
	cases := map[string]struct {
		args              []string
		oldResourceMap    map[string]*schema.Resource
		newResourceMap    map[string]*schema.Resource
		expectedResources []string
		expectError       bool
	}{
		"empty resource map": {
			args:              []string{"12345"},
			oldResourceMap:    map[string]*schema.Resource{},
			newResourceMap:    map[string]*schema.Resource{},
			expectedResources: nil,
		},
		"resource isn't changed": {
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
			expectedResources: nil,
		},
		"resource is changed": {
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
			expectedResources: []string{"google_x_resource"},
		},
		"multiple resources are changed": {
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
			expectedResources: []string{"google_x_resource", "google_z_resource"},
		},
		"multiple resources but not all are changed": {
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
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			expectedResources: []string{"google_x_resource"},
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			o := changedSchemaResourcesOptions{
				computeSchemaDiff: func() diff.SchemaDiff {
					return diff.ComputeSchemaDiff(tc.oldResourceMap, tc.newResourceMap)
				},
				stdout: &buf,
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
			var gotResources []string
			if err = json.Unmarshal(out, &gotResources); err != nil {
				t.Fatalf("Unable to unmarshal labels (%q): %s", out, err)
			}

			less := func(a, b string) bool { return a < b }
			if (len(tc.expectedResources) > 0 || len(gotResources) > 0) && !cmp.Equal(tc.expectedResources, gotResources, cmpopts.SortSlices(less)) {
				t.Errorf("Unexpected final labels. Want %q, got %q", tc.expectedResources, gotResources)
			}
		})
	}
}
