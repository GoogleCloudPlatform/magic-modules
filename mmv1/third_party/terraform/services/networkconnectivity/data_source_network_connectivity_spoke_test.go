package networkconnectivity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/networkconnectivity"
)

func TestAccDataSourceNetworkConnectivitySpoke_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivitySpokeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkConnectivitySpoke_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_network_connectivity_spoke.default", "google_network_connectivity_spoke.primary"),
					acctest.CheckDataSourceStateMatchesResourceState("data.google_network_connectivity_spoke.default_optional", "google_network_connectivity_spoke.primary"),
				),
			},
		},
	})
}

func TestAccDataSourceNetworkConnectivitySpoke_linkedVpcNetwork(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivitySpokeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkConnectivitySpoke_linkedVpcNetwork(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_network_connectivity_spoke.default", "google_network_connectivity_spoke.primary"),
					acctest.CheckDataSourceStateMatchesResourceState("data.google_network_connectivity_spoke.default_no_project", "google_network_connectivity_spoke.primary"),
				),
			},
		},
	})
}

func testAccDataSourceNetworkConnectivitySpoke_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub-%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_compute_network" "network" {
  name                    = "tf-test-network-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_router" "router" {
  name    = "tf-test-router-%{random_suffix}"
  region  = "%{region}"
  network = google_compute_network.network.name
  bgp {
    asn = 16550
  }
}

resource "google_compute_interconnect_attachment" "interconnect_attachment" {
  name                     = "tf-test-ia-%{random_suffix}"
  edge_availability_domain = "AVAILABILITY_DOMAIN_1"
  type                     = "PARTNER"
  router                   = google_compute_router.router.id
  mtu                      = 1500
  region                   = "%{region}"
}

resource "google_network_connectivity_spoke" "primary" {
  name        = "tf-test-spoke-%{random_suffix}"
  location    = "%{region}"
  description = "A sample spoke with a linked Interconnect Attachment"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_interconnect_attachments {
    uris                       = [google_compute_interconnect_attachment.interconnect_attachment.self_link]
    site_to_site_data_transfer = true
    include_import_ranges      = ["ALL_IPV4_RANGES"]
  }
}

data "google_network_connectivity_spoke" "default" {
  name     = google_network_connectivity_spoke.primary.name
  location = google_network_connectivity_spoke.primary.location
  project  = google_network_connectivity_spoke.primary.project
}

data "google_network_connectivity_spoke" "default_optional" {
  name = google_network_connectivity_spoke.primary.name
}
`, context)
}

func testAccDataSourceNetworkConnectivitySpoke_linkedVpcNetwork(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network" {
  name                    = "tf-test-network-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub-%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-two = "value-one"
  }
}

resource "google_network_connectivity_spoke" "primary" {
  name        = "tf-test-spoke-%{random_suffix}"
  location    = "global"
  description = "A sample spoke with a linked VPC network"
  labels = {
    label-one = "value-one"
  }
  hub = google_network_connectivity_hub.basic_hub.id
  linked_vpc_network {
    exclude_export_ranges = [
      "198.51.100.0/24",
      "10.10.0.0/16"
    ]
    include_export_ranges = [
      "198.51.100.0/23",
      "10.0.0.0/8"
    ]
    uri = google_compute_network.network.self_link
  }
}

data "google_network_connectivity_spoke" "default" {
  name     = google_network_connectivity_spoke.primary.name
  location = google_network_connectivity_spoke.primary.location
  project  = google_network_connectivity_spoke.primary.project
}

data "google_network_connectivity_spoke" "default_no_project" {
  name     = google_network_connectivity_spoke.primary.name
  location = google_network_connectivity_spoke.primary.location
}
`, context)
}
