package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeInterconnect_computeInterconnectMacsecTest(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInterconnectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInterconnect_computeInterconnectCreate(context),
			},
			{
				ResourceName:            "google_compute_interconnect.example-interconnect",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
			{
				Config: testAccComputeInterconnect_computeInterconnectEnableMacsec(context),
			},
			{
				ResourceName:            "google_compute_interconnect.example-interconnect",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccComputeInterconnect_computeInterconnectCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_compute_interconnect" "example-interconnect" {
  name                 = "tf-test-example-interconnect%{random_suffix}"
  customer_name        = "internal_customer" # Special customer only available for Google testing.
  interconnect_type    = "DEDICATED"
  link_type            = "LINK_TYPE_ETHERNET_100G_LR"
  location             = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.id}/global/interconnectLocations/z2z-us-east4-zone1-lciadl-z" # Special location only available for Google testing.
  requested_link_count = 1
  admin_enabled        = true
  description          = "example description"
  macsec_enabled       = false
  noc_contact_email    = "user@example.com"
  requested_features   = ["IF_MACSEC"]
  labels = {
    mykey = "myvalue"
  }
}
`, context)
}

func testAccComputeInterconnect_computeInterconnectEnableMacsec(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_compute_interconnect" "example-interconnect" {
  name                 = "tf-test-example-interconnect%{random_suffix}"
  customer_name        = "internal_customer" # Special customer only available for Google testing.
  interconnect_type    = "DEDICATED"
  link_type            = "LINK_TYPE_ETHERNET_100G_LR"
  location             = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.id}/global/interconnectLocations/z2z-us-east4-zone1-lciadl-z" # Special location only available for Google testing.
  requested_link_count = 1
  admin_enabled        = true
  description          = "example description"
  macsec_enabled       = true
  noc_contact_email    = "user@example.com"
  requested_features   = ["IF_MACSEC"]
  macsec {
    pre_shared_keys {
      name = "test-key"
      start_time = "2023-07-01T21:00:01.000Z"
    }
    fail_open = true
  }
  labels = {
    mykey = "newvalue"
  }
}
`, context)
}

func TestAccComputeInterconnect_computeInterconnectCrossSiteNetworkTest(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"location":      "z2z-us-west8-zone2-tzphxd-z",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInterconnectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInterconnect_crossSiteInterconnectBasic(context),
			},
			{
				ResourceName:            "google_compute_interconnect.example-interconnect-cross-site-network",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
			{
				Config: testAccComputeInterconnect_crossSiteInterconnectUpdate(context),
			},
			{
				ResourceName:            "google_compute_interconnect.example-interconnect-cross-site-network",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_interconnect.example-interconnect-cross-site-network", "description", "example description updated")),
			},
		},
	})
}

func testAccComputeInterconnect_crossSiteInterconnectBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_compute_cross_site_network" "example-cross-site-network" {
  name        = "{{index $.Vars "cross_site_network"}}"
  description = "Example cross site network"
  provider    = google-beta
}

resource "google_compute_interconnect" "example-interconnect-cross-site-network" {
  name                 = "tf-test-example-interconnect%{random_suffix}"
  customer_name        = "internal_customer" # Special customer only available for Google testing.
  interconnect_type    = "DEDICATED"
  requested_features  = [google_compute_cross_site_network.example-cross-site-network.name]
  cross_site_network   = google_compute_cross_site_network.example-cross-site-network.name
  depends_on = [
    google_compute_cross_site_network.example-cross-site-network
  ]
  available_features = [google_compute_cross_site_network.example-cross-site-network.name]
  link_type            = "LINK_TYPE_ETHERNET_10G_LR"
  location             = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.id}/global/interconnectLocations/z2z-us-west8-zone2-tzphxd-z" # Special location only available for Google testing.
  requested_link_count = 1
  description          = "example description"

}
`, context)
}
func testAccComputeInterconnect_crossSiteInterconnectUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_compute_cross_site_network" "example-cross-site-network" {
  name        = "{{index $.Vars "cross_site_network"}}"
  description = "Example cross site network"
  provider    = google-beta
}

resource "google_compute_interconnect" "example-interconnect-cross-site-network" {
  # This resource is used to test cross-site network interconnects.
  name                 = "tf-test-example-interconnect%{random_suffix}"
  customer_name        = "internal_customer" # Special customer only available for Google testing.
  interconnect_type    = "DEDICATED"
  requested_features  = [google_compute_cross_site_network.example-cross-site-network.name]
  cross_site_network   = google_compute_cross_site_network.example-cross-site-network.name
  depends_on = [
    google_compute_cross_site_network.example-cross-site-network
  ]
  available_features = [google_compute_cross_site_network.example-cross-site-network.name]
  link_type            = "LINK_TYPE_ETHERNET_10G_LR"
  location             = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.id}/global/interconnectLocations/z2z-us-west8-zone2-tzphxd-z" # Special location only available for Google testing.
  requested_link_count = 1
  description          = "example description updated"
}
`, context)
}
