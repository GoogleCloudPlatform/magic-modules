package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeTargetHttpProxy_httpFiltersUpdate(t *testing.T) {
	t.Parallel()

	target := fmt.Sprintf("thttp-test-%s", acctest.RandString(t, 10))
	backend := fmt.Sprintf("thttp-test-%s", acctest.RandString(t, 10))
	hc := fmt.Sprintf("thttp-test-%s", acctest.RandString(t, 10))
	urlmap1 := fmt.Sprintf("thttp-test-%s", acctest.RandString(t, 10))
	urlmap2 := fmt.Sprintf("thttp-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeTargetHttpProxyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeTargetHttpProxy_httpFiltersBasic(target, backend, hc, urlmap1, urlmap2),
			},
			{
				ResourceName:      "google_compute_target_http_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeTargetHttpProxy_httpFiltersUpdate(target, backend, hc, urlmap1, urlmap2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_target_http_proxy.foobar", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_compute_target_http_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeTargetHttpProxy_httpFiltersBasic(target, backend, hc, urlmap1, urlmap2 string) string {
	return fmt.Sprintf(`
resource "google_compute_target_http_proxy" "foobar" {
  description  = "Resource created for Terraform acceptance testing"
  name         = "%s"
  url_map      = google_compute_url_map.foobar1.self_link
  http_filters = []
}

resource "google_compute_global_forwarding_rule" "default" {
  name                  = "%s-fr"
  target                = google_compute_target_http_proxy.foobar.self_link
  port_range            = "80"
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
  ip_address            = "0.0.0.0"
}

resource "google_compute_backend_service" "foobar" {
  name                  = "%s"
  health_checks         = [google_compute_health_check.zero.self_link]
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
}

resource "google_compute_health_check" "zero" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1

  http_health_check {
    port = 80
  }
}

resource "google_compute_url_map" "foobar1" {
  name            = "%s"
  default_service = google_compute_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.foobar.self_link
  }
}

resource "google_compute_url_map" "foobar2" {
  name            = "%s"
  default_service = google_compute_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.foobar.self_link
  }
}
`, target, target, backend, hc, urlmap1, urlmap2)
}

func testAccComputeTargetHttpProxy_httpFiltersUpdate(target, backend, hc, urlmap1, urlmap2 string) string {
	return fmt.Sprintf(`
resource "google_compute_target_http_proxy" "foobar" {
  description  = "Resource created for Terraform acceptance testing"
  name         = "%s"
  url_map      = google_compute_url_map.foobar2.self_link
  http_filters = []
}

resource "google_compute_global_forwarding_rule" "default" {
  name                  = "%s-fr"
  target                = google_compute_target_http_proxy.foobar.self_link
  port_range            = "80"
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
  ip_address            = "0.0.0.0"
}

resource "google_compute_backend_service" "foobar" {
  name                  = "%s"
  health_checks         = [google_compute_health_check.zero.self_link]
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
}

resource "google_compute_health_check" "zero" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1

  http_health_check {
    port = 80
  }
}

resource "google_compute_url_map" "foobar1" {
  name            = "%s"
  default_service = google_compute_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.foobar.self_link
  }
}

resource "google_compute_url_map" "foobar2" {
  name            = "%s"
  default_service = google_compute_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.foobar.self_link
  }
}
`, target, target, backend, hc, urlmap1, urlmap2)
}
