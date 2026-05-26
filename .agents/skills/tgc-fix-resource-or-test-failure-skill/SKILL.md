---
name: fixing-tgc-resource-or-test-failures
description: Guides the agent through fixing resource conversion or integration test failures in TGC. Use when encountering failures in Phase 6 of the TGC Main Loop.
---

# fixing-tgc-resource-or-test-failures

This skill guides the agent through fixing resource conversion or integration test failures.

1. **Initial Error Reporting**: If a failure is detected or provided, you **MUST** report it using the following template before proceeding.
   
   #### Error Report Template:
   - **Failed Command**: `[The command that failed]`
   - **Detailed Logs**:
     ```
     [Paste the full, relevant error logs here]
     ```
   - **Mandatory Skill Read**: `[Yes/No]` (Confirm you have read `tgc-fix-handwritten-resources-tests-skill` or `tgc-fix-integration-tests-skill`)
   - **Analysis**: `[Analyze the cause of the failure]`
   - **Proposed Solution**: `[Outline the solution and ask for user approval before applying it]`

2. **Identification**:
   - **Identify Resource Type**: Determine whether the failing resource is a **handwritten** or **generated** resource.
     - *Generated Resources*: Follow standard generation patterns from MMv1.
     - *Handwritten Resources*: Often located in or referencing `tgc_next` or having custom manual implementations.
   - **Check Known Patterns**: Refer to the appropriate skill based on the resource type:
     - For handwritten resources: See `tgc-fix-handwritten-resources-tests-skill`.
     - For generated resources: See `tgc-fix-integration-tests-skill`.

3. **Fix**:
   - Apply fixes in MMv1.
   - **DO NOT** modify templates in `mmv1/templates/terraform`. You may modify `mmv1/templates/tgc_next`.
   - **DO NOT** add new fields to `mmv1/api/resource/custom_code.go` unless guided by the user.
   - **DO NOT** remove existing custom_code.

4. **Verification**:
   - Self-check that all failures are reported.
   - Proceed to Phase 3 (Generate Code) and Phase 4 (Unit Testing) after applying fixes.
