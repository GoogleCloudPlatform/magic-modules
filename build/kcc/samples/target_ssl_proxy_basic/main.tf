resource "google_compute_target_ssl_proxy" "default" {
  name             = "test-proxy-${local.name_suffix}"
  backend_service  = "${google_compute_backend_service.default.self_link}"
  ssl_certificates = ["${google_compute_ssl_certificate.default.self_link}"]
}

resource "google_compute_ssl_certificate" "default" {
  name        = "default-cert-${local.name_suffix}"
  private_key = "${file("../static/ssl_cert/test.key")}"
  certificate = "${file("../static/ssl_cert/test.crt")}"
}

resource "google_compute_backend_service" "default" {
  name          = "backend-service-${local.name_suffix}"
  protocol      = "SSL"
  health_checks = ["${google_compute_health_check.default.self_link}"]
}

resource "google_compute_health_check" "default" {
  name               = "health-check-${local.name_suffix}"
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "443"
  }
}
