resource "google_compute_forwarding_rule" "test-1" {
  all_ports              = true
  allow_global_access    = false
  description            = "test description"
  ip_address             = "10.128.0.62"
  ip_protocol            = "TCP"
  is_mirroring_collector = false
  load_balancing_scheme  = "INTERNAL_MANAGED"
  name                   = "test-1"
  network_tier           = "PREMIUM"
  port_range             = "80-82"
  region                 = "us-central1"
  subnetwork             = "projects/myproj/regions/us-central1/subnetworks/default"
  target                 = "projects/myproj/regions/us-central1/targetHttpProxies/test1-target-proxy"
}

resource "google_compute_forwarding_rule" "test-2" {
  all_ports              = false
  allow_global_access    = false
  backend_service        = "projects/myproj/regions/us-central1/backendServices/test-bs-1"
  ip_address             = "projects/myproj/regions/us-central1/addresses/test-ip-1"
  ip_protocol            = "TCP"
  is_mirroring_collector = false
  load_balancing_scheme  = "EXTERNAL"
  name                   = "test-2"
  ports                  = ["80", "81"]
  region                 = "us-central1"
}
