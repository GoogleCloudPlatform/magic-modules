package backupdr

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceBackupDRDataSource() *schema.Resource {
	dsSchema := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"createTime": {
			Type:     schema.TypeString,
			Required: true,
		},
		"updateTime": {
			Type:     schema.TypeString,
			Required: true,
		},
		"backupCount": {
			Type:     schema.TypeString,
			Required: true,
		},
		"etag": {
			Type:     schema.TypeString,
			Required: true,
		},
		"state": {
			Type:     schema.TypeString,
			Required: true,
		},
		"totalStoredBytes": {
			Type:     schema.TypeString,
			Required: true,
		},
		"dataSourceBackupApplianceApplication": {
			Type:     schema.TypeMap,
			Required: true,
		},
	}
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "location", "project", "dataSourceId", "backupVaultId")

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
	url, err := tpgresource.ReplaceVars(d, config, "{{BackupDRBasePath}}projects/{{project}}/locations/{{location}}/backupVaults/{{backupVaultId}}/dataSources/{{dataSourceId}}")
	if err != nil {
		return err
	}
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error reading BackupVault: %s", err)
	}

	if err := d.Set("name", flattenBackupDRDataSourceName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading DataSource: %s", err)
	}

	return nil
}

func flattenBackupDRDataSourceName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}
