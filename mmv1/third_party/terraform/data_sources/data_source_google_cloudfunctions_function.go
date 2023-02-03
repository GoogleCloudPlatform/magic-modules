package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleCloudFunctionsFunction() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := DatasourceSchemaFromResourceSchema(resourceCloudFunctionsFunction().Schema)

	// Set 'Required' schema elements
	AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudFunctionsFunctionRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudFunctionsFunctionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := GetRegion(d, config)
	if err != nil {
		return err
	}

	cloudFuncId := &cloudFunctionId{
		Project: project,
		Region:  region,
		Name:    d.Get("name").(string),
	}

	d.SetId(cloudFuncId.cloudFunctionId())

	err = resourceCloudFunctionsRead(d, meta)
	if err != nil {
		return err
	}

	return nil
}
