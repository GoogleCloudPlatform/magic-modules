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
* `--only-migration` (Optional): Run only the migration steps (examples -> samples conversion, copy and migrate templates). Do not sort keys or format string quotes. Cannot be combined with `--only-format`.
* `--only-format` (Optional): Run only the formatting steps (sort keys, strip string quotes). Do not migrate examples to samples or copy templates. Cannot be combined with `--only-migration`.
* `--explicit-config-path` (Optional): Force writing `config_path` in step definition mappings during migration, even if they match the default generated template path or were not explicitly defined in the original YAML file.
* `--eap` (Optional): Enable EAP private overrides repository migration. Forces EAP folder layout resolution and forces `--explicit-config-path` to be true.
* `--skip-open-pr` (Optional): Skip files modified by active open PRs updated in the last N days (configured by `--skip-open-pr-days`).
* `--skip-open-pr-days <days>` (Optional): Number of days of open PR history to verify when checking open PRs (defaults to `60`).

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
To migrate a single public product YAML file safely only if there is no active PR updated in the last 30 days touching it:
```bash
go run main.go convert-resource-template \
  --skip-open-pr \
  --skip-open-pr-days 30 \
  -f mmv1/products/hypercomputecluster/Cluster.yaml \
  <path_to_magic_modules_repository>
```
* **Output**:
  ```
  Fetching open PRs updated in the last 30 days from GitHub...
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

### Example 6: Run Formatting Only (No Migration)
To re-format, sort keys, and strip quotes on one or more resource files without migrating examples to samples:
```bash
go run main.go convert-resource-template \
  --only-format \
  -f mmv1/products/dns/ManagedZone.yaml \
  <path_to_magic_modules_repository>
```

### Example 7: Run Migration Only (No Formatting/Sorting)
To migrate templates and configurations to samples, keeping the original YAML key ordering and quote formatting:
```bash
go run main.go convert-resource-template \
  --only-migration \
  -f mmv1/products/dns/ManagedZone.yaml \
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
