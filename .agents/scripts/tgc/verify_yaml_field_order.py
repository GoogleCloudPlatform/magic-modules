#!/usr/bin/env python3
import sys
import os
import re

def parse_struct_tags(go_file_path, target_struct_name):
    """Parses a specific struct's tags from a Go file, resolving inlined structs recursively."""
    tags = []
    struct_pattern = re.compile(r'type\s+' + re.escape(target_struct_name) + r'\s+struct\s*\{')
    
    with open(go_file_path, "r") as f:
        lines = f.readlines()

    start_idx = -1
    for idx, line in enumerate(lines):
        if struct_pattern.search(line):
            start_idx = idx
            break

    if start_idx == -1:
        return []

    brace_count = 1
    inline_pattern = re.compile(r'([a-zA-Z0-9_]+)\s+`yaml:",inline"`')
    tag_pattern = re.compile(r'`yaml:"([^,"]+)')

    for line in lines[start_idx + 1:]:
        line = line.strip()
        if not line:
            continue
        
        brace_count += line.count("{")
        brace_count -= line.count("}")
        if brace_count <= 0:
            break # End of struct

        # Check for inline structs
        inline_match = inline_pattern.search(line)
        if inline_match:
            embedded_struct = inline_match.group(1)
            tags.extend(parse_struct_tags(go_file_path, embedded_struct))
            continue

        # Check for standard tag
        tag_match = tag_pattern.search(line)
        if tag_match:
            tags.append(tag_match.group(1))

    return tags

def extract_yaml_keys(yaml_file_path):
    """Extracts top-level keys in order from a YAML file."""
    keys = []
    # Match top-level keys (no leading whitespace, word character followed by colon)
    key_pattern = re.compile(r'^([a-zA-Z0-9_-]+)\s*:')

    with open(yaml_file_path, "r") as f:
        for line in f:
            match = key_pattern.match(line)
            if match:
                keys.append(match.group(1))
    return keys

def get_git_modified_keys(yaml_file_path):
    """Uses git diff to find top-level keys that were added or modified in the current changes."""
    import subprocess
    modified_keys = set()
    try:
        # Run git diff -U0 to get only the changed lines without context
        res = subprocess.run(
            ["git", "diff", "-U0", yaml_file_path],
            capture_output=True,
            text=True,
            check=True
        )
        # Match added/modified top-level keys (lines starting with '+' but not '+++' and containing key definition)
        added_line_pattern = re.compile(r'^\+\s*([a-zA-Z0-9_-]+)\s*:')
        for line in res.stdout.splitlines():
            if line.startswith("+++"):
                continue
            match = added_line_pattern.match(line)
            if match:
                modified_keys.add(match.group(1))
    except Exception as e:
        # Fallback if git is not initialized or failed
        pass
    return modified_keys

def main():
    if len(sys.argv) < 3:
        print("Usage: verify_yaml_field_order.py <go-struct-file> <yaml-file> [struct-name]", file=sys.stderr)
        sys.exit(1)

    go_file = sys.argv[1]
    yaml_file = sys.argv[2]
    struct_name = sys.argv[3] if len(sys.argv) > 3 else "Resource"

    if not os.path.exists(go_file):
        print(f"Error: Go file not found at {go_file}", file=sys.stderr)
        sys.exit(1)

    if not os.path.exists(yaml_file):
        print(f"Error: YAML file not found at {yaml_file}", file=sys.stderr)
        sys.exit(1)

    canonical_tags = parse_struct_tags(go_file, struct_name)
    yaml_keys = extract_yaml_keys(yaml_file)

    # Filter lists to only elements present in both for direct relative comparison
    yaml_keys_set = set(yaml_keys)
    canonical_tags_set = set(canonical_tags)
    
    filtered_canonical = [tag for tag in canonical_tags if tag in yaml_keys_set]
    filtered_yaml = [key for key in yaml_keys if key in canonical_tags_set]

    # Determine git-modified keys to restrict check to newly added/modified fields
    modified_keys = get_git_modified_keys(yaml_file)
    
    # If no keys are modified (e.g. during a dry run or clean environment), fallback to checking everything
    if not modified_keys:
        modified_keys = yaml_keys_set

    # Check relative ordering for newly modified keys
    has_violations = False
    violations = []

    canonical_index = {tag: i for i, tag in enumerate(filtered_canonical)}
    yaml_index = {key: i for i, key in enumerate(filtered_yaml)}

    for key in filtered_yaml:
        if key in modified_keys:
            # Compare relative position of this modified key against all other keys
            for other_key in filtered_yaml:
                if other_key == key:
                    continue
                
                expected_before = canonical_index[key] < canonical_index[other_key]
                actual_before = yaml_index[key] < yaml_index[other_key]
                
                if expected_before != actual_before:
                    has_violations = True
                    rel = "before" if expected_before else "after"
                    violations.append(f"⚠️  New/modified field '{key}' should be positioned {rel} '{other_key}'")

    if has_violations:
        print(f"❌ FIELD ORDERING VERIFICATION FAILED for {yaml_file}:", file=sys.stderr)
        print("\nIncorrect relative order detected for new or modified TGC fields:\n", file=sys.stderr)
        for violation in sorted(list(set(violations))):
            print(f"  {violation}", file=sys.stderr)
        
        print("\n[Expected relative order of all keys from Go Struct]:", file=sys.stderr)
        for i, tag in enumerate(filtered_canonical):
            print(f"  {i+1:2d}. {tag}", file=sys.stderr)
        sys.exit(1)

    print(f"✅ Field ordering verification passed for new/modified fields in {yaml_file}.")
    sys.exit(0)

if __name__ == "__main__":
    main()
