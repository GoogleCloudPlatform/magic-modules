package google_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	google "internal/terraform-provider-google"
)

func TestAccDataSourceSqlDatabase_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": google.RandString(t, 10),
	}

	google.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { google.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: google.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlDatabase_basic(context),
				Check: resource.ComposeTestCheckFunc(
					google.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_sql_database.qa",
						"google_sql_database.db",
						map[string]struct{}{
							"deletion_policy": {},
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceSqlDatabase_basic(context map[string]interface{}) string {
	return google.Nprintf(`
resource "google_sql_database_instance" "main" {
  name             = "tf-test-instance-%{random_suffix}"
  database_version = "POSTGRES_14"
  region           = "us-central1"

  settings {
    tier = "db-f1-micro"
  }

  deletion_protection = false
}

resource "google_sql_database" "db" {
	name = "tf-test-db-%{random_suffix}"
	instance = google_sql_database_instance.main.name
	depends_on = [
		google_sql_database_instance.main
	]
}

data "google_sql_database" "qa" {
	name = google_sql_database.db.name
    instance = google_sql_database_instance.main.name
	depends_on = [
		google_sql_database.db
  	]
}
`, context)
}
