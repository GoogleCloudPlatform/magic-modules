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

### 2. Triage & Identification (The Data Loss Check)
You MUST perform a strict **Triage & Identification** process before proposing or applying any changes. Do NOT use `is_missing_in_cai: true` as a shortcut to pass tests if the CAI asset actually supports the field.

#### [NEW] Automated Diagnostics Tool (`diagnose_test_failure.py`)
To automate this triage, run the integration test with `WRITE_FILES=1` enabled:
```bash
export WRITE_FILES=1
TF_CLI_CONFIG_FILE="${PWD}/tf-dev-override.tfrc" GO111MODULE=on go test -v ./test/services/<service> -run="<TestName>"
```
The Go test assertion runner will automatically execute `./test/diagnose_test_failure.py` on failure, producing a side-by-side report comparing original CAI assets, HCL exports, and roundtrip configurations. Review this report in the standard test logs.

#### Step A: Verify CAI Asset Support
Inspect the raw CAI asset file (`Test.json` or `Test_export.json`) generated during the failed integration test run (or the automated diagnostics report).
* **Rule**: If the missing property (e.g., `"gpuDirectStrategy": "RDMA"`) is present in the CAI asset, it is **NOT** missing in CAI.
* **Constraint**: Using `is_missing_in_cai: true` is **strictly forbidden**. You MUST implement full schema mapping, expander, and flattener logic.
* **Condition**: Only use `is_missing_in_cai: true` if you have verified that the field is permanently omitted from CAI schemas and CAI asset payloads across all tests for this resource.

#### Step B: Identify Resource Type (Handwritten vs. Generated)
Determine whether the failing resource is **handwritten** or **generated** by checking these conditions:
* **It is Handwritten if:**
  - There are custom code files ending in `_cai2hcl.go` and `_tfplan2cai.go` for the resource in `mmv1/third_party/tgc_next/` (even if a YAML file exists in `products/`).
* **It is Generated if:**
  - A YAML file exists in `products/` AND does NOT contain `exclude_resource: true`, AND
  - There are no custom code files ending in `_cai2hcl.go` or `_tfplan2cai.go` in `mmv1/third_party/tgc_next/` for this resource.

#### Step C: Select Troubleshooting Skill
Based on the resource type, read the corresponding skill before designing a fix:
* **Handwritten Resources**: Refer to `tgc-fix-handwritten-resources-tests-skill`.
* **Generated Resources**: Refer to `tgc-fix-integration-tests-skill`.

### 3. Fix & Implementation Playbook
* **Replicate Standard Provider Patterns**:
  1. Locate how the standard provider template (e.g., `node_config.go.tmpl`) maps the field.
  2. Replicate that exact Go/YAML logic directly into the resource's mappings to ensure 100% parity.
* **DO NOT** modify templates in `mmv1/templates/terraform` or shared helper files/templates in `mmv1/third_party/terraform`. You may modify templates in `mmv1/templates/tgc_next`.
* **DO NOT** modify legacy TGC code in `mmv1/third_party/cai2hcl` or `mmv1/third_party/tgc`. All changes must be made to TGC Next code in `mmv1/third_party/tgc_next`.
* **DO NOT** add new fields to `mmv1/api/resource/custom_code.go` unless guided by the user.
* **DO NOT** modify `mmv1/api/resource.go` to hardcode or manually append ignored fields for tests. Use resource YAML property renaming with `api_name` instead to align with the actual Terraform schema.
* **DO NOT** remove existing custom_code.

### 4. Verification
- Self-check that all failures are reported.
- Proceed to Phase 4 (Generate Code), Phase 5 (Unit Testing), Phase 6 (Integration Testing) after applying fixes.