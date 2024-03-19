package rules

import (
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceSchemaRule_RemovingAField(t *testing.T) {
	for _, tc := range resourceSchemaRule_RemovingAField_TestCases {
		tc.check(resourceSchemaRule_RemovingAField, t)
	}
}

type resourceSchemaTestCase struct {
	name           string
	resourceDiff   diff.ResourceDiff
	expectedFields []string
}

var resourceSchemaRule_RemovingAField_TestCases = []resourceSchemaTestCase{
	{
		name: "control",
		resourceDiff: diff.ResourceDiff{
			Fields: map[string]diff.FieldDiff{
				"field-a": diff.FieldDiff{
					Old: &schema.Schema{Description: "beep", Optional: true},
					New: &schema.Schema{Description: "beep", Optional: true},
				},
			},
		},
		expectedFields: []string{},
	},
	{
		name: "adding a field",
		resourceDiff: diff.ResourceDiff{
			Fields: map[string]diff.FieldDiff{
				"field-a": diff.FieldDiff{
					Old: nil,
					New: &schema.Schema{Description: "beep", Optional: true},
				},
			},
		},
		expectedFields: []string{},
	},
	{
		name: "removing a field",
		resourceDiff: diff.ResourceDiff{
			Fields: map[string]diff.FieldDiff{
				"field-a": diff.FieldDiff{
					Old: &schema.Schema{Description: "beep", Optional: true},
					New: nil,
				},
			},
		},
		expectedFields: []string{"field-a"},
	},
	{
		name: "removing multiple fields",
		resourceDiff: diff.ResourceDiff{
			Fields: map[string]diff.FieldDiff{
				"field-a": diff.FieldDiff{
					Old: &schema.Schema{Description: "beep", Optional: true},
					New: nil,
				},
				"field-b": diff.FieldDiff{
					Old: &schema.Schema{Description: "beep", Optional: true},
					New: nil,
				},
			},
		},
		expectedFields: []string{"field-a", "field-b"},
	},
}

func (tc *resourceSchemaTestCase) check(rule ResourceSchemaRule, t *testing.T) {
	fields := rule.IsRuleBreak(tc.resourceDiff)
	less := func(a, b string) bool { return a < b }
	if !cmp.Equal(fields, tc.expectedFields, cmpopts.SortSlices(less)) {
		t.Errorf("Test `%s` failed: wanted %v , got %v", tc.name, tc.expectedFields, fields)
	}
}
