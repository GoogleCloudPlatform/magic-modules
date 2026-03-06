---
subcategory: "Apigee"
description: |-
  Get info about a Google Apigee Instance.
---

# google_apigee_instance
Get information about a Google Apigee Instance.

## Example Usage
```hcl
data "google_apigee_instance" "my_instance" { 
    name = "my-instance-name" 
    org_id = "organizations/my-org-id"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Apigee instance. [3]
* `org_id` - (Required) The Apigee Organization associated with the instance, in the format `organizations/{{org_name}}`. [3]

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `host` - The hostname or IP address of the exposed Apigee endpoint. [3]
* `port` - The port number of the exposed Apigee endpoint. [3]
* `service_attachment` - The PSC service attachment for the instance. [3]
* `location` - The GCP region where the instance resides. [3]
* `ip_range` - The IP range used by the instance. [3]