package gkeonprem_test

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

func TestAccGkeonpremVmwareCluster_vmwareClusterUpdateBasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremVmwareClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateMetalLbStart(context),
			},
			{
				ResourceName:            "google_gkeonprem_vmware_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateMetalLb(context),
			},
			{
				ResourceName:            "google_gkeonprem_vmware_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccGkeonpremVmwareCluster_vmwareClusterUpdateF5Lb(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremVmwareClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateF5LbStart(context),
			},
			{
				ResourceName:      "google_gkeonprem_vmware_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateF5lb(context),
			},
			{
				ResourceName:      "google_gkeonprem_vmware_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGkeonpremVmwareCluster_vmwareClusterUpdateManualLb(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremVmwareClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateManualLbStart(context),
			},
			{
				ResourceName:      "google_gkeonprem_vmware_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateManualLb(context),
			},
			{
				ResourceName:      "google_gkeonprem_vmware_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateMetalLbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster"
    on_prem_version = "1.13.1-gke.35"
    annotations = {
      env = "test"
    }
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/12"]
      pod_address_cidr_blocks = ["192.168.0.0/16"]
      dhcp_ip_config {
        enabled = true
      }
    }
    control_plane_node {
       cpus = 4
       memory = 8192
       replicas = 1
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.5"
        ingress_vip = "10.251.135.19"
      }
      metal_lb_config {
        address_pools {
          pool = "ingress-ip"
          manual_assign = "true"
          addresses = ["10.251.135.19"]
          avoid_buggy_ips = true
        }
        address_pools {
          pool = "lb-test-ip"
          manual_assign = "true"
          addresses = ["10.251.135.19"]
          avoid_buggy_ips = true
        }
      }
    }
  }
`, context)
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateMetalLb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster updated"
    on_prem_version = "1.13.1-gke.36"
    annotations = {
      env = "test-update"
    }
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/16"]
      pod_address_cidr_blocks = ["192.168.0.0/20"]
      dhcp_ip_config {
        enabled = true
      }
    }
    control_plane_node {
       cpus = 5
       memory = 4098
       replicas = 3
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.6"
        ingress_vip = "10.251.135.20"
      }
      metal_lb_config {
        address_pools {
          pool = "ingress-ip-updated"
          manual_assign = "false"
          addresses = ["10.251.135.20"]
          avoid_buggy_ips = false
        }
        address_pools {
          pool = "lb-test-ip-updated"
          manual_assign = "false"
          addresses = ["10.251.135.20"]
          avoid_buggy_ips = false
        }
      }
    }
  }
`, context)
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateF5LbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster"
    on_prem_version = "1.13.1-gke.35"
    annotations = {}
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/12"]
      pod_address_cidr_blocks = ["192.168.0.0/16"]
      dhcp_ip_config {
        enabled = true
      }
      control_plane_v2_config {
        control_plane_ip_block {
          ips {
            hostname = "test-hostname"
            ip = "10.0.0.1"
          }
          netmask="10.0.0.1/32"
          gateway="test-gateway"
        }
      }
    }
    control_plane_node {
       cpus = 4
       memory = 8192
       replicas = 1
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.5"
        ingress_vip = "10.251.135.19"
      }
      f5_config {
        address = "10.0.0.1"
        partition = "test-partition"
        snat_pool = "test-snap-pool"
      }
    }
  }
`, context)
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateF5lb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster"
    on_prem_version = "1.13.1-gke.35"
    annotations = {}
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/12"]
      pod_address_cidr_blocks = ["192.168.0.0/16"]
      dhcp_ip_config {
        enabled = true
      }
      control_plane_v2_config {
        control_plane_ip_block {
          ips {
            hostname = "test-hostname-updated"
            ip = "10.0.0.2"
          }
          netmask="10.0.0.2/32"
          gateway="test-gateway-updated"
        }
      }
    }
    control_plane_node {
       cpus = 4
       memory = 8192
       replicas = 1
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.5"
        ingress_vip = "10.251.135.19"
      }
      f5_config {
        address = "10.0.0.2"
        partition = "test-partition-updated"
        snat_pool = "test-snap-pool-updated"
      }
    }
  }
`, context)
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateManualLbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster"
    on_prem_version = "1.13.1-gke.35"
    annotations = {}
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/12"]
      pod_address_cidr_blocks = ["192.168.0.0/16"]
      host_config {
        dns_servers = ["10.254.41.1"]
        ntp_servers = ["216.239.35.8"]
        dns_search_domains = ["test-domain"]
      }
      static_ip_config {
        ip_blocks {
          netmask = "255.255.252.0"
          gateway = "10.251.31.254"
          ips {
            ip = "10.251.30.153"
            hostname = "test-hostname1"
          }
          ips {
            ip = "10.251.31.206"
            hostname = "test-hostname2"
          }
          ips {
            ip = "10.251.31.193"
            hostname = "test-hostname3"
          }
          ips { 
            ip = "10.251.30.230"
            hostname = "test-hostname4"
          }
        }
      }
    }
    control_plane_node {
       cpus = 4
       memory = 8192
       replicas = 1
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.5"
        ingress_vip = "10.251.135.19"
      }
      manual_lb_config {
        ingress_http_node_port = 30005
        ingress_https_node_port = 30006
        control_plane_node_port = 30007
        konnectivity_server_node_port = 30008
      }
    }
    vcenter {
      resource_pool = "test-resource-pool"
      datastore = "test-datastore"
      datacenter = "test-datacenter"
      cluster = "test-cluster"
      folder = "test-folder"
      ca_cert_data = "test-ca-cert-data"
      storage_policy_name = "test-storage-policy-name"
    }
    dataplane_v2 {
      dataplane_v2_enabled = true
      windows_dataplane_v2_enabled = true
      advanced_networking = true
    }
    vm_tracking_enabled = true
    enable_control_plane_v2 = true
    enable_advanced_cluster = true
    disable_bundled_ingress = true
    upgrade_policy {
      control_plane_only = true
    }
    authorization {
      admin_users {
        username = "testuser@gmail.com"
      }
    }
    anti_affinity_groups {
      aag_config_disabled = true
    }
    auto_repair_config {
      enabled = true
    }
  }
`, context)
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateManualLb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster"
    on_prem_version = "1.13.1-gke.35"
    annotations = {}
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/12"]
      pod_address_cidr_blocks = ["192.168.0.0/16"]
      host_config {
        dns_servers = ["10.254.41.1"]
        ntp_servers = ["216.239.35.8"]
        dns_search_domains = ["test-domain"]
      }
      static_ip_config {
        ip_blocks {
          netmask = "255.255.252.1"
          gateway = "10.251.31.255"
          ips {
            ip = "10.251.30.154"
            hostname = "test-hostname1-updated"
          }
          ips {
            ip = "10.251.31.206"
            hostname = "test-hostname2"
          }
          ips {
            ip = "10.251.31.193"
            hostname = "test-hostname3"
          }
          ips { 
            ip = "10.251.30.230"
            hostname = "test-hostname4"
          }
        }
      }
    }
    control_plane_node {
       cpus = 4
       memory = 8192
       replicas = 1
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.5"
        ingress_vip = "10.251.135.19"
      }
      manual_lb_config {
        ingress_http_node_port = 30006
        ingress_https_node_port = 30007
        control_plane_node_port = 30008
        konnectivity_server_node_port = 30009
      }
    }
    vcenter {
      resource_pool = "test-resource-pool-updated"
      datastore = "test-datastore-updated"
      datacenter = "test-datacenter-updated"
      cluster = "test-cluster-updated"
      folder = "test-folder-updated"
      ca_cert_data = "test-ca-cert-data-updated"
      storage_policy_name = "test-storage-policy-name-updated"
    }
    dataplane_v2 {
      dataplane_v2_enabled = true
      windows_dataplane_v2_enabled = true
      advanced_networking = true
    }
    vm_tracking_enabled = false
    disable_bundled_ingress = false
    upgrade_policy {
      control_plane_only = true
    }
    authorization {
      admin_users {
        username = "testuser-updated@gmail.com"
      }
    }
    anti_affinity_groups {
      aag_config_disabled = true
    }
    auto_repair_config {
      enabled = true
    }
  }
`, context)
}

func TestAccGkeOnPremVmwareCluster_tags(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "vmwctest-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "vmwctest-tagvalue", tagKey)

	// Fetch project from environment variable
	project := os.Getenv("GOOGLE_PROJECT")
	if project == "" {
		t.Skip("Skipping test: GOOGLE_PROJECT environment variable not set")
	}

	// Fetch location from environment variable, default to us-west1
	location := os.Getenv("GOOGLE_REGION")
	if location == "" {
		location = "us-west1"
	}

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           org,
		"tagKey":        tagKey,
		"tagValue":      tagValue,
		"project":       project,
		"location":      location,
		// Using the same admin cluster project as other tests, but a dynamic membership name
		"admin_cluster_membership": fmt.Sprintf("projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"),
		// ** IMPORTANT: These IPs must be valid and available in your test environment **
		"control_plane_vip": "10.251.133.7", // Example, adjust as needed
		"ingress_vip":       "10.251.135.21", // Example, adjust as needed
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremVmwareClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonPremVmwareClusterTagsConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_gkeonprem_vmware_cluster.default", "tags.%"),
					testAccCheckGkeOnPremVmwareClusterHasTagBindings(t),
				),
			},
			{
				ResourceName:            "google_gkeonprem_vmware_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

// HELPER FUNCTION to generate HCL for the tags test
func testAccGkeonPremVmwareClusterTagsConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gkeonprem_vmware_cluster" "default" {
  name                       = "tf-test-vmwc-tags-%{random_suffix}"
  project                    = "%{project}"
  location                   = "%{location}"
  admin_cluster_membership   = "%{admin_cluster_membership}"
  on_prem_version            = "1.13.1-gke.35" # Specify a valid version for your env

  network_config {
    service_address_cidr_blocks = ["10.96.0.0/12"]
    pod_address_cidr_blocks     = ["192.168.0.0/16"]
    dhcp_ip_config {
      enabled = true
    }
  }

  control_plane_node {
     cpus     = 4
     memory   = 8192
     replicas = 1
  }

  load_balancer {
    vip_config {
      control_plane_vip = "%{control_plane_vip}"
      ingress_vip       = "%{ingress_vip}"
    }
    metal_lb_config {
      address_pools {
        pool          = "ingress-ip"
        manual_assign = true
        addresses     = ["%{ingress_vip}"]
      }
    }
  }

  authorization {
    admin_users {
      username = "testacc@example.com"
    }
  }

  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}

// HELPER FUNCTION to check for tag bindings
func testAccCheckGkeOnPremVmwareClusterHasTagBindings(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_gkeonprem_vmware_cluster" {
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

			// 3. Get the tag bindings from the GKE OnPrem VMware Cluster.
			// ID format: projects/{project}/locations/{location}/vmwareClusters/{name}
			parentURL := fmt.Sprintf("//gkeonprem.googleapis.com/%s", rs.Primary.ID)
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
				return fmt.Errorf("expected tag value %s (from %q) not found in bindings for %s. No bindings returned", fullTagValueName, configuredTagValueNamespacedName, rs.Primary.ID)
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
