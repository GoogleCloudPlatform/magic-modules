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

resource "google_compute_region_commitment" "foobar" {
  name = "my-full-commitment"
  description = "some description"
  plan = "THIRTY_SIX_MONTH"
  type = "MEMORY_OPTIMIZED"
  category = "MACHINE"
  auto_renew = true
  region  = "us-east1"
  resources {
      type = "VCPU"
      amount = "4"
  }
  resources {
      type = "MEMORY"
      amount = "9"
  }
}