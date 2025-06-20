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

resource "google_datastream_connection_profile" "default" {
    display_name          = "Connection profile"
    location              = "us-central1"
    connection_profile_id = "my-profile"

    gcs_profile {
        bucket    = "my-bucket"
        root_path = "/path"
    }
}