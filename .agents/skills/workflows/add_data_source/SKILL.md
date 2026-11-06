---
name: add-data-source-workflow
description: "Workflow for checking if a corresponding Terraform data source exists and adding or updating a data source (resource-backed or resourceless) in Magic Modules. Invoked directly or from add-fields-workflow and new-resource-workflow."
---

# `add-data-source-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.

## Scope & Purpose

This workflow verifies whether a corresponding Terraform data source exists for a target GCP resource or read-only API endpoint, and adds or updates the handwritten data source, acceptance tests, and registry documentation in `mmv1/third_party/terraform/`.

- **When Invoked**:
  - **Directly**: When asked to add, audit, or update a data source.
  - **From `new-resource-workflow`**: After adding a new resource, to ensure a corresponding data source exists.
  - **From `add-fields-workflow`**: After adding fields to an existing resource, to ensure the corresponding data source exists and accurately exposes the new fields.
- **Out of Scope**:
  - Direct edits in downstream `terraform-provider-google` or `terraform-provider-google-beta` repositories.
  - Modifying MMv1 generator templates (`mmv1/templates/terraform/`).
  - Opting MMv1 resources into Terraform `list` block generation (use `.agents/skills/workflows/add_list_resource/SKILL.md`).

---

## Task Definition

The task is complete only when all applicable deliverables below are satisfied:

1. **Existence Verified First**: Checked whether `google_<product>_<resource>` already has a data source before creating any files.
   - If invoked directly and the data source already exists, report its path to the user and confirm whether they want updates before editing.
   - If invoked from `add-fields-workflow` and the data source already exists, verify that new fields propagate to the data source schema and `Read` state, updating `Read` / manual schemas only if required, and run the data source acceptance test.
2. **Implementation Complete** (when adding or updating a data source):
   - Schema, `Read` function, and `registry.Schema` registration implemented in `mmv1/third_party/terraform/services/<product>/data_source_<product>_<resource>.go` (or `.go.tmpl` with a `{{ if ne $.TargetVersionName "ga" -}}` guard if `min_version: beta`).
   - Documentation added or updated in `mmv1/third_party/terraform/website/docs/d/<product>_<resource>.html.markdown`.
   - Acceptance tests added or updated in `mmv1/third_party/terraform/services/<product>/data_source_<product>_<resource>_test.go` (or `.go.tmpl`) verifying attribute parity and provider-default fallback for optional scope fields.
3. **Verified End-to-End**: Pre-generation static checks (`run-pre-gen-checks`), downstream provider generation and compilation (`generate-provider`), provider validation (`validate-provider-changes`), and live acceptance tests (`run-acctests`) all pass without weakening baseline test coverage.

---

## Context Requirements & Knowledge Index

Consult `.agents/knowledge/index.md` first to locate canonical contributor guides and best practices rather than duplicating or guessing provider rules.

- **Strictly Required Context (Always Read)**:
  - `.agents/knowledge/index.md` — index of contributor guides, field rules, testing standards, and documentation style guides.
  - `docs/content/develop/add-handwritten-datasource.md` — baseline procedure for resource-backed and resourceless data sources, testing, and documentation.
  - Target resource schema (`mmv1/products/<product>/<Resource>.yaml` or `mmv1/third_party/terraform/services/<product>/resource_<product>_<resource>.go`), or the GCP REST Discovery / OpenAPI specification when no backing resource exists.
  - Existing `data_source_*` files in `mmv1/third_party/terraform/services/<product>/` and `mmv1/third_party/terraform/website/docs/d/` to check existence and match package naming conventions.
- **Conditional Context (Read Only When Needed)**:
  - Downstream generated `google/services/<product>/resource_<product>_<resource>.go` — open only when inspecting generated `ReplaceVars` templates, flatteners, `virtual_fields` defaults, or `post_read` hooks.
  - `docs/content/document/handwritten-docs-style-guide.md` and `docs/content/test/test.md` — open via `.agents/knowledge/index.md` when writing custom attribute tables or acceptance tests.
- **Unnecessary Context (Do Not Load)**:
  - Unrelated `<product>` service packages, bulk-reading all files in `.agents/knowledge/` or `docs/content/`, MMv1 compiler internals (`mmv1/provider/`, `mmv1/api/`), or `.agents/skills/workflows/add_data_source/references/evaluations.md` (only needed when running skill evaluations).

---

## Model Behavior: Non-Obvious Provider Invariants

Follow `docs/content/develop/add-handwritten-datasource.md` for standard schema wrapping, `init()` registration, test scaffolding, and documentation structure, matching the naming style of surrounding files in `mmv1/third_party/terraform/services/<product>/`. Apply the following constraints that are not enforced automatically by `DatasourceSchemaFromResourceSchema` or `resource<Name>Read`:

1. **`project`, `location`, `region`, and `zone` MUST Always Be Optional**:
   - Never mark `project`, `location`, `region`, or `zone` as `Required` on a data source (do not copy older data sources that mistakenly marked `zone` required). Mark only the resource's primary identifier(s) and non-defaultable parent IDs as `Required`.
   - When omitted in HCL, fall back to provider defaults via `tpgresource.GetProject`, `tpgresource.GetLocation`, `tpgresource.GetRegion`, or `tpgresource.GetZone` (or the service default such as `"global"`).
   - **Explicitly `d.Set` `location` and `url_param_only` inputs before `ReplaceVars` / `resource<Name>Read`**:
     - `BuildReplacementFunc` in `mmv1/third_party/terraform/tpgresource/utils.go` auto-resolves `{{project}}`, `{{region}}`, and `{{zone}}` from provider config, **not** `{{location}}`.
     - Generated `resource<Name>Read` functions do not flatten `url_param_only: true` fields from API responses.
     - Resolve and `d.Set(...)` `location` plus any required `url_param_only` or composite parent fields prior to calling `ReplaceVars` and `resource<Name>Read`.
     - **Code References**:
       - Standard optional `location` & `project`: `mmv1/third_party/terraform/services/memorystore/data_source_memorystore_instance.go`
       - Merging friendlier parent/scope inputs (`tpgresource.MergeSchemas`): `mmv1/third_party/terraform/services/alloydb/data_source_alloydb_database_instance.go` and `mmv1/third_party/terraform/services/lustre/data_source_lustre_instance.go`

2. **Always Populate Data Source Labels and Annotations**:
   - Managed resource `Read` functions populate `effective_labels` / `effective_annotations` and leave `labels` / `annotations` empty on a fresh data source state.
   - Whenever the resource has `labels` (`KeyValueLabels`), call `tpgresource.SetDataSourceLabels(d)` after `resource<Name>Read(d, meta)` (and `tpgresource.SetDataSourceAnnotations(d)` when `annotations` are present).
   - **Code Reference**: `dataSourceMemorystoreInstanceRead` in `mmv1/third_party/terraform/services/memorystore/data_source_memorystore_instance.go`.

3. **Handle `virtual_fields` and `ignore_read` Properties**:
   - **`virtual_fields` with `default_value`**: In `mmv1/templates/terraform/resource.go.tmpl`, `resource<Name>Read` sets unset virtual fields to their `default_value` (such as `deletion_protection = true`), causing mismatches when test resources set `deletion_protection = false` for cleanup. Reset or clear client-side virtual fields after `resource<Name>Read`, or ignore strictly client-side fields in `acctest.CheckDataSourceStateMatchesResourceStateWithIgnores` with an adjacent comment justifying why the API does not return them. For `output: true` virtual fields populated via `post_read`, ensure inputs read by the hook are set before calling `resource<Name>Read`.
     - **Code References**: `mmv1/third_party/terraform/services/alloydb/data_source_alloydb_cluster.go`, `mmv1/third_party/terraform/services/compute/data_source_google_compute_storage_pool.go`, and `mmv1/third_party/terraform/services/sql/data_source_sql_database.go`.
   - **`ignore_read: true` fields with `effective_*` counterparts**: Generated flatteners skip `ignore_read: true` properties; copy the value from its `effective_*` field after `resource<Name>Read`.
     - **Code Reference**: `SetDataSourceReservedIpRange` in `mmv1/third_party/terraform/services/redis/data_source_redis_instance.go`.

4. **Fail Explicitly When the Target Resource Is Not Found (`d.Id() == ""`)**:
   - On HTTP 404, `transport_tpg.HandleNotFoundError` inside `resource<Name>Read` clears `d.SetId("")` and returns `nil`. Always set `d.SetId(id)` before calling `resource<Name>Read` and return a `not found` error if `d.Id() == ""` afterwards.
   - **Code Reference**: `dataSourceMemorystoreInstanceRead` in `mmv1/third_party/terraform/services/memorystore/data_source_memorystore_instance.go`.

5. **Resourceless Data Sources (No Backing Terraform Resource)**:
   - When querying a read-only GCP object with no corresponding Terraform resource, define the schema directly and implement `Read` using `transport_tpg.SendRequest`. Resolve provider defaults for `project` and `zone`/`region`/`location`, remap reserved Terraform attribute names (`id` -> `resource_id`, `count` -> `<type>_count`), set `d.SetId(...)`, assert attributes directly in acceptance tests (`TestCheckResourceAttrSet` / `TestCheckResourceAttr`), and document all exported attributes in `website/docs/d/`.
   - **Code References**: `mmv1/third_party/terraform/services/compute/data_source_google_compute_reservation_block.go` and `mmv1/third_party/terraform/website/docs/d/compute_reservation_block.html.markdown`.

---

## Baseline-First Workflow Steps

Start with the minimal baseline implementation rather than over-engineering for speculative edge cases; let compilation and live `CheckDataSourceStateMatchesResourceState` runs surface any resource-specific discrepancies.

### Step 1: Check Existence
- Search `mmv1/third_party/terraform/services/<product>/` (`data_source_*`) and `mmv1/third_party/terraform/website/docs/d/` for an existing data source.
- Follow the decision criteria in **Task Definition (#1)** above before creating or modifying files.

### Step 2: Inspect the Target Resource or API & Implement Baseline
- Inspect the target YAML/Go resource definition (or REST API schema if resourceless) for URL parameters, `labels`/`annotations`, `virtual_fields`, `ignore_read` properties, and `min_version: beta`.
- Implement the baseline data source, acceptance test (covering both explicit scope arguments and omitted provider-default fallback), and documentation per `docs/content/develop/add-handwritten-datasource.md` and **Model Behavior** above.

### Step 3: Generate, Build, and Test Empirically
1. Run static checks via `.agents/skills/utils/run-pre-gen-checks/SKILL.md`.
2. Sync, generate, and compile the downstream provider(s) (`ga` and/or `beta`) via `.agents/skills/operations/repo-sync/SKILL.md` and `.agents/skills/operations/generate-provider/SKILL.md`.
3. Validate schema/documentation changes via `.agents/skills/utils/validate-provider-changes/SKILL.md`.
4. Run acceptance tests via `.agents/skills/utils/run-acctests/SKILL.md`. If `CheckDataSourceStateMatchesResourceState` or API requests fail on a resource-specific edge case, inspect the diff/log and apply a targeted fix before re-running.

---

## Evaluation Summary

Task accuracy and skill effectiveness are evaluated against a baseline-vs-skill benchmark protocol covering all four data source task classes (standard resource-backed, `virtual_fields` / `ignore_read` handling, resourceless `SendRequest`, and `add-fields-workflow` audits):

- **Pass Criteria**:
  1. Existence check performed before creating files; zero duplicate data sources.
  2. `project`, `location`, `region`, and `zone` are all `Optional` and work when omitted in HCL.
  3. `SetDataSourceLabels` / `SetDataSourceAnnotations` and `if d.Id() == ""` 404 checks are present where applicable.
  4. `CheckDataSourceStateMatchesResourceState` passes without weakening coverage (any ignored client-side field has an adjacent justification comment).
  5. Downstream build, pre-gen checks, provider validation, and live `TF_ACC=1` tests pass.
- **Full Benchmark & A/B Evaluation Framework**: See `.agents/skills/workflows/add_data_source/references/evaluations.md`.
