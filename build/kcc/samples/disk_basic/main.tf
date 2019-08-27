resource "google_compute_disk" "default" {
  name  = "test-disk-${local.name_suffix}"
  type  = "pd-ssd"
  zone  = "us-central1-a"
  image = "debian-8-jessie-v20170523"
  labels = {
    environment = "dev"
  }
  physical_block_size_bytes = 4096
}
