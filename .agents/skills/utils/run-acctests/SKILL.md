---
name: run-acctests
description: "Executes acceptance tests (testacc) for a specific resource or suite in the GA or beta provider and outputs verbose debug logs."
---

# `run-acctests`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Prerequisites
* Know the provider (`beta` or `ga`), service (e.g. `compute`), and test name (e.g. `TestAccComputeInstance_basic`).
* The provider must be generated from your current changes (`generate-provider` with the same `VERSION`).
* PR CI only runs beta acceptance tests, so run GA locally whenever GA behavior matters.

## Execution Steps

### 1. Select the provider

```bash
VERSION="beta" # or "ga"
if [ "$VERSION" = "ga" ]; then
  REPO=terraform-provider-google; PKG_DIR=google
else
  REPO=terraform-provider-google-beta; PKG_DIR=google-beta
fi
cd "${GOPATH:-$HOME/go}/src/github.com/hashicorp/$REPO" && pwd
```

### 2. Run the test
Stream `TF_LOG=DEBUG` output to a log file so failures can be parsed.

```bash
TF_LOG=DEBUG make testacc TEST=./$PKG_DIR/services/<SERVICE_NAME> TESTARGS='-run=<TEST_NAME>$$' > test_output_$VERSION.log 2>&1
```

`no tests to run` on GA usually means the test isn't generated for GA (still guarded or `min_version: beta`). That is a finding, not a pass.

### 3. Verification & Handoff
* If the test succeeds, return to your primary workflow.
* If the test fails, do **NOT** attempt a blind fix immediately. You MUST invoke the `parse-debug-logs` skill on the log file to understand the API failure before proposing a fix.
