package compute

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tpgresource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	compute "google.golang.org/api/compute/v0.beta"
)

// ComputeInstanceAssetType is the CAI asset type name for compute instance.
const ComputeInstanceAssetType string = "compute.googleapis.com/Instance"

// ComputeInstanceSchemaName is the TF resource schema name for compute instance.
const ComputeInstanceSchemaName string = "google_compute_instance"

// ComputeInstanceConverter for compute instance resource.
type ComputeInstanceConverter struct {
	name   string
	schema map[string]*schema.Schema
}

// NewComputeInstanceConverter returns an HCL converter for compute instance.
func NewComputeInstanceConverter(provider *schema.Provider) models.Converter {
	schema := provider.ResourcesMap[ComputeInstanceSchemaName].Schema

	return &ComputeInstanceConverter{
		name:   ComputeInstanceSchemaName,
		schema: schema,
	}
}

// Convert converts asset to HCL resource blocks.
func (c *ComputeInstanceConverter) Convert(asset caiasset.Asset) ([]*models.TerraformResourceBlock, error) {
	var blocks []*models.TerraformResourceBlock
	block, err := c.convertResourceData(asset)
	if err != nil {
		return nil, err
	}
	blocks = append(blocks, block)
	return blocks, nil
}

func (c *ComputeInstanceConverter) convertResourceData(asset caiasset.Asset) (*models.TerraformResourceBlock, error) {
	if asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}

	project := utils.ParseFieldValue(asset.Name, "projects")

	var instance *compute.Instance
	if err := utils.DecodeJSON(asset.Resource.Data, &instance); err != nil {
		return nil, err
	}

	hclData := make(map[string]interface{})

	if instance.CanIpForward {
		hclData["can_ip_forward"] = instance.CanIpForward
	}
	hclData["machine_type"] = tpgresource.GetResourceNameFromSelfLink(instance.MachineType)
	hclData["network_performance_config"] = flattenNetworkPerformanceConfig(instance.NetworkPerformanceConfig)

	// Set the networks
	networkInterfaces, _, _, err := flattenNetworkInterfaces(instance.NetworkInterfaces, project)
	if err != nil {
		return nil, err
	}
	hclData["network_interface"] = networkInterfaces

	if instance.Tags != nil {
		hclData["tags"] = tpgresource.ConvertStringArrToInterface(instance.Tags.Items)
	}

	hclData["labels"] = utils.RemoveTerraformAttributionLabel(instance.Labels)
	hclData["service_account"] = flattenServiceAccounts(instance.ServiceAccounts)
	hclData["resource_policies"] = instance.ResourcePolicies

	bootDisk, ads, scratchDisks := flattenDisks(instance.Disks, instance.Name)
	hclData["boot_disk"] = bootDisk
	hclData["attached_disk"] = ads
	hclData["scratch_disk"] = scratchDisks

	hclData["scheduling"] = flattenScheduling(instance.Scheduling)
	hclData["guest_accelerator"] = flattenGuestAccelerators(instance.GuestAccelerators)
	hclData["shielded_instance_config"] = flattenShieldedVmConfig(instance.ShieldedInstanceConfig)
	hclData["enable_display"] = flattenEnableDisplay(instance.DisplayDevice)
	hclData["min_cpu_platform"] = instance.MinCpuPlatform

	// Only convert the field when its value is not default false
	if instance.DeletionProtection {
		hclData["deletion_protection"] = instance.DeletionProtection
	}
	hclData["zone"] = tpgresource.GetResourceNameFromSelfLink(instance.Zone)
	hclData["name"] = instance.Name
	hclData["description"] = instance.Description
	hclData["hostname"] = instance.Hostname
	hclData["confidential_instance_config"] = flattenConfidentialInstanceConfig(instance.ConfidentialInstanceConfig)
	hclData["advanced_machine_features"] = flattenAdvancedMachineFeatures(instance.AdvancedMachineFeatures)
	hclData["reservation_affinity"] = flattenReservationAffinity(instance.ReservationAffinity)
	hclData["key_revocation_action_type"] = instance.KeyRevocationActionType

	// TODO: convert details from the boot disk assets (separate disk assets) into initialize_params in cai2hcl?
	// It needs to integrate the disk assets into instance assets with the resolver.

	ctyVal, err := utils.MapToCtyValWithSchema(hclData, c.schema)
	if err != nil {
		return nil, err
	}
	return &models.TerraformResourceBlock{
		Labels: []string{c.name, instance.Name},
		Value:  ctyVal,
	}, nil

}

func flattenDisks(disks []*compute.AttachedDisk, instanceName string) ([]map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {
	attachedDisks := []map[string]interface{}{}
	bootDisk := []map[string]interface{}{}
	scratchDisks := []map[string]interface{}{}
	for _, disk := range disks {
		if disk.Boot {
			bootDisk = flattenBootDisk(disk, instanceName)
		} else if disk.Type == "SCRATCH" {
			scratchDisks = append(scratchDisks, flattenScratchDisk(disk))
		} else {
			di := map[string]interface{}{
				"source":      tpgresource.ConvertSelfLinkToV1(disk.Source),
				"device_name": disk.DeviceName,
				"mode":        disk.Mode,
			}
			if key := disk.DiskEncryptionKey; key != nil {
				if key.KmsKeyName != "" {
					// The response for crypto keys often includes the version of the key which needs to be removed
					// format: projects/<project>/locations/<region>/keyRings/<keyring>/cryptoKeys/<key>/cryptoKeyVersions/1
					di["kms_key_self_link"] = strings.Split(disk.DiskEncryptionKey.KmsKeyName, "/cryptoKeyVersions")[0]
				}
			}
			attachedDisks = append(attachedDisks, di)
		}
	}

	// Remove nils from map in case there were disks in the config that were not present on read;
	// i.e. a disk was detached out of band
	ads := []map[string]interface{}{}
	for _, d := range attachedDisks {
		if d != nil {
			ads = append(ads, d)
		}
	}
	return bootDisk, ads, scratchDisks
}

func flattenBootDisk(disk *compute.AttachedDisk, instanceName string) []map[string]interface{} {
	result := map[string]interface{}{}

	if !disk.AutoDelete {
		result["auto_delete"] = false
	}

	if !strings.Contains(disk.DeviceName, "persistent-disk-") {
		result["device_name"] = disk.DeviceName
	}

	if disk.Mode != "READ_WRITE" {
		result["mode"] = disk.Mode
	}

	if disk.DiskEncryptionKey != nil {
		if disk.DiskEncryptionKey.KmsKeyName != "" {
			// The response for crypto keys often includes the version of the key which needs to be removed
			// format: projects/<project>/locations/<region>/keyRings/<keyring>/cryptoKeys/<key>/cryptoKeyVersions/1
			result["kms_key_self_link"] = strings.Split(disk.DiskEncryptionKey.KmsKeyName, "/cryptoKeyVersions")[0]
		}
	}

	// Don't convert the field with the default value
	if disk.Interface != "SCSI" {
		result["interface"] = disk.Interface
	}

	if !strings.HasSuffix(disk.Source, instanceName) {
		result["source"] = tpgresource.ConvertSelfLinkToV1(disk.Source)
	}

	if len(result) == 0 {
		return nil
	}

	return []map[string]interface{}{result}
}

func flattenScratchDisk(disk *compute.AttachedDisk) map[string]interface{} {
	result := map[string]interface{}{
		"size": disk.DiskSizeGb,
	}

	if !strings.Contains(disk.DeviceName, "persistent-disk-") {
		result["device_name"] = disk.DeviceName
	}

	// Don't convert the field with the default value
	if disk.Interface != "SCSI" {
		result["interface"] = disk.Interface
	}

	return result
}
