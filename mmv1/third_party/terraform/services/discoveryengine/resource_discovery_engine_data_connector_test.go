package discoveryengine_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccDiscoveryEngineDataConnector_discoveryengineDataconnectorJiraBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"kms_key_name":  acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us", "tftest-shared-key-dataconnector-0").CryptoKey.Name,
		"client_id":     "tf-test-client-id",
		"client_secret": "tf-test-client-secret",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorJiraBasicExample_basic(context),
			},
			{
				ResourceName:            "google_discovery_engine_data_connector.jira-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collection_display_name", "collection_id", "location", "params"},
			},
			{
				Config: testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorJiraBasicExample_update(context),
			},
			{
				ResourceName:            "google_discovery_engine_data_connector.jira-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collection_display_name", "collection_id", "location", "params"},
			},
		},
	})
}

func testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorJiraBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_connector" "jira-basic" {
  location                  = "global"
  collection_id             = "tf-test-collection-id%{random_suffix}"
  collection_display_name   = "tf-test-dataconnector-jira"
  data_source               = "jira"
  params = {
      instance_id           = "33db20a3-dc45-4305-a505-d70b68599840"
      instance_uri          = "https://vaissptbots1.atlassian.net/"
      client_secret         = "%{client_secret}"
      client_id             = "%{client_id}"
      refresh_token         = "fill-in-the-blank"
  }
  refresh_interval          = "86400s"
  entities {
      entity_name           = "project"
  }
  entities {
      entity_name           = "issue"
  }
  entities {
      entity_name           = "attachment"
  }
  entities {
      entity_name           = "comment"
  }
  entities {
      entity_name           = "worklog"
  }
  static_ip_enabled         = true
}
`, context)
}

func testAccDiscoveryEngineDataConnector_discoveryengineDataconnectorJiraBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_connector" "jira-basic" {
  location                  = "global"
  collection_id             = "tf-test-collection-id%{random_suffix}"
  collection_display_name   = "tf-test-dataconnector-jira"
  data_source               = "jira"
  params = {
      instance_id           = "33db20a3-dc45-4305-a505-d70b68599840"
      instance_uri          = "https://vaissptbots1.atlassian.net/"
      client_secret         = "%{client_secret}"
      client_id             = "%{client_id}"
      refresh_token         = "fill-in-the-blank"
  }
  refresh_interval          = "86400s"
  entities {
      entity_name           = "project"
  }
  entities {
      entity_name           = "issue"
  }
  entities {
      entity_name           = "attachment"
  }
  entities {
      entity_name           = "comment"
  }
  entities {
      entity_name           = "worklog"
  }
  static_ip_enabled         = true
  kms_key_name              = "%{kms_key_name}"
}
`, context)
}
