---
subcategory: "Cloud DNS"
description: |-
  List Google Cloud DNS record sets in a managed zone for use with terraform query
  and .tfquery.hcl files.
---

# google_dns_record_set (list)

Lists Cloud DNS **record sets** in a managed zone for use with
[`terraform query`](https://developer.hashicorp.com/terraform/cli/commands/query) and
**`.tfquery.hcl`** files. Results correspond to existing
[`google_dns_record_set`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/dns_record_set)
managed resources.

For how list resources work in this provider, file layout, Terraform version requirements, and
shared `list` block arguments, refer to the guide
[Use list resources with terraform query (Google Cloud provider)](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_list_resources_with_terraform_query).

## Example

```hcl
list "google_dns_record_set" "all_in_zone" {
  provider = google

  config {
    managed_zone = "my-managed-zone"

    # Optional. Defaults to the provider project when omitted.
    # project = "other-project"

    # Optional filters.
    # name = "www.example.com."
    # type = "A"
  }
}
```

Run `terraform query` from the directory that contains the `.tfquery.hcl` file.

## Configuration (`config` block)

* `managed_zone` - (Required) Managed zone name to list record sets from.
* `project` - (Optional) Project ID to list record sets from. If unset, the provider's
  configured default project is used (same idea as the managed resource).
* `name` - (Optional) Filter results to a specific DNS record name.
* `type` - (Optional) Filter results to a specific DNS record type.

## Results

By default each result includes **resource identity** for `google_dns_record_set` (see
[Resource identity](https://developer.hashicorp.com/terraform/language/resources/identities)):

* `project` - Project ID when applicable.
* `managed_zone` - Managed zone name.
* `name` - DNS record name.
* `type` - DNS record type.

With `include_resource = true` on the `list` block, results also include the full resource-style
attributes documented for the managed
[`google_dns_record_set` resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/dns_record_set#attributes-reference)
(for example `ttl`, `rrdatas`, and `routing_policy` where present in state).
