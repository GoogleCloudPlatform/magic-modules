package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMemcacheInstance_update(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("%d", tf.RandInt(t))
	name := fmt.Sprintf("tf-test-%s", prefix)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMemcacheInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemcacheInstance_update(prefix, name),
			},
			{
				ResourceName:      "google_memcache_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMemcacheInstance_update2(prefix, name),
			},
			{
				ResourceName:      "google_memcache_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMemcacheInstance_update(prefix, name string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "network" {
  provider = google-beta
  name = "tf-test%s"
}

resource "google_compute_global_address" "service_range" {
  provider = google-beta
  name          = "tf-test%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.network.id
}

resource "google_service_networking_connection" "private_service_connection" {
  provider = google-beta
  network                 = google_compute_network.network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.service_range.name]
}

resource "google_memcache_instance" "test" {
  provider = google-beta
  name = "%s"
  region = "us-central1"
  authorized_network = google_service_networking_connection.private_service_connection.network

  node_config {
    cpu_count      = 1
    memory_size_mb = 1024
  }
  node_count = 1
}
`, prefix, prefix, name)
}

func testAccMemcacheInstance_update2(prefix, name string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "network" {
  provider = google-beta
  name = "tf-test%s"
}

resource "google_compute_global_address" "service_range" {
  provider = google-beta
  name          = "tf-test%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.network.id
}

resource "google_service_networking_connection" "private_service_connection" {
  provider = google-beta
  network                 = google_compute_network.network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.service_range.name]
}

resource "google_memcache_instance" "test" {
  provider = google-beta
  name = "%s"
  region = "us-central1"
  authorized_network = google_service_networking_connection.private_service_connection.network

  node_config {
    cpu_count      = 1
    memory_size_mb = 2048
  }
  node_count = 2
}
`, prefix, prefix, name)
}
