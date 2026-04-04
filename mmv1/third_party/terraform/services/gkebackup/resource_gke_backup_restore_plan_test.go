package gkebackup_test

import (
	"fmt"
	"testing"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccGKEBackupRestorePlan_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":             envvar.GetTestProjectFromEnv(),
		"deletion_protection": false,
		"network_name":        acctest.BootstrapSharedTestNetwork(t, "gke-cluster"),
		"subnetwork_name":     acctest.BootstrapSubnet(t, "gke-cluster", acctest.BootstrapSharedTestNetwork(t, "gke-cluster")),
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

func TestAccGkeBackupRestorePlan_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "gkebackup-rp-tagkey", map[string]interface{}{})
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestOrganizationTagValue(t, "gkebackup-rp-tagvalue", tagKey),
		"cluster_name":  "tf-test-cluster-" + acctest.RandString(t, 10),
		"location": "us-east1",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeBackupRestorePlanTags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_gke_backup_restore_plan.restore_plan", "tags.%"),
					testAccCheckGkeBackupRestorePlanHasTagBindings(t),
				),
			},
			{
				ResourceName:            "google_gke_backup_restore_plan.restore_plan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels", "tags"},
			},
		},
	})
}

func testAccCheckGkeBackupRestorePlanHasTagBindings(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_gke_backup_restore_plan" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

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
				return fmt.Errorf("invalid resource ID format for GKE Backup Plan: %s", rs.Primary.ID)
			}
			location := parts[3]

			parentURL := fmt.Sprintf("//gkebackup.googleapis.com/%s", rs.Primary.ID)
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

func testAccGkeBackupRestorePlanTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-gke-backup-cluster-%{random_suffix}"
  location           = "us-east1"
  initial_node_count = 1
  deletion_protection = false
  network            = "samathews-test"
}

resource "google_gke_backup_backup_plan" "basic" {
  name = "tf-test-restore-plan%{random_suffix}"
  cluster = google_container_cluster.primary.id
  location = "us-east1"
  backup_config {
    include_volume_data = true
    include_secrets = true
    all_namespaces = true
  }
}

resource "google_gke_backup_restore_plan" "restore_plan" {
  name = "tf-test-restore-plan%{random_suffix}"
  location = "us-east1"
  backup_plan = google_gke_backup_backup_plan.basic.id
  cluster = google_container_cluster.primary.id
  tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
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
