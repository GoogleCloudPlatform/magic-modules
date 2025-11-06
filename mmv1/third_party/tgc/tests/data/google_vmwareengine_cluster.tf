resource "google_vmwareengine_cluster" "main" {
  provider = google-beta
  name = "gg-asset-cl-38930-c6db"
  # Add parent = "projects/{{.Provider.project}}/locations/us-central1-a/privateClouds/gg-asset-pc-38930-c6db" when parent issue is fixed
  parent = ""

  node_type_configs {
    node_type_id = "standard-72"
    node_count = 1
    custom_core_count = 32
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