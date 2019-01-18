variable "ssl_policy" {
  type = "map"
}

variable "topic" {
  type = "map"
}

variable "subscription" {
  type = "map"
}

variable "managed_zone" {
	type = "map"
}

variable "record_set" {
	type = "map"
}

variable "instance_group_manager" {
  type = "map"
}

variable "autoscaler" {
  type = "map"
}

variable "target_pool" {
  type = "map"
}

resource "google_compute_ssl_policy" "custom-ssl-policy" {
  name            = "${var.ssl_policy["name"]}"
  min_tls_version = "${var.ssl_policy["min_tls_version"]}"
  profile         = "${var.ssl_policy["profile"]}"
  custom_features = ["${var.ssl_policy["custom_feature"]}", "${var.ssl_policy["custom_feature2"]}"]
  project         = "${var.gcp_project_id}"
}

resource "google_pubsub_topic" "topic" {
  project = "${var.gcp_project_id}"
  name    = "${var.topic["name"]}"
}

resource "google_pubsub_subscription" "default" {
  project              = "${var.gcp_project_id}"
  name                 = "${var.subscription["name"]}"
  topic                = "${google_pubsub_topic.topic.name}"
  ack_deadline_seconds = "${var.subscription["ack_deadline_seconds"]}"
}

resource "google_dns_managed_zone" "prod" {
  name        = "${var.managed_zone["name"]}"
  dns_name    = "${var.managed_zone["dns_name"]}"
  description = "${var.managed_zone["description"]}"

  labels = {
    key = "${var.managed_zone["label_value"]}"
  }
  project = "${var.gcp_project_id}"
}

resource "google_dns_record_set" "a" {
  name = "${var.record_set["name"]}"
  managed_zone = "${google_dns_managed_zone.prod.name}"
  type = "${var.record_set["type"]}"
  ttl  = "${var.record_set["ttl"]}"

  rrdatas = ["${var.record_set["rrdatas1"]}", "${var.record_set["rrdatas2"]}"]
  project = "${var.gcp_project_id}"
}

resource "google_compute_instance_group_manager" "gcp-inspec-igm" {
  project           = "${var.gcp_project_id}"
  zone              = "${var.gcp_zone}"
  name              = "${var.instance_group_manager["name"]}"
  instance_template = "${google_compute_instance_template.default.self_link}"
  base_instance_name        = "${var.instance_group_manager["base_instance_name"]}"
  target_pools = []
  target_size  = 0
  named_port {
    name = "${var.instance_group_manager["named_port_name"]}"
    port = "${var.instance_group_manager["named_port_port"]}"
  }
}

resource "google_compute_autoscaler" "gcp-inspec-autoscaler" {
  project = "${var.gcp_project_id}"
  name    = "${var.autoscaler["name"]}"
  zone    = "${var.gcp_zone}"
  target  = "${google_compute_instance_group_manager.gcp-inspec-igm.self_link}"

  autoscaling_policy = {
    max_replicas    = "${var.autoscaler["max_replicas"]}"
    min_replicas    = "${var.autoscaler["min_replicas"]}"
    cooldown_period = "${var.autoscaler["cooldown_period"]}"

    cpu_utilization {
      target = "${var.autoscaler["cpu_utilization_target"]}"
    }
  }
}

resource "google_compute_target_pool" "gcp-inspec-target-pool" {
  project = "${var.gcp_project_id}"
  name = "${var.target_pool["name"]}"
  session_affinity = "${var.target_pool["session_affinity"]}"
  
  instances = [
    "${var.gcp_zone}/${var.gcp_ext_vm_name}",
  ]
}