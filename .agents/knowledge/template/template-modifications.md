---
name: template-modifications
description: Safety rules and blast radius precautions for modifying Magic Modules engine templates (mmv1/templates/).
topics: [template]
task_types: [bug-fix, field-add, new-resource]
source: authored
status: draft
last_verified: 2026-08-06
---

# Template Modifications & Blast Radius

Templates in `mmv1/templates/terraform/` control code generation globally across hundreds of resources in both `terraform-provider-google` and `terraform-provider-google-beta`. Modifications to these files have provider-wide impact.

## Rules for Template Edits

### 1. Human-in-the-Loop Approval Before Editing Templates
* **Rule:** If a fix or feature requires altering an engine template rather than product YAML files, halt and obtain explicit user approval before making the change.
* **Why:** Engine templates alter generated code across multiple unrelated services. A change intended for one resource can unintentionally break or alter dozens of others.

### 2. Mandatory Provider-Wide Generation & Build Verification
* **Rule:** Whenever modifying any file in `mmv1/templates/`, you must run full downstream generation (`make provider`), full provider compilation (`make build`), and inspect `git diff --stat` across downstream services.
* **Why:** Running acceptance tests or compiling an individual package only validates that specific service. Broken code or syntax regressions in unrelated services will only be caught by generating and building the entire provider binary.

## Do NOT Use For
* Resource-specific customizations that can be achieved via standard custom code hooks (`custom_code: ...` in product YAMLs such as `pre_create`, `pre_update`, `encoder`, `flattener`). Only modify core templates when the underlying generation engine itself lacks support for the required pattern.
