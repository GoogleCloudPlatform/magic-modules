package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleCloudFunctions2Function() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := DatasourceSchemaFromResourceSchema(ResourceCloudfunctions2function().Schema)

	// Set 'Required' schema elements
	AddRequiredFieldsToSchema(dsSchema, "name", "location")

	// Set 'Optional' schema elements
	AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudFunctions2FunctionRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudFunctions2FunctionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := GetProject(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/functions/%s", project, d.Get("location").(string), d.Get("name").(string)))

	err = resourceCloudfunctions2functionRead(d, meta)
	if err != nil {
		return err
	}

	return nil
}
