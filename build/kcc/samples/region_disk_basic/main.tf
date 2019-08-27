resource "google_compute_region_disk" "regiondisk" {
  name = "my-region-disk-${local.name_suffix}"
  snapshot = "${google_compute_snapshot.snapdisk.self_link}"
  type = "pd-ssd"
  region = "us-central1"
  physical_block_size_bytes = 4096

  replica_zones = ["us-central1-a", "us-central1-f"]
}

resource "google_compute_disk" "disk" {
  name = "my-disk-${local.name_suffix}"
  image = "debian-cloud/debian-9"
  size = 50
  type = "pd-ssd"
  zone = "us-central1-a"
}

resource "google_compute_snapshot" "snapdisk" {
  name = "my-snapshot-${local.name_suffix}"
  source_disk = "${google_compute_disk.disk.name}"
  zone = "us-central1-a"
}
