# Evaluation Framework: `add-data-source-workflow`

This document defines how to measure model accuracy, task completion rate, and instruction efficiency for `.agents/skills/workflows/add_data_source/SKILL.md`.

## 1. Baseline-First A/B Protocol

Before adding new rules or edge cases to `SKILL.md`, run candidate prompts in two isolated git worktrees pinned to the same pre-PR commit:

1. **Control Run (Baseline — Without Skill)**:
   - Execute the task prompt using only the default repository context (`AGENTS.md` + `.agents/knowledge/index.md` + `docs/content/develop/add-handwritten-datasource.md`).
   - Record what the model already handles naturally (e.g., standard `schema.Resource` boilerplate, `init()` registration, basic `SendRequest` JSON unmarshaling) versus where it fails (e.g., skipping MMv1 YAML `datasource_experimental` checks, marking `location`/`zone` as `Required`, calling `ReplaceVars` with `{{location}}` before `d.Set("location", location)`, losing `Sensitive = true` after `DatasourceSchemaFromResourceSchema`, missing `SetDataSourceLabels`, or only resetting *or* ignoring client-side `virtual_fields` instead of doing both).
2. **Treatment Run (With `add-data-source-workflow`)**:
   - Execute the identical prompt with `add_data_source/SKILL.md` enabled.
   - Verify that every instruction in `SKILL.md` improves pass rate over baseline without over-constraining code style or inflating context token usage on unnecessary files.

> **Pruning Rule**: If the Control Run passes an edge case reliably on its own—or fails only once before immediately self-correcting from compiler/test feedback—do **not** bloat `SKILL.md` with extra instructions for that edge case.

---

## 2. Representative Benchmark Suite

Evaluate across these historical task archetypes (using `git checkout <commit>~1` in an ephemeral branch to test clean recreation):

| ID | Archetype | Historical Target | Key Pitfalls Tested |
| :--- | :--- | :--- | :--- |
| **EVAL-1** | Standard Resource-Backed Data Source & `virtual_fields` | `google_oracle_database_odb_dns_zone` (`mmv1/third_party/terraform/services/oracledatabase/data_source_oracle_database_odb_dns_zone.go` & `data_source_oracle_database_odb_dns_zone_test.go`) | Optional `project` + `location` (`AddOptionalFieldsToSchema`); resolving and setting `location` + `url_param_only` `odbnetwork` with error checks before `ReplaceVars` and `Read`; `SetDataSourceLabels`; `d.Id() == ""` 404 check with saved `id`; clearing client-side `deletion_protection` via `d.Set("deletion_protection", nil)` **and** excluding it via commented `CheckDataSourceStateMatchesResourceStateWithIgnores` (plus omitted `project`/`location` provider-default test coverage). |
| **EVAL-2** | Sensitive Fields & `post_read` Virtual Fields | `google_redis_cluster_auth_token` / `google_compute_reservation` (`mmv1/products/compute/Reservation.yaml` `block_names` via `post_read`) | Restoring `dsSchema["<field>"].Sensitive = true` stripped by `DatasourceSchemaFromResourceSchema`; setting inputs consumed by `post_read` hooks prior to `resource<Name>Read`. |
| **EVAL-3** | Resourceless & Plural Data Sources (`SendRequest`) | `google_storage_control_project_intelligence_finding` (`mmv1/third_party/terraform/services/storagecontrol/data_source_storage_control_project_intelligence_finding.go`), `google_compute_zones` (`mmv1/third_party/terraform/services/compute/data_source_google_compute_zones.go`), and `google_cloud_quotas_quota_info` (`mmv1/third_party/terraform/services/cloudquotas/data_source_google_cloud_quotas_quota_info.go`) | `transport_tpg.BaseUrl(Product, config)`; `HandleDataSourceNotFoundError`; `UseJSONNumber: true`; optional `project` + `location`/`region` default resolution and `d.Set`; pagination across `nextPageToken`; asserting `>= 1` returned item and scalar/nested attributes in `_test.go`; complete `website/docs/d/` documentation. |
| **EVAL-4** | Existing Data Source Audit for New Fields | `google_redis_instance` (`mmv1/third_party/terraform/services/redis/data_source_redis_instance.go`) | Detecting existing generated (`datasource` / `datasource_experimental`) or handwritten data source first (zero duplicate files); handling `ignore_read: true` field backed by `effective_*` (`SetDataSourceReservedIpRange`); running existing data source acceptance tests. |

---

## 3. Objective Grading Rubric

Score each evaluation run across six binary checks (all 6 must pass for task completion):

1. **Existence Check & PR Scope Accuracy**:
   - Checked MMv1 YAML (`datasource` / `datasource_experimental`) plus `mmv1/third_party/terraform/services/<product>/` and `website/docs/d/` prior to creating files; created zero duplicate data sources; kept PR scoped to a single data source.
2. **Scope Optionality & Default Resolution**:
   - Defaultable scope fields (`project`, `location`, `region`, `zone`) configured as `Optional` via `AddOptionalFieldsToSchema` (not `Required` / `AddRequiredFieldsToSchema`).
   - `location` and `url_param_only` inputs resolved and set (checking `d.Set` errors) before `ReplaceVars` / `resource<Name>Read`, so omitting `project` and `location`/`region`/`zone` in HCL resolves cleanly from provider config without empty URL segments (`//`) or empty state attributes.
3. **State Completeness, Sensitivity & 404 Handling**:
   - `dsSchema["<field>"].Sensitive = true` restored for every sensitive resource attribute.
   - `tpgresource.SetDataSourceLabels(d)` / `SetDataSourceAnnotations(d)` called whenever the backing resource has `KeyValueLabels` / `KeyValueAnnotations`.
   - Returns a non-nil error naming the target ID when HTTP 404 clears `d.Id() == ""` (or via `HandleDataSourceNotFoundError`).
4. **Test Coverage Integrity**:
   - Resource-backed tests assert state parity via `CheckDataSourceStateMatchesResourceState` (any client-side `virtual_fields` entry with `default_value` is **both** cleared via `d.Set(..., nil)` in `Read` **and** excluded via `CheckDataSourceStateMatchesResourceStateWithIgnores` with an adjacent comment justifying the API behavior; never hardcodes `deletion_protection = false` in production `Read`).
   - Resourceless and plural tests assert `id`, key computed fields, and `>= 1` list item via `TestCheckResourceAttr*`.
5. **Context & Inner-Loop Discipline**:
   - Opened `.agents/knowledge/index.md`, `docs/content/develop/add-handwritten-datasource.md`, and target service files without bulk-reading unrelated product directories or MMv1 compiler internals; iterated against a single downstream before running final `VERSION=ga` + `VERSION=beta` verification.
6. **Build, Validation, and `TF_ACC=1` Gate**:
   - `run-pre-gen-checks`, downstream compilation (`VERSION=ga` and `VERSION=beta` for GA; `VERSION=beta` plus guarded `VERSION=ga` build for Beta), `validate-provider-changes`, and `TF_ACC=1` acceptance tests all exit `0`.
