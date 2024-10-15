package backupdr_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBackupDRDataSource_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBackupDRDataSource_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_backup_dr_data_source.foo", "google_backup_dr_backup_vault.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBackupDRDataSource_basic(context map[string]interface{}) string {
	return fmt.Sprintf(`
data "google_project" "project" {
}

resource "google_backup_dr_backup_vault" "foo" {
  backup_vault_id = "bv-test"
  backup_minimum_enforced_retention_duration = "100000s"
  location = "us-central1"
  provider = google
}

data "google_backup_dr_data_source" "foo" {
  name = "tf-test-data-source-%{random_suffix}"
  project = data.google_project.project.project_id
  location      = "us-central1"
  backup_vault_id = "bv-test"
  data_source_id = "ds-test"
}

`, context)
}
