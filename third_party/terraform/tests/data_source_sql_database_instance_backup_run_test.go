package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSqlDatabaseInstanceBackupRun_basic(t *testing.T) {
	t.Parallel()

	instance := BootstrapSharedSQLInstanceBackupRun(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlDatabaseInstanceBackupRun_basic(instance),
				Check:  resource.TestMatchResourceAttr("data.google_sql_database_instance_backup_run.backup", "status", regexp.MustCompile("SUCCESSFUL")),
			},
		},
	})
}

func TestAccDataSourceSqlDatabaseInstanceBackupRun_notFound(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceSqlDatabaseInstanceBackupRun_notFound(context),
				ExpectError: regexp.MustCompile("No backups found for SQL Database Instance"),
			},
		},
	})
}

func testAccDataSourceSqlDatabaseInstanceBackupRun_basic(instance string) string {
	return fmt.Sprintf(`
data "google_sql_database_instance_backup_run" "backup" {
	instance = "%s"
	most_recent = true
}
`, instance)
}

func testAccDataSourceSqlDatabaseInstanceBackupRun_notFound(context map[string]interface{}) string {
	return Nprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "tf-test-instance-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  settings {
	tier = "db-f1-micro"
	backup_configuration {
		enabled            = "false"
	}
  }

  deletion_protection = false
}

data "google_sql_database_instance_backup_run" "backup" {
	instance = google_sql_database_instance.instance.name
	most_recent = true
	depends_on = [google_sql_database_instance.instance]
}
`, context)
}
