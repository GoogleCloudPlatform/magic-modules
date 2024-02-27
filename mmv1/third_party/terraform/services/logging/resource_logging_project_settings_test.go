package logging_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingProjectSettings_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id": envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSettings_onlyRequired(context),
			},
			{
				ResourceName:            "google_logging_project_settings.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingProjectSettings_full(context),
			},
			{
				ResourceName:            "google_logging_project_settings.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func testAccLoggingProjectSettings_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_logging_project_settings" "example" {
  project                = "%{project_name}"
  kms_service_account_id = data.google_logging_project_settings.settings.logging_service_account_id
}

data "google_logging_project_settings" "settings" {
  project = "%{project_name}"
}
`, context)
}

func testAccLoggingProjectSettings_onlyRequired(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_logging_project_settings" "example" {
  project = "%{project_name}"
}
`, context)
}
