resource "google_monitoring_group" "parent" {
  display_name = "New Test SubGroup"
  filter = "resource.metadata.region=\"europe-west2\""
}

resource "google_monitoring_group" "subgroup" {
  display_name = "New Test SubGroup"
  filter = "resource.metadata.region=\"europe-west2\""
  parent_name =  "${google_monitoring_group.parent.name}"
}
