package ces_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck" // Add this import
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCESTool_cesToolClientFunctionBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCESToolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCESTool_cesToolClientFunctionBasicExample_full(context),
			},
			{
				ResourceName:            "google_ces_tool.ces_tool_client_function_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app", "location", "tool_id"},
			},
			{
				Config: testAccCESTool_cesToolClientFunctionBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_ces_tool.ces_tool_client_function_basic", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_ces_tool.ces_tool_client_function_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app", "location", "tool_id"},
			},
		},
	})
}

func testAccCESTool_cesToolClientFunctionBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "my-app" {
    location     = "us"
    display_name = "tf-test-my-app%{random_suffix}"
    app_id       = "tf-test-app-id%{random_suffix}"
    time_zone_settings {   
        time_zone = "America/Los_Angeles"
    }
}
resource "google_ces_tool" "ces_tool_client_function_basic" {
    location     = "us"
    app          = google_ces_app.my-app.name
    tool_id      = "tf_test_ces_tool_basic1%{random_suffix}"
    execution_type = "SYNCHRONOUS"
    client_function {
        name = "tf_test_ces_tool_client_function_basic%{random_suffix}"
        description = "example-description"
        parameters {
            description = "schema description"
            type        = "ARRAY"
            nullable    = true
            required = ["some_property"]
            enum = ["VALUE_A", "VALUE_B"]
            ref = "#/defs/MyDefinition"
            unique_items = true
            defs = jsonencode({
                SimpleString = {
                type        = "STRING"
                description = "A simple string definition"
            }})
            any_of = jsonencode([
                {
                type        = "STRING"
                description = "any_of option 1: string"
                },])
            default = jsonencode(
                false)
            prefix_items = jsonencode([
                {
                type        = "ARRAY"
                description = "prefix item 1"
                },])
            additional_properties = jsonencode(
                {
                type        = "BOOLEAN"
                })
            properties = jsonencode({
                name = {
                type        = "STRING"
                description = "A name"
            }})
            items = jsonencode({
                type        = "ARRAY"
                description = "An array"
            })
        }
        response {
            description = "schema description"
            type        = "ARRAY"
            nullable    = true
            required = ["some_property"]
            enum = ["VALUE_A", "VALUE_B"]
            ref = "#/defs/MyDefinition"
            unique_items = true
            defs = jsonencode({
                SimpleString = {
                type        = "STRING"
                description = "A simple string definition"
            }})
            any_of = jsonencode([
                {
                type        = "STRING"
                description = "any_of option 1: string"
                },])
            default = jsonencode(
                false)
            prefix_items = jsonencode([
                {
                type        = "ARRAY"
                description = "prefix item 1"
                },])
            additional_properties = jsonencode(
                {
                type        = "BOOLEAN"
                })
            properties = jsonencode({
                name = {
                type        = "STRING"
                description = "A name"
            }})
            items = jsonencode({
                type        = "ARRAY"
                description = "An array"
            })
        }
    }
}
`, context)
}

func testAccCESTool_cesToolClientFunctionBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "my-app" {
    location     = "us"
    display_name = "tf-test-my-app%{random_suffix}"
    app_id       = "tf-test-app-id%{random_suffix}"
    time_zone_settings {   
        time_zone = "America/Los_Angeles"
    }
}
resource "google_ces_tool" "ces_tool_client_function_basic" {
    location     = "us"
    app          = google_ces_app.my-app.name
    tool_id      = "tf_test_ces_tool_basic1%{random_suffix}"
    execution_type = "SYNCHRONOUS"
    client_function {
        name = "tf_test_ces_tool_client_function_basic%{random_suffix}"
        description = "example-description-updated"
        parameters {
            description = "schema description"
            type        = "ARRAY"
            nullable    = true
            required = ["some_property"]
            enum = ["VALUE_A", "VALUE_B"]
            ref = "#/defs/MyDefinition"
            unique_items = true
            defs = jsonencode({
                SimpleString = {
                type        = "STRING"
                description = "A simple string definition"
            }})
            any_of = jsonencode([
                {
                type        = "STRING"
                description = "any_of option 1: string"
                },])
            default = jsonencode(
                false)
            prefix_items = jsonencode([
                {
                type        = "ARRAY"
                description = "prefix item 1"
                },])
            additional_properties = jsonencode(
                {
                type        = "BOOLEAN"
                })
            properties = jsonencode({
                name = {
                type        = "STRING"
                description = "A name"
            }})
            items = jsonencode({
                type        = "ARRAY"
                description = "An array"
            })
        }
        response {
            description = "schema description"
            type        = "ARRAY"
            nullable    = true
            required = ["some_property"]
            enum = ["VALUE_A", "VALUE_B"]
            ref = "#/defs/MyDefinition"
            unique_items = true
            defs = jsonencode({
                SimpleString = {
                type        = "STRING"
                description = "A simple string definition"
            }})
            any_of = jsonencode([
                {
                type        = "STRING"
                description = "any_of option 1: string"
                },])
            default = jsonencode(
                false)
            prefix_items = jsonencode([
                {
                type        = "ARRAY"
                description = "prefix item 1"
                },])
            additional_properties = jsonencode(
                {
                type        = "BOOLEAN"
                })
            properties = jsonencode({
                name = {
                type        = "STRING"
                description = "A name"
            }})
            items = jsonencode({
                type        = "ARRAY"
                description = "An array"
            })
        }
    }
}
`, context)
}

func TestAccCESTool_cesToolDataStoreToolEngineSourceBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCESToolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCESTool_cesToolDataStoreToolEngineSourceBasicExample_full(context),
			},
			{
				ResourceName:            "google_ces_tool.ces_tool_data_store_tool_engine_source_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app", "location", "tool_id"},
			},
			{
				Config: testAccCESTool_cesToolDataStoreToolEngineSourceBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_ces_tool.ces_tool_data_store_tool_engine_source_basic", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_ces_tool.ces_tool_data_store_tool_engine_source_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app", "location", "tool_id"},
			},
		},
	})
}

func testAccCESTool_cesToolDataStoreToolEngineSourceBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "test_data_store" {
  location                    = "global"
  data_store_id               = "tf_test_data_store_id%{random_suffix}"
  display_name                = "Structured datastore"
  industry_vertical           = "GENERIC"
  content_config              = "NO_CONTENT"
  solution_types              = ["SOLUTION_TYPE_CHAT"]
}
resource "google_discovery_engine_data_store" "test_data_store_2" {
  location                    = google_discovery_engine_data_store.test_data_store.location
  data_store_id               = "tf_test_data_store_id_2%{random_suffix}"
  display_name                = "Structured datastore 2"
  industry_vertical           = "GENERIC"
  content_config              = "NO_CONTENT"
  solution_types              = ["SOLUTION_TYPE_CHAT"]
}
resource "google_discovery_engine_chat_engine" "primary" {
  engine_id = "tf_test_engine_id%{random_suffix}"
  collection_id = "default_collection"
  location = google_discovery_engine_data_store.test_data_store.location
  display_name = "Chat engine 2"
  industry_vertical = "GENERIC"
  data_store_ids = [google_discovery_engine_data_store.test_data_store.data_store_id, google_discovery_engine_data_store.test_data_store_2.data_store_id]
  common_config {
    company_name = "test-company"
  }
  chat_engine_config {
    agent_creation_config {
    business = "test business name"
    default_language_code = "en"
    time_zone = "America/Los_Angeles"
    }
  }
}
resource "google_ces_app" "my-app" {
    location     = "us"
    display_name = "tf-test-my-app%{random_suffix}"
    app_id       = "tf-test-app-id%{random_suffix}"
    time_zone_settings {   
        time_zone = "America/Los_Angeles"
    }
}
resource "google_ces_tool" "ces_tool_data_store_tool_engine_source_basic" {
    location       = "us"
    app            = google_ces_app.my-app.name
    tool_id        = "tf_test_ces_tool_basic2%{random_suffix}"
    execution_type = "SYNCHRONOUS"
    data_store_tool {
        name = "example-tool"
        description = "example-description"
        boost_specs {
            data_stores = [
                google_discovery_engine_data_store.test_data_store_2.name,
            ]
            spec {
                condition_boost_specs {
                    condition = "(lang_code: ANY(\"en\", \"fr\"))"
                    boost = 1
                    boost_control_spec {
                        field_name = "example-field"
                        attribute_type = "NUMERICAL"
                        interpolation_type = "LINEAR"
                        control_points {
                            attribute_value = 1
                            boost_amount = 1
                        }
                    }
                }
            }
        }
        modality_configs {
            modality_type = "TEXT"
            rewriter_config {
                model_settings {
                    model = "gemini-2.5-flash"
                    temperature = 1
                }
                prompt = "example-prompt"
                disabled = false
            }
            summarization_config {
                model_settings {
                    model = "gemini-2.5-flash"
                    temperature = 1
                }
                prompt = "example-prompt"
                disabled = false
            }
            grounding_config {
                grounding_level = 3
                disabled = false
            }
        }

        engine_source {
            engine = google_discovery_engine_chat_engine.primary.name
            data_store_sources {
                filter = "example_field: ANY(\"specific_example\")"
                data_store {
                    name = google_discovery_engine_data_store.test_data_store_2.name
                }
            }
            filter = "example_field: ANY(\"specific_example\")"
        }
        max_results = 5
    }
}
`, context)
}

func testAccCESTool_cesToolDataStoreToolEngineSourceBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "test_data_store" {
  location                    = "global"
  data_store_id               = "tf_test_data_store_id%{random_suffix}"
  display_name                = "Structured datastore"
  industry_vertical           = "GENERIC"
  content_config              = "NO_CONTENT"
  solution_types              = ["SOLUTION_TYPE_CHAT"]
}
resource "google_discovery_engine_data_store" "test_data_store_2" {
  location                    = google_discovery_engine_data_store.test_data_store.location
  data_store_id               = "tf_test_data_store_id_2%{random_suffix}"
  display_name                = "Structured datastore 2"
  industry_vertical           = "GENERIC"
  content_config              = "NO_CONTENT"
  solution_types              = ["SOLUTION_TYPE_CHAT"]
}
resource "google_discovery_engine_chat_engine" "primary" {
  engine_id = "tf_test_engine_id%{random_suffix}"
  collection_id = "default_collection"
  location = google_discovery_engine_data_store.test_data_store.location
  display_name = "Chat engine 2"
  industry_vertical = "GENERIC"
  data_store_ids = [google_discovery_engine_data_store.test_data_store.data_store_id, google_discovery_engine_data_store.test_data_store_2.data_store_id]
  common_config {
    company_name = "test-company"
  }
  chat_engine_config {
    agent_creation_config {
    business = "test business name"
    default_language_code = "en"
    time_zone = "America/Los_Angeles"
    }
  }
}
resource "google_ces_app" "my-app" {
    location     = "us"
    display_name = "tf-test-my-app%{random_suffix}"
    app_id       = "tf-test-app-id%{random_suffix}"
    time_zone_settings {   
        time_zone = "America/Los_Angeles"
    }
}
resource "google_ces_tool" "ces_tool_data_store_tool_engine_source_basic" {
    location       = "us"
    app            = google_ces_app.my-app.name
    tool_id        = "tf_test_ces_tool_basic2%{random_suffix}"
    execution_type = "SYNCHRONOUS"
    data_store_tool {
        name = "example-tool"
        description = "example-description-updated"
        boost_specs {
            data_stores = [
                google_discovery_engine_data_store.test_data_store_2.name,
            ]
            spec {
                condition_boost_specs {
                    condition = "(lang_code: ANY(\"en\", \"fr\"))"
                    boost = 1
                    boost_control_spec {
                        field_name = "example-field"
                        attribute_type = "NUMERICAL"
                        interpolation_type = "LINEAR"
                        control_points {
                            attribute_value = 1
                            boost_amount = 1
                        }
                    }
                }
            }
        }
        modality_configs {
            modality_type = "TEXT"
            rewriter_config {
                model_settings {
                    model = "gemini-2.5-flash"
                    temperature = 1
                }
                prompt = "example-prompt"
                disabled = false
            }
            summarization_config {
                model_settings {
                    model = "gemini-2.5-flash"
                    temperature = 1
                }
                prompt = "example-prompt"
                disabled = false
            }
            grounding_config {
                grounding_level = 3
                disabled = false
            }
        }

        engine_source {
            engine = google_discovery_engine_chat_engine.primary.name
            data_store_sources {
                filter = "example_field: ANY(\"updated_example\")"
                data_store {
                    name = google_discovery_engine_data_store.test_data_store_2.name
                }
            }
            filter = "example_field: ANY(\"updated_example\")"
        }
        max_results = 5
    }
}
`, context)
}

func TestAccCESTool_cesToolGoogleSearchToolBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCESToolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCESTool_cesToolGoogleSearchToolBasicExample_full(context),
			},
			{
				ResourceName:            "google_ces_tool.ces_tool_google_search_tool_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app", "location", "tool_id"},
			},
			{
				Config: testAccCESTool_cesToolGoogleSearchToolBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_ces_tool.ces_tool_google_search_tool_basic", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_ces_tool.ces_tool_google_search_tool_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app", "location", "tool_id"},
			},
		},
	})
}

func testAccCESTool_cesToolGoogleSearchToolBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "my-app" {
    location     = "us"
    display_name = "tf-test-my-app%{random_suffix}"
    app_id       = "tf-test-app-id%{random_suffix}"
    time_zone_settings {   
        time_zone = "America/Los_Angeles"
    }
}
resource "google_ces_tool" "ces_tool_google_search_tool_basic" {
    location       = "us"
    app            = google_ces_app.my-app.name
    tool_id        = "tf_test_ces_tool_basic4%{random_suffix}"
    execution_type = "SYNCHRONOUS"
    google_search_tool {
        name            = "example-tool"
        description     = "example-description"
        exclude_domains = ["example.com", "example2.com"]
    }
}
`, context)
}

func testAccCESTool_cesToolGoogleSearchToolBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "my-app" {
    location     = "us"
    display_name = "tf-test-my-app%{random_suffix}"
    app_id       = "tf-test-app-id%{random_suffix}"
    time_zone_settings {   
        time_zone = "America/Los_Angeles"
    }
}
resource "google_ces_tool" "ces_tool_google_search_tool_basic" {
    location       = "us"
    app            = google_ces_app.my-app.name
    tool_id        = "tf_test_ces_tool_basic4%{random_suffix}"
    execution_type = "SYNCHRONOUS"
    google_search_tool {
        name            = "example-tool"
        description     = "example-description-updated"
        exclude_domains = ["example.com", "example2.com"]
    }
}
`, context)
}

func TestAccCESTool_cesToolPythonFunctionBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCESToolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCESTool_cesToolPythonFunctionBasicExample_full(context),
			},
			{
				ResourceName:            "google_ces_tool.ces_tool_python_function_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app", "location", "tool_id"},
			},
			{
				Config: testAccCESTool_cesToolPythonFunctionBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_ces_tool.ces_tool_python_function_basic", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_ces_tool.ces_tool_python_function_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app", "location", "tool_id"},
			},
		},
	})
}

func testAccCESTool_cesToolPythonFunctionBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "my-app" {
    location     = "us"
    display_name = "tf-test-my-app%{random_suffix}"
    app_id       = "tf-test-app-id%{random_suffix}"
    time_zone_settings {   
        time_zone = "America/Los_Angeles"
    }
}
resource "google_ces_tool" "ces_tool_python_function_basic" {
    location       = "us"
    app            = google_ces_app.my-app.name
    tool_id        = "tf_test_ces_tool_basic5%{random_suffix}"
    execution_type = "SYNCHRONOUS"
    python_function {
        name = "example_function"
        python_code = "def example_function() -> int: return 0"
    }
}
`, context)
}

func testAccCESTool_cesToolPythonFunctionBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_ces_app" "my-app" {
    location     = "us"
    display_name = "tf-test-my-app%{random_suffix}"
    app_id       = "tf-test-app-id%{random_suffix}"
    time_zone_settings {   
        time_zone = "America/Los_Angeles"
    }
}
resource "google_ces_tool" "ces_tool_python_function_basic" {
    location       = "us"
    app            = google_ces_app.my-app.name
    tool_id        = "tf_test_ces_tool_basic5%{random_suffix}"
    execution_type = "SYNCHRONOUS"
    python_function {
        name = "example_function_updated"
        python_code = "def example_function_updated() -> int: return 0"
    }
}
`, context)
}
