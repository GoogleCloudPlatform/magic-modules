package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/discoveryengine"
)

func TestAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDiscoveryEngineSearchEngineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_basic(context),
			},
			{
				ResourceName:            "google_discovery_engine_search_engine.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"engine_id", "collection_id", "location", "kms_key_name"},
			},
			{
				Config: testAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_update(context),
			},
			{
				ResourceName:            "google_discovery_engine_search_engine.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"engine_id", "collection_id", "location", "kms_key_name"},
			},
		},
	})
}

func testAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "basic" {
    location                    = "global"
    data_store_id               = "tf-test-example-datastore%{random_suffix}"
    display_name                = "tf-test-structured-datastore"
    industry_vertical           = "GENERIC"
    content_config              = "NO_CONTENT"
    solution_types              = ["SOLUTION_TYPE_SEARCH"]
    create_advanced_site_search = false
    }
resource "google_discovery_engine_data_store" "second" {
    location                    = "global"
    data_store_id               = "tf-test-example2-datastore%{random_suffix}"
    display_name                = "tf-test-structured-datastore2"
    industry_vertical           = "GENERIC"
    content_config              = "NO_CONTENT"
    solution_types              = ["SOLUTION_TYPE_SEARCH"]
    create_advanced_site_search = false
    }
resource "google_discovery_engine_search_engine" "basic" {
  engine_id = "tf-test-example-engine-id%{random_suffix}"
  collection_id = "default_collection"
  location = google_discovery_engine_data_store.basic.location
  display_name = "Example Display Name"
  data_store_ids = [google_discovery_engine_data_store.basic.data_store_id, google_discovery_engine_data_store.second.data_store_id]
  industry_vertical = google_discovery_engine_data_store.basic.industry_vertical
  common_config {
    company_name = "Example Company Name"
  }
  app_type = "APP_TYPE_INTRANET"
  search_engine_config {
    search_tier = "SEARCH_TIER_ENTERPRISE"
    required_subscription_tier = "SUBSCRIPTION_TIER_ENTERPRISE"
    search_add_ons = ["SEARCH_ADD_ON_LLM"]
  }
  features = {
    "agent-sharing-without-admin-approval" = "FEATURE_STATE_ON"
    "disable-agent-sharing" = "FEATURE_STATE_OFF"
  }
  knowledge_graph_config {
    enable_cloud_knowledge_graph = false
    enable_private_knowledge_graph = true
  }
}
`, context)
}

func testAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "basic" {
    location                    = "global"
    data_store_id               = "tf-test-example-datastore%{random_suffix}"
    display_name                = "tf-test-structured-datastore"
    industry_vertical           = "GENERIC"
    content_config              = "NO_CONTENT"
    solution_types              = ["SOLUTION_TYPE_SEARCH"]
    create_advanced_site_search = false
    }
resource "google_discovery_engine_data_store" "second" {
    location                    = "global"
    data_store_id               = "tf-test-example2-datastore%{random_suffix}"
    display_name                = "tf-test-structured-datastore2"
    industry_vertical           = "GENERIC"
    content_config              = "NO_CONTENT"
    solution_types              = ["SOLUTION_TYPE_SEARCH"]
    create_advanced_site_search = false
    }
resource "google_discovery_engine_search_engine" "basic" {
  engine_id = "tf-test-example-engine-id%{random_suffix}"
  collection_id = "default_collection"
  location = google_discovery_engine_data_store.basic.location
  display_name = "Updated Example Display Name"
  data_store_ids = [google_discovery_engine_data_store.basic.data_store_id]
  industry_vertical = google_discovery_engine_data_store.basic.industry_vertical
  disable_analytics = true
  common_config {
    company_name = "Updated Example Company Name"
  }
  app_type = "APP_TYPE_INTRANET"
  search_engine_config {
    search_tier = "SEARCH_TIER_STANDARD"
    required_subscription_tier = "SUBSCRIPTION_TIER_ENTERPRISE"
    search_add_ons = ["SEARCH_ADD_ON_LLM"]
  }
  features = {
    feedback = "FEATURE_STATE_OFF"
    "agent-sharing-without-admin-approval" = "FEATURE_STATE_ON"
    "disable-agent-sharing" = "FEATURE_STATE_OFF"
  }
  knowledge_graph_config {
    enable_cloud_knowledge_graph = false
    cloud_knowledge_graph_types = ["foobar"]
    enable_private_knowledge_graph = true
    feature_config {
      disable_private_kg_query_understanding = true
      disable_private_kg_enrichment = true
      disable_private_kg_auto_complete = true
      disable_private_kg_query_ui_chips = true
    }
  }
}
`, context)
}

func TestDiscoveryEngineSearchEngineFeaturesDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Key, Old, New      string
		ExpectDiffSuppress bool
	}{
		"server provided feature in state when config is empty": {
			Key:                "features.enable-end-user-sharing-with-groups",
			Old:                "FEATURE_STATE_OFF",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"server provided feature in config": {
			Key:                "features.enable-end-user-sharing-with-groups",
			Old:                "FEATURE_STATE_OFF",
			New:                "FEATURE_STATE_ON",
			ExpectDiffSuppress: false,
		},
		"other feature in state when config is empty": {
			Key:                "features.agent-sharing-without-admin-approval",
			Old:                "FEATURE_STATE_OFF",
			New:                "",
			ExpectDiffSuppress: false,
		},
		"features map item count delta": {
			Key:                "features.%",
			Old:                "3",
			New:                "2",
			ExpectDiffSuppress: true,
		},
		"other map item count delta": {
			Key:                "other_map.%",
			Old:                "3",
			New:                "2",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if discoveryengine.DiscoveryEngineSearchEngineFeaturesDiffSuppress(tc.Key, tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, key=%q %q => %q expect DiffSuppress to return %t", tn, tc.Key, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}
