package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeRegionNetworkPolicy_regionNetworkPolicyUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionNetworkPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionNetworkPolicy_regionNetworkPolicyCreate(context),
			},
			{
				ResourceName:            "google_compute_region_network_policy.policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccComputeRegionNetworkPolicy_regionNetworkPolicyCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_policy" "policy" {
  name = "tf-test-tf-test-policy%{random_suffix}"
  description = "Terraform test"
}
`, context)
}

func testAccComputeRegionNetworkPolicy_regionNetworkPolicyUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_policy" "policy" {
  name = "tf-test-tf-test-policy%{random_suffix}"
  description = "Terraform test update"
}
`, context)
}
