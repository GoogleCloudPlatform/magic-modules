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

resource "google_logging_project_bucket_config" "basic" {
    project    = "{{.Provider.project}}"
    location  = "global"
    retention_days = 30
    bucket_id = "_Default"
}