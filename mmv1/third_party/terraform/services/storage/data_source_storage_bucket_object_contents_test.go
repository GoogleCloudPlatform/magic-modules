package storage_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceStorageBucketObjectContents_Basic(t *testing.T) {
	bucket := "tf-bucket-object-contents-" + acctest.RandString(t, 10)

	content1 := "hello world"
	content2 := "goodbye world"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStorageBucketObjectContents_Basic(bucket, content1, content2),
				Check: resource.ComposeTestCheckFunc(
					// ensure list exists
					resource.TestCheckResourceAttrSet(
						"data.google_storage_bucket_object_contents.this",
						"bucket_objects.#",
					),

					// deterministic ordering (alphabetical by name)
					resource.TestCheckResourceAttr(
						"data.google_storage_bucket_object_contents.this",
						"bucket_objects.0.name",
						"object-1",
					),
					resource.TestCheckResourceAttr(
						"data.google_storage_bucket_object_contents.this",
						"bucket_objects.1.name",
						"object-2",
					),

					// content
					resource.TestCheckResourceAttr(
						"data.google_storage_bucket_object_contents.this",
						"bucket_objects.0.content",
						content1,
					),
					resource.TestCheckResourceAttr(
						"data.google_storage_bucket_object_contents.this",
						"bucket_objects.1.content",
						content2,
					),

					// base64
					resource.TestCheckResourceAttr(
						"data.google_storage_bucket_object_contents.this",
						"bucket_objects.0.content_base64",
						base64.StdEncoding.EncodeToString([]byte(content1)),
					),
					resource.TestCheckResourceAttr(
						"data.google_storage_bucket_object_contents.this",
						"bucket_objects.1.content_base64",
						base64.StdEncoding.EncodeToString([]byte(content2)),
					),
				),
			},
		},
	})
}

func testAccDataSourceStorageBucketObjectContents_Basic(bucket, content1, content2 string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "this" {
  force_destroy               = true
  location                    = "US"
  name                        = "%s"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object1" {
  bucket  = google_storage_bucket.this.name
  content = "%s"
  name    = "object-1"
}

resource "google_storage_bucket_object" "object2" {
  bucket  = google_storage_bucket.this.name
  content = "%s"
  name    = "object-2"
}

data "google_storage_bucket_object_contents" "this" {
  bucket = google_storage_bucket.this.name

  depends_on = [
    google_storage_bucket_object.object1,
    google_storage_bucket_object.object2,
  ]
}
`, bucket, content1, content2)
}
