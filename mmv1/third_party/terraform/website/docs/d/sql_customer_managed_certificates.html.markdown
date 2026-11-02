---
subcategory: "Cloud SQL"
description: |-
  List customer-managed certificate bundles configured on a Cloud SQL instance.
---

# google_sql_customer_managed_certificates

Use this data source to list customer-managed SSL/TLS certificate bundles configured on a Cloud SQL instance.

## Example Usage

```hcl
data "google_sql_customer_managed_certificates" "all" {
  instance = google_sql_database_instance.default.name
}

output "certificate_bundles" {
  value = data.google_sql_customer_managed_certificates.all.customer_managed_certificates
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (Required) The name of the Cloud SQL instance.

* `project` - (Optional) The ID of the project in which the instance belongs. If it is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `id` - An identifier for the data source with format `projects/{{project}}/instances/{{instance}}/customerManagedCertificates`.

* `customer_managed_certificates` - A list of customer-managed certificate bundles configured on the instance. Each entry contains:
  * `bundle` - Customer managed certificate bundle name.
  * `intent` - Intent for the customer managed certificate bundle configuration.
  * `primary_dns` - DNS name associated with the primary instance in a replication topology.
  * `primary_ipv4_address` - IPv4 address for the primary DNS name.
  * `file_paths` - File paths of the PEM-encoded certificates and key on the instance.
  * `server_ca_cert_file_path` - File path of the PEM-encoded peer server CA certificate on the instance.
  * `client_cert_file_path` - File path of the PEM-encoded peer client certificate on the instance.
  * `client_key_file_path` - File path of the PEM-encoded peer client private key on the instance.
