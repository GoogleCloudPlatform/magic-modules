package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeRouter() *schema.Resource {
	dsSchema := DatasourceSchemaFromResourceSchema(ResourceComputeRouter().Schema)
	AddRequiredFieldsToSchema(dsSchema, "name")
	AddRequiredFieldsToSchema(dsSchema, "network")
	AddOptionalFieldsToSchema(dsSchema, "region")
	AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceComputeRouterRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeRouterRead(d *schema.ResourceData, meta interface{}) error {
	routerName := d.Get("name").(string)

	d.SetId(routerName)
	return resourceComputeRouterRead(d, meta)
}
