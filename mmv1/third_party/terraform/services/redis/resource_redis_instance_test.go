package redis_test

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
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

func TestAccRedisInstance_deletionprotection(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstance_deletionprotection(name, "us-central1", true),
			},
			{
				ResourceName:            "google_redis_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "deletion_protection"},
			},
			{
				Config:      testAccRedisInstance_deletionprotection(name, "us-west2", true),
				ExpectError: regexp.MustCompile("deletion_protection"),
			},
			{
				Config: testAccRedisInstance_deletionprotection(name, "us-central1", false),
			},
		},
	})
}

func testAccRedisInstance_deletionprotection(name string, region string, deletionProtection bool) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  region       = "%s"
  display_name   = "tf-test-instance"
  memory_size_gb = 1
  deletion_protection = %t

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
`, name, region, deletionProtection)
}

func TestAccRedisInstance_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "redis-instance-tagkey", map[string]interface{}{})
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestOrganizationTagValue(t, "redis-instance-tagvalue", tagKey),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceTags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_redis_instance.test", "tags.%"),
					checkRedisInstanceTags(t),
				),
			},
			{
				ResourceName:            "google_redis_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func checkRedisInstanceTags(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_redis_instance" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			// 1. Get the configured tag key and value from the state.
			var configuredTagValueNamespacedName string
			var tagKeyNamespacedName, tagValueShortName string
			for key, val := range rs.Primary.Attributes {
				if strings.HasPrefix(key, "tags.") && key != "tags.%" {
					tagKeyNamespacedName = strings.TrimPrefix(key, "tags.")
					tagValueShortName = val
					if tagValueShortName != "" {
						configuredTagValueNamespacedName = fmt.Sprintf("%s/%s", tagKeyNamespacedName, tagValueShortName)
						break
					}
				}
			}

			if configuredTagValueNamespacedName == "" {
				return fmt.Errorf("could not find a configured tag value in the state for resource %s", rs.Primary.ID)
			}

			// Check if placeholders are still present.
			if strings.Contains(configuredTagValueNamespacedName, "%{") {
				return fmt.Errorf("tag namespaced name contains unsubstituted variables: %q. Ensure the context map in the test step is populated", configuredTagValueNamespacedName)
			}

			// 2. Describe the tag value using the namespaced name to get its full resource name.
			safeNamespacedName := url.QueryEscape(configuredTagValueNamespacedName)
			describeTagValueURL := fmt.Sprintf("https://cloudresourcemanager.googleapis.com/v3/tagValues/namespaced?name=%s", safeNamespacedName)

			respDescribe, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    describeTagValueURL,
				UserAgent: config.UserAgent,
			})

			if err != nil {
				return fmt.Errorf("error describing tag value using namespaced name %q: %v", configuredTagValueNamespacedName, err)
			}

			fullTagValueName, ok := respDescribe["name"].(string)
			if !ok || fullTagValueName == "" {
				return fmt.Errorf("tag value details (name) not found in response for namespaced name: %q, response: %v", configuredTagValueNamespacedName, respDescribe)
			}

			// 3. Get the tag bindings from the Redis Instance.
			parts := strings.Split(rs.Primary.ID, "/")
			if len(parts) != 6 {
				return fmt.Errorf("invalid resource ID format: %s", rs.Primary.ID)
			}
			project := parts[1]
			location := parts[3]
			instance_id := parts[5]

			parentURL := fmt.Sprintf("//redis.googleapis.com/projects/%s/locations/%s/instances/%s", project, location, instance_id)
			listBindingsURL := fmt.Sprintf("https://%s-cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", location, url.QueryEscape(parentURL))

			resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    listBindingsURL,
				UserAgent: config.UserAgent,
			})

			if err != nil {
				return fmt.Errorf("error calling TagBindings API: %v", err)
			}

			tagBindingsVal, exists := resp["tagBindings"]
			if !exists {
				tagBindingsVal = []interface{}{}
			}

			tagBindings, ok := tagBindingsVal.([]interface{})
			if !ok {
				return fmt.Errorf("'tagBindings' is not a slice in response for resource %s. Response: %v", rs.Primary.ID, resp)
			}

			// 4. Perform the comparison.
			foundMatch := false
			for _, binding := range tagBindings {
				bindingMap, ok := binding.(map[string]interface{})
				if !ok {
					continue
				}
				if bindingMap["tagValue"] == fullTagValueName {
					foundMatch = true
					break
				}
			}

			if !foundMatch {
				return fmt.Errorf("expected tag value %s (from namespaced %q) not found in tag bindings for resource %s. Bindings: %v", fullTagValueName, configuredTagValueNamespacedName, rs.Primary.ID, tagBindings)
			}

			t.Logf("Successfully found matching tag binding for %s with tagValue %s", rs.Primary.ID, fullTagValueName)
		}

		return nil
	}
}

func testAccRedisInstanceTags(context map[string]interface{}) string {
	return acctest.Nprintf(`	
resource "google_redis_instance" "test" {
  name = "tf-test-instance-%{random_suffix}"
  memory_size_gb = 5
  tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
}`, context)
}
