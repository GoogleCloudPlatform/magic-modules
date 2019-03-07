package google

import (
	"fmt"
	"reflect"
)

func GetComputeInstanceCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	if obj, err := GetComputeInstanceApiObject(d, config); err == nil {
		return Asset{
			Name: fmt.Sprintf("//compute.googleapis.com/%s", obj["selfLink"]),
			Type: "google.compute.Instance",
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Instance",
				Data:                 obj,
			},
		}, nil
	} else {
		return Asset{}, err
	}
}

func GetComputeInstanceApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	canIpForwardProp, err := expandComputeInstanceCanIpForward(d.Get("can_ip_forward"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("can_ip_forward"); !isEmptyValue(reflect.ValueOf(canIpForwardProp)) && (ok || !reflect.DeepEqual(v, canIpForwardProp)) {
		obj["canIpForward"] = canIpForwardProp
	}
	disksProp, err := expandComputeInstanceDisks(d.Get("disks"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("disks"); !isEmptyValue(reflect.ValueOf(disksProp)) && (ok || !reflect.DeepEqual(v, disksProp)) {
		obj["disks"] = disksProp
	}
	guestAcceleratorsProp, err := expandComputeInstanceGuestAccelerators(d.Get("guest_accelerators"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("guest_accelerators"); !isEmptyValue(reflect.ValueOf(guestAcceleratorsProp)) && (ok || !reflect.DeepEqual(v, guestAcceleratorsProp)) {
		obj["guestAccelerators"] = guestAcceleratorsProp
	}
	labelFingerprintProp, err := expandComputeInstanceLabelFingerprint(d.Get("label_fingerprint"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("label_fingerprint"); !isEmptyValue(reflect.ValueOf(labelFingerprintProp)) && (ok || !reflect.DeepEqual(v, labelFingerprintProp)) {
		obj["labelFingerprint"] = labelFingerprintProp
	}
	metadataProp, err := expandComputeInstanceMetadata(d.Get("metadata"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("metadata"); !isEmptyValue(reflect.ValueOf(metadataProp)) && (ok || !reflect.DeepEqual(v, metadataProp)) {
		obj["metadata"] = metadataProp
	}
	machineTypeProp, err := expandComputeInstanceMachineType(d.Get("machine_type"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("machine_type"); !isEmptyValue(reflect.ValueOf(machineTypeProp)) && (ok || !reflect.DeepEqual(v, machineTypeProp)) {
		obj["machineType"] = machineTypeProp
	}
	minCpuPlatformProp, err := expandComputeInstanceMinCpuPlatform(d.Get("min_cpu_platform"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("min_cpu_platform"); !isEmptyValue(reflect.ValueOf(minCpuPlatformProp)) && (ok || !reflect.DeepEqual(v, minCpuPlatformProp)) {
		obj["minCpuPlatform"] = minCpuPlatformProp
	}
	nameProp, err := expandComputeInstanceName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	networkInterfacesProp, err := expandComputeInstanceNetworkInterfaces(d.Get("network_interfaces"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("network_interfaces"); !isEmptyValue(reflect.ValueOf(networkInterfacesProp)) && (ok || !reflect.DeepEqual(v, networkInterfacesProp)) {
		obj["networkInterfaces"] = networkInterfacesProp
	}
	schedulingProp, err := expandComputeInstanceScheduling(d.Get("scheduling"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("scheduling"); !isEmptyValue(reflect.ValueOf(schedulingProp)) && (ok || !reflect.DeepEqual(v, schedulingProp)) {
		obj["scheduling"] = schedulingProp
	}
	serviceAccountsProp, err := expandComputeInstanceServiceAccounts(d.Get("service_accounts"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("service_accounts"); !isEmptyValue(reflect.ValueOf(serviceAccountsProp)) && (ok || !reflect.DeepEqual(v, serviceAccountsProp)) {
		obj["serviceAccounts"] = serviceAccountsProp
	}
	statusProp, err := expandComputeInstanceStatus(d.Get("status"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("status"); !isEmptyValue(reflect.ValueOf(statusProp)) && (ok || !reflect.DeepEqual(v, statusProp)) {
		obj["status"] = statusProp
	}
	tagsProp, err := expandComputeInstanceTags(d.Get("tags"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("tags"); !isEmptyValue(reflect.ValueOf(tagsProp)) && (ok || !reflect.DeepEqual(v, tagsProp)) {
		obj["tags"] = tagsProp
	}
	zoneProp, err := expandComputeInstanceZone(d.Get("zone"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("zone"); !isEmptyValue(reflect.ValueOf(zoneProp)) && (ok || !reflect.DeepEqual(v, zoneProp)) {
		obj["zone"] = zoneProp
	}

	return obj, nil
}

func expandComputeInstanceCanIpForward(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisks(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedAutoDelete, err := expandComputeInstanceDisksAutoDelete(original["auto_delete"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAutoDelete); val.IsValid() && !isEmptyValue(val) {
			transformed["autoDelete"] = transformedAutoDelete
		}

		transformedBoot, err := expandComputeInstanceDisksBoot(original["boot"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedBoot); val.IsValid() && !isEmptyValue(val) {
			transformed["boot"] = transformedBoot
		}

		transformedDeviceName, err := expandComputeInstanceDisksDeviceName(original["device_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDeviceName); val.IsValid() && !isEmptyValue(val) {
			transformed["deviceName"] = transformedDeviceName
		}

		transformedDiskEncryptionKey, err := expandComputeInstanceDisksDiskEncryptionKey(original["disk_encryption_key"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDiskEncryptionKey); val.IsValid() && !isEmptyValue(val) {
			transformed["diskEncryptionKey"] = transformedDiskEncryptionKey
		}

		transformedIndex, err := expandComputeInstanceDisksIndex(original["index"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedIndex); val.IsValid() && !isEmptyValue(val) {
			transformed["index"] = transformedIndex
		}

		transformedInitializeParams, err := expandComputeInstanceDisksInitializeParams(original["initialize_params"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedInitializeParams); val.IsValid() && !isEmptyValue(val) {
			transformed["initializeParams"] = transformedInitializeParams
		}

		transformedInterface, err := expandComputeInstanceDisksInterface(original["interface"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedInterface); val.IsValid() && !isEmptyValue(val) {
			transformed["interface"] = transformedInterface
		}

		transformedMode, err := expandComputeInstanceDisksMode(original["mode"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedMode); val.IsValid() && !isEmptyValue(val) {
			transformed["mode"] = transformedMode
		}

		transformedSource, err := expandComputeInstanceDisksSource(original["source"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSource); val.IsValid() && !isEmptyValue(val) {
			transformed["source"] = transformedSource
		}

		transformedType, err := expandComputeInstanceDisksType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !isEmptyValue(val) {
			transformed["type"] = transformedType
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeInstanceDisksAutoDelete(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksBoot(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksDeviceName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksDiskEncryptionKey(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedRawKey, err := expandComputeInstanceDisksDiskEncryptionKeyRawKey(original["raw_key"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRawKey); val.IsValid() && !isEmptyValue(val) {
		transformed["rawKey"] = transformedRawKey
	}

	transformedRsaEncryptedKey, err := expandComputeInstanceDisksDiskEncryptionKeyRsaEncryptedKey(original["rsa_encrypted_key"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRsaEncryptedKey); val.IsValid() && !isEmptyValue(val) {
		transformed["rsaEncryptedKey"] = transformedRsaEncryptedKey
	}

	transformedSha256, err := expandComputeInstanceDisksDiskEncryptionKeySha256(original["sha256"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSha256); val.IsValid() && !isEmptyValue(val) {
		transformed["sha256"] = transformedSha256
	}

	return transformed, nil
}

func expandComputeInstanceDisksDiskEncryptionKeyRawKey(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksDiskEncryptionKeyRsaEncryptedKey(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksDiskEncryptionKeySha256(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksIndex(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksInitializeParams(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDiskName, err := expandComputeInstanceDisksInitializeParamsDiskName(original["disk_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskName); val.IsValid() && !isEmptyValue(val) {
		transformed["diskName"] = transformedDiskName
	}

	transformedDiskSizeGb, err := expandComputeInstanceDisksInitializeParamsDiskSizeGb(original["disk_size_gb"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskSizeGb); val.IsValid() && !isEmptyValue(val) {
		transformed["diskSizeGb"] = transformedDiskSizeGb
	}

	transformedDiskType, err := expandComputeInstanceDisksInitializeParamsDiskType(original["disk_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskType); val.IsValid() && !isEmptyValue(val) {
		transformed["diskType"] = transformedDiskType
	}

	transformedSourceImage, err := expandComputeInstanceDisksInitializeParamsSourceImage(original["source_image"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSourceImage); val.IsValid() && !isEmptyValue(val) {
		transformed["sourceImage"] = transformedSourceImage
	}

	transformedSourceImageEncryptionKey, err := expandComputeInstanceDisksInitializeParamsSourceImageEncryptionKey(original["source_image_encryption_key"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSourceImageEncryptionKey); val.IsValid() && !isEmptyValue(val) {
		transformed["sourceImageEncryptionKey"] = transformedSourceImageEncryptionKey
	}

	return transformed, nil
}

func expandComputeInstanceDisksInitializeParamsDiskName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksInitializeParamsDiskSizeGb(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksInitializeParamsDiskType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseZonalFieldValue("diskTypes", v.(string), "project", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for disk_type: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeInstanceDisksInitializeParamsSourceImage(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksInitializeParamsSourceImageEncryptionKey(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedRawKey, err := expandComputeInstanceDisksInitializeParamsSourceImageEncryptionKeyRawKey(original["raw_key"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRawKey); val.IsValid() && !isEmptyValue(val) {
		transformed["rawKey"] = transformedRawKey
	}

	transformedSha256, err := expandComputeInstanceDisksInitializeParamsSourceImageEncryptionKeySha256(original["sha256"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSha256); val.IsValid() && !isEmptyValue(val) {
		transformed["sha256"] = transformedSha256
	}

	return transformed, nil
}

func expandComputeInstanceDisksInitializeParamsSourceImageEncryptionKeyRawKey(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksInitializeParamsSourceImageEncryptionKeySha256(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksInterface(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksMode(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceDisksSource(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseZonalFieldValue("disks", v.(string), "project", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for source: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeInstanceDisksType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGuestAccelerators(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedAcceleratorCount, err := expandComputeInstanceGuestAcceleratorsAcceleratorCount(original["accelerator_count"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAcceleratorCount); val.IsValid() && !isEmptyValue(val) {
			transformed["acceleratorCount"] = transformedAcceleratorCount
		}

		transformedAcceleratorType, err := expandComputeInstanceGuestAcceleratorsAcceleratorType(original["accelerator_type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAcceleratorType); val.IsValid() && !isEmptyValue(val) {
			transformed["acceleratorType"] = transformedAcceleratorType
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeInstanceGuestAcceleratorsAcceleratorCount(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceGuestAcceleratorsAcceleratorType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceLabelFingerprint(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceMetadata(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandComputeInstanceMachineType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseZonalFieldValue("machineTypes", v.(string), "project", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for machine_type: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeInstanceMinCpuPlatform(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceNetworkInterfaces(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedAccessConfigs, err := expandComputeInstanceNetworkInterfacesAccessConfigs(original["access_configs"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAccessConfigs); val.IsValid() && !isEmptyValue(val) {
			transformed["accessConfigs"] = transformedAccessConfigs
		}

		transformedAliasIpRanges, err := expandComputeInstanceNetworkInterfacesAliasIpRanges(original["alias_ip_ranges"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAliasIpRanges); val.IsValid() && !isEmptyValue(val) {
			transformed["aliasIpRanges"] = transformedAliasIpRanges
		}

		transformedName, err := expandComputeInstanceNetworkInterfacesName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !isEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedNetwork, err := expandComputeInstanceNetworkInterfacesNetwork(original["network"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNetwork); val.IsValid() && !isEmptyValue(val) {
			transformed["network"] = transformedNetwork
		}

		transformedNetworkIP, err := expandComputeInstanceNetworkInterfacesNetworkIP(original["network_ip"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNetworkIP); val.IsValid() && !isEmptyValue(val) {
			transformed["networkIP"] = transformedNetworkIP
		}

		transformedSubnetwork, err := expandComputeInstanceNetworkInterfacesSubnetwork(original["subnetwork"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubnetwork); val.IsValid() && !isEmptyValue(val) {
			transformed["subnetwork"] = transformedSubnetwork
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeInstanceNetworkInterfacesAccessConfigs(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedName, err := expandComputeInstanceNetworkInterfacesAccessConfigsName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !isEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedNatIP, err := expandComputeInstanceNetworkInterfacesAccessConfigsNatIP(original["nat_ip"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNatIP); val.IsValid() && !isEmptyValue(val) {
			transformed["natIP"] = transformedNatIP
		}

		transformedType, err := expandComputeInstanceNetworkInterfacesAccessConfigsType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !isEmptyValue(val) {
			transformed["type"] = transformedType
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeInstanceNetworkInterfacesAccessConfigsName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceNetworkInterfacesAccessConfigsNatIP(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseRegionalFieldValue("addresses", v.(string), "project", "region", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for nat_ip: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeInstanceNetworkInterfacesAccessConfigsType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceNetworkInterfacesAliasIpRanges(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedIpCidrRange, err := expandComputeInstanceNetworkInterfacesAliasIpRangesIpCidrRange(original["ip_cidr_range"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedIpCidrRange); val.IsValid() && !isEmptyValue(val) {
			transformed["ipCidrRange"] = transformedIpCidrRange
		}

		transformedSubnetworkRangeName, err := expandComputeInstanceNetworkInterfacesAliasIpRangesSubnetworkRangeName(original["subnetwork_range_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubnetworkRangeName); val.IsValid() && !isEmptyValue(val) {
			transformed["subnetworkRangeName"] = transformedSubnetworkRangeName
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeInstanceNetworkInterfacesAliasIpRangesIpCidrRange(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceNetworkInterfacesAliasIpRangesSubnetworkRangeName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceNetworkInterfacesName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceNetworkInterfacesNetwork(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("networks", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for network: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeInstanceNetworkInterfacesNetworkIP(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceNetworkInterfacesSubnetwork(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseRegionalFieldValue("subnetworks", v.(string), "project", "region", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for subnetwork: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeInstanceScheduling(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAutomaticRestart, err := expandComputeInstanceSchedulingAutomaticRestart(original["automatic_restart"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAutomaticRestart); val.IsValid() && !isEmptyValue(val) {
		transformed["automaticRestart"] = transformedAutomaticRestart
	}

	transformedOnHostMaintenance, err := expandComputeInstanceSchedulingOnHostMaintenance(original["on_host_maintenance"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOnHostMaintenance); val.IsValid() && !isEmptyValue(val) {
		transformed["onHostMaintenance"] = transformedOnHostMaintenance
	}

	transformedPreemptible, err := expandComputeInstanceSchedulingPreemptible(original["preemptible"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPreemptible); val.IsValid() && !isEmptyValue(val) {
		transformed["preemptible"] = transformedPreemptible
	}

	return transformed, nil
}

func expandComputeInstanceSchedulingAutomaticRestart(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceSchedulingOnHostMaintenance(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceSchedulingPreemptible(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceServiceAccounts(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedEmail, err := expandComputeInstanceServiceAccountsEmail(original["email"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedEmail); val.IsValid() && !isEmptyValue(val) {
			transformed["email"] = transformedEmail
		}

		transformedScopes, err := expandComputeInstanceServiceAccountsScopes(original["scopes"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedScopes); val.IsValid() && !isEmptyValue(val) {
			transformed["scopes"] = transformedScopes
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeInstanceServiceAccountsEmail(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceServiceAccountsScopes(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceStatus(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceTags(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedFingerprint, err := expandComputeInstanceTagsFingerprint(original["fingerprint"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFingerprint); val.IsValid() && !isEmptyValue(val) {
		transformed["fingerprint"] = transformedFingerprint
	}

	transformedItems, err := expandComputeInstanceTagsItems(original["items"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedItems); val.IsValid() && !isEmptyValue(val) {
		transformed["items"] = transformedItems
	}

	return transformed, nil
}

func expandComputeInstanceTagsFingerprint(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceTagsItems(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeInstanceZone(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("zones", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for zone: %s", err)
	}
	return f.RelativeLink(), nil
}
