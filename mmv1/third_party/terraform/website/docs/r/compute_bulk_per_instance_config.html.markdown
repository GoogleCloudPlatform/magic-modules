---
subcategory: "Compute Engine"
description: |-
  A config defined for multiple managed instances that belong to an instance group manager
---

# google_compute_bulk_per_instance_config

A config defined for multiple managed instances that belong to an instance group manager. It preserves the instance name
across instance group manager operations.

It's recommended to use it with `lifecycle.ignore_changes[target_size]` on instance group manager.

To get more information about BulkPerInstanceConfig, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/instanceGroupManagers)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/compute/docs/instance-groups/bulk-create-instances-in-mig)

## Argument Reference

The following arguments are supported:


* `per_instance_configs` -
  (Required)
  The list of per-instance configs.
  Structure is [documented below](#nested_per_instance_configs).

* `instance_group_manager` -
  (Required)
  The instance group manager this bulk per instance config is part of.


* `name` -
  (Optional)
  The name for this bulk per-instance config.

* `zone` -
  (Optional)
  Zone where the containing instance group manager is located

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

* `deletion_policy` - (Optional) Whether Terraform will be prevented from destroying the resource. Defaults to DELETE.
	When a 'terraform destroy' or 'terraform apply' would delete the resource,
	the command will fail if this field is set to "PREVENT" in Terraform state.
	When set to "ABANDON", the command will remove the resource from Terraform
	management without updating or deleting the resource in the API.
	When set to "DELETE", deleting the resource is allowed.


<a name="nested_per_instance_configs"></a>The `per_instance_configs` block supports:

* `name` -
  (Required)
  The name for this per instance config's instance.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{project}}/{{zone}}/{{instance_group_manager}}/{{name}}`


## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
