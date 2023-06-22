---
title: "Best practices"
weight: 25
aliases:
  - /docs/best-practices
---

# Best practices

This document describes best practices for contributors using Magic Modules to ensure an internally-consistent UX for the Google Terraform provider, as well as consistency when compared to other Terraform providers.

## Mark immutable fields and resources as immutable {#immutability}

[Terraform's core purpose](https://developer.hashicorp.com/terraform/intro#how-does-terraform-work) is to align real infrastructure with user-provided configuration. This means that fields or resources which cannot be updated in place must be deleted and recreated if the field or resource configuration changes.

Therefore, immutable fields and resources must be marked as immutable, using `immutable: true` for MMv1 resources or `ForceNew` in handwritten resources. If a field or resource can never be updated in-place and is not marked with `ForceNew`, that is considered a bug in the provider.

### Support conditional mutability if possible {#conditional_mutability}

Some fields or resources may be possible to update in place, but only under specific conditions. In these cases, you can treat the field as updatable - that is, do not mark it as ForceNew; instead, implement standard update functionality. Then, call `diff.ForceNew` inside a [`CustomizeDiff`](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/customizing-differences) if the appropriate conditions to allow update in place are not met. Any `CustomizeDiff` function like this must be thoroughly unit tested. Making a field conditionally updatable like this is considered a good and useful enhancement in cases where recreation is costly and conditional updates do not introduce undue complexity.

In complex cases, it is better to mark the field `ForceNew` to ensure that users can apply their configurations successfully.

### Mitigate data loss risk via deletion_protection {#deletion_protection}

Some resources, such as databases, have a significant risk of unrecoverable data loss if the resource is accidentally deleted due to a change to a ForceNew field. For these resources, the best practice is to add a `deletion_protection` field that defaults to `true`, which prevents the resource from being deleted if enabled. Although it is a small breaking change, for users, the benefits of `deletion_protection` defaulting to `true` outweigh the cost.

APIs also sometimes add `deletion_protection` fields, which will generally default to `false` for backwards-compatibility reasons. Any `deletion_protection` API field added to an existing Terraform resource must match the API default initially. The default may be set to `true` in the next major release. For new Terraform resources, any `deletion_protection` field should default to `true` in Terraform regardless of the API default.

A resource can have up to two `deletion_protection` fields (with different names): one that represents a field in the API, and one that is only in Terraform. This could happen because the API added its field after `deletion_protection` already existed in Terraform; it could also happen because a separate field was added in Terraform to make sure that `deletion_protection` is enabled by default. In either case, they should be reconciled into a single field (that defaults to `true`) in the next major release.

Resources that do not have a significant risk of unrecoverable data loss or similar critical concern will not be given `deletion_protection` fields.

{{< hint info >}}
**Note:** The previous best practice was a field called `force_delete` that defaulted to `false`. This is still present on some resources for backwards-compatibility reasons, but `deletion_protection` is preferred going forward.
{{< /hint >}}

## Control actions at delete time with `deletion_policy` {#deletion_policy}

Some resources need to let users control the actions taken add deletion time. For these resources, the best practice is to add a `deletion_policy` enum field that defaults to an empty string and allows special values that control the deletion behavior.

One common example is `ABANDON`, which is useful if the resource is safe to delete from Terraform but could cause problems if deleted from the API - for example, `google_bigtable_gc_policy` deletion can fail in replicated instances. `ABANDON` indicates that attempts to delete the resource should remove it from state without actually deleting it.

See [magic-modules#13107](https://github.com/hashicorp/terraform-provider-google/pull/13107) for an example of adding a `deletion_policy` field to an existing resource.

### DiffSuppressFunc

Terraform allows fields to specify a
[DiffSuppressFunc](https://www.terraform.io/plugin/sdkv2/schemas/schema-behaviors#diffsuppressfunc),
which allows you to ignore diffs in cases where the two values are
**functionally identical**. This is generally useful when the API returns a
normalized value - for example by standardizing the case.

Note: The *preferred* behavior for APIs is to always return the value that the
user sent. DiffSuppressFunc is a workaround for APIs that don't.

The Terraform provider comes with a set of
["common diff suppress functions"](https://github.com/hashicorp/terraform-provider-google/blob/main/google/common_diff_suppress.go).
These fit frequent needs like ignoring whitespace at the beginning and end of a
string, or ignoring case differences.

If you need to define a custom diff specifically for your resource, you can do
so in a "constants" file, which is a `.go.erb` file in
[mmv1/templates/terraform/constants](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/constants)
named `<product>_<resource>.erb`. You can then declare this custom code in
`ResourceName.yaml`:

```yaml
--- !ruby/object:Api::Resource
  name: ResourceName

  custom_code: !ruby/object:Provider::Terraform::CustomCode
    constants: templates/terraform/constants/product_resource_name.go.erb
```

Once you have chosen a DiffSuppressFunc, you can declare it as an override on
your resource:

```yaml
--- !ruby/object:Api::Resource
  name: ResourceName

  properties:
    - !ruby/object:Api::Type::String
      name: "myField"
      diff_suppress_func: 'tpgresource.CaseDiffSuppress'
```

The value of diff_suppress_func can be any valid DiffSuppressFunc, including the
result of a function call. For example:

```yaml
diff_suppress_func: 'tpgresource.OptionalPrefixSuppress("folders/")'
```

Please make sure to add thorough unit tests (in addition to basic integration
tests) for your diff suppress func.

Example: DomainMapping (DomainMappingLabelDiffSuppress)

-   [DomainMapping.yaml resource file](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/products/cloudrun/DomainMapping.yaml)
    -   [`custom_code`](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/products/cloudrun/DomainMapping.yaml#L40)
    -   [`diff_suppress_func: 'resourceBigQueryDatasetAccessRoleDiffSuppress'`](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/products/bigquery/DatasetAccess.yaml#L112)
-   [constants file](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/templates/terraform/constants/cloud_run_domain_mapping.go.erb)
-   [unit tests](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/third_party/terraform/tests/resource_cloud_run_domain_mapping_test.go#L9)

### `custom_expand` and `custom_flatten` 
```
  # Overrides the default "expander" for the field. Expanders convert data from
  # a terraform representation to an API representation.
  # custom_expand: 'templates/terraform/custom_expand/PRODUCT_RESOURCE_FIELD.go.erb'

  # Overrides the default "flattener" for the field. Flatteners convert data
  # from an API representation to a terraform representation.
  # custom_expand: 'templates/terraform/custom_flatten/PRODUCT_RESOURCE_FIELD.go.erb'
```
