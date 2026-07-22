---
name: test-failure-decision-tree
description: "Classification decision tree and remediation strategies for diagnosing Terraform acceptance test failures across all workflows."
---

# `test-failure-decision-tree`

This document provides a standardized classification catalog for diagnosing, isolating, and fixing failing Terraform acceptance tests in `magic-modules`. It is designed to be consulted by any workflow (`test_fix`, `new_resource`, `add_list_resource`, `default`) whenever `make testacc` fails.

---

## Decision Matrix Catalog

| Scenario ID | Symptom / Error Pattern | Root Cause Category | Primary File Location & Remedy |
| :--- | :--- | :--- | :--- |
| **Scenario 1** | State drift, `ImportStateVerify` mismatch, or non-empty plan after Apply | State Normalization / Default Mismatch | Edit `mmv1/products/<product>/<Resource>.yaml` to add `diff_suppress_func`, adjust state reader, or handle API defaults |
| **Scenario 2** | HTTP 400 `InvalidArgument`, `Unknown Field`, or field serialization error | API Request Payload Schema Mismatch | Compare POST request JSON against API proto schema; fix field `camelCase`/`snake_case` mapping or `send_empty_value` |
| **Scenario 3** | HTTP 404 immediately post-Create or concurrent operation conflict | LRO / Eventual Consistency Timing | Configure `autogen_async` in YAML or add polling waiter logic in handwritten code overrides |
| **Scenario 4** | Pre-requisite resource failure, name collision, or test setup error | Sample HCL Test Template Bug | Update `..._test.go.tmpl` to use dynamic random string suffixes (`ResourceIdVars`), or add target version guard (`{{- if ne $.TargetVersionName "ga" }}`) for beta-only resources |
| **Scenario 5** | gRPC Code 13/14 (`RESOURCE_EXHAUSTED`), HTTP 429 (`Quota limit exceeded`/`Quota exhausted`), HTTP 500/502/503 (`An internal error has occurred`), or API backend crash | Non-Remediable GCP Backend / Quota Environment Error | **Early Exit & Handoff:** Do NOT edit `magic-modules`. Report failure as a GCP service-side internal error or test environment quota exhaustion requiring human investigation by the Terraform team. |

---

## Detailed Remediation Recipes

### Scenario 1: State Drift / Non-Empty Plan
* **Inspection:** Check the GET response JSON and `outline.txt`.
* **Remedies:**
  - Add `diff_suppress_func` in `mmv1/products/<product>/<Resource>.yaml` for equivalent string/IP/JSON formats.
  - Set `default_from_api: true` if the API sets default values post-create.
  - Fix custom flatteners in `templates/terraform/custom_flatten/` if state parsing fails.

### Scenario 2: API 400 Payload Mismatch
* **Inspection:** Compare `01_POST_request.json` against the GCP REST/gRPC schema.
* **Remedies:**
  - Verify field naming in YAML (`api_name` vs `name`).
  - Add `send_empty_value: true` if empty strings/zero-values are required by the API.

### Scenario 3: LRO & Eventual Consistency Timing
* **Inspection:** Check timing between `POST` and subsequent `GET` requests in `outline.txt`.
* **Remedies:**
  - Configure `async:` block in YAML for long-running operations.
  - Implement polling waiters in handwritten overrides (`mmv1/third_party/terraform/services/<service>/`).

### Scenario 4: Test Template Setup & Name Collisions
* **Inspection:** Check test configuration in `..._test.go.tmpl` or sample HCL templates.
* **Remedies:**
  - Replace static resource names with `{{index $.ResourceIdVars "var_name"}}` and define `resource_id_vars` in product YAML.
  - If test depends on beta-only resources (e.g. `google_project_service_identity`), wrap test with `{{- if ne $.TargetVersionName "ga" }}` ... `{{- end }}`.

### Scenario 5: Non-Remediable Backend / Quota Errors
* **Action:** Early exit immediately. Do not modify `magic-modules`. Report failure to human maintainers.

---

## Guidelines for Adding New Scenarios

When a new failure pattern is identified:
1. Assign a new sequential Scenario ID (e.g., Scenario 6).
2. Document the symptom pattern, root cause category, and primary file location / remedy.
3. Keep the catalog updated so all workflows can consult the latest patterns.
