package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageFolder_storageFolderBasic(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageFolder_storageBucket(bucketName, true) + testAccStorageFolder_storageFolder(true),
			},
			{
				ResourceName:            "google_storage_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "recursive", "force_destroy"},
			},
		},
	})
}

func testAccStorageFolder_storageBucket(bucketName string, forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "EU"
  uniform_bucket_level_access = true
  hierarchical_namespace {
	enabled = true
  }
  force_destroy = %t
}
`, bucketName, forceDestroy)
}

func testAccStorageFolder_storageFolder(forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "folder/name/"
  recursive     = true 
  force_destroy = %t
}
`, forceDestroy)
}
