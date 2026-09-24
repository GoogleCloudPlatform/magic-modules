---
name: promote-to-ga-workflow
description: "Workflow for promoting beta-only resources, fields, data sources, or whole products to the GA (`google`) provider."
---

# `promote-to-ga-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.

The mechanics are in `docs/content/develop/promote-to-ga.md`, and past failures are in
`.agents/knowledge/promotion/ga-promotion-pitfalls.md`. Read both before editing. Promotions vary widely
(MMv1, handwritten, whole product, single nested field), so use judgment.

## Starting point

You're told what to promote and that it is GA in the API.
- "GA API" means whatever the product's `ga` version calls. That is often `v1`, but not always (e.g.
  Cloud Run gates beta features by launch stage within one version).
- If GA status isn't established, or live API behavior contradicts it (a field is rejected or omitted, or
  a request returns 404), stop and report the evidence. Lagging public docs are not evidence either way.

## Done means

1. Everything in scope generates into and works in `terraform-provider-google`. Everything out of scope
   stays beta-only.
2. Beta behavior is unchanged unless intended. In particular, beta must not silently switch API version.
3. Every promoted field is covered by a test that runs in GA, and it passes live on GA and on beta.
4. The PR has correct release notes and the GA test results.

## Approach

- **Inventory first.** Search for the resource/field names and for nearby `TargetVersionName` /
  `min_version`. Misses usually happen in files without the resource's name: shared test helpers,
  sibling resources' schema helpers, sweepers, data sources, IAM, and other resources' test configs.
- **Classify** each element as: promote, stays beta, or blocked. Blocked means a beta-only test dependency
  or missing from the GA API. For blocked items, present the options rather than silently picking one.
- **Cleanup scope:** fix pre-existing GA gaps in the promoted resource's own files. Report gaps elsewhere.

## Verification

PR CI only acceptance-tests **beta**, and VCR may skip silently when the beta output didn't change. The
GA tests are yours to run.

- Run [`run-pre-gen-checks`](../../utils/run-pre-gen-checks/SKILL.md). Then [`repo-sync`](../../operations/repo-sync/SKILL.md)
  and [`generate-provider`](../../operations/generate-provider/SKILL.md) for **both** `ga` and `beta`, and
  build both.
- Check the diffs:
  - GA should look like the promoted surface arriving.
  - Beta should be empty or explainable.
- In GA, run `make test-compile TEST=./google/services/<svc>`. `make build` skips test files.
- Run every test touching the promoted surface on GA and on beta, with
  [`run-acctests`](../../utils/run-acctests/SKILL.md) (`VERSION=ga|beta`), `qa-test-runner`, or
  `test-fixer` (`target_provider: both`).
- If a test can't run (credentials, org resources, quota), say so to the user and in the PR.

## PR

- Release notes (`docs/content/code-review/release-notes.md`):
  - resources and data sources: `new-resource` / `new-datasource` with `` `google_x` (ga) ``, one block
    each, including each IAM resource
  - fields: `enhancement` with `` product: added `field` field to `google_x` resource (ga) ``
- Description: what was promoted, what stayed beta and why, and the GA/beta test results.
- No internal tracker IDs, internal links or unannounced dates in commits, PR text, comments or release notes.
- Push with [`prepare-pr-link`](../../operations/prepare-pr-link/SKILL.md).

Report any pitfall you hit that isn't in the knowledge entry.
