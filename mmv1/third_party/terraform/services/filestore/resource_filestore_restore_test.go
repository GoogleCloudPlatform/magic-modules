package filestore_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccFilestoreInstance_restore(t *testing.T) {
	t.Parallel()

	srcInstancetName := fmt.Sprintf("tf-fs-inst-source-%d", acctest.RandInt(t))
	restoreInstanceName := fmt.Sprintf("tf-fs-inst-restored-%d", acctest.RandInt(t))
	backupName := fmt.Sprintf("tf-fs-bkup-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstanceRestore_restore(srcInstancetName, restoreInstanceName, backupName),
			},
			{
				ResourceName:            "google_filestore_instance.instance_source",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccFilestoreInstanceRestore_restore(srcInstancetName, restoreInstanceName, backupName string) string {
	return fmt.Sprintf(`
	resource "google_filestore_instance" "instance_source" {
		name        = "%s"
		location    = "us-central1-b"
		tier        = "BASIC_HDD"
		description = "An instance created during testing."
	  
		file_shares {
		  capacity_gb = 1024
		  name        = "volume1"
		}
	  
		networks {
		  network      = "default"
		  modes        = ["MODE_IPV4"]
		  connect_mode = "DIRECT_PEERING"
		}
	}

	resource "google_filestore_instance" "instance_restored" {
		name        = "%s"
		location    = "us-central1-b"
		tier        = "BASIC_HDD"
		description = "An instance created during testing."
	  
		file_shares {
		  capacity_gb = 1024
		  name        = "volume1"
		  source_backup = google_filestore_backup.backup.id
		}
	  
		networks {
		  network      = "default"
		  modes        = ["MODE_IPV4"]
		  connect_mode = "DIRECT_PEERING"
		}
	}
	  
	resource "google_filestore_backup" "backup" {
		name        = "%s"
		location    = "us-central1"
		source_instance   = google_filestore_instance.instance_source.id
		source_file_share = "volume1"
	  
		description = "This is a filestore backup for the test instance"
	}
	  
	`, srcInstancetName, restoreInstanceName, backupName)
}

func TestAccFilestoreInstance_restoreBackupDR(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	backupName := fmt.Sprintf("tf-test-backup-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_restoreBackupDR(instanceName, backupName),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location"},
			},
		},
	})
}

func testAccFilestoreInstance_restoreBackupDR(instanceName string, backupName string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

locals {
  instance_name = "%s"
  backup_name = "%s"
}

resource "google_filestore_instance" "source_instance" {
  name     = "tf-source-instance-${local.instance_name}"
  location = "us-central1"
  tier     = "REGIONAL"

  file_shares {
    capacity_gb = 1024
    name        = "share"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
}

resource "google_backup_dr_backup_vault" "backup_vault" {
   location ="us-central1"
   backup_vault_id    = "tf-backup-vault-${local.backup_name}"
   description = "This is a second backup vault built by Terraform."
   backup_minimum_enforced_retention_duration = "100000s"
   force_update = "true"
   force_delete = "true"
   allow_missing = "true"
}

resource "google_backup_dr_backup_plan" "backup_plan" {
 location       = "us-central1"
 backup_plan_id = "tf-backup-plan-${local.backup_name}"
 resource_type  = "file.googleapis.com/Instance"
 backup_vault   = google_backup_dr_backup_vault.backup_vault.name

 backup_rules {
   rule_id                = "rule-1"
   backup_retention_days  = 2

    standard_schedule {
     recurrence_type     = "HOURLY"
     hourly_frequency    = 6
     time_zone           = "UTC"

     backup_window {
       start_hour_of_day = 12
       end_hour_of_day   = 18
     }
    }
 }
}

resource "google_backup_dr_backup_plan_association" "backcup_association" {
 location = "us-central1"
 backup_plan_association_id = "tf-backup-plan-association-${local.backup_name}"
 resource = google_filestore_instance.source_instance.id
 resource_type= "file.googleapis.com/Instance"
 backup_plan = google_backup_dr_backup_plan.backup_plan.name
 depends_on = [ google_filestore_instance.source_instance ]
}

data "google_backup_dr_data_source_references" "all_refs" {
	project       = data.google_project.project.project_id
	location      = "us-central1"
	resource_type = "file.googleapis.com/Instance"
	depends_on    = [google_backup_dr_backup_plan_association.backcup_association]
}


resource "google_filestore_instance" "instance" {
  name     = "tf-restored-instance-${local.instance_name}"
  location = "us-central1"
  tier     = "REGIONAL"

  file_shares {
    capacity_gb = 1024
    name        = "share"
    source_backupdr_backup = data.google_backup_dr_data_source_references.all_refs.data_source_references[0].name
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }

  depends_on = [data.google_backup_dr_data_source_references.all_refs]  
}
`, instanceName, backupName)
}
