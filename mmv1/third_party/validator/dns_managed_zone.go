package google

import (
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const DNSManagedZoneAssetType string = "dns.googleapis.com/ManagedZone"

func resourceConverterDNSManagedZone() ResourceConverter {
	return ResourceConverter{
		AssetType: DNSManagedZoneAssetType,
		Convert:   GetDNSManagedZoneCaiObject,
	}
}

func GetDNSManagedZoneCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//dns.googleapis.com/projects/{{project}}/managedZones/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetDNSManagedZoneApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: DNSManagedZoneAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/dns/v1/rest",
				DiscoveryName:        "ManagedZone",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetDNSManagedZoneApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	descriptionProp, err := expandDNSManagedZoneDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	dnsNameProp, err := expandDNSManagedZoneDnsName(d.Get("dns_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("dns_name"); !isEmptyValue(reflect.ValueOf(dnsNameProp)) && (ok || !reflect.DeepEqual(v, dnsNameProp)) {
		obj["dnsName"] = dnsNameProp
	}
	dnssecConfigProp, err := expandDNSManagedZoneDnssecConfig(d.Get("dnssec_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("dnssec_config"); !isEmptyValue(reflect.ValueOf(dnssecConfigProp)) && (ok || !reflect.DeepEqual(v, dnssecConfigProp)) {
		obj["dnssecConfig"] = dnssecConfigProp
	}
	nameProp, err := expandDNSManagedZoneName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	labelsProp, err := expandDNSManagedZoneLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	visibilityProp, err := expandDNSManagedZoneVisibility(d.Get("visibility"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("visibility"); !isEmptyValue(reflect.ValueOf(visibilityProp)) && (ok || !reflect.DeepEqual(v, visibilityProp)) {
		obj["visibility"] = visibilityProp
	}
	privateVisibilityConfigProp, err := expandDNSManagedZonePrivateVisibilityConfig(d.Get("private_visibility_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("private_visibility_config"); ok || !reflect.DeepEqual(v, privateVisibilityConfigProp) {
		obj["privateVisibilityConfig"] = privateVisibilityConfigProp
	}
	forwardingConfigProp, err := expandDNSManagedZoneForwardingConfig(d.Get("forwarding_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("forwarding_config"); !isEmptyValue(reflect.ValueOf(forwardingConfigProp)) && (ok || !reflect.DeepEqual(v, forwardingConfigProp)) {
		obj["forwardingConfig"] = forwardingConfigProp
	}
	peeringConfigProp, err := expandDNSManagedZonePeeringConfig(d.Get("peering_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("peering_config"); !isEmptyValue(reflect.ValueOf(peeringConfigProp)) && (ok || !reflect.DeepEqual(v, peeringConfigProp)) {
		obj["peeringConfig"] = peeringConfigProp
	}

	return obj, nil
}

func expandDNSManagedZoneDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZoneDnsName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZoneDnssecConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedKind, err := expandDNSManagedZoneDnssecConfigKind(original["kind"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedKind); val.IsValid() && !isEmptyValue(val) {
		transformed["kind"] = transformedKind
	}

	transformedNonExistence, err := expandDNSManagedZoneDnssecConfigNonExistence(original["non_existence"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNonExistence); val.IsValid() && !isEmptyValue(val) {
		transformed["nonExistence"] = transformedNonExistence
	}

	transformedState, err := expandDNSManagedZoneDnssecConfigState(original["state"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedState); val.IsValid() && !isEmptyValue(val) {
		transformed["state"] = transformedState
	}

	transformedDefaultKeySpecs, err := expandDNSManagedZoneDnssecConfigDefaultKeySpecs(original["default_key_specs"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDefaultKeySpecs); val.IsValid() && !isEmptyValue(val) {
		transformed["defaultKeySpecs"] = transformedDefaultKeySpecs
	}

	return transformed, nil
}

func expandDNSManagedZoneDnssecConfigKind(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZoneDnssecConfigNonExistence(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZoneDnssecConfigState(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZoneDnssecConfigDefaultKeySpecs(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedAlgorithm, err := expandDNSManagedZoneDnssecConfigDefaultKeySpecsAlgorithm(original["algorithm"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAlgorithm); val.IsValid() && !isEmptyValue(val) {
			transformed["algorithm"] = transformedAlgorithm
		}

		transformedKeyLength, err := expandDNSManagedZoneDnssecConfigDefaultKeySpecsKeyLength(original["key_length"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedKeyLength); val.IsValid() && !isEmptyValue(val) {
			transformed["keyLength"] = transformedKeyLength
		}

		transformedKeyType, err := expandDNSManagedZoneDnssecConfigDefaultKeySpecsKeyType(original["key_type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedKeyType); val.IsValid() && !isEmptyValue(val) {
			transformed["keyType"] = transformedKeyType
		}

		transformedKind, err := expandDNSManagedZoneDnssecConfigDefaultKeySpecsKind(original["kind"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedKind); val.IsValid() && !isEmptyValue(val) {
			transformed["kind"] = transformedKind
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDNSManagedZoneDnssecConfigDefaultKeySpecsAlgorithm(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZoneDnssecConfigDefaultKeySpecsKeyLength(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZoneDnssecConfigDefaultKeySpecsKeyType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZoneDnssecConfigDefaultKeySpecsKind(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZoneName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZoneLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandDNSManagedZoneVisibility(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZonePrivateVisibilityConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		// The API won't remove the the field unless an empty network array is sent.
		transformed := make(map[string]interface{})
		emptyNetwork := make([]interface{}, 0)
		transformed["networks"] = emptyNetwork
		return transformed, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedNetworks, err := expandDNSManagedZonePrivateVisibilityConfigNetworks(original["networks"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNetworks); val.IsValid() && !isEmptyValue(val) {
		transformed["networks"] = transformedNetworks
	}

	return transformed, nil
}

func expandDNSManagedZonePrivateVisibilityConfigNetworks(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedNetworkUrl, err := expandDNSManagedZonePrivateVisibilityConfigNetworksNetworkUrl(original["network_url"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNetworkUrl); val.IsValid() && !isEmptyValue(val) {
			transformed["networkUrl"] = transformedNetworkUrl
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDNSManagedZonePrivateVisibilityConfigNetworksNetworkUrl(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	if v == nil || v.(string) == "" {
		return "", nil
	} else if strings.HasPrefix(v.(string), "https://") {
		return v, nil
	}
	url, err := replaceVars(d, config, "{{ComputeBasePath}}"+v.(string))
	if err != nil {
		return "", err
	}
	return ConvertSelfLinkToV1(url), nil
}

func expandDNSManagedZoneForwardingConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedTargetNameServers, err := expandDNSManagedZoneForwardingConfigTargetNameServers(original["target_name_servers"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTargetNameServers); val.IsValid() && !isEmptyValue(val) {
		transformed["targetNameServers"] = transformedTargetNameServers
	}

	return transformed, nil
}

func expandDNSManagedZoneForwardingConfigTargetNameServers(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedIpv4Address, err := expandDNSManagedZoneForwardingConfigTargetNameServersIpv4Address(original["ipv4_address"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedIpv4Address); val.IsValid() && !isEmptyValue(val) {
			transformed["ipv4Address"] = transformedIpv4Address
		}

		transformedForwardingPath, err := expandDNSManagedZoneForwardingConfigTargetNameServersForwardingPath(original["forwarding_path"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedForwardingPath); val.IsValid() && !isEmptyValue(val) {
			transformed["forwardingPath"] = transformedForwardingPath
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDNSManagedZoneForwardingConfigTargetNameServersIpv4Address(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZoneForwardingConfigTargetNameServersForwardingPath(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDNSManagedZonePeeringConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedTargetNetwork, err := expandDNSManagedZonePeeringConfigTargetNetwork(original["target_network"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTargetNetwork); val.IsValid() && !isEmptyValue(val) {
		transformed["targetNetwork"] = transformedTargetNetwork
	}

	return transformed, nil
}

func expandDNSManagedZonePeeringConfigTargetNetwork(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedNetworkUrl, err := expandDNSManagedZonePeeringConfigTargetNetworkNetworkUrl(original["network_url"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNetworkUrl); val.IsValid() && !isEmptyValue(val) {
		transformed["networkUrl"] = transformedNetworkUrl
	}

	return transformed, nil
}

func expandDNSManagedZonePeeringConfigTargetNetworkNetworkUrl(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	if v == nil || v.(string) == "" {
		return "", nil
	} else if strings.HasPrefix(v.(string), "https://") {
		return v, nil
	}
	url, err := replaceVars(d, config, "{{ComputeBasePath}}"+v.(string))
	if err != nil {
		return "", err
	}
	return ConvertSelfLinkToV1(url), nil
}
