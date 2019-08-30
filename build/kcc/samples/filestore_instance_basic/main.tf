resource "google_filestore_instance" "instance" {
  name = "test-instance"
  zone = "us-central1-b"
  tier = "PREMIUM"

  file_shares {
    capacity_gb = 2660
    name        = "share1"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
}
