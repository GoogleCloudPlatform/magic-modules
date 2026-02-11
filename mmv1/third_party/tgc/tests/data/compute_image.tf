terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 4.51.0"
    }
  }
}

provider "google" {
  project = "{{.Provider.project}}"
}

resource "google_compute_image" "gg-asset-42745-fa77" {
  name        = "gg-asset-42745-fa77"
  project     = "{{.Provider.project}}"
  source_disk = "https://www.googleapis.com/compute/v1/projects/{{.Provider.project}}/zones/us-central1-a/disks/gg-asset-source-disk-42745-fa77"
  family      = "gg-asset-family-42745-fa77"
  description = "Description for gg-asset-42745-fa77"
  labels = {
    "gg-asset-label-42745-fa77" = "gg-asset-value-42745-fa77"
  }
  licenses = [
    "https://www.googleapis.com/compute/v1/projects/vm-options/global/licenses/enable-vmx"
  ]
  guest_os_features {
    type = "UEFI_COMPATIBLE"
  }
  guest_os_features {
    type = "MULTI_IP_SUBNET"
  }
  storage_locations = [
    "us-central1",
  ]
  disk_size_gb = 20
}
