resource "google_compute_region_backend_service" "default" {
  name                            = "region-backend-service"
  region                          = "us-central1"
  health_checks                   = ["${google_compute_health_check.default.self_link}"]
  connection_draining_timeout_sec = 10
  session_affinity                = "CLIENT_IP"
}

resource "google_compute_health_check" "default" {
  name               = "health-check"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = "80"
  }
}
