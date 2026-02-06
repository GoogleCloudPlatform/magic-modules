package apigee_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccApigeeSpace_handwritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"space_name":      "test-space-" + acctest.RandString(t, 10),
		"display_name":    "Test Space",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApigeeSpaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeSpace_handwrittenConfig(context),
			},
			{
				ResourceName:      "google_apigee_space.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApigeeSpace_handwrittenConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "%{org_id}-%{space_name}"
  name            = "%{org_id}-%{space_name}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "apigee" {
  project = google_project.project.project_id
  service = "apigee.googleapis.com"
}

resource "google_compute_network" "apigee_network" {
  name    = "apigee-network"
  project = google_project.project.project_id
}

resource "google_compute_global_address" "apigee_range" {
  name          = "apigee-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.apigee_network.id
  project       = google_project.project.project_id
}

resource "google_service_networking_connection" "apigee_vpc_connection" {
  network                 = google_compute_network.apigee_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.apigee_range.name]
}

resource "google_apigee_organization" "apigee_org" {
  analytics_region   = "us-central1"
  project_id         = google_project.project.project_id
  authorized_network = google_compute_network.apigee_network.id
  depends_on         = [
    google_service_networking_connection.apigee_vpc_connection,
    google_project_service.apigee,
  ]
}

resource "google_apigee_space" "primary" {
  org_id       = google_apigee_organization.apigee_org.id
  space_id     = "%{space_name}"
  name         = "organizations/${google_apigee_organization.apigee_org.name}/spaces/%{space_name}"
  display_name = "%{display_name}"
}
`, context)
}
