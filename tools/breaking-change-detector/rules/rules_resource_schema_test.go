package rules

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceSchemaRule_RemovingAField(t *testing.T) {
	for _, tc := range resourceSchemaRule_RemovingAField_TestCases {
		tc.check(resourceSchemaRule_RemovingAField, t)
	}
}

type resourceSchemaTestCase struct {
	name               string
	oldResourceSchema  map[string]*schema.Schema
	newResourceSchema  map[string]*schema.Schema
	expectedViolations int
}

var resourceSchemaRule_RemovingAField_TestCases = []resourceSchemaTestCase{
	{
		name: "control",
		oldResourceSchema: map[string]*schema.Schema{
			"field-a": {Description: "beep", Optional: true},
			"field-b": {Description: "beep", Optional: true},
		},
		newResourceSchema: map[string]*schema.Schema{
			"field-a": {Description: "beep", Optional: true},
			"field-b": {Description: "beep", Optional: true},
		},
		expectedViolations: 0,
	},
	{
		name: "adding a field",
		oldResourceSchema: map[string]*schema.Schema{
			"field-a": {Description: "beep", Optional: true},
			"field-b": {Description: "beep", Optional: true},
		},
		newResourceSchema: map[string]*schema.Schema{
			"field-a": {Description: "beep", Optional: true},
			"field-b": {Description: "beep", Optional: true},
			"field-c": {Description: "beep", Optional: true},
		},
		expectedViolations: 0,
	},
	{
		name: "renaming a field",
		oldResourceSchema: map[string]*schema.Schema{
			"field-a": {Description: "beep", Optional: true},
			"field-b": {Description: "beep", Optional: true},
		},
		newResourceSchema: map[string]*schema.Schema{
			"field-a": {Description: "beep", Optional: true},
			"field-d": {Description: "beep", Optional: true},
		},
		expectedViolations: 1,
	},
	{
		name: "removing a field",
		oldResourceSchema: map[string]*schema.Schema{
			"field-a": {Description: "beep", Optional: true},
			"field-b": {Description: "beep", Optional: true},
		},
		newResourceSchema: map[string]*schema.Schema{
			"field-a": {Description: "beep", Optional: true},
		},
		expectedViolations: 1,
	},
	{
		name: "renaming a field and removing a field",
		oldResourceSchema: map[string]*schema.Schema{
			"field-a": {Description: "beep", Optional: true},
			"field-b": {Description: "beep", Optional: true},
		},
		newResourceSchema: map[string]*schema.Schema{
			"field-z": {Description: "beep", Optional: true},
		},
		expectedViolations: 2,
	},
}

func (tc *resourceSchemaTestCase) check(rule ResourceSchemaRule, t *testing.T) {
	violations := rule.isRuleBreak(tc.oldResourceSchema, tc.newResourceSchema)
	if tc.expectedViolations != len(violations) {
		t.Errorf("Test `%s` failed: expected %d violations, got %d", tc.name, tc.expectedViolations, len(violations))
	}
}
