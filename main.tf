#terraform {
#  required_providers {
#    google-beta = {
#      version = "~> 4.0.0"
#    }
#  }
#}
resource "google_service_directory_namespace" "example" {
  provider     = google-beta
  project      = "lcaggioni-sandbox"
  namespace_id = "namespace"
  location     = "us-central1"
}

resource "google_service_directory_service" "example" {
  provider   = google-beta
  service_id = "example-service"
  namespace  = google_service_directory_namespace.example.id
}

resource "google_service_directory_endpoint" "example" {
  provider    = google-beta
  endpoint_id = "example-endpointi-2"
  service     = google_service_directory_service.example.id

  metadata = {
    stage  = "prod"
    region = "us-central1"
  }

  address = "1.2.3.44"
  port    = 5353
  network = "projects/629017717406/locations/global/networks/to-onprem"
}
