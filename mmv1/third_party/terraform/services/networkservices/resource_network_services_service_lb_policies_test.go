package networkservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkServicesLBPolicies_update(t *testing.T) {
	t.Parallel()

	policyName := fmt.Sprintf("tf-test-lb-policy-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesServiceLbPoliciesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesLBPolicies_basic(policyName),
			},
			{
				ResourceName:            "google_network_services_service_lb_policies.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesLBPolicies_update(policyName),
			},
			{
				ResourceName:            "google_network_services_service_lb_policies.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkServicesLBPolicies_basic(policyName string) string {
	return fmt.Sprintf(`
resource "google_network_services_service_lb_policies" "foobar" {
  name        = "%s"
  location    = "global"
  description = "my description"
}
`, policyName)
}

func testAccNetworkServicesLBPolicies_update(policyName string) string {
	return fmt.Sprintf(`
resource "google_network_services_service_lb_policies" "foobar" {
  name                     = "%s"
  location                 = "global"
  description              = "my description"
  load_balancing_algorithm = "SPRAY_TO_REGION"

  auto_capacity_drain {
    enable = true
  }

  failover_config {
    failover_health_threshold = 70
  }
  
  labels = {
    foo = "bar"
  }
}
`, policyName)
}
