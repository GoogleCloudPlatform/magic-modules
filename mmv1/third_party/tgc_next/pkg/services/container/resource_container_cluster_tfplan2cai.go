package container

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/converters/cai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/transport"

	"google.golang.org/api/container/v1"
)

func ContainerClusterTfplan2CaiConverter() cai.Tfplan2caiConverter {
	return cai.Tfplan2caiConverter{
		Convert: GetContainerCluster,
	}
}

func GetContainerCluster(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]caiasset.Asset, error) {
	name, err := cai.AssetName(d, config, "//container.googleapis.com/projects/{{project}}/locations/{{location}}/clusters/{{name}}")
	if err != nil {
		return []caiasset.Asset{}, err
	}
	if data, err := GetContainerClusterData(d, config); err == nil {
		location, _ := tpgresource.GetLocation(d, config)
		return []caiasset.Asset{
			{
				Name: name,
				Type: ContainerClusterAssetType,
				Resource: &caiasset.AssetResource{
					Version:              "v1",
					DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/container/v1/rest",
					DiscoveryName:        "Cluster",
					Data:                 data,
					Location:             location,
				},
			},
		}, nil
	} else {
		return []caiasset.Asset{}, err
	}
}

func GetContainerClusterData(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	cluster, err := expandContainerCluster(project, d, config)
	if err != nil {
		return nil, err
	}

	return cai.JsonMap(cluster)
}

func expandContainerCluster(project string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*container.Cluster, error) {
	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return nil, err
	}

	clusterName := d.Get("name").(string)

	// TODO: check if needUpdateAfterCreate is needed
	ipAllocationBlock, _, err := expandIPAllocationPolicy(d.Get("ip_allocation_policy"), d, d.Get("networking_mode").(string), d.Get("enable_autopilot").(bool), config)
	if err != nil {
		return nil, err
	}

	var workloadPolicyConfig *container.WorkloadPolicyConfig
	if allowed := d.Get("allow_net_admin").(bool); allowed {
		workloadPolicyConfig = &container.WorkloadPolicyConfig{
			AllowNetAdmin: allowed,
		}
	}

	cluster := &container.Cluster{
		Name:                        clusterName,
		InitialNodeCount:            int64(d.Get("initial_node_count").(int)),
		MaintenancePolicy:           expandMaintenancePolicy(d, config),
		ControlPlaneEndpointsConfig: expandControlPlaneEndpointsConfig(d),
		InitialClusterVersion:       d.Get("min_master_version").(string),
		ClusterIpv4Cidr:             d.Get("cluster_ipv4_cidr").(string),
		Description:                 d.Get("description").(string),
		LegacyAbac: &container.LegacyAbac{
			Enabled:         d.Get("enable_legacy_abac").(bool),
			ForceSendFields: []string{"Enabled"},
		},
		LoggingService:        d.Get("logging_service").(string),
		MonitoringService:     d.Get("monitoring_service").(string),
		NetworkPolicy:         expandNetworkPolicy(d.Get("network_policy")),
		AddonsConfig:          expandClusterAddonsConfig(d.Get("addons_config")),
		EnableKubernetesAlpha: d.Get("enable_kubernetes_alpha").(bool),
		IpAllocationPolicy:    ipAllocationBlock,
		PodAutoscaling:        expandPodAutoscaling(d.Get("pod_autoscaling")),
		SecretManagerConfig:   expandSecretManagerConfig(d.Get("secret_manager_config")),

		Autoscaling:         expandClusterAutoscaling(d.Get("cluster_autoscaling"), d),
		BinaryAuthorization: expandBinaryAuthorization(d.Get("binary_authorization")),
		Autopilot: &container.Autopilot{
			Enabled:              d.Get("enable_autopilot").(bool),
			WorkloadPolicyConfig: workloadPolicyConfig,
			ForceSendFields:      []string{"Enabled"},
		},
		ReleaseChannel:       expandReleaseChannel(d.Get("release_channel")),
		GkeAutoUpgradeConfig: expandGkeAutoUpgradeConfig(d.Get("gke_auto_upgrade_config")),

		EnableTpu: d.Get("enable_tpu").(bool),
		NetworkConfig: &container.NetworkConfig{
			EnableIntraNodeVisibility:            d.Get("enable_intranode_visibility").(bool),
			DefaultSnatStatus:                    expandDefaultSnatStatus(d.Get("default_snat_status")),
			DatapathProvider:                     d.Get("datapath_provider").(string),
			EnableCiliumClusterwideNetworkPolicy: d.Get("enable_cilium_clusterwide_network_policy").(bool),
			PrivateIpv6GoogleAccess:              d.Get("private_ipv6_google_access").(string),
			InTransitEncryptionConfig:            d.Get("in_transit_encryption_config").(string),
			EnableL4ilbSubsetting:                d.Get("enable_l4_ilb_subsetting").(bool),
			DisableL4LbFirewallReconciliation:    d.Get("disable_l4_lb_firewall_reconciliation").(bool),
			DnsConfig:                            expandDnsConfig(d.Get("dns_config")),
			GatewayApiConfig:                     expandGatewayApiConfig(d.Get("gateway_api_config")),
			EnableMultiNetworking:                d.Get("enable_multi_networking").(bool),
			DefaultEnablePrivateNodes:            expandDefaultEnablePrivateNodes(d),
			EnableFqdnNetworkPolicy:              d.Get("enable_fqdn_network_policy").(bool),
			NetworkPerformanceConfig:             expandNetworkPerformanceConfig(d.Get("network_performance_config")),
		},
		MasterAuth:           expandMasterAuth(d.Get("master_auth")),
		NotificationConfig:   expandNotificationConfig(d.Get("notification_config")),
		ConfidentialNodes:    expandConfidentialNodes(d.Get("confidential_nodes")),
		ResourceLabels:       tpgresource.ExpandStringMap(d, "effective_labels"),
		NodePoolAutoConfig:   expandNodePoolAutoConfig(d.Get("node_pool_auto_config")),
		CostManagementConfig: expandCostManagementConfig(d.Get("cost_management_config")),
		EnableK8sBetaApis:    expandEnableK8sBetaApis(d.Get("enable_k8s_beta_apis"), nil),
	}

	v := d.Get("enable_shielded_nodes")
	cluster.ShieldedNodes = &container.ShieldedNodes{
		Enabled:         v.(bool),
		ForceSendFields: []string{"Enabled"},
	}

	if v, ok := d.GetOk("default_max_pods_per_node"); ok {
		cluster.DefaultMaxPodsConstraint = expandDefaultMaxPodsConstraint(v)
	}

	// Only allow setting node_version on create if it's set to the equivalent master version,
	// since `InitialClusterVersion` only accepts valid master-style versions.
	if v, ok := d.GetOk("node_version"); ok {
		// ignore -gke.X suffix for now. if it becomes a problem later, we can fix it.
		mv := strings.Split(cluster.InitialClusterVersion, "-")[0]
		nv := strings.Split(v.(string), "-")[0]
		if mv != nv {
			return nil, fmt.Errorf("node_version and min_master_version must be set to equivalent values on create")
		}
	}

	if v, ok := d.GetOk("node_locations"); ok {
		locationsSet := v.(*schema.Set)
		if locationsSet.Contains(location) {
			return nil, fmt.Errorf("when using a multi-zonal cluster, node_locations should not contain the original 'zone'")
		}

		// GKE requires a full list of node locations
		// but when using a multi-zonal cluster our schema only asks for the
		// additional zones, so append the cluster location if it's a zone
		if tpgresource.IsZone(location) {
			locationsSet.Add(location)
		}
		cluster.Locations = tpgresource.ConvertStringSet(locationsSet)
	}

	if v, ok := d.GetOk("network"); ok {
		network, err := tpgresource.ParseNetworkFieldValue(v.(string), d, config)
		if err != nil {
			return nil, err
		}
		cluster.Network = network.RelativeLink()
	}

	if v, ok := d.GetOk("subnetwork"); ok {
		subnetwork, err := tpgresource.ParseRegionalFieldValue("subnetworks", v.(string), "project", "location", "location", d, config, true) // variant of ParseSubnetworkFieldValue
		if err != nil {
			return nil, err
		}
		cluster.Subnetwork = subnetwork.RelativeLink()
	}

	nodePoolsCount := d.Get("node_pool.#").(int)
	if nodePoolsCount > 0 {
		// TODO: implement expandNodePool

		// nodePools := make([]*container.NodePool, 0, nodePoolsCount)
		// for i := 0; i < nodePoolsCount; i++ {
		// 	prefix := fmt.Sprintf("node_pool.%d.", i)
		// 	nodePool, err := expandNodePool(d, prefix)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	nodePools = append(nodePools, nodePool)
		// }
		// cluster.NodePools = nodePools
	} else {
		// Node Configs have default values that are set in the expand function,
		// but can only be set if node pools are unspecified.
		cluster.NodeConfig = expandNodeConfig(d, "", []interface{}{})
	}

	if v, ok := d.GetOk("node_pool_defaults"); ok {
		cluster.NodePoolDefaults = expandNodePoolDefaults(v)
	}

	if v, ok := d.GetOk("node_config"); ok {
		cluster.NodeConfig = expandNodeConfig(d, "", v)
	}

	if v, ok := d.GetOk("authenticator_groups_config"); ok {
		cluster.AuthenticatorGroupsConfig = expandAuthenticatorGroupsConfig(v)
	}

	if v, ok := d.GetOk("private_cluster_config.0.master_ipv4_cidr_block"); ok {
		cluster.PrivateClusterConfig = expandPrivateClusterConfigMasterIpv4CidrBlock(v, cluster)
	}

	if v, ok := d.GetOk("vertical_pod_autoscaling"); ok {
		cluster.VerticalPodAutoscaling = expandVerticalPodAutoscaling(v)
	}

	if v, ok := d.GetOk("service_external_ips_config"); ok {
		cluster.NetworkConfig.ServiceExternalIpsConfig = expandServiceExternalIpsConfig(v)
	}

	if v, ok := d.GetOk("mesh_certificates"); ok {
		cluster.MeshCertificates = expandMeshCertificates(v)
	}

	if v, ok := d.GetOk("database_encryption"); ok {
		cluster.DatabaseEncryption = expandDatabaseEncryption(v)
	}

	if v, ok := d.GetOk("workload_identity_config"); ok {
		cluster.WorkloadIdentityConfig = expandWorkloadIdentityConfig(v)
	}

	if v, ok := d.GetOk("identity_service_config"); ok {
		cluster.IdentityServiceConfig = expandIdentityServiceConfig(v)
	}

	if v, ok := d.GetOk("resource_usage_export_config"); ok {
		cluster.ResourceUsageExportConfig = expandResourceUsageExportConfig(v)
	}

	if v, ok := d.GetOk("logging_config"); ok {
		cluster.LoggingConfig = expandContainerClusterLoggingConfig(v)
	}

	if v, ok := d.GetOk("monitoring_config"); ok {
		cluster.MonitoringConfig = expandMonitoringConfig(v)
	}

	if v, ok := d.GetOk("fleet"); ok {
		cluster.Fleet = expandFleet(v)
	}

	if v, ok := d.GetOk("user_managed_keys_config"); ok {
		cluster.UserManagedKeysConfig = expandUserManagedKeysConfig(v)
	}

	if err := validateNodePoolAutoConfig(cluster); err != nil {
		return nil, err
	}

	if v, ok := d.GetOk("security_posture_config"); ok {
		cluster.SecurityPostureConfig = expandSecurityPostureConfig(v)
	}

	if v, ok := d.GetOk("enterprise_config"); ok {
		cluster.EnterpriseConfig = expandEnterpriseConfig(v)
	}

	if v, ok := d.GetOk("anonymous_authentication_config"); ok {
		cluster.AnonymousAuthenticationConfig = expandAnonymousAuthenticationConfig(v)
	}

	if v, ok := d.GetOk("rbac_binding_config"); ok {
		cluster.RbacBindingConfig = expandRBACBindingConfig(v)
	}
	return cluster, nil
}

func expandEnterpriseConfig(configured interface{}) *container.EnterpriseConfig {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}

	ec := &container.EnterpriseConfig{}
	enterpriseConfig := l[0].(map[string]interface{})
	if v, ok := enterpriseConfig["cluster_tier"]; ok {
		ec.ClusterTier = v.(string)
	}

	if v, ok := enterpriseConfig["desired_tier"]; ok {
		ec.DesiredTier = v.(string)
	}
	return ec
}

func expandClusterAddonsConfig(configured interface{}) *container.AddonsConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
	ac := &container.AddonsConfig{}

	if v, ok := config["http_load_balancing"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.HttpLoadBalancing = &container.HttpLoadBalancing{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
	}

	if v, ok := config["horizontal_pod_autoscaling"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.HorizontalPodAutoscaling = &container.HorizontalPodAutoscaling{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
	}

	if v, ok := config["network_policy_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.NetworkPolicyConfig = &container.NetworkPolicyConfig{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
	}

	if v, ok := config["gcp_filestore_csi_driver_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.GcpFilestoreCsiDriverConfig = &container.GcpFilestoreCsiDriverConfig{
			Enabled:         addon["enabled"].(bool),
			ForceSendFields: []string{"Enabled"},
		}
	}

	if v, ok := config["cloudrun_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.CloudRunConfig = &container.CloudRunConfig{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
		if addon["load_balancer_type"] != "" {
			ac.CloudRunConfig.LoadBalancerType = addon["load_balancer_type"].(string)
		}
	}

	if v, ok := config["dns_cache_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.DnsCacheConfig = &container.DnsCacheConfig{
			Enabled:         addon["enabled"].(bool),
			ForceSendFields: []string{"Enabled"},
		}
	}

	if v, ok := config["gce_persistent_disk_csi_driver_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.GcePersistentDiskCsiDriverConfig = &container.GcePersistentDiskCsiDriverConfig{
			Enabled:         addon["enabled"].(bool),
			ForceSendFields: []string{"Enabled"},
		}
	}
	if v, ok := config["gke_backup_agent_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.GkeBackupAgentConfig = &container.GkeBackupAgentConfig{
			Enabled:         addon["enabled"].(bool),
			ForceSendFields: []string{"Enabled"},
		}
	}
	if v, ok := config["config_connector_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.ConfigConnectorConfig = &container.ConfigConnectorConfig{
			Enabled:         addon["enabled"].(bool),
			ForceSendFields: []string{"Enabled"},
		}
	}
	if v, ok := config["gcs_fuse_csi_driver_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.GcsFuseCsiDriverConfig = &container.GcsFuseCsiDriverConfig{
			Enabled:         addon["enabled"].(bool),
			ForceSendFields: []string{"Enabled"},
		}
	}

	if v, ok := config["stateful_ha_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.StatefulHaConfig = &container.StatefulHAConfig{
			Enabled:         addon["enabled"].(bool),
			ForceSendFields: []string{"Enabled"},
		}
	}

	if v, ok := config["ray_operator_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.RayOperatorConfig = &container.RayOperatorConfig{
			Enabled:         addon["enabled"].(bool),
			ForceSendFields: []string{"Enabled"},
		}
		if v, ok := addon["ray_cluster_logging_config"]; ok && len(v.([]interface{})) > 0 {
			loggingConfig := v.([]interface{})[0].(map[string]interface{})
			ac.RayOperatorConfig.RayClusterLoggingConfig = &container.RayClusterLoggingConfig{
				Enabled:         loggingConfig["enabled"].(bool),
				ForceSendFields: []string{"Enabled"},
			}
		}
		if v, ok := addon["ray_cluster_monitoring_config"]; ok && len(v.([]interface{})) > 0 {
			loggingConfig := v.([]interface{})[0].(map[string]interface{})
			ac.RayOperatorConfig.RayClusterMonitoringConfig = &container.RayClusterMonitoringConfig{
				Enabled:         loggingConfig["enabled"].(bool),
				ForceSendFields: []string{"Enabled"},
			}
		}
	}

	if v, ok := config["parallelstore_csi_driver_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.ParallelstoreCsiDriverConfig = &container.ParallelstoreCsiDriverConfig{
			Enabled:         addon["enabled"].(bool),
			ForceSendFields: []string{"Enabled"},
		}
	}

	if v, ok := config["lustre_csi_driver_config"]; ok && len(v.([]interface{})) > 0 {
		lustreConfig := v.([]interface{})[0].(map[string]interface{})
		ac.LustreCsiDriverConfig = &container.LustreCsiDriverConfig{
			Enabled:         lustreConfig["enabled"].(bool),
			ForceSendFields: []string{"Enabled"},
		}

		// Check for enable_legacy_lustre_port
		if val, ok := lustreConfig["enable_legacy_lustre_port"]; ok {
			ac.LustreCsiDriverConfig.EnableLegacyLustrePort = val.(bool)
			ac.LustreCsiDriverConfig.ForceSendFields = append(ac.LustreCsiDriverConfig.ForceSendFields, "EnableLegacyLustrePort")
		}
	}

	return ac
}

func expandPodCidrOverprovisionConfig(configured interface{}) *container.PodCIDROverprovisionConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &container.PodCIDROverprovisionConfig{
		Disable:         config["disabled"].(bool),
		ForceSendFields: []string{"Disable"},
	}
}

func expandPodIpv4RangeNames(configured interface{}) []string {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	var ranges []string
	for _, rawRange := range l {
		ranges = append(ranges, rawRange.(string))
	}
	return ranges
}

func expandAdditionalIpRangesConfigs(configured interface{}, d tpgresource.TerraformResourceData, c *transport_tpg.Config) ([]*container.AdditionalIPRangesConfig, error) {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	var additionalIpRangesConfig []*container.AdditionalIPRangesConfig
	for _, rawConfig := range l {
		config := rawConfig.(map[string]interface{})
		subnetwork, err := tpgresource.ParseSubnetworkFieldValue(config["subnetwork"].(string), d, c)
		if err != nil {
			return nil, err
		}
		additionalIpRangesConfig = append(additionalIpRangesConfig, &container.AdditionalIPRangesConfig{
			Subnetwork:        subnetwork.RelativeLink(),
			PodIpv4RangeNames: expandPodIpv4RangeNames(config["pod_ipv4_range_names"]),
		})
	}

	return additionalIpRangesConfig, nil
}

func expandIPAllocationPolicy(configured interface{}, d tpgresource.TerraformResourceData, networkingMode string, autopilot bool, c *transport_tpg.Config) (*container.IPAllocationPolicy, []*container.AdditionalIPRangesConfig, error) {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		if networkingMode == "VPC_NATIVE" {
			return nil, nil, nil
		}
		return &container.IPAllocationPolicy{
			UseIpAliases:    false,
			UseRoutes:       true,
			StackType:       "IPV4",
			ForceSendFields: []string{"UseIpAliases"},
		}, nil, nil
	}

	config := l[0].(map[string]interface{})
	stackType := config["stack_type"].(string)

	// We expand and return additional_ip_ranges_config separately because
	// this field is OUTPUT_ONLY for ClusterCreate RPCs. Instead, during the
	// Terraform Create flow, we follow the CreateCluster (without
	// additional_ip_ranges_config populated) with an UpdateCluster (_with_
	// additional_ip_ranges_config populated).
	additionalIpRangesConfigs, err := expandAdditionalIpRangesConfigs(config["additional_ip_ranges_config"], d, c)
	if err != nil {
		return nil, nil, err
	}

	return &container.IPAllocationPolicy{
		UseIpAliases:               networkingMode == "VPC_NATIVE" || networkingMode == "",
		ClusterIpv4CidrBlock:       config["cluster_ipv4_cidr_block"].(string),
		ServicesIpv4CidrBlock:      config["services_ipv4_cidr_block"].(string),
		ClusterSecondaryRangeName:  config["cluster_secondary_range_name"].(string),
		ServicesSecondaryRangeName: config["services_secondary_range_name"].(string),
		ForceSendFields:            []string{"UseIpAliases"},
		UseRoutes:                  networkingMode == "ROUTES",
		StackType:                  stackType,
		PodCidrOverprovisionConfig: expandPodCidrOverprovisionConfig(config["pod_cidr_overprovision_config"]),
		AutoIpamConfig:             expandAutoIpamConfig(config["auto_ipam_config"]),
		NetworkTierConfig:          expandNetworkTierConfig(config["network_tier_config"]),
	}, additionalIpRangesConfigs, nil
}

func expandNetworkTierConfig(configured interface{}) *container.NetworkTierConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
	return &container.NetworkTierConfig{
		NetworkTier: config["network_tier"].(string),
	}
}

func expandAutoIpamConfig(configured interface{}) *container.AutoIpamConfig {
	l, ok := configured.([]interface{})
	if !ok || len(l) == 0 || l[0] == nil {
		return nil
	}

	return &container.AutoIpamConfig{
		Enabled: l[0].(map[string]interface{})["enabled"].(bool),
	}
}

func expandMaintenancePolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) *container.MaintenancePolicy {
	// We have to perform a full Get() as part of this, to get the fingerprint.  We can't do this
	// at any other time, because the fingerprint update might happen between plan and apply.
	// We can omit error checks, since to have gotten this far, a project is definitely configured.
	project, _ := tpgresource.GetProject(d, config)
	location, _ := tpgresource.GetLocation(d, config)
	clusterName := d.Get("name").(string)
	name := containerClusterFullName(project, location, clusterName)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil
	}
	clusterGetCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.Get(name)
	if config.UserProjectOverride {
		clusterGetCall.Header().Add("X-Goog-User-Project", project)
	}
	cluster, _ := clusterGetCall.Do()
	resourceVersion := ""
	exclusions := make(map[string]container.TimeWindow)
	if cluster != nil && cluster.MaintenancePolicy != nil {
		// If the cluster doesn't exist or if there is a read error of any kind, we will pass in an empty
		// resourceVersion.  If there happens to be a change to maintenance policy, we will fail at that
		// point.  This is a compromise between code cleanliness and a slightly worse user experience in
		// an unlikely error case - we choose code cleanliness.
		resourceVersion = cluster.MaintenancePolicy.ResourceVersion

		// Having a MaintenancePolicy doesn't mean that you need MaintenanceExclusions, but if they were set,
		// they need to be assigned to exclusions.
		if cluster.MaintenancePolicy.Window != nil && cluster.MaintenancePolicy.Window.MaintenanceExclusions != nil {
			exclusions = cluster.MaintenancePolicy.Window.MaintenanceExclusions
		}
	}

	configured := d.Get("maintenance_policy")
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return &container.MaintenancePolicy{
			ResourceVersion: resourceVersion,
			Window: &container.MaintenanceWindow{
				MaintenanceExclusions: exclusions,
			},
		}
	}
	maintenancePolicy := l[0].(map[string]interface{})

	if maintenanceExclusions, ok := maintenancePolicy["maintenance_exclusion"]; ok {
		for k := range exclusions {
			delete(exclusions, k)
		}
		for _, me := range maintenanceExclusions.(*schema.Set).List() {
			exclusion := me.(map[string]interface{})
			exclusions[exclusion["exclusion_name"].(string)] = container.TimeWindow{
				StartTime: exclusion["start_time"].(string),
				EndTime:   exclusion["end_time"].(string),
			}
			if exclusionOptions, ok := exclusion["exclusion_options"]; ok && len(exclusionOptions.([]interface{})) > 0 {
				meo := exclusionOptions.([]interface{})[0].(map[string]interface{})
				mex := exclusions[exclusion["exclusion_name"].(string)]
				mex.MaintenanceExclusionOptions = &container.MaintenanceExclusionOptions{
					Scope:           meo["scope"].(string),
					ForceSendFields: []string{"Scope"},
				}
				if len(meo["end_time_behavior"].(string)) > 0 {
					mex.MaintenanceExclusionOptions.EndTimeBehavior = meo["end_time_behavior"].(string)
					mex.EndTime = ""
				}
				exclusions[exclusion["exclusion_name"].(string)] = mex
			}
		}
	}

	if dailyMaintenanceWindow, ok := maintenancePolicy["daily_maintenance_window"]; ok && len(dailyMaintenanceWindow.([]interface{})) > 0 {
		dmw := dailyMaintenanceWindow.([]interface{})[0].(map[string]interface{})
		startTime := dmw["start_time"].(string)
		return &container.MaintenancePolicy{
			Window: &container.MaintenanceWindow{
				MaintenanceExclusions: exclusions,
				DailyMaintenanceWindow: &container.DailyMaintenanceWindow{
					StartTime: startTime,
				},
			},
			ResourceVersion: resourceVersion,
		}
	}
	if recurringWindow, ok := maintenancePolicy["recurring_window"]; ok && len(recurringWindow.([]interface{})) > 0 {
		rw := recurringWindow.([]interface{})[0].(map[string]interface{})
		return &container.MaintenancePolicy{
			Window: &container.MaintenanceWindow{
				MaintenanceExclusions: exclusions,
				RecurringWindow: &container.RecurringTimeWindow{
					Window: &container.TimeWindow{
						StartTime: rw["start_time"].(string),
						EndTime:   rw["end_time"].(string),
					},
					Recurrence: rw["recurrence"].(string),
				},
			},
			ResourceVersion: resourceVersion,
		}
	}
	return nil
}

func expandClusterAutoscaling(configured interface{}, d tpgresource.TerraformResourceData) *container.ClusterAutoscaling {
	l, ok := configured.([]interface{})
	enableAutopilot := false
	if v, ok := d.GetOk("enable_autopilot"); ok && v == true {
		enableAutopilot = true
	}
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		if enableAutopilot {
			return nil
		}
		return &container.ClusterAutoscaling{
			EnableNodeAutoprovisioning: false,
			ForceSendFields:            []string{"EnableNodeAutoprovisioning"},
		}
	}

	config := l[0].(map[string]interface{})

	// Conditionally provide an empty list to preserve a legacy 2.X behaviour
	// when `enabled` is false and resource_limits is unset, allowing users to
	// explicitly disable the feature. resource_limits don't work when node
	// auto-provisioning is disabled at time of writing. This may change API-side
	// in the future though, as the feature is intended to apply to both node
	// auto-provisioning and node autoscaling.
	var resourceLimits []*container.ResourceLimit
	if limits, ok := config["resource_limits"]; ok {
		resourceLimits = make([]*container.ResourceLimit, 0)
		if lmts, ok := limits.([]interface{}); ok {
			for _, v := range lmts {
				limit := v.(map[string]interface{})
				resourceLimits = append(resourceLimits,
					&container.ResourceLimit{
						ResourceType: limit["resource_type"].(string),
						// Here we're relying on *not* setting ForceSendFields for 0-values.
						Minimum: int64(limit["minimum"].(int)),
						Maximum: int64(limit["maximum"].(int)),
					})
			}
		}
	}
	var defaultCCConfig *container.DefaultComputeClassConfig
	if defaultCCEnabled, ok := config["default_compute_class_enabled"]; ok {
		defaultCCConfig = &container.DefaultComputeClassConfig{
			Enabled: defaultCCEnabled.(bool),
		}
	}
	return &container.ClusterAutoscaling{
		EnableNodeAutoprovisioning:       config["enabled"].(bool),
		ResourceLimits:                   resourceLimits,
		DefaultComputeClassConfig:        defaultCCConfig,
		AutoscalingProfile:               config["autoscaling_profile"].(string),
		AutoprovisioningNodePoolDefaults: expandAutoProvisioningDefaults(config["auto_provisioning_defaults"], d),
		AutoprovisioningLocations:        tpgresource.ConvertStringArr(config["auto_provisioning_locations"].([]interface{})),
	}
}

func expandAutoProvisioningDefaults(configured interface{}, d tpgresource.TerraformResourceData) *container.AutoprovisioningNodePoolDefaults {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return &container.AutoprovisioningNodePoolDefaults{}
	}
	config := l[0].(map[string]interface{})

	npd := &container.AutoprovisioningNodePoolDefaults{
		OauthScopes:     tpgresource.ConvertStringArr(config["oauth_scopes"].([]interface{})),
		ServiceAccount:  config["service_account"].(string),
		DiskSizeGb:      int64(config["disk_size"].(int)),
		DiskType:        config["disk_type"].(string),
		ImageType:       config["image_type"].(string),
		BootDiskKmsKey:  config["boot_disk_kms_key"].(string),
		Management:      expandManagement(config["management"]),
		UpgradeSettings: expandUpgradeSettings(config["upgrade_settings"]),
	}

	if v, ok := config["shielded_instance_config"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		npd.ShieldedInstanceConfig = &container.ShieldedInstanceConfig{
			EnableSecureBoot:          conf["enable_secure_boot"].(bool),
			EnableIntegrityMonitoring: conf["enable_integrity_monitoring"].(bool),
		}
	}

	cpu := config["min_cpu_platform"].(string)
	// the only way to unset the field is to pass "automatic" as its value
	if cpu == "" {
		cpu = "automatic"
	}
	npd.MinCpuPlatform = cpu

	return npd
}

func expandUpgradeSettings(configured interface{}) *container.UpgradeSettings {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return &container.UpgradeSettings{}
	}
	config := l[0].(map[string]interface{})

	upgradeSettings := &container.UpgradeSettings{
		MaxSurge:          int64(config["max_surge"].(int)),
		MaxUnavailable:    int64(config["max_unavailable"].(int)),
		Strategy:          config["strategy"].(string),
		BlueGreenSettings: expandBlueGreenSettings(config["blue_green_settings"]),
	}

	return upgradeSettings
}

func expandBlueGreenSettings(configured interface{}) *container.BlueGreenSettings {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return &container.BlueGreenSettings{}
	}
	config := l[0].(map[string]interface{})

	blueGreenSettings := &container.BlueGreenSettings{
		NodePoolSoakDuration:  config["node_pool_soak_duration"].(string),
		StandardRolloutPolicy: expandStandardRolloutPolicy(config["standard_rollout_policy"]),
	}

	return blueGreenSettings
}

func expandStandardRolloutPolicy(configured interface{}) *container.StandardRolloutPolicy {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return &container.StandardRolloutPolicy{}
	}

	config := l[0].(map[string]interface{})
	standardRolloutPolicy := &container.StandardRolloutPolicy{
		BatchPercentage:   config["batch_percentage"].(float64),
		BatchNodeCount:    int64(config["batch_node_count"].(int)),
		BatchSoakDuration: config["batch_soak_duration"].(string),
	}

	return standardRolloutPolicy
}

func expandManagement(configured interface{}) *container.NodeManagement {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return nil
	}
	config := l[0].(map[string]interface{})

	mng := &container.NodeManagement{
		AutoUpgrade:    config["auto_upgrade"].(bool),
		AutoRepair:     config["auto_repair"].(bool),
		UpgradeOptions: expandUpgradeOptions(config["upgrade_options"]),
	}

	return mng
}

func expandUpgradeOptions(configured interface{}) *container.AutoUpgradeOptions {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return &container.AutoUpgradeOptions{}
	}
	config := l[0].(map[string]interface{})

	upgradeOptions := &container.AutoUpgradeOptions{
		AutoUpgradeStartTime: config["auto_upgrade_start_time"].(string),
		Description:          config["description"].(string),
	}

	return upgradeOptions
}

func expandAuthenticatorGroupsConfig(configured interface{}) *container.AuthenticatorGroupsConfig {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}
	result := &container.AuthenticatorGroupsConfig{}
	config := l[0].(map[string]interface{})
	if securityGroup, ok := config["security_group"]; ok {
		result.Enabled = true
		result.SecurityGroup = securityGroup.(string)
	}
	return result
}

func expandSecurityPostureConfig(configured interface{}) *container.SecurityPostureConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	spc := &container.SecurityPostureConfig{}
	spConfig := l[0].(map[string]interface{})
	if v, ok := spConfig["mode"]; ok {
		spc.Mode = v.(string)
	}

	if v, ok := spConfig["vulnerability_mode"]; ok {
		spc.VulnerabilityMode = v.(string)
	}
	return spc
}

func expandNotificationConfig(configured interface{}) *container.NotificationConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return &container.NotificationConfig{
			Pubsub: &container.PubSub{
				Enabled: false,
			},
		}
	}

	notificationConfig := l[0].(map[string]interface{})
	if v, ok := notificationConfig["pubsub"]; ok {
		if len(v.([]interface{})) > 0 {
			pubsub := notificationConfig["pubsub"].([]interface{})[0].(map[string]interface{})

			nc := &container.NotificationConfig{
				Pubsub: &container.PubSub{
					Enabled: pubsub["enabled"].(bool),
					Topic:   pubsub["topic"].(string),
				},
			}

			if vv, ok := pubsub["filter"]; ok && len(vv.([]interface{})) > 0 {
				filter := vv.([]interface{})[0].(map[string]interface{})
				eventType := filter["event_type"].([]interface{})
				nc.Pubsub.Filter = &container.Filter{
					EventType: tpgresource.ConvertStringArr(eventType),
				}
			}

			return nc
		}
	}

	return &container.NotificationConfig{
		Pubsub: &container.PubSub{
			Enabled: false,
		},
	}
}

func expandBinaryAuthorization(configured interface{}) *container.BinaryAuthorization {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return &container.BinaryAuthorization{
			Enabled:         false,
			ForceSendFields: []string{"Enabled"},
		}
	}
	config := l[0].(map[string]interface{})
	return &container.BinaryAuthorization{
		Enabled:        config["enabled"].(bool),
		EvaluationMode: config["evaluation_mode"].(string),
	}
}

func expandMasterAuth(configured interface{}) *container.MasterAuth {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	masterAuth := l[0].(map[string]interface{})
	result := &container.MasterAuth{}

	if v, ok := masterAuth["client_certificate_config"]; ok {
		if len(v.([]interface{})) > 0 {
			clientCertificateConfig := masterAuth["client_certificate_config"].([]interface{})[0].(map[string]interface{})

			result.ClientCertificateConfig = &container.ClientCertificateConfig{
				IssueClientCertificate: clientCertificateConfig["issue_client_certificate"].(bool),
			}
		}
	}

	return result
}

func expandMasterAuthorizedNetworksConfig(d tpgresource.TerraformResourceData) *container.MasterAuthorizedNetworksConfig {
	v := d.Get("master_authorized_networks_config").([]interface{})
	if len(v) == 0 {
		// TF doesn't have an explicit enabled field for authorized networks, it is assumed to be enabled based
		// on whether the master_authorized_networks_conifg is present at all. The GKE API pays attention to the
		// field presence of authorized_networks_config, so it's important to explicitly include enabled = false
		// to allow disabling this during updates.
		return &container.MasterAuthorizedNetworksConfig{
			Enabled: false,
		}
	}

	result := &container.MasterAuthorizedNetworksConfig{
		Enabled: true,
	}
	if v, ok := d.GetOk("master_authorized_networks_config.0.cidr_blocks"); ok {
		result.CidrBlocks = expandManCidrBlocks(v)
	}
	if v, ok := d.GetOkExists("master_authorized_networks_config.0.gcp_public_cidrs_access_enabled"); ok {
		result.GcpPublicCidrsAccessEnabled = v.(bool)
		result.ForceSendFields = append(result.ForceSendFields, "GcpPublicCidrsAccessEnabled")
	}
	if v, ok := d.GetOkExists("master_authorized_networks_config.0.private_endpoint_enforcement_enabled"); ok {
		result.PrivateEndpointEnforcementEnabled = v.(bool)
		result.ForceSendFields = append(result.ForceSendFields, "PrivateEndpointEnforcementEnabled")
	}
	return result
}

func expandAnonymousAuthenticationConfig(configured interface{}) *container.AnonymousAuthenticationConfig {
	l, ok := configured.([]interface{})
	if len(l) == 0 || l[0] == nil || !ok {
		return nil
	}

	anonAuthConfig := l[0].(map[string]interface{})
	result := container.AnonymousAuthenticationConfig{}

	if v, ok := anonAuthConfig["mode"]; ok {
		if mode, ok := v.(string); ok && mode != "" {
			result.Mode = mode
		}
	}
	return &result
}

func expandManCidrBlocks(configured interface{}) []*container.CidrBlock {
	config, ok := configured.(*schema.Set)
	if !ok {
		return nil
	}
	cidrBlocks := config.List()
	result := make([]*container.CidrBlock, 0)
	for _, v := range cidrBlocks {
		cidrBlock := v.(map[string]interface{})
		result = append(result, &container.CidrBlock{
			CidrBlock:   cidrBlock["cidr_block"].(string),
			DisplayName: cidrBlock["display_name"].(string),
		})
	}
	return result
}

func expandNetworkPolicy(configured interface{}) *container.NetworkPolicy {
	result := &container.NetworkPolicy{}
	l := configured.([]interface{})
	if len(l) == 0 {
		return result
	}
	config := l[0].(map[string]interface{})
	if enabled, ok := config["enabled"]; ok && enabled.(bool) {
		result.Enabled = true
		if provider, ok := config["provider"]; ok {
			result.Provider = provider.(string)
		}
	}
	return result
}

// Most of the contents of PrivateClusterConfig have been deprecated in the underlying API and replaced by ControlPlaneEndpointsConfig.
// This function primarily handles the sole remaining undeprecated field, master_ipv4_cidr_block.
// Unfortunately, since the private_cluster_config.enable_private_nodes proto field is not marked optional, we can't just leave it
// unset, as that would implicitly use the value false, and it must match the value of network_config.default_enable_private_nodes.
// This function is intended to be called only during cluster creation, after the network_config field is been configured.
// This is possible because master_ipv4_cidr_block is immutable.
func expandPrivateClusterConfigMasterIpv4CidrBlock(configured interface{}, c *container.Cluster) *container.PrivateClusterConfig {
	v := configured.(string)

	return &container.PrivateClusterConfig{
		MasterIpv4CidrBlock: v,
		EnablePrivateNodes:  c.NetworkConfig.DefaultEnablePrivateNodes,
		ForceSendFields:     []string{"MasterIpv4CidrBlock"},
	}
}

func expandDefaultEnablePrivateNodes(d tpgresource.TerraformResourceData) bool {
	b, ok := d.GetOk("private_cluster_config.0.enable_private_nodes")
	if ok {
		v, _ := b.(bool)
		return v
	}
	return false
}

func expandControlPlaneEndpointsConfig(d tpgresource.TerraformResourceData) *container.ControlPlaneEndpointsConfig {
	dns := &container.DNSEndpointConfig{}
	if v := d.Get("control_plane_endpoints_config.0.dns_endpoint_config.0.allow_external_traffic"); v != nil {
		dns.AllowExternalTraffic = v.(bool)
		dns.ForceSendFields = []string{"AllowExternalTraffic"}
	}

	if v := d.Get("control_plane_endpoints_config.0.dns_endpoint_config.0.enable_k8s_tokens_via_dns"); v != nil {
		dns.EnableK8sTokensViaDns = v.(bool)
		dns.ForceSendFields = []string{"EnableK8sTokensViaDns"}
	}

	if v := d.Get("control_plane_endpoints_config.0.dns_endpoint_config.0.enable_k8s_certs_via_dns"); v != nil {
		dns.EnableK8sCertsViaDns = v.(bool)
		dns.ForceSendFields = []string{"EnableK8sCertsViaDns"}
	}

	ip := &container.IPEndpointsConfig{
		Enabled:         true,
		ForceSendFields: []string{"Enabled"},
	}
	if v := d.Get("control_plane_endpoints_config.0.ip_endpoints_config.#"); v != 0 {
		ip.Enabled = d.Get("control_plane_endpoints_config.0.ip_endpoints_config.0.enabled").(bool)

		if !ip.Enabled {
			return &container.ControlPlaneEndpointsConfig{
				DnsEndpointConfig: dns,
				IpEndpointsConfig: ip,
			}
		}
	}
	if v := d.Get("private_cluster_config.0.enable_private_endpoint"); v != nil {
		ip.EnablePublicEndpoint = !v.(bool)
		ip.ForceSendFields = append(ip.ForceSendFields, "EnablePublicEndpoint")
	}
	if v := d.Get("private_cluster_config.0.private_endpoint_subnetwork"); v != nil {
		ip.PrivateEndpointSubnetwork = v.(string)
		ip.ForceSendFields = append(ip.ForceSendFields, "PrivateEndpointSubnetwork")
	}
	if v := d.Get("private_cluster_config.0.master_global_access_config.0.enabled"); v != nil {
		ip.GlobalAccess = v.(bool)
		ip.ForceSendFields = append(ip.ForceSendFields, "GlobalAccess")
	}
	ip.AuthorizedNetworksConfig = expandMasterAuthorizedNetworksConfig(d)

	return &container.ControlPlaneEndpointsConfig{
		DnsEndpointConfig: dns,
		IpEndpointsConfig: ip,
	}
}

func expandVerticalPodAutoscaling(configured interface{}) *container.VerticalPodAutoscaling {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &container.VerticalPodAutoscaling{
		Enabled: config["enabled"].(bool),
	}
}

func expandServiceExternalIpsConfig(configured interface{}) *container.ServiceExternalIPsConfig {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &container.ServiceExternalIPsConfig{
		Enabled:         config["enabled"].(bool),
		ForceSendFields: []string{"Enabled"},
	}
}

func expandMeshCertificates(configured interface{}) *container.MeshCertificates {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &container.MeshCertificates{
		EnableCertificates: config["enable_certificates"].(bool),
		ForceSendFields:    []string{"EnableCertificates"},
	}
}

func expandDatabaseEncryption(configured interface{}) *container.DatabaseEncryption {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &container.DatabaseEncryption{
		State:   config["state"].(string),
		KeyName: config["key_name"].(string),
	}
}

func expandReleaseChannel(configured interface{}) *container.ReleaseChannel {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &container.ReleaseChannel{
		Channel: config["channel"].(string),
	}
}

func expandGkeAutoUpgradeConfig(configured interface{}) *container.GkeAutoUpgradeConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &container.GkeAutoUpgradeConfig{
		PatchMode: config["patch_mode"].(string),
	}
}

func expandDefaultSnatStatus(configured interface{}) *container.DefaultSnatStatus {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &container.DefaultSnatStatus{
		Disabled:        config["disabled"].(bool),
		ForceSendFields: []string{"Disabled"},
	}

}

func expandWorkloadIdentityConfig(configured interface{}) *container.WorkloadIdentityConfig {
	l := configured.([]interface{})
	v := &container.WorkloadIdentityConfig{}

	// this API considers unset and set-to-empty equivalent. Note that it will
	// always return an empty block given that we always send one, but clusters
	// not created in TF will not always return one (and may return nil)
	if len(l) == 0 || l[0] == nil {
		return v
	}

	config := l[0].(map[string]interface{})
	v.WorkloadPool = config["workload_pool"].(string)

	return v
}

func expandIdentityServiceConfig(configured interface{}) *container.IdentityServiceConfig {
	l := configured.([]interface{})
	v := &container.IdentityServiceConfig{}

	config := l[0].(map[string]interface{})
	v.Enabled = config["enabled"].(bool)

	return v
}

func expandPodAutoscaling(configured interface{}) *container.PodAutoscaling {
	if configured == nil {
		return nil
	}

	podAutoscaling := &container.PodAutoscaling{}

	configs := configured.([]interface{})

	if len(configs) == 0 || configs[0] == nil {
		return nil
	}

	config := configs[0].(map[string]interface{})

	if v, ok := config["hpa_profile"]; ok {
		podAutoscaling.HpaProfile = v.(string)
	}

	return podAutoscaling
}

func expandSecretManagerConfig(configured interface{}) *container.SecretManagerConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
	sc := &container.SecretManagerConfig{
		Enabled:         config["enabled"].(bool),
		ForceSendFields: []string{"Enabled"},
	}
	if autoRotation, ok := config["rotation_config"]; ok {
		if autoRotationList, ok := autoRotation.([]interface{}); ok {
			if len(autoRotationList) > 0 {
				autoRotationConfig := autoRotationList[0].(map[string]interface{})
				if rotationInterval, ok := autoRotationConfig["rotation_interval"].(string); ok && rotationInterval != "" {
					sc.RotationConfig = &container.RotationConfig{
						Enabled:          autoRotationConfig["enabled"].(bool),
						RotationInterval: rotationInterval,
					}
				} else {
					sc.RotationConfig = &container.RotationConfig{
						Enabled: autoRotationConfig["enabled"].(bool),
					}
				}
			}
		}
	}
	return sc
}

func expandDefaultMaxPodsConstraint(v interface{}) *container.MaxPodsConstraint {
	if v == nil {
		return nil
	}

	return &container.MaxPodsConstraint{
		MaxPodsPerNode: int64(v.(int)),
	}
}

func expandCostManagementConfig(configured interface{}) *container.CostManagementConfig {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}

	config := l[0].(map[string]interface{})
	return &container.CostManagementConfig{
		Enabled:         config["enabled"].(bool),
		ForceSendFields: []string{"Enabled"},
	}
}

func expandResourceUsageExportConfig(configured interface{}) *container.ResourceUsageExportConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return &container.ResourceUsageExportConfig{}
	}

	resourceUsageConfig := l[0].(map[string]interface{})

	result := &container.ResourceUsageExportConfig{
		EnableNetworkEgressMetering: resourceUsageConfig["enable_network_egress_metering"].(bool),
		ConsumptionMeteringConfig: &container.ConsumptionMeteringConfig{
			Enabled:         resourceUsageConfig["enable_resource_consumption_metering"].(bool),
			ForceSendFields: []string{"Enabled"},
		},
		ForceSendFields: []string{"EnableNetworkEgressMetering"},
	}
	if _, ok := resourceUsageConfig["bigquery_destination"]; ok {
		destinationArr := resourceUsageConfig["bigquery_destination"].([]interface{})
		if len(destinationArr) > 0 && destinationArr[0] != nil {
			bigqueryDestination := destinationArr[0].(map[string]interface{})
			if _, ok := bigqueryDestination["dataset_id"]; ok {
				result.BigqueryDestination = &container.BigQueryDestination{
					DatasetId: bigqueryDestination["dataset_id"].(string),
				}
			}
		}
	}
	return result
}

func expandDnsConfig(configured interface{}) *container.DNSConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
	return &container.DNSConfig{
		AdditiveVpcScopeDnsDomain: config["additive_vpc_scope_dns_domain"].(string),
		ClusterDns:                config["cluster_dns"].(string),
		ClusterDnsScope:           config["cluster_dns_scope"].(string),
		ClusterDnsDomain:          config["cluster_dns_domain"].(string),
	}
}

func expandNetworkPerformanceConfig(configured interface{}) *container.ClusterNetworkPerformanceConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
	return &container.ClusterNetworkPerformanceConfig{
		TotalEgressBandwidthTier: config["total_egress_bandwidth_tier"].(string),
	}
}

func expandGatewayApiConfig(configured interface{}) *container.GatewayAPIConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
	return &container.GatewayAPIConfig{
		Channel: config["channel"].(string),
	}
}

func expandFleet(configured interface{}) *container.Fleet {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
	return &container.Fleet{
		Project:        config["project"].(string),
		MembershipType: config["membership_type"].(string),
	}
}

func expandUserManagedKeysConfig(configured interface{}) *container.UserManagedKeysConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
	umkc := &container.UserManagedKeysConfig{
		ClusterCa:                     config["cluster_ca"].(string),
		EtcdApiCa:                     config["etcd_api_ca"].(string),
		EtcdPeerCa:                    config["etcd_peer_ca"].(string),
		AggregationCa:                 config["aggregation_ca"].(string),
		ControlPlaneDiskEncryptionKey: config["control_plane_disk_encryption_key"].(string),
		GkeopsEtcdBackupEncryptionKey: config["gkeops_etcd_backup_encryption_key"].(string),
	}
	if v, ok := config["service_account_signing_keys"]; ok {
		sk := v.(*schema.Set)
		skss := tpgresource.ConvertStringSet(sk)
		if len(skss) > 0 {
			umkc.ServiceAccountSigningKeys = skss
		}
	}
	if v, ok := config["service_account_verification_keys"]; ok {
		vk := v.(*schema.Set)
		vkss := tpgresource.ConvertStringSet(vk)
		if len(vkss) > 0 {
			umkc.ServiceAccountVerificationKeys = vkss
		}
	}
	return umkc
}

func expandEnableK8sBetaApis(configured interface{}, enabledAPIs []string) *container.K8sBetaAPIConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
	result := &container.K8sBetaAPIConfig{}
	if v, ok := config["enabled_apis"]; ok {
		notEnabledAPIsSet := v.(*schema.Set)
		for _, enabledAPI := range enabledAPIs {
			if notEnabledAPIsSet.Contains(enabledAPI) {
				notEnabledAPIsSet.Remove(enabledAPI)
			}
		}

		result.EnabledApis = tpgresource.ConvertStringSet(notEnabledAPIsSet)
	}

	return result
}

func expandContainerClusterLoggingConfig(configured interface{}) *container.LoggingConfig {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}

	var components []string
	if l[0] != nil {
		config := l[0].(map[string]interface{})
		components = tpgresource.ConvertStringArr(config["enable_components"].([]interface{}))
	}

	return &container.LoggingConfig{
		ComponentConfig: &container.LoggingComponentConfig{
			EnableComponents: components,
		},
	}
}

func expandMonitoringConfig(configured interface{}) *container.MonitoringConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	mc := &container.MonitoringConfig{}
	config := l[0].(map[string]interface{})
	if v, ok := config["enable_components"]; ok {
		enable_components := v.([]interface{})
		mc.ComponentConfig = &container.MonitoringComponentConfig{
			EnableComponents: tpgresource.ConvertStringArr(enable_components),
		}
	}
	if v, ok := config["managed_prometheus"]; ok && len(v.([]interface{})) > 0 {
		managed_prometheus := v.([]interface{})[0].(map[string]interface{})
		mc.ManagedPrometheusConfig = &container.ManagedPrometheusConfig{
			Enabled: managed_prometheus["enabled"].(bool),
		}
		if autoMonitoring, ok := managed_prometheus["auto_monitoring_config"]; ok {
			if autoMonitoringList, ok := autoMonitoring.([]interface{}); ok {
				if len(autoMonitoringList) > 0 {
					autoMonitoringConfig := autoMonitoringList[0].(map[string]interface{})
					if scope, ok := autoMonitoringConfig["scope"].(string); ok && scope != "" {
						mc.ManagedPrometheusConfig.AutoMonitoringConfig = &container.AutoMonitoringConfig{
							Scope: scope,
						}
					}
				}
			}
		}
	}

	if v, ok := config["advanced_datapath_observability_config"]; ok && len(v.([]interface{})) > 0 {
		advanced_datapath_observability_config := v.([]interface{})[0].(map[string]interface{})
		mc.AdvancedDatapathObservabilityConfig = &container.AdvancedDatapathObservabilityConfig{
			EnableMetrics:   advanced_datapath_observability_config["enable_metrics"].(bool),
			EnableRelay:     advanced_datapath_observability_config["enable_relay"].(bool),
			ForceSendFields: []string{"EnableRelay"},
		}
	}

	return mc
}

func expandContainerClusterAuthenticatorGroupsConfig(configured interface{}) *container.AuthenticatorGroupsConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
	result := &container.AuthenticatorGroupsConfig{}
	if securityGroup, ok := config["security_group"]; ok {
		if securityGroup == nil || securityGroup.(string) == "" {
			result.Enabled = false
		} else {
			result.Enabled = true
			result.SecurityGroup = securityGroup.(string)
		}
	}
	return result
}

func expandNodePoolDefaults(configured interface{}) *container.NodePoolDefaults {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return nil
	}
	nodePoolDefaults := &container.NodePoolDefaults{}
	config := l[0].(map[string]interface{})
	if v, ok := config["node_config_defaults"]; ok && len(v.([]interface{})) > 0 {
		nodePoolDefaults.NodeConfigDefaults = expandNodeConfigDefaults(v)
	}
	return nodePoolDefaults
}

func flattenNodePoolDefaults(c *container.NodePoolDefaults) []map[string]interface{} {
	if c == nil {
		return nil
	}

	result := make(map[string]interface{})
	if c.NodeConfigDefaults != nil {
		result["node_config_defaults"] = flattenNodeConfigDefaults(c.NodeConfigDefaults)
	}

	return []map[string]interface{}{result}
}

func expandNodePoolAutoConfig(configured interface{}) *container.NodePoolAutoConfig {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return nil
	}

	npac := &container.NodePoolAutoConfig{}
	config := l[0].(map[string]interface{})

	if v, ok := config["node_kubelet_config"]; ok {
		npac.NodeKubeletConfig = expandKubeletConfig(v)
	}

	if v, ok := config["network_tags"]; ok && len(v.([]interface{})) > 0 {
		npac.NetworkTags = expandNodePoolAutoConfigNetworkTags(v)
	}

	if v, ok := config["resource_manager_tags"]; ok && len(v.(map[string]interface{})) > 0 {
		npac.ResourceManagerTags = expandResourceManagerTags(v)
	}

	return npac
}

func expandNodePoolAutoConfigNetworkTags(configured interface{}) *container.NetworkTags {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	nt := &container.NetworkTags{}
	config := l[0].(map[string]interface{})

	if v, ok := config["tags"]; ok && len(v.([]interface{})) > 0 {
		nt.Tags = tpgresource.ConvertStringArr(v.([]interface{}))
	}
	return nt
}

func expandRBACBindingConfig(configured interface{}) *container.RBACBindingConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
	return &container.RBACBindingConfig{
		EnableInsecureBindingSystemUnauthenticated: config["enable_insecure_binding_system_unauthenticated"].(bool),
		EnableInsecureBindingSystemAuthenticated:   config["enable_insecure_binding_system_authenticated"].(bool),
		ForceSendFields:                            []string{"EnableInsecureBindingSystemUnauthenticated", "EnableInsecureBindingSystemAuthenticated"},
	}
}
