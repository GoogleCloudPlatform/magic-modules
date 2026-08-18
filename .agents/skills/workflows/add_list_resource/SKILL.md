---
name: add-list-resource-workflow
description: "Opt resource into MMv1 list-resource generation by setting `generate_list_resource: true` and validate it locally with generated tests. Invoke when the user asks to add list-resource support for a specific MMv1 resource or to enable `generate_list_resource` for an eligible resource."
---

# `add-list-resource-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

This workflow produces a change scoped to **one product** that flips `generate_list_resource: true` on every eligible MMv1 resource in that product, generates the downstream code, and runs the generated list-query tests. Do **one product per change**, with as many eligible resources as pass.

Consult `.agents/knowledge/index.md` for the topics this task touches and open the relevant sources.

## Prerequisites

* You are in the `magic-modules` root directory.
* `$GOPATH` is set and `terraform-provider-google` is checked out at `$GOPATH/src/github.com/hashicorp/terraform-provider-google` (or another known path — confirm with the user).
* Remote repositories are configured (e.g. `upstream` pointing to `GoogleCloudPlatform/magic-modules` and personal fork remote such as `origin`).
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

If **any** test fails, do not patch the generator or the YAML to suppress the failure. Report the failing resources to the user. The user decides whether to (a) drop those resources and re-generate, or (b) abort entirely. Never silently ship a change that has failing list-query tests.

## Handoff & Guardrails

* **One product per change.** Bundle every eligible resource in the product into a single change. Do not split a product across multiple changes unless the user explicitly asks.
* **Never edit the generator** (`mmv1/api/`, `mmv1/provider/`, `mmv1/templates/terraform/list_resource*`) from this workflow. If the generator misbehaves, stop and escalate to the user.
* **Never commit downstream provider files** to the magic-modules branch.
* On any failure during generate/build, abort and report the exact failing command and output. On test failures, drop the failing resource(s) (or abort) per the user's choice — never ship a change with failing list-query tests.
