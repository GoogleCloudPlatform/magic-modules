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

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"location":      "global",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2OrganizationMuteConfig(context, "A test organization mute config", "severity = \"LOW\"", "organization_mute_test1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test1", "description", "A test organization mute config"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test1", "filter", "severity = \"LOW\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test1", "mute_config_id", fmt.Sprintf("tf-test-my-config-%s", context["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test1", "location", context["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test1", "organization", fmt.Sprintf("organizations/%s", context["org_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_organization_mute_config.organization_mute_test1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization", "location"},
			},
			{
				Config: testAccSecurityCenterV2OrganizationMuteConfig(context, "A test organization mute config with high severity", "severity = \"HIGH\"", "organization_mute_test2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test2", "description", "A test organization mute config with high severity"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test2", "filter", "severity = \"HIGH\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test2", "mute_config_id", fmt.Sprintf("tf-test-my-config-%s", context["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test2", "location", context["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.organization_mute_test2", "organization", fmt.Sprintf("organizations/%s", context["org_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_organization_mute_config.organization_mute_test2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization", "location"},
			},
		},
	})
}

func testAccSecurityCenterV2OrganizationMuteConfig(context map[string]interface{}, description, filter, resourceName string) string {
	return acctest.Nprintf(`
resource "google_scc_v2_organization_mute_config" "%s" {
  description    = "%s"
  filter         = "%s"
  mute_config_id = "tf-test-my-config-%{random_suffix}"
  location       = "%{location}"
  organization   = "organizations/%{org_id}"
  type           = "STATIC"
}
`, resourceName, description, filter, context)
}
