---
name: verify-changes-workflow
description: "Workflow to run all CI pipeline checks (pre-generation static checks, downstream provider generation & unit tests, provider change validations, and acceptance tests for modified services) with early exit on failure."
---

# `verify-changes-workflow`

This workflow runs all validation and test checks equivalent to the Magic Modules CI pipeline (`.ci` and `.github`). The workflow executes checks in three sequential phases and **stops immediately** if any step in a phase fails.

---

## Prerequisites

- You must be in the `magic-modules` root directory.
- `go` (v1.26+) and `git` must be installed and configured.

---

## Execution Steps

### Phase 1: Pre-Generation & Static Checks

These checks run directly against `magic-modules` without generating downstream providers.

#### 1. Go Formatting Check (`gofmt`)
Verify that all `.go` files in the repository follow standard formatting rules:
```bash
GOFMT_OUTPUT="$(gofmt -l .)"
if [ -n "$GOFMT_OUTPUT" ]; then
  echo "The following files are not formatted properly:" >&2
  echo "$GOFMT_OUTPUT" >&2
  exit 1
fi
```

#### 2. Template Validation Checks (`tools/template-check`)
Check for invalid version guards and unlinked new templates:
```bash
# Version guard check
(
  cd tools/template-check
  tmpls=$(git diff --name-only --diff-filter=d origin/main ../../*.tmpl | sed 's=^=../../=g')
  if [ -n "$tmpls" ]; then
    go run main.go version-guard --file-list "${tmpls//$'\n'/,}"
  fi
)

# Unused template check
(
  cd tools/template-check
  newtmplfiles=$(git diff --name-only --diff-filter=A origin/main HEAD -- ../../mmv1 | grep .tmpl | sed 's=^=../../=g')
  if [ -n "$newtmplfiles" ]; then
    go run main.go unused-tmpl --file-list "${newtmplfiles//$'\n'/,}"
  fi
)
```

#### 3. MMv1 Core Unit Tests
Run unit tests for `mmv1`:
```bash
(cd mmv1 && go test ./...)
```

#### 4. Tool Unit Tests
Run unit tests for internal tools (`go-changelog`, `issue-labeler`, `template-check`, `test-reader`):
```bash
(cd tools/go-changelog && go test ./...)
(cd tools/issue-labeler && go test ./...)
(cd tools/template-check && go test ./...)
(cd tools/test-reader && go test ./...)
```

> 🛑 **SHORT-CIRCUIT GUARD:** If any check in **Phase 1** fails, **STOP IMMEDIATELY**. Report the failure details to the user and do not proceed to Phase 2.

---

### Phase 2: Downstream Provider Generation & Unit Tests

This phase generates downstream provider code and runs provider build, unit test, and breaking change / missing test / missing doc validation checks in parallel in isolated scratch environments.

#### 1. Downstream Generation & Unit Tests
Invoke the `build-and-test-downstreams` skill ([build-and-test-downstreams/SKILL.md](../../utils/build-and-test-downstreams/SKILL.md)) to generate GA, Beta, and docs-examples downstreams and execute `go build`, `make testnolint`, `make lint`, and `make docscheck` in parallel:
```bash
./.agents/skills/utils/build-and-test-downstreams/scripts/build_and_test_downstreams.sh
```

#### 2. Provider Changes Validation (Diff Processor)
Invoke the `validate-provider-changes` skill ([validate-provider-changes/SKILL.md](../../utils/validate-provider-changes/SKILL.md)) to check for breaking changes, missing acceptance tests, and missing documentation across GA and Beta outputs:
```bash
./.agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh
```

> 🛑 **SHORT-CIRCUIT GUARD:** If `build_and_test_downstreams.sh` or `validate_provider_changes.sh` fails, **STOP IMMEDIATELY**. Report the failure details to the user and do not proceed to Phase 3.

---

### Phase 3: Acceptance Tests (Modified Services Only)

This phase executes acceptance tests for GCP services modified by the changes.

#### 1. Identify Modified Services
Identify which services were touched by checking the diff on `google-beta/services/` and `google/services/`.

#### 2. Execute Acceptance Tests
Invoke the `run-acctests` skill ([run-acctests/SKILL.md](../../utils/run-acctests/SKILL.md)) to run acceptance tests sequentially (first Beta, then GA):
```bash
# Execute acceptance tests sequentially (Beta, then GA) for modified service
./.agents/skills/utils/run-acctests/scripts/run_acctests.sh <SERVICE_NAME> [TEST_NAME]
```

> 🛑 **SHORT-CIRCUIT GUARD:** If any acceptance test fails:
> 1. **STOP IMMEDIATELY**. Do not proceed to subsequent tests or PR creation.
> 2. Invoke the `parse-debug-logs` skill on `scratch/acctest-<version>/logs/test_output_<version>.log` (where `<version>` is `beta` or `ga`) to analyze the failure.
> 3. Present the diagnostic report to the user and propose remediation steps.

---

## Summary of Handoff & Guardrails

- All three phases must complete successfully for verification to pass.
- Never bypass a failure or ignore failing checks.
- If any step fails to execute (including due to missing permissions or permission timeouts), the workflow MUST be marked as FAILED and stop immediately. Do NOT report success if any step was skipped or failed.
- If all phases pass, present a summary of executed checks and confirm that changes are verified and ready for PR creation.
