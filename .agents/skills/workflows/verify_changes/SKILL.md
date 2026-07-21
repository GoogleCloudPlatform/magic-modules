---
name: verify-changes-workflow
description: "Workflow to run all CI pipeline checks (pre-generation static checks, downstream provider generation & unit tests, provider change validations, and acceptance tests for modified services) and short-circuit on failure."
---

# `verify-changes-workflow`

This workflow runs all validation and test checks equivalent to the Magic Modules CI pipeline (`.ci` and `.github`). The workflow executes checks in three sequential phases and **short-circuits (stops immediately)** if any step in a phase fails.

---

## Prerequisites

- You must be in the `magic-modules` root directory.
- `go` (v1.26+) and `yamllint` must be installed and configured.
- The `$GOPATH` environment variable must be set correctly.
- Downstream provider repositories must be cloned in `$GOPATH/src/github.com/`:
  - `hashicorp/terraform-provider-google`
  - `hashicorp/terraform-provider-google-beta`
  - `terraform-google-modules/docs-examples`

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

#### 2. YAML Linting (`mmv1/products`)
Lint modified YAML product definitions (or all product YAML files if full validation is desired):
```bash
yamlfiles=$(git diff --name-only --diff-filter=d origin/main -- mmv1/products)
if [ -z "$yamlfiles" ]; then
  yamlfiles=$(find mmv1/products -name "*.yaml" -o -name "*.yml" | tr '\n' ' ')
fi
if [ -n "$yamlfiles" ]; then
  yamllint -c .yamllint $yamlfiles
fi
```

#### 3. Template Validation Checks (`tools/template-check`)
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

#### 4. MMv1 Core Unit Tests
Run unit tests for `mmv1`:
```bash
(cd mmv1 && go test ./...)
```

#### 5. Tool Unit Tests
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

This phase generates downstream provider code and runs provider build, unit test, and breaking change / missing test / missing doc validation checks.

#### 1. Code Generation
Generate downstream target repositories from Magic Modules:
```bash
# Generate GA Provider
make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google"

# Generate Beta Provider
make provider VERSION=beta OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta"

# Generate Docs Examples (tf-oics)
make tf-oics OUTPUT_PATH="$GOPATH/src/github.com/terraform-google-modules/docs-examples"
```

#### 2. GA Provider Build & Unit Tests (`terraform-provider-google`)
```bash
cd "$GOPATH/src/github.com/hashicorp/terraform-provider-google"
go build
make testnolint
make lint
make docscheck
```

#### 3. Beta Provider Build & Unit Tests (`terraform-provider-google-beta`)
```bash
cd "$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta"
go build
make testnolint
make lint
make docscheck
```

#### 4. Provider Changes Validation (Diff Processor)
Execute the `validate-provider-changes` skill ([validate-provider-changes/SKILL.md](../../utils/validate-provider-changes/SKILL.md)) to check for breaking changes, missing acceptance tests, and missing documentation across GA and Beta outputs:
```bash
cd "$MAGIC_MODULES_PATH" # Return to magic-modules root
./.agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh
```

> 🛑 **SHORT-CIRCUIT GUARD:** If any generation, compilation, unit test, lint check, docs check, or `validate-provider-changes` check fails, **STOP IMMEDIATELY**. Report the failure details to the user and do not proceed to Phase 3.

---

### Phase 3: Acceptance Tests (Modified Services Only)

This phase executes acceptance tests for GCP services modified by the changes.

#### 1. Identify Modified Services
Identify which services were touched in `mmv1/products/<service>` or in downstream provider service directories (`google-beta/services/<service>`).

#### 2. Execute Acceptance Tests
Run the acceptance tests for the modified service using the `run-acctests` skill ([run-acctests/SKILL.md](../../utils/run-acctests/SKILL.md)):
```bash
cd "$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta"
TF_LOG=DEBUG make testacc TEST=./google-beta/services/<SERVICE_NAME> TESTARGS='-run=<TEST_NAME>$$' > test_output.log 2>&1
```

> 🛑 **SHORT-CIRCUIT GUARD:** If any acceptance test fails:
> 1. **STOP IMMEDIATELY**. Do not proceed to subsequent tests or PR creation.
> 2. Invoke the `parse-debug-logs` skill on `test_output.log` to analyze the failure.
> 3. Present the diagnostic report to the user and propose remediation steps.

---

## Summary of Handoff & Guardrails

- All three phases must complete successfully for verification to pass.
- Never bypass a failure or ignore failing checks.
- If all phases pass, present a summary of executed checks and confirm that changes are verified and ready for PR creation.
