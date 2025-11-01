---
subcategory: "Cloud Identity"
layout: "google"
page_title: "Google: google_cloud_identity_policy"
sidebar_current: "docs-google-data-cloud-identity-policy"
description: |-
  Use this data source to retrieve a Cloud Identity policy.
---

# google_cloud_identity_policy

Use this data source to retrieve a Cloud Identity policy.

## Example Usage

```hcl
data "google_cloud_identity_policy" "test" {
  name = "policies/{policy_id}"
}

// The customer the policy belongs to
output "policy_customer" {
  value = data.google_cloud_identity_policy.test.customer
}

// The CEL query of the policy
output "policy_query_query" {
  value = data.google_cloud_identity_policy.test.policy_query[0].query
}

// The org unit the policy applies to
output "policy_query_org_unit" {
  value = data.google_cloud_identity_policy.test.policy_query[0].org_unit
}

// The group the policy applies to
output "policy_query_group" {
  value = data.google_cloud_identity_policy.test.policy_query[0].group
}

// The sort order of the policy
output "policy_query_sort_order" {
  value = data.google_cloud_identity_policy.test.policy_query[0].sort_order
}

// The setting of the policy as a JSON string
output "policy_setting" {
  value = data.google_cloud_identity_policy.test.setting
}
```

## Argument Reference

The following arguments are supported:

*   `name` - (Required) The resource name of the policy to retrieve. Format: `policies/{policy_id}`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

*   `name` - The resource name of the policy.

*   `customer` - The customer that the policy belongs to.

*   `policy_query` - A list containing the CEL query that defines which entities the policy applies to. Structure is documented below.

*   `setting` - The setting configured by this policy, represented as a JSON string.

*   `type` - The type of the policy.

---

The `policy_query` block contains:

*   `query` - The query that defines which entities the policy applies to.

*   `group` - The group that the policy applies to.

*   `org_unit` - The org unit that the policy applies to.

*   `sort_order` - The sort order of the policy.
