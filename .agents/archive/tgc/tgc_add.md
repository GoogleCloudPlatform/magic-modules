# TGC Resource Addition Workflow

This document defines the main loop for adding a new resource to TGC.

## Required Skills
Before proceeding with the workflow, ensure you are familiar with and read the following skills when prompted in the phases:
- [tgc-sync-provider](skills/tgc/tgc-sync-provider/SKILL.md) (Phase 1)
- [tgc-new-generated-resource-skill](skills/tgc/tgc-add-new-generated-resource-skill/SKILL.md) (Phase 2)
- [tgc-build-skill](skills/tgc/tgc-build-skill/SKILL.md) (Phase 3)
- [tgc-run-integration-tests-skill](skills/tgc/tgc-run-integration-tests-skill/SKILL.md) (Phase 4)
- [tgc-fix-integration-tests-skill](skills/tgc/tgc-fix-integration-tests-skill/SKILL.md) (Phase 5)

## The Workflow

### 1. Session Setup
- **Initialize Task List**: You MUST create a `task.md` file containing all workflow phases (Session Setup, Implementation, Generate Code, Unit Testing, Integration Testing, Finalization) as uncompleted tasks before performing any other actions. If you encounter an integration test failure later, you MUST include explicit sub-tasks in the Fix phase of `task.md` to trace the failure backwards (inspecting `Test_roundtrip.tf`, `Test_roundtrip.json`, `Test_export.tf`, `Test.json`).
- **Set Environment**: Ensure `TGC_DIR` environment variable is set to the absolute path of your active TGC downstream workspace.
  ```bash
  export TGC_DIR=/path/to/downstream/workspace
  export PATH=/usr/local/go/bin:/opt/homebrew/bin:$PATH
  ```
- **Use Skill**: You MUST read the [tgc-sync-provider](skills/tgc/tgc-sync-provider/SKILL.md) skill. **Ask the user** which synchronization method to use (Aligning to Base Commit vs. Fast-Forward to Latest) and follow their choice to synchronize the downstream repository before proceeding to Phase 2.

### 2. Implementation (Parent Agent)
- **Read Skill**: Read [tgc-new-generated-resource-skill](skills/tgc/tgc-add-new-generated-resource-skill/SKILL.md) for guidance on adding a new resource.
- **Define Resource**: Add or modify the resource definition in Magic Modules (`mmv1/products/...`).
- **Field Ordering**: Ensure fields in YAML files follow the order defined in `mmv1/api/resource.go`.
- **Redundant Overrides**: Ensure you do not add redundant `cai_asset_name_format` overrides that match `id_format`. Run `.agents/scripts/tgc/check_redundant_cai_overrides.py <yaml-file>` to verify.

### 3. Generate Code & Unit Testing
- **Read Skill**: Read [tgc-build-skill](skills/tgc/tgc-build-skill/SKILL.md) to project changes to the downstream repository.
- **Selective Unit Testing**: The build pipeline (`build_tgc.sh`) automatically executes the selective unit test runner during code generation. If any changed unit tests fail, the build will block immediately.
- If build, unit test, or dependency errors occur, stop and immediately report the error in the conversation using the required template.

### 4. Integration Testing
- **Run All Tests**: For a new resource, you MUST run the top-level test (e.g., `TestAccGKEHub2Feature`) instead of a specific subtest, as the top-level test will cover all of the subtests.
- **Read Skill**: Read [tgc-run-integration-tests-skill](skills/tgc/tgc-run-integration-tests-skill/SKILL.md) for guidance on running integration tests.

> [!NOTE]
> If no tests are generated for the resource (e.g., the `test/services/<product>` directory is missing or empty), refer to [tgc-fix-integration-tests-skill](skills/tgc/tgc-fix-integration-tests-skill/SKILL.md) (Troubleshooting Playbook Item 11) for guidance on forcing test generation using `tgc_tests` and bootstrapping data files with `WRITE_FILES=true`.

### 5. Fixes
- If tests fail, **Read Skill**: 
  Read [tgc-fix-integration-tests-skill](skills/tgc/tgc-fix-integration-tests-skill/SKILL.md).
  *(You MUST execute `view_file` on the corresponding skill before proposing a solution or implementation plan).* apply fixes in MMv1.
- **Trace Failure**: Follow the "Tracing Failures Backwards" protocol in the Playbook to isolate the stage where data was lost before proposing a solution.
- **Report Failure**: Report the failure using the template in [TGC_WORKFLOWS.md](TGC_WORKFLOWS.md).
- **[MANDATORY] Stop and wait for user approval before applying the fix.**
- **Apply Fix**: Apply fixes in Magic Modules (`mmv1/`). DO NOT make changes directly in the downstream repository.
- **DON'T** change the schema of a resource (e.g., making a Required field Optional or a Set) to fix conversion failures, unless explicitly guided by the user.
- **Repeat Loop**: After ANY fix, you MUST repeat the full verification loop (Step 3: Generate Code & Unit Testing, Step 4: Integration Testing).

### 6. Finalization
- Ask the user if the task is complete and if you should proceed with committing.
- Include a summary of any failures encountered using the template specified in [TGC_WORKFLOWS.md](TGC_WORKFLOWS.md).
- Commit changes under `mmv1/` folder only.
- Exclude scratch files from commits.
