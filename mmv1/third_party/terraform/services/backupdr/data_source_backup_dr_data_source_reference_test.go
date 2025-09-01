package backupdr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBackupDRDataSourceReference_basic(t *testing.T) {
	t.Parallel()

	// IMPORTANT: Replace these placeholders with values from an EXISTING
	// DataSourceReference in your 'kushallunkad-consmer' project.
	// This test WILL FAIL if this DSR does not exist.
	preExistingDSRID := "3de95d556c1f37169eb3529f300d2b7fc82cbbdd" // REPLACE THIS
	preExistingLocation := "us-central1"                           // REPLACE THIS if different
	// expectedBackupVault := "projects/658319736396/locations/us-central1/backupVaults/default-vault-us-central1"                                                     // REPLACE THIS - Expected value from the existing DSR
	expectedDataSource := "projects/658319736396/locations/us-central1/backupVaults/default-vault-us-central1/dataSources/2132a3489855e0947e2c23fa2982f2bc1c5d6060" // REPLACE THIS - Expected value

	context := map[string]interface{}{
		"dsr_id":   preExistingDSRID,
		"location": preExistingLocation,
		"project":  "kushallunkad-consumer", // Project to test against
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t), // Use Beta Factories
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudBackupDRDataSourceReference_Config(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_backup_dr_data_source_reference.dsr_get", "data_source_reference_id", preExistingDSRID),
					resource.TestCheckResourceAttr("data.google_backup_dr_data_source_reference.dsr_get", "location", preExistingLocation),
					resource.TestCheckResourceAttr("data.google_backup_dr_data_source_reference.dsr_get", "project", context["project"].(string)),
					// Validate the output properties based on the pre-existing resource
					resource.TestCheckResourceAttrSet("data.google_backup_dr_data_source_reference.dsr_get", "name"),
					// resource.TestCheckResourceAttr("data.google_backup_dr_data_source_reference.dsr_get", "backup_vault", expectedBackupVault),
					resource.TestCheckResourceAttr("data.google_backup_dr_data_source_reference.dsr_get", "data_source", expectedDataSource),
					resource.TestCheckResourceAttrSet("data.google_backup_dr_data_source_reference.dsr_get", "create_time"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudBackupDRDataSourceReference_Config(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_backup_dr_data_source_reference" "dsr_get" {
  provider                 = google-beta
  data_source_reference_id = "%{dsr_id}"
  location                 = "%{location}"
  project                  = "%{project}"
}
`, context)
}
