---
title: "MMv1 metadata reference"
weight: 35
---

# MMv1 metadata reference

This page documents all properties for metadata. Metadata does not impact the provider itself, but is used by Google internally for coverage metrics.

## Required

### `resource`

The name of the Terraform resource e.g., "google_cloudfunctions2_function".

### `generation_type`

The generation method used to create the Terraform resource e.g., "mmv1", "dcl", "handwritten".

## Optional

### `api_service_name`

The base name of the API used for this resource e.g., "cloudfunctions.googleapis.com".

### `api_version`

The version of the API used for this resource e.g., "v2".

### `api_resource_type_kind`

The API "resource type kind" used for this resource e.g., "Function".

### `cai_asset_name_format`

The custom CAI asset name format for this resource is typically specified (e.g., //cloudsql.googleapis.com/projects/{{project}}/instances/{{name}}). If this format is not provided, the Terraform resource ID format is used instead.

### `api_variant_patterns`

The API URL patterns used by this resource that represent variants e.g., "folders/{folder}/feeds/{feed}". Each pattern must match the value defined in the API exactly. The use of `api_variant_patterns` is only meaningful when the resource type has multiple parent types available.

### `fields`

The list of fields used by this resource. Each field can contain the following attributes:

- `api_field`: The name of the field in the REST API, including the path e.g., "buildConfig.source.storageSource.bucket".
- `field`: The name of the field in Terraform, including the path e.g., "build_config.source.storage_source.bucket". Defaults to the value of `api_field` converted to snake_case.
- `provider_only`: If true, the field is only present in the provider. This primarily applies for virtual fields and url-only parameters. When set to true, `field` should be set and `api_field` should be left empty. Default: `false`.
- `json`: If true, this is a JSON field which "covers" all child API fields. As a special case, JSON fields which cover an entire resource can have `api_field` set to `*`.
