---
subcategory: "Oracle Database"
description: |-
  List all CloudVmClusters.
---

# google_oracle_database_cloud_vm_clusters

List all CloudVmClusters.

## Example Usage

```hcl
data "google_oracle_database_cloud_vm_clusters" "my_vmclusters"{
  location = "us-east4"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location of the resource.

- - -
* `project` - (Optional) The project to which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `CloudVmClusters` - A list of CloudVmClusters.

See [google_oracle_database_cloud_vm_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_oracle_database_cloud_vm_cluster#argument-reference) resource for details of the available attributes.
