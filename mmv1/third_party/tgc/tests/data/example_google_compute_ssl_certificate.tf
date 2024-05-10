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
  private_key = file("/etc/apache2/ssl/apache.key")
  certificate = file("/etc/apache2/ssl/apache.crt")
  lifecycle {
    create_before_destroy = true
  }
}