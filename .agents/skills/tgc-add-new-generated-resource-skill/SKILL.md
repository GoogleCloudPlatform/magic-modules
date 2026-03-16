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

### 2. Handle Custom Code

In the `Resource.yaml` file, search for any custom code blocks and apply the appropriate solutions based on the type of custom code:

#### A. `custom_flatten`

- **Issue 1 - Zero Values**: The field is populated in Terraform state by a `custom_flatten` function in Magic Modules. These functions often use `d.Get()`, which reads from the Terraform state. During `CAI -> HCL`, `d.Get()` returns zero values, so the field appears missing or has an incorrect zero value. This can also happen if the custom flattener explicitly sets a value to nil.
  - **Solution**: Add `tgc_ignore_terraform_custom_flatten: true` to the field's definition in the YAML. This tells TGC not to execute the Magic Module's `custom_flatten` function for this field during the conversion. Instead, TGC will rely on direct mapping from the CAI JSON.
  - **Examples**:
    - `filestore/Instance.yaml`
    - `apphub/ServiceProjectAttachment.yaml`
    - `datastream/ConnectionProfile.yaml`
    - `secretmanager/SecretVersion.yaml`

- **Issue 2 - Panics/Type Mismatches**: A field's `custom_flatten` generated code panics or errors out during `CAI -> HCL` conversion because of a type mismatch (e.g., CAI returns an empty map `{}` but HCL expects a boolean `true`), or due to unsupported `d.Set()` / `d.Get()` operations that fail during CAI mapping.
  - **Solution**: Add `tgc_ignore_terraform_custom_flatten: true` to the field definition in the YAML to skip executing the custom flatten function and avoid the panic. If the field *still* needs complex data transformation from CAI format to HCL, use a `tgc_decoder`.
    - Add `tgc_decoder: 'templates/tgc_next/decoders/resource_name.go.tmpl'` under `custom_code` in the YAML.
    - A `tgc_decoder` gives you direct access to safely read the raw `res map[string]interface{}` and write to the output `hclData map[string]interface{}` bypassing strict Terraform schemas and validations.
  - **Example**:
    - `bigtable/AppProfile.yaml`

#### B. `custom_expand`

- **Issue - API Calls During Expansion**: The process of converting HCL to the CAI format, which likely happens within `tfplan2cai`, is attempting to call the GCP API to validate or fetch details about a resource (e.g., to resolve an instance name to an ID). The API call fails because API calls should be avoided during custom expander functions.
  - **Solution**: Override the default expander logic for the field causing the API call.
    - Identify the field in the YAML (e.g., `nextHopInstance` in `Route.yaml`).
    - Add a `custom_tgc_expand` entry pointing to a custom Go template.
  - **Example**: `custom_tgc_expand: 'templates/terraform/custom_expand/resourceref_with_validation.go.tmpl'`

#### C. `decoder`

- **Issue 1 - Zero Value Defaults**: The resource's default Terraform decoder (used to convert API responses to Terraform state) contains logic (`d.Get()` or similar functions) that sets the field to a zero value during the `CAI -> HCL` conversion since there's no prior Terraform state. The value from the CAI asset is ignored.
  - **Solution**: Add `tgc_ignore_terraform_decoder: true` to the resource or field in the YAML. This prevents TGC from using the standard Magic Modules decoder template for this item during `CAI -> HCL`.
  - **Examples**:
    - `spanner/Instance.yaml`
    - `spanner/Database.yaml`
      - *Enabled custom_code `tgc_ignore_terraform_decoder: true`*: This prevents the standard Terraform decoder from running during TGC conversion. The standard decoder attempts to make API calls (to fetch `authString`), which causes "client is nil" errors in the TGC environment where no authorized client is available. It also avoids "project: required field is not set" errors that originated from the decoder's validation logic.
    - `dataproc/Batch.yaml`
    - `kms/KeyRing.yaml`

- **Issue 2 - API Calls During CAI2HCL**: The process of converting the CAI format to HCL, which likely happens within `cai2hcl`, is attempting to call the GCP API. The API call fails because API calls should be avoided during `cai2hcl`.
  - **Solution**: Add `tgc_ignore_terraform_decoder: true` to the resource in the YAML.
  - **Examples**:
    - `redis/Cluster.yaml`
    - `redis/Instance.yaml`