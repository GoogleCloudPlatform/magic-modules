package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func suppressPartnerProviderPairingKey(k, old, new string, d *schema.ResourceData) bool {
	/*
		This prevents something like pairing_key from triggering a destroy and update cycle
		on update when used like:
		pairing_key = google_compute_interconnect_attachment.attachment_partner.pairing_key
		This happens because pairing_key gets set to "" after a successful creation when
		type is `PARTNER_PROVIDER` but subsequent calls try to add the pairing_key
	*/
	attachmentType := d.Get("type")
	if state, ok := d.GetOk("state"); ok {
		if attachmentType == "PARTNER_PROVIDER" && state == "ACTIVE" {
			return true
		}
	}

	return false
}
