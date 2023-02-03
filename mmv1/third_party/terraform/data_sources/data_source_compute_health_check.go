package google

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceGoogleComputeHealthCheck() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := DatasourceSchemaFromResourceSchema(resourceComputeHealthCheck().Schema)

	// Set 'Required' schema elements
	AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeHealthCheckRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	id, err := ReplaceVars(d, meta.(*Config), "projects/{{project}}/global/healthChecks/{{name}}")
	if err != nil {
		return err
	}
	d.SetId(id)

	return resourceComputeHealthCheckRead(d, meta)
}
