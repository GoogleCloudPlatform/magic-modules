package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleProjectService() *schema.Resource {

	dsSchema := DatasourceSchemaFromResourceSchema(resourceGoogleProjectService().Schema)
	AddRequiredFieldsToSchema(dsSchema, "service")
	AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleProjectServiceRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleProjectServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := ReplaceVars(d, config, "{{project}}/{{service}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceGoogleProjectServiceRead(d, meta)
}
