---
name: add-list-resource-workflow
description: "Opt resource into MMv1 list-resource generation by setting `generate_list_resource: true`, validate it locally, and open a one-resource PR against GoogleCloudPlatform/magic-modules. Invoke when the user asks to add list-resource support for a specific MMv1 resource or to enable `generate_list_resource` for an eligible resource."
---

# `add-list-resource-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

This workflow produces a single PR scoped to **one product** that flips `generate_list_resource: true` on every eligible MMv1 resource in that product, generates the downstream code, runs the generated list-query tests, and opens the PR. Do **one product per PR**, with as many eligible resources as pass.

Consult `.agents/knowledge/index.md` for the topics this task touches and open the relevant sources.

## Prerequisites

* You are in the `magic-modules` root directory.
* `$GOPATH` is set and `terraform-provider-google` is checked out at `$GOPATH/src/github.com/hashicorp/terraform-provider-google` (or another known path — confirm with the user).
* The magic-modules `mau` remote points at the user's personal fork (`git remote get-url mau`). Confirm with the user if it is missing.
* `gh` CLI is authenticated as the user (`gh auth status`).
* The user has named **one** target product (e.g. `compute`). Resources will be selected automatically by the eligibility scan.

## Eligibility Check

A resource is eligible when **the generated list-query test can run unattended**. In practice that means:

1. The resource must not be excluded from identity or read generation (`exclude_identity_generation: true` or `exclude_read: true`) — the generator hard-fails on these.
2. The first example must not have `exclude_test: true`, since the generated query test reuses its config.
3. Every scope parameter in the list URL (path params other than `project`/`region`/`zone`/`location` — e.g. `disk`, `instance`, `parent`) needs special handling, because the query test template currently only auto-injects `project`/`region`/`zone` into the test context map. **Resources whose list URL contains other scope params are not eligible until the template is extended.**

Required *body* fields (set at create time) do **not** affect list eligibility. Only path scope params on the list/collection URL matter.

Run the eligibility scan across the whole product and produce a candidate list. Report it back to the user before editing any YAML.

```bash
PRODUCT=<product>   # e.g. compute

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
    # Inspect list URL scope params. Prefer collection_url_key construction over base_url,
    # but base_url is a reliable proxy for the list collection URL.
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

Stop and present the candidate list to the user before making any edits. Do not attempt to remove `required: true` from properties or to remove `exclude_identity_generation` to force eligibility.

## Execution Steps

### 1. Sync and branch

```bash
git fetch upstream main
BRANCH="add-${PRODUCT}-list-resources"   # e.g. add-compute-list-resources
git checkout -b "$BRANCH" upstream/main
```

If the working tree is dirty, stash before checkout and warn the user.

### 2. Edit each eligible resource's YAML

For **every** resource the user approved from the eligibility scan, insert `generate_list_resource: true` as a top-level key in its YAML. Place it adjacent to other top-level booleans such as `immutable:` or `has_self_link:` for readability. Do not touch any other fields.

```bash
# Example placement (manual edit, repeated per resource)
# ...
# has_self_link: true
# immutable: true
# generate_list_resource: true
# timeouts:
# ...
```

### 3. Generate the downstream provider

```bash
PROVIDER_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google"

# Stop if downstream has uncommitted work
( cd "$PROVIDER_PATH" && git status --porcelain ) && echo "Confirm clean before continuing"

make provider VERSION=ga OUTPUT_PATH="$PROVIDER_PATH" PRODUCT=<PRODUCT>
```

Expected new files in the downstream per opted-in resource:
* `google/services/<product>/list_<resource>.go`
* `google/services/<product>/list_<resource>_generated_test.go`
* `website/docs/list-resources/<terraform_name>.html.markdown`

### 4. Build and test

```bash
cd "$PROVIDER_PATH"
go build ./...

# Run every generated list-query test for the product in one go.
# Test name format: TestAcc<ResourceName>ListQuery_generated
TF_ACC=1 go test -v -timeout 120m \
  ./google/services/<product> \
  -run 'ListQuery_generated$' | tee /tmp/list_query_test.out
```

Required environment for a live run: `GOOGLE_PROJECT`, `GOOGLE_REGION`, `GOOGLE_ZONE`, `GOOGLE_CREDENTIALS` (or ADC). Confirm with the user before consuming GCP resources.

If **any** test fails, do not patch the generator or the YAML to suppress the failure. Report the failing resources to the user. The user decides whether to (a) drop those resources from this PR and re-generate, or (b) abort the PR entirely. Never silently ship a PR that has failing list-query tests.

### 5. Open the PR

Stage only the YAML changes in magic-modules; the downstream provider edits are throwaway artifacts and must not be committed in magic-modules.

```bash
cd <magic-modules-root>
git add mmv1/products/<PRODUCT>/
git commit -m "<product>: add list resources for <N> resources"
git push mau "$BRANCH"
```

Write the PR body to `/tmp/pr_body.md` before invoking `gh pr create`. Include one `release-note:new-list-resource` block per resource and the trimmed `--- PASS:` lines from the test output.

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

Open the PR from the `mau` fork against `GoogleCloudPlatform/magic-modules:main`:

```bash
gh pr create \
  --repo GoogleCloudPlatform/magic-modules \
  --base main \
  --head "$(gh api user -q .login):$BRANCH" \
  --title "<product>: add list resources" \
  --body-file /tmp/pr_body.md
```

### 6. Request a reviewer

After the PR is created, post a comment so `modular-magician` assigns a human reviewer:

```bash
gh pr comment <PR_NUMBER> \
  --repo GoogleCloudPlatform/magic-modules \
  --body "@modular-magician reassign-reviewer"
```

## Handoff & Guardrails

* **One product per PR.** Bundle every eligible resource in the product into a single PR. Do not split a product across multiple PRs unless the user explicitly asks.
* **Never edit the generator** (`mmv1/api/`, `mmv1/provider/`, `mmv1/templates/terraform/list_resource*`) from this workflow. If the generator misbehaves, stop and escalate to the user.
* **Never commit downstream provider files** to the magic-modules branch.
* On any failure during generate/build, abort and report the exact failing command and output. On test failures, drop the failing resource(s) from the PR (or abort) per the user's choice — never ship a PR with failing list-query tests.
