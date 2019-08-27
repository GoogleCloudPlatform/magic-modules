resource "google_compute_image" "example" {
  name = "example-image-${local.name_suffix}"

  raw_disk {
    source = "https://storage.googleapis.com/bosh-cpi-artifacts/bosh-stemcell-3262.4-google-kvm-ubuntu-trusty-go_agent-raw.tar.gz"
  }
}
