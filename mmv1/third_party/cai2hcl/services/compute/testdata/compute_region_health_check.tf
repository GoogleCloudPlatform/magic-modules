resource "google_compute_region_health_check" "hc" {
  check_interval_sec = 5
  healthy_threshold  = 2

  log_config {
    enable = false
  }

  name   = "hc"
  region = "us-central1"

  tcp_health_check {
    port         = 80
    proxy_header = "NONE"
  }

  timeout_sec         = 5
  type                = "TCP"
  unhealthy_threshold = 2
}

resource "google_compute_region_health_check" "hc" {
  check_interval_sec = 5
  description        = "descr"
  healthy_threshold  = 2

  log_config {
    enable = false
  }

  name   = "hc"
  region = "us-central1"

  tcp_health_check {
    port         = 8
    proxy_header = "PROXY_V1"
    request      = "a"
    response     = "b"
  }

  timeout_sec         = 5
  type                = "TCP"
  unhealthy_threshold = 2
}
