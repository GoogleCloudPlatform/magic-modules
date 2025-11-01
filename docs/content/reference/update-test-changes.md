---
title: "Test Template Migration Guide"
weight: 50
---
# Test Template Migration Guide

The test generation framework in Magic Modules has been updated to natively support multi-step tests (create and update) directly from the resource's YAML definition. This guide is for experienced contributors who are familiar with the previous `examples` based workflow and provides a concise overview of what has changed.

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
```

New `samples` format

```yaml
samples:
  - name: "pubsub_topic_basic"
    primary_resource_id: "example"
    steps:
      - name: "create" # The name for this step
        prefixed_vars:
          topic_name: "example-topic"
```

### Create and Update Test Comparison

Previously, update tests had to be handwritten. Now, they can be defined by adding a second step to a sample.

Old examples format (Create only)

```yaml
examples:
  - name: "pubsub_topic_full"
    primary_resource_id: "example"
    vars:
      topic_name: "example-topic"
      label_key: "key-one"
      label_value: "value-one"
```
(An update test would have required a separate handwritten Go test file.)

New samples format (Create and Update)

```yaml
samples:
  - name: "pubsub_topic_update"
    primary_resource_id: "example"
    steps:
      - name: "create"
        prefixed_vars:
          topic_name: "example-topic"
        vars:
          label_key: "key-one"
          label_value: "value-one"
      - name: "update"
        prefixed_vars:
          topic_name: "example-topic" 
        vars:
          label_key: "key-one-updated"
          label_value: "value-one-updated"
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
| `examples.primary_resource_name` | `sample.primary_resource_name` | Remains at the sample level. |
| `examples.bootstrap_iam` | `sample.bootstrap_iam` | Remains at the sample level. |
| `examples.min_version` | `sample.min_version` | Remains at the sample level. |
| `examples.exclude_test` | `sample.exclude_test` | Remains at the sample level. |
| `examples.region_override` | `sample.region_override` | Remains at the sample level. |
| `examples.skip_vcr` | `sample.skip_vcr` | Remains at the sample level. |
| `examples.skip_test` | `sample.skip_test` | Remains at the sample level. |
| `examples.external_providers` | `sample.external_providers` | Remains at the sample level. |
| `examples.tgc_skip_test` | `sample.tgc_skip_test` | Remains at the sample level. |

---

### Fields Mapped to `Step`
These fields are now configured within each individual `Step` of a `Sample`.

| Old Field | New Location | Notes |
| :--- | :--- | :--- |
| `examples.vars` | `step.prefixed_vars` | The entire `vars` map is moved into `prefixed_vars` within each step. |
| `examples.test_env_vars` | `step.test_env_vars` | Moved to the step level. |
| `examples.test_vars_overrides` | `step.test_vars_overrides` | Moved to the step level. |
| `examples.oics_vars_overrides` | `step.oics_vars_overrides` | Moved to the step level. |
| `examples.ignore_read_extra` | `step.ignore_read_extra` | Moved to the step level. |
| `examples.exclude_docs` | `step.exclude_docs` | Moved to the step level. |
| `examples.exclude_import_test` | `step.exclude_import_test` | Moved to the step level. |
| `examples.config_path` | `step.config_path` | Path is updated to the new service-specific directory within the step. |

---

### New Fields
This field is new in the `steps` object and has no direct equivalent in the old `examples` structure.

| Old Field | New Location | Notes |
| :--- | :--- | :--- |
| *(N/A)* | `step.vars` | Newly added at the step level. Values are copied directly to tests. |
| *(N/A)* | `step.min_version` | Newly added to set a version for a specific step |


## Template `.tf.`tmpl File Changes

The location for template files has moved from `templates/terraform/examples/` to a service-specific directory under `templates/terraform/samples/services/`.

Additionally, the variable object passed into the templates has been updated. `$.PrefixedVars` will append `tf-test` prefixes and random string suffixes, which is used for resource identifiers in most cases. `$.Vars` will apply plain values from the YAML configuration.

### Example .tf.tmpl variables

Old template `pubsub_topic_basic.tf.tmpl` (in `templates/terraform/examples/`)

```tf
resource "google_pubsub_topic" "{{$.PrimaryResourceId}}" {
  name = "{{index $.Vars "topic_name"}}"

  labels = {
    foo = "bar"
  }
}
```

New template `pubsub_topic_basic.tf.tmpl` (in `templates/terraform/samples/services/pubsub/`)

```tf
resource "google_pubsub_topic" "{{$.PrimaryResourceId}}" {
  name = "{{index $.PrefixedVars "topic_name"}}"

  labels = {
    "{{index $.Vars "label_key"}}" = "{{index $.Vars "label_value"}}"
  }
}
```