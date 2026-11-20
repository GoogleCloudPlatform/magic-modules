data "google_project" "project" {
}

resource "google_project_service" "observability_api" {
  project            = data.google_project.project.project_id
  service            = "observability.googleapis.com"
  disable_on_destroy = false
}

# Wait for the project to be created and recognized by the Observability API
resource "time_sleep" "wait_for_project_propagation" {
  create_duration = "90s"
  depends_on      = [google_project_service.observability_api]
}

resource "google_observability_project_settings" "global" {
  location                 = "global"
  project                  = data.google_project.project.project_id
  default_storage_location = "us"
  depends_on = [time_sleep.wait_for_project_propagation]
}

# Buffer for the new identity to propagate to global IAM indexes
resource "time_sleep" "wait_for_sa_propagation" {
  create_duration = "30s"
  depends_on      = [
    google_observability_project_settings.global
  ]
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-observability.iam.gserviceaccount.com"

  depends_on = [
    time_sleep.wait_for_sa_propagation
  ]
}

resource "google_observability_bucket" "primary" {
  bucket_id    = "_Trace"
  location     = "us"
  project      = data.google_project.project.project_id
  display_name = "%{display_name}"
  description  = "%{description}"

  cmek_settings {
    kms_key = "%{kms_key}"
  }

  depends_on = [
    google_kms_crypto_key_iam_member.crypto_key
  ]
}
