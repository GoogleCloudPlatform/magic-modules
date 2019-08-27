resource "google_monitoring_group" "basic" {
  display_name = "New Test Group-${local.name_suffix}"

  filter = "resource.metadata.region=\"europe-west2\""
}
