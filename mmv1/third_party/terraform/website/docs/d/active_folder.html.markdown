---
subcategory: "Cloud Platform"
description: |-
  Get a folder within GCP.
---

# google\_active\_folder

Get an active folder within GCP by `display_name` and `parent`.

## Example Usage

```tf
data "google_active_folder" "department1" {
  display_name = "Department 1"
  parent       = "organizations/1234567"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The folder's display name.

* `parent` - (Required) The resource name of the parent Folder or Organization.

* `rest_method` - (Optional) The method to use to search for the folder. Both `LIST` and `SEARCH` are supported. The default value is `LIST` as SEARCH is [eventually consistent](https://cloud.google.com/resource-manager/reference/rest/v3/folders/search#:~:text=eventually%20consistent) while `LIST` is [strongly consistent](https://cloud.google.com/resource-manager/reference/rest/v3/folders/list#:~:text=list()%20provides%20a-,strongly%20consistent,-view%20of%20the).

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `name` - The resource name of the Folder. This uniquely identifies the folder.
