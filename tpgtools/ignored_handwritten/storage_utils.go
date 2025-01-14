package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	storage "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/storage"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func forceDestroyBucketObjects(d *schema.ResourceData, config *transport_tpg.Config, bucket *storage.Bucket) error {
	return nil
}
