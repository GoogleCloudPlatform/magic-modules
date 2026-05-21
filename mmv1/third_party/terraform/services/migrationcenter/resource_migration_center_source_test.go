package migrationcenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMigrationCenterSource_migrationSourceUpdate(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"source_id":     "tf-test-source-test" + randomSuffix,
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMigrationCenterSourceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMigrationCenterSource_migrationSourceStart(context),
			},
			{
				ResourceName:            "google_migration_center_source.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "source_id"},
			},
			{
				Config: testAccMigrationCenterSource_migrationSourceUpdate(context),
			},
			{
				ResourceName:            "google_migration_center_source.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "source_id"},
			},
		},
	})
}

func testAccMigrationCenterSource_migrationSourceStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_migration_center_source" "default" {
  location     = "us-central1"
  source_id    = "%{source_id}"
  description  = "Terraform integration test description"
  display_name = "Terraform integration test display"
  priority     = 10
  type         = "SOURCE_TYPE_CUSTOM"
  managed      = false
}
`, context)
}

func testAccMigrationCenterSource_migrationSourceUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_migration_center_source" "default" {
  location     = "us-central1"
  source_id    = "%{source_id}"
  description  = "Updated Terraform integration test description"
  display_name = "Updated integration test display"
  priority     = 15
  type         = "SOURCE_TYPE_CUSTOM"
  managed      = false
}
`, context)
}
