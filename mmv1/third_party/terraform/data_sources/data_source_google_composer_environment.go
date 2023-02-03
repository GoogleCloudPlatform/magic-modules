package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComposerEnvironment() *schema.Resource {
	dsSchema := DatasourceSchemaFromResourceSchema(resourceComposerEnvironment().Schema)

	// Set 'Required' schema elements
	AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleComposerEnvironmentRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComposerEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := GetRegion(d, config)
	if err != nil {
		return err
	}
	envName := d.Get("name").(string)

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/environments/%s", project, region, envName))

	return resourceComposerEnvironmentRead(d, meta)
}
