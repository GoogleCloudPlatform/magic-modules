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

variable "trigger" {
  type = "map"
}

variable "health_check" {
  type = "map"
}

variable "backend_service" {
  type = "map"
}

variable "http_health_check" {
  type = "map"
}

variable "https_health_check" {
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

resource "google_cloudbuild_trigger" "gcp-inspec-cloudbuild-trigger" {
  project = "${var.gcp_project_id}"
  trigger_template {
    branch_name = "${var.trigger["trigger_template_branch"]}"
    project     = "${var.trigger["trigger_template_project"]}"
    repo_name   = "${var.trigger["trigger_template_repo"]}"
  }
  filename = "${var.trigger["filename"]}"
}

resource "google_compute_health_check" "gcp-inspec-health-check" {
 project = "${var.gcp_project_id}"
 name = "${var.health_check["name"]}"

 timeout_sec = "${var.health_check["timeout_sec"]}"
 check_interval_sec = "${var.health_check["check_interval_sec"]}"

 tcp_health_check {
   port = "${var.health_check["tcp_health_check_port"]}"
 }
}

resource "google_compute_backend_service" "gcp-inspec-backend-service" {
  project     = "${var.gcp_project_id}"
  name        = "${var.backend_service["name"]}"
  description = "${var.backend_service["description"]}"
  port_name   = "${var.backend_service["port_name"]}"
  protocol    = "${var.backend_service["protocol"]}"
  timeout_sec = "${var.backend_service["timeout_sec"]}"
  enable_cdn  = "${var.backend_service["enable_cdn"]}"

  backend {
    group = "${google_compute_instance_group_manager.gcp-inspec-igm.instance_group}"
  }

  health_checks = ["${google_compute_health_check.gcp-inspec-health-check.self_link}"]
}

resource "google_compute_http_health_check" "gcp-inspec-http-health-check" {
  project      = "${var.gcp_project_id}"
  name         = "${var.http_health_check["name"]}"
  request_path = "${var.http_health_check["request_path"]}"

  timeout_sec        = "${var.http_health_check["timeout_sec"]}"
  check_interval_sec = "${var.http_health_check["check_interval_sec"]}"
}

resource "google_compute_https_health_check" "gcp-inspec-https-health-check" {
  project      = "${var.gcp_project_id}"
  name         = "${var.https_health_check["name"]}"
  request_path = "${var.https_health_check["request_path"]}"

  timeout_sec         = "${var.https_health_check["timeout_sec"]}"
  check_interval_sec  = "${var.https_health_check["check_interval_sec"]}"
  unhealthy_threshold = "${var.https_health_check["unhealthy_threshold"]}"
}