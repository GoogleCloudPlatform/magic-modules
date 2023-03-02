---
subcategory: "Alloydb"
description: |-
  Fetches the details of a location.
---

# google\_alloydb\_location

Use this data source to get information about a particular location.

## Example Usage


```hcl
data "google_alloydb_location" "qa"{
    location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (required) The canonical id of the location. For example: `us-east1`.

* `project` - (optional) The ID of the project.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `name` - Resource name for the location, which may vary between implementations. For example: "projects/example-project/locations/us-east1".

* `location_id` - The canonical id for this location. For example: "us-east1"..

* `display_name` - The friendly name for this location, typically a nearby city name. For example, "Tokyo".

* `labels` - Cross-service attributes for the location. For example `{"cloud.googleapis.com/region": "us-east1"}`.

* `metadata` - Service-specific metadata. For example the available capacity at the given location.
