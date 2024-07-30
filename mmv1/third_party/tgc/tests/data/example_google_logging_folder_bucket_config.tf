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

resource "google_folder" "default" {
	display_name = "some-folder-name"
	parent = "organizations/{{.OrgID}}"
	deletion_protection = false
}
  
resource "google_logging_folder_bucket_config" "basic" {
	folder = google_folder.default.name
	location = "global"
	retention_days = 30
	bucket_id      = "_Default"
  
	index_configs {
		field_path = "jsonPayload.request.status"
		type = "INDEX_TYPE_STRING"
	}
}