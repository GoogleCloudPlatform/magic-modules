---
name: add-data-source-workflow
description: "Workflow for checking if a corresponding Terraform data source exists and adding or updating a single data source (resource-backed, plural list, or resourceless) in Magic Modules."
---

# `add-data-source-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.

## Scope & Purpose

This workflow verifies whether a corresponding Terraform data source exists for a target GCP resource or read-only API endpoint, and adds or updates a single handwritten data source, acceptance tests, and registry documentation in `mmv1/third_party/terraform/`.

- **When Invoked**:
  - **Directly**: When asked (or routed via a `new-datasource` / `data-source` GitHub issue) to add, audit, or update a data source.
  - **Existing Data Source Field Audit**: When auditing an *existing* handwritten data source to ensure newly added resource fields propagate cleanly to data source schema/state.
  - Note: `new_resource` and `add_fields` do **not** automatically invoke this workflow.
- **PR Scoping Guardrail (1 Resource or Data Source per PR)**:
  - Keep every PR scoped to **one** resource or **one** data source (`docs/content/code-review/create-pr.md`). Never bundle a brand-new data source into the same PR as a new resource or field addition, and keep singular (`google_<product>_<resource>`) and plural (`google_<product>_<resources>`) data sources in separate PRs unless a maintainer explicitly asks to combine them. When both singular and plural data sources are needed, land the singular data source first so the plural data source can reuse its schema and flatteners.
- **Out of Scope**:
  - Direct edits in downstream `terraform-provider-google` or `terraform-provider-google-beta` repositories.
  - Modifying MMv1 generator templates (`mmv1/templates/terraform/`).
  - Opting MMv1 resources into Terraform `list` resource block generation (use `.agents/skills/workflows/add_list_resource/SKILL.md`).

---

## Task Definition

The task is complete only when all applicable deliverables below are satisfied:

1. **Existence Verified First**: Checked whether `google_<product>_<resource>` already has a generated or handwritten data source before creating any files:
   - Inspect `mmv1/products/<product>/*.yaml` for `datasource:` or `datasource_experimental:` with `generate: true` (for example `mmv1/products/networkservices/Gateway.yaml`).
   - Search `mmv1/third_party/terraform/services/<product>/data_source_*` and `mmv1/third_party/terraform/website/docs/d/`.
   - If the data source already exists, report its path to the user and confirm whether they want updates before editing (or, when auditing an existing manual-schema data source for missing fields, update its schema/`Read`/docs/tests as requested).
2. **Implementation Complete** (when adding or updating a data source):
   - Schema, `Read` function, and `registry.Schema` registration implemented in `mmv1/third_party/terraform/services/<product>/data_source_[google_]<tf_name_without_google_>.go` matching existing package file naming (or `.go.tmpl` with `package <product>` kept outside `{{- if ne $.TargetVersionName "ga" }} ... {{- end }}` when `min_version: beta`).
   - Documentation added or updated in `mmv1/third_party/terraform/website/docs/d/<tf_name_without_google_>.html.markdown` (never prefixed with `google_`, and copied verbatim without `.tmpl` rendering — include Beta callouts, `provider = google-beta`, and `google-beta` Registry links directly in `.html.markdown` when `min_version: beta`).
   - Acceptance tests added or updated in `mmv1/third_party/terraform/services/<product>/data_source_[google_]<tf_name_without_google_>_test.go` (or `_test.go.tmpl` with `acctest.ProtoV5ProviderBetaFactories(t)` and `provider = google-beta` when `min_version: beta`) verifying attribute parity and provider-default fallback for optional scope fields.
3. **Verified End-to-End**: Pre-generation static checks (`run-pre-gen-checks`), downstream compilation (`generate-provider` — including mandatory `VERSION=ga` verification for GA data sources), provider validation (`validate-provider-changes`), and acceptance tests (`run-acctests`, following user-confirmation guardrails before running billable `TF_ACC=1` tests) all pass without weakening baseline test coverage.

---

## Context Requirements & Knowledge Index

Consult `.agents/knowledge/index.md` first to locate canonical contributor guides and best practices rather than duplicating or guessing provider rules.

- **Strictly Required Context (Always Read)**:
  - `.agents/knowledge/index.md` — index of contributor guides, field rules, testing standards, and documentation style guides.
  - `docs/content/develop/add-handwritten-datasource.md` — baseline procedure for resource-backed and resourceless data sources, testing, and documentation.
  - Target resource schema (`mmv1/products/<product>/<Resource>.yaml`, `mmv1/products/<product>/product.yaml`, or `mmv1/third_party/terraform/services/<product>/resource_<product>_<resource>.go[.tmpl]`), or the GCP REST Discovery / OpenAPI specification when no backing resource exists.
  - Existing `data_source_*` files in `mmv1/third_party/terraform/services/<product>/` and `mmv1/third_party/terraform/website/docs/d/` to check existence and match package naming conventions.
- **Conditional Context (Read Only When Needed)**:
  - Downstream generated `google/services/<product>/resource_<product>_<resource>.go` (or `google-beta/...`) — open only when inspecting generated `ReplaceVars` templates, flatteners, `Sensitive` flags, `virtual_fields` defaults, `d.Identity()` usage, or `post_read` hooks.
  - `docs/content/document/handwritten-docs-style-guide.md` and `docs/content/test/test.md` — open via `.agents/knowledge/index.md` when writing custom attribute documentation, Beta warnings, or acceptance tests.
- **Unnecessary Context (Do Not Load)**:
  - Unrelated `<product>` service packages, bulk-reading all files in `.agents/knowledge/` or `docs/content/`, MMv1 compiler internals (`mmv1/provider/`, `mmv1/api/`), or `.agents/skills/workflows/add_data_source/references/evaluations.md` (only needed when running skill evaluations).

---

## Model Behavior: Non-Obvious Provider Invariants

Follow `docs/content/develop/add-handwritten-datasource.md` for standard schema wrapping, `init()` registration, test scaffolding, and documentation structure, matching the naming style of surrounding files in `mmv1/third_party/terraform/services/<product>/`. Apply the following constraints that are not enforced automatically by `DatasourceSchemaFromResourceSchema` or `resource<Name>Read`:

1. **Provider-Defaultable Scope Fields (`project`, `location`, `region`, `zone`) Should Be Optional**:
   - On new data sources, configure defaultable scope fields (`project`, `location`, `region`, `zone`) as `Optional` via `tpgresource.AddOptionalFieldsToSchema` rather than `Required` / `AddRequiredFieldsToSchema` (do not copy legacy data sources that marked `zone` or `location` required unless no provider/service default applies). Mark only the resource's primary identifier(s) and non-defaultable parent IDs as `Required`. Note that `AddOptionalFieldsToSchema` and `AddRequiredFieldsToSchema` panic on nil pointers at `init()` if a field name is absent from `Resource<Name>().Schema`.
   - When omitted in HCL, fall back to provider defaults via `tpgresource.GetProject`, `tpgresource.GetLocation`, `tpgresource.GetRegion`, or `tpgresource.GetZone` (or the service default such as `"global"`).
   - **Resolve and `d.Set` `location` and `url_param_only` inputs BEFORE `ReplaceVars` and `resource<Name>Read`**:
     - `BuildReplacementFunc` in `mmv1/third_party/terraform/tpgresource/utils.go` auto-resolves `{{project}}`, `{{region}}`, and `{{zone}}` from provider config, **not** `{{location}}`. Calling `ReplaceVars` with `{{location}}` before setting `location` produces an empty segment (`locations//`) when `location` is omitted in HCL.
     - Generated `resource<Name>Read` functions do not flatten `url_param_only: true` fields from API responses.
     - Always resolve `project`/`location` and call `if err := d.Set("location", location); err != nil { return ... }` (plus any required `url_param_only` or composite parent fields) **before** constructing the ID / calling `ReplaceVars` and `resource<Name>Read`.
     - When adding custom lookup fields via `tpgresource.MergeSchemas(dsSchema, customFields)`, pass `customFields` second so overlapping keys are not overwritten back to `Computed`-only.
     - **Code References**:
       - Standard optional `location` & `project` (sets `location` before `ReplaceVars` with error check, checks `d.Id() == ""`, sets labels, resets `deletion_protection` to `nil`): `mmv1/third_party/terraform/services/oracledatabase/data_source_oracle_database_odb_dns_zone.go`
       - Merging custom lookup fields (`tpgresource.MergeSchemas`): `mmv1/third_party/terraform/services/monitoring/data_source_monitoring_service.go`

2. **Re-apply `Sensitive = true` Stripped by `DatasourceSchemaFromResourceSchema`**:
   - `tpgresource.DatasourceSchemaFromResourceSchema` (`mmv1/third_party/terraform/tpgresource/datasource_helpers.go`) copies only `Description`, `Type`, `Set`, and `Elem`, **stripping `Sensitive: true`** at both top-level and nested `Elem.(*schema.Resource).Schema` levels.
   - Audit the backing resource for `sensitive: true` properties (or `effective_*` counterparts of sensitive fields) and explicitly restore `dsSchema["<field>"].Sensitive = true` (or on the nested `Elem.(*schema.Resource).Schema` entry) so secrets are never exposed in plaintext plan/output logs.

3. **Always Populate Data Source Labels and Annotations & Guard Shared `Read` Paths**:
   - Managed resource `Read` functions populate `effective_labels` / `effective_annotations` and leave `labels` / `annotations` empty on a fresh data source state. Whenever the resource has `labels` (`KeyValueLabels`), call `tpgresource.SetDataSourceLabels(d)` after `resource<Name>Read(d, meta)` (and `tpgresource.SetDataSourceAnnotations(d)` when `annotations` are present). Note that `acctest.CheckDataSourceStateMatchesResourceState` automatically skips `labels.*`, `terraform_labels.*`, `effective_labels.*`, and `deletion_policy`, but does **not** skip `annotations.*` / `effective_annotations.*`.
   - **Never call `d.Identity()` unconditionally in shared `Read` code**: Data sources do not declare `Identity: &schema.ResourceIdentity` unless explicitly configured, so unguarded `d.Identity()` calls inside handwritten/custom `Read` helpers invoked by a data source will error.
   - **Code Reference**: `dataSourceOracleDatabaseOdbDnsZoneRead` in `mmv1/third_party/terraform/services/oracledatabase/data_source_oracle_database_odb_dns_zone.go`.

4. **Handle `virtual_fields`, Write-Only Fields, and `ignore_read` Properties**:
   - **Client-side `virtual_fields` with `default_value` (Clear AND Ignore)**: In `mmv1/templates/terraform/resource.go.tmpl`, `resource<Name>Read` populates unset `virtual_fields` with their `default_value` (such as `deletion_protection = true`). **Never** hardcode a fake value like `d.Set("deletion_protection", false)` in production `Read`, and do **not** treat clearing vs. ignoring as alternatives:
     1. Always clear strictly client-side `virtual_fields` after `resource<Name>Read` via `d.Set("<field>", nil)` (`docs/content/develop/client-side-fields.md`) so the data source never leaks a fabricated default into user state.
     2. When the acceptance test resource sets a non-empty value (such as `deletion_protection = false` for cleanup), **also** pass `"<field>"` to `acctest.CheckDataSourceStateMatchesResourceStateWithIgnores` with an adjacent comment explaining that the field is client-side only and not returned by the API (note `deletion_policy` from Universal Deletion Policy is already skipped automatically by `test_utils.go.tmpl`).
     3. For `output: true` `virtual_fields` populated via `post_read` (such as `mmv1/products/compute/Reservation.yaml` `block_names` in `mmv1/templates/terraform/post_read/compute_reservetion.go.tmpl`), ensure all inputs read by the hook are set before calling `resource<Name>Read`.
     - **Code Reference**: `mmv1/third_party/terraform/services/oracledatabase/data_source_oracle_database_odb_dns_zone.go` and `mmv1/third_party/terraform/services/oracledatabase/data_source_oracle_database_odb_dns_zone_test.go`.
   - **Write-only (`_wo`) / `ignore_read` inputs without `effective_*` counterparts**: If a resource schema contains a write-only or `ignore_read: true` input without an `effective_*` counterpart that `resource<Name>Read` does not `d.Set`, remove it from the data source schema via `tpgresource.DeleteFieldsFromSchema(dsSchema, "<field>")` (only do this when `Read` does not `d.Set` the key, otherwise `d.Set` will fail with `Invalid address to set`).
     - **Code Reference**: `mmv1/third_party/terraform/services/certificatemanager/data_source_google_certificate_manager_certificates.go`.
   - **`ignore_read: true` fields with `effective_*` counterparts**: Generated flatteners skip `ignore_read: true` properties; copy the value from its `effective_*` field after `resource<Name>Read`.
     - **Code Reference**: `SetDataSourceReservedIpRange` in `mmv1/third_party/terraform/services/redis/data_source_redis_instance.go`.

5. **Fail Explicitly When the Target Resource Is Not Found (`d.Id() == ""`)**:
   - On HTTP 404, `transport_tpg.HandleNotFoundError` inside `resource<Name>Read` clears `d.SetId("")` and returns `nil`. Capture `id` in a local variable, call `d.SetId(id)` before `resource<Name>Read`, and check `if d.Id() == "" { return fmt.Errorf("%s not found", id) }` afterwards (never interpolate `d.Id()` in the error message after it was cleared to `""`).
   - **Code Reference**: `dataSourceOracleDatabaseOdbDnsZoneRead` in `mmv1/third_party/terraform/services/oracledatabase/data_source_oracle_database_odb_dns_zone.go`.

6. **Resourceless & Plural (`list`) Data Sources**:
   - **Endpoint & Transport Hygiene**: Build URLs from `transport_tpg.BaseUrl(Product, config)` (or `tpgresource.ReplaceVars(d, config, "{{<Product>BasePath}}...")` when a package lacks a generated `Product` variable) so custom endpoints and regional REP overrides work seamlessly, set `UseJSONNumber: true` on `&schema.Resource{}` when responses contain 64-bit integers or numbers, and use `transport_tpg.SendRequest` paired with `transport_tpg.HandleDataSourceNotFoundError` on HTTP 404. Only wrap side-effect-free synchronous read RPCs as data sources; if an API exposes read-only metadata exclusively via an asynchronous `POST` operation, poll the operation to completion and document the refresh behavior prominently with a `~> **Note:**` callout.
   - **Resourceless Data Sources**: Resolve optional `project` and `zone`/`region`/`location` defaults (`d.Set` them back into state), remap reserved Terraform schema names (`id` -> `resource_id`, `count` -> `<type>_count`), set a deterministic `d.SetId(...)`, assert attributes in acceptance tests (`TestCheckResourceAttr` / `TestCheckResourceAttrSet`), and manually verify every exported attribute is documented in `website/docs/d/` (`detect-missing-docs` only checks file existence for data sources).
     - **Code References**: `mmv1/third_party/terraform/services/storagecontrol/data_source_storage_control_project_intelligence_finding.go` (`data_source_storage_control_project_intelligence_finding_test.go`) and `mmv1/third_party/terraform/services/cloudquotas/data_source_google_cloud_quotas_quota_info.go` (`data_source_google_cloud_quotas_quota_info_test.go`).
   - **Plural (`list`) Data Sources**: Implement the singular data source first when applicable so the plural data source can reuse the singular schema (`DatasourceSchemaFromResourceSchema` or shared item schema) and item flattener. Paginate all pages (using `transport_tpg.ListPages` when backed by a `*schema.Resource` + flattener, or a `nextPageToken` loop over `transport_tpg.SendRequest` with `transport_tpg.Is429RetryableQuotaError`), include `filter` in `d.SetId(...)` when `filter` is set so filtered and unfiltered data source blocks do not collide, and assert `>= 1` returned item in acceptance tests.
     - **Code References**: `mmv1/third_party/terraform/services/networkmanagement/data_source_network_management_connectivity_tests.go` and `mmv1/third_party/terraform/services/compute/data_source_google_compute_zones.go` (`data_source_google_compute_zones_test.go`).

7. **Beta-Only Data Sources (`min_version: beta`)**:
   - In `.go.tmpl` and `_test.go.tmpl` files, keep `package <product>` (or `package <product>_test`) **outside** the `{{- if ne $.TargetVersionName "ga" }} ... {{- end }}` guard (`docs/content/develop/add-resource.md`, validated by `tools/template-check`), and escape `ReplaceVars` placeholders inside `.go.tmpl` as `{{"{{"}}project{{"}}"}}`.
   - Use `acctest.ProtoV5ProviderBetaFactories(t)` in `_test.go.tmpl` and add `provider = google-beta` to every resource and data source block in HCL examples and tests.
   - Because `mmv1/third_party/terraform/website/docs/d/` is copied verbatim without `.tmpl` rendering, write the standard Beta warning banner (`~> **Warning:** This data source is in beta, and should be used with the terraform-provider-google-beta provider...` per `docs/content/document/handwritten-docs-style-guide.md`) and `google-beta` Registry links directly in `.html.markdown`.

---

## Baseline-First Workflow Steps

Start with the minimal baseline implementation rather than over-engineering for speculative edge cases; let compilation and `CheckDataSourceStateMatchesResourceState` runs surface any resource-specific discrepancies.

### Step 1: Check Existence & Confirm Single-Target PR Scope
- Check `mmv1/products/<product>/*.yaml` for `datasource:` / `datasource_experimental:` (`generate: true`), and search `mmv1/third_party/terraform/services/<product>/` (`data_source_*`) and `mmv1/third_party/terraform/website/docs/d/` for an existing data source.
- Follow **Task Definition (#1)** and the **1 Resource or Data Source per PR** guardrail before creating or modifying files.

### Step 2: Inspect the Target Resource or API & Implement Baseline
- Inspect the target YAML/Go resource definition (or REST API schema and `product.yaml` if resourceless) for URL parameters, `sensitive: true` fields, `labels`/`annotations`, `virtual_fields`, `ignore_read` properties, and `min_version: beta`.
- Implement the baseline data source, acceptance test (covering both explicit scope arguments and omitted provider-default fallback), and documentation per `docs/content/develop/add-handwritten-datasource.md` and **Model Behavior** above.

### Step 3: Generate, Build, and Test Empirically
1. Run static checks via `.agents/skills/utils/run-pre-gen-checks/SKILL.md`.
2. During inner-loop iteration, generate and test the single relevant downstream (`beta` for `min_version: beta`, or `ga` for GA data sources) via `.agents/skills/operations/repo-sync/SKILL.md` and `.agents/skills/operations/generate-provider/SKILL.md` rather than regenerating both on every edit.
3. Before completion, **always** generate and verify `VERSION=ga` (`terraform-provider-google`) in addition to `VERSION=beta` for GA data sources (upstream `VCR-test` runs against `google-beta` and will not catch unguarded beta references or `/v1/` endpoint differences in GA), and run `.agents/skills/utils/validate-provider-changes/SKILL.md` (note that `diff-processor` checks breaking changes on data sources, skips data sources in `detect-missing-tests`, and checks only file existence in `detect-missing-docs`, so review resourceless `website/docs/d/` attribute coverage manually).
4. Run acceptance tests via `.agents/skills/utils/run-acctests/SKILL.md` (confirming with the user first before running live billable `TF_ACC=1` tests against GCP, testing both `ga` and `beta` for GA data sources or `beta` for Beta-only data sources). If `CheckDataSourceStateMatchesResourceState` or API requests fail on a resource-specific edge case, inspect the diff/log and apply a targeted fix before re-running.

---

## Evaluation Summary

Task accuracy and skill effectiveness are evaluated against a baseline-vs-skill benchmark protocol covering resource-backed, `virtual_fields` / `Sensitive` / `ignore_read`, resourceless / plural `SendRequest`, and existing-data-source audit tasks:

- **Pass Criteria**:
  1. Existence check (both MMv1 YAML `datasource` / `datasource_experimental` and handwritten `data_source_*`) performed before creating files; zero duplicate data sources; 1 data source per PR.
  2. Defaultable scope fields (`project`, `location`, `region`, `zone`) are `Optional`, set before `ReplaceVars` with error checks when required by URL templates, and work when omitted in HCL.
  3. `Sensitive = true` restored on sensitive attributes, `SetDataSourceLabels` / `SetDataSourceAnnotations` called when applicable, and `if d.Id() == ""` (or `HandleDataSourceNotFoundError`) 404 checks return a clear error.
  4. Client-side `virtual_fields` are **both** cleared via `d.Set("<field>", nil)` in `Read` **and** ignored via `CheckDataSourceStateMatchesResourceStateWithIgnores` with an adjacent justification comment; plural data sources assert `>= 1` item.
  5. Downstream build (`VERSION=ga` and `VERSION=beta` for GA; `VERSION=beta` with clean `VERSION=ga` guard for Beta), pre-gen checks, provider validation, and acceptance tests pass.
- **Full Benchmark & A/B Evaluation Framework**: See `.agents/skills/workflows/add_data_source/references/evaluations.md`.
