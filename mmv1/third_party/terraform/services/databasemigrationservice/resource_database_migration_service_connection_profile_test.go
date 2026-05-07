package databasemigrationservice_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseMigrationServiceConnectionProfile_update(t *testing.T) {
	t.Parallel()

	suffix := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServiceConnectionProfile_basic(suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_database_migration_service_connection_profile.default", "role", "SOURCE"),
				),
			},
			{
				ResourceName:            "google_database_migration_service_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "mysql.0.password", "labels", "terraform_labels"},
			},
			{
				Config: testAccDatabaseMigrationServiceConnectionProfile_update(suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_database_migration_service_connection_profile.default", "role", "DESTINATION"),
				),
			},
			{
				ResourceName:            "google_database_migration_service_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "mysql.0.password", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccDatabaseMigrationServiceConnectionProfile_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_database_migration_service_connection_profile" "default" {
	location = "us-central1"
	connection_profile_id = "tf-test-dbms-connection-profile%{random_suffix}"
	display_name          = "tf-test-dbms-connection-profile-display%{random_suffix}"
	role                  = "SOURCE"
	labels	= { 
		foo = "bar" 
	}
	mysql {
	  host = "10.20.30.40"
	  port = 3306
	  username = "tf-test-dbms-test-user%{random_suffix}"
	  password = "tf-test-dbms-test-pass%{random_suffix}"
	}
}
`, context)
}

func testAccDatabaseMigrationServiceConnectionProfile_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_database_migration_service_connection_profile" "default" {
	location = "us-central1"
	connection_profile_id = "tf-test-dbms-connection-profile%{random_suffix}"
	display_name          = "tf-test-dbms-connection-profile-updated-display%{random_suffix}"
	role                  = "DESTINATION"
	labels	= { 
		bar = "foo" 
	}
	mysql {
	  host = "10.20.30.50"
	  port = 3306
	  username = "tf-test-update-dbms-test-user%{random_suffix}"
	  password = "tf-test-update-dbms-test-pass%{random_suffix}"
	}
}
`, context)
}

func TestAccDatabaseMigrationServiceConnectionProfile_databaseMigrationServiceConnectionProfileAlloydb(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "vpc-network-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDatabaseMigrationServiceConnectionProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServiceConnectionProfile_databaseMigrationServiceConnectionProfileAlloydb(context),
			},
			{
				ResourceName:            "google_database_migration_service_connection_profile.alloydbprofile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "alloydb.0.settings.0.initial_user.0.password", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccDatabaseMigrationServiceConnectionProfile_databaseMigrationServiceConnectionProfileAlloydb(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_database_migration_service_connection_profile" "alloydbprofile" {
  location = "us-central1"
  connection_profile_id = "tf-test-my-profileid%{random_suffix}"
  display_name = "tf-test-my-profileid%{random_suffix}_display"
  labels = { 
    foo = "bar" 
  }
  alloydb {
    cluster_id = "tf-test-dbmsalloycluster%{random_suffix}"
    settings {
      initial_user {
        user = "alloyuser%{random_suffix}"
        password = "alloypass%{random_suffix}"
      }
      vpc_network = data.google_compute_network.default.id
      labels  = { 
        alloyfoo = "alloybar" 
      }
      primary_instance_settings {
        id = "priminstid"
        machine_config {
          cpu_count = 2
        }
        database_flags = { 
        }
        labels = { 
          alloysinstfoo = "allowinstbar" 
        }
      }
    }
  }
}
`, context)
}

func TestAccDatabaseMigrationServiceConnectionProfile_postgresqlPrivateSsl(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServiceConnectionProfile_postgresqlPrivateSsl(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_database_migration_service_connection_profile.default", "role", "SOURCE"),
					resource.TestCheckResourceAttr("google_database_migration_service_connection_profile.default", "postgresql.0.ssl.0.type", "REQUIRED"),
					resource.TestCheckResourceAttrPair("google_database_migration_service_connection_profile.default", "postgresql.0.private_connectivity.0.private_connection", "google_database_migration_service_private_connection.default", "id"),
				),
			},
			{
				ResourceName:            "google_database_migration_service_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "postgresql.0.password", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccDatabaseMigrationServiceConnectionProfile_postgresqlPrivateSsl(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_database_migration_service_private_connection" "default" {
	display_name          = "dbms_pc"
	location              = "us-central1"
	private_connection_id = "tf-test-my-connection-ssl%{random_suffix}"

	psc_interface_config {
		network_attachment = google_compute_network_attachment.default.id
	}
}

resource "google_compute_network_attachment" "default" {
  name                  = "tf-test-attachment-ssl%{random_suffix}"
  region                = "us-central1"
  connection_preference = "ACCEPT_AUTOMATIC"
  subnetworks           = [google_compute_subnetwork.default.id]
}

resource "google_compute_network" "default" {
  name = "tf-test-network-ssl%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-subnet-ssl%{random_suffix}"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.default.id
}

resource "google_database_migration_service_connection_profile" "default" {
	location = "us-central1"
	connection_profile_id = "tf-test-dbms-cp-ssl%{random_suffix}"
	display_name          = "tf-test-dbms-cp-ssl-display%{random_suffix}"
	role                  = "SOURCE"

	postgresql {
	  host = "10.20.30.40"
	  port = 5432
	  username = "tf-test-dbms-user%{random_suffix}"
	  password = "tf-test-dbms-pass%{random_suffix}"
	  database = "tf-test-db"
	  ssl {
	  	type = "REQUIRED"
	  }
	  private_connectivity {
	  	private_connection = google_database_migration_service_private_connection.default.id
	  }
	}
}
`, context)
}

func TestAccDatabaseMigrationServiceConnectionProfile_postgresqlDestinationCloudSql(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServiceConnectionProfile_postgresqlDestinationCloudSql(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_database_migration_service_connection_profile.default", "role", "DESTINATION"),
					resource.TestCheckResourceAttrPair("google_database_migration_service_connection_profile.default", "postgresql.0.cloud_sql_id", "google_sql_database_instance.postgres", "name"),
				),
			},
			{
				ResourceName:            "google_database_migration_service_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "postgresql.0.password", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccDatabaseMigrationServiceConnectionProfile_postgresqlDestinationCloudSql(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_database_migration_service_connection_profile" "default" {
	location = "us-central1"
	connection_profile_id = "tf-test-dbms-cp-dest%{random_suffix}"
	display_name          = "tf-test-dbms-cp-dest-display%{random_suffix}"
	role                  = "DESTINATION"

	postgresql {
	  username = "tf-test-dbms-user%{random_suffix}"
	  password = "tf-test-dbms-pass%{random_suffix}"
	  cloud_sql_id = google_sql_database_instance.postgres.name
	}
}

resource "google_sql_database_instance" "postgres" {
  name             = "tf-test-clouddb-%{random_suffix}"
  database_version = "POSTGRES_12"
  settings {
    tier = "db-f1-micro"
  }
  deletion_protection =  false
}
`, context)
}
