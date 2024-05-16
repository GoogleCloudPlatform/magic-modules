package rules

import (
	"strings"
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
	{
		// TODO: detect this as violation b/300515447
		name:     "field added as required",
		oldField: nil,
		newField: &schema.Schema{
			Description: "beep",
			Required:    true,
		},
		expectedViolation: false,
	},
	{
		name: "field removed",
		oldField: &schema.Schema{
			Description: "beep",
			Optional:    true,
		},
		newField:          nil,
		expectedViolation: false,
	},
}

// !! min max ?
// isRuleBreak: fieldRule_OptionalComputedToOptional_func,

func TestFieldRule_ChangingType(t *testing.T) {
	for _, tc := range fieldRule_ChangingTypeTestCases {
		tc.check(fieldRule_ChangingType, t)
	}
}

var fieldRule_ChangingTypeTestCases = []fieldTestCase{
	{
		name: "control",
		oldField: &schema.Schema{
			Type: schema.TypeBool,
		},
		newField: &schema.Schema{
			Type: schema.TypeBool,
		},
		expectedViolation: false,
	},
	{
		name:     "field added",
		oldField: nil,
		newField: &schema.Schema{
			Type: schema.TypeBool,
		},
		expectedViolation: false,
	},
	{
		name: "field removed",
		oldField: &schema.Schema{
			Type: schema.TypeBool,
		},
		newField:          nil,
		expectedViolation: false,
	},
	{
		name: "field sub-element type control",
		oldField: &schema.Schema{
			Type: schema.TypeList,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		newField: &schema.Schema{
			Type: schema.TypeList,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		expectedViolation: false,
	},
	{
		name: "field transition bool -> string",
		oldField: &schema.Schema{
			Type: schema.TypeBool,
		},
		newField: &schema.Schema{
			Type: schema.TypeString,
		},
		expectedViolation: true,
	},
	{
		name: "field transition string -> int",
		oldField: &schema.Schema{
			Type: schema.TypeString,
		},
		newField: &schema.Schema{
			Type: schema.TypeInt,
		},
		expectedViolation: true,
	},
	{
		name: "field transition sub-element type ",
		oldField: &schema.Schema{
			Type: schema.TypeList,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		newField: &schema.Schema{
			Type: schema.TypeList,
			Elem: &schema.Schema{Type: schema.TypeBool},
		},
		expectedViolation: true,
	},
}

func TestFieldRule_DefaultModification(t *testing.T) {
	for _, tc := range fieldRule_DefaultModificationTestCases {
		tc.check(fieldRule_DefaultModification, t)
	}
}

var fieldRule_DefaultModificationTestCases = []fieldTestCase{
	{
		name: "control",
		oldField: &schema.Schema{
			Default: "same",
		},
		newField: &schema.Schema{
			Default: "same",
		},
		expectedViolation: false,
	},
	{
		name:              "control - no default",
		oldField:          &schema.Schema{},
		newField:          &schema.Schema{},
		expectedViolation: false,
	},
	{
		name: "default value change - string",
		oldField: &schema.Schema{
			Default: "1",
		},
		newField: &schema.Schema{
			Default: "2",
		},
		expectedViolation: true,
	},
	{
		name: "default value change - int",
		oldField: &schema.Schema{
			Default: 1,
		},
		newField: &schema.Schema{
			Default: 2,
		},
		expectedViolation: true,
	},
	{
		name: "default value change - bool",
		oldField: &schema.Schema{
			Default: false,
		},
		newField: &schema.Schema{
			Default: true,
		},
		expectedViolation: true,
	},
	{
		name: "default value change - float",
		oldField: &schema.Schema{
			Default: 1.2,
		},
		newField: &schema.Schema{
			Default: 3.4,
		},
		expectedViolation: true,
	},
	{
		name: "default value change - transitions int -> string",
		oldField: &schema.Schema{
			Default: 1,
		},
		newField: &schema.Schema{
			Default: "1",
		},
		expectedViolation: true,
	},
	{
		name: "default value change - transitions bool -> string",
		oldField: &schema.Schema{
			Default: false,
		},
		newField: &schema.Schema{
			Default: "false",
		},
		expectedViolation: true,
	},
	{
		name:     "field added",
		oldField: nil,
		newField: &schema.Schema{
			Default: "same",
		},
		expectedViolation: false,
	},
	{
		name: "field removed",
		oldField: &schema.Schema{
			Default: "same",
		},
		newField:          nil,
		expectedViolation: false,
	},
}

func TestFieldRule_BecomingComputedOnly(t *testing.T) {
	for _, tc := range fieldRule_BecomingComputedOnlyTestCases {
		tc.check(fieldRule_BecomingComputedOnly, t)
	}
}

var fieldRule_BecomingComputedOnlyTestCases = []fieldTestCase{
	{
		name: "control - already computed",
		oldField: &schema.Schema{
			Computed: true,
		},
		newField: &schema.Schema{
			Computed: true,
		},
		expectedViolation: false,
	},
	{
		name: "control - not computed",
		oldField: &schema.Schema{
			Optional: true,
		},
		newField: &schema.Schema{
			Optional: true,
		},
		expectedViolation: false,
	},
	{
		name: "control - computed + optional",
		oldField: &schema.Schema{
			Optional: true,
			Computed: true,
		},
		newField: &schema.Schema{
			Optional: true,
			Computed: true,
		},
		expectedViolation: false,
	},
	{
		name: "transition computed + optional -> computed",
		oldField: &schema.Schema{
			Optional: true,
			Computed: true,
		},
		newField: &schema.Schema{
			Computed: true,
		},
		expectedViolation: true,
	},
	{
		name: "transition optional -> computed",
		oldField: &schema.Schema{
			Optional: true,
		},
		newField: &schema.Schema{
			Computed: true,
		},
		expectedViolation: true,
	},
	{
		name: "transition required -> computed",
		oldField: &schema.Schema{
			Required: true,
		},
		newField: &schema.Schema{
			Computed: true,
		},
		expectedViolation: true,
	},
	{
		name:     "added computed field",
		oldField: nil,
		newField: &schema.Schema{
			Computed: true,
		},
		expectedViolation: false,
	},
	{
		name: "removed computed field",
		oldField: &schema.Schema{
			Computed: true,
		},
		newField:          nil,
		expectedViolation: false,
	},
}

func TestFieldRule_OptionalComputedToOptional(t *testing.T) {
	for _, tc := range fieldRule_OptionalComputedToOptionalTestCases {
		tc.check(fieldRule_OptionalComputedToOptional, t)
	}
}

var fieldRule_OptionalComputedToOptionalTestCases = []fieldTestCase{
	{
		name: "control - static",
		oldField: &schema.Schema{
			Computed: true,
			Optional: true,
		},
		newField: &schema.Schema{
			Computed: true,
			Optional: true,
		},
		expectedViolation: false,
	},
	{
		name: "control - not relevant",
		oldField: &schema.Schema{
			Required: true,
		},
		newField: &schema.Schema{
			Optional: true,
		},
		expectedViolation: false,
	},
	{
		name: "transition o+c -> o",
		oldField: &schema.Schema{
			Optional: true,
			Computed: true,
		},
		newField: &schema.Schema{
			Optional: true,
		},
		expectedViolation: true,
	},
	{
		name:     "field added",
		oldField: nil,
		newField: &schema.Schema{
			Optional: true,
		},
		expectedViolation: false,
	},
	{
		name: "field removed",
		oldField: &schema.Schema{
			Optional: true,
			Computed: true,
		},
		newField:          nil,
		expectedViolation: false,
	},
}

func TestFieldRule_GrowingMin(t *testing.T) {
	for _, tc := range fieldRule_GrowingMinTestCases {
		tc.check(fieldRule_GrowingMin, t)
	}
}

var fieldRule_GrowingMinTestCases = []fieldTestCase{
	{
		name: "control:min - static",
		oldField: &schema.Schema{
			MinItems: 1,
		},
		newField: &schema.Schema{
			MinItems: 1,
		},
		expectedViolation: false,
	},
	{
		name:              "control - unset",
		oldField:          &schema.Schema{},
		newField:          &schema.Schema{},
		expectedViolation: false,
	},
	{
		name: "control - min shrinks",
		oldField: &schema.Schema{
			MinItems: 20,
		},
		newField: &schema.Schema{
			MinItems: 2,
		},
		expectedViolation: false,
	},
	{
		name: "min grows",
		oldField: &schema.Schema{
			MinItems: 2,
		},
		newField: &schema.Schema{
			MinItems: 20,
		},
		expectedViolation: true,
	},
	{
		name:     "field added",
		oldField: nil,
		newField: &schema.Schema{
			MinItems: 1,
		},
		expectedViolation: false,
	},
	{
		name: "field removed",
		oldField: &schema.Schema{
			MinItems: 1,
		},
		newField:          nil,
		expectedViolation: false,
	},
	{
		name:     "min unset to defined",
		oldField: &schema.Schema{},
		newField: &schema.Schema{
			MinItems: 2,
		},
		expectedViolation: true,
	},
}

func TestFieldRule_ShrinkingMax(t *testing.T) {
	for _, tc := range fieldRule_ShrinkingMaxTestCases {
		tc.check(fieldRule_ShrinkingMax, t)
	}
}

var fieldRule_ShrinkingMaxTestCases = []fieldTestCase{
	{
		name: "control:max - static",
		oldField: &schema.Schema{
			MaxItems: 2,
		},
		newField: &schema.Schema{
			MaxItems: 2,
		},
		expectedViolation: false,
	},
	{
		name:              "control - unset",
		oldField:          &schema.Schema{},
		newField:          &schema.Schema{},
		expectedViolation: false,
	},
	{
		name: "control - max grows",
		oldField: &schema.Schema{
			MaxItems: 2,
		},
		newField: &schema.Schema{
			MaxItems: 20,
		},
		expectedViolation: false,
	},
	{
		name: "max shrinks",
		oldField: &schema.Schema{
			MaxItems: 20,
		},
		newField: &schema.Schema{
			MaxItems: 2,
		},
		expectedViolation: true,
	},
	{
		name:     "field added",
		oldField: nil,
		newField: &schema.Schema{
			MaxItems: 2,
		},
		expectedViolation: false,
	},
	{
		name: "field removed",
		oldField: &schema.Schema{
			MaxItems: 2,
		},
		newField:          nil,
		expectedViolation: false,
	},
	{
		name:     "max unset to defined",
		oldField: &schema.Schema{},
		newField: &schema.Schema{
			MaxItems: 2,
		},
		expectedViolation: true,
	},
}

func (tc *fieldTestCase) check(rule FieldRule, t *testing.T) {
	breakage := rule.isRuleBreak(tc.oldField, tc.newField, MessageContext{})

	violation := breakage != nil
	if breakage != nil && strings.Contains(breakage.Message, "{{") {
		t.Errorf("Test `%s` failed: replacements for `{{<val>}}` not successful ", tc.name)
	}
	if tc.expectedViolation != violation {
		t.Errorf("Test `%s` failed: expected %v violations, got %v", tc.name, tc.expectedViolation, violation)
	}
}

func TestBreakingMessage(t *testing.T) {
	breakageMessage := fieldRule_OptionalComputedToOptional.IsRuleBreak(
		&schema.Schema{
			Optional: true,
			Computed: true,
		},
		&schema.Schema{
			Optional: true,
		},
		MessageContext{
			Resource: "a",
			Field:    "b",
		},
	)

	if !strings.Contains(breakageMessage.Message, "Field `b` transitioned from optional+computed to optional `a`") {
		t.Errorf("Test `%s` failed: replacements for `{{<val>}}` not successful ", "TestBreakingMessage")
	}

}
