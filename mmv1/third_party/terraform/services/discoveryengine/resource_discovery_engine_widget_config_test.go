package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDiscoveryEngineWidgetConfig_discoveryengineWidgetconfigBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineWidgetConfig_discoveryengineWidgetconfigBasicExample_basic(context),
			},
			{
				ResourceName:      "google_discovery_engine_widget_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDiscoveryEngineWidgetConfig_discoveryengineWidgetconfigBasicExample_update(context),
			},
			{
				ResourceName:      "google_discovery_engine_widget_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDiscoveryEngineWidgetConfig_discoveryengineWidgetconfigBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "basic" {
  location                    = "global"
  data_store_id               = "tf-test-data-store-id%{random_suffix}"
  display_name                = "tf-test-datastore"
  industry_vertical           = "GENERIC"
  content_config              = "NO_CONTENT"
  solution_types              = ["SOLUTION_TYPE_SEARCH"]
  create_advanced_site_search = false
}

resource "google_discovery_engine_search_engine" "basic" {
  engine_id                   = "tf-test-engine-id%{random_suffix}"
  collection_id               = "default_collection"
  location                    = google_discovery_engine_data_store.basic.location
  display_name                = "tf-test-engine"
  data_store_ids              = [google_discovery_engine_data_store.basic.data_store_id]
  industry_vertical           = "GENERIC"
  app_type                    = "APP_TYPE_INTRANET"
  search_engine_config {
  }
}

resource "google_discovery_engine_widget_config" "basic" {
  location = google_discovery_engine_search_engine.basic.location
  engine_id = google_discovery_engine_search_engine.basic.engine_id
  access_settings {
    enable_web_app = true
  }
}
`, context)
}

func testAccDiscoveryEngineWidgetConfig_discoveryengineWidgetconfigBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "basic" {
  location                    = "global"
  data_store_id               = "tf-test-data-store-id%{random_suffix}"
  display_name                = "tf-test-datastore"
  industry_vertical           = "GENERIC"
  content_config              = "NO_CONTENT"
  solution_types              = ["SOLUTION_TYPE_SEARCH"]
  create_advanced_site_search = false
}

resource "google_discovery_engine_search_engine" "basic" {
  engine_id                   = "tf-test-engine-id%{random_suffix}"
  collection_id               = "default_collection"
  location                    = google_discovery_engine_data_store.basic.location
  display_name                = "tf-test-engine"
  data_store_ids              = [google_discovery_engine_data_store.basic.data_store_id]
  industry_vertical           = "GENERIC"
  app_type                    = "APP_TYPE_INTRANET"
  search_engine_config {
  }
}

resource "google_discovery_engine_widget_config" "basic" {
  location = google_discovery_engine_search_engine.basic.location
  engine_id = google_discovery_engine_search_engine.basic.engine_id
  access_settings {
    enable_web_app = true
    workforce_identity_pool_provider = "locations/global/workforcePools/workforce-pool-id/providers/workforce-pool-provider"
  }
}
`, context)
}
