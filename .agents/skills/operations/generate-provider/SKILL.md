---
name: generate-provider
description: "Executes code generation from Magic Modules into downstream Terraform providers and optionally compiles the binary."
---

# `generate-provider`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Prerequisites
* You must be in the `magic-modules` root directory.
* The `$GOPATH` environment variable must be set correctly.
* You must have verified there are no unsaved or uncommitted changes in the downstream provider directory that would be overwritten by code generation.

## Execution Steps

Set your variables and run the commands in sequence.

```bash
# 1. Define Variables
VERSION="beta" # or "ga"
PROVIDER_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta" # or "terraform-provider-google" for GA

# Verify we are in /Users/camthornton/magic-modules
pwd

# Verify no uncommitted changes exist in the downstream provider directory. 
# Stop if there are uncommitted changes and alert the user.
cd "$PROVIDER_PATH"
git status --porcelain

# 2. Code Generation (Run from Magic Modules root)
cd /Users/camthornton/magic-modules
make provider VERSION=$VERSION OUTPUT_PATH="$PROVIDER_PATH"

# 3. Verify Changes (Run from Downstream Provider root)
cd "$PROVIDER_PATH"

# Verify that git status precisely matches the scope of your local Magic Modules edits.
# If you see "surprise" diffs in unrelated resources, ask the user if you should run `sync-provider` first.
git status --porcelain

# 4. Compile Binary (Already in Downstream Provider root)
make build
```

## Handoff & Guardrails

-   Verify successful compilation.
-   **CRITICAL SAFEGUARD:** Stop here! Do NOT automatically proceed to running tests unless the user explicitly requested you to do so in their initial prompt. If testing was not requested, ask the user for their next steps.
