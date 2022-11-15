package google

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleStorageBucket() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceStorageBucket().Schema)

	addRequiredFieldsToSchema(dsSchema, "name")

	return &schema.Resource{
		Read:   dataSourceGoogleStorageBucketRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleStorageBucketRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	// Get the bucket and acl
	bucket := d.Get("name").(string)

	res, err := config.NewStorageClient(userAgent).Buckets.Get(bucket).Do()
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Read bucket %v at location %v\n\n", res.Name, res.SelfLink)

	setStorageBucket(d, config, res, bucket, userAgent)

	d.Set("project", strconv.Itoa(int(res.ProjectNumber)))
	return nil
}
