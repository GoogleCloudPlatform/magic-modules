package vmwareengine_test

import (
	"os"
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

func TestAccVmwareengineNetworkPolicy_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":               "me-west1", // region with allocated quota
		"random_suffix":        acctest.RandString(t, 10),
		"org_id":               envvar.GetTestOrgFromEnv(t),
		"billing_account":      envvar.GetTestBillingAccountFromEnv(t),
		"vmwareengine_project": os.Getenv("GOOGLE_VMWAREENGINE_PROJECT"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckVmwareengineNetworkPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareengineNetworkPolicy_config(context, "description1", "192.168.0.0/26", false, false),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_vmwareengine_network_policy.ds", "google_vmwareengine_network_policy.vmw-engine-network-policy"),
				),
			},
			{
				ResourceName:            "google_vmwareengine_network_policy.vmw-engine-network-policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time"},
			},
			{
				Config: testAccVmwareengineNetworkPolicy_config(context, "description2", "192.168.1.0/26", true, true),
			},
			{
				ResourceName:            "google_vmwareengine_network_policy.vmw-engine-network-policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time"},
			},
		},
	})
}

func testAccVmwareengineNetworkPolicy_config(context map[string]interface{}, description string, edgeServicesCidr string, internetAccess bool, externalIp bool) string {
	context["internet_access"] = internetAccess
	context["external_ip"] = externalIp
	context["edge_services_cidr"] = edgeServicesCidr
	context["description"] = description

	return acctest.Nprintf(`
resource "google_vmwareengine_network" "network-policy-nw" {
  project           = "%{vmwareengine_project}"
  name              = "tf-test-sample-nw%{random_suffix}"
  location          = "global" 
  type              = "STANDARD"
  description       = "VMwareEngine standard network sample"
}

resource "google_vmwareengine_network_policy" "vmw-engine-network-policy" {
  project = "%{vmwareengine_project}"
  location = "%{region}"
  name = "tf-test-sample-network-policy%{random_suffix}"
  description = "%{description}" 
  internet_access {
    enabled = "%{internet_access}"
  }
  external_ip {
    enabled = "%{external_ip}"
  }
  edge_services_cidr = "%{edge_services_cidr}"
  vmware_engine_network = google_vmwareengine_network.network-policy-nw.id
}

data "google_vmwareengine_network_policy" "ds" {
  project = "%{vmwareengine_project}"
  name = google_vmwareengine_network_policy.vmw-engine-network-policy.name
  location = "%{region}"
}
`, context)
}

func TestAccVmwareengineNetworkPolicy_tags(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "venp-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "venp-tagvalue", tagKey)
	// Network Policies are regional. We need a VmwareEngineNetwork to associate with.
	// Assuming a VEN test setup function or a pre-existing one.
	venContext := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"location":      "global",
	}
	venConfig := testAccVmwareengineNetworkBasic(venContext) // Example helper for VEN

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           org,
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"region":        "me-west1",
		"ven_name":      "tf-test-ven-" + venContext["random_suffix"].(string),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineNetworkPolicyDestroyProducer(t), // Assuming this exists
		Steps: []resource.TestStep{
			{
				Config: venConfig + testAccVmwareengineNetworkPolicyTags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_vmwareengine_network_policy.default", "tags.%"),
					testAccCheckVmwareengineNetworkPolicyHasTagBindings(t),
				),
			},
			{
				ResourceName:            "google_vmwareengine_network_policy.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "vmware_engine_network", "tags"},
			},
		},
	})
}

// Example basic VEN config helper
func testAccVmwareengineNetworkBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "ven" {
  name        = "tf-test-ven-%{random_suffix}"
  location    = "%{location}"
  type        = "STANDARD"
  description = "Terraform test network for policy"
}
`, context)
}

func testAccVmwareengineNetworkPolicyTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network_policy" "default" {
  name        = "tf-test-ven-policy-%{random_suffix}"
  location    = "%{region}"
  vmware_engine_network = google_vmwareengine_network.ven.id
  edge_services_cidr  = "192.168.1.0/26"  # Required: Provide a valid, non-overlapping CIDR range

  internet_access {
    enabled = true
  }
  external_ip {
    enabled = true
  }

  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}

func testAccCheckVmwareengineNetworkPolicyHasTagBindings(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vmwareengine_network_policy" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			// 1. Get the configured tag key and value from the state.
			var configuredTagValueNamespacedName string
			for key, val := range rs.Primary.Attributes {
				if strings.HasPrefix(key, "tags.") && key != "tags.#" {
					tagKeyNamespacedName := strings.TrimPrefix(key, "tags.")
					tagValueShortName := val
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
				return fmt.Errorf("error describing tag value %q: %v", configuredTagValueNamespacedName, err)
			}
			fullTagValueName, ok := respDescribe["name"].(string)
			if !ok || fullTagValueName == "" {
				return fmt.Errorf("tag value name not found for %q: %v", configuredTagValueNamespacedName, respDescribe)
			}

			// 3. Get the tag bindings from the VMware Engine Network Policy.
			// ID format: projects/{project}/locations/{location}/networkPolicies/{networkPolicyId}
			parts := strings.Split(rs.Primary.ID, "/")
			if len(parts) != 6 {
				return fmt.Errorf("invalid resource ID format: %s", rs.Primary.ID)
			}
			location := parts[3] // This will be a region

			parentURL := fmt.Sprintf("//vmwareengine.googleapis.com/%s", rs.Primary.ID)
			// Network Policy is regional, so TagBindings API is always regional.
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
				return fmt.Errorf("'tagBindings' is not a slice for %s: %v", rs.Primary.ID, resp)
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
				return fmt.Errorf("expected tag value %s (from %q) not found in bindings for %s. Bindings: %v", fullTagValueName, configuredTagValueNamespacedName, rs.Primary.ID, tagBindings)
			}
			t.Logf("Successfully found matching tag binding for %s with tagValue %s", rs.Primary.ID, fullTagValueName)
		}
		return nil
	}
}
