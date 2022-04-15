package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSettingsOrganizationSetting_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id": getTestOrgTargetFromEnv(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSettingsOrganizationSetting_start(context),
			},
			{
				ResourceName:      "google_resource_settings_organization_setting.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceSettingsOrganizationSetting_update(context),
			},
			{
				ResourceName:      "google_resource_settings_organization_setting.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceSettingsOrganizationSetting_end(context),
			},
			{
				ResourceName:      "google_resource_settings_organization_setting.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceSettingsOrganizationSetting_start(context map[string]interface{}) string {
	return Nprintf(`
resource "google_resource_settings_organization_setting" "default" {
  organization_id = "%{org_id}"
  name = "iam-serviceAccountKeyExpiry"
  local_value {
    string_value = "24hours"
  }
}
`, context)
}

func testAccResourceSettingsOrganizationSetting_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_resource_settings_organization_setting" "default" {
  organization_id = "%{org_id}"
  name = "iam-serviceAccountKeyExpiry"
  local_value {
    string_value = "7days"
  }
}
`, context)
}

func testAccResourceSettingsOrganizationSetting_end(context map[string]interface{}) string {
	return Nprintf(`
resource "google_resource_settings_organization_setting" "default" {
  organization_id = "%{org_id}"
  name = "iam-serviceAccountKeyExpiry"
  local_value {
    string_value = "never"
  }
}
`, context)
}
