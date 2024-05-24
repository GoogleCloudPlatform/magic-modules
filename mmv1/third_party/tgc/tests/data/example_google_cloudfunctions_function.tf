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

resource "google_cloudfunctions_function" "function" {
  name        = "function-test"
  description = "My function"
  runtime     = "nodejs14"

  available_memory_mb   = 128
  source_archive_bucket = "validator_bucket_local"
  source_archive_object = "sample.zip"
  trigger_http          = true
  timeout               = 60
  entry_point           = "helloGCS"
  labels = {
    my-label = "my-label-value"
  }
  vpc_connector = "projects/tf-deployer-2"
  vpc_connector_egress_settings = "ALL_TRAFFIC"

  environment_variables = {
    MY_ENV_VAR = "my-env-var-value"
  }

  region = "australia-southeast1"
}
