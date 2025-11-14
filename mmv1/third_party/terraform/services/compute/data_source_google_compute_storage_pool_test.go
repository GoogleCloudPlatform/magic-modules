// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeStoragePool_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeStoragePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeStoragePool_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_storage_pool.test-storage-pool", "google_compute_storage_pool.my-storage-pool-data"),
				),
			},
		},
	})
}
func testAccDataSourceComputeStoragePool_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "test-storage-pool" {
  name                         = "tf-test-storage-pool-%{random_suffix}"
  zone                         = "us-central1-a"
  pool_provisioned_capacity_gb = "10240"
  pool_provisioned_throughput  = "1024"
  storage_pool_type            = "hyperdisk-throughput"
}

data "google_compute_storage_pool" "my-storage-pool-data" {
  name = google_compute_storage_pool.test-storage-pool.name
  zone = "us-central1-a"
}
`, context)
}
