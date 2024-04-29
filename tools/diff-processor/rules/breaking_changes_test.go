package rules

import (
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestComputeBreakingChanges(t *testing.T) {
	cases := []struct {
		name               string
		oldResourceMap     map[string]*schema.Resource
		newResourceMap     map[string]*schema.Resource
		expectedViolations int
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
			expectedViolations: 0,
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
			expectedViolations: 0,
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
			expectedViolations: 0,
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
			newResourceMap:     map[string]*schema.Resource{},
			expectedViolations: 1,
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
			expectedViolations: 1,
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
			expectedViolations: 1,
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
			expectedViolations: 2,
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
			expectedViolations: 3,
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
			expectedViolations: 1,
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
			expectedViolations: 1,
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
			expectedViolations: 1,
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
			expectedViolations: 1,
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
			if tc.expectedViolations != len(violations) {
				t.Errorf("Test `%s` failed: expected %d violations, got %d", tc.name, tc.expectedViolations, len(violations))
			}
		})
	}
}
