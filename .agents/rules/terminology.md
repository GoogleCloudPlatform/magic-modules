---
trigger: always_on
description: Vocabulary definitions for standard actions in this repository
---

# Terminology Rules

## 1. Sync
- **Definition:** Aligning the Git histories of `magic-modules` and a downstream provider using the `sync-provider` skill. Setting up a clean baseline for verification.

## 2. Generate
- **Definition:** Running `make provider` from Magic Modules to generate code into the downstream repository.
- **When to use:** When you have modified local YAML/templates and need to see their output in Go code.

## 3. Build
- **Definition:** Running `make build` in a downstream provider to compile the Go binary. This could be ambiguous with "generate", so if you are not sure, clarifiy with the user.
- **When to use:** When handwriting Go changes or testing compilation.
