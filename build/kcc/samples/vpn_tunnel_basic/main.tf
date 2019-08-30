resource "google_compute_vpn_tunnel" "tunnel1" {
  name          = "tunnel1-${local.name_suffix}"
  peer_ip       = "15.0.0.120"
  shared_secret = "a secret message"

  target_vpn_gateway = "${google_compute_vpn_gateway.target_gateway.self_link}"

  depends_on = [
    "google_compute_forwarding_rule.fr_esp",
    "google_compute_forwarding_rule.fr_udp500",
    "google_compute_forwarding_rule.fr_udp4500",
  ]
}

resource "google_compute_vpn_gateway" "target_gateway" {
  name    = "vpn1-${local.name_suffix}"
  network = "${google_compute_network.network1.self_link}"
}

resource "google_compute_network" "network1" {
  name       = "network1-${local.name_suffix}"
}

resource "google_compute_address" "vpn_static_ip" {
  name   = "vpn-static-ip-${local.name_suffix}"
}

resource "google_compute_forwarding_rule" "fr_esp" {
  name        = "fr-esp-${local.name_suffix}"
  ip_protocol = "ESP"
  ip_address  = "${google_compute_address.vpn_static_ip.address}"
  target      = "${google_compute_vpn_gateway.target_gateway.self_link}"
}

resource "google_compute_forwarding_rule" "fr_udp500" {
  name        = "fr-udp500-${local.name_suffix}"
  ip_protocol = "UDP"
  port_range  = "500"
  ip_address  = "${google_compute_address.vpn_static_ip.address}"
  target      = "${google_compute_vpn_gateway.target_gateway.self_link}"
}

resource "google_compute_forwarding_rule" "fr_udp4500" {
  name        = "fr-udp4500-${local.name_suffix}"
  ip_protocol = "UDP"
  port_range  = "4500"
  ip_address  = "${google_compute_address.vpn_static_ip.address}"
  target      = "${google_compute_vpn_gateway.target_gateway.self_link}"
}

resource "google_compute_route" "route1" {
  name       = "route1-${local.name_suffix}"
  network    = "${google_compute_network.network1.name}"
  dest_range = "15.0.0.0/24"
  priority   = 1000

  next_hop_vpn_tunnel = "${google_compute_vpn_tunnel.tunnel1.self_link}"
}
