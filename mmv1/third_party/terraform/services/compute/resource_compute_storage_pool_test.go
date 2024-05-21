package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeStoragePool_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"storage_pool_name": fmt.Sprintf("tf-test-storage-pool-%s", acctest.RandString(t, 10)),
		"storage_pool_type": fmt.Sprintf("hyperdisk-throughput"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeStoragePool_basic(context),
			},
			{
				ResourceName:            "google_compute_storage_pool.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_pool_type", "zone"},
			},
			{
				Config: testAccComputeStoragePool_update(context),
			},
			{
				ResourceName:            "google_compute_storage_pool.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_pool_type", "zone"},
				Destroy:                 true,
			},
		},
	})
}

func TestAccComputeStoragePool_fromStoragePoolTypeUrl(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"storage_pool_name": fmt.Sprintf("tf-test-storage-pool-%s", acctest.RandString(t, 10)),
		"storage_pool_type": fmt.Sprintf("projects/%s/zones/us-central1-a/storagePoolTypes/hyperdisk-balanced", envvar.GetTestProjectFromEnv()),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeStoragePool_hyperdiskBalanced(context),
			},
			{
				ResourceName:            "google_compute_storage_pool.hdb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_pool_type", "zone"},
				Destroy:                 true,
			},
		},
	})
}

func testAccComputeStoragePool_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "foobar" {
  name  = "%{storage_pool_name}"
  zone  = "us-central1-a"
  description = "testing storage pool basic"
  storage_pool_type = "%{storage_pool_type}"
  capacity_provisioning_type = "ADVANCED"
  performance_provisioning_type = "STANDARD"
  pool_provisioned_capacity_gb = 10240
  pool_provisioned_throughput = 140
}
`, context)
}

func testAccComputeStoragePool_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "foobar" {
  name  = "%{storage_pool_name}"
  zone  = "us-central1-a"
  description = "testing storage pool update"
  storage_pool_type = "%{storage_pool_type}"
  capacity_provisioning_type = "ADVANCED"
  performance_provisioning_type = "STANDARD"
  pool_provisioned_capacity_gb = 11264
  pool_provisioned_throughput = 120
}
`, context)
}

func testAccComputeStoragePool_hyperdiskBalanced(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "hdb" {
  name  = "%{storage_pool_name}"
  zone  = "us-central1-a"
  storage_pool_type = "%{storage_pool_type}"
  capacity_provisioning_type = "ADVANCED"
  pool_provisioned_capacity_gb = 10240
  pool_provisioned_iops = 10000
  pool_provisioned_throughput = 1024
}
`, context)
}
