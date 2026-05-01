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

**Reference**: `.agents/skills/tgc-add-new-generated-resource-skill/SKILL.md`

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

Run this command to check for missing tests compared to metadata:
```bash
.agents/skills/tgc-run-integration-tests-skill/scripts/check_missing_tests.sh <ResourceName> <GeneratedTestFilePath>
```

**CAUTION**: Verify that **not all** of the tests are skipped (e.g., check for `[no tests to run]` or full `SKIP` in the test output).



### Step 5: Fix Integration Tests (If Failed)
If the integration tests fail, analyze the logs generated in Step 4 and apply the fixes specified in the `tgc-fix-integration-tests-skill` playbook.

**Reference**: `.agents/skills/tgc-fix-integration-tests-skill/SKILL.md`
Refer to the **Troubleshooting Playbook** and **Examples** in that file for common solutions (e.g., handling missing requested fields with decoders).

### Step 6: Start Over
After applying any fix in Step 5 (whether in a YAML file, a Go template, or a decoder), you **MUST start over from Step 2**.
1. Return to **Step 2 (Build TGC)** to compile your fixes into the binary.
2. Proceed to **Step 4 (Run Integration Tests)** to verify if the test now passes.

Repeat this `Build -> Test -> Fix` loop until all integration tests for the resource pass with exit code `0`. Once they pass, commit your changes.

---

## Workflow Checklist Template
When starting to add or fix a resource, copy this template into your `task.md` file to track progress:

```markdown
- [ ] Step 1: Add/Modify Resource in MMv1 <!-- id: 1 -->
- [/] Step 2: Build TGC binary <!-- id: 2 -->
- [ ] Step 3: Run Unit Tests <!-- id: 3 -->
- [ ] Step 4: Run Integration Tests (with WRITE_FILES=true) <!-- id: 4 -->
  - [ ] Verify generated tests exist in `test/services/<service>/` (Cite file name) <!-- id: 7 -->
  - [ ] Verify the added resource has CAI asset data in `tests_metadata_*.json` files in the `test` directory (Cite file and date) <!-- id: 8 -->
  - [ ] Verify if any tests for the added resource in `tests_metadata_*.json` are missing in the generated test file (refer to Case 16 in `tgc-fix-integration-tests-skill/troubleshooting_playbook.md` if missing due to excluded examples) <!-- id: 9 -->
  - [ ] Verify that **not all** of the generated tests were skipped or reported "no tests to run" in the output (refer to Case 11 or Case 16 in `tgc-fix-integration-tests-skill/troubleshooting_playbook.md` if missing)<!-- id: 10 -->
- [ ] Step 5: Fix failures & restart from Step 2 <!-- id: 5 -->
- [ ] Step 6: Commit changes after green tests <!-- id: 6 -->
```

## Final Status Reporting
When completing a task, always:
1. **Show the final state of this checklist** in your final response or notification to the user.
2. Use `[s]` or `[N/A]` for steps that were skipped (e.g., Step 5 if tests passed without fixes), to distinguish from uncompleted steps `[ ]`.
3. Summarize the results in `walkthrough.md` as shown in the checklist.


## Critical Rules
- **DO NOT** run integration tests after a fix without rebuilding the TGC binary first (Step 2).
- **ALWAYS** set `WRITE_FILES=true` when running integration tests to generate CAI asset data for verification.
- **DO NOT** manually modify `tests_metadata_*.json` files in downstream TGC repo. These are read-only assets from GCS.
- **DO NOT** remove `exclude_test: true` from existing examples in resource YAML files to make integration tests pass. Instead, add a new example or use `tgc_tests` if appropriate.
