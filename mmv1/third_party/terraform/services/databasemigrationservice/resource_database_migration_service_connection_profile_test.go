package databasemigrationservice_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
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
			},
			{
				ResourceName:            "google_database_migration_service_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "mysql.0.password"},
			},
			{
				Config: testAccDatabaseMigrationServiceConnectionProfile_update(suffix),
			},
			{
				ResourceName:            "google_database_migration_service_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "mysql.0.password"},
			},
		},
	})
}

func TestAccDatabaseMigrationServiceConnectionProfile_Postgres_PSC(t *testing.T) {
	t.Parallel()

	instanceName := "tf-test-" + acctest.RandString(t, 10)
	projectId := "psctestproject" + acctest.RandString(t, 10)
	orgId := envvar.GetTestOrgFromEnv(t)
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	certName := "sqlcert" + acctest.RandString(t, 10)
	userName := "username" + acctest.RandString(t, 10)
	passWord := "password" + acctest.RandString(t, 10)
	profileName := "dbmsprofile" + acctest.RandString(t, 10)
	profileDisplay:= "profiledisplay" + acctest.RandString(t, 10)



	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServiceConnectionProfile_Postgres_PSC(instanceName, projectId, orgId, billingAccount, suffix),
				Check:  resource.ComposeTestCheckFunc(verifyPscOperation("google_sql_database_instance.instance", true, true, []string{envvar.GetTestProjectFromEnv()})),
			},
			{
				ResourceName:            "google_database_migration_service_connection_profile.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/", projectId),
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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

func testAccDatabaseMigrationServiceConnectionProfile_Postgres_PSC(instanceName string, projectId string, orgId string, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "testproject" {
  name                = "%s"
  project_id          = "%s"
  org_id              = "%s"
  billing_account     = "%s"
}

resource "google_sql_database_instance" "postgresqldb" {
  project             = google_project.testproject.project_id
  name                = "%s"
  region              = "us-south1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
		psc_config {
			psc_enabled = true
			allowed_consumer_projects = ["%s"]
		}
		ipv4_enabled = false
    }
	backup_configuration {
		enabled = true
		binary_log_enabled = true
	}
	availability_type = "REGIONAL"
  }
}

resource "google_sql_ssl_cert" "sql_client_cert" {
  common_name = "%s"
  instance    = google_sql_database_instance.postgresqldb.name

  depends_on = [google_sql_database_instance.postgresqldb]
}

resource "google_sql_user" "sqldb_user" {
  name     = "%s"
  instance = google_sql_database_instance.postgresqldb.name
  password = %s"


  depends_on = [google_sql_ssl_cert.sql_client_cert]
}

resource "google_database_migration_service_connection_profile" "dbms_profile" {
  location = "us-central1"
  connection_profile_id = "%s"
  display_name          = "%s"
  labels = { 
    foo = "bar" 
  }
  postgresql {
    host = google_sql_database_instance.postgresqldb.ip_address.0.ip_address
    port = 5432
	username = "%s"
	password = "%s"
    ssl {
      client_key = google_sql_ssl_cert.sql_client_cert.private_key
      client_certificate = google_sql_ssl_cert.sql_client_cert.cert
      ca_certificate = google_sql_ssl_cert.sql_client_cert.server_ca_cert
    }
    cloud_sql_id = "%s"
    private_service_connect_connectivity {
      service_attachment = google_sql_database_instance.postgresqldb.psc_service_attachment_link
    }
  }

`, projectId, projectId, orgId, billingAccount, instanceName, projectId, certName, userName, passWord, profileName, profileDisplay, userName, passWord, instanceName)
}
