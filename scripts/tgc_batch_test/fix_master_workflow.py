import json
import os

with open('matched_by_product.json') as f:
    mappings = json.load(f)

# The auto-excluder list from the crashed Python run
excluded = [
    "google_bigtable_app_profile",
    "google_bigquery_data_transfer_config",
    "google_network_management_connectivity_test",
    "google_network_services_gateway",
    "google_eventarc_channel",
    "google_eventarc_trigger",
    "google_compute_network_firewall_policy",
    "google_compute_region_security_policy",
    "google_compute_ssl_certificate",
    "google_netapp_volume",
    "google_netapp_volume_replication",
    "google_bigquery_dataset_access",
    "google_bigquery_routine",
    "google_bigquery_table",
    "google_compute_instance",
    "google_compute_instance_template"
]

for product in mappings:
    filtered = []
    for res in mappings[product]:
        if res['target'] not in excluded:
            filtered.append(res)
    mappings[product] = filtered

products = [p for p in sorted(mappings.keys()) if mappings[p]]

script_lines = []

script_lines.append("#!/bin/bash")
script_lines.append("export WRITE_FILES=true")
script_lines.append("export GOPATH=$HOME/go")
script_lines.append("MM_DIR=\"/Users/zhenhuali/Documents/workspace/magic-modules\"")
script_lines.append("TGC_DIR=\"$GOPATH/src/github.com/GoogleCloudPlatform/terraform-google-conversion\"")
script_lines.append("")
script_lines.append("rm -f $MM_DIR/tgc_report.txt")
script_lines.append("echo \"=== Starting TGC Master Workflow ===\" > $MM_DIR/tgc_report.txt")
script_lines.append("")

script_lines.append("echo \"=== Clean up leftover modified files ===\"")
script_lines.append("cd \"$MM_DIR\" || exit 1")
script_lines.append("git checkout -- mmv1/products/")
script_lines.append("")
script_lines.append("echo \"=== Prep: Disabling parallel testing in Makefile ===\"")
script_lines.append("sed -i '' 's/-p 8 -parallel 8 //g' \"$MM_DIR/mmv1/third_party/tgc_next/Makefile\"")
script_lines.append("")

script_lines.append("echo \"=== Enabling All Valid Resources ===\"")
script_lines.append("cd \"$MM_DIR\" || exit 1")
    
for product in products:
    for res in mappings[product]:
        filepath = f"mmv1/products/{product}/{res['file']}.yaml"
        target_res = res['target']
        
        # Check and inject include_in_tgc_next: true
        script_lines.append(f"if [ -f \"{filepath}\" ]; then")
        script_lines.append(f"    if ! grep -q \"include_in_tgc_next: true\" \"{filepath}\"; then")
        script_lines.append(f"        sed -i '' '/^name:/a\\")
        script_lines.append(f"include_in_tgc_next: true\\")
        script_lines.append(f"' \"{filepath}\"")
        script_lines.append(f"    fi")
        script_lines.append("else")
        script_lines.append(f"    echo \"WARNING: File not found {filepath}\"")
        script_lines.append("fi")
script_lines.append("")

script_lines.append("echo \"=== Generating and Compiling All TGC Code ===\"")
script_lines.append("cd \"$MM_DIR\" || exit 1")
script_lines.append("make clean-tgc OUTPUT_PATH=\"$TGC_DIR\"")
script_lines.append("make tgc OUTPUT_PATH=\"$TGC_DIR\"")
script_lines.append("if [ $? -ne 0 ]; then")
script_lines.append("    echo \"FATAL: Global build failed! Check for remaining broken resources.\"")
script_lines.append("    exit 1")
script_lines.append("fi")
script_lines.append("")
script_lines.append("cd \"$TGC_DIR\" || exit 1")
script_lines.append("make mod-clean")
script_lines.append("make build")
script_lines.append("if [ $? -ne 0 ]; then")
script_lines.append("    echo \"FATAL: Global build failed! Check for remaining broken resources.\"")
script_lines.append("    exit 1")
script_lines.append("fi")
script_lines.append("")

for product in products:
    script_lines.append(f"echo \"=======================================\"")
    script_lines.append(f"echo \"Testing Product: {product}\"")
    script_lines.append(f"echo \"=======================================\"")
    
    script_lines.append(f"echo \"---> Running Integration Tests for {product}\"")
    script_lines.append("cd \"$TGC_DIR\" || exit 1")
    
    script_lines.append(f"echo \"=== Summary for {product.capitalize()} Resources ===\" >> $MM_DIR/tgc_report.txt")
    script_lines.append(f"echo \"Passed/Executed Resources:\" >> $MM_DIR/tgc_report.txt")
    
    for res in mappings[product]:
        target_res = res['target']
        # Convert snake_case to CamelCase (google_compute_forwarding_rule -> ComputeForwardingRule)
        parts = target_res.split('_')[1:] # Skip 'google'
        camel_name = ''.join(p.capitalize() for p in parts)
        test_func = f"TestAcc{camel_name}"
        
        log_file = f"{test_func}.log"
        script_lines.append(f"echo \"Running test: {test_func}...\"")
        
        # We only pass if ANY subtests actually passed (ignoring skips)
        script_lines.append(f"make test-integration-local TESTPATH=\"./test/services/{product}\" TESTARGS=\"-run=^{test_func}\" > \"{log_file}\" 2>&1")
        script_lines.append(f"if grep -q \"--- PASS: {test_func}\" \"{log_file}\"; then")
        # Check if subtests actually ran vs all skipped
        script_lines.append(f"    if grep -q \"--- PASS: {test_func}/\" \"{log_file}\"; then")
        script_lines.append(f"        PASS_COUNT=$(grep -c \"--- PASS: {test_func}/\" \"{log_file}\")")
        script_lines.append(f"        echo \"✅ {test_func} ($PASS_COUNT tests passed)\" >> $MM_DIR/tgc_report.txt")
        script_lines.append("    else")
        script_lines.append(f"        echo \"🟡 {test_func}: ALL EXAMPLES SKIPPED (Top level passed but no subtests ran)\" >> $MM_DIR/tgc_report.txt_failures")
        script_lines.append("    fi")
        script_lines.append("else")
        script_lines.append(f"    FAIL_COUNT=$(grep -c \"--- FAIL:\" \"{log_file}\")")
        script_lines.append(f"    if [ \"$FAIL_COUNT\" -gt 0 ]; then")
        script_lines.append(f"        echo \"❌ {test_func} FAILED ($FAIL_COUNT failures)\" >> $MM_DIR/tgc_report.txt_failures")
        script_lines.append("    else")
        script_lines.append(f"        echo \"⏳ {test_func} TIMED OUT OR UNKNOWN ERROR\" >> $MM_DIR/tgc_report.txt_failures")
        script_lines.append("    fi")
        script_lines.append("fi")
    script_lines.append(f"echo \"\" >> $MM_DIR/tgc_report.txt")
    
    script_lines.append(f"echo \"Skipped/Failed Details:\" >> $MM_DIR/tgc_report.txt")
    script_lines.append(f"if [ -f $MM_DIR/tgc_report.txt_failures ]; then")
    script_lines.append(f"    cat $MM_DIR/tgc_report.txt_failures >> $MM_DIR/tgc_report.txt")
    script_lines.append(f"    rm $MM_DIR/tgc_report.txt_failures")
    script_lines.append("else")
    script_lines.append(f"    echo \"No tests totally skipped or failed.\" >> $MM_DIR/tgc_report.txt")
    script_lines.append("fi")
    
script_lines.append("echo \"=======================================\"")
script_lines.append("echo \"=== Cleanup: Restoring Workspace ===\"")
script_lines.append("cd \"$MM_DIR\" || exit 1")
script_lines.append("git checkout -- mmv1/products/")
script_lines.append("echo \"Master workflow finished! Results in tgc_report.txt\"")

script_lines.append("cat $MM_DIR/tgc_report.txt")

script_lines.append("")

# Write to next version script
version = 1
while os.path.exists(f'tgc_master_workflow{version}.sh'):
    version += 1

out_file = f'tgc_master_workflow{version}.sh'
with open(out_file, 'w') as f:
    f.write('\n'.join(script_lines) + '\n')

os.chmod(out_file, 0o755)
print(f"Generated {out_file} configured for {sum(len(v) for v in mappings.values())} mapped resources.")
