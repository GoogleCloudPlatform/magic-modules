#!/usr/bin/env python3
import sys
import os
import re

def find_test_file_and_body(service_dir, test_name):
    """
    Searches the service directory for Go test files containing the test_name function.
    Returns (file_path, function_body, receiver_name) or (None, None, None).
    """
    if not os.path.exists(service_dir):
        return None, None, None
        
    # Pattern to match the test function definition and capture the receiver name (usually 't')
    func_pattern = re.compile(
        r'func\s+' + re.escape(test_name) + r'\s*\(\s*([a-zA-Z0-9_]+)\s+\*testing\.T\s*\)'
    )
    
    for filename in os.listdir(service_dir):
        if not (filename.endswith("_test.go") or filename.endswith("_test.go.tmpl")):
            continue
        file_path = os.path.join(service_dir, filename)
        with open(file_path, "r", encoding="utf-8") as f:
            lines = f.readlines()
            
        for idx, line in enumerate(lines):
            match = func_pattern.search(line)
            if match:
                receiver_name = match.group(1)
                # Extract function body
                body = extract_go_function_body(lines, idx)
                return file_path, body, receiver_name
                
    return None, None, None

def extract_go_function_body(lines, start_idx):
    brace_count = 0
    body_lines = []
    found_start = False
    for line in lines[start_idx:]:
        if not found_start:
            if '{' in line:
                brace_count += line.count('{')
                brace_count -= line.count('}')
                found_start = True
                body_lines.append(line)
                if brace_count == 0:
                    break
            continue
        
        brace_count += line.count('{')
        brace_count -= line.count('}')
        body_lines.append(line)
        if brace_count <= 0:
            break
            
    return "".join(body_lines)

def check_yaml_tests(yaml_file_path):
    """
    Parses the YAML file, checks test configurations, and validates them against handwritten test bodies.
    """
    # 1. Read file content
    with open(yaml_file_path, "r", encoding="utf-8") as f:
        content = f.read()
        
    # Check if resource is included in TGC Next
    if not re.search(r'include_in_tgc_next\s*:\s*true', content):
        return [] # Skipped
        
    # Extract service name (product name) from file path
    # Path is like mmv1/products/<product>/<Resource>.yaml
    path_parts = yaml_file_path.split(os.sep)
    if len(path_parts) < 3:
        return []
    product_name = path_parts[-2]
    service_dir = os.path.join("mmv1", "third_party", "terraform", "services", product_name)
    
    # 2. Extract tgc_tests block
    # Simple YAML block extractor for tgc_tests
    # Look for tgc_tests: followed by indented lines starting with - name:
    tgc_tests_match = re.search(r'tgc_tests\s*:\s*\n((?:\s+-\s*name\s*:\s*[^\n]+\n?)+)', content)
    if not tgc_tests_match:
        return [] # No tgc_tests configured explicitly, or empty
        
    test_entries = tgc_tests_match.group(1)
    test_names = re.findall(r'-\s*name\s*:\s*[\'"]?([^\'"\n]+)[\'"]?', test_entries)
    
    errors = []
    for full_test_name in test_names:
        # Split test name into parent and subtest
        parts = full_test_name.split("/")
        parent_test_name = parts[0]
        has_subtest_in_name = len(parts) > 1
        
        # Locate handwritten test file and body
        file_path, body, receiver_name = find_test_file_and_body(service_dir, parent_test_name)
        if not body:
            # Test function not found in service directory.
            # (It could be defined in a common/shared package, but typically they are in the service directory).
            continue
            
        # Check if the function body uses receiver.Run("
        # e.g., t.Run(
        subtest_pattern = re.compile(re.escape(receiver_name) + r'\.Run\s*\(')
        has_subtests_in_body = bool(subtest_pattern.search(body))
        
        if has_subtests_in_body and not has_subtest_in_name:
            errors.append(
                f"❌ Test verification failed for {yaml_file_path}:\n"
                f"  Test name '{full_test_name}' is configured as a top-level test, but the handwritten\n"
                f"  test function '{parent_test_name}' in {file_path} contains subtests (uses {receiver_name}.Run).\n"
                f"  -> You MUST configure it with a specific subtest name, e.g. '{parent_test_name}/<subtest>'"
            )
        elif not has_subtests_in_body and has_subtest_in_name:
            errors.append(
                f"❌ Test verification failed for {yaml_file_path}:\n"
                f"  Test name '{full_test_name}' contains a subtest suffix, but the handwritten\n"
                f"  test function '{parent_test_name}' in {file_path} does NOT contain any subtests (no {receiver_name}.Run found).\n"
                f"  -> Please configure it without a subtest name suffix, e.g. '{parent_test_name}'"
            )
            
    return errors

def main():
    if len(sys.argv) < 2:
        print("Usage: verify_test_configs.py <yaml-file>", file=sys.stderr)
        sys.exit(1)
        
    yaml_file = sys.argv[1]
    if not os.path.exists(yaml_file):
        print(f"Error: YAML file not found at {yaml_file}", file=sys.stderr)
        sys.exit(1)
        
    errors = check_yaml_tests(yaml_file)
    if errors:
        for err in errors:
            print(err, file=sys.stderr)
        sys.exit(1)
        
    print(f"✅ Test config verification passed for {yaml_file}.")
    sys.exit(0)

if __name__ == "__main__":
    main()
