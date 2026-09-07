---
name: add-fields-workflow
description: "Workflow for adding new fields to existing MMv1 or handwritten resources."
---

# `add-fields-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.

## Prerequisites

- You must be operating in the `magic-modules` root directory.
- You must know the target resource (e.g., `google_compute_instance`) and the API field(s) to add.

---

## Execution Steps

### 1. Context & Schema Investigation

Before you begin implementation, read the following documents:

- `docs/content/develop/add-fields.md` — Complete procedural guide for adding fields to MMv1 and handwritten resources.
- `docs/content/reference/field.md` — Comprehensive MMv1 field configuration reference.
- `docs/content/test/test.md` - Complete procedural guide for adding tests.
- `docs/content/reference/sample.md` - Comprehensive MMv1 sample configuration reference.
- `.agents/knowledge/index.md` for general information.

**Resource Definition Locations**:

- MMv1 generated resources: `mmv1/products/<product>/<resource>.yaml`
- Handwritten resources: `mmv1/third_party/terraform/services/<product>/resource_<product>_<resource>.go`

Implement logic according to the API behavior and terraform best practices. Your code should match the style of the resource and product.

### 2. Add fields

Add the requested field(s) to the requested resource(s).

### 3. Add tests for the fields

Field tests should exist according to the following rules:

- Every field must be present in at least one test step.
- Every mutable field must have its value altered in an update step.
- It's preferable for optional fields to be missing in at least one step to ensure that they're truly optional.

Because you're adding fields to an existing resource, try to modify existing tests if possible rather than adding new tests.

**Test Locations**:

- Generated : `mmv1/templates/terraform/samples/services/<PRODUCT>/`
- Handwritten resources: `mmv1/third_party/terraform/services/<PRODUCT>/resource_<product>_<resource>_test.go`

### 4. Run Pre-Gen Checks

- Use [run-pre-gen-checks](.agents/skills/utils/run-pre-gen-checks/SKILL.md). If issues are found, analyze and fix them.

### 5. Generate Provider

- Use [generate-provider](.agents/skills/operations/generate-provider/SKILL.md).
- Confirm that the provider was generated successfully.

### 6. Test and Debug

- Use [`repo-sync`](`.agents/skills/operations/repo-sync/`) to ensure the downstream repositories are in sync with magic-modules.
- Use [`generate-provider`](.agents/skills/operations/generate-provider/`) to generate the provider code into the downstream repositories.
- Invoke [`qa-test-runner`](.agents/skills/workflows/qa-test-runner/) to run all acceptance tests for the modified resource.
- Invoke [`test-fixer`](.agents/skills/operations/test-fixer/`) to fix any issues found by the `qa-test-runner` skill. Return to step 4 after any changes.

---

## The Loop

Repeat steps 4-6 as needed until the primary task is complete.
