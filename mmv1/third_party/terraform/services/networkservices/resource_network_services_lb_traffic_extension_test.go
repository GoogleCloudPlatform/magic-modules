package networkservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetworkServicesLbTrafficExtension_networkServicesLbTrafficExtensionBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesLbTrafficExtensionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesLbTrafficExtension_basic(context),
			},
			{
				ResourceName:            "google_network_services_lb_traffic_extension.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesLbTrafficExtension_update(context),
			},
			{
				ResourceName:            "google_network_services_lb_traffic_extension.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkServicesLbTrafficExtension_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
# Internal HTTP load balancer with a managed instance group backend
# VPC network
resource "google_compute_network" "ilb_network" {
  name                    = "tf-test-l7-ilb-network%{random_suffix}"
  provider                = google-beta
  auto_create_subnetworks = false
}

# proxy-only subnet
resource "google_compute_subnetwork" "proxy_subnet" {
  name          = "tf-test-l7-ilb-proxy-subnet%{random_suffix}"
  provider      = google-beta
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-west1"
  purpose       = "REGIONAL_MANAGED_PROXY"
  role          = "ACTIVE"
  network       = google_compute_network.ilb_network.id
}

# backend subnet
resource "google_compute_subnetwork" "ilb_subnet" {
  name          = "tf-test-l7-ilb-subnet%{random_suffix}"
  provider      = google-beta
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-west1"
  network       = google_compute_network.ilb_network.id
}

# forwarding rule
resource "google_compute_forwarding_rule" "default" {
  name                  = "tf-test-l7-ilb-forwarding-rule%{random_suffix}"
  provider              = google-beta
  region                = "us-west1"
  depends_on            = [google_compute_subnetwork.proxy_subnet]
  ip_protocol           = "TCP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  port_range            = "80"
  target                = google_compute_region_target_http_proxy.default.id
  network               = google_compute_network.ilb_network.id
  subnetwork            = google_compute_subnetwork.ilb_subnet.id
  network_tier          = "PREMIUM"
}

# HTTP target proxy
resource "google_compute_region_target_http_proxy" "default" {
  name     = "tf-test-l7-ilb-target-http-proxy%{random_suffix}"
  provider = google-beta
  region   = "us-west1"
  url_map  = google_compute_region_url_map.default.id
}

# URL map
resource "google_compute_region_url_map" "default" {
  name            = "tf-test-l7-ilb-regional-url-map%{random_suffix}"
  provider        = google-beta
  region          = "us-west1"
  default_service = google_compute_region_backend_service.default.id
}

# backend service
resource "google_compute_region_backend_service" "default" {
  name                  = "tf-test-l7-ilb-backend-subnet%{random_suffix}"
  provider              = google-beta
  region                = "us-west1"
  protocol              = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec           = 10
  health_checks         = [google_compute_region_health_check.default.id]
  backend {
    group           = google_compute_region_instance_group_manager.mig.instance_group
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 1.0
  }
}

# instance template
resource "google_compute_instance_template" "instance_template" {
  name         = "tf-test-l7-ilb-mig-template%{random_suffix}"
  provider     = google-beta
  machine_type = "e2-small"
  tags         = ["http-server"]

  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id
    access_config {
      # add external ip to fetch packages
    }
  }

  disk {
    source_image = "debian-cloud/debian-10"
    auto_delete  = true
    boot         = true
  }

  # install nginx and serve a simple web page
  metadata = {
    startup-script = <<-EOF1
      #! /bin/bash
      set -euo pipefail

      export DEBIAN_FRONTEND=noninteractive
      apt-get update
      apt-get install -y nginx-light jq

      NAME=$(curl -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/hostname")
      IP=$(curl -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip")
      METADATA=$(curl -f -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/attributes/?recursive=True" | jq 'del(.["startup-script"])')

      cat <<EOF > /var/www/html/index.html
      <pre>
      Name: $NAME
      IP: $IP
      Metadata: $METADATA
      </pre>
      EOF
    EOF1
  }

  lifecycle {
    create_before_destroy = true
  }
}

# health check
resource "google_compute_region_health_check" "default" {
  name     = "tf-test-l7-ilb-hc%{random_suffix}"
  provider = google-beta
  region   = "us-west1"

  http_health_check {
    port_specification = "USE_SERVING_PORT"
  }
}

# MIG
resource "google_compute_region_instance_group_manager" "mig" {
  name     = "tf-test-l7-ilb-mig1%{random_suffix}"
  provider = google-beta
  region   = "us-west1"

  base_instance_name = "vm"
  target_size        = 2

  version {
    instance_template = google_compute_instance_template.instance_template.id
    name              = "primary"
  }
}

# allow all access from IAP and health check ranges
resource "google_compute_firewall" "fw-iap" {
  name          = "tf-test-l7-ilb-fw-allow-iap-hc%{random_suffix}"
  provider      = google-beta
  direction     = "INGRESS"
  network       = google_compute_network.ilb_network.id
  source_ranges = ["130.211.0.0/22", "35.191.0.0/16", "35.235.240.0/20"]

  allow {
    protocol = "tcp"
  }
}

# allow http from proxy subnet to backends
resource "google_compute_firewall" "fw-ilb-to-backends" {
  name          = "tf-test-l7-ilb-fw-allow-ilb-to-backends%{random_suffix}"
  provider      = google-beta
  direction     = "INGRESS"
  network       = google_compute_network.ilb_network.id
  source_ranges = ["10.0.0.0/24"]
  target_tags   = ["http-server"]

  allow {
    protocol = "tcp"
    ports    = ["80", "443", "8080"]
  }
}

# test instance
resource "google_compute_instance" "vm-test" {
  name         = "tf-test-l7-ilb-test-vm%{random_suffix}"
  provider     = google-beta
  zone         = "us-west1-b"
  machine_type = "e2-small"

  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id
  }

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }
}

resource "google_network_services_lb_traffic_extension" "default" {
  provider = google-beta
  name     = "tf-test-l7-ilb-traffic-ext%{random_suffix}"
  description = "my traffic extension"
  location = "us-west1"

  load_balancing_scheme = "INTERNAL_MANAGED"
  forwarding_rules      = [google_compute_forwarding_rule.default.self_link]

  extension_chains {
      name = "chain1"

      match_condition {
          cel_expression = "request.host == 'example.com'"
      }

      extensions {
          name      = "ext11"
          authority = "ext11.com"
          service   = google_compute_region_backend_service.callouts_backend.self_link
          timeout   = "0.1s"
          fail_open = false

          supported_events = ["REQUEST_HEADERS"]
          forward_headers = ["custom-header"]
      }
  }

  labels = {
    foo = "bar"
  }
}

# Traffic Extension Backend Instance
resource "google_compute_instance" "callouts_instance" {
  provider = google-beta

  name = "tf-test-l7-ilb-callouts-ins%{random_suffix}"
  zone = "us-west1-a"

  machine_type = "e2-small"
  labels = {
    "container-vm" = "cos-stable-109-17800-147-54"
  }
  tags         = ["allow-ssh","load-balanced-backend"]

  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id
    access_config {
        # add external ip to fetch packages
    }
  }
  boot_disk {
    auto_delete  = true
    initialize_params {
      type = "pd-standard"
      size = 10
      image = "https://www.googleapis.com/compute/v1/projects/cos-cloud/global/images/cos-stable-109-17800-147-54"
    }
  }

  # Initialize an Envoy's Ext Proc gRPC API based on a docker container
  metadata = {
    gce-container-declaration = "# DISCLAIMER:\n# This container declaration format is not a public API and may change without\n# notice. Please use gcloud command-line tool or Google Cloud Console to run\n# Containers on Google Compute Engine.\n\nspec:\n  containers:\n  - image: us-docker.pkg.dev/service-extensions/ext-proc/service-callout-basic-example-python:latest\n    name: callouts-vm\n    securityContext:\n      privileged: false\n    stdin: false\n    tty: false\n    volumeMounts: []\n  restartPolicy: Always\n  volumes: []\n"
    google-logging-enabled = "true"
  }
  lifecycle {
    create_before_destroy = true
  }
}

// callouts instance group
resource "google_compute_instance_group" "callouts_instance_group" {
  provider    = google-beta
  name        = "tf-test-l7-ilb-callouts-ins-group%{random_suffix}"
  description = "Terraform test instance group"

  instances = [
    google_compute_instance.callouts_instance.id,
  ]

  named_port {
    name = "http"
    port = "80"
  }

  named_port {
    name = "grpc"
    port = "443"
  }

  zone = "us-west1-a"
}

# callout health check
resource "google_compute_region_health_check" "callouts_health_check" {
  provider = google-beta
  name     = "tf-test-l7-ilb-callouts-hc%{random_suffix}"
  region   = "us-west1"
  http_health_check {
    port = 80
  }
}

# callout backend service
resource "google_compute_region_backend_service" "callouts_backend" {
  provider              = google-beta
  name                  = "tf-test-l7-ilb-callouts-backend%{random_suffix}"
  region                = "us-west1"
  protocol              = "HTTP2"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec           = 10
  port_name             = "grpc"
  health_checks         = [google_compute_region_health_check.callouts_health_check.id]

  backend {
    group           = google_compute_instance_group.callouts_instance_group.id
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 1.0
  }
}
`, context)
}

func testAccNetworkServicesLbTrafficExtension_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
# Internal HTTP load balancer with a managed instance group backend
# VPC network
resource "google_compute_network" "ilb_network" {
  name                    = "tf-test-l7-ilb-network%{random_suffix}"
  provider                = google-beta
  auto_create_subnetworks = false
}

# proxy-only subnet
resource "google_compute_subnetwork" "proxy_subnet" {
  name          = "tf-test-l7-ilb-proxy-subnet%{random_suffix}"
  provider      = google-beta
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-west1"
  purpose       = "REGIONAL_MANAGED_PROXY"
  role          = "ACTIVE"
  network       = google_compute_network.ilb_network.id
}

# backend subnet
resource "google_compute_subnetwork" "ilb_subnet" {
  name          = "tf-test-l7-ilb-subnet%{random_suffix}"
  provider      = google-beta
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-west1"
  network       = google_compute_network.ilb_network.id
}

# forwarding rule
resource "google_compute_forwarding_rule" "default" {
  name                  = "tf-test-l7-ilb-forwarding-rule%{random_suffix}"
  provider              = google-beta
  region                = "us-west1"
  depends_on            = [google_compute_subnetwork.proxy_subnet]
  ip_protocol           = "TCP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  port_range            = "80"
  target                = google_compute_region_target_http_proxy.default.id
  network               = google_compute_network.ilb_network.id
  subnetwork            = google_compute_subnetwork.ilb_subnet.id
  network_tier          = "PREMIUM"
}

# HTTP target proxy
resource "google_compute_region_target_http_proxy" "default" {
  name     = "tf-test-l7-ilb-target-http-proxy%{random_suffix}"
  provider = google-beta
  region   = "us-west1"
  url_map  = google_compute_region_url_map.default.id
}

# URL map
resource "google_compute_region_url_map" "default" {
  name            = "tf-test-l7-ilb-regional-url-map%{random_suffix}"
  provider        = google-beta
  region          = "us-west1"
  default_service = google_compute_region_backend_service.default.id
}

# backend service
resource "google_compute_region_backend_service" "default" {
  name                  = "tf-test-l7-ilb-backend-subnet%{random_suffix}"
  provider              = google-beta
  region                = "us-west1"
  protocol              = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec           = 10
  health_checks         = [google_compute_region_health_check.default.id]
  backend {
    group           = google_compute_region_instance_group_manager.mig.instance_group
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 1.0
  }
}

# instance template
resource "google_compute_instance_template" "instance_template" {
  name         = "tf-test-l7-ilb-mig-template%{random_suffix}"
  provider     = google-beta
  machine_type = "e2-small"
  tags         = ["http-server"]

  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id
    access_config {
      # add external ip to fetch packages
    }
  }

  disk {
    source_image = "debian-cloud/debian-10"
    auto_delete  = true
    boot         = true
  }

  # install nginx and serve a simple web page
  metadata = {
    startup-script = <<-EOF1
      #! /bin/bash
      set -euo pipefail

      export DEBIAN_FRONTEND=noninteractive
      apt-get update
      apt-get install -y nginx-light jq

      NAME=$(curl -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/hostname")
      IP=$(curl -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip")
      METADATA=$(curl -f -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/attributes/?recursive=True" | jq 'del(.["startup-script"])')

      cat <<EOF > /var/www/html/index.html
      <pre>
      Name: $NAME
      IP: $IP
      Metadata: $METADATA
      </pre>
      EOF
    EOF1
  }

  lifecycle {
    create_before_destroy = true
  }
}

# health check
resource "google_compute_region_health_check" "default" {
  name     = "tf-test-l7-ilb-hc%{random_suffix}"
  provider = google-beta
  region   = "us-west1"

  http_health_check {
    port_specification = "USE_SERVING_PORT"
  }
}

# MIG
resource "google_compute_region_instance_group_manager" "mig" {
  name     = "tf-test-l7-ilb-mig1%{random_suffix}"
  provider = google-beta
  region   = "us-west1"

  base_instance_name = "vm"
  target_size        = 2

  version {
    instance_template = google_compute_instance_template.instance_template.id
    name              = "primary"
  }
}

# allow all access from IAP and health check ranges
resource "google_compute_firewall" "fw-iap" {
  name          = "tf-test-l7-ilb-fw-allow-iap-hc%{random_suffix}"
  provider      = google-beta
  direction     = "INGRESS"
  network       = google_compute_network.ilb_network.id
  source_ranges = ["130.211.0.0/22", "35.191.0.0/16", "35.235.240.0/20"]

  allow {
    protocol = "tcp"
  }
}

# allow http from proxy subnet to backends
resource "google_compute_firewall" "fw-ilb-to-backends" {
  name          = "tf-test-l7-ilb-fw-allow-ilb-to-backends%{random_suffix}"
  provider      = google-beta
  direction     = "INGRESS"
  network       = google_compute_network.ilb_network.id
  source_ranges = ["10.0.0.0/24"]
  target_tags   = ["http-server"]

  allow {
    protocol = "tcp"
    ports    = ["80", "443", "8080"]
  }
}

# test instance
resource "google_compute_instance" "vm-test" {
  name         = "tf-test-l7-ilb-test-vm%{random_suffix}"
  provider     = google-beta
  zone         = "us-west1-b"
  machine_type = "e2-small"

  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id
  }

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }
}

resource "google_network_services_lb_traffic_extension" "default" {
  provider = google-beta
  name     = "tf-test-l7-ilb-traffic-ext%{random_suffix}"
  description = "my traffic extension"
  location = "us-west1"

  load_balancing_scheme = "INTERNAL_MANAGED"
  forwarding_rules      = [google_compute_forwarding_rule.default.self_link]

  extension_chains {
      name = "chain1"

      match_condition {
          cel_expression = "request.host == 'example.com'"
      }

      extensions {
          name      = "ext11"
          authority = "ext11.com"
          service   = google_compute_region_backend_service.callouts_backend.self_link
          timeout   = "0.1s"
          fail_open = false

          supported_events = ["REQUEST_HEADERS"]
          forward_headers = ["custom-header"]
      }
  }

  labels = {
    foo = "bar"
  }
}

# Traffic Extension Backend Instance
resource "google_compute_instance" "callouts_instance" {
  provider = google-beta

  name = "tf-test-l7-ilb-callouts-ins%{random_suffix}"
  zone = "us-west1-a"

  machine_type = "e2-small"
  labels = {
    "container-vm" = "cos-stable-109-17800-147-54"
  }
  tags         = ["allow-ssh","load-balanced-backend"]

  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id
    access_config {
        # add external ip to fetch packages
    }
  }
  boot_disk {
    auto_delete  = true
    initialize_params {
      type = "pd-standard"
      size = 10
      image = "https://www.googleapis.com/compute/v1/projects/cos-cloud/global/images/cos-stable-109-17800-147-54"
    }
  }

  # Initialize an Envoy's Ext Proc gRPC API based on a docker container
  metadata = {
    gce-container-declaration = "# DISCLAIMER:\n# This container declaration format is not a public API and may change without\n# notice. Please use gcloud command-line tool or Google Cloud Console to run\n# Containers on Google Compute Engine.\n\nspec:\n  containers:\n  - image: us-docker.pkg.dev/service-extensions/ext-proc/service-callout-basic-example-python:latest\n    name: callouts-vm\n    securityContext:\n      privileged: false\n    stdin: false\n    tty: false\n    volumeMounts: []\n  restartPolicy: Always\n  volumes: []\n"
    google-logging-enabled = "true"
  }
  lifecycle {
    create_before_destroy = true
  }
}

// callouts instance group
resource "google_compute_instance_group" "callouts_instance_group" {
  provider    = google-beta
  name        = "tf-test-l7-ilb-callouts-ins-group%{random_suffix}"
  description = "Terraform test instance group"

  instances = [
    google_compute_instance.callouts_instance.id,
  ]

  named_port {
    name = "http"
    port = "80"
  }

  named_port {
    name = "grpc"
    port = "443"
  }

  zone = "us-west1-a"
}

# callout health check
resource "google_compute_region_health_check" "callouts_health_check" {
  provider = google-beta
  name     = "tf-test-l7-ilb-callouts-hc%{random_suffix}"
  region   = "us-west1"
  http_health_check {
    port = 80
  }
}

# callout backend service
resource "google_compute_region_backend_service" "callouts_backend" {
  provider              = google-beta
  name                  = "tf-test-l7-ilb-callouts-backend%{random_suffix}"
  region                = "us-west1"
  protocol              = "HTTP2"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec           = 10
  port_name             = "grpc"
  health_checks         = [google_compute_region_health_check.callouts_health_check.id]

  backend {
    group           = google_compute_instance_group.callouts_instance_group.id
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 1.0
  }
}
`, context)
}
