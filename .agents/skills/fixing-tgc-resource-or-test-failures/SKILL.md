---
name: fixing-tgc-resource-or-test-failures
description: Guide for Phase 6 of TGC Main Loop (Fix). Use to understand how to fix failures.
---

# fixing-tgc-resource-or-test-failures

## When to Use This Skill
- When you need to triage and fix a resource conversion or integration test failure in TGC.

## How to Use It

### 1. Initial Error Reporting
If a failure is detected or provided, you **MUST** report it using the following template before proceeding.

#### Error Report Template:
- **Failed Command**: `[The command that failed]`
- **Detailed Logs**:
  ```
  [Paste the full, relevant error logs here]
  ```
- **Mandatory Skill Read**: `[Yes/No]` (Confirm you have read `tgc-fix-handwritten-resources-tests-skill` or `tgc-fix-integration-tests-skill`)
- **Analysis**: `[Analyze the cause of the failure]`
- **Proposed Solution**: `[Outline the solution and ask for user approval before applying it]`

### 2. Identification
- **Identify Resource Type**: Determine whether the failing resource is a **handwritten** or **generated** resource by checking these conditions:
  - **It is Handwritten if:**
    - There are custom code files ending in `_cai2hcl.go` and `_tfplan2cai.go` for the resource in `mmv1/third_party/tgc_next/` (even if a YAML file exists in `products/`).
  - **It is Generated if:**
    - A YAML file exists in `products/` AND does NOT contain `exclude_resource: true`, AND
    - There are no custom code files ending in `_cai2hcl.go` or `_tfplan2cai.go` in `mmv1/third_party/tgc_next/` for this resource.
- **Check Known Patterns**: Refer to the appropriate skill based on the resource type:
  - For handwritten resources: See `tgc-fix-handwritten-resources-tests-skill`.
  - For generated resources: See `tgc-fix-integration-tests-skill`.

### 3. Fix
- Apply fixes in MMv1.
- **DO NOT** modify templates in `mmv1/templates/terraform`. You may modify `mmv1/templates/tgc_next`.
- **DO NOT** add new fields to `mmv1/api/resource/custom_code.go` unless guided by the user.
- **DO NOT** remove existing custom_code.

### 4. Verification
- Self-check that all failures are reported.
- Proceed to Phase 4 (Generate Code), Phase 5 (Unit Testing), Phase 6 (Integration Testing) after applying fixes.