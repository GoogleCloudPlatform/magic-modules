---
name: tgc-resource-addition-workflow-skill
description: End-to-end workflow for adding a new resource to TGC. Includes building, unit testing, integration testing, and fixing.
---

# tgc-resource-addition-workflow-skill

When you are tasked with adding a new resource to TGC (Terraform Google Conversion) or fixing a resource conversion, follow this end-to-end workflow. This skill glues together the various individual TGC skills into a single loop.

## The TGC Development Loop

Follow these steps sequentially. If you make a change to fix a test (Step 5), you must restart the loop from Step 2.

### Step 1: Add or Modify the Resource
Add or modify the resource definition (YAML config, templates, test data, etc.) within `magic-modules` (`mmv1/`). This is your baseline implementation.

### Step 2: Build TGC
You must rebuild the downstream generated code and compile the TGC binary from your Magic Modules changes using the `tgc-build-skill`.

**Reference**: `.agents/skills/tgc-build-skill/SKILL.md`
```bash
./.agents/skills/tgc-build-skill/scripts/build_tgc.sh
```

### Step 3: Run Unit Tests
Run the unit tests to ensure no core conversion logic was fundamentally broken. Use the `tgc-run-unit-tests-skill`.

**Reference**: `.agents/skills/tgc-run-unit-tests-skill/SKILL.md`
```bash
make test-unit-local TEST=./pkg/...
```

### Step 4: Run Integration Tests
Run the integration tests for your specific resource using the `tgc-run-integration-tests-skill`. 

**Reference**: `.agents/skills/tgc-run-integration-tests-skill/SKILL.md`
```bash
export WRITE_FILES=true
.agents/skills/tgc-run-integration-tests-skill/scripts/run_integration_test.sh <test-path> <test-name>
```

### Step 5: Fix Integration Tests (If Failed)
If the integration tests fail, analyze the logs generated in Step 4 and apply the fixes specified in the `tgc-fix-integration-tests-skill` playbook.

**Reference**: `.agents/skills/tgc-fix-integration-tests-skill/SKILL.md`

Common fixes applied here include:
- Supressing mutually exclusive fields via a custom Go `tgc_decoder`.
- Adding fields that CAI drops (e.g. default values) to the ignore lists (via `is_missing_in_cai: true` in the YAML).
- Editing `TGCTestIgnorePropertiesToStrings` in `mmv1/api/resource.go` to ignore entire blocks natively missing from CAI payload arrays (like `dynamic`).

### Step 6: Start Over
After applying any fix in Step 5 (whether in a YAML file, a Go template, or a decoder), you **MUST start over from Step 2**.
1. Return to **Step 2 (Build TGC)** to compile your fixes into the binary.
2. Proceed to **Step 4 (Run Integration Tests)** to verify if the test now passes.

Repeat this `Build -> Test -> Fix` loop until all integration tests for the resource pass with exit code `0`. Once they pass, commit your changes.
