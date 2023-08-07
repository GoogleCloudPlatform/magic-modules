package google

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

/*
 * Restore tests are kept separate from other cluster tests because they require an instance and a backup to exist
 */

// Restore tests depend on instances and backups being taken, which can take up to 10 minutes. Since the instance doesn't change in between tests,
// we condense everything into individual test cases.
// 1. Create the source cluster, instance, and backup
// 2. Restore from the backup directly
// 3. Determine the
func TestAccAlloydbCluster_restore(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedTestNetwork(t, "alloydbinstance-mandatory"),
	}

	time.Sleep(10000 * time.Millisecond)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbClusterAndInstanceAndBackup(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.source",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackup(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.restored_from_backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location", "restore_source_backup"},
			},
			{
				Config: testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.restored_from_point_in_time",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location", "restore_source_cluster", "restore_point_in_time"},
			},
			{
				Config: testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime_AllowDestroy(context),
			},
		},
	})
}

func testAccAlloydbClusterAndInstanceAndBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = data.google_compute_network.default.id
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = data.google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
`, context)
}

func testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

resource "google_alloydb_cluster" "restored_from_backup" {
  cluster_id            = "tf-test-alloydb-backup-restored-cluster-%{random_suffix}"
  location              = "us-central1"
  network               = data.google_compute_network.default.id
  restore_source_backup = google_alloydb_backup.default.name

  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = data.google_compute_network.default.id
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = data.google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
`, context)
}

// The source cluster, instance, and backup should all exist prior to this being invoked. Otherwise the PITR restore will not succeed
// due to the time being too early.
func testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

resource "google_alloydb_cluster" "restored_from_backup" {
  cluster_id            = "tf-test-alloydb-backup-restored-cluster-%{random_suffix}"
  location              = "us-central1"
  network               = data.google_compute_network.default.id
  restore_source_backup = google_alloydb_backup.default.name

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_alloydb_cluster" "restored_from_point_in_time" {
  cluster_id             = "tf-test-alloydb-pitr-restored-cluster-%{random_suffix}"
  location               = "us-central1"
  network                = data.google_compute_network.default.id
  restore_source_cluster = google_alloydb_cluster.source.name
  restore_point_in_time  = google_alloydb_backup.default.update_time

  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = data.google_compute_network.default.id
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = data.google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
`, context)
}

// The source cluster, instance, and backup should all exist prior to this being invoked. Otherwise the PITR restore will not succeed
// due to the time being too early.
func testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime_AllowDestroy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

resource "google_alloydb_cluster" "restored_from_backup" {
  cluster_id            = "tf-test-alloydb-backup-restored-cluster-%{random_suffix}"
  location              = "us-central1"
  network               = data.google_compute_network.default.id
  restore_source_backup = google_alloydb_backup.default.name
}

resource "google_alloydb_cluster" "restored_from_point_in_time" {
  cluster_id             = "tf-test-alloydb-pitr-restored-cluster-%{random_suffix}"
  location               = "us-central1"
  network                = data.google_compute_network.default.id
  restore_source_cluster = google_alloydb_cluster.source.name
  restore_point_in_time  = google_alloydb_backup.default.update_time
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = data.google_compute_network.default.id
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = data.google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
`, context)
}

// // Validates that updating fields on a restored cluster works fine and doesn't require re-creating the cluster
// func TestAccAlloydbCluster_restore_updateRestoredCluster(t *testing.T) {
// }

// // Validates that updating the restore source or point in time requires re-creating the cluster.
// // This encompasses both updates to the fields and removing the fields entirely.
// func TestAccAlloydbCluster_restore_cannotUpdateRestoreSource(t *testing.T) {
// }

// // Validates that only one restore source can be provided
// func TestAccAlloydbCluster_restore_onlyOneSourceAllowed(t *testing.T) {
// }

// // Validates that pointInTime and sourceCluster must come together
// func TestAccAlloydbCluster_restore_sourceClusterAndPointInTimeRequired(t *testing.T) {
// }
