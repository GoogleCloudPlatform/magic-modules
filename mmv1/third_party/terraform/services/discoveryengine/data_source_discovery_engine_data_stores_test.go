// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleDiscoveryEngineDataStores_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDiscoveryEngineDataStoreDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleDiscoveryEngineDataStores_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_discovery_engine_data_stores.all", "data_stores.#"),
					resource.TestCheckResourceAttrSet("data.google_discovery_engine_data_stores.all", "data_stores.0.name"),
					resource.TestCheckResourceAttrSet("data.google_discovery_engine_data_stores.all", "data_stores.0.data_store_id"),
					resource.TestCheckResourceAttrSet("data.google_discovery_engine_data_stores.all", "data_stores.0.display_name"),
					resource.TestCheckResourceAttrSet("data.google_discovery_engine_data_stores.all", "data_stores.0.industry_vertical"),
					resource.TestCheckResourceAttrSet("data.google_discovery_engine_data_stores.all", "data_stores.0.content_config"),
					resource.TestCheckResourceAttrSet("data.google_discovery_engine_data_stores.all", "data_stores.0.create_time"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleDiscoveryEngineDataStores_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "res" {
  location                     = "global"
  data_store_id                = "tf-test-ds-id%{random_suffix}"
  display_name                 = "tf-test-ds-list-datasource"
  industry_vertical            = "GENERIC"
  content_config               = "NO_CONTENT"
  solution_types               = ["SOLUTION_TYPE_SEARCH"]
  create_advanced_site_search  = false
  skip_default_schema_creation = false
}

data "google_discovery_engine_data_stores" "all" {
  location = "global"

  depends_on = [google_discovery_engine_data_store.res]
}
`, context)
}
