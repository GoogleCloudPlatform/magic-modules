package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeInterconnect_computeInterconnectBasicTestExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		CheckDestroy:             testAccCheckComputeInterconnectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInterconnect_computeInterconnectBasicTestExample_basic(context),
			},
			{
				ResourceName:            "google_compute_interconnect.example-interconnect",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
			{
				Config: testAccComputeInterconnect_computeInterconnectBasicTestExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_interconnect.example-interconnect", plancheck.ResourceActionUpdate),
					},
				},
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

func testAccComputeInterconnect_computeInterconnectBasicTestExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  provider = google-beta
}

resource "google_compute_interconnect" "example-interconnect" {
  provider             = google-beta	
  name                 = "tf-test-example-interconnect%{random_suffix}"
  customer_name        = "internal_customer" # Special customer only available for Google testing.
  interconnect_type    = "DEDICATED"
  link_type            = "LINK_TYPE_ETHERNET_10G_LR"
  location             = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.name}/global/interconnectLocations/z2z-us-east4-zone1-lciadl-a" # Special location only available for Google testing.
  requested_link_count = 1
  admin_enabled        = true
  description          = "example description"
  macsec_enabled       = false
  noc_contact_email    = "user@example.com"
  requested_features   = []
  labels = {
    mykey = "myvalue"
  }
}
`, context)
}

func testAccComputeInterconnect_computeInterconnectBasicTestExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  provider = google-beta
}

resource "google_compute_interconnect" "example-interconnect" {
  provider             = google-beta
  name                 = "tf-test-example-interconnect%{random_suffix}"
  customer_name        = "internal_customer" # Special customer only available for Google testing.
  interconnect_type    = "DEDICATED"
  link_type            = "LINK_TYPE_ETHERNET_10G_LR"
  location             = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.name}/global/interconnectLocations/z2z-us-east4-zone1-lciadl-a" # Special location only available for Google testing.
  requested_link_count = 1
  admin_enabled        = true
  description          = "example description"
  macsec_enabled       = false
  noc_contact_email    = "user@example.com"
  requested_features   = []
  labels = {
    mykey = "myvalue"
  }
  aaiEnabled = true
  applicationAwareInterconnect = {
    profileDescription = "application awareness config with BandwidthPercentage policy."
	bandwidthPercentagePolicy = {
	  bandwidthPercentages = [
	 	 {
	  		trafficClass = "TC1"
			percentage   = 20
		 },
		 {
	  		trafficClass = "TC2"
			percentage   = 20
		 },
		 {
	  		trafficClass = "TC3"
			percentage   = 20
		 },
		 {
	  		trafficClass = "TC4"
			percentage   = 20
		 },
		 {
	  		trafficClass = "TC5"
			percentage   = 10
		 },
		 {
	  		trafficClass = "TC6"
			percentage   = 10
		 }
	  ]
	}
  }
}
`, context)
}
