package compute

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tgcresource"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tpgresource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ComputeInstanceCai2hclConverter for compute instance resource.
type ComputeInstanceCai2hclConverter struct {
	name   string
	schema map[string]*schema.Schema
}

// NewComputeInstanceCai2hclConverter returns an HCL converter for compute instance.
func NewComputeInstanceCai2hclConverter(provider *schema.Provider) models.Cai2hclConverter {
	schema := provider.ResourcesMap[ComputeInstanceSchemaName].Schema

	return &ComputeInstanceCai2hclConverter{
		name:   ComputeInstanceSchemaName,
		schema: schema,
	}
}

// Convert converts asset to HCL resource blocks.
func (c *ComputeInstanceCai2hclConverter) Convert(asset caiasset.Asset) ([]*models.TerraformResourceBlock, error) {
	var blocks []*models.TerraformResourceBlock
	block, err := c.convertResourceData(asset)
	if err != nil {
		return nil, err
	}
	blocks = append(blocks, block)
	return blocks, nil
}

func (c *ComputeInstanceCai2hclConverter) convertResourceData(asset caiasset.Asset) (*models.TerraformResourceBlock, error) {
	if asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}

	project := tgcresource.ParseFieldValue(asset.Name, "projects")
	data := asset.Resource.Data

	hclData := make(map[string]interface{})

	if canIpForward, ok := data["canIpForward"].(bool); ok && canIpForward {
		hclData["can_ip_forward"] = true
	}
	if machineType, ok := data["machineType"].(string); ok && machineType != "" {
		hclData["machine_type"] = tpgresource.GetResourceNameFromSelfLink(machineType)
	}
	hclData["network_performance_config"] = flattenNetworkPerformanceConfigTgcNext(data["networkPerformanceConfig"])

	// Set the networks
	networkInterfaces, _, _, err := flattenNetworkInterfacesTgcNext(data["networkInterfaces"], project)
	if err != nil {
		return nil, err
	}
	hclData["network_interface"] = networkInterfaces

	if tags, ok := data["tags"].(map[string]interface{}); ok && tags != nil {
		if items, ok := tags["items"].([]interface{}); ok {
			hclData["tags"] = items
		}
	}

	md := flattenMetadataBetaTgcNext(data["metadata"])
	if startupScript, ok := md["startup-script"].(string); ok && startupScript != "" {
		hclData["metadata_startup_script"] = startupScript
		delete(md, "startup-script")
	}

	hclData["service_account"] = flattenServiceAccountsTgcNext(data["serviceAccounts"])
	hclData["resource_policies"] = data["resourcePolicies"]

	instanceName, _ := data["name"].(string)
	bootDisk, ads, scratchDisks := flattenDisksTgcNext(data["disks"], instanceName)
	hclData["boot_disk"] = bootDisk
	hclData["attached_disk"] = ads
	hclData["scratch_disk"] = scratchDisks

	hclData["scheduling"] = flattenSchedulingTgcNext(data["scheduling"])
	hclData["guest_accelerator"] = flattenGuestAcceleratorsTgcNext(data["guestAccelerators"])
	hclData["shielded_instance_config"] = flattenShieldedVmConfigTgcNext(data["shieldedInstanceConfig"])
	hclData["enable_display"] = flattenEnableDisplayTgcNext(data["displayDevice"])
	hclData["min_cpu_platform"] = data["minCpuPlatform"]

	// Only convert the field when its value is not default false
	if deletionProtection, ok := data["deletionProtection"].(bool); ok && deletionProtection {
		hclData["deletion_protection"] = true
	}

	if project != "" {
		hclData["project"] = project
	}
	if zone, ok := data["zone"].(string); ok && zone != "" {
		hclData["zone"] = tpgresource.GetResourceNameFromSelfLink(zone)
	}
	hclData["name"] = data["name"]
	hclData["description"] = data["description"]
	hclData["hostname"] = data["hostname"]
	hclData["confidential_instance_config"] = flattenConfidentialInstanceConfigTgcNext(data["confidentialInstanceConfig"])
	hclData["advanced_machine_features"] = flattenAdvancedMachineFeaturesTgcNext(data["advancedMachineFeatures"])
	hclData["reservation_affinity"] = flattenReservationAffinityTgcNext(data["reservationAffinity"])
	if keyRevocationActionType, ok := data["keyRevocationActionType"].(string); ok {
		hclData["key_revocation_action_type"] = strings.TrimSuffix(keyRevocationActionType, "_ON_KEY_REVOCATION")
	}
	hclData["instance_encryption_key"] = flattenComputeInstanceEncryptionKeyTgcNext(data["instanceEncryptionKey"])

	ctyVal, err := utils.MapToCtyValWithSchema(hclData, c.schema)
	if err != nil {
		return nil, err
	}
	return &models.TerraformResourceBlock{
		Labels: []string{c.name, instanceName},
		Value:  ctyVal,
	}, nil
}

func flattenNetworkPerformanceConfigTgcNext(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	resp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return []map[string]interface{}{{
		"total_egress_bandwidth_tier": resp["totalEgressBandwidthTier"],
	}}
}

func flattenNetworkInterfacesTgcNext(v interface{}, project string) ([]map[string]interface{}, string, string, error) {
	if v == nil {
		return nil, "", "", nil
	}
	networkInterfaces, ok := v.([]interface{})
	if !ok {
		return nil, "", "", nil
	}

	flattened := make([]map[string]interface{}, len(networkInterfaces))
	var internalIP, externalIP string

	for i, rawIface := range networkInterfaces {
		iface, ok := rawIface.(map[string]interface{})
		if !ok {
			continue
		}

		var ac []map[string]interface{}
		ac, externalIP = flattenAccessConfigsTgcNext(iface["accessConfigs"])

		flattened[i] = map[string]interface{}{
			"network_ip":                  iface["networkIP"],
			"access_config":               ac,
			"alias_ip_range":              flattenAliasIpRangeTgcNext(iface["aliasIpRanges"]),
			"nic_type":                    iface["nicType"],
			"stack_type":                  iface["stackType"],
			"igmp_query":                  iface["igmpQuery"],
			"ipv6_access_config":          flattenIpv6AccessConfigsTgcNext(iface["ipv6AccessConfigs"]),
			"ipv6_address":                iface["ipv6Address"],
			"internal_ipv6_prefix_length": iface["internalIpv6PrefixLength"],
		}

		if vlan, ok := iface["vlan"].(float64); ok && vlan != 0 {
			flattened[i]["vlan"] = int(vlan)
		}
		if networkAttachment, ok := iface["networkAttachment"].(string); ok && networkAttachment != "" {
			flattened[i]["network_attachment"] = networkAttachment
		}

		if network, ok := iface["network"].(string); ok && network != "" {
			flattened[i]["network"] = tpgresource.ConvertSelfLinkToV1(network)
		}
		if subnetwork, ok := iface["subnetwork"].(string); ok && subnetwork != "" {
			flattened[i]["subnetwork"] = tpgresource.ConvertSelfLinkToV1(subnetwork)
			subnetProject := tgcresource.ParseFieldValue(subnetwork, "projects")
			if subnetProject != project {
				flattened[i]["subnetwork_project"] = subnetProject
			}
		}

		if stackType, ok := iface["stackType"].(string); ok && stackType != "IPV4_ONLY" {
			flattened[i]["stack_type"] = stackType
		}

		if queueCount, ok := iface["queueCount"].(float64); ok && queueCount != 0 {
			flattened[i]["queue_count"] = int(queueCount)
		}

		if internalIP == "" {
			if ip, ok := iface["networkIP"].(string); ok {
				internalIP = ip
			}
		}

		if accessConfigs, ok := iface["accessConfigs"].([]interface{}); ok && len(accessConfigs) > 0 {
			if firstAc, ok := accessConfigs[0].(map[string]interface{}); ok {
				if sp, ok := firstAc["securityPolicy"].(string); ok && sp != "" {
					flattened[i]["security_policy"] = sp
				}
			}
		} else if ipv6AccessConfigs, ok := iface["ipv6AccessConfigs"].([]interface{}); ok && len(ipv6AccessConfigs) > 0 {
			if firstAc, ok := ipv6AccessConfigs[0].(map[string]interface{}); ok {
				if sp, ok := firstAc["securityPolicy"].(string); ok && sp != "" {
					flattened[i]["security_policy"] = sp
				}
			}
		}
	}
	return flattened, internalIP, externalIP, nil
}

func flattenAccessConfigsTgcNext(v interface{}) ([]map[string]interface{}, string) {
	if v == nil {
		return nil, ""
	}
	accessConfigs, ok := v.([]interface{})
	if !ok {
		return nil, ""
	}
	flattened := make([]map[string]interface{}, len(accessConfigs))
	natIP := ""
	for i, rawAc := range accessConfigs {
		ac, ok := rawAc.(map[string]interface{})
		if !ok {
			continue
		}
		flattened[i] = map[string]interface{}{
			"nat_ip":       ac["natIP"],
			"network_tier": ac["networkTier"],
		}
		if setPublicPtr, ok := ac["setPublicPtr"].(bool); ok && setPublicPtr {
			flattened[i]["public_ptr_domain_name"] = ac["publicPtrDomainName"]
		}
		if natIP == "" {
			if ip, ok := ac["natIP"].(string); ok {
				natIP = ip
			}
		}
	}
	return flattened, natIP
}

func flattenAliasIpRangeTgcNext(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	ranges, ok := v.([]interface{})
	if !ok {
		return nil
	}
	rangesSchema := make([]map[string]interface{}, 0, len(ranges))
	for _, rawRange := range ranges {
		ipRange, ok := rawRange.(map[string]interface{})
		if !ok {
			continue
		}
		rangesSchema = append(rangesSchema, map[string]interface{}{
			"ip_cidr_range":         ipRange["ipCidrRange"],
			"subnetwork_range_name": ipRange["subnetworkRangeName"],
		})
	}
	return rangesSchema
}

func flattenIpv6AccessConfigsTgcNext(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	ipv6AccessConfigs, ok := v.([]interface{})
	if !ok {
		return nil
	}
	flattened := make([]map[string]interface{}, len(ipv6AccessConfigs))
	for i, rawAc := range ipv6AccessConfigs {
		ac, ok := rawAc.(map[string]interface{})
		if !ok {
			continue
		}
		flattened[i] = map[string]interface{}{
			"network_tier": ac["networkTier"],
			"name":         ac["name"],
		}
		if publicPtr, ok := ac["publicPtrDomainName"].(string); ok && publicPtr != "" {
			flattened[i]["public_ptr_domain_name"] = publicPtr
		}
		flattened[i]["external_ipv6"] = ac["externalIpv6"]
		flattened[i]["external_ipv6_prefix_length"] = ac["externalIpv6PrefixLength"]
	}
	return flattened
}

func flattenServiceAccountsTgcNext(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	serviceAccounts, ok := v.([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, len(serviceAccounts))
	for i, rawSa := range serviceAccounts {
		serviceAccount, ok := rawSa.(map[string]interface{})
		if !ok {
			continue
		}
		scopes := serviceAccount["scopes"]
		if scopes == nil {
			scopes = []interface{}{}
		}
		result[i] = map[string]interface{}{
			"email":  serviceAccount["email"],
			"scopes": scopes,
		}
	}
	return result
}

func flattenDisksTgcNext(v interface{}, instanceName string) ([]map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {
	if v == nil {
		return nil, nil, nil
	}
	disks, ok := v.([]interface{})
	if !ok {
		return nil, nil, nil
	}

	attachedDisks := []map[string]interface{}{}
	bootDisk := []map[string]interface{}{}
	scratchDisks := []map[string]interface{}{}
	for _, rawDisk := range disks {
		disk, ok := rawDisk.(map[string]interface{})
		if !ok {
			continue
		}
		isBoot, _ := disk["boot"].(bool)
		diskType, _ := disk["type"].(string)
		if isBoot {
			bootDisk = flattenBootDiskTgcNext(disk, instanceName)
		} else if diskType == "SCRATCH" {
			scratchDisks = append(scratchDisks, flattenScratchDiskTgcNext(disk))
		} else {
			di := map[string]interface{}{
				"device_name": disk["deviceName"],
				"mode":        disk["mode"],
			}
			if source, ok := disk["source"].(string); ok && source != "" {
				di["source"] = tpgresource.ConvertSelfLinkToV1(source)
			}
			if key, ok := disk["diskEncryptionKey"].(map[string]interface{}); ok && key != nil {
				if kmsKeyName, ok := key["kmsKeyName"].(string); ok && kmsKeyName != "" {
					di["kms_key_self_link"] = strings.Split(kmsKeyName, "/cryptoKeyVersions")[0]
				}
			}
			attachedDisks = append(attachedDisks, di)
		}
	}

	ads := []map[string]interface{}{}
	for _, d := range attachedDisks {
		if d != nil {
			ads = append(ads, d)
		}
	}
	return bootDisk, ads, scratchDisks
}

func flattenBootDiskTgcNext(v interface{}, instanceName string) []map[string]interface{} {
	if v == nil {
		return nil
	}
	disk, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	result := map[string]interface{}{}

	if autoDelete, ok := disk["autoDelete"].(bool); ok && !autoDelete {
		result["auto_delete"] = false
	}

	if deviceName, ok := disk["deviceName"].(string); ok && !strings.Contains(deviceName, "persistent-disk-") {
		result["device_name"] = deviceName
	}

	if mode, ok := disk["mode"].(string); ok && mode != "READ_WRITE" {
		result["mode"] = mode
	}

	if key, ok := disk["diskEncryptionKey"].(map[string]interface{}); ok && key != nil {
		if kmsKeyName, ok := key["kmsKeyName"].(string); ok && kmsKeyName != "" {
			result["kms_key_self_link"] = strings.Split(kmsKeyName, "/cryptoKeyVersions")[0]
		}
		if rsaEncryptedKey, ok := key["rsaEncryptedKey"].(string); ok && rsaEncryptedKey != "" {
			result["disk_encryption_key_rsa"] = rsaEncryptedKey
		}
		if rawKey, ok := key["rawKey"].(string); ok && rawKey != "" {
			result["disk_encryption_key_raw"] = rawKey
		}
	}

	if source, ok := disk["source"].(string); ok && source != "" {
		result["source"] = tpgresource.ConvertSelfLinkToV1(source)
	}
	result["guest_os_features"] = flattenComputeInstanceGuestOsFeaturesTgcNext(disk["guestOsFeatures"])

	// The interface property was missing in boot disk mapping, leading to integration test failure for nvme/scsi interface options.
	if diskInterface, ok := disk["interface"].(string); ok && diskInterface != "" {
		result["interface"] = diskInterface
	}

	if len(result) == 0 {
		return nil
	}

	return []map[string]interface{}{result}
}

func flattenComputeInstanceGuestOsFeaturesTgcNext(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	features, ok := v.([]interface{})
	if !ok {
		return nil
	}
	var result []interface{}
	for _, rawFeature := range features {
		feature, ok := rawFeature.(map[string]interface{})
		if !ok {
			continue
		}
		if t, ok := feature["type"].(string); ok && t != "" {
			result = append(result, t)
		}
	}
	return result
}

func flattenScratchDiskTgcNext(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	disk, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	result := map[string]interface{}{
		"size": disk["diskSizeGb"],
	}

	if deviceName, ok := disk["deviceName"].(string); ok && !strings.Contains(deviceName, "persistent-disk-") {
		result["device_name"] = deviceName
	}

	result["interface"] = disk["interface"]

	return result
}

func flattenSchedulingTgcNext(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	resp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	schedulingMap := make(map[string]interface{}, 0)

	if terminationTime, ok := resp["terminationTime"].(string); ok && terminationTime != "" {
		schedulingMap["termination_time"] = terminationTime
	}

	if ita, ok := resp["instanceTerminationAction"].(string); ok && ita != "" {
		schedulingMap["instance_termination_action"] = ita
	}

	if minNodeCpus, ok := resp["minNodeCpus"].(float64); ok && minNodeCpus != 0 {
		schedulingMap["min_node_cpus"] = int(minNodeCpus)
	}

	if ohm, ok := resp["onHostMaintenance"].(string); ok && ohm != "" {
		schedulingMap["on_host_maintenance"] = ohm
	}

	if ar, ok := resp["automaticRestart"].(bool); ok && !ar {
		schedulingMap["automatic_restart"] = false
	}

	if preemptible, ok := resp["preemptible"].(bool); ok && preemptible {
		schedulingMap["preemptible"] = true
	}

	if na, ok := resp["nodeAffinities"].([]interface{}); ok && len(na) > 0 {
		nodeAffinities := []map[string]interface{}{}
		for _, rawNa := range na {
			item, ok := rawNa.(map[string]interface{})
			if !ok {
				continue
			}
			nodeAffinities = append(nodeAffinities, map[string]interface{}{
				"key":      item["key"],
				"operator": item["operator"],
				"values":   tpgresource.ConvertStringArrToInterface(convertToStringSliceTgcNext(item["values"])),
			})
		}
		schedulingMap["node_affinities"] = nodeAffinities
	}

	if pm, ok := resp["provisioningModel"].(string); ok && pm != "" {
		schedulingMap["provisioning_model"] = pm
	}

	if ad, ok := resp["availabilityDomain"].(float64); ok && ad != 0 {
		schedulingMap["availability_domain"] = int(ad)
	}

	if mrd, ok := resp["maxRunDuration"].(map[string]interface{}); ok && mrd != nil {
		schedulingMap["max_run_duration"] = flattenComputeMaxRunDurationTgcNext(mrd)
	}

	if oisa, ok := resp["onInstanceStopAction"].(map[string]interface{}); ok && oisa != nil {
		schedulingMap["on_instance_stop_action"] = flattenOnInstanceStopActionTgcNext(oisa)
	}

	if lsrt, ok := resp["localSsdRecoveryTimeout"].(map[string]interface{}); ok && lsrt != nil {
		schedulingMap["local_ssd_recovery_timeout"] = flattenComputeLocalSsdRecoveryTimeoutTgcNext(lsrt)
	}

	if len(schedulingMap) == 0 {
		return nil
	}

	return []map[string]interface{}{schedulingMap}
}

func convertToStringSliceTgcNext(v interface{}) []string {
	if v == nil {
		return nil
	}
	rawList, ok := v.([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, len(rawList))
	for i, val := range rawList {
		if s, ok := val.(string); ok {
			result[i] = s
		}
	}
	return result
}

func flattenComputeMaxRunDurationTgcNext(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	mrd, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"seconds": mrd["seconds"],
			"nanos":   mrd["nanos"],
		},
	}
}

func flattenOnInstanceStopActionTgcNext(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	oisa, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"discard_local_ssd": oisa["discardLocalSsd"],
		},
	}
}

func flattenComputeLocalSsdRecoveryTimeoutTgcNext(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	lsrt, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"seconds": lsrt["seconds"],
			"nanos":   lsrt["nanos"],
		},
	}
}

func flattenGuestAcceleratorsTgcNext(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	accelerators, ok := v.([]interface{})
	if !ok {
		return nil
	}
	acceleratorsSchema := make([]map[string]interface{}, len(accelerators))
	for i, rawAccelerator := range accelerators {
		accelerator, ok := rawAccelerator.(map[string]interface{})
		if !ok {
			continue
		}
		var acceleratorType string
		if at, ok := accelerator["acceleratorType"].(string); ok {
			acceleratorType = tpgresource.GetResourceNameFromSelfLink(at)
		}
		acceleratorsSchema[i] = map[string]interface{}{
			"count": accelerator["acceleratorCount"],
			"type":  acceleratorType,
		}
	}
	return acceleratorsSchema
}

func flattenShieldedVmConfigTgcNext(v interface{}) []map[string]bool {
	if v == nil {
		return nil
	}
	shieldedVmConfig, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	result := map[string]bool{}
	if secureBoot, ok := shieldedVmConfig["enableSecureBoot"].(bool); ok {
		result["enable_secure_boot"] = secureBoot
	}
	if vtpm, ok := shieldedVmConfig["enableVtpm"].(bool); ok {
		result["enable_vtpm"] = vtpm
	}
	if integrityMonitoring, ok := shieldedVmConfig["enableIntegrityMonitoring"].(bool); ok {
		result["enable_integrity_monitoring"] = integrityMonitoring
	}

	return []map[string]bool{result}
}

func flattenEnableDisplayTgcNext(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	displayDevice, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	if enableDisplay, ok := displayDevice["enableDisplay"].(bool); ok && enableDisplay {
		return true
	}
	return nil
}

func flattenConfidentialInstanceConfigTgcNext(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	resp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return []map[string]interface{}{{
		"enable_confidential_compute": resp["enableConfidentialCompute"],
		"confidential_instance_type":  resp["confidentialInstanceType"],
	}}
}

func flattenAdvancedMachineFeaturesTgcNext(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	resp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return []map[string]interface{}{{
		"enable_nested_virtualization": resp["enableNestedVirtualization"],
		"threads_per_core":             resp["threadsPerCore"],
		"visible_core_count":           resp["visibleCoreCount"],
		"performance_monitoring_unit":  resp["performanceMonitoringUnit"],
		"enable_uefi_networking":       resp["enableUefiNetworking"],
		"turbo_mode":                   resp["turboMode"],
	}}
}

func flattenReservationAffinityTgcNext(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	affinity, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	consumeReservationType, ok := affinity["consumeReservationType"].(string)
	if !ok {
		return nil
	}
	crt := strings.ReplaceAll(consumeReservationType, "_ALLOCATION", "_RESERVATION")
	flattened := map[string]interface{}{
		"type": crt,
	}

	if crt == "SPECIFIC_RESERVATION" {
		flattened["specific_reservation"] = []map[string]interface{}{{
			"key":    affinity["key"],
			"values": affinity["values"],
		}}
	}

	return []map[string]interface{}{flattened}
}

func flattenComputeInstanceEncryptionKeyTgcNext(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	resp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return []map[string]interface{}{{
		"kms_key_self_link": resp["kmsKeyName"],
	}}
}

func flattenMetadataBetaTgcNext(v interface{}) map[string]interface{} {
	metadataMap := make(map[string]interface{})
	if v == nil {
		return metadataMap
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return metadataMap
	}
	items, ok := m["items"].([]interface{})
	if !ok {
		return metadataMap
	}
	for _, rawItem := range items {
		item, ok := rawItem.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := item["key"].(string)
		value, _ := item["value"].(string)
		if key != "" {
			metadataMap[key] = value
		}
	}
	return metadataMap
}
