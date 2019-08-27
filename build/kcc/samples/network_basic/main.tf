resource "google_compute_network" "vpc_network" {
  name = "vpc-network-${local.name_suffix}"
}
