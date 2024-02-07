// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetappVolume_netappVolumeBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "gcnv-network-config-1", acctest.ServiceNetworkWithParentService("netapp.servicenetworking.goog")),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappVolumeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappVolume_volumeBasicExample_basic(context),
			},
			{
				ResourceName:            "google_netapp_volume.test_volume",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"restore_parameters", "location", "name", "deletion_policy", "labels", "terraform_labels"},
			}, {
				Config: testAccNetappVolume_volumeBasicExample_full(context),
			},
			{
				ResourceName:            "google_netapp_volume.test_volume",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"restore_parameters", "location", "name", "deletion_policy", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappVolume_volumeBasicExample_update(context),
			},
			{
				ResourceName:            "google_netapp_volume.test_volume",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"restore_parameters", "location", "name", "deletion_policy", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappVolume_volumeBasicExample_updatesnapshot(context),
			},
			{
				ResourceName:            "google_netapp_volume.test_volume",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"restore_parameters", "location", "name", "deletion_policy", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappVolume_volumeBasicExample_createclonevolume(context),
			},
			{
				ResourceName:            "google_netapp_volume.test_volume_clone",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"restore_parameters", "location", "name", "deletion_policy", "labels", "terraform_labels"},
			},
			// backupConfig is an upcoming feature. Needs more work/testing before release.
			// 	{
			// 		Config: testAccNetappVolume_volumeBasicExample_setupBackup(context),
			// 	},
			// 	{
			// 		ResourceName:            "google_netapp_volume.test_volume",
			// 		ImportState:             true,
			// 		ImportStateVerify:       true,
			// 		ImportStateVerifyIgnore: []string{"restore_parameters", "location", "name", "deletion_policy", "labels", "terraform_labels"},
			// 	},
			// 	{
			// 		Config: testAccNetappVolume_volumeBasicExample_updateBackup(context),
			// 	},
			// 	{
			// 		ResourceName:            "google_netapp_volume.test_volume",
			// 		ImportState:             true,
			// 		ImportStateVerify:       true,
			// 		ImportStateVerifyIgnore: []string{"restore_parameters", "location", "name", "deletion_policy", "labels", "terraform_labels"},
			// 	},
		},
	})
}

func testAccNetappVolume_volumeBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_storage_pool" "default" {
	name = "tf-test-test-pool%{random_suffix}"
	location = "us-west2"
	service_level = "PREMIUM"
	capacity_gib = "2048"
	network = data.google_compute_network.default.id
}

resource "google_netapp_volume" "test_volume" {
	location = "us-west2"
	name = "tf-test-test-volume%{random_suffix}"
	capacity_gib = "100"
	share_name = "tf-test-test-volume%{random_suffix}"
	storage_pool = google_netapp_storage_pool.default.name
	protocols = ["NFSV3"]
}

data "google_compute_network" "default" {
	name = "%{network_name}"
}
`, context)
}

func testAccNetappVolume_volumeBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_storage_pool" "default" {
	name = "tf-test-test-pool%{random_suffix}"
	location = "us-west2"
	service_level = "PREMIUM"
	capacity_gib = "2048"
	network = data.google_compute_network.default.id
}
	
resource "google_netapp_storage_pool" "default2" {
	name = "tf-test-pool%{random_suffix}"
	location = "us-west2"
	service_level = "EXTREME"
	capacity_gib = "2048"
	network = data.google_compute_network.default.id
}
	  
resource "google_netapp_volume" "test_volume" {
	location = "us-west2"
	name = "tf-test-test-volume%{random_suffix}"
	capacity_gib = "100"
	share_name = "tf-test-test-volume%{random_suffix}"
	storage_pool = google_netapp_storage_pool.default.name
	protocols = ["NFSV3"]
	smb_settings = []
	unix_permissions = "0770"
	labels = {
		key= "test"
		value= "pool"
	}
	description = "This is a test description"
	snapshot_directory = false
	security_style = "UNIX"
	kerberos_enabled = false
	export_policy {
		rules {
			access_type           = "READ_ONLY"
			allowed_clients       = "0.0.0.0/0"
			has_root_access       = "false"
			kerberos5_read_only   = false
			kerberos5_read_write  = false
			kerberos5i_read_only  = false
			kerberos5i_read_write = false
			kerberos5p_read_only  = false
			kerberos5p_read_write = false
			nfsv3                 = true
			nfsv4                 = false
		}
		rules {
			access_type           = "READ_WRITE"
			allowed_clients       = "10.2.3.4,10.2.3.5"
			has_root_access       = "true"
			kerberos5_read_only   = false
			kerberos5_read_write  = false
			kerberos5i_read_only  = false
			kerberos5i_read_write = false
			kerberos5p_read_only  = false
			kerberos5p_read_write = false
			nfsv3                 = true
			nfsv4                 = false
		}
	}
	restricted_actions = []
	snapshot_policy {
		daily_schedule {
			snapshots_to_keep = 2
		}
		enabled = true
		hourly_schedule {
			snapshots_to_keep = 2
		}
		monthly_schedule {
			snapshots_to_keep = 4
		}
		weekly_schedule {
			snapshots_to_keep = 2
		}
	}
}

data "google_compute_network" "default" {
	name = "%{network_name}"
}
	  `, context)
}

func testAccNetappVolume_volumeBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_netapp_storage_pool" "default" {
	name = "tf-test-test-pool%{random_suffix}"
	location = "us-west2"
	service_level = "PREMIUM"
	capacity_gib = "2048"
	network = data.google_compute_network.default.id
}
  
resource "google_netapp_storage_pool" "default2" {
	name = "tf-test-pool%{random_suffix}"
	location = "us-west2"
	service_level = "EXTREME"
	capacity_gib = "2048"
	network = data.google_compute_network.default.id
}

resource "google_netapp_volume" "test_volume" {
	location = "us-west2"
	name = "tf-test-test-volume%{random_suffix}"
	capacity_gib = "200"
	share_name = "tf-test-test-volume%{random_suffix}"
	storage_pool = google_netapp_storage_pool.default2.name
	protocols = ["NFSV3"]
	smb_settings = []
	unix_permissions = "0740"
	labels = {}
	description = ""
	snapshot_directory = true
	security_style = "UNIX"
	kerberos_enabled = false
	export_policy {
		rules {
			access_type           = "READ_WRITE"
			allowed_clients       = "0.0.0.0/0"
			has_root_access       = "true"
			kerberos5_read_only   = false
			kerberos5_read_write  = false
			kerberos5i_read_only  = false
			kerberos5i_read_write = false
			kerberos5p_read_only  = false
			kerberos5p_read_write = false
			nfsv3                 = true
			nfsv4                 = false
		}
	}
	# Delete protection only gets active after an NFS client mounts.
	# Setting it here is save, volume can still be deleted.
	deletion_policy = "FORCE"
	snapshot_policy {
		enabled = true
		daily_schedule {
			hour              = 1
			minute            = 2
			snapshots_to_keep = 1
		}
		hourly_schedule {
			minute            = 10
			snapshots_to_keep = 1
		}
		monthly_schedule {
			days_of_month     = "2"
			hour              = 3
			minute            = 4
			snapshots_to_keep = 1
		}
		weekly_schedule {
			day               = "Monday"
			hour              = 5
			minute            = 6
			snapshots_to_keep = 1
		}
	}
}

data "google_compute_network" "default" {
	name = "%{network_name}"
}
  `, context)
}

func testAccNetappVolume_volumeBasicExample_updatesnapshot(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_netapp_storage_pool" "default2" {
	name = "tf-test-pool%{random_suffix}"
	location = "us-west2"
	service_level = "EXTREME"
	capacity_gib = "2048"
	network = data.google_compute_network.default.id
}
	
resource "google_netapp_volume" "test_volume" {
	location = "us-west2"
	name = "tf-test-test-volume%{random_suffix}"
	capacity_gib = "200"
	share_name = "tf-test-test-volume%{random_suffix}"
	storage_pool = google_netapp_storage_pool.default2.name
	protocols = ["NFSV3"]
	smb_settings = []
	unix_permissions = "0740"
	labels = {}
	description = ""
	snapshot_directory = true
	security_style = "UNIX"
	kerberos_enabled = false
	export_policy {
		rules {
			access_type           = "READ_WRITE"
			allowed_clients       = "0.0.0.0/0"
			has_root_access       = "true"
			kerberos5_read_only   = false
			kerberos5_read_write  = false
			kerberos5i_read_only  = false
			kerberos5i_read_write = false
			kerberos5p_read_only  = false
			kerberos5p_read_write = false
			nfsv3                 = true
			nfsv4                 = false
		}
	}
	# Delete protection only gets active after an NFS client mounts.
	# Setting it here is save, volume can still be deleted.
	restricted_actions = ["DELETE"]
	deletion_policy = "FORCE"
}

resource "google_netapp_volume_snapshot" "test-snapshot" {
	depends_on = [google_netapp_volume.test_volume]
	location = google_netapp_volume.test_volume.location
	volume_name = google_netapp_volume.test_volume.name
	name = "test-snapshot"
}

data "google_compute_network" "default" {
	name = "%{network_name}"
}
	`, context)
}

// Tests creating a new volume (clone) from a snapshot created from existing volume
func testAccNetappVolume_volumeBasicExample_createclonevolume(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_netapp_storage_pool" "default2" {
	name = "tf-test-pool%{random_suffix}"
	location = "us-west2"
	service_level = "EXTREME"
	capacity_gib = "2048"
	network = data.google_compute_network.default.id
}
	
resource "google_netapp_volume" "test_volume" {
	location = "us-west2"
	name = "tf-test-test-volume%{random_suffix}"
	capacity_gib = "200"
	share_name = "tf-test-test-volume%{random_suffix}"
	storage_pool = google_netapp_storage_pool.default2.name
	protocols = ["NFSV3"]
	smb_settings = []
	unix_permissions = "0740"
	labels = {}
	description = ""
	snapshot_directory = true
	security_style = "UNIX"
	kerberos_enabled = false
	export_policy {
		rules {
			access_type           = "READ_WRITE"
			allowed_clients       = "0.0.0.0/0"
			has_root_access       = "true"
			kerberos5_read_only   = false
			kerberos5_read_write  = false
			kerberos5i_read_only  = false
			kerberos5i_read_write = false
			kerberos5p_read_only  = false
			kerberos5p_read_write = false
			nfsv3                 = true
			nfsv4                 = false
		}
	}
	# Delete protection only gets active after an NFS client mounts.
	# Setting it here is save, volume can still be deleted.
	restricted_actions = ["DELETE"]
	deletion_policy = "FORCE"
}

resource "google_netapp_volume_snapshot" "test-snapshot" {
	depends_on = [google_netapp_volume.test_volume]
	location = google_netapp_volume.test_volume.location
	volume_name = google_netapp_volume.test_volume.name
	name = "test-snapshot"
}

resource "google_netapp_volume" "test_volume_clone" {
	location = "us-west2"
	name = "tf-test-test-volume-clone%{random_suffix}"
	capacity_gib = "200"
	share_name = "tf-test-test-volume-clone%{random_suffix}"
	storage_pool = google_netapp_storage_pool.default2.name
	protocols = ["NFSV3"]
	deletion_policy = "FORCE"
	restore_parameters {
		source_snapshot = google_netapp_volume_snapshot.test-snapshot.id
	}
}

data "google_compute_network" "default" {
	name = "%{network_name}"
}
	`, context)
}

// Tests adding backup vault and policy to the volume
func testAccNetappVolume_volumeBasicExample_setupBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_netapp_storage_pool" "default2" {
	name = "tf-test-pool%{random_suffix}"
	location = "us-west2"
	service_level = "EXTREME"
	capacity_gib = "2048"
	network = data.google_compute_network.default.id
}
	
resource "google_netapp_volume" "test_volume" {
	location = "us-west2"
	name = "tf-test-test-volume%{random_suffix}"
	capacity_gib = "200"
	share_name = "tf-test-test-volume%{random_suffix}"
	storage_pool = google_netapp_storage_pool.default2.name
	protocols = ["NFSV3"]
	security_style = "UNIX"
	# Delete protection only gets active after an NFS client mounts.
	# Setting it here is save, volume can still be deleted.
	restricted_actions = ["DELETE"]
	deletion_policy = "FORCE"
	backup_config {
		backup_policies = [
		  google_netapp_backup_policy.backup-policy.id
		]
		backup_vault = google_netapp_backup_vault.backup-vault.id
		scheduled_backup_enabled = true
	  }
}

resource "google_netapp_backup_vault" "backup-vault" {
	location = "us-west2"
	name = "tf-test-vault%{random_suffix}"
}
  
resource "google_netapp_backup_policy" "backup-policy" {
	name          		 = "tf-test-backup-policy%{random_suffix}"
	location 			 = "us-west2"
	daily_backup_limit   = 2
	weekly_backup_limit  = 1
	monthly_backup_limit = 1
	enabled = true
}

data "google_compute_network" "default" {
	name = "%{network_name}"
}
	`, context)
}

// Updates backup settings. Cannot update backup vault, since only one allowed per location
func testAccNetappVolume_volumeBasicExample_updateBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_netapp_storage_pool" "default2" {
	name = "tf-test-pool%{random_suffix}"
	location = "us-west2"
	service_level = "EXTREME"
	capacity_gib = "2048"
	network = data.google_compute_network.default.id
}
	
resource "google_netapp_volume" "test_volume" {
	location = "us-west2"
	name = "tf-test-test-volume%{random_suffix}"
	capacity_gib = "200"
	share_name = "tf-test-test-volume%{random_suffix}"
	storage_pool = google_netapp_storage_pool.default2.name
	protocols = ["NFSV3"]
	security_style = "UNIX"
	# Delete protection only gets active after an NFS client mounts.
	# Setting it here is save, volume can still be deleted.
	restricted_actions = ["DELETE"]
	deletion_policy = "FORCE"
	backup_config {
		backup_policies = [
		  google_netapp_backup_policy.backup-policy2.id
		]
		backup_vault = google_netapp_backup_vault.backup-vault.id
		scheduled_backup_enabled = false
	  }
}

resource "google_netapp_backup_vault" "backup-vault" {
	location = "us-west2"
	name = "tf-test-vault%{random_suffix}"
}

resource "google_netapp_backup_policy" "backup-policy2" {
	name          		 = "tf-test-backup-policy2%{random_suffix}"
	location 			 = "us-west2"
	daily_backup_limit   = 4
	weekly_backup_limit  = 2
	monthly_backup_limit = 2
	enabled = true
}

data "google_compute_network" "default" {
	name = "%{network_name}"
}
	`, context)
}
