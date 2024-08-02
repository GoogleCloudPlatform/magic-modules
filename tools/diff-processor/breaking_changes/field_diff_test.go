package breaking_changes

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

func TestFieldBecomingRequired(t *testing.T) {
	for _, tc := range FieldBecomingRequiredTestCases {
		tc.check(FieldBecomingRequired, t)
	}
}

var FieldBecomingRequiredTestCases = []fieldTestCase{
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
// isRuleBreak: FieldOptionalComputedToOptional_func,

func TestFieldChangingType(t *testing.T) {
	for _, tc := range FieldChangingTypeTestCases {
		tc.check(FieldChangingType, t)
	}
}

var FieldChangingTypeTestCases = []fieldTestCase{
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

func TestFieldDefaultModification(t *testing.T) {
	for _, tc := range FieldDefaultModificationTestCases {
		tc.check(FieldDefaultModification, t)
	}
}

var FieldDefaultModificationTestCases = []fieldTestCase{
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

func TestFieldBecomingComputedOnly(t *testing.T) {
	for _, tc := range FieldBecomingComputedOnlyTestCases {
		tc.check(FieldBecomingComputedOnly, t)
	}
}

var FieldBecomingComputedOnlyTestCases = []fieldTestCase{
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

func TestFieldOptionalComputedToOptional(t *testing.T) {
	for _, tc := range FieldOptionalComputedToOptionalTestCases {
		tc.check(FieldOptionalComputedToOptional, t)
	}
}

var FieldOptionalComputedToOptionalTestCases = []fieldTestCase{
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

func TestFieldGrowingMin(t *testing.T) {
	for _, tc := range FieldGrowingMinTestCases {
		tc.check(FieldGrowingMin, t)
	}
}

var FieldGrowingMinTestCases = []fieldTestCase{
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

func TestFieldShrinkingMax(t *testing.T) {
	for _, tc := range FieldShrinkingMaxTestCases {
		tc.check(FieldShrinkingMax, t)
	}
}

var FieldShrinkingMaxTestCases = []fieldTestCase{
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

func TestFieldAddingSubfieldToConfigModeAttr(t *testing.T) {
	for _, tc := range FieldAddingSubfieldToConfigModeAttrTestCases {
		tc.check(FieldAddingSubfieldToConfigModeAttr, t)
	}
}

var FieldAddingSubfieldToConfigModeAttrTestCases = []fieldTestCase{
	{
		name: "no new subfields",
		oldField: &schema.Schema{
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "beep",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"field_one": {},
				},
			},
		},
		newField: &schema.Schema{
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "beep",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"field_one": {},
				},
			},
		},
		expectedViolation: false,
	},
	{
		name: "adding a subfield with no SchemaConfigModeAttr",
		oldField: &schema.Schema{
			Description: "beep",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"field_one": {},
				},
			},
		},
		newField: &schema.Schema{
			Description: "beep",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"field_one": {},
					"field_two": {},
				},
			},
		},
		expectedViolation: false,
	},
	{
		name: "adding a field with SchemaConfigModeAttr",
		oldField: &schema.Schema{
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "beep",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"field_one": {},
				},
			},
		},
		newField: &schema.Schema{
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "beep",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"field_one": {},
					"field_two": {},
				},
			},
		},
		expectedViolation: true,
	},
}

func (tc *fieldTestCase) check(rule FieldDiffRule, t *testing.T) {
	messages := rule.Messages("resource", "field", tc.oldField, tc.newField)

	violation := len(messages) > 0
	if tc.expectedViolation != violation {
		t.Errorf("Test `%s` failed: expected %v violations, got %v", tc.name, tc.expectedViolation, violation)
	}
}
