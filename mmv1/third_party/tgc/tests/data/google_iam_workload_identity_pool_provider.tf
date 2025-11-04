terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 4.54.0"
    }
  }
}

provider "google" {
  project = "{{.Provider.project}}"
}

resource "google_iam_workload_identity_pool" "gg_asset_44602_7df7" {
  project                   = "{{.Provider.project}}"
  workload_identity_pool_id = "gg-asset-44602-7df7"
  display_name              = "gg-asset-44602-7df7"
  description               = "Workload Identity Pool for gg-asset-44602-7df7"
}

resource "google_iam_workload_identity_pool_provider" "gg_asset_44602_7df7" {
  project                            = "{{.Provider.project}}"
  workload_identity_pool_id          = google_iam_workload_identity_pool.gg_asset_44602_7df7.workload_identity_pool_id
  workload_identity_pool_provider_id = "gg-asset-44602-7df7"
  display_name                       = "gg-asset-44602-7df7"
  description                        = "AWS provider for gg-asset-44602-7df7"
  aws {
    account_id = "111111111111"
  }
}

resource "google_iam_workload_identity_pool" "gg_asset_45050_bcd4" {
  project                    = "{{.Provider.project}}"
  workload_identity_pool_id  = "gg-asset-45050-bcd4"
  display_name               = "gg-asset-45050-bcd4"
  description                = "Workload Identity Pool for gg-asset-45050-bcd4"
  disabled                   = false
}

resource "google_iam_workload_identity_pool_provider" "gg_asset_45050_bcd4" {
  project                            = "{{.Provider.project}}"
  workload_identity_pool_id          = google_iam_workload_identity_pool.gg_asset_45050_bcd4.workload_identity_pool_id
  workload_identity_pool_provider_id = "gg-asset-45050-bcd4"
  display_name                       = "gg-asset-45050-bcd4"
  description                        = "OIDC provider for gg-asset-45050-bcd4"
  disabled                           = true
  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
  }
  attribute_condition = "assertion.repository_owner == 'google'"

  oidc {
    issuer_uri = "https://oidc.gg-asset-45050-bcd4.com"
  }
}
