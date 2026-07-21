---
name: run-acctests
description: "Generate downstream providers (GA, Beta), run acceptance tests for a specific service (and optionally, a specific test) in the downstream provider, and output verbose debug logs. Use this skill as a final check when the user wants to check if downstream providers pass tests."
---

# `run-acctests`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.
> This skill generates the downstream provider from local Magic Modules code into an isolated scratch directory, compiles the binary, runs the specified acceptance test with `TF_LOG=DEBUG`, and streams output to `scratch/test_output.log`.

## Prerequisites

- You must be operating in the `magic-modules` root directory.
- You must know the target version (`beta` or `ga`), service path or service name (e.g. `compute` or `./google-beta/services/compute`), and optional test name pattern (e.g. `TestAccComputeInstance_basic`).

## Execution Steps

### 1. Execute Acceptance Test Runner

Run the acceptance test runner script:

```bash
# General usage:
# ./.agents/skills/utils/run-acctests/scripts/run_acctests.sh [beta|ga] <service_name_or_path> [test_name_pattern]

# Example: Run a specific acceptance test in beta
./.agents/skills/utils/run-acctests/scripts/run_acctests.sh beta compute TestAccComputeInstance_basic

# Example: Run all acceptance tests for a service in GA
./.agents/skills/utils/run-acctests/scripts/run_acctests.sh ga storage
```

### 2. Verification & Handoff

* If the test succeeds, return to your primary workflow.
* If the test fails, do **NOT** attempt a blind fix immediately. You MUST invoke the `parse-debug-logs` skill on `scratch/test_output.log` to analyze the API failure before proposing a fix.
