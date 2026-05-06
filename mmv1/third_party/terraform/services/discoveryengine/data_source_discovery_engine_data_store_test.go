// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleDiscoveryEngineDataStore_basic(t *testing.T) {
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
				Config: testAccDataSourceGoogleDiscoveryEngineDataStore_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_discovery_engine_data_store.ds",
						"google_discovery_engine_data_store.res",
						[]string{
							"create_advanced_site_search",
							"skip_default_schema_creation",
						},
					),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleDiscoveryEngineDataStore_byDisplayName(t *testing.T) {
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
				Config: testAccDataSourceGoogleDiscoveryEngineDataStore_byDisplayName(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_discovery_engine_data_store.ds",
						"google_discovery_engine_data_store.res",
						[]string{
							"create_advanced_site_search",
							"skip_default_schema_creation",
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceGoogleDiscoveryEngineDataStore_byDisplayName(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "res" {
  location                     = "global"
  data_store_id                = "tf-test-ds-name%{random_suffix}"
  display_name                 = "tf-test-ds-byname%{random_suffix}"
  industry_vertical            = "GENERIC"
  content_config               = "NO_CONTENT"
  solution_types               = ["SOLUTION_TYPE_SEARCH"]
  create_advanced_site_search  = false
  skip_default_schema_creation = false
}

data "google_discovery_engine_data_store" "ds" {
  location     = google_discovery_engine_data_store.res.location
  display_name = google_discovery_engine_data_store.res.display_name
}
`, context)
}

func testAccDataSourceGoogleDiscoveryEngineDataStore_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "res" {
  location                     = "global"
  data_store_id                = "tf-test-ds-id%{random_suffix}"
  display_name                 = "tf-test-ds-datasource"
  industry_vertical            = "GENERIC"
  content_config               = "NO_CONTENT"
  solution_types               = ["SOLUTION_TYPE_SEARCH"]
  create_advanced_site_search  = false
  skip_default_schema_creation = false
}

data "google_discovery_engine_data_store" "ds" {
  location      = google_discovery_engine_data_store.res.location
  data_store_id = google_discovery_engine_data_store.res.data_store_id
}
`, context)
}
