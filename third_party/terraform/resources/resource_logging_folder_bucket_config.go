package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var loggingFolderBucketConfigSchema = map[string]*schema.Schema{
	"folder": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		StateFunc: func(val interface{}) string {
			folder := val.(string)
			if !strings.HasPrefix(folder, "folder") {
				folder = "folders/" + folder
			}
			return folder
		},
	},
}

func folderBucketConfigID(d *schema.ResourceData, config *Config) (string, error) {
	folder := d.Get("folder").(string)
	location := d.Get("location").(string)
	bucketID := d.Get("bucket_id").(string)

	if !strings.HasPrefix(folder, "folder") {
		folder = "folders/" + folder
	}

	id := fmt.Sprintf("%s/locations/%s/buckets/%s", folder, location, bucketID)
	return id, nil
}

// Create Logging Bucket config
func ResourceLoggingFolderBucketConfig() *schema.Resource {
	return ResourceLoggingBucketConfig("folder", loggingFolderBucketConfigSchema, folderBucketConfigID)
}
