package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccOracledatabaseExadbVmCluster_update(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracledatabaseExadbVmCluster_full(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_oracle_database_exadb_vm_cluster.my_exadb_vm_cluster", "properties.0.node_count", "1"),
				),
			},
			{
				ResourceName:      "google_oracle_database_exadb_vm_cluster.my_exadb_vm_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOracledatabaseExadbVmCluster_update(t),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_oracle_database_exadb_vm_cluster.my_exadb_vm_cluster", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_oracle_database_exadb_vm_cluster.my_exadb_vm_cluster", "properties.0.node_count", "2"),
				),
			},
			{
				ResourceName:      "google_oracle_database_exadb_vm_cluster.my_exadb_vm_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOracledatabaseExadbVmCluster_full(t *testing.T) string {
	return fmt.Sprintf(acctest.Nprintf(t, `
resource "google_oracle_database_exadb_vm_cluster" "my_exadb_vm_cluster"{
    exadb_vm_cluster_id = "%s"
    display_name = "%s displayname"
    location = "europe-west2"
    project = "%s"
    odb_network = "%s"
    odb_subnet = "%s"
    backup_odb_subnet = "%s"
    labels = {
        "label-one" = "value-one"
    }
    properties {
        ssh_public_keys = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCz1X2744t+6vRLmE5u6nHi6/QWh8bQDgHmd+OIxRQIGA/IWUtCs2FnaCNZcqvZkaeyjk5v0lTA/n+9jvO42Ipib53athrfVG8gRt8fzPL66C6ZqHq+6zZophhrCdfJh/0G4x9xJh5gdMprlaCR1P8yAaVvhBQSKGc4SiIkyMNBcHJ5YTtMQMTfxaB4G1sHZ6SDAY9a6Cq/zNjDwfPapWLsiP4mRhE5SSjJX6l6EYbkm0JeLQg+AbJiNEPvrvDp1wtTxzlPJtIivthmLMThFxK7+DkrYFuLvN5AHUdo9KTDLvHtDCvV70r8v0gafsrKkM/OE9Jtzoo0e1N/5K/ZdyFRbAkFT4QSF3nwpbmBWLf2Evg//YyEuxnz4CwPqFST2mucnrCCGCVWp1vnHZ0y30nM35njLOmWdRDFy5l27pKUTwLp02y3UYiiZyP7d3/u5pKiN4vC27VuvzprSdJxWoAvluOiDeRh+/oeQDowxoT/Oop8DzB9uJmjktXw8jyMW2+Rpg+ENQqeNgF1OGlEzypaWiRskEFlkpLb4v/s3ZDYkL1oW0Nv/J8LTjTOTEaYt2Udjoe9x2xWiGnQixhdChWuG+MaoWffzUgx1tsVj/DBXijR5DjkPkrA1GA98zd3q8GKEaAdcDenJjHhNYSd4+rE9pIsnYn7fo5X/tFfcQH1XQ== nobody@google.com"]
        time_zone {
            id = "UTC"
        }
        grid_image_id = "ocid1.dbpatch.oc1.uk-london-1.anwgiljrt5t4sqqa7anvfhtjk3kukfffjqwjyu2fv435wlcw3hzto6iqyngq"
        node_count = 1
        enabled_ecpu_count_per_node = 8
        vm_file_system_storage {
            size_in_gbs_per_node = 220
        }
        exascale_db_storage_vault = google_oracle_database_exascale_db_storage_vault.exascaleDbStorageVaults.id
        hostname_prefix = "hostname8"
        shape_attribute = "SMART_STORAGE"
        data_collection_options {
            is_diagnostics_events_enabled = "true"
            is_health_monitoring_enabled  = "true"
            is_incident_logs_enabled      = "true"
        }
        license_model = "LICENSE_INCLUDED"
        scan_listener_port_tcp = 1521
        additional_ecpu_count_per_node = 8
        cluster_name = "example"
    }

    deletion_protection = false
}

resource "google_oracle_database_exascale_db_storage_vault" "exascaleDbStorageVaults"{
  exascale_db_storage_vault_id = "%s"
  display_name = "%s displayname"
  location = "europe-west2"
  project = "%s"
  properties {
    exascale_db_storage_details {
        total_size_gbs = 512
    }
  }

  deletion_protection = false
}
`, acctest.RandString(t, 10), acctest.RandString(t, 10), getTestProjectFromEnv(), acctest.RandString(t, 10), acctest.RandString(t, 10), acctest.RandString(t, 10), acctest.RandString(t, 10), acctest.RandString(t, 10), getTestProjectFromEnv()))
}

func testAccOracledatabaseExadbVmCluster_update(t *testing.T) string {
	return fmt.Sprintf(Nprintf(t, `
resource "google_oracle_database_exadb_vm_cluster" "my_exadb_vm_cluster"{
    exadb_vm_cluster_id = "%s"
    display_name = "%s displayname"
    location = "europe-west2"
    project = "%s"
    odb_network = "%s"
    odb_subnet = "%s"
    backup_odb_subnet = "%s"
    labels = {
        "label-one" = "value-one"
    }
    properties {
        ssh_public_keys = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCz1X2744t+6vRLmE5u6nHi6/QWh8bQDgHmd+OIxRQIGA/IWUtCs2FnaCNZcqvZkaeyjk5v0lTA/n+9jvO42Ipib53athrfVG8gRt8fzPL66C6ZqHq+6zZophhrCdfJh/0G4x9xJh5gdMprlaCR1P8yAaVvhBQSKGc4SiIkyMNBcHJ5YTtMQMTfxaB4G1sHZ6SDAY9a6Cq/zNjDwfPapWLsiP4mRhE5SSjJX6l6EYbkm0JeLQg+AbJiNEPvrvDp1wtTxzlPJtIivthmLMThFxK7+DkrYFuLvN5AHUdo9KTDLvHtDCvV70r8v0gafsrKkM/OE9Jtzoo0e1N/5K/ZdyFRbAkFT4QSF3nwpbmBWLf2Evg//YyEuxnz4CwPqFST2mucnrCCGCVWp1vnHZ0y30nM35njLOmWdRDFy5l27pKUTwLp02y3UYiiZyP7d3/u5pKiN4vC27VuvzprSdJxWoAvluOiDeRh+/oeQDowxoT/Oop8DzB9uJmjktXw8jyMW2+Rpg+ENQqeNgF1OGlEzypaWiRskEFlkpLb4v/s3ZDYkL1oW0Nv/J8LTjTOTEaYt2Udjoe9x2xWiGnQixhdChWuG+MaoWffzUgx1tsVj/DBXijR5DjkPkrA1GA98zd3q8GKEaAdcDenJjHhNYSd4+rE9pIsnYn7fo5X/tFfcQH1XQ== nobody@google.com"]
        time_zone {
            id = "UTC"
        }
        grid_image_id = "ocid1.dbpatch.oc1.uk-london-1.anwgiljrt5t4sqqa7anvfhtjk3kukfffjqwjyu2fv435wlcw3hzto6iqyngq"
        node_count = 2
        enabled_ecpu_count_per_node = 8
        vm_file_system_storage {
            size_in_gbs_per_node = 220
        }
        exascale_db_storage_vault = google_oracle_database_exascale_db_storage_vault.exascaleDbStorageVaults.id
        hostname_prefix = "hostname6"
        shape_attribute = "SMART_STORAGE"
        data_collection_options {
            is_diagnostics_events_enabled = "true"
            is_health_monitoring_enabled  = "true"
            is_incident_logs_enabled      = "true"
        }
        license_model = "LICENSE_INCLUDED"
        scan_listener_port_tcp = 1521
        additional_ecpu_count_per_node = 8
        cluster_name = "example"
    }

    deletion_protection = false
}

resource "google_oracle_database_exascale_db_storage_vault" "exascaleDbStorageVaults"{
  exascale_db_storage_vault_id = "%s"
  display_name = "%s displayname"
  location = "europe-west2"
  project = "%s"
  properties {
    exascale_db_storage_details {
        total_size_gbs = 512
    }
  }

  deletion_protection = false
}
`, acctest.RandString(t, 10), acctest.RandString(t, 10), getTestProjectFromEnv(), acctest.RandString(t, 10), acctest.RandString(t, 10), acctest.RandString(t, 10), acctest.RandString(t, 10), acctest.RandString(t, 10), getTestProjectFromEnv()))
}
