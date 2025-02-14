package breaking_changes

import (
	"sort"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestRemovingAFieldMessages(t *testing.T) {
	for _, tc := range resourceSchemaRule_RemovingAField_TestCases {
		gotMessages := RemovingAFieldMessages("resource", tc.resourceDiff)

		if len(gotMessages) != len(tc.expectedFields) {
			t.Errorf("RemovingAFieldMessages(%v) got %d messages; want %d", tc.name, len(gotMessages), len(tc.expectedFields))
			continue
		}
		wantFields := tc.expectedFields
		sort.Strings(wantFields)
		sort.Strings(gotMessages)
		for i, field := range wantFields {
			if !strings.Contains(gotMessages[i], field) {
				t.Errorf("RemovingAFieldMessages(%v) got message %q; want field %q", tc.name, gotMessages[i], field)
			}
		}
	}
}

func TestAddingExactlyOneOfMessages(t *testing.T) {
	for _, tc := range resourceSchemaRule_AddingExactlyOneOf_TestCases {
		gotMessages := AddingExactlyOneOfMessages("resource", tc.resourceDiff)
		if len(gotMessages) != len(tc.expectedFields) {
			t.Errorf("AddingExactlyOneOfMessages(%v) got %d messages; want %d", tc.name, len(gotMessages), len(tc.expectedFields))
			continue
		}
		wantFields := tc.expectedFields
		sort.Strings(wantFields)
		sort.Strings(gotMessages)
		for i, field := range wantFields {
			if !strings.Contains(gotMessages[i], field) {
				t.Errorf("AddingExactlyOneOfMessages(%v) got message %q; want field %q", tc.name, gotMessages[i], field)
			}
		}
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

var resourceSchemaRule_AddingExactlyOneOf_TestCases = []resourceSchemaTestCase{
	{
		name: "control",
		resourceDiff: diff.ResourceDiff{
			FieldSets: diff.ResourceFieldSetsDiff{
				Old: diff.ResourceFieldSets{
					ExactlyOneOf: map[string]diff.FieldSet{
						"field-a,field-b": {"field-a": {}, "field-b": {}},
					},
				},
				New: diff.ResourceFieldSets{
					ExactlyOneOf: map[string]diff.FieldSet{
						"field-a,field-b": {"field-a": {}, "field-b": {}},
					},
				},
			},
			Fields: map[string]diff.FieldDiff{
				"field-a": {
					Old: &schema.Schema{Description: "beep", Optional: true},
					New: &schema.Schema{Description: "beep", Optional: true},
				},
				"field-b": {
					Old: &schema.Schema{Description: "boop", Optional: true},
					New: &schema.Schema{Description: "boop", Optional: true},
				},
			},
		},
	},
	{
		name: "adding an existing field to exactly-one-of",
		resourceDiff: diff.ResourceDiff{
			FieldSets: diff.ResourceFieldSetsDiff{
				Old: diff.ResourceFieldSets{
					ExactlyOneOf: map[string]diff.FieldSet{
						"field-a,field-b": {"field-a": {}, "field-b": {}},
					},
				},
				New: diff.ResourceFieldSets{
					ExactlyOneOf: map[string]diff.FieldSet{
						"field-a,field-b,field-c": {"field-a": {}, "field-b": {}, "field-c": {}},
					},
				},
			},
			Fields: map[string]diff.FieldDiff{
				"field-a": {
					Old: &schema.Schema{Description: "beep", Optional: true},
					New: &schema.Schema{Description: "beep", Optional: true},
				},
				"field-b": {
					Old: &schema.Schema{Description: "boop", Optional: true},
					New: &schema.Schema{Description: "boop", Optional: true},
				},
				"field-c": {
					Old: &schema.Schema{Description: "beep", Optional: true},
					New: &schema.Schema{Description: "beep", Optional: true},
				},
			},
		},
		expectedFields: []string{"field-c"},
	},
	{
		name: "adding new exactly-one-of with an existing field",
		resourceDiff: diff.ResourceDiff{
			FieldSets: diff.ResourceFieldSetsDiff{
				Old: diff.ResourceFieldSets{},
				New: diff.ResourceFieldSets{
					ExactlyOneOf: map[string]diff.FieldSet{
						"field-a": {"field-a": {}},
					},
				},
			},
			Fields: map[string]diff.FieldDiff{
				"field-a": {
					Old: &schema.Schema{Description: "beep", Optional: true},
					New: &schema.Schema{Description: "beep", Optional: true},
				},
			},
		},
		expectedFields: []string{"field-a"},
	},
	{
		name: "adding a new field to exactly-one-of",
		resourceDiff: diff.ResourceDiff{
			FieldSets: diff.ResourceFieldSetsDiff{
				Old: diff.ResourceFieldSets{
					ExactlyOneOf: map[string]diff.FieldSet{
						"field-a,field-b": {"field-a": {}, "field-b": {}},
					},
				},
				New: diff.ResourceFieldSets{
					ExactlyOneOf: map[string]diff.FieldSet{
						"field-a,field-b,field-c": {"field-a": {}, "field-b": {}, "field-c": {}},
					},
				},
			},
			Fields: map[string]diff.FieldDiff{
				"field-a": {
					Old: &schema.Schema{Description: "beep", Optional: true},
					New: &schema.Schema{Description: "beep", Optional: true},
				},
				"field-b": {
					Old: &schema.Schema{Description: "boop", Optional: true},
					New: &schema.Schema{Description: "boop", Optional: true},
				},
			},
		},
	},
}
