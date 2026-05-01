---
name: tgc-run-integration-tests-skill
description: Run integration tests for TGC. Use when you need to run integration tests for TGC.
---

# tgc-run-integration-tests-skill

When you need to run integration tests for TGC, use this skill.

## When to Use This Skill

- Use this when running integration tests for TGC.
- This is helpful when you need to validate conversion logic for resources you have just added or modified.

---

## How to Use It

If you added or modified a resource, you must run its corresponding integration tests.

1. Identify the target test name and its specific service directory. 
   - *Example*: Target `TestAccAlloydbBackup` located in `./test/services/alloydb`.

2. Run the test using the script from the `scripts` file, passing the test path and test name:
   ```bash
   .agents/skills/tgc-run-integration-tests-skill/scripts/run_integration_test.sh <test-path> <test-name>
   ```
   **Example**:
   ```bash
   .agents/skills/tgc-run-integration-tests-skill/scripts/run_integration_test.sh ./test/services/alloydb TestAccAlloydbBackup
   ```

> [!IMPORTANT]
> You MUST use the `./` prefix for the `<test-path>` (e.g., `./test/services/alloydb`) to ensure `go test` interprets it as a local directory rather than a standard library package.

> ### Error Reporting
   If a failure is detected or provided, you **MUST** report it using the template specified in GEMINI.md before proceeding.
