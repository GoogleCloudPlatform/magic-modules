resource "google_compute_health_check" "health-check-1" {
  check_interval_sec = 5
  healthy_threshold  = 2

  log_config {
    enable = false
  }

  name = "health-check-1"

  tcp_health_check {
    port         = 80
    proxy_header = "NONE"
  }

  timeout_sec         = 5
  type                = "TCP"
  unhealthy_threshold = 2
}
