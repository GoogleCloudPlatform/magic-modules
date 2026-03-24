package memcache_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMemcacheInstance_update(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("%d", acctest.RandInt(t))
	name := fmt.Sprintf("tf-test-%s", prefix)
	network := acctest.BootstrapSharedServiceNetworkingConnection(t, "memcache-instance-update-1")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemcacheInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemcacheInstance_update(prefix, name, network),
			},
			{
				ResourceName:            "google_memcache_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reserved_ip_range_id"},
			},
			{
				Config: testAccMemcacheInstance_update2(prefix, name, network),
			},
			{
				ResourceName:      "google_memcache_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					// Check if effective_maintenance_version is present and not empty
					testAccMemcacheInstance_checkMaintenanceVersionIsNotEmpty("effective_maintenance_version"),
				),
			},
		},
	})
}

func testAccMemcacheInstance_update(prefix, name, network string) string {
	return fmt.Sprintf(`
resource "google_memcache_instance" "test" {
  name = "%s"
  region = "us-central1"
  authorized_network = data.google_compute_network.memcache_network.id

  node_config {
    cpu_count      = 1
    memory_size_mb = 1024
  }
  node_count = 1

  memcache_parameters {
    params = {
      "listen-backlog" = "2048"
      "max-item-size" = "8388608"
    }
  }
  reserved_ip_range_id = ["tf-bootstrap-addr-memcache-instance-update-1"]
}

data "google_compute_network" "memcache_network" {
  name = "%s"
}
`, name, network)
}

func testAccMemcacheInstance_update2(prefix, name, network string) string {
	return fmt.Sprintf(`
resource "google_memcache_instance" "test" {
  name = "%s"
  region = "us-central1"
  authorized_network = data.google_compute_network.memcache_network.id

  node_config {
    cpu_count      = 1
    memory_size_mb = 1024
  }
  node_count = 2

  memcache_parameters {
    params = {
      "listen-backlog" = "2048"
      "max-item-size" = "8388608"
    }
  }

  memcache_version = "MEMCACHE_1_6_15"
}

data "google_compute_network" "memcache_network" {
  name = "%s"
}
`, name, network)
}

func TestAccMemcacheInstance_deletionprotection(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("%d", acctest.RandInt(t))
	name := fmt.Sprintf("tf-test-%s", prefix)
	network := acctest.BootstrapSharedServiceNetworkingConnection(t, "memcache-instance-update-1")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemcacheInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemcacheInstanceConfig(prefix, name, network, "us-central1", true), // deletion_protection = true
			},
			{
				ResourceName:            "google_memcache_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reserved_ip_range_id", "deletion_protection"},
			},
			{
				Config:      testAccMemcacheInstanceConfig(prefix, name, network, "us-west2", true), // deletion_protection = true
				ExpectError: regexp.MustCompile("deletion_protection"),
			},
			{
				Config: testAccMemcacheInstanceConfig(prefix, name, network, "us-central1", false), // deletion_protection = false
			},
			{
				ResourceName:            "google_memcache_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reserved_ip_range_id", "deletion_protection"},
			},
		},
	})
}

func testAccMemcacheInstanceConfig(prefix, name, network, region string, deletionProtection bool) string {
	return fmt.Sprintf(`
resource "google_memcache_instance" "test" {
  name = "%s"
  region = "%s"
  authorized_network = data.google_compute_network.memcache_network.id
  deletion_protection = %t
  node_config {
    cpu_count      = 1
    memory_size_mb = 1024
  }
  node_count = 1
  memcache_parameters {
    params = {
      "listen-backlog" = "2048"
      "max-item-size" = "8388608"
    }
  }
  reserved_ip_range_id = ["tf-bootstrap-addr-memcache-instance-update-1"]
}
data "google_compute_network" "memcache_network" {
  name = "%s"
}
`, name, region, deletionProtection, network)
}

// Helper function to check if a maintenance version attribute is not empty.
func testAccMemcacheInstance_checkMaintenanceVersionIsNotEmpty(attr string) resource.TestCheckFunc {
	return func(s *resource.State) error {
		rs, ok := s.RootResourceStateMode.Resources["google_memcache_instance.test"]
		if !ok {
			return fmt.Errorf("root resource not found")
		}

		attrValue, ok := rs.Attributes[attr]
		if !ok {
			return fmt.Errorf("attribute %s not found", attr)
		}
		if attrValue == "" {
			return fmt.Errorf("attribute %s is empty", attr)
		}
		// Optional: Add more specific checks if needed, e.g., format validation.
		// For now, just checking for presence and non-empty is sufficient per the request.
		return nil
	}
}
