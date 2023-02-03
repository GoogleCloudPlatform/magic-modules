package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeSslPolicy() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := DatasourceSchemaFromResourceSchema(resourceComputeSslPolicy().Schema)

	// Set 'Required' schema elements
	AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   datasourceComputeSslPolicyRead,
		Schema: dsSchema,
	}
}

func datasourceComputeSslPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	policyName := d.Get("name").(string)

	d.SetId(fmt.Sprintf("projects/%s/global/sslPolicies/%s", project, policyName))

	return resourceComputeSslPolicyRead(d, meta)
}
