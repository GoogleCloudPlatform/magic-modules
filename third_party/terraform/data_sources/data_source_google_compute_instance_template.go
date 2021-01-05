package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeInstanceTemplate() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeInstanceTemplate().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "project")

	// Set 'Optional' schema elements
	// TODO: add/handle filter
	addOptionalFieldsToSchema(dsSchema, "name")

	return &schema.Resource{
		Read:   datasourceComputeInstanceTemplateRead,
		Schema: dsSchema,
	}
}

func datasourceComputeInstanceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Get("name").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.SetId("projects/" + project + "/global/instanceTemplates/" + name)

	return resourceComputeInstanceTemplateRead(d, meta)
}
