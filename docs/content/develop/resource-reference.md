---
title: "MMv1 resource reference"
weight: 32
aliases:
  - /reference/resource-reference
  - /reference/iam-policy-reference
---

# MMv1 resource reference

This page documents commonly-used properties for resources. For a full list of
available properties, see [resource.go ↗](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/resource.go).

## Basic

### `name`

API resource name.

### `description`

Resource description. Used in documentation.

Example:

```yaml
description: |
  This is a multi-line description
  of a resource. All multi-line descriptions must follow
  this format of using a vertical bar followed by a line-break,
  with the remaining description being indented.
```

### `references`

Links to reference documentation for a resource. Contains two attributes:

- `guides`: Link to quickstart in the API's Guides section
- `api`: Link to the REST API reference for the resource

Example:

```yaml
references:
  guides:
    'Create and connect to a database': 'https://cloud.google.com/alloydb/docs/quickstart/create-and-connect'
  api: 'https://cloud.google.com/alloydb/docs/reference/rest/v1/projects.locations.backups'
```

### `min_version: beta`
Marks the field (and any subfields) as beta-only. Ensure a beta version block
is present in provider.yaml.

### `docs`
Inserts styled markdown into the header of the resource's page in the provider
documentation. Can contain two attributes:

- `warning`: Warning text which will be displayed at the top of the resource docs on a yellow background.
- `note`: Note text which will be displayed at the top of the resource docs on a blue background.

Example:

```yaml
docs:
  warning: |
    This is a multi-line warning and will be
    displayed on a yellow background.
  note: |
    This is a multi-line note and will be
    displayed on a blue background.
```


## API interactions

### `base_url`

URL for the resource's [standard List method](https://google.aip.dev/132).
Terraform field names enclosed in double curly braces are replaced with
the field values from the resource at runtime.

```yaml
base_url: 'projects/{{project}}/locations/{{location}}/resourcenames'
```

### `self_link`

URL for the resource's [standard Get method](https://google.aip.dev/131).
Terraform field names enclosed in double curly braces are replaced with
the field values from the resource at runtime.

```yaml
self_link: 'projects/{{project}}/locations/{{location}}/resourcenames/{{name}}'
```

### `immutable`

If true, the resource and all its fields are considered immutable - that is,
only creatable, not updatable. Individual fields can override this if they
have a custom update method in the API.

See [Best practices: ForceNew](https://googlecloudplatform.github.io/magic-modules/best-practices/#forcenew) for more information.

Default: `false`

Example:

```yaml
immutable: true
```

### `timeouts`

Overrides one or more timeouts, in minutes. All timeouts default to 20.

Example:

```yaml
timeouts:
  insert_minutes: 40
  update_minutes: 40
  delete_minutes: 40
```

### `create_url`

URL for the resource's [standard Create method](https://google.aip.dev/133), including query parameters.
Terraform field names enclosed in double curly braces are replaced with
the field values from the resource at runtime.

Example:

```yaml
create_url: 'projects/{{project}}/locations/{{location}}/resourcenames?resourceId={{name}}'
```

### `create_verb`

Overrides the HTTP verb used to create a new resource.
Allowed values: `'POST'`, `'PUT'`, `'PATCH'`.

Default: `'POST'`

```yaml
create_verb: 'PATCH'
```

### `update_url`
Overrides the URL for the resource's [standard Update method](https://google.aip.dev/134).
If unset, the [`self_link` URL](#self_link) is used by default.
Terraform field names enclosed in double curly braces are replaced with
the field values from the resource at runtime.

```yaml
update_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{name}}'
```

### `update_verb`

The HTTP verb used to update a resource. Allowed values: `'POST'`, `'PUT'`, `'PATCH'`.

Default: `'PUT'`.

Example:

```yaml
update_verb: 'PATCH'
```

### `update_mask`

If true, the resource sets an `updateMask` query parameter listing modified
fields when updating the resource. If false, it doesn't.

Default: `false`

Example:

```yaml
update_mask: true
```

### `delete_url`

Overrides the URL for the resource's [standard Delete method](https://google.aip.dev/135).
If unset, the [`self_link` URL](#self_link) is used by default.
Terraform field names enclosed in double curly braces are replaced with
the field values from the resource at runtime.

Example:

```yaml
delete_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{name}}'
```

### `delete_verb`
Overrides the HTTP verb used to delete a resource.
Allowed values: `'POST'`, `'PUT'`, `'PATCH'`, `'DELETE'`.

Default: `'DELETE'`

Example:

```yaml
delete_verb: 'POST'
```

### `autogen_async`

If true, code for handling long-running operations is generated along with
the resource. If false, that code isn't generated and must be handwritten.

Default: `false`

```yaml
autogen_async: true
```

### `async`

Sets parameters for handling operations returned by the API. Can contain several attributes:

- `actions`: Overrides which API calls return operations. Default: `['create', 'update', 'delete']`
- `operation.base_url`: This should always be set to `'{{op_id}}'` unless you know that's wrong.
- `result.resource_inside_response`: If true, the provider sets the resource's Terraform ID after
  the resource is created, taking into account values that are set by the API at create time. This
  is only possible when the completed operation's JSON includes the created resource in the
  "response" field. If false, the provider sets the resource's Terraform ID before the resource is
  created, based only on the resource configuration. Default: `false`.

Example:

```yaml
async:
  actions: ['create', 'update', 'delete']
  operation:
    base_url: '{{op_id}}'
  result:
    resource_inside_response: true
```

## IAM resources

### `iam_policy`

Allows configuration of generated IAM resources. Supports the following common
attributes – for a full reference, see
[iam_policy.rb ↗](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/resource/iam_policy.rb):

- `parent_resource_attribute`: Name of the field on the terraform IAM resources
  which references the parent resource.
- `method_name_separator`: Character preceding setIamPolicy in the full URL for
  the API method. Usually `:`.
- `fetch_iam_policy_verb`: HTTP method for getIamPolicy. Usually `:POST`.
  Allowed values: `:GET`, `:POST`. Default: `:GET`
- `set_iam_policy_verb`: HTTP method for getIamPolicy. Usually `:POST`.
  Allowed values: :POST, :PUT. Default: :POST
- `import_format`: Must match the parent resource's `import_format` (or `self_link` if
  `import_format` is unset), but with the `parent_resource_attribute`
  value substituted for the final field.
- `allowed_iam_role`: Valid IAM role that can be set by generated tests. Default: `'roles/viewer'`
- `iam_conditions_request_type`: If IAM conditions are supported, set this attribute to indicate how the
  conditions should be passed to the API. Allowed values: `:QUERY_PARAM`,
  `:REQUEST_BODY`, `:QUERY_PARAM_NESTED`. Note: `:QUERY_PARAM_NESTED` should
  only be used if the query param field contains a `.`
- `min_version: beta`: Marks IAM support as beta-only.

Example:

```yaml
iam_policy:
  parent_resource_attribute: 'cloud_function'
  method_name_separator: ':'
  fetch_iam_policy_verb: :POST
  import_format:
    - 'projects/{{project}}/locations/{{location}}/resourcenames/{{cloud_function}}',
    - '{{cloud_function}}'
  allowed_iam_role: 'roles/viewer'
  iam_conditions_request_type: :REQUEST_BODY
  min_version: beta
```

## Resource behavior

### `custom_code`

Injects arbitrary logic into a generated resource. For more information, see [Add custom resource code]({{< ref "/develop/custom-code" >}}).

### `mutex`

All resources (of all kinds) that share a mutex value will block rather than
executing concurrent API requests. Terraform field names enclosed in double
curly braces are replaced with the field values from the resource at runtime.

Example:

```yaml
mutex: 'alloydb/instance/{{name}}'
```

## Fields

### `virtual_fields`

Contains a list of [virtual_fields]({{< ref "/develop/client-side-fields" >}}). By convention,
these should be fields that do not get sent to the API, and are instead used to modify
the behavior of a Terraform resource such as `deletion_protection`.

### `parameters`

Contains a list of [fields]({{< ref "/develop/field-reference" >}}). By convention,
these should be the fields that are part URL parameters such as `location` and `name`.

### `properties`

Contains a list of [fields]({{< ref "/develop/field-reference" >}}). By convention,
these should be fields that aren't part of the URL parameters.

Example:

```yaml
properties:
  - name: 'fieldOne'
    type: String
```
