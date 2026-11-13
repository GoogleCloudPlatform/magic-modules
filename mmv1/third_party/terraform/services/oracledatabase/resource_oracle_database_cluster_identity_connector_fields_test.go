package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/oracledatabase"
)

func TestAccOracleDatabaseCloudVmCluster_identityConnectorDisabled(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project":                         "oasis-terraform-testing-prod",
		"cloud_exadata_infrastructure_id": fmt.Sprintf("ofake-tf-test-exadata-ic-dis-%s", randomSuffix),
		"cloud_vm_cluster_id":             fmt.Sprintf("ofake-tf-test-vmcluster-ic-dis-%s", randomSuffix),
		"hostname_prefix":                 fmt.Sprintf("h%s", randomSuffix),
		"odb_network":                     "projects/oasis-terraform-testing-prod/locations/europe-west2/odbNetworks/tf-test-permanent-odbnetwork",
		"odb_subnet":                      "projects/oasis-terraform-testing-prod/locations/europe-west2/odbNetworks/tf-test-permanent-odbnetwork/odbSubnets/tf-test-permanent-client-odbsubnet",
		"backup_odb_subnet":               "projects/oasis-terraform-testing-prod/locations/europe-west2/odbNetworks/tf-test-permanent-odbnetwork/odbSubnets/tf-test-permanent-backup-odbsubnet",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOracleDatabaseCloudVmClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseCloudVmCluster_identityConnectorDisabledConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_oracle_database_cloud_vm_cluster.vmcluster", "identity_connector_disabled", "true"),
					resource.TestCheckResourceAttr("google_oracle_database_cloud_vm_cluster.vmcluster", "identity_connector.#", "0"),
				),
			},
			{
				ResourceName:      "google_oracle_database_cloud_vm_cluster.vmcluster",
				ImportState:       true,
				ImportStateVerify: true,
				// cloud_vm_cluster_id, location, deletion_protection, identity_connector_disabled, properties.0.gi_version, and properties.0.hostname_prefix are URL-only, virtual, or INPUT_ONLY API fields not returned on GET.
				ImportStateVerifyIgnore: []string{"cloud_vm_cluster_id", "deletion_protection", "identity_connector_disabled", "location", "properties.0.gi_version", "properties.0.hostname_prefix"},
			},
			{
				ResourceName:       "google_oracle_database_cloud_vm_cluster.vmcluster",
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				ImportStateKind:    resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func TestAccOracleDatabaseExadbVmCluster_identityConnectorDisabled(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project":                      "oasis-terraform-testing-prod",
		"exascale_db_storage_vault_id": fmt.Sprintf("ofake-tf-test-vault-ic-dis-%s", randomSuffix),
		"exadb_vm_cluster_id":          fmt.Sprintf("ofake-tf-test-exadb-ic-dis-%s", randomSuffix),
		"hostname_prefix":              fmt.Sprintf("h%s", randomSuffix),
		"odb_network":                  "projects/oasis-terraform-testing-prod/locations/europe-west2/odbNetworks/tf-test-permanent-odbnetwork",
		"odb_subnet":                   "projects/oasis-terraform-testing-prod/locations/europe-west2/odbNetworks/tf-test-permanent-odbnetwork/odbSubnets/tf-test-permanent-client-odbsubnet",
		"backup_odb_subnet":            "projects/oasis-terraform-testing-prod/locations/europe-west2/odbNetworks/tf-test-permanent-odbnetwork/odbSubnets/tf-test-permanent-backup-odbsubnet",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOracleDatabaseExadbVmClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseExadbVmCluster_identityConnectorDisabledConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_oracle_database_exadb_vm_cluster.exadb", "identity_connector_disabled", "true"),
					resource.TestCheckResourceAttr("google_oracle_database_exadb_vm_cluster.exadb", "identity_connector.#", "0"),
				),
			},
			{
				ResourceName:      "google_oracle_database_exadb_vm_cluster.exadb",
				ImportState:       true,
				ImportStateVerify: true,
				// exadb_vm_cluster_id, location, deletion_protection, identity_connector_disabled, and properties.0.hostname_prefix are URL-only, virtual, or INPUT_ONLY API fields not returned on GET.
				ImportStateVerifyIgnore: []string{"deletion_protection", "exadb_vm_cluster_id", "identity_connector_disabled", "location", "properties.0.hostname_prefix"},
			},
			{
				ResourceName:       "google_oracle_database_exadb_vm_cluster.exadb",
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				ImportStateKind:    resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func TestAccDataSourceOracleDatabaseCloudVmCluster_identityConnectorOutputs(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: `
data "google_oracle_database_cloud_vm_cluster" "default_vmcluster" {
  cloud_vm_cluster_id = "ofake-do-not-delete-tf-vmcluster"
  project             = "oasis-terraform-testing-prod"
  location            = "us-east4"
}

data "google_oracle_database_cloud_vm_clusters" "all_vmclusters" {
  project  = "oasis-terraform-testing-prod"
  location = "us-east4"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_oracle_database_cloud_vm_cluster.default_vmcluster", "identity_connector.#", "1"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_vm_cluster.default_vmcluster", "identity_connector.0.service_agent_email"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_vm_cluster.default_vmcluster", "identity_connector.0.connection_state"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_vm_cluster.default_vmcluster", "identity_connector.0.workload_identity_pool"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_vm_cluster.default_vmcluster", "identity_connector.0.identity_connector_id"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_vm_clusters.all_vmclusters", "cloud_vm_clusters.#"),
				),
			},
		},
	})
}

func testAccOracleDatabaseCloudVmCluster_identityConnectorDisabledConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_oracle_database_cloud_exadata_infrastructure" "infra" {
  cloud_exadata_infrastructure_id = "%{cloud_exadata_infrastructure_id}"
  display_name                    = "%{cloud_exadata_infrastructure_id} displayname"
  location                        = "europe-west2"
  project                         = "%{project}"

  properties {
    shape         = "Exadata.X9M"
    compute_count = "2"
    storage_count = "3"
  }

  deletion_protection = false
}

resource "google_oracle_database_cloud_vm_cluster" "vmcluster" {
  cloud_vm_cluster_id         = "%{cloud_vm_cluster_id}"
  display_name                = "%{cloud_vm_cluster_id} displayname"
  location                    = "europe-west2"
  project                     = "%{project}"
  exadata_infrastructure      = google_oracle_database_cloud_exadata_infrastructure.infra.id
  odb_network                 = "%{odb_network}"
  odb_subnet                  = "%{odb_subnet}"
  backup_odb_subnet           = "%{backup_odb_subnet}"
  identity_connector_disabled = true

  properties {
    license_type    = "LICENSE_INCLUDED"
    ssh_public_keys = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCz1X2744t+6vRLmE5u6nHi6/QWh8bQDgHmd+OIxRQIGA/IWUtCs2FnaCNZcqvZkaeyjk5v0lTA/n+9jvO42Ipib53athrfVG8gRt8fzPL66C6ZqHq+6zZophhrCdfJh/0G4x9xJh5gdMprlaCR1P8yAaVvhBQSKGc4SiIkyMNBcHJ5YTtMQMTfxaB4G1sHZ6SDAY9a6Cq/zNjDwfPapWLsiP4mRhE5SSjJX6l6EYbkm0JeLQg+AbJiNEPvrvDp1wtTxzlPJtIivthmLMThFxK7+DkrYFuLvN5AHUdo9KTDLvHtDCvV70r8v0gafsrKkM/OE9Jtzoo0e1N/5K/ZdyFRbAkFT4QSF3nwpbmBWLf2Evg//YyEuxnz4CwPqFST2mucnrCCGCVWp1vnHZ0y30nM35njLOmWdRDFy5l27pKUTwLp02y3UYiiZyP7d3/u5pKiN4vC27VuvzprSdJxWoAvluOiDeRh+/oeQDowxoT/Oop8DzB9uJmjktXw8jyMW2+Rpg+ENQqeNgF1OGlEzypaWiRskEFlkpLb4v/s3ZDYkL1oW0Nv/J8LTjTOTEaYt2Udjoe9x2xWiGnQixhdChWuG+MaoWffzUgx1tsVj/DBXijR5DjkPkrA1GA98zd3q8GKEaAdcDenJjHhNYSd4+rE9pIsnYn7fo5X/tFfcQH1XQ== nobody@google.com"]
    cpu_core_count  = "4"
    gi_version      = "19.0.0.0"
    hostname_prefix = "%{hostname_prefix}"
  }

  deletion_protection = false
}
`, context)
}

func testAccOracleDatabaseExadbVmCluster_identityConnectorDisabledConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_oracle_database_exascale_db_storage_vault" "vault" {
  exascale_db_storage_vault_id = "%{exascale_db_storage_vault_id}"
  display_name                 = "%{exascale_db_storage_vault_id} displayname"
  location                     = "europe-west2"
  project                      = "%{project}"

  properties {
    exascale_db_storage_details {
      total_size_gbs = 512
    }
  }

  deletion_protection = false
}

resource "google_oracle_database_exadb_vm_cluster" "exadb" {
  exadb_vm_cluster_id         = "%{exadb_vm_cluster_id}"
  display_name                = "%{exadb_vm_cluster_id} displayname"
  location                    = "europe-west2"
  project                     = "%{project}"
  odb_network                 = "%{odb_network}"
  odb_subnet                  = "%{odb_subnet}"
  backup_odb_subnet           = "%{backup_odb_subnet}"
  identity_connector_disabled = true

  properties {
    ssh_public_keys = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCz1X2744t+6vRLmE5u6nHi6/QWh8bQDgHmd+OIxRQIGA/IWUtCs2FnaCNZcqvZkaeyjk5v0lTA/n+9jvO42Ipib53athrfVG8gRt8fzPL66C6ZqHq+6zZophhrCdfJh/0G4x9xJh5gdMprlaCR1P8yAaVvhBQSKGc4SiIkyMNBcHJ5YTtMQMTfxaB4G1sHZ6SDAY9a6Cq/zNjDwfPapWLsiP4mRhE5SSjJX6l6EYbkm0JeLQg+AbJiNEPvrvDp1wtTxzlPJtIivthmLMThFxK7+DkrYFuLvN5AHUdo9KTDLvHtDCvV70r8v0gafsrKkM/OE9Jtzoo0e1N/5K/ZdyFRbAkFT4QSF3nwpbmBWLf2Evg//YyEuxnz4CwPqFST2mucnrCCGCVWp1vnHZ0y30nM35njLOmWdRDFy5l27pKUTwLp02y3UYiiZyP7d3/u5pKiN4vC27VuvzprSdJxWoAvluOiDeRh+/oeQDowxoT/Oop8DzB9uJmjktXw8jyMW2+Rpg+ENQqeNgF1OGlEzypaWiRskEFlkpLb4v/s3ZDYkL1oW0Nv/J8LTjTOTEaYt2Udjoe9x2xWiGnQixhdChWuG+MaoWffzUgx1tsVj/DBXijR5DjkPkrA1GA98zd3q8GKEaAdcDenJjHhNYSd4+rE9pIsnYn7fo5X/tFfcQH1XQ== nobody@google.com"]
    time_zone {
      id = "UTC"
    }
    grid_image_id               = "ocid1.dbpatch.oc1.uk-london-1.anwgiljrt5t4sqqa7anvfhtjk3kukfffjqwjyu2fv435wlcw3hzto6iqyngq"
    node_count                  = 1
    enabled_ecpu_count_per_node = 8
    vm_file_system_storage {
      size_in_gbs_per_node = 220
    }
    exascale_db_storage_vault = google_oracle_database_exascale_db_storage_vault.vault.id
    hostname_prefix           = "%{hostname_prefix}"
    shape_attribute           = "SMART_STORAGE"
  }

  deletion_protection = false
}
`, context)
}
