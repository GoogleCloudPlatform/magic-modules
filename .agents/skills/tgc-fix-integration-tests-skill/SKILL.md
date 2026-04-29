---
name: tgc-fix-integration-tests-skill
description: Fix integration tests for TGC. Use when you need to fix integration tests for TGC.
---

# tgc-fix-integration-tests-skill

When you need to fix integration tests for TGC, use this skill.

## When to Use This Skill

- Use this when fixing integration tests for TGC.
- This is helpful when you need to address build and testing failures for any Terraform Google Conversion (TGC) resource.

---

## How to Use It

To fix integration tests, follow the guidelines and playbook below:

### 1. General Rules for Fixing Tests
When troubleshooting and resolving test failures, adhere to these constraints:
- **DON'T** modify the templates in `mmv1/templates/terraform`. It is **only** allowed to modify the templates in `mmv1/templates/tgc_next`.
- **DON'T** add `ignore_read_extra` to the example in `Resource.yaml`.
- **DON'T** add new fields to `mmv1/api/resource/custom_code.go` unless explicitly guided by the user.
- **DON'T** remove any existing `custom_code`, including any constants.
- **DO** add a comment for each fix in the YAML file or other files to explain the root cause and the solution.
- **DON'T** use `d.Set` in custom decoders for `cai2hcl`. Conversion in `cai2hcl` is a direct mapping from CAI asset data maps to HCL maps without involving Terraform state. Mutate the data map directly.
- **DON'T** use `is_missing_in_cai: true` if the missing fields are present in other raw JSON files of other tests.
- **DO** trace the value of a failing field through `Test_export.tf`, `Test_roundtrip.json`, and `Test_roundtrip.tf` to identify the exact stage where data is lost or mutated, rather than guessing based on the final error message.

---

### 2. Test Discovery and Naming Conventions

When tests are not being generated or run for a resource, ensure that the handwritten test files in `mmv1/third_party/terraform/services/<product>/` follow the expected naming convention: `resource_<product>_<resource_name>_test.go`.

If the file name does not match, the generator function `addTestsFromHandwrittenTests` will fail to find it and will log a message like `no handwritten test file found for <resource>`. Renaming the file to match the expected convention allows the generator to discover it.

---

### 3. Troubleshooting Playbook

The detailed troubleshooting playbook has been moved to a separate file to keep this skill description concise.

Please refer to [Troubleshooting Playbook](file:///Users/zhenhuali/Documents/workspace/tgc-supported-resources/.agents/skills/tgc-fix-integration-tests-skill/troubleshooting_playbook.md) for solutions to common test failures.