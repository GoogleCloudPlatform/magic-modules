package cmd

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSchemaDiffCmdRun(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		oldResourceMap map[string]*schema.Resource
		newResourceMap map[string]*schema.Resource
		want           simpleSchemaDiff
	}{
		{
			name:           "empty resource map",
			args:           []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{},
			newResourceMap: map[string]*schema.Resource{},
			want:           simpleSchemaDiff{},
		},
		{
			name:           "resource is added",
			args:           []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{},
			newResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			want: simpleSchemaDiff{
				AddedResources: []string{"google_x_resource"},
			},
		},
		{
			name:           "multiple resources are added",
			args:           []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{},
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
			want: simpleSchemaDiff{
				AddedResources: []string{"google_x_resource", "google_z_resource"},
			},
		},
		{
			name: "resource is removed",
			args: []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{
				"google_x_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{},
			want: simpleSchemaDiff{
				RemovedResources: []string{"google_x_resource"},
			},
		},
		{
			name: "multiple resources are removed",
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
			newResourceMap: map[string]*schema.Resource{},
			want: simpleSchemaDiff{
				RemovedResources: []string{"google_x_resource", "google_z_resource"},
			},
		},
		{
			name: "resource isn't changed",
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
			want: simpleSchemaDiff{},
		},
		{
			name: "resource is changed",
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
			want: simpleSchemaDiff{
				ModifiedResources: []string{"google_x_resource"},
			},
		},
		{
			name: "multiple resources are changed",
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
			want: simpleSchemaDiff{
				ModifiedResources: []string{"google_x_resource", "google_z_resource"},
			},
		},
		{
			name: "multiple resources but not all are changed",
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
			want: simpleSchemaDiff{
				ModifiedResources: []string{"google_x_resource"},
			},
		},
		{
			name: "multiple resources are added, changed, or removed",
			args: []string{"12345"},
			oldResourceMap: map[string]*schema.Resource{
				"google_y_resource": {
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
				"google_y_resource": {
					Schema: map[string]*schema.Schema{
						"field_a": {Description: "beep", Optional: true},
						"field_b": {Description: "beep", Optional: true},
						"field_c": {Description: "beep", Optional: true},
					},
				},
			},
			want: simpleSchemaDiff{
				AddedResources:    []string{"google_x_resource"},
				ModifiedResources: []string{"google_y_resource"},
				RemovedResources:  []string{"google_z_resource"},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			o := schemaDiffOptions{
				computeSchemaDiff: func() diff.SchemaDiff {
					return diff.ComputeSchemaDiff(tc.oldResourceMap, tc.newResourceMap)
				},
				stdout: &buf,
			}

			err := o.run()
			if err != nil {
				t.Errorf("Error running command: %s", err)
			}

			out := make([]byte, buf.Len())
			buf.Read(out)
			var got simpleSchemaDiff
			if err = json.Unmarshal(out, &got); err != nil {
				t.Fatalf("Unable to unmarshal simple diff (%q): %s", out, err)
			}

			if !cmp.Equal(tc.want, got) {
				t.Errorf("Unexpected simple diff. Want %q, got %q", tc.want, got)
			}
		})
	}
}
