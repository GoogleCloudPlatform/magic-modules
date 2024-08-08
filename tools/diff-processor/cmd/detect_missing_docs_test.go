package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDetectMissingDocs(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd = %v", err)
	}
	cases := []struct {
		name           string
		oldResourceMap map[string]*schema.Resource
		newResourceMap map[string]*schema.Resource
		repoPath       string
		want           []MissingDocsResource
	}{
		{
			name: "no new fields",
			oldResourceMap: map[string]*schema.Resource{
				"google_x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Optional: true},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Optional: true},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			want: []MissingDocsResource{},
		},
		{
			name: "one new resource",
			newResourceMap: map[string]*schema.Resource{
				"google_x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep"},
						"field-b": {Description: "beep"},
					},
				},
			},
			oldResourceMap: map[string]*schema.Resource{},
			want: []MissingDocsResource{
				{
					Resource: "google_x",
					Fields:   []string{"field-a", "field-b"},
				},
			},
		},
		{
			name: "one new field",
			newResourceMap: map[string]*schema.Resource{
				"google_x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Optional: true},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			oldResourceMap: map[string]*schema.Resource{
				"google_x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Required: true},
					},
				},
			},
			want: []MissingDocsResource{
				{
					Resource: "google_x",
					Fields:   []string{"field-b"},
				},
			},
		},
		{
			name: "one new resource with doc",
			newResourceMap: map[string]*schema.Resource{
				"google_x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep"},
						"field-b": {Description: "beep"},
					},
				},
			},
			oldResourceMap: map[string]*schema.Resource{},
			repoPath:       filepath.Join(cwd, "testdata"),
			want:           []MissingDocsResource{},
		},
		{
			name: "one new field with doc",
			newResourceMap: map[string]*schema.Resource{
				"google_x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Optional: true},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			oldResourceMap: map[string]*schema.Resource{
				"google_x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Required: true},
					},
				},
			},
			repoPath: filepath.Join(cwd, "testdata"),
			want:     []MissingDocsResource{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			o := detectMissingDocsOptions{
				computeSchemaDiff: func() diff.SchemaDiff {
					return diff.ComputeSchemaDiff(tc.oldResourceMap, tc.newResourceMap)
				},
				stdout: &buf,
			}

			err := o.run([]string{tc.repoPath})
			if err != nil {
				t.Fatalf("Error running command: %s", err)
			}

			out := make([]byte, buf.Len())
			buf.Read(out)

			var got []MissingDocsResource
			if err = json.Unmarshal(out, &got); err != nil {
				t.Fatalf("Failed to unmarshall output: %s", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Unexpected result. Want %+v, got %+v. ", tc.want, got)
			}
		})
	}
}
