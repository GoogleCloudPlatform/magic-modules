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

func TestAccVmwareengineNetworkPeering_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
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
		CheckDestroy: testAccCheckVmwareengineNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareengineNetworkPeering_config(context, "Sample description."),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_vmwareengine_network_peering.ds", "google_vmwareengine_network_peering.vmw-engine-network-peering"),
				),
			},
			{
				ResourceName:            "google_vmwareengine_network_peering.vmw-engine-network-peering",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccVmwareengineNetworkPeering_config(context, "Updated description."),
			},
			{
				ResourceName:            "google_vmwareengine_network_peering.vmw-engine-network-peering",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}

func testAccVmwareengineNetworkPeering_config(context map[string]interface{}, description string) string {
	context["description"] = description
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "network-peering-nw" {
  project           = "%{vmwareengine_project}"
  name              = "tf-test-sample-nw%{random_suffix}"
  location          = "global"
  type              = "STANDARD"
}

resource "google_vmwareengine_network" "network-peering-peer-nw" {
  project           = "%{vmwareengine_project}"
  name              = "tf-test-peer-nw%{random_suffix}"
  location          = "global"
  type              = "STANDARD"
}

resource "google_vmwareengine_network_peering" "vmw-engine-network-peering" {
  project = "%{vmwareengine_project}"
  name = "tf-test-sample-network-peering%{random_suffix}"
  description = "%{description}"
  vmware_engine_network = google_vmwareengine_network.network-peering-nw.id
  peer_network = google_vmwareengine_network.network-peering-peer-nw.id
  peer_network_type = "VMWARE_ENGINE_NETWORK"
}

data "google_vmwareengine_network_peering" "ds" {
  project = "%{vmwareengine_project}"
  name = google_vmwareengine_network_peering.vmw-engine-network-peering.name
}
`, context)
}

func TestAccVmwareengineNetworkPeering_tags(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "vennp-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "vennp-tagvalue", tagKey)

	venSuffix1 := acctest.RandString(t, 10)
	venSuffix2 := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           org,
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"ven_suffix1":   venSuffix1,
		"ven_suffix2":   venSuffix2,
		"ven_name1":     "tf-test-ven-" + venSuffix1,
		"ven_name2":     "tf-test-ven-" + venSuffix2,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineNetworkPeeringDestroyProducer(t), // Assuming this exists
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareengineNetworkPeeringTags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_vmwareengine_network_peering.default", "tags.%"),
					testAccCheckVmwareengineNetworkPeeringHasTagBindings(t),
				),
			},
			{
				ResourceName:            "google_vmwareengine_network_peering.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "vmware_engine_network", "peer_network", "tags"},
			},
		},
	})
}

func testAccVmwareengineNetworkPeeringTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "ven1" {
  name        = "%{ven_name1}"
  location    = "global"
  type        = "STANDARD"
  description = "Terraform test network 1 for peering"
}

resource "google_vmwareengine_network" "ven2" {
  name        = "%{ven_name2}"
  location    = "global"
  type        = "STANDARD"
  description = "Terraform test network 2 for peering"
}

resource "google_vmwareengine_network_peering" "default" {
  name                  = "tf-test-ven-peering-%{random_suffix}"
  vmware_engine_network = google_vmwareengine_network.ven1.id
  peer_network          = google_vmwareengine_network.ven2.id
  peer_network_type     = "VMWARE_ENGINE_NETWORK"
  # description = "Optional description"

  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}

func testAccCheckVmwareengineNetworkPeeringHasTagBindings(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vmwareengine_network_peering" {
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

			// 3. Get the tag bindings from the VMware Engine Network Peering.
			// ID format: projects/{project}/locations/{location}/networkPeerings/{networkPeeringId}
			// Network Peerings are global for Vmware Engine
			parentURL := fmt.Sprintf("//vmwareengine.googleapis.com/%s", rs.Primary.ID)
			listBindingsURL := fmt.Sprintf("https://cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", url.QueryEscape(parentURL))

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

