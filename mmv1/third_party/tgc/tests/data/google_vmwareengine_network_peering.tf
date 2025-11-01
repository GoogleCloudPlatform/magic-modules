resource "google_compute_network" "peered_network" {
  name                    = "network-peering-test-nw"
  auto_create_subnetworks = false
}

resource "google_vmwareengine_network" "vmware_network" {
  name        = "network-peering-test-ven"
  location    = "global"
  type        = "STANDARD"
  description = "VMware Engine network for peering test"
}

resource "google_vmwareengine_network_peering" "peering" {
  name                  = "network-peering-test"
  peer_network          = google_compute_network.peered_network.id
  peer_network_type     = "STANDARD"
  vmware_engine_network = google_vmwareengine_network.vmware_network.id
}
