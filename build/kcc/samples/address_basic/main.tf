resource "google_compute_address" "ip_address" {
  name = "my-address-${local.name_suffix}"
}
