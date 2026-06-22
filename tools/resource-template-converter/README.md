# Resource Template Converter Tool

The `resource-template-converter` is designed to convert Magic Modules resource YAML configurations (from the legacy `examples:` structure to the modern `samples:` block) and migrate associated Terraform templates (`.tf.tmpl`).
It also reformats product YAML files.

---

## Quick Start & Execution Methods

You can run the tool directly without building it (recommended), or build it as a binary.

### Method A: Direct Execution (Recommended)
No compilation step is required. Run the tool directly using `go run` from the `tools/resource-template-converter` directory:
```bash
cd tools/resource-template-converter
go run main.go convert-resource-template [flags] <path_to_repository_root>
```

### Method B: Building a Binary
If you prefer to build a standalone binary executable:
```bash
cd tools/resource-template-converter
go build -o bin/convert-resource-template main.go

# Run the built binary:
./bin/convert-resource-template convert-resource-template [flags] <path_to_repository_root>
```

---

## Usage Instructions

### Positional Arguments
1. `<path_to_repository_root>` (Required): Path to the root of the repository you want to process (e.g., `<path_to_magic_modules_repository>`).

### Flags
* `-f, --file <paths>` (Optional): Comma-separated list of resource YAML file paths to convert (relative or absolute). If omitted, the tool walks and processes all YAML files in the `products/` directory.
* `-p, --product <product_names>` (Optional): Comma-separated list of product directories to convert (e.g., `vertexai` or `vertexai,pubsublite`). If specified, only YAML files under the matching product directories are walked and processed. Cannot be specified together with `--file`.
* `-F, --skip-file <paths>` (Optional): Comma-separated list of resource YAML file paths to skip from migration.
* `-P, --skip-product <product_names>` (Optional): Comma-separated list of product directories to skip from migration.
* `--skip-open-pr` (Optional): Fetch open PRs updated in the last 2 months from GitHub. Any matching YAML files modified in those PRs will be skipped.

---

## Real-World Examples (Using `go run`)

All examples below are executed from the `tools/resource-template-converter` directory.

### Example 1: File Migration (Single or Multiple)
To migrate one or more specific resource YAML files:
```bash
go run main.go convert-resource-template \
  -f mmv1/products/vertexai/Dataset.yaml,mmv1/products/pubsublite/Topic.yaml \
  <path_to_magic_modules_repository>
```

### Example 2: Single File Migration with PR Verification
To migrate a single public product YAML file safely only if there is no active PR updated in the last 2 months touching it:
```bash
go run main.go convert-resource-template \
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

### Example 3: Product Directory Migration (Single or Multiple)
To migrate all resource YAML files in one or more specific product directories (e.g., `pubsublite,essentialcontacts`):
```bash
go run main.go convert-resource-template \
  -p pubsublite,essentialcontacts \
  <path_to_magic_modules_repository>
```

### Example 4: Bulk Repository Conversion (Excluding Active PR Files)
To safely bulk-migrate all product files in the repository:
```bash
go run main.go convert-resource-template \
  --skip-open-pr \
  <path_to_magic_modules_repository>
```

### Example 5: Bulk Conversion with Skip Filters
To safely bulk-migrate the entire repository but exclude specific products (e.g. `compute`) and specific resource files (e.g. `mmv1/products/dns/ManagedZone.yaml`):
```bash
go run main.go convert-resource-template \
  -P compute \
  -F mmv1/products/dns/ManagedZone.yaml \
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
