package backupdr_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBackupDRFetchBackupsForResourceType_basic(t *testing.T) {
	t.Parallel()

	projectDsName := "data.google_project.project"
	var projectID string
	context := map[string]interface{}{
		"location":      "us-central1",
		"resource_type": "sqladmin.googleapis.com/Instance",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBackupDRFetchBackupsForResourceType_basic(context),
				Check: func(s *terraform.State) error {
					project, ok := s.RootModule().Resources[projectDsName]
					if !ok {
						return fmt.Errorf("project data source not found: %s", projectDsName)
					}
					projectID = project.Primary.Attributes["project_id"]

					return resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.google_backup_dr_fetch_backups.default", "project", projectID),
						resource.TestCheckResourceAttr("data.google_backup_dr_fetch_backups.default", "location", context["location"].(string)),
						resource.TestCheckResourceAttr("data.google_backup_dr_fetch_backups.default", "resource_type", context["resource_type"].(string)),
						resource.TestCheckResourceAttrSet("data.google_backup_dr_fetch_backups.default", "backups.#"),
					)(s)
				},
			},
		},
	})
}

func testAccDataSourceGoogleBackupDRFetchBackupsForResourceType_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}
resource "google_service_account" "default" {
 account_id   = "tf-test-my-custom-%{random_suffix}"
 display_name = "Custom SA for VM Instance"
}
// Prerequisite resources to create a DataSource and potential Backups
resource "google_sql_database_instance" "instance" {
 name             = "default-%{random_suffix}"
 database_version = "MYSQL_8_0"
 region          = "us-central1"
 deletion_protection = false
 settings {
   tier = "db-f1-micro"
   availability_type = "ZONAL"
   activation_policy = "ALWAYS"
 }
}
resource "google_backup_dr_backup_vault" "my-backup-vault" {
   location ="%{location}"
   backup_vault_id    = "tf-test-bv-%{random_suffix}"
   description = "This is a second backup vault built by Terraform."
   backup_minimum_enforced_retention_duration = "100000s"
   labels = {
     foo = "bar1"
     bar = "baz1"
   }
   annotations = {
     annotations1 = "bar1"
     annotations2 = "baz1"
   }
   force_update = "true"
   force_delete = "true"
   allow_missing = "true"
}
resource "google_backup_dr_backup_plan" "foo" {
 location       = "%{location}"
 backup_plan_id = "tf-test-bp-test-%{random_suffix}"
 resource_type  = "%{resource_type}"
 backup_vault   = google_backup_dr_backup_vault.my-backup-vault.name
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
 location = "%{location}"
 backup_plan_association_id = "tf-test-bpa-test-%{random_suffix}"
 resource = "projects/${data.google_project.project.project_id}/instances/${google_sql_database_instance.instance.name}"
 resource_type= "%{resource_type}"
 backup_plan = google_backup_dr_backup_plan.foo.name
 depends_on = [ google_sql_database_instance.instance ]
}

resource "time_sleep" "wait_for_data_source" {
  depends_on = [google_backup_dr_backup_plan_association.bpa]
  create_duration = "120s"
}

// The actual data source under test
data "google_backup_dr_fetch_backups" "default" {
  project         = data.google_project.project.project_id
  location        = "%{location}"
  backup_vault_id = google_backup_dr_backup_vault.my-backup-vault.backup_vault_id
  data_source_id  = element(split("/", google_backup_dr_backup_plan_association.bpa.data_source), 7)
  resource_type   = "%{resource_type}"
  depends_on      = [time_sleep.wait_for_data_source]
}
`, context)
}
