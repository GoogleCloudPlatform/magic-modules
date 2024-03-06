resource "google_compute_backend_service" "bs-1" {
  affinity_cookie_ttl_sec = 0

  backend {
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 1
    group           = "projects/cn-fe-playground/zones/us-central1-a/instanceGroups/alicjab-instance-group-1"
    max_utilization = 0.8
  }

  cdn_policy {
    cache_key_policy {
      include_host         = true
      include_protocol     = true
      include_query_string = true
    }

    cache_mode                   = "CACHE_ALL_STATIC"
    client_ttl                   = 3600
    default_ttl                  = 3600
    max_ttl                      = 86400
    negative_caching             = false
    serve_while_stale            = 0
    signed_url_cache_max_age_sec = 0
  }

  connection_draining_timeout_sec = 300
  creation_timestamp              = "2023-06-09T04:35:59.474-07:00"
  description                     = "bs-1 description"
  enable_cdn                      = false
  fingerprint                     = "m1r6cXyt2rI="
  health_checks                   = ["projects/cn-fe-playground/global/healthChecks/alicjab-health-check"]
  load_balancing_scheme           = "EXTERNAL_MANAGED"
  locality_lb_policy              = "ROUND_ROBIN"

  log_config {
    enable = false
  }

  name             = "bs-1"
  port_name        = "http"
  protocol         = "HTTP"
  security_policy  = "projects/cn-fe-playground/global/securityPolicies/example-policy"
  session_affinity = "NONE"
  timeout_sec      = 30
}
