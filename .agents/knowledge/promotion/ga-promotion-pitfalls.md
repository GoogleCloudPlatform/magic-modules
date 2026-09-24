---
name: ga-promotion-pitfalls
description: Where beta->GA promotions break (CI blind spots, API drift, beta-only test deps, missed guards) and how to avoid it.
topics: [promotion]
task_types: [ga-promotion, test-fix]
source: "authored: analysis of ~50 promotion PRs, their reverts, and GA-nightly failures (2026-09)"
status: draft
last_verified: 2026-09-24
---

# GA promotion pitfalls

Covers what goes wrong beyond `docs/content/develop/promote-to-ga.md`. It is not a checklist: partial
promotions deliberately keep beta references.

## CI blind spots
- **Presubmit VCR runs only against beta.** Beta's schema is a superset of GA's, so HCL that references
  things missing from GA still passes.
- **VCR may not run at all:** it skips silently when the beta diff has no Go changes, which is common for
  promotions.
- **Other checks are beta-only or non-blocking:** the missing-test and missing-doc detectors inspect beta
  only, and the GA unit/compile job doesn't block merge.
- So the first live GA run is the nightly after merge, which is where most promotion breakages surface.
  Run the GA tests yourself.

## GA API ≠ beta API
- **Check every promoted field, nested sub-field and enum value in the GA API**, not just the resource.
  Sibling fields and enum values are often still beta-only.
- **GA may rename fields.** Renaming the TF field is a breaking change, so map it with `api_name` or an
  encoder/decoder instead.
- **Don't promote before the GA API is publicly serving.** Early promotions have been reverted, and
  reverting an unreleased resource still trips the breaking-change check.
- **Lagging public docs don't prove non-GA.** Removing a field from GA is breaking for users.

## product.yaml and API versions
- **Add `ga`; never delete `beta`.** A missing version falls back to the closest one, so beta would
  silently start calling the GA URL. Beta-only products don't generate into GA at all.
- **Versioned URLs hide outside product.yaml:**
  - resource `base_url`/`self_link`
  - `/beta/` in custom code (use `transport_tpg.BaseUrl`)
  - `references.api` links
  - sample names containing "beta"
  - `required_providers { source = "hashicorp/google-beta" }` in samples
- **Whole-product promotions:** add the service to `.teamcity/components/inputs/services_ga.kt`. The PR
  check only fires for new product.yaml files.

## Test configs
- **No beta-only dependencies anywhere in the config:** resources, fields on other resources, data
  sources, provider arguments.
  - `google_project_service_identity` is a common one.
  - Options: promote the dependency first, restructure the test, or keep the sample `min_version: beta`
    with a comment explaining why.
  - Whichever you choose, make sure a GA test still covers the promoted fields.
- **Delete `provider = google-beta`; don't swap it for `provider = google`.** This includes `data` blocks.
- **Stale beta artifacts permadiff** once the API is GA:
  - Cloud Run `launch_stage = "BETA"` / launch-stage annotations
  - `compute/beta/` in test assertions
- **Hardcoded names collide** once both nightlies run the test. Use `random_suffix`.

## Handwritten code and templates
- **Guards hide outside the resource's files:**
  - whole-file test guards
  - `services/<svc>/bootstrap_test_utils.go[.tmpl]`
  - sweepers
  - handwritten resources and data sources: they self-register via `registry.Schema`, so a whole-file
    guard is the registration guard
  - handwritten data source docs
  - TGC `mmv1/third_party/tgc/resource_converters.go.tmpl`
  - allowlists such as `addonsConfigKeys` and ForceSendFields
- **Delete `{{ else }}` GA-only branches**; don't just unwrap them.
- **Shared schema helpers span sibling resources** (e.g. instance, instance_template, region template,
  from_template). Promote the field in all of them.
- **`d.Set` of a still-beta key compiles in GA but panics at runtime.** Watch for this in partial
  promotions.
- **A `.go.tmpl` with no guards left → rename it to `.go` and gofmt it.** Template sources hide formatting
  and import issues. Also unescape `{{"{{"}}...{{"}}"}}` literals.

## Release notes and bots
- Past PRs mix release-note formats. Follow the docs (see the workflow's PR section).
- Promotions look "new" to the bots, so expect:
  - missing service labels
  - `override-multiple-resources`
  - the "GA-only additions" manual-verification note, which you answer with your GA test results

## Quick checks (use judgment)
```bash
git -C "$TPG" grep -nE 'provider\s*=\s*google-beta|ProtoV5ProviderBetaFactories|launch.stage.*BETA|/beta/' -- google/services/<svc>
git grep -n 'TargetVersionName' -- mmv1/third_party/terraform/services/<svc>
git -C "$TPG" diff | grep -E '^\+func TestAcc'   # new GA tests: run them live
```
