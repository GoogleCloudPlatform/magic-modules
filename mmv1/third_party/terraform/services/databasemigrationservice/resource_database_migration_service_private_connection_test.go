package databasemigrationservice_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseMigrationServicePrivateConnection_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServicePrivateConnection_basic(context),
			},
			{
				ResourceName:            "google_database_migration_service_private_connection.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_connection_id", "location", "labels", "terraform_labels", "create_without_validation"},
			},
		},
	})
}

func testAccDatabaseMigrationServicePrivateConnection_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_database_migration_service_private_connection" "default" {
	display_name          = "dbms_pc"
	location              = "us-central1"
	private_connection_id = "tf-test-my-connection%{random_suffix}"

	labels = {
		key = "value"
	}

	vpc_peering_config {
		vpc_name = google_compute_network.default.id
		subnet = "10.0.0.0/29"
	}
}

resource "google_compute_network" "default" {
  name = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}

func TestAccDatabaseMigrationServicePrivateConnection_psc(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServicePrivateConnection_psc(context),
			},
			{
				ResourceName:            "google_database_migration_service_private_connection.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_connection_id", "location", "labels", "terraform_labels", "create_without_validation"},
			},
		},
	})
}

func testAccDatabaseMigrationServicePrivateConnection_psc(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_database_migration_service_private_connection" "default" {
	display_name          = "dbms_pc"
	location              = "us-central1"
	private_connection_id = "tf-test-my-connection%{random_suffix}"

	labels = {
		key = "value"
	}

	psc_interface_config {
		network_attachment = google_compute_network_attachment.default.id
	}
}

resource "google_compute_network_attachment" "default" {
  name                  = "tf-test-attachment%{random_suffix}"
  region                = "us-central1"
  connection_preference = "ACCEPT_AUTOMATIC"
  subnetworks           = [google_compute_subnetwork.default.id]
}

resource "google_compute_network" "default" {
  name = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-subnet%{random_suffix}"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.default.id
}
`, context)
}
