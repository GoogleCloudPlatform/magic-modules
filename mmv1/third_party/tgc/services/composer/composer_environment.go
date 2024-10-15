package composer

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const ComposerEnvironmentAssetType string = "composer.googleapis.com/Environment"

func ResourceConverterComposerEnvironment() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: ComposerEnvironmentAssetType,
		Convert:   GetComposerEnvironmentCaiObject,
	}
}

func GetComposerEnvironmentCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//compute.googleapis.com/projects/{{project}}/locations/{{location}}/environments/{{name}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetComposerEnvironmentApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: ComposerEnvironmentAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://composer.googleapis.com/$discovery/rest?version=v1",
				DiscoveryName:        "Environment",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetComposerEnvironmentApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	nameProp, err := expandComputeEnvironmentName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	labelsProp, err := expandComputeEnvironmentLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	regionProp, err := expandComputeEnvironmentRegion(d.Get("region"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("region"); !tpgresource.IsEmptyValue(reflect.ValueOf(regionProp)) && (ok || !reflect.DeepEqual(v, regionProp)) {
		obj["region"] = regionProp
	}

	configProp, err := expandComputeEnvironmentConfig(d.Get("config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("config"); !tpgresource.IsEmptyValue(reflect.ValueOf(configProp)) && (ok || !reflect.DeepEqual(v, configProp)) {
		obj["config"] = configProp
	}

	storageConfigProp, err := expandComputeEnvironmentStorageConfig(d.Get("storage_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("storage_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(storageConfigProp)) && (ok || !reflect.DeepEqual(v, storageConfigProp)) {
		obj["storageConfig"] = storageConfigProp
	}

	return obj, nil
}

func expandComputeEnvironmentName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandComputeEnvironmentRegion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentStorageConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedNodeCount, err := expandComputeEnvironmentConfigNodeCount(original["node_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNodeCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["nodeCount"] = transformedNodeCount
	}

	transformedNodeConfig, err := expandComputeEnvironmentConfigNodeConfig(original["node_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNodeConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["nodeConfig"] = transformedNodeConfig
	}

	transformedRecoveryConfig, err := expandComputeEnvironmentConfigRecoveryConfig(original["recovery_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRecoveryConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["recoveryConfig"] = transformedRecoveryConfig
	}

	transformedSoftwareConfig, err := expandComputeEnvironmentConfigSoftwareConfig(original["software_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSoftwareConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["softwareConfig"] = transformedSoftwareConfig
	}

	transformedPrivateEnvironmentConfig, err := expandComputeEnvironmentConfigPrivateEnvironmentConfig(original["private_environment_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPrivateEnvironmentConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["privateEnvironmentConfig"] = transformedPrivateEnvironmentConfig
	}

	transformedWebServerNetworkAccessControl, err := expandComputeEnvironmentConfigWebServerNetworkAccessControl(original["web_server_network_access_control"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedWebServerNetworkAccessControl); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["webServerNetworkAccessControl"] = transformedWebServerNetworkAccessControl
	}

	transformedDatabaseConfig, err := expandComputeEnvironmentConfigDatabaseConfig(original["database_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDatabaseConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["databaseConfig"] = transformedDatabaseConfig
	}

	transformedWebServerConfig, err := expandComputeEnvironmentConfigWebServerConfig(original["web_server_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedWebServerConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["webServerConfig"] = transformedWebServerConfig
	}

	transformedEncryptionConfig, err := expandComputeEnvironmentConfigEncryptionConfig(original["encryption_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEncryptionConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["encryptionConfig"] = transformedEncryptionConfig
	}

	transformedMaintenanceWindow, err := expandComputeEnvironmentConfigMaintenanceWindow(original["maintenance_window"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMaintenanceWindow); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["maintenanceWindow"] = transformedMaintenanceWindow
	}

	transformedWorkloadsConfig, err := expandComputeEnvironmentConfigWorkloadsConfig(original["workloads_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedWorkloadsConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["workloadsConfig"] = transformedWorkloadsConfig
	}

	transformedDataRetentionConfig, err := expandComputeEnvironmentConfigDataRetentionConfig(original["data_retention_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDataRetentionConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["dataRetentionConfig"] = transformedDataRetentionConfig
	}

	transformedMasterAuthorizedNetworksConfig, err := expandComputeEnvironmentConfigMasterAuthorizedNetworksConfig(original["master_authorized_networks_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMasterAuthorizedNetworksConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["masterAuthorizedNetworksConfig"] = transformedMasterAuthorizedNetworksConfig
	}

	transformedResilienceMode, err := expandComputeEnvironmentConfigResilienceMode(original["resilience_mode"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedResilienceMode); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["resilienceMode"] = transformedResilienceMode
	}

	return transformed, nil
}

func expandComputeEnvironmentConfigNodeCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedZone, err := expandComputeEnvironmentConfigNodeConfigZone(original["zone"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedZone); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["zone"] = transformedZone
	}

	transformedMachineType, err := expandComputeEnvironmentConfigNodeConfigMachineType(original["machine_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMachineType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["machineType"] = transformedMachineType
	}

	transformedNetwork, err := expandComputeEnvironmentConfigNodeConfigNetwork(original["network"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNetwork); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["network"] = transformedNetwork
	}

	transformedSubnetwork, err := expandComputeEnvironmentConfigNodeConfigSubnetwork(original["subnetwork"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSubnetwork); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["subnetwork"] = transformedSubnetwork
	}

	transformedDiskSizeGb, err := expandComputeEnvironmentConfigNodeConfigDiskSizeGb(original["disk_size_gb"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskSizeGb); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskSizeGb"] = transformedDiskSizeGb
	}

	transformedServiceAccount, err := expandComputeEnvironmentConfigNodeConfigServiceAccount(original["service_account"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedServiceAccount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["serviceAccount"] = transformedServiceAccount
	}

	transformedIpAllocationPolicy, err := expandComputeEnvironmentConfigNodeConfigIpAllocationPolicy(original["ip_allocation_policy"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIpAllocationPolicy); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["ipAllocationPolicy"] = transformedIpAllocationPolicy
	}

	transformedOauthScopes, err := expandComputeEnvironmentConfigNodeConfigOauthScopes(original["oauth_scopes"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOauthScopes); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["oauthScopes"] = transformedOauthScopes
	}

	transformedMaxPodsPerNode, err := expandComputeEnvironmentConfigNodeConfigMaxPodsPerNode(original["max_pods_per_node"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMaxPodsPerNode); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["maxPodsPerNode"] = transformedMaxPodsPerNode
	}

	transformedEnableIpMasqAgent, err := expandComputeEnvironmentConfigNodeConfigEnableIpMasqAgent(original["enable_ip_masq_agent"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnableIpMasqAgent); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enableIpMasqAgent"] = transformedEnableIpMasqAgent
	}

	transformedTags, err := expandComputeEnvironmentConfigNodeConfigTags(original["tags"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTags); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["tags"] = transformedTags
	}

	return transformed, nil
}

func expandComputeEnvironmentConfigNodeConfigOauthScopes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfigTags(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfigZone(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfigMachineType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfigNetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfigSubnetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfigServiceAccount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfigIpAllocationPolicy(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfigEnableIpMasqAgent(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfigMaxPodsPerNode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfigDiskSizeGb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigRecoveryConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigSoftwareConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigPrivateEnvironmentConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigWebServerNetworkAccessControl(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigDatabaseConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigWebServerConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigEncryptionConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigMaintenanceWindow(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigWorkloadsConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigDataRetentionConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigMasterAuthorizedNetworksConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEnabled, err := expandComputeEnvironmentConfigNodeConfigEnabled(original["enabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enabled"] = transformedEnabled
	}

	transformedCidrBlocks, err := expandComputeEnvironmentConfigNodeConfigCidrBlocks(original["cidr_blocks"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCidrBlocks); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["cidrBlocks"] = transformedCidrBlocks
	}

	return transformed, nil
}

func expandComputeEnvironmentConfigNodeConfigEnabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandComputeEnvironmentConfigNodeConfigCidrBlocks(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandComputeEnvironmentConfigResilienceMode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
