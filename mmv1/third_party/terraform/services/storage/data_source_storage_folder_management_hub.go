package storage

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleStorageFolderManagementHub() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceStorageFolderManagementHub().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	return &schema.Resource{
		Read:   dataSourceGoogleStorageFolderManagementHubRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleStorageFolderManagementHubRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "folders/{{name}}/locations/global/managementHub")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	err = resourceStorageFolderManagementHubRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
