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

These are mistakes that commonly slip past the procedure in `docs/content/develop/promote-to-ga.md`. This
is not a checklist to apply blindly. Some items won't apply to your promotion, and partial promotions
deliberately keep some beta references.

## CI won't catch GA problems

- **Presubmit acceptance tests (VCR) only run against `google-beta`.** The beta schema contains
  everything GA has plus more. A test config that uses something missing from GA therefore still passes
  in CI.
- **VCR may not run at all.** It skips without comment when the generated beta code has no Go changes,
  which is common for a promotion.
- **The missing-test and missing-doc checks** only look at the beta diff.
- **The GA build-and-unit-test job** doesn't block merging.
- **Consequence:** the first real GA test run is the nightly test run after merge, which is where most
  promotion breakages show up. Run the GA tests yourself before opening the PR.

## The GA API can differ from the beta API

- **Check every promoted field against the GA API version,** including nested sub-fields and enum values,
  not just the top-level resource. Individual fields and enum values are often still beta-only even when
  the resource is GA. A field that the GA API ignores causes a permanent diff. A field it rejects causes a
  400 error.
- **The GA API may rename a field.** Renaming the Terraform field to match is a breaking change for users.
  Keep the Terraform name and map it to the new API name, using `api_name` or an encoder/decoder.
- **Don't promote before the GA API is actually serving publicly.** Otherwise the promotion has to be
  reverted. Reverting a resource that hasn't been released yet still trips the breaking-change check.
- **Missing or outdated public docs don't prove a feature isn't GA.** Removing a field from the `google`
  provider because of lagging docs is a breaking change for users.

## product.yaml and hardcoded API versions

- **Add a `ga` entry to `versions`, but never delete the `beta` entry.** When a provider version has no
  entry of its own, it falls back to the closest one. Deleting `beta` would silently point `google-beta`
  at the GA URL. It would also make future beta-only fields impossible. Products with only a `beta` entry
  aren't generated into `google` at all.
- **API versions also hide outside `product.yaml`.** Look for these and update them:
  - a resource-level `base_url` or `self_link` that contains the beta version
  - literal `/beta/` URLs in custom code (use `transport_tpg.BaseUrl` instead)
  - `references.api` links to beta documentation
  - sample names containing "beta"
  - `required_providers { google = { source = "hashicorp/google-beta" } }` blocks in sample configs
- **Promoting a whole product also means adding the service to the GA TeamCity list**
  (`.teamcity/components/inputs/services_ga.kt`). The PR check for this only runs when a `product.yaml` is
  newly added, so it won't remind you.

## Test configs

- **A GA test can't use anything that is still beta-only.** Check the whole config, not just the
  promoted resource. That covers other resources, fields on other resources, data sources and provider
  arguments.
  - `google_project_service_identity` is a frequent culprit.
  - Your options are to promote the dependency first, restructure the test to avoid it, or keep that
    sample beta-only (`min_version: beta`) with a comment explaining why.
  - Whichever you choose, make sure some GA test still covers the promoted fields.
- **Delete `provider = google-beta` lines; don't replace them with `provider = google`.** Resources use the
  default provider when `provider` is absent, and reviewers ask for the line to be removed. Check `data`
  blocks too.
- **Leftover beta settings cause permanent diffs once the API is GA.** Examples:
  - Cloud Run's `launch_stage = "BETA"` and its launch-stage annotation
  - `compute/beta/` URLs hardcoded in test assertions
- **Hardcoded resource names can collide** once the test runs in both the GA and beta nightly runs. Use
  `random_suffix`.

## Handwritten code and templates

- **Version guards (`{{- if ne $.TargetVersionName "ga" }}`) also hide outside the resource's own files.**
  Check:
  - guards wrapping an entire test file, which keep every test in it out of GA
  - shared test helpers in `services/<service>/bootstrap_test_utils.go` (or `.go.tmpl`)
  - sweepers
  - handwritten resources and data sources. These register themselves through `registry.Schema` in
    their own file, so a guard around the whole file also hides the registration.
  - handwritten data source docs
  - TGC converters in `mmv1/third_party/tgc/resource_converters.go.tmpl`
  - allowlists of field names, such as container's `addonsConfigKeys` or ForceSendFields
- **Delete `{{ else }}` branches that only exist for GA.** Don't just remove the surrounding `if`, or the
  GA-only code will also run in beta.
- **Some schema helpers are shared by sibling resources.** For example, compute instance,
  instance_template, region_instance_template and instance_from_template share code. Promote the field
  in every resource that uses the helper, not just one.
- **Setting a still-beta field in Read or a flattener (`d.Set`) compiles in GA but panics at runtime.**
  Watch for this in partial promotions, where a parent block is promoted but a child field stays beta.
- **Once a `.go.tmpl` file has no guards left, rename it to `.go` and run `gofmt`.** Template files
  aren't checked for formatting or unused imports, so problems surface after the rename. Also convert
  escaped template literals like `{{"{{"}}...{{"}}"}}` back to plain text.

## Release notes and CI bots

- Past promotion PRs use inconsistent release-note formats. Follow the format in the workflow's PR
  section, which matches `docs/content/code-review/release-notes.md`.
- To the CI bots a promotion looks like new code, so expect these notices:
  - missing service labels on the newly-GA resources
  - a request for the `override-multiple-resources` label when promoting several resources
  - a "Manual Verification Required (GA-only additions)" note, which you answer with your GA test results

## Quick checks

These searches often help. Use judgment, since partial promotions legitimately keep some matches.

```bash
# Leftover beta references in the generated GA code for the service:
git -C "$TPG" grep -nE 'provider\s*=\s*google-beta|ProtoV5ProviderBetaFactories|launch.stage.*BETA|/beta/' -- google/services/<service>
# Version guards remaining in the handwritten sources:
git grep -n 'TargetVersionName' -- mmv1/third_party/terraform/services/<service>
# Tests that are new to GA (CI won't run these, so run them live):
git -C "$TPG" diff | grep -E '^\+func TestAcc'
```
