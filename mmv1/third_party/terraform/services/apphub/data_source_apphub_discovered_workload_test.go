package apphub_test

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceApphubDiscoveredWorkload_basic(t *testing.T) {
    t.Parallel()

    context := map[string]interface{}{
    	"random_suffix": acctest.RandString(t, 10),
    }

    acctest.VcrTest(t, resource.TestCase{
        PreCheck:                 func() { acctest.AccTestPreCheck(t) },
        ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
        Steps: []resource.TestStep{
            {
                Config: testDataSourceApphubDiscoveredWorkload_basic(context),
            },
        },
    })
}

func testDataSourceApphubDiscoveredWorkload_basic(context map[string]interface{}) string {
    return acctest.Nprintf(`
data "google_apphub_discovered_workload" "catalog-workload" {
  provider = google
  location = "us-central1"
  workload_uri = "${replace(google_compute_region_instance_group_manager.mig.instance_group, "https://www.googleapis.com/compute/v1", "//compute.googleapis.com")}"
  depends_on = [google_compute_region_instance_group_manager.mig]
}

# VPC network
resource "google_compute_network" "ilb_network" {
  name                    = "l7-ilb-network-%{random_suffix}"
  project                 = "tf-project-5-416803"
  auto_create_subnetworks = false
}

# backend subnet
resource "google_compute_subnetwork" "ilb_subnet" {
  name          = "l7-ilb-subnetwork-%{random_suffix}"
  project       = "tf-project-5-416803"
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-central1"
  network       = google_compute_network.ilb_network.id
  depends_on = [google_compute_network.ilb_network]
}

# instance template
resource "google_compute_instance_template" "instance_template" {
  name         = "l7-ilb-mig-template-%{random_suffix}"
  project               = "tf-project-5-416803"
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
  depends_on = [google_compute_subnetwork.ilb_subnet, google_compute_network.ilb_network]
}

resource "google_compute_region_instance_group_manager" "mig" {
  name     = "l7-ilb-mig1-%{random_suffix}"
  project               = "tf-project-5-416803"
  region   = "us-central1"
  version {
    instance_template = google_compute_instance_template.instance_template.id
    name              = "primary"
  }
  base_instance_name = "vm"
  target_size        = 2
  depends_on = [google_compute_instance_template.instance_template]
}
`, context)
}

