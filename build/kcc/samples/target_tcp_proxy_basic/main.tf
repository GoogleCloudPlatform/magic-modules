resource "google_compute_target_tcp_proxy" "default" {
  name            = "test-proxy-${local.name_suffix}"
  backend_service = "${google_compute_backend_service.default.self_link}"
}

resource "google_compute_backend_service" "default" {
  name          = "backend-service-${local.name_suffix}"
  protocol      = "TCP"
  timeout_sec   = 10

  health_checks = ["${google_compute_health_check.default.self_link}"]
}

resource "google_compute_health_check" "default" {
  name               = "health-check-${local.name_suffix}"
  timeout_sec        = 1
  check_interval_sec = 1

  tcp_health_check {
    port = "443"
  }
}
