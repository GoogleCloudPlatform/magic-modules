resource "google_compute_target_instance" "default" {
  name        = "target-${local.name_suffix}"
  instance    = "${google_compute_instance.target-vm.self_link}"
}

data "google_compute_image" "vmimage" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance" "target-vm" {
  name         = "target-vm-${local.name_suffix}"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params{
      image = "${data.google_compute_image.vmimage.self_link}"
    }
  }

  network_interface {
    network = "default"
  }
}
