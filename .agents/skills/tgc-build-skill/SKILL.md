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

If you have modified YAML configurations or Go templates inside of `<magic-modules-path>/mmv1`, you must generate those changes into the downstream repository. Always run `make clean-tgc` first to remove any lingering downstream ghost files.

1. Clean the downstream target repository:
   ```bash
   make clean-tgc OUTPUT_PATH="$GOPATH/src/github.com/GoogleCloudPlatform/terraform-google-conversion"
   ```

2. Build to the downstream repository by explicitly providing the `OUTPUT_PATH`.

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

---

## Automation Script

You can execute the script `build_tgc.sh` to automate both the code generation and compilation phases. This is the recommended and fastest way to ensure your codebase is fully synchronized.

1. Ensure your current working directory is the root of the `magic-modules` repository.
2. Execute the script from the `.agents/skills` directory:
   ```bash
   ./.agents/skills/tgc-build-skill/scripts/build_tgc.sh
   ```

This will automatically:
- Run `make tgc` targeting the standard `$GOPATH` convention path for `terraform-google-conversion`.
- Navigate into the respective directory.
- Effectively run `make mod-clean && make build`.
