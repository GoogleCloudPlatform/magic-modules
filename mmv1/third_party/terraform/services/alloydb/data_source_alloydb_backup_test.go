package alloydb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/alloydb"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/services/servicenetworking"
)

func TestAccAlloydbBackupDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  servicenetworking.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbBackupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbBackupDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_alloydb_backup.default",
						"google_alloydb_backup.default",
					),
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_alloydb_backup.without_optional",
						"google_alloydb_backup.default",
					),
				),
			},
		},
	})
}

func testAccAlloydbBackupDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }

  deletion_protection = false
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.default.name

  description = "example description"
  type        = "ON_DEMAND"
  labels = {
    "label" = "key"
  }
  annotations = {
    "annotation" = "value"
  }
  depends_on = [google_alloydb_instance.default]
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

data "google_alloydb_backup" "default" {
  backup_id = google_alloydb_backup.default.backup_id
  location  = google_alloydb_backup.default.location
  project   = google_alloydb_backup.default.project
}

data "google_alloydb_backup" "without_optional" {
  backup_id = google_alloydb_backup.default.backup_id
}
`, context)
}
