resource "google_compute_target_ssl_proxy" "default" {
  name             = "test-proxy"
  backend_service  = "${google_compute_backend_service.default.self_link}"
  ssl_certificates = ["${google_compute_ssl_certificate.default.self_link}"]
}

resource "google_compute_ssl_certificate" "default" {
  name        = "default-cert"
  private_key = "${file("path/to/private.key")}"
  certificate = "${file("path/to/certificate.crt")}"
}

resource "google_compute_backend_service" "default" {
  name          = "backend-service"
  protocol      = "SSL"
  health_checks = ["${google_compute_health_check.default.self_link}"]
}

resource "google_compute_health_check" "default" {
  name               = "health-check"
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "443"
  }
}
