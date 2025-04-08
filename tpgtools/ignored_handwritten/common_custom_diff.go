package google

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func IsShrinkageIpCidr(old, new, _ interface{}) bool {
	_, oldCidr, oldErr := net.ParseCIDR(old.(string))
	_, newCidr, newErr := net.ParseCIDR(new.(string))

	if oldErr != nil || newErr != nil {
		// This should never happen. The ValidateFunc on the field ensures it.
		return false
	}

	oldStart, oldEnd := cidr.AddressRange(oldCidr)

	if newCidr.Contains(oldStart) && newCidr.Contains(oldEnd) {
		// This is a CIDR range expansion, no need to ForceNew, we have an update method for it.
		return false
	}

	return true
}

func resourceComputeSubnetworkSecondaryIpRangeSetStyleDiff(diff *schema.ResourceDiff, meta interface{}) error {
	keys := diff.GetChangedKeysPrefix("secondary_ip_range")
	if len(keys) == 0 {
		return nil
	}
	oldCount, newCount := diff.GetChange("secondary_ip_range.#")
	var count int
	// There could be duplicates - worth continuing even if the counts are unequal.
	if oldCount.(int) < newCount.(int) {
		count = newCount.(int)
	} else {
		count = oldCount.(int)
	}

	if count < 1 {
		return nil
	}
	old := make([]interface{}, 0, count)
	new := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		o, n := diff.GetChange(fmt.Sprintf("secondary_ip_range.%d", i))

		if o != nil {
			old = append(old, o)
		}
		if n != nil {
			new = append(new, n)
		}
	}

	oldSet := schema.NewSet(schema.HashResource(ResourceComputeSubnetwork().Schema["secondary_ip_range"].Elem.(*schema.Resource)), old)
	newSet := schema.NewSet(schema.HashResource(ResourceComputeSubnetwork().Schema["secondary_ip_range"].Elem.(*schema.Resource)), new)

	if oldSet.Equal(newSet) {
		if err := diff.Clear("secondary_ip_range"); err != nil {
			return err
		}
	}

	return nil
}
