package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeHaVpnGateway() *schema.Resource {
	dsSchema := DatasourceSchemaFromResourceSchema(ResourceComputeHaVpnGateway().Schema)

	// Set 'Required' schema elements
	AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	AddOptionalFieldsToSchema(dsSchema, "project")
	AddOptionalFieldsToSchema(dsSchema, "region")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeHaVpnGatewayRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeHaVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Get("name").(string)

	project, err := GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := GetRegion(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/vpnGateways/%s", project, region, name))

	return resourceComputeHaVpnGatewayRead(d, meta)
}
