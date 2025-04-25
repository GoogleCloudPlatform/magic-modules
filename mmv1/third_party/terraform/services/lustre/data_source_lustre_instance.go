package lustre

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func DataSourceLustreInstance() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceLustreInstance().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceLustreInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceLustreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*tpgresource.Config)

	// Set the ID
	d.SetId(fmt.Sprintf("projects/%s/locations/%s/instances/%s", d.Get("project"), d.Get("region"), d.Get("name")))

	err := resourceLustreInstanceRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", d.Id())
	}

	return nil
}
