package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/discoveryengine"
)

func TestAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_update(t *testing.T) {
	// Skips this update test due to duration and flakiness.
	t.Skip()

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
				ImportStateVerifyIgnore: []string{"collection_display_name", "collection_id", "location", "params", "update_time", "action_config.0.action_params", "action_config.0.create_bap_connection"},
			},
			{
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorServicenowBasicExample_update(context),
			},
			{
				ResourceName:            "google_discovery_engine_data_connector.servicenow-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collection_display_name", "collection_id", "location", "params", "update_time", "action_config.0.action_params", "action_config.0.create_bap_connection"},
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
  data_source_version          = 3
  params = {
    auth_type                  = "OAUTH_PASSWORD_GRANT"
    instance_uri               = "https://gcpconnector1.service-now.com/"
    client_id                  = "SECRET_MANAGER_RESOURCE_NAME"
    client_secret              = "SECRET_MANAGER_RESOURCE_NAME"
    static_ip_enabled          = "false"
    user_account               = "connectorsuserqa@google.com"
    password                   = "SECRET_MANAGER_RESOURCE_NAME"
  }
  refresh_interval             = "86400s"
  entities {
    entity_name                = "catalog"
    key_property_mappings = {
      title       = "title"
      description = "short_description"
    }
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
  destination_configs {
    key = "url"
    destinations {
      host = "https://gcpconnector1.service-now.com/"
      port = 123
    }
    params.                    = jsonencode({
      "destination_type": "private"
    })
  }
  incremental_refresh_interval = "21600s"
  connector_modes              = ["DATA_INGESTION", "ACTIONS"]
  sync_mode                    = "PERIODIC"
  auto_run_disabled            = true
  incremental_sync_disabled    = true
  action_config {
    action_params = {
      instance_uri  = "https://example.atlassian.net"
      instance_id   = "unused"
      client_id     = "unused"
      client_secret = "unused"
      auth_type     = "OAUTH"
    }
    create_bap_connection = true
  }
  bap_config {
    supported_connector_modes = ["ACTIONS"]
    enabled_actions = [
      "create_issue",
      "update_issue",
      "change_issue_status",
      "create_comment",
      "update_comment",
      "upload_attachment",
    ]
  }
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
  destination_configs {
    key = "url"
    destinations {
      host = "https://gcpconnector1.service-now.com/"
      port = 123
    }
    params                     = jsonencode({
      "destination_type": "private"
    })
  }
  incremental_refresh_interval = "21600s"
  connector_modes              = ["DATA_INGESTION", "ACTIONS"]
  sync_mode                    = "PERIODIC"
  auto_run_disabled            = false
  incremental_sync_disabled    = false
  action_config {
    action_params = {
      instance_uri  = "https://example.atlassian.net"
      instance_id   = "unused"
      client_id     = "unused"
      client_secret = "unused"
      auth_type     = "OAUTH"
    }
    create_bap_connection = true
  }
  bap_config {
    supported_connector_modes = ["ACTIONS"]
    enabled_actions = [
      "create_issue",
      "update_issue",
      "change_issue_status",
      "create_comment",
      "update_comment",
      "upload_attachment",
    ]
  }
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
