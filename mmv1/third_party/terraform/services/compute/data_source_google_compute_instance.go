package compute

import (
	"fmt"
	"strconv"

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

	url := fmt.Sprintf("%sprojects/%s/zones/%s/instances/%s", config.ComputeBasePath, project, zone, name)
	instance, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("Instance %s", name), id)
	}

	var md map[string]string
	if metadataRaw, ok := instance["metadata"].(map[string]interface{}); ok {
		md = map[string]string{}
		if items, ok := metadataRaw["items"].([]interface{}); ok {
			for _, itemRaw := range items {
				if kv, ok := itemRaw.(map[string]interface{}); ok {
					k, _ := kv["key"].(string)
					v, _ := kv["value"].(string)
					md[k] = v
				}
			}
		}
	}
	if err = d.Set("metadata", md); err != nil {
		return fmt.Errorf("error setting metadata: %s", err)
	}

	if err := d.Set("can_ip_forward", instance["canIpForward"]); err != nil {
		return fmt.Errorf("Error setting can_ip_forward: %s", err)
	}
	machineTypeStr, _ := instance["machineType"].(string)
	if err := d.Set("machine_type", tpgresource.GetResourceNameFromSelfLink(machineTypeStr)); err != nil {
		return fmt.Errorf("Error setting machine_type: %s", err)
	}
	if err := d.Set("hostname", instance["hostname"]); err != nil {
		return fmt.Errorf("Error setting hostname: %s", err)
	}

	// Set the networks
	// Use the first external IP found for the default connection info.
	nisRaw, _ := instance["networkInterfaces"].([]interface{})
	networkInterfaces, _, internalIP, externalIP, err := flattenNetworkInterfaces(d, config, nisRaw)
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
	if metadataRaw, ok := instance["metadata"].(map[string]interface{}); ok {
		if err := d.Set("metadata_fingerprint", metadataRaw["fingerprint"]); err != nil {
			return fmt.Errorf("Error setting metadata_fingerprint: %s", err)
		}
	}

	// Set the tags fingerprint if there is one.
	if tagsRaw, ok := instance["tags"].(map[string]interface{}); ok && tagsRaw != nil {
		if err := d.Set("tags_fingerprint", tagsRaw["fingerprint"]); err != nil {
			return fmt.Errorf("Error setting tags_fingerprint: %s", err)
		}
		if err := d.Set("tags", tagsRaw["items"]); err != nil {
			return fmt.Errorf("Error setting tags: %s", err)
		}
	}

	if err := d.Set("labels", instance["labels"]); err != nil {
		return err
	}

	if err := d.Set("terraform_labels", instance["labels"]); err != nil {
		return err
	}

	if lf, ok := instance["labelFingerprint"].(string); ok && lf != "" {
		if err := d.Set("label_fingerprint", lf); err != nil {
			return fmt.Errorf("Error setting label_fingerprint: %s", err)
		}
	}

	attachedDisks := []map[string]interface{}{}
	scratchDisks := []map[string]interface{}{}
	for _, diskRaw := range instance["disks"].([]interface{}) {
		disk := diskRaw.(map[string]interface{})
		boot, _ := disk["boot"].(bool)
		diskType, _ := disk["type"].(string)
		diskSource, _ := disk["source"].(string)
		if boot {
			err = d.Set("boot_disk", flattenBootDisk(d, disk, config))
			if err != nil {
				return err
			}
		} else if diskType == "SCRATCH" {
			var diskSizeGb int64
			if s, ok := disk["diskSizeGb"].(string); ok {
				diskSizeGb, _ = strconv.ParseInt(s, 10, 64)
			} else if f, ok := disk["diskSizeGb"].(float64); ok {
				diskSizeGb = int64(f)
			}
			scratchDisks = append(scratchDisks, map[string]interface{}{
				"device_name": disk["deviceName"],
				"interface":   disk["interface"],
				"size":        diskSizeGb,
			})
		} else {
			di := map[string]interface{}{
				"source":      tpgresource.ConvertSelfLinkToV1(diskSource),
				"device_name": disk["deviceName"],
				"mode":        disk["mode"],
			}
			if keyRaw, ok := disk["diskEncryptionKey"].(map[string]interface{}); ok && keyRaw != nil {
				di["disk_encryption_key_sha256"] = keyRaw["sha256"]
				di["kms_key_self_link"] = keyRaw["kmsKeyName"]
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

	sasRaw, _ := instance["serviceAccounts"].([]interface{})
	err = d.Set("service_account", flattenServiceAccounts(sasRaw))
	if err != nil {
		return err
	}

	schedulingRaw, _ := instance["scheduling"].(map[string]interface{})
	err = d.Set("scheduling", flattenScheduling(schedulingRaw))
	if err != nil {
		return err
	}

	gasRaw, _ := instance["guestAccelerators"].([]interface{})
	err = d.Set("guest_accelerator", flattenGuestAccelerators(gasRaw))
	if err != nil {
		return err
	}

	err = d.Set("scratch_disk", scratchDisks)
	if err != nil {
		return err
	}

	siRaw, _ := instance["shieldedInstanceConfig"].(map[string]interface{})
	err = d.Set("shielded_instance_config", flattenShieldedVmConfig(siRaw))
	if err != nil {
		return err
	}

	ddRaw, _ := instance["displayDevice"].(map[string]interface{})
	err = d.Set("enable_display", flattenEnableDisplay(ddRaw))
	if err != nil {
		return err
	}

	if err := d.Set("attached_disk", ads); err != nil {
		return fmt.Errorf("Error setting attached_disk: %s", err)
	}
	if err := d.Set("cpu_platform", instance["cpuPlatform"]); err != nil {
		return fmt.Errorf("Error setting cpu_platform: %s", err)
	}
	if err := d.Set("min_cpu_platform", instance["minCpuPlatform"]); err != nil {
		return fmt.Errorf("Error setting min_cpu_platform: %s", err)
	}
	if err := d.Set("deletion_protection", instance["deletionProtection"]); err != nil {
		return fmt.Errorf("Error setting deletion_protection: %s", err)
	}
	selfLinkStr, _ := instance["selfLink"].(string)
	if err := d.Set("self_link", tpgresource.ConvertSelfLinkToV1(selfLinkStr)); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	instanceId, _ := instance["id"].(string)
	if err := d.Set("instance_id", instanceId); err != nil {
		return fmt.Errorf("Error setting instance_id: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	zoneStr, _ := instance["zone"].(string)
	if err := d.Set("zone", tpgresource.GetResourceNameFromSelfLink(zoneStr)); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}
	if err := d.Set("current_status", instance["status"]); err != nil {
		return fmt.Errorf("Error setting current_status: %s", err)
	}
	if err := d.Set("name", instance["name"]); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("key_revocation_action_type", instance["keyRevocationActionType"]); err != nil {
		return fmt.Errorf("Error setting key_revocation_action_type: %s", err)
	}
	if err := d.Set("creation_timestamp", instance["creationTimestamp"]); err != nil {
		return fmt.Errorf("Error setting creation_timestamp: %s", err)
	}

	instanceName, _ := instance["name"].(string)
	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, tpgresource.GetResourceNameFromSelfLink(zoneStr), instanceName))
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
