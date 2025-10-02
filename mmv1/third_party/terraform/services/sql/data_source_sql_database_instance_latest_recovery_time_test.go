package sql_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceSqlDatabaseInstanceLatestRecoveryTime_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}
	resourceName := "data.google_sql_database_instance_latest_recovery_time.default"

	expectedError := regexp.MustCompile(`.*No backups found for instance.* and deletion time seconds: .*`)

	if acctest.IsVcrEnabled() {
		expectedError = regexp.MustCompile(`.*Error 400.*`)
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlDatabaseInstanceLatestRecoveryTime_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance"),
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "latest_recovery_time"),
				),
			},
			{
				// On non-deleted instance should return error containing both instance name and deletion time
				Config:      testAccDataSourceSqlDatabaseInstanceLatestRecoveryTime_withDeletionTime(context),
				ExpectError: expectedError,
			},
		},
	})
}

func testAccDataSourceSqlDatabaseInstanceLatestRecoveryTime_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "main" {
  name             = "tf-test-instance-%{random_suffix}"
  database_version = "POSTGRES_14"
  region           = "us-central1"

  settings {
    tier = "db-g1-small"
    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
      start_time                     = "20:55"
      transaction_log_retention_days = "3"
    }
  }

  deletion_protection = false
}

resource "time_sleep" "wait_for_instance" {
  // Wait 30 seconds after the instance is created
  depends_on = [google_sql_database_instance.main]
  create_duration = "330s"
}

data "google_sql_database_instance_latest_recovery_time" "default" {
  instance = google_sql_database_instance.main.name
  depends_on = [time_sleep.wait_for_instance]
}
`, context)
}

func testAccDataSourceSqlDatabaseInstanceLatestRecoveryTime_withDeletionTime(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "main" {
  name             = "tf-test-instance-%{random_suffix}"
  database_version = "POSTGRES_14"
  region           = "us-central1"

  settings {
    tier = "db-g1-small"
    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
      start_time                     = "20:55"
      transaction_log_retention_days = "3"
    }
  }

  deletion_protection = false
}

resource "time_sleep" "wait_for_instance" {
  // Wait 30 seconds after the instance is created
  depends_on = [google_sql_database_instance.main]
  create_duration = "330s"
}

data "google_sql_database_instance_latest_recovery_time" "default" {
  instance = google_sql_database_instance.main.name
  source_instance_deletion_time = "2025-06-20T17:23:59.648821586Z"
  depends_on = [time_sleep.wait_for_instance]
}
`, context)
}
