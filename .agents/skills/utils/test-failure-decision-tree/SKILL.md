---
name: test-failure-decision-tree
description: "Classification decision tree and remediation strategies for diagnosing Terraform acceptance test failures across all workflows."
---

# `test-failure-decision-tree`

This document provides a standardized classification catalog for diagnosing, isolating, and fixing failing Terraform acceptance tests in `magic-modules`. It is designed to be consulted by any workflow (`test_fix`, `new_resource`, `add_list_resource`, `default`) whenever `make testacc` fails.

---

## Global Architectural Rule: Generated vs. Handwritten

Every Terraform resource and acceptance test in `magic-modules` is either **Generated (MMv1)** or **Handwritten (Go overrides)** as documented in `docs/content/test/test.md`. Before applying any remediation recipe in this catalog, always verify whether the source is generated (schemas in `mmv1/products/`, tests in `mmv1/templates/terraform/samples/`) or handwritten (Go code and tests in `mmv1/third_party/terraform/services/`) so you edit the correct source files. Consult `docs/content/develop/diffs.md` for diff resolutions, `docs/content/test/test.md` for test conventions, and `docs/content/test/run-tests.md` for compilation and execution guides.

---

## Global Handwritten Test Naming Convention & Execution Rule

For **generated acceptance tests** (MMv1 samples in `mmv1/templates/terraform/samples/`), Magic Modules automatically generates test function names (`TestAcc...`) during `make provider`.

However, when adding any new **handwritten Go test function** (`resource_PRODUCT_RESOURCE_test.go` or other `*_test.go` files under `mmv1/third_party/terraform/services/<service>/`), always verify that the test function name matches its execution tier:
- **Handwritten Unit Tests** (in-memory logic, diff suppressors, custom decoders/encoders, flatteners/expanders, validation): MUST be named `Test<FunctionName>` (e.g., `TestDiscoveryEngineSearchEngineFeaturesDiffSuppress`) **without** the `TestAcc` prefix.
- **Handwritten Acceptance Tests** (live API CRUD tests, VCR-recorded tests): MUST be named `TestAcc<Resource>_<Scenario>`.

`TestAcc` is strictly reserved for acceptance tests run via `make testacc`. Naming a handwritten unit test with `TestAcc` causes it to be filtered out of standard unit test runs and misclassified by CI.

---

## Decision Matrix Catalog

| Scenario ID | Symptom / Error Pattern | Root Cause Category | Primary File Location & Remedy |
| :--- | :--- | :--- | :--- |
| **Scenario 1** | gRPC Code 8 (`RESOURCE_EXHAUSTED`) / HTTP 429 (`Quota limit exceeded`/`Quota exhausted`), gRPC Code 13 (`INTERNAL`) / Code 14 (`UNAVAILABLE`) / HTTP 500/502/503 (`An internal error has occurred`), or API backend crash | Non-Remediable GCP Backend / Quota Environment Error | **Early Exit & Handoff:** Do NOT edit `magic-modules`. Report failure as a GCP service-side internal error or test environment quota exhaustion requiring human investigation by the Terraform team. *(Excludes HTTP 403 SERVICE_DISABLED -> see Scenario 8)* |
| **Scenario 2** | State drift, `ImportStateVerify` mismatch, non-empty plan after Apply, or server-injected metadata diffs (e.g. `run.googleapis.com/...`) | State Normalization / Default Mismatch / Server-Injected Annotation | Edit product YAML `mmv1/products/<product>/<Resource>.yaml` (for generated resources) or edit handwritten Go code `mmv1/third_party/terraform/services/<service>/resource_<resource_name>.go` (for handwritten resources) to add diff suppress (including regex suppression), adjust state reader, or handle defaults |
| **Scenario 3** | HTTP 400 `InvalidArgument`, `Unknown Field`, or field serialization error | API Request Payload Schema Mismatch | Compare POST request JSON against API proto schema; fix field `camelCase`/`snake_case` mapping, decoder/expander, or `send_empty_value` |
| **Scenario 4** | HTTP 404 immediately post-Create, eventual consistency delay, or 409 conflict | LRO / Eventual Consistency Timing | Configure `async:` / `autogen_async` in YAML, add retry handlers (`transport_tpg.SendRequestWithTimeout`), or polling waiter logic in handwritten overrides |
| **Scenario 5** | Resource name collision (`alreadyExists`), static string reuse, or test dependency setup bug | Acceptance Test Naming & Setup Bug | Edit sample `.tf.tmpl` + resource YAML (for generated tests) or edit Go test `.go` (for handwritten tests) to use dynamic random strings and bootstrap dependencies |
| **Scenario 6** | HTTP 400 `InvalidParameter: ... Not valid version` or outdated image/version string | Outdated / Deprecated API Parameter Version | Update outdated version strings in sample `.tf.tmpl` or handwritten test context maps to currently supported API versions |
| **Scenario 7** | HTTP 400 `managedZoneDnsNameNotAvailable` or domain allocation limit error across concurrent tests | Shared Resource / Domain Allocation Conflict | Consolidate multiple separate test functions into a single acceptance test function that reuses the managed domain/resource |
| **Scenario 8** | HTTP 403 `SERVICE_DISABLED` or error indicating API service has not been used in project before / is disabled | API Service Not Enabled in Target Project | For shared CI test projects (`ci-test-project-nightly-ga`/`beta`): run `gcloud services enable <service.googleapis.com> --project=<ci-project>`. For test-created or secondary projects: use Terraform resources (`google_project_service` / Go bootstrap utilities) in test config. |
| **Scenario 9** | `Requested interaction not found` in VCR `REPLAYING` mode after `RECORDING` passed | Non-Deterministic HTTP Request Ordering or Payload | For user-facing doc samples (`mmv1/templates/terraform/samples/` without `exclude_basic_doc: true`): set `skip_vcr: true` with a `# reason` comment in product YAML. For test-only configs: order concurrent resources with `depends_on` or replace unseeded randomness with `acctest.RandString(t, 10)`. |
| **Scenario 10** | `Error waiting for Creating/Updating/Deleting <Resource>: timeout while waiting for state to become ...` (e.g. `'done: true'` for standard GCP LROs, or `'DONE'` for Compute) | Long-Running Operation (LRO) Exceeding Resource Timeout | Retrieve the operation ID/URL from debug logs and poll/query the operation to find its total length (e.g. `metadata.endTime` - `metadata.createTime` or `endTime` - `startTime`); bump the resource timeouts (`timeouts:` in product YAML or `schema.ResourceTimeout` in handwritten Go + docs) to an appropriate higher interval (`20`, `30`, or a multiple of `30` minutes such as `60`, `90`, `120`). |

---

## Detailed Remediation Recipes

### Scenario 1: Non-Remediable Backend / Quota Errors
* **Inspection:** Check for GCP backend crashes, quota exhaustion (gRPC Code 8 `RESOURCE_EXHAUSTED` / `Quota limit exceeded` / HTTP 429), or internal errors (gRPC Code 13 `INTERNAL` / Code 14 `UNAVAILABLE` / HTTP 500/502/503).
  - **Exclusion / Actionability Check:** Before classifying a gRPC Code 13 (`INTERNAL`) or HTTP 500 error as a non-remediable Scenario 1 crash, inspect the error message carefully. If the message references specific schema fields, parameters, or value constraints (e.g. `invalid value`, `missing field foo`, `cannot be parsed`), it is likely a client-side payload error that triggered a backend crash. Classify it as **Scenario 3 (Payload Mismatch)** instead of exiting.
  - **Note:** Do NOT classify HTTP 403 `SERVICE_DISABLED` or API Not Enabled errors as Scenario 1; those must be classified as **Scenario 8 (API Service Not Enabled)** for automated remediation.
* **CI Policy Alignment:** In nightly CI failure tracking (`magician`), Quota errors are classified as **Terraform Team Owned Errors** (`teamOwnedErrorTypes: "Quota"`), which automatically label tickets with `service/terraform` and assign them to the Release Shepherd.
* **Action:** Early exit immediately. Do not modify `magic-modules`. Report failure to human maintainers as a non-remediable test environment or service-side error.

### Scenario 2: State Drift / Non-Empty Plan (Permadiff)
* **Inspection:** Check the GET response JSON, `stdout` plan diff, and `outline.txt`.
* **Documentation Reference:** Consult `docs/content/develop/diffs.md` for tabbed code examples of MMv1 and Handwritten diff suppression, default values, and `ignore_read`.
* **Remedies:**
  - **API-Side Default (`default_from_api: true`)**: For existing/dynamic fields where the API returns a server-side default if unset, set `default_from_api: true` in product YAML (converts to `Optional: true, Computed: true` in schema).
  - **Client-Side Default (`default_value: ...`)**: Use `default_value` only for static, stable default values that match the API default.
  - **API Empty-if-Default Flattener (`default_if_empty.tmpl`)**: If the API returns empty/nil when default value is sent, use `custom_flatten: 'templates/terraform/custom_flatten/default_if_empty.tmpl'` (or check `v == nil || IsEmptyValue` in handwritten flatteners).
  - **Diff Suppression (`diff_suppress_func`)**:
    - Use `tpgresource.CaseDiffSuppress` for capitalization normalization differences.
    - Use `tpgresource.DurationDiffSuppress` for duration format differences ("60.0s" vs "60s" for `google.protobuf.Duration` fields).
    - Use `tpgresource.ProjectNumberDiffSuppress` for project ID vs project number diffs.
  - **Unreturned / Secret Fields (`ignore_read: true`)**: For fields never returned in GET responses (passwords, secrets), set `ignore_read: true` and `sensitive: true` in YAML. In handwritten tests, set flattened value from `d.Get(...)` and add the field to `ImportStateVerifyIgnore`.
  - **List/Array Element Reordering**: Use custom flatteners calling `tpgresource.SortStringsByConfigOrder` or `tpgresource.SortMapsByConfigOrder` so API response ordering matches user config order.
  - **Avoid Creation DiffSuppress Traps**: Remove erroneous `DiffSuppressFunc: EmptyOrDefaultStringSuppress(...)` if it suppresses configured values on resource creation, causing state drift.
  - **Server-Injected Map Entries / Metadata Drift (Annotations, Labels, Features, Properties)**: For plan diffs caused by server-added annotations, labels, features, or properties not present in user config (e.g., `run.googleapis.com/...-disabled` or default feature flags):
    - **Do NOT use `ignore_read: true`** on `KeyValuePairs` / map attributes, as that disables drift detection for legitimate user-initiated modifications.
    - **Instead, define a custom `diff_suppress_func`** in `custom_code.constants` (`constants/<resource>.go.tmpl`) that returns `true` when the key matches a server-provided key AND `new == ""`, and also returns `true` for the map item count delta (`strings.Contains(k, "map_field.%")`). For all other keys or explicit user modifications (`new != ""`), return `false`.
    - **YAML Syntax Rule (`custom_code`)**: When adding `constants:` or other hooks to product YAML files, check if a `custom_code:` block already exists (e.g. for `encoder:`/`decoder:`). Merge into the existing block instead of adding a duplicate `custom_code:` root key.
  - **Handwritten Resources:** If the resource is fully handwritten (code located under `mmv1/third_party/terraform/services/<service>/resource_<resource_name>.go`), directly edit the schema definition or CRUD method implementation in the Go file to add `DiffSuppressFunc`, adjust the Read function, or handle defaults in the resource map.

### Scenario 3: API 400 Payload Mismatch & Decoder/Expander Errors
* **Inspection:** Compare `01_POST_request.json` against the GCP REST/gRPC schema or inspect custom decoder logic.
* **Remedies:**
  - Verify field naming in YAML (`api_name` vs `name`).
  - Add `send_empty_value: true` if empty strings/zero-values are required by the API.
  - Fix custom decoders or expanders in `templates/terraform/decoders/` or `templates/terraform/expanders/` if array indexing or struct mapping corrupts request payloads or state.

### Scenario 4: LRO & Eventual Consistency Timing
* **Inspection:** Check timing between `POST` and subsequent `GET` requests in `outline.txt`, or inspect 409/404 retry backtraces.
* **Remedies:**
  - **Long-Running Operations**: Configure `autogen_async: true` and `async:` settings in YAML for long-running operations.
  - **APIs Lacking LROs**: Use post-create sleep templates: `templates/terraform/post_create/sleep.go.tmpl`, `sleep_2_min.go.tmpl`, or `sleep_5_min.go.tmpl` in `custom_code.post_create` to allow API propagation.
  - **Handwritten Overrides**: Add retry wrappers (`transport_tpg.SendRequestWithTimeout` / `time.Sleep`) in handwritten overrides (`mmv1/third_party/terraform/services/<service>/`).

### Scenario 5: Acceptance Test Naming, Randomization & Setup Bug
* **Inspection:**
  - Check whether the failing test is **Generated** (`mmv1/templates/terraform/samples/<service>/<sample>.tf.tmpl`) or **Handwritten** (`mmv1/third_party/terraform/services/<service>/<resource>_test.go`).
  - **Beta-Only Resource Leaking into GA Check:** If a test fails in GA (`target_provider: "ga"`) with a lockfile/provider error (e.g. `provider registry.terraform.io/hashicorp/google-beta: required by this configuration`), or if the tested resource has `min_version: beta` in `mmv1/products/<product>/<Resource>.yaml`, immediately check if the handwritten test file in `mmv1/third_party/terraform/services/<product>/` is named `*.go` instead of `*.go.tmpl`. Do **NOT** attempt to promote the resource to GA or modify provider HCL blocks.
* **Documentation Reference:** Consult `docs/content/test/test.md` for tabbed code examples of MMv1 and Handwritten test naming, randomization, dependency bootstrapping, and VCR mode skipping.
* **Remedies for Generated Tests:**
  - Replace static resource names with `{{index $.ResourceIdVars "var_name"}}`.
  - Define `resource_id_vars` in the resource's product YAML (`mmv1/products/<product>/<Resource>.yaml`) to generate unique per-test identifiers.
  - **Bootstrapping Dependencies**: Use `bootstrap_iam` in YAML for project/org IAM bindings; use `kms.BootstrapKMSKey(t)` for KMS keys; use `tpgcompute.BootstrapSharedTestNetwork(t, "identifier")` for networks.
  - **VCR Skipping**: Set `skip_vcr: true` for tests incompatible with HTTP recording. Set `skip_func: acctest.SkipTestUntil(t, "YYYY-MM-DD")` for unreleased features.
* **Remedies for Handwritten Tests:**
  - Replace static resource names with dynamic random strings using `acctest.RandString(t, 10)` in test contexts (`"instance_name": fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))`) or `%{random_suffix}` variables.
  - **Bootstrapping Dependencies**: Use `kms.BootstrapKMSKey(t)`, `tpgcompute.BootstrapSharedTestNetwork(t, ...)`, and `resourcemanager.BootstrapIamMembers(t, ...)`.
  - **VCR Skipping**: Call `acctest.SkipIfVcr(t)` or `acctest.SkipTestUntil(t, "YYYY-MM-DD")`.
  - **Beta-Only Resource Version Guard**: If a test depends on beta-only resources (`min_version: beta`), rename the file from `.go` to `.go.tmpl` and wrap the test implementation with `{{ if ne $.TargetVersionName "ga" }}` ... `{{ end }}` (do not generate an `{{ else }}` placeholder comment in GA).

### Scenario 6: Outdated / Deprecated API Parameter Version
* **Inspection:** Check API error message for invalid parameter versions (e.g., `oggoracle:21.18...` or deprecated image families).
* **Remedies:**
  - Update outdated version strings in sample `.tf.tmpl` or handwritten test context maps to valid, active supported versions.

### Scenario 7: Shared Resource / Domain Conflict (Test Consolidation)
* **Inspection:** Check if multiple test functions attempt to register identical domain names or global resources concurrently.
* **Remedies:**
  - Merge related test steps into a single test function so the shared domain or managed zone is created once and reused across steps.

### Scenario 8: API Service Not Enabled (SERVICE_DISABLED / HTTP 403)
* **Inspection:** Check whether the disabled API error (`Error 403: ... API has not been used in project ... before or it is disabled`) targets our shared CI test project (`ci-test-project-nightly-ga`, `ci-test-project-nightly-beta`, or standard shared test runner project) versus a secondary or custom test project created within the test itself.
* **CI Policy Alignment:** In nightly CI failure tracking (`magician`), API enablement errors are classified as **Terraform Team Owned Errors** (`teamOwnedErrorTypes: "API enablement (Test environment)"`), which automatically label tickets with `service/terraform` and assign them to the Release Shepherd.
* **Remedies:**
  - **Shared CI Test Project (`ci-test-project-nightly-ga` / `ci-test-project-nightly-beta`)**: Automatically execute `gcloud services enable <service.googleapis.com> --project=<ci-project>` via `run_command` to enable the API directly in the target CI project. Do **not** modify test code in `magic-modules` for shared CI project API enablement.
  - **Test-Created / Secondary Project**: If the disabled API is for a secondary or test-created project, modify the test configuration (`mmv1/templates/terraform/samples/<service>/<sample>.tf.tmpl` or handwritten Go test `mmv1/third_party/terraform/services/<service>/<resource>_test.go.tmpl`) to include a `google_project_service` resource or call test bootstrapping helper functions (`EnableServiceUsageProjectServices` / `bootstrap_test_utils`) so that the Terraform test itself enables the required API in that project.

### Scenario 9: VCR Replay Mismatch (`Requested interaction not found` in `REPLAYING` mode after `RECORDING` passed)
* **Inspection:** Diff HTTP sequences between `recording_test/<pkg>.log` and `replaying_after_recording_test/<pkg>.log` (from `gs://ci-vcr-cassettes/.../build-log/`). `NewVcrMatcherFunc` (`acctest/vcr_utils.go`) matches the first unused cassette interaction with identical `(Method, normalized URL, reflect.DeepEqual(JSON body))` — JSON arrays in request bodies are order-sensitive.
* **Common Root Causes:**
  1. **Batched requests with parallel resources:** e.g., multiple `google_project_service` resources where goroutine scheduling changes the `serviceIds` array order in `services:batchEnable`.
  2. **Concurrent fine-grained resources mutating the same parent:** e.g., multiple `google_*_iam_member` or `google_*_iam_binding` resources on the same parent, where read-modify-write scheduling alters intermediate `PUT`/`POST` bodies and etags.
  3. **Unseeded randomness or timestamps:** `random_*` provider, `uuid()`, `timestamp()`, or `id.UniqueId()` instead of `acctest.RandString(t, 10)`.
* **Remedies (Doc Samples vs. Test-Only Configs):**
  - **User-Facing Documentation Samples** (`mmv1/templates/terraform/samples/**/*.tf.tmpl` without `exclude_basic_doc: true`): Do **not** add test-only `depends_on` chains that clutter public docs. Instead, set `skip_vcr: true` in the sample block of `mmv1/products/<product>/<Resource>.yaml` with a `# reason` comment above it.
  - **Test-Only Configs** (handwritten `*_test.go` or samples with `exclude_basic_doc: true`): Fix non-determinism directly by chaining parallel resources that share a batch or parent policy with `depends_on`, or replacing unseeded randomness with `acctest.RandString(t, 10)`. Only use `acctest.SkipIfVcr(t)` / `skip_vcr: true` (with a `# reason` comment) when non-determinism cannot be ordered (e.g., clients bypassing the shared VCR HTTP transport or `ExternalProviders`).

### Scenario 10: Operation Timeout (`timeout while waiting for state to become ...`)
* **Inspection:**
  1. **Retrieve Operation ID & URL:** Inspect the parsed debug log (`outline.txt`) or raw debug log from `gs://nightly-test-data/...` to locate the operation resource name / URL returned by the initial `POST`/`PATCH`/`PUT`/`DELETE` request and polled by subsequent `GET` requests (e.g., `https://<service>.googleapis.com/v1/projects/<project>/locations/<location>/operations/<op_id>` for standard GCP LROs, or `https://compute.googleapis.com/compute/v1/projects/<project>/<scope>/operations/<op_id>` for Compute). Record the operation's creation/start timestamp (`metadata.createTime` / `createTime` for standard LROs, or `insertTime` / `startTime` for Compute).
  2. **Poll / Query the Operation for Total Length:** Poll or query the operation directly via REST API (`curl -H "Authorization: Bearer $(gcloud auth print-access-token)" "<operation_url>"`) or `gcloud` (`gcloud <product> operations describe ...`) to check whether it finished (`"done": true` for standard LROs, `"status": "DONE"` for Compute, or the target resource state) and retrieve its completion timestamp (`metadata.endTime` / `updateTime` / `endTime`).
  3. **Calculate Total Length:** Compute the total elapsed duration from creation/start (`metadata.createTime` / `createTime` / `startTime` / `insertTime`) to completion (`metadata.endTime` / `updateTime` / `endTime`).
     - If the operation ran in a temporary test-created project that was deleted during test teardown (returning `404`/`403`), check recent nightly history and sibling tests for the same resource to gauge typical duration.
     - If the operation remains in-flight (`"done": false` / `"status": "RUNNING"`) indefinitely on the backend or terminated with a backend internal error, classify as **Scenario 1 (Non-Remediable Backend Error)** and report it rather than inflating timeouts.
* **Remedies:**
  - **Choose Target Interval:** Bump the resource timeout up to an appropriate higher interval that comfortably exceeds the total operation length — generally **`20`, `30`, or a multiple of `30` minutes** (e.g., `30`, `60`, `90`, `120`).
  - **MMv1 Generated Resources** (`mmv1/products/<product>/<Resource>.yaml`): Update `timeouts:` (`insert_minutes`, `update_minutes`, `delete_minutes` — all default to `20`) for the timed-out action(s). If `async.operation.timeouts` is also present in the resource YAML, update it to match.
  - **Handwritten Resources** (`mmv1/third_party/terraform/services/<service>/resource_<resource>.go`): Update `Timeouts: &schema.ResourceTimeout{ Create: schema.DefaultTimeout(30 * time.Minute), ... }` (or Framework `timeouts.Create(ctx, 30*time.Minute)`) and update the corresponding `## Timeouts` section in `mmv1/third_party/terraform/website/docs/r/<resource>.html.markdown`.

---

## Guidelines for Adding New Scenarios

When a new failure pattern is identified:
1. Assign a new sequential Scenario ID (e.g., Scenario 11).
2. Document the symptom pattern, root cause category, and primary file location / remedy.
3. Keep the catalog updated so all workflows can consult the latest patterns.
