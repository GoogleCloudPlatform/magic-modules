package apigee_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"
)

func TestAccApigeeDeveloperApp_apigeeDeveloperApp_full(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckApigeeDeveloperAppDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeDeveloperApp_apigeeDeveloperApp_full(context),
			},
			{
				ResourceName:            "google_apigee_developer_app.apigee_developer_app",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"developer_email", "org_id"},
			},
			{
				Config: testAccApigeeDeveloperApp_apigeeDeveloperApp_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_apigee_developer_app.apigee_developer_app", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_apigee_developer_app.apigee_developer_app",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"developer_email", "org_id"},
			},
		},
	})
}

func testAccApigeeDeveloperApp_apigeeDeveloperApp_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "time_sleep" "wait_60_seconds" {
  create_duration = "60s"
  depends_on = [google_project.project]
}

resource "google_project_service" "apigee" {
  project = google_project.project.project_id
  service = "apigee.googleapis.com"
  depends_on = [time_sleep.wait_60_seconds]
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
  depends_on = [google_project_service.apigee]
}

resource "google_project_service" "servicenetworking" {
  project = google_project.project.project_id
  service = "servicenetworking.googleapis.com"
  depends_on = [google_project_service.compute]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"
  depends_on = [google_project_service.servicenetworking]
}

resource "google_compute_network" "apigee_network" {
  name       = "apigee-network"
  project    = google_project.project.project_id
  depends_on = [time_sleep.wait_120_seconds]
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
  depends_on              = [google_project_service.servicenetworking]
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

resource "google_apigee_instance" "apigee_instance" {
  name               = "tf-test%{random_suffix}"
  location           = "us-central1"
  org_id             = google_apigee_organization.apigee_org.id
  peering_cidr_range = "SLASH_22"
}

resource "google_apigee_developer" "apigee_developer" {
  email      = "tf-test%{random_suffix}@acme.com"
  first_name = "John"
  last_name  = "Doe"
  user_name  = "john.doe"
  org_id     = google_apigee_organization.apigee_org.id
  depends_on = [
    google_apigee_instance.apigee_instance
  ]
}

resource "google_apigee_developer_app" "apigee_developer_app" {
  name              = "tf-test%{random_suffix}"
  developer_email   = google_apigee_developer.apigee_developer.email
  org_id            = google_apigee_organization.apigee_org.id
  callback_url      = "http://localhost"
	status						= "revoked"
	key_expires_in		= 900000

  attributes {
    name  = "sample_name"
    value = "sample_value"
  }
}
`, context)
}

func testAccApigeeDeveloperApp_apigeeDeveloperApp_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "time_sleep" "wait_60_seconds" {
  create_duration = "60s"
  depends_on = [google_project.project]
}

resource "google_project_service" "apigee" {
  project = google_project.project.project_id
  service = "apigee.googleapis.com"
  depends_on = [time_sleep.wait_60_seconds]
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
  depends_on = [google_project_service.apigee]
}

resource "google_project_service" "servicenetworking" {
  project = google_project.project.project_id
  service = "servicenetworking.googleapis.com"
  depends_on = [google_project_service.compute]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"
  depends_on = [google_project_service.servicenetworking]
}

resource "google_compute_network" "apigee_network" {
  name       = "apigee-network"
  project    = google_project.project.project_id
  depends_on = [time_sleep.wait_120_seconds]
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
  depends_on              = [google_project_service.servicenetworking]
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

resource "google_apigee_instance" "apigee_instance" {
  name               = "tf-test%{random_suffix}"
  location           = "us-central1"
  org_id             = google_apigee_organization.apigee_org.id
  peering_cidr_range = "SLASH_22"
}

resource "google_apigee_developer" "apigee_developer" {
  email      = "tf-test%{random_suffix}@acme.com"
  first_name = "John"
  last_name  = "Doe"
  user_name  = "john.doe"
  org_id     = google_apigee_organization.apigee_org.id
  depends_on = [
    google_apigee_instance.apigee_instance
  ]
}

resource "google_apigee_developer_app" "apigee_developer_app" {
  name              = "tf-test%{random_suffix}"
  developer_email   = google_apigee_developer.apigee_developer.email
  org_id            = google_apigee_organization.apigee_org.id
  callback_url      = "http://localhost:1234"

  attributes {
    name  = "updated_name"
    value = "updated_value"
  }
}
`, context)
}
