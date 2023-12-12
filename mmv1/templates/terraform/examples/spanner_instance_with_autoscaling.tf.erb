resource "google_spanner_instance" "example" {
  config       = "regional-us-central1"
  display_name = "Test Spanner Instance"
  autoscaling_config {
    autoscaling_limits {
      // Define the minimum and maximum compute capacity allocated to the instance
      // Either use nodes or processing units to specify the limits,
      // but should use the same unit to set both the min_limit and max_limit.
      max_processing_units            = 3000 // OR max_nodes  = 3
      min_processing_units            = 2000 // OR min_nodes = 2
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 75
      storage_utilization_percent           = 90
    }
  }
  labels = {
    "foo" = "bar"
  }
}
