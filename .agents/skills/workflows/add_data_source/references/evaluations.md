# Evaluation Framework: `add-data-source-workflow`

This document defines how to measure model accuracy, task completion rate, and instruction efficiency for `.agents/skills/workflows/add_data_source/SKILL.md`.

## 1. Baseline-First A/B Protocol

Before adding new rules or edge cases to `SKILL.md`, run candidate prompts in two isolated git worktrees pinned to the same pre-PR commit:

1. **Control Run (Baseline — Without Skill)**:
   - Execute the task prompt using only the default repository context (`AGENTS.md` + `.agents/knowledge/index.md` + `docs/content/develop/add-handwritten-datasource.md`).
   - Record what the model already handles naturally (e.g., standard `schema.Resource` boilerplate, `init()` registration, basic `SendRequest` JSON unmarshaling) versus where it fails (e.g., marking `location`/`zone` as `Required`, omitting `d.Set("location", location)` before `ReplaceVars`, missing `SetDataSourceLabels`, failing `CheckDataSourceStateMatchesResourceState` on `deletion_protection`).
2. **Treatment Run (With `add-data-source-workflow`)**:
   - Execute the identical prompt with `add_data_source/SKILL.md` enabled.
   - Verify that every instruction in `SKILL.md` improves pass rate over baseline without over-constraining code style or inflating context token usage on unnecessary files.

> **Pruning Rule**: If the Control Run passes an edge case reliably on its own—or fails only once before immediately self-correcting from compiler/test feedback—do **not** bloat `SKILL.md` with extra instructions for that edge case.

---

## 2. Representative Benchmark Suite

Evaluate across these four historical task archetypes (using `git checkout <commit>~1` in an ephemeral branch to test clean recreation):

| ID | Archetype | Historical Target | Key Pitfalls Tested |
| :--- | :--- | :--- | :--- |
| **EVAL-1** | Standard Resource-Backed Data Source | `google_memorystore_instance` (`mmv1/third_party/terraform/services/memorystore/data_source_memorystore_instance.go`) | Optional `project` + `location`; setting `location` (`url_param_only`) in state; `SetDataSourceLabels`; `d.Id() == ""` 404 check; provider-default fallback test. |
| **EVAL-2** | Resource with `virtual_fields` Default | `google_alloydb_cluster` (`mmv1/third_party/terraform/services/alloydb/data_source_alloydb_cluster.go`) or `google_compute_storage_pool` (`mmv1/third_party/terraform/services/compute/data_source_google_compute_storage_pool.go`) | Handling client-side `deletion_protection` (`default_value: true` in YAML) so `CheckDataSourceStateMatchesResourceState*` passes without un-commented ignores. |
| **EVAL-3** | Resourceless Data Source (Custom `Read`) | `google_compute_reservation_block` (`mmv1/third_party/terraform/services/compute/data_source_google_compute_reservation_block.go`) | Custom `map[string]*schema.Schema`; optional `zone` + `project` default resolution; `transport_tpg.SendRequest`; avoiding reserved `id`/`count` names; full attribute docs. |
| **EVAL-4** | Existing Data Source Audit via `add-fields-workflow` | `google_redis_instance` (`mmv1/third_party/terraform/services/redis/data_source_redis_instance.go`) | Detecting existing data source first (zero duplicate files); handling `ignore_read: true` field backed by `effective_*` (`SetDataSourceReservedIpRange`); running existing data source acceptance tests. |

---

## 3. Objective Grading Rubric

Score each evaluation run across six binary checks (all 6 must pass for task completion):

1. **Existence Check Accuracy**:
   - Searched `mmv1/third_party/terraform/services/<product>/` and `website/docs/d/` prior to creating files; created zero duplicate data sources when one already existed.
2. **Scope Optionality & Default Resolution**:
   - Zero occurrences of `Required: true` or `AddRequiredFieldsToSchema` on `project`, `location`, `region`, or `zone`.
   - Omitting `project` and `location`/`region`/`zone` in HCL resolves cleanly from provider config without empty URL segments (`//`) or empty state attributes.
3. **State Completeness & 404 Handling**:
   - `tpgresource.SetDataSourceLabels(d)` / `SetDataSourceAnnotations(d)` called whenever the backing resource has `KeyValueLabels` / `KeyValueAnnotations`.
   - Returns a non-nil error when the target resource returns HTTP 404 (`d.Id() == ""` or `res == nil`).
4. **Test Coverage Integrity**:
   - Resource-backed tests assert state parity via `CheckDataSourceStateMatchesResourceState` (any client-side virtual field excluded via `CheckDataSourceStateMatchesResourceStateWithIgnores` has an adjacent comment justifying the API behavior).
   - Resourceless tests assert `id` and key computed fields via `TestCheckResourceAttr*`.
5. **Context Discipline**:
   - Opened `.agents/knowledge/index.md`, `docs/content/develop/add-handwritten-datasource.md`, and target service files without bulk-reading unrelated product directories or MMv1 compiler internals.
6. **Build, Validation, and Live `TF_ACC=1` Gate**:
   - `run-pre-gen-checks`, downstream `make build` (`ga` and `beta`), `validate-provider-changes`, and `TF_ACC=1` acceptance tests all exit `0`.
