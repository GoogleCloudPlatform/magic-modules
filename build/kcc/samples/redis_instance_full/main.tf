resource "google_redis_instance" "cache" {
  name           = "ha-memory-cache"
  tier           = "STANDARD_HA"
  memory_size_gb = 1

  location_id             = "us-central1-a"
  alternative_location_id = "us-central1-f"

  authorized_network = "${google_compute_network.auto-network.self_link}"

  redis_version     = "REDIS_3_2"
  display_name      = "Terraform Test Instance"
  reserved_ip_range = "192.168.0.0/29"

  labels = {
    my_key    = "my_val"
    other_key = "other_val"
  }
}

resource "google_compute_network" "auto-network" {
  name = "authorized-network"
}
