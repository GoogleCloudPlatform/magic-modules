package backupdr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBackupDRDataSource_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	stepChecks := func(wantName string, wantState string) []resource.TestCheckFunc {
		stepCheck := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr("data.google_backup_dr_data_source.foo", "name", wantName),
			resource.TestCheckResourceAttr("data.google_backup_dr_data_source.foo", "state", wantState),
		}
		return stepCheck
	}

	expectedName := "projects/liyunhuang-consumer/locations/us-central1/backupVaults/bv-test/dataSources/ds-test"
	expectedState := "ACTIVE"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBackupDRDataSource_basic(context),
				Check:  resource.ComposeTestCheckFunc(stepChecks(expectedName, expectedState)...),
			},
		},
	})
}

func testAccDataSourceGoogleBackupDRDataSource_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

data "google_backup_dr_data_source" "foo" {
  project = data.google_project.project.project_id
  location      = "us-central1"
  backup_vault_id = "bv-test"
  data_source_id = "ds-test"
}

`, context)
}
