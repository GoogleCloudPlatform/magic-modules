---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_billing_subaccount"
sidebar_current: "docs-google-billing-subaccount"
description: |-
 Allows management of a Google Cloud Billing Subaccount.
---

# google\_billing\_subaccount

Allows creation and management of a Google Cloud Billing Subaccount.

!> **WARNING:** Deleting this Terraform resource will not delete or close the billing subaccount.

```hcl
resource "google_billing_subaccount" "subaccount" {
    display_name = "My Billing Account"
    master_billing_account = "012345-567890-ABCDEF"
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available billing accounts.
The given filters must match exactly one billing account whose data will be exported as attributes.
The following arguments are supported:

* `display_name` (Required) - The display name of the billing account.
* `master_billing_account` (Required) - The name of the master billing account that the subaccount
  will be created under in the form `{billing_account_id}` or `billingAccounts/{billing_account_id}`.
* `rename_on_destroy` (Optional) - If `true` the billing account display_name will be changed to
  "Terraform Destroyed" along with a timestamp.  If `false` this will not occur.  Default is `false`.

## Attributes Reference

The following additional attributes are exported:

* `open` - `true` if the billing account is open, `false` if the billing account is closed.
* `name` - The resource name of the billing account in the form `billingAccounts/{billing_account_id}`.
* `billing_account_id` - The billing account id.
