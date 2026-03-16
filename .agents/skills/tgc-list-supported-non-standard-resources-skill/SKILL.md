---
name: tgc-list-supported-non-standard-resources-skill
description: Find all supported non-standard resources in TGC. Use when you need to find all supported non-standard resources in TGC.
---

# tgc-list-supported-non-standard-resources-skill

When you need to find all supported non-standard resources in TGC, use this skill.

## When to Use This Skill

- Use this when searching for all supported non-standard resources inside TGC.
- This is helpful when you need established examples of how to implement custom conversion logic (`custom_code`) for a new resource.

---

## How to Use It

1. **Standard Resources** only have the following flags in their `mmv1/products/.../Resource.yaml` files:
   - `include_in_tgc_next: true`
   - `is_missing_in_cai: true`
   - `tgc_skip_test: true`
   - `tgc_tests: true`

2. **Non-Standard Resources** declare additional TGC-specific fields in their `mmv1/products/.../Resource.yaml` files. 
   - *Examples of non-standard fields*: 
     - `tgc_encoder`
     - `tgc_decoder`
     - `tgc_ignore_terraform_encoder`
     - `tgc_ignore_terraform_decoder`
     - `custom_tgc_expand`
     - `custom_tgc_flatten`
     - `tgc_ignore_terraform_custom_flatten`

3. To find all non-standard resources, search recursively inside the `mmv1/products/` directory for any of the non-standard fields listed in step 2.
