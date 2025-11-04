
resource "google_vmwareengine_network" "main" {
  name        = "gg-asset-vmw-net-00563-8560"
  location    = "global"
  type        = "STANDARD"
  description = "Standard VMware Engine Network"
  project     = "{{.Provider.project}}"
}

resource "google_vmwareengine_private_cloud" "main" {
  name        = "gg-asset-pc-00563-8560"
  location    = "us-central1-a"
  description = "Standard Private Cloud"
  project     = "{{.Provider.project}}"

  network_config {
    management_cidr       = "192.168.0.0/24"
    vmware_engine_network = google_vmwareengine_network.main.id
  }

  management_cluster {
    cluster_id = "gg-asset-mgmt-cl-00563-8560"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 3
    }
  }
}
