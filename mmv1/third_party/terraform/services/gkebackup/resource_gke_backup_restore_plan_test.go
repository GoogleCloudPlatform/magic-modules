package gkebackup_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/container"
	_ "github.com/hashicorp/terraform-provider-google/google/services/gkebackup"
	"github.com/hashicorp/terraform-provider-google/google/services/tags"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccGKEBackupRestorePlan_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":             envvar.GetTestProjectFromEnv(),
		"deletion_protection": false,
		"network_name":        tpgcompute.BootstrapSharedTestNetwork(t, "gke-cluster"),
		"subnetwork_name":     tpgcompute.BootstrapSubnet(t, "gke-cluster", tpgcompute.BootstrapSharedTestNetwork(t, "gke-cluster")),
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEBackupRestorePlan_full(context),
			},
			{
				ResourceName:            "google_gke_backup_restore_plan.restore_plan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
			{
				Config: testAccGKEBackupRestorePlan_update(context),
			},
			{
				ResourceName:            "google_gke_backup_restore_plan.restore_plan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccGKEBackupRestorePlan_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-restore-plan%{random_suffix}-cluster"
  location           = "us-central1"
  initial_node_count = 1
  workload_identity_config {
    workload_pool = "%{project}.svc.id.goog"
  }
  addons_config {
    gke_backup_agent_config {
      enabled = true
    }
  }
  deletion_protection  = "%{deletion_protection}"
  network       = "%{network_name}"
  subnetwork    = "%{subnetwork_name}"
}

resource "google_gke_backup_backup_plan" "basic" {
  name = "tf-test-restore-plan%{random_suffix}"
  cluster = google_container_cluster.primary.id
  location = "us-central1"
  backup_config {
    include_volume_data = true
    include_secrets = true
    all_namespaces = true
  }
}

resource "google_gke_backup_restore_plan" "restore_plan" {
  name = "tf-test-restore-plan%{random_suffix}"
  location = "us-central1"
  backup_plan = google_gke_backup_backup_plan.basic.id
  cluster = google_container_cluster.primary.id
  restore_config {
    all_namespaces = true
    namespaced_resource_restore_mode = "MERGE_SKIP_ON_CONFLICT"
    volume_data_restore_policy = "RESTORE_VOLUME_DATA_FROM_BACKUP"
    cluster_resource_restore_scope {
      all_group_kinds = true
    }
    cluster_resource_conflict_policy = "USE_EXISTING_VERSION"
    restore_order {
        group_kind_dependencies {
            satisfying {
                resource_group = "stable.example.com"
                resource_kind = "kindA"
            }
            requiring {
                resource_group = "stable.example.com"
                resource_kind = "kindB"
            }
        }
        group_kind_dependencies {
            satisfying {
                resource_group = "stable.example.com"
                resource_kind = "kindB"
            }
            requiring {
                resource_group = "stable.example.com"
                resource_kind = "kindC"
            }
        }
    }
    volume_data_restore_policy_bindings {
        policy = "RESTORE_VOLUME_DATA_FROM_BACKUP"
        volume_type = "GCE_PERSISTENT_DISK"
    }
  }
}
`, context)
}

func testAccGKEBackupRestorePlan_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-restore-plan%{random_suffix}-cluster"
  location           = "us-central1"
  initial_node_count = 1
  workload_identity_config {
    workload_pool = "%{project}.svc.id.goog"
  }
  addons_config {
    gke_backup_agent_config {
      enabled = true
    }
  }
  deletion_protection  = "%{deletion_protection}"
  network       = "%{network_name}"
  subnetwork    = "%{subnetwork_name}"
}

resource "google_gke_backup_backup_plan" "basic" {
  name = "tf-test-restore-plan%{random_suffix}"
  cluster = google_container_cluster.primary.id
  location = "us-central1"
  backup_config {
    include_volume_data = true
    include_secrets = true
    all_namespaces = true
  }
}

resource "google_gke_backup_restore_plan" "restore_plan" {
  name = "tf-test-restore-plan%{random_suffix}"
  location = "us-central1"
  backup_plan = google_gke_backup_backup_plan.basic.id
  cluster = google_container_cluster.primary.id
  restore_config {
    all_namespaces = true
    namespaced_resource_restore_mode = "MERGE_REPLACE_VOLUME_ON_CONFLICT"
    volume_data_restore_policy = "RESTORE_VOLUME_DATA_FROM_BACKUP"
    cluster_resource_restore_scope {
      all_group_kinds = true
    }
    cluster_resource_conflict_policy = "USE_EXISTING_VERSION"
    restore_order {
        group_kind_dependencies {
            satisfying {
                resource_group = "stable.example.com"
                resource_kind = "kindA"
            }
            requiring {
                resource_group = "stable.example.com"
                resource_kind = "kindB"
            }
        }
        group_kind_dependencies {
            satisfying {
                resource_group = "stable.example.com"
                resource_kind = "kindB"
            }
            requiring {
                resource_group = "stable.example.com"
                resource_kind = "kindC"
            }
        }
        group_kind_dependencies {
            satisfying {
                resource_group = "stable.example.com"
                resource_kind = "kindC"
            }
            requiring {
                resource_group = "stable.example.com"
                resource_kind = "kindD"
            }
        }
    }
    volume_data_restore_policy_bindings {
      policy = "REUSE_VOLUME_HANDLE_FROM_BACKUP"
      volume_type = "GCE_PERSISTENT_DISK"
    }
  }
}
`, context)
}

func TestAccGKEBackupRestorePlan_tags(t *testing.T) {
	t.Parallel()

	tagKey := tags.BootstrapSharedTestOrganizationTagKey(t, "gkebackup-restoreplan-tagkey", map[string]interface{}{})

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"org":           envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  tpgcompute.BootstrapSharedTestNetwork(t, "gke-cluster"),
		"tagKey":        tagKey,
		"tagValue":      tags.BootstrapSharedTestOrganizationTagValue(t, "gkebackup-restoreplan-tagvalue", tagKey),
		"region":        "us-east1",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEBackupRestorePlan_tags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_gke_backup_restore_plan.restore_plan", "tags.%"),
					testAccCheckGKEBackupRestorePlanHasTagBindings(t),
				),
			},
			{
				ResourceName:            "google_gke_backup_restore_plan.restore_plan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccGKEBackupRestorePlan_tags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_subnetwork" "primary" {
  name          = "tf-test-subnet%{random_suffix}"
  ip_cidr_range = "10.2.0.0/20"
  region        = "%{region}"
  network       = "%{network_name}"
}

resource "google_container_cluster" "primary" {
  name               = "tf-test-testcluster%{random_suffix}"
  location           = "%{region}"
  initial_node_count = 1
  workload_identity_config {
    workload_pool = "%{project}.svc.id.goog"
  }
  addons_config {
    gke_backup_agent_config {
      enabled = true
    }
  }
  deletion_protection = false
  network       = "%{network_name}"
  subnetwork    = google_compute_subnetwork.primary.name
}

resource "google_gke_backup_backup_plan" "basic" {
  name = "tf-test-restore-plan%{random_suffix}"
  cluster = google_container_cluster.primary.id
  location = "%{region}"
  backup_config {
    include_volume_data = true
    include_secrets = true
    all_namespaces = true
  }
}

resource "google_gke_backup_restore_plan" "restore_plan" {
  name = "tf-test-restore-plan%{random_suffix}"
  location = "%{region}"
  backup_plan = google_gke_backup_backup_plan.basic.id
  cluster = google_container_cluster.primary.id
  restore_config {
    all_namespaces = true
    namespaced_resource_restore_mode = "MERGE_SKIP_ON_CONFLICT"
    volume_data_restore_policy = "RESTORE_VOLUME_DATA_FROM_BACKUP"
    cluster_resource_restore_scope {
      all_group_kinds = true
    }
    cluster_resource_conflict_policy = "USE_EXISTING_VERSION"
  }
  tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}

func testAccCheckGKEBackupRestorePlanHasTagBindings(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_gke_backup_restore_plan" {
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

			// 3. Check if tags are public (returned in resource GET response)
			parts := strings.Split(rs.Primary.ID, "/")
			if len(parts) != 6 {
				return fmt.Errorf("invalid resource ID format: %s", rs.Primary.ID)
			}
			project := parts[1]
			location := parts[3]
			planName := parts[5]

			resourceURL := fmt.Sprintf("https://gkebackup.googleapis.com/v1/projects/%s/locations/%s/restorePlans/%s", project, location, planName)
			respResource, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    resourceURL,
				UserAgent: config.UserAgent,
			})
			if err != nil {
				return fmt.Errorf("error calling RestorePlan GET API: %v", err)
			}

			if tagsVal, exists := respResource["tags"]; exists {
				t.Logf("Tags are public and returned in RestorePlan GET response: %v", tagsVal)
			} else {
				t.Logf("Tags are NOT returned in RestorePlan GET response")
			}

			// 4. Get the tag bindings from TagBindings API.
			parentURL := fmt.Sprintf("//gkebackup.googleapis.com/projects/%s/locations/%s/restorePlans/%s", project, location, planName)
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

			// 5. Perform the comparison.
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
