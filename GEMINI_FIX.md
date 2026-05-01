# TGC Test Fix Workflow

This document defines the main loop for fixing a resource conversion or integration test failure in TGC.

## Required Skills
Before proceeding with the workflow, ensure you are familiar with and read the following skills when prompted in the phases:
- `sync-provider` (Phase 1)
- `fixing-tgc-resource-or-test-failures` (Phase 3)
- `tgc-fix-handwritten-resources-tests-skill` (Phase 3 - if resource is handwritten)
- `tgc-fix-integration-tests-skill` (Phase 3 - if resource is generated)
- `tgc-build-skill` (Phase 4)
- `tgc-run-unit-tests-skill` (Phase 5)
- `tgc-run-integration-tests-skill` (Phase 6)

## The Workflow

### 1. Session Setup
- **Initialize Task List**: You MUST create a `task.md` file containing all workflow phases (Session Setup, Triage & Isolate, Run All Integration Tests, Fix, Generate Code, Unit Testing, Integration Testing, Finalization) as uncompleted tasks before performing any other actions. If the failure involves an integration test, you MUST include explicit sub-tasks in the Fix phase of `task.md` to trace the failure backwards (inspecting `Test_roundtrip.tf`, `Test_roundtrip.json`, `Test_export.tf`, `Test.json`). You MUST also add a specific task item in Phase 3 to run all integration tests for the resource.
- **Set Environment**: Ensure `TGC_DIR` environment variable is set to the absolute path of your active TGC downstream workspace.
  ```bash
  export TGC_DIR=/path/to/downstream/workspace
  export PATH=/usr/local/go/bin:/opt/homebrew/bin:$PATH
  ```
- **Use Skill**: You MUST read the `sync-provider` skill. **Ask the user** which synchronization method to use (Aligning to Base Commit vs. Fast-Forward to Latest) and follow their choice to synchronize the downstream repository before proceeding to Phase 2.

### 2. Triage & Isolate
- **Generate Code**: Before running tests, you MUST generate code to ensure the downstream repository is up to date with Magic Modules.
  ```bash
  ./.agents/skills/tgc-build-skill/scripts/build_tgc.sh $TGC_DIR
  ```

### 3. Run All Integration Tests
- **Run All Tests**: Before analyzing the failure or proposing a solution, you MUST run all integration tests for the affected resource (e.g., by running the top-level test instead of a specific subtest) to identify all potential failures and get a complete picture.
> [!CRITICAL]
> Do not propose a solution or plan based only on a user-reported subtest failure. You MUST run the top-level test first to discover all failures. 
- **Read Skill**: Read `tgc-run-integration-tests-skill` for guidance on running integration tests.

### 4. Fix (Parent Agent)
- **Use Skill**: You MUST read the `fixing-tgc-resource-or-test-failures` skill.
- **Trace Failure**: Follow the "Tracing Failures Backwards" protocol in the Playbook to isolate the stage where data was lost before proposing a solution.
- **Report Failure**: Report the failure using the template in `GEMINI.md`.
- **[MANDATORY] Stop and wait for user approval before applying the fix.**
- **Apply Fix**: Apply fixes in Magic Modules (`mmv1/`). DO NOT make changes directly in the downstream repository.
- **DON'T** change the schema of a resource (e.g., making a Required field Optional or a Set) to fix conversion failures, unless explicitly guided by the user.
- **Repeat Loop**: After ANY fix, you MUST repeat the full verification loop (Step 5: Generate Code, Step 6: Unit Testing, Step 7: Integration Testing).**

### 5. Generate Code
- **Read Skill**: Read `tgc-build-skill` to project changes to the downstream repository.
- If build or dependency errors occur, stop and immediately report the error in the conversation using the required template.

### 6. Unit Testing
- **Read Skill**: Read `tgc-run-unit-tests-skill` for guidance on running unit tests.

### 7. Integration Testing
- **Read Skill**: Read `tgc-run-integration-tests-skill` for guidance on running integration tests.

### 8. Finalization
- Ask the user if the task is complete and if you should proceed with committing.
- Include a summary of any failures encountered using the template specified in GEMINI.md.
- Commit changes under `mmv1/` folder only.
- Exclude scratch files from commits.
