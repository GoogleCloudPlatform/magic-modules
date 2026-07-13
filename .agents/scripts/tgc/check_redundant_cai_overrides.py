#!/usr/bin/env python3
import sys
import os
import re

def extract_field_value(yaml_file_path, field_name):
    """Extracts the string value of a top-level field in a YAML file using regex."""
    # Matches field_name followed by optional quotes, capturing value
    pattern = re.compile(r'^' + re.escape(field_name) + r'\s*:\s*[\'"]?([^\'"]+)[\'"]?')
    with open(yaml_file_path, "r") as f:
        for line in f:
            match = pattern.match(line.strip())
            if match:
                return match.group(1).strip()
    return None

def main():
    if len(sys.argv) < 2:
        print("Usage: check_redundant_cai_overrides.py <resource-yaml-file>", file=sys.stderr)
        sys.exit(1)

    resource_yaml = sys.argv[1]
    if not os.path.exists(resource_yaml):
        print(f"Error: Resource YAML file not found at {resource_yaml}", file=sys.stderr)
        sys.exit(1)

    cai_override = extract_field_value(resource_yaml, "cai_asset_name_format")
    if not cai_override:
        # No override specified, so it's perfectly clean!
        print(f"✅ No cai_asset_name_format override specified in {os.path.basename(resource_yaml)} (Clean).")
        sys.exit(0)

    id_format = extract_field_value(resource_yaml, "id_format")
    if not id_format:
        print(f"Warning: id_format not found in {resource_yaml}. Skipping redundancy check.", file=sys.stderr)
        sys.exit(0)

    # Standardize by stripping any leading/trailing slashes and standard prefixes
    def normalize(val):
        val = val.strip().strip("/")
        # Remove potential leading domain/service name if full URL was passed
        val = re.sub(r'^(https?://)?([a-zA-Z0-9.-]+\.googleapis\.com/)?', '', val)
        return val

    norm_cai_override = normalize(cai_override)
    norm_id_format = normalize(id_format)

    if norm_cai_override == norm_id_format:
        print(f"❌ REDUNDANT CAI OVERRIDE DETECTED in {resource_yaml}:", file=sys.stderr)
        print(f"  cai_asset_name_format: '{cai_override}'", file=sys.stderr)
        print(f"  id_format:             '{id_format}'", file=sys.stderr)
        print("\nThis override is identical to the default computed format and can be safely removed.", file=sys.stderr)
        sys.exit(1)

    print(f"✅ Verification passed: cai_asset_name_format override is custom and non-redundant.")
    sys.exit(0)

if __name__ == "__main__":
    main()
