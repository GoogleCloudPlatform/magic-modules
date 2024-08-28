package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDetectMissingDocs(t *testing.T) {
	cases := []struct {
		name           string
		oldResourceMap map[string]*schema.Resource
		newResourceMap map[string]*schema.Resource
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
			name: "multiple new fields missing doc",
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
					Fields: []detector.MissingDocField{
						{
							Field:   "field-a",
							Section: "Arguments Reference",
						},
						{
							Field:   "field-b",
							Section: "Arguments Reference",
						},
					},
				},
			},
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

			err := o.run([]string{t.TempDir()})
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
