// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDiscoveryEngineTargetSite_discoveryengineTargetsiteBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineTargetSite_discoveryengineTargetsiteBasicExample_basic(context),
			},
			{
				ResourceName:            "google_discovery_engine_target_site.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_store_id", "location", "project", "provided_uri_pattern", "target_site_id"},
			},
			{
				Config: testAccDiscoveryEngineTargetSite_discoveryengineTargetsiteBasicExample_update(context),
			},
			{
				ResourceName:            "google_discovery_engine_target_site.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_store_id", "location", "project", "provided_uri_pattern", "target_site_id"},
			},
		},
	})
}

func testAccDiscoveryEngineTargetSite_discoveryengineTargetsiteBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_target_site" "basic" {
  location                    = google_discovery_engine_data_store.basic.location
  data_store_id               = google_discovery_engine_data_store.basic.data_store_id
  provided_uri_pattern        = "http://cloud.google.com/docs/*"
  type                        = "INCLUDE"
  exact_match                 = false
}

resource "google_discovery_engine_data_store" "basic" {
  location                     = "global"
  data_store_id                = "tf-test-data-store-id%{random_suffix}"
  display_name                 = "tf-test-website-datastore"
  industry_vertical            = "GENERIC"
  content_config               = "PUBLIC_WEBSITE"
  solution_types               = ["SOLUTION_TYPE_SEARCH"]
  create_advanced_site_search  = false
  skip_default_schema_creation = false
}
`, context)
}

func testAccDiscoveryEngineTargetSite_discoveryengineTargetsiteBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_target_site" "basic" {
  location                    = google_discovery_engine_data_store.basic.location
  data_store_id               = google_discovery_engine_data_store.basic.data_store_id
  provided_uri_pattern        = "https://cloud.google.com/generative-ai-app-builder/docs/*"
  type                        = "INCLUDE"
  exact_match                 = false
}

resource "google_discovery_engine_data_store" "basic" {
  location                     = "global"
  data_store_id                = "tf-test-data-store-id%{random_suffix}"
  display_name                 = "tf-test-website-datastore"
  industry_vertical            = "GENERIC"
  content_config               = "PUBLIC_WEBSITE"
  solution_types               = ["SOLUTION_TYPE_SEARCH"]
  create_advanced_site_search  = false
  skip_default_schema_creation = false
}
`, context)
}
