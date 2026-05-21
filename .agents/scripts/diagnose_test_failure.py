#!/usr/bin/env python3
import sys
import os
import re
import json
import glob
import time

def snake_to_camel(snake_str):
    components = snake_str.split('_')
    return components[0] + ''.join(x.title() for x in components[1:])

def search_key_in_dict(d, key_lower):
    """Recursively search for a key (case-insensitive) in a nested dictionary/list."""
    if isinstance(d, dict):
        for k, v in d.items():
            if k.lower() == key_lower:
                return True, k
            found, actual_k = search_key_in_dict(v, key_lower)
            if found:
                return True, actual_k
    elif isinstance(d, list):
        for item in d:
            found, actual_k = search_key_in_dict(item, key_lower)
            if found:
                return True, actual_k
    return False, None

def audit_historical_metadata(cai_type, camel_field):
    """Audit historical nightly run metadata JSON files for the missing field."""
    meta_dir = "/Users/zhenhuali/go/src/github.com/GoogleCloudPlatform/tgc-supported-resources/test_mata"
    pattern = os.path.join(meta_dir, "tests_metadata_*.json")
    files = glob.glob(pattern)
    
    checked_payloads_count = 0
    field_found_count = 0
    occurrences = []
    
    for filepath in sorted(files):
        filename = os.path.basename(filepath)
        if os.path.getsize(filepath) < 100:
            continue
            
        try:
            with open(filepath, 'r') as f:
                data = json.load(f)
        except Exception:
            continue
            
        # Iterate through the test runs in the metadata file
        for test_name, steps in data.items():
            for step_num, step_payload in steps.items():
                res_metadata = step_payload.get("resource_metadata", {})
                for addr, res_meta in res_metadata.items():
                    cai_data = res_meta.get("cai_data", {})
                    if not cai_data:
                        cai_data = res_meta.get("cai", {})
                        
                    if not isinstance(cai_data, dict):
                        continue
                        
                    for cai_addr, cai_payload in cai_data.items():
                        if not isinstance(cai_payload, dict):
                            continue
                        asset = cai_payload.get("cai_asset", {})
                        if not isinstance(asset, dict):
                            continue
                            
                        if asset.get("asset_type") == cai_type:
                            resource = asset.get("resource", {})
                            if not isinstance(resource, dict):
                                continue
                            data_block = resource.get("data", {})
                            
                            if data_block:
                                checked_payloads_count += 1
                                found, actual_k = search_key_in_dict(data_block, camel_field.lower())
                                if found:
                                    field_found_count += 1
                                    occurrences.append(f"{filename} ({test_name} step {step_num})")
                                    
    return checked_payloads_count, field_found_count, occurrences

def main():
    if len(sys.argv) < 3:
        print("Usage: diagnose_test_failure.py <log_file> <test_path>")
        sys.exit(1)
        
    log_file = sys.argv[1]
    test_path = sys.argv[2]
    
    if not os.path.exists(log_file):
        print(f"Log file not found: {log_file}")
        sys.exit(1)
        
    # 1. Parse log file for "missing fields"
    missing_fields = set()
    missing_fields_pattern = re.compile(r'missing fields:\s*\[(.*?)\]')
    
    with open(log_file, 'r') as f:
        for line in f:
            m = missing_fields_pattern.search(line)
            if m:
                fields = [f.strip() for f in m.group(1).split(',')]
                for field in fields:
                    if field:
                        missing_fields.add(field)
                        
    if not missing_fields:
        return
        
    print("\n" + "="*60)
    print("               TGC INTEGRATION TEST DIAGNOSTICS")
    print("="*60)
    print(f"Detected failed test with missing fields: {list(missing_fields)}")
    
    # 2. Find recently modified JSON files in test_path
    json_pattern = os.path.join(test_path, "*.json")
    json_files = glob.glob(json_pattern)
    
    now = time.time()
    recent_json_files = [f for f in json_files if now - os.path.getmtime(f) < 600]
    
    if not recent_json_files:
        print("Warning: No recently generated CAI JSON asset files found in test path.")
        print("Make sure WRITE_FILES=true is enabled to capture asset dumps.")
        print("="*60 + "\n")
        return
        
    # Try to extract the CAI asset type from the first valid local JSON file
    cai_type = None
    for jfile in recent_json_files:
        filename = os.path.basename(jfile)
        if "roundtrip" in filename or "reexport" in filename or "_attrs" in filename:
            continue
        try:
            with open(jfile, 'r') as f:
                content = json.load(f)
            if isinstance(content, list) and len(content) > 0:
                cai_type = content[0].get("asset_type")
            elif isinstance(content, dict):
                cai_type = content.get("asset_type")
            if cai_type:
                break
        except Exception:
            pass
            
    for field in missing_fields:
        camel_field = snake_to_camel(field)
        print(f"\nAnalyzing missing field: '{field}' (CAI name: '{camel_field}')...")
        
        missing_in_all_local_cai = True
        checked_local_files_count = 0
        
        for jfile in recent_json_files:
            filename = os.path.basename(jfile)
            if "roundtrip" in filename or "reexport" in filename or "_attrs" in filename:
                continue
                
            try:
                with open(jfile, 'r') as f:
                    content = json.load(f)
                
                checked_local_files_count += 1
                found_in_cai, actual_key = search_key_in_dict(content, camel_field.lower())
                
                if found_in_cai:
                    missing_in_all_local_cai = False
                    print(f"  [CURRENT RUN: FOUND] Field exists in CAI asset '{filename}' as '{actual_key}'")
                else:
                    print(f"  [CURRENT RUN: ABSENT] Field is missing from CAI asset '{filename}'")
                    
            except Exception as e:
                print(f"  Error reading {filename}: {e}")
                
        if checked_local_files_count == 0:
            continue
            
        print("\nChecking historical nightly metadata runs...")
        if cai_type:
            print(f"  Target CAI Asset Type: '{cai_type}'")
            hist_checked, hist_found, occurrences = audit_historical_metadata(cai_type, camel_field)
            
            if hist_checked > 0:
                print(f"  Checked {hist_checked} historical asset payloads across past nightly runs.")
                if hist_found == 0:
                    print(f"  [HISTORICAL RUNS: ABSENT] Field was NEVER found in any historical runs.")
                    is_systematic_omission = True
                else:
                    print(f"  [HISTORICAL RUNS: FOUND] Field was found in {hist_found} runs:")
                    for occ in occurrences[:5]:
                        print(f"    - {occ}")
                    if len(occurrences) > 5:
                        print(f"    - and {len(occurrences) - 5} other occurrences.")
                    is_systematic_omission = False
            else:
                print("  No historical data found matching this CAI Asset Type.")
                is_systematic_omission = missing_in_all_local_cai
        else:
            print("  Could not extract CAI Asset Type to perform historical search.")
            is_systematic_omission = missing_in_all_local_cai
            
        print("-"*60)
        if is_systematic_omission:
            print(f"DIAGNOSTIC RESULT for '{field}':")
            print(f"  --> The field is systematically ABSENT from BOTH current and historical CAI asset payloads.")
            print(f"  --> This proves a permanent omission on the CAIS API service side.")
            print("\nRECOMMENDED ACTION:")
            print(f"  >> Set 'is_missing_in_cai: true' on '{camel_field}' in the Magic Modules YAML configuration.")
        else:
            if missing_in_all_local_cai:
                print(f"DIAGNOSTIC RESULT for '{field}':")
                print(f"  --> The field is absent in the current run, but WAS present in historical runs.")
                print(f"  --> This indicates a transient API deployment issue, environment drift, or conditional logic.")
                print("\nRECOMMENDED ACTION:")
                print(f"  >> Verify if the target Google Cloud project has the feature flag/preview API enabled.")
            else:
                print(f"DIAGNOSTIC RESULT for '{field}':")
                print(f"  --> The field exists in the current CAI asset, but was lost during HCL conversion.")
                print("\nRECOMMENDED ACTION:")
                print(f"  >> Inspect 'cai2hcl' flatteners, encoders, or decoders in mmv1/templates/tgc_next/.")
                print(f"  >> A custom flattener/decoder may be required to map the field to HCL.")
            
    print("="*60 + "\n")

if __name__ == "__main__":
    main()
