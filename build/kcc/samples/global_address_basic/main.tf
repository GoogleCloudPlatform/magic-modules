resource "google_compute_global_address" "default" {
  name = "global-appserver-ip-${local.name_suffix}"
}
