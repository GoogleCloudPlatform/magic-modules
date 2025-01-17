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


resource "google_cloud_tasks_queue" "default" {
  name = "cloud-tasks-queue-test"
  location = "us-central1"
  rate_limits {
   max_dispatches_per_second = 10  
  }
}
