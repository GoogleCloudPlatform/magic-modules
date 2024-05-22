package securitycenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterOrganizationMuteConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":       envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:testAccCheckSecurityCenterOrganizationMuteConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterOrganizationMuteConfig_basic(context),
			},
			{
				ResourceName:"google_scc_organization_mute_config.default",
				ImportState:true,
				ImportStateVerify:true,
				ImportStateVerifyIgnore: []string{
					"parent",
					"organizationmuteConfigId",
				},
			},
		},
	})
}

func testAccSecurityCenterOrganizationMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_organization_mute_config" "default" {
  organizationmuteConfigId = "tf-test-my-config%{random_suffix}"
  parent             = "organizations/%{org_id}"
  description        = "A test organization mute config"
  filter             = "resource_name = \"projects/my-project/resource-type/findings\""
}
`, context)
}
