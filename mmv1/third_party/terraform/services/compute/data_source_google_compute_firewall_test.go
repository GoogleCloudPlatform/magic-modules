package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

func TestAccComputeFirewallDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewallDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_firewall.default", "google_compute_firewall.default"),
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_firewall.with_project", "google_compute_firewall.default"),
				),
			},
		},
	})
}

func testAccComputeFirewallDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "tf-test-firewall-net-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "default" {
  name          = "tf-test-firewall-%{random_suffix}"
  network       = google_compute_network.default.name
  description   = "Test firewall rule for data source"
  source_ranges = ["10.0.0.0/8"]

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["80", "8080", "1000-2000"]
  }
}

data "google_compute_firewall" "default" {
  name = google_compute_firewall.default.name
}

data "google_compute_firewall" "with_project" {
  name    = google_compute_firewall.default.name
  project = google_compute_firewall.default.project
}
`, context)
}
