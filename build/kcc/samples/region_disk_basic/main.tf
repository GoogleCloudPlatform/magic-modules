resource "google_compute_region_disk" "regiondisk" {
  name = "my-region-disk"
  snapshot = "${google_compute_snapshot.snapdisk.self_link}"
  type = "pd-ssd"
  region = "us-central1"
  physical_block_size_bytes = 4096

  replica_zones = ["us-central1-a", "us-central1-f"]
}

resource "google_compute_disk" "disk" {
  name = "my-disk"
  image = "debian-cloud/debian-9"
  size = 50
  type = "pd-ssd"
  zone = "us-central1-a"
}

resource "google_compute_snapshot" "snapdisk" {
  name = "my-snapshot"
  source_disk = "${google_compute_disk.disk.name}"
  zone = "us-central1-a"
}
