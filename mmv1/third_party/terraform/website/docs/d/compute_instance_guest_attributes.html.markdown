---
subcategory: "Compute Engine"
description: |-
  Get GCE instance's guest attributes
---

# google_compute_instance_guest_attributes

Get information about a VM instance resource within GCE. For more information see
[the official documentation](https://cloud.google.com/compute/docs/instances)
and
[API](https://cloud.google.com/compute/docs/reference/latest/instances).

Get information about VM's guest attrubutes. For more information see [the official documentation]()
and
[API]()

## Example Usage - get all attributes from a single namespace

```hcl
data "google_compute_instance_guest_attributes" "appserver_ga" {
  name       = "primary-application-server"
  zone       = "us-central1-a"
  query_path = "variables/"
}
```

## Example Usage - get a specific variable

```hcl
data "google_compute_instance_guest_attributes" "appserver_ga" {
  name         = "primary-application-server"
  zone         = "us-central1-a"
  variable_key = "variables/key1"
}
```

## Argument Reference

The following arguments are supported:

* `self_link` - (Optional) The self link of the instance. One of `name` or `self_link` must be provided.

* `name` - (Optional) The name of the instance. One of `name` or `self_link` must be provided.

---

* `project` - (Optional) The ID of the project in which the resource belongs.
    If `self_link` is provided, this value is ignored.  If neither `self_link`
    nor `project` are provided, the provider project is used.

* `zone` - (Optional) The zone of the instance. If `self_link` is provided, this
    value is ignored.  If neither `self_link` nor `zone` are provided, the
    provider zone is used.

* `query_path` -

* `variable_key` -

## Attributes Reference

* `query_value` - Structure is [documented below](#nested_query_value).

* `variable_value` -

---

<a name="nested_query_value"></a>The `query_value` block supports:

* `key` -

* `namespace` -

* `value` -