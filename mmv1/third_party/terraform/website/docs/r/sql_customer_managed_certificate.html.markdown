---
subcategory: "Cloud SQL"
description: |-
  Creates and manages a customer-managed certificate bundle on a Cloud SQL instance.
---

# google_sql_customer_managed_certificate

Creates and manages a customer-managed SSL/TLS certificate bundle on a Cloud SQL instance (for example, to configure outbound mTLS for PostgreSQL logical replication (`PG_LOGICAL`)).

For more information, see the [Cloud SQL documentation](https://cloud.google.com/sql/docs/postgres/replication/configure-logical-replication).

## Example Usage - SQL Customer Managed Certificate

```hcl
resource "google_sql_database_instance" "default" {
  name                = "pg-logical-subscriber"
  region              = "us-central1"
  database_version    = "POSTGRES_15"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_customer_managed_certificate" "bundle1" {
  bundle               = "bundle1"
  instance             = google_sql_database_instance.default.name
  intent               = "PG_LOGICAL"
  primary_dns          = "primary.example.com"
  primary_ipv4_address = "10.0.0.1"
  server_ca_cert       = file("path/to/server-ca.pem")
  client_cert          = file("path/to/client-cert.pem")
  client_key           = file("path/to/client-key.pem")
}
```

## Argument Reference

The following arguments are supported:

* `bundle` - (Required) Customer managed certificate bundle name. Constrained to 1-10 alphanumeric characters (`[a-zA-Z0-9]{1,10}`). Changing this forces a new resource to be created.

* `instance` - (Required) The name of the Cloud SQL instance. This does not include the project ID. Changing this forces a new resource to be created.

* `intent` - (Optional) Intent for the customer managed certificate bundle configuration. Possible values are `"PG_LOGICAL"` and `"NATIVE_REPLICATION"`. Defaults to `"PG_LOGICAL"`.

* `server_ca_cert` - (Optional) Input only. The PEM-encoded peer server CA certificate content. At least one of `server_ca_cert`, `client_cert`, or `client_key` must be specified.

* `client_cert` - (Optional) Input only. The PEM-encoded peer client certificate content. At least one of `server_ca_cert`, `client_cert`, or `client_key` must be specified.

* `client_key` - (Optional) Input only. The PEM-encoded peer client private key content. At least one of `server_ca_cert`, `client_cert`, or `client_key` must be specified.

* `primary_dns` - (Optional) DNS name associated with the primary instance in a replication topology.

* `primary_ipv4_address` - (Optional) IPv4 address for the primary DNS name.

* `deletion_policy` - (Optional) The deletion policy for the customer-managed certificate bundle. Setting `"ABANDON"` allows the resource to be abandoned rather than deleted. Setting `"PREVENT"` prevents the resource from being deleted. Possible values are: `"PREVENT"`, `"ABANDON"`, `"DELETE"`. Defaults to `"DELETE"`.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - An identifier for the resource with format `projects/{{project}}/instances/{{instance}}/customerManagedCertificates/{{bundle}}`.

* `file_paths` - Output only. File paths of the PEM-encoded certificates and key on the instance.

* `server_ca_cert_file_path` - Output only. File path of the PEM-encoded peer server CA certificate on the instance.

* `client_cert_file_path` - Output only. File path of the PEM-encoded peer client certificate on the instance.

* `client_key_file_path` - Output only. File path of the PEM-encoded peer client private key on the instance.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 10 minutes.
- `update` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

Customer-managed certificate bundles can be imported using any of these accepted formats:

* `projects/{{project}}/instances/{{instance}}/customerManagedCertificates/{{bundle}}`
* `instances/{{instance}}/customerManagedCertificates/{{bundle}}`
* `{{project}}/{{instance}}/{{bundle}}`
* `{{instance}}/{{bundle}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import a customer-managed certificate bundle using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/instances/{{instance}}/customerManagedCertificates/{{bundle}}"
  to = google_sql_customer_managed_certificate.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), import using one of the formats above. For example:

```
$ terraform import google_sql_customer_managed_certificate.default projects/{{project}}/instances/{{instance}}/customerManagedCertificates/{{bundle}}
```
