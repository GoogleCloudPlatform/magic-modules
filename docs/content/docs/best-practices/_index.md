---
title: "Best practices"
weight: 25
---

# Best practices

The following is a list of best practices that contributions are expected to follow in order to ensure a consistent UX for the Google Terraform provider internally and also compared to other Terraform providers.

Terraform is a tool for provisioning and managing your infrastructure using human-readable configurations. The core Terraform workflow is:

1. Write a configuration
2. Plan the create / update / delete operations needed to get the real infrastructure to align with the configuration.
3. Apply (execute) the plan after approval

See [`How does Terraform work?`](https://developer.hashicorp.com/terraform/intro#how-does-terraform-work) from Hashicorp for a more detailed breakdown.

## ForceNew

[`ForceNew`](https://developer.hashicorp.com/terraform/intro#how-does-terraform-work) in a Terraform resource schema attribute that indicates that a field is immutable, i.e. that a change to the field requires the resource to be destroyed and recreated.

This is necessary and required for cases where a field can't be updated in-place, so that the core workflow of aligning real infrastructure with configuration can be achieved. Terraform assumes that all fields in a resource are either marked as `ForceNew` or that updates are handled through the `Update` function of the resource's implementation. If a field or resource can never be updated in-place and is not marked with `ForceNew`, that is considered a bug in the provider.

In addition, some fields or resources may be possible to update in place, but only under specific conditions. In these cases, you can call `diff.ForceNew` inside a [`CustomizeDiff`](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/customizing-differences) function to force recreation when pertinent. This is considered a good and useful enhancement in cases where it doesn't introduce undue complexity. Any `CustomizeDiff` function like this should be thoroughly unit tested.

### Mitigating data loss risk

For some resources, such as databases, there is a significant risk of data loss if the resource is accidentally deleted due to a change to a ForceNew field. However, because this is a core Terraform workflow, it is not possible (or desirable) to break user expectations by deviating from this best practice. Instead, we provide mitigations to prevent accidental deletion.

New resources that have a high risk of data loss should include a `deletion_policy` "virtual" enum field that defaults to an empty string and also allows `DELETE` (if deletion is possible) and (optionally if deletion is possible) `ABANDON` or other values. Terraform should reject all attempts to delete the resource unless "deletion_protection" is explicitly set to `DELETE`. If `ABANDON` is specified, then attempts to delete the resource should remove it from state without actually deleting it. See [magic-modules#13107](https://github.com/hashicorp/terraform-provider-google/pull/13107) for an example of adding a `deletion_policy` field to an existing resource.

Due to the impact of data loss, this is not currently considered a breaking change; however, that may change in the future. Ideally this field should be added when the resource is first created.

The previous best practices were:

- A virtual field called `force_delete` that defaulted to `false`
- A virtual field called `deletion_protection` that defaulted to `true`

These are **no longer considered best practices** because they are not extensible and are not as explicit; however, for backwards-compatibility reasons, they are still supported on existing resources.
