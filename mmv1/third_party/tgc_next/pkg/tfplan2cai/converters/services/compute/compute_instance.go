package compute

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/api/googleapi"

	compute "google.golang.org/api/compute/v0.beta"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/converters/cai"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const ComputeInstanceAssetType string = "compute.googleapis.com/Instance"
const ComputeDiskAssetType string = "compute.googleapis.com/Disk"

func ResourceConverterComputeInstance() cai.ResourceConverter {
	return cai.ResourceConverter{
		Convert: GetComputeInstanceAndDisksCaiObjects,
	}
}

func GetComputeInstanceAndDisksCaiObjects(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]caiasset.Asset, error) {
	if instanceAsset, err := GetComputeInstanceCaiObject(d, config); err == nil {
		assets := []caiasset.Asset{instanceAsset}
		if diskAsset, err := GetComputeDiskCaiObject(d, config); err == nil {
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
		return caiasset.Asset{
			Name: name,
			Type: ComputeInstanceAssetType,
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Instance",
				Data:                 data,
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

	sch := d.Get("scheduling").([]interface{})
	var scheduling *compute.Scheduling
	if len(sch) == 0 {
		// TF doesn't do anything about defaults inside of nested objects, so if
		// scheduling hasn't been set, then send it with its default values.
		scheduling = &compute.Scheduling{
			AutomaticRestart: googleapi.Bool(true),
		}
	} else {
		prefix := "scheduling.0"
		scheduling = &compute.Scheduling{
			AutomaticRestart:  googleapi.Bool(d.Get(prefix + ".automatic_restart").(bool)),
			Preemptible:       d.Get(prefix + ".preemptible").(bool),
			OnHostMaintenance: d.Get(prefix + ".on_host_maintenance").(string),
			ProvisioningModel: d.Get(prefix + ".provisioning_model").(string),
			ForceSendFields:   []string{"AutomaticRestart", "Preemptible"},
		}
	}

	params, err := expandParams(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating params: %s", err)
	}

	metadata, err := resourceInstanceMetadata(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating metadata: %s", err)
	}

	PartnerMetadata, err := resourceInstancePartnerMetadata(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating partner metadata: %s", err)
	}

	networkInterfaces, err := expandNetworkInterfaces(d, config)
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
	return &compute.Instance{
		CanIpForward:               d.Get("can_ip_forward").(bool),
		Description:                d.Get("description").(string),
		Disks:                      disks,
		MachineType:                machineTypeUrl,
		Metadata:                   metadata,
		PartnerMetadata:            PartnerMetadata,
		Name:                       d.Get("name").(string),
		Zone:                       d.Get("zone").(string),
		NetworkInterfaces:          networkInterfaces,
		NetworkPerformanceConfig:   networkPerformanceConfig,
		Tags:                       resourceInstanceTags(d),
		Params:                     params,
		Labels:                     tpgresource.ExpandLabels(d),
		ServiceAccounts:            expandServiceAccounts(d.Get("service_account").([]interface{})),
		GuestAccelerators:          accels,
		MinCpuPlatform:             d.Get("min_cpu_platform").(string),
		Scheduling:                 scheduling,
		DeletionProtection:         d.Get("deletion_protection").(bool),
		Hostname:                   d.Get("hostname").(string),
		ForceSendFields:            []string{"CanIpForward", "DeletionProtection"},
		ConfidentialInstanceConfig: expandConfidentialInstanceConfig(d),
		AdvancedMachineFeatures:    expandAdvancedMachineFeatures(d),
		ShieldedInstanceConfig:     expandShieldedVmConfigs(d),
		DisplayDevice:              expandDisplayDevice(d),
		ResourcePolicies:           tpgresource.ConvertStringArr(d.Get("resource_policies").([]interface{})),
		ReservationAffinity:        reservationAffinity,
		KeyRevocationActionType:    d.Get("key_revocation_action_type").(string),
	}, nil
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
		Source: sourceLink,
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

	if v, ok := d.GetOk("boot_disk.0.disk_encryption_key_raw"); ok {
		if v != "" {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				RawKey: v.(string),
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

	if v, ok := d.GetOk("boot_disk.0.source"); ok {
		source, err := tpgresource.ParseDiskFieldValue(v.(string), d, config)
		if err != nil {
			return nil, err
		}
		disk.Source = source.RelativeLink()
	}

	if _, ok := d.GetOk("boot_disk.0.initialize_params"); ok {
		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.size"); ok {
			disk.DiskSizeGb = int64(v.(int))
		}
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

func GetComputeDiskCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (caiasset.Asset, error) {
	name, err := cai.AssetName(d, config, "//compute.googleapis.com/projects/{{project}}/zones/{{zone}}/disks/{{name}}")
	if err != nil {
		return caiasset.Asset{}, err
	}
	if data, err := GetComputeDiskData(d, config); err == nil {
		return caiasset.Asset{
			Name: name,
			Type: ComputeDiskAssetType,
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Disk",
				Data:                 data,
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
