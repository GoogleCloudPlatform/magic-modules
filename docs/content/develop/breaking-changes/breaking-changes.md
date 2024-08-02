---
title: "Types of breaking changes"
weight: 10
aliases:
- /docs/content/develop/breaking-changes
- /reference/breaking-change-detector
- /develop/breaking-changes

---

# Types of breaking changes

A "breaking change" is any change that requires an end user to modify a
previously-valid configuration after a provider upgrade. In this context,
a "valid configuration" is one that:

- Is considered syntactically correct by `terraform validate`
- Does not return an error during `terraform apply`
- Creates, updates, deletes, or does not modify resources
- Only manages resources that have not been altered with other tools,
  such as `gcloud` or Cloud Console.

This document lists many types of breaking changes but may not be entirely
comprehensive. Some types of changes that would normally be "breaking" may
have specific mitigating circumstances that make them non-breaking.

For more information, see
[Make a breaking change]({{< ref "/develop/breaking-changes/make-a-breaking-change" >}}).

## Provider-level breaking changes

* <a name="provider-config-fundamental"></a>Changing fundamental provider behavior such as:
  * authentication
  * environment variable usage
  * restricting retry behavior

## Resource-level breaking changes

* <a name="resource-map-resource-removal-or-rename"></a>Removing or renaming a resource
  or datasource
* <a name="resource-id"></a> Changing resource ID format
  * Terraform uses resource ID to read resource state from the API. Modification of
    the ID format will break the ability to parse the IDs from any deployments.
* <a name="resource-import-format"></a> Removing or altering resource import ID formats
  * Automation written by end users may rely on specific import formats.
* Changes to default resource behavior
  *  Changing resource deletion behavior
    * In limited cases changes may be permissible if the prior behavior could **never** succeed.
    * Changing resource deletion to skip deleting the resource by default if delete was previously called
    * Changing resource deletion to specify a force flag
  * Adding a new field with a default different from the API default
    * If an API default is expected to change- a breaking change for the API- use `default_from_api` which will avoid sending a value and safely take the server default in Terraform

## Field-level breaking changes

* <a name="resource-schema-field-removal-or-rename"></a>Removing or renaming a field 
* <a name="field-changing-type"></a> Changing field output type
  * Between primitive types, like changing a String to an Integer
  * Between complex types like changing a List to a Set.
  * Changing the field type between primitive and complex data
    types is not possible. For this scenario, field renames are preferred.
* <a name="field-optional-to-required"></a> Making an optional field required or adding a new required field
* <a name="field-becoming-computed"></a> Making a settable field read-only
  * For MMv1 resources, adding `output: true` to an existing field.
  * For handwritten resources, adding `Computed: true` to a field that does not have `Optional: true` set.
* <a name="field-oc-to-c"></a> Removing support for API-side defaults
  * For MMv1 resources, removing `default_from_api: true`.
  * For handwritten resources, altering a field schema with `Computed: true` + `Optional: true`
    to only have `Optional: true`.
* <a name="field-changing-default-value"></a> Adding or changing a default value
  * Default values in Terraform are used to replace null values in configuration at
    plan/apply time and **do not** respect previously-configured values by the user.
    Even in major releases, these changes are often undesirable, as their impact is extremely broad.

    When a default is changed, every user that has not specified an explicit value in their
    configuration will see Terraform propose changing the value of the field **including**
    if the change will destroy and recreate the resource due to changing an immutable value.
    Default changes in the provider are comparable in impact to default changes in an API,
    and modifying examples and modules may achieve the intended effect with a smaller blast radius.
* <a name="field-changing-data-format"></a> Modifying how field data is stored in state
  * For example, changing the case of a value returned by the API in a flattener or decorder
* <a name="field-removing-diff-suppress"></a> Removing diff suppression from a field.
  * For MMv1 resources, removing `diff_suppress_func` from a field.
  * For handwritten resources, removing `DiffSuppressFunc` from a field.
* <a name="field-adding-subfield-to-config-mode-attr"></a> Adding a subfield to
  a SchemaConfigModeAttr field.
  * Subfields of SchemaConfigModeAttr fields are treated as required even if the schema says they are optional.
* Removing update support from a field.

### Making validation more strict

* <a name="field-growing-min"></a> Increasing the minimum number of items in an array
  * For MMv1 resources, increasing `min_size` on an Array field.
  * For handwritten resources, increasing `MinItems` on an Array field.
* <a name="field-shrinking-max"></a> Decreasing the maximum number of items in an array
  * For MMv1 resources, decreasing `max_size` on an Array field.
  * For handwritten resources, decreasing `MaxItems` on an Array field.
* Adding validation to a field that previously had no validation
  * For MMv1 resources, adding `validation` to a field.
  * For handwritten resources, adding `ValidateFunc` to a field.

