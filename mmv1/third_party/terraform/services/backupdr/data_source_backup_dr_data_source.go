package backupdr

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceBackupDRDataSource() *schema.Resource {
	dsSchema := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"create_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"update_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"backup_count": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"etag": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"total_stored_bytes": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"data_source_backup_appliance_application": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"application_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"backup_appliance": {
						Type:     schema.TypeString,
						Required: true,
					},
					"appliance_id": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"application_id": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"hostname": {
						Type:     schema.TypeString,
						Required: true,
					},
					"host_id": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		"data_source_gcp_resource": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"gcp_resourcename": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"location": {
						Type:     schema.TypeString,
						Required: true,
					},
					"type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"compute_instance_data_source_properties": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Required: true,
								},
								"description": {
									Type:     schema.TypeString,
									Required: true,
								},
								"machine_type": {
									Type:     schema.TypeString,
									Required: true,
								},
								"total_disk_count": {
									Type:     schema.TypeInt,
									Required: true,
								},
								"total_disk_size_gb": {
									Type:     schema.TypeInt,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		"location": {
			Type:     schema.TypeString,
			Required: true,
		},
		"project": {
			Type:     schema.TypeString,
			Required: true,
		},
		"data_source_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"backup_vault_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	log.Printf("Schema declared")

	return &schema.Resource{
		Read:   DataSourceBackupDRDataSourceRead,
		Schema: dsSchema,
	}
}

func DataSourceBackupDRDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	if len(location) == 0 {
		return fmt.Errorf("Cannot determine location: set location in this data source or at provider-level")
	}

	billingProject := project
	url, err := tpgresource.ReplaceVars(d, config, "{{BackupDRBasePath}}projects/{{project}}/locations/{{location}}/backupVaults/{{backup_vault_id}}/dataSources/{{data_source_id}}")
	log.Printf("url retrieved")
	if err != nil {
		return err
	}
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	log.Printf("project got... going to send request")
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})

	log.Printf("get request initiated")

	if err != nil {
		return fmt.Errorf("Error reading BackupVault: %s", err)
	}

	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error reading DataSource: %s", err)
	}

	if err := d.Set("create_time", res["create_time"]); err != nil {
		return fmt.Errorf("Error reading DataSource: %s", err)
	}

	if err := d.Set("update_time", res["update_time"]); err != nil {
		return fmt.Errorf("Error reading DataSource: %s", err)
	}

	if err := d.Set("backup_count", res["backup_count"]); err != nil {
		return fmt.Errorf("Error reading DataSource: %s", err)
	}

	if err := d.Set("etag", res["etag"]); err != nil {
		return fmt.Errorf("Error reading DataSource: %s", err)
	}

	if err := d.Set("state", res["state"]); err != nil {
		return fmt.Errorf("Error reading DataSource: %s", err)
	}

	if err := d.Set("total_stored_bytes", res["total_stored_bytes"]); err != nil {
		return fmt.Errorf("Error reading DataSource: %s", err)
	}

	d.SetId(res["name"].(string))

	return nil
}
