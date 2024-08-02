package breaking_changes

import (
	"sort"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestComputeBreakingChanges(t *testing.T) {
	cases := []struct {
		name           string
		oldResourceMap map[string]*schema.Resource
		newResourceMap map[string]*schema.Resource
		wantViolations []BreakingChange
	}{
		{
			name: "control",
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
		},
		{
			name: "adding resources",
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
				"google-y": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep", Optional: true},
					},
				},
			},
		},
		{
			name: "adding fields",
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
						"field-c": {Description: "beep", Optional: true},
					},
				},
			},
		},
		{
			name: "resource missing",
			oldResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {Description: "beep"},
						"field-b": {Description: "beep"},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{},
			wantViolations: []BreakingChange{
				{
					Message:                "Resource `google-x` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#resource-map-resource-removal-or-rename",
				},
			},
		},
		{
			name: "field missing",
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
					},
				},
			},
			wantViolations: []BreakingChange{
				{
					Message:                "Field `field-b` within resource `google-x` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#resource-schema-field-removal-or-rename",
				},
			},
		},
		{
			name: "optional field to required",
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
						"field-a": {Description: "beep", Required: true},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			wantViolations: []BreakingChange{
				{
					Message:                "Field `field-a` changed from optional to required on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#field-optional-to-required",
				},
			},
		},
		{
			name: "field missing and optional to required",
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
						"field-a": {Description: "beep", Required: true},
					},
				},
			},
			wantViolations: []BreakingChange{
				{
					Message:                "Field `field-a` changed from optional to required on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#field-optional-to-required",
				},
				{
					Message:                "Field `field-b` within resource `google-x` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#resource-schema-field-removal-or-rename",
				},
			},
		},
		{
			name: "field missing, resource missing, and optional to required",
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
			wantViolations: []BreakingChange{
				{
					Message:                "Field `field-a` changed from optional to required on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#field-optional-to-required",
				},
				{
					Message:                "Field `field-b` within resource `google-x` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#resource-schema-field-removal-or-rename",
				},
				{
					Message:                "Resource `google-y` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#resource-map-resource-removal-or-rename",
				},
			},
		},
		{
			name: "removing a subfield",
			oldResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {
							Description: "beep",
							Optional:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sub-field-1": {Description: "beep", Optional: true},
									"sub-field-2": {Description: "beep", Optional: true},
								},
							},
						},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {
							Description: "beep",
							Optional:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sub-field-1": {Description: "beep", Optional: true},
								},
							},
						},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			wantViolations: []BreakingChange{
				{
					Message:                "Field `field-a.sub-field-2` within resource `google-x` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#resource-schema-field-removal-or-rename",
				},
			},
		},
		{
			name: "subfield max shrinking",
			oldResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {
							Description: "beep",
							Optional:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sub-field-1": {Description: "beep", Optional: true, MaxItems: 100},
								},
							},
						},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {
							Description: "beep",
							Optional:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sub-field-1": {Description: "beep", Optional: true, MaxItems: 25},
								},
							},
						},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			wantViolations: []BreakingChange{
				{
					Message:                "Field `field-a.sub-field-1` MinItems went from 100 to 25 on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#field-shrinking-max",
				},
			},
		},
		{
			name: "subfield max shrinking",
			oldResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {
							Description: "beep",
							Optional:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sub-field-1": {Description: "beep", Optional: true, MaxItems: 100},
								},
							},
						},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {
							Description: "beep",
							Optional:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sub-field-1": {Description: "beep", Optional: true, MaxItems: 25},
								},
							},
						},
						"field-b": {Description: "beep", Optional: true},
					},
				},
			},
			wantViolations: []BreakingChange{
				{
					Message:                "Field `field-a.sub-field-1` MinItems went from 100 to 25 on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#field-shrinking-max",
				},
			},
		},
		{
			name: "min growing",
			oldResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {
							Description: "beep",
							Optional:    true,
							MinItems:    1,
						},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google-x": {
					Schema: map[string]*schema.Schema{
						"field-a": {
							Description: "beep",
							Optional:    true,
							MinItems:    4,
						},
					},
				},
			},
			wantViolations: []BreakingChange{
				{
					Message:                "Field `field-a` MinItems went from 1 to 4 on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes/breaking-changes#field-growing-min",
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			schemaDiff := diff.ComputeSchemaDiff(tc.oldResourceMap, tc.newResourceMap)
			violations := ComputeBreakingChanges(schemaDiff)
			for _, v := range violations {
				if strings.Contains(v.Message, "{{") || strings.Contains(v.Message, "}}") {
					t.Errorf("Test `%s` failed: found unreplaced characters in string - %s", tc.name, v)
				}
			}
			sort.Slice(violations, func(i, j int) bool {
				return violations[i].Message < violations[j].Message
			})
			if diff := cmp.Diff(tc.wantViolations, violations); diff != "" {
				t.Errorf("Test `%s` failed: violation diff(-want, +got) = %s", tc.name, diff)
			}
		})
	}
}
