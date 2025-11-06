package backupdr

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/googleapi"
)

func DataSourceGoogleCloudBackupDRDataSourceReference() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleCloudBackupDRDataSourceReferenceRead,
		Schema: map[string]*schema.Schema{
			"data_source_reference_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The `id` of the data source reference.",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location of the data source reference.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the project in which the resource belongs.",
			},
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
	}
}

func dataSourceGoogleCloudBackupDRDataSourceReferenceRead(d *schema.ResourceData, meta interface{}) error {
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
	dataSourceReferenceId := d.Get("data_source_reference_id").(string)
	url := fmt.Sprintf("%sprojects/%s/locations/%s/dataSourceReferences/%s", config.BackupDRBasePath, project, location, dataSourceReferenceId)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			d.SetId("") // Resource not found
			return nil
		}
		return fmt.Errorf("Error reading DataSourceReference: %s", err)
	}

	if err := flattenDataSourceReference(d, res); err != nil {
		return err
	}

	d.SetId(res["name"].(string))
	return nil
}

func flattenDataSourceReference(d *schema.ResourceData, data map[string]interface{}) error {
	ref, err := flattenDataSourceReferenceToMap(data)
	if err != nil {
		return err
	}
	for k, v := range ref {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("Error setting %s: %s", k, err)
		}
	}
	return nil
}

func flattenDataSourceReferenceToMap(data map[string]interface{}) (map[string]interface{}, error) {
	ref := map[string]interface{}{
		"name":                data["name"],
		"data_source":         data["dataSource"],
		"backup_config_state": data["dataSourceBackupConfigState"],
	}
	if v, ok := data["dataSourceBackupCount"].(string); ok {
		if i, err := strconv.Atoi(v); err == nil {
			ref["backup_count"] = i
		}
	}
	if configInfo, ok := data["dataSourceBackupConfigInfo"].(map[string]interface{}); ok {
		ref["last_backup_state"] = configInfo["lastBackupState"]
		ref["last_successful_backup_time"] = configInfo["lastSuccessfulBackupConsistencyTime"]
	}
	if resourceInfo, ok := data["dataSourceGcpResourceInfo"].(map[string]interface{}); ok {
		ref["gcp_resource_name"] = resourceInfo["gcpResourcename"]
		ref["resource_type"] = resourceInfo["type"]
	}
	return ref, nil
}
