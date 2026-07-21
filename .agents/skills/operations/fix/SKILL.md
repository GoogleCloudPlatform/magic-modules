---
name: fix
description: "Plan a remediation strategy after a test failure."
---

# `fix`

Use this skill after the `qa-test-runner` subagent returns a failure report, or invoke the `test-fixer` subagent (`.agents/agents/test-fixer/`) to handle diagnosis, remediation, generation, and re-testing end-to-end.

## Execution Options

### Option A: End-to-End Test Fix Subagent (Automated)
1. Use the `intake-test-failure` skill (`.agents/skills/utils/intake-test-failure/`) to parse the raw failure input (GitHub issue URL, GCS log link, or direct text) into a **Normalized Failure Payload**.
2. Pass the normalized payload to the `test-fixer` subagent (`.agents/agents/test-fixer/`) to isolate the root cause, modify Magic Modules source files, generate downstream code, build, and re-run acceptance tests to verify pass.


### Option B: Interactive Remediation Planning
1. **Analyze Failure**: Read the report returned by `qa-test-runner`. Compare findings against standard troubleshooting patterns (e.g. perpetual diffs, normalization issues).
2. **Propose Strategy**: Propose a new strategy from reference guides or user input.
3. **Handoff**: Apply the change in Magic Modules, and transition to Step 3 (Generate) to compile and verify.

