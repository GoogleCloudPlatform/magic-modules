package compute

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/api/googleapi"

	compute "google.golang.org/api/compute/v0.beta"
)

func instanceSchedulingNodeAffinitiesElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"IN", "NOT_IN"}, false),
			},
			"values": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func expandAliasIpRanges(ranges []interface{}) []*compute.AliasIpRange {
	ipRanges := make([]*compute.AliasIpRange, 0, len(ranges))
	for _, raw := range ranges {
		data := raw.(map[string]interface{})
		ipRanges = append(ipRanges, &compute.AliasIpRange{
			IpCidrRange:         data["ip_cidr_range"].(string),
			SubnetworkRangeName: data["subnetwork_range_name"].(string),
		})
	}
	return ipRanges
}

func expandScheduling(v interface{}) (*compute.Scheduling, error) {
	if v == nil {
		// We can't set default values for lists.
		return &compute.Scheduling{
			AutomaticRestart: googleapi.Bool(true),
		}, nil
	}

	ls := v.([]interface{})
	if len(ls) == 0 {
		// We can't set default values for lists
		return &compute.Scheduling{
			AutomaticRestart: googleapi.Bool(true),
		}, nil
	}

	if len(ls) > 1 || ls[0] == nil {
		return nil, fmt.Errorf("expected exactly one scheduling block")
	}

	original := ls[0].(map[string]interface{})
	scheduling := &compute.Scheduling{
		ForceSendFields: make([]string, 0, 4),
	}

	if v, ok := original["automatic_restart"]; ok {
		scheduling.AutomaticRestart = googleapi.Bool(v.(bool))
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "AutomaticRestart")
	}

	if v, ok := original["preemptible"]; ok {
		scheduling.Preemptible = v.(bool)
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "Preemptible")
	}

	if v, ok := original["on_host_maintenance"]; ok {
		scheduling.OnHostMaintenance = v.(string)
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "OnHostMaintenance")
	}

	if v, ok := original["node_affinities"]; ok && v != nil {
		naSet := v.(*schema.Set).List()
		scheduling.NodeAffinities = make([]*compute.SchedulingNodeAffinity, len(ls))
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "NodeAffinities")
		for _, nodeAffRaw := range naSet {
			if nodeAffRaw == nil {
				continue
			}
			nodeAff := nodeAffRaw.(map[string]interface{})
			transformed := &compute.SchedulingNodeAffinity{
				Key:      nodeAff["key"].(string),
				Operator: nodeAff["operator"].(string),
				Values:   tpgresource.ConvertStringArr(nodeAff["values"].(*schema.Set).List()),
			}
			scheduling.NodeAffinities = append(scheduling.NodeAffinities, transformed)
		}
	}

	if v, ok := original["min_node_cpus"]; ok {
		scheduling.MinNodeCpus = int64(v.(int))
	}
	if v, ok := original["provisioning_model"]; ok {
		scheduling.ProvisioningModel = v.(string)
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "ProvisioningModel")
	}
	if v, ok := original["instance_termination_action"]; ok {
		scheduling.InstanceTerminationAction = v.(string)
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "InstanceTerminationAction")
	}
	if v, ok := original["availability_domain"]; ok && v != nil {
		scheduling.AvailabilityDomain = int64(v.(int))
	}
	if v, ok := original["max_run_duration"]; ok {
		transformedMaxRunDuration, err := expandComputeMaxRunDuration(v)
		if err != nil {
			return nil, err
		}
		scheduling.MaxRunDuration = transformedMaxRunDuration
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "MaxRunDuration")
	}

	if v, ok := original["on_instance_stop_action"]; ok {
		transformedOnInstanceStopAction, err := expandComputeOnInstanceStopAction(v)
		if err != nil {
			return nil, err
		}
		scheduling.OnInstanceStopAction = transformedOnInstanceStopAction
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "OnInstanceStopAction")
	}

	if v, ok := original["host_error_timeout_seconds"]; ok {
		scheduling.HostErrorTimeoutSeconds = int64(v.(int))
		//host_error_timeout_seconds doesn't get removed correctly due to an API bug on instances.SetScheduling.
		//We need to set it to NullFields as a workaround because nil is rounded to 0
		if v == 0 || v == nil {
			scheduling.NullFields = append(scheduling.NullFields, "HostErrorTimeoutSeconds")
		} else {
			scheduling.ForceSendFields = append(scheduling.ForceSendFields, "HostErrorTimeoutSeconds")
		}
	}

	if v, ok := original["maintenance_interval"]; ok {
		scheduling.MaintenanceInterval = v.(string)
	}

	if v, ok := original["graceful_shutdown"]; ok {
		transformedGracefulShutdown, err := expandGracefulShutdown(v)
		if err != nil {
			return nil, err
		}
		scheduling.GracefulShutdown = transformedGracefulShutdown
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "GracefulShutdown")
	}

	if v, ok := original["local_ssd_recovery_timeout"]; ok {
		transformedLocalSsdRecoveryTimeout, err := expandComputeLocalSsdRecoveryTimeout(v)
		if err != nil {
			return nil, err
		}
		scheduling.LocalSsdRecoveryTimeout = transformedLocalSsdRecoveryTimeout
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "LocalSsdRecoveryTimeout")
	}
	return scheduling, nil
}

func expandComputeMaxRunDuration(v interface{}) (*compute.Duration, error) {
	l := v.([]interface{})
	duration := compute.Duration{}
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})

	transformedNanos, err := expandComputeMaxRunDurationNanos(original["nanos"])
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNanos); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		duration.Nanos = int64(transformedNanos.(int))
	}

	transformedSeconds, err := expandComputeMaxRunDurationSeconds(original["seconds"])
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSeconds); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		duration.Seconds = int64(transformedSeconds.(int))
	}

	return &duration, nil
}

func expandComputeMaxRunDurationNanos(v interface{}) (interface{}, error) {
	return v, nil
}

func expandComputeMaxRunDurationSeconds(v interface{}) (interface{}, error) {
	return v, nil
}

func expandComputeOnInstanceStopAction(v interface{}) (*compute.SchedulingOnInstanceStopAction, error) {
	l := v.([]interface{})
	onInstanceStopAction := compute.SchedulingOnInstanceStopAction{}
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})

	if d, ok := original["discard_local_ssd"]; ok {
		onInstanceStopAction.DiscardLocalSsd = d.(bool)
	} else {
		return nil, nil
	}

	return &onInstanceStopAction, nil
}

func expandComputeLocalSsdRecoveryTimeout(v interface{}) (*compute.Duration, error) {
	l := v.([]interface{})
	duration := compute.Duration{}
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})

	transformedNanos, err := expandComputeLocalSsdRecoveryTimeoutNanos(original["nanos"])
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNanos); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		duration.Nanos = int64(transformedNanos.(int))
	}

	transformedSeconds, err := expandComputeLocalSsdRecoveryTimeoutSeconds(original["seconds"])
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSeconds); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		duration.Seconds = int64(transformedSeconds.(int))
	}
	return &duration, nil
}

func expandComputeLocalSsdRecoveryTimeoutNanos(v interface{}) (interface{}, error) {
	return v, nil
}

func expandComputeLocalSsdRecoveryTimeoutSeconds(v interface{}) (interface{}, error) {
	return v, nil
}

func expandGracefulShutdown(v interface{}) (*compute.SchedulingGracefulShutdown, error) {
	l := v.([]interface{})
	gracefulShutdown := compute.SchedulingGracefulShutdown{}
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})

	originalMaxDuration := original["max_duration"].([]interface{})
	maxDuration, err := expandGracefulShutdownMaxDuration(originalMaxDuration)
	if err != nil {
		return nil, err
	}
	if maxDuration != nil {
		gracefulShutdown.MaxDuration = maxDuration
	}

	gracefulShutdown.Enabled = original["enabled"].(bool)
	gracefulShutdown.ForceSendFields = append(gracefulShutdown.ForceSendFields, "Enabled")
	return &gracefulShutdown, nil
}

func expandGracefulShutdownMaxDuration(v interface{}) (*compute.Duration, error) {
	l := v.([]interface{})
	duration := compute.Duration{}
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]

	maxDurationMap := raw.(map[string]interface{})
	transformedNanos := maxDurationMap["nanos"]
	transformedSeconds := maxDurationMap["seconds"]

	if val := reflect.ValueOf(transformedNanos); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		duration.Nanos = int64(transformedNanos.(int))
	}
	if val := reflect.ValueOf(transformedSeconds); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		duration.Seconds = int64(transformedSeconds.(int))
	}

	duration.ForceSendFields = append(duration.ForceSendFields, "Seconds")

	return &duration, nil
}

func expandAccessConfigs(configs []interface{}) []*compute.AccessConfig {
	acs := make([]*compute.AccessConfig, len(configs))
	for i, raw := range configs {
		acs[i] = &compute.AccessConfig{}
		acs[i].Type = "ONE_TO_ONE_NAT"
		if raw != nil {
			data := raw.(map[string]interface{})
			acs[i].NatIP = data["nat_ip"].(string)
			acs[i].NetworkTier = data["network_tier"].(string)
			if ptr, ok := data["public_ptr_domain_name"]; ok && ptr != "" {
				acs[i].SetPublicPtr = true
				acs[i].PublicPtrDomainName = ptr.(string)
			}
		}
	}
	return acs
}

func expandIpv6AccessConfigs(configs []interface{}) []*compute.AccessConfig {
	iacs := make([]*compute.AccessConfig, len(configs))
	for i, raw := range configs {
		iacs[i] = &compute.AccessConfig{}
		if raw != nil {
			data := raw.(map[string]interface{})
			iacs[i].NetworkTier = data["network_tier"].(string)
			if ptr, ok := data["public_ptr_domain_name"]; ok && ptr != "" {
				iacs[i].PublicPtrDomainName = ptr.(string)
			}
			if eip, ok := data["external_ipv6"]; ok && eip != "" {
				iacs[i].ExternalIpv6 = eip.(string)
			}
			if eipl, ok := data["external_ipv6_prefix_length"]; ok && eipl != "" {
				if strVal, ok := eipl.(string); ok {
					if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
						iacs[i].ExternalIpv6PrefixLength = intVal
					}
				}
			}
			if name, ok := data["name"]; ok && name != "" {
				iacs[i].Name = name.(string)
			}
			iacs[i].Type = "DIRECT_IPV6" // Currently only type supported
		}
	}
	return iacs
}

func expandNetworkInterfaces(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]*compute.NetworkInterface, error) {
	configs := d.Get("network_interface").([]interface{})
	ifaces := make([]*compute.NetworkInterface, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		var networkAttachment = ""
		network := data["network"].(string)
		subnetwork := data["subnetwork"].(string)
		if networkAttachmentObj, ok := data["network_attachment"]; ok {
			networkAttachment = networkAttachmentObj.(string)
		}
		// Checks if networkAttachment is not specified in resource, network or subnetwork have to be specified.
		if networkAttachment == "" && network == "" && subnetwork == "" {
			return nil, fmt.Errorf("exactly one of network, subnetwork, or network_attachment must be provided")
		}

		nf, err := tpgresource.ParseNetworkFieldValue(network, d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for network %q: %s", network, err)
		}

		subnetProjectField := fmt.Sprintf("network_interface.%d.subnetwork_project", i)
		sf, err := tpgresource.ParseSubnetworkFieldValueWithProjectField(subnetwork, subnetProjectField, d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for subnetwork %q: %s", subnetwork, err)
		}

		ifaces[i] = &compute.NetworkInterface{
			NetworkIP:         data["network_ip"].(string),
			Network:           nf.RelativeLink(),
			Subnetwork:        sf.RelativeLink(),
			AccessConfigs:     expandAccessConfigs(data["access_config"].([]interface{})),
			AliasIpRanges:     expandAliasIpRanges(data["alias_ip_range"].([]interface{})),
			NicType:           data["nic_type"].(string),
			StackType:         data["stack_type"].(string),
			QueueCount:        int64(data["queue_count"].(int)),
			Ipv6AccessConfigs: expandIpv6AccessConfigs(data["ipv6_access_config"].([]interface{})),
			Ipv6Address:       data["ipv6_address"].(string),
		}
	}
	return ifaces, nil
}

func expandServiceAccounts(configs []interface{}) []*compute.ServiceAccount {
	accounts := make([]*compute.ServiceAccount, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		accounts[i] = &compute.ServiceAccount{
			Email:  data["email"].(string),
			Scopes: tpgresource.CanonicalizeServiceScopes(tpgresource.ConvertStringSet(data["scopes"].(*schema.Set))),
		}

		if accounts[i].Email == "" {
			accounts[i].Email = "default"
		}
	}
	return accounts
}

func resourceInstanceTags(d tpgresource.TerraformResourceData) *compute.Tags {
	// Calculate the tags
	var tags *compute.Tags
	if v := d.Get("tags"); v != nil {
		vs := v.(*schema.Set)
		tags = new(compute.Tags)
		tags.Items = make([]string, vs.Len())
		for i, v := range vs.List() {
			tags.Items[i] = v.(string)
		}

		tags.Fingerprint = d.Get("tags_fingerprint").(string)
	}

	return tags
}

func expandShieldedVmConfigs(d tpgresource.TerraformResourceData) *compute.ShieldedInstanceConfig {
	if _, ok := d.GetOk("shielded_instance_config"); !ok {
		return nil
	}

	prefix := "shielded_instance_config.0"
	return &compute.ShieldedInstanceConfig{
		EnableSecureBoot:          d.Get(prefix + ".enable_secure_boot").(bool),
		EnableVtpm:                d.Get(prefix + ".enable_vtpm").(bool),
		EnableIntegrityMonitoring: d.Get(prefix + ".enable_integrity_monitoring").(bool),
		ForceSendFields:           []string{"EnableSecureBoot", "EnableVtpm", "EnableIntegrityMonitoring"},
	}
}

func expandConfidentialInstanceConfig(d tpgresource.TerraformResourceData) *compute.ConfidentialInstanceConfig {
	if _, ok := d.GetOk("confidential_instance_config"); !ok {
		return nil
	}

	prefix := "confidential_instance_config.0"
	return &compute.ConfidentialInstanceConfig{
		EnableConfidentialCompute: d.Get(prefix + ".enable_confidential_compute").(bool),
		ConfidentialInstanceType:  d.Get(prefix + ".confidential_instance_type").(string),
	}
}

func expandAdvancedMachineFeatures(d tpgresource.TerraformResourceData) *compute.AdvancedMachineFeatures {
	if _, ok := d.GetOk("advanced_machine_features"); !ok {
		return nil
	}

	prefix := "advanced_machine_features.0"
	return &compute.AdvancedMachineFeatures{
		EnableNestedVirtualization: d.Get(prefix + ".enable_nested_virtualization").(bool),
		ThreadsPerCore:             int64(d.Get(prefix + ".threads_per_core").(int)),
		TurboMode:                  d.Get(prefix + ".turbo_mode").(string),
		VisibleCoreCount:           int64(d.Get(prefix + ".visible_core_count").(int)),
		PerformanceMonitoringUnit:  d.Get(prefix + ".performance_monitoring_unit").(string),
		EnableUefiNetworking:       d.Get(prefix + ".enable_uefi_networking").(bool),
	}
}

func expandDisplayDevice(d tpgresource.TerraformResourceData) *compute.DisplayDevice {
	if _, ok := d.GetOk("enable_display"); !ok {
		return nil
	}
	return &compute.DisplayDevice{
		EnableDisplay:   d.Get("enable_display").(bool),
		ForceSendFields: []string{"EnableDisplay"},
	}
}

// Node affinity updates require a reboot
func schedulingHasChangeRequiringReboot(d *schema.ResourceData) bool {
	o, n := d.GetChange("scheduling")
	oScheduling := o.([]interface{})[0].(map[string]interface{})
	newScheduling := n.([]interface{})[0].(map[string]interface{})
	return hasNodeAffinitiesChanged(oScheduling, newScheduling) ||
		hasMaxRunDurationChanged(oScheduling, newScheduling) ||
		hasGracefulShutdownChangedWithReboot(d, oScheduling, newScheduling)
}

// Terraform doesn't correctly calculate changes on schema.Set, so we do it manually
// https://github.com/hashicorp/terraform-plugin-sdk/issues/98
func schedulingHasChangeWithoutReboot(d *schema.ResourceData) bool {
	if !d.HasChange("scheduling") {
		// This doesn't work correctly, which is why this method exists
		// But it is here for posterity
		return false
	}
	o, n := d.GetChange("scheduling")
	oScheduling := o.([]interface{})[0].(map[string]interface{})
	newScheduling := n.([]interface{})[0].(map[string]interface{})

	if schedulingHasChangeRequiringReboot(d) {
		return false
	}

	if oScheduling["automatic_restart"] != newScheduling["automatic_restart"] {
		return true
	}

	if oScheduling["preemptible"] != newScheduling["preemptible"] {
		return true
	}

	if oScheduling["on_host_maintenance"] != newScheduling["on_host_maintenance"] {
		return true
	}

	if oScheduling["provisioning_model"] != newScheduling["provisioning_model"] {
		return true
	}

	if oScheduling["instance_termination_action"] != newScheduling["instance_termination_action"] {
		return true
	}
	if oScheduling["availability_domain"] != newScheduling["availability_domain"] {
		return true
	}

	if oScheduling["host_error_timeout_seconds"] != newScheduling["host_error_timeout_seconds"] {
		return true
	}

	if hasGracefulShutdownChanged(oScheduling, newScheduling) {
		return true
	}

	return false
}

func hasGracefulShutdownChangedWithReboot(d *schema.ResourceData, oScheduling, nScheduling map[string]interface{}) bool {
	allow_stopping_for_update := d.Get("allow_stopping_for_update").(bool)
	if !allow_stopping_for_update {
		return false
	}
	return hasGracefulShutdownChanged(oScheduling, nScheduling)
}

func hasGracefulShutdownChanged(oScheduling, nScheduling map[string]interface{}) bool {
	oGrShut := oScheduling["graceful_shutdown"].([]interface{})
	nGrShut := nScheduling["graceful_shutdown"].([]interface{})

	if (len(oGrShut) == 0 || oGrShut[0] == nil) && (len(nGrShut) == 0 || nGrShut[0] == nil) {
		return false
	}
	if (len(oGrShut) == 0 || oGrShut[0] == nil) || (len(nGrShut) == 0 || nGrShut[0] == nil) {
		return true
	}

	oldGrShut := oGrShut[0].(map[string]interface{})
	newGrShut := nGrShut[0].(map[string]interface{})
	oldMaxDuration := oldGrShut["max_duration"].([]interface{})
	newMaxDuration := newGrShut["max_duration"].([]interface{})
	var oldMaxDurationMap map[string]interface{}
	var newMaxDurationMap map[string]interface{}

	if len(oldMaxDuration) > 0 && oldMaxDuration[0] != nil {
		oldMaxDurationMap = oldMaxDuration[0].(map[string]interface{})
	} else {
		oldMaxDurationMap = nil
	}

	if len(newMaxDuration) > 0 && newMaxDuration[0] != nil {
		newMaxDurationMap = newMaxDuration[0].(map[string]interface{})
	} else {
		newMaxDurationMap = nil
	}

	if oldGrShut["enabled"] != newGrShut["enabled"] {
		return true
	}
	if oldMaxDurationMap["seconds"] != newMaxDurationMap["seconds"] {
		return true
	}
	if oldMaxDurationMap["nanos"] != newMaxDurationMap["nanos"] {
		return true
	}

	return false
}

func hasMaxRunDurationChanged(oScheduling, nScheduling map[string]interface{}) bool {
	oMrd := oScheduling["max_run_duration"].([]interface{})
	nMrd := nScheduling["max_run_duration"].([]interface{})

	if (len(oMrd) == 0 || oMrd[0] == nil) && (len(nMrd) == 0 || nMrd[0] == nil) {
		return false
	}
	if (len(oMrd) == 0 || oMrd[0] == nil) || (len(nMrd) == 0 || nMrd[0] == nil) {
		return true
	}

	oldMrd := oMrd[0].(map[string]interface{})
	newMrd := nMrd[0].(map[string]interface{})

	if oldMrd["seconds"] != newMrd["seconds"] {
		return true
	}
	if oldMrd["nanos"] != newMrd["nanos"] {
		return true
	}

	return false
}

func hasNodeAffinitiesChanged(oScheduling, newScheduling map[string]interface{}) bool {
	oldNAs := oScheduling["node_affinities"].(*schema.Set).List()
	newNAs := newScheduling["node_affinities"].(*schema.Set).List()
	if len(oldNAs) != len(newNAs) {
		return true
	}
	for i := range oldNAs {
		oldNodeAffinity := oldNAs[i].(map[string]interface{})
		newNodeAffinity := newNAs[i].(map[string]interface{})
		if oldNodeAffinity["key"] != newNodeAffinity["key"] {
			return true
		}
		if oldNodeAffinity["operator"] != newNodeAffinity["operator"] {
			return true
		}

		// ConvertStringSet will sort the set into a slice, allowing DeepEqual
		if !reflect.DeepEqual(tpgresource.ConvertStringSet(oldNodeAffinity["values"].(*schema.Set)), tpgresource.ConvertStringSet(newNodeAffinity["values"].(*schema.Set))) {
			return true
		}
	}

	return false
}

func expandReservationAffinity(d tpgresource.TerraformResourceData) (*compute.ReservationAffinity, error) {
	_, ok := d.GetOk("reservation_affinity")
	if !ok {
		return nil, nil
	}

	prefix := "reservation_affinity.0"
	reservationAffinityType := d.Get(prefix + ".type").(string)

	affinity := compute.ReservationAffinity{
		ConsumeReservationType: reservationAffinityType,
		ForceSendFields:        []string{"ConsumeReservationType"},
	}

	_, hasSpecificReservation := d.GetOk(prefix + ".specific_reservation")
	if (reservationAffinityType == "SPECIFIC_RESERVATION") != hasSpecificReservation {
		return nil, fmt.Errorf("specific_reservation must be set when reservation_affinity is SPECIFIC_RESERVATION, and not set otherwise")
	}

	prefix = prefix + ".specific_reservation.0"
	if hasSpecificReservation {
		affinity.Key = d.Get(prefix + ".key").(string)
		affinity.ForceSendFields = append(affinity.ForceSendFields, "Key", "Values")

		for _, v := range d.Get(prefix + ".values").([]interface{}) {
			affinity.Values = append(affinity.Values, v.(string))
		}
	}

	return &affinity, nil
}

func expandNetworkPerformanceConfig(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*compute.NetworkPerformanceConfig, error) {
	configs, ok := d.GetOk("network_performance_config")
	if !ok {
		return nil, nil
	}

	npcSlice := configs.([]interface{})
	if len(npcSlice) > 1 {
		return nil, fmt.Errorf("cannot specify multiple network_performance_configs")
	}

	if len(npcSlice) == 0 || npcSlice[0] == nil {
		return nil, nil
	}
	npc := npcSlice[0].(map[string]interface{})
	return &compute.NetworkPerformanceConfig{
		TotalEgressBandwidthTier: npc["total_egress_bandwidth_tier"].(string),
	}, nil
}
