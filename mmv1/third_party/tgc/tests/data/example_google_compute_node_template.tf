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

resource "google_compute_node_template" "template" {
  name      = "soletenant-tmpl"
  region    = "us-central1"
  node_type = "n1-node-96-624"
}