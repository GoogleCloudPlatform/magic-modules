package rules

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type resourceInventoryTestCase struct {
	name     string
	old      *schema.Resource
	new      *schema.Resource
	expected bool
}

func TestResourceInventoryRule_RemovingAResource(t *testing.T) {
	for _, tc := range resourceInventoryRule_RemovingAResourceTestCases {
		tc.check(resourceInventoryRule_RemovingAResource, t)
	}
}

var resourceInventoryRule_RemovingAResourceTestCases = []resourceInventoryTestCase{
	{
		name:     "control",
		old:      &schema.Resource{},
		new:      &schema.Resource{},
		expected: false,
	},
	{
		name:     "resource added",
		old:      nil,
		new:      &schema.Resource{},
		expected: false,
	},
	{
		name:     "resource removed",
		old:      &schema.Resource{},
		new:      nil,
		expected: true,
	},
}

func (tc *resourceInventoryTestCase) check(rule ResourceInventoryRule, t *testing.T) {
	got := rule.isRuleBreak(tc.old, tc.new)
	if tc.expected != got {
		t.Errorf("Test `%s` failed: want %t, got %t", tc.name, tc.expected, got)
	}
}
