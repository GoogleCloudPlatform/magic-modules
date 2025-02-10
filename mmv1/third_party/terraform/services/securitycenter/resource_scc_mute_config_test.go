package securitycenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterMuteConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterMuteConfig_basic(context),
			},
			{
				ResourceName:      "google_scc_mute_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"mute_config_id",
				},
			},
			{
				Config: testAccSecurityCenterMuteConfig_update(context),
			},
			{
				ResourceName:      "google_scc_mute_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"mute_config_id",
				},
			},
		},
	})
}

func testAccSecurityCenterMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_scc_mute_config" "default" {
  mute_config_id = "tf-test-mute-config-%{random_suffix}"
  parent       	 = "organizations/%{org_id}"
  filter         = "category: \"OS_VULNERABILITY\""
  description    = "A Test Mute Config"
  type           = "DYNAMIC"
  expiry_time    = "2215-02-03T15:01:23Z"
}  
`, context)
}

func testAccSecurityCenterMuteConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_scc_mute_config" "default" {
  mute_config_id = "tf-test-mute-config-%{random_suffix}"
  parent       	 = "organizations/%{org_id}"
  filter         = "category: \"OS_VULNERABILITY\""
  description    = "An Updated Test Mute Config"
  type           = "DYNAMIC"
  expiry_time    = "2215-02-03T15:01:23Z"
}  
`, context)
}
