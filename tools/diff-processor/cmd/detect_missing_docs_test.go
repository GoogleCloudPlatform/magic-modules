package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDetectMissingDocs(t *testing.T) {
	cases := []struct {
		name           string
		oldResourceMap map[string]*schema.Resource
		newResourceMap map[string]*schema.Resource
		want           []MissingDocsInfo
	}{
		{
			name: "no new fields",
			oldResourceMap: map[string]*schema.Resource{
				"google_x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Computed: true, Optional: true},
						"field-b": {Description: "beep", Computed: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Computed: true, Optional: true},
						"field-b": {Description: "beep", Computed: true},
					},
				},
			},
			want: []MissingDocsInfo{},
		},
		{
			name:           "multiple new fields missing doc",
			oldResourceMap: map[string]*schema.Resource{},
			newResourceMap: map[string]*schema.Resource{
				"google_x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Computed: true, Optional: true},
						"field-b": {Description: "beep", Computed: true},
					},
				},
			},
			want: []MissingDocsInfo{
				{
					Name:     "google_x",
					FilePath: "/website/docs/r/x.html.markdown",
					Fields:   []string{"field-a", "field-b"},
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
				newResourceSchema: tc.newResourceMap,
				stdout:            &buf,
			}

			err := o.run([]string{t.TempDir()})
			if err != nil {
				t.Fatalf("Error running command: %s", err)
			}

			out := make([]byte, buf.Len())
			buf.Read(out)

			var got []MissingDocsInfo
			if err = json.Unmarshal(out, &got); err != nil {
				t.Fatalf("Failed to unmarshall output: %s", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Unexpected result. Want %+v, got %+v. ", tc.want, got)
			}
		})
	}
}
