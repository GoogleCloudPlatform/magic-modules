package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleKmsKeyRing() *schema.Resource {
	dsSchema := DatasourceSchemaFromResourceSchema(ResourceKMSKeyRing().Schema)
	AddRequiredFieldsToSchema(dsSchema, "name")
	AddRequiredFieldsToSchema(dsSchema, "location")
	AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleKmsKeyRingRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleKmsKeyRingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := GetProject(d, config)
	if err != nil {
		return err
	}

	keyRingId := kmsKeyRingId{
		Name:     d.Get("name").(string),
		Location: d.Get("location").(string),
		Project:  project,
	}
	d.SetId(keyRingId.keyRingId())

	return resourceKMSKeyRingRead(d, meta)
}
