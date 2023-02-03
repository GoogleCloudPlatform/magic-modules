package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleBeyondcorpAppConnection() *schema.Resource {

	dsSchema := DatasourceSchemaFromResourceSchema(resourceBeyondcorpAppConnection().Schema)

	AddRequiredFieldsToSchema(dsSchema, "name")

	AddOptionalFieldsToSchema(dsSchema, "project")
	AddOptionalFieldsToSchema(dsSchema, "region")

	return &schema.Resource{
		Read:   dataSourceGoogleBeyondcorpAppConnectionRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleBeyondcorpAppConnectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Get("name").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := GetRegion(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/appConnections/%s", project, region, name))

	return resourceBeyondcorpAppConnectionRead(d, meta)
}
