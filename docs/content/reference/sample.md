---
title: "MMv1 sample reference"
weight: 25
---

# MMv1 sample reference

A sample is a collection of one or more `steps`, where each step represents a distinct Terraform configuration and test step (e.g., create, update).

Each sample supports the following attributes at the top level, with more granular control inside each step.

---

## Sample Attributes (Top-Level)

These attributes are defined once for an entire sample.

* `name`: `snake_case` name for the overall sample. This is used for generating test names.
* `primary_resource_id`: The ID of the main resource under test for the entire sample. Tests use this to run additional checks automatically.
* `primary_resource_type`: Optional resource type override of the primary resource. Used for import assertion validations.
* `bootstrap_iam`: Specify member/role pairs that should exist before the test runs. This avoids race conditions on global IAM permissions. `{project_number}` and `{organization_id}` are replaced automatically.
* `min_version`: Sets a minimum provider version for the entire sample (e.g., `beta`). This can be overridden by the `min_version` attribute within a specific step.
* `exclude_test`: If `true`, no tests are generated for this entire sample.
* `exclude_basic_doc`: If `true`, excludes the first step of this sample from the generated documentation. By default, the first step is automatically included as a use case in the documentation. Use this if you want to skip it.
* `skip_vcr`: If `true`, skips VCR testing for the entire sample.
* `skip_test`: If not empty, the entire sample is skipped during tests. The value should be a link to a ticket explaining why.
* `skip_func`: A custom function call to run to determine if tests should be skipped.
* `region_override`: Overrides location/region identifiers specifically inside IAM assertion checks.
* `external_providers`: A list of external providers (e.g., `random`, `time`) needed for the sample.
* `tgc_skip_test`: Skips generated conversion tests specifically running inside the TGC (Terraform Google Conversion) suite (value should be a ticket link reason).

---

## Step Attributes

A sample contains a list of one or more `steps`. Each step has its own configuration and test-specific attributes.

* `name`: `snake_case` name of the individual step. This is used for generating test configuration function names and documentation headers.
* `config_path`: The path to the step's configuration file. If omitted, it defaults to `templates/terraform/samples/services/{{product}}/{{step_name}}.tf.tmpl`.
* `resource_id_vars`: Key/value pairs to inject into the configuration file. Reference them with `{{index $.ResourceIdVars "key"}}`. Values here automatically receive a `tf-test` (or `tf_test`) prefix and a random suffix to ensure they are picked up by resource sweepers for cleanup and to avoid collisions. **Note:** The value must contain at least one `-` or `_` to trigger this automatic prefixing. If a resource name does not allow hyphens, use an underscore (e.g., `my_resource`) to generate a `tf_test` prefix instead.
* `vars`: Key/value pairs that are copied directly to tests without a prefix. Reference with `{{index $.Vars "key"}}`. **Note:** This should ONLY be used for fields that vary between steps (e.g., to test update functionality). Constant values should be hardcoded directly in the `.tf.tmpl` file.
* `test_env_vars`: Key/value pairs that map variable names to environment variables for tests (e.g., `PROJECT_NAME`, `REGION`, `ORG_ID`).
* `test_vars_overrides`: Key/value pairs to override variables with literal values or function calls specifically for tests.
* `oics_vars_overrides`: Key/value pairs to override variables with literal values specifically for Open in Cloud Shell (OiCS) tutorial generation.
* `min_version`: Overrides the sample-level `min_version` for this specific step.
* `ignore_read_extra`: A list of properties to ignore during the import test for this step, typically for write-only fields.
* `exclude_import_test`: If `true`, no import test is generated for this specific step.
* `include_step_doc`: If `true`, forces this specific step to be included in the generated documentation. By default, only the first step of a sample is included in the documentation as a use case. Use this on later steps to showcase update scenarios or complex configurations. This will override a top-level `exclude_basic_doc` setting if applied to the first step.

---

## Example

```yaml
samples:
  - name: service_resource_update
    primary_resource_id: example
    bootstrap_iam:
      - member: "serviceAccount:service-{project_number}@gcp-sa-healthcare.iam.gserviceaccount.com"
        role: "roles/bigquery.dataEditor"
      - member: "serviceAccount:service-org-{organization_id}@gcp-sa-osconfig.iam.gserviceaccount.com"
        role: "roles/osconfig.serviceAgent"
    min_version: "beta"
    skip_vcr: true
    external_providers:
      - "time"
    steps:
      - name: service_resource_minimal # Matches templates/terraform/samples/services/{{product}}/service_resource_minimal.tf.tmpl
        vars: # Varies between steps to test update functionality
          description: "A minimal description" 
        resource_id_vars: # Used for resource id in the configuration file
          dataset_id: "my-dataset"
          network_name: "my-network"
        test_env_vars:
          org_id: "ORG_ID"
        test_vars_overrides:
          network_name: 'acctest.BootstrapSharedServiceNetworkingConnection(t, "service-resource-network-config")'
        ignore_read_extra:
          - 'foo'
        exclude_import_test: true
      - name: service_resource_update # Matches templates/terraform/samples/services/{{product}}/service_resource_update.tf.tmpl
        vars:
          description: "An updated description" # This value updates the description field
        resource_id_vars:
          dataset_id: "my-dataset"
          network_name: "my-network"
```
