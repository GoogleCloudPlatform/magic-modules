# First test without squash mode. Classic rules
# Runs with old provider version to make sure API changes don't break existing install base
# Will complain about storage-pool drift. Ignore it.
data "google_compute_network" "default" {
    name = "default"
}

resource "google_netapp_storage_pool" "default" {
  name = "my-tf-pool2"
  location = "us-central1"
  service_level = "PREMIUM"
  capacity_gib = "2048"
  network = data.google_compute_network.default.id
}

resource "google_netapp_volume" "my-nfsv3-volume" {
  location         = "us-central1"
  name             = "ok2-squash-nfsv3-volume"
  capacity_gib     = 100
  share_name       = "ok-squash-nfsv3-volume"
  storage_pool     = google_netapp_storage_pool.default.name
  protocols        = ["NFSV3"]
  unix_permissions = "0777"

  export_policy {
    rules {
      access_type     = "READ_WRITE"
      allowed_clients = "10.0.1.0/24"
      nfsv3           = true
      squash_mode     = "NO_ROOT_SQUASH"
    }
  }
}
