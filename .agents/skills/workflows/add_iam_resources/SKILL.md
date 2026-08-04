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

**IMPORTANT: You are only to modify the `iam_policy` block. Do not modify anything outside of this block, and do not modify any files other than the resource's metadata YAML file.** If you believe you need to make any modifications outside of this block, ask the user to confirm before proceeding. Offer to provide a detailed plan with the changes you are proposing to make.

### 3. Run Pre-Gen Checks
Use the `run-pre-gen-checks` skill to run pre-generation checks on the modified resource. If the checks fail, fix the issues and try again.

### 4. Generate Providers
Use the `generate-provider` skill to generate the provider with the new IAM resources. There should be two new files in each of the generated providers:

* `[product]/iam_[product]_[resource].go`
* `[product]/iam_[product]_[resource]_test.go`

Confirm that these files were created before proceeding to the next step.

### 5. Run Acceptance Tests
Use the `run-acctests` skill to run acceptance tests for the new IAM resources. You will only need to run the tests from the newly created IAM test files, not all of the test files for the provider. If the tests fail, you will need to fix the errors and try again, using the `run-acctests` skill again to run the tests.
