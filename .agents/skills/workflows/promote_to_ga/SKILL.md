---
name: promote-to-ga-workflow
description: "Workflow for promoting beta-only resources, fields, data sources, or whole products to the GA (`google`) provider."
---

# `promote-to-ga-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.

Follow `docs/content/develop/promote-to-ga.md`. This workflow only adds what that doc and PR CI won't
catch for you.

## Before you start

Confirm that what you're promoting is generally available in the API version the `google` provider will
call (the product's `ga` version). If that hasn't been established, or the live API contradicts it (for
example, the GA endpoint rejects or ignores a field), stop and report what you found.

## Verify beyond CI

PR CI generates and compiles the GA provider, but it only runs acceptance tests (VCR) against beta. VCR
may also skip entirely when the beta output didn't change, which is common for promotions. GA-only
failures otherwise first show up in the nightly tests after merge.

- Generate and build **both** providers ([`generate-provider`](../../operations/generate-provider/SKILL.md)
  with `VERSION=ga` and `VERSION=beta`).
- The beta downstream diff should be empty or explainable. An unexpected beta change means you changed
  beta behavior. For example, removing `beta` from `product.yaml` silently points `google-beta` at the GA
  URL.
- The GA code must not call beta endpoints. Tests still pass against a beta endpoint, so look for API
  versions hardcoded outside `product.yaml`, such as a resource `base_url` or URLs in custom code.
- Run every acceptance test that covers the promoted surface against GA and against beta
  ([`run-acctests`](../../utils/run-acctests/SKILL.md) with `VERSION=ga|beta`). Every promoted field must
  be exercised by at least one test that runs in GA.
- If a test can't run, tell the user and say so in the PR.

## PR

Include the GA and beta test results in the description. Keep internal tracker IDs, internal links and
unannounced dates out of commits, the PR, code comments and release notes.
