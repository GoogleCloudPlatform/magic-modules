package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeInstance() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeInstance().Schema)

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "name", "self_link", "project", "zone")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, zone, name, err := tpgresource.GetZonalResourcePropertiesFromSelfLinkOrSchema(d, config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, zone, name)

	instanceURL := fmt.Sprintf("%sprojects/%s/zones/%s/instances/%s", transport_tpg.BaseUrl(Product, config), project, zone, name)
	instance, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    instanceURL,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("Instance %s", name), id)
	}

	metadataMap, _ := instance["metadata"].(map[string]interface{})
	md := flattenMetadataFromApiMap(metadataMap)
	if err = d.Set("metadata", md); err != nil {
		return fmt.Errorf("error setting metadata: %s", err)
	}

	if err := d.Set("can_ip_forward", instance["canIpForward"]); err != nil {
		return fmt.Errorf("Error setting can_ip_forward: %s", err)
	}
	machineType, _ := instance["machineType"].(string)
	if err := d.Set("machine_type", tpgresource.GetResourceNameFromSelfLink(machineType)); err != nil {
		return fmt.Errorf("Error setting machine_type: %s", err)
	}
	hostname, _ := instance["hostname"].(string)
	if err := d.Set("hostname", hostname); err != nil {
		return fmt.Errorf("Error setting hostname: %s", err)
	}

	// Set the networks
	// Use the first external IP found for the default connection info.
	niRaw, _ := instance["networkInterfaces"].([]interface{})
	networkInterfaces, _, internalIP, externalIP, err := flattenNetworkInterfaces(d, config, niRaw)
	if err != nil {
		return err
	}
	if err := d.Set("network_interface", networkInterfaces); err != nil {
		return err
	}

	// Fall back on internal ip if there is no external ip.  This makes sense in the situation where
	// terraform is being used on a cloud instance and can therefore access the instances it creates
	// via their internal ips.
	sshIP := externalIP
	if sshIP == "" {
		sshIP = internalIP
	}

	// Initialize the connection info
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": sshIP,
	})

	// Set the metadata fingerprint if there is one.
	if metadataMap != nil {
		fp, _ := metadataMap["fingerprint"].(string)
		if err := d.Set("metadata_fingerprint", fp); err != nil {
			return fmt.Errorf("Error setting metadata_fingerprint: %s", err)
		}
	}

	// Set the tags fingerprint if there is one.
	if tags, ok := instance["tags"].(map[string]interface{}); ok && tags != nil {
		fp, _ := tags["fingerprint"].(string)
		if err := d.Set("tags_fingerprint", fp); err != nil {
			return fmt.Errorf("Error setting tags_fingerprint: %s", err)
		}
		var tagItems []string
		if items, ok := tags["items"].([]interface{}); ok {
			for _, item := range items {
				if s, ok := item.(string); ok {
					tagItems = append(tagItems, s)
				}
			}
		}
		if err := d.Set("tags", tpgresource.ConvertStringArrToInterface(tagItems)); err != nil {
			return fmt.Errorf("Error setting tags: %s", err)
		}
	}

	if labels, ok := instance["labels"].(map[string]interface{}); ok {
		if err := d.Set("labels", labels); err != nil {
			return err
		}
		if err := d.Set("terraform_labels", labels); err != nil {
			return err
		}
	}

	if lf, ok := instance["labelFingerprint"].(string); ok && lf != "" {
		if err := d.Set("label_fingerprint", lf); err != nil {
			return fmt.Errorf("Error setting label_fingerprint: %s", err)
		}
	}

	attachedDisks := []map[string]interface{}{}
	scratchDisks := []map[string]interface{}{}
	instanceDisks, _ := instance["disks"].([]interface{})
	for _, diskRaw := range instanceDisks {
		disk, ok := diskRaw.(map[string]interface{})
		if !ok {
			continue
		}
		isBoot, _ := disk["boot"].(bool)
		diskType, _ := disk["type"].(string)
		if isBoot {
			err = d.Set("boot_disk", flattenBootDisk(d, disk, config))
			if err != nil {
				return err
			}
		} else if diskType == "SCRATCH" {
			scratchDisks = append(scratchDisks, flattenScratchDisk(disk))
		} else {
			diskSource, _ := disk["source"].(string)
			deviceName, _ := disk["deviceName"].(string)
			mode, _ := disk["mode"].(string)
			di := map[string]interface{}{
				"source":      tpgresource.ConvertSelfLinkToV1(diskSource),
				"device_name": deviceName,
				"mode":        mode,
			}
			if key, ok := disk["diskEncryptionKey"].(map[string]interface{}); ok && key != nil {
				if sha256, ok := key["sha256"].(string); ok {
					di["disk_encryption_key_sha256"] = sha256
				}
				if kmsKey, ok := key["kmsKeyName"].(string); ok {
					di["kms_key_self_link"] = kmsKey
				}
			}
			attachedDisks = append(attachedDisks, di)
		}
	}
	// Remove nils from map in case there were disks in the config that were not present on read;
	// i.e. a disk was detached out of band
	ads := []map[string]interface{}{}
	for _, ad := range attachedDisks {
		if ad != nil {
			ads = append(ads, ad)
		}
	}

	serviceAccounts, _ := instance["serviceAccounts"].([]interface{})
	err = d.Set("service_account", flattenServiceAccounts(serviceAccounts))
	if err != nil {
		return err
	}

	schedulingMap, _ := instance["scheduling"].(map[string]interface{})
	err = d.Set("scheduling", flattenScheduling(schedulingMap))
	if err != nil {
		return err
	}

	accelerators, _ := instance["accelerators"].([]interface{})
	err = d.Set("guest_accelerator", flattenGuestAccelerators(accelerators))
	if err != nil {
		return err
	}

	err = d.Set("scratch_disk", scratchDisks)
	if err != nil {
		return err
	}

	shieldedConfig, _ := instance["shieldedInstanceConfig"].(map[string]interface{})
	err = d.Set("shielded_instance_config", flattenShieldedVmConfig(shieldedConfig))
	if err != nil {
		return err
	}

	displayDevice, _ := instance["displayDevice"].(map[string]interface{})
	err = d.Set("enable_display", flattenEnableDisplay(displayDevice))
	if err != nil {
		return err
	}

	if err := d.Set("attached_disk", ads); err != nil {
		return fmt.Errorf("Error setting attached_disk: %s", err)
	}
	cpuPlatform, _ := instance["cpuPlatform"].(string)
	if err := d.Set("cpu_platform", cpuPlatform); err != nil {
		return fmt.Errorf("Error setting cpu_platform: %s", err)
	}
	minCpuPlatform, _ := instance["minCpuPlatform"].(string)
	if err := d.Set("min_cpu_platform", minCpuPlatform); err != nil {
		return fmt.Errorf("Error setting min_cpu_platform: %s", err)
	}
	if err := d.Set("deletion_protection", instance["deletionProtection"]); err != nil {
		return fmt.Errorf("Error setting deletion_protection: %s", err)
	}
	selfLink, _ := instance["selfLink"].(string)
	if err := d.Set("self_link", tpgresource.ConvertSelfLinkToV1(selfLink)); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	var instanceIDStr string
	if idVal, ok := instance["id"].(float64); ok {
		instanceIDStr = fmt.Sprintf("%d", int64(idVal))
	}
	if err := d.Set("instance_id", instanceIDStr); err != nil {
		return fmt.Errorf("Error setting instance_id: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	instZone, _ := instance["zone"].(string)
	if err := d.Set("zone", tpgresource.GetResourceNameFromSelfLink(instZone)); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}
	status, _ := instance["status"].(string)
	if err := d.Set("current_status", status); err != nil {
		return fmt.Errorf("Error setting current_status: %s", err)
	}
	instName, _ := instance["name"].(string)
	if err := d.Set("name", instName); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	keyRevocation, _ := instance["keyRevocationActionType"].(string)
	if err := d.Set("key_revocation_action_type", keyRevocation); err != nil {
		return fmt.Errorf("Error setting key_revocation_action_type: %s", err)
	}
	creationTimestamp, _ := instance["creationTimestamp"].(string)
	if err := d.Set("creation_timestamp", creationTimestamp); err != nil {
		return fmt.Errorf("Error setting creation_timestamp: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, tpgresource.GetResourceNameFromSelfLink(instZone), instName))
	return nil
}

func init() {
	registry.Schema{
		Name:        "google_compute_instance",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeInstance(),
	}.Register()
}
