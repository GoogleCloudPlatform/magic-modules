package securitycenterv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterV2OrganizationMuteConfig_basic(t *testing.T) {
	t.Parallel()

	contextBasic := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"location":      "global",
	}

	contextHighSeverity := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"location":      "us_central",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2OrganizationMuteConfig_basic(contextBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test1", "description", "A test organization mute config"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test1", "filter", "severity = \"LOW\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test1", "mute_config_id", fmt.Sprintf("tf-test-my-config-%s", contextBasic["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test1", "location", contextBasic["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test1", "parent", fmt.Sprintf("organizations/%s", contextBasic["org_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_organization_mute_config.organization_mute_test1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "mute_config_id"},
			},
			{
				Config: testAccSecurityCenterV2OrganizationMuteConfig_highSeverity(contextHighSeverity),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test2", "description", "A test organization mute config with high severity"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test2", "filter", "severity = \"HIGH\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test2", "mute_config_id", fmt.Sprintf("tf-test-my-config-%s", contextHighSeverity["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test2", "location", contextHighSeverity["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test2", "parent", fmt.Sprintf("organizations/%s", contextHighSeverity["org_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_organization_mute_config.organization_mute_test2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "mute_config_id"},
			},
		},
	})
}

func testAccSecurityCenterV2OrganizationMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_organization_mute_config" "organization_mute_test1" {
  description          = "A test organization mute config"
  filter               = "severity = \"LOW\""
  mute_config_id       = "tf-test-my-config-%{random_suffix}"
  location             = "global"
  parent               = "organizations/%{org_id}"
  type                 =  "STATIC"
}
`, context)
}

func testAccSecurityCenterV2OrganizationMuteConfig_highSeverity(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_organization_mute_config" "organization_mute_test2" {
  description          = "A test organization mute config with high severity"
  filter               = "severity = \"HIGH\""
  mute_config_id       = "tf-test-my-config-%{random_suffix}"
  location             = "global"
  parent               = "organizations/%{org_id}"
  type                 =  "STATIC"
}
`, context)
}
