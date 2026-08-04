---
name: run-acctests
description: "Executes acceptance tests (testacc) for a specific resource or suite in the provider and outputs verbose debug logs."
---

# `run-acctests`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Prerequisites
* You must be in the relevant provider root directory, e.g., `$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta`.
* You must know the specific service path (e.g., `./google-beta/services/compute`) and test name (e.g., `TestAccComputeInstance_basic`) you wish to run.

## Execution Steps

### 1. Verification

#### Verify Directory Structure
```bash
pwd # Verify we are in the generated provider directory
```

### 2. The Core Commands
Run the acceptance test with `TF_LOG=DEBUG` enabled, and stream the output to a log file (`test_output.log`) so it can be parsed later if it fails.

#### Execute Acceptance Test
```bash
# Replace <SERVICE_NAME> and <TEST_NAME> with the appropriate values
TF_LOG=DEBUG make testacc TEST=./google-beta/services/<SERVICE_NAME> TESTARGS='-run=<TEST_NAME>$$' > test_output.log 2>&1
```

### 3. Verification & Handoff
* If the test succeeds, return to your primary workflow.
* If the test fails, do **NOT** attempt a blind fix immediately. You MUST invoke the `parse-debug-logs` skill on `test_output.log` to understand the API failure before proposing a fix.
