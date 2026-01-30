package observability

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceObservabilityOrganizationSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceObservabilityOrganizationSettingsRead,
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The organization for which to retrieve settings.`,
			},
			"default_storage_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default storage location for new resources, e.g. buckets. Only valid for global location.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The location for which to retrieve settings.`,
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
				Description: `The service account used by Cloud Observability for this organization.`,
			},
		},
	}
}

func dataSourceObservabilityOrganizationSettingsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	organization := d.Get("organization").(string)
	location := d.Get("location").(string)

	url := fmt.Sprintf("%sorganizations/%s/locations/%s/settings", config.ObservabilityBasePath, organization, location)

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("ObservabilityOrganizationSettings %q", url), url)
	}

	d.SetId(res["name"].(string))

	if err := d.Set("organization", organization); err != nil {
		return fmt.Errorf("Error reading OrganizationSettings: %s", err)
	}
	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error reading OrganizationSettings: %s", err)
	}
	if v, ok := res["defaultStorageLocation"]; ok {
		if err := d.Set("default_storage_location", v); err != nil {
			return fmt.Errorf("Error reading OrganizationSettings: %s", err)
		}
	}
	if v, ok := res["kmsKeyName"]; ok {
		if err := d.Set("kms_key_name", v); err != nil {
			return fmt.Errorf("Error reading OrganizationSettings: %s", err)
		}
	}
	if v, ok := res["serviceAccountId"]; ok {
		if err := d.Set("service_account_id", v); err != nil {
			return fmt.Errorf("Error reading OrganizationSettings: %s", err)
		}
	}

	return nil
}
