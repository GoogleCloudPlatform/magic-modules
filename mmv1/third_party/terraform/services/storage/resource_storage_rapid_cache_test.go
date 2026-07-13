package storage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/storage"
)

func TestAccStorageRapidCache_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccStorageRapidCache_full(context),
			},
			{
				ResourceName:            "google_storage_rapid_cache.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "update_time", "cache_type"},
			},
			{
				Config: testAccStorageRapidCache_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_storage_rapid_cache.cache", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_storage_rapid_cache.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "update_time", "cache_type"},
			},
		},
	})
}

func testAccStorageRapidCache_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name = "tf-test-bucket-name%{random_suffix}"
  location = "US"
  force_destroy = "true"
}

resource "google_storage_rapid_cache" "cache" {
  bucket          = google_storage_bucket.bucket.name
  zone            = "us-central1-f"
  cache_type      = "rapid-cache"
  ttl             = "3601s"
  ingest_on_write = false
}
`, context)
}

func testAccStorageRapidCache_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name = "tf-test-bucket-name%{random_suffix}"
  location = "US"
  force_destroy = "true"
}

resource "google_storage_rapid_cache" "cache" {
  bucket           = google_storage_bucket.bucket.name
  zone             = "us-central1-f"
  admission_policy = "no-read-admission"
  cache_type       = "rapid-cache"
  ttl              = "3620s"
  ingest_on_write  = true
}
`, context)
}

func TestAccStorageRapidCache_cacheType(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccStorageRapidCache_full(context),
			},
			{
				ResourceName:            "google_storage_rapid_cache.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "update_time", "cache_type"},
			},
			{
				Config: testAccStorageRapidCache_cacheTypeUltra(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_storage_rapid_cache.cache", plancheck.ResourceActionReplace),
					},
				},
			},
			{
				ResourceName:            "google_storage_rapid_cache.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "update_time", "cache_type"},
			},
		},
	})
}

func testAccStorageRapidCache_cacheTypeUltra(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name = "tf-test-bucket-name%{random_suffix}"
  location = "US"
  force_destroy = "true"
  uniform_bucket_level_access = true
}

resource "google_storage_rapid_cache" "cache" {
  bucket          = google_storage_bucket.bucket.name
  zone            = "us-central1-a"
  cache_type      = "rapid-cache-ultra"
  ttl             = "3601s"
  ingest_on_write = false
}
`, context)
}
