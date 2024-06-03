package securitycenterv2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/config"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/securitycenter/v1"
)

func TestAccSecurityCenterv2OrganizationMuteConfig_basic(t *testing.T) {
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
		CheckDestroy:             testAccCheckSecurityCenterv2OrganizationMuteConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterv2OrganizationMuteConfig_basic(contextBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.default", "description", "A test organization mute config"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.default", "filter", "severity = \"LOW\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.default", "organization_mute_config_id", fmt.Sprintf("tf-test-my-config%s", contextBasic["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.default", "location", contextBasic["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.default", "parent", fmt.Sprintf("organizations/%s", contextBasic["org_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_organization_mute_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "organization_mute_config_id"},
			},
			{
				Config: testAccSecurityCenterv2OrganizationMuteConfig_highSeverity(contextHighSeverity),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.default", "description", "A test organization mute config with high severity"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.default", "filter", "severity = \"HIGH\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.default", "organization_mute_config_id", fmt.Sprintf("tf-test-my-config%s", contextHighSeverity["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.default", "location", contextHighSeverity["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_organization_mute_config.default", "parent", fmt.Sprintf("organizations/%s", contextHighSeverity["org_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_organization_mute_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "organization_mute_config_id"},
			},
		},
	})
}

func testAccSecurityCenterv2OrganizationMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_organization_mute_config" "default" {
  description          = "A test organization mute config"
  filter               = "severity = \"LOW\""
  organization_mute_config_id = "tf-test-my-config%{random_suffix}"
  location             = "%{location}"
  parent               = "organizations/%{org_id}"
}
`, context)
}

func testAccSecurityCenterv2OrganizationMuteConfig_highSeverity(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_organization_mute_config" "default" {
  description          = "A test organization mute config with high severity"
  filter               = "severity = \"HIGH\""
  organization_mute_config_id = "tf-test-my-config%{random_suffix}"
  location             = "%{location}"
  parent               = "organizations/%{org_id}"
}
`, context)
}

func testAccCheckSecurityCenterv2OrganizationMuteConfigDestroyProducer(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.Provider.Meta().(*config.Config)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_scc_v2_organization_mute_config" {
				continue
			}

			// Initialize Security Command Center Service
			sc, err := securitycenter.NewService(context.Background(), config.GoogleClientOptions...)
			if err != nil {
				return fmt.Errorf("Error creating Security Command Center client: %s", err)
			}

			// Get the organization mute config by name
			name := rs.Primary.ID

			_, err = sc.Organizations.MuteConfigs.Get(name).Do()
			if err == nil {
				return fmt.Errorf("Organization mute config %s still exists", name)
			}
			if !isGoogleAPINotFoundError(err) {
				return fmt.Errorf("Error while retrieving mute config: %s", err)
			}
		}

		return nil
	}
}

// isGoogleAPINotFoundError checks if an error is a Google API Not Found error
func isGoogleAPINotFoundError(err error) bool {
	apiErr, ok := err.(*googleapi.Error)
	return ok && apiErr.Code == 404
}
