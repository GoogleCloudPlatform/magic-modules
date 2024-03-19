package cmd

import (
	"bytes"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"testing"
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

			lines := strings.Split(string(out), "\n")
			nonemptyLines := []string{}
			for _, line := range lines {
				if line != "" {
					nonemptyLines = append(nonemptyLines, line)
				}
			}
			if len(nonemptyLines) != tc.expectedViolations {
				t.Errorf("Unexpected number of violations. Want %d, got %d. Output: %s", tc.expectedViolations, len(nonemptyLines), out)
			}
		})
	}
}
