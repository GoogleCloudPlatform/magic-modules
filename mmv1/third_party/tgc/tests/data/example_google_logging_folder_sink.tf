terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 4.54.0"
    }
  }
}

provider "google" {
  project = "{{.Provider.project}}"
}

resource "google_storage_bucket" "test_bucket" {
  name     = "tf-test-bucket-{{.Project.Number}}"
  location = "US"
  project  = "{{.Provider.project}}"
}

resource "google_logging_folder_sink" "test_sink" {
  name        = "tf-test-sink"
  folder      = "folders/{{.FolderID}}"
  destination = "storage.googleapis.com/${google_storage_bucket.test_bucket.name}"
  filter      = "severity >= ERROR"
  include_children = true

  exclusions {
    name        = "exclude-gce-activity"
    description = "Exclude GCE activity logs."
    filter      = "logName:\"logs/compute.googleapis.com%2Factivity_log\""
  }

  bigquery_options {
    use_partitioned_tables = true
  }
}
