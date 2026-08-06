---
name: add-list-resource-workflow
description: "Opt resource into MMv1 list-resource generation by setting `generate_list_resource: true` and validate it locally with generated tests. Invoke when the user asks to add list-resource support for a specific MMv1 resource or to enable `generate_list_resource` for an eligible resource."
---

# `add-list-resource-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

This workflow produces a change scoped to **one product** that flips `generate_list_resource: true` on every eligible MMv1 resource in that product, generates the downstream code, and runs the generated list-query tests. Do **one product per change**, with as many eligible resources as pass.

## Step -1 — Autonomous product selection (when no PRODUCT is provided)

Run this when the orchestrator has not specified a target product. Score every
product by eligible-resource count, exclude any with an open list-resource PR
or an existing fork branch, and return the highest-scoring candidate.

### -1a. Collect exclusions

```bash
# Open PRs against upstream that mention list resources
gh pr list \
  --repo GoogleCloudPlatform/magic-modules \
  --state open \
  --json title,headRefName \
  --limit 200 \
  > /tmp/open_prs.json

# Fork branches that follow the fixed naming convention
git fetch origin 2>/dev/null || true
git branch -r | grep 'origin/add-.*-list-resources' | \
  sed 's|.*/add-\(.*\)-list-resources|\1|' \
  > /tmp/fork_branches.txt
```

### -1b. Score every product

```bash
python3 - <<'PY'
import sys, glob, yaml, os, re, json, subprocess

AUTO_SCOPES = {"project", "region", "zone", "location"}

def eligible_count(product):
    count = 0
    for f in glob.glob(f"mmv1/products/{product}/*.yaml"):
        if f.endswith("/product.yaml"):
            continue
        try:
            d = yaml.safe_load(open(f).read())
        except Exception:
            continue
        if not isinstance(d, dict):
            continue
        if d.get("exclude") or d.get("exclude_resource"):
            continue
        if d.get("exclude_identity_generation") or d.get("exclude_read"):
            continue
        if d.get("generate_list_resource"):
            continue
        ex = d.get("examples") or d.get("samples") or []
        if not ex or not isinstance(ex[0], dict):
            continue
        if ex[0].get("exclude_test"):
            continue
        list_url = d.get("base_url") or ""
        bad = [s for s in re.findall(r"{{\s*(\w+)\s*}}", list_url) if s not in AUTO_SCOPES]
        if bad:
            continue
        count += 1
    return count

pr_result = subprocess.run(
    ["gh","pr","list","--repo","GoogleCloudPlatform/magic-modules",
     "--state","open","--json","title,headRefName","--limit","200"],
    capture_output=True, text=True)
prs = json.loads(pr_result.stdout) if pr_result.returncode == 0 else []

try:
    fork_branches = open("/tmp/fork_branches.txt").read().split()
except FileNotFoundError:
    fork_branches = []

excluded = set(fork_branches)
for pr in prs:
    b = pr["headRefName"].lower()
    if "list" in pr["title"].lower() or "list" in b:
        m = re.search(r"add-(.+?)-list-resources", b)
        if m:
            excluded.add(m.group(1))

products = sorted([d for d in os.listdir("mmv1/products/") if os.path.isdir(f"mmv1/products/{d}")])
scored = [(eligible_count(p), p) for p in products if p not in excluded]
scored.sort(reverse=True)

if not scored or scored[0][0] == 0:
    print(json.dumps({"product":"","candidates":0,
        "reason":"all products already covered or have zero eligible resources"}))
else:
    n, best = scored[0]
    print(json.dumps({"product": best, "candidates": n,
        "reason": f"highest eligible-resource count ({n}) with no open PR or fork branch"}))
PY
```

### -1c. Respond

Respond **ONLY** with the raw JSON object the Python script printed.
No markdown fences. No prose before or after. No extra keys.

Example success: `{"product":"accessapproval","candidates":177,"reason":"highest eligible-resource count (3) with no open PR or fork branch"}`
Example all-done: `{"product":"","candidates":0,"reason":"all products already covered or have zero eligible resources"}`

---

## Step 0 — Read the oracle before doing anything else

```bash
cat .agents/knowledge/list-resource-oracle.md
```

The oracle is a catalog of every failure pattern encountered in previous list-resource runs — wrong
`collection_url_key` values, missing template imports, region/zone test mismatches, bare-array API
responses, and more. Reading it before you start means you can fix those issues proactively instead
of being sent back by the Validator. Do not skip this step.

Consult `.agents/knowledge/index.md` for any other topics this task touches.

## Prerequisites

* You are in the `magic-modules` root directory.
* `$GOPATH` is set and `terraform-provider-google` is checked out at `$GOPATH/src/github.com/hashicorp/terraform-provider-google` (or another known path — confirm with the user).
* Remote repositories are configured (e.g. `upstream` pointing to `GoogleCloudPlatform/magic-modules` and personal fork remote such as `origin`).
* `PRODUCT` is set — either provided by the caller or selected in Step -1. Resources will be selected automatically by the eligibility scan.

## Eligibility Check

A resource is eligible when **the generated list-query test can run unattended**. In practice that means:

1. The resource must not be excluded from identity or read generation (`exclude_identity_generation: true` or `exclude_read: true`) — the generator hard-fails on these. See oracle **P-10**.
2. The first example must not have `exclude_test: true`, since the generated query test reuses its config.
3. Every scope parameter in the list URL (path params other than `project`/`region`/`zone`/`location` — e.g. `disk`, `instance`, `parent`) needs special handling. See oracle **P-11**.

Required *body* fields (set at create time) do **not** affect list eligibility. Only path scope params on the list/collection URL matter.

Run the eligibility scan across the whole product and produce a candidate list before editing any YAML. When proceeding, print the final opted-in resource names explicitly so downstream verification can match them to generated tests.

```bash
python3 - "$PRODUCT" <<'PY'
import sys, glob, yaml, os, re
product = sys.argv[1]
AUTO_SCOPES = {"project", "region", "zone", "location"}
candidates, skipped = [], []
for f in sorted(glob.glob(f"mmv1/products/{product}/*.yaml")):
    if f.endswith("/product.yaml"):
        continue
    try:
        d = yaml.safe_load(open(f).read())
    except Exception as e:
        skipped.append((f, f"yaml parse error: {e}")); continue
    if not isinstance(d, dict):
        continue
    name = d.get("name") or os.path.basename(f)
    if d.get("exclude") or d.get("exclude_resource"):
        skipped.append((name, "excluded resource")); continue
    if d.get("exclude_identity_generation") or d.get("exclude_read"):
        skipped.append((name, "exclude_identity_generation or exclude_read")); continue
    if d.get("generate_list_resource"):
        skipped.append((name, "already opted in")); continue
    ex = d.get("examples") or d.get("samples") or []
    if not ex or not isinstance(ex[0], dict):
        skipped.append((name, "no examples/samples")); continue
    first = ex[0]
    if first.get("exclude_test"):
        skipped.append((name, "first example has exclude_test")); continue
    list_url = d.get("base_url") or ""
    scope_params = re.findall(r"{{\s*(\w+)\s*}}", list_url)
    bad_scope = [s for s in scope_params if s not in AUTO_SCOPES]
    if bad_scope:
        skipped.append((name, f"list URL has unsupported scope param(s): {bad_scope}")); continue
    candidates.append((name, f))

print("CANDIDATES:")
for n, f in candidates:
    print(f"  - {n}  ({f})")
print("\nSKIPPED:")
for n, r in skipped:
    print(f"  - {n}: {r}")
PY
```

Stop and present the candidate list before making any edits. Do not attempt to remove `required: true`
from properties or to remove `exclude_identity_generation` to force eligibility.

## Execution Steps

### 1. Sync and branch

The branch name is fixed — the same branch is reused across every run for this product. If it already
exists in the fork, check it out; otherwise create it from `upstream/main`.

```bash
git fetch upstream main
BRANCH="add-${PRODUCT}-list-resources"

if git ls-remote --exit-code origin "$BRANCH" > /dev/null 2>&1; then
  git fetch origin "$BRANCH"
  git checkout "$BRANCH"
else
  git checkout -b "$BRANCH" upstream/main
fi
```

If the working tree is dirty, stash before checkout.

### 2. Apply oracle checks to each candidate before editing YAML

For each candidate resource, consult the oracle for known issues **before** writing the YAML change:

- Does the resource's `base_url` response key match the camelized resource name? If not, you will need
  `collection_url_key`. See oracle **P-01**.
- Does the resource have any `Integer`-typed identity property? See oracle **P-02**.
- Does the sample config hardcode a `region`, `zone`, or `location` string? If so, move it to `vars`
  or `test_context_vars`. See oracle **P-04** and **P-05**.
- Is `region` or `zone` currently in `resource_id_vars`? Move it. See oracle **P-05**.
- Does the list endpoint return a bare JSON array instead of a wrapped object? See oracle **P-08**.

Resolve every applicable oracle issue in the same YAML edit. Do not add `generate_list_resource: true`
and leave known issues for the Validator to catch.

### 3. Edit each eligible resource's YAML

For **every** approved candidate, insert `generate_list_resource: true` as a top-level key, adjacent
to other top-level booleans such as `immutable:` or `has_self_link:`. Apply any `collection_url_key`,
`vars`, or `test_context_vars` fixes identified in Step 2 at the same time.

```bash
# Example — minimum addition:
# generate_list_resource: true

# Example — with collection_url_key fix (oracle P-01):
# generate_list_resource: true
# collection_url_key: 'associations'

# Example — with region moved to vars (oracle P-04/P-05):
# samples:
#   - name: my_sample
#     steps:
#       - name: my_config
#         vars:
#           region: 'us-central1'
#         resource_id_vars:
#           resource_name: 'my-resource'
```

### 4. Run the pre-generation oracle

```bash
./.agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh
```

This runs YAML lint, gofmt, template checks, and mmv1 unit tests. Fix every failure before
proceeding to generation. If a failure matches an oracle pattern, apply the documented fix — do not
work around it.

### 5. Generate the downstream provider

```bash
PROVIDER_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google"

# Verify downstream is clean before generating (oracle P-12)
( cd "$PROVIDER_PATH" && git status --porcelain ) && echo "WARNING: downstream has uncommitted work — stash it before continuing"

make provider VERSION=ga OUTPUT_PATH="$PROVIDER_PATH"
```

Expected new files in the downstream per opted-in resource:
* `google/services/<product>/list_<resource>.go`
* `google/services/<product>/list_<resource>_generated_test.go`
* `website/docs/list-resources/<terraform_name>.html.markdown`

If generation fails, cross-reference the error against the oracle before attempting a fix:
- `undefined: strconv` → **P-02**
- `imported and not used: "types"` → **P-03**
- Compile error on a list response key → **P-01**
- Any other template rendering failure → stop and escalate; do not edit the generator

### 6. Build and verify

```bash
cd "$PROVIDER_PATH"
go build ./...
```

A compile error here that matches an oracle pattern must be fixed in the YAML (or via an oracle-branch
template fix), not by editing the generated `.go` file directly.

#### Quota failures

If a test fails with a GCP quota error (e.g. `Quota 'NETWORKS' exceeded`, `RESOURCE_EXHAUSTED`,
`rateLimitExceeded`, or `quota exceeded`), **this is a GCP infrastructure constraint, not a code
bug**. Do not drop the resource from the PR or retry indefinitely.

Respond to the orchestrator with a **`QUOTA_FAIL`** status so the loop can distinguish this from
a real code failure:

```json
{"status":"QUOTA_FAIL","feedback":"<ResourceName>: Quota exceeded — <exact quota error message>"}
```

The resource's YAML and generated code are correct. The PR should still be opened. The GCP
project quota will be raised on production infrastructure before the full run.

#### Other test failures

If **any** test fails for a non-quota reason, do not patch the generator or the YAML to suppress
the failure. Report the failing resources to the user. The user decides whether to (a) drop those
resources from this PR and re-generate, or (b) abort the PR entirely. Never silently ship a PR
that has failing list-query tests. When reporting test results programmatically, return valid JSON
only with no surrounding prose or markdown fences.

### 7. Commit only the YAML changes

Stage only the files in `mmv1/products/<PRODUCT>/`. The downstream provider edits are throwaway
artifacts — never commit them to the magic-modules branch (oracle **P-13**).

```bash
cd <magic-modules-root>
git diff --cached --name-only   # confirm only mmv1/products/<PRODUCT>/ files are staged
git add mmv1/products/<PRODUCT>/
git commit -m "<product>: add list resources for <N> resources"
git push --force-with-lease origin "$BRANCH"
```

### 8. Handoff to the Validator

The Validator will independently fetch the branch, run all four oracle checks, and return a structured
JSON verdict. If it returns `FAIL`, the feedback will include an oracle pattern ID (e.g. `P-04`) and
the exact command output. Look up that pattern in `.agents/knowledge/list-resource-oracle.md`, apply
the documented fix, and push an updated commit to the same branch.

### 9. Open the PR and request review

Write the PR body to `/tmp/pr_body.md`, create the PR, print the PR URL clearly, then ask
`modular-magician` to reassign a reviewer.

```markdown
Adds list-resource generation for the following <product> resources:

- `google_<product>_<resource_a>`
- `google_<product>_<resource_b>`
- ...

```release-note:new-list-resource
`google_<product>_<resource_a>`
```

```release-note:new-list-resource
`google_<product>_<resource_b>`
```

<details><summary>Local test output</summary>

```
<paste trimmed `--- PASS: TestAcc...ListQuery_generated` lines and the final `PASS / ok` summary>
```

</details>
```

```bash
gh pr create \
  --repo GoogleCloudPlatform/magic-modules \
  --base main \
  --head "$(gh api user -q .login):$BRANCH" \
  --title "<product>: add list resources" \
  --body-file /tmp/pr_body.md

gh pr comment <PR_NUMBER> \
  --repo GoogleCloudPlatform/magic-modules \
  --body "@modular-magician reassign-reviewer"
```

## Handoff & Guardrails

* **One product per change.** Bundle every eligible resource in the product into a single change. Do not split a product across multiple changes unless the user explicitly asks.
* **Never edit the generator** (`mmv1/api/`, `mmv1/provider/`, `mmv1/templates/terraform/list_resource*`) from this workflow. If the generator misbehaves, stop and escalate to the user.
* **Never commit downstream provider files** to the magic-modules branch (oracle **P-13**).
* On any failure during generate/build, abort and report the exact failing command and output. On test failures, drop the failing resource(s) (or abort) per the user's choice — never ship a change with failing list-query tests.
* **Never remove an exclusion flag to force eligibility.** See oracle **P-10**.
* **Never weaken a check to make it pass.** If a resource cannot pass the oracle, drop it from the batch and report it to the user.
