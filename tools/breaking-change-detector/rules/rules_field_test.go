package rules

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type fieldTestCase struct {
	name              string
	oldField          *schema.Schema
	newField          *schema.Schema
	expectedViolation bool
}

func TestFieldRule_BecomingRequired(t *testing.T) {
	for _, tc := range fieldRule_BecomingRequiredTestCases {
		tc.check(fieldRule_BecomingRequired, t)
	}
}

var fieldRule_BecomingRequiredTestCases = []fieldTestCase{
	{
		name: "control",
		oldField: &schema.Schema{
			Description: "beep",
			Optional:    true,
		},
		newField: &schema.Schema{
			Description: "beep",
			Optional:    true,
		},
		expectedViolation: false,
	},
	{
		name: "optional to optional + computed",
		oldField: &schema.Schema{
			Description: "beep",
			Optional:    true,
		},
		newField: &schema.Schema{
			Description: "beep",
			Optional:    true,
			Computed:    true,
		},
		expectedViolation: false,
	},
	{
		name: "optional to required",
		oldField: &schema.Schema{
			Description: "beep",
			Optional:    true,
		},
		newField: &schema.Schema{
			Description: "beep",
			Required:    true,
		},
		expectedViolation: true,
	},
	{
		name: "optional computed to required",
		oldField: &schema.Schema{
			Description: "beep",
			Computed:    true,
			Optional:    true,
		},
		newField: &schema.Schema{
			Description: "beep",
			Required:    true,
		},
		expectedViolation: true,
	},
}

func (tc *fieldTestCase) check(rule FieldRule, t *testing.T) {
	violation := rule.isRuleBreak(tc.oldField, tc.newField)
	if tc.expectedViolation != violation {
		t.Errorf("Test `%s` failed: expected %v violations, got %v", tc.name, tc.expectedViolation, violation)
	}
}
