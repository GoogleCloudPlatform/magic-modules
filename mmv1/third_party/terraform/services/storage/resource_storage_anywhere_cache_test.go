// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccStorageAnywhereCache_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckStorageAnywhereCacheDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccStorageAnywhereCache_full(context),
			},
			{
				ResourceName:            "google_storage_anywhere_cache.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket"},
			},
			{
				Config: testAccStorageAnywhereCache_update(context),
			},
			{
				ResourceName:            "google_storage_anywhere_cache.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket"},
			},
		},
	})
}

func testAccStorageAnywhereCache_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
}

resource "google_storage_bucket_iam_binding" "binding" {
  bucket = google_storage_bucket.bucket.name
  role = "roles/storage.admin"
  members = [
    "allUsers",
  ]
  depends_on = [ google_storage_bucket.bucket ]
}

resource "time_sleep" "wait_4000_seconds" {
  depends_on = [google_storage_bucket.bucket]
  destroy_duration = "4000s"
}

resource "google_storage_anywhere_cache" "cache" {
    bucket = google_storage_bucket.bucket.name
    zone = "us-central1-f"
	admission_policy = "admit-on-first-miss"
	ttl = "90000s"
    depends_on = [ google_storage_bucket_iam_binding.binding ]
}
`, context)
}

func testAccStorageAnywhereCache_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
}

resource "google_storage_bucket_iam_binding" "binding" {
  bucket = google_storage_bucket.bucket.name
  role = "roles/storage.admin"
  members = [
    "allUsers",
  ]
  depends_on = [ google_storage_bucket.bucket ]
}

resource "time_sleep" "wait_4000_seconds" {
  depends_on = [google_storage_bucket.bucket]
  destroy_duration = "4000s"
}

resource "google_storage_anywhere_cache" "cache" {
    bucket = google_storage_bucket.bucket.name
    zone = "us-central1-f"
	admission_policy = "admit-on-second-miss"
	ttl = "100000s"
    depends_on = [ google_storage_bucket_iam_binding.binding ]
}
`, context)
}
