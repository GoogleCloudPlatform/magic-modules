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
								"max_qps": knownvalue.StringExact("100"),
							}),
						),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.servicenow-basic",
							tfjsonpath.New("refresh_interval"),
							knownvalue.StringExact("172800s"),
						),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.servicenow-basic",
							tfjsonpath.New("auto_run_disabled"),
							knownvalue.Bool(false),
						),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.servicenow-basic",
							tfjsonpath.New("incremental_sync_disabled"),
							knownvalue.Bool(false),
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
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_updateJsonParams(context, "200"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_discovery_engine_data_connector.servicenow-basic", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.servicenow-basic",
							tfjsonpath.New("json_params"),
							knownvalue.StringExact(`{"max_qps":"200"}`),
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
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_updateJsonParams(context, "300"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_discovery_engine_data_connector.servicenow-basic", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.servicenow-basic",
							tfjsonpath.New("json_params"),
							knownvalue.StringExact(`{"max_qps":"300"}`),
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
  data_source                  = "servicenow"
  tag                          = "tf-test-tag-%{random_suffix}"
  params = {
    auth_type                  = "OAUTH_PASSWORD_GRANT"
    instance_uri               = "https://gcpconnector1.service-now.com/"
    client_id                  = "SECRET_MANAGER_RESOURCE_NAME"
    client_secret              = "SECRET_MANAGER_RESOURCE_NAME"
    user_account               = "connectorsuserqa@google.com"
    password                   = "SECRET_MANAGER_RESOURCE_NAME"
  }
  refresh_interval             = "86400s"
  entities {
    entity_name                = "catalog"
    params = jsonencode({
      "inclusion_filters" : {
        "knowledgeBaseSysId" : [
          "123"
        ]
      }
    })
  }
  entities {
    entity_name = "incident"
    params = jsonencode({
      "inclusion_filters" : {
        "knowledgeBaseSysId" : [
          "123"
        ]
      }
    })
  }
  entities {
    entity_name = "knowledge_base"
    params = jsonencode({
      "inclusion_filters" : {
        "knowledgeBaseSysId" : [
          "123"
        ]
      }
    })
  }
  static_ip_enabled            = false
  incremental_refresh_interval = "21600s"
  connector_modes              = ["DATA_INGESTION"]
  sync_mode                    = "PERIODIC"
  auto_run_disabled            = true
  incremental_sync_disabled    = true
}
`, context)
}

func testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "time_sleep" "wait_1_hour" {
  create_duration = "3s"
}

resource "google_discovery_engine_data_connector" "servicenow-basic" {
  depends_on                   = [time_sleep.wait_1_hour]
  location                     = "global"
  collection_id                = "tf-test-collection-id%{random_suffix}"
  collection_display_name      = "tf-test-dataconnector-servicenow"
  data_source                  = "servicenow"
  tag                          = "tf-test-tag-%{random_suffix}"
  params = {
    max_qps                    = "100"
  }
  refresh_interval             = "172800s"
  entities {
    entity_name                = "catalog"
    params                     = jsonencode({
      "inclusion_filters": {
        "knowledgeBaseSysId": [
          "456"
        ]
      }
    })
  }
  entities {
    entity_name                = "incident"
    params                     = jsonencode({
      "inclusion_filters": {
        "knowledgeBaseSysId": [
          "456"
        ]
      }
    })
  }
  entities {
    entity_name                = "knowledge_base"
    params                     = jsonencode({
      "inclusion_filters": {
        "knowledgeBaseSysId": [
          "456"
        ]
      }
    })
  }
  static_ip_enabled            = false
  incremental_refresh_interval = "21600s"
  connector_modes              = ["DATA_INGESTION"]
  sync_mode                    = "PERIODIC"
  auto_run_disabled            = false
  incremental_sync_disabled    = false
}
`, context)
}

func testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_updateJsonParams(context map[string]interface{}, maxQps string) string {
	context["max_qps"] = maxQps
	return acctest.Nprintf(`
resource "time_sleep" "wait_1_hour" {
  create_duration = "3s"
}

resource "google_discovery_engine_data_connector" "servicenow-basic" {
  depends_on                   = [time_sleep.wait_1_hour]
  location                     = "global"
  collection_id                = "tf-test-collection-id%{random_suffix}"
  collection_display_name      = "tf-test-dataconnector-servicenow"
  data_source                  = "servicenow"
  tag                          = "tf-test-tag-%{random_suffix}"
  json_params = jsonencode({
    max_qps                    = "%{max_qps}"
  })
  refresh_interval             = "172800s"
  entities {
    entity_name                = "catalog"
    params                     = jsonencode({
      "inclusion_filters": {
        "knowledgeBaseSysId": [
          "456"
        ]
      }
    })
  }
  entities {
    entity_name                = "incident"
    params                     = jsonencode({
      "inclusion_filters": {
        "knowledgeBaseSysId": [
          "456"
        ]
      }
    })
  }
  entities {
    entity_name                = "knowledge_base"
    params                     = jsonencode({
      "inclusion_filters": {
        "knowledgeBaseSysId": [
          "456"
        ]
      }
    })
  }
  static_ip_enabled            = false
  incremental_refresh_interval = "21600s"
  connector_modes              = ["DATA_INGESTION"]
  sync_mode                    = "PERIODIC"
  auto_run_disabled            = false
  incremental_sync_disabled    = false
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
