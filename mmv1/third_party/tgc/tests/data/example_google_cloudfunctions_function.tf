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

resource "google_vpc_access_connector" "connector" {
  name          = "vpc-con-cf"
  ip_cidr_range = "10.8.0.0/28"
  network       = "default"
  region = "us-east1"
}

resource "google_cloudfunctions_function" "function" {
  name        = "my-cf"
  description = "My CloudFunction"
  runtime     = "nodejs14"

  available_memory_mb   = 128
  source_archive_bucket = "validator_bucket_local"
  source_archive_object = "sample.zip"
  trigger_http          = true
  timeout               = 60
  entry_point           = "helloGCS"
  labels = {
    my-cf-label-value = "my-cf-label-value"
  }

  ingress_settings = "ALLOW_INTERNAL_ONLY"
  vpc_connector = google_vpc_access_connector.connector.name
  vpc_connector_egress_settings = "PRIVATE_RANGES_ONLY"

  environment_variables = {
    MY_CF_ENV = "my-cf-env"
  }

  region = "us-east1"
}
