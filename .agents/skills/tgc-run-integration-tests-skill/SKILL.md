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

1. Ensure the debug output directory exists in the downstream path:
   ```bash
   mkdir -p <terraform-google-conversion-path>/debug_output/raw_logs
   ```

2. Set the following environment variable in the TGC repository (or prepend it directly to your test command):
   ```bash
   export WRITE_FILES=true
   ```

3. Identify the target test name and its specific service directory. 
   - *Example*: Target `TestAccAlloydbBackup` located in `./test/services/alloydb`.

4. Run the test, redirecting both standard output and standard error to a log file:
   ```bash
   make test-integration-local TESTPATH=./test/services/alloydb TESTARGS='-run=TestAccAlloydbBackup' > debug_output/raw_logs/alloydbBackup.log 2>&1
   ```

> **Note**: Every time you run an integration test, save the logs to a unique file so you don't overwrite the output of previous runs.