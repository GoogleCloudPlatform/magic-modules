---
title: "Fix a permadiff"
weight: 60
---

# Fix a permadiff

Permadiffs are an extremely common class of errors that users experience. They manifest as diffs at plan time on fields that a user has not modified in their configuration. They can also show up as test failures with the error message: "After applying this test step, the plan was not empty."

In a general sense, permadiffs are caused by the API returning a different value for the field than what the user sent, which causes Terraform to try to re-send the same request, which gets the same response, which continues to result in the user seeing a diff. In general, APIs that return exactly what the user sent are more friendly for Terraform or other declarative tooling. However, many GCP APIs normalize inputs, have server-side defaults that are returned to the user, do not return all the fields set on a resource, or return data in a different format in some other way.

This page outlines best practices for working around various types of permadiffs in the `google` and `google-beta` providers.

## API returns default value for unset field {#default}

For new fields, if possible, set a client-side default that matches the API default. This will prevent the diff and will allow users to accurately see what the end state will be if the field is not set in their configuration. A client-side default should only be used if the API sets the same default value in all cases and the default value will be stable over time. Changing a client-side default is a [breaking change]({{< ref "/develop/breaking-changes/breaking-changes" >}}).

{{< tabs "default_value" >}}
{{< tab "MMv1" >}}
```yaml
default_value: DEFAULT_VALUE
```

In the providers, this will be converted to:

```go
"field": {
    // ...
    Default: "DEFAULT_VALUE",
}
```

See [SDKv2 Schema Behaviors - Default ↗](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#default) for more information.
{{< /tab >}}
{{< tab "Handwritten" >}}
```go
"field": {
    // ...
    Default: "DEFAULT_VALUE",
}
```

See [SDKv2 Schema Behaviors - Default ↗](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#default) for more information.
{{< /tab >}}
{{< /tabs >}}

For existing fields (or new fields that are not eligible for a client-side default), mark the field as having an API-side default. If the field is not set (or is set to an "empty" value such as zero, false, or an empty string) the provider will treat the most recent value returned by the API as the value for the field, and will send that value for the field on subsequent requests. The field will show as `(known after apply)` in plans and it will not be possible for the user to explicitly set the field to an "empty" value.

{{< tabs "default_from_api" >}}
{{< tab "MMv1" >}}
```yaml
default_from_api: true
```

In the providers, this will be converted to:

```go
"field": {
    // ...
    Optional: true,
    Computed: true,
}
```

See [SDKv2 Schema Behaviors - Optional ↗](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#optional) and [SDKv2 Schema Behaviors - Computed ↗](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#computed) for more information.
{{< /tab >}}
{{< tab "Handwritten" >}}
```go
"field": {
    // ...
    Optional: true,
    Computed: true,
}
```

See [SDKv2 Schema Behaviors - Optional ↗](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#optional) and [SDKv2 Schema Behaviors - Computed ↗](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#computed) for more information.
{{< /tab >}}
{{< /tabs >}}

## API returns an empty value if default value is sent {#default_if_empty}

Use a flattener to store the default value in state if the response has an empty (or unset) value.

{{< tabs "default_if_empty" >}}
{{< tab "MMv1" >}}
Use the standard `default_if_empty` flattener.

```yaml
custom_flatten: 'templates/terraform/custom_flatten/default_if_empty.erb'
```
{{< /tab >}}
{{< tab "Handwritten" >}}
```go
func flattenResourceNameFieldName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
    if v == nil || tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
        return "DEFAULT_VALUE"
    }
    // Any other necessary logic goes here.
    return v
}
```
{{< /tab >}}
{{< /tabs >}}

## API normalizes a value {#normalized-value}

In cases where the API normalizes and returns a value in a simple, predictable way (such as capitalizing the value) add a diff suppress function for the field to suppress the diff.

The `tpgresource` package in each provider supplies diff suppress functions for the following common cases:

- `tpgresource.CaseDiffSuppress`: Suppress diffs from capitalization differences between the user's configuration and the API.
- `tpgresource.DurationDiffSuppress`: Suppress diffs from duration format differences such as "60.0s" vs "60s". This is necessary for [`Duration`](https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#duration) API fields.
- `tpgresource.ProjectNumberDiffSuppress`: Suppress diffs caused by the provider sending a project ID and the API returning a project number.


{{< tabs "diff_suppress_simple" >}}

{{< tab "MMv1" >}}
```yaml
# Use a built-in function
diff_suppress_func: tpgresource.CaseDiffSuppress

# Reference a resource-specific function
diff_suppress_func: resourceNameFieldNameDiffSuppress
```

Define resource-specific functions in a [`custom_code.constants`](https://googlecloudplatform.github.io/magic-modules/develop/custom-code/#add-reusable-variables-and-functions) file.

```go
func resourceNameFieldNameDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
    // Separate function for easier unit testing
    return resourceNameFieldNameDiffSuppressLogic(old, new)
}

func resourceNameFieldNameDiffSuppressLogic(old, new) bool {
    // Diff suppression logic. Returns true if the diff should be suppressed - that is, if the
    // old and new values should be considered "the same".
}
```

See [SDKv2 Schema Behaviors - DiffSuppressFunc ↗](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#diffsuppressfunc) for more information.
{{< /tab >}}
{{< tab "Handwritten" >}}
Define resource-specific functions in your service package, for example at the top of the related resource file.

```go
func resourceNameFieldNameDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
    // Separate function for easier unit testing
    return resourceNameFieldNameDiffSuppressLogic(old, new)
}

func resourceNameFieldNameDiffSuppressLogic(old, new) bool {
    // Diff suppression logic. Returns true if the diff should be suppressed - that is, if the
    // old and new values should be considered "the same".
}
```

Reference diff suppress functions from the field definition.

```go
"field": {
    // ...
    DiffSuppressFunc: resourceNameFieldNameDiffSuppress,
}
```

See [SDKv2 Schema Behaviors - DiffSuppressFunc ↗](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#diffsuppressfunc) for more information.
{{< /tab >}}
{{< /tabs >}}

## API field that is never included in the response {#ignore_read}

This is common for fields that store credentials or similar information. Such fields should also be marked as [`sensitive`]({{< ref "/develop/field-reference#sensitive" >}}).

In the flattener for the field, return the value of the field in the user's configuration.

{{< tabs "ignore_read" >}}
{{< tab "MMv1" >}}
On top-level fields, this can be done with:

```yaml
ignore_read: true
```

For nested fields, `ignore_read` is [not currently supported](https://github.com/hashicorp/terraform-provider-google/issues/12410), so this must be implemented with a [custom flattener]({{< ref "/develop/custom-code#custom_flatten" >}}). You will also need to add the field to `ignore_read_extra` on any examples that are used to generate tests; this will cause tests to ignore the field when checking that the values in the API match the user's configuration.

```go
func flatten<%= prefix -%><%= titlelize_property(property) -%>(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
    // We want to ignore read on this field, but cannot because it is nested
    return d.Get("path.0.to.0.nested.0.field")
}
```

```yaml
examples:
 - !ruby/object:Provider::Terraform::Examples
   # example configuration
   ignore_read_extra:
     - "path.0.to.0.nested.0.field"
```
{{< /tab >}}
{{< tab "Handwritten" >}}
Use `d.Get` to set the flattened value to be the same as the user-configured value (instead of a value from the API).

```go
func flattenParentField(d *schema.ResourceData, disk *compute.AttachedDisk, config *transport_tpg.Config) []map[string]interface{} {
    result := map[string]interface{}{
        "nested_field": d.Get("path.0.to.0.parent_field.0.nested_field")
    }
    return []map[string]interface{}{result}
}
```

In tests, add the field to `ImportStateVerifyIgnore` on any relevant import steps.

```
{
    ResourceName:            "google_product_resource.default",
    ImportState:             true,
    ImportStateVerify:       true,
    ImportStateVerifyIgnore: []string{""path.0.to.0.parent_field.0.nested_field"},
},
```
{{< /tab >}}
{{< /tabs >}}

## API returns a list in a different order than was sent {#list-order}

For an Array of nested objects, convert it to a Set – this is a [breaking change]({{< ref "/develop/breaking-changes/breaking-changes" >}}) and can only happen in a major release.

For an Array of simple values (such as strings or ints), rewrite the value in the flattener to match the order in the user's configuration. This will also simplify diffs if new values are added or removed.

{{< tabs "diff_suppress_list" >}}

{{< tab "MMv1" >}}
Add a [custom flattener]({{< ref "/develop/custom-code#custom_flatten" >}}) for the field.

```go
func flatten<%= prefix -%><%= titlelize_property(property) -%>(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
    configValue := d.Get("path.0.to.0.parent_field.0.nested_field").([]string)

    ret := []string{}
    // Add values from v to ret to match order in configValue and put any new strings at the end

    return ret
}
```
{{< /tab >}}
{{< tab "Handwritten" >}}
Define resource-specific functions in your service package, for example at the top of the related resource file.

```go
func flattenResourceNameFieldName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
    configValue := d.Get("path.0.to.0.parent_field.0.nested_field").([]string)

    ret := []string{}
    // Add values from v to ret to match order in configValue and put any new strings at the end

    return ret
}
```
{{< /tab >}}
{{< /tabs >}}
