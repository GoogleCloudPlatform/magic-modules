---
name: update-dependencies
description: "Update or add Go module dependencies for the downstream Terraform providers (terraform-provider-google and terraform-provider-google-beta) and sync go.mod and go.sum back to mmv1/third_party/terraform/."
---

# `update-dependencies`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Overview & Golden Rule

The Magic Modules repository (`mmv1/third_party/terraform/`) does **not** contain a complete Go module: resource code is generated from `mmv1/products/` YAML schemas at build time, and many handwritten files are `.go.tmpl` templates.

**NEVER run `go get` or `go mod tidy` directly inside `mmv1/third_party/terraform/`.** Doing so misses imports that only exist in `.go.tmpl` templates or generated `.go` files, stripping needed packages or demoting direct dependencies to `// indirect` and breaking downstream `go mod tidy` / `go mod vendor`.

Instead, always run Go module commands inside a generated downstream provider (`terraform-provider-google` or `terraform-provider-google-beta`) and copy the resulting `go.mod` and `go.sum` back to `mmv1/third_party/terraform/`. See [docs/content/develop/update-dependencies.md](../../../../docs/content/develop/update-dependencies.md).

## Prerequisites

* You must have access to the `magic-modules` repository and a downstream provider repository (`terraform-provider-google` or `terraform-provider-google-beta`).
* `$GOPATH` must be set in your shell environment (`export GOPATH="${GOPATH:-$(go env GOPATH)}"`).
* If using an isolated `git worktree` as `OUTPUT_PATH`, the worktree directory basename **must** be `terraform-provider-google` or `terraform-provider-google-beta` so `mmv1`'s output-folder and copyright-header checks match standard provider checkouts.

## Execution Steps

### 1. Determine Update Scope

Unless the user explicitly specifies a narrower target (such as bumping a single package or running `go mod tidy` only), **default to updating all module and test dependencies (`go get -u ./... && go mod tidy`)**:
* **Default — All module & test dependencies**: `go get -u ./...` + `go mod tidy`
* **Targeted package bump** (when a specific package is requested, e.g., `google.golang.org/api`): `go get <package>@<version>` (or `@latest`) + `go mod tidy`
* **Tidy only** (when requested): `go mod tidy`

### 2. Generate the Downstream Provider and Update Dependencies

```bash
export GOPATH="${GOPATH:-$(go env GOPATH)}"
export PATH="$GOPATH/bin:$PATH"

MAGIC_MODULES_PATH="$(pwd)"
PROVIDER_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google"

# 1. Ensure downstream checkout is clean and synced to upstream/main
cd "$PROVIDER_PATH"
git status --porcelain

# 2. Generate the provider from Magic Modules into the downstream checkout
cd "$MAGIC_MODULES_PATH"
make provider VERSION=ga OUTPUT_PATH="$PROVIDER_PATH"

# 3. Update dependencies inside the generated downstream provider
cd "$PROVIDER_PATH"
# Default: Update all module and test dependencies
go get -u ./...
# Or, if a specific package was requested:
# go get google.golang.org/api@latest

go mod tidy
```

### 3. Sync `go.mod` and `go.sum` Back to Magic Modules

Copy the updated `go.mod` and `go.sum` files back to `mmv1/third_party/terraform/`:

```bash
cp "$PROVIDER_PATH/go.mod" "$MAGIC_MODULES_PATH/mmv1/third_party/terraform/go.mod"
cp "$PROVIDER_PATH/go.sum" "$MAGIC_MODULES_PATH/mmv1/third_party/terraform/go.sum"
```

> **Important (Beta Worktree Caveat):** `mmv1/third_party/terraform/go.mod` **must** declare `module github.com/hashicorp/terraform-provider-google` on line 1. During `VERSION=beta` generation, `mmv1` rewrites the module declaration to `github.com/hashicorp/terraform-provider-google-beta` and will fatal if the source `mmv1/third_party/terraform/go.mod` already contains `terraform-provider-google-beta`. If you ran `go get` / `go mod tidy` inside `terraform-provider-google-beta`, rewrite line 1 before copying:
> ```bash
> sed '1s|module github.com/hashicorp/terraform-provider-google-beta|module github.com/hashicorp/terraform-provider-google|' \
>   "$PROVIDER_PATH/go.mod" > "$MAGIC_MODULES_PATH/mmv1/third_party/terraform/go.mod"
> ```

### 4. Verification & Handoff

1. **Re-generate and verify module integrity** in both `ga` and `beta` (or the target provider):
   ```bash
   cd "$MAGIC_MODULES_PATH"
   make provider VERSION=ga OUTPUT_PATH="$PROVIDER_PATH"

   cd "$PROVIDER_PATH"
   go mod tidy
   git diff --exit-code go.mod go.sum
   go mod vendor && rm -rf vendor
   ```
2. **Verify compilation and unit tests** in the downstream provider:
   ```bash
   cd "$PROVIDER_PATH"
   go build -o /dev/null .
   make lint
   make testnolint
   ```
3. **Separate Pull Request:** Per [docs/content/develop/update-dependencies.md](../../../../docs/content/develop/update-dependencies.md), submit dependency upgrades in a standalone PR containing *only* the `go.mod` and `go.sum` changes (plus any code adjustments required by the upgrade) with a `release-note:none` block.
