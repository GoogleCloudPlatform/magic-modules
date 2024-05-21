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

resource "google_gke_hub_membership" "membership" {
  membership_id = "basic"
  location = "us-central1-a"
  endpoint {
    gke_cluster {
      resource_link = ""
    }
  }
}