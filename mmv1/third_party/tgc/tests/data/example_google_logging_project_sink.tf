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

resource "google_logging_project_sink" "my-sink" {
  name   = "my-sink"
  project = "{{.Provider.project}}"
  destination = "bigquery.googleapis.com/projects/{{.Provider.project}}/datasets/my_dataset"
  filter = "severity >= ERROR"
  description = "A sink for errors"
  disabled = false

  exclusions {
    name        = "exclude-debug"
    description = "Exclude debug logs"
    filter      = "severity < INFO"
    disabled    = false
  }

  bigquery_options {
    use_partitioned_tables = true
  }
}
