---
name: tgc-add-new-generated-resource-skill
description: Add a new generated resource to TGC. Use when you need to add a new generated resource to TGC.
---

# tgc-add-new-generated-resource-skill

When you need to add a new generated resource to TGC, use this skill.

## When to Use This Skill

- Use this when adding a new generated resource to TGC.
- This is helpful when you need to understand the structural steps and configurations needed to expose a generated resource to the Terraform Google Conversion (TGC) library.

---

## How to Use It

If you added or modified a generated resource, follow the steps below carefully.

### 1. Map and Enable

- **Mapping**: Use a script or command to locate `mmv1/products/.../Resource.yaml` for each `google_` type.
  - **Search pattern**: Look for `name: 'resource_name'` or perform a filename match.
- **Enabling**: For each found YAML file:
  - Ensure `include_in_tgc_next: true` is present at the **top-level**.
  - Place it in the proper order according to the progression of fields in the `mmv1/api/resource.go` file.

### 2. Check for URL Parameters and Asset Name Format

- If the resource has parameters marked as `url_param_only: true` and `required: true`, verify if they can be extracted from the CAI asset name during `cai2hcl` conversion.
- If the `self_link` in the YAML file is just `{{name}}` or does not contain all the required parameters in its pattern, you MUST specify `cai_asset_name_format` at the top-level to define the pattern for extraction.
- Example: `cai_asset_name_format: 'projects/{{project}}/locations/{{location}}/notificationConfigs/{{config_id}}'`

### 3. Check for Custom Flatteners

- When enabling an existing generated resource for TGC, check if any fields use `custom_flatten`.
- If the custom flattener uses `d.Get(...)` instead of reading from the passed value `v` (common in shared templates), it will return empty values during `cai2hcl` conversion because there is no Terraform state.
- If this causes issues (e.g., dropping required fields), consider adding `tgc_ignore_terraform_custom_flatten: true` to the field's definition in the YAML to use the default mapping.

### 4. Skipping Tests Safely for TGC

- Tests generated from examples and handwritten tests in `third_party` are shared with the standard Google Provider. DO NOT use `exclude_test: true` in examples or rename handwritten tests to skip them for TGC, as this will affect the Google Provider as well!
- To skip a test generated from an **example** for TGC only: Add `tgc_skip_test: 'Reason for skipping'` to the example definition in the resource's YAML file.
- To skip a **handwritten test** for TGC only: Add the test name to the `tgc_tests` section at the top-level of the resource's YAML file with `skip: 'Reason for skipping'`. This prevents the generator from creating duplicates and applies the skip.

### 5. Handling Missing CAI Data vs Schema Requirements

- If the CAIS API does not return certain fields, they will be missing in the input CAI asset files for tests.
- If the resource schema requires at least one of several blocks to be specified, and CAIS returns an empty block (which `cai2hcl` drops), it may fail validation (`Invalid combination of arguments`).
- **Solution**: Implement a custom `tgc_decoder` in `mmv1/templates/tgc_next/decoders/` to inject minimal valid data or an empty map to satisfy the schema when data is missing in CAI.

### Troubleshooting Build Failures

### Missing Package Dependency in Shared Templates
- **Symptom**: `go mod tidy` or compilation fails after generation because a package (e.g., `compute`) is not found in the TGC environment.
- **Cause**: Shared templates in `mmv1/templates/terraform/constants` may contain hardcoded imports or functions relying on packages not available in TGC.
- **Solution**: Wrap the problematic code in the template with a compiler condition to exclude it for TGC generation. You can use the helper method `IsTgcCompiler`:
  ```tmpl
  {{- if not $.ResourceMetadata.ProductMetadata.IsTgcCompiler }}
  // Code to exclude for TGC (only included for standard Terraform provider)
  {{- end }}
  ```
  *Note: The exact path to `IsTgcCompiler` may vary depending on the template's context (e.g., `$.IsTgcCompiler` or `$.ProductMetadata.IsTgcCompiler`).*

### No Tests Generated Failure
- **Symptom**: `Error generating resource tests: No TGC tests for resource <ResourceName>`
- **Action**: This commonly happens when all examples in the YAML are excluded or there is a file naming mismatch. Please refer to **Item 11** in the [Troubleshooting Playbook](file:///Users/zhenhuali/Documents/workspace/feature-a/.agents/skills/tgc-fix-integration-tests-skill/troubleshooting_playbook.md) for detailed causes and solutions.
