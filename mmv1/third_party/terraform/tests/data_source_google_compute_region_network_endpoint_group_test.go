package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRegionNetworkEndpointGroup(t *testing.T) {
	// Randomness in instance template
	skipIfVcr(t)
	t.Parallel()
	poolName := "tf-test-pool-" + randString(t, 6)
	dataSourceName := "tf-test-ds-" + randString(t, 6)
	rnegName := "tf-test-rneg" + randString(t, 6)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRegionNetworkEndpointGroup_basic(poolName, rnegName, dataSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_region_network_endpoint_group.data_source", "name", rnegName),
					resource.TestCheckResourceAttr("data.google_compute_region_network_endpoint_group.data_source", "project", getTestProjectFromEnv()),
					resource.TestCheckResourceAttrSet("data.google_compute_region_network_endpoint_group.data_source", "self_link"),
					resource.TestCheckResourceAttr("data.google_compute_region_network_endpoint_group.data_source", "region", "us-central1")),
			},
		},
	})
}

func testAccDataSourceRegionNetworkEndpointGroup_basic(poolName string, rigName string, dataSourceName string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "default" {
  name = "%s"
}

data "google_compute_image" "debian" {
  project = "debian-cloud"
  name    = "debian-11-bullseye-v20220719"
}

resource "google_compute_instance_template" "default" {
  machine_type = "e2-medium"
  disk {
    source_image = data.google_compute_image.debian.self_link
  }
  network_interface {
    access_config {
    }
    network = "default"
  }
}

resource "google_compute_region_instance_group_manager" "default" {
  name               = "%s"
  base_instance_name = "foo"
  version {
    instance_template = google_compute_instance_template.default.self_link
    name              = "primary"
  }
  region       = "us-central1"
  target_pools = [google_compute_target_pool.default.self_link]
  target_size  = 1

  named_port {
    name = "web"
    port = 80
  }
  wait_for_instances = true
}

data "google_compute_region_network_endpoint_group" "data_source" {
    name = "%s"
    self_link = google_compute_region_instance_group_manager.default.instance_group
    region       = "us-central1"
}
`, poolName, rigName, dataSourceName)
}
