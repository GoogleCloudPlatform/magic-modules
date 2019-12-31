---
subcategory: "Cloud SQL"
layout: "google"
page_title: "Google: google_sql_ca_certs"
sidebar_current: "docs-google-datasource-sql-ca-certs"
description: |-
  Get all of the trusted Certificate Authorities (CAs) for the specified SQL database instance.
---

# google\_sql\_ca\_certs

Get all of the trusted Certificate Authorities (CAs) for the specified SQL database instance. For more information see the
[official documentation](https://cloud.google.com/sql/)
and
[API](https://cloud.google.com/sql/docs/mysql/admin-api/rest/v1beta4/instances/listServerCas).


## Example Usage

```hcl
data "google_sql_ca_certs" "ca_certs" {
  instance = "primary-database-server"
}
```

## Argument Reference

The following arguments are supported:

* `instance_self_link` - (Optional) The self link of the instance. One of `name` or `instance_self_link` must be provided.

* `instance` - (Optional) The name of the instance. One of `name` or `instance_self_link` must be provided.

---

* `project` - (Optional) The ID of the project in which the resource belongs.
    If `instance_self_link` is provided, this value is ignored.  If `project` is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `active_version` - The boot disk for the instance. Structure is documented below.

* `certs` - A list of server CA certificates for the instance. Each contains:
  * `cert` - The CA certificate used to connect to the SQL instance via SSL.
  * `common_name` - The CN valid for the CA cert.
  * `create_time` - Creation time of the CA cert.
  * `expiration_time` - Expiration time of the CA cert.
  * `sha1_fingerprint` - SHA1 fingerprint of the CA cert.
