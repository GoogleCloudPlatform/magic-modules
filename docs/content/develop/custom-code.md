---
title: "Add custom resource code"
weight: 39
---

# Add custom resource code

This document covers how to add "custom code" to [MMv1 resources]({{< ref "/get-started/how-magic-modules-works#mmv1" >}}). Custom code can be used to add arbitrary logic to a resource while still generating most of the code; it allows for a balance between maintainability and supporting real-worlds APIs that deviate from what MMv1 can support. Custom code should only be added if the desired behavior can't be achieved otherwise.

Most custom code attributes are strings that contain a path to a template file relative to the `mmv1` directory. For example:

```yaml
custom_code:
  # References mmv1/templates/terraform/custom_delete/resource_name_custom_delete.go.tmpl
  custom_delete: templates/terraform/custom_delete/resource_name_custom_delete.go.tmpl
```

By convention, the template files are stored in a directory matching the type of custom code, and the name of the file includes the resource (and, if relevant, field) impacted by the custom code. Like handwritten resource and test code, custom code is written as go templates which render go code.

When in doubt about the behavior of custom code, write the custom code, [generate the providers]({{< ref "/get-started/generate-providers" >}}), and inspect what changed in the providers using `git diff`.

The following sections describe types of custom code in more detail.

## Add reusable variables and functions

```yaml
custom_code:
  constants: templates/terraform/constants/PRODUCT_RESOURCE.go.tmpl
```

Use `custom_code.constants` to inject top-level code in a resource file. This is useful for anything that should be referenced from other parts of the resource, such as:

- Constants
- Regexes compiled at build time
- Functions, such as [diff suppress functions]({{<ref "/develop/field-reference#diff_suppress_func" >}}),
  [validation functions]({{<ref "/develop/field-reference#validation" >}}),
  CustomizeDiff functions, and so on.
- Methods

Any custom functions added should have thorough [unit tests]({{< ref "/develop/test/test#add-unit-tests" >}}).

## Modify the API request or response

API requests and responses can be modified in the following order:

1. Modify the API request value for a specific field
2. Modify the API request data for an entire resource
3. Modify the API response data for an entire resource
4. Modify the API response value for a specific field

These are described in more detail in the following sections.

### Modify the API request value for a specific field {#custom_expand}

```yaml
- name: 'FIELD'
  type: String
  custom_expand: 'templates/terraform/custom_expand/PRODUCT_RESOURCE_FIELD.go.tmpl'
```

Set `custom_expand` on a field to inject code that modifies the value to send to the API for that field. Custom expanders run _before_ any [`encoder` or `update_encoder`]({{< ref "#encoder" >}}). The referenced file must include the function signature for the expander. For example:

```erb
func expand{{$.GetPrefix}}{{$.TitlelizeProperty}}(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
  if v == nil {
    return nil, nil
  }

  return base64.StdEncoding.EncodeToString([]byte(v.(string))), nil
}
```

The parameters the function receives are:

- `v`: The value for the field
- `d`:  Terraform resource data. Use `d.Get("field_name")` to get a field's current value.
- `config`: Config object. Can be used to make API calls.

The function returns a final value that will be sent to the API.

### Modify the API request data for an entire resource {#encoder}

```yaml
custom_code:
  encoder: templates/terraform/encoder/PRODUCT_RESOURCE.go.tmpl
  update_encoder: templates/terraform/update_encoder/PRODUCT_RESOURCE.go.tmpl
```

Use `custom_code.encoder` to inject code that modifies the data that will be sent in the API request. This is useful if the API expects the data to be in a significantly different structure than Terraform does - for example, if the API expects the entire object to be nested under a key, or a particular field must never be sent to the API. The encoder will run _after_ any [`custom_expand`]({{< ref "#custom_expand" >}}) code.

The encoder code will be wrapped in a function like:

```go
func resourceProductResourceEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
  // Your code will be injected here.
}
```

The parameters the function receives are:

- `d`: Terraform resource data. Use `d.Get("field_name")` to get a field's current value.
- `meta`: Can be cast to a Config object (which can make API calls) using `meta.(*transport_tpg.Config)`
- `obj`: The data that will be sent to the API. 

The function returns data that will be sent to the API and an optional error.

If the Create and Update methods for the resource need different logic, set `custom_code.update_encoder` to override the logic for update only. It is otherwise the same as `custom_code.encoder`.


### Modify the API response data for an entire resource {#decoder}

```yaml
custom_code:
  decoder: templates/terraform/decoder/PRODUCT_RESOURCE.go.tmpl
```


Use `custom_code.decoder` to inject code that modifies the data recieved from an API response. This is useful if the API returns data in a significantly different structure than what Terraform expects - for example, if the API returns the entire object nested under a key, or uses a different name for a field in the response than in the request. The decoder will run _before_ any [`custom_flatten`]({{< ref "#custom_flatten" >}}) code.

The decoder code will be wrapped in a function like:

```go
func resourceProductResourceDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
    // Your code will be injected here.
}
```

The parameters the function receives are:

- `d`: Terraform resource data. Use `d.Get("field_name")` to get a field's current value.
- `meta`: Can be cast to a Config object (which can make API calls) using `meta.(*transport_tpg.Config)`
- `res`: The data ("response") returned by the API. 

The function returns data that will be set in Terraform state and an optional error.

### Modify the API response value for a specific field {#custom_flatten}

```yaml
- name: 'FIELD'
  type: String
  custom_flatten: 'templates/terraform/custom_flatten/PRODUCT_RESOURCE_FIELD.go.tmpl'
```

Set `custom_flatten` on a field to inject code that modifies the value returned by the API prior to storing it in Terraform state. Custom flatteners run _after_ any [`decoder`]({{< ref "#encoder" >}}). The referenced file must include the function signature for the flattener. For example:

```erb
func flatten{{$.GetPrefix}}{{$.TitlelizeProperty}}(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
  if v == nil {
    return "0"
  }
  return v
}
```

The parameters the function receives are:

- `v`: The value for the field
- `d`:  Terraform resource data. Use `d.Get("field_name")` to get a field's current value.
- `config`: Config object. Can be used to make API calls.

The function returns a final value that will be stored in Terraform state for the field, which will be compared with the user's configuration to determine if there is a diff.

## Inject code before / after CRUD operations and Import {#pre_post_injection}

```yaml
custom_code:
  pre_create: templates/terraform/pre_create/PRODUCT_RESOURCE.go.tmpl
  post_create: templates/terraform/post_create/PRODUCT_RESOURCE.go.tmpl

  pre_read: templates/terraform/pre_read/PRODUCT_RESOURCE.go.tmpl

  pre_update: templates/terraform/pre_update/PRODUCT_RESOURCE.go.tmpl
  post_update: templates/terraform/post_update/PRODUCT_RESOURCE.go.tmpl

  pre_delete: templates/terraform/pre_delete/PRODUCT_RESOURCE.go.tmpl
  post_delete: templates/terraform/post_delete/PRODUCT_RESOURCE.go.tmpl

  post_import: templates/terraform/post_import/PRODUCT_RESOURCE.go.tmpl
```

CRUD operations can be modified with pre/post hooks. This code will be injected directly into the relevant CRUD method as close as possible to the related API call and will have access to any variables that are present when it runs. `pre_create` and `pre_update` run after any [`encoder`]({{< ref "#encoder" >}}). Some example use cases:

- Use `post_create` to set an update-only field after create finishes.
- Use `pre_delete` to detach a disk before deleting it.
- Use `post_import` to parse attributes from the import ID and call `d.Set("field")` so that the resource can be read from the API.

### Custom create error handling

```yaml
custom_code:
  post_create_failure: templates/terraform/post_create_failure/PRODUCT_RESOURCE.go.tmpl
```

Use `custom_code.post_create_failure` to inject code that runs if a Create request to the API returns an error.

The post_create_failure code will be wrapped in a function like:

```go
func resourceProductResourcePostCreateFailure(d *schema.ResourceData, meta interface{}) {
    // Your code will be injected here.
}
```

The parameters the function receives are:

- `d`: Terraform resource data. Use `d.Get("field_name")` to get a field's current value.
- `meta`: Can be cast to a Config object (which can make API calls) using `meta.(*transport_tpg.Config)`

## Replace entire CRUD methods

```yaml
custom_code:
  custom_create: templates/terraform/custom_create/PRODUCT_RESOURCE.go.tmpl
  custom_update: templates/terraform/custom_update/PRODUCT_RESOURCE.go.tmpl
  custom_delete: templates/terraform/custom_delete/PRODUCT_RESOURCE.go.tmpl
  custom_import: templates/terraform/custom_import/PRODUCT_RESOURCE.go.tmpl
```

Custom methods replace the entire contents of the Create, Update, Delete, or Import methods. For example:

```go
func resourceProductResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
    // Your code will be injected here.
}
```

Custom methods are similar to handwritten code and should be avoided if possible. If you have to replace two or more methods, the resource should be handwritten instead.

## Add extra fields to a resource

Use `custom_code.extra_schema_entry` to add additional fields to a resource. Do not use `extra_schema_entry` unless there is no other option. The extra fields are injected at the end of the resource's [`Schema` field](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-types).  They should be formatted as entries in the map. For example:

```go
"foo": &schema.Schema{ ... },
```

Any fields added in this way will need to be have documentation manually added using the top-level `docs` field:

```yaml
docs:
  optional_properties: |
    * `FIELD_NAME` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) FIELD_DESCRIPTION
```

See [Add documentation (Handwritten)]({{< ref "/develop/resource#add-documentation" >}}) for more information about what to include in the field documentation.

## What's next?

- [Add tests]({{< ref "/develop/test/test.md" >}})
- [Run tests]({{< ref "/develop/test/run-tests.md" >}})
