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
  name        = "soletenant-group"
  zone        = "us-central1-f"
  description = "example google_compute_node_group for Terraform Google Provider"

  initial_size          = 1
  node_template = google_compute_node_template.soletenant-tmpl.id
}