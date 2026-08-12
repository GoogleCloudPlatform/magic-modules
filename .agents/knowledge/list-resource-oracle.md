---
name: list-resource-oracle
description: "Living knowledge base of every failure pattern, gotcha, and fix encountered during list-resource generation work. Read this BEFORE starting any list-resource task. Every entry represents a real retry that happened — follow the documented fix to avoid repeating it."
---

# List-Resource Oracle

> **Note to AI Agents:** Read this document in full before beginning any list-resource task. Every
> entry below is sourced from a real Coder retry. Following these patterns eliminates the most common
> causes of validation failure before they occur.

This oracle covers three surfaces:

1. **Eligibility** — which resources can and cannot receive `generate_list_resource: true`
2. **Generator bugs** — issues in templates that require a fix on the oracle branch before generation works
3. **Test configuration** — how the generated list-query test must be wired to pass

When a new failure pattern is found, add it here on the `update-list-resource-oracle` branch following
the format below and open a PR to `GoogleCloudPlatform/magic-modules:main`.

---

## Pattern Catalog

---

### P-01 — `collection_url_key` mismatch: generator uses wrong response key

**Symptom:**
`make provider` succeeds and `go build ./...` passes, but the generated `List*s` function reads from the
wrong key in the API JSON response. At runtime (or in the generated test) the list returns empty or panics.
May also surface as a compile-time `undefined` if the camelized default key doesn't match any field in the
generated struct.

**Root cause:**
The MMv1 generator derives the list response key by camelizing the resource name. When the API uses a
different key name (e.g. `associations` instead of `networkFirewallPolicyAssociations`, or `apiProduct`
instead of `apiProducts`), the generated code silently reads from a key that doesn't exist in the response.

**Real examples:**
- `NetworkFirewallPolicyAssociation` — default would be `networkFirewallPolicyAssociations`; API returns
  key `associations`. Fix: `collection_url_key: 'associations'`
- `NetworkFirewallPolicyPacketMirroringRule` — default would be `networkFirewallPolicyPacketMirroringRules`;
  API returns key `packetMirroringRules`. Fix: `collection_url_key: 'packetMirroringRules'`
- `ApigeeApiProduct` — default would be `apiProducts`; API returns key `apiProduct`.
  Fix: `collection_url_key: 'apiProduct'`

**Fix:**
Add `collection_url_key: '<exact-api-response-key>'` as a top-level YAML field alongside
`generate_list_resource: true`. To find the correct key, inspect the real API response:

```bash
# Use gcloud or curl to GET the list endpoint and inspect the top-level key names
curl -s -H "Authorization: Bearer $(gcloud auth print-access-token)" \
  "https://<service>.googleapis.com/<version>/projects/$PROJECT/<resource-plural>" | python3 -m json.tool | head -20
```

**Do NOT:**
Assume the camelized resource name matches the API key. Verify against the actual API response or the
resource's proto definition before setting `generate_list_resource: true`.

---

### P-02 — `strconv` import missing when resource has Integer identity property

**Symptom:**
`go build ./...` fails with:

```
undefined: strconv
```

in the generated `list_<resource>.go` file.

**Root cause:**
The `list_resource.go.tmpl` template unconditionally emits code that calls `strconv.FormatInt` when an
identity property is of type `Integer`, but the `import` block previously did not include `"strconv"`
conditionally. This is a generator template bug — it affects any resource whose identity field (the field
used to uniquely identify a list item) is typed `Integer`.

**Real example:**
Any resource with an `Integer`-typed identity property (e.g. a priority field used as the resource ID).
Fix was to gate the `"strconv"` import on whether any identity property has `type: Integer`.

**How to detect before generation:**
Check the resource YAML for identity properties with integer types:

```bash
grep -A2 "type: Integer" mmv1/products/<product>/<Resource>.yaml
```

If any identity property is `Integer` typed and the template fix (P-02) is not yet in the generator,
exclude this resource from the batch and open an oracle PR first.

**Fix (oracle branch — generator template):**
In `mmv1/templates/terraform/list_resource.go.tmpl`, gate the `strconv` import:

```go
{{- $hasIntegerIdentity := false -}}
{{- range $id := $.IdentityProperties -}}
  {{- if eq $id.Type "Integer" -}}
    {{- $hasIntegerIdentity = true -}}
  {{- end -}}
{{- end -}}

import (
    "context"
    "fmt"
    {{- if $hasIntegerIdentity }}
    "strconv"
    {{- end }}
    ...
)
```

**Do NOT:**
Add `strconv` unconditionally to work around the missing import — that produces an "imported and not used"
compile error on resources with no Integer identity properties.

---

### P-03 — `types` import emitted unconditionally, fails for resources with no scope properties

**Symptom:**
`go build ./...` fails with:

```
"github.com/hashicorp/terraform-plugin-framework/types" imported and not used
```

in the generated `list_<resource>.go` file.

**Root cause:**
The `list_resource.go.tmpl` template unconditionally emits `"github.com/hashicorp/terraform-plugin-framework/types"`
in the import block. Resources with no list scope properties (i.e. `ListScopeProperties` is empty) have an
empty `ListModel` struct and never reference `types.String`, making the import unused.

**Real example:**
`ApigeeAddonsConfig` and `ApigeeOrganization` — both have no scope params in their list URL beyond the
implicit org identifier, so `ListScopeProperties` is empty.

**Fix (oracle branch — generator template):**
In `mmv1/templates/terraform/list_resource.go.tmpl`, gate the types import:

```go
{{- if $.ListScopeProperties }}
"github.com/hashicorp/terraform-plugin-framework/types"
{{- end }}
```

**Do NOT:**
Add a blank `_` import alias or remove the struct fields to suppress the error — fix the template.

---

### P-04 — Region/zone in sample config is hardcoded, mismatches test context injection

**Symptom:**
The generated list-query test (`TestAcc<Resource>ListQuery_generated`) fails with a region or zone
mismatch error at test runtime:

```
Error: region "us-east1" does not match the configured test region "us-central1"
```

or the test creates a resource in one region but the list query targets a different region.

**Root cause:**
The sample `.tf.tmpl` file hardcodes a region string (e.g. `region = "us-east1"`). The generated query
test injects `envvar.GetTestRegionFromEnv()` for the region scope parameter, which reads from the
`GOOGLE_REGION` environment variable. If the hardcoded region doesn't match the env var, the resource is
created in one region but listed in another, causing the test to return zero results or a 404.

**Real examples:**
- `NetworkEdgeSecurityService` — sample had `region = "us-east1"` hardcoded; test env uses `us-central1`.
  Fix: moved region to `vars: region: 'us-east1'` so the query test reads the same value.
- `RegionTargetTcpProxy` — sample had `region = "europe-west4"` hardcoded; test region env is different.
  Fix: moved region to `test_context_vars: region: envvar.GetTestRegionFromEnv()` so both the resource
  creation and the list query use the same runtime-resolved region.
- `RegionBackendBucket` — region was in `resource_id_vars` (wrong map); needed to move to `vars`.

**Fix:**
When the sample config contains a hardcoded region or zone that must match the test environment, move it
out of the HCL string into the YAML `vars` or `test_context_vars` map so the query test template can
consume the same value.

Use `vars` when the value is a static string you want injected consistently:
```yaml
samples:
  - name: my_sample
    steps:
      - name: my_config
        vars:
          region: 'us-central1'
```

Use `test_context_vars` when the value should be resolved at test runtime from an env var:
```yaml
samples:
  - name: my_sample
    steps:
      - name: my_config
        test_context_vars:
          region: envvar.GetTestRegionFromEnv()
```

Then in the `.tf.tmpl` file, reference the value via the appropriate template accessor:
- For `vars`: `{{index $.Vars "region"}}`
- For `test_context_vars` / `resource_id_vars`: `%{region}` (interpolated by the test framework)

**Do NOT:**
Use `resource_id_vars` for region or zone — that map is for resource name identifiers, not location
parameters. Values in `resource_id_vars` are not injected into the query test context map.

---

### P-05 — `resource_id_vars` used for region/location instead of `vars`

**Symptom:**
Region or zone value is set under `resource_id_vars` in the YAML sample definition. The generated query
test ignores it and falls back to `envvar.GetTestRegionFromEnv()`, causing a region mismatch (see P-04).

**Root cause:**
`resource_id_vars` is specifically for values used to construct the resource's identity (its import ID /
self-link), not for injecting values into the test context. The query test template only reads from `vars`
and `test_context_vars` when building the `contextMap` passed to `config.TestContextFunc`.

**Real example:**
`RegionBackendBucket` had `region: 'us-central1'` under `resource_id_vars`. The query test template did
not pick it up. Fix: moved to `vars: region: 'us-central1'`.

**Fix:**
Audit the sample YAML. Any `region`, `zone`, or `location` entry under `resource_id_vars` that is also
referenced in the `.tf.tmpl` as a location parameter must be moved to `vars` (static) or
`test_context_vars` (env-var-resolved).

---

### P-06 — `TestContextVars` override not respected; auto-injection overwrites resource-specific value

**Symptom:**
The query test injects `envvar.GetTestRegionFromEnv()` for the `region` scope parameter even when
`test_context_vars.region` is explicitly set in the YAML, because the template emits the auto-inject
block unconditionally before checking `test_context_vars`.

**Root cause:**
In `query_test_file.go.tmpl`, the scope-parameter injection loop checked `test_context_vars` after
emitting the auto-inject block, resulting in duplicate map entries. Go's map literal syntax rejects
duplicate keys, causing a compile error, or the last entry silently wins (wrong value).

The fix is to gate the auto-inject block with:
```
{{- if not (index $step.TestContextVars $n) }}
  ... auto-inject region/zone/project ...
{{- end }}
```

**Fix (oracle branch — generator template):**
In `mmv1/templates/terraform/samples/base_configs/query_test_file.go.tmpl`, wrap every auto-injection
branch inside a `{{- if not (index $step.TestContextVars $n) }}` ... `{{- end }}` guard so that an
explicit `test_context_vars` entry always wins over the default env-var injection.

**Do NOT:**
Remove `test_context_vars` from the YAML as a workaround — that just hides the template bug and leaves
the resource with the wrong region in future test runs.

---

### P-07 — `firewall_policy` reference uses `.name` instead of `.id` in sample

**Symptom:**
The list-query test fails at resource creation time:

```
Error: Invalid value for "firewall_policy": must be a self-link or id, not a name
```

**Root cause:**
The sample `.tf.tmpl` references the firewall policy via `.name` (just the short name string). The
`firewall_policy` field on rules and associations requires the full self-link or the resource `id`
attribute, which includes the project and path prefix.

**Real example:**
`NetworkFirewallPolicyPacketMirroringRule` sample had:
```hcl
firewall_policy = google_compute_network_firewall_policy.basic_network_firewall_policy.name
```
Fix:
```hcl
firewall_policy = google_compute_network_firewall_policy.basic_network_firewall_policy.id
```

**Fix:**
In the sample `.tf.tmpl`, for any field that the API expects as a self-link or full resource URL, use
the `.id` attribute of the referenced resource rather than `.name`. Check the API documentation for the
field type (`string` vs `resource reference`).

---

### P-08 — Bare-array API response requires `list_response_is_array: true`

**Symptom:**
`go build ./...` fails or the generated `List*s` function panics at runtime because it tries to extract
a named key from the response body, but the API returns a top-level JSON array with no wrapper object:

```json
["item1", "item2"]
```

instead of:

```json
{"items": ["item1", "item2"]}
```

**Root cause:**
The standard `ListPages` helper expects a wrapper object. Resources whose list endpoint returns a raw JSON
array are not compatible with the default generation path.

**Real examples:**
- `ApigeeEnvironmentKeyvaluemaps` — `GET {{env_id}}/keyvaluemaps` returns `["kvm1", "kvm2", ...]`
- `ApigeeTargetServer` — `GET {{env_id}}/targetservers` returns `[{...}, {...}]`

**Fix (oracle branch — generator + YAML):**
This requires both a template fix and a YAML flag:
1. In the resource YAML, add `list_response_is_array: true`.
2. In the generator templates (`resource.go`, `list_resource_method.go.tmpl`), gate on this flag to
   emit `ListArrayPages` calls instead of `ListPages`.
3. The `ListArrayPages` transport helper must exist in `transport/transport.go` — add it if missing,
   following the same seed-state isolation, rate-limit retry, and `Flattener` + `Callback` structure
   as `ListPages` but decoding `[]interface{}` directly.

These are generator-level changes. Add them to the `update-list-resource-oracle` branch. Do NOT set
`generate_list_resource: true` on these resources until the generator fix is merged upstream.

---

### P-09 — Custom `org` / `name` decoder needed when list items use different field for identity

**Symptom:**
The generated list data source returns items where the identity field (used to construct the resource
address) is populated from a different field than expected — e.g. the API response uses `"organization"`
but the resource identity uses `"name"`.

**Root cause:**
The list generator reads the identity field directly from the response item. When the API's list response
uses a different field name than the resource's identity field, the generated code reads from the wrong
key and produces empty or incorrect resource addresses.

**Real examples:**
- `ApigeeAddonsConfig` — list items have `"organization"` field; resource identity uses `"org"`.
  Fix: custom decoder that copies `res["organization"]` → `d.Set("org", ...)`.
- `ApigeeOrganization` — list items have `"organization"` field; resource identity uses `"name"`.
  Fix: custom decoder that copies `res["organization"]` → `res["name"]` (no-op on direct reads).

**Fix:**
Add a custom decoder template at
`mmv1/templates/terraform/decoders/<product>_<resource_snake>.go.tmpl` that remaps the field, and
reference it in the resource YAML under `custom_code.decoder`. This is a per-resource fix — it does
not affect the generator template.

---

### P-10 — Eligibility scan: do not force-remove exclusion flags to manufacture eligibility

**Symptom:**
A resource has `exclude_identity_generation: true` or `exclude_read: true`. The generator hard-fails
when `generate_list_resource: true` is also set on these resources.

**Root cause:**
These flags exist because the resource has genuine API constraints (e.g. the identity cannot be
auto-generated, or the resource is write-only and never returned by GET). Removing them to pass the
eligibility scan breaks the resource's existing behaviour.

**Fix:**
Do not touch these flags. Mark the resource as ineligible and move on. If the constraint is incorrect
and the API has changed, that is a separate PR — not part of the list-resource batch.

**Do NOT:**
Remove `exclude_identity_generation: true` or `exclude_read: true` to force a resource through the
eligibility scan. This will break the build or produce a non-functional data source.

---

### P-11 — Eligibility scan: non-standard scope params in list URL block generation

**Symptom:**
A resource's `base_url` contains template params beyond `project`, `region`, `zone`, `location` — for
example `{{disk}}`, `{{instance}}`, `{{parent}}`, `{{env_id}}`. The eligibility scan marks it as
ineligible with `"list URL has unsupported scope param(s): ['disk']"`.

**Root cause:**
The query test template only auto-injects `project`, `region`, `zone`, and `location` into the test
context map. Resources whose list URL requires additional path parameters (parent resource IDs) cannot
have their query test run unattended without extending the template to inject those values.

**Fix:**
Do not set `generate_list_resource: true` on these resources. If you want to make a resource with a
non-standard scope param eligible, first extend the `AUTO_SCOPES` set and the query test template
injection logic in an oracle PR, get that merged upstream, then include the resource in a subsequent
list-resource batch.

**Do NOT:**
Add the resource to the batch and hope the test passes. It will fail at test runtime with a missing
context variable error.

---

### P-12 — Downstream provider has uncommitted work before generation

**Symptom:**
`make provider` succeeds but subsequent `git status --porcelain` in the downstream provider shows
unexpected diffs in files unrelated to the current product. The Validator's `validate-provider-changes`
oracle reports spurious breaking changes or missing tests on resources you did not touch.

**Root cause:**
The downstream `terraform-provider-google` clone had uncommitted or unstaged changes from a previous
run before generation started. The generator overwrites those files, mixing old work into the diff.

**Fix:**
Before running `make provider`, always verify the downstream is clean:

```bash
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
git status --porcelain
```

If there are unexpected changes, stash or reset them before generating:

```bash
git stash   # or: git checkout .
```

Then re-run `make provider`.

**Do NOT:**
Proceed with generation when the downstream has uncommitted work. The oracle's `validate-provider-changes`
check will report false positives and the Validator will return FAIL with misleading feedback.

---

### P-13 — `generate_list_resource: true` committed alongside downstream generated files

**Symptom:**
The Validator's `validate-provider-changes` oracle fails with unexpected breaking changes, or the PR
reviewer rejects the PR because it contains generated `.go` files under
`google/services/<product>/list_*.go`.

**Root cause:**
The `add_list_resource` SKILL's guardrail says to stage only `mmv1/products/<product>/` YAML changes
in the magic-modules commit. Generated downstream files are throwaway artifacts. If they are committed
to the magic-modules branch, the PR fails CI and review.

**Fix:**
Before committing, verify staged files:

```bash
git status --porcelain
git diff --cached --name-only
```

Only `mmv1/products/<product>/*.yaml` should appear. If any `*.go`, `*.html.markdown`, or downstream
paths are staged, unstage them:

```bash
git restore --staged .
git add mmv1/products/<product>/
git commit -m "<product>: add list resources"
```

---

## Adding New Entries

When the Coder agent is retried due to a Validator FAIL, the failure pattern MUST be added here before
the oracle branch is merged. Follow this format exactly:

```markdown
### P-NN — Short title describing the failure

**Symptom:**
The exact error message or observable behaviour that triggered the retry.

**Root cause:**
Why it happens — which template, YAML field, or test configuration causes it.

**Real example:**
The specific resource(s) and commit(s) where this was first seen.

**Fix:**
The exact change required. Include file paths and code snippets where relevant.

**Do NOT:**
The shortcut that makes the symptom disappear without fixing the underlying cause.
```

Entries are numbered sequentially. Do not renumber existing entries. If a fix is superseded by a later
generator change, mark the entry `[RESOLVED in upstream as of <date>]` but do not delete it — historical
context is valuable.

---

### P-14 — HTTP 403 on acceptance test: diagnose before dropping

**Symptom:**
The generated list-query test (`TestAcc<Resource>ListQuery_generated`) fails with an HTTP 403
(Permission Denied) when calling the resource's list endpoint.

**Root cause:**
A bare 403 has four distinct root causes with different correct resolutions. Dropping the resource
immediately discards resources that are perfectly valid but happen to hit an environment or IAM gap:

1. **Required GCP API not enabled** in the test project (environment gap — resource is supportable).
2. **Org-scoped resource** — the list URL contains `/organizations/{org_id}` and `GOOGLE_PROJECT`
   alone is insufficient; org-level credentials are needed.
3. **Alpha/private feature or allowlist required** — the resource is gated behind an allowlist,
   alpha program, or feature flag not active in the test project.
4. **IAM permission gap** — the test service account lacks the required `<service>.<resource>.list`
   IAM role, but the role is a standard one that can be granted.

**Fix (Validator diagnosis protocol):**
Run these checks in order; stop at the first that matches:

1. Check API enablement:
   ```bash
   gcloud services list --project=$GOOGLE_PROJECT --enabled | grep -i '<service>'
   ```
   Not enabled → **HARD DROP** with reason `"API 403 — required GCP API not enabled: <service>"`.

2. Check `base_url` / `list_url` in the resource YAML for `/organizations/`:
   Not present → **HARD DROP** with reason `"API 403 — resource is org-scoped, requires GOOGLE_ORG"`.

3. Check the 403 response body in `outline.txt` for `allowlist`, `alpha`, `private feature`:
   Matched → **HARD DROP** with reason `"API 403 — requires allowlisting/alpha: <message>"`.

4. Check test SA IAM roles; compare to required `<service>.<resource>.list` permission:
   Missing standard role → return FAIL (environment issue, not a resource issue) so the SA can be
   granted the role and the test re-run. Do NOT drop.
   Requires org/special access → **HARD DROP** with reason `"API 403 — requires privileged access"`.

**Do NOT:**
Drop a resource on a bare 403 without running the four-step diagnosis. Dropping on an IAM gap
permanently excludes a valid resource from all future batches.

---

### P-15 — HTTP 404 on acceptance test: diagnose before dropping

**Symptom:**
The generated list-query test (`TestAcc<Resource>ListQuery_generated`) fails with an HTTP 404
(Not Found) when calling the resource's list endpoint, OR the resource creation step itself returns
404 before the list query is even reached.

**Root cause:**
A bare 404 has four distinct root causes. The most common is a region/zone mismatch (P-04/P-05),
which is YAML-fixable. Dropping on a 404 without diagnosis throws away fixable resources.

1. **Region/zone mismatch** — resource created in a different region/zone than the list query
   targets. The list returns 404 because it queries the wrong scope. This is oracle P-04/P-05.
2. **Resource type genuinely unavailable** — the resource type does not exist in the given
   project/region (e.g. a regional resource in a zone that doesn't support it).
3. **Wrong API version** — the resource uses an `alpha` or `beta` base URL not available in the
   test project.
4. **List URL differs from resource URL** — the list endpoint is at a different path than the
   CRUD endpoint. The generated `list_<resource>.go` may have derived the wrong list URL.

**Fix (Validator diagnosis protocol):**
Run these checks in order:

1. Region/zone mismatch — grep `outline.txt` for mismatched region/zone values between the POST
   (create) and GET (list) requests. If mismatch found: **YAML-fix P-04**, return FAIL to Coder.

2. Resource creation status — grep `outline.txt` for the creation HTTP status. If 404/403 on
   creation: resource type unavailable in this project/region. **HARD DROP** with reason
   `"API 404 — resource type not available in project/region: <url>"`.

3. API version — grep resource YAML for `min_version: beta` or `alpha` in base URL. If alpha/beta:
   **HARD DROP** with reason `"API 404 — resource requires alpha/beta API not available in test project"`.

4. List URL mismatch — if creation succeeded (201) but list returned 404, compare the list URL
   in the generated `list_<resource>.go` to the API documentation. If different: **YAML-fix**,
   return FAIL noting that `list_url` must be set explicitly in the resource YAML.

5. None of the above: **HARD DROP** with reason
   `"API 404 — root cause undetermined after diagnosis. See /tmp/debug_<resource>/outline.txt"`.

**Do NOT:**
Drop a resource on a bare 404 before running diagnosis step 1 (region mismatch check). The most
common 404 in list-query tests is P-04 — a fixable YAML issue, not a permanent ineligibility.

---

### P-16 — Generator template changes must be in a separate PR from YAML batch changes

**Symptom:**
The Coder or Validator modified a core generator template file (e.g. `list_resource.go.tmpl`,
`query_test_file.go.tmpl`) as part of the same branch that adds `generate_list_resource: true` to
a batch of YAML files. The PR is rejected in review because template changes affect all products
and require independent review.

**Root cause:**
Generator template changes (`mmv1/templates/terraform/*.go.tmpl`, excluding per-resource `examples/`
and `decoders/`) affect the generated output for every product — not just the one being batched. They
carry higher review risk and must be evaluated independently so reviewers can assess cross-product
impact. Bundling them with a YAML batch PR makes the diff harder to review and raises the chance of
an accidental regression in another product being silently merged.

**Files that trigger this rule (must be in a separate PR):**
- `mmv1/templates/terraform/list_resource.go.tmpl`
- `mmv1/templates/terraform/list_resource_method.go.tmpl`
- `mmv1/templates/terraform/samples/base_configs/query_test_file.go.tmpl`
- Any other `mmv1/templates/terraform/*.go.tmpl` not in `examples/` or `decoders/`

**Files that are allowed in the YAML batch PR (do NOT split these):**
- `mmv1/products/<product>/*.yaml` — the YAML edits
- `mmv1/templates/terraform/examples/<product>_*.tf.tmpl` — per-resource sample config fixes
- `mmv1/templates/terraform/decoders/<product>_*.go.tmpl` — per-resource custom decoders
- `mmv1/third_party/terraform/services/<product>/` — handwritten custom code

**Fix (Validator Step 6f):**
1. Identify template files in the diff (Step 1b).
2. Strip them from the YAML batch branch (`git checkout upstream/main -- <file>`).
3. Re-test any resource that relied on the template fix; if it no longer passes, drop it with
   reason `"requires generator template fix (P-NN) before it can be included in a YAML batch PR"`.
4. Create a separate branch from `upstream/main` containing only the template changes.
5. Open a PR for the template branch with `release-note:none`.
6. Reference the template PR in the YAML batch PR body as a dependency.

**Do NOT:**
Include both YAML batch changes and generator template changes in the same PR. The template PR
must merge and its CI must pass before the YAML batch PR that depends on it is considered ready.

---

### P-17 — Scope mismatch: list query filter does not match where the resource was created

**Symptom:**
Recording log contains `no query results found after filtering` or the test's
`ExpectLengthAtLeast(..., 1)` assertion fails because the list returns 0 items.

**Root cause:**
The `base_url` (used as the list collection URL) resolves a scope parameter (e.g. `{{location}}`,
`{{region}}`) to a different value than the one used when the resource was created in Step 1 of the
test. The list query searches a scope that doesn't contain the resource it just created.

**Real example:**
Any resource where `test_context_vars.location` or `vars.region` is set to a hardcoded value that
differs from the `envvar.GetTestRegionFromEnv()` value injected by the auto-scope logic.

**Fix:**
Verify `base_url` and `id_format` use the same scope tokens. If the resource is scoped to a
non-standard location (e.g. a hardcoded project or region in `test_vars_overrides`), ensure the list
query scope variables are captured from the created resource state (via `listScope.Capture`), not
from environment variables. See also P-04 and P-05 for the YAML-level fix.

**Do NOT:**
Add `skip_vcr: true` or `t.Skip(...)` to work around this. See P-19.

---

### P-18 — Resource has no list endpoint: 501, 404, or generator `has no list endpoint` error

**Symptom:**
The generator fails with `has no list endpoint`, OR the recording log contains HTTP 404 or 501
when the list URL is called.

**Root cause:**
The API does not support listing this resource type at the collection URL implied by `base_url`.

**Fix:**
Remove `generate_list_resource: true` from the YAML **only when** the API has confirmed no list
endpoint (verified by reading the API reference or observing a 501/404 on the collection URL).
This is the **sole** legitimate reason to remove the flag after it has been set.

**Do NOT:**
Remove the flag for compile errors, scope mismatches, or VCR cassette issues. See P-19.

---

### P-19 — `skip_vcr: true` or `t.Skip` is never a fix for list-resource tests

**Symptom:**
A fix attempt added `skip_vcr: true` to a sample in a resource YAML, or added `acctest.SkipIfVcr(t)`
or `t.Skip(...)` to a generated test file or the `query_test_file.go.tmpl` template.

**Root cause:**
This is the bug itself, not a fix. `skip_vcr: true` permanently opts the resource out of VCR CI.
It hides the real failure without resolving it.

**Fix:**
Revert the `skip_vcr: true` addition immediately. Diagnose the actual recording failure using P-17
or P-18. If the recording passed and only the replaying section shows failures, see P-20.

**Do NOT:**
Use `skip_vcr: true`, `acctest.SkipIfVcr(t)`, or any test-skip mechanism as a "fix" for a
CI failure. These are forbidden by the ABSOLUTE PROHIBITIONS in the CI-fix prompt.

---

### P-20 — Replaying "no cassette found" when recording passed

**Symptom:**
The orchestrator prints `replaying failures detected (recording passed — not a cassette issue)`.
The replaying section of the CI log contains `no cassette found on disk for TestAcc...ListQuery_generated`.

**Root cause:**
The recording run **succeeded** and wrote cassettes. The replaying run could not find them. This is
almost always caused by a non-deterministic component in the test context map that produces a different
cassette filename in replay than in recording. Common culprits: `acctest.RandString(t, N)` used in a
value that feeds into the VCR cassette path, or a test name that includes a random suffix.

**Real example (oracledatabase, 2026-08-12):**
All 10 oracledatabase `*ListQuery_generated` tests failed with `no cassette found` in the replaying
phase. The recording had passed. Prior fix attempts incorrectly diagnosed this as P-17 and applied
`skip_vcr: true` to 13 YAML samples — a direct P-19 violation. The `skip_vcr` additions must be
reverted. The real fix is to identify which context map value produces a non-deterministic cassette
path and make it deterministic (or seed it from the VCR cassette name).

**Fix:**
1. Confirm the orchestrator log says `recording passed`. Do NOT treat this as a recording failure.
2. Check the generated test's context map for values derived from `acctest.RandString`. Any such
   value that appears in the cassette filename must be either: (a) replaced with a deterministic
   value, or (b) anchored to the test name via `acctest.VcrRecorderPath(t)`.
3. Do NOT add `skip_vcr: true` or any skip mechanism.

**Do NOT:**
Treat P-20 as P-17 (scope mismatch). The recording passed — the resource was found. The problem is
cassette path resolution in replay, not a filter or scope issue.
