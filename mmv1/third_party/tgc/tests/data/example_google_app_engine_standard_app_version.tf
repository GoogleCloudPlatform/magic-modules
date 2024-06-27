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

resource "google_app_engine_standard_app_version" "my_app_v1" {
  version_id = "v1"
  service    = "default"
  runtime    = "python39" 

  entrypoint {
    shell = "python3 world.py"
  }

  deployment {
    zip {
      source_url = "https://storage.googleapis.com/bucket-app-engine/world.zip"
    }
  }
}