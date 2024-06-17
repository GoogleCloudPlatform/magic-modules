package securitycenterv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterOrganizationMuteConfig_createUpdateDelete(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"location":      "global",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecurityCenterOrganizationMuteConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterOrganizationMuteConfig(context),
			},
			{
				ResourceName:            "google_scc_organization_mute_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organizationmuteConfigId", "parent"},
			},
			{
				Config: testAccSecurityCenterOrganizationMuteConfig_update(context),
			},
			{
				ResourceName:            "google_scc_organization_mute_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organizationmuteConfigId", "parent"},
			},
			{
				Config: testAccSecurityCenterOrganizationMuteConfig_delete(context),
			},
		},
	})
}

func testAccSecurityCenterOrganizationMuteConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_organization_mute_config" "default" {
  mute_config_id = "tf-test-my-config-%{random_suffix}"
  parent         = "organizations/%{org_id}"
  filter         = "category: \"OS_VULNERABILITY\""
  location       = "%{location}"
  description    = "My Mute Config"
  type           = "STATIC"
}
`, context)
}

func testAccSecurityCenterOrganizationMuteConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_organization_mute_config" "default" {
  mute_config_id = "tf-test-my-config-%{random_suffix}"
  parent         = "organizations/%{org_id}"
  filter         = "category: \"WEB_SECURITY\""
  location       = "%{location}"
  description    = "My Mute Config Updated"
  type           = "STATIC"
}
`, context)
}

func testAccSecurityCenterOrganizationMuteConfig_delete(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_organization_mute_config" "default" {
  mute_config_id = "tf-test-my-config-%{random_suffix}"
  parent         = "organizations/%{org_id}"
  filter         = "category: \"OS_VULNERABILITY\""
  location       = "%{location}" 
  description    = "My Mute Config"
  type           = "STATIC"
}
`, context)
}
