package migrationcenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMigrationCenterSettings_settingsUpdate(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMigrationCenterSettings_settingsStart(),
			},
			{
				ResourceName:            "google_migration_center_settings.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccMigrationCenterSettings_settingsUpdateFalse(),
			},
			{
				ResourceName:            "google_migration_center_settings.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccMigrationCenterSettings_settingsUpdateTrue(),
			},
			{
				ResourceName:            "google_migration_center_settings.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccMigrationCenterSettings_settingsStart() string {
	return `
resource "google_migration_center_settings" "default" {
  location              = "us-central1"
  disable_cloud_logging = true
}
`
}

func testAccMigrationCenterSettings_settingsUpdateFalse() string {
	return `
resource "google_migration_center_settings" "default" {
  location              = "us-central1"
  disable_cloud_logging = false
}
`
}

func testAccMigrationCenterSettings_settingsUpdateTrue() string {
	return `
resource "google_migration_center_settings" "default" {
  location              = "us-central1"
  disable_cloud_logging = true
}
`
}
