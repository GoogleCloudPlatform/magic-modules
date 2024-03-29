package datalossprevention_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/envvar"
)

func TestAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigBasicExample(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
    inspect_templates = ["FAKE"]
}
`, context)
}
