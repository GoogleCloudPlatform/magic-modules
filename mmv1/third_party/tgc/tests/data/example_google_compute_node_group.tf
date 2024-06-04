terraform {
  required_providers {
    google = {
      source = "hashicorp/google-beta"
      version = "~> {{.Provider.version}}"
    }
  }
}

provider "google" {
  {{if .Provider.credentials }}credentials = "{{.Provider.credentials}}"{{end}}
}

resource "google_compute_node_template" "soletenant-tmpl" {
  name      = "soletenant-tmpl"
  region    = "us-central1"
  node_type = "n1-node-96-624"
}

resource "google_compute_node_group" "nodes" {
  project     = "{{.Provider.project}}"
  provider    = google-beta
  name        = "soletenant-group"
  zone        = "us-central1-f"
  description = "example google_compute_node_group for Terraform Google Provider"
  maintenance_policy = "RESTART_IN_PLACE"
  maintenance_window {
    start_time = "08:00"
  }

  initial_size          = 1
  node_template = google_compute_node_template.soletenant-tmpl.id
  autoscaling_policy {
    mode      = "ONLY_SCALE_OUT"
    min_nodes = 1
    max_nodes = 10
  }
}