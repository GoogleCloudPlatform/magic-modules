---
title: "Labels and annotations"
weight: 30
---

# Add labels and annotations support

{{< hint info >}}
**Note:** This page covers best practices guidance for the Terraform provider for Google Cloud, which is used to ensure a consistent UX for Terraform users across providers or GCP users across the Google provider. Generally, this guidance should be followed and exceptions should be clearly demarcated / discussed.
{{< /hint >}}

The new labels model and the new annotations model are introduced in [Terraform provider for Google Cloud 5.0.0](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/version_5_upgrade#provider).

There are now three label-related fields with the new labels model:
* The `labels` field is now non-authoritative and only manages the label keys defined in your configuration for the resource.
* The `terraform_labels` cannot be specified directly by the user. It merges the labels defined in the resource's configuration and the default labels configured in the provider block. If the same label key exists on both the resource level and provider level, the value on the resource will override the provider-level default.
* The output-only `effective_labels` will list all the labels present on the resource in GCP, including the labels configured through Terraform, the system, and other clients.

There are now two annotation-related fields with the new annotations model:
* The `annotations` field is now non-authoritative and only manages the annotation keys defined in your configuration for the resource.
* The output-only `effective_annotations` will list all the annotations present on the resource in GCP, including the annotations configured through Terraform, the system, and other clients.

This document describes how to add `labels` and `annotations` field to resources to support the new models.

## Labels support
When adding a new `labels` field, please make the changes below to support the new labels model. Otherwise, it has to wait for the next major release to make the changes.

### MMv1 resources

1. Use the type `KeyValueLabels` for the standard resource `labels` field. The standard resource `labels` field could be the top level `labels` field or the nested `labels` field inside the top level `metadata` field. Don't add `default_from_api: true` to this field or don't use this type for other `labels` fields in the resource. `KeyValueLabels` will add all of changes required for the new model automatically.

```yaml
 - name: 'labels'
   type: KeyValueLabels
   description: |
   The labels associated with this dataset. You can use these to
   organize and group your datasets.
```
2. In the handwritten acceptance tests, add `labels` and `terraform_labels` to `ImportStateVerifyIgnore` if `labels` field is in the configuration.

```go
ImportStateVerifyIgnore: []string{"labels", "terraform_labels"}, 
```
3. In the corresponding data source, after the resource read method, call the function `tpgresource.SetDataSourceLabels(d)` to make `labels` and `terraform_labels` have all of the labels on the resource.

```go
err = resourceArtifactRegistryRepositoryRead(d, meta)
if err != nil {
   return err
}

if err := tpgresource.SetDataSourceLabels(d); err != nil {
   return err
}
```

### Handwritten resources

1. Add `tpgresource.SetLabelsDiff`  to `CustomizeDiff` of the resource.
```go
CustomizeDiff: customdiff.All(
   tpgresource.SetLabelsDiff,
),
```
2. Add `labels` field and add more attributes (such as `ForceNew: true,`, `Set: schema.HashString,`) to this field if necessary.
```go
"labels": {
   Type:     schema.TypeMap,
   Optional: true,
   Elem:     &schema.Schema{Type: schema.TypeString},
   Description: `A set of key/value label pairs to assign to the project.
   
   **Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
   Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
},
```
3. Add output only field `terraform_labels` and add more attributes (such as `Set: schema.HashString,`) to this field if necessary. Don't add `ForceNew:true,` to this field.
```go
"terraform_labels": {
   Type:        schema.TypeMap,
   Computed:    true,
   Description: `The combination of labels configured directly on the resource and default labels configured on the provider.`,
   Elem:        &schema.Schema{Type: schema.TypeString},
},
```
4. Add output only field `effective_labels` and add more attributes (such as `ForceNew: true,`, `Set: schema.HashString,`) to this field if necessary.
```go
"effective_labels": {
   Type:        schema.TypeMap,
   Computed:    true,
   Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
   Elem:        &schema.Schema{Type: schema.TypeString},
},
```
5. In the create method, use the value of `effective_labels` in API request.
6. In the update method, use the value of `effective_labels` in API request.
7. In the read mehtod, set `labels`, `terraform_labels` and `effective_labels` to state.
```go
if err := tpgresource.SetLabels(res.Labels, d, "labels"); err != nil {
   return fmt.Errorf("Error setting labels: %s", err)
}
if err := tpgresource.SetLabels(res.Labels, d, "terraform_labels"); err != nil {
   return fmt.Errorf("Error setting terraform_labels: %s", err)
}
if err := d.Set("effective_labels", res.Labels); err != nil {
   return fmt.Errorf("Error setting effective_labels: %s", err)
}
```
8. In the handwritten acceptance tests, add `labels` and `terraform_labels` to `ImportStateVerifyIgnore`.
9. In the corresponding data source, after the resource read method, call the function `tpgresource.SetDataSourceLabels(d)` to make `labels` and `terraform_labels` have all of the labels on the resource.
10. Add the documentation for these label-related fields.

## Annotations support
When adding a new `annotations` field, please make the changes below below to support the new annotations model. Otherwise, it has to wait for the next major release to make the breaking changes.

### MMv1 resources

1. Use the type `KeyValueAnnotations` for the standard resource `annotations` field. The standard resource `annotations` field could be the top level `annotations` field or the nested `annotations` field inside the top level `metadata` field. Don't add `default_from_api: true` to this field or don't use this type for other `annotations` fields in the resource. `KeyValueAnnotations` will add all of changes required for the new model automatically.

```yaml
- name: 'annotations'
  type: KeyValueAnnotations
  description: |
   Client-specified annotations. This is distinct from labels.
```
2. In the handwritten acceptance tests, add `annotations` to `ImportStateVerifyIgnore` if `annotations` field is in the configuration.

```go
ImportStateVerifyIgnore: []string{"annotations"},
```
3. In the corresponding data source, after the resource read method, call the function `tpgresource.SetDataSourceAnnotations(d)` to make `annotations` have all of the annotations on the resource.

```go
err = resourceSecretManagerSecretRead(d, meta)
if err != nil {
   return err
}

if err := tpgresource.SetDataSourceLabels(d); err != nil {
   return err
}

if err := tpgresource.SetDataSourceAnnotations(d); err != nil {
   return err
}
```

### Handwritten resources

1. Add `tpgresource.SetAnnotationsDiff`  to `CustomizeDiff` of the resource.
2. Add `annotations` field and add more attributes (such as `ForceNew: true,`, `Set: schema.HashString,`) to this field if necessary.
3. Add output only field `effective_annotations` and add more attributes (such as `ForceNew: true,`, `Set: schema.HashString,`) to this field if necessary.
4. In the create method, use the value of `effective_annotations` in API request.
5. In the update method, use the value of `effective_annotations` in API request.
6. In the read mehtod, set `annotations`, and `effective_annotations` to state.
7. In the handwritten acceptance tests, add `annotations` to `ImportStateVerifyIgnore`.
8. In the corresponding data source, after the resource read method, call the function `tpgresource.SetDataSourceAnnotations(d)` to make `annotations` have all of the labels on the resource.
9. Add the documentation for these annotation-related fields.
