---
subcategory: "Cloud Identity"
layout: "google"
page_title: "Google: google_cloud_identity_list_policies"
sidebar_current: "docs-google-data-cloud-identity-list-policies"
description: |-
Use this data source to list Cloud Identity policies.
---

# google_cloud_identity_list_policies

Use this data source to list Cloud Identity policies.

## Example Usage

```hcl
data "google_cloud_identity_list_policies" "all" {
    # Example filter (optional)"
    # filter = "customer == \"customers/my_customer\" &&
    setting.type.mathces('^settings/gmail\\..*$')"
}

output "first_policy_name" {
    value = data.google_cloud_identity_list_policies.all.policies[0].name
}

output "first_policy_customer" {
    value = data.google_cloud_identity_list_policies.all.policies[0].customer
    }

