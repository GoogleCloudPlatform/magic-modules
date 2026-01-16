---
subcategory: "Cloud Platform"
description: |-
  Allows management of a Google Cloud Platform service account Key
---

# google_service_account_key

Creates and manages service account keys, which allow the use of a service account with Google Cloud.

-> **Warning**: This resource persists a sensitive credential in plaintext in the [remote state](https://www.terraform.io/language/state/sensitive-data) used by Terraform.
Please take appropriate measures to protect your remote state.

* [API documentation](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts.keys)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/iam/docs/creating-managing-service-account-keys)


## Example Usage, creating a new Key

```hcl
resource "google_service_account" "myaccount" {
  account_id   = "myaccount"
  display_name = "My Service Account"
}

resource "google_service_account_key" "mykey" {
  service_account_id = google_service_account.myaccount.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}
```

## Example Usage, creating and regularly rotating a key

```hcl
resource "google_service_account" "myaccount" {
  account_id   = "myaccount"
  display_name = "My Service Account"
}

# note this requires the terraform to be run regularly
resource "time_rotating" "mykey_rotation" {
  rotation_days = 30
}

resource "google_service_account_key" "mykey" {
  service_account_id = google_service_account.myaccount.name

  keepers = {
    rotation_time = time_rotating.mykey_rotation.rotation_rfc3339
  }
}
```

## Example Usage, uploading a user-managed key with write-only attribute

This example shows how to upload a user-managed public key without storing sensitive data in Terraform state.
This follows [GCP's recommended approach](https://cloud.google.com/iam/docs/best-practices-for-managing-service-account-keys)
for high-security environments where you create your own key pair and only upload the public key to GCP.

```hcl
resource "google_service_account" "myaccount" {
  account_id   = "myaccount"
  display_name = "My Service Account"
}

resource "tls_private_key" "sa_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "sa_cert" {
  private_key_pem = tls_private_key.sa_key.private_key_pem

  subject {
    common_name = "myaccount"
  }

  validity_period_hours = 87600 # 10 years

  allowed_uses = [
    "digital_signature",
  ]
}

resource "google_service_account_key" "mykey" {
  service_account_id         = google_service_account.myaccount.name
  public_key_data_wo         = base64encode(tls_self_signed_cert.sa_cert.cert_pem)
  public_key_data_wo_version = 1 # Increment to rotate the key
}

# Store the private key securely in Secret Manager
resource "google_secret_manager_secret" "sa_private_key" {
  secret_id = "my-sa-private-key"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "sa_private_key" {
  secret      = google_secret_manager_secret.sa_private_key.id
  secret_data = tls_private_key.sa_key.private_key_pem
}
```

## Example Usage, save key in Kubernetes secret - DEPRECATED

```hcl
# Workload Identity is the recommended way of accessing Google Cloud APIs from pods.
# https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity

resource "google_service_account" "myaccount" {
  account_id   = "myaccount"
  display_name = "My Service Account"
}

resource "google_service_account_key" "mykey" {
  service_account_id = google_service_account.myaccount.name
}

resource "kubernetes_secret" "google-application-credentials" {
  metadata {
    name = "google-application-credentials"
  }
  data = {
    "credentials.json" = base64decode(google_service_account_key.mykey.private_key)
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_account_id` - (Required) The Service account id of the Key. This can be a string in the format
`{ACCOUNT}` or `projects/{PROJECT_ID}/serviceAccounts/{ACCOUNT}`. If the `{ACCOUNT}`-only syntax is used, either
the **full** email address of the service account or its name can be specified as a value, in which case the project will
automatically be inferred from the account. Otherwise, if the `projects/{PROJECT_ID}/serviceAccounts/{ACCOUNT}`
syntax is used, the `{ACCOUNT}` specified can be the full email address of the service account or the service account's
unique id. Substituting `-` as a wildcard for the `{PROJECT_ID}` will infer the project from the account.

* `key_algorithm` - (Optional) The algorithm used to generate the key. KEY_ALG_RSA_2048 is the default algorithm.
Valid values are listed at
[ServiceAccountPrivateKeyType](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts.keys#ServiceAccountKeyAlgorithm)
(only used on create)

* `public_key_type` (Optional) The output format of the public key requested. TYPE_X509_PEM_FILE is the default output format.

* `private_key_type` (Optional) The output format of the private key. TYPE_GOOGLE_CREDENTIALS_FILE is the default output format.

* `public_key_data` (Optional) Public key data to create a service account key for given service account. The expected format for this field is a base64 encoded X509_PEM and it conflicts with `key_algorithm`, `private_key_type`, and `public_key_data_wo`.

* `public_key_data_wo` (Optional) Write-only version of `public_key_data`. Public key data to create a service account key for given service account, without storing the value in Terraform state. The expected format for this field is a base64 encoded X509_PEM. Conflicts with `key_algorithm`, `private_key_type`, and `public_key_data`. Must be used with `public_key_data_wo_version`. Requires Terraform 1.11+.

* `public_key_data_wo_version` (Optional) Version number for `public_key_data_wo`. Increment this value to trigger recreation of the service account key with the write-only public key data. Must be used with `public_key_data_wo`.

* `keepers` (Optional) Arbitrary map of values that, when changed, will trigger a new key to be generated.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - an identifier for the resource with format `projects/{{project}}/serviceAccounts/{{account}}/keys/{{key}}`

* `name` - The name used for this key pair

* `public_key` - The public key, base64 encoded

* `private_key` - The private key in JSON format, base64 encoded. This is what you normally get as a file when creating
service account keys through the CLI or web console. This is only populated when creating a new key.

* `valid_after` - The key can be used after this timestamp. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".

* `valid_before` - The key can be used before this timestamp.
A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".

## Import

This resource does not support import.
