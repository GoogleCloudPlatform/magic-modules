# Troubleshooting Reference

Failure classes agents hit during Magic Modules tasks, mapped to fix strategies. Consult this after
`parse-debug-logs` or a `qa-test-runner` report to classify the failure before invoking the `fix` skill.

Each entry: **Symptom** (what the logs actually say) → **Cause** → **Fix** → **Do NOT** (the shortcut that
makes the failure disappear without fixing anything — these shortcuts are checked for and will fail review).

---

## 1. Generated provider fails to compile after a schema change

**Symptom:** `go build` errors in the generated downstream referencing the resource you touched — undefined
symbols, type mismatches (`cannot use ... as ...`), or a template rendering error during generation itself.

**Cause:** the MMv1 YAML change produced Go the rest of the file can't consume: a field type that doesn't
match custom code expectations, a renamed field still referenced by `custom_code` hooks, or a missing
required property on a new nested object.

**Fix:** read the first compile error only (the rest cascade). Diff the generated file against the previous
generation to see exactly what the YAML change did. Check whether the resource has `custom_code:` entries
(custom expanders/flatteners in `mmv1/templates/`) that reference the old shape, and update them together
with the schema.

**Do NOT:** delete or comment out the custom code to get a compile; it exists for a reason (find it via
`git log` on the template file).

## 2. Test fails with a permanent diff (permadiff)

**Symptom:** acceptance test fails at the plan-check step after apply: "After applying this test step, the
plan was not empty" — the same field shows a diff on every plan.

**Cause:** the API returns a normalized or server-computed form of what was sent (case changes, unit
conversion, defaults filled in, list reordering), so state never matches config.

**Fix:** identify what the API actually returns (the debug log's HTTP response body). Legitimate fixes, in
preference order: model the field correctly (e.g. mark it `output: true` if it is genuinely
server-controlled), add `default_from_api: true` **when the API fills a default the user didn't set**, or a
`diff_suppress_func` that suppresses only the specific normalization. Each of these requires a
`justified:` comment explaining the API behavior (see the baseline-guardrails check).

**Do NOT:** add `ignore_read` to make the diff vanish — that stops the provider reading the field at all
and hides real drift. `ignore_read` is for fields the API genuinely never returns (e.g. secrets), nothing
else.

## 3. Test fails at create with 403 / permission / quota errors

**Symptom:** the create call returns 403, `PERMISSION_DENIED`, `quotaExceeded`, or an org-policy violation
before any resource logic runs.

**Cause:** environment, not code — the test project lacks an API enablement, IAM role, or quota, or the
resource needs an org-level feature.

**Fix:** classify as **environment** and report it; do not iterate on schema or test code. Note the exact
missing permission/API from the error body in your report. If the test previously passed in nightlies, say
so (that distinguishes a broken environment from a test that never worked).

**Do NOT:** mark the test skipped to move on. Environment failures are surfaced, not silenced.

## 4. Test fails at import verification

**Symptom:** `ImportStateVerify` step fails: attributes differ between the imported state and the original.

**Cause:** a field doesn't survive the read path — the API doesn't return it, returns it under a different
name/shape, or the import ID doesn't carry enough to reconstruct it.

**Fix:** check the read/flatten path for the field (is it actually read?). If the API genuinely never
returns the value (write-only/secret fields), `ImportStateVerifyIgnore` with a `justified:` comment naming
the API behavior is the correct, documented escape hatch.

**Do NOT:** add fields to `ImportStateVerifyIgnore` because the values "look equivalent" — that's a
normalization bug (see permadiff) or a read-path bug, and ignoring it hides broken import for users.

## 5. VCR cassette mismatch

**Symptom:** test passes live but fails in VCR replay (or vice versa): request not found in cassette,
or a recorded response no longer matches current behavior.

**Cause:** the recorded HTTP interactions are stale relative to the code path (new fields in the request
body change the request signature), or the test has nondeterminism (random names not registered with the
VCR name normalizer).

**Fix:** if the code change legitimately changes requests, the cassette needs re-recording (nightly/CI
concern — note it in the PR body). For nondeterminism, use the framework's random-name helpers rather than
raw randomness.

**Do NOT:** gate the test with a VCR skip to avoid the mismatch; a test that only runs in one mode ages
into a test that runs in none.

## 6. "Update" test that doesn't update

**Symptom:** not a failure — the test is green. The update step re-applies a config identical to the
create step, or a mutable field has no update step at all.

**Cause:** test authored to satisfy "an update step exists" rather than "the update path is exercised."

**Fix:** the second config must differ in every mutable field under test; after update, import-and-verify
again. Server-populated (`output: true`) fields must have their values asserted, not just be present.

**Do NOT:** treat green as done. This class is invisible to CI today and is the most common human review
catch; the test-adequacy check exists specifically to find it.

---

*This reference grows via retros: when a new failure class is diagnosed in a session, add an entry in the
same Symptom/Cause/Fix/Do-NOT shape in the PR that fixes it.*
