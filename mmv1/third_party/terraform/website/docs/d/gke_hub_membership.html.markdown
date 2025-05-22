---
subcategory: "GKEHub"
description: |-
  Retrieves the details of a GKE Hub Membership.
---

# `google_gke_hub_membership`

Retrieves the details of a specific GKE Hub Membership. Use this data source to retrieve the membership's configuration and state.

## Example Usage

```hcl
data "google_gke_hub_membership" "example" {
  project       = "my-project-id"
  location      = "global"
  membership_id = "my-membership-id" # GKE Cluster's name
}
```

## Argument Reference

The following arguments are supported:

* `membership_id` - (Required) The GKE Hub Membership id or GKE Cluster's name.

* `location` - (Required) The location for the GKE Hub Membership.
    Currently only `global` is supported.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

See [google_gke_hub_membership](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/gke_hub_membership) resource for details of the available attributes.