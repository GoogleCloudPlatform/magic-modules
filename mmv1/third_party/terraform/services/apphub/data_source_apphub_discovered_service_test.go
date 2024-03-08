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
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
		"host_project":    envvar.GetTestProjectFromEnv(),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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

# Enable Compute API
resource "google_project_service" "compute_service_project" {
  project_id = google_project.service_project.project_id
  service = "compute.googleapis.com"
}

resource "google_apphub_service_project_attachment" "service_project_attachment" {
  service_project_attachment_id = google_project.service_project_full.project_id
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
  name                    = "ilb_network-%{random_suffix}"
  project                 = google_project.service_project.project_id
  auto_create_subnetworks = false
}

# backend subnet
resource "google_compute_subnetwork" "ilb_subnet" {
  name          			 = "ilb_subnet-%{random_suffix}"
  project       			 = google_project.service_project.project_id
  ip_cidr_range 			 = "10.0.1.0/24"
  region        			 = "us-central1"
  network       			 = google_compute_network.ilb_network.id
}

# forwarding rule
resource "google_compute_forwarding_rule" "forwarding_rule" {
  name                  = "forwarding_rule-%{random_suffix}"
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
  name                  = "backend_service-%{random_suffix}"
  project               = google_project.service_project.project_id
  region                = "us-central1"
  health_checks         = [google_compute_health_check.default.id]
}
    
# health check
resource "google_compute_health_check" "default" {
  name     					 		= "health_check-%{random_suffix}"
  project  					 		= google_project.service_project.project_id
  check_interval_sec 		= 1
  timeout_sec        		= 1

  tcp_health_check {
    port = "80"
  }
}
`, context)
}
