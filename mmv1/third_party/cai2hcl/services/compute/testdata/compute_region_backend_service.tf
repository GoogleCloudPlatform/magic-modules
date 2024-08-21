resource "google_compute_region_backend_service" "bs-1" {
  backend {
    balancing_mode = "CONNECTION"
    failover       = false
    group          = "projects/myproj/zones/us-central1-a/instanceGroups/ig-1"
  }

  connection_draining_timeout_sec = 30
  description                     = "bs-1 description"
  health_checks                   = ["projects/myproj/global/healthChecks/hc-1"]
  load_balancing_scheme           = "INTERNAL"

  log_config {
    enable      = true
    sample_rate = 0.2
  }

  name             = "bs-1"
  network          = "projects/myproj/global/networks/default"
  protocol         = "TCP"
  region           = "us-central1"
  session_affinity = "NONE"
}

resource "google_compute_region_backend_service" "bs-2" {
  backend {
    balancing_mode  = "CONNECTION"
    capacity_scaler = 0.1
    group           = "projects/myproj/zones/us-central1-c/networkEndpointGroups/neg-1"
    max_connections = 2
  }

  circuit_breakers {
    max_retries = 1
  }

  connection_draining_timeout_sec = 300
  health_checks                   = ["projects/myproj/regions/us-central1/healthChecks/hc-1"]
  load_balancing_scheme           = "EXTERNAL_MANAGED"
  locality_lb_policy              = "RING_HASH"

  log_config {
    enable = false
  }

  name = "bs-2"

  outlier_detection {
    base_ejection_time {
      nanos   = 0
      seconds = 30
    }

    consecutive_errors                    = 5
    consecutive_gateway_failure           = 3
    enforcing_consecutive_errors          = 0
    enforcing_consecutive_gateway_failure = 100
    enforcing_success_rate                = 100

    interval {
      nanos   = 0
      seconds = 1
    }

    max_ejection_percent        = 50
    success_rate_minimum_hosts  = 5
    success_rate_request_volume = 100
    success_rate_stdev_factor   = 1900
  }

  protocol         = "TCP"
  region           = "us-central1"
  session_affinity = "CLIENT_IP"
  timeout_sec      = 30
}
