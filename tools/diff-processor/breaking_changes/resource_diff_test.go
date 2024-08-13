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
