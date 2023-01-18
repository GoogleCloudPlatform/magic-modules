package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeGlobalForwardingRule() *schema.Resource {
	dsSchema := DatasourceSchemaFromResourceSchema(ResourceComputeGlobalForwardingRule().Schema)

	// Set 'Required' schema elements
	AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeGlobalForwardingRuleRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeGlobalForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Get("name").(string)

	project, err := GetProject(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/global/forwardingRules/%s", project, name))

	return resourceComputeGlobalForwardingRuleRead(d, meta)
}
