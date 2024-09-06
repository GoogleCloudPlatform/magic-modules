package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/breaking_changes"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestBreakingChangesCmd(t *testing.T) {
	cases := map[string]struct {
		oldResourceMap     map[string]*schema.Resource
		newResourceMap     map[string]*schema.Resource
		expectedViolations int
	}{
		"no breaking changes": {
			oldResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Optional: true},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Optional: true},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			expectedViolations: 0,
		},
		"resource missing": {
			oldResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep"},
						"field-b": {Description: "beep"},
					},
				},
			},
			newResourceMap:     map[string]*schema.Resource{},
			expectedViolations: 1,
		},
		"field missing, resource missing, and optional to required": {
			oldResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Optional: true},
						"field-b": {Description: "beep", Optional: true},
					},
				},
				"google-y": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Required: true},
					},
				},
			},
			expectedViolations: 3,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			o := breakingChangesOptions{
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

			var got []breaking_changes.BreakingChange
			if err = json.Unmarshal(out, &got); err != nil {
				t.Fatalf("Failed to unmarshall output: %s", err)
			}

			if len(got) != tc.expectedViolations {
				t.Errorf("Unexpected number of violations. Want %d, got %d. Output: %s", tc.expectedViolations, len(got), out)
			}
		})
	}
}
