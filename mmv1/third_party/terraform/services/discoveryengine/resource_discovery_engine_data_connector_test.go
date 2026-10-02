package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/discoveryengine"
)

func TestAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_update(t *testing.T) {
	// TODO(b/560162779): DO NOT SUBMIT - Uncomment t.Skip() before marking PR ready for review!
	// Skips this update test due to duration and flakiness.
	// t.Skip()

	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_basic(context),
			},
			{
				ResourceName:            "google_discovery_engine_data_connector.servicenow-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action_state", "auto_run_disabled", "collection_display_name", "collection_id", "errors", "incremental_sync_disabled", "location", "params", "state", "sync_mode", "update_time"},
			},
			{
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_discovery_engine_data_connector.servicenow-basic", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.servicenow-basic",
							tfjsonpath.New("params"),
							knownvalue.MapExact(map[string]knownvalue.Check{
								"auth_type":     knownvalue.StringExact("OAUTH"),
								"client_id":     knownvalue.StringExact("client_id_1"),
								"client_secret": knownvalue.StringExact("client_secret_1"),
								"tenant_id":     knownvalue.StringExact("tenant_id_2"),
							}),
						),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.servicenow-basic",
							tfjsonpath.New("refresh_interval"),
							knownvalue.StringExact("172800s"),
						),
					},
				},
			},
			{
				ResourceName:            "google_discovery_engine_data_connector.servicenow-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action_state", "auto_run_disabled", "collection_display_name", "collection_id", "errors", "incremental_sync_disabled", "location", "params", "state", "sync_mode", "update_time"},
			},
			{
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_updateJsonParams(context, "tenant_id_3"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_discovery_engine_data_connector.servicenow-basic", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.servicenow-basic",
							tfjsonpath.New("json_params"),
							knownvalue.StringExact(`{"auth_type":"OAUTH","client_id":"client_id_1","client_secret":"client_secret_1","tenant_id":"tenant_id_3"}`),
						),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.servicenow-basic",
							tfjsonpath.New("params"),
							knownvalue.Null(),
						),
					},
				},
			},
			{
				ResourceName:            "google_discovery_engine_data_connector.servicenow-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action_state", "auto_run_disabled", "collection_display_name", "collection_id", "errors", "incremental_sync_disabled", "location", "json_params", "state", "sync_mode", "update_time"},
			},
			{
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_updateJsonParams(context, "tenant_id_4"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_discovery_engine_data_connector.servicenow-basic", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.servicenow-basic",
							tfjsonpath.New("json_params"),
							knownvalue.StringExact(`{"auth_type":"OAUTH","client_id":"client_id_1","client_secret":"client_secret_1","tenant_id":"tenant_id_4"}`),
						),
					},
				},
			},
			{
				ResourceName:            "google_discovery_engine_data_connector.servicenow-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action_state", "auto_run_disabled", "collection_display_name", "collection_id", "errors", "incremental_sync_disabled", "location", "json_params", "state", "sync_mode", "update_time"},
			},
		},
	})
}

func testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_discovery_engine_data_connector" "servicenow-basic" {
  location                     = "global"
  collection_id                = "tf-test-collection-id%{random_suffix}"
  collection_display_name      = "tf-test-dataconnector-servicenow"
  data_source                  = "onedrive_federated_search"
  tag                          = "tf-test-tag-%{random_suffix}"
  params = {
    auth_type                  = "OAUTH"
    client_id                  = "client_id_1"
    client_secret              = "client_secret_1"
    tenant_id                  = "tenant_id_1"
  }
  refresh_interval             = "86400s"
  entities {
    entity_name                = "file"
    params = jsonencode({
      "custom_id" : "123"
    })
  }
  static_ip_enabled            = false
  incremental_refresh_interval = "21600s"
  connector_modes              = ["FEDERATED"]
  auto_run_disabled            = true
  incremental_sync_disabled    = true
}
`, context)
}

func testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "time_sleep" "wait_1_hour" {
  create_duration = "60s"
}

resource "google_discovery_engine_data_connector" "servicenow-basic" {
  depends_on                   = [time_sleep.wait_1_hour]
  location                     = "global"
  collection_id                = "tf-test-collection-id%{random_suffix}"
  collection_display_name      = "tf-test-dataconnector-servicenow"
  data_source                  = "onedrive_federated_search"
  tag                          = "tf-test-tag-%{random_suffix}"
  params = {
    auth_type                  = "OAUTH"
    client_id                  = "client_id_1"
    client_secret              = "client_secret_1"
    tenant_id                  = "tenant_id_2"
  }
  refresh_interval             = "172800s"
  entities {
    entity_name                = "file"
    params                     = jsonencode({
      "custom_id": "456"
    })
  }
  static_ip_enabled            = false
  incremental_refresh_interval = "21600s"
  connector_modes              = ["FEDERATED"]
  auto_run_disabled            = true
  incremental_sync_disabled    = true
}
`, context)
}

func testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_updateJsonParams(context map[string]interface{}, tenantId string) string {
	context["tenant_id"] = tenantId
	return acctest.Nprintf(`
resource "time_sleep" "wait_1_hour" {
  create_duration = "60s"
}

resource "google_discovery_engine_data_connector" "servicenow-basic" {
  depends_on                   = [time_sleep.wait_1_hour]
  location                     = "global"
  collection_id                = "tf-test-collection-id%{random_suffix}"
  collection_display_name      = "tf-test-dataconnector-servicenow"
  data_source                  = "onedrive_federated_search"
  tag                          = "tf-test-tag-%{random_suffix}"
  json_params = jsonencode({
    auth_type                  = "OAUTH"
    client_id                  = "client_id_1"
    client_secret              = "client_secret_1"
    tenant_id                  = "%{tenant_id}"
  })
  refresh_interval             = "172800s"
  entities {
    entity_name                = "file"
    params                     = jsonencode({
      "custom_id": "456"
    })
  }
  static_ip_enabled            = false
  incremental_refresh_interval = "21600s"
  connector_modes              = ["FEDERATED"]
  auto_run_disabled            = true
  incremental_sync_disabled    = true
}
`, context)
}

func TestDiscoveryEngineDataConnector_DataConnectorEntitiesParamsDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"Old empty JSON": {
			Old:                "{}",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"New empty JSON": {
			Old:                "",
			New:                "{}",
			ExpectDiffSuppress: true,
		},
		"Diff not supressed": {
			Old:                "123",
			New:                "",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if discoveryengine.DataConnectorJsonStructFieldsDiffSuppress("entities_params_diff_supress", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, %q => %q expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}
