---
name: validate-list-resource
description: "Validate a list-resource branch produced by the Coder agent: fetch the branch, run all oracle checks, run acceptance tests, write handwritten tests if needed, drop hard-ineligible resources, open the PR, and append a run entry to the oracle log. Invoke this skill when acting as the Validator agent in the list-resource autonomous loop."
---

# `validate-list-resource`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current task.

You are the **Validator** in the list-resource autonomous loop. The Coder has pushed a branch to the
fork. Your job is to independently verify that branch is correct, run acceptance tests, drop any
resources that are provably unsupportable, write any missing handwritten tests, open the PR, and
record this run in the oracle log for future agents.

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

## Step 1 — Fetch and check out the branch

```bash
git fetch $FORK_REMOTE $BRANCH
git checkout $BRANCH
```

Verify the diff is scoped only to `mmv1/products/$PRODUCT/` YAML files:

```bash
git diff upstream/main --name-only
```

If any generated `.go` or `.html.markdown` files are present in the diff, fail immediately — that
violates the guardrail documented in oracle **P-13**:

```json
{"status":"FAIL","feedback":"P-13: downstream generated files committed to magic-modules branch. Files: <list>. Strip them with git rm --cached and force-push.","pr_url":""}
```

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
- **DROPPED** — resources removed from this PR with a recorded reason (see drop table below)

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

### 6b — Classify every failure: YAML-fix (return to Coder) vs hard drop (fix here)

When a test exits non-zero, parse the log first:

```bash
python3 .agents/scripts/tf_debug_parser.py /tmp/test_<resource>.log \
  --extract-dir /tmp/debug_<resource>
```

Read the generated `outline.txt` and any referenced JSON payload files. Then classify using this
table. The two columns are mutually exclusive — pick the first row that matches:

| Observed failure | Classification | Action |
|---|---|---|
| Wrong response key / list returns 0 items | YAML-fix — P-01 | Return FAIL to Coder with exact key and P-01 |
| Region/zone mismatch between create and list | YAML-fix — P-04/P-05/P-06 | If only sample `.tf.tmpl` is wrong: fix inline, regenerate, re-run. If YAML `vars` map is wrong: return FAIL to Coder. |
| `firewall_policy` uses `.name` not `.id` | YAML-fix — P-07 | Fix sample inline, regenerate, re-run |
| Identity field name mismatch in response | YAML-fix — P-09 | Return FAIL to Coder with exact field names |
| API returns HTTP **403** (permission denied) | **HARD DROP** | Drop with reason: `"API 403 — resource requires elevated permissions not available in test project"` |
| API returns HTTP **404** (resource type unknown) | **HARD DROP** | Drop with reason: `"API 404 — resource type not available in this project/region"` |
| List URL has non-standard scope param (not `project`/`region`/`zone`/`location`) | **HARD DROP** | Drop with reason: `"list URL requires unsupported scope param <param> — oracle P-11"` |
| Compile error in the generated `list_*.go` file | **HARD DROP** | Drop with reason: `"generator produced invalid Go — escalate to generator fix oracle branch"` |
| Bare-array API response (no wrapper object) | **HARD DROP** | Drop with reason: `"bare-array API response — oracle P-08, requires generator fix"` |
| Custom org/name decoder needed | **HARD DROP** | Drop with reason: `"identity field mismatch requires custom decoder — oracle P-09"` |
| Test flaked / timed out on first run | Retry once (re-run 6a for this resource only) | If it flakes again → **HARD DROP** with reason: `"transient failure after retry"` |
| Any other failure with no clear cause | **HARD DROP** | Drop with reason: `"unresolvable failure — see /tmp/test_<resource>.log"` |

> **Rule:** Never return a YAML-fixable failure AND a hard-drop in the same verdict. Resolve all
> hard drops first (Step 6e), then return FAIL for any remaining YAML issues.

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

1. Remove `generate_list_resource: true` (and `collection_url_key` if present) from its YAML:

```bash
# Edit mmv1/products/$PRODUCT/<ResourcePascalCase>.yaml
# Remove the generate_list_resource line and any adjacent collection_url_key line
```

2. After removing all dropped resources, regenerate to ensure the downstream is consistent:

```bash
make provider VERSION=ga OUTPUT_PATH=$GOPATH/src/github.com/hashicorp/terraform-provider-google PRODUCT=$PRODUCT
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google && go build ./...
```

3. Amend the commit to reflect the reduced YAML set and force-push:

```bash
cd <magic-modules-root>
git add mmv1/products/$PRODUCT/
git commit --amend --no-edit
git push --force-with-lease $FORK_REMOTE $BRANCH
```

4. If after drops PASSING is empty (every resource was dropped), **do not open a PR**. Skip to
   Step 8 (oracle log) and return:

```json
{"status":"FAIL","feedback":"All resources were dropped. No PR opened. Reasons: <list>","pr_url":""}
```

---

## Step 7 — Open the PR (only if PASSING is non-empty and all checks passed)

Follow `.agents/skills/operations/create-pr/SKILL.md` exactly.

Key requirements:
- Write the PR body to `/tmp/pr_body.txt` using a single-quoted HEREDOC (never inline backticks in
  `--body` — the create-pr SKILL explains why this silently strips release-note blocks).
- One `release-note:new-list-resource` block per resource in PASSING.
- Include a **Dropped Resources** section listing each dropped resource and its reason — this is the
  primary record for reviewers and for the oracle log.
- Open against `GoogleCloudPlatform/magic-modules:main` from `$FORK_REMOTE:$BRANCH`.
- After `gh pr create`, run `gh pr view` to verify the release-note block rendered.
- Post `@modular-magician reassign-reviewer` as a PR comment.

PR body template:

```
Adds list-resource generation for the following $PRODUCT resources:

- `google_<product>_<resource_a>`
- `google_<product>_<resource_b>`

```release-note:new-list-resource
google_<product>_<resource_a>
```

```release-note:new-list-resource
google_<product>_<resource_b>
```

## Dropped Resources

The following resources were evaluated but excluded from this PR:

| Resource | Reason |
|---|---|
| `google_<product>_<resource_x>` | API 403 — resource requires elevated permissions |
| `google_<product>_<resource_y>` | list URL requires unsupported scope param `disk` — oracle P-11 |
```

---

## Step 8 — Append run entry to oracle log

After Step 7 (whether or not a PR was opened), append a `## Run:` entry to the oracle run log.
This is the permanent record of what was tried, what passed, and what was dropped — future agents
read this before starting a new batch for the same product.

### 8a — Ensure the oracle branch exists and is checked out

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

### 8b — Create or append to the run log

```bash
mkdir -p .agents/knowledge/list-resource
TODAY=$(date +%Y-%m-%d)
PR_URL="<captured from Step 7, or 'none — all resources dropped'>"
```

**If the file does not exist yet** (first run ever), create it with the full header:

```markdown
---
name: list-resource-patterns
description: "Append-only log of every add-list-resource run: what passed, what was dropped, and why. Read before starting a new batch for any product that appears here."
topics: [list-resource]
source: agent-generated — written by validate-list-resource skill after each run
---

# List-Resource Run Log

```

Then append the first `## Run:` section below the header.

**If the file already exists**, append only the new `## Run:` section at the end:

```markdown
## Run: $PRODUCT — $TODAY

**PR:** $PR_URL

**Passing resources:**
- `google_${PRODUCT}_<resource_a>`
- `google_${PRODUCT}_<resource_b>`

**Dropped resources:**
| Resource | Reason |
|---|---|
| `google_${PRODUCT}_<resource_x>` | <reason> |

**Observations:**
- <any non-obvious finding from this run, e.g. "all zone-scoped resources in this product needed GOOGLE_ZONE set explicitly">
- <if nothing notable: "none">

```

### 8c — Commit and push to oracle branch

```bash
git add "$ORACLE_FILE"
git commit -m "list-resource-patterns: $PRODUCT run $TODAY"
git push $FORK_REMOTE "$ORACLE_BRANCH"
```

Then return to the product branch:

```bash
git checkout $BRANCH
```

---

## Response format

Respond ONLY with a raw JSON object — no markdown fences, no prose before or after.

All checks passed, PR opened, oracle log updated:
```json
{"status":"PASS","feedback":"","pr_url":"https://github.com/GoogleCloudPlatform/magic-modules/pull/<N>"}
```

Any check requires a YAML fix by the Coder (do NOT open a PR):
```json
{"status":"FAIL","feedback":"P-NN: <step name>:\n<exact command>\n<full output>","pr_url":""}
```

All resources dropped, no PR opened, oracle log still updated:
```json
{"status":"FAIL","feedback":"All resources dropped. Reasons: <resource>: <reason>; ...","pr_url":""}
```

Always include the oracle pattern ID in the feedback when one matches.
