---
title: "Metadata (meta.yaml) reference"
weight: 35
---

# Metadata (meta.yaml) reference

This page documents the properties available in meta.yaml files. These files make it easy to access data about the provider implementation without needing to compile the provider. They are used for coverage metrics and CAI-related overrides.

Meta.yaml files are auto-generated for MMv1 generated resources.

DCL and Handwritten resources will have handwritten meta.yaml files in the appropriate service directory in [mmv1/third_party/terraform/services/](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services). The file name will be `resource_PRODUCT_RESOURCE_meta.yaml(.tmpl)`. For example, `resource_compute_instance_meta.yaml.tmpl` Handwritten meta.yaml files with a `.tmpl` extension can use version guards (`{{- if ne $.TargetVersionName "ga" }}...{{- else}}...{{- end}}`) to exclude beta fields from the `google` provider.

All resources and fields must be present in meta.yaml files for the provider(s) they're available in.

## Required

### `resource`

The name of the Terraform resource. For example, "google_cloudfunctions2_function".

### `generation_type`

The generation method used to create the Terraform resource. For example, "mmv1", "dcl", "handwritten".

### `api_service_name`

The base name of the API used for this resource. For example, "cloudfunctions.googleapis.com".

### `api_version`

The version of the API used for this resource. For example, "v2".

### `api_resource_type_kind`

The API "resource type kind" used for this resource. For example, "Function".

## Optional

### `cai_asset_name_format`

The custom CAI asset name format for this resource is typically specified (for example, //cloudsql.googleapis.com/projects/{{project}}/instances/{{name}}). This should only have a value if it's different than the Terraform resource ID format.

### `api_variant_patterns`

The API URL patterns used by this resource that represent variants. For example, "folders/{folder}/feeds/{feed}". Each pattern must match the value defined in the API exactly. The use of `api_variant_patterns` is only meaningful when the resource type has multiple parent types available.

### `fields`

The list of fields used by this resource. Each field can contain the following attributes:

- `api_field`: Required for fields that aren't provider-only. The name of the field in the REST API, including the path. For example, "buildConfig.source.storageSource.bucket".
- `field`: The name of the field in Terraform, including the path. For example, "build_config.source.storage_source.bucket". Must be provided if and only if the field is provider-only or the Terraform field name can't be derived from the API name.
- `provider_only`: If true, the field is only present in the provider. This primarily applies for virtual fields and url-only parameters. When set to true, `field` should be set and `api_field` should be left empty. Default: `false`.
- `json`: If true, this is a JSON field which "covers" all child API fields. As a special case, JSON fields which cover an entire resource can have `api_field` set to `*`.
