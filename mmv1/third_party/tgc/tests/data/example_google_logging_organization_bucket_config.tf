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

resource "google_logging_organization_bucket_config" "basic" {
  organization = "12345"
  location = "global"
  retention_days = 30
  bucket_id = "_Default"

  index_configs {
    field_path = "jsonPayload.request.status"
    type = "INDEX_TYPE_STRING"
  }
}