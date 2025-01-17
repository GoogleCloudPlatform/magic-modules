package dataproc

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const DataprocClusterAssetType string = "dataproc.googleapis.com/Cluster"

func ResourceConverterDataprocCluster() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: DataprocClusterAssetType,
		Convert:   GetDataprocClusterCaiObject,
	}
}

func GetDataprocClusterCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//compute.googleapis.com/projects/{{project}}/regions/{{region}}/clusters/{{name}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetDataprocClusterApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: DataprocClusterAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://dataproc.googleapis.com/$discovery/rest?version=v1",
				DiscoveryName:        "Cluster",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetDataprocClusterApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	projectIdProp, err := expandDataprocClusterProjectId(d.Get("project"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("project"); !tpgresource.IsEmptyValue(reflect.ValueOf(projectIdProp)) && (ok || !reflect.DeepEqual(v, projectIdProp)) {
		obj["projectId"] = projectIdProp
	}

	clusterNameProp, err := expandDataprocClusterName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(clusterNameProp)) && (ok || !reflect.DeepEqual(v, clusterNameProp)) {
		obj["clusterName"] = clusterNameProp
	}

	configProp, err := expandDataprocClusterConfig(d.Get("cluster_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("cluster_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(configProp)) && (ok || !reflect.DeepEqual(v, configProp)) {
		obj["config"] = configProp
	}

	virtualClusterConfigProp, err := expandDataprocVirtualClusterConfig(d.Get("virtual_cluster_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("virtual_cluster_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(virtualClusterConfigProp)) && (ok || !reflect.DeepEqual(v, virtualClusterConfigProp)) {
		obj["virtualClusterConfig"] = virtualClusterConfigProp
	}

	labelsProp, err := expandDataprocClusterLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	return obj, nil
}

func expandDataprocClusterProjectId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedConfigBucket, err := expandDataprocClusterConfigBucket(original["staging_bucket"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedConfigBucket); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["configBucket"] = transformedConfigBucket
	}

	transformedTempBucket, err := expandDataprocClusterTempBucket(original["temp_bucket"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTempBucket); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["tempBucket"] = transformedTempBucket
	}

	transformedGceClusterConfig, err := expandDataprocClusterConfigGceClusterConfig(original["gce_cluster_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["gceClusterConfig"] = transformedGceClusterConfig
	}

	transformedMasterConfig, err := expandDataprocClusterConfigMasterConfig(original["master_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["masterConfig"] = transformedMasterConfig
	}

	transformedWorkerConfig, err := expandDataprocClusterConfigWorkerConfig(original["worker_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["workerConfig"] = transformedWorkerConfig
	}

	transformedSecondaryWorkerConfig, err := expandDataprocClusterConfigSecondaryWorkerConfig(original["preemptible_worker_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["secondaryWorkerConfig"] = transformedSecondaryWorkerConfig
	}

	transformedSoftwareConfig, err := expandDataprocClusterConfigSoftwareConfig(original["software_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["softwareConfig"] = transformedSoftwareConfig
	}

	transformedSecurityConfig, err := expandDataprocClusterConfigSecurityConfig(original["security_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["securityConfig"] = transformedSecurityConfig
	}

	transformedAutoscalingConfig, err := expandDataprocClusterConfigAutoscalingConfig(original["autoscaling_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["autoscalingConfig"] = transformedAutoscalingConfig
	}

	transformedNodeInitializationAction, err := expandDataprocClusterConfigNodeInitializationAction(original["initialization_action"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["initializationActions"] = transformedNodeInitializationAction
	}

	transformedEncryptionConfig, err := expandDataprocClusterConfigEncryptionConfig(original["encryption_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["encryptionConfig"] = transformedEncryptionConfig
	}

	transformedLifecycleConfig, err := expandDataprocClusterConfigLifecycleConfig(original["lifecycle_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["lifecycleConfig"] = transformedLifecycleConfig
	}

	transformedEndpointConfig, err := expandDataprocClusterConfigEndpointConfig(original["endpoint_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["endpointConfig"] = transformedEndpointConfig
	}

	transformedDataprocMetricConfig, err := expandDataprocClusterConfigDataprocMetricConfig(original["dataproc_metric_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["dataprocMetricConfig"] = transformedDataprocMetricConfig
	}

	transformedAuxiliaryNodeGroups, err := expandDataprocClusterConfigAuxiliaryNodeGroups(original["auxiliary_node_groups"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["auxiliaryNodeGroups"] = transformedAuxiliaryNodeGroups
	}

	transformedMetastoreConfig, err := expandDataprocClusterConfigMetastoreConfig(original["metastore_config"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["metastoreConfig"] = transformedMetastoreConfig
	}

	return transformed, nil
}

func expandDataprocClusterConfigBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterTempBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedZone, err := expandDataprocClusterConfigGceClusterConfigZone(original["zone"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedZone); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["zoneUri"] = transformedZone
	}

	transformedNetwork, err := expandDataprocClusterConfigGceClusterConfigNetwork(original["network"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNetwork); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["networkUri"] = transformedNetwork
	}

	transformedSubnetwork, err := expandDataprocClusterConfigGceClusterConfigSubnetwork(original["subnetwork"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSubnetwork); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["subnetworkUri"] = transformedSubnetwork
	}

	transformedServiceAccount, err := expandDataprocClusterConfigGceClusterConfigServiceAccount(original["service_account"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedServiceAccount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["serviceAccount"] = transformedServiceAccount
	}

	transformedServiceAccountScopes, err := expandDataprocClusterConfigGceClusterConfigServiceAccountScopes(original["service_account_scopes"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedServiceAccountScopes); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["serviceAccountScopes"] = transformedServiceAccountScopes
	}

	transformedTags, err := expandDataprocClusterConfigGceClusterConfigTags(original["tags"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTags); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["tags"] = transformedTags
	}

	transformedInternalIpOnly, err := expandDataprocClusterConfigGceClusterConfigInternalIpOnly(original["internal_ip_only"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedInternalIpOnly); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["internalIpOnly"] = transformedInternalIpOnly
	}

	transformedMetadata, err := expandDataprocClusterConfigGceClusterConfigMetadata(original["metadata"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMetadata); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["metadata"] = transformedMetadata
	}

	transformedReservationAffinity, err := expandDataprocClusterConfigGceClusterConfigReservationAffinity(original["reservation_affinity"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedReservationAffinity); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["reservationAffinity"] = transformedReservationAffinity
	}

	transformedNodeGroupAffinity, err := expandDataprocClusterConfigGceClusterConfigNodeGroupAffinity(original["node_group_affinity"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNodeGroupAffinity); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["nodeGroupAffinity"] = transformedNodeGroupAffinity
	}

	transformedShieldedInstanceConfig, err := expandDataprocClusterConfigGceClusterConfigShieldedInstanceConfig(original["shielded_instance_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedShieldedInstanceConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["shieldedInstanceConfig"] = transformedShieldedInstanceConfig
	}

	return transformed, nil
}

func expandDataprocClusterConfigGceClusterConfigZone(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigNetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigSubnetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigServiceAccount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigServiceAccountScopes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigTags(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigInternalIpOnly(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigMetadata(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandDataprocClusterConfigGceClusterConfigReservationAffinity(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedConsumeReservationType, err := expandDataprocClusterConfigGceClusterConfigReservationAffinityConsumeReservationType(original["consume_reservation_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedConsumeReservationType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["consumeReservationType"] = transformedConsumeReservationType
	}

	transformedKey, err := expandDataprocClusterConfigGceClusterConfigReservationAffinityKey(original["key"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedKey); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["key"] = transformedKey
	}

	transformedValues, err := expandDataprocClusterConfigGceClusterConfigReservationAffinityValues(original["values"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedValues); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["values"] = transformedValues
	}

	return transformed, nil
}

func expandDataprocClusterConfigGceClusterConfigReservationAffinityConsumeReservationType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigReservationAffinityKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigReservationAffinityValues(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigNodeGroupAffinity(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedNodeGroupUri, err := expandDataprocClusterConfigGceClusterConfigNodeGroupAffinityNodeGroupUri(original["node_group_uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNodeGroupUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["nodeGroupUri"] = transformedNodeGroupUri
	}

	return transformed, nil
}

func expandDataprocClusterConfigGceClusterConfigNodeGroupAffinityNodeGroupUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigShieldedInstanceConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEnableSecureBoot, err := expandDataprocClusterConfigGceClusterConfigShieldedInstanceConfigEnableSecureBoot(original["enable_secure_boot"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnableSecureBoot); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enableSecureBoot"] = transformedEnableSecureBoot
	}

	transformedEnableVtpm, err := expandDataprocClusterConfigGceClusterConfigShieldedInstanceConfigEnableVtpm(original["enable_vtpm"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnableVtpm); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enableVtpm"] = transformedEnableVtpm
	}

	transformedEnableIntegrityMonitoring, err := expandDataprocClusterConfigGceClusterConfigShieldedInstanceConfigEnableIntegrityMonitoring(original["enable_integrity_monitoring"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnableIntegrityMonitoring); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enableIntegrityMonitoring"] = transformedEnableIntegrityMonitoring
	}

	return transformed, nil
}

func expandDataprocClusterConfigGceClusterConfigShieldedInstanceConfigEnableSecureBoot(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigShieldedInstanceConfigEnableVtpm(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigGceClusterConfigShieldedInstanceConfigEnableIntegrityMonitoring(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigMasterConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedNumInstances, err := expandDataprocClusterConfigMasterConfigNumInstances(original["num_instances"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNumInstances); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["numInstances"] = transformedNumInstances
	}

	transformedMachineType, err := expandDataprocClusterConfigMasterConfigMachineType(original["machine_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMachineType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["machineType"] = transformedMachineType
	}

	transformedMinCpuPlatform, err := expandDataprocClusterConfigMasterConfigMinCpuPlatform(original["min_cpu_platform"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinCpuPlatform); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minCpuPlatform"] = transformedMinCpuPlatform
	}

	transformedImageUri, err := expandDataprocClusterConfigMasterConfigImageUri(original["image_uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedImageUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["imageUri"] = transformedImageUri
	}

	transformedDiskConfig, err := expandDataprocClusterConfigMasterConfigDiskConfig(original["disk_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskConfig"] = transformedDiskConfig
	}

	transformedAcceleratorConfig, err := expandDataprocClusterConfigMasterConfigAcceleratorConfig(original["accelerators"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAcceleratorConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["accelerators"] = transformedAcceleratorConfig
	}

	return transformed, nil
}

func expandDataprocClusterConfigMasterConfigNumInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigMasterConfigMachineType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigMasterConfigMinCpuPlatform(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigMasterConfigImageUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigMasterConfigDiskConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBootDiskType, err := expandDataprocClusterConfigMasterConfigDiskConfigBootDiskType(original["boot_disk_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBootDiskType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bootDiskType"] = transformedBootDiskType
	}

	transformedBootDiskSizeGb, err := expandDataprocClusterConfigMasterConfigDiskConfigBootDiskSizeGb(original["boot_disk_size_gb"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBootDiskSizeGb); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bootDiskSizeGb"] = transformedBootDiskSizeGb
	}

	transformedNumLocalSsds, err := expandDataprocClusterConfigMasterConfigDiskConfigNumLocalSsds(original["num_local_ssds"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNumLocalSsds); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["numLocalSsds"] = transformedNumLocalSsds
	}

	transformedLocalSsdInterface, err := expandDataprocClusterConfigMasterConfigDiskConfigLocalSsdInterface(original["local_ssd_interface"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocalSsdInterface); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["localSsdInterface"] = transformedLocalSsdInterface
	}

	return transformed, nil
}

func expandDataprocClusterConfigMasterConfigDiskConfigBootDiskType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigMasterConfigDiskConfigBootDiskSizeGb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigMasterConfigDiskConfigNumLocalSsds(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigMasterConfigDiskConfigLocalSsdInterface(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigMasterConfigAcceleratorConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()

	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedAcceleratorTypeUri, err := expandDataprocClusterConfigMasterConfigAcceleratorConfigAcceleratorTypeUri(original["accelerator_type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAcceleratorTypeUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["acceleratorTypeUri"] = transformedAcceleratorTypeUri
		}

		transformedAcceleratorCount, err := expandDataprocClusterConfigMasterConfigAcceleratorConfigAcceleratorCount(original["accelerator_count"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAcceleratorCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["acceleratorCount"] = transformedAcceleratorCount
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandDataprocClusterConfigMasterConfigAcceleratorConfigAcceleratorTypeUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigMasterConfigAcceleratorConfigAcceleratorCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigWorkerConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedNumInstances, err := expandDataprocClusterConfigWorkerConfigNumInstances(original["num_instances"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNumInstances); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["numInstances"] = transformedNumInstances
	}

	transformedMachineType, err := expandDataprocClusterConfigWorkerConfigMachineType(original["machine_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMachineType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["machineType"] = transformedMachineType
	}

	transformedMinCpuPlatform, err := expandDataprocClusterConfigWorkerConfigMinCpuPlatform(original["min_cpu_platform"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinCpuPlatform); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minCpuPlatform"] = transformedMinCpuPlatform
	}

	transformedMinNumInstances, err := expandDataprocClusterConfigWorkerConfigMinNumInstances(original["min_num_instances"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinNumInstances); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minNumInstances"] = transformedMinNumInstances
	}

	transformedImageUri, err := expandDataprocClusterConfigWorkerConfigImageUri(original["image_uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedImageUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["imageUri"] = transformedImageUri
	}

	transformedDiskConfig, err := expandDataprocClusterConfigWorkerConfigDiskConfig(original["disk_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskConfig"] = transformedDiskConfig
	}

	transformedAcceleratorConfig, err := expandDataprocClusterConfigWorkerConfigAcceleratorConfig(original["accelerators"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAcceleratorConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["accelerators"] = transformedAcceleratorConfig
	}

	return transformed, nil
}

func expandDataprocClusterConfigWorkerConfigNumInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigWorkerConfigMachineType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigWorkerConfigMinCpuPlatform(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigWorkerConfigMinNumInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigWorkerConfigImageUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigWorkerConfigDiskConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBootDiskType, err := expandDataprocClusterConfigWorkerConfigDiskConfigBootDiskType(original["boot_disk_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBootDiskType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bootDiskType"] = transformedBootDiskType
	}

	transformedBootDiskSizeGb, err := expandDataprocClusterConfigWorkerConfigDiskConfigBootDiskSizeGb(original["boot_disk_size_gb"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBootDiskSizeGb); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bootDiskSizeGb"] = transformedBootDiskSizeGb
	}

	transformedNumLocalSsds, err := expandDataprocClusterConfigWorkerConfigDiskConfigNumLocalSsds(original["num_local_ssds"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNumLocalSsds); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["numLocalSsds"] = transformedNumLocalSsds
	}

	return transformed, nil
}

func expandDataprocClusterConfigWorkerConfigDiskConfigBootDiskType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigWorkerConfigDiskConfigBootDiskSizeGb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigWorkerConfigDiskConfigNumLocalSsds(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigWorkerConfigAcceleratorConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()

	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedAcceleratorTypeUri, err := expandDataprocClusterConfigWorkerConfigAcceleratorConfigAcceleratorTypeUri(original["accelerator_type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAcceleratorTypeUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["acceleratorTypeUri"] = transformedAcceleratorTypeUri
		}

		transformedAcceleratorCount, err := expandDataprocClusterConfigWorkerConfigAcceleratorConfigAcceleratorCount(original["accelerator_count"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAcceleratorCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["acceleratorCount"] = transformedAcceleratorCount
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandDataprocClusterConfigWorkerConfigAcceleratorConfigAcceleratorTypeUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigWorkerConfigAcceleratorConfigAcceleratorCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigSecondaryWorkerConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedNumInstances, err := expandDataprocClusterConfigSecondaryWorkerConfigNumInstances(original["num_instances"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNumInstances); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["numInstances"] = transformedNumInstances
	}

	transformedPreemptibility, err := expandDataprocClusterConfigSecondaryWorkerConfigPreemptibility(original["preemptibility"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPreemptibility); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["preemptibility"] = transformedPreemptibility
	}

	transformedDiskConfig, err := expandDataprocClusterConfigSecondaryWorkerConfigDiskConfig(original["disk_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskConfig"] = transformedDiskConfig
	}

	transformedInstanceFlexibilityPolicy, err := expandDataprocClusterConfigSecondaryWorkerConfigInstanceFlexibilityPolicy(original["instance_flexibility_policy"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedInstanceFlexibilityPolicy); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["instanceFlexibilityPolicy"] = transformedInstanceFlexibilityPolicy
	}

	return transformed, nil
}

func expandDataprocClusterConfigSecondaryWorkerConfigNumInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigSecondaryWorkerConfigPreemptibility(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigSecondaryWorkerConfigDiskConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBootDiskType, err := expandDataprocClusterConfigSecondaryWorkerConfigDiskConfigBootDiskType(original["boot_disk_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBootDiskType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bootDiskType"] = transformedBootDiskType
	}

	transformedBootDiskSizeGb, err := expandDataprocClusterConfigSecondaryWorkerConfigDiskConfigBootDiskSizeGb(original["boot_disk_size_gb"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBootDiskSizeGb); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bootDiskSizeGb"] = transformedBootDiskSizeGb
	}

	transformedNumLocalSsds, err := expandDataprocClusterConfigSecondaryWorkerConfigDiskConfigNumLocalSsds(original["num_local_ssds"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNumLocalSsds); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["numLocalSsds"] = transformedNumLocalSsds
	}

	return transformed, nil
}

func expandDataprocClusterConfigSecondaryWorkerConfigDiskConfigBootDiskType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigSecondaryWorkerConfigDiskConfigBootDiskSizeGb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigSecondaryWorkerConfigDiskConfigNumLocalSsds(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigSecondaryWorkerConfigInstanceFlexibilityPolicy(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigSoftwareConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedImageVersion, err := expandDataprocClusterConfigSoftwareConfigImageVersion(original["image_version"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedImageVersion); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["imageVersion"] = transformedImageVersion
	}

	transformedProperties, err := expandDataprocClusterConfigSoftwareConfigProperties(original["override_properties"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedProperties); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["properties"] = transformedProperties
	}

	transformedOptionalComponents, err := expandDataprocClusterConfigSoftwareConfigOptionalComponents(original["optional_components"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOptionalComponents); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["optionalComponents"] = transformedOptionalComponents
	}

	return transformed, nil
}

func expandDataprocClusterConfigSoftwareConfigImageVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigSoftwareConfigProperties(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandDataprocClusterConfigSoftwareConfigOptionalComponents(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandDataprocClusterConfigSecurityConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigAutoscalingConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigNodeInitializationAction(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigEncryptionConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigLifecycleConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigEndpointConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigDataprocMetricConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMetrics, err := expandDataprocClusterConfigDataprocMetricConfigMetrics(original["metrics"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMetrics); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["metrics"] = transformedMetrics
	}

	return transformed, nil
}

func expandDataprocClusterConfigDataprocMetricConfigMetrics(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedMetricSource, err := expandDataprocClusterConfigDataprocMetricConfigMetricsMetricSource(original["metric_source"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedMetricSource); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["metricSource"] = transformedMetricSource
		}

		transformedMetricOverrides, err := expandDataprocClusterConfigDataprocMetricConfigMetricsMetricOverrides(original["metric_overrides"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedMetricOverrides); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["metricOverrides"] = transformedMetricOverrides
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandDataprocClusterConfigDataprocMetricConfigMetricsMetricSource(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigDataprocMetricConfigMetricsMetricOverrides(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandDataprocClusterConfigAuxiliaryNodeGroups(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedNodeGroup, err := expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroup(original["node_group"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNodeGroup); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["nodeGroup"] = transformedNodeGroup
	}

	return transformed, nil
}

func expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroup(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedName, err := expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupName(original["name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["name"] = transformedName
	}

	transformedRoles, err := expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupRoles(original["roles"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRoles); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["roles"] = transformedRoles
	}

	transformedNodeGroupConfig, err := expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfig(original["node_group_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNodeGroupConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["nodeGroupConfig"] = transformedNodeGroupConfig
	}

	return transformed, nil
}

func expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupRoles(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedNumInstances, err := expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfigNumInstances(original["num_instances"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNumInstances); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["numInstances"] = transformedNumInstances
	}

	transformedMachineType, err := expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfigMachineType(original["machine_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMachineType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["machineType"] = transformedMachineType
	}

	transformedMinCpuPlatform, err := expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfigMinCpuPlatform(original["min_cpu_platform"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinCpuPlatform); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minCpuPlatform"] = transformedMinCpuPlatform
	}

	transformedDiskConfig, err := expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfigDiskConfig(original["disk_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskConfig"] = transformedDiskConfig
	}

	transformedAccelerators, err := expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfigAccelerators(original["accelerators"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAccelerators); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["accelerators"] = transformedAccelerators
	}

	return transformed, nil
}

func expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfigNumInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfigMachineType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfigMinCpuPlatform(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfigDiskConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfigAuxiliaryNodeGroupsNodeGroupNodeGroupConfigAccelerators(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandDataprocClusterConfigMetastoreConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandDataprocVirtualClusterConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
