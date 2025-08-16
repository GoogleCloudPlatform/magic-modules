package redis_test

import (
	"context"
	"fmt"
	"testing"

	redis "cloud.google.com/go/redis/apiv1"
	redispb "cloud.google.com/go/redis/apiv1/redispb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccRedisInstance_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstance_update(name, true),
			},
			{
				ResourceName:            "google_redis_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccRedisInstance_update2(name, true),
			},
			{
				ResourceName:            "google_redis_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccRedisInstance_update2(name, false),
			},
		},
	})
}

// Validate that read replica is enabled on the instance without having to recreate
func TestAccRedisInstance_updateReadReplicasMode(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceReadReplicasUnspecified(name, true),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstanceReadReplicasEnabled(name, true),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstanceReadReplicasUnspecified(name, false),
			},
		},
	})
}

/* Validate that read replica is enabled on the instance without recreate
 * and secondaryIp is auto provisioned when passed as 'auto' */
func TestAccRedisInstance_updateReadReplicasModeWithAutoSecondaryIp(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceReadReplicasUnspecified(name, true),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstanceReadReplicasEnabledWithAutoSecondaryIP(name, true),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstanceReadReplicasUnspecified(name, false),
			},
		},
	})
}

func testAccRedisInstanceReadReplicasUnspecified(name string, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "redissss"
  memory_size_gb = 5
	tier = "STANDARD_HA"
  region         = "us-central1"
	%s
  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
}
`, name, lifecycleBlock)
}

func testAccRedisInstanceReadReplicasEnabled(name string, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "redissss"
  memory_size_gb = 5
  tier = "STANDARD_HA"
  region         = "us-central1"
	%s
  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
  read_replicas_mode = "READ_REPLICAS_ENABLED"
  secondary_ip_range = "10.79.0.0/28"
	}
`, name, lifecycleBlock)
}

func testAccRedisInstanceReadReplicasEnabledWithAutoSecondaryIP(name string, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "redissss"
  memory_size_gb = 5
  tier = "STANDARD_HA"
  region         = "us-central1"
	%s
  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
  read_replicas_mode = "READ_REPLICAS_ENABLED"
  secondary_ip_range = "auto"
}
`, name, lifecycleBlock)
}

func TestAccRedisInstance_regionFromLocation(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	// Pick a zone that isn't in the provider-specified region so we know we
	// didn't fall back to that one.
	region := "us-west1"
	zone := "us-west1-a"
	if envvar.GetTestRegionFromEnv() == "us-west1" {
		region = "us-central1"
		zone = "us-central1-a"
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstance_regionFromLocation(name, zone),
				Check:  resource.TestCheckResourceAttr("google_redis_instance.test", "region", region),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRedisInstance_redisInstanceAuthEnabled(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstance_redisInstanceAuthEnabled(context),
			},
			{
				ResourceName:            "google_redis_instance.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccRedisInstance_redisInstanceAuthDisabled(context),
			},
			{
				ResourceName:            "google_redis_instance.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func TestAccRedisInstance_downgradeRedisVersion(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstance_redis5(name),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstance_redis4(name),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRedisInstance_update(name string, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "pre-update"
  memory_size_gb = 1
  region         = "us-central1"
	%s

  labels = {
    my_key    = "my_val"
    other_key = "other_val"
  }

  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
  redis_version = "REDIS_4_0"
}
`, name, lifecycleBlock)
}

func testAccRedisInstance_update2(name string, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "post-update"
  memory_size_gb = 1
	%s

  labels = {
    my_key    = "my_val"
    other_key = "new_val"
  }

  redis_configs = {
    maxmemory-policy       = "noeviction"
    notify-keyspace-events = ""
  }
  redis_version = "REDIS_5_0"
}
`, name, lifecycleBlock)
}

func testAccRedisInstance_regionFromLocation(name, zone string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  memory_size_gb = 1
  location_id    = "%s"
}
`, name, zone)
}

func testAccRedisInstance_redisInstanceAuthEnabled(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_instance" "cache" {
  name           = "tf-test-memory-cache%{random_suffix}"
  memory_size_gb = 1
  auth_enabled = true
}
`, context)
}

func testAccRedisInstance_redisInstanceAuthDisabled(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_instance" "cache" {
  name           = "tf-test-memory-cache%{random_suffix}"
  memory_size_gb = 1
  auth_enabled = false
}
`, context)
}

func testAccRedisInstance_redis5(name string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "redissss"
  memory_size_gb = 1
  region         = "us-central1"

  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
  redis_version = "REDIS_5_0"
}
`, name)
}

func testAccRedisInstance_redis4(name string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "redissss"
  memory_size_gb = 1
  region         = "us-central1"

  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
  redis_version = "REDIS_4_0"
}
`, name)
}

func TestAccRedisInstance_tags(t *testing.T) {
	t.Parallel()

	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "redis-instances-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "redis-instances-tagvalue", tagKey)

	testContext := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"random_suffix": acctest.RandString(t, 10),
	}
	resourceName := "google_redis_instance.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceTags(testContext),
				Check: resource.ComposeTestCheckFunc(
					checkRedisInstanceTags(resourceName, testContext),
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

func testAccRedisInstanceTags(testContext map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_redis_instance" "test" {
	  name = "tf-test-instance-%{random_suffix}"
	  memory_size_gb = 5
	  tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, testContext)
}

// This function gets the instance via the Redis API and inspects its tags.
func checkRedisInstanceTags(resourceName string, testContext map[string]interface{}) resource.TestCheckFunc {
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

		// Create a Redis client
		redisClient, err := redis.NewCloudRedisClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create redis client: %v", err)
		}
		defer redisClient.Close()

		// Construct the request to get the instance details
		req := &redispb.GetInstanceRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/instances/%s", project, location, instanceName),
		}

		// Get the Redis instance
		instance, err := redisClient.GetInstance(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get redis instance '%s': %v", req.Name, err)
		}

		// Check the instance's labels for the expected tag
		// In the Redis API, tags are represented as labels.
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
