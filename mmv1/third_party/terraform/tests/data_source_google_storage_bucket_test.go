package google_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	google "internal/terraform-provider-google"
)

func TestAccDataSourceGoogleStorageBucket_basic(t *testing.T) {
	t.Parallel()

	bucket := "tf-bucket-" + google.RandString(t, 10)

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { google.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: google.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStorageBucketConfig(bucket),
				Check: resource.ComposeTestCheckFunc(
					google.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_storage_bucket.bar", "google_storage_bucket.foo", map[string]struct{}{"force_destroy": {}}),
				),
			},
		},
	})
}

func testAccDataSourceGoogleStorageBucketConfig(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "foo" {
  name     = "%s"
  location = "US"
}

data "google_storage_bucket" "bar" {
  name = google_storage_bucket.foo.name
  depends_on = [
    google_storage_bucket.foo,
  ]
}
`, bucketName)
}
