---
name: tgc-run-unit-tests-skill
description: Run unit tests for TGC. Use when you need to run unit tests for TGC.
---

# tgc-run-unit-tests-skill

When you need to run unit tests for TGC, use this skill.

## When to Use This Skill

- Use this when running unit tests for TGC.
- This is helpful when you need to validate conversion logic for resources you have just added or modified.

---

## How to Use It

If you added or modified a resource, you must run its corresponding unit tests.

1. Run the test, redirecting both standard output and standard error to a log file:
   ```bash
   make test-unit-local TEST=./pkg/...
   ```