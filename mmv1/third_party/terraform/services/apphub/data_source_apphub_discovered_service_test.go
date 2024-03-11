package apphub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestDataSourceApphubDiscoveredService_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testDataSourceApphubDiscoveredService_basic(context),
			},
		},
	})
}

func testDataSourceApphubDiscoveredService_basic(context map[string]interface{}) string {
	return acctest.Nprintf(
		`
resource "google_project" "service_project" {
	project_id ="apphub-service-project-%{random_suffix}"
	name = "Service Project"
	org_id = "%{org_id}"
}

resource "time_sleep" "wait_120s_for_service_project" {
  depends_on = [google_project.service_project]
  create_duration = "120s"
}

# Enable Compute API
resource "google_project_service" "compute_service_project" {
  project = google_project.service_project.project_id
  service = "compute.googleapis.com"
  depends_on = [time_sleep.wait_120s_for_service_project]
}

resource "time_sleep" "wait_120s_for_compute_api" {
  depends_on = [google_project_service.compute_service_project]
  create_duration = "120s"
}

resource "google_apphub_service_project_attachment" "service_project_attachment" {
  service_project_attachment_id = google_project.service_project.project_id
  depends_on = [time_sleep.wait_120s_for_service_project]
}

# discovered service block
data "google_apphub_discovered_service" "catalog-service" {
  location = "us-central1"
  # ServiceReference | Application Hub | Google Cloud
  # Using this reference means that this resource will not be provisioned until the forwarding rule is fully created
  service_uri = "//compute.googleapis.com/${google_compute_forwarding_rule.forwarding_rule.id}"
	depends_on = [google_apphub_service_project_attachment.service_project_attachment]
}

# VPC network
resource "google_compute_network" "ilb_network" {
  name                    = "ilb-network-%{random_suffix}"
  project                 = google_project.service_project.project_id
  auto_create_subnetworks = false
  depends_on = [time_sleep.wait_120s_for_compute_api]
}

# backend subnet
resource "google_compute_subnetwork" "ilb_subnet" {
  name          			 = "ilb-subnet-%{random_suffix}"
  project       			 = google_project.service_project.project_id
  ip_cidr_range 			 = "10.0.1.0/24"
  region        			 = "us-central1"
  network       			 = google_compute_network.ilb_network.id
}

# forwarding rule
resource "google_compute_forwarding_rule" "forwarding_rule" {
  name                  = "forwarding-rule-%{random_suffix}"
  project               = google_project.service_project.project_id
  region                = "us-central1"
  ip_version            = "IPV4"
  load_balancing_scheme = "INTERNAL"
  all_ports             = true
  backend_service       = google_compute_region_backend_service.backend.id
  network               = google_compute_network.ilb_network.id
  subnetwork            = google_compute_subnetwork.ilb_subnet.id
}

# backend service
resource "google_compute_region_backend_service" "backend" {
  name                  = "backend-service-%{random_suffix}"
  project               = google_project.service_project.project_id
  region                = "us-central1"
  health_checks         = [google_compute_health_check.default.id]
}
    
# health check
resource "google_compute_health_check" "default" {
  name     					 		= "health-check-%{random_suffix}"
  project  					 		= google_project.service_project.project_id
  check_interval_sec 		= 1
  timeout_sec        		= 1

  tcp_health_check {
    port = "80"
  }
  depends_on = [time_sleep.wait_120s_for_compute_api]
}
`, context)
}
