package google

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func expandComputeRouteNextHopInstance(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) *string {
	if v == "" {
		return nil
	}
	val, err := tpgresource.ParseZonalFieldValue("instances", v.(string), "project", "next_hop_instance_zone", d, config, true)
	if err != nil {
		return nil
	}
	nextInstance, err := config.clientCompute.Instances.Get(val.Project, val.Zone, val.Name).Do()
	if err != nil {
		return nil
	}
	return &nextInstance.SelfLink
}

func expandComputeVpnTunnelRegion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) *string {
	if v == "" {
		return nil
	}
	if reg, ok := v.(string); ok {
		return &reg
	}

	f, err := tpgresource.ParseRegionalFieldValue("targetVpnGateways", d.Get("target_vpn_gateway").(string), "project", "region", "zone", d, config, true)
	if err != nil {
		return nil
	}
	return &f.Region
}
