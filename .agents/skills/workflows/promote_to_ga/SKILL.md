---
name: promote-to-ga-workflow
description: "Workflow for promoting beta-only resources, fields, data sources, or whole products to the GA (`google`) provider."
---

# `promote-to-ga-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.

Before editing, read two things:
- `docs/content/develop/promote-to-ga.md`, the step-by-step procedure.
- `.agents/knowledge/promotion/ga-promotion-pitfalls.md`, the common mistakes the procedure doesn't cover.

Promotions vary widely: MMv1 or handwritten, a whole product or a single nested field. Treat this as
guidance and use judgment for your case.

## Starting point

You should be told what to promote, and that it is generally available (GA) in the underlying API.

- The "GA API" is whatever API version the product's `ga` entry in `product.yaml` points to. This is often
  `v1` (versus `v1beta1` for beta), but not always. For example, Cloud Run serves beta features from the
  same API version and gates them with a launch-stage setting. Work out what the `google` provider will
  actually call before judging availability.
- If nobody has established GA status, or the live API contradicts it, stop and report what you found.
  Contradictions include the GA endpoint rejecting a field, leaving it out of responses, or returning 404
  for the resource. Public documentation often lags the API, so missing docs don't prove anything either
  way.

## What "done" means

1. Everything in scope is generated into `terraform-provider-google` and works there. Anything out of
   scope stays beta-only through `min_version: beta` or version guards.
2. The `google-beta` provider behaves the same as before, unless you intended a change. In particular, it
   must not silently switch to a different API version.
3. Every promoted field is exercised by at least one test that runs in the `google` provider. Those tests
   pass live against both `google` and `google-beta`.
4. The PR has correct release notes and reports the GA test results.

## Approach

- **Inventory before editing.** Search the repo for the resource and field names, and for
  `TargetVersionName` version guards and `min_version` settings near them. Promotions most often miss
  files that don't contain the resource's name:
  - shared test helpers
  - schema helpers shared with sibling resources
  - sweepers
  - data sources and IAM resources
  - test configs of *other* resources that use this one

  Build your edit list from this inventory rather than from the doc's checklist alone.
- **Classify each item** as one of three:
  - **Promote.**
  - **Stays beta.** The item is intentionally left out of the promotion.
  - **Blocked.** Its test depends on something that is still beta-only, or the GA API doesn't support it.

  For blocked items, lay out the options to the user instead of silently picking one. The pitfalls entry
  lists the usual options.
- **Cleanup scope.** If you find existing GA problems in the promoted resource's own files, such as a test
  still wrapped in a beta-only guard, fix them. Problems elsewhere should be reported, not fixed.

## Verification

PR CI runs acceptance tests (VCR) only against the beta provider. If the beta output didn't change, which
is common for promotions, VCR may skip entirely without saying so. Nobody else will run the GA tests
before merge, so you must.

- Run [`run-pre-gen-checks`](../../utils/run-pre-gen-checks/SKILL.md).
- Use [`repo-sync`](../../operations/repo-sync/SKILL.md) to confirm both downstream repos are in sync.
- Use [`generate-provider`](../../operations/generate-provider/SKILL.md) to generate **both** providers
  (`VERSION=ga` and `VERSION=beta`), and build both.
- Review both downstream diffs:
  - The GA diff should look like the promoted resource or fields arriving: new resource files, schema
    entries, tests and docs.
  - The beta diff should be empty or easy to explain. An unexpected beta diff means you changed beta
    behavior.
- In the GA repo, run `make test-compile TEST=./google/services/<service>`. `make build` doesn't compile
  test files, and problems in `.go.tmpl` sources only appear once they're generated.
- Run every acceptance test that touches the promoted surface against GA, then against beta. You can use
  [`run-acctests`](../../utils/run-acctests/SKILL.md) (`VERSION=ga` or `beta`), the `qa-test-runner`
  subagent, or the `test-fixer` subagent with `target_provider: both`.
- If a test can't run (missing credentials, org-level resources, quota), tell the user and say so in the
  PR. Don't imply coverage you didn't get.

## PR

- Write release notes following `docs/content/code-review/release-notes.md`:
  - For resources and data sources, add one `new-resource` or `new-datasource` block each, with the body
    `` `google_x` (ga) ``. Each IAM resource gets its own block.
  - For fields, use `enhancement` with `` product: added `field` field to `google_x` resource (ga) ``.
- In the description, cover:
  - what was promoted
  - what intentionally stayed beta, and why
  - the GA and beta test results
- Don't put internal tracker IDs, internal links or unannounced launch dates in commits, the PR, code
  comments or release notes.
- Push the branch and create the PR link with [`prepare-pr-link`](../../operations/prepare-pr-link/SKILL.md).

If you hit a problem that isn't covered in the pitfalls entry, mention it in your final report so the
entry can be updated.
