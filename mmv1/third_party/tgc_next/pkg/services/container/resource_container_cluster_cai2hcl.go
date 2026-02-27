package container

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"google.golang.org/api/container/v1"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/transport"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ContainerClusterCai2hclConverter struct {
	name   string
	schema map[string]*schema.Schema
}

func NewContainerClusterCai2hclConverter(provider *schema.Provider) models.Cai2hclConverter {
	return &ContainerClusterCai2hclConverter{
		name:   ContainerClusterSchemaName,
		schema: provider.ResourcesMap[ContainerClusterSchemaName].Schema,
	}
}

func (c *ContainerClusterCai2hclConverter) Convert(asset caiasset.Asset) ([]*models.TerraformResourceBlock, error) {
	var blocks []*models.TerraformResourceBlock
	block, err := c.convertResourceData(asset)
	if err != nil {
		return nil, err
	}
	blocks = append(blocks, block)
	return blocks, nil
}

func (c *ContainerClusterCai2hclConverter) convertResourceData(asset caiasset.Asset) (*models.TerraformResourceBlock, error) {
	if asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}

	var err error
	config := transport.NewConfig()

	// This is a fake resource used to get fake d
	// d.Get will return empty map, instead of nil
	fakeResource := &schema.Resource{
		Schema: c.schema,
	}
	d := fakeResource.TestResourceData()

	hclData := make(map[string]interface{})

	outputFields := map[string]struct{}{}
	an := strings.Replace(asset.Name, "/zones/", "/locations/", 1)
	utils.ParseUrlParamValuesFromAssetName(an, "//container.googleapis.com/projects/{{project}}/locations/{{location}}/clusters/{{cluster}}/nodePools/{{name}}", outputFields, hclData)

	hclData["name"] = asset.Resource.Data["name"]
	hclData["description"] = asset.Resource.Data["description"]
	hclData["security_posture_config"] = flattenSecurityPostureConfig(asset.Resource.Data["securityPostureConfig"])
	hclData["enterprise_config"] = flattenEnterpriseConfig(asset.Resource.Data["enterpriseConfig"])
	hclData["anonymous_authentication_config"] = flattenAnonymousAuthenticationConfig(asset.Resource.Data["anonymousAuthenticationConfig"])
	hclData["notification_config"] = flattenNotificationConfig(asset.Resource.Data["notificationConfig"])
	hclData["binary_authorization"] = flattenBinaryAuthorization(asset.Resource.Data["binaryAuthorization"])
	hclData["network_policy"] = flattenNetworkPolicy(asset.Resource.Data["networkPolicy"])
	hclData["addons_config"] = flattenClusterAddonsConfig(asset.Resource.Data["addonsConfig"])
	hclData["node_pool"], err = flattenContainerClusterNodePools(d, config, asset.Resource.Data["nodePools"])
	hclData["authenticator_groups_config"] = flattenAuthenticatorGroupsConfig(asset.Resource.Data["authenticatorGroupsConfig"])
	hclData["control_plane_endpoints_config"] = flattenControlPlaneEndpointsConfig(asset.Resource.Data["controlPlaneEndpointsConfig"])

	privateClusterConfig := asset.Resource.Data["privateClusterConfig"]
	controlPlaneEndpointsConfig := asset.Resource.Data["controlPlaneEndpointsConfig"]
	networkConfig := asset.Resource.Data["networkConfig"]
	hclData["private_cluster_config"] = flattenPrivateClusterConfig(controlPlaneEndpointsConfig, privateClusterConfig, networkConfig)

	hclData["vertical_pod_autoscaling"] = flattenVerticalPodAutoscaling(asset.Resource.Data["verticalPodAutoscaling"])
	hclData["release_channel"] = flattenReleaseChannel(asset.Resource.Data["releaseChannel"])
	hclData["gke_auto_upgrade_config"] = flattenGkeAutoUpgradeConfig(asset.Resource.Data["gkeAutoUpgradeConfig"])

	if nc, ok := networkConfig.(map[string]interface{}); ok {
		hclData["default_snat_status"] = flattenDefaultSnatStatus(nc["defaultSnatStatus"])
		hclData["service_external_ips_config"] = flattenServiceExternalIpsConfig(nc["serviceExternalIpsConfig"])
		hclData["dns_config"] = flattenDnsConfig(nc["dnsConfig"])
		hclData["network_performance_config"] = flattenNetworkPerformanceConfig(nc["networkPerformanceConfig"])
		hclData["gateway_api_config"] = flattenGatewayApiConfig(nc["gatewayApiConfig"])
	}

	hclData["workload_identity_config"] = flattenWorkloadIdentityConfig(asset.Resource.Data["workloadIdentityConfig"])
	hclData["identity_service_config"] = flattenIdentityServiceConfig(asset.Resource.Data["identityServiceConfig"])

	if ipAlloc, ok := asset.Resource.Data["ipAllocationPolicy"].(map[string]interface{}); ok {
		hclData["pod_cidr_overprovision_config"] = flattenPodCidrOverprovisionConfig(ipAlloc["podCidrOverprovisionConfig"])
	}

	ipPolicy, err := flattenIPAllocationPolicy(asset.Resource.Data, nil, nil)
	if err != nil {
		return nil, err
	}
	hclData["ip_allocation_policy"] = ipPolicy

	if ipAlloc, ok := asset.Resource.Data["ipAllocationPolicy"].(map[string]interface{}); !ok || ipAlloc == nil {
		hclData["networking_mode"] = "ROUTES"
	} else if useIpAliases, ok := ipAlloc["useIpAliases"].(bool); !ok || !useIpAliases {
		hclData["networking_mode"] = "ROUTES"
	} else {
		hclData["networking_mode"] = "VPC_NATIVE"
	}

	hclData["maintenance_policy"] = flattenMaintenancePolicy(asset.Resource.Data["maintenancePolicy"])
	hclData["master_auth"] = flattenMasterAuth(asset.Resource.Data["masterAuth"])
	hclData["cluster_autoscaling"] = flattenClusterAutoscaling(asset.Resource.Data["autoscaling"])
	hclData["master_authorized_networks_config"] = flattenMasterAuthorizedNetworksConfig(asset.Resource.Data["masterAuthorizedNetworksConfig"])
	hclData["pod_autoscaling"] = flattenPodAutoscaling(asset.Resource.Data["podAutoscaling"])
	hclData["secret_manager_config"] = flattenSecretManagerConfig(asset.Resource.Data["secretManagerConfig"])
	hclData["resource_usage_export_config"] = flattenResourceUsageExportConfig(asset.Resource.Data["resourceUsageExportConfig"])
	hclData["mesh_certificates"] = flattenMeshCertificates(asset.Resource.Data["meshCertificates"])
	hclData["cost_management_config"] = flattenManagementConfig(asset.Resource.Data["costManagementConfig"])
	hclData["database_encryption"] = flattenDatabaseEncryption(asset.Resource.Data["databaseEncryption"])
	hclData["fleet"] = flattenFleet(asset.Resource.Data["fleet"])
	hclData["user_managed_keys_config"] = flattenUserManagedKeysConfig(asset.Resource.Data["userManagedKeysConfig"])
	hclData["enable_k8s_beta_apis"] = flattenEnableK8sBetaApis(asset.Resource.Data["enableK8sBetaApis"])
	hclData["logging_config"] = flattenContainerClusterLoggingConfig(asset.Resource.Data["loggingConfig"])
	hclData["monitoring_config"] = flattenMonitoringConfig(asset.Resource.Data["monitoringConfig"])
	hclData["node_pool_auto_config"] = flattenNodePoolAutoConfig(asset.Resource.Data["nodePoolAutoConfig"])
	hclData["rbac_binding_config"] = flattenRBACBindingConfig(asset.Resource.Data["rbacBindingConfig"])

	ctyVal, err := utils.MapToCtyValWithSchema(hclData, c.schema)
	if err != nil {
		return nil, err
	}
	// name is likely string, safe cast or fallback
	name, _ := asset.Resource.Data["name"].(string)
	return &models.TerraformResourceBlock{
		Labels: []string{c.name, name},
		Value:  ctyVal,
	}, nil
}

func flattenSecurityPostureConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	spc, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"mode":               spc["mode"],
		"vulnerability_mode": spc["vulnerabilityMode"],
	}

	return []map[string]interface{}{transformed}
}

func flattenEnterpriseConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	ec, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"cluster_tier": ec["clusterTier"],
		"desired_tier": ec["desiredTier"],
	}

	return []map[string]interface{}{transformed}
}

func flattenAnonymousAuthenticationConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	aac, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"mode": aac["enabled"],
	}

	return []map[string]interface{}{transformed}
}

func flattenAdditionalPodRangesConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	ipAllocationPolicy, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	result := make(map[string]interface{})

	if aprc, ok := ipAllocationPolicy["additionalPodRangesConfig"].(map[string]interface{}); ok && aprc != nil {
		if names, ok := aprc["podRangeNames"].([]interface{}); ok && len(names) > 0 {
			result["pod_range_names"] = names
		} else {
			return nil
		}
	} else {
		return nil
	}

	return []map[string]interface{}{result}
}

func isEnablePrivateEndpointPSCCluster(cluster *container.Cluster) bool {
	// EnablePrivateEndpoint not provided
	if cluster == nil || cluster.PrivateClusterConfig == nil {
		return false
	}
	// Not a PSC cluster
	if cluster.PrivateClusterConfig.EnablePrivateNodes || len(cluster.PrivateClusterConfig.MasterIpv4CidrBlock) > 0 {
		return false
	}
	// PSC Cluster with EnablePrivateEndpoint
	if cluster.PrivateClusterConfig.EnablePrivateEndpoint {
		return true
	}
	return false
}

func isEnablePDCSI(cluster *container.Cluster) bool {
	if cluster.AddonsConfig == nil || cluster.AddonsConfig.GcePersistentDiskCsiDriverConfig == nil {
		return true // PDCSI is enabled by default.
	}
	return cluster.AddonsConfig.GcePersistentDiskCsiDriverConfig.Enabled
}

func flattenNotificationConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{}
	if val, ok := c["pubsub"].(map[string]interface{}); ok {
		pubsub := map[string]interface{}{
			"enabled": val["enabled"],
			"topic":   val["topic"],
		}
		if filter, ok := val["filter"].(map[string]interface{}); ok {
			pubsub["filter"] = []map[string]interface{}{
				{
					"event_type": filter["eventType"],
				},
			}
		}
		transformed["pubsub"] = []map[string]interface{}{pubsub}
	}

	return []map[string]interface{}{transformed}
}

func flattenBinaryAuthorization(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"enabled":         c["enabled"],
		"evaluation_mode": c["evaluationMode"],
	}

	return []map[string]interface{}{transformed}
}

func flattenNetworkPolicy(v interface{}) []map[string]interface{} {
	if v == nil {
		// Explicitly set the network policy to the default.
		return []map[string]interface{}{
			{
				"enabled":  false,
				"provider": "PROVIDER_UNSPECIFIED",
			},
		}
	}

	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"enabled":  c["enabled"],
		"provider": c["provider"],
	}

	return []map[string]interface{}{transformed}
}

func flattenClusterAddonsConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]interface{})

	if val, ok := c["horizontalPodAutoscaling"].(map[string]interface{}); ok {
		result["horizontal_pod_autoscaling"] = []map[string]interface{}{
			{
				"disabled": val["disabled"],
			},
		}
	}
	if val, ok := c["httpLoadBalancing"].(map[string]interface{}); ok {
		result["http_load_balancing"] = []map[string]interface{}{
			{
				"disabled": val["disabled"],
			},
		}
	}
	if val, ok := c["networkPolicyConfig"].(map[string]interface{}); ok {
		result["network_policy_config"] = []map[string]interface{}{
			{
				"disabled": val["disabled"],
			},
		}
	}

	if val, ok := c["gcpFilestoreCsiDriverConfig"].(map[string]interface{}); ok {
		result["gcp_filestore_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": val["enabled"],
			},
		}
	}

	if val, ok := c["cloudRunConfig"].(map[string]interface{}); ok {
		cloudRunConfig := map[string]interface{}{
			"disabled": val["disabled"],
		}
		// Currently we only allow setting load_balancer_type to LOAD_BALANCER_TYPE_INTERNAL
		if lbType, ok := val["loadBalancerType"].(string); ok && lbType == "LOAD_BALANCER_TYPE_INTERNAL" {
			cloudRunConfig["load_balancer_type"] = "LOAD_BALANCER_TYPE_INTERNAL"
		}
		result["cloudrun_config"] = []map[string]interface{}{cloudRunConfig}
	}

	if val, ok := c["dnsCacheConfig"].(map[string]interface{}); ok {
		result["dns_cache_config"] = []map[string]interface{}{
			{
				"enabled": val["enabled"],
			},
		}
	}

	if val, ok := c["gcePersistentDiskCsiDriverConfig"].(map[string]interface{}); ok {
		result["gce_persistent_disk_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": val["enabled"],
			},
		}
	}
	if val, ok := c["gkeBackupAgentConfig"].(map[string]interface{}); ok {
		result["gke_backup_agent_config"] = []map[string]interface{}{
			{
				"enabled": val["enabled"],
			},
		}
	}
	if val, ok := c["configConnectorConfig"].(map[string]interface{}); ok {
		result["config_connector_config"] = []map[string]interface{}{
			{
				"enabled": val["enabled"],
			},
		}
	}
	if val, ok := c["gcsFuseCsiDriverConfig"].(map[string]interface{}); ok {
		result["gcs_fuse_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": val["enabled"],
			},
		}
	}
	if val, ok := c["statefulHaConfig"].(map[string]interface{}); ok {
		result["stateful_ha_config"] = []map[string]interface{}{
			{
				"enabled": val["enabled"],
			},
		}
	}
	if val, ok := c["rayOperatorConfig"].(map[string]interface{}); ok {
		rayConfig := []map[string]interface{}{
			{
				"enabled": val["enabled"],
			},
		}
		if logging, ok := val["rayClusterLoggingConfig"].(map[string]interface{}); ok {
			rayConfig[0]["ray_cluster_logging_config"] = []map[string]interface{}{{
				"enabled": logging["enabled"],
			}}
		}
		if monitoring, ok := val["rayClusterMonitoringConfig"].(map[string]interface{}); ok {
			rayConfig[0]["ray_cluster_monitoring_config"] = []map[string]interface{}{{
				"enabled": monitoring["enabled"],
			}}
		}
		result["ray_operator_config"] = rayConfig
	}
	if val, ok := c["parallelstoreCsiDriverConfig"].(map[string]interface{}); ok {
		result["parallelstore_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": val["enabled"],
			},
		}
	}
	if val, ok := c["lustreCsiDriverConfig"].(map[string]interface{}); ok {
		result["lustre_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled":                   val["enabled"],
				"enable_legacy_lustre_port": val["enableLegacyLustrePort"],
			},
		}
	}

	return []map[string]interface{}{result}
}

func flattenAuthenticatorGroupsConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"security_group": c["securityGroup"],
	}

	return []map[string]interface{}{transformed}
}

func flattenControlPlaneEndpointsConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"dns_endpoint_config": flattenDnsEndpointConfig(c["dnsEndpointConfig"]),
		"ip_endpoints_config": flattenIpEndpointsConfig(c["ipEndpointsConfig"]),
	}

	return []map[string]interface{}{transformed}
}

func flattenDnsEndpointConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	dns, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"endpoint":                  dns["endpoint"],
		"allow_external_traffic":    dns["allowExternalTraffic"],
		"enable_k8s_tokens_via_dns": dns["enableK8sTokensViaDns"],
		"enable_k8s_certs_via_dns":  dns["enableK8sCertsViaDns"],
	}

	return []map[string]interface{}{transformed}
}

func flattenIpEndpointsConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	ip, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"enabled": ip["enabled"],
	}

	return []map[string]interface{}{transformed}
}

// Most of PrivateClusterConfig has moved to ControlPlaneEndpointsConfig.
// Most of PrivateClusterConfig has moved to ControlPlaneEndpointsConfig.
func flattenPrivateClusterConfig(cpec, pcc, nc interface{}) []map[string]interface{} {
	if cpec == nil && pcc == nil && nc == nil {
		return nil
	}

	r := map[string]interface{}{}
	if cpec != nil {
		if c, ok := cpec.(map[string]interface{}); ok {
			// Note the change in semantics from private to public endpoint.
			if ipEndpointsConfig, ok := c["ipEndpointsConfig"].(map[string]interface{}); ok {
				if v, ok := ipEndpointsConfig["enablePublicEndpoint"].(bool); ok {
					r["enable_private_endpoint"] = !v
				}
				r["private_endpoint"] = ipEndpointsConfig["privateEndpoint"]
				r["private_endpoint_subnetwork"] = ipEndpointsConfig["privateEndpointSubnetwork"]
				r["public_endpoint"] = ipEndpointsConfig["publicEndpoint"]
				if v, ok := ipEndpointsConfig["globalAccess"].(bool); ok {
					r["master_global_access_config"] = []map[string]interface{}{
						{
							"enabled": v,
						},
					}
				}
			}
		}
	}
	// This is the only field that is canonically still in the PrivateClusterConfig message.
	if pcc != nil {
		if c, ok := pcc.(map[string]interface{}); ok {
			r["peering_name"] = c["peeringName"]
			r["master_ipv4_cidr_block"] = c["masterIpv4CidrBlock"]
		}
	}
	if nc != nil {
		if c, ok := nc.(map[string]interface{}); ok {
			r["enable_private_nodes"] = c["defaultEnablePrivateNodes"]
		}
	}

	return []map[string]interface{}{r}
}

func flattenVerticalPodAutoscaling(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"enabled": c["enabled"],
	}

	return []map[string]interface{}{transformed}
}

func flattenReleaseChannel(v interface{}) []map[string]interface{} {
	if v == nil {
		// Explicitly set the release channel to the UNSPECIFIED.
		return []map[string]interface{}{
			{
				"channel": "UNSPECIFIED",
			},
		}
	}

	if c, ok := v.(map[string]interface{}); ok && c["channel"] != "" {
		transformed := map[string]interface{}{
			"channel": c["channel"],
		}
		return []map[string]interface{}{transformed}
	}

	// Explicitly set the release channel to the UNSPECIFIED.
	return []map[string]interface{}{
		{
			"channel": "UNSPECIFIED",
		},
	}
}

func flattenGkeAutoUpgradeConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{}
	if val, ok := c["patchMode"].(string); ok && val != "" {
		transformed["patch_mode"] = val
	}

	return []map[string]interface{}{transformed}
}

func flattenDefaultSnatStatus(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"disabled": c["disabled"],
	}

	return []map[string]interface{}{transformed}
}

func flattenWorkloadIdentityConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"workload_pool": c["workloadPool"],
	}

	return []map[string]interface{}{transformed}
}

func flattenIdentityServiceConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"enabled": c["enabled"],
	}

	return []map[string]interface{}{transformed}
}

func flattenPodCidrOverprovisionConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"disabled": c["disable"],
	}

	return []map[string]interface{}{transformed}
}

func flattenAdditionalIpRangesConfigs(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.([]interface{})
	if !ok {
		return nil
	}
	if len(c) == 0 {
		return nil
	}

	var outRanges []map[string]interface{}
	for _, raw := range c {
		rangeConfig, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		outRangeConfig := map[string]interface{}{
			"subnetwork":           rangeConfig["subnetwork"],
			"pod_ipv4_range_names": rangeConfig["podIpv4RangeNames"],
		}
		outRanges = append(outRanges, outRangeConfig)
	}

	return outRanges
}

func flattenNetworkTierConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}

	ntc, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"network_tier": ntc["networkTier"],
	}

	return []map[string]interface{}{transformed}
}

func flattenIPAllocationPolicy(v interface{}, d *schema.ResourceData, config *transport.Config) ([]map[string]interface{}, error) {
	// If IP aliasing isn't enabled, none of the values in this block can be set.
	if v == nil {
		if d != nil {
			if err := d.Set("networking_mode", "ROUTES"); err != nil {
				return nil, fmt.Errorf("Error setting networking_mode: %s", err)
			}
		}
		return nil, nil
	}

	p, ok := v.(map[string]interface{})
	if !ok {
		return nil, nil // Should not happen if v is not nil and correct type is passed
	}

	useIpAliases, ok := p["useIpAliases"].(bool)
	if !ok || !useIpAliases {
		if d != nil {
			if err := d.Set("networking_mode", "ROUTES"); err != nil {
				return nil, fmt.Errorf("Error setting networking_mode: %s", err)
			}
		}
		return nil, nil
	}

	if d != nil {
		if err := d.Set("networking_mode", "VPC_NATIVE"); err != nil {
			return nil, fmt.Errorf("Error setting networking_mode: %s", err)
		}
	}

	// handle older clusters that return JSON null
	// corresponding to "STACK_TYPE_UNSPECIFIED" due to GKE declining to backfill
	// equivalent to default_if_empty
	stackType, _ := p["stackType"].(string)
	if stackType == "" {
		stackType = "IPV4"
	}

	return []map[string]interface{}{
		{
			"cluster_ipv4_cidr_block":       p["clusterIpv4CidrBlock"],
			"services_ipv4_cidr_block":      p["servicesIpv4CidrBlock"],
			"cluster_secondary_range_name":  p["clusterSecondaryRangeName"],
			"services_secondary_range_name": p["servicesSecondaryRangeName"],
			"stack_type":                    stackType,
			"pod_cidr_overprovision_config": flattenPodCidrOverprovisionConfig(p["podCidrOverprovisionConfig"]),
			"additional_pod_ranges_config":  flattenAdditionalPodRangesConfig(p),
			"additional_ip_ranges_config":   flattenAdditionalIpRangesConfigs(p["additionalIpRangesConfigs"]),
			"auto_ipam_config":              flattenAutoIpamConfig(p["autoIpamConfig"]),
			"network_tier_config":           flattenNetworkTierConfig(p["networkTierConfig"]),
		},
	}, nil
}

func flattenAutoIpamConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	aic, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"enabled": aic["enabled"],
	}

	return []map[string]interface{}{transformed}
}

func flattenMaintenancePolicy(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	mp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	window, ok := mp["window"].(map[string]interface{})
	if !ok || window == nil {
		return nil
	}

	exclusions := []map[string]interface{}{}
	if maintenanceExclusions, ok := window["maintenanceExclusions"].(map[string]interface{}); ok && maintenanceExclusions != nil {
		for wName, wVal := range maintenanceExclusions {
			windowVal, ok := wVal.(map[string]interface{})
			if !ok {
				continue
			}
			exclusion := map[string]interface{}{
				"start_time":     windowVal["startTime"],
				"exclusion_name": wName,
			}

			if opts, ok := windowVal["maintenanceExclusionOptions"].(map[string]interface{}); ok && opts != nil {
				scope := "NO_UPGRADES"
				if s, ok := opts["scope"].(string); ok && s != "" {
					scope = s
				}
				exclusion["exclusion_options"] = []map[string]interface{}{
					{
						"scope": scope,
					},
				}
				if endTime, _ := windowVal["endTime"].(string); endTime != "" {
					exclusion["end_time"] = endTime
				}
			} else {
				if endTime, _ := windowVal["endTime"].(string); endTime != "" {
					exclusion["end_time"] = endTime
				}
			}
			exclusions = append(exclusions, exclusion)
		}
	}

	transformed := map[string]interface{}{
		"maintenance_exclusion": exclusions,
	}

	if dailyMaintenanceWindow, ok := window["dailyMaintenanceWindow"].(map[string]interface{}); ok && dailyMaintenanceWindow != nil {
		transformed["daily_maintenance_window"] = []map[string]interface{}{
			{
				"start_time": dailyMaintenanceWindow["startTime"],
				"duration":   dailyMaintenanceWindow["duration"],
			},
		}
	}

	if recurringWindow, ok := window["recurringWindow"].(map[string]interface{}); ok && recurringWindow != nil {
		// recurringWindow has nested window
		rwWindow, _ := recurringWindow["window"].(map[string]interface{})

		windowMap := map[string]interface{}{}
		if rwWindow != nil {
			windowMap["start_time"] = rwWindow["startTime"]
			windowMap["end_time"] = rwWindow["endTime"]
		}
		windowMap["recurrence"] = recurringWindow["recurrence"]

		transformed["recurring_window"] = []map[string]interface{}{windowMap}
	}

	return []map[string]interface{}{transformed}
}

func flattenMasterAuth(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}

	a, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"client_certificate":     a["clientCertificate"],
		"client_key":             a["clientKey"],
		"cluster_ca_certificate": a["clusterCaCertificate"],
	}

	// No version of the GKE API returns the client_certificate_config value.
	// Instead, we need to infer whether or not it was set based on the
	// client cert being returned from the API or not.
	// Previous versions of the provider didn't record anything in state when
	// a client cert was enabled, only setting the block when it was false.
	if cert, ok := a["clientCertificate"].(string); ok {
		transformed["client_certificate_config"] = []map[string]interface{}{
			{
				"issue_client_certificate": len(cert) != 0,
			},
		}
	} else {
		// if clientCertificate is missing or not string, maybe default to false?
		// Original code: len(a["clientCertificate"].(string)) != 0
		// This would panic if nil or not string.
		// Safe verification:
		transformed["client_certificate_config"] = []map[string]interface{}{
			{
				"issue_client_certificate": false,
			},
		}
	}

	return []map[string]interface{}{transformed}
}

func flattenClusterAutoscaling(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	a, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := make(map[string]interface{})

	if val, ok := a["enableNodeAutoprovisioning"].(bool); ok && val {
		resourceLimits := make([]interface{}, 0)
		if rl, ok := a["resourceLimits"].([]interface{}); ok {
			for _, item := range rl {
				if limit, ok := item.(map[string]interface{}); ok {
					resourceLimits = append(resourceLimits, map[string]interface{}{
						"resource_type": limit["resourceType"],
						"minimum":       limit["minimum"],
						"maximum":       limit["maximum"],
					})
				}
			}
		}
		transformed["resource_limits"] = resourceLimits
		transformed["enabled"] = true
		transformed["auto_provisioning_defaults"] = flattenAutoProvisioningDefaults(a["autoprovisioningNodePoolDefaults"])
		transformed["auto_provisioning_locations"] = a["autoprovisioningLocations"]
	} else {
		transformed["enabled"] = false
	}
	transformed["autoscaling_profile"] = a["autoscalingProfile"]
	if dccc, ok := a["defaultComputeClassConfig"].(map[string]interface{}); ok {
		transformed["default_compute_class_enabled"] = dccc["enabled"]
	}

	return []map[string]interface{}{transformed}
}

func flattenAutoProvisioningDefaults(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	a, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["oauth_scopes"] = a["oauthScopes"]
	transformed["service_account"] = a["serviceAccount"]
	transformed["disk_size"] = a["diskSizeGb"]
	transformed["disk_type"] = a["diskType"]
	transformed["image_type"] = a["imageType"]
	transformed["min_cpu_platform"] = a["minCpuPlatform"]
	transformed["boot_disk_kms_key"] = a["bootDiskKmsKey"]
	transformed["shielded_instance_config"] = flattenShieldedInstanceConfig(a["shieldedInstanceConfig"])
	transformed["management"] = flattenManagement(a["management"])
	transformed["upgrade_settings"] = flattenUpgradeSettings(a["upgradeSettings"])

	return []map[string]interface{}{transformed}
}

func flattenUpgradeSettings(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	a, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["max_surge"] = a["maxSurge"]
	transformed["max_unavailable"] = a["maxUnavailable"]
	transformed["strategy"] = a["strategy"]
	transformed["blue_green_settings"] = flattenBlueGreenSettings(a["blueGreenSettings"])

	return []map[string]interface{}{transformed}
}

func flattenBlueGreenSettings(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	a, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["node_pool_soak_duration"] = a["nodePoolSoakDuration"]
	transformed["standard_rollout_policy"] = flattenStandardRolloutPolicy(a["standardRolloutPolicy"])

	return []map[string]interface{}{transformed}
}

func flattenStandardRolloutPolicy(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	a, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["batch_percentage"] = a["batchPercentage"]
	transformed["batch_node_count"] = a["batchNodeCount"]
	transformed["batch_soak_duration"] = a["batchSoakDuration"]

	return []map[string]interface{}{transformed}
}

func flattenManagement(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	a, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["auto_upgrade"] = a["autoUpgrade"]
	transformed["auto_repair"] = a["autoRepair"]
	transformed["upgrade_options"] = flattenUpgradeOptions(a["upgradeOptions"])

	return []map[string]interface{}{transformed}
}

func flattenUpgradeOptions(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	a, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["auto_upgrade_start_time"] = a["autoUpgradeStartTime"]
	transformed["description"] = a["description"]

	return []map[string]interface{}{transformed}
}

func flattenMasterAuthorizedNetworksConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	if val, ok := c["enabled"].(bool); !ok || !val {
		return nil
	}

	transformed := make(map[string]interface{})
	cidrBlocks := make([]interface{}, 0)
	if cb, ok := c["cidrBlocks"].([]interface{}); ok {
		for _, item := range cb {
			if v, ok := item.(map[string]interface{}); ok {
				cidrBlocks = append(cidrBlocks, map[string]interface{}{
					"cidr_block":   v["cidrBlock"],
					"display_name": v["displayName"],
				})
			}
		}
	}

	transformed["cidr_blocks"] = schema.NewSet(schema.HashResource(cidrBlockConfig), cidrBlocks)
	transformed["gcp_public_cidrs_access_enabled"] = c["gcpPublicCidrsAccessEnabled"]
	transformed["private_endpoint_enforcement_enabled"] = c["privateEndpointEnforcementEnabled"]

	return []map[string]interface{}{transformed}
}

func flattenPodAutoscaling(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"hpa_profile": c["hpaProfile"],
	}

	return []map[string]interface{}{transformed}
}

func flattenSecretManagerConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return []map[string]interface{}{
			{
				"enabled": false,
			},
		}
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	result["enabled"] = c["enabled"]

	rotationList := []map[string]interface{}{}
	if rotationConfig, ok := c["rotationConfig"].(map[string]interface{}); ok && rotationConfig != nil {
		rotationConfigMap := map[string]interface{}{
			"enabled": rotationConfig["enabled"],
		}
		if interval, ok := rotationConfig["rotationInterval"].(string); ok && interval != "" {
			rotationConfigMap["rotation_interval"] = interval
		}
		rotationList = append(rotationList, rotationConfigMap)
	}
	result["rotation_config"] = rotationList
	return []map[string]interface{}{result}
}

func flattenResourceUsageExportConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	enableResourceConsumptionMetering := false
	if cmc, ok := c["consumptionMeteringConfig"].(map[string]interface{}); ok && cmc != nil {
		if val, ok := cmc["enabled"].(bool); ok && val {
			enableResourceConsumptionMetering = true
		}
	}

	transformed := map[string]interface{}{
		"enable_network_egress_metering":       c["enableNetworkEgressMetering"],
		"enable_resource_consumption_metering": enableResourceConsumptionMetering,
	}

	if bqd, ok := c["bigqueryDestination"].(map[string]interface{}); ok && bqd != nil {
		transformed["bigquery_destination"] = []map[string]interface{}{
			{"dataset_id": bqd["datasetId"]},
		}
	}

	return []map[string]interface{}{transformed}
}

func flattenServiceExternalIpsConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"enabled": c["enabled"],
	}

	return []map[string]interface{}{transformed}
}

func flattenMeshCertificates(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"enable_certificates": c["enableCertificates"],
	}

	return []map[string]interface{}{transformed}
}

func flattenManagementConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"enabled": c["enabled"],
	}

	return []map[string]interface{}{transformed}
}

func flattenDatabaseEncryption(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"state":    c["state"],
		"key_name": c["keyName"],
	}

	return []map[string]interface{}{transformed}
}

func flattenDnsConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"additive_vpc_scope_dns_domain": c["additiveVpcScopeDnsDomain"],
		"cluster_dns":                   c["clusterDns"],
		"cluster_dns_scope":             c["clusterDnsScope"],
		"cluster_dns_domain":            c["clusterDnsDomain"],
	}

	return []map[string]interface{}{transformed}
}

func flattenNetworkPerformanceConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"total_egress_bandwidth_tier": c["totalEgressBandwidthTier"],
	}

	return []map[string]interface{}{transformed}
}

func flattenGatewayApiConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"channel": c["channel"],
	}

	return []map[string]interface{}{transformed}
}

func flattenFleet(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	membership, _ := c["membership"].(string)

	// Parse membership_id and membership_location from full membership name.
	var membership_id, membership_location string
	membershipRE := regexp.MustCompile(`^(//[a-zA-Z0-9\.\-]+)?/?projects/([^/]+)/locations/([a-zA-Z0-9\-]+)/memberships/([^/]+)$`)
	if match := membershipRE.FindStringSubmatch(membership); match != nil {
		membership_id = match[4]
		membership_location = match[3]
	}

	transformed := map[string]interface{}{
		"project":             c["project"],
		"membership":          membership,
		"membership_id":       membership_id,
		"membership_location": membership_location,
		"pre_registered":      c["preRegistered"],
		"membership_type":     c["membershipType"],
	}

	return []map[string]interface{}{transformed}
}

func flattenUserManagedKeysConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	f := map[string]interface{}{
		"cluster_ca":                        c["clusterCa"],
		"etcd_api_ca":                       c["etcdApiCa"],
		"etcd_peer_ca":                      c["etcdPeerCa"],
		"aggregation_ca":                    c["aggregationCa"],
		"control_plane_disk_encryption_key": c["controlPlaneDiskEncryptionKey"],
		"gkeops_etcd_backup_encryption_key": c["gkeopsEtcdBackupEncryptionKey"],
	}
	allEmpty := true
	for _, v := range f {
		if s, ok := v.(string); ok && s != "" {
			allEmpty = false
		}
	}

	if keys, ok := c["serviceAccountSigningKeys"].([]string); ok && len(keys) != 0 {
		f["service_account_signing_keys"] = keys
		allEmpty = false
	}
	if keys, ok := c["serviceAccountVerificationKeys"].([]string); ok && len(keys) != 0 {
		f["service_account_verification_keys"] = keys
		allEmpty = false
	}
	if allEmpty {
		return nil
	}
	return []map[string]interface{}{f}
}

func flattenEnableK8sBetaApis(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"enabled_apis": c["enabledApis"],
	}

	return []map[string]interface{}{transformed}
}

func flattenContainerClusterLoggingConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{}
	if componentConfig, ok := c["componentConfig"].(map[string]interface{}); ok && componentConfig != nil {
		transformed["enable_components"] = componentConfig["enableComponents"]
	}

	return []map[string]interface{}{transformed}
}

func flattenMonitoringConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := make(map[string]interface{})
	if componentConfig, ok := c["componentConfig"].(map[string]interface{}); ok && componentConfig != nil {
		transformed["enable_components"] = componentConfig["enableComponents"]
	}
	if managedPrometheusConfig, ok := c["managedPrometheusConfig"].(map[string]interface{}); ok && managedPrometheusConfig != nil {
		transformed["managed_prometheus"] = flattenManagedPrometheusConfig(managedPrometheusConfig)
	}
	if advancedDatapathObservabilityConfig, ok := c["advancedDatapathObservabilityConfig"].(map[string]interface{}); ok && advancedDatapathObservabilityConfig != nil {
		transformed["advanced_datapath_observability_config"] = flattenAdvancedDatapathObservabilityConfig(advancedDatapathObservabilityConfig)
	}

	return []map[string]interface{}{transformed}
}

func flattenAdvancedDatapathObservabilityConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"enable_metrics": c["enableMetrics"],
		"enable_relay":   c["enableRelay"],
	}

	return []map[string]interface{}{transformed}
}

func flattenManagedPrometheusConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["enabled"] = c["enabled"]

	autoMonitoringList := []map[string]interface{}{}
	if amc, ok := c["autoMonitoringConfig"].(map[string]interface{}); ok && amc != nil {
		if scope, ok := amc["scope"].(string); ok && scope != "" {
			autoMonitoringMap := map[string]interface{}{
				"scope": scope,
			}
			autoMonitoringList = append(autoMonitoringList, autoMonitoringMap)
		}
	}

	transformed["auto_monitoring_config"] = autoMonitoringList

	return []map[string]interface{}{transformed}
}

func flattenNodePoolAutoConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := make(map[string]interface{})
	if nodeKubeletConfig, ok := c["nodeKubeletConfig"].(map[string]interface{}); ok && nodeKubeletConfig != nil {
		transformed["node_kubelet_config"] = flattenNodePoolAutoConfigNodeKubeletConfig(nodeKubeletConfig)
	}
	if networkTags, ok := c["networkTags"].(map[string]interface{}); ok && networkTags != nil {
		transformed["network_tags"] = flattenNodePoolAutoConfigNetworkTags(networkTags)
	}
	if resourceManagerTags, ok := c["resourceManagerTags"].(map[string]interface{}); ok && resourceManagerTags != nil {
		transformed["resource_manager_tags"] = flattenResourceManagerTags(resourceManagerTags)
	}
	if linuxNodeConfig, ok := c["linuxNodeConfig"].(map[string]interface{}); ok && linuxNodeConfig != nil {
		transformed["linux_node_config"] = []map[string]interface{}{
			{"cgroup_mode": linuxNodeConfig["cgroupMode"]},
		}
	}

	return []map[string]interface{}{transformed}
}

func flattenNodePoolAutoConfigNetworkTags(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := make(map[string]interface{})
	if tags, ok := c["tags"].([]interface{}); ok {
		transformed["tags"] = tags
	}

	return []map[string]interface{}{transformed}
}

func flattenContainerClusterNodePools(d *schema.ResourceData, config *transport.Config, v interface{}) ([]map[string]interface{}, error) {
	if v == nil {
		return nil, nil
	}
	nodePools, ok := v.([]interface{})
	if !ok {
		return nil, nil
	}
	result := make([]map[string]interface{}, 0, len(nodePools))
	for _, np := range nodePools {
		if npMap, ok := np.(map[string]interface{}); ok {
			nodePool, err := flattenNodePool(d, config, npMap, "")
			if err != nil {
				return nil, err
			}
			result = append(result, nodePool)
		}
	}
	return result, nil
}

func flattenRBACBindingConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"enable_insecure_binding_system_authenticated":   c["enableInsecureBindingSystemAuthenticated"],
		"enable_insecure_binding_system_unauthenticated": c["enableInsecureBindingSystemUnauthenticated"],
	}

	return []map[string]interface{}{transformed}
}
