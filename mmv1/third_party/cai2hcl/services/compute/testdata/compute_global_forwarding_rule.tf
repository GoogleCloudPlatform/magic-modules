resource "google_compute_global_forwarding_rule" "fr" {
  ip_protocol           = "TCP"
  ip_version            = "IPV4"
  load_balancing_scheme = "EXTERNAL"
  name                  = "fr"
  port_range            = "443"
  target                = "projects/myproj/global/targetSslProxies/tp"
}

resource "google_compute_global_forwarding_rule" "fr" {
  ip_address            = "projects/cn-fe-playground/global/addresses/ipaddr"
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  name                  = "fr"
  port_range            = "25"
  target                = "projects/cn-fe-playground/global/targetTcpProxies/tp"
}
