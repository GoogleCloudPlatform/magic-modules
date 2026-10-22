package discoveryengine_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/discoveryengine"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDiscoveryEngineDataConnector_discoveryengineDataconnectorOnedriveBasicExample_update(t *testing.T) {
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
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorOnedriveBasicExample_basic(context),
				Check:  testAccCheckDiscoveryEngineDataConnectorParamsTenantId(t, "google_discovery_engine_data_connector.onedrive-basic", "tenant_id_1"),
			},
			{
				ResourceName:            "google_discovery_engine_data_connector.onedrive-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_run_disabled", "collection_display_name", "collection_id", "incremental_sync_disabled", "location", "params", "state", "update_time"},
			},
			{
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorOnedriveBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_discovery_engine_data_connector.onedrive-basic", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.onedrive-basic",
							tfjsonpath.New("params"),
							knownvalue.MapExact(map[string]knownvalue.Check{
								"auth_type":     knownvalue.StringExact("OAUTH"),
								"client_id":     knownvalue.StringExact("client_id_1"),
								"client_secret": knownvalue.StringExact("client_secret_1"),
								"tenant_id":     knownvalue.StringExact("tenant_id_2"),
							}),
						),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.onedrive-basic",
							tfjsonpath.New("refresh_interval"),
							knownvalue.StringExact("172800s"),
						),
					},
				},
				Check: testAccCheckDiscoveryEngineDataConnectorParamsTenantId(t, "google_discovery_engine_data_connector.onedrive-basic", "tenant_id_2"),
			},
			{
				ResourceName:            "google_discovery_engine_data_connector.onedrive-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_run_disabled", "collection_display_name", "collection_id", "incremental_sync_disabled", "location", "params", "update_time"},
			},
			{
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorOnedriveBasicExample_updateJsonParams(context, "tenant_id_3"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_discovery_engine_data_connector.onedrive-basic", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.onedrive-basic",
							tfjsonpath.New("json_params"),
							knownvalue.StringExact(`{"auth_type":"OAUTH","client_id":"client_id_1","client_secret":"client_secret_1","tenant_id":"tenant_id_3"}`),
						),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.onedrive-basic",
							tfjsonpath.New("params"),
							knownvalue.MapExact(map[string]knownvalue.Check{}),
						),
					},
				},
				Check: testAccCheckDiscoveryEngineDataConnectorParamsTenantId(t, "google_discovery_engine_data_connector.onedrive-basic", "tenant_id_3"),
			},
			{
				ResourceName:            "google_discovery_engine_data_connector.onedrive-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_run_disabled", "collection_display_name", "collection_id", "incremental_sync_disabled", "json_params", "location", "params", "update_time"},
			},
			{
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorOnedriveBasicExample_updateJsonParams(context, "tenant_id_4"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_discovery_engine_data_connector.onedrive-basic", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(
							"google_discovery_engine_data_connector.onedrive-basic",
							tfjsonpath.New("json_params"),
							knownvalue.StringExact(`{"auth_type":"OAUTH","client_id":"client_id_1","client_secret":"client_secret_1","tenant_id":"tenant_id_4"}`),
						),
					},
				},
				Check: testAccCheckDiscoveryEngineDataConnectorParamsTenantId(t, "google_discovery_engine_data_connector.onedrive-basic", "tenant_id_4"),
			},
			{
				ResourceName:            "google_discovery_engine_data_connector.onedrive-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_run_disabled", "collection_display_name", "collection_id", "incremental_sync_disabled", "json_params", "location", "params", "update_time"},
			},
		},
	})
}

func testAccCheckDiscoveryEngineDataConnectorParamsTenantId(t *testing.T, resourceName, expectedTenantId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		config := acctest.GoogleProviderConfig(t)
		url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(discoveryengine.Product, config)+"projects/{{project}}/locations/{{location}}/collections/{{collection_id}}/dataConnector")
		if err != nil {
			return err
		}

		billingProject := ""
		if config.BillingProject != "" {
			billingProject = config.BillingProject
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return err
		}

		params, ok := res["params"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected params map in API response, got: %#v", res["params"])
		}
		gotTenantId, ok := params["tenant_id"].(string)
		if !ok || gotTenantId != expectedTenantId {
			return fmt.Errorf("expected API params.tenant_id to be %q, got %q", expectedTenantId, gotTenantId)
		}
		return nil
	}
}

func testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorOnedriveBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_discovery_engine_data_connector" "onedrive-basic" {
  location                     = "global"
  collection_id                = "tf-test-collection-id%{random_suffix}"
  collection_display_name      = "tf-test-dataconnector-onedrive"
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

func testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorOnedriveBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "time_sleep" "wait_1_minute" {
  create_duration = "60s"
}

resource "google_discovery_engine_data_connector" "onedrive-basic" {
  depends_on                   = [time_sleep.wait_1_minute]
  location                     = "global"
  collection_id                = "tf-test-collection-id%{random_suffix}"
  collection_display_name      = "tf-test-dataconnector-onedrive"
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

func testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorOnedriveBasicExample_updateJsonParams(context map[string]interface{}, tenantId string) string {
	context["tenant_id"] = tenantId
	return acctest.Nprintf(`
resource "time_sleep" "wait_1_minute" {
  create_duration = "60s"
}

resource "google_discovery_engine_data_connector" "onedrive-basic" {
  depends_on                   = [time_sleep.wait_1_minute]
  location                     = "global"
  collection_id                = "tf-test-collection-id%{random_suffix}"
  collection_display_name      = "tf-test-dataconnector-onedrive"
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

func TestUnitDiscoveryEngineDataConnector_flattenImportAndReadHydration(t *testing.T) {
	t.Run("Import hydration with empty prior state", func(t *testing.T) {
		d := discoveryengine.ResourceDiscoveryEngineDataConnector().Data(nil)

		res := map[string]interface{}{
			"name":            "projects/test-project/locations/global/collections/test-coll/dataConnector",
			"dataSource":      "jira",
			"refreshInterval": "86400s",
			"autoRunDisabled": true,
			"actionConfig": map[string]interface{}{
				"isActionConfigured": true,
				"actionParams": map[string]interface{}{
					"auth_type":    "OAUTH",
					"instance_uri": "https://example.atlassian.net",
				},
				"createBapConnection": true,
			},
		}

		err := discoveryengine.ResourceDiscoveryEngineDataConnectorFlatten(d, nil, res, nil, "test-project", "test-ua", "test-project", "test-url", nil)
		if err != nil {
			t.Fatalf("unexpected error from ResourceDiscoveryEngineDataConnectorFlatten: %v", err)
		}

		if got := d.Get("auto_run_disabled"); got != true {
			t.Errorf("expected auto_run_disabled = true on import, got %v", got)
		}
		if got := d.Get("action_config.0.action_params.auth_type"); got != "OAUTH" {
			t.Errorf("expected action_config.0.action_params.auth_type = OAUTH on import, got %v", got)
		}
		if got := d.Get("action_config.0.create_bap_connection"); got != true {
			t.Errorf("expected action_config.0.create_bap_connection = true on import, got %v", got)
		}
	})

	t.Run("Read hydration preserves state secrets in action_params", func(t *testing.T) {
		d := discoveryengine.ResourceDiscoveryEngineDataConnector().Data(nil)
		if err := d.Set("action_config", []interface{}{
			map[string]interface{}{
				"action_params": map[string]interface{}{
					"auth_type":     "OAUTH",
					"instance_uri":  "https://old.atlassian.net",
					"client_secret": "SECRET_MANAGER_RESOURCE_NAME",
				},
				"create_bap_connection": true,
			},
		}); err != nil {
			t.Fatalf("failed to seed prior state: %v", err)
		}

		res := map[string]interface{}{
			"name":            "projects/test-project/locations/global/collections/test-coll/dataConnector",
			"dataSource":      "jira",
			"refreshInterval": "86400s",
			"autoRunDisabled": true,
			"actionConfig": map[string]interface{}{
				"isActionConfigured": true,
				"actionParams": map[string]interface{}{
					"auth_type":    "OAUTH",
					"instance_uri": "https://updated.atlassian.net",
				},
			},
		}

		err := discoveryengine.ResourceDiscoveryEngineDataConnectorFlatten(d, nil, res, nil, "test-project", "test-ua", "test-project", "test-url", nil)
		if err != nil {
			t.Fatalf("unexpected error from ResourceDiscoveryEngineDataConnectorFlatten: %v", err)
		}

		if got := d.Get("action_config.0.action_params.client_secret"); got != "SECRET_MANAGER_RESOURCE_NAME" {
			t.Errorf("expected client_secret preserved in state, got %v", got)
		}
		if got := d.Get("action_config.0.action_params.instance_uri"); got != "https://updated.atlassian.net" {
			t.Errorf("expected updated instance_uri from API, got %v", got)
		}
		if got := d.Get("action_config.0.create_bap_connection"); got != true {
			t.Errorf("expected create_bap_connection preserved from state, got %v", got)
		}
	})
}
