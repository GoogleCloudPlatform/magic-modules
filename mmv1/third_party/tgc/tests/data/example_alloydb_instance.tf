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


resource "google_alloydb_cluster" "default" {
  cluster_id = "alloydb-cluster"
  location   = "us-central1"
  network_config {
    network = "default"
  }
  
  initial_user {
    password = "alloydb-cluster"
  }

  deletion_protection = false
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.cluster_id
  instance_id   = "alloydb-instance"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}
