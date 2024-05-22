package securitycenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterProjectMuteConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":  envvar.GetTestProjectFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:testAccCheckSecurityCenterProjectMuteConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterProjectMuteConfig_basic(context),
			},
			{
				ResourceName:"google_scc_project_mute_config.default",
				ImportState:true,
				ImportStateVerify:true,
				ImportStateVerifyIgnore: []string{
					"parent",
					"projectmuteConfigId",
				},
			},
		},
	})
}

func testAccSecurityCenterProjectMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_project_mute_config" "default" {
  projectmuteConfigId = "tf-test-my-config%{random_suffix}"
  location           = "us-central1"
  parent             = "projects/%{project_id}"
  description        = "A test project mute config"
  filter             = "severity = \"MEDIUM\""
}
`, context)
}
