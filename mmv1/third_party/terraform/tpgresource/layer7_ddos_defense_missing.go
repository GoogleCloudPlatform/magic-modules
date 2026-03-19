package tpgresource

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SuppressLayer7DdosDefenseMissing suppresses the diff when the layer_7_ddos_defense_config
// block is removed from the Terraform configuration, but the API still returns the block
// with enable set to false.
func SuppressLayer7DdosDefenseMissing(k, old, new string, d *schema.ResourceData) bool {
	// Only suppress when Terraform is trying to remove/omit the field or block
	if new != "" && new != "0" {
		return false
	}

	// We focus strictly on the status of layer_7_ddos_defense_config.
	// If it was already enabled in the state (API), we SHOULD NOT suppress it.
	enableVal, _ := d.GetChange("adaptive_protection_config.0.layer_7_ddos_defense_config.0.enable")
	if enable, ok := enableVal.(bool); ok && enable {
		return false
	}

	// Suppress diffs for the layer_7_ddos_defense_config field and any of its children.
	if strings.Contains(k, "layer_7_ddos_defense_config") {
		return true
	}

	// If evaluating the parent block (adaptive_protection_config.#), suppress it only if
	// no other features (like auto_deploy_config) are intentionally configured.
	if k == "adaptive_protection_config.#" {
		if _, ok := d.GetOk("adaptive_protection_config.0.auto_deploy_config"); !ok {
			return true
		}
	}

	return false
}
