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

---

### 2. Troubleshooting Playbook

This playbook helps diagnose and fix common issues encountered when running TGC tests. These tests often involve converting between CAI (Cloud Asset Inventory) JSON and HCL (Terraform).

#### Debugging the Conversion Workflow

These tests check the accuracy of the conversions between Cloud Asset Inventory (CAI) format and Terraform HCL. The workflow consists of the following stages:

1. **Initial Conversion (CAI to HCL)**: The input CAI asset file (`Test.json`) is converted into a Terraform configuration file (`Test_export.tf`) using the `cai2hcl` tool.
   * **Verify**:
     - Does `Test_export.tf` correctly represent the resources defined in `Test.json`?
     - Does the Terraform decoder cause the issue?
     - Does the Terraform flattener cause the issue?

2. **Plan and Convert (HCL to CAI)**: The generated `Test_export.tf` is used to create a Terraform plan, which is then converted into a CAI asset format file (`Test_roundtrip.json`) using the `tfplan2cai` tool.
   * **Verify**:
     - Does the Terraform plan on `Test_export.tf` run without errors? 
     - Does `Test_roundtrip.json` accurately reflect the planned state?
     - Does the Terraform encoder cause the issue?
     - Does the Terraform expander cause the issue?

3. **Round-trip Conversion (CAI to HCL)**: The CAI asset file from the previous step (`Test_roundtrip.json`) is converted back into a Terraform configuration file (`Test_roundtrip.tf`) using `cai2hcl`.
   * **Verify**:
     - Is `Test_roundtrip.tf` semantically equivalent to `Test_export.tf`?
     - **CAI to HCL**: Does the Terraform decoder/flattener cause the issue?
     - **HCL to CAI**: Does the Terraform encoder/expander cause the issue?

> **Tip**: When a test fails, inspect the intermediate files (`Test_export.tf`, `Test_roundtrip.json`, `Test_roundtrip.tf`) available in the specific service folder under `test/services` to know which stage of this workflow is introducing the error or unexpected diff.

---

### Common Test Failures & Solutions

#### 1. Conflicting Fields in HCL
- **Symptom**: Error message like `Invalid combination of arguments ... "next_hop_ilb": only one of ... can be specified, but next_hop_ilb,next_hop_ip were specified.`
- **Cause**: The CAI asset to HCL conversion (`cai2hcl`) produced HCL where multiple mutually exclusive fields are set.
- **Solution**: Implement custom logic in a `tgc_decoder` to enforce exclusivity. This code runs during the `CAI -> HCL` conversion.
  1. Add `tgc_decoder: 'templates/tgc_next/decoders/xxxx.go.tmpl'` to the `custom_code` section for the resource in its Magic Modules YAML file.
  2. Create the `.go.tmpl` file in the `templates/tgc_next/decoders` folder.
  3. In the template, write Go code to inspect the fields and unset the lower-priority ones.
- **Example**:
  - **Failing test**: `TestAccComputeRoute_routeIlbExample/step1`
  - **Error**: `only one of next_hop_gateway,next_hop_ilb,next_hop_instance,next_hop_ip,next_hop_vpn_tunnel can be specified, but next_hop_ilb,next_hop_ip were specified.`
  - **Debug**: Check the test JSON file `TestAccComputeRoute_routeIlbExample_step1.json`. `nextHopIlb` and `nextHopIp` both exist. After `cai2hcl` conversion, both `next_hop_ilb` and `next_hop_ip` exist in the converted HCL in the `TestAccComputeRoute_routeIlbExample_step1_export.tf` file.
  - **Solution**: Add `tgc_decoder: 'templates/tgc_next/decoders/compute_route.go.tmpl'` to the `custom_code` section for the resource in its Magic Modules YAML file `Route.yaml`. Create the `.go.tmpl` file under `mmv1/templates/tgc_next/decoders/compute_route.go.tmpl`. In the template, write Go code to inspect the fields and unset the lower-priority ones. If `nextHopIlb` is present, unset `nextHopIp`.

#### 2. Required Fields Missing After CAI to HCL Conversion
- **Symptom**: Terraform error like `At least 1 "trust_anchors" blocks are required.`
- **Cause**: The Terraform schema requires a field or block, but it's not present in the CAI asset and thus not in the HCL generated by the default `cai2hcl` process.
- **Solution**: Use a `tgc_decoder` to add the missing field/block during `CAI -> HCL`. This might involve setting default values or deriving them.
  1. Add `tgc_decoder` to the resource's YAML file.
  2. Create the `.go.tmpl` file in the `mmv1/templates/tgc_next/decoders` folder.
  3. Implement logic in the `.go.tmpl` to populate the required field.
- **Example 1**:
  - **Failing test**: `TestAccIAMBetaWorkloadIdentityPoolProvider_x509`
  - **Error message**: `At least 1 "trust_anchors" blocks are required.`
  - **Debug**: In `TestAccIAMBetaWorkloadIdentityPoolProvider_x509_step1.json`, the field `pemCertificate` doesn’t exist in the CAI asset as it is sensitive information.
    ```json
    "x509": {
      "trustStore": {
        "trustAnchors": [{}]
      }
    }
    ```
    As a result, in `TestAccIAMBetaWorkloadIdentityPoolProvider_x509_step1_export.tf`, the converted HCL doesn’t have the field `pem_certificate` after `cai2hcl` conversion.
    ```hcl
    x509 {
      trust_store {}
    }
    ```
  - **Solution**: Add `tgc_decoder` to the resource's YAML file. Create the `.go.tmpl` file `mmv1/templates/tgc_next/decoders/iam_workload_identity_pool_provider.go.tmpl`. Implement logic in the template to set `pemCertificate` to `"unknown"` if it is missing in the CAI asset.
- **Example 2**:
  - **Failing test**: `TestAccCertificateManagerCertificate_certificateManagerClientAuthCertificateExample`
  - **Error message**: `"self_managed": one of managed,self_managed must be specified`
  - **Debug**: Field `selfManaged` is empty in the CAI asset data in `TestAccCertificateManagerCertificate_certificateManagerSelfManagedCertificateExample_step1.json`. The field `self_managed` is missing in the converted HCL after `cai2hcl` conversion. The top-level `pemCertificate` exists in the CAI asset data in the JSON file and can be used for `selfManaged.pemCertificate` instead.
  - **Solution**: Add `tgc_decoder` to the resource's YAML file `Certificate.yaml`. Create the `.go.tmpl` file `mmv1/templates/tgc_next/decoders/iam_workload_identity_pool_provider.go.tmpl`. Implement logic to set `pemCertificate` to the top-level `pemCertificate` (as it is missing in the CAI asset block) and set `pemPrivateKey` to `"unknown"`.

#### 3. Error Converting Round-Trip Config (HCL -> CAI) - API Calls
- **Symptom**: Error message like `error when converting the round-trip config: &fmt.wrapError{msg:"tfplan2ai converting: googleapi: Error 404: The resource 'projects/ci-test-project-nightly-ga/zones/us-central1-b/instances/tf-test-jqz6svzykj' was not found..."}`
- **Cause**: The process of converting HCL to the CAI format (`custom_expand`), which likely happens within `tfplan2cai`, is attempting to call the GCP API to validate or fetch details about a resource (e.g., to resolve an instance name to an ID), but the resource doesn't exist in the test project during this phase.
- **Solution**: Override the default expansion logic for the field causing the API call.
  1. Identify the field in the YAML (e.g., `nextHopInstance` in `Route.yaml`).
  2. Add a `custom_tgc_expand` entry pointing to a custom Go template.
- **Example**: `custom_tgc_expand: 'templates/terraform/custom_expand/resourceref_with_validation.go.tmpl'`

#### 4. Argument Required, But No Definition Found - Encoder Issue
- **Symptom**: Error message like `The argument "name" is required, but no definition was found.`
- **Cause**: The field is in the CAI asset data. The field is present in the HCL but is being dropped when TGC converts the HCL to the CAI JSON representation. This often means the default Terraform provider's "encoder" function for this resource in Magic Modules drops the field.
- **Solution**: Add `tgc_ignore_terraform_encoder: true` to the field's definition in the resource's YAML file. This tells TGC to use its own logic (or custom TGC encoders) for this field when going from HCL to CAI, instead of relying on the Terraform provider's encoder.
- **Example**: `tgc_ignore_terraform_encoder: true` in `Subscription.yaml`

#### 5. Schema Mismatch: CAI vs. GET API Response
- **Symptom**: 
  - `At least 1 "feed_output_config" blocks are required.` (for Cloud Asset Feed)
  - `The argument "config" is required, but no definition was found.` (for Spanner Instance)
  - Step 1 (`CAI -> HCL`) passes, but the round-trip TF file is missing the required block/argument.
- **Cause**: The field is in the CAI asset data, but the JSON structure of the asset from CAI is different from the structure returned by the resource's GET API call (which the Terraform provider's encoder/flattener is often based on).
- **Solution**: Use `tgc_ignore_terraform_encoder: true` on the field in the YAML that has the structural mismatch. This allows custom TGC logic to handle the mapping.
- **Examples**:
  - **Cloud Asset FolderFeed**: `FolderFeed.yaml`. Error message: `At least 1 "feed_output_config" blocks are required.`. The Terraform encoder puts the resource data into the `feed` field, but CAI doesn’t expect the `feed` field. Use `tgc_ignore_terraform_encoder: true` on the field in the YAML to skip the shared encoder during `tfplan2cai` conversion.
  - **Spanner Instance**: `Instance.yaml`.

#### 6. Ignored Terraform Decoder Needed
- **Symptom**: Build error like `undefined: tpgresource.ParseImportId`
- **Cause**: `tpgresource.ParseImportId` is inside a standard Terraform decoder template. Also, the resource's default Terraform decoder (used to convert API responses to Terraform state) contains logic (e.g., `d.Get()`, parsing import IDs) that is not applicable or fails during the `CAI -> HCL` conversion, as there's no prior Terraform state.
- **Solution**: Add `tgc_ignore_terraform_decoder: true` to the resource or field in the YAML. This prevents TGC from using the standard Magic Modules decoder template for this item during `CAI -> HCL`.
- **Examples**:
  - **Instance**: `Instance.yaml`. Enable `tgc_ignore_terraform_decoder: true`. This prevents the standard Terraform decoder from running during TGC conversion. The standard decoder attempts to make API calls (to fetch `authString`), which causes "client is nil" errors in the TGC environment where no authorized client is available. It also avoids "project: required field is not set" errors that originated from the decoder's validation logic.

#### 7. Missing Fields in HCL from CAI - Custom Flatten Issues
- **Symptom**: Error message like `missing fields in resource ... after cai2hcl conversion: [field_name]` or a type panic during `cai2hcl` conversion (e.g., `interface conversion: interface {} is bool, not map[string]interface {}`).
  - **Example 1**: `TestAccFilestoreInstance_reservedIpRange_update_step1: missing fields in resource google_filestore_instance.instance after cai2hcl conversion: [networks.reserved_ip_range]`
  - **Example 2**: `TestAccApphubServiceProjectAttachment_serviceProjectAttachmentFullExample_step1: missing fields in resource google_apphub_service_project_attachment.example2 after cai2hcl conversion: [service_project]`
- **Cause**: The field is in the CAI asset data. It is populated in Terraform state by a `custom_flatten` function in Magic Modules. These functions often use `d.Get()`, which reads from the Terraform state. During `CAI -> HCL`, `d.Get()` returns zero values, so the field appears missing or has an incorrect zero value. This can also cause a panic if there's a type mismatch.
- **Debug**: Check if the field exists in the CAI asset data by inspecting the test's JSON file (e.g., `TestAccFilestoreInstance_reservedIpRange_update_step1.json`). If the field is present, examine any custom code associated with the resource that might be unsetting it or panicking.
- **Solution**: Add `tgc_ignore_terraform_custom_flatten: true` to the field's definition in the YAML. This tells TGC not to execute the Magic Module's `custom_flatten` function for this field during `CAI -> HCL` conversion, avoiding the panic or zero-value drop. If complex data transformation from CAI format to HCL is needed, use a `tgc_decoder`.
  - Add `tgc_decoder: 'templates/tgc_next/decoders/resource_name.go.tmpl'` under `custom_code` in the YAML. A `tgc_decoder` gives you direct access to safely read the raw `res map[string]interface{}` and write to the output `hclData map[string]interface{}` bypassing strict Terraform schemas and validations.
- **Examples**:
  - Filestore Instance: `Instance.yaml`
  - Apphub Service Project Attachment: `ServiceProjectAttachment.yaml`
  - Bigtable AppProfile: `AppProfile.yaml`

#### 8. Field Not Present in CAI Asset
- **Symptom**: Error message like `missing fields in resource ... after cai2hcl conversion: [field_name]`.
  - **Example**: `TestAccComputeBackendBucket_backendBucketSecurityPolicyExample_step1: missing fields in resource google_compute_backend_bucket.image_backend after cai2hcl conversion: [edge_security_policy]`
- **Cause**: The field is part of the Terraform provider's schema but is not included in the asset data provided by Cloud Asset Inventory.
- **Debug**: Check if the field exists in the CAI asset data by inspecting the test's JSON file. If the field is absent, it is not being provided in the CAI asset data.
- **Solution**: Mark the field with `is_missing_in_cai: true` in the resource YAML file. This informs TGC that the field is not expected to be in the CAI input.
- **Example**: `is_missing_in_cai: true` in the field `edgeSecurityPolicy` in `BackendBucket.yaml`.

#### 9. Incorrect CAI Asset Name Format
- **Symptom**: Error message like `A required argument (e.g., "folder") is missing in the generated HCL.`
- **Cause**: TGC parses the CAI asset's `name` field to extract resource identifiers. The default CAI asset name format is a computed field. The parameters fields cannot be extracted from the CAI asset name during `cai2hcl` and are not generated in the `TEST_export.tf` file.
- **Solution**: Specify the correct format using `cai_asset_name_format` in the resource's YAML. Use `{{field_name}}` placeholders.
- **Example**: 
  - **Failing test**: `TestAccCloudAssetFolderFeed_cloudAssetFolderFeedExample`
  - **Error message**: `The argument "folder" is required, but no definition was found.`
  - **Debug**: Check if `folder` is a parameter in `FolderFeed.yaml`. Inspect `cloudasset_folder_feed_cai2hcl.go` to check if the CAI asset name format is correct. TGC parses the CAI asset's `name` field, and the default format is `{{name}}`. The `folder` field is a parameter and cannot be extracted from the CAI asset name during `cai2hcl`.
  - **Solution**: Add `cai_asset_name_format` in `FolderFeed.yaml` to support field extraction from the asset name as `feed_id` is not computed: `cai_asset_name_format: '//cloudasset.googleapis.com/folders/{{folder}}/feeds/{{feed_id}}'`

#### 10. Add Handwritten Subtests
- **Symptom**: Message like `TestAccMonitoringAlertPolicy: test steps are unavailable`
- **Cause**: The top-level test `TestAccMonitoringAlertPolicy` doesn’t exist. Instead, subtests (e.g., `TestAccMonitoringAlertPolicy/basic`) exist.
- **Solution**: Add `tgc_tests` with the subtests to the resource YAML file.
- **Example**: [AlertPolicy.yaml reference](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/products/monitoring/AlertPolicy.yaml#L42)

#### 11. Undefined Function
- **Symptom**: Message like `undefined: json`
- **Cause**: The `json` package is not imported.
- **Solution**: Add `"encoding/json"` import path and add `_ = json.Unmarshal` to `mmv1/templates/tgc_next/services/resource.go.tmpl` to ensure the import is preserved.
- **Example**: [PR #16178 reference](https://github.com/GoogleCloudPlatform/magic-modules/pull/16178)

#### 12. Identifying Resources with Overlapping CAI Asset Types
- **Symptom**: A single CAI asset type maps to multiple distinct Terraform resources (e.g., `compute.googleapis.com/SslCertificate` mapping to `google_compute_ssl_certificate`, `google_compute_managed_ssl_certificate`, `google_compute_region_ssl_certificate`), and the automated HCL converter assigns the wrong converter.
- **Cause**: The TGC `cai2hcl` generator automatically attempts to identify resources by extracting distinct segments from their URL paths. If resources are differentiated by JSON data payload properties (e.g., `"type": "MANAGED"`) rather than distinct URL structures, the structural URL matching will fail to correctly route the asset.
- **Solution**: Inject a manual type-checking override block directly into the `ConvertResource` function within `mmv1/templates/tgc_next/cai2hcl/convert_resource.go.tmpl` to explicitly handle the problematic asset type.
- **Example**:
  - **Failing Resource**: `google_compute_managed_ssl_certificate`
  - **Error**: Assets are consistently mapped to the standard `google_compute_ssl_certificate` incorrectly.
  - **Solution**: Add hardcoded logic to `mmv1/templates/tgc_next/cai2hcl/convert_resource.go.tmpl` ahead of the automatic `IdentityParams` loop:
  ```go
		{{- if eq $resourceType "compute.googleapis.com/SslCertificate" }}
			if strings.Contains(asset.Name, "regions") {
				converter = ConverterMap[asset.Type]["ComputeRegionSslCertificate"]
			} else if typeVal, ok := asset.Resource.Data["type"]; ok && typeVal == "MANAGED" {
				converter = ConverterMap[asset.Type]["ComputeManagedSslCertificate"]
			} else {
				converter = ConverterMap[asset.Type]["ComputeSslCertificate"]
			}
		{{- else }}
  ```