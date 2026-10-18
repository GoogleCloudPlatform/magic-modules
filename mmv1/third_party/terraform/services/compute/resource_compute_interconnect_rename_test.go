package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

func TestAccComputeInterconnect_rename(t *testing.T) {
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
				Config: testAccComputeInterconnect_renameStart(context),
			},
			{
				ResourceName:            "google_compute_interconnect.example-interconnect",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
			{
				Config: testAccComputeInterconnect_renameUpdated(context),
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

func testAccComputeInterconnect_renameStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_compute_interconnect" "example-interconnect" {
  name                 = "tf-test-interconnect-%{random_suffix}"
  customer_name        = "internal_customer"
  interconnect_type    = "DEDICATED"
  link_type            = "LINK_TYPE_ETHERNET_10G_LR"
  location             = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.project_id}/global/interconnectLocations/z2z-us-east4-zone1-lciadl-a"
  requested_link_count = 1
}
`, context)
}

func testAccComputeInterconnect_renameUpdated(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_compute_interconnect" "example-interconnect" {
  name                 = "tf-test-interconnect-renamed-%{random_suffix}"
  customer_name        = "internal_customer"
  interconnect_type    = "DEDICATED"
  link_type            = "LINK_TYPE_ETHERNET_10G_LR"
  location             = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.project_id}/global/interconnectLocations/z2z-us-east4-zone1-lciadl-a"
  requested_link_count = 1
}
`, context)
}
