# TGC Test Fix Workflow

This document defines the main loop for fixing a resource conversion or integration test failure in TGC.

## Required Skills
Before proceeding with the workflow, ensure you are familiar with and read the following skills when prompted in the phases:
- `sync-provider` (Phase 1)
- `tgc-fix-handwritten-resources-tests-skill` (Phase 3 - if resource is handwritten)
- `tgc-fix-integration-tests-skill` (Phase 3 - if resource is generated)
- `fixing-tgc-resource-or-test-failures` (Phase 3)

## The Workflow

### 1. Session Setup
- **Initialize Task List**: You MUST create a `task.md` file containing all workflow phases (Session Setup, Triage & Isolate, Fix, Generate Code, Unit Testing, Integration Testing, Finalization) as uncompleted tasks before performing any other actions.
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
- **Run All Tests**: Before analyzing the failure or proposing a solution, you MUST run all integration tests for the affected resource (e.g., by running the top-level test instead of a specific subtest) to identify all potential failures and get a complete picture.
  ```bash
  .agents/skills/tgc-run-integration-tests-skill/scripts/run_integration_test.sh <test-path> <test-name>
  ```

### 3. Fix (Parent Agent)
- **Use Skill**: You MUST read the `fixing-tgc-resource-or-test-failures` skill.
- **Trace Failure**: Follow the "Tracing Failures Backwards" protocol in the Playbook to isolate the stage where data was lost before proposing a solution.

### 4. Generate Code
- **Generate Code**: Use the automation script `./.agents/skills/tgc-build-skill/scripts/build_tgc.sh $TGC_DIR` to project changes to the downstream repository.
- If build or dependency errors occur, stop and immediately report the error.

### 5. Unit Testing
- Run the following command for folders pkg and test. **Do NOT scope to specific services or tests; all unit tests in./pkg must be run:**
  ```bash
  make test-local TEST=./test
  make test-local TEST=./pkg/...
  ```

### 6. Integration Testing
- Verify the fix by running the test(s) again using the script:
  ```bash
  .agents/skills/tgc-run-integration-tests-skill/scripts/run_integration_test.sh <test-path> <test-name>
  ```
> [!IMPORTANT]
> **After ANY fix applied, you MUST repeat the full verification loop:**
> - Step 4 (Generate Code)
> - Step 5 (Unit Testing)
> - Step 6 (Integration Testing)
> Do not skip any of these steps to ensure no new regressions are introduced.

### 7. Finalization
- Ask the user if the task is complete and if you should proceed with committing.
- Commit changes under `mmv1/` folder only.
- Exclude scratch files from commits.
