package backupdr

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudBackupDRListDataSourceReferences() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleCloudBackupDRListDataSourceReferencesRead,
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
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The filter to apply to the list results.",
			},
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The order to sort results by.",
			},
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
						"data_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The underlying data source resource.",
						},
						"backup_config_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the backup config for the data source.",
						},
						"backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of backups for the data source.",
						},
						"last_backup_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the last backup.",
						},
						"last_successful_backup_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last time a successful backup was made.",
						},
						"gcp_resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The GCP resource name for the data source.",
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

func dataSourceGoogleCloudBackupDRListDataSourceReferencesRead(d *schema.ResourceData, meta interface{}) error {
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
	url := fmt.Sprintf("%sprojects/%s/locations/%s/dataSourceReferences", config.BackupDRBasePath, project, location)

	params := make(map[string]string)
	if v, ok := d.GetOk("filter"); ok {
		params["filter"] = v.(string)
	}
	if v, ok := d.GetOk("order_by"); ok {
		// API expects `orderBy` as the query parameter name
		params["orderBy"] = v.(string)
	}

	// Attach query params to the URL so the API receives filter/orderBy
	if len(params) > 0 {
		var err error
		url, err = transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return fmt.Errorf("Error adding query params to URL: %s", err)
		}
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error listing DataSourceReferences: %s", err)
	}

	items, ok := res["dataSourceReferences"].([]interface{})
	if !ok {
		items = make([]interface{}, 0)
	}

	flattenedDataSourceReferences, err := flattenDataSourceReferencesList(items)
	if err != nil {
		return err
	}

	if err := d.Set("data_source_references", flattenedDataSourceReferences); err != nil {
		return fmt.Errorf("Error setting data_source_references: %s", err)
	}
	id := fmt.Sprintf("projects/%s/locations/%s/dataSourceReferences", config.BackupDRBasePath, project, location)
	d.SetId(id)
	return nil
}

func flattenDataSourceReferencesList(items []interface{}) ([]map[string]interface{}, error) {
	references := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot cast item to map[string]interface{}")
		}

		ref := map[string]interface{}{
			"name":                data["name"],
			"data_source":         data["dataSource"],
			"backup_config_state": data["dataSourceBackupConfigState"],
		}

		// The API returns backup count as a string, so we parse it to an integer.
		if v, ok := data["dataSourceBackupCount"].(string); ok {
			if i, err := strconv.Atoi(v); err == nil {
				ref["backup_count"] = i
			}
		}

		// Flatten the nested dataSourceBackupConfigInfo object.
		if configInfo, ok := data["dataSourceBackupConfigInfo"].(map[string]interface{}); ok {
			ref["last_backup_state"] = configInfo["lastBackupState"]
			ref["last_successful_backup_time"] = configInfo["lastSuccessfulBackupConsistencyTime"]
		}

		if resourceInfo, ok := data["dataSourceGcpResourceInfo"].(map[string]interface{}); ok {
			ref["gcp_resource_name"] = resourceInfo["gcpResourcename"]
			ref["resource_type"] = resourceInfo["type"]
		}

		references = append(references, ref)
	}
	return references, nil
}
