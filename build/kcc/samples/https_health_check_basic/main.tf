resource "google_compute_https_health_check" "default" {
  name         = "authentication-health-check-${local.name_suffix}"
  request_path = "/health_check"

  timeout_sec        = 1
  check_interval_sec = 1
}
