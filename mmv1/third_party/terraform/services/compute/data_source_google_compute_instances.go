package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeInstances() *schema.Resource {
	resourceSchema := ResourceComputeInstance().Schema

	instanceSchema := map[string]*schema.Schema{
		"name":                {Type: schema.TypeString, Computed: true},
		"description":         {Type: schema.TypeString, Computed: true},
		"instance_id":         {Type: schema.TypeString, Computed: true},
		"zone":                {Type: schema.TypeString, Computed: true},
		"machine_type":        {Type: schema.TypeString, Computed: true},
		"self_link":           {Type: schema.TypeString, Computed: true},
		"project":             {Type: schema.TypeString, Computed: true},
		"current_status":      {Type: schema.TypeString, Computed: true},
		"cpu_platform":        {Type: schema.TypeString, Computed: true},
		"min_cpu_platform":    {Type: schema.TypeString, Computed: true},
		"can_ip_forward":      {Type: schema.TypeBool, Computed: true},
		"deletion_protection": {Type: schema.TypeBool, Computed: true},
		"hostname":            {Type: schema.TypeString, Computed: true},
		"label_fingerprint":   {Type: schema.TypeString, Computed: true},
		"creation_timestamp":  {Type: schema.TypeString, Computed: true},
		"enable_display":      {Type: schema.TypeBool, Computed: true},
		"labels": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tags": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"metadata": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"network_interface": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network":    {Type: schema.TypeString, Computed: true},
					"subnetwork": {Type: schema.TypeString, Computed: true},
					"network_ip": {Type: schema.TypeString, Computed: true},
					"nic_type":   {Type: schema.TypeString, Computed: true},
					"stack_type": {Type: schema.TypeString, Computed: true},
					"access_config": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"nat_ip":       {Type: schema.TypeString, Computed: true},
								"network_tier": {Type: schema.TypeString, Computed: true},
							},
						},
					},
				},
			},
		},
		"boot_disk": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"auto_delete":                {Type: schema.TypeBool, Computed: true},
					"device_name":                {Type: schema.TypeString, Computed: true},
					"mode":                       {Type: schema.TypeString, Computed: true},
					"source":                     {Type: schema.TypeString, Computed: true},
					"disk_encryption_key_sha256": {Type: schema.TypeString, Computed: true},
					"kms_key_self_link":          {Type: schema.TypeString, Computed: true},
				},
			},
		},
		"attached_disk": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"source":                     {Type: schema.TypeString, Computed: true},
					"device_name":                {Type: schema.TypeString, Computed: true},
					"mode":                       {Type: schema.TypeString, Computed: true},
					"disk_encryption_key_sha256": {Type: schema.TypeString, Computed: true},
					"kms_key_self_link":          {Type: schema.TypeString, Computed: true},
				},
			},
		},
		"scratch_disk": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_name": {Type: schema.TypeString, Computed: true},
					"interface":   {Type: schema.TypeString, Computed: true},
					"size":        {Type: schema.TypeInt, Computed: true},
				},
			},
		},
	}

	// These blocks are reused directly from the resource schema (converted to
	// read-only) so that the shared flatten* helpers used below -- which build
	// their output to match that schema -- stay valid here even as fields are
	// added to the resource.
	for _, key := range []string{"service_account", "scheduling", "guest_accelerator", "shielded_instance_config"} {
		instanceSchema[key] = &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: tpgresource.DatasourceSchemaFromResourceSchema(resourceSchema[key].Elem.(*schema.Resource).Schema),
			},
		}
	}

	return &schema.Resource{
		Read: dataSourceGoogleComputeInstancesRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The zone to list instances in, e.g. us-central1-a. If not provided, instances from all zones in the project are returned.`,
			},

			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A filter expression, as described in https://cloud.google.com/sdk/gcloud/reference/topic/filters, used to restrict the instances returned.`,
			},

			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `A list of instances matching the given filter, zone and project. See the google_compute_instance datasource for further documentation of the fields below. Fields that require an additional API call per instance (such as boot_disk.initialize_params) are not populated here, to keep this datasource usable against projects with many instances.`,
				Elem: &schema.Resource{
					Schema: instanceSchema,
				},
			},
		},
	}
}

func dataSourceGoogleComputeInstancesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	// Read straight off the schema rather than using tpgresource.GetZone: leaving
	// zone unset here means "every zone in the project", which is the point of
	// this datasource. Falling back to a provider-level default zone would
	// silently narrow that to a single zone for anyone who has one configured,
	// and GetZone returns an error entirely when neither is set.
	zone := d.Get("zone").(string)
	filter := d.Get("filter").(string)

	var baseURL string
	if zone != "" {
		baseURL, err = tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/instances")
	} else {
		baseURL, err = tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/aggregated/instances")
	}
	if err != nil {
		return err
	}

	instances := make([]map[string]interface{}, 0)

	pageToken := ""
	for {
		params := map[string]string{}
		if filter != "" {
			params["filter"] = filter
		}
		if pageToken != "" {
			params["pageToken"] = pageToken
		}

		url, err := transport_tpg.AddQueryParams(baseURL, params)
		if err != nil {
			return err
		}

		resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Instances : %s", project))
		}

		if zone != "" {
			for _, raw := range getInterfaceSlice(resp["items"]) {
				if instance, ok := raw.(map[string]interface{}); ok {
					instances = append(instances, flattenComputeInstancesListItem(instance, project))
				}
			}
		} else {
			for _, rawScope := range getMap(resp["items"]) {
				scope, ok := rawScope.(map[string]interface{})
				if !ok {
					continue
				}
				for _, raw := range getInterfaceSlice(scope["instances"]) {
					if instance, ok := raw.(map[string]interface{}); ok {
						instances = append(instances, flattenComputeInstancesListItem(instance, project))
					}
				}
			}
		}

		pageToken, _ = resp["nextPageToken"].(string)
		if pageToken == "" {
			break
		}
	}

	if err := d.Set("instances", instances); err != nil {
		return fmt.Errorf("Error retrieving instances: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	id := fmt.Sprintf("projects/%s/instances", project)
	if zone != "" {
		id = fmt.Sprintf("projects/%s/zones/%s/instances", project, zone)
	}
	if filter != "" {
		id = fmt.Sprintf("%s/filter=%s", id, filter)
	}
	d.SetId(id)

	return nil
}

func flattenComputeInstancesListItem(instance map[string]interface{}, project string) map[string]interface{} {
	bootDisks, attachedDisks, scratchDisks := flattenComputeInstancesListDisks(getInterfaceSlice(instance["disks"]))

	result := map[string]interface{}{
		"name":                stringOrEmpty(instance["name"]),
		"description":         stringOrEmpty(instance["description"]),
		"instance_id":         stringOrEmpty(instance["id"]),
		"zone":                tpgresource.GetResourceNameFromSelfLink(stringOrEmpty(instance["zone"])),
		"machine_type":        tpgresource.GetResourceNameFromSelfLink(stringOrEmpty(instance["machineType"])),
		"self_link":           tpgresource.ConvertSelfLinkToV1(stringOrEmpty(instance["selfLink"])),
		"project":             project,
		"current_status":      stringOrEmpty(instance["status"]),
		"cpu_platform":        stringOrEmpty(instance["cpuPlatform"]),
		"min_cpu_platform":    stringOrEmpty(instance["minCpuPlatform"]),
		"can_ip_forward":      boolOrFalse(instance["canIpForward"]),
		"deletion_protection": boolOrFalse(instance["deletionProtection"]),
		"hostname":            stringOrEmpty(instance["hostname"]),
		"label_fingerprint":   stringOrEmpty(instance["labelFingerprint"]),
		"creation_timestamp":  stringOrEmpty(instance["creationTimestamp"]),
		"labels":              instance["labels"],
		"metadata":            flattenMetadataBeta(getMap(instance["metadata"])),
		"network_interface":   flattenComputeInstancesListNetworkInterfaces(getInterfaceSlice(instance["networkInterfaces"])),
		"boot_disk":           bootDisks,
		"attached_disk":       attachedDisks,
		"scratch_disk":        scratchDisks,
		"service_account":     flattenServiceAccounts(getInterfaceSlice(instance["serviceAccounts"])),
		"scheduling":          flattenScheduling(getMap(instance["scheduling"])),
		"guest_accelerator":   flattenGuestAccelerators(getInterfaceSlice(instance["guestAccelerators"])),
	}

	if tags := getMap(instance["tags"]); tags != nil {
		result["tags"] = tags["items"]
	}

	if sic := getMap(instance["shieldedInstanceConfig"]); sic != nil {
		result["shielded_instance_config"] = flattenShieldedVmConfig(sic)
	}

	if dd := getMap(instance["displayDevice"]); dd != nil {
		result["enable_display"] = flattenEnableDisplay(dd)
	}

	return result
}

func flattenComputeInstancesListNetworkInterfaces(rawInterfaces []interface{}) []map[string]interface{} {
	interfaces := make([]map[string]interface{}, 0, len(rawInterfaces))
	for _, raw := range rawInterfaces {
		iface, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		accessConfigs := make([]map[string]interface{}, 0)
		for _, rawAC := range getInterfaceSlice(iface["accessConfigs"]) {
			ac, ok := rawAC.(map[string]interface{})
			if !ok {
				continue
			}
			accessConfigs = append(accessConfigs, map[string]interface{}{
				"nat_ip":       stringOrEmpty(ac["natIP"]),
				"network_tier": stringOrEmpty(ac["networkTier"]),
			})
		}

		interfaces = append(interfaces, map[string]interface{}{
			"network":       tpgresource.ConvertSelfLinkToV1(stringOrEmpty(iface["network"])),
			"subnetwork":    tpgresource.ConvertSelfLinkToV1(stringOrEmpty(iface["subnetwork"])),
			"network_ip":    stringOrEmpty(iface["networkIP"]),
			"nic_type":      stringOrEmpty(iface["nicType"]),
			"stack_type":    stringOrEmpty(iface["stackType"]),
			"access_config": accessConfigs,
		})
	}
	return interfaces
}

func flattenComputeInstancesListDisks(rawDisks []interface{}) (bootDisks, attachedDisks, scratchDisks []map[string]interface{}) {
	for _, raw := range rawDisks {
		disk, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		var encryptionKeySha256, kmsKeySelfLink string
		if key := getMap(disk["diskEncryptionKey"]); key != nil {
			encryptionKeySha256 = stringOrEmpty(key["sha256"])
			kmsKeySelfLink = stringOrEmpty(key["kmsKeyName"])
		}

		switch {
		case boolOrFalse(disk["boot"]):
			bootDisks = append(bootDisks, map[string]interface{}{
				"auto_delete":                boolOrFalse(disk["autoDelete"]),
				"device_name":                stringOrEmpty(disk["deviceName"]),
				"mode":                       stringOrEmpty(disk["mode"]),
				"source":                     tpgresource.ConvertSelfLinkToV1(stringOrEmpty(disk["source"])),
				"disk_encryption_key_sha256": encryptionKeySha256,
				"kms_key_self_link":          kmsKeySelfLink,
			})
		case stringOrEmpty(disk["type"]) == "SCRATCH":
			scratchDisks = append(scratchDisks, map[string]interface{}{
				"device_name": stringOrEmpty(disk["deviceName"]),
				"interface":   stringOrEmpty(disk["interface"]),
				"size":        getInt(disk["diskSizeGb"]),
			})
		default:
			attachedDisks = append(attachedDisks, map[string]interface{}{
				"source":                     tpgresource.ConvertSelfLinkToV1(stringOrEmpty(disk["source"])),
				"device_name":                stringOrEmpty(disk["deviceName"]),
				"mode":                       stringOrEmpty(disk["mode"]),
				"disk_encryption_key_sha256": encryptionKeySha256,
				"kms_key_self_link":          kmsKeySelfLink,
			})
		}
	}
	return
}

func stringOrEmpty(v interface{}) string {
	s, _ := v.(string)
	return s
}

func boolOrFalse(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

func getMap(v interface{}) map[string]interface{} {
	m, _ := v.(map[string]interface{})
	return m
}

func init() {
	registry.Schema{
		Name:        "google_compute_instances",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeInstances(),
	}.Register()
}
