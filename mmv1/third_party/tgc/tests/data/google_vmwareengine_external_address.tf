resource "google_vmwareengine_network" "main" {
  name        = "gg-asset-vmw-net-03361-811b"
  location    = "global"
  type        = "STANDARD"
  description = "VMware Engine Network"
}

resource "google_vmwareengine_network_policy" "main" {
  name                  = "gg-asset-np-03361-811b"
  location              = "us-central1"
  vmware_engine_network = google_vmwareengine_network.main.id
  edge_services_cidr    = "192.168.30.0/26"
  description           = "Network policy to enable external IP access"

  internet_access {
    enabled = true
  }

  external_ip {
    enabled = true
  }
}

resource "google_vmwareengine_private_cloud" "main" {
  name        = "gg-asset-pc-03361-811b"
  location    = "us-central1-a"
  type        = "TIME_LIMITED"
  description = "Private Cloud for External Address testing"

  network_config {
    management_cidr       = "192.168.0.0/24"
    vmware_engine_network = google_vmwareengine_network.main.id
  }

  management_cluster {
    cluster_id = "gg-asset-mgmt-cl-03361-811b"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 1
    }
  }
}

resource "google_vmwareengine_external_address" "main" {
  name        = "gg-asset-ext-addr-03361-811b"
  parent      = google_vmwareengine_private_cloud.main.id
  internal_ip = "10.100.0.10"
  description = "External address for testing"
  depends_on  = [google_vmwareengine_network_policy.main]
}
