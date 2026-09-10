---
trigger: always_on
description: Always-on system prompt for TGC development
---

---
trigger: always_on
description: Enforce TGC development Rules
---

# Environment Setup Rule
Before proceeding to Phase 2 (Implementation) or running any tests, the agent MUST execute a command to verify that `TGC_DIR` is set to the active downstream TGC directory and print its value in the chat. Failure to do so is a violation of process.

The agent MUST run the `tgc-sync-provider` skill during Phase 1 (Session Setup) before Phase 2, to ensure the downstream repository is aligned with Magic Modules.

# TGC development Rules

As an AI agent operating in this repository, you must **ALWAYS** follow these steps before attempting to add a new resource/field to TGC:

1. **[MANDATORY WORKFLOW ENTRYPOINT]**: At the very start of any session (Phase 1: Session Setup), the agent MUST execute the `view_file` tool on `.agents/TGC_WORKFLOWS.md`, and either `.agents/tgc_add.md` or `.agents/tgc_fix.md` depending on the task. The agent MUST print a summary of these entrypoints in their first response and list them as checked items in `task.md`. Failure to call `view_file` on these files in Phase 1 is a fatal process violation.

2. In the magic-modules repository, don't run command `go test` or `go mod tidy`.

3. In the downstream TGC repository, don't run command `go test`.


4. To fix the failed TGC integration tests
   - **don't** modify the templates in `mmv1/templates/terraform`. It is allowed to modify the templates in `mmv1/templates/tgc_next`.
   - **don't** add ignore_read_extra to example in Resource.yaml
   - **don't** add new fields to mmv1/api/resource/custom_code.go unless it is guided by the user
   - **don't** remove any existing custom_code, including any constants
   - **don't** change the schema of a resource (e.g., making a Required field Optional or a Set) to fix conversion failures, unless explicitly guided by the user.
   - **don't** remove `exclude_test: true` from examples in resource YAML files to force test generation for TGC, as this may break tests in the standard Terraform provider.


5. Only commit files under the `mmv1` folder in the branch, and exclude scratch files like `.txt`, `.py`, and `.sh` from commits.

6. DO NOT make changes directly in the downstream repository (`terraform-google-conversion`). All changes must be driven through Magic Modules (`mmv1/`).

7. You must strictly follow the sequence of phases defined in `.agents/tgc_add.md` or `.agents/tgc_fix.md` depending on whether you are adding a resource or fixing a failure (Session Setup -> Implementation -> Unit Testing -> Integration Testing). Code generation (Phase 2) MUST be performed before unit tests (Phase 3), and unit tests MUST be performed before integration tests (Phase 5). Structure your `task.md` to reflect these phases.

8. For any failure (build, unit test, integration test, or verification), stop execution immediately after reporting the error using the template in .agents/TGC_WORKFLOWS.md. Do not proceed with applying any fixes or running subsequent steps until the user has explicitly approved the proposed solution in the chat.

9. **Prioritize process reporting**: If the user's request involves a test failure (input or discovered), you must follow the reporting template specified in `.agents/TGC_WORKFLOWS.md` before proceeding with planning or execution.

10. **Integration Test Subtests**: When running integration tests using `run_integration_test.sh`, if the test name contains an underscore (e.g., `TestAccContainerCluster_withAutopilotClusterPolicy`), it is likely a subtest. You MUST verify if it expects the format `ParentTest/SubTest` (e.g., `TestAccContainerCluster/TestAccContainerCluster_withAutopilotClusterPolicy`) and pass it accordingly as documented in `.agents/TGC_WORKFLOWS.md`. For new resources, you MUST run the top-level test (e.g., `TestAccGKEHub2Feature`) instead of a specific subtest, as the top-level test will cover all of the subtests.

11. **Mandatory Skill Reading for Specialized Tasks**: Before proposing an implementation plan or making code changes for a resource identified as handwritten or generated, you MUST:
    - Identify the resource type.
    - Add a specific task item to `task.md` to read the corresponding skill (e.g., `tgc-fix-handwritten-resources-tests-skill`).
    - Execute the `view_file` tool on that skill and mark the task as completed in `task.md`. **Note: Notwithstanding general Planning Mode instructions, you are authorized and required to create `task.md` during Phase 1 (Session Setup) to track these mandatory steps.**
    Failure to perform the reading step strictly before planning is a violation of process.

12. **Skill Reading Before Proposing Solutions**: The agent MUST NOT propose a solution in the error report or create an implementation plan until it has executed the `view_file` tool on the mandatory skill corresponding to the resource type (either `tgc-fix-handwritten-resources-tests-skill` or `tgc-fix-integration-tests-skill`).

13. **Field Ordering in YAML Files**: When adding or modifying fields in Magic Modules YAML files (e.g., `include_in_tgc_next`, `tgc_decoder`), you MUST ensure they follow the order of fields defined in the `Resource` struct in `mmv1/api/resource.go`.

14. **Tracing Evidence for Integration Tests**: For **integration test failures**, the agent MUST NOT propose a solution or apply a fix until it has explicitly cited evidence from the intermediate files (`Test_export.tf`, `Test_roundtrip.json`, etc.) in the chat to demonstrate tracing.

15. **Mandatory Tracing Checklist**: When fixing integration test failures, the agent MUST include explicit sub-items in `task.md` to track the verification of `Test_roundtrip.tf`, `Test_roundtrip.json`, and `Test_export.tf`.

16. **Run All Failed Tests Before Analysis**: When fixing integration test failures, the agent MUST run all failed integration tests for the affected resource before analyzing the failure or proposing a solution.