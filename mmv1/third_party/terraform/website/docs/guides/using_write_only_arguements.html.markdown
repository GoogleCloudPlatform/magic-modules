---
page_title: "Use write-only arguments in the Google Cloud provider"
description: |-
  How to use write-only arguments in the Google Cloud provider
---

# Write-only arguments in the Google Cloud provider

The Google Cloud provider has introduced new write-only arguments for a more secure way to manage data. These arguments accept values from configuration and will not be stored in plan or state providing an additional layer of security and control over data access.

For more information, see the [official HashiCorp documentation for Write-only Arguments](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/write-only-arguments).

The Google Cloud provider has added the following write-only arguments:
- [`google_sql_user: password_wo`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_user#password-1)
- [`google_secret_manager_secret_version: secret_data_wo`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret_version#secret_data_wo)
- [`google_bigquery_data_transfer_config: sensitive_params.secret_access_key_wo`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/bigquery_data_transfer_config#secret_access_key_wo)

These were chosen due to them being marked as sensitive already in the provider. Although sensitive arguments do not appear in `terraform plan`, they are still stored in the Terraform state. Write-only arguments allow users to access and use data in their configurations without that data being stored in Terraform state.

## Use the Google Cloud provider's new write-only arguments

The following sections show how to use the new write-only argument in the Google Cloud provider.

### Applying a write-only argument

The following example shows how to apply a write-only argument. All write-only arguments are marked with the `wo` suffix and can not be used with the argument that it's mirroring. For example, `secret_data_wo` can not be used with `secret_data`.

```hcl
resource "google_sql_database_instance" "instance" {
  name                = "main-instance"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}
resource "google_sql_user" "user1" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "gmail.com"
  password_wo = "test_password"
  password_wo_version = 1
}
```

During `terraform plan` you will see that the write-only argument is marked appropriately:

```
  # google_sql_user.user1 will be created
  + resource "google_sql_user" "user1" {
      + host                    = "gmail.com"
      + id                      = (known after apply)
      + instance                = "main-instance"
      + name                    = "admin"
      + password_wo             = (write-only attribute)
      + password_wo_version     = 1
      + project                 = "hc-terraform-testing"
      + sql_server_user_details = (known after apply)
    }
```

Upon `terrform apply` you will see in `terraform.tfstate` that the write-only attribute from the configuration is not reflected in the state:

```hcl
...
      "mode": "managed",
      "type": "google_sql_user",
      "name": "user1",
      "provider": "provider[\"registry.terraform.io/hashicorp/google\"]",
      "instances": [
        {
          "schema_version": 1,
          "attributes": {
            "deletion_policy": null,
            "host": "gmail.com",
            "id": "admin/gmail.com/main-instance",
            "instance": "main-instance",
            "name": "admin",
            "password": null,
            "password_policy": [],
            "password_wo": null, // write-only attribute is not stored in state
            "password_wo_version": 1,
            "project": "hc-terraform-testing",
            "sql_server_user_details": [],
            "timeouts": null,
            "type": ""
          },
```

Any value that is set for a write-only argument is nulled out before the RPC response is sent to Terraform.

### Updating write-only arguments

Since write-only arguments are not stored in the Terraform state, they cannot be updated by just changing the value in the configuration due to the attribute being nulled out.

In order to update a write-only attribute we must change the write-only attribute's version.

```hcl
resource "google_sql_user" "user1" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "gmail.com"
  password_wo = "updated_password" // updated password
  password_wo_version = 2 // updated version
}
```

A `terraform apply` of this configuration will allow you to update the write-only argument despite the new value not being shown in the plan.

```hcl
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # google_sql_user.user1 will be updated in-place
  ~ resource "google_sql_user" "user1" {
        id                      = "admin/gmail.com/main-instance"
        name                    = "admin"
      ~ password_wo_version     = 1 -> 2
        # (6 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```
