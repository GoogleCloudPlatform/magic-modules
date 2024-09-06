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


resource "google_apikeys_key" "primary" {
  name         = "key"
  display_name = "sample-key"
  project      = "{{.Provider.project}}"

  restrictions {
    android_key_restrictions {
      allowed_applications {
        package_name     = "com.example.app123"
        sha1_fingerprint = "1699466a142d4682a5f91b50fdf400f2358e2b0b"
      }
    }

    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET"]
    }
  }
}
