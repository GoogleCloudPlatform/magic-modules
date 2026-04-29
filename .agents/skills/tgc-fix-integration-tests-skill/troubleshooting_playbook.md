# Troubleshooting Playbook

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

#### Tracing Failures Backwards (Recommended Approach)

When a test fails at the final plan or validation stage, do not guess the cause based on the error message alone. Systematically trace the failure backwards through the intermediate files:

1. **Inspect Step 3 Output (`Test_roundtrip.tf`)**: 
   * Check if the failing field is present and correct. 
   * If it is missing or invalid, the error was introduced in an earlier stage.

2. **Inspect Step 2 Output (`Test_roundtrip.json`)**: 
   * Check if the field exists in the `resource.data` block. 
   * If it is missing, either Step 2 (`tfplan2cai`) dropped it, or the input to Step 2 (the HCL from Step 1) was invalid/empty.

3. **Inspect Step 1 Output (`Test_export.tf`)**: 
   * Check if the field is present and has the expected value. 
   * If it is empty or missing here (while present in the source JSON), **Step 1 (`cai2hcl`) failed the mapping**. This is often caused by `custom_flatten` functions reading from unpopulated state.

4. **Inspect Original Input (`Test.json`)**: 
   * Verify if the data was actually present in the source CAI asset. 
   * If yes, and it was lost in Step 1, you have isolated the problem to the `cai2hcl` mapping (e.g., needing `tgc_ignore_terraform_custom_flatten`).

By isolating the exact file where the data disappears, you avoid dead-ends and immediately know whether to look at flatteners, encoders, or decoders.

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

- **Example 3**:
  - **Failing test**: `TestAccDatabaseMigrationServiceConnectionProfile_databaseMigrationServiceConnectionProfileCloudsqlExample`
  - **Error message**: `"host": required field is not set` (triggered by `required_with` constraints).
  - **Debug**: CAI asset does not provide `host`, `username`, `port`, or `password` when using `cloudSqlId` or `alloydbClusterId`. But Terraform schema requires them due to `required_with` constraints if any of them are set (or if we attempt to provide dummy values for some of them!).
  - **Solution**: Add `tgc_decoder` to the resource's YAML file. Implement CONDITIONAL logic to set missing required fields to `"unknown"` (for strings) or `0` (for integers like `port`) ONLY IF related fields like `alloydbClusterId` or `cloudSqlId` are present.
- **Key Takeaway**: When handling `required_with` constraints in decoders, ensure you provide values for ALL inter-dependent fields, using appropriate types (e.g., `0` for integers).

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
- **Alternative Symptom**: Validation error during `terraform plan` on the round-trip file (Step 3) such as `"field": one of ... must be specified`. This happens when the field belongs to a required group (like `exactly_one_of`) and was generated as empty in Step 1, causing it to be dropped entirely in Step 2 (`tfplan2cai`).
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
  - Colab NotebookExecution: `NotebookExecution.yaml`

#### 8. Field Not Present in CAI Asset
- **Symptom**: Error message like `missing fields in resource ... after cai2hcl conversion: [field_name]`.
  - **Example**: `TestAccComputeBackendBucket_backendBucketSecurityPolicyExample_step1: missing fields in resource google_compute_backend_bucket.image_backend after cai2hcl conversion: [edge_security_policy]`
- **Cause**: The field is part of the Terraform provider's schema but is not included in the asset data provided by Cloud Asset Inventory.
- **Debug**: Check if the field exists in the CAI asset data by inspecting the test's JSON file. If the field is absent, it is not being provided in the CAI asset data.
- **Solution**: Mark the field with `is_missing_in_cai: true` in the resource YAML file. This informs TGC that the field is not expected to be in the CAI input. Only add `is_missing_in_cai: true` if the field is missing in **ALL** of the resource's CAI asset JSON files. If the field exists in some CAI files but is missing in others, do not use this flag.
- **Example**: `is_missing_in_cai: true` for the field `edgeSecurityPolicy` in `BackendBucket.yaml`.

#### 9. Incorrect CAI Asset Name Format
- **Symptom**: Error message like `A required argument (e.g., "folder") is missing in the generated HCL.`
- **Cause**: TGC parses the CAI asset's `name` field to extract resource identifiers. The default CAI asset name format falls back to the resource ID format (derived from `self_link` or `base_url`). If `self_link` is just `{{name}}`, parameters cannot be extracted from the CAI asset name during `cai2hcl` and are not generated in the `TEST_export.tf` file.
- **Solution**: Specify the correct format using `cai_asset_name_format` in the resource's YAML if the default ID format does not contain the required parameters. Use `{{field_name}}` placeholders.
- **Example**: 
  - **Failing test**: `TestAccCloudAssetFolderFeed_cloudAssetFolderFeedExample`
  - **Error message**: `The argument "folder" is required, but no definition was found.`
  - **Debug**: Check if `folder` is a parameter in `FolderFeed.yaml`. Inspect `cloudasset_folder_feed_cai2hcl.go` to check if the CAI asset name format is correct. If `self_link` is just `{{name}}`, the default format is `{{name}}` and the `folder` field cannot be extracted.
  - **Solution**: Add `cai_asset_name_format` in `FolderFeed.yaml` to support field extraction from the asset name: `cai_asset_name_format: '//cloudasset.googleapis.com/folders/{{folder}}/feeds/{{feed_id}}'`

#### 10. Add Handwritten Subtests
- **Symptom**: Message like `TestAccMonitoringAlertPolicy: test steps are unavailable`
- **Cause**: The top-level test `TestAccMonitoringAlertPolicy` doesn’t exist. Instead, subtests (e.g., `TestAccMonitoringAlertPolicy/basic`) exist.
- **Solution**: Add `tgc_tests` with the subtests to the resource YAML file.
- **Example**: [AlertPolicy.yaml reference](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/products/monitoring/AlertPolicy.yaml#L42)

#### 11. Test Name Mismatch in Metadata or Tests Not Generated
- **Symptom**: The test name in `tests_metadata_*.json` does not match the newly generated test name, or the test is not generated at all (e.g., due to excluded examples).
- **Cause**: 
  - The resource was recently added or tests were refactored, causing the metadata to point to a test name that doesn't match the standard generated pattern.
  - Or, the only examples in the YAML are excluded (`exclude_test: true`), so no tests are generated by default.
    > [!CAUTION]
    > DO NOT remove `exclude_test: true` from examples to force generation. These examples are also used to generate tests for the standard Terraform provider, and removing it may break provider tests or cause unwanted test generation there.
  - **Special Case**: The test function in the handwritten file starts with a lowercase `t` (e.g., `testAcc...`) because it is a helper or subtest called by a top-level test. The generator regex `func (TestAcc[^(]+)` only matches uppercase `TestAcc`, so it fails to discover it.
- **Solution**: Explicitly list the expected test name (or subtest name like `TestAcc.../subtest`) in `tgc_tests` in the resource's YAML file. This forces the generator to include that specific test case in the generated test file.
- **Note on Missing Data Files**: If the `.tf` and `.json` files are missing in TGC because examples were excluded, running the test with `WRITE_FILES=true` will automatically generate these golden files in the TGC repository during execution.
- **Verification via Metadata Cache**: The name in `tgc_tests` must match the key used in the GCS metadata to successfully retrieve the test data. You can verify if a key is valid by searching for it in the cached metadata files in the `test_mata/` directory at the root of the downstream repository. For example: `grep -r "TestAccAccessApprovalSettings/folder" test_mata/`. If it finds matches like `"TestAccAccessApprovalSettings/folder": {`, then it's a valid key!
- **Example 1**: In `MessageBus.yaml` (forcing generation when examples are excluded):
  ```yaml
  tgc_tests:
    - name: 'TestAccEventarcMessageBus/basic'
  ```

- **Example 2**: In `FolderSettings.yaml` (handling lowercase `testAcc` helper function):
  The test function was named `testAccAccessApprovalFolderSettings` (lowercase `t`) because it was a helper function called by a parent test `TestAccAccessApprovalSettings` to ensure serial execution (due to the hierarchical nature of Access Approval settings).
  The solution was to explicitly list the test name (including the subtest) in the `tgc_tests` field in the YAML file:
  ```yaml
  tgc_tests:
    - name: 'TestAccAccessApprovalSettings/folder'
  ```
  This forced the generator to include it.

#### 12. Undefined Function
- **Symptom**: Message like `undefined: json`
- **Cause**: The `json` package is not imported.
- **Solution**: Add `"encoding/json"` import path and add `_ = json.Unmarshal` to `mmv1/templates/tgc_next/services/resource.go.tmpl` to ensure the import is preserved.
- **Example**: [PR #16178 reference](https://github.com/GoogleCloudPlatform/magic-modules/pull/16178)

#### 13. Identifying Resources with Overlapping CAI Asset Types
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

#### 14. CAI Resource Kind Mismatch
- **Symptom**: Error message stating a resource "is supported in either `tfplan2cai` or `cai2hcl` within tgc, but not in both."
- **Cause**: This happens due to a namespace collision in the `cai2hcl` generator. TGC builds a two-way routing map to know which Terraform resource corresponds to which CAI Asset Type. If two resources (e.g., `BackendService` and `RegionBackendService`) both generate the exact same CAI asset type path (e.g. `compute.googleapis.com/BackendService`), they will overwrite each other's mapping inside `ConverterMap`. This leaves one of the resources registered in `tfplan2cai` but completely orphaned in `cai2hcl`.
- **Solution**: Explicitly define a unique asset kind suffix using the `cai_resource_kind` parameter in the resource's `.yaml` file. This ensures `cai2hcl` registers a distinct dictionary key for mapping incoming CAI payloads back to the explicit Terraform resource variant.
- **Example**: In `RegionBackendService.yaml` or `GlobalForwardingRule.yaml`, specify the exact CAI suffix name to avoid colliding with `BackendService.yaml` or `ForwardingRule.yaml`:
  ```yaml
  include_in_tgc_next: true
  cai_resource_kind: 'GlobalForwardingRule' # or 'RegionBackendService'
  ```

#### 15. Undefined Functions in Shared Utilities
- **Symptom**: Functions like `expandToLongForm` or `expandToRegionalLongForm` are reported as undefined during the TGC build.
- **Cause**: The file defining these shared utilities (e.g., `mmv1/third_party/terraform/services/eventarc/eventarc_utils.go`) is not being compiled into the TGC binary.
- **Solution**: Add the file to the CopyCommonFiles or CompileCommonFiles list in `mmv1/provider/terraform_tgc_next.go` (or equivalent provider file) to ensure it is included in the build.

#### 16. No Tests Generated Due to Excluded Examples
- **Symptom**: No integration tests are generated for a resource, leading to messages like "SKIPPED: No tests generated".
- **Cause**: All examples defined for the resource in the Magic Modules YAML file are marked with `exclude_test: true`.
- **Solution**: This is expected if the resource cannot be tested automatically in the shared environment. If tests are desired, add a new example that does not have `exclude_test: true` and supports setup/teardown in tests.
- **Example**: `eventarc/Enrollment` was skipped because its only example was marked with `exclude_test: true`. In this case, verified tests cannot be generated automatically.
