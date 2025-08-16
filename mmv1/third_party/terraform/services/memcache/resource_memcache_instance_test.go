package memcache_test

import (
	"context"
	"fmt"
	"testing"

	memcache "cloud.google.com/go/memcache/apiv1"
	memcachepb "cloud.google.com/go/memcache/apiv1/memcachepb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
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

func TestAccMemcacheInstance_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "memcache-instances-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "memcache-instances-tagvalue", tagKey)

	testContext := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"random_suffix": acctest.RandString(t, 10),
	}
	resourceName := "google_memcache_instance.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemcacheInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemcacheInstanceTags(testContext),
				Check: resource.ComposeTestCheckFunc(
					checkMemcacheInstanceTags(resourceName, testContext),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func testAccMemcacheInstanceTags(testContext map[string]interface{}) string {
	return acctest.Nprintf(`
	provider "google" {
  project                 = "kshitij-memcached-test"
  user_project_override   = true
}
	resource "google_memcache_instance" "test" {
	  name = "tf-test-instance-%{random_suffix}"
	  node_count = 1
	  region = "us-central1"
	  node_config {
	    cpu_count = 1
	    memory_size_mb = 1024
	  }
	  tags = {
	    "%{org}/%{tagKey}" = "%{tagValue}"
	  }
	}`, testContext)
}

// This function gets the instance via the Memcache API and inspects its tags.
func checkMemcacheInstanceTags(resourceName string, testContext map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// Get resource attributes from state
		project := rs.Primary.Attributes["project"]
		location := rs.Primary.Attributes["region"]
		instanceName := rs.Primary.Attributes["name"]

		// Construct the expected full tag key
		expectedTagKey := fmt.Sprintf("%s/%s", testContext["org"], testContext["tagKey"])
		expectedTagValue := fmt.Sprintf("%s", testContext["tagValue"])

		// This `ctx` variable is now a `context.Context` object
		ctx := context.Background()

		// Create a Memcache client
		memcacheClient, err := memcache.NewCloudMemcacheClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create memcache client: %v", err)
		}
		defer memcacheClient.Close()

		// Construct the request to get the instance details
		req := &memcachepb.GetInstanceRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/instances/%s", project, location, instanceName),
		}

		// Get the Memcache instance
		instance, err := memcacheClient.GetInstance(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get memcache instance '%s': %v", req.Name, err)
		}

		// Check the instance's labels for the expected tag
		// In the Memcache API, tags are represented as labels.
		labels := instance.GetLabels()
		if labels == nil {
			return fmt.Errorf("expected labels not found on instance '%s'", req.Name)
		}

		if actualValue, ok := labels[expectedTagKey]; ok {
			if actualValue == expectedTagValue {
				// The tag was found with the correct value. Success!
				return nil
			}
			return fmt.Errorf("tag key '%s' found with incorrect value. Expected: %s, Got: %s", expectedTagKey, expectedTagValue, actualValue)
		}

		// If we reach here, the tag key was not found.
		return fmt.Errorf("expected tag key '%s' not found on instance '%s'", expectedTagKey, req.Name)
	}
}
