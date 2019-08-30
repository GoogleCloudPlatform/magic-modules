resource "google_compute_snapshot" "snapshot" {
	name = "my-snapshot"
	source_disk = "${google_compute_disk.persistent.name}"
	zone = "us-central1-a"
	labels = {
		my_label = "value"
	}
}

data "google_compute_image" "debian" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "persistent" {
	name = "debian-disk"
	image = "${data.google_compute_image.debian.self_link}"
	size = 10
	type = "pd-ssd"
	zone = "us-central1-a"
}
