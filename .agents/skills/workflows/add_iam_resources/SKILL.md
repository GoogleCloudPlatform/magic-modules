---
name: add_iam_resources
description: "Creates IAM resources for an existing MMv1 resource that does not yet have IAM resources."
---

# `add_iam_resources`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Prerequisites
Your prompt will have included a resource to target, e.g., "Create IAM resources for google_container_analysis_note". The target resource must not already have IAM resources, and its API must support IAM.

## Execution Steps

### 1. Verification
For the target resource, confirm its location in the MMv1 codebase, and whether IAM support is already implemented.

* Find the resources metadata YAML file in `products/[product]/[resource].yaml`. The resource already supports IAM if there is a top-level `iam_policy` key.
* Use the `verify-api-schema` skill to check if the resource's API supports IAM by checking if it has `getIamPolicy` and `setIamPolicy` methods. If you can't find the resource's discovery document you might be able to use the API documentation instead. Look for a `references.api` field in the resource's metadata file, and then look up the API documentation for that service. If you still can't find the resource's IAM methods, ask the user to provide a link to the documentation and try again.

### 2. Create IAM Resources
Create the resources by adding the `iam_policy` block to the resource's metadata YAML file. Consult the documentation for the resource's API to determine the correct IAM policy configuration. The available fields for `iam_policy` are documented in `docs/content/reference/resource.md`.

* If the parent resource URL path contains multiple hierarchical path parameters (e.g., `catalogs/{{catalog}}/namespaces/{{namespace}}/tables/{{name}}`), ensure that `additional_parent_resource_attributes` is configured under `iam_policy` (e.g. `additional_parent_resource_attributes: [catalog, namespace]`) alongside `parent_resource_attribute: name` so that all parent key identifiers are properly captured in IAM bindings/policies.
* Note that `allowed_iam_role` is optional and defaults to `'roles/viewer'` if omitted. If configuring a service-specific role (e.g., `roles/apphub.admin`), consult public Google Cloud IAM documentation using `read_url_content` or `search_web`, or inspect existing product YAMLs. Do not use internal search tools.
* Do **not** set `iam_conditions_request_type` solely because `options.requestedPolicyVersion` appears in the Discovery Document—any service that embeds the standard `google.iam.v1.IAMPolicy` proto includes `options.requestedPolicyVersion` in its schema even when the backend rejects IAM Conditions (`version: 3`) on `setIamPolicy` with `400 INVALID_ARGUMENT`. Only configure `iam_conditions_request_type` if public IAM documentation confirms that the resource supports IAM Conditions and you verify that the generated `_withCondition` acceptance tests pass against the live API.

**IMPORTANT: You are only to modify the `iam_policy` block. Do not modify anything outside of this block, and do not modify any files other than the resource's metadata YAML file.** If you believe you need to make any modifications outside of this block, ask the user to confirm before proceeding. Offer to provide a detailed plan with the changes you are proposing to make.

### 3. Run Pre-Gen Checks
Use the `run-pre-gen-checks` skill to run pre-generation checks on the modified resource. If the checks fail, fix the issues and try again.

### 4. Generate Providers
Use the `generate-provider` skill to generate the provider with the new IAM resources. There should be two new code/test files in each of the generated providers:

* `[product]/iam_[product]_[resource].go`
* `[product]/iam_[product]_[resource]_generated_test.go`

Confirm that these files were created before proceeding to the next step.
* Note: MMv1 only generates `iam_[product]_[resource]_generated_test.go` if the resource has at least one sample with `exclude_test: false` (default). If all samples in the resource YAML have `exclude_test: true` and the user confirms removing `exclude_test: true` from a sample, MMv1 will also generate `[product]/resource_[product]_[resource]_generated_test.go`. When that happens, inspect why the sample test was previously excluded (e.g., soft-delete requiring `custom_code.test_check_destroy`, or `ignore_read_extra` on the sample step), ensure the generated sample test passes, and deduplicate any overlapping basic tests or custom `CheckDestroy` helpers in `mmv1/third_party/terraform/services/[product]/`.

### 5. Run Acceptance Tests
Use the `run-acctests` skill to run acceptance tests for the new IAM resources:
* List all `func TestAcc` functions across every newly created or modified `*_test.go` file (`grep '^func TestAcc' ...`) before running tests.
* Run **all** generated test functions in `iam_[product]_[resource]_generated_test.go` (using a prefix regex like `-run=TestAcc<Resource>Iam` without a trailing `$` anchor so `_withCondition` and `_withAndWithoutCondition` variants are never skipped), as well as any newly generated sample tests in `resource_[product]_[resource]_generated_test.go`.
* If the tests fail, fix the errors and try again using the `run-acctests` skill.
