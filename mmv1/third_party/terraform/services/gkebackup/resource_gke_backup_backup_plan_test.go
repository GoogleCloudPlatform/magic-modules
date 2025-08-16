package gkebackup_test

import (
	"testing"
	"context"
	gkebackup "cloud.google.com/go/gkebackup/apiv1"
	gkebackuppb "cloud.google.com/go/gkebackup/apiv1/gkebackuppb"
	"fmt"
	
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGKEBackupBackupPlan_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":         envvar.GetTestProjectFromEnv(),
		"random_suffix":   acctest.RandString(t, 10),
		"network_name":    acctest.BootstrapSharedTestNetwork(t, "gke-cluster"),
		"subnetwork_name": acctest.BootstrapSubnet(t, "gke-cluster", acctest.BootstrapSharedTestNetwork(t, "gke-cluster")),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEBackupBackupPlanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEBackupBackupPlan_basic(context),
			},
			{
				ResourceName:            "google_gke_backup_backup_plan.backupplan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccGKEBackupBackupPlan_permissive(context),
			},
			{
				ResourceName:            "google_gke_backup_backup_plan.backupplan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccGKEBackupBackupPlan_full(context),
			},
			{
				ResourceName:            "google_gke_backup_backup_plan.backupplan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccGKEBackupBackupPlan_rpo_daily_window(context),
			},
			{
				ResourceName:            "google_gke_backup_backup_plan.backupplan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccGKEBackupBackupPlan_rpo_weekly_window(context),
			},
			{
				ResourceName:            "google_gke_backup_backup_plan.backupplan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccGKEBackupBackupPlan_full(context),
			},
			{
				ResourceName:            "google_gke_backup_backup_plan.backupplan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccGKEBackupBackupPlan_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-testcluster%{random_suffix}"
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
  deletion_protection = false
  network       = "%{network_name}"
  subnetwork    = "%{subnetwork_name}"
}

resource "google_gke_backup_backup_plan" "backupplan" {
  name = "tf-test-testplan%{random_suffix}"
  cluster = google_container_cluster.primary.id
  location = "us-central1"
  backup_config {
    include_volume_data = false
    include_secrets = false
    all_namespaces = true
  }
  labels = {
    "some-key-1": "some-value-1"
  }
}
`, context)
}

func testAccGKEBackupBackupPlan_permissive(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-testcluster%{random_suffix}"
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
  deletion_protection = false
  network       = "%{network_name}"
  subnetwork    = "%{subnetwork_name}"
}

resource "google_gke_backup_backup_plan" "backupplan" {
  name = "tf-test-testplan%{random_suffix}"
  cluster = google_container_cluster.primary.id
  location = "us-central1"
  backup_config {
    include_volume_data = false
    include_secrets = false
    all_namespaces = true
    permissive_mode = true
  }
  labels = {
    "some-key-1": "some-value-1"
  }
}
`, context)
}

func testAccGKEBackupBackupPlan_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-testcluster%{random_suffix}"
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
  deletion_protection = false
  network       = "%{network_name}"
  subnetwork    = "%{subnetwork_name}"
}
	
resource "google_gke_backup_backup_plan" "backupplan" {
  name = "tf-test-testplan%{random_suffix}"
  cluster = google_container_cluster.primary.id
  location = "us-central1"
  retention_policy {
    backup_delete_lock_days = 30
    backup_retain_days = 180
  }
  backup_schedule {
    cron_schedule = "0 9 * * 1"
  }
  backup_config {
    include_volume_data = true
    include_secrets = true
    selected_applications {
      namespaced_names {
        name = "app1"
        namespace = "ns1"
      }
      namespaced_names {
        name = "app2"
        namespace = "ns2"
      }
    }
  }
  labels = {
    "some-key-2": "some-value-2"
  }
}
`, context)
}

func testAccGKEBackupBackupPlan_rpo_daily_window(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-testcluster%{random_suffix}"
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
  deletion_protection = false
  network       = "%{network_name}"
  subnetwork    = "%{subnetwork_name}"
}
	
resource "google_gke_backup_backup_plan" "backupplan" {
  name = "tf-test-testplan%{random_suffix}"
  cluster = google_container_cluster.primary.id
  location = "us-central1"
  retention_policy {
    backup_delete_lock_days = 30
    backup_retain_days = 180
  }
  backup_schedule {
    paused = true
    rpo_config {
      target_rpo_minutes=1440
      exclusion_windows {
        start_time  {
          hours = 12
        }
        duration = "7200s"
        daily = true
      }
      exclusion_windows {
        start_time  {
          hours = 8
          minutes = 40
          seconds = 1
        }
        duration = "3600s"
        single_occurrence_date {
          year = 2024
          month = 3
          day = 16
        }
      }
    }
  }
  backup_config {
    include_volume_data = true
    include_secrets = true
    selected_applications {
      namespaced_names {
        name = "app1"
        namespace = "ns1"
      }
      namespaced_names {
        name = "app2"
        namespace = "ns2"
      }
    }
  }
  labels = {
    "some-key-2": "some-value-2"
  }
}
`, context)
}

func testAccGKEBackupBackupPlan_rpo_weekly_window(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-testcluster%{random_suffix}"
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
  deletion_protection = false
  network       = "%{network_name}"
  subnetwork    = "%{subnetwork_name}"
}
	
resource "google_gke_backup_backup_plan" "backupplan" {
  name = "tf-test-testplan%{random_suffix}"
  cluster = google_container_cluster.primary.id
  location = "us-central1"
  retention_policy {
    backup_delete_lock_days = 30
    backup_retain_days = 180
  }
  backup_schedule {
    paused = true
    rpo_config {
      target_rpo_minutes=1400
      exclusion_windows {
        start_time  {
          hours = 1
          minutes = 23
        }
        duration = "1800s"
        days_of_week {
          days_of_week = ["MONDAY", "THURSDAY"]
        }
      }
      exclusion_windows {
        start_time  {
          hours = 12
        }
        duration = "3600s"
        single_occurrence_date {
          year = 2024
          month = 3
          day = 17
        }
      }
      exclusion_windows {
        start_time  {
          hours = 8
          minutes = 40
        }
        duration = "600s"
        single_occurrence_date {
          year = 2024
          month = 3
          day = 18
        }
      }
    }
  }
  backup_config {
    include_volume_data = true
    include_secrets = true
    selected_applications {
      namespaced_names {
        name = "app1"
        namespace = "ns1"
      }
      namespaced_names {
        name = "app2"
        namespace = "ns2"
      }
    }
  }
  labels = {
    "some-key-2": "some-value-2"
  }
}
`, context)
}

func TestAccGKEBackupBackupPlan_tags(t *testing.T) {
	t.Parallel()

	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "gkebackup-bptagkey", map[string]interface{}{})
    	tagValue := acctest.BootstrapSharedTestOrganizationTagValue(t, "gkebackup-bptagvalue", tagKey)
    	prjLabel := "gke-tags"
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
	resourceName := "google_gke_backup_backup_plan.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEBackupBackupPlanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEBackupBackupPlanTags(testContext),
				Check: resource.ComposeTestCheckFunc(
					checkGKEBackupBackupPlanTags(resourceName, testContext),
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

func testAccGKEBackupBackupPlanTags(testContext map[string]interface{}) string {
    return acctest.Nprintf(`
resource "google_container_cluster" "primary_tags" {
  name               = "tf-test-tags-%{random_suffix}"
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

resource "google_gke_backup_backup_plan" "test" {
  project  = "%{project}"
  name     = "tf-test-plan-tags-%{random_suffix}"
  cluster  = google_container_cluster.primary_tags.id
  location = "us-central1"

  # Basic backup config
  backup_config {
    include_volume_data = false
    include_secrets     = false
    all_namespaces      = true
  }

  # The tags to be tested
  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, testContext)
}

// This function gets the backup plan via the Gkebackup API and inspects its tags.
func checkGKEBackupBackupPlanTags(resourceName string, testContext map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// Get resource attributes from state
		project := rs.Primary.Attributes["project"]
		location := rs.Primary.Attributes["location"]
		backupPlanName := rs.Primary.Attributes["name"]

		// Construct the expected full tag key
		expectedTagKey := fmt.Sprintf("%s/%s", testContext["org"], testContext["tagKey"])
		expectedTagValue := fmt.Sprintf("%s", testContext["tagValue"])

		// This `ctx` variable is now a `context.Context` object
		ctx := context.Background()

		// Create a Gkebackup client
		gkebackupClient, err := gkebackup.NewBackupForGKEClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create gkebackup client: %v", err)
		}
		defer gkebackupClient.Close()

		// Construct the request to get the backup plan details
		req := &gkebackuppb.GetBackupPlanRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/backupPlans/%s", project, location, backupPlanName),
		}

		// Get the Gkebackup backupplan
		backupPlan, err := gkebackupClient.GetBackupPlan(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get backup plan '%s': %v", req.Name, err)
		}

		// Check the instance's labels for the expected tag
		// In the Gkebackup API, tags are represented as labels.
		labels := backupPlan.GetLabels()
		if labels == nil {
			return fmt.Errorf("expected labels not found on backup plan '%s'", req.Name)
		}

		if actualValue, ok := labels[expectedTagKey]; ok {
			if actualValue == expectedTagValue {
				// The tag was found with the correct value. Success!
				return nil
			}
			return fmt.Errorf("tag key '%s' found with incorrect value. Expected: %s, Got: %s", expectedTagKey, expectedTagValue, actualValue)
		}

		// If we reach here, the tag key was not found.
		return fmt.Errorf("expected tag key '%s' not found on backup plan '%s'", expectedTagKey, req.Name)
	}
}
