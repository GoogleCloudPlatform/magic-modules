import sys
import os

def add_missing_tests(yaml_path, missing_tests_path):
    if not os.path.exists(missing_tests_path):
        print(f"File not found: {missing_tests_path}")
        return

    with open(missing_tests_path, 'r') as f:
        missing_tests = [line.strip() for line in f if line.strip()]

    if not missing_tests:
        print("No missing tests to add.")
        return

    if not os.path.exists(yaml_path):
        print(f"YAML file not found: {yaml_path}")
        return

    with open(yaml_path, 'r') as f:
        lines = f.readlines()

    # Find if tgc_tests already exists
    tgc_tests_index = -1
    for i, line in enumerate(lines):
        if line.strip().startswith('tgc_tests:'):
            tgc_tests_index = i
            break

    if tgc_tests_index != -1:
        # tgc_tests exists. Let's find existing tests to avoid duplication
        existing_tests = set()
        j = tgc_tests_index + 1
        while j < len(lines) and (lines[j].startswith(' ') or lines[j].startswith('\t') or lines[j].strip() == ''):
            if "name:" in lines[j]:
                test_name = lines[j].split("name:")[1].strip().strip("'").strip('"')
                existing_tests.add(test_name)
            j += 1
        
        # Insert missing tests that are not already there
        insert_index = tgc_tests_index + 1
        for test in missing_tests:
            if test not in existing_tests:
                lines.insert(insert_index, f"  - name: '{test}'\n")
                insert_index += 1
    else:
        # Append tgc_tests: to the end of the file
        lines.append('\ntgc_tests:\n')
        for test in missing_tests:
            lines.append(f"  - name: '{test}'\n")

    with open(yaml_path, 'w') as f:
        f.writelines(lines)
    print(f"Added {len(missing_tests)} tests to {yaml_path}")

if __name__ == '__main__':
    if len(sys.argv) < 3:
        print("Usage: python add_missing_tests.py <YAMLFilePath> <MissingTestsFilePath>")
        sys.exit(1)
    yaml_path = sys.argv[1]
    missing_tests_path = sys.argv[2]
    add_missing_tests(yaml_path, missing_tests_path)
