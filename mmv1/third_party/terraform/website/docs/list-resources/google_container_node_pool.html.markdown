---
subcategory: "Kubernetes (GKE) Engine"
description: |-
  List Google Kubernetes Engine node pools in a cluster for use with terraform query
  and .tfquery.hcl files.
---

# google_container_node_pool (list)

Lists GKE node pools in a cluster for use with
[`terraform query`](https://developer.hashicorp.com/terraform/cli/commands/query) and
`.tfquery.hcl` files. Results correspond to existing
[`google_container_node_pool`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_node_pool)
managed resources.

For how list resources work in this provider, Terraform version requirements, and shared
`list` block arguments, refer to the guide
[Use list resources with terraform query (Google Cloud provider)](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_list_resources_with_terraform_query).

## Example

```hcl
list "google_container_node_pool" "all" {
  provider = google

  config {
    # Optional. Defaults to the provider project when omitted.
    # project = "my-project"
    location = "us-central1-a"
    cluster  = "my-cluster"
  }
}
```

Run `terraform query` from the directory that contains the `.tfquery.hcl` file.

## Configuration (`config` block)

* `project` - (Optional) Project ID containing the cluster. If unset, the provider project is used.

* `location` - (Required) Cluster location (zone or region).

* `cluster` - (Required) Cluster name.

## Results

By default each result includes resource identity for `google_container_node_pool`:

* `project` - Project ID.

* `location` - Cluster location.

* `cluster` - Cluster name.

* `name` - Node pool name.

With `include_resource = true` on the `list` block, results also include the full resource-style
attributes documented for the managed
[`google_container_node_pool` resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_node_pool#attributes-reference).
