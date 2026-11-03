package storage_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/storage"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
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
	bucketName := fmt.Sprintf("tf-test-managed-folder-%s", acctest.RandString(t, 10))
	zone := "us-central1-f"
	var opSelfLink string

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolder_rapidCacheBucket(bucketName),
				Check: func(s *terraform.State) error {
					config := acctest.GoogleProviderConfig(t)
					url := fmt.Sprintf("%sb/%s/anywhereCaches/", transport_tpg.BaseUrl(storage.Product, config), bucketName)
					res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
						Config:    config,
						Method:    "POST",
						RawURL:    url,
						UserAgent: config.UserAgent,
						Body: map[string]interface{}{
							"zone": zone,
						},
					})
					if err != nil {
						return fmt.Errorf("error initiating anywhereCache for test: %w", err)
					}
					if sl, ok := res["selfLink"].(string); ok {
						opSelfLink = sl
					}
					return nil
				},
			},
			{
				Config: testAccStorageManagedFolder_rapidCacheConfig(bucketName, zone, "enabled"),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config: testAccStorageManagedFolder_rapidCacheConfig(bucketName, zone, "unspecified"),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config: testAccStorageManagedFolder_rapidCacheConfigCleared(bucketName),
				Check: func(s *terraform.State) error {
					if opSelfLink == "" {
						return nil
					}
					config := acctest.GoogleProviderConfig(t)
					_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
						Config:    config,
						Method:    "POST",
						RawURL:    fmt.Sprintf("%s/cancel", opSelfLink),
						UserAgent: config.UserAgent,
					})
					if err != nil {
						return fmt.Errorf("error cancelling anywhereCache LRO: %w", err)
					}
					cacheURL := fmt.Sprintf("%sb/%s/anywhereCaches/%s", transport_tpg.BaseUrl(storage.Product, config), bucketName, zone)
					for i := 0; i < 30; i++ {
						_, getErr := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
							Config:    config,
							Method:    "GET",
							RawURL:    cacheURL,
							UserAgent: config.UserAgent,
						})
						if getErr != nil {
							break
						}
						time.Sleep(2 * config.SleepUnit)
					}
					return nil
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

func testAccStorageManagedFolder_rapidCacheBucket(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "US"
  uniform_bucket_level_access = true
  force_destroy               = true
}
`, bucketName)
}

func testAccStorageManagedFolder_rapidCacheConfig(bucketName, zone, ingestOnWrite string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "US"
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_storage_managed_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "datasets/"
  force_destroy = true
  rapid_cache_config {
    policies {
      rapid_cache_id  = "%s"
      ingest_on_write = "%s"
    }
  }
}
`, bucketName, zone, ingestOnWrite)
}

func testAccStorageManagedFolder_rapidCacheConfigCleared(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "US"
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_storage_managed_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "datasets/"
  force_destroy = true
}
`, bucketName)
}
