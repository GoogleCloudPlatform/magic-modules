package apihub_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/apihub"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccApihubPluginInstance_apihubPluginInstanceSystemPlugin(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApihubPluginInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApihubPluginInstance_systemPlugin(context),
			},
			{
				ResourceName:            "google_apihub_plugin_instance.system_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable", "location", "plugin", "plugin_instance_id"},
			},
			{
				Config: testAccApihubPluginInstance_systemPluginUpdate(context),
			},
			{
				ResourceName:            "google_apihub_plugin_instance.system_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable", "location", "plugin", "plugin_instance_id"},
			},
		},
	})
}

func TestAccApihubPluginInstance_apihubPluginInstanceAdditionalConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApihubPluginInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApihubPluginInstance_additionalConfig(context),
			},
			{
				ResourceName:            "google_apihub_plugin_instance.custom_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"additional_config", "disable", "location", "plugin", "plugin_instance_id"},
			},
			{
				Config: testAccApihubPluginInstance_additionalConfigUpdate(context),
			},
			{
				ResourceName:            "google_apihub_plugin_instance.custom_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"additional_config", "disable", "location", "plugin", "plugin_instance_id"},
			},
		},
	})
}

func testAccApihubPluginInstance_systemPlugin(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "plugin_sa" {
  account_id   = "tf-test-pi-sa-%{random_suffix}"
  display_name = "Test Service Account for Plugin Instance"
}

resource "google_apihub_plugin_instance" "system_instance" {
  location           = "us-central1"
  plugin             = "system-edge-private-cloud"
  plugin_instance_id = "tf-test-pi-sys-%{random_suffix}"
  display_name       = "Test System Plugin Instance %{random_suffix}"
  source_project_id  = "org-%{random_suffix}"
  disable            = false

  actions {
    action_id       = "sync-metadata"
    service_account = google_service_account.plugin_sa.email
  }

  actions {
    action_id       = "sync-runtime-data"
    service_account = google_service_account.plugin_sa.email
  }

  source_environments_config {
    key                    = "prod"
    source_environment     = "prod"
    source_environment_uri = "https://example.com/prod"
    create_time            = "2024-01-01T00:00:00Z"
    update_time            = "2024-01-02T00:00:00Z"
  }
}
`, context)
}

func testAccApihubPluginInstance_systemPluginUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "plugin_sa" {
  account_id   = "tf-test-pi-sa-%{random_suffix}"
  display_name = "Test Service Account for Plugin Instance"
}

resource "google_apihub_plugin_instance" "system_instance" {
  location           = "us-central1"
  plugin             = "system-edge-private-cloud"
  plugin_instance_id = "tf-test-pi-sys-%{random_suffix}"
  display_name       = "Test System Plugin Instance %{random_suffix}"
  source_project_id  = "org-%{random_suffix}"
  disable            = false

  actions {
    action_id       = "sync-metadata"
    service_account = google_service_account.plugin_sa.email
  }

  actions {
    action_id       = "sync-runtime-data"
    service_account = google_service_account.plugin_sa.email
  }

  source_environments_config {
    key                = "staging"
    source_environment = "staging"
  }
}
`, context)
}

func testAccApihubPluginInstance_additionalConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_apihub_plugin" "custom_plugin" {
  location        = "us-central1"
  plugin_id       = "tf-test-plugin-%{random_suffix}"
  display_name    = "Test Custom Plugin %{random_suffix}"
  description     = "Custom plugin for testing additional_config"
  plugin_category = "API_GATEWAY"

  hosting_service {
    service_uri = "https://apihub-phs-stub-${data.google_project.project.number}.us-central1.run.app"
  }

  actions_config {
    id           = "custom-action"
    display_name = "Custom Action"
    description  = "Custom action for testing"
    trigger_mode = "API_HUB_ON_DEMAND_TRIGGER"
  }

  config_template {
    auth_config_template {
      supported_auth_types = ["NO_AUTH"]
    }
    additional_config_template {
      id               = "string-val"
      description      = "String config variable"
      value_type       = "STRING"
      required         = false
      validation_regex = "^[a-zA-Z0-9]{3,20}$"
    }
    additional_config_template {
      id          = "int-val"
      description = "Integer config variable"
      value_type  = "INT"
      required    = false
    }
    additional_config_template {
      id          = "bool-val"
      description = "Boolean config variable"
      value_type  = "BOOL"
      required    = false
    }
    additional_config_template {
      id          = "secret-val"
      description = "Secret config variable"
      value_type  = "SECRET"
      required    = false
    }
    additional_config_template {
      id          = "enum-val"
      description = "Enum config variable"
      value_type  = "ENUM"
      required    = false
      enum_options {
        id           = "Option1"
        display_name = "Option1"
        description  = "Description for Option1"
      }
      enum_options {
        id           = "Option2"
        display_name = "Option2"
      }
    }
    additional_config_template {
      id          = "multi-select-val"
      description = "Multi select config variable"
      value_type  = "MULTI_SELECT"
      required    = false
      multi_select_options {
        id           = "Select1"
        display_name = "Select1"
        description  = "Description for Select1"
      }
      multi_select_options {
        id           = "Select2"
        display_name = "Select2"
      }
    }
    additional_config_template {
      id          = "multi-string-val"
      description = "Multi string config variable"
      value_type  = "MULTI_STRING"
      required    = false
    }
    additional_config_template {
      id          = "multi-int-val"
      description = "Multi int config variable"
      value_type  = "INT"
      required    = false
    }
  }
}

resource "google_apihub_plugin_instance" "custom_instance" {
  location           = "us-central1"
  plugin             = google_apihub_plugin.custom_plugin.plugin_id
  plugin_instance_id = "tf-test-pi-cfg-%{random_suffix}"
  display_name       = "Test Custom Plugin Instance %{random_suffix}"
  disable            = false

  auth_config {
    auth_type = "NO_AUTH"
  }

  actions {
    action_id = "custom-action"
  }

  additional_config {
    key          = "string-val"
    string_value = "testval123"
  }

  additional_config {
    key       = "int-val"
    int_value = 42
  }

  additional_config {
    key        = "bool-val"
    bool_value = true
  }

  additional_config {
    key = "secret-val"
    secret_value {
      secret_version = "projects/%{project}/secrets/test-secret/versions/1"
    }
  }

  additional_config {
    key = "enum-val"
    enum_value {
      id           = "Option1"
      display_name = "Option1"
      description  = "Description for Option1"
    }
  }

  additional_config {
    key = "multi-select-val"
    multi_select_values {
      values {
        id           = "Select1"
        display_name = "Select1"
        description  = "Description for Select1"
      }
    }
  }

  additional_config {
    key = "multi-string-val"
    multi_string_values {
      values = ["val1", "val2"]
    }
  }

  additional_config {
    key = "multi-int-val"
    multi_int_values {
      values = [10, 20]
    }
  }
}
`, context)
}

func testAccApihubPluginInstance_additionalConfigUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_apihub_plugin" "custom_plugin" {
  location        = "us-central1"
  plugin_id       = "tf-test-plugin-%{random_suffix}"
  display_name    = "Test Custom Plugin %{random_suffix}"
  description     = "Custom plugin for testing additional_config"
  plugin_category = "API_GATEWAY"

  hosting_service {
    service_uri = "https://apihub-phs-stub-${data.google_project.project.number}.us-central1.run.app"
  }

  actions_config {
    id           = "custom-action"
    display_name = "Custom Action"
    description  = "Custom action for testing"
    trigger_mode = "API_HUB_ON_DEMAND_TRIGGER"
  }

  config_template {
    auth_config_template {
      supported_auth_types = ["NO_AUTH"]
    }
    additional_config_template {
      id               = "string-val"
      description      = "String config variable"
      value_type       = "STRING"
      required         = false
      validation_regex = "^[a-zA-Z0-9]{3,20}$"
    }
    additional_config_template {
      id          = "int-val"
      description = "Integer config variable"
      value_type  = "INT"
      required    = false
    }
    additional_config_template {
      id          = "bool-val"
      description = "Boolean config variable"
      value_type  = "BOOL"
      required    = false
    }
    additional_config_template {
      id          = "secret-val"
      description = "Secret config variable"
      value_type  = "SECRET"
      required    = false
    }
    additional_config_template {
      id          = "enum-val"
      description = "Enum config variable"
      value_type  = "ENUM"
      required    = false
      enum_options {
        id           = "Option1"
        display_name = "Option1"
        description  = "Description for Option1"
      }
      enum_options {
        id           = "Option2"
        display_name = "Option2"
      }
    }
    additional_config_template {
      id          = "multi-select-val"
      description = "Multi select config variable"
      value_type  = "MULTI_SELECT"
      required    = false
      multi_select_options {
        id           = "Select1"
        display_name = "Select1"
        description  = "Description for Select1"
      }
      multi_select_options {
        id           = "Select2"
        display_name = "Select2"
      }
    }
    additional_config_template {
      id          = "multi-string-val"
      description = "Multi string config variable"
      value_type  = "MULTI_STRING"
      required    = false
    }
    additional_config_template {
      id          = "multi-int-val"
      description = "Multi int config variable"
      value_type  = "INT"
      required    = false
    }
  }
}

resource "google_apihub_plugin_instance" "custom_instance" {
  location           = "us-central1"
  plugin             = google_apihub_plugin.custom_plugin.plugin_id
  plugin_instance_id = "tf-test-pi-cfg-%{random_suffix}"
  display_name       = "Test Custom Plugin Instance Updated %{random_suffix}"
  disable            = false

  auth_config {
    auth_type = "NO_AUTH"
  }

  actions {
    action_id = "custom-action"
  }

  additional_config {
    key          = "string-val"
    string_value = "updatedval456"
  }

  additional_config {
    key       = "int-val"
    int_value = 99
  }

  additional_config {
    key        = "bool-val"
    bool_value = false
  }

  additional_config {
    key = "secret-val"
    secret_value {
      secret_version = "projects/%{project}/secrets/test-secret/versions/2"
    }
  }

  additional_config {
    key = "enum-val"
    enum_value {
      id           = "Option2"
      display_name = "Option2"
    }
  }

  additional_config {
    key = "multi-select-val"
    multi_select_values {
      values {
        id           = "Select2"
        display_name = "Select2"
      }
    }
  }

  additional_config {
    key = "multi-string-val"
    multi_string_values {
      values = ["val3"]
    }
  }

  additional_config {
    key = "multi-int-val"
    multi_int_values {
      values = [30]
    }
  }
}
`, context)
}

func testAccCheckApihubPluginInstanceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_apihub_plugin_instance" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)
			url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(apihub.Product, config)+"projects/{{project}}/locations/{{location}}/plugins/{{plugin}}/instances/{{plugin_instance_id}}")
			if err != nil {
				return err
			}

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ApihubPluginInstance still exists at %s", url)
			}
		}

		return nil
	}
}
