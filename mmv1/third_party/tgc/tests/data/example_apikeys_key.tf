provider "google" {
  project     = "tf-deployer-2"
  region      = "us-central1"
}

resource "google_apikeys_key" "primary" {
  name         = "key"
  display_name = "sample-key"
  project      = "tf-deployer-2"

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
