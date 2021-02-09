package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func suppressPartnerProviderPairingKey(k, old, new string, d *schema.ResourceData) bool {
	attachmentType := d.Get("type")
	if state, ok := d.GetOk("state"); ok {
		if attachmentType == "PARTNER_PROVIDER" && state == "ACTIVE" {
			return true
		}
	}

	return false
}
