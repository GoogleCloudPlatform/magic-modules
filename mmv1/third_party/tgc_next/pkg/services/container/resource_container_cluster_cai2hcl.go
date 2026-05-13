package container

import (
	"fmt"
	"slices"
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
	hclData["location"] = asset.Resource.Data["location"]
	if v := asset.Resource.Data["network"]; v != nil && v != "default" {
		hclData["network"] = v
	}
	hclData["subnetwork"] = asset.Resource.Data["subnetwork"]
	hclData["initial_node_count"] = asset.Resource.Data["initialNodeCount"]
	if legacyAbac, ok := asset.Resource.Data["legacyAbac"].(map[string]interface{}); ok {
		if v := legacyAbac["enabled"]; v != nil && v != false {
			hclData["enable_legacy_abac"] = v
		}
	}
	if autopilot, ok := asset.Resource.Data["autopilot"].(map[string]interface{}); ok {
		if enabled, _ := autopilot["enabled"].(bool); enabled {
			hclData["enable_autopilot"] = true
		}
		if workloadPolicyConfig, ok := autopilot["workloadPolicyConfig"].(map[string]interface{}); ok {
			hclData["allow_net_admin"] = workloadPolicyConfig["allowNetAdmin"]
		}

		if clusterPolicyConfig, ok := autopilot["clusterPolicyConfig"].(map[string]interface{}); ok {
			policyConfig := map[string]interface{}{}
			if v, ok := clusterPolicyConfig["noStandardNodePools"]; ok {
				policyConfig["no_standard_node_pools"] = v
			}
			if v, ok := clusterPolicyConfig["noSystemImpersonation"]; ok {
				policyConfig["no_system_impersonation"] = v
			}
			if v, ok := clusterPolicyConfig["noSystemMutation"]; ok {
				policyConfig["no_system_mutation"] = v
			}
			if v, ok := clusterPolicyConfig["noUnsafeWebhooks"]; ok {
				policyConfig["no_unsafe_webhooks"] = v
			}
			if len(policyConfig) > 0 {
				hclData["autopilot_cluster_policy_config"] = []map[string]interface{}{policyConfig}
			}
		}
		if privilegedAdmissionConfig, ok := autopilot["privilegedAdmissionConfig"].(map[string]interface{}); ok {
			hclData["autopilot_privileged_admission"] = privilegedAdmissionConfig["allowlistPaths"]
		}
	}
	enableAutopilot, _ := hclData["enable_autopilot"].(bool)
	if shieldedNodes, ok := asset.Resource.Data["shieldedNodes"].(map[string]interface{}); ok && !enableAutopilot {
		if v := shieldedNodes["enabled"]; v != nil && v != true {
			hclData["enable_shielded_nodes"] = v
		}
	}
	hclData["confidential_nodes"] = flattenConfidentialNodes(asset.Resource.Data["confidentialNodes"])
	if locations, ok := asset.Resource.Data["locations"].([]interface{}); ok {
		idx := -1
		for i, location := range locations {
			if locationString, ok := location.(string); ok {
				if hclLocation, ok := hclData["location"].(string); ok && locationString == hclLocation {
					idx = i
				}
			}
		}
		if idx != -1 {
			locations = slices.Delete(locations, idx, idx+1)
		}
		if len(locations) != 0 {
			hclData["node_locations"] = locations
		}
	}
	hclData["logging_service"] = asset.Resource.Data["loggingService"]
	hclData["monitoring_service"] = asset.Resource.Data["monitoringService"]
	hclData["node_config"] = flattenNodeConfig(asset.Resource.Data["nodeConfig"], nil)
	hclData["description"] = asset.Resource.Data["description"]
	hclData["security_posture_config"] = flattenSecurityPostureConfig(asset.Resource.Data["securityPostureConfig"])
	hclData["enterprise_config"] = flattenEnterpriseConfig(asset.Resource.Data["enterpriseConfig"])
	hclData["anonymous_authentication_config"] = flattenAnonymousAuthenticationConfig(asset.Resource.Data["anonymousAuthenticationConfig"])
	hclData["notification_config"] = flattenNotificationConfig(asset.Resource.Data["notificationConfig"])
	hclData["binary_authorization"] = flattenBinaryAuthorization(asset.Resource.Data["binaryAuthorization"])
	if !enableAutopilot {
		hclData["network_policy"] = flattenNetworkPolicy(asset.Resource.Data["networkPolicy"])
	}
	hclData["addons_config"] = flattenClusterAddonsConfig(asset.Resource.Data["addonsConfig"], enableAutopilot)
	if !enableAutopilot {
		hclData["node_pool"], err = flattenContainerClusterNodePools(d, config, asset.Resource.Data["nodePools"])
	}
	hclData["node_pool_defaults"] = flattenNodePoolDefaults(asset.Resource.Data["nodePoolDefaults"])
	hclData["authenticator_groups_config"] = flattenAuthenticatorGroupsConfig(asset.Resource.Data["authenticatorGroupsConfig"])
	hclData["control_plane_endpoints_config"] = flattenControlPlaneEndpointsConfig(asset.Resource.Data["controlPlaneEndpointsConfig"])

	privateClusterConfig := asset.Resource.Data["privateClusterConfig"]
	controlPlaneEndpointsConfig := asset.Resource.Data["controlPlaneEndpointsConfig"]
	networkConfig := asset.Resource.Data["networkConfig"]
	if v := asset.Resource.Data["enableKubernetesAlpha"]; v != nil && v != false {
		hclData["enable_kubernetes_alpha"] = v
	}
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
		if v := nc["enableCiliumClusterwideNetworkPolicy"]; v != nil && v != false {
			hclData["enable_cilium_clusterwide_network_policy"] = v
		}
		if !enableAutopilot {
			hclData["enable_intranode_visibility"] = nc["enableIntraNodeVisibility"]
		}
		hclData["private_ipv6_google_access"] = nc["privateIpv6GoogleAccess"]
		hclData["datapath_provider"] = nc["datapathProvider"]
		if v := nc["enableMultiNetworking"]; v != nil && v != false {
			hclData["enable_multi_networking"] = v
		}
		hclData["enable_l4_ilb_subsetting"] = nc["enableL4ilbSubsetting"]
		if v := nc["disableL4LbFirewallReconciliation"]; v != nil && v != false {
			hclData["disable_l4_lb_firewall_reconciliation"] = v
		}
		hclData["in_transit_encryption_config"] = nc["inTransitEncryptionConfig"]
		if v := nc["enableFqdnNetworkPolicy"]; v != nil && v != false {
			hclData["enable_fqdn_network_policy"] = v
		}
	}

	hclData["workload_identity_config"] = flattenWorkloadIdentityConfig(asset.Resource.Data["workloadIdentityConfig"])
	hclData["identity_service_config"] = flattenIdentityServiceConfig(asset.Resource.Data["identityServiceConfig"])

	if ipAlloc, ok := asset.Resource.Data["ipAllocationPolicy"].(map[string]interface{}); ok {
		hclData["pod_cidr_overprovision_config"] = flattenPodCidrOverprovisionConfig(ipAlloc["podCidrOverprovisionConfig"])

		ipPolicy, err := flattenIPAllocationPolicy(ipAlloc, nil, nil)
		if err != nil {
			return nil, err
		}

		hclData["ip_allocation_policy"] = ipPolicy
	}

	if ipAlloc, ok := asset.Resource.Data["ipAllocationPolicy"].(map[string]interface{}); !ok || ipAlloc == nil {
		hclData["networking_mode"] = "ROUTES"
	} else if useIpAliases, ok := ipAlloc["useIpAliases"].(bool); !ok || !useIpAliases {
		hclData["networking_mode"] = "ROUTES"
	} else {
		hclData["networking_mode"] = "VPC_NATIVE"
	}

	hclData["maintenance_policy"] = flattenMaintenancePolicy(asset.Resource.Data["maintenancePolicy"])
	hclData["master_auth"] = flattenMasterAuth(asset.Resource.Data["masterAuth"])
	hclData["cluster_autoscaling"] = flattenClusterAutoscaling(asset.Resource.Data["autoscaling"], enableAutopilot)
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
	transformed := make(map[string]interface{})
	if ec["desiredTier"] != nil {
		transformed["desired_tier"] = ec["desiredTier"]
	}

	if len(transformed) == 0 {
		return nil
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
		"mode": aac["mode"],
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
	if !ok || len(c) == 0 {
		return nil
	}

	transformed := map[string]interface{}{}
	if val, ok := c["pubsub"].(map[string]interface{}); ok {
		enabled, ok := val["enabled"].(bool)
		if !ok {
			enabled = false
		}
		pubsub := map[string]interface{}{
			"enabled": enabled,
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
	if !ok || len(c) == 0 {
		return nil
	}

	transformed := make(map[string]interface{})
	if c["enabled"] != nil {
		transformed["enabled"] = c["enabled"]
	}
	if c["evaluationMode"] != nil {
		transformed["evaluation_mode"] = c["evaluationMode"]
	}

	if len(transformed) == 0 {
		return nil
	}

	return []map[string]interface{}{transformed}
}

func flattenNetworkPolicy(v interface{}) []map[string]interface{} {
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	enabled := false
	if v, ok := c["enabled"]; ok && v != nil {
		enabled = v.(bool)
	}

	transformed := map[string]interface{}{
		"enabled":  enabled,
		"provider": c["provider"],
	}

	return []map[string]interface{}{transformed}
}

func flattenClusterAddonsConfig(v interface{}, enableAutopilot bool) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]interface{})

	if val, ok := c["horizontalPodAutoscaling"].(map[string]interface{}); ok {
		disabled := false
		if v, ok := val["disabled"]; ok && v != nil {
			disabled = v.(bool)
		}
		result["horizontal_pod_autoscaling"] = []map[string]interface{}{
			{
				"disabled": disabled,
			},
		}
	}
	if val, ok := c["httpLoadBalancing"].(map[string]interface{}); ok {
		disabled := false
		if v, ok := val["disabled"]; ok && v != nil {
			disabled = v.(bool)
		}
		result["http_load_balancing"] = []map[string]interface{}{
			{
				"disabled": disabled,
			},
		}
	}
	if val, ok := c["networkPolicyConfig"].(map[string]interface{}); ok && !enableAutopilot {
		disabled := false
		if v, ok := val["disabled"]; ok && v != nil {
			disabled = v.(bool)
		}
		result["network_policy_config"] = []map[string]interface{}{
			{
				"disabled": disabled,
			},
		}
	}

	if val, ok := c["gcpFilestoreCsiDriverConfig"].(map[string]interface{}); ok {
		enabled := false
		if v, ok := val["enabled"]; ok && v != nil {
			enabled = v.(bool)
		}
		result["gcp_filestore_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": enabled,
			},
		}
	}

	if val, ok := c["cloudRunConfig"].(map[string]interface{}); ok {
		disabled := false
		if v, ok := val["disabled"]; ok && v != nil {
			disabled = v.(bool)
		}
		cloudRunConfig := map[string]interface{}{
			"disabled": disabled,
		}
		// Currently we only allow setting load_balancer_type to LOAD_BALANCER_TYPE_INTERNAL
		if lbType, ok := val["loadBalancerType"].(string); ok && lbType == "LOAD_BALANCER_TYPE_INTERNAL" {
			cloudRunConfig["load_balancer_type"] = "LOAD_BALANCER_TYPE_INTERNAL"
		}
		result["cloudrun_config"] = []map[string]interface{}{cloudRunConfig}
	}

	if val, ok := c["dnsCacheConfig"].(map[string]interface{}); ok && !enableAutopilot {
		enabled := false
		if v, ok := val["enabled"]; ok && v != nil {
			enabled = v.(bool)
		}

		result["dns_cache_config"] = []map[string]interface{}{
			{
				"enabled": enabled,
			},
		}
	}

	if val, ok := c["gcePersistentDiskCsiDriverConfig"].(map[string]interface{}); ok {
		enabled := false
		if v, ok := val["enabled"]; ok && v != nil {
			enabled = v.(bool)
		}

		result["gce_persistent_disk_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": enabled,
			},
		}
	}
	if val, ok := c["gkeBackupAgentConfig"].(map[string]interface{}); ok {
		enabled := false
		if v, ok := val["enabled"]; ok && v != nil {
			enabled = v.(bool)
		}

		result["gke_backup_agent_config"] = []map[string]interface{}{
			{
				"enabled": enabled,
			},
		}
	}
	if val, ok := c["configConnectorConfig"].(map[string]interface{}); ok {
		enabled := false
		if v, ok := val["enabled"]; ok && v != nil {
			enabled = v.(bool)
		}

		result["config_connector_config"] = []map[string]interface{}{
			{
				"enabled": enabled,
			},
		}
	}
	if val, ok := c["gcsFuseCsiDriverConfig"].(map[string]interface{}); ok {
		enabled := false
		if v, ok := val["enabled"]; ok && v != nil {
			enabled = v.(bool)
		}

		result["gcs_fuse_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": enabled,
			},
		}
	}
	if val, ok := c["statefulHaConfig"].(map[string]interface{}); ok && !enableAutopilot {
		enabled := false
		if v, ok := val["enabled"]; ok && v != nil {
			enabled = v.(bool)
		}

		result["stateful_ha_config"] = []map[string]interface{}{
			{
				"enabled": enabled,
			},
		}
	}
	if val, ok := c["sliceControllerConfig"].(map[string]interface{}); ok {
		enabled := false
		if v, ok := val["enabled"]; ok && v != nil {
			enabled = v.(bool)
		}

		result["slice_controller_config"] = []map[string]interface{}{
			{
				"enabled": enabled,
			},
		}
	}
	if val, ok := c["rayOperatorConfig"].(map[string]interface{}); ok {
		enabled := false
		if v, ok := val["enabled"]; ok && v != nil {
			enabled = v.(bool)
		}

		rayConfig := []map[string]interface{}{
			{
				"enabled": enabled,
			},
		}
		if logging, ok := val["rayClusterLoggingConfig"].(map[string]interface{}); ok {
			enabled := false
			if v, ok := logging["enabled"]; ok && v != nil {
				enabled = v.(bool)
			}

			rayConfig[0]["ray_cluster_logging_config"] = []map[string]interface{}{{
				"enabled": enabled,
			}}
		}
		if monitoring, ok := val["rayClusterMonitoringConfig"].(map[string]interface{}); ok {
			enabled := false
			if v, ok := monitoring["enabled"]; ok && v != nil {
				enabled = v.(bool)
			}

			rayConfig[0]["ray_cluster_monitoring_config"] = []map[string]interface{}{{
				"enabled": enabled,
			}}
		}
		result["ray_operator_config"] = rayConfig
	}
	if val, ok := c["parallelstoreCsiDriverConfig"].(map[string]interface{}); ok {
		enabled := false
		if v, ok := val["enabled"]; ok && v != nil {
			enabled = v.(bool)
		}

		result["parallelstore_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": enabled,
			},
		}
	}
	if val, ok := c["lustreCsiDriverConfig"].(map[string]interface{}); ok {
		enabled := false
		if v, ok := val["enabled"]; ok && v != nil {
			enabled = v.(bool)
		}

		result["lustre_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled":                   enabled,
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

	secGrp := ""
	if val, ok := c["securityGroup"].(string); ok {
		secGrp = val
	}

	transformed := map[string]interface{}{
		"security_group": secGrp,
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
				if v, ok := ipEndpointsConfig["privateEndpointSubnetwork"]; ok && v != nil {
					r["private_endpoint_subnetwork"] = v
				}

				enabled := false
				if v, ok := ipEndpointsConfig["globalAccess"]; ok && v != nil {
					enabled = v.(bool)
				}

				r["master_global_access_config"] = []map[string]interface{}{
					{
						"enabled": enabled,
					},
				}
			}
		}
	}
	// This is the only field that is canonically still in the PrivateClusterConfig message.
	if pcc != nil {
		if c, ok := pcc.(map[string]interface{}); ok {
			if v, ok := c["masterIpv4CidrBlock"]; ok && v != nil {
				r["master_ipv4_cidr_block"] = v
			}
		}
	}
	if nc != nil {
		if c, ok := nc.(map[string]interface{}); ok {
			if v, ok := c["defaultEnablePrivateNodes"]; ok && v != nil {
				r["enable_private_nodes"] = v
			}
		}
	}

	if len(r) == 0 {
		return nil
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

	if c, ok := v.(map[string]interface{}); ok {
		if ch, ok := c["channel"].(string); ok && ch != "" {
			transformed := map[string]interface{}{
				"channel": ch,
			}
			return []map[string]interface{}{transformed}
		}
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
	if !ok || len(c) == 0 {
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

	workloadPool, ok := c["workloadPool"].(string)
	if !ok || workloadPool == "" {
		return nil
	}

	transformed := map[string]interface{}{
		"workload_pool": workloadPool,
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
	if !ok || len(c) == 0 {
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
			"status":               rangeConfig["status"],
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

	transformed := map[string]interface{}{
		"cluster_secondary_range_name":  p["clusterSecondaryRangeName"],
		"services_secondary_range_name": p["servicesSecondaryRangeName"],
		"pod_cidr_overprovision_config": flattenPodCidrOverprovisionConfig(p["podCidrOverprovisionConfig"]),
		"additional_pod_ranges_config":  flattenAdditionalPodRangesConfig(p),
		"additional_ip_ranges_config":   flattenAdditionalIpRangesConfigs(p["additionalIpRangesConfigs"]),
		"auto_ipam_config":              flattenAutoIpamConfig(p["autoIpamConfig"]),
		"network_tier_config":           flattenNetworkTierConfig(p["networkTierConfig"]),
	}
	if stackType != "" && stackType != "IPV4" {
		transformed["stack_type"] = stackType
	}
	conflicts := false
	for _, field := range ipAllocationRangeFields {
		fieldParts := strings.Split(field, ".")
		if val := transformed[fieldParts[len(fieldParts)-1]]; val != nil {
			conflicts = true
		}
	}
	if !conflicts {
		transformed["cluster_ipv4_cidr_block"] = p["clusterIpv4CidrBlock"]
		transformed["services_ipv4_cidr_block"] = p["servicesIpv4CidrBlock"]
	}
	return []map[string]interface{}{transformed}, nil
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

	transformed := map[string]interface{}{}

	if window, ok := mp["window"].(map[string]interface{}); ok && window != nil {
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
					if endTimeBehavior, ok := opts["endTimeBehavior"].(string); ok && endTimeBehavior != "" {
						exclusion["exclusion_options"] = []map[string]interface{}{
							{
								"scope":             scope,
								"end_time_behavior": endTimeBehavior,
							},
						}
					} else {
						exclusion["exclusion_options"] = []map[string]interface{}{
							{
								"scope": scope,
							},
						}
						if endTime, _ := windowVal["endTime"].(string); endTime != "" {
							exclusion["end_time"] = endTime
						}
					}
				} else {
					if endTime, _ := windowVal["endTime"].(string); endTime != "" {
						exclusion["end_time"] = endTime
					}
				}
				exclusions = append(exclusions, exclusion)
			}
		}

		transformed["maintenance_exclusion"] = exclusions

		if dailyMaintenanceWindow, ok := window["dailyMaintenanceWindow"].(map[string]interface{}); ok && dailyMaintenanceWindow != nil {
			transformed["daily_maintenance_window"] = []map[string]interface{}{
				{
					"start_time": dailyMaintenanceWindow["startTime"],
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
	}

	if disruptionBudget, ok := mp["disruptionBudget"].(map[string]interface{}); ok && disruptionBudget != nil {
		dbMap := map[string]interface{}{}
		if val, ok := disruptionBudget["minorVersionDisruptionInterval"]; ok && val != nil {
			dbMap["minor_version_disruption_interval"] = val
		}
		if val, ok := disruptionBudget["patchVersionDisruptionInterval"]; ok && val != nil {
			dbMap["patch_version_disruption_interval"] = val
		}
		transformed["disruption_budget"] = []map[string]interface{}{dbMap}
	}

	if len(transformed) == 0 {
		return nil
	}

	_, hasDaily := transformed["daily_maintenance_window"]
	_, hasRecurring := transformed["recurring_window"]
	if !hasDaily && !hasRecurring {
		return nil
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

	transformed := make(map[string]interface{}, 0)

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

func flattenClusterAutoscaling(v interface{}, enableAutopilot bool) []map[string]interface{} {
	if v == nil {
		return nil
	}
	a, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := make(map[string]interface{})

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
		if !enableAutopilot {
			transformed["resource_limits"] = resourceLimits
		}
	}

	if val, ok := a["enableNodeAutoprovisioning"].(bool); ok && val {
		if !enableAutopilot {
			transformed["enabled"] = true
		}
		transformed["auto_provisioning_defaults"] = flattenAutoProvisioningDefaults(a["autoprovisioningNodePoolDefaults"])
	}
	if v := a["autoprovisioningLocations"]; v != nil {
		transformed["auto_provisioning_locations"] = v
	}
	if v := a["autoscalingProfile"]; v != nil && v != "BALANCED" {
		transformed["autoscaling_profile"] = v
	}
	if dccc, ok := a["defaultComputeClassConfig"].(map[string]interface{}); ok {
		transformed["default_compute_class_enabled"] = dccc["enabled"]
	}

	if len(transformed) == 0 {
		return nil
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
	if v := a["serviceAccount"]; v != nil && v != "default" {
		transformed["service_account"] = v
	}
	if v := a["diskSizeGb"]; v != nil && v != float64(100) {
		transformed["disk_size"] = v
	}
	if v := a["diskType"]; v != nil && v != "pd-standard" {
		transformed["disk_type"] = v
	}
	if v := a["imageType"]; v != nil && v != "COS_CONTAINERD" {
		transformed["image_type"] = v
	}
	if v := a["minCpuPlatform"]; v != nil {
		transformed["min_cpu_platform"] = v
	}
	if v := a["bootDiskKmsKey"]; v != nil {
		transformed["boot_disk_kms_key"] = v
	}
	if v := flattenShieldedInstanceConfig(a["shieldedInstanceConfig"]); v != nil {
		transformed["shielded_instance_config"] = v
	}
	if v := flattenManagement(a["management"]); v != nil {
		transformed["management"] = v
	}
	if v := flattenUpgradeSettings(a["upgradeSettings"]); v != nil {
		transformed["upgrade_settings"] = v
	}

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
	if val := a["maxSurge"]; val != nil {
		transformed["max_surge"] = val
	}
	if val := a["maxUnavailable"]; val != nil {
		transformed["max_unavailable"] = val
	}
	if val := a["strategy"]; val != nil {
		transformed["strategy"] = val
	}
	if val := flattenBlueGreenSettings(a["blueGreenSettings"]); val != nil {
		transformed["blue_green_settings"] = val
	}

	if len(transformed) == 0 {
		return nil
	}

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
	if val := a["nodePoolSoakDuration"]; val != nil {
		transformed["node_pool_soak_duration"] = val
	}
	if val := flattenStandardRolloutPolicy(a["standardRolloutPolicy"]); val != nil {
		transformed["standard_rollout_policy"] = val
	}

	if len(transformed) == 0 {
		return nil
	}

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
	if v := a["batchSoakDuration"]; v != nil && v != "0s" {
		transformed["batch_soak_duration"] = v
	}

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
	if v := a["autoUpgrade"]; v != nil {
		transformed["auto_upgrade"] = v
	}
	if v := a["autoRepair"]; v != nil {
		transformed["auto_repair"] = v
	}
	if len(transformed) == 0 {
		return nil
	}

	return []map[string]interface{}{transformed}
}

func flattenUpgradeOptions(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	if _, ok := v.(map[string]interface{}); !ok {
		return nil
	}

	// upgrade_options block is Computed-only.
	// Specifically, auto_upgrade_start_time and description are Computed.
	// We omit them here to prevent 'Value for unconfigurable attribute' errors.
	transformed := make(map[string]interface{})

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

	transformed["cidr_blocks"] = cidrBlocks
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

	hpaProfile, ok := c["hpaProfile"].(string)
	if !ok || hpaProfile == "" || hpaProfile == "HPA_PROFILE_UNSPECIFIED" {
		return nil
	}

	transformed := map[string]interface{}{
		"hpa_profile": hpaProfile,
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
	if val, ok := c["enabled"]; ok && val != nil {
		result["enabled"] = val
	} else {
		result["enabled"] = false
	}

	rotationList := []map[string]interface{}{}
	if rotationConfig, ok := c["rotationConfig"].(map[string]interface{}); ok && rotationConfig != nil {
		rotationConfigMap := map[string]interface{}{}
		if rVal, ok := rotationConfig["enabled"]; ok && rVal != nil {
			rotationConfigMap["enabled"] = rVal
		} else {
			rotationConfigMap["enabled"] = false
		}

		if interval, ok := rotationConfig["rotationInterval"].(string); ok && interval != "" {
			rotationConfigMap["rotation_interval"] = interval
		}
		rotationList = append(rotationList, rotationConfigMap)
	}

	if len(rotationList) > 0 {
		result["rotation_config"] = rotationList
	}
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

	transformed := map[string]interface{}{}
	if v := c["enableNetworkEgressMetering"]; v != nil && v != false {
		transformed["enable_network_egress_metering"] = v
	}
	if enableResourceConsumptionMetering != true {
		transformed["enable_resource_consumption_metering"] = enableResourceConsumptionMetering
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

	enabled := false
	if v, ok := c["enabled"]; ok && v != nil {
		enabled = v.(bool)
	}

	transformed := map[string]interface{}{
		"enabled": enabled,
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

	enabled := false
	if val, ok := c["enabled"].(bool); ok {
		enabled = val
	}

	transformed := map[string]interface{}{
		"enabled": enabled,
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
		"cluster_dns_scope":             c["clusterDnsScope"],
		"cluster_dns_domain":            c["clusterDnsDomain"],
	}
	if v := c["clusterDns"]; v != nil && v != "PROVIDER_UNSPECIFIED" {
		transformed["cluster_dns"] = v
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

	transformed := map[string]interface{}{}
	if val, ok := c["project"]; ok && val != "" {
		transformed["project"] = val
	}
	if val, ok := c["membershipType"]; ok && val != "" {
		transformed["membership_type"] = val
	}

	if len(transformed) == 0 {
		return nil
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

	if keys, ok := c["serviceAccountSigningKeys"].([]interface{}); ok && len(keys) != 0 {
		f["service_account_signing_keys"] = keys
		allEmpty = false
	}
	if keys, ok := c["serviceAccountVerificationKeys"].([]interface{}); ok && len(keys) != 0 {
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
		if len(componentConfig) == 0 {
			return nil
		}
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
	if advancedDatapathObservabilityConfig, ok := c["advancedDatapathObservabilityConfig"].(map[string]interface{}); ok && len(advancedDatapathObservabilityConfig) != 0 {
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

	enableMetrics, ok := c["enableMetrics"].(bool)
	if !ok {
		enableMetrics = false
	}

	enableRelay, ok := c["enableRelay"].(bool)
	if !ok {
		enableRelay = false
	}

	transformed := map[string]interface{}{
		"enable_metrics": enableMetrics,
		"enable_relay":   enableRelay,
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

func flattenNodePoolDefaults(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	if nodeConfigDefaults, ok := c["nodeConfigDefaults"]; ok && nodeConfigDefaults != nil {
		result["node_config_defaults"] = flattenNodeConfigDefaults(nodeConfigDefaults)
	}

	return []map[string]interface{}{result}
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
		transformed["linux_node_config"] = flattenNodePoolAutoConfigLinuxNodeConfig(linuxNodeConfig)
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
