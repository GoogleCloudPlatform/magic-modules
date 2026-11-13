---
name: cross-schema-field-sync
description: Keep paired *Version.yaml schemas and TGC is_missing_in_cai flags synchronized when adding fields.
topics: [field, tgc]
task_types: [field-add, new-resource, code-review]
source: docs/content/develop/add-fields.md + docs/content/convert/add-new-resource-tgc.md
status: approved
last_verified: 2026-11-06
---

# Cross-Schema Field Sync (`*Version.yaml` & `is_missing_in_cai`)

When adding fields to `mmv1/products/<product>/<Resource>.yaml`, check sibling schemas and Cloud Asset Inventory (CAI) converter flags before generating.

## 1. Paired `<Resource>Version.yaml` Schemas

Some products pair a resource with a sibling `<Resource>Version.yaml` that embeds the parent resource's configuration message:

- **Rule:** Check the same product directory for `<Resource>Version.yaml`. If the version API's embedded message exposes the new field, mirror it using **that version schema's existing mutability pattern**:
  - Read-only snapshot blocks (e.g., `snapshot` in `mmv1/products/ces/AppVersion.yaml`): mark the mirrored property and all nested children `output: true`.
  - Configurable version input blocks (e.g., `tool` in `mmv1/products/dialogflowcx/ToolVersion.yaml`): match the parent's `required` / optional input attributes.
- Skip request-only parameters (`url_param_only: true`) or write-only secrets omitted from version payloads.

## 2. `is_missing_in_cai: true` on `include_in_tgc_next: true` Resources

- **Rule:** If `Resource.yaml` sets `include_in_tgc_next: true` and a newly added field does not yet exist in the upstream CAI asset proto for that resource, set `is_missing_in_cai: true` on the field so `terraform-google-conversion` skips it cleanly.
- **Do NOT** set `is_missing_in_cai: true` when `include_in_tgc_next: true` is absent or when the CAI schema already includes the field.
