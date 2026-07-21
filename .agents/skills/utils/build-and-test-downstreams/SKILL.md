---
name: build-and-test-downstreams
description: "Generate downstreams (GA, Beta, and docs-examples) and run unit tests, linting, and docscheck. Use this skill when the user wants to check if downstreams pass tests or to check for errors during generation."
---

# `build-and-test-downstreams`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.
> This skill generates downstream providers into isolated scratch directories and runs compilation (`go build`), unit tests (`make testnolint`), linting (`make lint`), and documentation checks (`make docscheck`) concurrently.

## Prerequisites

- You must be operating in the `magic-modules` root directory.
- `go` (v1.26+) and `git` must be installed and configured.

## Execution Steps

### 1. Run Downstream Build and Unit Test Runner

Execute the script to generate and test all downstream target providers in parallel:

```bash
./.agents/skills/utils/build-and-test-downstreams/scripts/build_and_test_downstreams.sh
```

### 2. Output and Verification

- The script caches upstream downstream provider repos in `scratch/provider-cache-*`.
- Downstream provider builds and unit tests execute concurrently in `scratch/downstream-build/downstream-*`.
- Console output streams the build and test logs for GA, Beta, and docs-examples targets.
- If any build or unit test step fails, the script exits with status `1`.
