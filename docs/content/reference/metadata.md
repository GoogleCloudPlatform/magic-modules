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

### `api_variant_patterns`

The API URL patterns used by this resource that represent variants e.g., "folders/{folder}/feeds/{feed}". Each pattern must match the value defined in the API exactly. The use of `api_variant_patterns` is only meaningful when the resource type has multiple parent types available.

### `fields`

The list of fields used by this resource. Each field can contain the following attributes:

- `field`: The name of the field in Terraform, including the path e.g., "build_config.source.storage_source.bucket"
- `api_field`: The name of the field in the API, including the path e.g., "build_config.source.storage_source.bucket". Defaults to the value of `field`.
- `provider_only`: If true, the field is only present in the provider. This primarily applies for virtual fields and url-only parameters. When set to true, `api_field` should be left empty, as it will be ignored. Default: `false`.
