---
title: "Best practices"
weight: 25
aliases:
  - /docs/best-practices
---

# Best practices

The following is a list of best practices that contributions are expected to follow in order to ensure a consistent UX for the Google Terraform provider internally and also compared to other Terraform providers.

## ForceNew

[`ForceNew`](https://developer.hashicorp.com/terraform/intro#how-does-terraform-work) in a Terraform resource schema attribute that indicates that a field is immutable – that is, that a change to the field requires the resource to be destroyed and recreated.

This is necessary and required for cases where a field can't be updated in-place, so that [Terraform's core workflow](https://developer.hashicorp.com/terraform/intro#how-does-terraform-work) of aligning real infrastructure with configuration can be achieved. If a field or resource can never be updated in-place and is not marked with `ForceNew`, that is considered a bug in the provider.

Some fields or resources may be possible to update in place, but only under specific conditions. In these cases, you can treat the field as updatable - that is, do not mark it as ForceNew; instead, implement standard update functionality. Then, call `diff.ForceNew` inside a [`CustomizeDiff`](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/customizing-differences) if the appropriate conditions to allow update in place are not met. Any `CustomizeDiff` function like this must be thoroughly unit tested. Making a field conditionally updatable like this is considered a good and useful enhancement in cases where recreation is costly and conditional updates do not introduce undue complexity.

In complex cases, it is better to mark the field `ForceNew` to ensure that users can apply their configurations successfully.

### Mitigating data loss risk via deletion_protection

Some resources, such as databases, have a significant risk of unrecoverable data loss if the resource is accidentally deleted due to a change to a ForceNew field. For these resources, the best practice is to add a `deletion_protection` field that defaults to `true`, which prevents the resource from being deleted if enabled. Although it is a small breaking change, for users, the benefits of `deletion_protection` defaulting to `true` outweigh the cost.

APIs also sometimes add `deletion_protection` fields, which will generally default to `false` for backwards-compatibility reasons. Any `deletion_protection` API field added to an existing Terraform resource must match the API default initially. The default may be set to `true` in the next major release. For new Terraform resources, any `deletion_protection` field should default to `true` in Terraform regardless of the API default.

A resource can have up to two `deletion_protection` fields (with different names): one that represents a field in the API, and one that is only in Terraform. This could happen because the API added its field after `deletion_protection` already existed in Terraform; it could also happen because a separate field was added in Terraform to make sure that `deletion_protection` is enabled by default. In either case, they should be reconciled into a single field (that defaults to `true`) in the next major release.

Resources that do not have a significant risk of unrecoverable data loss or similar critical concern will not be given `deletion_protection` fields.

{{< hint info >}}
**Note:** The previous best practice was a field called `force_delete` that defaulted to `false`. This is still present on some resources for backwards-compatibility reasons, but `deletion_protection` is preferred going forward.
{{< /hint >}}

## Deletion policy

Some resources need to let users control the actions taken add deletion time. For these resources, the best practice is to add a `deletion_policy` enum field that defaults to an empty string and allows special values that control the deletion behavior.

One common example is `ABANDON`, which is useful if the resource is safe to delete from Terraform but could cause problems if deleted from the API - for example, `google_bigtable_gc_policy` deletion can fail in replicated instances. `ABANDON` indicates that attempts to delete the resource should remove it from state without actually deleting it.

See [magic-modules#13107](https://github.com/hashicorp/terraform-provider-google/pull/13107) for an example of adding a `deletion_policy` field to an existing resource.
