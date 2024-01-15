package discoveryengine_test

import (
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"

	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	// "github.com/hashicorp/terraform-provider-google/google/acctest"
	// "github.com/hashicorp/terraform-provider-google/google/tpgresource"
	// transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
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
				ImportStateVerifyIgnore: []string{"engine_id", "collection_id", "location"},
			},
			{
				Config: testAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_update(context),
			},
			{
				ResourceName:            "google_discovery_engine_search_engine.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"engine_id", "collection_id", "location"},
			},
		},
	})
}

func testAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_search_engine" "basic" {
  engine_id = "tf-test-example-engine-id%{random_suffix}"
  collection_id = "default_collection"
  location = "global"
  display_name = "Example Display Name"
  industry_vertical = "GENERIC"
  data_store_ids = ["test-demo-data_1705085526185"]
  solution_type = "SOLUTION_TYPE_SEARCH"
  common_config {
    company_name = "Example Company Name"
  }
  search_engine_config {
    search_tier = "SEARCH_TIER_ENTERPRISE"
    search_add_ons = ["SEARCH_ADD_ON_LLM"]
  }
}
`, context)
}

func testAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_search_engine" "basic" {
  engine_id = "tf-test-example-engine-id%{random_suffix}"
  collection_id = "default_collection"
  location = "global"
  display_name = "Updated Example Display Name"
  industry_vertical = "GENERIC"
  data_store_ids = ["test-demo-data_1705085526185"]
  solution_type = "SOLUTION_TYPE_SEARCH"
  common_config {
    company_name = "Updated Example Company Name"
  }
  search_engine_config {
    search_tier = "SEARCH_TIER_STANDARD"
    search_add_ons = ["SEARCH_ADD_ON_LLM"]
  }
}
`, context)
}
