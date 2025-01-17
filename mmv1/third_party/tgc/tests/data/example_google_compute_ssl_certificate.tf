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

resource "google_compute_ssl_certificate" "webserver_cert" {
  name         = "prod-webserver-cert"
  private_key  = base64encode("-----BEGIN RSA PRIVATE KEY...")
  certificate  = base64encode("-----BEGIN CERTIFICATE...")
  lifecycle {
    create_before_destroy = true
  }
}