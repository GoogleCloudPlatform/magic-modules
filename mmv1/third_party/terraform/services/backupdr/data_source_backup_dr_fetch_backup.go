package backupdr

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudBackupDRFetchBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleCloudBackupDRFetchBackupsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the project in which the resource belongs.",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location of the backups.",
			},
			"backup_vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Backup Vault.",
			},
			"data_source_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Data Source.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the resource to fetch backups for (e.g., compute.googleapis.com/Instance).",
			},
			"backups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of backups matching the resource type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource name of the backup.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the backup was created.",
						},
						"consistency_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The point in time when the backup was taken.",
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleCloudBackupDRFetchBackupsRead(d *schema.ResourceData, meta interface{}) error {
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
	backupVaultId := d.Get("backup_vault_id").(string)
	dataSourceId := d.Get("data_source_id").(string)
	resourceType := d.Get("resource_type").(string)

	params := make(map[string]string)
	params["resourceType"] = resourceType
	var allItems []interface{}

	for {
		url, err := tpgresource.ReplaceVars(d, config, "{{BackupDRBasePath}}projects/{{project}}/locations/{{location}}/backupVaults/{{backup_vault_id}}/dataSources/{{data_source_id}}/backups:fetchForResourceType")
		if err != nil {
			return err
		}

		url, err = transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return fmt.Errorf("Error fetching backups for resource type: %s", err)
		}

		items, ok := res["backups"].([]interface{})
		if ok {
			allItems = append(allItems, items...)
		}

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	flattenedBackups, err := flattenFetchBackups(allItems)
	if err != nil {
		return err
	}

	if err := d.Set("backups", flattenedBackups); err != nil {
		return fmt.Errorf("Error setting backups: %s", err)
	}

	id := fmt.Sprintf("projects/%s/locations/%s/backupVaults/%s/dataSources/%s/backupTypes/%s", project, location, backupVaultId, dataSourceId, resourceType)
	d.SetId(id)
	return nil
}

func flattenFetchBackups(items []interface{}) ([]map[string]interface{}, error) {
	backups := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot cast item to map[string]interface{}")
		}
		backup := map[string]interface{}{
			"name":             data["name"],
			"create_time":      data["createTime"],
			"consistency_time": data["consistencyTime"],
		}
		backups = append(backups, backup)
	}
	return backups, nil
}
