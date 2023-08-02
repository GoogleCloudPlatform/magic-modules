---
subcategory: "Dataproc Cluster"
description: |-
  Get information about a Google Dataproc Cluster.
---

# google\_dataproc\_cluster

Get information about a Google Dataproc Cluster. For more information see
the [official documentation](https://cloud.google.com/dataproc/docs/)
and [API](https://cloud.google.com/dataproc/docs/apis).

## Example Usage

```hcl
data "google_dataproc_cluster" "my_cluster" {
  cluster_name  = "my-cluster"
  region        = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required) The name of the dataproc cluster.

* `region` - (Required) The location of the dataproc cluster. eg us-central1

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_dataproc_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/dataproc_cluster#argument-reference) resource for details of the available attributes.
