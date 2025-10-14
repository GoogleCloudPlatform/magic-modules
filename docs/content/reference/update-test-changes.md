--
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

## Template `.tf.`tmpl File Changes

The location for template files has moved from `templates/terraform/examples/` to a service-specific directory under `templates/terraform/samples/services/`.

Additionally, the variable object passed into the templates has been updated. $.PrefixedVars will append tf-test prefixes and random string suffixes, which is used for resource identifiers in most cases. $.Vars will apply plain values from the YAML configuration.

### Example .tf.tmpl variables
Old template `pubsub_topic_basic.tf.tmpl` (in `templates/terraform/examples/`)

```
resource "google_pubsub_topic" "{{$.PrimaryResourceId}}" {
  name = "{{index $.Vars "topic_name"}}"

  labels = {
    foo = "bar"
  }
}
```

New template `pubsub_topic_basic.tf.tmpl` (in `templates/terraform/samples/services/pubsub/`)

```
resource "google_pubsub_topic" "{{$.PrimaryResourceId}}" {
  name = "{{index $.PrefixedVars "topic_name"}}"

  labels = {
    "{{index $.Vars "label_key"}}" = "{{index $.Vars "label_value"}}"
  }
}
```