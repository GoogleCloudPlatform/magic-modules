package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var loggingProjectBucketConfigSchema = map[string]*schema.Schema{
	"project": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		StateFunc: func(val interface{}) string {
			project := val.(string)
			if !strings.HasPrefix(project, "project") {
				project = "projects/" + project
			}
			return project
		},
	},
}

func projectBucketConfigID(d *schema.ResourceData, config *Config) (string, error) {
	project := d.Get("project").(string)
	location := d.Get("location").(string)
	bucketID := d.Get("bucket_id").(string)

	if !strings.HasPrefix(project, "project") {
		project = "projects/" + project
	}

	id := fmt.Sprintf("%s/locations/%s/buckets/%s", project, location, bucketID)
	return id, nil
}

// Create Logging Bucket config
func ResourceLoggingProjectBucketConfig() *schema.Resource {
	return ResourceLoggingBucketConfig("project", loggingProjectBucketConfigSchema, projectBucketConfigID)
}
