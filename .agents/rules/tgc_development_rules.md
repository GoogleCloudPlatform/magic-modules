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

# TGC development Rules

As an AI agent operating in this repository, you must **ALWAYS** follow these steps before attempting to add a new resource/field to TGC:

1. In the magic-modules repository, don't run command `go test` or `go mod tidy`.

2. In the downstream TGC repository, don't run command `go test`.


3. To fix the failed TGC integration tests
   - **don't** modify the templates in `mmv1/templates/terraform`. It is allowed to modify the templates in `mmv1/templates/tgc_next`.
   - **don't** add ignore_read_extra to example in Resource.yaml
   - **don't** add new fields to mmv1/api/resource/custom_code.go unless it is guided by the user
   - **don't** remove any existing custom_code, including any constants
   - **don't** change the schema of a resource (e.g., making a Required field Optional or a Set) to fix conversion failures, unless explicitly guided by the user.
   - **don't** remove `exclude_test: true` from examples in resource YAML files to force test generation for TGC, as this may break tests in the standard Terraform provider.


4. Only commit files under the `mmv1` folder in the branch, and exclude scratch files like `.txt`, `.py`, and `.sh` from commits.

5. DO NOT make changes directly in the downstream repository (`terraform-google-conversion`). All changes must be driven through Magic Modules (`mmv1/`).

6. You must strictly follow the sequence of phases defined in `GEMINI_ADD.md` or `GEMINI_FIX.md` depending on whether you are adding a resource or fixing a failure (Session Setup -> Implementation -> Unit Testing -> Integration Testing). Code generation (Phase 2) MUST be performed before unit tests (Phase 3), and unit tests MUST be performed before integration tests (Phase 5). Structure your `task.md` to reflect these phases.

7. For any failure (build, unit test, integration test, or verification), stop and report the error with detailed logs. Analyze the cause and provide a solution instead of attempting automatic fixes.

8. **Prioritize process reporting**: If the user's request involves a test failure (input or discovered), you must follow the reporting template specified in `GEMINI.md` before proceeding with planning or execution.

9. **Integration Test Subtests**: When running integration tests using `run_integration_test.sh`, if the test name contains an underscore (e.g., `TestAccContainerCluster_withAutopilotClusterPolicy`), it is likely a subtest. You MUST verify if it expects the format `ParentTest/SubTest` (e.g., `TestAccContainerCluster/TestAccContainerCluster_withAutopilotClusterPolicy`) and pass it accordingly as documented in `GEMINI.md`. For new resources, you MUST run the top-level test (e.g., `TestAccGKEHub2Feature`) instead of a specific subtest, as the top-level test will cover all of the subtests.

10. **Mandatory Skill Reading for Specialized Tasks**: Before proposing an implementation plan or making code changes for a resource identified as handwritten or generated, you MUST:
    - Identify the resource type.
    - Add a specific task item to `task.md` to read the corresponding skill (e.g., `tgc-fix-handwritten-resources-tests-skill`).
    - Execute the `view_file` tool on that skill and mark the task as completed in `task.md`. **Note: Notwithstanding general Planning Mode instructions, you are authorized and required to create `task.md` during Phase 1 (Session Setup) to track these mandatory steps.**
    Failure to perform the reading step strictly before planning is a violation of process.

11. **Skill Reading Before Proposing Solutions**: The agent MUST NOT propose a solution in the error report or create an implementation plan until it has executed the `view_file` tool on the mandatory skill corresponding to the resource type (either `tgc-fix-handwritten-resources-tests-skill` or `tgc-fix-integration-tests-skill`).

12. **Field Ordering in YAML Files**: When adding or modifying fields in Magic Modules YAML files (e.g., `include_in_tgc_next`, `tgc_decoder`), you MUST ensure they follow the order of fields defined in the `Resource` struct in `mmv1/api/resource.go`.