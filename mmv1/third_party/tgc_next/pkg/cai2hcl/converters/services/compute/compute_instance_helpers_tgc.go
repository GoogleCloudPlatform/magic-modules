package compute

import (
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tpgresource"

	compute "google.golang.org/api/compute/v0.beta"
)

func flattenAliasIpRangeTgc(ranges []*compute.AliasIpRange) []map[string]interface{} {
	rangesSchema := make([]map[string]interface{}, 0, len(ranges))
	for _, ipRange := range ranges {
		rangesSchema = append(rangesSchema, map[string]interface{}{
			"ip_cidr_range":         ipRange.IpCidrRange,
			"subnetwork_range_name": ipRange.SubnetworkRangeName,
		})
	}
	return rangesSchema
}

func flattenSchedulingTgc(resp *compute.Scheduling) []map[string]interface{} {
	schedulingMap := make(map[string]interface{}, 0)

	// gracefulShutdown is not in the cai asset, so graceful_shutdown is skipped.

	if resp.InstanceTerminationAction != "" {
		schedulingMap["instance_termination_action"] = resp.InstanceTerminationAction
	}

	if resp.MinNodeCpus != 0 {
		schedulingMap["min_node_cpus"] = resp.MinNodeCpus
	}

	schedulingMap["on_host_maintenance"] = resp.OnHostMaintenance

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

	schedulingMap["provisioning_model"] = resp.ProvisioningModel

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

func flattenNetworkInterfacesTgc(networkInterfaces []*compute.NetworkInterface, project string) ([]map[string]interface{}, string, string, error) {
	flattened := make([]map[string]interface{}, len(networkInterfaces))
	var internalIP, externalIP string

	for i, iface := range networkInterfaces {
		var ac []map[string]interface{}
		ac, externalIP = flattenAccessConfigs(iface.AccessConfigs)

		flattened[i] = map[string]interface{}{
			"network_ip":                  iface.NetworkIP,
			"access_config":               ac,
			"alias_ip_range":              flattenAliasIpRangeTgc(iface.AliasIpRanges),
			"nic_type":                    iface.NicType,
			"stack_type":                  iface.StackType,
			"ipv6_access_config":          flattenIpv6AccessConfigs(iface.Ipv6AccessConfigs),
			"ipv6_address":                iface.Ipv6Address,
			"network":                     tpgresource.ConvertSelfLinkToV1(iface.Network),
			"subnetwork":                  tpgresource.ConvertSelfLinkToV1(iface.Subnetwork),
			"internal_ipv6_prefix_length": iface.InternalIpv6PrefixLength,
		}

		subnetProject := utils.ParseFieldValue(iface.Subnetwork, "projects")
		if subnetProject != project {
			flattened[i]["subnetwork_project"] = subnetProject
		}

		// The field name is computed, no it is not converted.

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

func flattenServiceAccountsTgc(serviceAccounts []*compute.ServiceAccount) []map[string]interface{} {
	result := make([]map[string]interface{}, len(serviceAccounts))
	for i, serviceAccount := range serviceAccounts {
		scopes := serviceAccount.Scopes
		if len(scopes) == 0 {
			scopes = []string{}
		}
		result[i] = map[string]interface{}{
			"email":  serviceAccount.Email,
			"scopes": scopes,
		}
	}
	return result
}

func flattenGuestAcceleratorsTgc(accelerators []*compute.AcceleratorConfig) []map[string]interface{} {
	acceleratorsSchema := make([]map[string]interface{}, len(accelerators))
	for i, accelerator := range accelerators {
		acceleratorsSchema[i] = map[string]interface{}{
			"count": accelerator.AcceleratorCount,
			"type":  tpgresource.GetResourceNameFromSelfLink(accelerator.AcceleratorType),
		}
	}
	return acceleratorsSchema
}

func flattenReservationAffinityTgc(affinity *compute.ReservationAffinity) []map[string]interface{} {
	if affinity == nil {
		return nil
	}

	// The values of ConsumeReservationType in cai assets are NO_ALLOCATION, SPECIFIC_ALLOCATION, ANY_ALLOCATION
	crt := strings.ReplaceAll(affinity.ConsumeReservationType, "_ALLOCATION", "_RESERVATION")
	flattened := map[string]interface{}{
		"type": crt,
	}

	if crt == "SPECIFIC_RESERVATION" {
		flattened["specific_reservation"] = []map[string]interface{}{{
			"key":    affinity.Key,
			"values": affinity.Values,
		}}
	}

	return []map[string]interface{}{flattened}
}
