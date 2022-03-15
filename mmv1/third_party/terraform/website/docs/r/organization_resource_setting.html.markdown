---
subcategory: "ResourceSettings"
layout: "google"
page_title: "Google: google_organization_resource_setting"
sidebar_current: "docs-google-organization-resource-setting"
description: |-
  Manages an Organization Resource Setting on Google Cloud.
---

# google_organization_resource_setting

Manages Resource Setting at the Organization level. For more information see
[the official documentation](https://cloud.google.com/resource-manager/docs/resource-settings/overview) and
[list of available settings](https://cloud.google.com/resource-manager/docs/resource-settings/manage-resource-settings).

~> **Note:** This resource can have side effects. Because Resource Settings cannot be deleted, the deletion of the Terraform resource will unset the localValue of the Setting. 

## Example Usage

```hcl
resource "google_organization_resource_setting" "keys" {
  organization_id             = "my-org"
  setting_name                = "iam-serviceAccountKeyExpiry"

  local_value {
    string_value = "1hours"
  }
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required) The organization the setting will apply to.
    Changing this forces a new resource to be created.

* `setting_name` - (Required) A unique ID for the resource.
    Changing this forces a new resource to be created.

* `local_value` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

<a name="nested_local_value"></a>The `local_value` block supports:

* `boolean_value` - (Optional) - Holds the value for the local value field with bool type.

* `string_value` - (Optional) - Holds the value for a local value field with string type.

* `enum_value` - (Optional) - The display name of the enum value.

* `duration_value` - (Optional) - Defines this value as being a Duration.
    A duration in seconds with up to nine fractional digits, terminated by 's'.
    Example: "3.5s".

