---
title: "MMv1 resource reference"
weight: 10
aliases:
  - /reference/resource-reference
  - /reference/iam-policy-reference
  - /develop/resource-reference
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
only creatable, not updatable. Individual fields can override this for themselves and
their subfields with [`update_url`]({{< ref "/reference/field#update_url" >}})
if they have a custom update method in the API.

See [Best practices: Immutable fields]({{< ref "/best-practices/immutable-fields/" >}}) for more information.

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

### `exclude_delete`
If true, deleting the resource will only remove it from the Terraform state and will not call an API. If false, deleting the resource will run the standard deletion behavior and/or any [custom code]({{< ref "/develop/custom-code" >}}) related to deletion.
This should be used if the resource can never be deleted in the API, and there is no other reasonable action to take on deletion. See [Deletion behaviors]({{< ref "/best-practices/deletion-behaviors" >}}) for more information.

```yaml
exclude_delete: true
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

### `error_retry_predicates`

An array of function names that determine whether an error is retryable.

```yaml
error_retry_predicates:
  - 'transport_tpg.IamMemberMissing'
```

### `error_abort_predicates`

An array of function names that determine whether an error is not retryable.

```yaml
error_abort_predicates:
  - 'transport_tpg.Is429QuotaError'
```

## IAM resources

### `iam_policy`

Allows configuration of generated IAM resources. Supports the following common
attributes – for a full reference, see
[iam_policy.rb ↗](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/resource/iam_policy.go):

- `parent_resource_attribute`: Name of the field on the terraform IAM resources
  which references the parent resource.
- `method_name_separator`: Character preceding setIamPolicy in the full URL for
  the API method. Usually `:`.
- `fetch_iam_policy_verb`: HTTP method for getIamPolicy. Usually `'POST'`.
  Allowed values: `'GET'`, `'POST'`. Default: `'GET'`
- `set_iam_policy_verb`: HTTP method for getIamPolicy. Usually `'POST'`.
  Allowed values: `'POST'`, `'PUT'`. Default: `'POST'`
- `import_format`: Must match the parent resource's `import_format` (or `self_link` if
  `import_format` is unset), but with the `parent_resource_attribute`
  value substituted for the final field.
- `allowed_iam_role`: Valid IAM role that can be set by generated tests. Default: `'roles/viewer'`
- `iam_conditions_request_type`: If IAM conditions are supported, set this attribute to indicate how the
  conditions should be passed to the API. Allowed values: `'QUERY_PARAM'`,
  `'REQUEST_BODY'`, `'QUERY_PARAM_NESTED'`. Note: `'QUERY_PARAM_NESTED'` should
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

## Sweeper

Sweepers are a testing infrastructure mechanism that automatically clean up resources created during tests. They run before tests start and can be run manually to clean up dangling resources. Sweepers help prevent test failures due to resource quota limits and reduce cloud infrastructure costs by removing test resources that were not properly cleaned up.

Sweeper generation is enabled by default, except in the following conditions which require customization here:

- Resources with custom deletion code
- Resources with parent-child relationships (unless the parent relationship is configured)
- Resources with complex URL parameters that aren't simple region/project parameters

Define the sweeper block in a resource to override these exclusions and enable sweeper generation for that resource.

### `exclude_sweeper`

If set to `true`, no sweeper will be generated for this resource. This is useful for resources that cannot or should not be automatically cleaned up.

Default: `false`

Example:

```yaml
exclude_sweeper: true
```

### `sweeper`

Configures how test resources are swept (cleaned up) after tests. The sweeper system helps ensure resources created during tests are properly removed, even when tests fail unexpectedly. All fields within the `sweeper` block are optional, with reasonable defaults provided when not specified. See [sweeper.go ↗](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/resource/sweeper.go) for the implementation.

- `identifier_field`: Specifies which field in the API resource object should be used to identify resources for deletion. If not specified, defaults to "name" if present in the resource, otherwise falls back to "id".

- `prefixes`: Specifies name prefixes that identify resources eligible for sweeping. Resources whose names start with any of these prefixes will be deleted. By default, resources with the `tf-test-` prefix are automatically eligible for sweeping even if no prefixes are specified.

- `url_substitutions`: Allows customizing URL parameters when listing resources. Each map entry represents a set of key-value pairs to substitute in the URL template. This is commonly used to specify regions to sweep in. If not specified, the sweeper will only run in the default region (us-central1) and zone (us-central1-a).

- `dependencies`: Lists other resource types that must be swept before this one. This ensures proper cleanup order for resources with dependencies. If not specified, no dependencies are assumed.

- `parent`: Configures sweeping for resources that depend on parent resources (like a nodepool that belongs to a cluster).

  Required fields:
  - `resource_type`: The type of the parent resource (for example, "google_container_cluster")
  - `child_field`: The field in your resource that references the parent (for example, "cluster")
  - At least one of `parent_field` or `template` is required

  Options for getting parent reference:
  - `parent_field`: The field from parent to use (typically "name" or "id")
  - `template`: A template string like "projects/{{project}}/locations/{{location}}/clusters/{{value}}" 
  
  Options for processing the parent field value:
  - `parent_field_extract_name`: When set to true, extracts just the resource name from a self_link URL by taking the portion after the last slash. This is useful when the parent field contains a fully-qualified resource URL (like "https://www.googleapis.com/compute/v1/projects/my-project/global/networks/my-network" or "projects/my-project/zones/us-central1-a/instances/my-instance") but you only need the final resource name component ("my-network" or "my-instance").
  
  - `parent_field_regex`: A regex pattern with a capture group to extract a specific portion of the parent field value. This is useful when you need more control over extracting parts of complex resource identifiers. The pattern must contain at least one capture group (in parentheses), and the first capture group's match will be used as the extracted value.

- `query_string`: Allows appending additional query parameters to the resource's delete URL when performing delete operations. Format should include the starting character, for example, "?force=true" or "&verbose=true". If not specified, no additional query parameters are added.

- `ensure_value`: Specifies a field that must be set to a specific value before deletion. Used for resources that have fields like 'deletionProtectionEnabled' that must be explicitly disabled before the resource can be deleted. All fields within the `ensure_value` block are required except `include_full_resource`:
  
  - `field`: The API field name that needs to be updated before deletion. Can include dot notation for nested fields (for example, "settings.deletionProtectionEnabled").
  
  - `value`: The required value that `field` must be set to before deletion. For boolean fields use "true" or "false", for integers use string representation, for string fields use the exact string value required.
  
  - `include_full_resource`: Determines whether to send the entire resource object with the updated field (true) or to send just the field that needs updating (false). Some APIs require the full resource to be sent in update operations. Default: `false`.

Examples:

Basic sweeper configuration:

```yaml
sweeper:
  prefixes:
    - "tf-test-"
    - "tmp-"
```

Sweeper with parent-child relationship (basic):

```yaml
sweeper:
  dependencies: # sweep google_compute_instance before attempting to sweep this resource
    - "google_compute_instance"
  parent:
    resource_type: "google_container_cluster"
    parent_field: "name"
    child_field: "cluster"
```

Sweeper with parent_field_extract_name:

```yaml
sweeper:
  parent:
    resource_type: "google_compute_network"
    parent_field: "selfLink"  # Contains: "projects/my-project/global/networks/my-network"
    parent_field_extract_name: true  # Extracts just "my-network"
    child_field: "network"
```

Sweeper with parent template and pre-deletion field update:

```yaml
sweeper:
  parent:
    resource_type: "google_container_cluster"
    template: "projects/{{project}}/locations/{{location}}/clusters/{{value}}"
    parent_field: "displayName"  # Provides the value for {{value}}
    child_field: "cluster"
  ensure_value:
    field: "deletionProtection"
    value: "false"
    include_full_resource: false

```

Sweeper with URL substitutions for multiple regions:

```yaml
sweeper:
  url_substitutions:
    - collection_id: default_collection
      region: global
    - collection_id: default_collection
      region: eu
```

Sweeper with URL substitutions specifying only regions:

```yaml
sweeper:
  identifier_field: "displayName"
  url_substitutions:
    - region: "us-central1"
    - region: "us-east1"
    - region: "europe-west1"
```

## Fields

### `virtual_fields`

Contains a list of [virtual_fields]({{< ref "/develop/client-side-fields" >}}). By convention,
these should be fields that do not get sent to the API, and are instead used to modify
the behavior of a Terraform resource such as `deletion_protection`.

### `parameters`

Contains a list of [fields]({{< ref "/reference/field" >}}). By convention,
these should be the fields that are part URL parameters such as `location` and `name`.

### `properties`

Contains a list of [fields]({{< ref "/reference/field" >}}). By convention,
these should be fields that aren't part of the URL parameters.

Example:

```yaml
properties:
  - name: 'fieldOne'
    type: String
```

## Examples

### `examples`

A list of configurations that are used to generate documentation and tests. Each example supports the following common
attributes – for a full reference, see
[examples.go ↗](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/resource/examples.go):

- `name`: snake_case name of the example. This corresponds to the configuration file in
  [mmv1/templates/terraform/examples](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/examples) (excluding the `.go.tmpl` suffix) and is used to generate the test name and the documentation header.
- `primary_resource_id`: The id of the resource under test. This is used by tests to automatically run additional checks.
  Configuration files should reference this to avoid getting out of sync. For example:
  `resource "google_compute_address" ""{{$.PrimaryResourceId}}" {`
- `bootstrap_iam`: specify member/role pairs that should always exist. `{project_number}` will be replaced with the
  default project's project number, and `{organization_id}` will be replaced with the "target" test organization's ID. This avoids race conditions when modifying the global IAM permissions.
  Permissions attached to resources created _in_ a test should instead be provisioned with standard terraform resources.
- `vars`: Key/value pairs of variables to inject into the configuration file. These can be referenced in the configuration file
  with `{{index $.Vars "key"}}`. All resource IDs (even for resources not under test) should be declared with variables that
  contain a `-` or `_`; this will ensure that, in tests, the resources are created with a `tf-test` prefix to allow automatic cleanup of dangling resources and a random suffix to avoid name collisions.
- `test_env_vars`: Key/value pairs of variable names and special values indicating variables that should be pulled from the
  environment during tests. These will receive a neutral default value in documentation. Common special values include:
  `PROJECT_NAME`, `REGION`, `ORG_ID`, `ORG_TARGET` (a separate test org for testing certain org-level resources such as IAM), `BILLING_ACCT`, `SERVICE_ACCT` (the test runner service account).
- `test_vars_overrides`: Key/value pairs of literal overrides for variables used in tests. This can be used to call functions to
  generate or determine a variable's value.
- `min_version`: Set this to `beta` if the resource is in the `google` provider but the example will only work with the
  `google-beta` provider (for example, because it includes a beta-only field.)
- `ignore_read_extra`: Properties to not check on import. This should be used in cases where a property will not be set on import,
  for example write-only fields.
- `exclude_test`: If set to `true`, no test will be generated based on this example.
- `exclude_docs`: If set to `true`, no documentation will be generated based on this example.
- `exclude_import_test`: If set to `true`, no import test will be generated for this example.
- `skip_vcr`: See [Skip tests in VCR replaying mode]({{< ref "/test/test#skip-vcr" >}}) for more information about this flag.
- `skip_test`: If not empty, the test generated based on this example will always be skipped. In most cases, the value should be a
  link to a ticket explaining the issue that needs to be resolved before the test can be unskipped.
- `external_providers`: A list of external providers that are needed for the testcase. This does add some latency to the testcase,
  so only use if necessary. Common external providers: `random`, `time`.

Example:

```yaml
examples:
  - name: service_resource_basic
    primary_resource_id: example
    bootstrap_iam:
      - member: "serviceAccount:service-{project_number}@gcp-sa-healthcare.iam.gserviceaccount.com"
        role: "roles/bigquery.dataEditor"
      - member: "serviceAccount:service-org-{organization_id}@gcp-sa-osconfig.iam.gserviceaccount.com"
        role: "roles/osconfig.serviceAgent"
    vars:
      dataset_id: "my-dataset"
      network_name: "my-network"
    test_env_vars:
      org_id: "ORG_ID"
    test_vars_overrides:
      network_name: 'acctest.BootstrapSharedServiceNetworkingConnection(t, "service-resource-network-config")'
    min_version: "beta"
    ignore_read_extra: 
      - 'foo'
    exclude_test: true
    exclude_docs: true
    exclude_import_test: true
    skip_vcr: true
    skip_test: "https://github.com/hashicorp/terraform-provider-google/issues/20574"
    external_providers:
      - "time"
```