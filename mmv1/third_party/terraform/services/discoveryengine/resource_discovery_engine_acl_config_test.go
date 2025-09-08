package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDiscoveryEngineAclConfig_discoveryengineAclconfigBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineAclConfig_discoveryengineAclconfigBasicExample_basic(context),
			},
			{
				ResourceName:            "google_discovery_engine_acl_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccDiscoveryEngineAclConfig_discoveryengineAclconfigBasicExample_update(context),
			},
			{
				ResourceName:            "google_discovery_engine_acl_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccDiscoveryEngineAclConfig_discoveryengineAclconfigBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_acl_config" "basic" {
  location = "global"
  idp_config {
    idp_type = "THIRD_PARTY"
    external_idp_config {
      workforce_pool_name = "locations/global/workforcePools/cloud-console-pool-manual"
    }
  }
}
`, context)
}

func testAccDiscoveryEngineAclConfig_discoveryengineAclconfigBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_acl_config" "basic" {
  location = "global"
  idp_config {
    idp_type = "GSUITE"
  }
}
`, context)
}
