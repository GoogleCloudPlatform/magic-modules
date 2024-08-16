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
  location = "us-west1"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
  authority {
    issuer = "https://container.googleapis.com/v1/${google_container_cluster.primary.id}"
  }

  labels = {
    env = "test"
  }
}

resource "google_container_cluster" "primary" {
  name               = "basic-cluster"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection  = "true"
  network       = "default"
  subnetwork    = "default"
}