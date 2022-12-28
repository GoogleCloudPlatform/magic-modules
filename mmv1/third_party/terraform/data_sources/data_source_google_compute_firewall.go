package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeFirewall() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeFirewall().Schema)

	addRequiredFieldsToSchema(dsSchema, "name")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudFirewallRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudFirewallRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := replaceVars(d, config, "projects/{{project}}/global/firewalls/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	d.SetId(id)
	return resourceComputeFirewallRead(d, meta)
}
