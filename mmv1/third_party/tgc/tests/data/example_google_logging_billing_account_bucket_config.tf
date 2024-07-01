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

resource "google_logging_billing_account_bucket_config" "basic" {
  billing_account = "{{.Project.BillingAccountName}}"
  location        = "global"
  retention_days  = 30
  bucket_id       = "_Default"
}