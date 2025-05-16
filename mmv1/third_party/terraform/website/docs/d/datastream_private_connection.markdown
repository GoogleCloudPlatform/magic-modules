---
subcategory: "Datastream"
description: |-
  Get information about a Google Cloud Datastream Private Connection.
---

# google_datastream_private_connection

Get information about a Google Cloud Datastream Private Connection. For more information see
the [official documentation](https://cloud.google.com/datastream/docs/private-connectivity)
and [API](https://cloud.google.com/datastream/docs/reference/rest/v1/projects.locations.privateConnections).

## Example Usage

```hcl
data "google_datastream_private_connection" "default" {
  private_connection_id = "my-connection"
}
```

## Argument Reference

The following arguments are supported:

* `private_connection_id` - (Required) The ID of the Datastream Private Connection.

- - -

* `location` -
  (Required)
  The canonical id of the location. For example: us-east1.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.