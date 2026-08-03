---
name: fix
description: "Plan a remediation strategy after a test failure."
---

# `fix`

Use this skill after the `qa-test-runner` subagent returns a failure report, or invoke the `test-fixer` subagent (`.agents/agents/test-fixer/`) to handle diagnosis, remediation, generation, and re-testing end-to-end.

## Execution Options

### Option A: End-to-End Test Fix Subagent (Automated)
1. Use the `intake-test-failure` skill (`.agents/skills/utils/intake-test-failure/`) to parse the raw failure input (GitHub issue URL, GCS log link, or direct text) into a **Normalized Failure Payload**.
2. Pass the normalized payload to the `test-fixer` subagent (`.agents/agents/test-fixer/`) to classify symptoms against `.agents/skills/utils/test-failure-decision-tree/SKILL.md` (all catalog scenarios), modify Magic Modules source files, generate downstream code, build, and re-run acceptance tests to verify pass.


### Option B: Interactive Remediation Planning
1. **Analyze Failure**: Read the report returned by `qa-test-runner`. Classify symptoms against the decision tree catalog in `.agents/skills/utils/test-failure-decision-tree/SKILL.md` (all catalog scenarios) to identify the failure scenario and root cause category.
2. **Propose Strategy**: Propose a specific remediation strategy from `.agents/skills/utils/test-failure-decision-tree/SKILL.md` or user input.
3. **Handoff**: Apply the change in Magic Modules, and transition to Step 3 (Generate) to compile and verify.

## Remediation Scope Guardrails

* **Strict Evidence-Based Scoping**: When applying fixes in Magic Modules, modify ONLY the specific field(s), resource(s), or configuration(s) directly proven by the test failure logs, backtrace, or plan diff to be causing the failure.
* **No Assumption-Based Expansion**: Do NOT modify sibling, adjacent, or similarly-named fields on assumption without empirical evidence from the test output. Keep PR diffs strictly scoped to proven root causes.
* **Breaking Change Compliance & User Notification**: Consult all files in `docs/content/breaking-changes/` (`breaking-changes.md` and `make-a-breaking-change.md`) before making any schema or behavioral modifications. If a breaking change needs to be made, explicitly state it to the user and list out the reason why, referencing the applicable policy in `docs/content/breaking-changes/`.

