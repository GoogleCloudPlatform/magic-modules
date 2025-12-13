package vmwareengine_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"net/url"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccVmwareenginePrivateCloud_vmwareEnginePrivateCloudUpdate(t *testing.T) {
	acctest.SkipIfVcr(t)
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
		CheckDestroy: testAccCheckVmwareenginePrivateCloudDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testVmwareenginePrivateCloudCreateConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						[]string{
							"deletion_delay_hours",
							"send_deletion_delay_hours_if_zero",
						}),
					testAccCheckGoogleVmwareengineNsxCredentialsMeta("data.google_vmwareengine_nsx_credentials.nsx-ds"),
					testAccCheckGoogleVmwareengineVcenterCredentialsMeta("data.google_vmwareengine_vcenter_credentials.vcenter-ds"),
				),
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "deletion_delay_hours", "send_deletion_delay_hours_if_zero"},
			},

			{
				Config: testVmwareenginePrivateCloudUpdateNodeConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						[]string{
							"deletion_delay_hours",
							"send_deletion_delay_hours_if_zero",
						}),
				),
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "deletion_delay_hours", "send_deletion_delay_hours_if_zero"},
			},

			{
				Config: testVmwareenginePrivateCloudUpdateAutoscaleConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						[]string{
							"deletion_delay_hours",
							"send_deletion_delay_hours_if_zero",
						}),
				),
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "deletion_delay_hours", "send_deletion_delay_hours_if_zero"},
			},

			{
				Config: testVmwareenginePrivateCloudDelayedDeleteConfig(context),
			},
			{
				ResourceName:            "google_vmwareengine_network.vmw-engine-nw",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},

			{
				Config: testVmwareenginePrivateCloudUndeleteConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						[]string{
							"deletion_delay_hours",
							"send_deletion_delay_hours_if_zero",
						}),
				),
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "deletion_delay_hours", "send_deletion_delay_hours_if_zero"},
			},

			{
				Config: testVmwareengineSubnetImportConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_vmwareengine_subnet.subnet-ds", "google_vmwareengine_subnet.vmw-engine-subnet"),
				),
			},
			{
				ResourceName:            "google_vmwareengine_subnet.vmw-engine-subnet",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},

			{
				Config: testVmwareengineSubnetUpdateConfig(context),
			},
			{
				ResourceName:            "google_vmwareengine_subnet.vmw-engine-subnet",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
		},
	})
}

func testVmwareenginePrivateCloudCreateConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudConfig(context, "sample description", "TIME_LIMITED", 1, 0) + testVmwareengineVcenterNSXCredentailsConfig(context)
}

func testVmwareenginePrivateCloudUpdateNodeConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudConfig(context, "sample updated description", "STANDARD", 3, 8) + testVmwareengineVcenterNSXCredentailsConfig(context)
}

func testVmwareenginePrivateCloudUpdateAutoscaleConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudAutoscaleConfig(context, "sample updated description", "", 3, 8) + testVmwareengineVcenterNSXCredentailsConfig(context)
}

func testVmwareenginePrivateCloudDelayedDeleteConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudDeletedConfig(context)
}

func testVmwareenginePrivateCloudUndeleteConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudAutoscaleConfig(context, "sample updated description", "STANDARD", 3, 0) + testVmwareengineVcenterNSXCredentailsConfig(context)
}

func testVmwareengineSubnetImportConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudAutoscaleConfig(context, "sample updated description", "STANDARD", 3, 0) + testVmwareengineSubnetConfig(context, "192.168.1.0/26")
}

func testVmwareengineSubnetUpdateConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudAutoscaleConfig(context, "sample updated description", "STANDARD", 3, 0) + testVmwareengineSubnetConfig(context, "192.168.2.0/26")
}

func testVmwareenginePrivateCloudConfig(context map[string]interface{}, description, pcType string, nodeCount, delayHours int) string {
	context["node_count"] = nodeCount
	context["delay_hrs"] = delayHours
	context["description"] = description
	context["type"] = pcType
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "vmw-engine-nw" {
  project = "%{vmwareengine_project}"
  name              = "tf-test-pc-nw-%{random_suffix}"
  location          = "global"
  type              = "STANDARD"
  description       = "PC network description."
}

resource "google_vmwareengine_private_cloud" "vmw-engine-pc" {
  project = "%{vmwareengine_project}"
  location = "%{region}-b"
  name = "tf-test-sample-pc%{random_suffix}"
  description = "%{description}"
  type = "%{type}"
  deletion_delay_hours = "%{delay_hrs}"
  send_deletion_delay_hours_if_zero = true
  network_config {
    management_cidr = "192.168.0.0/24"
    vmware_engine_network = google_vmwareengine_network.vmw-engine-nw.id
  }
  management_cluster {
    cluster_id = "tf-test-sample-mgmt-cluster-custom-core-count%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count = "%{node_count}"
      custom_core_count = 32
    }
  }
}

data "google_vmwareengine_private_cloud" "ds" {
    project = "%{vmwareengine_project}"
	location = "%{region}-b"
	name = "tf-test-sample-pc%{random_suffix}"
	depends_on = [
   	google_vmwareengine_private_cloud.vmw-engine-pc,
  ]
}
`, context)
}

func testVmwareenginePrivateCloudAutoscaleConfig(context map[string]interface{}, description, pcType string, nodeCount, delayHours int) string {
	context["node_count"] = nodeCount
	context["delay_hrs"] = delayHours
	context["description"] = description
	context["type"] = pcType
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "vmw-engine-nw" {
  project = "%{vmwareengine_project}"
  name              = "tf-test-pc-nw-%{random_suffix}"
  location          = "global"
  type              = "STANDARD"
  description       = "PC network description."
}

resource "google_vmwareengine_private_cloud" "vmw-engine-pc" {
  project = "%{vmwareengine_project}"
  location = "%{region}-b"
  name = "tf-test-sample-pc%{random_suffix}"
  description = "%{description}"
  type = "%{type}"
  deletion_delay_hours = "%{delay_hrs}"
  send_deletion_delay_hours_if_zero = true
  network_config {
    management_cidr = "192.168.0.0/24"
    vmware_engine_network = google_vmwareengine_network.vmw-engine-nw.id
  }
  management_cluster {
    cluster_id = "tf-test-sample-mgmt-cluster-custom-core-count%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count = "%{node_count}"
      custom_core_count = 32
    }
    autoscaling_settings {
      autoscaling_policies {
        autoscale_policy_id = "autoscaling-policy"
        node_type_id = "standard-72"
        scale_out_size = 1
        cpu_thresholds {
          scale_out = 80
          scale_in  = 15
        }
        consumed_memory_thresholds {
          scale_out = 75
          scale_in  = 20
        }
        storage_thresholds {
          scale_out = 80
          scale_in  = 20
        }
      }
      min_cluster_node_count = 3
      max_cluster_node_count = 8
      cool_down_period = "1800s"
    }
  }
}

data "google_vmwareengine_private_cloud" "ds" {
    project = "%{vmwareengine_project}"
	location = "%{region}-b"
	name = "tf-test-sample-pc%{random_suffix}"
	depends_on = [
   	google_vmwareengine_private_cloud.vmw-engine-pc,
  ]
}
`, context)
}

func testVmwareenginePrivateCloudDeletedConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "vmw-engine-nw" {
  project = "%{vmwareengine_project}"
  name              = "tf-test-pc-nw-%{random_suffix}"
  location          = "global"
  type              = "STANDARD"
  description       = "PC network description."
}
`, context)
}

func testVmwareengineVcenterNSXCredentailsConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_vmwareengine_nsx_credentials" "nsx-ds" {
	parent =  google_vmwareengine_private_cloud.vmw-engine-pc.id
}

data "google_vmwareengine_vcenter_credentials" "vcenter-ds" {
	parent =  google_vmwareengine_private_cloud.vmw-engine-pc.id
}
`, context)
}

func testVmwareengineSubnetConfig(context map[string]interface{}, ipCidrRange string) string {
	context["ip_cidr_range"] = ipCidrRange
	return acctest.Nprintf(`
resource "google_vmwareengine_subnet" "vmw-engine-subnet" {
  name = "service-2"
  parent =  google_vmwareengine_private_cloud.vmw-engine-pc.id
  ip_cidr_range = "%{ip_cidr_range}"
}

data "google_vmwareengine_subnet" "subnet-ds" {
  name = "service-2"
  parent = google_vmwareengine_private_cloud.vmw-engine-pc.id
  depends_on = [
    google_vmwareengine_subnet.vmw-engine-subnet,
  ]
}
`, context)
}

func testAccCheckGoogleVmwareengineNsxCredentialsMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find nsx credentials data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["username"]
		if !ok {
			return fmt.Errorf("can't find 'username' attribute in data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["password"]
		if !ok {
			return fmt.Errorf("can't find 'password' attribute in data source: %s", n)
		}
		return nil
	}
}

func testAccCheckGoogleVmwareengineVcenterCredentialsMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vcenter credentials data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["username"]
		if !ok {
			return fmt.Errorf("can't find 'username' attribute in data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["password"]
		if !ok {
			return fmt.Errorf("can't find 'password' attribute in data source: %s", n)
		}
		return nil
	}
}

func testAccCheckVmwareenginePrivateCloudDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vmwareengine_private_cloud" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}
			config := acctest.GoogleProviderConfig(t)
			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{VmwareengineBasePath}}projects/{{project}}/locations/{{location}}/privateClouds/{{name}}")
			if err != nil {
				return err
			}
			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}
			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				pcState, ok := res["state"]
				if !ok {
					return fmt.Errorf("Unable to fetch state for existing VmwareenginePrivateCloud %s", url)
				}
				if pcState.(string) != "DELETED" {
					return fmt.Errorf("VmwareenginePrivateCloud still exists at %s", url)
				}
			}
		}
		return nil
	}
}

func TestAccVmwareenginePrivateCloud_tags(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "venpc-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "venpc-tagvalue", tagKey)

	venSuffix := acctest.RandString(t, 10)
	venName := "tf-test-ven-" + venSuffix
	pcSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": pcSuffix,
		"org":           org,
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"ven_suffix":    venSuffix,
		"ven_name":      venName,
		"zone":          "me-west1-b",
		// Management CIDR must not overlap with any on-prem or VPC subnets
		"management_cidr": "192.168.101.0/24",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareenginePrivateCloudDestroyProducer(t), // Assuming this exists
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareenginePrivateCloudTags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_vmwareengine_private_cloud.default", "tags.%"),
					testAccCheckVmwareenginePrivateCloudHasTagBindings(t),
				),
			},
			{
				ResourceName:      "google_vmwareengine_private_cloud.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name",
					"location",
					"network_config",
					"management_cluster",
					"tags",
				},
			},
		},
	})
}

func testAccVmwareenginePrivateCloudTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "ven" {
  name        = "%{ven_name}"
  location    = "global"
  type        = "STANDARD"
  description = "Terraform test network for PC"
}

resource "google_vmwareengine_private_cloud" "default" {
  name     = "tf-test-pc-%{random_suffix}"
  location = "%{zone}"

  network_config {
    vmware_engine_network = google_vmwareengine_network.ven.id
    management_cidr     = "%{management_cidr}"
  }

  management_cluster {
    cluster_id = "tf-test-cluster-%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 3
    }
    # Add other node_type_configs blocks here if needed for other node types
  }
  # description = "Optional description"

  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }

}
`, context)
}

func testAccCheckVmwareenginePrivateCloudHasTagBindings(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vmwareengine_private_cloud" {
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

			// 3. Get the tag bindings from the VMware Engine Private Cloud.
			// ID format: projects/{project}/locations/{zone}/privateClouds/{privateCloudId}
			parts := strings.Split(rs.Primary.ID, "/")
			if len(parts) != 6 {
				return fmt.Errorf("invalid resource ID format: %s", rs.Primary.ID)
			}
			zone := parts[3]
			region, err := zoneToRegion(zone)
			if err != nil {
				return fmt.Errorf("could not determine region from zone %s: %v", zone, err)
			}

			parentURL := fmt.Sprintf("//vmwareengine.googleapis.com/%s", rs.Primary.ID)
			// Private Cloud is zonal, but TagBindings API is regional.
			listBindingsURL := fmt.Sprintf("https://%s-cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", region, url.QueryEscape(parentURL))

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

// zoneToRegion extracts the region from a zone name.
// Example: "us-central1-a" -> "us-central1"
func zoneToRegion(zone string) (string, error) {
	parts := strings.Split(zone, "-")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid zone format: %s", zone)
	}
	return strings.Join(parts[0:2], "-"), nil
}

