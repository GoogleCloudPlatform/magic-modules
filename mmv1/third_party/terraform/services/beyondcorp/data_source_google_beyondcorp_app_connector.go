package beyondcorp

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleBeyondcorpAppConnector() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceBeyondcorpAppConnector().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "region")

	return &schema.Resource{
		Read:   dataSourceGoogleBeyondcorpAppConnectorRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleBeyondcorpAppConnectorRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	name := d.Get("name").(string)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("projects/%s/locations/%s/appConnectors/%s", project, region, name)
	d.SetId(id)

	err = resourceBeyondcorpAppConnectorRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
