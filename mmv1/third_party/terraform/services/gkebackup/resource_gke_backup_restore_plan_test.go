package gkebackup_test

import (
	"testing"
	"context"
	"fmt"

	gkebackup "cloud.google.com/go/gkebackup/apiv1"
	gkebackuppb "cloud.google.com/go/gkebackup/apiv1/gkebackuppb"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
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

func TestAccGKEBackupRestorePlan_tags(t *testing.T) {
	t.Parallel()

	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "gkebackup-rptagkey", map[string]interface{}{})
	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "gkebackup-rptagvalue", tagKey)
	prjLabel := "gke-rp-tags"
	networkName := acctest.BootstrapSharedTestNetwork(t, prjLabel)

	testContext := map[string]interface{}{
		"project":         envvar.GetTestProjectFromEnv(),
		"org":             envvar.GetTestOrgFromEnv(t),
		"tagKey":          tagKey,
		"tagValue":        tagValue,
		"random_suffix":   acctest.RandString(t, 10),
		"network_name":    networkName,
		"subnetwork_name": acctest.BootstrapSubnet(t, prjLabel, networkName),
	}
	resourceName := "google_gke_backup_restore_plan.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		// You'll need a destroy check function for RestorePlans
		// CheckDestroy:             testAccCheckGKEBackupRestorePlanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEBackupRestorePlanTags(testContext),
				Check: resource.ComposeTestCheckFunc(
					checkGKEBackupRestorePlanTags(resourceName, testContext),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func testAccGKEBackupRestorePlanTags(testContext map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary_rp_tags" {
  name               = "tf-test-rp-tags-%{random_suffix}"
  location           = "us-central1"
  initial_node_count = 1
  project            = "%{project}"

  # Enable Backup for GKE on the cluster
  addons_config {
    gke_backup_agent_config {
      enabled = true
    }
  }
  network    = "%{network_name}"
  subnetwork = "%{subnetwork_name}"
  deletion_protection = false
}

resource "google_gke_backup_backup_plan" "rp_test_dep" {
  project  = "%{project}"
  name     = "tf-test-bp-for-rp-%{random_suffix}"
  cluster  = google_container_cluster.primary_rp_tags.id
  location = "us-central1"
  backup_config {
    include_volume_data = true
    include_secrets = true
    all_namespaces = true
  }
}

resource "google_gke_backup_restore_plan" "test" {
  project  = "%{project}"
  name     = "tf-test-rp-tags-%{random_suffix}"
  location = "us-central1"
  backup_plan = google_gke_backup_backup_plan.rp_test_dep.id
  cluster = google_container_cluster.primary_rp_tags.id

  restore_config {
    all_namespaces = true
    cluster_resource_restore_scope {
      all_group_kinds = true
    }
    volume_data_restore_policy = "RESTORE_VOLUME_DATA_FROM_BACKUP"
    cluster_resource_conflict_policy = "USE_EXISTING_VERSION"
    namespaced_resource_restore_mode = "DELETE_AND_RESTORE"
  }

  # The tags to be tested
  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, testContext)
}

// This function gets the restore plan via the Gkebackup API and inspects its tags.
func checkGKEBackupRestorePlanTags(resourceName string, testContext map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// Get resource attributes from state
		project := rs.Primary.Attributes["project"]
		location := rs.Primary.Attributes["location"]
		restorePlanName := rs.Primary.Attributes["name"]

		// Construct the expected full tag key
		expectedTagKey := fmt.Sprintf("%s/%s", testContext["org"], testContext["tagKey"])
		expectedTagValue := fmt.Sprintf("%s", testContext["tagValue"])

		ctx := context.Background()

		// Create a Gkebackup client
		gkebackupClient, err := gkebackup.NewBackupForGKEClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create gkebackup client: %v", err)
		}
		defer gkebackupClient.Close()

		// Construct the request to get the restore plan details
		req := &gkebackuppb.GetRestorePlanRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/restorePlans/%s", project, location, restorePlanName),
		}

		// Get the Gkebackup restoreplan
		restorePlan, err := gkebackupClient.GetRestorePlan(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get restore plan '%s': %v", req.Name, err)
		}

		// Check the instance's labels for the expected tag
		// In the Gkebackup API, tags are represented as labels.
		labels := restorePlan.GetLabels()
		if labels == nil {
			return fmt.Errorf("expected labels not found on restore plan '%s'", req.Name)
		}

		if actualValue, ok := labels[expectedTagKey]; ok {
			if actualValue == expectedTagValue {
				// The tag was found with the correct value. Success!
				return nil
			}
			return fmt.Errorf("tag key '%s' found with incorrect value on restore plan '%s'. Expected: %s, Got: %s", expectedTagKey, req.Name, expectedTagValue, actualValue)
		}

		// If we reach here, the tag key was not found.
		return fmt.Errorf("expected tag key '%s' not found on restore plan '%s'", expectedTagKey, req.Name)
	}
}
