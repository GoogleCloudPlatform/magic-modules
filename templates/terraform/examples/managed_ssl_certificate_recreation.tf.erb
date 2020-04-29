// This example allows the list of managed domains to be modified and will
// recreate the ssl certificate and update the target https proxy correctly

resource "google_compute_target_https_proxy" "default" {
  provider = google-beta
  name             = "test-proxy"
  url_map          = google_compute_url_map.default.self_link
  ssl_certificates = [google_compute_managed_ssl_certificate.cert.self_link]
}

locals {
  managed_domains = list("test.example.com")
}

resource "random_id" "certificate" {
  byte_length = 4
  prefix      = "issue6147-cert-"

  keepers = {
    domains = join(",", local.managed_domains)
  }
}

resource "google_compute_managed_ssl_certificate" "cert" {
  provider = google-beta
  name     = random_id.certificate.hex

  lifecycle {
    create_before_destroy = true
  }

  managed {
    domains = local.managed_domains
  }
}

resource "google_compute_url_map" "default" {
  provider = google-beta
  name            = "url-map"
  description     = "a description"
  default_service = google_compute_backend_service.default.self_link
  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }
  path_matcher {
    name            = "allpaths"
    default_service = google_compute_backend_service.default.self_link
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.default.self_link
    }
  }
}

resource "google_compute_backend_service" "default" {
  provider = google-beta
  name          = "backend-service"
  port_name     = "http"
  protocol      = "HTTP"
  timeout_sec   = 10
  health_checks = [google_compute_http_health_check.default.self_link]
}

resource "google_compute_http_health_check" "default" {
  provider = google-beta
  name               = "http-health-check"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
