---
title: "Immutable fields"
weight: 10
aliases:
  - /docs/best-practices
  - /best-practices
---

# Immutable fields

{{< hint info >}}
**Note:** This page covers best practices guidance for the Terraform provider for Google Cloud, which is used to ensure a consistent UX for Terraform users across providers or GCP users across the Google provider. Generally, this guidance should be followed and exceptions should be clearly demarcated / discussed.
{{< /hint >}}

[`ForceNew`](https://developer.hashicorp.com/terraform/intro#how-does-terraform-work) in a Terraform resource schema attribute that indicates that a field is immutable – that is, that a change to the field requires the resource to be destroyed and recreated.

This is necessary and required for cases where a field can't be updated in-place, so that [Terraform's core workflow](https://developer.hashicorp.com/terraform/intro#how-does-terraform-work) of aligning real infrastructure with configuration can be achieved. If a field or resource can never be updated in-place and is not marked with `ForceNew`, that is considered a bug in the provider.

Some fields or resources may be possible to update in place, but only under specific conditions. In these cases, you can treat the field as updatable - that is, do not mark it as ForceNew; instead, implement standard update functionality. Then, call `diff.ForceNew` inside a [`CustomizeDiff`](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/customizing-differences) if the appropriate conditions to allow update in place are not met. Any `CustomizeDiff` function like this must be thoroughly unit tested. Making a field conditionally updatable like this is considered a good and useful enhancement in cases where recreation is costly and conditional updates do not introduce undue complexity.

In complex cases, it is better to mark the field `ForceNew` to ensure that users can apply their configurations successfully.

## Safeguarding against deletion

See [Deletion behaviors]({{< ref "/best-practices/deletion-behaviors" >}}) for some mitigations against accidental deletion or other means to safeguard against deletion.
