---
name: deprecate-resource-or-field-workflow
description: "Workflow for deprecating existing resources or fields on the main branch prior to a major release."
---

# `deprecate-resource-or-field-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.

This workflow governs adding deprecation notices and establishing forwards-compatibility on the `main` branch prior to a major release.

## Prerequisites

- You must be operating on the `main` branch in the `magic-modules` root directory.
- You must know the target resource (e.g., `google_data_loss_prevention_job_trigger`) and the field (or entire resource) to deprecate.
- You must know the replacement path (if applicable) or the reason for deprecation.

---

## Execution Steps

### 1. Context & Guidance

Before beginning implementation, consult:
- `docs/content/breaking-changes/make-a-breaking-change.md` — Section *Add deprecations and warnings to the main branch*.
- `docs/content/reference/field.md` and `docs/content/reference/resource.md` — Configuration references.
- `.agents/knowledge/index.md` — General knowledge index.

### 2. Forwards-Compatibility & Replacement

- **Required Fields**: If a required field is being deprecated for future removal, make it optional (`required: true` $\rightarrow$ `optional: true` in MMv1 YAML or `Required: false, Optional: true` in Go) so configurations can begin omitting it prior to the major release.
- **Renames / Replacements**: Ensure the replacement field or resource is implemented, tested, and available on `main` before or alongside the deprecation.

### 3. Apply Deprecation Notice

Follow the conventions in `docs/content/breaking-changes/make-a-breaking-change.md`:
- **MMv1**: Set `deprecation_message` on the target property or resource in `mmv1/products/<product>/<Resource>.yaml`.
- **Handwritten**: Set `Deprecated` (fields) or `DeprecationMessage` (resources) in Go schema (`mmv1/third_party/terraform/services/<product>/...`) and update the handwritten documentation markdown file.

### 4. Run Pre-Gen Checks

- Execute [run-pre-gen-checks](.agents/skills/utils/run-pre-gen-checks/SKILL.md) to verify Go formatting, YAML linting, template validation, and unit tests.

### 5. Generate Provider & Verify Build

- Execute [generate-provider](.agents/skills/operations/generate-provider/SKILL.md).
- Run `make build` in the downstream provider repository to verify compilation.

### 6. Verification Testing

- Invoke [qa-test-runner](.agents/skills/operations/qa-test-runner/SKILL.md) to verify acceptance tests pass (`PASS`).

### 7. PR Creation & Release Note

Open a PR targeting `main`:
- Title: `<product>: deprecate <field_name|resource_name>`
- Body:
  ```markdown
  ```release-note:deprecation
  <product>: deprecated `<field_name>` on `google_<resource_name>`. Use `<replacement>` instead.
  ```
  *(or `<product>: deprecated `google_<resource_name>` resource. Use `google_<replacement>` instead.`)*
  ```
