package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRegionNetworkEndpointGroup_basic(t *testing.T) {
	t.Parallel()
	region := "us-central1"
	rnegName := "tf-test-rneg" + randString(t, 6)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRegionNetworkEndpointGroup_basic(rnegName, region),
				Check:  checkDataSourceStateMatchesResourceStateWithIgnores("data.google_compute_region_network_endpoint_group.data_source", "google_compute_instance_template.default", map[string]struct{}{"name": {}, "region": {}, "self_link": {}}),
			},
		},
	})
}

func testAccDataSourceRegionNetworkEndpointGroup_basic(rnegName string, region string) string {
	return fmt.Sprintf(`
  // Cloud Run Example
  resource "google_compute_region_network_endpoint_group" "cloudrun_neg" {
    name                  = "%{rnegName}"
    network_endpoint_type = "SERVERLESS"
    region                = "%{region}"
    cloud_run {
      service = google_cloud_run_service.cloudrun_neg.name
    }
  }

  resource "google_cloud_run_service" "cloudrun_neg" {
    name     = "cloudrun-neg"
    location = "us-central1"

    template {
      spec {
        containers {
          image = "us-docker.pkg.dev/cloudrun/container/hello"
        }
      }
    }

    traffic {
      percent         = 100
      latest_revision = true
    }

data "google_compute_region_network_endpoint_group" "data_source" {
    self_link = google_compute_region_network_endpoint_group.cloudrun_neg.instance_group
}
`, map[string]interface{}{"rnegName": rnegName, "region": region})
}
