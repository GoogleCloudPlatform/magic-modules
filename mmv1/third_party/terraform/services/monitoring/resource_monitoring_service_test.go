package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/monitoring"
)

func TestAccMonitoringService_basic(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSlo_cloudEndpoints(randomSuffix, "an-endpoint"),
			},
			{
				ResourceName:      "google_monitoring_service.srv",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitoringSlo_cloudEndpoints(randomSuffix, "another-endpoint"),
			},
			{
				ResourceName:      "google_monitoring_service.srv",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringSlo_cloudEndpoints(randSuffix, endpoint string) string {
	return fmt.Sprintf(`
resource "google_monitoring_service" "srv" {
	service_id = "tf-test-srv-%s"
	display_name = "My Basic CloudEnpoints Service"
	basic_service {
		service_type  = "CLOUD_ENDPOINTS"
		service_labels = {
			service = "%s"
		}
	}
}
`, randSuffix, endpoint)
}

func TestAccMonitoringService_cloudRun(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringService_cloudRun(randomSuffix, "My Cloud Run Service", "cloud-run-svc-1", "us-central1"),
			},
			{
				ResourceName:      "google_monitoring_service.srv",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitoringService_cloudRun(randomSuffix, "My Updated Cloud Run Service", "cloud-run-svc-1", "us-central1"),
			},
			{
				ResourceName:      "google_monitoring_service.srv",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitoringService_cloudRun(randomSuffix, "My Updated Cloud Run Service", "cloud-run-svc-2", "us-east1"),
			},
			{
				ResourceName:      "google_monitoring_service.srv",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringService_cloudRun(randSuffix, displayName, serviceName, location string) string {
	return fmt.Sprintf(`
resource "google_monitoring_service" "srv" {
	service_id   = "tf-test-srv-%s"
	display_name = "%s"
	cloud_run {
		service_name = "%s"
		location     = "%s"
	}
}
`, randSuffix, displayName, serviceName, location)
}
