---
name: tgc-add-field-to-handwritten-resource-skill
description: Add a new field to an existing handwritten resource in TGC. Use when you need to extend a handwritten resource.
---

# tgc-add-field-to-handwritten-resource-skill

When you need to add a new field to an existing handwritten resource in the TGC Next pipeline, use this skill.

## When to Use This Skill

- Use this when a resource in TGC is handwritten (it has `include_in_tgc_next: true` and custom files under `mmv1/third_party/tgc_next/pkg/services/`), and you need to support a new configuration option or nested attribute that is currently missing from the conversion.

---

## How to Use It

To add a new field, you must update the schema definition, the HCL flattener, and the CAI expander, then compile and verify round-trip conversion.

### Phase 1: Locate the Handwritten Resource Files
Locate the Go files under `mmv1/third_party/tgc_next/pkg/services/<product>/`:
- `resource_<resource_name>.go` (Schema definition)
- `resource_<resource_name>_cai2hcl.go` (Flattener)
- `resource_<resource_name>_tfplan2cai.go` (Expander)

### Phase 2: Implement the Changes

#### Step 1: Update the Schema (`resource_<resource_name>.go`)
Add the schema definition of the new field to the schema definition map inside the resource function.

- Copy the exact schema definition from the Terraform Google Provider.
- Ensure type, description, required/optional/computed flags, and any nested fields match.
- If the new field belongs to a nested block (e.g. `addons_config`), find that block's schema definition and append the field to it.

*Example:*
```go
// In resource_<resource_name>.go
"new_field_name": {
	Type:     schema.TypeBool,
	Optional: true,
	Default:  false,
},
```

#### Step 2: Update the Flattener (`resource_<resource_name>_cai2hcl.go`)
Retrieve the value from the CAI asset payload and map it to the HCL schema.

- CAI keys are usually camelCase (e.g., `newFieldName`), whereas HCL keys are snake_case (`new_field_name`).
- If the field is part of a nested block, ensure the parent block's flattener function is updated.
- **Optimization (Default Values)**: Do not convert the field (i.e. do not write it to `hclData`) if the retrieved value is the same as the default value defined in the schema (e.g., if it's `false` for boolean, `""` for string, or matches standard zero-values). This keeps the exported HCL clean and prevents unnecessary diffs.
- **Interface Parameter Type**: Any helper function written to flatten a field or a nested block MUST accept `v interface{}` as its parameter instead of a concrete type. Note that while standard Terraform provider flatteners often accept concrete client API structs (e.g. `*container.StandardRolloutPolicy`), TGC `cai2hcl` flatteners always decode raw JSON maps and must accept `interface{}` to allow safe type-assertion.

*Example:*
```go
// In resource_<resource_name>_cai2hcl.go
if newFieldVal, ok := apiObj["newFieldName"]; ok {
	// Only set the field if it is not the default value (false)
	if boolVal, ok := newFieldVal.(bool); ok && boolVal {
		hclData["new_field_name"] = boolVal
	}
}

// Example helper flattener function accepting interface{}
func flattenNestedConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	
	// Safe type assertion to a map representing the JSON/API object
	obj, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	
	transformed := make(map[string]interface{})
	if val, ok := obj["someApiField"]; ok {
		transformed["some_hcl_field"] = val
	}
	
	return []map[string]interface{}{transformed}
}
```

#### Step 3: Update the Expander (`resource_<resource_name>_tfplan2cai.go`)
Read the value from the Terraform planned resource data and map it to the API format.

- Use `d.GetOk("field_path")` or `d.Get("field_path")` to retrieve the HCL configuration value.
- Write the value to the API map under the correct camelCase API key.

*Example:*
```go
// In resource_<resource_name>_tfplan2cai.go
if v, ok := d.GetOk("new_field_name"); ok {
	obj["newFieldName"] = v
}
```

### Phase 3: Code Generation & Verification
1. **Regenerate TGC downstream code**:
   ```bash
   make tgc OUTPUT_PATH="/path/to/your/terraform-google-conversion"
   ```
2. **Compile downstream and run integration tests**:
   ```bash
   cd /path/to/your/terraform-google-conversion
   make mod-clean
   make build
   
   # Run the integration tests for the resource to verify roundtrip correctness
   make test-integration-local TESTPATH=./test/services/<product> TESTARGS='-run=TestAcc<ResourceName>'
   ```
3. **Verify the intermediate files**:
   Set `WRITE_FILES=true` environment variable before running tests. Inspect the generated test files under downstream `test/services/<product>/`:
   - `Test_export.tf`: Check that `new_field_name` is correctly flattened from the CAI payload.
   - `Test_roundtrip.json`: Check that `newFieldName` is correctly expanded back into the CAI JSON payload.
   - `Test_roundtrip.tf`: Check that `new_field_name` is present in the final roundtrip configuration.
