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

variable "instance_template" {
  type = "map"
}

variable "global_address" {
  type = "map"
}

variable "url_map" {
  type = "map"
}

variable "http_proxy" {
  type = "map"
}

variable "global_forwarding_rule" {
  type = "map"
}

variable "target_tcp_proxy" {
  type = "map"
}

variable "regional_cluster" {
  type = "map"
}

variable "route" {
  type = "map"
}

variable "router" {
  type = "map"
}

variable "snapshot" {
  type = "map"
}

variable "https_proxy" {
  type = "map"
}

variable "ssl_certificate" {
  type = "map"
}

variable "dataset" {
  type = "map"
}

variable "bigquery_table" {
  type = "map"
}

variable "repository" {
  type = "map"
}

variable "folder" {
  type = "map"
}

variable "gcp_organization_id" {
  type = "string"
  default = "none"
}

variable "cloudfunction" {
  type = "map"
}

variable "backend_bucket" {
  type = "map"
}

variable "gcp_cloud_function_region" {}

variable "regional_node_pool" {
  type = "map"
}

variable "region_backend_service_health_check" {
  type = "map"
}

variable "region_backend_service" {
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

resource "google_compute_health_check" "gcp-inspec-region-backend-service-hc" {
 project = "${var.gcp_project_id}"
 name = "${var.region_backend_service_health_check["name"]}"

 timeout_sec = "${var.region_backend_service_health_check["timeout_sec"]}"
 check_interval_sec = "${var.region_backend_service_health_check["check_interval_sec"]}"

 tcp_health_check {
   port = "${var.region_backend_service_health_check["tcp_health_check_port"]}"
 }
}

resource "google_compute_region_backend_service" "gcp-inspec-region-backend-service" {
  project     = "${var.gcp_project_id}"
  region      = "${var.gcp_location}"
  name        = "${var.region_backend_service["name"]}"
  description = "${var.region_backend_service["description"]}"
  protocol    = "${var.region_backend_service["protocol"]}"
  timeout_sec = "${var.region_backend_service["timeout_sec"]}"

  health_checks = ["${google_compute_health_check.gcp-inspec-region-backend-service-hc.self_link}"]
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

resource "google_compute_instance_template" "gcp-inspec-instance-template" {
  project     = "${var.gcp_project_id}"
  name        = "${var.instance_template["name"]}"
  description = "${var.instance_template["description"]}"

  tags = ["${var.instance_template["tag"]}"]

  instance_description = "${var.instance_template["instance_description"]}"
  machine_type         = "${var.instance_template["machine_type"]}"
  can_ip_forward       = "${var.instance_template["can_ip_forward"]}"

  scheduling {
    automatic_restart   = "${var.instance_template["scheduling_automatic_restart"]}"
    on_host_maintenance = "${var.instance_template["scheduling_on_host_maintenance"]}"
  }

  // Create a new boot disk from an image
  disk {
    source_image = "${var.instance_template["disk_source_image"]}"
    auto_delete  = "${var.instance_template["disk_auto_delete"]}"
    boot         = "${var.instance_template["disk_boot"]}"
  }

  network_interface {
    network = "${var.instance_template["network_interface_network"]}"
  }

  service_account {
    scopes = ["${var.instance_template["service_account_scope"]}"]
  }
}

resource "google_compute_global_address" "gcp-inspec-global-address" {
  project = "${var.gcp_project_id}"
  name = "${var.global_address["name"]}"
  ip_version = "${var.global_address["ip_version"]}"
}

resource "google_compute_url_map" "gcp-inspec-url-map" {
  project     = "${var.gcp_project_id}"
  name        = "${var.url_map["name"]}"
  description = "${var.url_map["description"]}"

  default_service = "${google_compute_backend_service.gcp-inspec-backend-service.self_link}"

  host_rule {
    hosts        = ["${var.url_map["host_rule_host"]}"]
    path_matcher = "${var.url_map["path_matcher_name"]}"
  }

  path_matcher {
    name            = "${var.url_map["path_matcher_name"]}"
    default_service = "${google_compute_backend_service.gcp-inspec-backend-service.self_link}"

    path_rule {
      paths   = ["${var.url_map["path_rule_path"]}"]
      service = "${google_compute_backend_service.gcp-inspec-backend-service.self_link}"
    }
  }

  test {
    service = "${google_compute_backend_service.gcp-inspec-backend-service.self_link}"
    host    = "${var.url_map["test_host"]}"
    path    = "${var.url_map["test_path"]}"
  }
}

resource "google_compute_target_http_proxy" "gcp-inspec-http-proxy" {
  project     = "${var.gcp_project_id}"
  name        = "${var.http_proxy["name"]}"
  url_map     = "${google_compute_url_map.gcp-inspec-url-map.self_link}"
  description = "${var.http_proxy["description"]}"
}

resource "google_compute_global_forwarding_rule" "gcp-inspec-global-forwarding-rule" {
  project    = "${var.gcp_project_id}"
  name       = "${var.global_forwarding_rule["name"]}"
  target     = "${google_compute_target_http_proxy.gcp-inspec-http-proxy.self_link}"
  port_range = "${var.global_forwarding_rule["port_range"]}"
}

resource "google_compute_backend_service" "gcp-inspec-tcp-backend-service" {
  project       = "${var.gcp_project_id}"
  name          = "${var.target_tcp_proxy["tcp_backend_service_name"]}"
  protocol      = "TCP"
  timeout_sec   = 10

  health_checks = ["${google_compute_health_check.gcp-inspec-health-check.self_link}"]
}

resource "google_compute_target_tcp_proxy" "gcp-inspec-target-tcp-proxy" {
  project         = "${var.gcp_project_id}"
  name            = "${var.target_tcp_proxy["name"]}"
  proxy_header    = "${var.target_tcp_proxy["proxy_header"]}"
  backend_service = "${google_compute_backend_service.gcp-inspec-tcp-backend-service.self_link}"
}

resource "google_container_cluster" "gcp-inspec-regional-cluster" {
  project = "${var.gcp_project_id}"
  name = "${var.regional_cluster["name"]}"
  region = "${var.gcp_location}"
  initial_node_count = 1
  remove_default_node_pool = true

  maintenance_policy {
    daily_maintenance_window {
      start_time = "23:00"
    }
  }
}

resource "google_compute_route" "gcp-inspec-route" {
  project     = "${var.gcp_project_id}"
  name        = "${var.route["name"]}"
  dest_range  = "${var.route["dest_range"]}"
  network     = "${google_compute_network.inspec-gcp-network.name}"
  next_hop_ip = "${var.route["next_hop_ip"]}"
  priority    = "${var.route["priority"]}"
  # google_compute_route depends on next_hop_ip belonging to a subnetwork
  # of the named network in this block. Since inspec-gcp-network does not
  # automatically create subnetworks, we need to create a dependency so
  # the route is not created before the subnetwork 
  depends_on  = ["google_compute_subnetwork.inspec-gcp-subnetwork"]
}

resource "google_compute_router" "gcp-inspec-router" {
  project = "${var.gcp_project_id}"
  name    = "${var.router["name"]}"
  network = "${google_compute_network.inspec-gcp-network.name}"
  bgp {
    asn               = "${var.router["bgp_asn"]}"
    advertise_mode    = "${var.router["bgp_advertise_mode"]}"
    advertised_groups = ["${var.router["bgp_advertised_group"]}"]
    advertised_ip_ranges {
      range = "${var.router["bgp_advertised_ip_range1"]}"
    }
    advertised_ip_ranges {
      range = "${var.router["bgp_advertised_ip_range2"]}"
    }
  }
}

resource "google_compute_snapshot" "gcp-inspec-snapshot" {
  project = "${var.gcp_project_id}"
  name = "${var.snapshot["name"]}"
  source_disk = "${google_compute_disk.generic_compute_disk.name}"
  zone = "${var.gcp_zone}"
  # Depends on the instance of the disk we are using. Allow instance to spin up
  # Before snapshotting the disk to avoid resourceInUse errors
  depends_on  = ["google_compute_instance.generic_external_vm_instance_data_disk"]
}

resource "google_compute_ssl_certificate" "gcp-inspec-ssl-certificate" {
  project     = "${var.gcp_project_id}"
  name        = "${var.ssl_certificate["name"]}"
  private_key = "${var.ssl_certificate["private_key"]}"
  certificate = "${var.ssl_certificate["certificate"]}"
  description = "${var.ssl_certificate["description"]}"
}

resource "google_compute_target_https_proxy" "gcp-inspec-https-proxy" {
  project     = "${var.gcp_project_id}"
  name        = "${var.https_proxy["name"]}"
  url_map     = "${google_compute_url_map.gcp-inspec-url-map.self_link}"
  description = "${var.https_proxy["description"]}"
  ssl_certificates = ["${google_compute_ssl_certificate.gcp-inspec-ssl-certificate.self_link}"]
}

resource "google_bigquery_dataset" "gcp-inspec-dataset" {
  project                     = "${var.gcp_project_id}"
  dataset_id                  = "${var.dataset["dataset_id"]}"
  friendly_name               = "${var.dataset["friendly_name"]}"
  description                 = "${var.dataset["description"]}"
  location                    = "${var.dataset["location"]}"
  default_table_expiration_ms = "${var.dataset["default_table_expiration_ms"]}"

  access {
    role          = "${var.dataset["access_writer_role"]}"
    special_group = "${var.dataset["access_writer_special_group"]}"
  }

  access {
    role          = "OWNER"
    special_group = "projectOwners"
  }
}

resource "google_bigquery_table" "gcp-inspec-bigquery-table" {
  project    = "${var.gcp_project_id}"
  dataset_id = "${google_bigquery_dataset.gcp-inspec-dataset.dataset_id}"
  table_id   = "${var.bigquery_table["table_id"]}"

  time_partitioning {
    type = "${var.bigquery_table["time_partitioning_type"]}"
  }

  description = "${var.bigquery_table["description"]}"
  expiration_time = "${var.bigquery_table["expiration_time"]}"
}

resource "google_sourcerepo_repository" "gcp-inspec-sourcerepo-repository" {
  project = "${var.gcp_project_id}"
  name = "${var.repository["name"]}"
}

resource "google_folder" "inspec-gcp-folder" {
  count = "${var.gcp_organization_id == "none" ? 0 : var.gcp_enable_privileged_resources}"
  display_name = "${var.folder["display_name"]}"
  parent       = "${var.gcp_organization_id}"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.js.zip"
  bucket = "${google_storage_bucket.generic-storage-bucket.name}"
  source = "../configuration/index.js.zip"
}

resource "google_cloudfunctions_function" "function" {
  project               = "${var.gcp_project_id}"
  region                = "${var.gcp_cloud_function_region}"
  name                  = "${var.cloudfunction["name"]}"
  description           = "${var.cloudfunction["description"]}"
  available_memory_mb   = "${var.cloudfunction["available_memory_mb"]}"
  source_archive_bucket = "${google_storage_bucket.generic-storage-bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  trigger_http          = "${var.cloudfunction["trigger_http"]}"
  timeout               = "${var.cloudfunction["timeout"]}"
  entry_point           = "${var.cloudfunction["entry_point"]}"
  runtime               = "nodejs6"

  environment_variables = {
    MY_ENV_VAR = "${var.cloudfunction["env_var_value"]}"
  }
}

resource "google_compute_backend_bucket" "image_backend" {
  project     = "${var.gcp_project_id}"
  name        = "${var.backend_bucket["name"]}"
  description = "${var.backend_bucket["description"]}"
  bucket_name = "${google_storage_bucket.generic-storage-bucket.name}"
  enable_cdn  = "${var.backend_bucket["enable_cdn"]}"
}

resource "google_container_node_pool" "inspec-gcp-regional-node-pool" {
  project    = "${var.gcp_project_id}"
  name       = "${var.regional_node_pool["name"]}"
  region     = "${var.gcp_location}"
  cluster    = "${google_container_cluster.gcp-inspec-regional-cluster.name}"
  node_count = "${var.regional_node_pool["node_count"]}"
}