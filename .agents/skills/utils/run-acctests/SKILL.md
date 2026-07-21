---
name: run-acctests
description: "Generate downstream providers and run acceptance tests sequentially (Beta, then GA) for a specific service (and optionally, a test pattern), outputting verbose debug logs and short-circuiting on Beta failure. Use this skill when checking if downstream providers pass acceptance tests."
---

# `run-acctests`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.
> This skill generates downstream providers from local Magic Modules code into isolated scratch directories, compiles binaries, and runs acceptance tests sequentially (first Beta, then GA) with `TF_LOG=DEBUG`. If Beta acceptance tests fail, the script **short-circuits (stops immediately)** and skips the GA test run. Output logs are written to `scratch/acctest-<version>/logs/test_output_<version>.log`.

## Prerequisites

- You must be operating in the `magic-modules` root directory.
- You must know the target service name or service path (e.g. `compute` or `./google-beta/services/compute`) and optional test name pattern (e.g. `TestAccComputeInstance_basic`).

## Execution Steps

### 1. Execute Acceptance Test Runner

Run the acceptance test runner script:

```bash
# General usage:
# ./.agents/skills/utils/run-acctests/scripts/run_acctests.sh <service_name_or_path> [test_name_pattern]

# Example: Run a specific acceptance test sequentially (Beta then GA)
./.agents/skills/utils/run-acctests/scripts/run_acctests.sh compute TestAccComputeInstance_basic

# Example: Run all acceptance tests for a service sequentially (Beta then GA)
./.agents/skills/utils/run-acctests/scripts/run_acctests.sh storage
```

### 2. Verification & Handoff

* If the tests succeed, return to your primary workflow.
* If any test fails, do **NOT** attempt a blind fix immediately. You MUST invoke the `parse-debug-logs` skill on `scratch/acctest-<version>/logs/test_output_<version>.log` (where `<version>` is `beta` or `ga`) to analyze the API failure before proposing a fix.
