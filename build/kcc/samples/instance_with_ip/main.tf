resource "google_compute_address" "static" {
  name = "ipv4-address-${local.name_suffix}"
}

data "google_compute_image" "debian_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_instance" "instance_with_ip" {
	name         = "vm-instance-${local.name_suffix}"
	machine_type = "f1-micro"
	zone         = "us-central1-a"

	boot_disk {
		initialize_params{
			image = "${data.google_compute_image.debian_image.self_link}"
		}
	}

	network_interface {
		network = "default"
		access_config {
			nat_ip = "${google_compute_address.static.address}"
		}
	}
}
