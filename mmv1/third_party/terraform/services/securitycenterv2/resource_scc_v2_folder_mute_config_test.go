package securitycenterv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterV2FolderMuteConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2FolderMuteConfig_basic(context),
			},
			{
				ResourceName:      "google_scc_v2_folder_mute_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"folder", "location", "mute_config_id",
				},
			},
			{
				Config: testAccSecurityCenterV2FolderMuteConfig_update(context),
			},
			{
				ResourceName:      "google_scc_v2_folder_mute_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"folder", "location", "mute_config_id",
				},
			},
		},
	})
}

func testAccSecurityCenterV2FolderMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_folder_mute_config" "default" {
  mute_config_id = "tf-test-config-%{random_suffix}"
  folder   = "%{folder}"
  location       = "global"
  description    = "A test folder mute config"
  filter         = "severity = \"HIGH\""
  type           = "STATIC"
}
`, context)
}

func testAccSecurityCenterV2FolderMuteConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_folder_mute_config" "default" {
  mute_config_id = "tf-test-config-%{random_suffix}"
  folder   = "%{folder}"
  location       = "global"
  description    = "An updated test folder mute config"
  filter         = "severity = \"HIGH\""
  type           = "STATIC"
}
`, context)
}
