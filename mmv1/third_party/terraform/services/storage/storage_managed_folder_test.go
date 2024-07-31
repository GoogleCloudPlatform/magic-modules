package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageManagedFolder_storageManagedFolderUpdate(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-managed-folder-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolder_storageManagedFolderUpdate(name, false),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config: testAccStorageManagedFolder_storageManagedFolderUpdate(name, true),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
		},
	})
}

func testAccStorageManagedFolder_storageManagedFolderUpdate(name string, forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_managed_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "managed/folder/name/"
  force_destroy = %t
}
`, name, forceDestroy)
}
