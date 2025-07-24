package networkservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkServicesWasmPlugin_wasmPluginLogConfigUpdate(t *testing.T) {
	acctest.SkipIfVcr(t) // Test requires a existing container image that contains the plugin code, published in an Artifact Registry repository.
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"test_project_id": envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesWasmPluginDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesWasmPlugin_wasmPluginBasicCreate(context),
			},
			{
				ResourceName:            "google_network_services_wasm_plugin.wasm_plugin",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "name", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesWasmPlugin_wasmPluginLogConfigUpdate(context),
			},
			{
				ResourceName:            "google_network_services_wasm_plugin.wasm_plugin",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "name", "terraform_labels"},
			},
		},
	})
}

func TestAccNetworkServicesWasmPlugin_wasmPluginVersionUpdate(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"test_project_id": envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesWasmPluginDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesWasmPlugin_wasmPluginVersionCreate(context),
			},
			{
				ResourceName:            "google_network_services_wasm_plugin.wasm_plugin",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "name", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesWasmPlugin_wasmPluginVersionUpdate(context),
			},
			{
				ResourceName:            "google_network_services_wasm_plugin.wasm_plugin",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "name", "terraform_labels"},
			},
		},
	})
}

func TestAccNetworkServicesWasmPlugin_wasmPluginConfigUpdate(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"test_project_id": envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesWasmPluginDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesWasmPlugin_wasmPluginBasicCreate(context),
			},
			{
				ResourceName:            "google_network_services_wasm_plugin.wasm_plugin",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "name", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesWasmPlugin_wasmPluginConfigDataUpdate(context),
			},
			{
				ResourceName:            "google_network_services_wasm_plugin.wasm_plugin",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "name", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesWasmPlugin_wasmPluginConfigUriUpdate(context),
			},
			{
				ResourceName:            "google_network_services_wasm_plugin.wasm_plugin",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "name", "terraform_labels"},
			},
		},
	})
}

func TestAccNetworkServicesWasmPlugin_wasmPluginLocation(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"test_project_id": envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesWasmPluginDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesWasmPlugin_wasmPluginLocationCreate(context),
			},
			{
				ResourceName:            "google_network_services_wasm_plugin.wasm_plugin",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "name", "terraform_labels"},
			},
			// We shoud test steps for regional plugins once available
		},
	})
}

func testAccNetworkServicesWasmPlugin_wasmPluginBasicCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_services_wasm_plugin" "wasm_plugin" {
  name        = "tf-test-my-wasm-plugin%{random_suffix}"
  description = "my wasm plugin"

  main_version_id = "v1"

  labels = {
    test_label =  "test_value"
  }
  log_config {
    enable =  true
    sample_rate = 1
    min_log_level =  "WARN"
  }

  versions {
    version_name = "v1"
    description = "v1 version of my wasm plugin"
    image_uri = "us-central1-docker.pkg.dev/%{test_project_id}/svextensionplugin/my-wasm-plugin:prod"

    labels = {
      test_label =  "test_value"
    }
  }
}
`, context)
}

func testAccNetworkServicesWasmPlugin_wasmPluginLogConfigUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_services_wasm_plugin" "wasm_plugin" {
  name        = "tf-test-my-wasm-plugin%{random_suffix}"
  description = "my wasm plugin"

  main_version_id = "v1"

  labels = {
    test_label2 =  "test_value2"
  }
  log_config {
    enable =  true
    sample_rate = 0.5
    min_log_level =  "ERROR"
  }

  versions {
    version_name = "v1"
    description = "v1 version of my wasm plugin"
    image_uri = "us-central1-docker.pkg.dev/%{test_project_id}/svextensionplugin/my-wasm-plugin:prod"

    labels = {
      test_label =  "test_value"
    }
  }
}
`, context)
}

func testAccNetworkServicesWasmPlugin_wasmPluginVersionCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_services_wasm_plugin" "wasm_plugin" {
  name        = "tf-test-my-wasm-plugin%{random_suffix}"
  description = "my wasm plugin"

  main_version_id = "v2"

  labels = {
    test_label =  "test_value"
  }
  log_config {
    enable =  true
    sample_rate = 1
    min_log_level =  "WARN"
  }

  versions {
    version_name = "v1"
    description = "v1 version of my wasm plugin"
    image_uri = "us-central1-docker.pkg.dev/%{test_project_id}/svextensionplugin/my-wasm-plugin:prod"

    labels = {
      test_label =  "test_value"
    }
  }
  versions {
    version_name = "v2"
    description = "v2 version of my wasm plugin"
    image_uri = "us-central1-docker.pkg.dev/%{test_project_id}/svextensionplugin/my-wasm-plugin:prod"

    labels = {
      test_label =  "test_value"
    }
  }
}
`, context)
}

func testAccNetworkServicesWasmPlugin_wasmPluginVersionUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_services_wasm_plugin" "wasm_plugin" {
  name        = "tf-test-my-wasm-plugin%{random_suffix}"
  description = "my wasm plugin"

  main_version_id = "v2"

  labels = {
    test_label =  "test_value"
  }
  log_config {
    enable =  true
    sample_rate = 1
    min_log_level =  "WARN"
  }

  versions {
    version_name = "v2"
    description = "v2 version of my wasm plugin"
    image_uri = "us-central1-docker.pkg.dev/%{test_project_id}/svextensionplugin/my-wasm-plugin:prod"

    labels = {
      test_label =  "test_value"
    }
  }
  versions {
    version_name = "v3"
    description = "v3 version of my wasm plugin"
    image_uri = "us-central1-docker.pkg.dev/%{test_project_id}/svextensionplugin/my-wasm-plugin:prod"

    labels = {
      test_label =  "test_value"
    }
  }
}
`, context)
}

func testAccNetworkServicesWasmPlugin_wasmPluginConfigDataUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_network_services_wasm_plugin" "wasm_plugin" {
  name        = "tf-test-my-wasm-plugin%{random_suffix}"
  description = "my wasm plugin"

  main_version_id = "v2"

  labels = {
    test_label =  "test_value"
  }
  log_config {
    enable =  true
    sample_rate = 1
    min_log_level =  "WARN"
  }

  versions {
    version_name = "v2"
    description = "v2 version of my wasm plugin"
    image_uri = "us-central1-docker.pkg.dev/%{test_project_id}/svextensionplugin/my-wasm-plugin:prod"
    plugin_config_data = base64encode("WasmPluginConfigDataTestValue%{random_suffix}")

    labels = {
      test_label =  "test_value"
    }
  }
}
`, context)
}

func testAccNetworkServicesWasmPlugin_wasmPluginConfigUriUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_network_services_wasm_plugin" "wasm_plugin" {
  name        = "tf-test-my-wasm-plugin%{random_suffix}"
  description = "my wasm plugin"

  main_version_id = "v3"

  labels = {
    test_label =  "test_value"
  }
  log_config {
    enable =  true
    sample_rate = 1
    min_log_level =  "WARN"
  }

  versions {
    version_name = "v3"
    description = "v3 version of my wasm plugin"
    image_uri = "us-central1-docker.pkg.dev/%{test_project_id}/svextensionplugin/my-wasm-plugin:prod"
    plugin_config_uri = "us-central1-docker.pkg.dev/%{test_project_id}/svextensionplugin/wasm-plugin-config-secret:prod"

    labels = {
      test_label =  "test_value"
    }
  }
}
`, context)
}

func testAccNetworkServicesWasmPlugin_wasmPluginLocationCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_services_wasm_plugin" "wasm_plugin" {
  name        = "tf-test-my-wasm-plugin%{random_suffix}"
  description = "my wasm plugin"
  location 		= "global"

  main_version_id = "v1"

  labels = {
    test_label =  "test_value"
  }
  log_config {
    enable =  true
    sample_rate = 1
    min_log_level =  "WARN"
  }

  versions {
    version_name = "v1"
    description = "v1 version of my wasm plugin"
    image_uri = "us-central1-docker.pkg.dev/%{test_project_id}/svextensionplugin/my-wasm-plugin:prod"

    labels = {
      test_label =  "test_value"
    }
  }
}
`, context)
}
