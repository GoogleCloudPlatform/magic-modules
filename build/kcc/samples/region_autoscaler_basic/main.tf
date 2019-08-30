resource "google_compute_region_autoscaler" "foobar" {
  name   = "my-region-autoscaler-${local.name_suffix}"
  region = "us-central1"
  target = "${google_compute_region_instance_group_manager.foobar.self_link}"

  autoscaling_policy {
    max_replicas    = 5
    min_replicas    = 1
    cooldown_period = 60

    cpu_utilization {
      target = 0.5
    }
  }
}

resource "google_compute_instance_template" "foobar" {
  name           = "my-instance-template-${local.name_suffix}"
  machine_type   = "n1-standard-1"
  can_ip_forward = false

  tags = ["foo", "bar"]

  disk {
    source_image = "${data.google_compute_image.debian_9.self_link}"
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_target_pool" "foobar" {
  name = "my-target-pool-${local.name_suffix}"
}

resource "google_compute_region_instance_group_manager" "foobar" {
  name   = "my-region-igm-${local.name_suffix}"
  region = "us-central1"

  instance_template  = "${google_compute_instance_template.foobar.self_link}"

  target_pools       = ["${google_compute_target_pool.foobar.self_link}"]
  base_instance_name = "foobar"
}

data "google_compute_image" "debian_9" {
	family  = "debian-9"
	project = "debian-cloud"
}
