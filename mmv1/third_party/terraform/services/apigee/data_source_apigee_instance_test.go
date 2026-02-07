package apigee_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceApigeeInstance_basic(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"instance_name": "tf-test-" + acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApigeeInstanceConfig(context),
			},
			{
				Config: testAccDataSourceApigeeInstanceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.google_apigee_instance.test", "host",
						"google_apigee_instance.test", "host"),
					resource.TestCheckResourceAttrPair(
						"data.google_apigee_instance.test", "service_attachment",
						"google_apigee_instance.test", "service_attachment"),
				),
			},
		},
	})
}

func testAccDataSourceApigeeInstanceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "net" {
  name = "net-%{random_suffix}"
}

resource "google_compute_global_address" "range" {
  name          = "range-%{random_suffix}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.net.id
}

resource "google_service_networking_connection" "con" {
  network                 = google_compute_network.net.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.range.name]
}

resource "google_apigee_organization" "org" {
  project_id                           = acctest.GetTestProject()
  analytics_region                     = "us-central1"
  runtime_type                         = "CLOUD"
  authorized_network                   = google_compute_network.net.id
  depends_on                           = [google_service_networking_connection.con]
}

resource "google_apigee_instance" "test" {
  name     = "%{instance_name}"
  location = "%{location}"
  org_id   = google_apigee_organization.org.id
}

data "google_apigee_instance" "test" {
  name   = google_apigee_instance.test.name
  org_id = google_apigee_instance.test.org_id
}
`, context)
}
