package networkservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkServicesLBPolicies_update(t *testing.T) {
	t.Parallel()

	gatewayName := fmt.Sprintf("tf-test-gateway-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesServiceLbPoliciesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesLBPolicies_basic(gatewayName),
			},
			{
				ResourceName:            "google_network_services_service_lb_policies.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesLBPolicies_update(gatewayName),
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

func testAccNetworkServicesLBPolicies_basic(gatewayName string) string {
	return fmt.Sprintf(`
resource "google_network_services_service_lb_policies" "foobar" {
  name        = "%s"
  location    = "global"
  description = "my description"
}
`, gatewayName)
}

func testAccNetworkServicesLBPolicies_update(gatewayName string) string {
	return fmt.Sprintf(`
resource "google_network_services_service_lb_policies" "foobar" {
  name                     = "%s"
  location                 = "global"
  description              = "my description"
  load_balancing_algorithm = "SPRAY_TO_REGION"
  
  labels = {
    foo = "bar"
  }
}
`, gatewayName)
}
