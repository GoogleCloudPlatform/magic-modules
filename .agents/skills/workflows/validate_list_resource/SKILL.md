---
name: validate-list-resource
description: "Validate a list-resource branch produced by the Coder agent: fetch the branch, run all oracle checks, run acceptance tests, write handwritten tests if needed, drop hard-ineligible resources, split template changes into a separate PR, open the YAML-only PR, and append a run entry to the oracle log. Invoke this skill when acting as the Validator agent in the list-resource autonomous loop."
---

# `validate-list-resource`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current task.

You are the **Validator** in the list-resource autonomous loop. The Coder has pushed a branch to the
fork. Your job is to independently verify that branch is correct, run acceptance tests, drop any
resources that are provably unsupportable, write any missing handwritten tests, enforce the
template-change separation rule, open the PR, and record this run in the oracle log.

If anything requires a YAML fix the Coder must make, return a structured JSON FAIL verdict so the
orchestrator can feed the exact failure back to the Coder. If a resource hits a **hard ineligibility**
(see Step 6b drop table), remove it yourself — do not loop it back to the Coder.

**Read the oracle first.** Before running any check, open
`.agents/knowledge/list-resource-oracle.md` and scan every pattern. The Coder may have already hit
one of these — knowing the patterns lets you give precise, actionable feedback instead of a raw error
dump.

---

## Prerequisites

* You are in the `magic-modules` root directory.
* `FORK_REMOTE`, `BRANCH`, and `PRODUCT` are provided by the orchestrator in the prompt.
* `GOPATH` is set and `terraform-provider-google` is checked out at
  `$GOPATH/src/github.com/hashicorp/terraform-provider-google`.
* `gh` CLI is authenticated (`gh auth status`).
* GCP credentials are available: `GOOGLE_PROJECT`, `GOOGLE_REGION`, `GOOGLE_ZONE`,
  `GOOGLE_APPLICATION_CREDENTIALS` (service account key file path) — required for acceptance tests.

---

## Step 0 — Read the oracle

```bash
cat .agents/knowledge/list-resource-oracle.md
```

Scan every pattern before evaluating any failure. When a check fails, cross-reference the output
against the catalog and include the matching pattern ID (e.g. `P-04`) in the feedback you return.

---

## Step 1 — Fetch, check out, and audit the diff

```bash
git fetch $FORK_REMOTE $BRANCH
git checkout $BRANCH
git diff upstream/main --name-only
```

### 1a — Downstream generated file guard (P-13)

If any file matching `google/services/` or ending in `.html.markdown` appears in the diff, fail
immediately:

```json
{"status":"FAIL","feedback":"P-13: downstream generated files committed to magic-modules branch. Files: <list>. Strip them with git rm --cached and force-push.","pr_url":""}
```

### 1b — Template change audit (P-16)

Check whether the diff touches any Go template files:

```bash
git diff upstream/main --name-only | grep 'mmv1/templates/terraform/'
```

**If template files are present in the diff**, record them in a `TEMPLATE_CHANGES` list. These files
**must not** be included in the YAML batch PR — they require their own PR (see Step 6f). Do not fail
here; continue the validation so you can determine whether the template change was actually necessary
for any of the resources. You will act on `TEMPLATE_CHANGES` in Step 6f.

Allowed in the YAML batch PR (do NOT flag these):
- `mmv1/products/$PRODUCT/*.yaml` — the YAML edits
- `mmv1/templates/terraform/examples/${PRODUCT}_*.tf.tmpl` — sample configs (P-04/P-05/P-06 fixes)
- `mmv1/templates/terraform/decoders/` — custom per-resource decoders (P-09 fix, custom code)
- `mmv1/third_party/terraform/services/$PRODUCT/` — handwritten custom code

Must be in a separate PR (flag these):
- `mmv1/templates/terraform/list_resource*.go.tmpl` — list resource Go templates
- `mmv1/templates/terraform/samples/base_configs/query_test_file.go.tmpl` — query test template
- Any other `mmv1/templates/terraform/*.go.tmpl` file not listed above

---

## Step 2 — Pre-generation oracle (Check 1)

```bash
./.agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh
```

Runs in parallel: submodule guard · gofmt · template validation (version-guard, unused-tmpl) ·
YAML lint · mmv1 unit tests · tools and .ci/magician unit tests. Exit 0 = pass.

Cross-reference failures against the oracle before reporting:
- YAML lint error on a new field → P-01, P-02, P-03
- gofmt failure → Coder edited a `.go` file it should not have touched

---

## Step 3 — Verify downstream is clean, then regenerate and build (Check 2)

Before generating, confirm the downstream has no uncommitted work (oracle **P-12**):

```bash
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
git status --porcelain
```

If any files are listed, return FAIL citing P-12 with the dirty file list before proceeding further.

Then generate and build:

```bash
cd <magic-modules-root>
make provider VERSION=ga OUTPUT_PATH=$GOPATH/src/github.com/hashicorp/terraform-provider-google PRODUCT=$PRODUCT
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google && go build ./...
```

Common compile failures mapped to oracle patterns:
- `undefined: strconv` → **P-02** (Integer identity property, missing gated import)
- `imported and not used: "types"` → **P-03** (no scope params, unconditional import)
- `undefined` on a list response key → **P-01** (wrong `collection_url_key`)
- panic or nil deref on list response → **P-08** (bare-array API response)

---

## Step 4 — Validate provider changes (Check 3)

```bash
cd <magic-modules-root>
./.agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh
```

Checks GA and Beta providers for: breaking schema changes · missing acceptance tests · missing
documentation. Exit 0 = pass.

If breaking changes appear on resources you did not touch: downstream was dirty → **P-12**.
If they appear on resources you did touch: the Coder introduced an incompatible schema change.

---

## Step 5 — Verify generated files exist (Check 4)

For every resource that gained `generate_list_resource: true` in the branch diff, confirm both
generated files are present in the downstream:

```bash
PROVIDER=$GOPATH/src/github.com/hashicorp/terraform-provider-google

for f in $(git diff upstream/main --name-only | grep "mmv1/products/$PRODUCT/" | grep -v product.yaml); do
  resource=$(basename "$f" .yaml | python3 -c "
import sys, re
s = sys.stdin.read().strip()
s = re.sub(r'(?<=[a-z0-9])(?=[A-Z])', '_', s).lower()
print(s)
")
  for suffix in '' '_generated_test'; do
    target="$PROVIDER/google/services/$PRODUCT/list_${resource}${suffix}.go"
    [ ! -f "$target" ] && echo "MISSING: $target"
  done
done
```

Any `MISSING` line = FAIL. Match against oracle patterns before reporting:
- Missing `.go` entirely → **P-01** (wrong key), **P-08** (bare array), **P-09** (identity mismatch)
- File present but empty → template rendering issue; report raw `make provider` output

---

## Step 6 — Run acceptance tests (Check 5)

Follow `.agents/skills/utils/run-acctests/SKILL.md` for each resource. Run from the downstream
provider root. Maintain two running lists throughout this step:

- **PASSING** — resources whose generated test exited 0
- **DROPPED** — resources removed from this PR with a recorded reason (see drop table in Step 6b)

### 6a — Run the generated list-query test for each resource

For each opted-in resource (snake_case name derived in Step 5):

```bash
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
TF_LOG=DEBUG make testacc \
  TEST=./google/services/$PRODUCT \
  TESTARGS='-run=TestAcc<ResourcePascalCase>ListQuery_generated$' \
  > /tmp/test_<resource>.log 2>&1
echo "exit: $?"
```

- Exit 0 → add to PASSING. Continue to next resource.
- Non-zero → diagnose using Step 6b before moving on.

### 6b — Classify every failure: YAML-fix · self-fix · 403/404 diagnosis · hard drop

When a test exits non-zero, parse the log first:

```bash
python3 .agents/scripts/tf_debug_parser.py /tmp/test_<resource>.log \
  --extract-dir /tmp/debug_<resource>
```

Read `outline.txt` and any referenced JSON payload files. Apply the classification table below in
order — use the **first row that matches**:

| Observed failure | Classification | Action |
|---|---|---|
| Wrong response key / list returns 0 items | YAML-fix — P-01 | Return FAIL to Coder with exact key and P-01 |
| Region/zone mismatch between create and list | YAML-fix — P-04/P-05/P-06 | If only sample `.tf.tmpl` wrong: fix inline + re-run. If YAML `vars` map wrong: return FAIL to Coder. |
| `firewall_policy` uses `.name` not `.id` | Self-fix — P-07 | Fix sample inline, regenerate, re-run |
| Identity field mismatch (response field ≠ identity field) | YAML-fix — P-09 | Return FAIL to Coder with exact field names |
| API returns HTTP **403** | Diagnose — see Step 6b-403 | See detailed protocol below |
| API returns HTTP **404** | Diagnose — see Step 6b-404 | See detailed protocol below |
| List URL has non-standard scope param | **HARD DROP** — P-11 | Drop: `"list URL requires unsupported scope param <param> — P-11"` |
| Compile error in generated `list_*.go` | **HARD DROP** | Drop: `"generator produced invalid Go — escalate to generator fix oracle branch"` |
| Bare-array API response | **HARD DROP** — P-08 | Drop: `"bare-array API response — P-08, requires generator fix before this resource can be included"` |
| Test flaked / timed out | Retry once (re-run 6a for this resource) | If flakes again → **HARD DROP**: `"transient failure after retry"` |
| Any other failure with no clear cause | **HARD DROP** | Drop: `"unresolvable failure — see /tmp/test_<resource>.log"` |

> **Rule:** Resolve all hard drops (Step 6e) before returning any YAML-fix FAIL. Never mix drops
> and FAIL feedback in the same verdict.

#### Step 6b-403 — HTTP 403 diagnosis protocol (P-14)

Do not drop on a bare 403. First run these checks in order; stop at the first one that resolves it:

1. **API enabled?** Verify the required GCP API is enabled in the test project:
   ```bash
   gcloud services list --project=$GOOGLE_PROJECT --enabled \
     | grep -i '<service-name>'
   ```
   If not enabled: this is an environment issue, not a resource issue. The resource may be supportable
   once the API is enabled. **HARD DROP** with reason:
   `"API 403 — required GCP API not enabled in test project: <service>. Enable it or use a project where it is active."`

2. **Org-level resource?** Inspect the resource's `base_url` in its YAML for `/organizations/` or
   `organizations/{{org_id}}`:
   ```bash
   grep -i 'base_url\|list_url' mmv1/products/$PRODUCT/<Resource>.yaml
   ```
   If the list URL requires an org ID: `GOOGLE_PROJECT` alone is insufficient.
   **HARD DROP** with reason:
   `"API 403 — resource is org-scoped and requires GOOGLE_ORG / org-level credentials not present in test environment."`

3. **Feature flag or allowlist?** Check the 403 response body in `outline.txt` for messages like
   `"requires allowlisting"`, `"not available in your organization"`, `"alpha/private feature"`:
   ```bash
   grep -i 'allowlist\|alpha\|private\|not available\|feature' /tmp/debug_<resource>/outline.txt
   ```
   If found: **HARD DROP** with reason:
   `"API 403 — resource requires allowlisting or alpha access not available in test project: <message>."`

4. **IAM permission gap?** If none of the above match, the service account likely lacks the required
   IAM role for the list operation. Check the API documentation for the required permission
   (usually `<service>.<resource>.list`) and verify the SA has it:
   ```bash
   gcloud projects get-iam-policy $GOOGLE_PROJECT \
     --flatten="bindings[].members" \
     --filter="bindings.members:$(gcloud config get-value account)" \
     --format="table(bindings.role)"
   ```
   If the SA is missing the role and it is a standard role (not org/allowlist-gated):
   this is a **YAML-fix opportunity** — return FAIL noting the required IAM role so the test
   environment can be corrected and the next loop can re-test. Do NOT drop.
   If the role is org-level or requires special access: **HARD DROP** with reason above.

#### Step 6b-404 — HTTP 404 diagnosis protocol (P-15)

Do not drop on a bare 404. Run these checks in order:

1. **Region/zone mismatch?** Check whether the resource was created in a different region/zone than
   the list query is targeting. This is the most common 404 cause for region/zone-scoped resources:
   ```bash
   grep -i 'region\|zone\|location' /tmp/debug_<resource>/outline.txt | head -20
   ```
   If the create request used a different region than the list request: this is a **YAML-fix** (P-04).
   Return FAIL to Coder with the exact mismatch.

2. **Resource created successfully?** Check `outline.txt` for the creation step's HTTP status:
   ```bash
   grep 'REQUEST_POST\|RESPONSE_POST\|201\|200\|404' /tmp/debug_<resource>/outline.txt | head -10
   ```
   If the resource was never created (creation itself returned 404 or 403): the resource type is
   unavailable in this project/region. **HARD DROP** with reason:
   `"API 404 — resource type not available in project $GOOGLE_PROJECT / region $GOOGLE_REGION. The API may not support this resource type in this region or it may require a specific project configuration."`

3. **Wrong API version?** Check the list URL in the resource YAML for `alpha` or `beta` base URL:
   ```bash
   grep 'base_url\|list_url\|min_version' mmv1/products/$PRODUCT/<Resource>.yaml
   ```
   If the resource uses an `alpha` or `beta` API endpoint that is not available in the test project:
   **HARD DROP** with reason:
   `"API 404 — resource uses alpha/beta API endpoint not available in test project: <url>."`

4. **Resource exists after creation?** If creation succeeded (201) but the list returns 404, the
   list endpoint itself may be in a different path than the resource creation endpoint. Inspect the
   list URL in the generated `list_<resource>.go` and compare to the API documentation. This is a
   **YAML-fix** (`list_url` field may need to be set explicitly):
   return FAIL with the exact list URL and the correct URL from docs.

5. **None of the above**: **HARD DROP** with reason:
   `"API 404 — could not determine root cause after diagnosis. See /tmp/debug_<resource>/outline.txt."`

### 6c — Write handwritten supplementary tests (when needed)

After the generated test passes, check whether a handwritten test already exists:

```bash
ls $GOPATH/src/github.com/hashicorp/terraform-provider-google/google/services/$PRODUCT/list_<resource>_test.go 2>/dev/null
```

Write a handwritten test **only when** at least one of these is true:
- `validate-provider-changes` (Step 4) flagged a missing acceptance test for this data source.
- The data source exposes filter parameters that the generated test does not exercise.
- The list URL has a non-standard scope that requires an integration test not covered by the
  generated template.

Handwritten test convention:
- Path: `google/services/$PRODUCT/list_<resource>_test.go` (downstream repo)
- Package: `<product>_test`
- Function: `TestAcc<ResourcePascalCase>ListQuery_<scenario>` (e.g. `_withFilter`)
- Must use `resource.ParallelTest` and `providertest.AccTestPreCheck`.

Run and confirm it passes before continuing:

```bash
TF_LOG=DEBUG make testacc \
  TEST=./google/services/$PRODUCT \
  TESTARGS='-run=TestAcc<ResourcePascalCase>ListQuery_<scenario>$' \
  > /tmp/test_<resource>_handwritten.log 2>&1
echo "exit: $?"
```

Parse any failure with `parse-debug-logs` before giving up.

### 6d — Commit sample fixes and force-push (if needed)

If you fixed any sample `.tf.tmpl` files during Step 6b (P-04/P-05/P-06/P-07):

```bash
cd <magic-modules-root>
git add mmv1/products/$PRODUCT/
git add mmv1/templates/terraform/examples/${PRODUCT}_*.tf.tmpl  # only if you changed samples
git commit --amend --no-edit
git push --force-with-lease $FORK_REMOTE $BRANCH
```

Do NOT stage generated `list_*.go` provider files — those stay downstream.

### 6e — Execute hard drops (if DROPPED list is non-empty)

For each resource in the DROPPED list:

1. Remove `generate_list_resource: true` (and `collection_url_key` if present) from its YAML.

2. After all drops, regenerate and verify the build is still clean:

```bash
make provider VERSION=ga OUTPUT_PATH=$GOPATH/src/github.com/hashicorp/terraform-provider-google PRODUCT=$PRODUCT
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google && go build ./...
```

3. Amend the commit and force-push:

```bash
cd <magic-modules-root>
git add mmv1/products/$PRODUCT/
git commit --amend --no-edit
git push --force-with-lease $FORK_REMOTE $BRANCH
```

4. If PASSING is now empty, skip to Step 8 and return:

```json
{"status":"FAIL","feedback":"All resources were dropped. No PR opened. Reasons: <resource>: <reason>; ...","pr_url":""}
```

### 6f — Template change separation (P-16)

If `TEMPLATE_CHANGES` (populated in Step 1b) is non-empty, enforce the separation rule now, before
opening the main PR.

**What belongs in the main YAML batch PR:**
- `mmv1/products/$PRODUCT/*.yaml`
- `mmv1/templates/terraform/examples/` sample `.tf.tmpl` fixes (P-04/P-05/P-06/P-07)
- `mmv1/templates/terraform/decoders/` custom decoder `.go.tmpl` files (P-09 custom code)
- `mmv1/third_party/terraform/` handwritten custom code

**What must be split into a separate template PR:**
- `mmv1/templates/terraform/list_resource*.go.tmpl`
- `mmv1/templates/terraform/samples/base_configs/query_test_file.go.tmpl`
- Any other core generator template (`.go.tmpl` files outside `examples/` and `decoders/`)

#### 6f-i — Strip template changes from the YAML branch

Remove the template file(s) from the current branch commit so the YAML batch PR stays clean:

```bash
# For each file in TEMPLATE_CHANGES:
git checkout upstream/main -- <template_file_path>
git add <template_file_path>
git commit --amend --no-edit
git push --force-with-lease $FORK_REMOTE $BRANCH
```

Any resource whose passing state **depended on** the template fix must now be re-tested without
the template change to confirm it still passes. If it now fails, classify the failure:
- If the failure is a hard drop (403, 404, bare array, scope param): drop it per Step 6e.
- If the failure is a YAML-fixable issue: return FAIL to Coder.
- If the resource only passes with the template change and cannot pass without it: the resource
  is **template-dependent** — add it to the DROPPED list with reason:
  `"resource requires generator template fix (P-NN) before it can be included in a YAML batch PR."`

#### 6f-ii — Open the template-fix PR on a separate branch

Create a dedicated branch for the template changes, targeting upstream main:

```bash
TEMPLATE_BRANCH="fix-list-resource-template-${PRODUCT}-$(date +%Y%m%d)"
git fetch upstream main
git checkout -b "$TEMPLATE_BRANCH" upstream/main

# Cherry-pick or re-apply only the template file changes:
for f in $TEMPLATE_CHANGES; do
  git checkout $BRANCH -- "$f"  # restore the template fix from the product branch
done

git add mmv1/templates/terraform/
git commit -m "list-resource templates: fix <issue> required for $PRODUCT batch"
git push $FORK_REMOTE "$TEMPLATE_BRANCH"
```

Then open a PR for the template branch:

```bash
gh pr create \
  --repo GoogleCloudPlatform/magic-modules \
  --base main \
  --head "$FORK_REMOTE:$TEMPLATE_BRANCH" \
  --title "list-resource templates: <fix description>" \
  --body "Template fix required for list-resource generation.

**Why separate from the YAML batch PR:**
Generator template changes affect all products and must be reviewed independently of any
single product's YAML batch. The YAML batch PR ($BRANCH) is held until this merges.

**Template files changed:**
$(for f in $TEMPLATE_CHANGES; do echo "- \`$f\`"; done)

**Oracle pattern:** P-NN — <pattern title>

\`\`\`release-note:none
\`\`\`
"
```

Capture the template PR URL. Include it in the oracle log (Step 8) and in the response JSON:

```json
{"status":"PASS","feedback":"Template changes split into separate PR: <template_pr_url>. YAML batch PR: <yaml_pr_url>.","pr_url":"<yaml_pr_url>","template_pr_url":"<template_pr_url>"}
```

---

## Step 7 — Open the YAML batch PR (only if PASSING is non-empty)

Follow `.agents/skills/operations/create-pr/SKILL.md` exactly.

Key requirements:
- Write the PR body to `/tmp/pr_body.txt` using a single-quoted HEREDOC.
- One `release-note:new-list-resource` block per resource in PASSING.
- Include a **Dropped Resources** section listing each dropped resource and its reason.
- If there is a pending template PR (from Step 6f), note it in the PR body:
  `"Note: this PR depends on <template_pr_url> merging first."`
- Open against `GoogleCloudPlatform/magic-modules:main` from `$FORK_REMOTE:$BRANCH`.
- After `gh pr create`, run `gh pr view` to verify release-note blocks rendered.
- Post `@modular-magician reassign-reviewer` as a PR comment.

PR body template:

```
Adds list-resource generation for the following $PRODUCT resources:

- `google_${PRODUCT}_<resource_a>`
- `google_${PRODUCT}_<resource_b>`

```release-note:new-list-resource
google_${PRODUCT}_<resource_a>
```

```release-note:new-list-resource
google_${PRODUCT}_<resource_b>
```

## Dropped Resources

The following resources were evaluated and excluded from this PR:

| Resource | Reason |
|---|---|
| `google_${PRODUCT}_<resource_x>` | API 403 — required GCP API not enabled in test project |
| `google_${PRODUCT}_<resource_y>` | list URL requires unsupported scope param `disk` — P-11 |
| `google_${PRODUCT}_<resource_z>` | resource requires generator template fix (P-02) — see <template_pr_url> |
```

---

## Step 8 — Append run entry to oracle log

After Step 7 (whether or not a PR was opened), append a `## Run:` entry to the oracle run log on
the `oracle/list-resource-patterns` branch. This is the permanent record of every run — future agents
read it before starting a new batch for any product that has appeared here.

### 8a — Switch to the oracle branch

```bash
ORACLE_BRANCH="oracle/list-resource-patterns"
ORACLE_FILE=".agents/knowledge/list-resource/list-resource-patterns.md"

if git ls-remote --exit-code $FORK_REMOTE "$ORACLE_BRANCH" > /dev/null 2>&1; then
  git fetch $FORK_REMOTE "$ORACLE_BRANCH"
  git checkout -B "$ORACLE_BRANCH" "$FORK_REMOTE/$ORACLE_BRANCH"
else
  git fetch upstream main
  git checkout -b "$ORACLE_BRANCH" upstream/main
fi
```

### 8b — Create or append

```bash
mkdir -p .agents/knowledge/list-resource
TODAY=$(date +%Y-%m-%d)
```

**First run (file does not exist):** create with frontmatter header then append the first `## Run:`.

```markdown
---
name: list-resource-patterns
description: "Append-only log of every add-list-resource run: what passed, what was dropped, and why. Read before starting a new batch for any product that appears here."
topics: [list-resource]
source: agent-generated — written by validate-list-resource skill after each run
---

# List-Resource Run Log
```

**Subsequent runs:** append only the new `## Run:` block at the end of the file.

```markdown
## Run: $PRODUCT — $TODAY

**PR:** <yaml_pr_url or "none — all resources dropped">
**Template PR:** <template_pr_url or "none">

**Passing resources:**
- `google_${PRODUCT}_<resource_a>`
- `google_${PRODUCT}_<resource_b>`

**Dropped resources:**
| Resource | Reason |
|---|---|
| `google_${PRODUCT}_<resource_x>` | <reason> |

**Observations:**
- <non-obvious finding, e.g. "403 on all org-scoped resources — test SA lacks org viewer">
- <or "none">
```

### 8c — Commit and push

```bash
git add "$ORACLE_FILE"
git commit -m "list-resource-patterns: $PRODUCT run $TODAY"
git push $FORK_REMOTE "$ORACLE_BRANCH"
git checkout $BRANCH
```

---

## Response format

Respond ONLY with a raw JSON object — no markdown fences, no prose before or after.

All checks passed, PR opened (no template split):
```json
{"status":"PASS","feedback":"","pr_url":"https://github.com/GoogleCloudPlatform/magic-modules/pull/<N>"}
```

All checks passed, PR opened, template changes split into separate PR:
```json
{"status":"PASS","feedback":"Template changes split into separate PR: <url>.","pr_url":"<yaml_pr_url>","template_pr_url":"<template_pr_url>"}
```

YAML fix required by Coder (no PR opened):
```json
{"status":"FAIL","feedback":"P-NN: <step>:\n<exact command>\n<full output>","pr_url":""}
```

All resources dropped (no PR opened, oracle log still written):
```json
{"status":"FAIL","feedback":"All resources dropped. Reasons: <resource>: <reason>; ...","pr_url":""}
```

Always include the oracle pattern ID when one matches.
