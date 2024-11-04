---
title: "Deletion behaviors"
weight: 20
---

# Deletion behaviors

{{< hint info >}}
**Note:** This page covers best practices guidance for the Terraform provider for Google Cloud, which is used to ensure a consistent UX for Terraform users across providers or GCP users across the Google provider. Generally, this guidance should be followed and exceptions should be clearly demarcated / discussed.
{{< /hint >}}

## Mitigating data loss risk via deletion_protection {#deletion_protection}

Some resources, such as databases, have a significant risk of unrecoverable data loss if the resource is accidentally deleted due to a change to a ForceNew field. For these resources, the best practice is to add a `deletion_protection` field that defaults to `true`, which prevents the resource from being deleted if enabled. Although it is a small breaking change, for users, the benefits of `deletion_protection` defaulting to `true` outweigh the cost.

APIs also sometimes add `deletion_protection` fields, which will generally default to `false` for backwards-compatibility reasons. Any `deletion_protection` API field added to an existing Terraform resource must match the API default initially. The default may be set to `true` in the next major release. For new Terraform resources, any `deletion_protection` field should default to `true` in Terraform regardless of the API default. When creating the corresponding Terraform field, the name
should match the API field name (i.e. it need not literally be named `deletion_protection` if the API uses something different) and should be the same field type (example: if the API field is an enum, so should the Terraform field).

A resource can have up to two `deletion_protection` fields (with different names): one that represents a field in the API, and one that is only in Terraform. This could happen because the API added its field after `deletion_protection` already existed in Terraform; it could also happen because a separate field was added in Terraform to make sure that `deletion_protection` is enabled by default. In either case, they should be reconciled into a single field (that defaults to enabled and whose name matches the API field) in the next major release.

Resources that do not have a significant risk of unrecoverable data loss or similar critical concern will not be given `deletion_protection` fields.

See [Client-side fields]({{< ref "/develop/client-side-fields" >}}) for information about adding `deletion_protection` fields.

{{< hint info >}}
**Note:** The previous best practice was a field called `force_delete` that defaulted to `false`. This is still present on some resources for backwards-compatibility reasons, but `deletion_protection` is preferred going forward.
{{< /hint >}}

## Deletion policy {#deletion_policy}

Some resources need to let users control the actions taken add deletion time. For these resources, the best practice is to add a `deletion_policy` enum field that defaults to an empty string and allows special values that control the deletion behavior.

One common example is `ABANDON`, which is useful if the resource is safe to delete from Terraform but could cause problems if deleted from the API - for example, `google_bigtable_gc_policy` deletion can fail in replicated instances. `ABANDON` indicates that attempts to delete the resource should remove it from state without actually deleting it.

See [Client-side fields]({{< ref "/develop/client-side-fields" >}}) for information about adding `deletion_policy` fields.
