---
title: "Best practices"
weight: 25
---

# Best practices

The following is a list of best practices that contributions are expected to follow in order to ensure a consistent UX for the Google Terraform provider internally and also compared to other Terraform providers.

## ForceNew

[`ForceNew`](https://developer.hashicorp.com/terraform/intro#how-does-terraform-work) in a Terraform resource schema attribute that indicates that a field is immutable – that is, that a change to the field requires the resource to be destroyed and recreated.

This is necessary and required for cases where a field can't be updated in-place, so that [Terraform's core workflow](https://developer.hashicorp.com/terraform/intro#how-does-terraform-work) of aligning real infrastructure with configuration can be achieved. If a field or resource can never be updated in-place and is not marked with `ForceNew`, that is considered a bug in the provider.

Some fields or resources may be possible to update in place, but only under specific conditions. In these cases, you can call `diff.ForceNew` inside a [`CustomizeDiff`](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/customizing-differences) function to force recreation. This is considered a good and useful enhancement in cases where it doesn't introduce undue complexity. Any `CustomizeDiff` function like this must be thoroughly unit tested.

### Mitigating data loss risk via deletion_protection

Some resources, such as databases, have a significant risk of unrecoverable data loss if the resource is accidentally deleted due to a change to a ForceNew field. For these resources, the best practice is to add a `deletion_protection` field that defaults to `true`. This can be client-side only or can be backed by a server-side field. The Terraform field should default to `true` even if the server-side field does not.

APIs also sometimes add `deletion_protection` fields, which will default to `false` for backwards-compatibility reasons. However, for Terraform, the benefits of `deletion_protection` defaulting to `true` outweigh the smaller cost of the breaking change. In cases where both server-side and client-side fields have been separately added to the provider, they should be reconciled into a single field in the next major release.

The previous best practice was a field called `force_delete` that defaulted to `false`. This is still present on some resources for backwards-compatibility reasons, but `deletion_protection` is preferred going forward.

Resources that do not have a significant risk of unrecoverable data loss or similar critical concern will not be given `deletion_protection` fields.

## Deletion policy

Some resources need to let users control the actions taken add deletion time. For these resources, the best practice is to add a `deletion_policy` enum field that defaults to an empty string and allows special values that control the deletion behavior.

One common example is `ABANDON`, which is useful if the resource is safe to delete from Terraform but could cause problems if deleted from the API - for example, `google_bigtable_gc_policy` deletion can fail in replicated instances. `ABANDON` indicates that attempts to delete the resource should remove it from state without actually deleting it.

See [magic-modules#13107](https://github.com/hashicorp/terraform-provider-google/pull/13107) for an example of adding a `deletion_policy` field to an existing resource.
