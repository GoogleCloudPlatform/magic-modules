---
name: fix
description: "Plan a remediation strategy after a test failure."
---

# `fix`

Use this skill after the `qa-verification` subagent returns a failure report. Your job is to analyze the failure and propose a new fix strategy (Remediation Planning). Do NOT run generation or testing yourself.

## Execution Steps

### 1. Analyze Failure
Read the report returned by `qa-verification`. Compare findings against standard troubleshooting patterns (e.g. perpetual diffs, normalization issues).

### 2. Propose a New Strategy
Ask the user if they have a suggested change, or if you should propose a new strategy from reference guides. Offer specific recommendations.

### 3. Handoff
Once a solution is agreed upon, apply the change manually, and then transition back to Step 3 (Generate) to compile the fix.
