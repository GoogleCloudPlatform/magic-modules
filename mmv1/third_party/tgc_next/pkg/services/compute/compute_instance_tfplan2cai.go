package compute

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/converters/cai"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/transport"
)

func ComputeInstanceTfplan2caiConverter() cai.Tfplan2caiConverter {
	return cai.Tfplan2caiConverter{
		Convert: GetComputeInstanceAndDisksCaiObjects,
	}
}

func GetComputeInstanceAndDisksCaiObjects(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]caiasset.Asset, error) {
	if instanceAsset, err := GetComputeInstanceCaiObject(d, config); err == nil {
		assets := []caiasset.Asset{instanceAsset}
		if diskAsset, err := GetComputeInstanceDiskCaiObject(d, config); err == nil {
			assets = append(assets, diskAsset)
			return assets, nil
		} else {
			return []caiasset.Asset{}, err
		}
	} else {
		return []caiasset.Asset{}, err
	}
}

func GetComputeInstanceCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (caiasset.Asset, error) {
	name, err := cai.AssetName(d, config, "//compute.googleapis.com/projects/{{project}}/zones/{{zone}}/instances/{{name}}")
	if err != nil {
		return caiasset.Asset{}, err
	}
	if data, err := GetComputeInstanceData(d, config); err == nil {
		location, _ := tpgresource.GetLocation(d, config)
		return caiasset.Asset{
			Name: name,
			Type: ComputeInstanceAssetType,
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Instance",
				Data:                 data,
				Location:             location,
			},
		}, nil
	} else {
		return caiasset.Asset{}, err
	}
}

func GetComputeInstanceData(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	instance, err := expandComputeInstance(project, d, config)
	if err != nil {
		return nil, err
	}

	return cai.JsonMap(instance)
}

func expandComputeInstance(project string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*compute.Instance, error) {
	// Get the machine type
	var machineTypeUrl string
	if mt, ok := d.GetOk("machine_type"); ok {
		machineType, err := tpgresource.ParseMachineTypesFieldValue(mt.(string), d, config)
		if err != nil {
			return nil, fmt.Errorf(
				"Error loading machine type: %s",
				err)
		}
		machineTypeUrl = machineType.RelativeLink()
	}

	// Build up the list of disks
	disks := []*compute.AttachedDisk{}
	if _, hasBootDisk := d.GetOk("boot_disk"); hasBootDisk {
		bootDisk, err := expandBootDisk(d, config, project)
		if err != nil {
			return nil, err
		}
		disks = append(disks, bootDisk)
	}

	if _, hasScratchDisk := d.GetOk("scratch_disk"); hasScratchDisk {
		scratchDisks, err := expandScratchDisks(d, config, project)
		if err != nil {
			return nil, err
		}
		disks = append(disks, scratchDisks...)
	}

	attachedDisksCount := d.Get("attached_disk.#").(int)

	for i := 0; i < attachedDisksCount; i++ {
		diskConfig := d.Get(fmt.Sprintf("attached_disk.%d", i)).(map[string]interface{})
		disk, err := expandAttachedDisk(diskConfig, d, config)
		if err != nil {
			return nil, err
		}

		disks = append(disks, disk)
	}

	scheduling, err := expandSchedulingTgc(d.Get("scheduling"))
	if err != nil {
		return nil, fmt.Errorf("error creating scheduling: %s", err)
	}

	params, err := expandParams(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating params: %s", err)
	}

	metadata, err := resourceInstanceMetadata(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating metadata: %s", err)
	}

	networkInterfaces, err := expandNetworkInterfacesTgc(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating network interfaces: %s", err)
	}

	networkPerformanceConfig, err := expandNetworkPerformanceConfig(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating network performance config: %s", err)
	}

	accels, err := expandInstanceGuestAccelerators(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating guest accelerators: %s", err)
	}

	reservationAffinity, err := expandReservationAffinity(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating reservation affinity: %s", err)
	}

	// Create the instance information
	instanceMap := map[string]interface{}{
		"canIpForward":               d.Get("can_ip_forward").(bool),
		"description":                d.Get("description").(string),
		"machineType":                machineTypeUrl,
		"name":                       d.Get("name").(string),
		"zone":                       d.Get("zone").(string),
		"networkInterfaces":          networkInterfaces,
		"networkPerformanceConfig":   networkPerformanceConfig,
		"tags":                       resourceInstanceTags(d),
		"labels":                     tpgresource.ExpandLabels(d),
		"serviceAccounts":            expandServiceAccounts(d.Get("service_account").([]interface{})),
		"minCpuPlatform":             d.Get("min_cpu_platform").(string),
		"scheduling":                 scheduling,
		"deletionProtection":         d.Get("deletion_protection").(bool),
		"hostname":                   d.Get("hostname").(string),
		"confidentialInstanceConfig": expandConfidentialInstanceConfig(d),
		"advancedMachineFeatures":    expandAdvancedMachineFeatures(d),
		"shieldedInstanceConfig":     expandShieldedVmConfigs(d),
		"displayDevice":              expandDisplayDevice(d),
		"reservationAffinity":        reservationAffinity,
		"keyRevocationActionType":    d.Get("key_revocation_action_type").(string),
		"instanceEncryptionKey":      expandComputeInstanceEncryptionKey(d),
		"metadata":                   metadata,
	}

	instanceBytes, err := json.Marshal(instanceMap)
	if err != nil {
		return nil, err
	}

	var instance compute.Instance
	if err := json.Unmarshal(instanceBytes, &instance); err != nil {
		return nil, err
	}

	// Re-assign typed slices/structs
	instance.Disks = disks
	instance.GuestAccelerators = accels
	instance.Params = params
	instance.ForceSendFields = []string{"CanIpForward", "DeletionProtection"}

	return &instance, nil
}

func expandAttachedDisk(diskConfig map[string]interface{}, d tpgresource.TerraformResourceData, meta interface{}) (*compute.AttachedDisk, error) {
	config := meta.(*transport_tpg.Config)

	s := diskConfig["source"].(string)
	var sourceLink string
	if strings.Contains(s, "regions/") {
		source, err := tpgresource.ParseRegionDiskFieldValue(s, d, config)
		if err != nil {
			return nil, err
		}
		sourceLink = source.RelativeLink()
	} else {
		source, err := tpgresource.ParseDiskFieldValue(s, d, config)
		if err != nil {
			return nil, err
		}
		sourceLink = source.RelativeLink()
	}

	disk := &compute.AttachedDisk{
		Source: fmt.Sprintf("https://www.googleapis.com/compute/v1/%s", sourceLink),
	}

	if v, ok := diskConfig["mode"]; ok {
		disk.Mode = v.(string)
	}

	if v, ok := diskConfig["device_name"]; ok {
		disk.DeviceName = v.(string)
	}

	keyValue, keyOk := diskConfig["disk_encryption_key_raw"]
	if keyOk {
		if keyValue != "" {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				RawKey: keyValue.(string),
			}
		}
	}

	keyValue, keyOk = diskConfig["disk_encryption_key_rsa"]
	if keyOk {
		if keyValue != "" {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				RsaEncryptedKey: keyValue.(string),
			}
		}
	}

	kmsValue, kmsOk := diskConfig["kms_key_self_link"]
	if kmsOk {
		if keyOk && keyValue != "" && kmsValue != "" {
			return nil, errors.New("Only one of kms_key_self_link and disk_encryption_key_raw can be set")
		}
		if kmsValue != "" {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				KmsKeyName: kmsValue.(string),
			}
		}
	}

	kmsServiceAccount, kmsServiceAccountOk := diskConfig["disk_encryption_service_account"]
	if kmsServiceAccountOk {
		if kmsServiceAccount != "" {
			if disk.DiskEncryptionKey == nil {
				disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
					KmsKeyServiceAccount: kmsServiceAccount.(string),
				}
			}
			disk.DiskEncryptionKey.KmsKeyServiceAccount = kmsServiceAccount.(string)
		}
	}
	return disk, nil
}

// See comment on expandInstanceTemplateGuestAccelerators regarding why this
// code is duplicated.
func expandInstanceGuestAccelerators(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]*compute.AcceleratorConfig, error) {
	configs, ok := d.GetOk("guest_accelerator")
	if !ok {
		return nil, nil
	}
	accels := configs.([]interface{})
	guestAccelerators := make([]*compute.AcceleratorConfig, 0, len(accels))
	for _, raw := range accels {
		data := raw.(map[string]interface{})
		if data["count"].(int) == 0 {
			continue
		}
		at, err := tpgresource.ParseAcceleratorFieldValue(data["type"].(string), d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot parse accelerator type: %v", err)
		}
		guestAccelerators = append(guestAccelerators, &compute.AcceleratorConfig{
			AcceleratorCount: int64(data["count"].(int)),
			AcceleratorType:  at.RelativeLink(),
		})
	}

	return guestAccelerators, nil
}

func expandParams(d tpgresource.TerraformResourceData) (*compute.InstanceParams, error) {
	if _, ok := d.GetOk("params.0.resource_manager_tags"); ok {
		params := &compute.InstanceParams{
			ResourceManagerTags: tpgresource.ExpandStringMap(d, "params.0.resource_manager_tags"),
		}
		return params, nil
	}

	return nil, nil
}

func expandBootDisk(d tpgresource.TerraformResourceData, config *transport_tpg.Config, project string) (*compute.AttachedDisk, error) {
	disk := &compute.AttachedDisk{
		AutoDelete: d.Get("boot_disk.0.auto_delete").(bool),
		Boot:       true,
	}

	if v, ok := d.GetOk("boot_disk.0.device_name"); ok {
		disk.DeviceName = v.(string)
	}

	if v, ok := d.GetOk("boot_disk.0.interface"); ok {
		disk.Interface = v.(string)
	}

	if v, ok := d.GetOk("boot_disk.0.guest_os_features"); ok {
		var typedGuestOsFeatures []*compute.GuestOsFeature
		if gofBytes, err := json.Marshal(expandComputeInstanceGuestOsFeatures(v)); err == nil {
			json.Unmarshal(gofBytes, &typedGuestOsFeatures)
		}
		disk.GuestOsFeatures = typedGuestOsFeatures
	}

	if v, ok := d.GetOk("boot_disk.0.disk_encryption_key_raw"); ok {
		if v != "" {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				RawKey: v.(string),
			}
		}
	}

	if v, ok := d.GetOk("boot_disk.0.disk_encryption_key_rsa"); ok {
		if v != "" {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				RsaEncryptedKey: v.(string),
			}
		}
	}

	if v, ok := d.GetOk("boot_disk.0.kms_key_self_link"); ok {
		if v != "" {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				KmsKeyName: v.(string),
			}
		}
	}

	if v, ok := d.GetOk("boot_disk.0.disk_encryption_service_account"); ok {
		if v != "" {
			disk.DiskEncryptionKey.KmsKeyServiceAccount = v.(string)
		}
	}

	// disk_encryption_key_sha256 is computed, so it is not converted.

	if v, ok := d.GetOk("boot_disk.0.source"); ok {
		var err error
		var source interface {
			RelativeLink() string
		}
		if strings.Contains(v.(string), "regions/") {
			source, err = tpgresource.ParseRegionDiskFieldValue(v.(string), d, config)
		} else {
			source, err = tpgresource.ParseDiskFieldValue(v.(string), d, config)
		}
		if err != nil {
			return nil, err
		}
		disk.Source = fmt.Sprintf("https://www.googleapis.com/compute/v1/%s", source.RelativeLink())
	}

	if _, ok := d.GetOk("boot_disk.0.initialize_params"); ok {
		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.size"); ok {
			disk.DiskSizeGb = int64(v.(int))
		}
	}

	if v, ok := d.GetOk("boot_disk.0.initialize_params.0.architecture"); ok {
		disk.Architecture = v.(string)
	}

	if v, ok := d.GetOk("boot_disk.0.mode"); ok {
		disk.Mode = v.(string)
	}

	return disk, nil
}

func expandScratchDisks(d tpgresource.TerraformResourceData, config *transport_tpg.Config, project string) ([]*compute.AttachedDisk, error) {
	diskType, err := readDiskType(config, d, "local-ssd")
	if err != nil {
		return nil, fmt.Errorf("Error loading disk type 'local-ssd': %s", err)
	}

	n := d.Get("scratch_disk.#").(int)
	scratchDisks := make([]*compute.AttachedDisk, 0, n)
	for i := 0; i < n; i++ {
		scratchDisks = append(scratchDisks, &compute.AttachedDisk{
			AutoDelete: true,
			Type:       "SCRATCH",
			Interface:  d.Get(fmt.Sprintf("scratch_disk.%d.interface", i)).(string),
			DeviceName: d.Get(fmt.Sprintf("scratch_disk.%d.device_name", i)).(string),
			DiskSizeGb: int64(d.Get(fmt.Sprintf("scratch_disk.%d.size", i)).(int)),
			InitializeParams: &compute.AttachedDiskInitializeParams{
				DiskType: diskType.RelativeLink(),
			},
		})
	}

	return scratchDisks, nil
}

func expandStoragePool(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	// ExpandStoragePoolUrl is generated by MMv1
	// return ExpandStoragePoolUrl(v, d, config)
	return nil, nil
}

func GetComputeInstanceDiskCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (caiasset.Asset, error) {
	name, err := cai.AssetName(d, config, "//compute.googleapis.com/projects/{{project}}/zones/{{zone}}/disks/{{name}}")
	if err != nil {
		return caiasset.Asset{}, err
	}
	if data, err := GetComputeDiskData(d, config); err == nil {
		location, _ := tpgresource.GetLocation(d, config)
		return caiasset.Asset{
			Name: name,
			Type: ComputeDiskAssetType,
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Disk",
				Data:                 data,
				Location:             location,
			},
		}, nil
	} else {
		return caiasset.Asset{}, err
	}
}

func GetComputeDiskData(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	diskApiObj, err := expandBootDisk(d, config, project)
	if err != nil {
		return nil, err
	}

	diskDetails, err := cai.JsonMap(diskApiObj)
	if err != nil {
		return nil, err
	}

	if v, ok := d.GetOk("boot_disk.0.initialize_params.0.type"); ok {
		diskTypeName := v.(string)
		diskType, err := readDiskType(config, d, diskTypeName)
		if err != nil {
			return nil, fmt.Errorf("Error loading disk type '%s': %s", diskTypeName, err)
		}
		diskDetails["DiskType"] = diskType.RelativeLink()
	}

	if v, ok := d.GetOk("boot_disk.0.initialize_params.0.image"); ok {
		diskDetails["SourceImage"] = v.(string)
	}

	if _, ok := d.GetOk("boot_disk.0.initialize_params.0.labels"); ok {
		diskDetails["Labels"] = tpgresource.ExpandStringMap(d, "boot_disk.0.initialize_params.0.labels")
	}

	if _, ok := d.GetOk("boot_disk.0.initialize_params.0.resource_policies"); ok {
		diskDetails["ResourcePolicies"] = tpgresource.ConvertStringArr(d.Get("boot_disk.0.initialize_params.0.resource_policies").([]interface{}))
	}

	if v, ok := d.GetOk("boot_disk.0.initialize_params.0.provisioned_iops"); ok {
		diskDetails["ProvisionedIops"] = int64(v.(int))
	}

	if v, ok := d.GetOk("boot_disk.0.initialize_params.0.provisioned_throughput"); ok {
		diskDetails["ProvisionedThroughput"] = int64(v.(int))
	}

	if v, ok := d.GetOk("boot_disk.0.initialize_params.0.enable_confidential_compute"); ok {
		diskDetails["EnableConfidentialCompute"] = v.(bool)
	}

	if v, ok := d.GetOk("boot_disk.0.initialize_params.0.storage_pool"); ok {
		storagePoolUrl, err := expandStoragePool(v, d, config)
		if err != nil {
			return nil, fmt.Errorf("Error resolving storage pool name '%s': '%s'", v.(string), err)
		}
		diskDetails["StoragePool"] = storagePoolUrl.(string)
	}

	return diskDetails, nil
}

func expandNetworkInterfacesTgc(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]*compute.NetworkInterface, error) {
	configs := d.Get("network_interface").([]interface{})
	ifaces := make([]*compute.NetworkInterface, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		network := data["network"].(string)
		subnetwork := data["subnetwork"].(string)

		ifaceMap := map[string]interface{}{
			"networkIP":                data["network_ip"].(string),
			"network":                  network,
			"subnetwork":               subnetwork,
			"accessConfigs":            expandAccessConfigs(data["access_config"].([]interface{})),
			"aliasIpRanges":            expandAliasIpRanges(data["alias_ip_range"].([]interface{})),
			"nicType":                  data["nic_type"].(string),
			"stackType":                data["stack_type"].(string),
			"queueCount":               strconv.FormatInt(int64(data["queue_count"].(int)), 10),
			"ipv6AccessConfigs":        expandIpv6AccessConfigs(data["ipv6_access_config"].([]interface{})),
			"ipv6Address":              data["ipv6_address"].(string),
			"internalIpv6PrefixLength": strconv.FormatInt(int64(data["internal_ipv6_prefix_length"].(int)), 10),
		}

		ifaceBytes, err := json.Marshal(ifaceMap)
		if err != nil {
			return nil, err
		}
		var iface compute.NetworkInterface
		if err := json.Unmarshal(ifaceBytes, &iface); err != nil {
			return nil, err
		}
		ifaces[i] = &iface
	}
	return ifaces, nil
}

func expandSchedulingTgc(v interface{}) (*compute.Scheduling, error) {
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
	schedulingMap := make(map[string]interface{})

	if v, ok := original["automatic_restart"]; ok {
		schedulingMap["automaticRestart"] = v.(bool)
	}

	if v, ok := original["preemptible"]; ok {
		schedulingMap["preemptible"] = v.(bool)
	}

	if v, ok := original["on_host_maintenance"]; ok {
		schedulingMap["onHostMaintenance"] = v.(string)
	}

	if v, ok := original["node_affinities"]; ok && v != nil {
		naSet := v.(*schema.Set).List()
		var nodeAffinities []interface{}
		for _, nodeAffRaw := range naSet {
			if nodeAffRaw == nil {
				continue
			}
			nodeAff := nodeAffRaw.(map[string]interface{})
			transformed := map[string]interface{}{
				"key":      nodeAff["key"].(string),
				"operator": nodeAff["operator"].(string),
				"values":   tpgresource.ConvertStringArr(nodeAff["values"].(*schema.Set).List()),
			}
			nodeAffinities = append(nodeAffinities, transformed)
		}
		schedulingMap["nodeAffinities"] = nodeAffinities
	}

	if v, ok := original["min_node_cpus"]; ok {
		schedulingMap["minNodeCpus"] = strconv.FormatInt(int64(v.(int)), 10)
	}
	if v, ok := original["provisioning_model"]; ok {
		schedulingMap["provisioningModel"] = v.(string)
	}
	if v, ok := original["instance_termination_action"]; ok {
		schedulingMap["instanceTerminationAction"] = v.(string)
	}
	if v, ok := original["availability_domain"]; ok && v != nil {
		schedulingMap["availabilityDomain"] = strconv.FormatInt(int64(v.(int)), 10)
	}
	if v, ok := original["max_run_duration"]; ok {
		transformedMaxRunDuration, err := expandComputeMaxRunDuration(v)
		if err != nil {
			return nil, err
		}
		schedulingMap["maxRunDuration"] = transformedMaxRunDuration
	}

	if v, ok := original["on_instance_stop_action"]; ok {
		transformedOnInstanceStopAction, err := expandComputeOnInstanceStopAction(v)
		if err != nil {
			return nil, err
		}
		schedulingMap["onInstanceStopAction"] = transformedOnInstanceStopAction
	}

	if v, ok := original["local_ssd_recovery_timeout"]; ok {
		transformedLocalSsdRecoveryTimeout, err := expandComputeLocalSsdRecoveryTimeout(v)
		if err != nil {
			return nil, err
		}
		schedulingMap["localSsdRecoveryTimeout"] = transformedLocalSsdRecoveryTimeout
	}
	if v, ok := original["termination_time"]; ok {
		schedulingMap["terminationTime"] = v.(string)
	}

	schedulingBytes, err := json.Marshal(schedulingMap)
	if err != nil {
		return nil, err
	}
	var sched compute.Scheduling
	if err := json.Unmarshal(schedulingBytes, &sched); err != nil {
		return nil, err
	}
	scheduling := &sched
	scheduling.ForceSendFields = []string{"AutomaticRestart", "Preemptible"}
	return scheduling, nil
}
