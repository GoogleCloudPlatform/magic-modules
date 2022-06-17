# Cloud Run service replicated across multiple GCP regions

######################################
# Terraform (Project-wide resources) #
######################################
data "google_project" "app" {
  project_id = var.project_id
}

resource "google_project_service" "compute_api" {
  project                    = data.google_project.app.project_id
  service                    = "compute.googleapis.com"
  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_project_service" "run_api" {
  project                    = data.google_project.app.project_id
  service                    = "run.googleapis.com"
  disable_dependent_services = true
  disable_on_destroy         = false
}

# [START cloudloadbalancing_run_multiregion]
############################
# Input variables (config) #
############################
variable "project_id" {
  type        = string
  description = "Your GCP Project ID"
}

variable "domain_name" {
  type        = string
  description = "Your domain name (e.g. `example.com`)"
}

variable "run_container_image" {
  type        = string
  description = "URL of the container image to run."
  default     = "us-docker.pkg.dev/cloudrun/container/hello"
}

variable "run_regions" {
  type        = list(string)
  description = "The regions to deploy this application to"
  default     = ["us-central1", "europe-west1"]
}

######################################
# Terraform (Load Balancer creation) #
######################################

resource "google_compute_global_address" "lb_default" {
  name    = "myservice-service-ip"
  project = data.google_project.app.project_id

  # Use an explicit depends_on clause to wait until API is enabled
  depends_on = [
    google_project_service.compute_api
  ]
}

resource "google_compute_backend_service" "lb_default" {
  name                  = "myservice-backend"
  project               = data.google_project.app.project_id
  load_balancing_scheme = "EXTERNAL_MANAGED"

  backend {
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 0.85 # TODO what does this do? docs are unclear! :(
    group           = google_compute_region_network_endpoint_group.lb_default[0].id
  }

  backend {
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 0.85 # TODO what does this do? docs are unclear! :(
    group           = google_compute_region_network_endpoint_group.lb_default[1].id
  }

  # Use an explicit depends_on clause to wait until API is enabled
  depends_on = [
    google_project_service.compute_api,
  ]
}


resource "google_compute_url_map" "lb_default" {
  name            = "myservice-lb-urlmap"
  project         = data.google_project.app.project_id
  default_service = google_compute_backend_service.lb_default.id

  path_matcher {
    name            = "allpaths"
    default_service = google_compute_backend_service.lb_default.id
    route_rules {
      priority = 1
      url_redirect {
        https_redirect         = true
        redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
      }
    }
  }
}

resource "google_compute_managed_ssl_certificate" "lb_default" {
  name    = "myservice-ssl-cert"
  project = data.google_project.app.project_id

  managed {
    domains = [var.domain_name]
  }
}

resource "google_compute_target_https_proxy" "lb_default" {
  name    = "myservice-https-proxy"
  project = data.google_project.app.project_id
  url_map = google_compute_url_map.lb_default.id
  ssl_certificates = [
    google_compute_managed_ssl_certificate.lb_default.name
  ]
  depends_on = [
    google_compute_managed_ssl_certificate.lb_default
  ]
}

resource "google_compute_global_forwarding_rule" "lb_default" {
  name                  = "myservice-lb-forwarding-rule"
  project               = data.google_project.app.project_id
  load_balancing_scheme = "EXTERNAL_MANAGED"
  target                = google_compute_target_https_proxy.lb_default.id
  ip_address            = google_compute_global_address.lb_default.id
  port_range            = "443"
  depends_on            = [google_compute_target_https_proxy.lb_default]
}

resource "google_compute_region_network_endpoint_group" "lb_default" {
  count                 = length(var.run_regions)
  project               = data.google_project.app.project_id
  name                  = "myservice-neg"
  network_endpoint_type = "SERVERLESS"
  region                = var.run_regions[count.index]
  cloud_run {
    service = google_cloud_run_service.run_default[count.index].name
  }
}

output "load_balancer_ip_addr" {
  value = google_compute_global_address.lb_default.address
}

#############################
# Terraform (Cloud Run app) #
#############################
resource "google_cloud_run_service" "run_default" {
  count    = length(var.run_regions)
  project  = data.google_project.app.project_id
  name     = "myservice-run-app-BANANAMAN"
  location = var.run_regions[count.index]

  template {
    spec {
      containers {
        image = var.run_container_image
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  # Use an explicit depends_on clause to wait until API is enabled
  depends_on = [
    google_project_service.run_api
  ]
}

data "google_iam_policy" "run_allow_unauthenticated" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "run_default" {
  count    = length(var.run_regions)
  project  = data.google_project.app.project_id
  location = google_cloud_run_service.run_default[count.index].location
  service  = google_cloud_run_service.run_default[count.index].name

  policy_data = data.google_iam_policy.run_allow_unauthenticated.policy_data
}

#############################
# Terraform (HTTP -> HTTPS) #
#############################
resource "google_compute_url_map" "https_default" {
  name    = "myservice-https-urlmap"
  project = data.google_project.app.project_id

  default_url_redirect {
    redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
    https_redirect         = true
    strip_query            = false
  }
}

resource "google_compute_target_http_proxy" "https_default" {
  name    = "myservice-http-proxy"
  project = data.google_project.app.project_id
  url_map = google_compute_url_map.https_default.id

  depends_on = [
    google_compute_url_map.https_default
  ]
}

resource "google_compute_global_forwarding_rule" "https_default" {
  name       = "myservice-https-forwarding-rule"
  project    = data.google_project.app.project_id
  target     = google_compute_target_http_proxy.https_default.id
  ip_address = google_compute_global_address.lb_default.id
  port_range = "80"
  depends_on = [google_compute_target_http_proxy.https_default]
}
# [END cloudloadbalancing_run_multiregion]
