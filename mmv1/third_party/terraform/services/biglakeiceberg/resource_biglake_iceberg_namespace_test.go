// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package biglakeiceberg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBiglakeIcebergIcebergNamespace_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBiglakeIcebergIcebergNamespaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeIcebergIcebergNamespace_basic(context),
			},
			{
				ResourceName:            "google_biglake_iceberg_namespace.my_namespace",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog"},
			},
			{
				Config: testAccBiglakeIcebergIcebergNamespace_update(context),
			},
			{
				ResourceName:            "google_biglake_iceberg_namespace.my_namespace",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog"},
			},
		},
	})
}

func testAccCheckBiglakeIcebergIcebergNamespaceDestroyProducer(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}

func testAccBiglakeIcebergIcebergNamespace_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "tf-test-iceberg-ns-%{random_suffix}"
  location      = "us-central1"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_biglake_iceberg_catalog" "catalog" {
    name = "tf-test-catalog-%{random_suffix}"
    catalog_type = "CATALOG_TYPE_GCS_BUCKET"
    depends_on = [
      google_storage_bucket.bucket
    ]
}

resource "google_biglake_iceberg_namespace" "my_namespace" {
  catalog   = google_biglake_iceberg_catalog.catalog.name
  namespace = ["accounting", "tax"]
  properties = {
    owner = "Hank"
  }
}
`, context)
}

func testAccBiglakeIcebergIcebergNamespace_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "tf-test-iceberg-ns-%{random_suffix}"
  location      = "us-central1"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_biglake_iceberg_catalog" "catalog" {
    name = "tf-test-catalog-%{random_suffix}"
    catalog_type = "CATALOG_TYPE_GCS_BUCKET"
    depends_on = [
      google_storage_bucket.bucket
    ]
}

resource "google_biglake_iceberg_namespace" "my_namespace" {
  catalog   = google_biglake_iceberg_catalog.catalog.name
  namespace = ["accounting", "tax"]
  properties = {
    owner = "Hank"
    dept  = "finance"
  }
}
`, context)
}