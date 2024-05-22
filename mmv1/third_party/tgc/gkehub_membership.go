package google

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const GkehubMembershipAssetType string = "gkehub.googleapis.com/Membership"

func resourceGkehubMembership() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: GkehubMembershipAssetType,
		Convert:   GetGkehubMembershipCaiObject,
	}
}

func GetGkehubMembershipCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//gkehub.googleapis.com/projects/{{project}}/locations/{{location}}/memberships/basic")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetGkehubMembershipApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: GkehubMembershipAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://gkehub.googleapis.com/$discovery/rest",
				DiscoveryName:        "Membership",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetGkehubMembershipApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	endpointProp, err := expandGkehubMembershipEndpoint(d.Get("endpoint"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("endpoint"); !tpgresource.IsEmptyValue(reflect.ValueOf(endpointProp)) && (ok || !reflect.DeepEqual(v, endpointProp)) {
		obj["endpoint"] = endpointProp
	}	
	nameProp, err := expandGkehubMembershipName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	labelsProp, err := expandGkehubMembershipLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	descriptionProp, err := expandGkehubMembershipDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	createTimeProp, err := expandGkehubMembershipCreateTime(d.Get("create_time"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("create_time"); !tpgresource.IsEmptyValue(reflect.ValueOf(createTimeProp)) && (ok || !reflect.DeepEqual(v, createTimeProp)) {
		obj["createTime"] = createTimeProp
	}
	updateTimeProp, err := expandGkehubMembershipUpdateTime(d.Get("update_time"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("update_time"); !tpgresource.IsEmptyValue(reflect.ValueOf(updateTimeProp)) && (ok || !reflect.DeepEqual(v, updateTimeProp)) {
		obj["updateTime"] = updateTimeProp
	}
	deleteTimeProp, err := expandGkehubMembershipDeleteTime(d.Get("delete_time"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("delete_time"); !tpgresource.IsEmptyValue(reflect.ValueOf(deleteTimeProp)) && (ok || !reflect.DeepEqual(v, deleteTimeProp)) {
		obj["deleteTime"] = deleteTimeProp
	}
	externalIdProp, err := expandGkehubMembershipExternalId(d.Get("external_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("external_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(externalIdProp)) && (ok || !reflect.DeepEqual(v, externalIdProp)) {
		obj["externalId"] = externalIdProp
	}
	lastConnectionTimeProp, err := expandGkehubMembershipLastConnectionTime(d.Get("last_connection_time"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("last_connection_time"); !tpgresource.IsEmptyValue(reflect.ValueOf(lastConnectionTimeProp)) && (ok || !reflect.DeepEqual(v, lastConnectionTimeProp)) {
		obj["lastConnectionTime"] = lastConnectionTimeProp
	}
	uniqueIdProp, err := expandGkehubMembershipUniqueId(d.Get("unique_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("unique_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(uniqueIdProp)) && (ok || !reflect.DeepEqual(v, uniqueIdProp)) {
		obj["uniqueId"] = uniqueIdProp
	}
	authorityProp, err := expandGkehubMembershipAuthority(d.Get("authority"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("authority"); !tpgresource.IsEmptyValue(reflect.ValueOf(authorityProp)) && (ok || !reflect.DeepEqual(v, authorityProp)) {
		obj["authority"] = authorityProp
	}

	return obj, nil
}

func expandGkehubMembershipMonitoringConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})
	
	transformedProjectId, err := expandGkehubMembershipProjectId(original["project_id"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedProjectId); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["projectId"] = transformedProjectId
	}
	transformedLocation, err := expandGkehubMembershipLocation(original["location"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocation); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["location"] = transformedLocation
	}
	transformedCluster, err := expandGkehubMembershipCluster(original["cluster"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCluster); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["cluster"] = transformedCluster
	}
	transformedKubernetesMetricsPrefix, err := expandGkehubMembershipKubernetesMetricsPrefix(original["kubernetes_metrics_prefix"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedKubernetesMetricsPrefix); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["kubernetesMetricsPrefix"] = transformedKubernetesMetricsPrefix
	}
	transformedClusterHash, err := expandGkehubMembershipClusterHash(original["cluster_hash"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClusterHash); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clusterHash"] = transformedClusterHash
	}

	return transformed, nil
}

func expandGkehubMembershipDeleteTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipClusterHash(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipKubernetesMetricsPrefix(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipCluster(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipLocation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipProjectId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipAuthority(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})
	
	transformedIssuer, err := expandGkehubMembershipIssuer(original["issuer"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIssuer); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["issuer"] = transformedIssuer
	}

	return transformed, nil
}

func expandGkehubMembershipOidcJwks(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipIdentityProvider(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipWorkloadIdentityPool(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipIssuer(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipUniqueId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipLastConnectionTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipExternalId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipCreateTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipState(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})
	
	transformedCode, err := expandGkehubMembershipCode(original["code"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCode); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["code"] = transformedCode
	}

	return transformed, nil
}

func expandGkehubMembershipCode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipEndpoint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})
	
	transformedGkeCluster, err := expandGkehubMembershipGkeCluster(original["gke_cluster"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGkeCluster); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["gkeCluster"] = transformedGkeCluster
	}

	return transformed, nil
}

func expandGkehubMembershipGoogleManaged(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipGkeCluster(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})
	
	transformedResourceLink, err := expandGkehubMembershipResourceLink(original["resource_link"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedResourceLink); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["resourceLink"] = transformedResourceLink
	}
	transformedClusterMissing, err := expandGkehubMembershipClusterMissing(original["cluster_missing"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClusterMissing); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clusterMissing"] = transformedClusterMissing
	}

	return transformed, nil
}

func expandGkehubMembershipOnPremCluster(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedResourceLink, err := expandGkehubMembershipResourceLink(original["resource_link"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedResourceLink); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["resourceLink"] = transformedResourceLink
	}
	transformedClusterMissing, err := expandGkehubMembershipClusterMissing(original["cluster_missing"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClusterMissing); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clusterMissing"] = transformedClusterMissing
	}
	transformedAdminCluster, err := expandGkehubMembershipAdminCluster(original["admin_cluster"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAdminCluster); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["adminCluster"] = transformedAdminCluster
	}
	transformedClusterType, err := expandGkehubMembershipClusterType(original["cluster_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClusterType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clusterType"] = transformedClusterType
	}

	return transformed, nil
}

func expandGkehubMembershipMultiCloudCluster(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})
	
	transformedResourceLink, err := expandGkehubMembershipResourceLink(original["resource_link"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedResourceLink); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["resourceLink"] = transformedResourceLink
	}
	transformedClusterMissing, err := expandGkehubMembershipClusterMissing(original["cluster_missing"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClusterMissing); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["clusterMissing"] = transformedClusterMissing
	}

	return transformed, nil
}

func expandGkehubMembershipEdgeCluster(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})
	
	transformedResourceLink, err := expandGkehubMembershipResourceLink(original["resource_link"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedResourceLink); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["resourceLink"] = transformedResourceLink
	}

	return transformed, nil
}

func expandGkehubMembershipApplianceCluster(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})
	
	transformedResourceLink, err := expandGkehubMembershipResourceLink(original["resource_link"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedResourceLink); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["resourceLink"] = transformedResourceLink
	}

	return transformed, nil
}

func expandGkehubMembershipKubernetesMetadata(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})
	
	transformedKubernetesApiServerVersion, err := expandGkehubMembershipKubernetesApiServerVersion(original["kubernetes_api_server_version	"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedKubernetesApiServerVersion); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["kubernetesApiServerVersion"] = transformedKubernetesApiServerVersion
	}
	transformedNodeProviderId, err := expandGkehubMembershipNodeProviderId(original["node_provider_id"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNodeProviderId); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["nodeProviderId"] = transformedNodeProviderId
	}
	transformedNodeCount, err := expandGkehubMembershipNodeCount(original["node_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNodeCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["nodeCount"] = transformedNodeCount
	}
	transformedVcpuCount, err := expandGkehubMembershipVcpuCount(original["vcpu_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedVcpuCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["vcpuCount"] = transformedVcpuCount
	}
	transformedMemoryMb, err := expandGkehubMembershipMemoryMb(original["memory_mb"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMemoryMb); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["memoryMb"] = transformedMemoryMb
	}
	transformedUpdateTime, err := expandGkehubMembershipUpdateTime(original["update_time"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUpdateTime); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["updateTime"] = transformedUpdateTime
	}

	return transformed, nil
}

func expandGkehubMembershipKubernetesResource(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})
	
	transformedMembershipCrManifest, err := expandGkehubMembershipCrManifest(original["membership_cr_manifest"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMembershipCrManifest); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["membershipCrManifest"] = transformedMembershipCrManifest
	}

	transformedMembershipResources, err := expandGkehubMembershipResources(original["membership_resources"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMembershipResources); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["membershipResources"] = transformedMembershipResources
	}

	transformedConnectResources, err := expandGkehubMembershipConnectResources(original["connect_resources"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedConnectResources); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["connectResources"] = transformedConnectResources
	}

	transformedResourceOptions, err := expandGkehubMembershipResourceOptions(original["resource_options"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedResourceOptions); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["ResourceOptions"] = transformedResourceOptions
	}

	return transformed, nil
}

func expandGkehubMembershipResourceOptions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedConnectVersion, err := expandGkehubMembershipConnectVersion(original["connect_version"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedConnectVersion); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["connectVersion"] = transformedConnectVersion
	}
	transformedV1beta1Crd, err := expandGkehubMembershipV1beta1Crd(original["b1beta1_crd"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedV1beta1Crd); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["v1beta1Crd"] = transformedV1beta1Crd
	}
	transformedK8sVersion, err := expandGkehubMembershipK8sVersion(original["k8s_version"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedK8sVersion); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["k8sVersion"] = transformedK8sVersion
	}
	
	return transformed, nil
}

func expandGkehubMembershipConnectVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipV1beta1Crd(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipK8sVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipResources(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedManifest, err := expandGkehubMembershipManifest(original["manifest"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedManifest); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["manifest"] = transformedManifest
		}

		transformedClusterScoped, err := expandGkehubMembershipClusterScoped(original["ClusterScoped"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedClusterScoped); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["clusterScoped"] = transformedClusterScoped
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandGkehubMembershipConnectResources(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedManifest, err := expandGkehubMembershipManifest(original["manifest"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedManifest); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["manifest"] = transformedManifest
		}

		transformedClusterScoped, err := expandGkehubMembershipClusterScoped(original["ClusterScoped"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedClusterScoped); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["clusterScoped"] = transformedClusterScoped
		}		

		req = append(req, transformed)
	}
	return req, nil
}

func expandGkehubMembershipClusterScoped(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipManifest(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipCrManifest(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipUpdateTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipMemoryMb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipVcpuCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipNodeCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipKubernetesApiServerVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipNodeProviderId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipResourceLink(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipClusterMissing(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipAdminCluster(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipClusterType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandGkehubMembershipLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
