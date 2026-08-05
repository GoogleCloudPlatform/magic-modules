---
name: validate-list-resource
description: "Validate a list-resource branch produced by the Coder agent: fetch the branch, run all oracle checks, run acceptance tests, write handwritten tests if needed, and open the PR if everything passes. Invoke this skill when acting as the Validator agent in the list-resource autonomous loop."
---

# `validate-list-resource`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current task.

You are the **Validator** in the list-resource autonomous loop. The Coder has pushed a branch to the
fork. Your job is to independently verify that branch is correct, run acceptance tests, write any
missing handwritten tests, and — only if everything passes — open the PR.
If anything fails, return a structured JSON verdict so the orchestrator can feed the exact failure back
to the Coder.

**Read the oracle first.** Before running any check, open
`.agents/knowledge/list-resource-oracle.md` and scan every pattern. The Coder may have already hit
one of these — knowing the patterns lets you give precise, actionable feedback instead of a raw error
dump.

---

## Prerequisites

* You are in the `magic-modules` root directory.
* `FORK_REMOTE` and `BRANCH` are provided by the orchestrator in the prompt.
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

Verify the diff is scoped only to `mmv1/products/<product>/` YAML files:

```bash
git diff upstream/main --name-only
```

If any generated `.go` or `.html.markdown` files are present in the diff, fail immediately — that
violates the guardrail documented in oracle **P-13**:

```
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
make provider VERSION=ga OUTPUT_PATH=$GOPATH/src/github.com/hashicorp/terraform-provider-google PRODUCT=<product>
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
PRODUCT=<product>
PROVIDER=$GOPATH/src/github.com/hashicorp/terraform-provider-google

for f in $(git diff upstream/main --name-only | grep "mmv1/products/$PRODUCT/" | grep -v product.yaml); do
  # Convert PascalCase YAML filename to snake_case (e.g. BackendBucket.yaml -> backend_bucket)
  resource=$(basename "$f" .yaml | python3 -c "
import sys, re
s = sys.stdin.read().strip()
s = re.sub(r'(?<=[a-z0-9])(?=[A-Z])', '_', s).lower()
print(s)
")
  for suffix in '' '_generated_test'; do
    target="$PROVIDER/google/services/$PRODUCT/list_${resource}${suffix}.go"
    if [ ! -f "$target" ]; then
      echo "MISSING: $target"
    fi
  done
done
```

Any `MISSING` line = FAIL. Match against oracle patterns before reporting:
- Missing `.go` entirely → **P-01** (wrong key), **P-08** (bare array), **P-09** (identity mismatch)
- File present but empty → template rendering issue, report raw `make provider` output

---

## Step 6 — Run acceptance tests (Check 5)

Follow `.agents/skills/utils/run-acctests/SKILL.md` for each resource that gained
`generate_list_resource: true`. Run from the downstream provider root.

### 6a — Run the generated list-query test for each resource

For each opted-in resource (snake_case name derived in Step 5):

```bash
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
TF_LOG=DEBUG make testacc \
  TEST=./google/services/<product> \
  TESTARGS='-run=TestAcc<ResourcePascalCase>ListQuery_generated$' \
  > /tmp/test_<resource>.log 2>&1
echo "exit: $?"
```

- If exit 0: test passed, continue to the next resource.
- If non-zero: the generated test failed. Follow Step 6b immediately before testing the next resource.

### 6b — Parse failures and fix

When a generated test fails, follow `.agents/skills/utils/parse-debug-logs/SKILL.md`:

```bash
python3 .agents/scripts/tf_debug_parser.py /tmp/test_<resource>.log \
  --extract-dir /tmp/debug_<resource>
```

Read the generated `outline.txt` and any referenced payload JSON files. Cross-reference against the
oracle patterns:

| Error pattern in outline | Oracle pattern |
|--------------------------|----------------|
| List returns 0 items / wrong key in response | **P-01** |
| Region/zone mismatch between resource creation and list query | **P-04**, **P-05**, **P-06** |
| `firewall_policy` value rejected | **P-07** |
| API returns bare array | **P-08** |
| Identity field name mismatch | **P-09** |

**If the fix is in the YAML** (wrong `collection_url_key`, `vars` vs `resource_id_vars`, etc.): the
Coder must fix it. Return FAIL with the exact oracle pattern ID and failure details in the feedback.

**If the fix is only in the sample `.tf.tmpl` or test scaffolding** (e.g. a hardcoded region in the
sample that does not affect the YAML `generate_list_resource: true` addition): fix it yourself,
regenerate (`make provider … PRODUCT=<product>`), re-run the test to confirm it passes, then amend
the commit and force-push.

### 6c — Write handwritten supplementary tests (when needed)

After the generated test passes, check whether a **handwritten** acceptance test for the list data
source already exists:

```bash
ls $GOPATH/src/github.com/hashicorp/terraform-provider-google/google/services/<product>/list_<resource>_test.go 2>/dev/null
```

Write a handwritten test **only when** at least one of these is true:
- The generated test only covers the happy path but the data source has filter parameters that need
  separate coverage.
- The resource's list URL has a non-standard scope param beyond `project`/`region`/`zone` that the
  generated test does not exercise.
- The `validate-provider-changes` oracle (Step 4) flagged a missing acceptance test for the new
  list data source.

Handwritten test file convention:
- Path: `google/services/<product>/list_<resource>_test.go`
- Package: `<product>_test`
- Test function name: `TestAcc<ResourcePascalCase>ListQuery_<scenario>` (e.g. `_withFilter`,
  `_withProject`)
- Must include `resource.ParallelTest(t, ...)` and the standard `providertest.AccTestPreCheck(t)`
  preamble.

After writing, run the handwritten test to confirm it passes:

```bash
TF_LOG=DEBUG make testacc \
  TEST=./google/services/<product> \
  TESTARGS='-run=TestAcc<ResourcePascalCase>ListQuery_<scenario>$' \
  > /tmp/test_<resource>_handwritten.log 2>&1
echo "exit: $?"
```

Parse any failure with `parse-debug-logs` before giving up.

### 6d — Commit test files and force-push

After all generated and handwritten tests pass:

1. Stage only the test files (do NOT stage generated `list_*.go` provider files — those stay
   downstream and are not committed to magic-modules):

```bash
cd <magic-modules-root>
git add mmv1/products/<product>/    # YAML changes only — already staged from Coder
# If you wrote or modified sample .tf.tmpl files to fix P-04/P-05/P-06:
git add mmv1/templates/terraform/examples/<product>_<resource>*.tf.tmpl
# Amend the Coder's commit to include sample fixes:
git commit --amend --no-edit
git push --force-with-lease $FORK_REMOTE $BRANCH
```

> **Note:** Handwritten test `.go` files live in the downstream provider repo
> (`$GOPATH/src/github.com/hashicorp/terraform-provider-google`), not in magic-modules. Do not
> commit them to the magic-modules branch. The magic-modules PR generates them during CI.

---

## Step 7 — Open the PR (only if ALL checks passed)

Follow `.agents/skills/operations/create-pr/SKILL.md` exactly.

Key requirements:
- Write the PR body to `/tmp/pr_body.txt` using a single-quoted HEREDOC. Never pass backticks inline
  to `--body` — the create-pr SKILL explains why this silently strips release-note blocks.
- One `release-note:new-list-resource` block per opted-in resource.
- Open against `GoogleCloudPlatform/magic-modules:main` from `$FORK_REMOTE:$BRANCH`.
- After `gh pr create`, run `gh pr view` to verify the release-note block rendered.
- Post `@modular-magician reassign-reviewer` as a PR comment.

---

## Response format

Respond ONLY with a raw JSON object — no markdown fences, no prose before or after.

All checks passed and PR opened:
```
{"status":"PASS","feedback":"","pr_url":"https://github.com/GoogleCloudPlatform/magic-modules/pull/<N>"}
```

Any check failed (do NOT open a PR in this case):
```
{"status":"FAIL","feedback":"P-NN: <check name>:\n<exact command>\n<full output>","pr_url":""}
```

Always include the oracle pattern ID in the feedback when one matches — it lets the Coder look up
the exact fix immediately without having to reason from the raw error.
