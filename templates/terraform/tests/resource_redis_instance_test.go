package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRedisInstance_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRedisInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRedisInstance_basic(name),
			},
			resource.TestStep{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRedisInstance_update(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRedisInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRedisInstance_update(name),
			},
			resource.TestStep{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccRedisInstance_update2(name),
			},
			resource.TestStep{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRedisInstance_full(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-test")
	network := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRedisInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRedisInstance_full(name, network),
			},
			resource.TestStep{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckRedisInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_redis_instance" {
			continue
		}

		redisIdParts := strings.Split(rs.Primary.ID, "/")
		if len(redisIdParts) != 3 {
			return fmt.Errorf("Unexpected resource ID %s, expected {project}/{region}/{name}", rs.Primary.ID)
		}

		project, region, inst := redisIdParts[0], redisIdParts[1], redisIdParts[2]

		name := fmt.Sprintf("projects/%s/locations/%s/instances/%s", project, region, inst)
		_, err := config.clientRedis.Projects.Locations.Get(name).Do()
		if err == nil {
			return fmt.Errorf("Redis instance still exists")
		}
	}

	return nil
}

func testAccRedisInstance_basic(name string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
	name           = "%s"
	memory_size_gb = 1
}`, name)
}

func testAccRedisInstance_update(name string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
	name           = "%s"
	display_name   = "pre-update"
	memory_size_gb = 1

	labels {
		my_key    = "my_val"
		other_key = "other_val"
	}

	redis_configs {
		maxmemory-policy       = "allkeys-lru"
		notify-keyspace-events = "KEA"
	}
}`, name)
}

func testAccRedisInstance_update2(name string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
	name           = "%s"
	display_name   = "post-update"
	memory_size_gb = 1

	labels {
		my_key    = "my_val"
		other_key = "new_val"
	}

	redis_configs {
		maxmemory-policy       = "noeviction"
		notify-keyspace-events = ""
	}
}`, name)
}

func testAccRedisInstance_full(name, network string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "test" {
	name = "%s"
}

resource "google_redis_instance" "test" {
	name           = "%s"
	tier           = "STANDARD_HA"
	memory_size_gb = 1

	authorized_network = "${google_compute_network.test.self_link}"

	region                  = "us-central1"
	location_id             = "us-central1-a"
	alternative_location_id = "us-central1-f"

	redis_version     = "REDIS_3_2"
	display_name      = "Terraform Test Instance"
	reserved_ip_range = "192.168.0.0/29"

	labels {
		my_key    = "my_val"
		other_key = "other_val"
	}

	redis_configs {
		maxmemory-policy       = "allkeys-lru"
		notify-keyspace-events = "KEA"
	}
}`, network, name)
}
