package composerservice

import (
	"fmt"

	google "terraform-provider-google/internal"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComposerEnvironment() *schema.Resource {
	dsSchema := google.DatasourceSchemaFromResourceSchema(ResourceComposerEnvironment().Schema)

	// Set 'Required' schema elements
	google.AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	google.AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleComposerEnvironmentRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComposerEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*google.Config)
	project, err := google.GetProject(d, config)
	if err != nil {
		return err
	}
	region, err := google.GetRegion(d, config)
	if err != nil {
		return err
	}
	envName := d.Get("name").(string)

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/environments/%s", project, region, envName))

	return resourceComposerEnvironmentRead(d, meta)
}
