# Resource Template Converter Tool

The `resource-template-converter` is designed to convert Magic Modules resource YAML configurations (from the legacy `examples:` structure to the modern `samples:` block) and migrate associated Terraform templates (`.tf.tmpl`).
It also reformats product YAML files.

---

## Quick Start

### 1. Build the Tool
From the root of the `magic-modules` repository, run:
```bash
cd tools/resource-template-converter
go build -o bin/convert-resource-template main.go
```

---

## Usage Instructions

```bash
./bin/convert-resource-template convert-resource-template [flags] <path_to_repository_root>
```

### Positional Arguments
1. `<path_to_repository_root>` (Required): Path to the root of the repository you want to process (e.g., `<path_to_magic_modules_repository>`).

### Flags
* `-f, --file <path>` (Optional): Path to a single resource YAML file to convert (relative or absolute). If omitted, the tool walks and processes all YAML files in the `products/` directory.
* `--skip-open-pr` (Optional): Fetch open PRs updated in the last 2 months from GitHub. Any matching YAML files modified in those PRs will be skipped.

---

## Real-World Examples

### Example 1: Single File Migration
To migrate a single resource YAML file:
```bash
./bin/convert-resource-template convert-resource-template \
  -f mmv1/products/vertexai/Dataset.yaml \
  <path_to_magic_modules_repository>
```

### Single File Migration with PR Verification
To migrate a single public product YAML file safely if active PR updated in the last 2 months touches this file:
```bash
./bin/convert-resource-template convert-resource-template \
  --skip-open-pr \
  -f mmv1/products/hypercomputecluster/Cluster.yaml \
  <path_to_magic_modules_repository>
```
* **Output**:
  ```
  Fetching open PRs updated in the last 2 months from GitHub...
  Successfully fetched open PRs. Found 835 modified files in open PRs.
  Skipping single target file mmv1/products/hypercomputecluster/Cluster.yaml: modified in active open PR(s) [17678 17610 17522]
  ```

### Example 3: Bulk Repository Conversion (Excluding Active PR Files)
To safely bulk-migrate all product files in the repository:
```bash
./bin/convert-resource-template convert-resource-template \
  --skip-open-pr \
  <path_to_magic_modules_repository>
```

---

## Architecture and Design

The tool is designed as a modular Go command line utility:
```
tools/resource-template-converter/
├── main.go                 # CLI entry point
├── cmd/
│   └── convert_resource_template.go   # Cobra command definitions & walking flow
├── copy/
│   └── copy.go             # TF template migration & replacement of $.Vars -> $.ResourceIdVars
├── github/
│   ├── github.go           # GitHub CLI API querying and path normalization
│   └── github_test.go      # Unit tests for path normalizations
└── migrate/
    ├── migrate.go          # Product YAML loader and AST-based transformer
    ├── format.go           # YAML formatting, node-sorting, and reflection helpers
    └── migrate_test.go     # Unit tests for YAML converters
```
