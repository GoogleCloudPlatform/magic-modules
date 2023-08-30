package rules

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type resourceInventoryTestCase struct {
	name               string
	oldResourceMap     map[string]*schema.Resource
	newResourceMap     map[string]*schema.Resource
	expectedViolations int
}

func TestResourceInventoryRule_RemovingAResource(t *testing.T) {
	for _, tc := range resourceInventoryRule_RemovingAResourceTestCases {
		tc.check(resourceInventoryRule_RemovingAResource, t)
	}
}

var resourceInventoryRule_RemovingAResourceTestCases = []resourceInventoryTestCase{
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
		name: "adding a resource",
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
		name: "resource renamed",
		oldResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
		},
		newResourceMap: map[string]*schema.Resource{
			"google-y": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
		},
		expectedViolations: 1,
	},
	{
		name: "resource renamed and another removed",
		oldResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
			"google-z": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
				},
			},
		},
		newResourceMap: map[string]*schema.Resource{
			"google-y": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
		},
		expectedViolations: 2,
	},
}

func (tc *resourceInventoryTestCase) check(rule ResourceInventoryRule, t *testing.T) {
	violations := rule.isRuleBreak(tc.oldResourceMap, tc.newResourceMap)
	if tc.expectedViolations != len(violations) {
		t.Errorf("Test `%s` failed: expected %d violations, got %d", tc.name, tc.expectedViolations, len(violations))
	}
}
