package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeHaVpnGateway() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(ResourceComputeHaVpnGateway().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")
	addOptionalFieldsToSchema(dsSchema, "region")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeHaVpnGatewayRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeHaVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	name := d.Get("name").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/vpnGateways/%s", project, region, name))

	return resourceComputeHaVpnGatewayRead(d, meta)
}
