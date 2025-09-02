package backupdr

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleBackupDRDataSourceReferences() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleBackupDRDataSourceReferencesRead,
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location to list the data source references from.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the project in which the resource belongs.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource type to get the data source references for. Examples include, "compute.googleapis.com/Instance", "sqladmin.googleapis.com/Instance".`,
			},

			// Output: a computed list of the data source references found
			"data_source_references": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of the data source references found.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleBackupDRDataSourceReferencesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	location := d.Get("location").(string)
	resourceType := d.Get("resource_type").(string)

	url := fmt.Sprintf("https://backupdr.googleapis.com/v1/projects/%s/locations/%s/dataSourceReferences:fetchForResourceType?resourceType=%s", project, location, resourceType)

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error reading DataSourceReferences: %s", err)
	}

	items, ok := res["dataSourceReferences"].([]interface{})
	if !ok {
		items = make([]interface{}, 0)
	}

	flattenedDataSourceReferences, err := flattenDataSourceReferences(items)
	if err != nil {
		return err
	}

	if err := d.Set("data_source_references", flattenedDataSourceReferences); err != nil {
		return fmt.Errorf("Error setting data_source_references: %s", err)
	}

	d.SetId(url)

	return nil
}

func flattenDataSourceReferences(items []interface{}) ([]map[string]interface{}, error) {
	references := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot cast item to map[string]interface{}")
		}
		references = append(references, map[string]interface{}{
			"name":          data["name"],
			"resource_type": data["resourceType"],
		})
	}
	return references, nil
}
