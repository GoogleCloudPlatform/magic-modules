package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGKEBackupBackupPlan_update(t *testing.T) {
	t.Parallel()

	random_suffix := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGKEBackupBackupPlanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEBackupBackupPlan_basic(random_suffix),
			},
			{
				ResourceName:      "google_gke_backup_backup_plan.backupplan",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEBackupBackupPlan_full(random_suffix),
			},
			{
				ResourceName:      "google_gke_backup_backup_plan.backupplan",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEBackupBackupPlan_basic(random_suffix string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  provider = google-beta
  name               = "testcluster-%s"
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
}
	
resource "google_gke_backup_backup_plan" "backupplan" {
  provider = google-beta
  name = "testplan%s"
  cluster = google_container_cluster.primary.id
    backup_config {
	  include_volume_data = false
	  include_secrets = false
	  all_namespaces = true
	}
}
`, random_suffix)
}

func testAccGKEBackupBackupPlan_full(random_suffix string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  provider = google-beta
  name               = "fullcluster-%s"
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
}
	
resource "google_gke_backup_backup_plan" "backupplan" {
  provider = google-beta
  name = "fullplan%s"
  cluster = google_container_cluster.primary.id
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
}
`, random_suffix)
}