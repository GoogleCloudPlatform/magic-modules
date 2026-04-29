---
name: tgc-build-skill
description: Build TGC from magic modules. Use when you need to build TGC from magic modules.
---

# tgc-build-skill

When you need to build TGC from magic modules, use this skill.

## When to Use This Skill

- Use this when building TGC from Magic Modules.
- This is helpful when you need to generate downstream code and compile the TGC binary after making configuration or template changes.

You must execute the script `build_tgc.sh` to automate both the code generation and compilation phases. This is the recommended and fastest way to ensure your codebase is fully synchronized.

1. Ensure your current working directory is the root of the `magic-modules` repository.
2. Execute the script from the `.agents/skills` directory:
   ```bash
   ./.agents/skills/tgc-build-skill/scripts/build_tgc.sh
   ```

This will automatically:
- Run `make tgc` targeting the standard `$GOPATH` convention path for `terraform-google-conversion`.
- Navigate into the respective directory.
- Effectively run `make mod-clean && make build`.
