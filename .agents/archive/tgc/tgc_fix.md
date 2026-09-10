# TGC Test Fix Workflow

This document defines the main loop for fixing a resource conversion or integration test failure in TGC.

## Required Skills
Before proceeding with the workflow, ensure you are familiar with and read the following skills when prompted in the phases:
- [tgc-sync-provider](skills/tgc/tgc-sync-provider/SKILL.md) (Phase 1)
- [fixing-tgc-resource-or-test-failures](skills/tgc/fixing-tgc-resource-or-test-failures/SKILL.md) (Phase 4)
- [tgc-fix-handwritten-resources-tests-skill](skills/tgc/tgc-fix-handwritten-resources-tests-skill/SKILL.md) (Phase 4 - if resource is handwritten)
- [tgc-fix-integration-tests-skill](skills/tgc/tgc-fix-integration-tests-skill/SKILL.md) (Phase 4 - if resource is generated)
- [tgc-build-skill](skills/tgc/tgc-build-skill/SKILL.md) (Phase 5)
- [tgc-run-integration-tests-skill](skills/tgc/tgc-run-integration-tests-skill/SKILL.md) (Phase 3 & Phase 6)

## The Workflow

### 1. Session Setup
- **Initialize Task List**: You MUST create a `task.md` file containing all workflow phases (Session Setup, Triage & Isolate, Run All Integration Tests, Fix, Generate Code, Unit Testing, Integration Testing, Finalization) as uncompleted tasks before performing any other actions. If the failure involves an integration test, you MUST include explicit sub-tasks in the Fix phase of `task.md` to systematically trace the failure backwards (inspecting `Test_roundtrip.tf`, `Test_roundtrip.json`, `Test_export.tf`, `Test.json`). You MUST also add a specific task item in Phase 3 to run all integration tests for the resource.
- **Set Environment**: Ensure `TGC_DIR` environment variable is set to the absolute path of your active TGC downstream workspace.
  ```bash
  export TGC_DIR=/path/to/downstream/workspace
  export PATH=/usr/local/go/bin:/opt/homebrew/bin:$PATH
  ```
- **Use Skill**: You MUST read the [tgc-sync-provider](skills/tgc/tgc-sync-provider/SKILL.md) skill. **Ask the user** which synchronization method to use (Aligning to Base Commit vs. Fast-Forward to Latest) and follow their choice to synchronize the downstream repository before proceeding to Phase 2.

### 2. Triage & Isolate (Systematic & Automated Diagnostics)
- **Generate Code**: Before running tests, you MUST generate code to ensure the downstream repository is up to date with Magic Modules.
  ```bash
  ./.agents/skills/tgc/tgc-build-skill/scripts/build_tgc.sh $TGC_DIR
  ```
- **Automated Diagnostics**: If an integration test fails, you should leverage the **Automated Diagnostic & Triage Tool** built into the test framework. Run the test with `WRITE_FILES=1` enabled:
  ```bash
  export WRITE_FILES=1
  TF_CLI_CONFIG_FILE="${PWD}/tf-dev-override.tfrc" GO111MODULE=on go test -v ./test/services/<service> -run="<TestName>"
  ```
  On failure, the test runner will automatically execute `./test/diagnose_test_failure.py` and print a structured diagnostics report identifying whether the bug is a `[cai2hcl FLATTENER BUG]` or a `[tfplan2cai EXPANDER BUG]`.
- **Manual Verification**: If automated diagnostics are unavailable, you must manually trace the data loss through `Test_roundtrip.tf`, `Test_roundtrip.json`, `Test_export.tf`, and `Test.json` to verify where the property was dropped.

### 3. Run Integration Tests
- **Run Failed Tests**: Before analyzing the failure or proposing a solution, you MUST run the failed integration tests for the affected resource to identify the specific failures.
- **Read Skill**: Read [tgc-run-integration-tests-skill](skills/tgc/tgc-run-integration-tests-skill/SKILL.md) for guidance on running integration tests.

### 4. Fix (Parent Agent)
- **Use Skill**: You MUST read the [fixing-tgc-resource-or-test-failures](skills/tgc/fixing-tgc-resource-or-test-failures/SKILL.md) skill.
- **Trace Failure**: Follow the "Tracing Failures Backwards" protocol in the Playbook to isolate the stage where data was lost before proposing a solution.
- **Report Failure**: Report the failure using the template in [TGC_WORKFLOWS.md](TGC_WORKFLOWS.md).
- **[MANDATORY] Stop and wait for user approval before applying the fix.**
- **Apply Fix**: Apply fixes in Magic Modules (`mmv1/`). DO NOT make changes directly in the downstream repository.
- **DON'T** change the schema of a resource (e.g., making a Required field Optional or a Set) to fix conversion failures, unless explicitly guided by the user.
- **Repeat Loop**: After ANY fix, you MUST repeat the full verification loop (Step 5: Generate Code & Unit Testing, Step 6: Integration Testing).

### 5. Generate Code & Unit Testing
- **Read Skill**: Read [tgc-build-skill](skills/tgc/tgc-build-skill/SKILL.md) to project changes to the downstream repository.
- **Selective Unit Testing**: The build pipeline (`build_tgc.sh`) automatically executes the selective unit test runner during code generation. If any changed unit tests fail, the build will block immediately.
- If build, unit test, or dependency errors occur, stop and immediately report the error in the conversation using the required template.

### 6. Integration Testing
- **Read Skill**: Read [tgc-run-integration-tests-skill](skills/tgc/tgc-run-integration-tests-skill/SKILL.md) for guidance on running integration tests.

### 7. Finalization
- Ask the user if the task is complete and if you should proceed with committing.
- Include a summary of any failures encountered using the template specified in [TGC_WORKFLOWS.md](TGC_WORKFLOWS.md).
- Commit changes under `mmv1/` folder only.
- Exclude scratch files from commits.
