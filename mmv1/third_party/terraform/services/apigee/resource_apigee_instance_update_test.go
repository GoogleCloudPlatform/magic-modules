package apigee_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccApigeeInstance_updateConsumerAcceptList(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix_1": acctest.RandString(t, 10),
		"random_suffix_2": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckApigeeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeInstance_basic(context),
			},
			{
				ResourceName:            "google_apigee_instance.apigee_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ip_range", "org_id"},
			},
			{
				Config: testAccApigeeInstance_updateConsumerAcceptList(context),
			},
			{
				ResourceName:            "google_apigee_instance.apigee_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ip_range", "org_id"},
			},
		},
	})
}

func testAccApigeeInstance_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project1" {
  project_id      = "tf-test%{random_suffix_1}"
  name            = "tf-test%{random_suffix_1}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "apigee" {
  project = google_project.project1.project_id
  service = "apigee.googleapis.com"
}

resource "google_project_service" "compute" {
  project = google_project.project1.project_id
  service = "compute.googleapis.com"
}

resource "google_project_service" "servicenetworking" {
  project = google_project.project1.project_id
  service = "servicenetworking.googleapis.com"
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"
  depends_on = [google_project_service.compute]
}

resource "google_compute_network" "apigee_network" {
  name       = "apigee-network"
  project    = google_project.project1.project_id
  depends_on = [time_sleep.wait_120_seconds]
}

resource "google_compute_global_address" "apigee_range" {
  name          = "apigee-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 21
  network       = google_compute_network.apigee_network.id
  project       = google_project.project1.project_id
}

resource "google_service_networking_connection" "apigee_vpc_connection" {
  network                 = google_compute_network.apigee_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.apigee_range.name]
  depends_on              = [google_project_service.servicenetworking]
}

resource "google_apigee_organization" "apigee_org" {
  analytics_region   = "us-central1"
  project_id         = google_project.project1.project_id
  authorized_network = google_compute_network.apigee_network.id
  billing_type       = "EVALUATION"
  depends_on = [
    google_service_networking_connection.apigee_vpc_connection,
    google_project_service.apigee,
  ]
}

resource "google_apigee_instance" "apigee_instance" {
  name     = "tf-test%{random_suffix_1}"
  location = "us-central1"
  org_id   = google_apigee_organization.apigee_org.id
  ip_range = "${google_compute_global_address.apigee_range.address}/22"
  consumer_accept_list = [
    google_project.project1.project_id,
  ]
}
`, context)
}

func testAccApigeeInstance_updateConsumerAcceptList(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project1" {
  project_id      = "tf-test%{random_suffix_1}"
  name            = "tf-test%{random_suffix_1}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project" "project2" {
  project_id      = "tf-test%{random_suffix_2}"
  name            = "tf-test%{random_suffix_2}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "apigee" {
  project = google_project.project1.project_id
  service = "apigee.googleapis.com"
}

resource "google_project_service" "compute" {
  project = google_project.project1.project_id
  service = "compute.googleapis.com"
}

resource "google_project_service" "servicenetworking" {
  project = google_project.project1.project_id
  service = "servicenetworking.googleapis.com"
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"
  depends_on = [google_project_service.compute]
}

resource "google_compute_network" "apigee_network" {
  name       = "apigee-network"
  project    = google_project.project1.project_id
  depends_on = [time_sleep.wait_120_seconds]
}

resource "google_compute_global_address" "apigee_range" {
  name          = "apigee-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 21
  network       = google_compute_network.apigee_network.id
  project       = google_project.project1.project_id
}

resource "google_service_networking_connection" "apigee_vpc_connection" {
  network                 = google_compute_network.apigee_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.apigee_range.name]
  depends_on              = [google_project_service.servicenetworking]
}

resource "google_apigee_organization" "apigee_org" {
  analytics_region   = "us-central1"
  project_id         = google_project.project1.project_id
  authorized_network = google_compute_network.apigee_network.id
  billing_type       = "EVALUATION"
  depends_on = [
    google_service_networking_connection.apigee_vpc_connection,
    google_project_service.apigee,
  ]
}

resource "google_apigee_instance" "apigee_instance" {
  name     = "tf-test%{random_suffix_1}"
  location = "us-central1"
  org_id   = google_apigee_organization.apigee_org.id
  ip_range = "${google_compute_global_address.apigee_range.address}/22"
  consumer_accept_list = [
    google_project.project1.project_id,
    google_project.project2.project_id,
  ]
}
`, context)
}
