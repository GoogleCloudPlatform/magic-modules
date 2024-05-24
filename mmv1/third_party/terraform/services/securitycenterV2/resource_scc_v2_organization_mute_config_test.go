package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/config"
	"google.golang.org/api/securitycenter/v2"
	"google.golang.org/api/googleapi"
)

func TestAccSecurityCenterV2OrganizationMuteConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecurityCenterV2OrganizationMuteConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2OrganizationMuteConfig_basic(context),
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

func testAccSecurityCenterV2OrganizationMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_organization_mute_config" "default" {
  organization_mute_config_id = "tf-test-my-config%{random_suffix}"
  location                    = "global"
  parent                      = "organizations/%{org_id}"
  description                 = "A test organization mute config"
  filter                      = "severity = \"LOW\""
}
`, context)
}

func testAccCheckSecurityCenterV2OrganizationMuteConfigDestroyProducer(t *testing.T) resource.TestCheckFunc {
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

			_, err = sc.Organizations.Locations.MuteConfigs.Get(name).Do()
			if err == nil {
				return fmt.Errorf("Organization mute config %s still exists", name)
			}
			if !isGoogleAPINotFoundError(err) {
				return fmt.Errorf("Error while retrieving organization mute config: %s", err)
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
