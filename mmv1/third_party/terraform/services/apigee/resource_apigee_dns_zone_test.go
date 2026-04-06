package apigee_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccApigeeDnsZone_apigeeDnsZoneBasicTest(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix":   randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckApigeeDnsZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeDnsZone_basicTest(context),
			},
			{
				ResourceName:            "google_apigee_dns_zone.apigee_dns_zone",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"org_id", "dns_zone_id"},
			},
		},
	})
}

func testAccApigeeDnsZone_basicTest(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "apigee" {
  project = google_project.project.project_id
  service = "apigee.googleapis.com"
}

resource "google_project_service" "compute" {
  project    = google_project.project.project_id
  service    = "compute.googleapis.com"
  depends_on = [google_project_service.apigee]
}

resource "google_project_service" "dns" {
  project    = google_project.project.project_id
  service    = "dns.googleapis.com"
  depends_on = [google_project_service.compute]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"
  depends_on      = [google_project_service.dns]
}

resource "google_compute_network" "apigee_network" {
  name       = "apigee-network"
  project    = google_project.project.project_id
  depends_on = [time_sleep.wait_120_seconds]
}

# Create the Apigee org first. This provisions the Apigee service agent SA
# (service-{project_number}@gcp-sa-apigee.iam.gserviceaccount.com).
resource "google_apigee_organization" "apigee_org" {
  analytics_region    = "us-central1"
  project_id          = google_project.project.project_id
  disable_vpc_peering = true
  depends_on = [
    google_project_service.apigee,
    google_project_service.compute,
    google_project_service.dns,
  ]
}

# Grant dns.peer to the Apigee service agent AFTER the org is created (the SA
# only exists once the org has been provisioned).
resource "google_project_iam_member" "apigee_dns_peer" {
  project    = google_project.project.project_id
  role       = "roles/dns.peer"
  member     = "serviceAccount:service-${google_project.project.number}@gcp-sa-apigee.iam.gserviceaccount.com"
  depends_on = [google_apigee_organization.apigee_org]
}

resource "google_apigee_dns_zone" "apigee_dns_zone" {
  dns_zone_id = "tf-test%{random_suffix}"
  org_id      = google_apigee_organization.apigee_org.id
  domain      = "foo.example.com"
  description = "Test DNS zone"
  peering_config {
    target_project_id = google_project.project.project_id
    target_network_id = google_compute_network.apigee_network.name
  }
  depends_on = [google_project_iam_member.apigee_dns_peer]
}
`, context)
}
