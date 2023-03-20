# DCL/tpgtools for MMv1 developers

The move from Magic Modules Version 1 (MMv1) to the MMv2 project means that how
you work on resources has changed. Instead of a single generator constructing
TPG resources directly, resources are added to the Declarative Client
Library (DCL) in an internal generator and those DCL resources are added to
TPG using this generator, `tpgtools`.

Naively the mapping is simple- definitons done in `api.yaml` belong to the DCL
and those in `terraform.yaml` to `tpgtools`. However `api.yaml` often wasn't
sufficient to build a model of a complete declarative resource. As such, many
overrides present in `terraform.yaml` have been hoisted to the DCL layer as well.

## api.yaml Definitions

All definitions in `api.yaml` are the responsibility of the DCL. They will be
exposed in `tpgtools` through the DCL's OpenAPI definitions. Some notable ones
include:

* Product:
    * `async`: Handled by the DCL. The DCL will expose a consistent interface
regardless of whether a given product/resource implements long-running
operations.
    * `versions`: See Versions below. Also http://b/159133756.
* Resource:
    * `min_version`: No longer present. See Versions below.
    * `input`: The DCL will annotate every field with `x-kubernetes-immutable`.
    * `identity`: The DCL exposes a series of `paths` in the OpenAPI that
correspond to the DCL's methods. `x-dcl-id` is also available, as a suggested
storage id for the resource.
    * `update_mask`: Handled by the DCL.
    * `*_url`: Handled by the DCL.
    * `nested_query`: Handled by the DCL. Fine-grained resources will appear as
top-level resources in many cases.
* Field:
    * `default`: Exposed in the DCL OpenAPI spec as `default`.
    * `description`: Exposed in the DCL OpenAPI spec as `description`.
    * `deprecation_message`: May be exposed by DCL in the future. See
Deprecation / Removed Messages below.
    * `output`: Exposed in the DCL OpenAPI spec as `readOnly: true`.
    * `input`: Exposed in the DCL OpenAPI spec as `x-kubernetes-immutable: true`.
    * `url_param_only`: Handled by the DCL. Exposed as normal fields.
    * `required`: Exposed in the DCL OpenAPI spec as `required` on the parent.
    * `update_*`, `fingerprint_name`, `send_empty_value`: Handled by the DCL.
    * `allow_empty_object`: Unknown. Likely handled by the DCL.
    * `min_version`, `exact_version`: See Versions below.
    * `conflicts`, `at_least_one_of`, `exactly_one_of`: Will be exposed by DCL
in the future. http://b/159243366
    * `new_type`: Must be handled with `tpgtools`.
    * `pattern`: Handled in the DCL.
* Array (Field)
    * `item_type`: Exposed in the DCL OpenAPI spec as `items`.
    * `min_size`: `max_size`: Will be exposed by DCL in the future.
* Enum (Field)
    * `values`: Exposed in the DCL OpenAPI spec as `enum`.
* Map (Field):
    * No examples present in the DCL yet. `tpgtools` must add support when
required.

### Versions

In MMv1 all fields are included in a single definition, annotated with a
`min_version` when needed. The MMv1 generator excludes fields that aren't
appropriate for a given version, and can determine the version of an individual
field.

With the DCL, fields aren't marked on a field-by-field basis. Instead, resource
definitions are packaged together into versions. For example in `compute` there
is a base `compute` directory containing GA definitions of resources, and a
`compute/beta` directory containing the beta definitions (which are a superset
of GA).

### Deprecation / Removed Messages

Deprecation messages are not currently implemented by the DCL, but may be used
in the future. When used, it's likely they'll be added with an OpenAPI
specification extension.

Removal messages are unlikely to be used in the DCL. Instead, the field will
be removed entirely.

However if TPG has deprecated/removed a field locally, adding deprecation or
removal messages will be a responsibility of `tpgtools`.

## terraform.yaml Definitions

* ResourceOverride:
    * `filename_override`: No equivalent. Could be implemented in `tpgtools`.
    * `legacy_name`: `"CUSTOM_RESOURCE_NAME"` in tpgtools.
    * `id_format`: `"CUSTOM_ID"` in `tpgtools`. Most existing custom formats
should be available directly through `x-dcl-id` from the DCL.
    * `import_format`: `"IMPORT_FORMAT"` in `tpgtools`.
    * `mutex`: `MUTEX` in `tpgtools`. Not supported in DCL yet, but likely
should be: http://b/164181606.
    * `examples`: Needs support in `tpgtools`: http://b/158508947. Will also
likely appear in the DCL in the future, but will need support in both.
    * `virtual_fields`: `"VIRTUAL_FIELD"` in `tpgtools`.
    * `autogen_async`: Handled in the DCL.
    * `exclude_import`: Requires a future `tpgtools` override.
    * `exclude_validator`: Terraform Validator is not supported by `tpgtools`.
    * `timeouts`: Handled by DCL. Can be implemented as a `tpgtools` override if
the two deviate.
    * `error_retry_predicates`: Handled by the DCL. The DCL retry *could* be
overriden in theory, but this should not be done.
    * `schema_version`: Requires a future `tpgtools` override.
    * `skip_sweeper`: Sweepers are not currently generated by `tpgtools`. We
need to determine whether the DCL will support them directly or not:
http://b/164185612.
    * `skip_delete`: Resources that can't be deleted should not have a delete
entry in the `paths` object.
    * `supports_indirect_user_project_override`: Will be supported by the DCL.
See http://164187236.
    * `read_error_transform`: Handled by the DCL.
* PropertyOverride:
    * `diff_suppress_func`: `"DIFF_SUPPRESS_FUNC"` in `tpgtools`. Some DSFs may
be added automatically based on DCL OpenAPI annotations.
    * `state_func`: Requires a future override in `tpgtools`.
    * `sensitive`: `x-dcl-sensitive` in DCL OpenAPI annotations.
    * `ignore_read`: `"IGNORE_READ"` in `tpgtools`. The DCL will return the
user-provided value in Apply.
    * `validation`: `"CUSTOM_VALIDATION"` in `tpgtools`. Will likely be exposed
by the DCL as well in the future.
    * `unordered_list`: Unknown. Possibly `x-dcl-list-type: set` in the DCL
OpenAPI spec.
    * `is_set`: `x-dcl-list-type: set` in the DCL OpenAPI spec.
    * `set_hash_func`: `"SET_HASH_FUNC"` in `tpgtools`.
    * `default_from_api`: `x-dcl-server-default` in the DCL OpenAPI spec.
    * `schema_config_mode_attr`: Requires a future override in `tpgtools`.
    * `at_least_one_of`, `exactly_one_of`: Will be exposed by
the DCL in the future. http://b/159243366
    * `update_mask_fields`: Handled by the DCL.
    * `key_expander`, `key_diff_suppress_func`: Unknown. Likely implemented in
`tpgtools` but no `Map` type fields have been implemented yet.
    * `flatten_object`: `COLLAPSED_OBJECT` in `tpgtools`.
    * `custom_expand`: Most are handled by the DCL. `"CUSTOM_STATE_GETTER"` in
`tpgtools`.
    * `custom_flatten`: Most are handled by the DCL. Not implemented, but will
be `"CUSTOM_STATE_SETTER"` in `tpgtools`.

## Custom Code (terraform.yaml)

Note: Custom code is entirely configured at the resource level.

* `extra_schema_entry`: `"VIRTUAL_FIELD"` in `tpgtools`, combined with other
overrides.
* `resource_definition`: `"CUSTOMIZE_DIFF"` in `tpgtools`, as well as a future
override to control schema version.
* `encoder`, `update_encoder`: Mostly handled in DCL. Some will be implemented
as `"CUSTOM_STATE_GETTER"`, likely in concert with `"VIRTUAL_FIELD"`.
* `decoder`: Mostly handled in DCL. Some will be implemented
as `"CUSTOM_STATE_SETTER"`, likely in concert with `"VIRTUAL_FIELD"`.
* `constants`: Written in a handwritten Go file.
* `pre_create`, `pre_update`, `post_update`, `post_delete`: Likely handled in
the DCL. Can be added to `tpgtools` if local behaviour is needed.
* `post_create`: Likely handled in the DCL. Local interactions, such as with
virtual fields can use `"POST_CREATE_FUNCTION"` in `tpgtools`.
* `post_create_error`: Likely handled in the DCL.
* `pre_delete`: Likely handled in the DCL. Local interactions, such as with
virtual fields can use `"PRE_DELETE_FUNCTION"` in `tpgtools`.
* `custom_create`, `custom_delete`: Handwrite the resource.
* `custom_import`: Requires a future `tpgtools` override.
* `post_import`: Could be added, but consider `custom_import`.
* `test_check_destroy`: Not relevant yet- no current support for examples.
