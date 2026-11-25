package networkservices_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/networkservices"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkServicesHttpRoute_update(t *testing.T) {
	t.Parallel()

	httpRouteName := fmt.Sprintf("tf-test-http-route-%s", acctest.RandString(t, 10))
	backendServiceName := fmt.Sprintf("tf-test-http-service-%s", acctest.RandString(t, 10))
	healthCheckName := fmt.Sprintf("tf-test-http-healthcheck-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesHttpRouteDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesHttpRoute_basic(httpRouteName, backendServiceName, healthCheckName),
			},
			{
				ResourceName:            "google_network_services_http_route.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesHttpRoute_update(httpRouteName, backendServiceName, healthCheckName),
			},
			{
				ResourceName:            "google_network_services_http_route.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesHttpRoute_updateOmitOptional(httpRouteName, backendServiceName, healthCheckName),
			},
			{
				ResourceName:            "google_network_services_http_route.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkServicesHttpRoute_backendService(backendServiceName, healthCheckName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foo" {
  name                  = "%s"
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
  health_checks         = [google_compute_health_check.bar.id]
}

resource "google_compute_health_check" "bar" {
  name = "%s"

  http_health_check {
    port = 80
  }
}
`, backendServiceName, healthCheckName)
}

func testAccNetworkServicesHttpRoute_basic(httpRouteName, backendServiceName, healthCheckName string) string {
	return testAccNetworkServicesHttpRoute_backendService(backendServiceName, healthCheckName) + fmt.Sprintf(`
resource "google_network_services_http_route" "foobar" {
  name        = "%s"
  description = "my description"
  hostnames   = ["example"]
  rules {
    matches {
      query_parameters {
        query_parameter = "key"
        exact_match     = "value"
      }
      full_path_match = "example"
    }
    action {
      idle_timeout = "30s"
      destinations {
        service_name = google_compute_backend_service.foo.id
        weight       = 1
        request_header_modifier {
          set    = { "version" = "1" }
          add    = { "minor-version" = "1" }
          remove = ["arg"]
        }
        response_header_modifier {
          set    = { "version" = "1" }
          add    = { "minor-version" = "1" }
          remove = ["removearg"]
        }
      }
      request_mirror_policy {
        destination {
          service_name = "new"
          weight       = 1
        }
        mirror_percent = 25.5
      }
    }
  }
  rules {
    matches {
      full_path_match = "/direct"
    }
    action {
      direct_response {
        status      = 200
        string_body = "example-body"
      }
    }
  }
}
`, httpRouteName)
}

func testAccNetworkServicesHttpRoute_update(httpRouteName, backendServiceName, healthCheckName string) string {
	return testAccNetworkServicesHttpRoute_backendService(backendServiceName, healthCheckName) + fmt.Sprintf(`
resource "google_network_services_http_route" "foobar" {
  name        = "%s"
  description = "update description"
  labels = {
    foo = "bar"
  }
  hostnames = ["example"]
  rules {
    matches {
      query_parameters {
        query_parameter = "key"
        exact_match     = "value"
      }
      full_path_match = "example"
    }
    action {
      idle_timeout = "60s"
      destinations {
        service_name = google_compute_backend_service.foo.id
        weight       = 1
        request_header_modifier {
          set    = { "version" = "2" }
          add    = { "minor-version" = "2" }
          remove = ["arg", "other-arg"]
        }
        response_header_modifier {
          set    = { "version" = "2" }
          add    = { "minor-version" = "2" }
          remove = ["removearg", "other-removearg"]
        }
      }
      request_mirror_policy {
        destination {
          service_name = "new"
          weight       = 1
        }
        mirror_percent = 75.5
      }
    }
  }
  rules {
    matches {
      full_path_match = "/direct"
    }
    action {
      direct_response {
        status     = 201
        bytes_body = "SGVsbG8gd29ybGQ="
      }
    }
  }
}
`, httpRouteName)
}

func testAccNetworkServicesHttpRoute_updateOmitOptional(httpRouteName, backendServiceName, healthCheckName string) string {
	return testAccNetworkServicesHttpRoute_backendService(backendServiceName, healthCheckName) + fmt.Sprintf(`
resource "google_network_services_http_route" "foobar" {
  name        = "%s"
  description = "update description"
  labels = {
    foo = "bar"
  }
  hostnames = ["example"]
  rules {
    matches {
      query_parameters {
        query_parameter = "key"
        exact_match     = "value"
      }
      full_path_match = "example"
    }
    action {
      destinations {
        service_name = google_compute_backend_service.foo.id
        weight       = 1
      }
      request_mirror_policy {
        destination {
          service_name = "new"
          weight       = 1
        }
      }
    }
  }
  rules {
    matches {
      full_path_match = "/direct"
    }
    action {
      direct_response {
        status = 204
      }
    }
  }
}
`, httpRouteName)
}
