package networkservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkServicesLbEdgeExtension_networkServicesLbEdgeExtensionBasicUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesLbEdgeExtensionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesWasmPlugin_artifactRegistryRepositorySetup(context),
				Check: resource.ComposeTestCheckFunc(
					// Upload the compiled plugin code to Artifact Registry
					testAccCheckNetworkServicesWasmPlugin_uploadCompiledCode(
						t,
						"google_artifact_registry_repository.test_repository",
						"my-wasm-plugin",
						"v1",
						"test-fixtures/compiled-package/plugin.wasm",
						"plugin.wasm",
					),
					testAccCheckNetworkServicesWasmPlugin_uploadCompiledCode(
						t,
						"google_artifact_registry_repository.test_repository",
						"my-wasm-plugin-2",
						"v1",
						"test-fixtures/compiled-package/plugin.wasm",
						"plugin.wasm",
					),
				),
			},
			{
				Config: testAccNetworkServicesLbEdgeExtension_networkServicesLbEdgeExtensionBasicCreate(context),
			},
			{
				ResourceName:            "google_network_services_lb_edge_extension.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "port_range", "target", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesLbEdgeExtension_networkServicesLbEdgeExtensionBasicUpdate(context),
			},
			{
				ResourceName:            "google_network_services_lb_edge_extension.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "port_range", "target", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkServicesLbEdgeExtension_networkServicesLbEdgeExtensionBasicCreate(context map[string]interface{}) string {
	return fmt.Sprint(testAccNetworkServicesWasmPlugin_artifactRegistryRepositorySetup(context), acctest.Nprintf(`
# forwarding rule
resource "google_compute_global_forwarding_rule" "default" {
  name                  = "tf-test-elb-forwarding-rule%{random_suffix}"
  target                = google_compute_target_http_proxy.default.id
  port_range            = "80"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  network_tier          = "PREMIUM"
}

resource "google_compute_target_http_proxy" "default" {
  name        = "tf-test-elb-target-http-proxy%{random_suffix}"
  description = "a description"
  url_map     = google_compute_url_map.default.id
}

resource "google_compute_url_map" "default" {
  name            = "tf-test-elb-url-map%{random_suffix}"
  description     = "a description"
  default_service = google_compute_backend_service.default.id

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name            = "allpaths"
    default_service = google_compute_backend_service.default.id

    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.default.id
    }
  }
}

resource "google_compute_backend_service" "default" {
  name                  = "tf-test-elb-backend-subnet%{random_suffix}"
  port_name             = "http"
  protocol              = "HTTP"
  timeout_sec           = 10
  load_balancing_scheme = "EXTERNAL_MANAGED"
}

resource "google_network_services_lb_edge_extension" "default" {
  name        = "tf-test-elb-edge-ext%{random_suffix}"
  description = "my edge extension"
  location    = "global"

  load_balancing_scheme = "EXTERNAL_MANAGED"
  forwarding_rules      = [google_compute_global_forwarding_rule.default.self_link]

  extension_chains {
    name = "chain1"

    match_condition {
      cel_expression = "request.host == 'example.com'"
    }

    extensions {
      name      = "ext11"
      service   = google_network_services_wasm_plugin.wasm-plugin.id
      fail_open = false
      supported_events = ["REQUEST_HEADERS"]
      forward_headers  = ["custom-header"]
    }
  }

  labels = {
    foo = "bar"
  }
}

resource "google_network_services_wasm_plugin" "wasm-plugin" {
  name        = "tf-test-elb-wasm-plugin%{random_suffix}"
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
    image_uri = "projects/%{project}/locations/us-central1/repositories/tf-test-repository-standard%{random_suffix}/genericArtifacts/my-wasm-plugin:v1"

    labels = {
      test_label =  "test_value"
    }
  }
}

resource "google_network_services_wasm_plugin" "wasm-plugin-2" {
  name        = "tf-test-elb-wasm-plugin-2%{random_suffix}"
  description = "my wasm plugin 2"

  main_version_id = "v1"

  labels = {
    test_label =  "test_value_2"
  }
  log_config {
    enable =  true
    sample_rate = 1
    min_log_level =  "WARN"
  }

  versions {
    version_name = "v1"
    description = "v1 version of my wasm plugin 2"
    image_uri = "projects/%{project}/locations/us-central1/repositories/tf-test-repository-standard%{random_suffix}/genericArtifacts/my-wasm-plugin-2:v1"

    labels = {
      test_label =  "test_value_2"
    }
  }
}
`, context))
}

func testAccNetworkServicesLbEdgeExtension_networkServicesLbEdgeExtensionBasicUpdate(context map[string]interface{}) string {
	return fmt.Sprint(testAccNetworkServicesWasmPlugin_artifactRegistryRepositorySetup(context), acctest.Nprintf(`
# forwarding rule
resource "google_compute_global_forwarding_rule" "default" {
  name                  = "tf-test-elb-forwarding-rule%{random_suffix}"
  target                = google_compute_target_http_proxy.default.id
  port_range            = "80"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  network_tier          = "PREMIUM"
}

resource "google_compute_target_http_proxy" "default" {
  name        = "tf-test-elb-target-http-proxy%{random_suffix}"
  description = "a description"
  url_map     = google_compute_url_map.default.id
}

resource "google_compute_url_map" "default" {
  name            = "tf-test-elb-url-map%{random_suffix}"
  description     = "a description"
  default_service = google_compute_backend_service.default.id

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name            = "allpaths"
    default_service = google_compute_backend_service.default.id

    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.default.id
    }
  }
}

resource "google_compute_backend_service" "default" {
  name                  = "tf-test-elb-backend-subnet%{random_suffix}"
  port_name             = "http"
  protocol              = "HTTP"
  timeout_sec           = 10
  load_balancing_scheme = "EXTERNAL_MANAGED"
}

resource "google_network_services_lb_edge_extension" "default" {
  name        = "tf-test-elb-edge-ext%{random_suffix}"
  description = "my edge extension"
  location    = "global"

  load_balancing_scheme = "EXTERNAL_MANAGED"
  forwarding_rules      = [google_compute_global_forwarding_rule.default.self_link]

  extension_chains {
    name = "chain1"

    match_condition {
      cel_expression = "request.host == 'example.com'"
    }

    extensions {
      name      = "ext11"
      service   = google_network_services_wasm_plugin.wasm-plugin.id
      fail_open = false
      supported_events = ["REQUEST_HEADERS"]
      forward_headers  = ["custom-header"]
    }
  }

  extension_chains {
    name = "chain2"

    match_condition {
      cel_expression = "request.host == 'example.com'"
    }

    extensions {
      name      = "ext12"
      service   = google_network_services_wasm_plugin.wasm-plugin-2.id
      fail_open = true
      supported_events = ["REQUEST_HEADERS"]
      forward_headers  = ["custom-header"]
    }
  }

  labels = {
    foo = "bar"
  }
}

resource "google_network_services_wasm_plugin" "wasm-plugin" {
  name        = "tf-test-elb-wasm-plugin%{random_suffix}"
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
    image_uri = "projects/%{project}/locations/us-central1/repositories/tf-test-repository-standard%{random_suffix}/genericArtifacts/my-wasm-plugin:v1"

    labels = {
      test_label =  "test_value"
    }
  }
}

resource "google_network_services_wasm_plugin" "wasm-plugin-2" {
  name        = "tf-test-elb-wasm-plugin-2%{random_suffix}"
  description = "my wasm plugin 2"

  main_version_id = "v1"

  labels = {
    test_label =  "test_value_2"
  }
  log_config {
    enable =  true
    sample_rate = 1
    min_log_level =  "WARN"
  }

  versions {
    version_name = "v1"
    description = "v1 version of my wasm plugin 2"
    image_uri = "projects/%{project}/locations/us-central1/repositories/tf-test-repository-standard%{random_suffix}/genericArtifacts/my-wasm-plugin-2:v1"

    labels = {
      test_label =  "test_value_2"
    }
  }
}
`, context))
}
