package apigee

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func DataSourceApigeeInstance() *schema.Resource {
	// Inherit schema from the resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceApigeeInstance().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name", "org_id")

	return &schema.Resource{
		Read:   dataSourceApigeeInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceApigeeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	orgId := d.Get("org_id").(string)
	name := d.Get("name").(string)

	instancePath := fmt.Sprintf("%s/instances/%s", orgId, name)
	d.SetId(instancePath)

	if err := resourceApigeeInstanceRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("Apigee Instance %q not found", instancePath)
	}

	// NOTE: If the resource had labels or annotations, we would call:
	// tpgresource.SetDataSourceLabels(d) here

	return nil
}
