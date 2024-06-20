package rules

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
		wantViolations []*BreakingChange
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
			wantViolations: []*BreakingChange{
				{
					Resource:               "google-x",
					Message:                "Resource `google-x` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#resource-map-resource-removal-or-rename",
					RuleTemplate:           "Resource {{resource}} was either removed or renamed",
					RuleDefinition:         "In terraform resources should be retained whenever possible. A removable of an resource will result in a configuration breakage wherever a dependency on that resource exists. Renaming or Removing a resources are functionally equivalent in terms of configuration breakages.",
					RuleName:               "Removing or Renaming an Resource",
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
			wantViolations: []*BreakingChange{
				{
					Resource:               "google-x",
					Field:                  "field-b",
					Message:                "Field `field-b` within resource `google-x` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#resource-schema-field-removal-or-rename",
					RuleTemplate:           "Field {{field}} within resource {{resource}} was either removed or renamed",
					RuleDefinition:         "In terraform fields should be retained whenever possible. A removable of an field will result in a configuration breakage wherever a dependency on that field exists. Renaming or Removing a field are functionally equivalent in terms of configuration breakages.",
					RuleName:               "Removing or Renaming an field",
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
			wantViolations: []*BreakingChange{
				{
					Resource:               "google-x",
					Field:                  "field-a",
					Message:                "Field `field-a` changed from optional to required on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#field-optional-to-required",
					RuleTemplate:           "Field {{field}} changed from optional to required on {{resource}}",
					RuleName:               "Field becoming Required Field",
					RuleDefinition:         "A field cannot become required as existing configs may not have this field defined. Thus, breaking configs in sequential plan or applies. If you are adding Required to a field so a block won't remain empty, this can cause two issues. First if it's a singular nested field the block may gain more fields later and it's not clear whether the field is actually required so it may be misinterpreted by future contributors. Second if users are defining empty blocks in existing configurations this change will break them. Consider these points in admittance of this type of change.",
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
			wantViolations: []*BreakingChange{
				{
					Resource:               "google-x",
					Field:                  "field-a",
					Message:                "Field `field-a` changed from optional to required on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#field-optional-to-required",
					RuleTemplate:           "Field {{field}} changed from optional to required on {{resource}}",
					RuleName:               "Field becoming Required Field",
					RuleDefinition:         "A field cannot become required as existing configs may not have this field defined. Thus, breaking configs in sequential plan or applies. If you are adding Required to a field so a block won't remain empty, this can cause two issues. First if it's a singular nested field the block may gain more fields later and it's not clear whether the field is actually required so it may be misinterpreted by future contributors. Second if users are defining empty blocks in existing configurations this change will break them. Consider these points in admittance of this type of change.",
				},
				{
					Resource:               "google-x",
					Field:                  "field-b",
					Message:                "Field `field-b` within resource `google-x` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#resource-schema-field-removal-or-rename",
					RuleTemplate:           "Field {{field}} within resource {{resource}} was either removed or renamed",
					RuleDefinition:         "In terraform fields should be retained whenever possible. A removable of an field will result in a configuration breakage wherever a dependency on that field exists. Renaming or Removing a field are functionally equivalent in terms of configuration breakages.",
					RuleName:               "Removing or Renaming an field",
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
			wantViolations: []*BreakingChange{
				{
					Resource:               "google-x",
					Field:                  "field-a",
					Message:                "Field `field-a` changed from optional to required on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#field-optional-to-required",
					RuleTemplate:           "Field {{field}} changed from optional to required on {{resource}}",
					RuleName:               "Field becoming Required Field",
					RuleDefinition:         "A field cannot become required as existing configs may not have this field defined. Thus, breaking configs in sequential plan or applies. If you are adding Required to a field so a block won't remain empty, this can cause two issues. First if it's a singular nested field the block may gain more fields later and it's not clear whether the field is actually required so it may be misinterpreted by future contributors. Second if users are defining empty blocks in existing configurations this change will break them. Consider these points in admittance of this type of change.",
				},
				{
					Resource:               "google-x",
					Field:                  "field-b",
					Message:                "Field `field-b` within resource `google-x` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#resource-schema-field-removal-or-rename",
					RuleTemplate:           "Field {{field}} within resource {{resource}} was either removed or renamed",
					RuleDefinition:         "In terraform fields should be retained whenever possible. A removable of an field will result in a configuration breakage wherever a dependency on that field exists. Renaming or Removing a field are functionally equivalent in terms of configuration breakages.",
					RuleName:               "Removing or Renaming an field",
				},
				{
					Resource:               "google-y",
					Message:                "Resource `google-y` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#resource-map-resource-removal-or-rename",
					RuleTemplate:           "Resource {{resource}} was either removed or renamed",
					RuleDefinition:         "In terraform resources should be retained whenever possible. A removable of an resource will result in a configuration breakage wherever a dependency on that resource exists. Renaming or Removing a resources are functionally equivalent in terms of configuration breakages.",
					RuleName:               "Removing or Renaming an Resource",
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
			wantViolations: []*BreakingChange{
				{
					Resource:               "google-x",
					Field:                  "field-a.sub-field-2",
					Message:                "Field `field-a.sub-field-2` within resource `google-x` was either removed or renamed",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#resource-schema-field-removal-or-rename",
					RuleTemplate:           "Field {{field}} within resource {{resource}} was either removed or renamed",
					RuleDefinition:         "In terraform fields should be retained whenever possible. A removable of an field will result in a configuration breakage wherever a dependency on that field exists. Renaming or Removing a field are functionally equivalent in terms of configuration breakages.",
					RuleName:               "Removing or Renaming an field",
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
			wantViolations: []*BreakingChange{
				{
					Resource:               "google-x",
					Field:                  "field-a.sub-field-1",
					Message:                "Field `field-a.sub-field-1` MinItems went from 100 to 25 on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#field-shrinking-max",
					RuleTemplate:           "Field {{field}} MinItems went from 100 to 25 on {{resource}}",
					RuleName:               "Shrinking Maximum Items",
					RuleDefinition:         "MaxItems cannot shrink. Otherwise existing terraform configurations that don't satisfy this rule will break.",
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
			wantViolations: []*BreakingChange{
				{
					Resource:               "google-x",
					Field:                  "field-a.sub-field-1",
					Message:                "Field `field-a.sub-field-1` MinItems went from 100 to 25 on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#field-shrinking-max",
					RuleTemplate:           "Field {{field}} MinItems went from 100 to 25 on {{resource}}",
					RuleName:               "Shrinking Maximum Items",
					RuleDefinition:         "MaxItems cannot shrink. Otherwise existing terraform configurations that don't satisfy this rule will break.",
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
			wantViolations: []*BreakingChange{
				{
					Resource:               "google-x",
					Field:                  "field-a",
					Message:                "Field `field-a` MinItems went from 1 to 4 on `google-x`",
					DocumentationReference: "https://googlecloudplatform.github.io/magic-modules/develop/breaking-changes#field-growing-min",
					RuleTemplate:           "Field {{field}} MinItems went from 1 to 4 on {{resource}}",
					RuleName:               "Growing Minimum Items",
					RuleDefinition:         "MinItems cannot grow. Otherwise existing terraform configurations that don't satisfy this rule will break.",
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
