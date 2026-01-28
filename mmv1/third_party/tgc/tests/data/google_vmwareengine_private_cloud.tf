resource "google_vmwareengine_private_cloud" "main" {
  name        = "gg-asset-pc-00563-8560"
  location    = "us-central1-a"
  description = "Standard Private Cloud"
  project     = "{{.Provider.project}}"

  network_config {
    management_cidr       = "192.168.0.0/24"
    vmware_engine_network = "projects/{{.Provider.project}}/locations/global/vmwareEngineNetworks/test_vmware_network"
  }

  management_cluster {
    cluster_id = "gg-asset-mgmt-cl-00563-8560"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 3
      custom_core_count = 36
    }
    node_type_configs {
      node_type_id = "standard-128"
      node_count   = 3
    }
    autoscaling_settings {
      autoscaling_policies {
        autoscale_policy_id = "autoscaling-policy"
        node_type_id = "standard-72"
        scale_out_size = 1
        cpu_thresholds {
          scale_out = 80
          scale_in  = 15
        }
        consumed_memory_thresholds {
          scale_out = 75
          scale_in  = 20
        }
        storage_thresholds {
          scale_out = 80
          scale_in  = 20
        }
      }
      min_cluster_node_count = 3
      max_cluster_node_count = 8
      cool_down_period = "1800s"
    }
  }
}