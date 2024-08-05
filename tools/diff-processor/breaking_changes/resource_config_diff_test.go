package breaking_changes

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type resourceInventoryTestCase struct {
	name           string
	old            *schema.Resource
	new            *schema.Resource
	wantViolations bool
}

func TestResourceInventoryRule_RemovingAResource(t *testing.T) {
	for _, tc := range resourceConfigRemovingAResourceTestCases {
		got := ResourceConfigRemovingAResource.Messages("resource", tc.old, tc.new)
		gotViolations := len(got) > 0
		if tc.wantViolations != gotViolations {
			t.Errorf("ResourceConfigRemovingAResource.Messages(%v) violations not expected. Got %v, want %v", tc.name, gotViolations, tc.wantViolations)
		}
	}
}

var resourceConfigRemovingAResourceTestCases = []resourceInventoryTestCase{
	{
		name:           "control",
		old:            &schema.Resource{},
		new:            &schema.Resource{},
		wantViolations: false,
	},
	{
		name:           "resource added",
		old:            nil,
		new:            &schema.Resource{},
		wantViolations: false,
	},
	{
		name:           "resource removed",
		old:            &schema.Resource{},
		new:            nil,
		wantViolations: true,
	},
}