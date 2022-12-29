package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleComputeFirewall_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudFirewall_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_compute_firewall.foo", "google_compute_firewall.default"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudFirewall_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_firewall" "default" {
  name    = "my-firewall%{random_suffix}"
  network = google_compute_network.default.name

  allow {
   protocol = "icmp"
  }

  allow {
   protocol = "tcp"
   ports    = ["80", "8080", "1000-2000"]
  }

  source_tags = ["web"]
}

resource "google_compute_network" "default" {
  name = "my-network%{random_suffix}"
}
	
data "google_compute_firewall" "foo" {
  name = google_compute_firewall.default.name
  project = google_compute_firewall.default.project
}`, context)

}
