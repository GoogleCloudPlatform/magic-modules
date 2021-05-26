package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecretManagerSecret() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceSecretManagerSecret().Schema)
	addRequiredFieldsToSchema(dsSchema, "secret_id")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceSecretManagerSecretRead,
		Schema: dsSchema,
	}
}

func dataSourceSecretManagerSecretRead(d *schema.ResourceData, meta interface{}) error {

	return resourceSecretManagerSecretRead(d, meta)

}
