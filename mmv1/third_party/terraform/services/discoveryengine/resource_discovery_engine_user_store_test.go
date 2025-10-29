package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDiscoveryEngineUserStore_discoveryengineUserstoreBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineUserStore_discoveryengineUserstoreBasicExample_basic(context),
			},
			{
				ResourceName:      "google_discovery_engine_user_store.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDiscoveryEngineUserStore_discoveryengineUserstoreBasicExample_update(context),
			},
			{
				ResourceName:      "google_discovery_engine_user_store.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDiscoveryEngineUserStore_discoveryengineUserstoreBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_license_config" "basic" {
  location = "global"
  license_config_id = "tf-test-license-config-id%{random_suffix}"
  license_count = 50
  subscription_tier = "SUBSCRIPTION_TIER_SEARCH_AND_ASSISTANT"
  start_date {
    year = 2099
    month = 1
    day = 1
  }
  end_date {
    year = 2100
    month = 1
    day = 1
  }
  subscription_term = "SUBSCRIPTION_TERM_ONE_YEAR"
}

data "google_project" "project" {}

resource "google_discovery_engine_user_store" "basic" {
  location = google_discovery_engine_license_config.basic.location
  default_license_config = "projects/${data.google_project.project.number}/locations/${google_discovery_engine_license_config.basic.location}/licenseConfigs/${google_discovery_engine_license_config.basic.license_config_id}"
  enable_license_auto_register = true
}
`, context)
}

func testAccDiscoveryEngineUserStore_discoveryengineUserstoreBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_license_config" "basic" {
  location = "global"
  license_config_id = "tf-test-license-config-id%{random_suffix}"
  license_count = 50
  subscription_tier = "SUBSCRIPTION_TIER_SEARCH_AND_ASSISTANT"
  start_date {
    year = 2099
    month = 1
    day = 1
  }
  end_date {
    year = 2100
    month = 1
    day = 1
  }
  subscription_term = "SUBSCRIPTION_TERM_ONE_YEAR"
}

resource "google_discovery_engine_user_store" "basic" {
  location = google_discovery_engine_license_config.basic.location
}
`, context)
}
