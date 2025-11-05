resource "google_vmwareengine_network_policy" "gg_asset_17806_51c7" {
  name                  = "gg-asset-17806-51c7"
  location              = "us-central1"
  project               = "{{.Provider.project}}"
  edge_services_cidr    = "192.168.30.0/26"
  vmware_engine_network = "projects/{{.Provider.project}}/locations/global/vmwareEngineNetworks/gg-asset-17806-51c7-network"
  internet_access {
    enabled = true
  }
  external_ip {
    enabled = true
  }
}
