package vmwareengine_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"strings"
)

func TestAccVmwareengineNetwork_vmwareEngineNetworkUpdate(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"region":          "me-west1", // region with allocated quota
		"random_suffix":   acctest.RandString(t, 10),
		"organization":    envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	configTemplate := vmwareEngineNetworkConfigTemplate(context)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckVmwareengineNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(configTemplate, "description1"),
			},
			{
				ResourceName:            "google_vmwareengine_network.default-nw",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
			{
				Config: fmt.Sprintf(configTemplate, "description2"),
			},
			{
				ResourceName:            "google_vmwareengine_network.default-nw",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
		},
	})
}

func vmwareEngineNetworkConfigTemplate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "default-nw" {
  project     = google_project_service.acceptance.project
  name        = "%{region}-default"
  location    = "%{region}"
  type        = "LEGACY"
  description = "%s"
}

# there can be only 1 Legacy network per region for a given project, so creating new project to isolate tests.
resource "google_project" "acceptance" {
  name            = "tf-test-%{random_suffix}"
  project_id      = "tf-test-%{random_suffix}"
  org_id          = "%{organization}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "acceptance" {
  project  = google_project.acceptance.project_id
  service  = "vmwareengine.googleapis.com"

  # Needed for CI tests for permissions to propagate, should not be needed for actual usage
  depends_on = [time_sleep.wait_60_seconds]
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [google_project.acceptance]

  create_duration = "60s"
}
`, context)
}

func TestAccVmwareengineNetwork_tags(t *testing.T) {
    t.Parallel()

    org := envvar.GetTestOrgFromEnv(t)
    tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "ven-tagkey", map[string]interface{}{})
    tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "ven-tagvalue", tagKey)

    context := map[string]interface{}{
        "random_suffix": acctest.RandString(t, 10),
        "org":           org,
        "tagKey":        tagKey,
        "tagValue":      tagValue,
        "location":      "global", // Or a specific region if testing regional VEN
    }

    acctest.VcrTest(t, resource.TestCase{
        PreCheck:                 func() { acctest.AccTestPreCheck(t) },
        ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
        CheckDestroy:             testAccCheckVmwareengineNetworkDestroyProducer(t),
        Steps: []resource.TestStep{
            {
                Config: testAccVmwareengineNetworkTags(context),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttrSet("google_vmwareengine_network.default", "tags.%"),
                    testAccCheckVmwareengineNetworkHasTagBindings(t),
                ),
            },
            {
                ResourceName:            "google_vmwareengine_network.default",
                ImportState:             true,
                ImportStateVerify:       true,
                ImportStateVerifyIgnore: []string{"name", "location", "tags"}, // terraform_labels may not apply
            },
        },
    })
}

func testAccVmwareengineNetworkTags(context map[string]interface{}) string {
    return acctest.Nprintf(`
resource "google_vmwareengine_network" "default" {
  name        = "tf-test-ven-%{random_suffix}"
  location    = "%{location}"
  type        = "STANDARD"
  description = "Terraform test network with tags"
  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}

func testAccCheckVmwareengineNetworkHasTagBindings(t *testing.T) func(s *terraform.State) error {
    return func(s *terraform.State) error {
        for name, rs := range s.RootModule().Resources {
            if rs.Type != "google_vmwareengine_network" {
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
                return fmt.Errorf("tag namespaced name contains unsubstituted variables: %q", configuredTagValueNamespacedName)
            }

            // 2. Describe the tag value to get its full resource name.
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

            // 3. Get the tag bindings from the VMware Engine Network.
            // ID format: projects/{project}/locations/{location}/vmwareEngineNetworks/{vmwareEngineNetworkId}
            parts := strings.Split(rs.Primary.ID, "/")
            if len(parts) != 6 {
                return fmt.Errorf("invalid resource ID format: %s", rs.Primary.ID)
            }
            location := parts[3]

            parentURL := fmt.Sprintf("//vmwareengine.googleapis.com/%s", rs.Primary.ID)

            var listBindingsURL string
            if location == "global" {
                listBindingsURL = fmt.Sprintf("https://cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", url.QueryEscape(parentURL))
            } else {
                listBindingsURL = fmt.Sprintf("https://%s-cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", location, url.QueryEscape(parentURL))
            }

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
