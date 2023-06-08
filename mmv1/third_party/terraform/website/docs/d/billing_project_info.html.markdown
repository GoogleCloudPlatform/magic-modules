---
subcategory: "Cloud Billing"
description: |-
  Get billing information about a Google Cloud project
---

# google\_billing\_project\_info

Use this data source to get the billing account linked to Google Cloud project.

```hcl
# Get folder by id
data "google_billing_project" "default" {
  project = "my-project-id"
}

output "my_project_billing_account" {
  value = data.google_billing_project.default.billing_account
}
```

## Argument Reference

* `project` - (Optional) The project ID. If it is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `billing_account` - The ID of the linked billing account.

