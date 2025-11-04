package backupdr_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBackupDRListDataSourceReferences_basic(t *testing.T) {
	t.Parallel()

	projectDsName := "data.google_project.project"
	var projectID string
	context := map[string]interface{}{
		"location":      "us-central1",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBackupDRListDataSourceReferences_basic(context),
				Check: func(s *terraform.State) error {
					project, ok := s.RootModule().Resources[projectDsName]
					if !ok {
						return fmt.Errorf("project data source not found: %s", projectDsName)
					}
					projectID = project.Primary.Attributes["project_id"]

					return resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.google_backup_dr_list_data_source_references.all", "project", projectID),
						resource.TestCheckResourceAttr("data.google_backup_dr_list_data_source_references.all", "location", context["location"].(string)),
						resource.TestCheckResourceAttrSet("data.google_backup_dr_list_data_source_references.all", "data_source_references.#"),
					)(s)
				},
			},
		},
	})
}

func testAccDataSourceGoogleBackupDRListDataSourceReferences_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_service_account" "default" {
 account_id   = "tf-test-my-custom-%{random_suffix}"
 display_name = "Custom SA for VM Instance"
}

resource "google_sql_database_instance" "instance" {
  name                = "tf-test-ds-list-%{random_suffix}"
  database_version    = "MYSQL_8_0"
  region              = "%{location}"
  deletion_protection = false
  settings {
    tier              = "db-n1-standard-1"
    availability_type = "ZONAL"
	activation_policy = "ALWAYS"
  }
}

resource "google_backup_dr_backup_vault" "vault" {
  location          = "%{location}"
  backup_vault_id   = "tf-test-bv-list-%{random_suffix}"
  description       = "Acceptance test vault"
  backup_minimum_enforced_retention_duration = "100000s"
  labels = {
	foo = "bar1"
	bar = "baz1"
  }
  annotations = {
	annotations1 = "bar1"
	annotations2 = "baz1"
  }
  force_update      = "true"
  force_delete      = "true"
  allow_missing     = "true"
}

resource "google_backup_dr_backup_plan" "plan" {
  location        = "%{location}"
  backup_plan_id  = "tf-test-bp-list-%{random_suffix}"
  resource_type   = "sqladmin.googleapis.com/Instance"
  backup_vault    = google_backup_dr_backup_vault.vault.name
   backup_rules {
   rule_id                = "rule-1"
   backup_retention_days  = 2


   standard_schedule {
     recurrence_type     = "HOURLY"
     hourly_frequency    = 6
     time_zone           = "UTC"


     backup_window {
       start_hour_of_day = 12
       end_hour_of_day   = 18
     }
   }
 }
}

resource "google_backup_dr_backup_plan_association" "bpa" {
  location                   = "%{location}"
  backup_plan_association_id = "tf-test-bpa-list-%{random_suffix}"
  resource                   = google_sql_database_instance.instance.self_link
  resource_type              = "sqladmin.googleapis.com/Instance"
  backup_plan                = google_backup_dr_backup_plan.plan.name
  depends_on                 = [google_sql_database_instance.instance]
}

data "google_backup_dr_list_data_source_references" "all" {
  project    = data.google_project.project.project_id
  location   = "%{location}"
  filter	 = "resource_type=sqladmin.googleapis.com/Instance"
  order_by   = "name"
  depends_on = [google_backup_dr_backup_plan_association.bpa]
}
`, context)
}
