package kms

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func DataSourceGoogleKmsAutokeyConfig() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceKMSAutokeyConfig().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "folder")

	return &schema.Resource{
		Read:   dataSourceGoogleKmsAutokeyConfigRead,
		Schema: dsSchema,
	}

}

func dataSourceGoogleKmsAutokeyConfigRead(d *schema.ResourceData, meta interface{}) error {
	configId := KmsAutokeyConfigId{
		Folder: d.Get("folder").(string),
	}
	id := configId.AutokeyConfigId()
	d.SetId(id)
	err := resourceKMSAutokeyConfigRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
