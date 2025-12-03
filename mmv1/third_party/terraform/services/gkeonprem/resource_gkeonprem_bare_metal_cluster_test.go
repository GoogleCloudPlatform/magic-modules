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

func TestAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremBareMetalClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateMetalLbStart(context),
			},
			{
				ResourceName:            "google_gkeonprem_bare_metal_cluster.cluster-metallb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateMetalLb(context),
			},
			{
				ResourceName:            "google_gkeonprem_bare_metal_cluster.cluster-metallb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccGkeonpremBareMetalCluster_bareMetalClusterUpdateManualLb(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremBareMetalClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateManualLbStart(context),
			},
			{
				ResourceName:      "google_gkeonprem_bare_metal_cluster.cluster-manuallb",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateManualLb(context),
			},
			{
				ResourceName:      "google_gkeonprem_bare_metal_cluster.cluster-manuallb",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBgpLb(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremBareMetalClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBgpLbStart(context),
			},
			{
				ResourceName:      "google_gkeonprem_bare_metal_cluster.cluster-bgplb",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBgpLb(context),
			},
			{
				ResourceName:      "google_gkeonprem_bare_metal_cluster.cluster-bgplb",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateMetalLbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-metallb" {
    name = "cluster-metallb%{random_suffix}"
    location = "us-west1"
    annotations = {
      env = "test"
    }
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/16"]
        pod_address_cidr_blocks = ["10.240.0.0/13"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.9"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 443
      }
      vip_config {
        control_plane_vip = "10.200.0.13"
        ingress_vip = "10.200.0.14"
      }
      metal_lb_config {
        address_pools {
          pool = "pool1"
          addresses = [
            "10.200.0.14/32",
            "10.200.0.15/32",
            "10.200.0.16/32",
            "10.200.0.17/32",
            "10.200.0.18/32",
            "fd00:1::f/128",
            "fd00:1::10/128",
            "fd00:1::11/128",
            "fd00:1::12/128"
          ]
        }
      }
    }
    storage {
      lvp_share_config {
        lvp_config {
          path = "/mnt/localpv-share"
          storage_class = "local-shared"
        }
        shared_path_pv_count = 5
      }
      lvp_node_mounts_config {
        path = "/mnt/localpv-disk"
        storage_class = "local-disks"
      }
    }
    security_config {
      authorization {
        admin_users {
          username = "admin@hashicorptest.com"
        }
      }
    }
  }
`, context)
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateMetalLb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-metallb" {
    name = "cluster-metallb%{random_suffix}"
    location = "us-west1"
    annotations = {
      env = "test-update"
    }
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/20"]
        pod_address_cidr_blocks = ["10.240.0.0/14"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.10"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 80
      }
      vip_config {
        control_plane_vip = "10.200.0.14"
        ingress_vip = "10.200.0.15"
      }
      metal_lb_config {
        address_pools {
          pool = "pool2"
          addresses = [
            "10.200.0.14/32",
            "10.200.0.15/32",
            "10.200.0.16/32",
            "10.200.0.17/32",
            "fd00:1::f/128",
            "fd00:1::10/128",
            "fd00:1::11/128"
          ]
        }
      }
    }
    storage {
      lvp_share_config {
        lvp_config {
          path = "/mnt/localpv-share-updated"
          storage_class = "local-shared-updated"
        }
        shared_path_pv_count = 6
      }
      lvp_node_mounts_config {
        path = "/mnt/localpv-disk-updated"
        storage_class = "local-disks-updated"
      }
    }
    security_config {
      authorization {
        admin_users {
          username = "admin-updated@hashicorptest.com"
        }
      }
    }
  }
`, context)
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateManualLbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-manuallb" {
    name = "cluster-manuallb%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/20"]
        pod_address_cidr_blocks = ["10.240.0.0/14"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.10"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 80
      }
      vip_config {
        control_plane_vip = "10.200.0.13"
        ingress_vip = "10.200.0.14"
      }
      metal_lb_config {
        address_pools {
          pool = "pool2"
          addresses = [
            "10.200.0.14/32",
            "10.200.0.15/32",
            "10.200.0.16/32",
            "10.200.0.17/32",
            "fd00:1::f/128",
            "fd00:1::10/128",
            "fd00:1::11/128"
          ]
        }
      }
    }
    storage {
      lvp_share_config {
        lvp_config {
          path = "/mnt/localpv-share"
          storage_class = "local-shared"
        }
        shared_path_pv_count = 6
      }
      lvp_node_mounts_config {
        path = "/mnt/localpv-disk"
        storage_class = "local-disks"
      }
    }
    security_config {
      authorization {
        admin_users {
          username = "admin@hashicorptest.com"
        }
      }
    }
    binary_authorization {
      evaluation_mode = "DISABLED"
    }
    upgrade_policy {
      policy = "SERIAL"
    }
  }
`, context)
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateManualLb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-manuallb" {
    name = "cluster-manuallb%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/20"]
        pod_address_cidr_blocks = ["10.240.0.0/14"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.10"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 80
      }
      vip_config {
        control_plane_vip = "10.200.0.14"
        ingress_vip = "10.200.0.15"
      }
      manual_lb_config {
        enabled = true
      }
    }
    storage {
      lvp_share_config {
        lvp_config {
          path = "/mnt/localpv-share-updated"
          storage_class = "local-shared-updated"
        }
        shared_path_pv_count = 6
      }
      lvp_node_mounts_config {
        path = "/mnt/localpv-disk-updated"
        storage_class = "local-disks-updated"
      }
    }
    security_config {
      authorization {
        admin_users {
          username = "admin-updated@hashicorptest.com"
        }
      }
    }
    binary_authorization {
      evaluation_mode = "PROJECT_SINGLETON_POLICY_ENFORCE"
    }
    upgrade_policy {
      policy = "CONCURRENT"
    }
  }
`, context)
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBgpLbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-bgplb" {
    name = "cluster-bgplb%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/20"]
        pod_address_cidr_blocks = ["10.240.0.0/14"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.10"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 80
      }
      vip_config {
        control_plane_vip = "10.200.0.13"
        ingress_vip = "10.200.0.14"
      }
      bgp_lb_config {
        asn = 123456
        bgp_peer_configs {
          asn = 123457
          ip_address = "10.0.0.1"
          control_plane_nodes = ["test-node"]
        }
        address_pools {
          pool = "pool1"
          addresses = [
            "10.200.0.14/32",
            "fd00:1::12/128"
          ]
          manual_assign = true
        }
        load_balancer_node_pool_config {
          node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.9"
            }
          }
        }
      }
    }
    storage {
      lvp_share_config {
        lvp_config {
          path = "/mnt/localpv-share"
          storage_class = "local-shared"
        }
        shared_path_pv_count = 6
      }
      lvp_node_mounts_config {
        path = "/mnt/localpv-disk"
        storage_class = "local-disks"
      }
    }
    security_config {
      authorization {
        admin_users {
          username = "admin@hashicorptest.com"
        }
      }
    }
  }
`, context)
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBgpLb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-bgplb" {
    name = "cluster-bgplb%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/20"]
        pod_address_cidr_blocks = ["10.240.0.0/14"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.10"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 80
      }
      vip_config {
        control_plane_vip = "10.200.0.14"
        ingress_vip = "10.200.0.15"
      }
      bgp_lb_config {
        asn = 123457
        bgp_peer_configs {
          asn = 123458
          ip_address = "10.0.0.2"
          control_plane_nodes = ["test-node-updated"]
        }
        address_pools {
          pool = "pool2"
          addresses = [
            "10.200.0.15/32",
            "fd00:1::16/128"
          ]
        }
        load_balancer_node_pool_config {
          node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.11"
            }
          }
        }
      }
    }
    storage {
      lvp_share_config {
        lvp_config {
          path = "/mnt/localpv-share-updated"
          storage_class = "local-shared-updated"
        }
        shared_path_pv_count = 6
      }
      lvp_node_mounts_config {
        path = "/mnt/localpv-disk-updated"
        storage_class = "local-disks-updated"
      }
    }
    security_config {
      authorization {
        admin_users {
          username = "admin-updated@hashicorptest.com"
        }
      }
    }
  }
`, context)
}

func TestAccGkeOnPremBareMetalCluster_tags(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "bmctest-tagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "bmctest-tagvalue", tagKey)
	
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
		"admin_cluster_membership": fmt.Sprintf("projects/%s/locations/global/memberships/gkeonprem-terraform-test", "870316890899"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		// *** Ensure this CheckDestroy function is implemented and uncommented ***
		CheckDestroy: testAccCheckGkeonpremBareMetalClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeOnPremBareMetalClusterTagsConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_gkeonprem_bare_metal_cluster.default", "tags.%"),
					testAccCheckGkeOnPremBareMetalClusterHasTagBindings(t),
				),
			},
			{
				ResourceName:            "google_gkeonprem_bare_metal_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"}, // Add any other fields that shouldn't be verified on import
			},
		},
	})
}

// testAccGkeOnPremBareMetalClusterTagsConfig returns the Terraform configuration string.
// ** REVIEW HARDCODED VALUES **: IPs and version must be valid in the test environment.
func testAccGkeOnPremBareMetalClusterTagsConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gkeonprem_bare_metal_cluster" "default" {
  name                      = "tf-test-bmc-%{random_suffix}"
  project                   = "%{project}"
  location                  = "%{location}"
  admin_cluster_membership  = "%{admin_cluster_membership}"
  bare_metal_version        = "1.15.0" # Ensure this version is appropriate

  control_plane {
    control_plane_node_pool_config { # Corrected block name
      node_pool_config {
        labels = {} // Changed from node_pool_labels
        node_configs {
          node_ip = "10.200.0.1" # Example IP - VERIFY THIS!
        }
      }
    }
  }

  load_balancer {
    vip_config {
      control_plane_vip = "10.200.0.10" # Example IP - VERIFY THIS!
      ingress_vip       = "10.200.0.11" # Example IP - VERIFY THIS!
    }
    port_config {
      control_plane_load_balancer_port = 443
    }
    # metal_lb_config or manual_lb_config or bgp_lb_config is required. Assuming MetalLB for this example.
    metal_lb_config {
      address_pools {
        pool = "pool1"
        addresses = [
          "10.200.0.11/32", // Ingress VIP must be in an address pool
          "10.200.0.12/32",
        ]
      }
    }
  }

  storage {
    lvp_share_config {
      lvp_config {
        path = "/mnt/localpv-share"
        storage_class = "local-shared"
      }
      shared_path_pv_count = 5
    }
    lvp_node_mounts_config {
      path = "/mnt/localpv-disk"
      storage_class = "local-disks"
    }
  }

  network_config {
    island_mode_cidr {
      service_address_cidr_blocks = ["10.96.0.0/12"]
      pod_address_cidr_blocks    = ["192.168.0.0/16"]
    }
  }

  security_config {
    authorization {
      admin_users {
        username = "testuser@example.com" # Added example admin user
      }
    }
  }

  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}

// testAccCheckGkeOnPremBareMetalClusterHasTagBindings checks if the correct tag bindings exist.
func testAccCheckGkeOnPremBareMetalClusterHasTagBindings(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_gkeonprem_bare_metal_cluster" {
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

			// 3. Get the tag bindings from the GKE OnPrem Bare Metal Cluster.
			// The ID format for google_gkeonprem_bare_metal_cluster is projects/{project}/locations/{location}/bareMetalClusters/{name}
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
				// No bindings found, return an error as we expect the tag to be there.
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
