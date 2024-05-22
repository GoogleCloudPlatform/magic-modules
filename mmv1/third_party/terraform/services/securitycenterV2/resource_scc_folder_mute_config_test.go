package securitycenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterFolderMuteConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":       envvar.GetTestOrgFromEnv(t),
		"folder_id":     acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:testAccCheckSecurityCenterFolderMuteConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterFolderMuteConfig_basic(context),
			},
			{
				ResourceName:"google_scc_folder_mute_config.default",
				ImportState:true,
				ImportStateVerify:true,
				ImportStateVerifyIgnore: []string{
					"parent",
					"foldermuteConfigId",
				},
			},
		},
	})
}

func testAccSecurityCenterFolderMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_folder_mute_config" "default" {
  foldermuteConfigId = "tf-test-my-config%{random_suffix}"
  location           = "us-central1"
  parent             = "folders/%{folder_id}"
  description        = "A test folder mute config"
  filter             = "severity = \"LOW\""
}
`, context)
}
