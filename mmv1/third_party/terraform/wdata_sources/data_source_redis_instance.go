package google

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSourceGoogleRedisInstance() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := DatasourceSchemaFromResourceSchema(ResourceRedisInstance().Schema)

	// Set 'Required' schema elements
	AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleRedisInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleRedisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	id, err := ReplaceVars(d, meta.(*Config), "projects/{{project}}/locations/{{region}}/instances/{{name}}")
	if err != nil {
		return err
	}
	d.SetId(id)

	return resourceRedisInstanceRead(d, meta)
}
