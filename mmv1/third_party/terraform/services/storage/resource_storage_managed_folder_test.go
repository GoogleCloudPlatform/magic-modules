package storage_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/storage"
)

func TestAccStorageManagedFolder_storageManagedFolderUpdate(t *testing.T) {
	t.Parallel()
	bucketName := fmt.Sprintf("tf-test-managed-folder-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolder_bucket(bucketName) + testAccStorageManagedFolder_managedFolder(false),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config:      testAccStorageManagedFolder_bucket(bucketName),
				ExpectError: regexp.MustCompile(`Error 409: The managed folder you tried to delete is not empty.`),
			},
			{
				Config: testAccStorageManagedFolder_bucket(bucketName) + testAccStorageManagedFolder_managedFolder(true),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config: testAccStorageManagedFolder_bucket(bucketName),
			},
		},
	})
}

func TestAccStorageManagedFolder_rapidCacheConfig(t *testing.T) {
	t.Parallel()
	sharedCache := storage.BootstrapSharedTestAnywhereCache(t, "managed-folder", "US-EAST1", "us-east1-b")
	folderName := fmt.Sprintf("datasets-%s/", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolder_rapidCacheConfig(sharedCache.Bucket, folderName, sharedCache.AnywhereCacheID, "enabled"),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config: testAccStorageManagedFolder_rapidCacheConfig(sharedCache.Bucket, folderName, sharedCache.AnywhereCacheID, "unspecified"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_storage_managed_folder.folder", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config: testAccStorageManagedFolder_rapidCacheConfigCleared(sharedCache.Bucket, folderName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_storage_managed_folder.folder", plancheck.ResourceActionUpdate),
					},
				},
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

func testAccStorageManagedFolder_bucket(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name       = "managed/folder/name/file.txt"
  content    = "This file will affect the folder being deleted if allowNonEmpty=false"
  bucket     = google_storage_bucket.bucket.name
}
`, bucketName)
}

func testAccStorageManagedFolder_managedFolder(forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_managed_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "managed/folder/name/"
  force_destroy = %t
}
`, forceDestroy)
}

func testAccStorageManagedFolder_rapidCacheConfig(bucketName, folderName, rapidCacheID, ingestOnWrite string) string {
	return fmt.Sprintf(`
resource "google_storage_managed_folder" "folder" {
  bucket        = "%s"
  name          = "%s"
  force_destroy = true
  rapid_cache_config {
    policies {
      rapid_cache_id  = "%s"
      ingest_on_write = "%s"
    }
  }
}
`, bucketName, folderName, rapidCacheID, ingestOnWrite)
}

func testAccStorageManagedFolder_rapidCacheConfigCleared(bucketName, folderName string) string {
	return fmt.Sprintf(`
resource "google_storage_managed_folder" "folder" {
  bucket        = "%s"
  name          = "%s"
  force_destroy = true
}
`, bucketName, folderName)
}
