package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
)

func expandAliasIpRanges(ranges []interface{}) []*computeBeta.AliasIpRange {
	ipRanges := make([]*computeBeta.AliasIpRange, 0, len(ranges))
	for _, raw := range ranges {
		data := raw.(map[string]interface{})
		ipRanges = append(ipRanges, &computeBeta.AliasIpRange{
			IpCidrRange:         data["ip_cidr_range"].(string),
			SubnetworkRangeName: data["subnetwork_range_name"].(string),
		})
	}
	return ipRanges
}

func expandAccessConfigs(configs []interface{}) []*computeBeta.AccessConfig {
	acs := make([]*computeBeta.AccessConfig, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})
		acs[i] = &computeBeta.AccessConfig{
			Type:        "ONE_TO_ONE_NAT",
			NatIP:       data["nat_ip"].(string),
			NetworkTier: data["network_tier"].(string),
		}
		if ptr, ok := data["public_ptr_domain_name"]; ok && ptr != "" {
			acs[i].SetPublicPtr = true
			acs[i].PublicPtrDomainName = ptr.(string)
		}
	}
	return acs
}

func expandNetworkInterfaces(d TerraformResourceData, config *Config) ([]*computeBeta.NetworkInterface, error) {
	configs := d.Get("network_interface").([]interface{})
	ifaces := make([]*computeBeta.NetworkInterface, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		network := data["network"].(string)
		subnetwork := data["subnetwork"].(string)

		// NOTE: Removed validation check on network and subnetwork.

		nf, err := ParseNetworkFieldValue(network, d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for network %q: %s", network, err)
		}

		subnetProjectField := fmt.Sprintf("network_interface.%d.subnetwork_project", i)
		sf, err := ParseSubnetworkFieldValueWithProjectField(subnetwork, subnetProjectField, d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for subnetwork %q: %s", subnetwork, err)
		}

		ifaces[i] = &computeBeta.NetworkInterface{
			NetworkIP:     data["network_ip"].(string),
			Network:       nf.RelativeLink(),
			Subnetwork:    sf.RelativeLink(),
			AccessConfigs: expandAccessConfigs(data["access_config"].([]interface{})),
			AliasIpRanges: expandAliasIpRanges(data["alias_ip_range"].([]interface{})),
		}

	}
	return ifaces, nil
}

func expandServiceAccounts(configs []interface{}) []*computeBeta.ServiceAccount {
	accounts := make([]*computeBeta.ServiceAccount, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		accounts[i] = &computeBeta.ServiceAccount{
			Email:  data["email"].(string),
			Scopes: canonicalizeServiceScopes(convertStringSet(data["scopes"].(*schema.Set))),
		}

		if accounts[i].Email == "" {
			accounts[i].Email = "default"
		}
	}
	return accounts
}

func resourceInstanceTags(d TerraformResourceData) *computeBeta.Tags {
	// Calculate the tags
	var tags *computeBeta.Tags
	if v := d.Get("tags"); v != nil {
		vs := v.(*schema.Set)
		tags = new(computeBeta.Tags)
		tags.Items = make([]string, vs.Len())
		for i, v := range vs.List() {
			tags.Items[i] = v.(string)
		}

		tags.Fingerprint = d.Get("tags_fingerprint").(string)
	}

	return tags
}
