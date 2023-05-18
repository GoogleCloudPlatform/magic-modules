package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeNetworkEdgeSecurityService_update(t *testing.T) {
	t.Parallel()

	nesName := fmt.Sprintf("tf-test-edge-security-services-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputenetworkEdgeSecurityServicesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkEdgeSecurityService_basic(nesName),
			},
			{
				ResourceName:      "google_compute_network_edge_security_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkEdgeSecurityService_basic(nesName string) string {
	return fmt.Sprintf(`
  resource "google_compute_security_policy" "policy" {
    name = "my-policy"
  
    rule {
      action   = "deny(403)"
      priority = "1000"
      match {
        versioned_expr = "SRC_IPS_V1"
        config {
          src_ip_ranges = ["9.9.9.0/24"]
        }
      }
      description = "Deny access to IPs in 9.9.9.0/24"
    }

    rule {
      action   = "allow"
      priority = "2147483647"
      match {
        versioned_expr = "SRC_IPS_V1"
        config {
          src_ip_ranges = ["*"]
        }
      }
      description = "default rule"
    }
  }
  
  resource "google_compute_network_edge_security_services" "foobar" {
    name     = "%s"
    region = "asia-southeast1"
    description  = "My basic resource using security policy"
    security_policy = google_compute_security_policy.policy.id
  }
`, nesName)
}
