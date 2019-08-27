resource "google_compute_network_endpoint_group" "neg" {
  name         = "my-lb-neg-${local.name_suffix}"
  network      = "${google_compute_network.default.self_link}"
  subnetwork   = "${google_compute_subnetwork.default.self_link}"
  default_port = "90"
  zone         = "us-central1-a"
}

resource "google_compute_network" "default" {
  name = "neg-network-${local.name_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "neg-subnetwork-${local.name_suffix}"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = "${google_compute_network.default.self_link}"
}
