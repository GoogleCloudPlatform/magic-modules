resource "google_compute_network_endpoint_group" "neg" {
  name         = "my-lb-neg"
  network      = "${google_compute_network.default.self_link}"
  subnetwork   = "${google_compute_subnetwork.default.self_link}"
  default_port = "90"
  zone         = "us-central1-a"
}

resource "google_compute_network" "default" {
  name = "neg-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "neg-subnetwork"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = "${google_compute_network.default.self_link}"
}
