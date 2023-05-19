package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeNetworkEdgeSecurityService_update(t *testing.T) {
	t.Parallel()

	pName := fmt.Sprintf("tf-test-security-policy-%s", RandString(t, 10))
	nesName := fmt.Sprintf("tf-test-edge-security-services-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputenetworkEdgeSecurityServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkEdgeSecurityService_basic(pName, nesName),
			},
			{
				ResourceName:      "google_compute_network_edge_security_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkEdgeSecurityService_update(pName, nesName),
			},
			{
				ResourceName:      "google_compute_network_edge_security_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkEdgeSecurityService_basic(pName, nesName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_security_policy" "foobar" {
  name        = "%s"
  description = "basic region security policy"
  type        = "CLOUD_ARMOR_NETWORK"
  region      = "asia-southeast1"
}

resource "google_compute_network_edge_security_service" "foobar" {
  name     = "%s"
  region = "asia-southeast1"
  description  = "My basic resource using security policy"
  security_policy = google_compute_region_security_policy.foobar.self_link
}
`, pName, nesName)
}

func testAccNetworkEdgeSecurityService_update(pName, nesName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_security_policy" "foobar" {
  name        = "%s"
  description = "basic region security policy"
  type        = "CLOUD_ARMOR_NETWORK"
  region      = "asia-southeast1"
}

resource "google_compute_network_edge_security_service" "foobar" {
  name     = "%s"
  region = "asia-southeast1"
  description  = "My basic updated resource using security policy"
  security_policy = google_compute_region_security_policy.foobar.self_link
}
`, pName, nesName)
}
