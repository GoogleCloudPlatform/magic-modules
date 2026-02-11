import glob
import os

print("=== FINAL INTEGRATION TEST REPORT ===")

logs = glob.glob(os.path.expanduser("~/go/src/github.com/GoogleCloudPlatform/terraform-google-conversion/*.log"))
pass_count = 0
fail_count = 0
skipped_count = 0
unknown = 0
total_tests = 0

passed_tests = []
failed_tests = []
skipped_tests = []

for log in logs:
    total_tests += 1
    test_name = os.path.basename(log).replace('.log', '')
    with open(log, 'r') as f:
        content = f.read()
    
    if f"--- PASS: {test_name}" in content:
        if f"--- PASS: {test_name}/" in content:
            passed_tests.append(test_name)
            pass_count += 1
        else:
            skipped_tests.append(test_name)
            skipped_count += 1
    elif "--- FAIL:" in content:
        failed_tests.append(test_name)
        fail_count += 1
    elif "No changes. Infrastructure is up-to-date" in content or "Apply complete!" in content:
        passed_tests.append(test_name)
        pass_count += 1
    elif "no tests to run" in content:
        skipped_tests.append(test_name)
        skipped_count += 1
    else:
        unknown += 1

print("\n=== SUMMARY ===")
print(f"Total Log Files Processed So Far: {total_tests}")
print(f"Passed: {pass_count}")
print(f"Failed: {fail_count}")

print("\n--- PASSED TESTS ---")
for p in passed_tests:
    print(f"✅ {p}")

print("\n--- FAILED TESTS ---")
for f in failed_tests:
    print(f"❌ {f}")
    
print("\n--- SKIPPED TESTS ---")
for s in skipped_tests:
    print(f"⏭️  {s}")
