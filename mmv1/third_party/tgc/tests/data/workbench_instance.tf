
resource "google_workbench_instance" "basic" {
  name = "workbench-instance-cai"
  location = "us-central1-a"
  gce_setup {
    machine_type = "e2-standard-2"
    shielded_instance_config {
      enable_secure_boot = true
      enable_vtpm = true
      enable_integrity_monitoring = true
    }
  }
}
