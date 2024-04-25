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

resource "google_compute_route" "my_route" {
  name         = "my-route"
  dest_range   = "10.1.0.0/16"
  next_hop_ip  = "10.0.0.1"
  network      = "my-network"
}