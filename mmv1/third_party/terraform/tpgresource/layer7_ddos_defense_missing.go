package tpgresource

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SuppressLayer7DdosDefenseMissing suppresses the diff when the layer_7_ddos_defense_config
// or its parent adaptive_protection_config block is removed from the Terraform configuration,
// but the API still returns the block with enable set to false.
func SuppressLayer7DdosDefenseMissing(k, old, new string, d *schema.ResourceData) bool {
	if !strings.HasPrefix(k, "adaptive_protection_config") {
		return false
	}

	if new != "" && new != "0" {
		return false
	}

	enable := d.Get("adaptive_protection_config.0.layer_7_ddos_defense_config.0.enable")

	if enable != nil && enable.(bool) == true {
		return false
	}

	if strings.Contains(k, "layer_7_ddos_defense_config") {
		return true
	}

	if k == "adaptive_protection_config.#" {
		autoDeployCount := d.Get("adaptive_protection_config.0.auto_deploy_config.#")
		if autoDeployCount == nil || autoDeployCount.(int) == 0 {
			return true
		}
	}

	return false
}
