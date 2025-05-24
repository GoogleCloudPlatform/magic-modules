package dialogflowcx_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDialogflowCXTool_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowCXTool_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_tool.my_tool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowCXTool_full(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_tool.my_tool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXTool_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_dialogflow_cx_agent" "agent_tool" {
		display_name             = "tf-test-%{random_suffix}"
		location                 = "global"
		default_language_code    = "en"
		time_zone                = "America/New_York"
		description              = "ageng for tool test"
	}
	resource "google_dialogflow_cx_tool" "my_tool" {
		parent       = google_dialogflow_cx_agent.agent_tool.id
		display_name = "Example"
		description  = "Example Description"
	}
`, context)
}

func testAccDialogflowCXTool_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_dialogflow_cx_agent" "agent" {
		display_name = "{{index $.Vars "agent_name"}}"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["fr","de","es"]
		time_zone = "America/New_York"
		description = "Example description."
		avatar_uri = "https://cloud.google.com/_static/images/cloud/icons/favicons/onecloud/super_cloud.png"
		enable_stackdriver_logging = true
		enable_spell_correction    = true
		speech_to_text_settings {
			enable_speech_adaptation = true
		}
		depends_on = [
			google_discovery_engine_data_store.my_datastore
		]
	}

	resource "google_dialogflow_cx_tool" "{{$.PrimaryResourceId}}" {
		parent       = google_dialogflow_cx_agent.agent.id
		display_name = "Example"
		description  = "Example Description"
		data_store_spec {
			data_store_connections {
				data_store_type = "UNSTRUCTURED"
				data_store = "projects/${data.google_project.project.number}/locations/global/collections/default_collection/dataStores/${google_discovery_engine_data_store.my_datastore.data_store_id}"
				document_processing_mode = "DOCUMENTS"
			}
			fallback_prompt {} 
		}
		depends_on = [
			google_discovery_engine_data_store.my_datastore,
			google_dialogflow_cx_agent.agent
		]
	}

	resource "google_discovery_engine_data_store" "my_datastore" {
		location          = "global"
		data_store_id     = "datastore-tool-full-test"
		display_name      = "datastore-tool-full-test"
		industry_vertical = "GENERIC"
		content_config    = "NO_CONTENT"
		solution_types    = ["SOLUTION_TYPE_CHAT"]
	}

	data "google_project" "project" {
	}
`, context)
}
