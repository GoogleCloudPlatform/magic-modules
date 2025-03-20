---
page_title: "Use external credentials in the Google Cloud provider with Terraform Stacks"
description: |-
  How to use external credentials in the Google Cloud provider with Terraform Stacks
---

# External Credentials in the Google Cloud provider with Terraform Stacks

Apart from using `access_token` and `credential` fields in the provider configuration, you can also use external credentials in the Google Cloud provider that are provided through a Workload Identity Federation (WIF) provider. This can be used to authenticate Terraform Stacks to provision resources in Google Cloud.

## Setting up a Workload Identity Federation (WIF) credentials

## Stacks Setup

A Terraform Stacks Project requires the following:

- A Workload Identity Federation (WIF) provider
- Components - `components.tfstacks.hcl`
- Deployment - `deployments.tfdeploy.hcl`

## Generating the Workload Identity Federation (WIF) credentials

```hcl
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.25.0"
    }
  }
}

provider "google" {
  region = "global"
}

variable "project_id" {
  type        = string
  description = "GCP Project ID"
}

# Create a service account for Terraform Stacks
resource "google_service_account" "terraform_stacks_sa" {
  account_id   = "terraform-stacks-sa"
  display_name = "Terraform Stacks Service Account"
  description  = "Service account used by Terraform Stacks for GCP resources"
}

# Create Workload Identity Pool
resource "google_iam_workload_identity_pool" "terraform_stacks_pool" {
  workload_identity_pool_id = "terraform-stacks-pool"
  display_name              = "Terraform Stacks Pool"
  description               = "Identity pool for Terraform Stacks authentication"
}

# Create Workload Identity Pool Provider
resource "google_iam_workload_identity_pool_provider" "terraform_stacks_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.terraform_stacks_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "terraform-stacks-provider"
  display_name                       = "Terraform Stacks Provider"
  description                        = "OIDC identity pool provider for Terraform Stacks"
  
  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
    "attribute.aud"        = "assertion.aud"
  }

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
    allowed_audiences = ["https://iam.googleapis.com/${google_iam_workload_identity_pool.terraform_stacks_pool.name}"]
  }
}

# Allow the Workload Identity Pool to impersonate the service account
resource "google_service_account_iam_binding" "workload_identity_binding" {
  service_account_id = google_service_account.terraform_stacks_sa.name
  role               = "roles/iam.workloadIdentityUser"
  
  members = [
    "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.terraform_stacks_pool.name}/*"
  ]
}

# Grant Storage Admin role to the service account (for bucket operations)
resource "google_project_iam_member" "sa_storage_admin" {
  project = var.project_id
  role    = "roles/storage.admin"
  member  = "serviceAccount:${google_service_account.terraform_stacks_sa.email}"
}

# Outputs to be used by Terraform Stacks
output "service_account_email" {
  value       = google_service_account.terraform_stacks_sa.email
  description = "Email of the service account to be used by Terraform Stacks"
}

output "audience" {
  value       = "https://iam.googleapis.com/${google_iam_workload_identity_pool.terraform_stacks_pool.name}"
  description = "The audience value to use when generating OIDC tokens"
}
```

## Terraform Stacks Setup with External Credentials

`deployments.tfdeploy.hcl`
```hcl
identity_token "jwt" {
  audience = ["hcp.workload.identity"]
}

deployment "staging" {
  inputs = {
    jwt = identity_token.jwt.jwt
  }
}
```

`components.tfstacks.hcl`
```hcl
required_providers {
  google = {
    source = "hashicorp/google"
    version = "6.25.0"
  }
}

provider "google" "this" {
  external_credentials {
    audience = "//iam.googleapis.com/projects/871647908372/locations/global/workloadIdentityPools/stacks-oidc-myz3/providers/stacks-oidc-myz3"
    service_account_email = "stacks-oidc-myz3@hc-terraform-testing.iam.gserviceaccount.com"
    identity_token_file = "./identity_token"
  }
}

variable "jwt" {
  type = string
}

component "storage_buckets" {
    source = "./buckets"

    inputs = {
        jwt = var.jwt
    }

    providers = {
        google    = provider.google.this
    }
}
```