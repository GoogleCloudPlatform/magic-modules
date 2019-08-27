resource "google_compute_firewall" "default" {
  name    = "test-firewall-${local.name_suffix}"
  network = "${google_compute_network.default.name}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["80", "8080", "1000-2000"]
  }

  source_tags = ["web"]
}

resource "google_compute_network" "default" {
  name = "test-network-${local.name_suffix}"
}
