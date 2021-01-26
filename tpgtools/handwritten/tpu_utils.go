package google

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// compareTpuNodeSchedulingConfig diff suppresses for the default
// scheduling, i.e. if preemptible is false, the API may either return no
// schedulingConfig or an empty schedulingConfig.
func compareTpuNodeSchedulingConfig(k, old, new string, d *schema.ResourceData) bool {
	if k == "scheduling_config.0.preemptible" {
		return old == "" && new == "false"
	}
	if k == "scheduling_config.#" {
		o, n := d.GetChange("scheduling_config.0.preemptible")
		return o.(bool) == n.(bool)
	}
	return false
}
