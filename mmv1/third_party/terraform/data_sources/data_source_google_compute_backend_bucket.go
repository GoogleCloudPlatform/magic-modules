package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeBackendBucket() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(ResourceComputeBackendBucket().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceComputeBackendBucketRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeBackendBucketRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	backendBucketName := d.Get("name").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/global/backendBuckets/%s", project, backendBucketName))

	return resourceComputeBackendBucketRead(d, meta)
}
