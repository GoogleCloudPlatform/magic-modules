resource "google_compute_network" "default" {
  name = "my-network-${local.name_suffix}"
}

resource "google_compute_subnetwork" "default" {
  name          = "my-subnet-${local.name_suffix}"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = "${google_compute_network.default.self_link}"
}

resource "google_compute_address" "internal_with_subnet_and_address" {
  name         = "my-internal-address-${local.name_suffix}"
  subnetwork   = "${google_compute_subnetwork.default.self_link}"
  address_type = "INTERNAL"
  address      = "10.0.42.42"
  region       = "us-central1"
}
