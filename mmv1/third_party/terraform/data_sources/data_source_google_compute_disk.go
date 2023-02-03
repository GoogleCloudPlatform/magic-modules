package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeDisk() *schema.Resource {

	dsSchema := DatasourceSchemaFromResourceSchema(resourceComputeDisk().Schema)
	AddRequiredFieldsToSchema(dsSchema, "name")
	AddOptionalFieldsToSchema(dsSchema, "project")
	AddOptionalFieldsToSchema(dsSchema, "zone")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeDiskRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeDiskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := ReplaceVars(d, config, "projects/{{project}}/zones/{{zone}}/disks/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceComputeDiskRead(d, meta)
}
