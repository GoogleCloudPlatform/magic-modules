package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleRuntimeconfigConfig() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceRuntimeconfigConfig().Schema)
	addRequiredFieldsToSchema(dsSchema, "name")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleRuntimeconfigConfigRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleRuntimeconfigConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := replaceVars(d, config, "projects/{{project}}/configs/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceRuntimeconfigConfigRead(d, meta)
}
