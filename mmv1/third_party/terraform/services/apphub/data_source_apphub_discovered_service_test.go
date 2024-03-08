package apphub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestDataSourceApphubDiscoveredService_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
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
	project_id ="<%= ctx[:vars]['service_project_attachment_id'] %>"
	name = "Service Project"
	org_id = "<%= ctx[:test_env_vars]['org_id'] %>"
}

resource "google_apphub_service_project_attachment" "service_project_attachment" {
  service_project_attachment_id = google_project.service_project.project_id
}

# discovered service block
data "google_apphub_discovered_service" "catalog-service" {
  provider = google
  location = "us-east1"
  # ServiceReference | Application Hub | Google Cloud
  # Using this reference means that this resource will not be provisioned until the forwarding rule is fully created
  service_uri = "//compute.googleapis.com/${google_compute_forwarding_rule.forwarding_rule.id}"
	depends_on = [google_apphub_service_project_attachment.service_project_attachment]
}

# VPC network
resource "google_compute_network" "ilb_network" {
  name                    = "<%= ctx[:vars]['ilb_network'] %>"
  project                 = google_project.service_project.project_id
  auto_create_subnetworks = false
}

# backend subnet
resource "google_compute_subnetwork" "ilb_subnet" {
  name          = "<%= ctx[:vars]['ilb_subnet'] %>"
  project       = google_project.service_project.project_id
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-east1"
  network       = google_compute_network.ilb_network.id
}

# forwarding rule
resource "google_compute_forwarding_rule" "forwarding_rule" {
  name                  ="<%= ctx[:vars]['forwarding_rule'] %>"
  project               = google_project.service_project.project_id
  region                = "us-east1"
  ip_version            = "IPV4"
  load_balancing_scheme = "INTERNAL"
  all_ports             = true
  backend_service       = google_compute_region_backend_service.backend.id
  network               = google_compute_network.ilb_network.id
  subnetwork            = google_compute_subnetwork.ilb_subnet.id
}

# backend service
resource "google_compute_region_backend_service" "backend" {
  name                  = "<%= ctx[:vars]['backend_service'] %>"
  project               = google_project.service_project.project_id
  region                = "us-east1"
  health_checks         = [google_compute_health_check.default.id]
}
    
# health check
resource "google_compute_health_check" "default" {
  name     = "<%= ctx[:vars]['health_check'] %>"
  project  = google_project.service_project.project_id
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = "80"
  }
}
`, context)
}
