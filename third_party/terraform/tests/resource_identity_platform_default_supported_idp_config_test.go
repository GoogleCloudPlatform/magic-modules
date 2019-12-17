package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityPlatformDefaultSupportedIdpConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigBasic(context),
			},
			{
				ResourceName:      "google_identity_platform_default_supported_idp_config.idp_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigUpdate(context),
			},
			{
				ResourceName:      "google_identity_platform_default_supported_idp_config.idp_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_identity_platform_default_supported_idp_config" "idp_config" {
  enabled = true
  client_id = "playgames.google.com"
  client_secret = "secret"
}
`, context)
}

func testAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_identity_platform_default_supported_idp_config" "idp_config" {
  enabled = false
  client_id = "playgames.google.com"
  client_secret = "anothersecret"
}
`, context)
}
