package memcache_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"net/url"
	"strings"
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
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "memcache-inst-tagkey", map[string]interface{}{})
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestOrganizationTagValue(t, "memcache-inst-tagvalue", tagKey),
		"region":        "us-central1",
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemcacheInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemcacheInstanceTags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_memcache_instance.default", "tags.%"),
					testAccCheckMemcacheInstanceHasTagBindings(t),
				),
			},
			{
				ResourceName:            "google_memcache_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "region", "labels", "tags"},
			},
		},
	})
}

func testAccCheckMemcacheInstanceHasTagBindings(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_memcache_instance" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			var configuredTagValueNamespacedName string
			var tagKeyNamespacedName, tagValueShortName string
			for key, val := range rs.Primary.Attributes {
				if strings.HasPrefix(key, "tags.") && key != "tags.#" {
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

			if strings.Contains(configuredTagValueNamespacedName, "%{") {
				return fmt.Errorf("tag namespaced name contains unsubstituted variables: %q. Ensure the context map in the test step is populated", configuredTagValueNamespacedName)
			}

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

			parts := strings.Split(rs.Primary.ID, "/")
			if len(parts) != 6 {
				return fmt.Errorf("invalid resource ID format: %s", rs.Primary.ID)
			}
			project := parts[1]
			location := parts[3] // This is the region
			instance_id := parts[5]

			parentURL := fmt.Sprintf("//memcache.googleapis.com/projects/%s/locations/%s/instances/%s", project, location, instance_id)
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

func testAccMemcacheInstanceTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  project                 = "kshitij-memcached-test"
  user_project_override   = true
}
resource "google_memcache_instance" "default" {
  name          = "tf-test-my-memcache-%{random_suffix}"
  region        =  "us-central1"
  node_count    = 1
  node_config {
    cpu_count   = 1
    memory_size_mb = 1024
  }
  labels = {
    env = "test"
  }
  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}`, context)
}

func TestAccMemcacheInstance_tagsUpdateDelete(t *testing.T) {
	t.Parallel()
	randSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": randSuffix,
		"org":           envvar.GetTestOrgFromEnv(t),
		"region":        "us-central1",
	}

	instanceResourceName := "google_memcache_instance.default"
	tagValue1ResourceName := "google_tags_tag_value.value1"
	tagValue2ResourceName := "google_tags_tag_value.value2"
	bindingResourceName := "google_tags_location_tag_binding.binding"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemcacheInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			// Step 1: Create instance and bind with tagValue1
			{
				Config: testAccMemcacheInstanceTagsConfig(context, 1), // Bind to value1
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(instanceResourceName, "name"),
					testAccCheckMemcacheInstanceTagBinding(t, instanceResourceName, tagValue1ResourceName, true),
					testAccCheckMemcacheInstanceTagBinding(t, instanceResourceName, tagValue2ResourceName, false),
				),
			},
			// Step 2: Import
			{
				ResourceName:            instanceResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "region", "labels"},
			},
			{
				ResourceName:      bindingResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Step 3: Update binding to tagValue2
			{
				Config: testAccMemcacheInstanceTagsConfig(context, 2), // Bind to value2
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMemcacheInstanceTagBinding(t, instanceResourceName, tagValue2ResourceName, true),
					testAccCheckMemcacheInstanceTagBinding(t, instanceResourceName, tagValue1ResourceName, false),
				),
			},
			// Step 4: Delete binding
			{
				Config: testAccMemcacheInstanceTagsConfig(context, 0), // No binding resource
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMemcacheInstanceNoTagBindings(t, instanceResourceName),
				),
			},
		},
	})
}

// testAccMemcacheInstanceTagsConfig generates the Terraform configuration string.
// version 1: bind to tagValue1
// version 2: bind to tagValue2
// version 0: no binding resource
func testAccMemcacheInstanceTagsConfig(context map[string]interface{}, version int) string {
	baseConfig := acctest.Nprintf(`
provider "google" {
  project                 = "kshitij-memcached-test"
  user_project_override   = true
}

resource "google_memcache_instance" "default" {
  name          = "tf-test-memcache-%{random_suffix}"
  region        =  "%{region}"
  node_count    = 1
  node_config {
    cpu_count   = 1
    memory_size_mb = 1024
  }
}

resource "google_tags_tag_key" "key" {
  parent      = "organizations/%{org}"
  short_name  = "tf-key-%{random_suffix}"
}

resource "google_tags_tag_value" "value1" {
  parent      = google_tags_tag_key.key.name
  short_name  = "tf-value1-%{random_suffix}"
}

resource "google_tags_tag_value" "value2" {
  parent      = google_tags_tag_key.key.name
  short_name  = "tf-value2-%{random_suffix}"
}
`, context)

	bindingConfig := ""
	switch version {
	case 1:
		bindingConfig = `
resource "google_tags_location_tag_binding" "binding" {
  parent    = google_memcache_instance.default.id
  tag_value = google_tags_tag_value.value1.name
  location  = google_memcache_instance.default.region
}
`
	case 2:
		bindingConfig = `
resource "google_tags_location_tag_binding" "binding" {
  parent    = google_memcache_instance.default.id
  tag_value = google_tags_tag_value.value2.name
  location  = google_memcache_instance.default.region
}
`
	case 0:
		// No binding resource
	}

	return baseConfig + bindingConfig
}

// testAccCheckMemcacheInstanceTagBinding checks if the given Memcache instance is bound to the specified tag value.
func testAccCheckMemcacheInstanceTagBinding(t *testing.T, instanceResourceName, tagValueResourceName string, expectBound bool) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		instanceRS, ok := s.RootModule().Resources[instanceResourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", instanceResourceName)
		}
		if instanceRS.Primary.ID == "" {
			return fmt.Errorf("resource ID not set: %s", instanceResourceName)
		}

		tagValueRS, ok := s.RootModule().Resources[tagValueResourceName]
		if !ok {
			return fmt.Errorf("tag value resource not found: %s", tagValueResourceName)
		}
		if tagValueRS.Primary.ID == "" {
			return fmt.Errorf("tag value resource ID not set: %s", tagValueResourceName)
		}
		expectedTagValueName := tagValueRS.Primary.ID

		config := acctest.GoogleProviderConfig(t)
		parts := strings.Split(instanceRS.Primary.ID, "/")
		if len(parts) != 6 {
			return fmt.Errorf("invalid instance resource ID format: %s", instanceRS.Primary.ID)
		}
		project := parts[1]
		location := parts[3]
		instanceID := parts[5]

		parentURL := fmt.Sprintf("//memcache.googleapis.com/projects/%s/locations/%s/instances/%s", project, location, instanceID)
		listBindingsURL := fmt.Sprintf("https://%s-cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", location, url.QueryEscape(parentURL))

		resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    listBindingsURL,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return fmt.Errorf("error calling TagBindings API for %s: %v", parentURL, err)
		}

		tagBindingsVal, exists := resp["tagBindings"]
		if !exists {
			tagBindingsVal = []interface{}{}
		}
		tagBindings, ok := tagBindingsVal.([]interface{})
		if !ok {
			return fmt.Errorf("'tagBindings' is not a slice in response for resource %s", instanceRS.Primary.ID)
		}

		foundMatch := false
		for _, binding := range tagBindings {
			bindingMap, ok := binding.(map[string]interface{})
			if !ok {
				continue
			}
			if currentTagValue, ok := bindingMap["tagValue"].(string); ok && currentTagValue == expectedTagValueName {
				foundMatch = true
				break
			}
		}

		if expectBound && !foundMatch {
			return fmt.Errorf("expected tag value %s (from %s) not found in tag bindings for resource %s. Bindings: %v", expectedTagValueName, tagValueResourceName, instanceRS.Primary.ID, tagBindings)
		}
		if !expectBound && foundMatch {
			return fmt.Errorf("unexpected tag value %s (from %s) found in tag bindings for resource %s", expectedTagValueName, tagValueResourceName, instanceRS.Primary.ID)
		}

		return nil
	}
}

// testAccCheckMemcacheInstanceNoTagBindings checks that the instance has no tag bindings.
func testAccCheckMemcacheInstanceNoTagBindings(t *testing.T, instanceResourceName string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		instanceRS, ok := s.RootModule().Resources[instanceResourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", instanceResourceName)
		}
		if instanceRS.Primary.ID == "" {
			return fmt.Errorf("resource ID not set: %s", instanceResourceName)
		}

		config := acctest.GoogleProviderConfig(t)
		parts := strings.Split(instanceRS.Primary.ID, "/")
		if len(parts) != 6 {
			return fmt.Errorf("invalid instance resource ID format: %s", instanceRS.Primary.ID)
		}
		project := parts[1]
		location := parts[3]
		instanceID := parts[5]

		parentURL := fmt.Sprintf("//memcache.googleapis.com/projects/%s/locations/%s/instances/%s", project, location, instanceID)
		listBindingsURL := fmt.Sprintf("https://%s-cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", location, url.QueryEscape(parentURL))

		resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    listBindingsURL,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return fmt.Errorf("error calling TagBindings API for %s: %v", parentURL, err)
		}

		tagBindingsVal, exists := resp["tagBindings"]
		if !exists || tagBindingsVal == nil {
			return nil // No bindings array means no bindings exist.
		}
		tagBindings, ok := tagBindingsVal.([]interface{})
		if !ok {
			return fmt.Errorf("'tagBindings' is not a slice in response for resource %s, got: %T", instanceRS.Primary.ID, tagBindingsVal)
		}

		if len(tagBindings) > 0 {
			return fmt.Errorf("expected no tag bindings for resource %s, but found: %v", instanceRS.Primary.ID, tagBindings)
		}

		return nil
	}
}

