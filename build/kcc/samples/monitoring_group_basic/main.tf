resource "google_monitoring_group" "basic" {
  display_name = "New Test Group"

  filter = "resource.metadata.region=\"europe-west2\""
}
