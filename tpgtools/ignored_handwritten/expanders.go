package google

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func expandStringArray(v interface{}) []string {
	arr, ok := v.([]string)

	if ok {
		return arr
	}

	if arr, ok := v.(*schema.Set); ok {
		return convertStringSet(arr)
	}

	return convertStringArr(v.([]interface{}))
}

func expandComputeRouteNextHopInstance(v interface{}, d TerraformResourceData, config *Config) *string {
	if v == "" {
		return nil
	}
	val, err := parseZonalFieldValue("instances", v.(string), "project", "next_hop_instance_zone", d, config, true)
	if err != nil {
		return nil
	}
	nextInstance, err := config.clientCompute.Instances.Get(val.Project, val.Zone, val.Name).Do()
	if err != nil {
		return nil
	}
	return &nextInstance.SelfLink
}

func expandComputeVpnTunnelRegion(v interface{}, d TerraformResourceData, config *Config) *string {
	if v == "" {
		return nil
	}
	if reg, ok := v.(string); ok {
		return &reg
	}

	f, err := parseRegionalFieldValue("targetVpnGateways", d.Get("target_vpn_gateway").(string), "project", "region", "zone", d, config, true)
	if err != nil {
		return nil
	}
	return &f.Region
}
