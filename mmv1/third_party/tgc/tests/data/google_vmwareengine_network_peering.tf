resource "google_vmwareengine_network_peering" "peering" {
  name                  = "network-peering-test"
  peer_network          = "projects/{{.Provider.Project}}/locations/global/vmwareEngineNetworks/network-peering-test-nw"
  peer_network_type     = "STANDARD"
  vmware_engine_network = "projects/{{.Provider.Project}}/locations/global/vmwareEngineNetworks/network-peering-test-ven"
}
