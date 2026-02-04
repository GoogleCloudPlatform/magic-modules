package container

import (
	"fmt"
	"regexp"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tpgresource"

	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/transport"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/container/v1"
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
	// TODO: use transport.NewConfig()
	// config := transport.NewConfig()

	// This is a fake resource used to get fake d
	// d.Get will return empty map, instead of nil
	// fakeResource := &schema.Resource{
	// 	Schema: c.schema,
	// }
	// d := fakeResource.TestResourceData()

	var cluster *container.Cluster
	if err := utils.DecodeJSON(asset.Resource.Data, &cluster); err != nil {
		return nil, err
	}

	hclData := make(map[string]interface{})

	hclData["name"] = cluster.Name
	hclData["description"] = cluster.Description
	hclData["security_posture_config"] = flattenSecurityPostureConfig(cluster.SecurityPostureConfig)
	hclData["enterprise_config"] = flattenEnterpriseConfig(cluster.EnterpriseConfig)
	hclData["anonymous_authentication_config"] = flattenAnonymousAuthenticationConfig(cluster.AnonymousAuthenticationConfig) // Assuming key
	hclData["notification_config"] = flattenNotificationConfig(cluster.NotificationConfig)
	hclData["binary_authorization"] = flattenBinaryAuthorization(cluster.BinaryAuthorization)
	hclData["network_policy"] = flattenNetworkPolicy(cluster.NetworkPolicy)
	hclData["addons_config"] = flattenClusterAddonsConfig(cluster.AddonsConfig)
	// TODO: node_pool
	hclData["authenticator_groups_config"] = flattenAuthenticatorGroupsConfig(cluster.AuthenticatorGroupsConfig)
	hclData["control_plane_endpoints_config"] = flattenControlPlaneEndpointsConfig(cluster.ControlPlaneEndpointsConfig)
	hclData["private_cluster_config"] = flattenPrivateClusterConfig(cluster.ControlPlaneEndpointsConfig, cluster.PrivateClusterConfig, cluster.NetworkConfig)
	hclData["vertical_pod_autoscaling"] = flattenVerticalPodAutoscaling(cluster.VerticalPodAutoscaling)
	hclData["release_channel"] = flattenReleaseChannel(cluster.ReleaseChannel)
	hclData["gke_auto_upgrade_config"] = flattenGkeAutoUpgradeConfig(cluster.GkeAutoUpgradeConfig) // key?
	hclData["default_snat_status"] = flattenDefaultSnatStatus(cluster.NetworkConfig.DefaultSnatStatus)
	hclData["workload_identity_config"] = flattenWorkloadIdentityConfig(cluster.WorkloadIdentityConfig, nil, nil)
	hclData["identity_service_config"] = flattenIdentityServiceConfig(cluster.IdentityServiceConfig, nil, nil)
	if cluster.IpAllocationPolicy != nil {
		hclData["pod_cidr_overprovision_config"] = flattenPodCidrOverprovisionConfig(cluster.IpAllocationPolicy.PodCidrOverprovisionConfig)
	}

	ipPolicy, err := flattenIPAllocationPolicy(cluster, nil, nil)
	if err != nil {
		return nil, err
	}
	hclData["ip_allocation_policy"] = ipPolicy

	if cluster.IpAllocationPolicy == nil || !cluster.IpAllocationPolicy.UseIpAliases {
		hclData["networking_mode"] = "ROUTES"
	} else {
		hclData["networking_mode"] = "VPC_NATIVE"
	}

	hclData["maintenance_policy"] = flattenMaintenancePolicy(cluster.MaintenancePolicy)
	hclData["master_auth"] = flattenMasterAuth(cluster.MasterAuth)
	hclData["cluster_autoscaling"] = flattenClusterAutoscaling(cluster.Autoscaling)
	hclData["master_authorized_networks_config"] = flattenMasterAuthorizedNetworksConfig(cluster.MasterAuthorizedNetworksConfig)
	hclData["pod_autoscaling"] = flattenPodAutoscaling(cluster.PodAutoscaling)
	hclData["secret_manager_config"] = flattenSecretManagerConfig(cluster.SecretManagerConfig)
	hclData["resource_usage_export_config"] = flattenResourceUsageExportConfig(cluster.ResourceUsageExportConfig)
	hclData["service_external_ips_config"] = flattenServiceExternalIpsConfig(cluster.NetworkConfig.ServiceExternalIpsConfig)
	hclData["mesh_certificates"] = flattenMeshCertificates(cluster.MeshCertificates)
	hclData["cost_management_config"] = flattenManagementConfig(cluster.CostManagementConfig)
	hclData["database_encryption"] = flattenDatabaseEncryption(cluster.DatabaseEncryption)
	hclData["dns_config"] = flattenDnsConfig(cluster.NetworkConfig.DnsConfig)
	hclData["network_performance_config"] = flattenNetworkPerformanceConfig(cluster.NetworkConfig.NetworkPerformanceConfig)
	hclData["gateway_api_config"] = flattenGatewayApiConfig(cluster.NetworkConfig.GatewayApiConfig)
	hclData["fleet"] = flattenFleet(cluster.Fleet)
	hclData["user_managed_keys_config"] = flattenUserManagedKeysConfig(cluster.UserManagedKeysConfig)
	hclData["enable_k8s_beta_apis"] = flattenEnableK8sBetaApis(cluster.EnableK8sBetaApis)
	hclData["logging_config"] = flattenContainerClusterLoggingConfig(cluster.LoggingConfig)
	hclData["monitoring_config"] = flattenMonitoringConfig(cluster.MonitoringConfig)
	hclData["node_pool_auto_config"] = flattenNodePoolAutoConfig(cluster.NodePoolAutoConfig)
	hclData["rbac_binding_config"] = flattenRBACBindingConfig(cluster.RbacBindingConfig)

	ctyVal, err := utils.MapToCtyValWithSchema(hclData, c.schema)
	if err != nil {
		return nil, err
	}
	return &models.TerraformResourceBlock{
		Labels: []string{c.name, cluster.Name},
		Value:  ctyVal,
	}, nil
}

func flattenSecurityPostureConfig(spc *container.SecurityPostureConfig) []map[string]interface{} {
	if spc == nil {
		return nil
	}
	result := make(map[string]interface{})

	result["mode"] = spc.Mode
	result["vulnerability_mode"] = spc.VulnerabilityMode

	return []map[string]interface{}{result}
}

func flattenEnterpriseConfig(ec *container.EnterpriseConfig) []map[string]interface{} {
	if ec == nil {
		return nil
	}
	result := make(map[string]interface{})

	result["cluster_tier"] = ec.ClusterTier
	result["desired_tier"] = ec.DesiredTier

	return []map[string]interface{}{result}
}

func flattenAnonymousAuthenticationConfig(aac *container.AnonymousAuthenticationConfig) []map[string]interface{} {
	if aac == nil {
		return nil
	}
	result := make(map[string]interface{})
	result["mode"] = aac.Mode
	return []map[string]interface{}{result}
}

func flattenAdditionalPodRangesConfig(ipAllocationPolicy *container.IPAllocationPolicy) []map[string]interface{} {
	if ipAllocationPolicy == nil {
		return nil
	}
	result := make(map[string]interface{})

	if aprc := ipAllocationPolicy.AdditionalPodRangesConfig; aprc != nil {
		if len(aprc.PodRangeNames) > 0 {
			result["pod_range_names"] = aprc.PodRangeNames
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

func flattenNotificationConfig(c *container.NotificationConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	if c.Pubsub.Filter != nil {
		filter := []map[string]interface{}{}
		if len(c.Pubsub.Filter.EventType) > 0 {
			filter = append(filter, map[string]interface{}{
				"event_type": c.Pubsub.Filter.EventType,
			})
		}

		return []map[string]interface{}{
			{
				"pubsub": []map[string]interface{}{
					{
						"enabled": c.Pubsub.Enabled,
						"topic":   c.Pubsub.Topic,
						"filter":  filter,
					},
				},
			},
		}
	}

	return []map[string]interface{}{
		{
			"pubsub": []map[string]interface{}{
				{
					"enabled": c.Pubsub.Enabled,
					"topic":   c.Pubsub.Topic,
				},
			},
		},
	}
}

func flattenBinaryAuthorization(c *container.BinaryAuthorization) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"enabled":         c.Enabled,
			"evaluation_mode": c.EvaluationMode,
		})
	}
	return result
}

func flattenNetworkPolicy(c *container.NetworkPolicy) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"enabled":  c.Enabled,
			"provider": c.Provider,
		})
	} else {
		// Explicitly set the network policy to the default.
		result = append(result, map[string]interface{}{
			"enabled":  false,
			"provider": "PROVIDER_UNSPECIFIED",
		})
	}
	return result
}

func flattenClusterAddonsConfig(c *container.AddonsConfig) []map[string]interface{} {
	result := make(map[string]interface{})
	if c == nil {
		return nil
	}
	if c.HorizontalPodAutoscaling != nil {
		result["horizontal_pod_autoscaling"] = []map[string]interface{}{
			{
				"disabled": c.HorizontalPodAutoscaling.Disabled,
			},
		}
	}
	if c.HttpLoadBalancing != nil {
		result["http_load_balancing"] = []map[string]interface{}{
			{
				"disabled": c.HttpLoadBalancing.Disabled,
			},
		}
	}
	if c.NetworkPolicyConfig != nil {
		result["network_policy_config"] = []map[string]interface{}{
			{
				"disabled": c.NetworkPolicyConfig.Disabled,
			},
		}
	}

	if c.GcpFilestoreCsiDriverConfig != nil {
		result["gcp_filestore_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": c.GcpFilestoreCsiDriverConfig.Enabled,
			},
		}
	}

	if c.CloudRunConfig != nil {
		cloudRunConfig := map[string]interface{}{
			"disabled": c.CloudRunConfig.Disabled,
		}
		if c.CloudRunConfig.LoadBalancerType == "LOAD_BALANCER_TYPE_INTERNAL" {
			// Currently we only allow setting load_balancer_type to LOAD_BALANCER_TYPE_INTERNAL
			cloudRunConfig["load_balancer_type"] = "LOAD_BALANCER_TYPE_INTERNAL"
		}
		result["cloudrun_config"] = []map[string]interface{}{cloudRunConfig}
	}

	if c.DnsCacheConfig != nil {
		result["dns_cache_config"] = []map[string]interface{}{
			{
				"enabled": c.DnsCacheConfig.Enabled,
			},
		}
	}

	if c.GcePersistentDiskCsiDriverConfig != nil {
		result["gce_persistent_disk_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": c.GcePersistentDiskCsiDriverConfig.Enabled,
			},
		}
	}
	if c.GkeBackupAgentConfig != nil {
		result["gke_backup_agent_config"] = []map[string]interface{}{
			{
				"enabled": c.GkeBackupAgentConfig.Enabled,
			},
		}
	}
	if c.ConfigConnectorConfig != nil {
		result["config_connector_config"] = []map[string]interface{}{
			{
				"enabled": c.ConfigConnectorConfig.Enabled,
			},
		}
	}
	if c.GcsFuseCsiDriverConfig != nil {
		result["gcs_fuse_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": c.GcsFuseCsiDriverConfig.Enabled,
			},
		}
	}
	if c.StatefulHaConfig != nil {
		result["stateful_ha_config"] = []map[string]interface{}{
			{
				"enabled": c.StatefulHaConfig.Enabled,
			},
		}
	}
	if c.RayOperatorConfig != nil {
		rayConfig := c.RayOperatorConfig
		result["ray_operator_config"] = []map[string]interface{}{
			{
				"enabled": rayConfig.Enabled,
			},
		}
		if rayConfig.RayClusterLoggingConfig != nil {
			result["ray_operator_config"].([]map[string]any)[0]["ray_cluster_logging_config"] = []map[string]interface{}{{
				"enabled": rayConfig.RayClusterLoggingConfig.Enabled,
			}}
		}
		if rayConfig.RayClusterMonitoringConfig != nil {
			result["ray_operator_config"].([]map[string]any)[0]["ray_cluster_monitoring_config"] = []map[string]interface{}{{
				"enabled": rayConfig.RayClusterMonitoringConfig.Enabled,
			}}
		}
	}
	if c.ParallelstoreCsiDriverConfig != nil {
		result["parallelstore_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled": c.ParallelstoreCsiDriverConfig.Enabled,
			},
		}
	}
	if c.LustreCsiDriverConfig != nil {
		lustreConfig := c.LustreCsiDriverConfig
		result["lustre_csi_driver_config"] = []map[string]interface{}{
			{
				"enabled":                   lustreConfig.Enabled,
				"enable_legacy_lustre_port": lustreConfig.EnableLegacyLustrePort,
			},
		}
	}

	return []map[string]interface{}{result}
}

func flattenAuthenticatorGroupsConfig(c *container.AuthenticatorGroupsConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"security_group": c.SecurityGroup,
		},
	}
}

func flattenControlPlaneEndpointsConfig(c *container.ControlPlaneEndpointsConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"dns_endpoint_config": flattenDnsEndpointConfig(c.DnsEndpointConfig),
			"ip_endpoints_config": flattenIpEndpointsConfig(c.IpEndpointsConfig),
		},
	}
}

func flattenDnsEndpointConfig(dns *container.DNSEndpointConfig) []map[string]interface{} {
	if dns == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"endpoint":                  dns.Endpoint,
			"allow_external_traffic":    dns.AllowExternalTraffic,
			"enable_k8s_tokens_via_dns": dns.EnableK8sTokensViaDns,
			"enable_k8s_certs_via_dns":  dns.EnableK8sCertsViaDns,
		},
	}
}

func flattenIpEndpointsConfig(ip *container.IPEndpointsConfig) []map[string]interface{} {
	if ip == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"enabled": ip.Enabled,
		},
	}
}

// Most of PrivateClusterConfig has moved to ControlPlaneEndpointsConfig.
func flattenPrivateClusterConfig(cpec *container.ControlPlaneEndpointsConfig, pcc *container.PrivateClusterConfig, nc *container.NetworkConfig) []map[string]interface{} {
	if cpec == nil && pcc == nil && nc == nil {
		return nil
	}

	r := map[string]interface{}{}
	if cpec != nil {
		// Note the change in semantics from private to public endpoint.
		r["enable_private_endpoint"] = !cpec.IpEndpointsConfig.EnablePublicEndpoint
		r["private_endpoint"] = cpec.IpEndpointsConfig.PrivateEndpoint
		r["private_endpoint_subnetwork"] = cpec.IpEndpointsConfig.PrivateEndpointSubnetwork
		r["public_endpoint"] = cpec.IpEndpointsConfig.PublicEndpoint
		r["master_global_access_config"] = []map[string]interface{}{
			{
				"enabled": cpec.IpEndpointsConfig.GlobalAccess,
			},
		}
	}
	// This is the only field that is canonically still in the PrivateClusterConfig message.
	if pcc != nil {
		r["peering_name"] = pcc.PeeringName
		r["master_ipv4_cidr_block"] = pcc.MasterIpv4CidrBlock
	}
	if nc != nil {
		r["enable_private_nodes"] = nc.DefaultEnablePrivateNodes
	}

	return []map[string]interface{}{r}
}

func flattenVerticalPodAutoscaling(c *container.VerticalPodAutoscaling) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"enabled": c.Enabled,
		},
	}
}

func flattenReleaseChannel(c *container.ReleaseChannel) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil && c.Channel != "" {
		result = append(result, map[string]interface{}{
			"channel": c.Channel,
		})
	} else {
		// Explicitly set the release channel to the UNSPECIFIED.
		result = append(result, map[string]interface{}{
			"channel": "UNSPECIFIED",
		})
	}
	return result
}

func flattenGkeAutoUpgradeConfig(c *container.GkeAutoUpgradeConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	result := []map[string]interface{}{}
	if c.PatchMode != "" {
		result = append(result, map[string]interface{}{
			"patch_mode": c.PatchMode,
		})
	}

	return result
}

func flattenDefaultSnatStatus(c *container.DefaultSnatStatus) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"disabled": c.Disabled,
		})
	}
	return result
}

func flattenWorkloadIdentityConfig(c *container.WorkloadIdentityConfig, d *schema.ResourceData, config *transport_tpg.Config) []map[string]interface{} {
	if c == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"workload_pool": c.WorkloadPool,
		},
	}
}

func flattenIdentityServiceConfig(c *container.IdentityServiceConfig, d *schema.ResourceData, config *transport_tpg.Config) []map[string]interface{} {
	if c == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"enabled": c.Enabled,
		},
	}
}

func flattenPodCidrOverprovisionConfig(c *container.PodCIDROverprovisionConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"disabled": c.Disable,
		},
	}
}

func flattenAdditionalIpRangesConfigs(c []*container.AdditionalIPRangesConfig) []map[string]interface{} {
	if len(c) == 0 {
		return nil
	}

	var outRanges []map[string]interface{}
	for _, rangeConfig := range c {
		outRangeConfig := map[string]interface{}{
			"subnetwork":           rangeConfig.Subnetwork,
			"pod_ipv4_range_names": rangeConfig.PodIpv4RangeNames,
		}
		outRanges = append(outRanges, outRangeConfig)
	}

	return outRanges
}

func flattenNetworkTierConfig(ntc *container.NetworkTierConfig) []map[string]interface{} {
	if ntc == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"network_tier": ntc.NetworkTier,
		},
	}
}

func flattenIPAllocationPolicy(c *container.Cluster, d *schema.ResourceData, config *transport_tpg.Config) ([]map[string]interface{}, error) {
	// If IP aliasing isn't enabled, none of the values in this block can be set.
	if c == nil || c.IpAllocationPolicy == nil || !c.IpAllocationPolicy.UseIpAliases {
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

	p := c.IpAllocationPolicy

	// handle older clusters that return JSON null
	// corresponding to "STACK_TYPE_UNSPECIFIED" due to GKE declining to backfill
	// equivalent to default_if_empty
	if p.StackType == "" {
		p.StackType = "IPV4"
	}

	return []map[string]interface{}{
		{
			"cluster_ipv4_cidr_block":       p.ClusterIpv4CidrBlock,
			"services_ipv4_cidr_block":      p.ServicesIpv4CidrBlock,
			"cluster_secondary_range_name":  p.ClusterSecondaryRangeName,
			"services_secondary_range_name": p.ServicesSecondaryRangeName,
			"stack_type":                    p.StackType,
			"pod_cidr_overprovision_config": flattenPodCidrOverprovisionConfig(p.PodCidrOverprovisionConfig),
			"additional_pod_ranges_config":  flattenAdditionalPodRangesConfig(c.IpAllocationPolicy),
			"additional_ip_ranges_config":   flattenAdditionalIpRangesConfigs(p.AdditionalIpRangesConfigs),
			"auto_ipam_config":              flattenAutoIpamConfig(p.AutoIpamConfig),
			"network_tier_config":           flattenNetworkTierConfig(p.NetworkTierConfig),
		},
	}, nil
}

func flattenAutoIpamConfig(aic *container.AutoIpamConfig) []map[string]interface{} {
	if aic == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"enabled": aic.Enabled,
		},
	}
}

func flattenMaintenancePolicy(mp *container.MaintenancePolicy) []map[string]interface{} {
	if mp == nil || mp.Window == nil {
		return nil
	}

	exclusions := []map[string]interface{}{}
	if mp.Window.MaintenanceExclusions != nil {
		for wName, window := range mp.Window.MaintenanceExclusions {
			exclusion := map[string]interface{}{
				"start_time":     window.StartTime,
				"exclusion_name": wName,
			}
			if window.MaintenanceExclusionOptions != nil {
				// When the scope is set to NO_UPGRADES which is the default value,
				// the maintenance exclusion returned by GCP will be empty.
				// This seems like a bug. To workaround this, assign NO_UPGRADES to the scope explicitly
				scope := "NO_UPGRADES"
				if window.MaintenanceExclusionOptions.Scope != "" {
					scope = window.MaintenanceExclusionOptions.Scope
				}
				if window.MaintenanceExclusionOptions.EndTimeBehavior != "" {
					exclusion["exclusion_options"] = []map[string]interface{}{
						{
							"scope":             scope,
							"end_time_behavior": window.MaintenanceExclusionOptions.EndTimeBehavior,
						},
					}
				} else {
					exclusion["exclusion_options"] = []map[string]interface{}{
						{
							"scope": scope,
						},
					}
					if window.EndTime != "" {
						exclusion["end_time"] = window.EndTime
					}
				}
			} else {
				if window.EndTime != "" {
					exclusion["end_time"] = window.EndTime
				}
			}
			exclusions = append(exclusions, exclusion)
		}
	}

	if mp.Window.DailyMaintenanceWindow != nil {
		return []map[string]interface{}{
			{
				"daily_maintenance_window": []map[string]interface{}{
					{
						"start_time": mp.Window.DailyMaintenanceWindow.StartTime,
						"duration":   mp.Window.DailyMaintenanceWindow.Duration,
					},
				},
				"maintenance_exclusion": exclusions,
			},
		}
	}
	if mp.Window.RecurringWindow != nil {
		return []map[string]interface{}{
			{
				"recurring_window": []map[string]interface{}{
					{
						"start_time": mp.Window.RecurringWindow.Window.StartTime,
						"end_time":   mp.Window.RecurringWindow.Window.EndTime,
						"recurrence": mp.Window.RecurringWindow.Recurrence,
					},
				},
				"maintenance_exclusion": exclusions,
			},
		}
	}
	return nil
}

func flattenMasterAuth(ma *container.MasterAuth) []map[string]interface{} {
	if ma == nil {
		return nil
	}
	masterAuth := []map[string]interface{}{
		{
			"client_certificate":     ma.ClientCertificate,
			"client_key":             ma.ClientKey,
			"cluster_ca_certificate": ma.ClusterCaCertificate,
		},
	}

	// No version of the GKE API returns the client_certificate_config value.
	// Instead, we need to infer whether or not it was set based on the
	// client cert being returned from the API or not.
	// Previous versions of the provider didn't record anything in state when
	// a client cert was enabled, only setting the block when it was false.
	masterAuth[0]["client_certificate_config"] = []map[string]interface{}{
		{
			"issue_client_certificate": len(ma.ClientCertificate) != 0,
		},
	}

	return masterAuth
}

func flattenClusterAutoscaling(a *container.ClusterAutoscaling) []map[string]interface{} {
	r := make(map[string]interface{})
	if a == nil {
		r["enabled"] = false
		return []map[string]interface{}{r}
	}

	if a.EnableNodeAutoprovisioning {
		resourceLimits := make([]interface{}, 0, len(a.ResourceLimits))
		for _, rl := range a.ResourceLimits {
			resourceLimits = append(resourceLimits, map[string]interface{}{
				"resource_type": rl.ResourceType,
				"minimum":       rl.Minimum,
				"maximum":       rl.Maximum,
			})
		}
		r["resource_limits"] = resourceLimits
		r["enabled"] = true
		r["auto_provisioning_defaults"] = flattenAutoProvisioningDefaults(a.AutoprovisioningNodePoolDefaults)
		r["auto_provisioning_locations"] = a.AutoprovisioningLocations
	} else {
		r["enabled"] = false
	}
	r["autoscaling_profile"] = a.AutoscalingProfile
	if a.DefaultComputeClassConfig != nil {
		r["default_compute_class_enabled"] = a.DefaultComputeClassConfig.Enabled
	}

	return []map[string]interface{}{r}
}

func flattenAutoProvisioningDefaults(a *container.AutoprovisioningNodePoolDefaults) []map[string]interface{} {
	r := make(map[string]interface{})
	r["oauth_scopes"] = a.OauthScopes
	r["service_account"] = a.ServiceAccount
	r["disk_size"] = a.DiskSizeGb
	r["disk_type"] = a.DiskType
	r["image_type"] = a.ImageType
	r["min_cpu_platform"] = a.MinCpuPlatform
	r["boot_disk_kms_key"] = a.BootDiskKmsKey
	r["shielded_instance_config"] = flattenShieldedInstanceConfig(a.ShieldedInstanceConfig)
	r["management"] = flattenManagement(a.Management)
	r["upgrade_settings"] = flattenUpgradeSettings(a.UpgradeSettings)

	return []map[string]interface{}{r}
}

func flattenUpgradeSettings(a *container.UpgradeSettings) []map[string]interface{} {
	if a == nil {
		return nil
	}
	r := make(map[string]interface{})
	r["max_surge"] = a.MaxSurge
	r["max_unavailable"] = a.MaxUnavailable
	r["strategy"] = a.Strategy
	r["blue_green_settings"] = flattenBlueGreenSettings(a.BlueGreenSettings)

	return []map[string]interface{}{r}
}

func flattenBlueGreenSettings(a *container.BlueGreenSettings) []map[string]interface{} {
	if a == nil {
		return nil
	}

	r := make(map[string]interface{})
	r["node_pool_soak_duration"] = a.NodePoolSoakDuration
	r["standard_rollout_policy"] = flattenStandardRolloutPolicy(a.StandardRolloutPolicy)

	return []map[string]interface{}{r}
}

func flattenStandardRolloutPolicy(a *container.StandardRolloutPolicy) []map[string]interface{} {
	if a == nil {
		return nil
	}

	r := make(map[string]interface{})
	r["batch_percentage"] = a.BatchPercentage
	r["batch_node_count"] = a.BatchNodeCount
	r["batch_soak_duration"] = a.BatchSoakDuration

	return []map[string]interface{}{r}
}

func flattenManagement(a *container.NodeManagement) []map[string]interface{} {
	if a == nil {
		return nil
	}
	r := make(map[string]interface{})
	r["auto_upgrade"] = a.AutoUpgrade
	r["auto_repair"] = a.AutoRepair
	r["upgrade_options"] = flattenUpgradeOptions(a.UpgradeOptions)

	return []map[string]interface{}{r}
}

func flattenUpgradeOptions(a *container.AutoUpgradeOptions) []map[string]interface{} {
	if a == nil {
		return nil
	}

	r := make(map[string]interface{})
	r["auto_upgrade_start_time"] = a.AutoUpgradeStartTime
	r["description"] = a.Description

	return []map[string]interface{}{r}
}

func flattenMasterAuthorizedNetworksConfig(c *container.MasterAuthorizedNetworksConfig) []map[string]interface{} {
	if c == nil || !c.Enabled {
		return nil
	}
	result := make(map[string]interface{})
	cidrBlocks := make([]interface{}, 0, len(c.CidrBlocks))
	for _, v := range c.CidrBlocks {
		cidrBlocks = append(cidrBlocks, map[string]interface{}{
			"cidr_block":   v.CidrBlock,
			"display_name": v.DisplayName,
		})
	}
	result["cidr_blocks"] = schema.NewSet(schema.HashResource(cidrBlockConfig), cidrBlocks)
	result["gcp_public_cidrs_access_enabled"] = c.GcpPublicCidrsAccessEnabled
	result["private_endpoint_enforcement_enabled"] = c.PrivateEndpointEnforcementEnabled
	return []map[string]interface{}{result}
}

func flattenPodAutoscaling(c *container.PodAutoscaling) []map[string]interface{} {
	config := make([]map[string]interface{}, 0, 1)

	if c == nil {
		return config
	}

	config = append(config, map[string]interface{}{
		"hpa_profile": c.HpaProfile,
	})
	return config
}

func flattenSecretManagerConfig(c *container.SecretManagerConfig) []map[string]interface{} {
	if c == nil {
		return []map[string]interface{}{
			{
				"enabled": false,
			},
		}
	}

	result := make(map[string]interface{})

	result["enabled"] = c.Enabled

	rotationList := []map[string]interface{}{}
	if c.RotationConfig != nil {
		rotationConfigMap := map[string]interface{}{
			"enabled": c.RotationConfig.Enabled,
		}
		if c.RotationConfig.RotationInterval != "" {
			rotationConfigMap["rotation_interval"] = c.RotationConfig.RotationInterval
		}
		rotationList = append(rotationList, rotationConfigMap)
	}
	result["rotation_config"] = rotationList
	return []map[string]interface{}{result}
}

func flattenResourceUsageExportConfig(c *container.ResourceUsageExportConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	enableResourceConsumptionMetering := false
	if c.ConsumptionMeteringConfig != nil && c.ConsumptionMeteringConfig.Enabled == true {
		enableResourceConsumptionMetering = true
	}

	return []map[string]interface{}{
		{
			"enable_network_egress_metering":       c.EnableNetworkEgressMetering,
			"enable_resource_consumption_metering": enableResourceConsumptionMetering,
			"bigquery_destination": []map[string]interface{}{
				{"dataset_id": c.BigqueryDestination.DatasetId},
			},
		},
	}
}

func flattenServiceExternalIpsConfig(c *container.ServiceExternalIPsConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"enabled": c.Enabled,
		},
	}
}

func flattenMeshCertificates(c *container.MeshCertificates) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"enable_certificates": c.EnableCertificates,
		},
	}
}

func flattenManagementConfig(c *container.CostManagementConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"enabled": c.Enabled,
		},
	}
}

func flattenDatabaseEncryption(c *container.DatabaseEncryption) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"state":    c.State,
			"key_name": c.KeyName,
		},
	}
}

func flattenDnsConfig(c *container.DNSConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"additive_vpc_scope_dns_domain": c.AdditiveVpcScopeDnsDomain,
			"cluster_dns":                   c.ClusterDns,
			"cluster_dns_scope":             c.ClusterDnsScope,
			"cluster_dns_domain":            c.ClusterDnsDomain,
		},
	}
}

func flattenNetworkPerformanceConfig(c *container.ClusterNetworkPerformanceConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"total_egress_bandwidth_tier": c.TotalEgressBandwidthTier,
		},
	}
}

func flattenGatewayApiConfig(c *container.GatewayAPIConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"channel": c.Channel,
		},
	}
}

func flattenFleet(c *container.Fleet) []map[string]interface{} {
	if c == nil {
		return nil
	}

	// Parse membership_id and membership_location from full membership name.
	var membership_id, membership_location string
	membershipRE := regexp.MustCompile(`^(//[a-zA-Z0-9\.\-]+)?/?projects/([^/]+)/locations/([a-zA-Z0-9\-]+)/memberships/([^/]+)$`)
	if match := membershipRE.FindStringSubmatch(c.Membership); match != nil {
		membership_id = match[4]
		membership_location = match[3]
	}

	return []map[string]interface{}{
		{
			"project":             c.Project,
			"membership":          c.Membership,
			"membership_id":       membership_id,
			"membership_location": membership_location,
			"pre_registered":      c.PreRegistered,
			"membership_type":     c.MembershipType,
		},
	}
}

func flattenUserManagedKeysConfig(c *container.UserManagedKeysConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	f := map[string]interface{}{
		"cluster_ca":                        c.ClusterCa,
		"etcd_api_ca":                       c.EtcdApiCa,
		"etcd_peer_ca":                      c.EtcdPeerCa,
		"aggregation_ca":                    c.AggregationCa,
		"control_plane_disk_encryption_key": c.ControlPlaneDiskEncryptionKey,
		"gkeops_etcd_backup_encryption_key": c.GkeopsEtcdBackupEncryptionKey,
	}
	allEmpty := true
	for _, v := range f {
		if v.(string) != "" {
			allEmpty = false
		}
	}
	if len(c.ServiceAccountSigningKeys) != 0 {
		f["service_account_signing_keys"] = schema.NewSet(schema.HashString, tpgresource.ConvertStringArrToInterface(c.ServiceAccountSigningKeys))
		allEmpty = false
	}
	if len(c.ServiceAccountVerificationKeys) != 0 {
		f["service_account_verification_keys"] = schema.NewSet(schema.HashString, tpgresource.ConvertStringArrToInterface(c.ServiceAccountVerificationKeys))
		allEmpty = false
	}
	if allEmpty {
		return nil
	}
	return []map[string]interface{}{f}
}

func flattenEnableK8sBetaApis(c *container.K8sBetaAPIConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"enabled_apis": c.EnabledApis,
		},
	}
}

func flattenContainerClusterLoggingConfig(c *container.LoggingConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"enable_components": c.ComponentConfig.EnableComponents,
		},
	}
}

func flattenMonitoringConfig(c *container.MonitoringConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	result := make(map[string]interface{})
	if c.ComponentConfig != nil {
		result["enable_components"] = c.ComponentConfig.EnableComponents
	}
	if c.ManagedPrometheusConfig != nil {
		result["managed_prometheus"] = flattenManagedPrometheusConfig(c.ManagedPrometheusConfig)
	}
	if c.AdvancedDatapathObservabilityConfig != nil {
		result["advanced_datapath_observability_config"] = flattenAdvancedDatapathObservabilityConfig(c.AdvancedDatapathObservabilityConfig)
	}

	return []map[string]interface{}{result}
}

func flattenAdvancedDatapathObservabilityConfig(c *container.AdvancedDatapathObservabilityConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"enable_metrics": c.EnableMetrics,
			"enable_relay":   c.EnableRelay,
		},
	}
}

func flattenManagedPrometheusConfig(c *container.ManagedPrometheusConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	result := make(map[string]interface{})
	result["enabled"] = c.Enabled

	autoMonitoringList := []map[string]interface{}{}
	if c.AutoMonitoringConfig != nil && c.AutoMonitoringConfig.Scope != "" {
		autoMonitoringMap := map[string]interface{}{
			"scope": c.AutoMonitoringConfig.Scope,
		}
		autoMonitoringList = append(autoMonitoringList, autoMonitoringMap)
	}

	result["auto_monitoring_config"] = autoMonitoringList

	return []map[string]interface{}{result}
}

func flattenNodePoolAutoConfig(c *container.NodePoolAutoConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	result := make(map[string]interface{})
	if c.NodeKubeletConfig != nil {
		result["node_kubelet_config"] = flattenNodePoolAutoConfigNodeKubeletConfig(c.NodeKubeletConfig)
	}
	if c.NetworkTags != nil {
		result["network_tags"] = flattenNodePoolAutoConfigNetworkTags(c.NetworkTags)
	}
	if c.ResourceManagerTags != nil {
		result["resource_manager_tags"] = flattenResourceManagerTags(c.ResourceManagerTags)
	}
	if c.LinuxNodeConfig != nil {
		result["linux_node_config"] = []map[string]interface{}{
			{"cgroup_mode": c.LinuxNodeConfig.CgroupMode},
		}
	}

	return []map[string]interface{}{result}
}

func flattenNodePoolAutoConfigNetworkTags(c *container.NetworkTags) []map[string]interface{} {
	if c == nil {
		return nil
	}

	result := make(map[string]interface{})
	if c.Tags != nil {
		result["tags"] = c.Tags
	}
	return []map[string]interface{}{result}
}

func flattenRBACBindingConfig(c *container.RBACBindingConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"enable_insecure_binding_system_authenticated":   c.EnableInsecureBindingSystemAuthenticated,
			"enable_insecure_binding_system_unauthenticated": c.EnableInsecureBindingSystemUnauthenticated,
		},
	}
}
