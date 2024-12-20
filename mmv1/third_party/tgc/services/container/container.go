// This file is written manually and is not based from terraform-provider-google.
// There is a huge potential for drift. The longer term plan is to have this
// file generated from the logic in terraform-provider-google. Please
// see https://github.com/GoogleCloudPlatform/magic-modules/pull/2485#issuecomment-545680059
// for the discussion.

package container

import (
	"fmt"
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const ContainerClusterAssetType string = "container.googleapis.com/Cluster"
const ContainerNodePoolAssetType string = "container.googleapis.com/NodePool"

func ResourceConverterContainerCluster() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: ContainerClusterAssetType,
		Convert:   GetContainerClusterCaiObject,
	}
}

func ResourceConverterContainerNodePool() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: ContainerNodePoolAssetType,
		Convert:   GetContainerNodePoolCaiObject,
	}
}

func expandContainerEnabledObject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if val := reflect.ValueOf(v); !val.IsValid() || tpgresource.IsEmptyValue(val) {
		return nil, nil
	}
	transformed := map[string]interface{}{
		"enabled": v,
	}
	return transformed, nil
}

func expandContainerClusterEnableLegacyAbac(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return expandContainerEnabledObject(v, d, config)
}

func expandContainerMaxPodsConstraint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if val := reflect.ValueOf(v); !val.IsValid() || tpgresource.IsEmptyValue(val) {
		return nil, nil
	}
	transformed := map[string]interface{}{
		"maxPodsPerNode": v,
	}
	return transformed, nil
}

func expandContainerClusterDefaultMaxPodsPerNode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return expandContainerMaxPodsConstraint(v, d, config)
}

func expandContainerNodePoolMaxPodsPerNode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return expandContainerMaxPodsConstraint(v, d, config)
}

func expandContainerClusterNetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	fv, err := tpgresource.ParseNetworkFieldValue(v.(string), d, config)
	if err != nil {
		return nil, err
	}
	return fv.RelativeLink(), nil
}

func expandContainerClusterSubnetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	fv, err := tpgresource.ParseNetworkFieldValue(v.(string), d, config)
	if err != nil {
		return nil, err
	}
	return fv.RelativeLink(), nil
}

func canonicalizeServiceScopesFromSet(scopesSet *schema.Set) (interface{}, error) {
	scopes := make([]string, scopesSet.Len())
	for i, scope := range scopesSet.List() {
		scopes[i] = tpgresource.CanonicalizeServiceScope(scope.(string))
	}
	return scopes, nil
}

func expandContainerClusterNodeConfigOauthScopes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	scopesSet := v.(*schema.Set)
	return canonicalizeServiceScopesFromSet(scopesSet)
}

func expandContainerNodePoolNodeConfigOauthScopes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	scopesSet := v.(*schema.Set)
	return canonicalizeServiceScopesFromSet(scopesSet)
}

func GetContainerClusterCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//container.googleapis.com/projects/{{project}}/locations/{{location}}/clusters/{{name}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetContainerClusterApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: ContainerClusterAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/container/v1/rest",
				DiscoveryName:        "Cluster",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetContainerClusterApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	enableKubernetesAlphaProp, err := expandContainerClusterEnableKubernetesAlpha(d.Get("enable_kubernetes_alpha"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("enable_kubernetes_alpha"); !tpgresource.IsEmptyValue(reflect.ValueOf(enableKubernetesAlphaProp)) && (ok || !reflect.DeepEqual(v, enableKubernetesAlphaProp)) {
		obj["enableKubernetesAlpha"] = enableKubernetesAlphaProp
	}

	nameProp, err := expandContainerClusterName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	descriptionProp, err := expandContainerClusterDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	initialNodeCountProp, err := expandContainerClusterInitialNodeCount(d.Get("initial_node_count"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("initial_node_count"); !tpgresource.IsEmptyValue(reflect.ValueOf(initialNodeCountProp)) && (ok || !reflect.DeepEqual(v, initialNodeCountProp)) {
		obj["initialNodeCount"] = initialNodeCountProp
	}
	nodeConfigProp, err := expandContainerClusterNodeConfig(d.Get("node_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("node_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(nodeConfigProp)) && (ok || !reflect.DeepEqual(v, nodeConfigProp)) {
		obj["nodeConfig"] = nodeConfigProp
	}
	masterAuthProp, err := expandContainerClusterMasterAuth(d.Get("master_auth"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("master_auth"); !tpgresource.IsEmptyValue(reflect.ValueOf(masterAuthProp)) && (ok || !reflect.DeepEqual(v, masterAuthProp)) {
		obj["masterAuth"] = masterAuthProp
	}
	loggingServiceProp, err := expandContainerClusterLoggingService(d.Get("logging_service"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("logging_service"); !tpgresource.IsEmptyValue(reflect.ValueOf(loggingServiceProp)) && (ok || !reflect.DeepEqual(v, loggingServiceProp)) {
		obj["loggingService"] = loggingServiceProp
	}
	monitoringServiceProp, err := expandContainerClusterMonitoringService(d.Get("monitoring_service"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("monitoring_service"); !tpgresource.IsEmptyValue(reflect.ValueOf(monitoringServiceProp)) && (ok || !reflect.DeepEqual(v, monitoringServiceProp)) {
		obj["monitoringService"] = monitoringServiceProp
	}
	networkProp, err := expandContainerClusterNetwork(d.Get("network"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("network"); !tpgresource.IsEmptyValue(reflect.ValueOf(networkProp)) && (ok || !reflect.DeepEqual(v, networkProp)) {
		obj["network"] = networkProp
	}
	privateClusterConfigProp, err := expandContainerClusterPrivateClusterConfig(d.Get("private_cluster_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("private_cluster_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(privateClusterConfigProp)) && (ok || !reflect.DeepEqual(v, privateClusterConfigProp)) {
		obj["privateClusterConfig"] = privateClusterConfigProp
	}
	workloadIdentityConfigProp, err := expandContainerClusterWorkloadIdentityConfig(d.Get("workload_identity_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("workload_identity_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(workloadIdentityConfigProp)) && (ok || !reflect.DeepEqual(v, workloadIdentityConfigProp)) {
		obj["workloadIdentityConfig"] = workloadIdentityConfigProp
	}
	clusterIpv4CidrProp, err := expandContainerClusterClusterIpv4Cidr(d.Get("cluster_ipv4_cidr"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("cluster_ipv4_cidr"); !tpgresource.IsEmptyValue(reflect.ValueOf(clusterIpv4CidrProp)) && (ok || !reflect.DeepEqual(v, clusterIpv4CidrProp)) {
		obj["clusterIpv4Cidr"] = clusterIpv4CidrProp
	}
	addonsConfigProp, err := expandContainerClusterAddonsConfig(d.Get("addons_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("addons_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(addonsConfigProp)) && (ok || !reflect.DeepEqual(v, addonsConfigProp)) {
		obj["addonsConfig"] = addonsConfigProp
	}
	subnetworkProp, err := expandContainerClusterSubnetwork(d.Get("subnetwork"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("subnetwork"); !tpgresource.IsEmptyValue(reflect.ValueOf(subnetworkProp)) && (ok || !reflect.DeepEqual(v, subnetworkProp)) {
		obj["subnetwork"] = subnetworkProp
	}
	locationsProp, err := expandContainerClusterNodeLocations(d.Get("node_locations"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("node_locations"); !tpgresource.IsEmptyValue(reflect.ValueOf(locationsProp)) && (ok || !reflect.DeepEqual(v, locationsProp)) {
		obj["locations"] = locationsProp
	}
	resourceLabelsProp, err := expandContainerClusterResourceLabels(d.Get("resource_labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("resource_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(resourceLabelsProp)) && (ok || !reflect.DeepEqual(v, resourceLabelsProp)) {
		obj["resourceLabels"] = resourceLabelsProp
	}
	labelFingerprintProp, err := expandContainerClusterLabelFingerprint(d.Get("label_fingerprint"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("label_fingerprint"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelFingerprintProp)) && (ok || !reflect.DeepEqual(v, labelFingerprintProp)) {
		obj["labelFingerprint"] = labelFingerprintProp
	}
	legacyAbacProp, err := expandContainerClusterEnableLegacyAbac(d.Get("enable_legacy_abac"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("enable_legacy_abac"); !tpgresource.IsEmptyValue(reflect.ValueOf(legacyAbacProp)) && (ok || !reflect.DeepEqual(v, legacyAbacProp)) {
		obj["legacyAbac"] = legacyAbacProp
	}
	networkPolicyProp, err := expandContainerClusterNetworkPolicy(d.Get("network_policy"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("network_policy"); !tpgresource.IsEmptyValue(reflect.ValueOf(networkPolicyProp)) && (ok || !reflect.DeepEqual(v, networkPolicyProp)) {
		obj["networkPolicy"] = networkPolicyProp
	}
	defaultMaxPodsConstraintProp, err := expandContainerClusterDefaultMaxPodsPerNode(d.Get("default_max_pods_per_node"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("default_max_pods_per_node"); !tpgresource.IsEmptyValue(reflect.ValueOf(defaultMaxPodsConstraintProp)) && (ok || !reflect.DeepEqual(v, defaultMaxPodsConstraintProp)) {
		obj["defaultMaxPodsConstraint"] = defaultMaxPodsConstraintProp
	}
	ipAllocationPolicyProp, err := expandContainerClusterIpAllocationPolicy(d.Get("ip_allocation_policy"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("ip_allocation_policy"); !tpgresource.IsEmptyValue(reflect.ValueOf(ipAllocationPolicyProp)) && (ok || !reflect.DeepEqual(v, ipAllocationPolicyProp)) {
		obj["ipAllocationPolicy"] = ipAllocationPolicyProp
	}
	initialClusterVersionProp, err := expandContainerClusterMinMasterVersion(d.Get("min_master_version"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("min_master_version"); !tpgresource.IsEmptyValue(reflect.ValueOf(initialClusterVersionProp)) && (ok || !reflect.DeepEqual(v, initialClusterVersionProp)) {
		obj["initialClusterVersion"] = initialClusterVersionProp
	}
	enableTpuProp, err := expandContainerClusterEnableTpu(d.Get("enable_tpu"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("enable_tpu"); !tpgresource.IsEmptyValue(reflect.ValueOf(enableTpuProp)) && (ok || !reflect.DeepEqual(v, enableTpuProp)) {
		obj["enableTpu"] = enableTpuProp
	}
	tpuIpv4CidrBlockProp, err := expandContainerClusterTPUIpv4CidrBlock(d.Get("tpu_ipv4_cidr_block"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("tpu_ipv4_cidr_block"); !tpgresource.IsEmptyValue(reflect.ValueOf(tpuIpv4CidrBlockProp)) && (ok || !reflect.DeepEqual(v, tpuIpv4CidrBlockProp)) {
		obj["tpuIpv4CidrBlock"] = tpuIpv4CidrBlockProp
	}
	masterAuthorizedNetworksConfigProp, err := expandContainerClusterMasterAuthorizedNetworksConfig(d.Get("master_authorized_networks_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("master_authorized_networks_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(masterAuthorizedNetworksConfigProp)) && (ok || !reflect.DeepEqual(v, masterAuthorizedNetworksConfigProp)) {
		obj["masterAuthorizedNetworksConfig"] = masterAuthorizedNetworksConfigProp
	}
	locationProp, err := expandContainerClusterLocation(d.Get("location"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("location"); !tpgresource.IsEmptyValue(reflect.ValueOf(locationProp)) && (ok || !reflect.DeepEqual(v, locationProp)) {
		obj["location"] = locationProp
	}
	kubectlPathProp, err := expandContainerClusterKubectlPath(d.Get("kubectl_path"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("kubectl_path"); !tpgresource.IsEmptyValue(reflect.ValueOf(kubectlPathProp)) && (ok || !reflect.DeepEqual(v, kubectlPathProp)) {
		obj["kubectlPath"] = kubectlPathProp
	}
	kubectlContextProp, err := expandContainerClusterKubectlContext(d.Get("kubectl_context"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("kubectl_context"); !tpgresource.IsEmptyValue(reflect.ValueOf(kubectlContextProp)) && (ok || !reflect.DeepEqual(v, kubectlContextProp)) {
		obj["kubectlContext"] = kubectlContextProp
	}
	databaseEncryptionProp, err := expandContainerClusterDatabaseEncryption(d.Get("database_encryption"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("database_encryption"); !tpgresource.IsEmptyValue(reflect.ValueOf(databaseEncryptionProp)) && (ok || !reflect.DeepEqual(v, databaseEncryptionProp)) {
		obj["databaseEncryption"] = databaseEncryptionProp
	}
	releaseChannelProp, err := expandContainerClusterReleaseChannel(d.Get("release_channel"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("release_channel"); !tpgresource.IsEmptyValue(reflect.ValueOf(releaseChannelProp)) && (ok || !reflect.DeepEqual(v, releaseChannelProp)) {
		obj["releaseChannel"] = releaseChannelProp
	}

	return obj, nil
}

func expandContainerClusterEnableKubernetesAlpha(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterPodSecurityPolicyConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEnabled, err := expandContainerClusterPodSecurityPolicyConfigEnabled(original["enabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enabled"] = transformedEnabled
	}

	return transformed, nil
}

func expandContainerClusterPodSecurityPolicyConfigEnabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterInitialNodeCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedWorkloadMetadataConfig, err := expandContainerClusterNodeConfigWorkloadMetadataConfig(original["workload_metadata_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedWorkloadMetadataConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["workloadMetadataConfig"] = transformedWorkloadMetadataConfig
	}

	transformedMachineType, err := expandContainerClusterNodeConfigMachineType(original["machine_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMachineType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["machineType"] = transformedMachineType
	}

	transformedDiskSizeGb, err := expandContainerClusterNodeConfigDiskSizeGb(original["disk_size_gb"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskSizeGb); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskSizeGb"] = transformedDiskSizeGb
	}

	transformedOauthScopes, err := expandContainerClusterNodeConfigOauthScopes(original["oauth_scopes"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOauthScopes); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["oauthScopes"] = transformedOauthScopes
	}

	transformedServiceAccount, err := expandContainerClusterNodeConfigServiceAccount(original["service_account"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedServiceAccount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["serviceAccount"] = transformedServiceAccount
	}

	transformedMetadata, err := expandContainerClusterNodeConfigMetadata(original["metadata"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMetadata); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["metadata"] = transformedMetadata
	}

	transformedImageType, err := expandContainerClusterNodeConfigImageType(original["image_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedImageType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["imageType"] = transformedImageType
	}

	transformedLabels, err := expandContainerClusterNodeConfigLabels(original["labels"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLabels); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["labels"] = transformedLabels
	}

	transformedLocalSsdCount, err := expandContainerClusterNodeConfigLocalSsdCount(original["local_ssd_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocalSsdCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["localSsdCount"] = transformedLocalSsdCount
	}

	transformedTags, err := expandContainerClusterNodeConfigTags(original["tags"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTags); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["tags"] = transformedTags
	}

	transformedPreemptible, err := expandContainerClusterNodeConfigPreemptible(original["preemptible"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPreemptible); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["preemptible"] = transformedPreemptible
	}

	transformedGuestAccelerator, err := expandContainerClusterNodeConfigGuestAccelerator(original["guest_accelerator"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGuestAccelerator); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["accelerators"] = transformedGuestAccelerator
	}

	transformedDiskType, err := expandContainerClusterNodeConfigDiskType(original["disk_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskType"] = transformedDiskType
	}

	transformedMinCpuPlatform, err := expandContainerClusterNodeConfigMinCpuPlatform(original["min_cpu_platform"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinCpuPlatform); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minCpuPlatform"] = transformedMinCpuPlatform
	}

	transformedTaint, err := expandContainerClusterNodeConfigTaint(original["taint"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTaint); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["taints"] = transformedTaint
	}

	return transformed, nil
}

func expandContainerClusterNodeConfigWorkloadMetadataConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMode, err := expandContainerClusterNodeConfigWorkloadMetadataConfigMode(original["mode"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMode); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["mode"] = transformedMode
	}

	return transformed, nil
}

func expandContainerClusterNodeConfigWorkloadMetadataConfigMode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigMachineType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigDiskSizeGb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigServiceAccount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigMetadata(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandContainerClusterNodeConfigImageType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandContainerClusterNodeConfigLocalSsdCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigTags(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigPreemptible(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigGuestAccelerator(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedCount, err := expandContainerClusterNodeConfigGuestAcceleratorCount(original["count"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["acceleratorCount"] = transformedCount
		}

		transformedType, err := expandContainerClusterNodeConfigGuestAcceleratorType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["acceleratorType"] = transformedType
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandContainerClusterNodeConfigGuestAcceleratorCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigGuestAcceleratorType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigDiskType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigMinCpuPlatform(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigTaint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedKey, err := expandContainerClusterNodeConfigTaintKey(original["key"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedKey); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["key"] = transformedKey
		}

		transformedValue, err := expandContainerClusterNodeConfigTaintValue(original["value"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedValue); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["value"] = transformedValue
		}

		transformedEffect, err := expandContainerClusterNodeConfigTaintEffect(original["effect"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedEffect); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["effect"] = transformedEffect
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandContainerClusterNodeConfigTaintKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigTaintValue(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeConfigTaintEffect(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterMasterAuth(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUsername, err := expandContainerClusterMasterAuthUsername(original["username"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUsername); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["username"] = transformedUsername
	}

	transformedPassword, err := expandContainerClusterMasterAuthPassword(original["password"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPassword); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["password"] = transformedPassword
	}

	transformedClientCertificateConfig, err := expandContainerClusterMasterAuthClientCertificateConfig(original["client_certificate_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClientCertificateConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clientCertificateConfig"] = transformedClientCertificateConfig
	}

	transformedClusterCaCertificate, err := expandContainerClusterMasterAuthClusterCaCertificate(original["cluster_ca_certificate"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClusterCaCertificate); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clusterCaCertificate"] = transformedClusterCaCertificate
	}

	transformedClientCertificate, err := expandContainerClusterMasterAuthClientCertificate(original["client_certificate"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClientCertificate); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clientCertificate"] = transformedClientCertificate
	}

	transformedClientKey, err := expandContainerClusterMasterAuthClientKey(original["client_key"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClientKey); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clientKey"] = transformedClientKey
	}

	return transformed, nil
}

func expandContainerClusterMasterAuthUsername(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterMasterAuthPassword(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterMasterAuthClientCertificateConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedIssueClientCertificate, err := expandContainerClusterMasterAuthClientCertificateConfigIssueClientCertificate(original["issue_client_certificate"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIssueClientCertificate); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["issueClientCertificate"] = transformedIssueClientCertificate
	}

	return transformed, nil
}

func expandContainerClusterMasterAuthClientCertificateConfigIssueClientCertificate(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterMasterAuthClusterCaCertificate(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterMasterAuthClientCertificate(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterMasterAuthClientKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterLoggingService(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterMonitoringService(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterWorkloadIdentityConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedWorkloadPool, err := expandContainerClusterWorkloadIdentityConfigWorkloadPool(original["workload_pool"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedWorkloadPool); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["workloadPool"] = transformedWorkloadPool
	}
	return transformed, nil
}

func expandContainerClusterWorkloadIdentityConfigWorkloadPool(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterPrivateClusterConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEnablePrivateNodes, err := expandContainerClusterPrivateClusterConfigEnablePrivateNodes(original["enable_private_nodes"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnablePrivateNodes); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enablePrivateNodes"] = transformedEnablePrivateNodes
	}

	transformedEnablePrivateEndpoint, err := expandContainerClusterPrivateClusterConfigEnablePrivateEndpoint(original["enable_private_endpoint"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnablePrivateEndpoint); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enablePrivateEndpoint"] = transformedEnablePrivateEndpoint
	}

	transformedMasterIpv4CidrBlock, err := expandContainerClusterPrivateClusterConfigMasterIpv4CidrBlock(original["master_ipv4_cidr_block"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMasterIpv4CidrBlock); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["masterIpv4CidrBlock"] = transformedMasterIpv4CidrBlock
	}

	transformedPrivateEndpoint, err := expandContainerClusterPrivateClusterConfigPrivateEndpoint(original["private_endpoint"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPrivateEndpoint); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["privateEndpoint"] = transformedPrivateEndpoint
	}

	transformedPublicEndpoint, err := expandContainerClusterPrivateClusterConfigPublicEndpoint(original["public_endpoint"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPublicEndpoint); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["publicEndpoint"] = transformedPublicEndpoint
	}

	return transformed, nil
}

func expandContainerClusterPrivateClusterConfigEnablePrivateNodes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterPrivateClusterConfigEnablePrivateEndpoint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterPrivateClusterConfigMasterIpv4CidrBlock(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterPrivateClusterConfigPrivateEndpoint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterPrivateClusterConfigPublicEndpoint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterClusterIpv4Cidr(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterAddonsConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedHttpLoadBalancing, err := expandContainerClusterAddonsConfigHttpLoadBalancing(original["http_load_balancing"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedHttpLoadBalancing); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["httpLoadBalancing"] = transformedHttpLoadBalancing
	}

	transformedHorizontalPodAutoscaling, err := expandContainerClusterAddonsConfigHorizontalPodAutoscaling(original["horizontal_pod_autoscaling"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedHorizontalPodAutoscaling); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["horizontalPodAutoscaling"] = transformedHorizontalPodAutoscaling
	}

	transformedNetworkPolicyConfig, err := expandContainerClusterAddonsConfigNetworkPolicyConfig(original["network_policy_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNetworkPolicyConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["networkPolicyConfig"] = transformedNetworkPolicyConfig
	}

	return transformed, nil
}

func expandContainerClusterAddonsConfigHttpLoadBalancing(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDisabled, err := expandContainerClusterAddonsConfigHttpLoadBalancingDisabled(original["disabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDisabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["disabled"] = transformedDisabled
	}

	return transformed, nil
}

func expandContainerClusterAddonsConfigHttpLoadBalancingDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterAddonsConfigHorizontalPodAutoscaling(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDisabled, err := expandContainerClusterAddonsConfigHorizontalPodAutoscalingDisabled(original["disabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDisabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["disabled"] = transformedDisabled
	}

	return transformed, nil
}

func expandContainerClusterAddonsConfigHorizontalPodAutoscalingDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterAddonsConfigNetworkPolicyConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDisabled, err := expandContainerClusterAddonsConfigNetworkPolicyConfigDisabled(original["disabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDisabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["disabled"] = transformedDisabled
	}

	return transformed, nil
}

func expandContainerClusterAddonsConfigNetworkPolicyConfigDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNodeLocations(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandContainerClusterResourceLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandContainerClusterLabelFingerprint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNetworkPolicy(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedProvider, err := expandContainerClusterNetworkPolicyProvider(original["provider"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedProvider); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["provider"] = transformedProvider
	}

	transformedEnabled, err := expandContainerClusterNetworkPolicyEnabled(original["enabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enabled"] = transformedEnabled
	}

	return transformed, nil
}

func expandContainerClusterNetworkPolicyProvider(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterNetworkPolicyEnabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterIpAllocationPolicy(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUseIpAliases, err := expandContainerClusterIpAllocationPolicyUseIpAliases(original["use_ip_aliases"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUseIpAliases); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["useIpAliases"] = transformedUseIpAliases
	}

	transformedCreateSubnetwork, err := expandContainerClusterIpAllocationPolicyCreateSubnetwork(original["create_subnetwork"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCreateSubnetwork); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["createSubnetwork"] = transformedCreateSubnetwork
	}

	transformedSubnetworkName, err := expandContainerClusterIpAllocationPolicySubnetworkName(original["subnetwork_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSubnetworkName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["subnetworkName"] = transformedSubnetworkName
	}

	transformedClusterSecondaryRangeName, err := expandContainerClusterIpAllocationPolicyClusterSecondaryRangeName(original["cluster_secondary_range_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClusterSecondaryRangeName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clusterSecondaryRangeName"] = transformedClusterSecondaryRangeName
	}

	transformedServicesSecondaryRangeName, err := expandContainerClusterIpAllocationPolicyServicesSecondaryRangeName(original["services_secondary_range_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedServicesSecondaryRangeName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["servicesSecondaryRangeName"] = transformedServicesSecondaryRangeName
	}

	transformedClusterIpv4CidrBlock, err := expandContainerClusterIpAllocationPolicyClusterIpv4CidrBlock(original["cluster_ipv4_cidr_block"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClusterIpv4CidrBlock); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clusterIpv4CidrBlock"] = transformedClusterIpv4CidrBlock
	}

	transformedNodeIpv4CidrBlock, err := expandContainerClusterIpAllocationPolicyNodeIpv4CidrBlock(original["node_ipv4_cidr_block"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNodeIpv4CidrBlock); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["nodeIpv4CidrBlock"] = transformedNodeIpv4CidrBlock
	}

	transformedServicesIpv4CidrBlock, err := expandContainerClusterIpAllocationPolicyServicesIpv4CidrBlock(original["services_ipv4_cidr_block"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedServicesIpv4CidrBlock); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["servicesIpv4CidrBlock"] = transformedServicesIpv4CidrBlock
	}

	transformedTPUIpv4CidrBlock, err := expandContainerClusterIpAllocationPolicyTPUIpv4CidrBlock(original["tpu_ipv4_cidr_block"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTPUIpv4CidrBlock); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["tpuIpv4CidrBlock"] = transformedTPUIpv4CidrBlock
	}

	return transformed, nil
}

func expandContainerClusterIpAllocationPolicyUseIpAliases(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterIpAllocationPolicyCreateSubnetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterIpAllocationPolicySubnetworkName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterIpAllocationPolicyClusterSecondaryRangeName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterIpAllocationPolicyServicesSecondaryRangeName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterIpAllocationPolicyClusterIpv4CidrBlock(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterIpAllocationPolicyNodeIpv4CidrBlock(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterIpAllocationPolicyServicesIpv4CidrBlock(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterIpAllocationPolicyTPUIpv4CidrBlock(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterMinMasterVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterEnableTpu(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterTPUIpv4CidrBlock(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterMasterAuthorizedNetworksConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	// enabled is always true as long as there is a master_authorized_networks_config config block.
	// There is no option in Terraform to disable that when master_authorized_networks_config is seen.
	transformed["enabled"] = true

	transformedCidrBlocks, err := expandContainerClusterMasterAuthorizedNetworksConfigCidrBlocks(original["cidr_blocks"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCidrBlocks); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["cidrBlocks"] = transformedCidrBlocks
	}

	return transformed, nil
}

func expandContainerClusterMasterAuthorizedNetworksConfigCidrBlocks(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedDisplayName, err := expandContainerClusterMasterAuthorizedNetworksConfigCidrBlocksDisplayName(original["display_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDisplayName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["displayName"] = transformedDisplayName
		}

		transformedCidrBlock, err := expandContainerClusterMasterAuthorizedNetworksConfigCidrBlocksCidrBlock(original["cidr_block"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedCidrBlock); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["cidrBlock"] = transformedCidrBlock
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandContainerClusterMasterAuthorizedNetworksConfigCidrBlocksDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterMasterAuthorizedNetworksConfigCidrBlocksCidrBlock(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterLocation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterKubectlPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterKubectlContext(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterDatabaseEncryption(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerClusterReleaseChannel(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func GetContainerNodePoolCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//container.googleapis.com/projects/{{project}}/locations/{{location}}/clusters/{{cluster}}/nodePools/{{name}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetContainerNodePoolApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: ContainerNodePoolAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/container/v1/rest",
				DiscoveryName:        "NodePool",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetContainerNodePoolApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandContainerNodePoolName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	configProp, err := expandContainerNodePoolNodeConfig(d.Get("node_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("node_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(configProp)) && (ok || !reflect.DeepEqual(v, configProp)) {
		obj["config"] = configProp
	}
	initialNodeCountProp, err := expandContainerNodePoolInitialNodeCount(d.Get("initial_node_count"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("initial_node_count"); !tpgresource.IsEmptyValue(reflect.ValueOf(initialNodeCountProp)) && (ok || !reflect.DeepEqual(v, initialNodeCountProp)) {
		obj["initialNodeCount"] = initialNodeCountProp
	}
	versionProp, err := expandContainerNodePoolVersion(d.Get("version"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("version"); !tpgresource.IsEmptyValue(reflect.ValueOf(versionProp)) && (ok || !reflect.DeepEqual(v, versionProp)) {
		obj["version"] = versionProp
	}
	autoscalingProp, err := expandContainerNodePoolAutoscaling(d.Get("autoscaling"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("autoscaling"); !tpgresource.IsEmptyValue(reflect.ValueOf(autoscalingProp)) && (ok || !reflect.DeepEqual(v, autoscalingProp)) {
		obj["autoscaling"] = autoscalingProp
	}
	managementProp, err := expandContainerNodePoolManagement(d.Get("management"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("management"); !tpgresource.IsEmptyValue(reflect.ValueOf(managementProp)) && (ok || !reflect.DeepEqual(v, managementProp)) {
		obj["management"] = managementProp
	}
	maxPodsConstraintProp, err := expandContainerNodePoolMaxPodsPerNode(d.Get("max_pods_per_node"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("max_pods_per_node"); !tpgresource.IsEmptyValue(reflect.ValueOf(maxPodsConstraintProp)) && (ok || !reflect.DeepEqual(v, maxPodsConstraintProp)) {
		obj["maxPodsConstraint"] = maxPodsConstraintProp
	}
	clusterProp, err := expandContainerNodePoolCluster(d.Get("cluster"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("cluster"); !tpgresource.IsEmptyValue(reflect.ValueOf(clusterProp)) && (ok || !reflect.DeepEqual(v, clusterProp)) {
		obj["cluster"] = clusterProp
	}
	locationProp, err := expandContainerNodePoolLocation(d.Get("location"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("location"); !tpgresource.IsEmptyValue(reflect.ValueOf(locationProp)) && (ok || !reflect.DeepEqual(v, locationProp)) {
		obj["location"] = locationProp
	}

	return obj, nil
}

func expandContainerNodePoolName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMachineType, err := expandContainerNodePoolNodeConfigMachineType(original["machine_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMachineType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["machineType"] = transformedMachineType
	}

	transformedDiskSizeGb, err := expandContainerNodePoolNodeConfigDiskSizeGb(original["disk_size_gb"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskSizeGb); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskSizeGb"] = transformedDiskSizeGb
	}

	transformedOauthScopes, err := expandContainerNodePoolNodeConfigOauthScopes(original["oauth_scopes"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOauthScopes); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["oauthScopes"] = transformedOauthScopes
	}

	transformedServiceAccount, err := expandContainerNodePoolNodeConfigServiceAccount(original["service_account"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedServiceAccount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["serviceAccount"] = transformedServiceAccount
	}

	transformedMetadata, err := expandContainerNodePoolNodeConfigMetadata(original["metadata"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMetadata); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["metadata"] = transformedMetadata
	}

	transformedImageType, err := expandContainerNodePoolNodeConfigImageType(original["image_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedImageType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["imageType"] = transformedImageType
	}

	transformedLabels, err := expandContainerNodePoolNodeConfigLabels(original["labels"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLabels); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["labels"] = transformedLabels
	}

	transformedLocalSsdCount, err := expandContainerNodePoolNodeConfigLocalSsdCount(original["local_ssd_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocalSsdCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["localSsdCount"] = transformedLocalSsdCount
	}

	transformedTags, err := expandContainerNodePoolNodeConfigTags(original["tags"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTags); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["tags"] = transformedTags
	}

	transformedPreemptible, err := expandContainerNodePoolNodeConfigPreemptible(original["preemptible"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPreemptible); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["preemptible"] = transformedPreemptible
	}

	transformedGuestAccelerator, err := expandContainerNodePoolNodeConfigGuestAccelerator(original["guest_accelerator"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGuestAccelerator); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["accelerators"] = transformedGuestAccelerator
	}

	transformedDiskType, err := expandContainerNodePoolNodeConfigDiskType(original["disk_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskType"] = transformedDiskType
	}

	transformedMinCpuPlatform, err := expandContainerNodePoolNodeConfigMinCpuPlatform(original["min_cpu_platform"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinCpuPlatform); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minCpuPlatform"] = transformedMinCpuPlatform
	}

	transformedTaint, err := expandContainerNodePoolNodeConfigTaint(original["taint"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTaint); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["taints"] = transformedTaint
	}

	return transformed, nil
}

func expandContainerNodePoolNodeConfigMachineType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigDiskSizeGb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigServiceAccount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigMetadata(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandContainerNodePoolNodeConfigImageType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandContainerNodePoolNodeConfigLocalSsdCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigTags(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigPreemptible(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigGuestAccelerator(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedCount, err := expandContainerNodePoolNodeConfigGuestAcceleratorCount(original["count"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["acceleratorCount"] = transformedCount
		}

		transformedType, err := expandContainerNodePoolNodeConfigGuestAcceleratorType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["acceleratorType"] = transformedType
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandContainerNodePoolNodeConfigGuestAcceleratorCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigGuestAcceleratorType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigDiskType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigMinCpuPlatform(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigTaint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedKey, err := expandContainerNodePoolNodeConfigTaintKey(original["key"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedKey); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["key"] = transformedKey
		}

		transformedValue, err := expandContainerNodePoolNodeConfigTaintValue(original["value"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedValue); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["value"] = transformedValue
		}

		transformedEffect, err := expandContainerNodePoolNodeConfigTaintEffect(original["effect"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedEffect); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["effect"] = transformedEffect
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandContainerNodePoolNodeConfigTaintKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigTaintValue(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigTaintEffect(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolInitialNodeCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolAutoscaling(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEnabled, err := expandContainerNodePoolAutoscalingEnabled(original["enabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enabled"] = transformedEnabled
	}

	transformedMinNodeCount, err := expandContainerNodePoolAutoscalingMinNodeCount(original["min_node_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinNodeCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minNodeCount"] = transformedMinNodeCount
	}

	transformedMaxNodeCount, err := expandContainerNodePoolAutoscalingMaxNodeCount(original["max_node_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMaxNodeCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["maxNodeCount"] = transformedMaxNodeCount
	}

	return transformed, nil
}

func expandContainerNodePoolAutoscalingEnabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolAutoscalingMinNodeCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolAutoscalingMaxNodeCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolManagement(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAutoUpgrade, err := expandContainerNodePoolManagementAutoUpgrade(original["auto_upgrade"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAutoUpgrade); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["autoUpgrade"] = transformedAutoUpgrade
	}

	transformedAutoRepair, err := expandContainerNodePoolManagementAutoRepair(original["auto_repair"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAutoRepair); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["autoRepair"] = transformedAutoRepair
	}

	return transformed, nil
}

func expandContainerNodePoolManagementAutoUpgrade(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolManagementAutoRepair(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolCluster(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	f, err := tpgresource.ParseGlobalFieldValue("clusters", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for cluster: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandContainerNodePoolLocation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
