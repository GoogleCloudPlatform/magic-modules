---
title: "Test Template Migration Guide"
weight: 50
---
# Test Template Migration Guide

The test generation framework in Magic Modules has been updated to natively support multi-step tests (create and update) directly from the resource's YAML definition. This guide is for experienced contributors who are familiar with the previous `examples` based workflow and provides a concise overview of what has changed.

## Transition Timeline

*   **Recommended Path:** The new `samples` configuration block is the recommended path for adding tests.
*   **Legacy Support:** The old `examples` block path is still supported until Mid May 2026.
*   **Action Required:** Contributors are encouraged to use `samples` for all new tests and migrate existing `examples` to `samples`.



## YAML Changes: From `examples` to `samples` and `steps`

The most significant change is the replacement of the `examples` block with a new `samples` block. Each `sample` can now contain one or more `steps`, where the first step corresponds to the resource creation and any subsequent steps correspond to updates.

This new structure eliminates the need for handwritten update tests for MMv1 resources, as they can now be generated automatically.

### Create Test Comparison

A simple create test is now defined as a sample with a single step.

Old `examples` format

```yaml
examples:
  - name: "pubsub_topic_basic"
    primary_resource_id: "example"
    vars:
      topic_name: "example-topic"
      field_1: "value-one"
      field_2: "value-two"
    test_vars_overrides:
      field_1: "value-one"
      field_2: "value-two"
```

New `samples` format

```yaml
samples:
  - name: "pubsub_topic_basic"
    primary_resource_id: "example"
    steps:
      - name: "pubsub_topic_basic" 
        resource_id_vars:
          topic_name: "example-topic"
        vars:
          field_1: "value-one"
          field_2: "value-two"
```

> [!NOTE]
> ### Why `test_vars_overrides` is no longer needed
> In the old `examples` block, any key placed under `vars` was automatically appended with a random suffix (for example, `example-topic-12345`). If you wanted to pass a plain value (without random suffixes), you had to use `test_vars_overrides`.
>
> In the new `samples` format, these two use cases are explicitly separated into different fields within a step:
> * **`resource_id_vars`**: **Use this exclusively for resource identifier variables (such as resource names or IDs).** It will automatically prepend a `tf-test` (or `tf_test`) prefix and append a random suffix. If a resource identifier doesn't support hyphens `-` or underscores `_`, use `test_vars_overrides` instead. For non-identifier variables, use `vars`.
> * **`vars`**: Use this for plain literal values that should be passed to the test exactly as written (replaces the need for `test_vars_overrides`). **Note:** This should ONLY be used for fields that vary between steps (for example, to test update functionality). Constant values should be hardcoded directly in the `.tf.tmpl` file.

### Update Test Comparison

Previously, update tests had to be handwritten. Now, they can be defined by adding a second step to a sample.

New samples format (Create and Update)

```yaml
samples:
  - name: "pubsub_topic_update"
    primary_resource_id: "example"
    steps:
      - name: "pubsub_topic_full"
        resource_id_vars:
          topic_name: "example-topic"
        vars:
          field_1: "value-one"
          field_2: "value-two"
      - name: "pubsub_topic_full"
        resource_id_vars:
          topic_name: "example-topic" 
        vars:
          field_1: "value-one-updated"
          field_2: "value-two-updated"
```

## YAML Field Migration

The migration moves fields from the old `Examples` structure to either the new top-level `Sample` or the nested `Step`. The following tables detail where each field has been moved.

### Fields Mapped to `Sample` (Top-Level)
These fields remain at the top level, moving from the old `examples` object to the new `sample` object.

| Old Field | New Location | Notes |
| :--- | :--- | :--- |
| `examples.name` | `sample.name` | Used for the overall sample name. A `step.name` is also used for each step. |
| `examples.primary_resource_id` | `sample.primary_resource_id` | Remains at the sample level. |
| `examples.primary_resource_type` | `sample.primary_resource_type` | Remains at the sample level. |
| `examples.bootstrap_iam` | `sample.bootstrap_iam` | Remains at the sample level. |
| `examples.min_version` | `sample.min_version` | Remains at the sample level. |
| `examples.exclude_test` | `sample.exclude_test` | Remains at the sample level. |
| `examples.region_override` | `sample.region_override` | Remains at the sample level. |
| `examples.skip_vcr` | `sample.skip_vcr` | Remains at the sample level. |
| `examples.skip_test` | `sample.skip_test` | Remains at the sample level. |
| `examples.skip_func` | `sample.skip_func` | Remains at the sample level. |
| `examples.external_providers` | `sample.external_providers` | Remains at the sample level. |
| `examples.tgc_skip_test` | `sample.tgc_skip_test` | Remains at the sample level. |

---

### Fields Mapped to `Step`
These fields are now configured within each individual `Step` of a `Sample`.

| Old Field | New Location | Notes |
| :--- | :--- | :--- |
| `examples.vars` | `step.resource_id_vars` | The entire `vars` map is moved into `resource_id_vars` within each step. |
| `examples.test_env_vars` | `step.test_env_vars` | Moved to the step level. |
| `examples.test_vars_overrides` | `step.test_vars_overrides` | Moved to the step level. |
| `examples.oics_vars_overrides` | `step.oics_vars_overrides` | Moved to the step level. |
| `examples.ignore_read_extra` | `step.ignore_read_extra` | Moved to the step level. |
| `examples.exclude_docs` | `sample.exclude_basic_doc` | Replaced by `sample.exclude_basic_doc` (or `step.include_step_doc` to override). |
| `examples.exclude_import_test` | `step.exclude_import_test` | Moved to the step level. |
| `examples.config_path` | `step.config_path` | Path is updated to the new service-specific directory within the step. |

---

### New Fields
This field is new in the `steps` object and has no direct equivalent in the old `examples` structure.

| Old Field | New Location | Notes |
| :--- | :--- | :--- |
| *(N/A)* | `step.vars` | Newly added at the step level. Values are copied directly to tests. |
| *(N/A)* | `step.min_version` | Newly added to set a version for a specific step |
| *(N/A)* | `step.include_step_doc` | Explicitly forcing a step's docs generation (overriding Sample block). |


## Template `.tf.`tmpl File Changes

The location for template files has moved from `templates/terraform/examples/` to a service-specific directory under `templates/terraform/samples/services/`.

Additionally, the variable object passed into the templates has been updated. `$.ResourceIdVars` will append `tf-test` prefixes and random string suffixes, which is used for resource identifiers in most cases. `$.Vars` will apply plain values from the YAML configuration.

### Example .tf.tmpl variables

Old template `pubsub_topic_basic.tf.tmpl` (in `templates/terraform/examples/`)

```tf
resource "google_pubsub_topic" "{{$.PrimaryResourceId}}" {
  name = "{{index $.Vars "topic_name"}}"

  field_1 = "{{index $.Vars "field_1"}}"
  field_2 = "{{index $.Vars "field_2"}}"
}
```

New template `pubsub_topic_basic.tf.tmpl` (in `templates/terraform/samples/services/pubsub/`)

```tf
resource "google_pubsub_topic" "{{$.PrimaryResourceId}}" {
  name = "{{index $.ResourceIdVars "topic_name"}}"

  field_1 = "{{index $.Vars "field_1"}}"
  field_2 = "{{index $.Vars "field_2"}}"
}
```