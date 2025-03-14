package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func testAccDataSourceComputeRegionBackendServiceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_region_backend_service" {
				continue
			}

			_, err := config.NewComputeClient(config.UserAgent).RegionBackendServices.Get(
				config.Project, rs.Primary.Attributes["region"], rs.Primary.Attributes["name"]).Do()
			if err == nil {
				return fmt.Errorf("Region Backend Service still exists")
			}
		}

		return nil
	}
}

func TestAccDataSourceComputeRegionBackendService_basic(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	region := "us-central1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccDataSourceComputeRegionBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeRegionBackendService_basic(serviceName, checkName, region),
				Check:  acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_region_backend_service.baz", "google_compute_region_backend_service.foobar"),
			},
		},
	})
}

func testAccDataSourceComputeRegionBackendService_basic(serviceName, checkName, region string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  name                  = "%s"
  description           = "foobar backend service"
  region                = "%s"
  protocol              = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  health_checks        = [google_compute_region_health_check.zero.self_link]
}

resource "google_compute_region_health_check" "zero" {
  name               = "%s"
  region            = "%s"
  http_health_check {
    port = 80
  }
}

data "google_compute_region_backend_service" "baz" {
  name   = google_compute_region_backend_service.foobar.name
  region = google_compute_region_backend_service.foobar.region
}
`, serviceName, region, checkName, region)
}

func TestAccDataSourceComputeRegionBackendService_withProject(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	region := "us-central1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccDataSourceComputeRegionBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeRegionBackendService_withProject(serviceName, checkName, region),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_region_backend_service.baz", "google_compute_region_backend_service.foobar"),
					resource.TestCheckResourceAttrSet("data.google_compute_region_backend_service.baz", "project"),
				),
			},
		},
	})
}

func testAccDataSourceComputeRegionBackendService_withProject(serviceName, checkName, region string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  name                  = "%s"
  description           = "foobar backend service"
  region                = "%s"
  protocol              = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  health_checks        = [google_compute_region_health_check.zero.self_link]
}

resource "google_compute_region_health_check" "zero" {
  name               = "%s"
  region            = "%s"
  http_health_check {
    port = 80
  }
}

data "google_compute_region_backend_service" "baz" {
  project = google_compute_region_backend_service.foobar.project
  name    = google_compute_region_backend_service.foobar.name
  region  = google_compute_region_backend_service.foobar.region
}
`, serviceName, region, checkName, region)
}
