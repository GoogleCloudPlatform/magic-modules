package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeStoragePool_computeStoragePool_update(t *testing.T) {
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
				Config: testAccComputeStoragePool_computeStoragePoolFullExample(context),
			},
			{
				ResourceName:            "google_compute_storage_pool.test-storage-pool-full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "labels", "terraform_labels", "zone"},
			},
			{
				Config: testAccComputeStoragePool_computeStoragePool_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_storage_pool.test-storage-pool-full", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_compute_storage_pool.test-storage-pool-full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "labels", "terraform_labels", "zone"},
			},
		},
	})
}

func testAccComputeStoragePool_computeStoragePool_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "test-storage-pool-full" {
  name                     = "tf-test-storage-pool-full%{random_suffix}"

  capacity_provisioning_type   = "STANDARD"
  pool_provisioned_capacity_gb = "11264"

  performance_provisioning_type = "STANDARD"
  pool_provisioned_iops         = "20000"
  pool_provisioned_throughput   = "2048"

  storage_pool_type = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.project_id}/zones/us-central1-a/storagePoolTypes/hyperdisk-balanced"

	deletion_protection = false
}

data "google_project" "project" {}
`, context)
}
