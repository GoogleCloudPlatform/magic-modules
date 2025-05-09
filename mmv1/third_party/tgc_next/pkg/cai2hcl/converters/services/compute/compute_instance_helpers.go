package compute

import (
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tpgresource"

	compute "google.golang.org/api/compute/v0.beta"
)

func flattenAliasIpRange(ranges []*compute.AliasIpRange) []map[string]interface{} {
	rangesSchema := make([]map[string]interface{}, 0, len(ranges))
	for _, ipRange := range ranges {
		rangesSchema = append(rangesSchema, map[string]interface{}{
			"ip_cidr_range":         ipRange.IpCidrRange,
			"subnetwork_range_name": ipRange.SubnetworkRangeName,
		})
	}
	return rangesSchema
}

func flattenScheduling(resp *compute.Scheduling) []map[string]interface{} {
	schedulingMap := make(map[string]interface{}, 0)

	if resp.InstanceTerminationAction != "" {
		schedulingMap["instance_termination_action"] = resp.InstanceTerminationAction
	}

	if resp.MinNodeCpus != 0 {
		schedulingMap["min_node_cpus"] = resp.MinNodeCpus
	}

	if resp.OnHostMaintenance != "MIGRATE" {
		schedulingMap["on_host_maintenance"] = resp.OnHostMaintenance
	}

	if resp.AutomaticRestart != nil && !*resp.AutomaticRestart {
		schedulingMap["automatic_restart"] = *resp.AutomaticRestart
	}

	if resp.Preemptible {
		schedulingMap["preemptible"] = resp.Preemptible
	}

	if resp.NodeAffinities != nil && len(resp.NodeAffinities) > 0 {
		nodeAffinities := []map[string]interface{}{}
		for _, na := range resp.NodeAffinities {
			nodeAffinities = append(nodeAffinities, map[string]interface{}{
				"key":      na.Key,
				"operator": na.Operator,
				"values":   tpgresource.ConvertStringArrToInterface(na.Values),
			})
		}
		schedulingMap["node_affinities"] = nodeAffinities
	}

	if resp.ProvisioningModel != "STANDARD" {
		schedulingMap["provisioning_model"] = resp.ProvisioningModel
	}

	if resp.AvailabilityDomain != 0 {
		schedulingMap["availability_domain"] = resp.AvailabilityDomain
	}

	if resp.MaxRunDuration != nil {
		schedulingMap["max_run_duration"] = flattenComputeMaxRunDuration(resp.MaxRunDuration)
	}

	if resp.OnInstanceStopAction != nil {
		schedulingMap["on_instance_stop_action"] = flattenOnInstanceStopAction(resp.OnInstanceStopAction)
	}

	if resp.HostErrorTimeoutSeconds != 0 {
		schedulingMap["host_error_timeout_seconds"] = resp.HostErrorTimeoutSeconds
	}

	if resp.MaintenanceInterval != "" {
		schedulingMap["maintenance_interval"] = resp.MaintenanceInterval
	}

	if resp.LocalSsdRecoveryTimeout != nil {
		schedulingMap["local_ssd_recovery_timeout"] = flattenComputeLocalSsdRecoveryTimeout(resp.LocalSsdRecoveryTimeout)
	}

	if len(schedulingMap) == 0 {
		return nil
	}

	return []map[string]interface{}{schedulingMap}
}

func flattenComputeMaxRunDuration(v *compute.Duration) []interface{} {
	if v == nil {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["nanos"] = v.Nanos
	transformed["seconds"] = v.Seconds
	return []interface{}{transformed}
}

func flattenOnInstanceStopAction(v *compute.SchedulingOnInstanceStopAction) []interface{} {
	if v == nil {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["discard_local_ssd"] = v.DiscardLocalSsd
	return []interface{}{transformed}
}

func flattenComputeLocalSsdRecoveryTimeout(v *compute.Duration) []interface{} {
	if v == nil {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["nanos"] = v.Nanos
	transformed["seconds"] = v.Seconds
	return []interface{}{transformed}
}

func flattenAccessConfigs(accessConfigs []*compute.AccessConfig) ([]map[string]interface{}, string) {
	flattened := make([]map[string]interface{}, len(accessConfigs))
	natIP := ""
	for i, ac := range accessConfigs {
		flattened[i] = map[string]interface{}{
			"nat_ip":       ac.NatIP,
			"network_tier": ac.NetworkTier,
		}
		if ac.SetPublicPtr {
			flattened[i]["public_ptr_domain_name"] = ac.PublicPtrDomainName
		}
		if natIP == "" {
			natIP = ac.NatIP
		}
		if ac.SecurityPolicy != "" {
			flattened[i]["security_policy"] = ac.SecurityPolicy
		}
	}
	return flattened, natIP
}

func flattenIpv6AccessConfigs(ipv6AccessConfigs []*compute.AccessConfig) []map[string]interface{} {
	flattened := make([]map[string]interface{}, len(ipv6AccessConfigs))
	for i, ac := range ipv6AccessConfigs {
		flattened[i] = map[string]interface{}{
			"network_tier": ac.NetworkTier,
		}
		flattened[i]["public_ptr_domain_name"] = ac.PublicPtrDomainName
		flattened[i]["external_ipv6"] = ac.ExternalIpv6
		flattened[i]["external_ipv6_prefix_length"] = strconv.FormatInt(ac.ExternalIpv6PrefixLength, 10)
		flattened[i]["name"] = ac.Name
		if ac.SecurityPolicy != "" {
			flattened[i]["security_policy"] = ac.SecurityPolicy
		}
	}
	return flattened
}

func flattenNetworkInterfaces(networkInterfaces []*compute.NetworkInterface, project string) ([]map[string]interface{}, string, string, error) {
	flattened := make([]map[string]interface{}, len(networkInterfaces))
	var internalIP, externalIP string

	for i, iface := range networkInterfaces {
		var ac []map[string]interface{}
		ac, externalIP = flattenAccessConfigs(iface.AccessConfigs)

		flattened[i] = map[string]interface{}{
			"network_ip":         iface.NetworkIP,
			"access_config":      ac,
			"alias_ip_range":     flattenAliasIpRange(iface.AliasIpRanges),
			"nic_type":           iface.NicType,
			"ipv6_access_config": flattenIpv6AccessConfigs(iface.Ipv6AccessConfigs),
			"ipv6_address":       iface.Ipv6Address,
		}

		if !strings.HasSuffix(iface.Network, "/default") {
			flattened[i]["network"] = tpgresource.ConvertSelfLinkToV1(iface.Network)
		}

		if !strings.HasSuffix(iface.Subnetwork, "/default") {
			flattened[i]["subnetwork"] = tpgresource.ConvertSelfLinkToV1(iface.Subnetwork)
		}

		subnetProject := utils.ParseFieldValue(iface.Subnetwork, "projects")
		if subnetProject != project {
			flattened[i]["subnetwork_project"] = subnetProject
		}

		if iface.StackType != "IPV4_ONLY" {
			flattened[i]["stack_type"] = iface.StackType
		}

		if iface.QueueCount != 0 {
			flattened[i]["queue_count"] = iface.QueueCount
		}

		if internalIP == "" {
			internalIP = iface.NetworkIP
		}

		if iface.NetworkAttachment != "" {
			networkAttachment, err := tpgresource.GetRelativePath(iface.NetworkAttachment)
			if err != nil {
				return nil, "", "", err
			}
			flattened[i]["network_attachment"] = networkAttachment
		}

		// the security_policy for a network_interface is found in one of its accessConfigs.
		if len(iface.AccessConfigs) > 0 && iface.AccessConfigs[0].SecurityPolicy != "" {
			flattened[i]["security_policy"] = iface.AccessConfigs[0].SecurityPolicy
		} else if len(iface.Ipv6AccessConfigs) > 0 && iface.Ipv6AccessConfigs[0].SecurityPolicy != "" {
			flattened[i]["security_policy"] = iface.Ipv6AccessConfigs[0].SecurityPolicy
		}
	}
	return flattened, internalIP, externalIP, nil
}

func flattenServiceAccounts(serviceAccounts []*compute.ServiceAccount) []map[string]interface{} {
	result := make([]map[string]interface{}, len(serviceAccounts))
	for i, serviceAccount := range serviceAccounts {
		result[i] = map[string]interface{}{
			"email":  serviceAccount.Email,
			"scopes": serviceAccount.Scopes,
		}
	}
	return result
}

func flattenGuestAccelerators(accelerators []*compute.AcceleratorConfig) []map[string]interface{} {
	acceleratorsSchema := make([]map[string]interface{}, len(accelerators))
	for i, accelerator := range accelerators {
		acceleratorsSchema[i] = map[string]interface{}{
			"count": accelerator.AcceleratorCount,
			"type":  accelerator.AcceleratorType,
		}
	}
	return acceleratorsSchema
}

func flattenConfidentialInstanceConfig(ConfidentialInstanceConfig *compute.ConfidentialInstanceConfig) []map[string]interface{} {
	if ConfidentialInstanceConfig == nil {
		return nil
	}

	return []map[string]interface{}{{
		"enable_confidential_compute": ConfidentialInstanceConfig.EnableConfidentialCompute,
		"confidential_instance_type":  ConfidentialInstanceConfig.ConfidentialInstanceType,
	}}
}

func flattenAdvancedMachineFeatures(AdvancedMachineFeatures *compute.AdvancedMachineFeatures) []map[string]interface{} {
	if AdvancedMachineFeatures == nil {
		return nil
	}
	return []map[string]interface{}{{
		"enable_nested_virtualization": AdvancedMachineFeatures.EnableNestedVirtualization,
		"threads_per_core":             AdvancedMachineFeatures.ThreadsPerCore,
		"turbo_mode":                   AdvancedMachineFeatures.TurboMode,
		"visible_core_count":           AdvancedMachineFeatures.VisibleCoreCount,
		"performance_monitoring_unit":  AdvancedMachineFeatures.PerformanceMonitoringUnit,
		"enable_uefi_networking":       AdvancedMachineFeatures.EnableUefiNetworking,
	}}
}

func flattenShieldedVmConfig(shieldedVmConfig *compute.ShieldedInstanceConfig) []map[string]bool {
	if shieldedVmConfig == nil {
		return nil
	}

	shieldedInstanceConfig := map[string]bool{}

	if shieldedVmConfig.EnableSecureBoot {
		shieldedInstanceConfig["enable_secure_boot"] = shieldedVmConfig.EnableSecureBoot
	}

	if !shieldedVmConfig.EnableVtpm {
		shieldedInstanceConfig["enable_vtpm"] = shieldedVmConfig.EnableVtpm
	}

	if !shieldedVmConfig.EnableIntegrityMonitoring {
		shieldedInstanceConfig["enable_integrity_monitoring"] = shieldedVmConfig.EnableIntegrityMonitoring
	}

	if len(shieldedInstanceConfig) == 0 {
		return nil
	}

	return []map[string]bool{shieldedInstanceConfig}
}

func flattenEnableDisplay(displayDevice *compute.DisplayDevice) interface{} {
	if displayDevice == nil {
		return nil
	}

	return displayDevice.EnableDisplay
}

func flattenReservationAffinity(affinity *compute.ReservationAffinity) []map[string]interface{} {
	if affinity == nil {
		return nil
	}

	flattened := map[string]interface{}{
		"type": affinity.ConsumeReservationType,
	}

	if affinity.ConsumeReservationType == "SPECIFIC_RESERVATION" {
		flattened["specific_reservation"] = []map[string]interface{}{{
			"key":    affinity.Key,
			"values": affinity.Values,
		}}
	}

	return []map[string]interface{}{flattened}
}

func flattenNetworkPerformanceConfig(c *compute.NetworkPerformanceConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"total_egress_bandwidth_tier": c.TotalEgressBandwidthTier,
		},
	}
}
