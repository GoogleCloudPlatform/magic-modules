package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeInterconnectAttachmentGroup_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"deletion_protection": false,
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInterconnectAttachmentGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInterconnectAttachmentGroup_basic(context),
			},
			{
				ResourceName:      "google_compute_interconnect_attachment_group.example-interconnect-attachment-group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeInterconnectAttachmentGroup_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_compute_interconnect" "example-interconnect" {
  name                 = "tf-test-example-interconnect%{random_suffix}"
  customer_name        = "internal_customer" # Special customer only available for Google testing.
  interconnect_type    = "DEDICATED"
  link_type            = "LINK_TYPE_ETHERNET_100G_LR"
  location             = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.project_id}/global/interconnectLocations/z2z-us-central1-zone1-tzcbfa-z" # Special location only available for Google testing.
  requested_link_count = 1
  admin_enabled        = true
  description          = "example description"
  noc_contact_email    = "user@example.com"
  labels = {
	mykey = "myvalue"
  }
}

resource "google_compute_network" "example-network" {
  name = "tf-test-example-network%{random_suffix}"
}

resource "google_compute_router" "example-router" {
  name    = "tf-test-example-router%{random_suffix}"
  network = google_compute_network.example-network.name
}

resource "google_compute_interconnect_attachment" "example-attachment" {
  name = "tf-test-example-attachment%{random_suffix}"
  router = google_compute_router.example-router.name
  interconnect = google_compute_interconnect.example-interconnect.id
  vlan_tag8021q = 5
}

resource "google_compute_interconnect_attachment_group" "example-interconnect-attachment-group" {
  name   = "tf-test-example-interconnect-attachment-group%{random_suffix}"
  intent {
    availability_sla = "NO_SLA"
  }
}
`, context)
}

func TestAccComputeInterconnectAttachmentGroup_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"deletion_protection": false,
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInterconnectAttachmentGroup_basic(context),
			},
			{
				ResourceName:      "google_compute_interconnect_attachment_group.example-interconnect-attachment-group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeInterconnectAttachmentGroup_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_interconnect_attachment_group.example-interconnect-attachment-group", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_compute_interconnect_attachment_group.example-interconnect-attachment-group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeInterconnectAttachmentGroup_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_compute_interconnect" "example-interconnect" {
  name                 = "tf-test-example-interconnect%{random_suffix}"
  customer_name        = "internal_customer" # Special customer only available for Google testing.
  interconnect_type    = "DEDICATED"
  link_type            = "LINK_TYPE_ETHERNET_100G_LR"
  location             = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.project_id}/global/interconnectLocations/z2z-us-central1-zone1-tzcbfa-z" # Special location only available for Google testing.
  requested_link_count = 1
  admin_enabled        = true
  description          = "example description"
  noc_contact_email    = "user@example.com"
  labels = {
	mykey = "myvalue"
  }
}

resource "google_compute_network" "example-network" {
  name = "tf-test-example-network%{random_suffix}"
}

resource "google_compute_router" "example-router" {
  name    = "tf-test-example-router%{random_suffix}"
  network = google_compute_network.example-network.name
}

resource "google_compute_interconnect_attachment" "example-attachment" {
  name = "tf-test-example-attachment%{random_suffix}"
  router = google_compute_router.example-router.name
  interconnect = google_compute_interconnect.example-interconnect.id
  vlan_tag8021q = 5
}

resource "google_compute_interconnect_attachment_group" "example-interconnect-attachment-group" {
  name   	  = "tf-test-example-interconnect-attachment-group%{random_suffix}"
  intent {
    availability_sla = "NO_SLA"
  }
  attachments {
	name = "my-attachment"
	attachment = google_compute_interconnect_attachment.example-attachment.name
  }
  description = "New description"
}
`, context)
}
