resource "google_dns_managed_zone" "private-zone" {
  name = "private-zone"
  dns_name = "private.example.com."
  description = "Example private DNS zone"
  labels = {
    foo = "bar"
  }

  visibility = "private"

  private_visibility_config {
    networks {
      network_url =  "${google_compute_network.network-1.self_link}"
    }
    networks {
      network_url =  "${google_compute_network.network-2.self_link}"
    }
  }
}

resource "google_compute_network" "network-1" {
  name = "network-1"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network-2" {
  name = "network-2"
  auto_create_subnetworks = false
}
