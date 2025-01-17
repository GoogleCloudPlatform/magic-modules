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


resource "google_kms_key_ring_import_job" "import-job" {
  key_ring = "keyring-id-example"
  import_job_id = "my-import-job"

  import_method = "RSA_OAEP_3072_SHA1_AES_256"
  protection_level = "SOFTWARE"
}