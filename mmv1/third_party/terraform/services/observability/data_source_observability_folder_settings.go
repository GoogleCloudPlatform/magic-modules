package observability

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceObservabilityFolderSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceObservabilityFolderSettingsRead,
		Schema: map[string]*schema.Schema{
			"folder": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The folder for which to retrieve settings.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The location for which to retrieve settings.`,
			},
			"default_storage_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default storage location for new resources, e.g. buckets. Only valid for global location.`,
			},
			"kms_key_name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The default Cloud KMS key to use for new resources. Only valid for regional locations.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource name of the settings.`,
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The service account used by Cloud Observability for this folder.`,
			},
		},
	}
}

func dataSourceObservabilityFolderSettingsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	folder := d.Get("folder").(string)
	location := d.Get("location").(string)

	url := fmt.Sprintf("%sfolders/%s/locations/%s/settings", config.ObservabilityBasePath, folder, location)

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("ObservabilityFolderSettings %q", url), url)
	}

	d.SetId(res["name"].(string))

	if err := d.Set("folder", folder); err != nil {
		return fmt.Errorf("Error reading FolderSettings: %s", err)
	}
	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error reading FolderSettings: %s", err)
	}
	if v, ok := res["defaultStorageLocation"]; ok {
		if err := d.Set("default_storage_location", v); err != nil {
			return fmt.Errorf("Error reading FolderSettings: %s", err)
		}
	}
	if v, ok := res["kmsKeyName"]; ok {
		if err := d.Set("kms_key_name", v); err != nil {
			return fmt.Errorf("Error reading FolderSettings: %s", err)
		}
	}
	if v, ok := res["serviceAccountId"]; ok {
		if err := d.Set("service_account_id", v); err != nil {
			return fmt.Errorf("Error reading FolderSettings: %s", err)
		}
	}

	return nil
}
