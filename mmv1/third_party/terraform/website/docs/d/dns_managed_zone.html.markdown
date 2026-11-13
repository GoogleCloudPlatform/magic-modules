---
subcategory: "Cloud DNS"
description: |-
  Provides access to the attributes of a zone within Google Cloud DNS
---

# google_dns_managed_zone

Provides access to a zone's attributes within Google Cloud DNS.
For more information see
[the official documentation](https://cloud.google.com/dns/zones/)
and
[API](https://cloud.google.com/dns/api/v1/managedZones).

```hcl
data "google_dns_managed_zone" "env_dns_zone" {
  name = "qa-zone"
}

resource "google_dns_record_set" "dns" {
  name = "my-address.${data.google_dns_managed_zone.env_dns_zone.dns_name}"
  type = "TXT"
  ttl  = 300

  managed_zone = data.google_dns_managed_zone.env_dns_zone.name

  rrdatas = ["test"]
}
```

## Argument Reference

* `name` - (Required) A unique name for the resource.

* `project` - (Optional) The ID of the project for the Google Cloud DNS zone.  If this is not provided the default project will be used.

## Attributes Reference

See [google_dns_managed_zone](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/dns_managed_zone) resource for details of all the available attributes.
