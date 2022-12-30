package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleComputeRegionInstanceGroupManager_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeRegionInstanceGroupManager_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_compute_region_instance_group_manager.foo", "google_compute_region_instance_group_manager.appserver"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleComputeRegionInstanceGroupManager_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_instance_template" "appserver" {
  name         = "template%{random_suffix}"
  machine_type = "e2-medium"
	
  disk {
	source_image = "debian-cloud/debian-11"
	auto_delete  = true
	disk_size_gb = 100
	boot         = true
  }
	
  network_interface {
	network = "default"
  }
	
  metadata = {
	foo = "bar"
  }
	
  can_ip_forward = true
}

resource "google_compute_region_instance_group_manager" "appserver" {
  name = "appserver-igm%{random_suffix}"
	
  base_instance_name         = "app"
  region                     = "us-central1"
  distribution_policy_zones  = ["us-central1-a"]
	
  version {
	instance_template = google_compute_instance_template.appserver.id
  }
	
	
  named_port {
	name = "custom"
	port = 8888
  }
  wait_for_instances = false
}
	
data "google_compute_region_instance_group_manager" "foo" {
   name = google_compute_region_instance_group_manager.appserver.name
}

`, context)

}
