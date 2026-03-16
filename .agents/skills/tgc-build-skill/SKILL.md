---
name: tgc-build-skill
description: Build TGC from magic modules. Use when you need to build TGC from magic modules.
---

# tgc-build-skill

When you need to build TGC from magic modules, use this skill.

## When to Use This Skill

- Use this when building TGC from Magic Modules.
- This is helpful when you need to generate downstream code and compile the TGC binary after making configuration or template changes.

---

## How to Use It

Follow these two phases carefully to build TGC accurately.

### Phase 1: Generating Code from Magic Modules (`make tgc`)

If you have modified YAML configurations or Go templates inside of `<magic-modules-path>/mmv1`, you must generate those changes into the downstream repository.

1. Build to the downstream repository by explicitly providing the `OUTPUT_PATH`.

   **Example**:
   ```bash
   make tgc OUTPUT_PATH="$GOPATH/src/github.com/GoogleCloudPlatform/terraform-google-conversion"
   ```

### Phase 2: Building the TGC Binary (`make build`)

1. Navigate to the downstream provider repository you are working on:
   ```bash
   cd "$GOPATH/src/github.com/GoogleCloudPlatform/terraform-google-conversion"
   # or `cd $OUTPUT_PATH` 
   ```
2. Compile the binary:
   ```bash
   make mod-clean && make build
   ```
3. **DO NOT** run `go mod tidy && make build`. It is guaranteed to fail in this context.