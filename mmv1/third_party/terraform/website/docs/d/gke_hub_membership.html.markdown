---
subcategory: "GKEHub"
description: |-
  Retrieves the details of a GKE Hub Membership.
---

# `google_gke_hub_membership`

Retrieves the details of a specific GKE Hub Membership. Use this data source to retrieve the membership's configuration and state.

## Example Usage

```hcl
data "google_gke_hub_membership" "example" {
  project       = "my-project-id"
  location      = "global"
  name          = "my-membership-id" # GKE Cluster's name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The GKE Cluster name of the GKE Hub Membership id.

* `location` - (Required) The location for the GKE Hub Membership.
    Currently only `global` is supported.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `name` - The full resource name of the membership, in the format: `projects/{project}/locations/{location}/memberships/{membership_id}`.

* `labels` - GCP labels for this Membership.

* `description` - A textual description of this membership.

* `endpoint` - (Output)
  Endpoint information to reach this member.
  Structure is documented below.

* `state` - (Output)
  State of the Membership resource.
  Structure is documented below.

* `authority` - (Output)
  Authority encodes how Google will recognize identities from this Membership.
  See the workload identity documentation for more details:
  https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity
  Structure is documented below.

* `unique_id` - A globally unique identifier for this Membership.

* `monitoring_config` - (Output)
  Monitoring configuration of the GKE Hub membership.
  Structure is documented below.

---

The `endpoint` block contains:

* `gke_cluster` - (Output) If this Membership is a GKE cluster, this field will be set.
  * `resource_link` - (Output) Self-link of the GKE cluster.
* `kubernetes_metadata` - (Output) Useful Kubernetes metadata.
  * `kubernetes_api_server_version` - (Output) Kubernetes API server version string.
  * `node_provider_id` - (Output) Node providerID as seen by the Kubernetes API server.
  * `node_count` - (Output) Node count as reported by Kubernetes.
  * `vcpu_count` - (Output) vCPU count as reported by Kubernetes.
  * `memory_mb` - (Output) The total memory on the nodes as reported by Kubernetes.
  * `update_time` - (Output) The time at which these details were last updated.

The `state` block contains:

* `code` - (Output) The current state of the Membership resource. (Code.CREATING, Code.READY, Code.DELETING, Code.UPDATING, Code.SERVICE_UPDATING)

The `authority` block contains:

* `issuer` - (Output) An identity provider that reflects the `issuer` in the workload identity pool.
* `workload_identity_pool` - (Output) The name of the workload identity pool in which `issuer` lives.
* `identity_provider` - (Output) The name of the identity provider.

The `monitoring_config` block contains:

* `project` - (Output) The project in which the resource belongs.
* `location` - (Output) The location of the monitoring configuration.
* `cluster` - (Output) The cluster name of the monitoring configuration.
* `kubernetes_metrics_prefix` - (Output) Kubernetes metrics prefix.