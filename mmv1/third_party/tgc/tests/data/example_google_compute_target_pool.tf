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

resource "google_compute_target_pool" "default" {
  name = "instance-target-pool"
  region="us-central1"
  instances = [
    "us-central1-a/myinstance1",
    "us-central1-b/myinstance2",
  ]
}