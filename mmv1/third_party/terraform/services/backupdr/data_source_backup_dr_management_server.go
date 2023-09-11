package backupdr

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudBackupDRService() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceBackupDRManagementServer().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "project","location")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudBackupDRServiceRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudBackupDRServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/managementServers")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceBackupDRManagementServerRead(d, meta)
}