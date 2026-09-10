#!/usr/bin/env python3
import os
import sys
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

def parse_tf_file(filepath):
    if not os.path.exists(filepath):
        return {}
    with open(filepath, 'r') as f:
        content = f.read()
    
    attributes = {}
    scope_stack = []
    
    lines = []
    for line in content.splitlines():
        line = line.strip()
        if not line or line.startswith('#') or line.startswith('//'):
            continue
        lines.append(line)
        
    for line in lines:
        if line.endswith('{'):
            block_decl = line[:-1].strip()
            parts = block_decl.split()
            if len(parts) == 1:
                block_name = parts[0]
            elif len(parts) > 2 and parts[0] in ('resource', 'data'):
                block_name = "" # ignore resource wrapper type
            else:
                block_name = parts[-1].replace('"', '')
            
            if block_name:
                scope_stack.append(block_name)
            else:
                scope_stack.append("")
        elif line == '}':
            if scope_stack:
                scope_stack.pop()
        elif '=' in line:
            parts = line.split('=', 1)
            key = parts[0].strip().replace('"', '')
            val = parts[1].strip().strip('"').strip("'")
            
            active_scopes = [s for s in scope_stack if s]
            if active_scopes:
                full_key = ".".join(active_scopes) + "." + key
            else:
                full_key = key
                
            attributes[full_key] = val
            
    return attributes

def flatten_json(obj, prefix=""):
    flat = {}
    if isinstance(obj, dict):
        for k, v in obj.items():
            new_prefix = f"{prefix}.{k}" if prefix else k
            flat.update(flatten_json(v, new_prefix))
    elif isinstance(obj, list):
        for i, item in enumerate(obj):
            new_prefix = f"{prefix}.{i}" if prefix else str(i)
            flat.update(flatten_json(item, new_prefix))
    else:
        flat[prefix] = str(obj)
    return flat

def fuzzy_match_keys(hcl_key, cai_keys):
    norm_hcl = hcl_key.lower().replace("_", "").replace(".", "")
    norm_hcl = norm_hcl.replace("nodeconfig", "config")
    
    for cai_key in cai_keys:
        norm_cai = cai_key.lower().replace("_", "").replace(".", "")
        if norm_hcl in norm_cai or norm_cai in norm_hcl:
            return cai_key
    return None

def audit_historical_metadata(cai_type, camel_field):
    """Audit historical nightly run metadata JSON files for the missing field."""
    tgc_dir = os.environ.get("TGC_DIR")
    if not tgc_dir and len(sys.argv) >= 3:
        test_path = sys.argv[2]
        parts = os.path.abspath(test_path).split(os.sep)
        if "test" in parts and "services" in parts:
            idx = parts.index("test")
            tgc_dir = os.sep.join(parts[:idx])
    if not tgc_dir:
        tgc_dir = "."
    meta_dir = os.path.join(tgc_dir, "test_meta")
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

def run_deep_comparative_diagnostics(test_name, test_path):
    input_cai_path = os.path.join(test_path, f"{test_name}.json")
    export_tf_path = os.path.join(test_path, f"{test_name}_export.tf")
    roundtrip_cai_path = os.path.join(test_path, f"{test_name}_roundtrip.json")
    roundtrip_tf_path = os.path.join(test_path, f"{test_name}_roundtrip.tf")
    
    if not os.path.exists(input_cai_path) or not os.path.exists(export_tf_path):
        return
        
    print("\n" + "="*80)
    print(f"    DETAILED PIPELINE DATA LOSS REPORT: {test_name}")
    print("="*80)
    
    has_roundtrip = os.path.exists(roundtrip_cai_path) and os.path.exists(roundtrip_tf_path)
    export_attrs = parse_tf_file(export_tf_path)
    
    with open(input_cai_path, 'r') as f:
        input_cai = json.load(f)
    input_cai_data = input_cai[0]["resource"]["data"] if len(input_cai) > 0 else {}
    input_cai_attrs = flatten_json(input_cai_data)
    
    failures_found = 0
    
    if not has_roundtrip:
        print("  --> TEST FAILED EARLY IN STEP 1 (cai2hcl / HCL Export phase).")
        print("  --> Analyzing why CAI fields were missing from the HCL export configuration...")
        
        for cai_key, cai_val in input_cai_attrs.items():
            if any(x in cai_key.lower() for x in ("project", "location", "zone", "name", "id")):
                continue
                
            matched_hcl_key = fuzzy_match_keys(cai_key, export_attrs.keys())
            
            if not matched_hcl_key:
                failures_found += 1
                print(f"\n[DATA LOSS DETECTED] Field: {cai_key}")
                print(f"  - Original CAI Value: '{cai_val}'")
                print(f"  - Exported HCL Value: (NOT FOUND)")
                print("  - Diagnosis:")
                print("    --> [cai2hcl FLATTENER BUG]: The field exists in the original CAI asset,")
                print("        but the flattener failed to map it into the exported HCL configuration.")
                if "." in cai_key:
                    print(f"    --> [STRUCTURAL MISMATCH WARNING]: The field '{cai_key}' is nested in the CAI asset,")
                    print("        but it may be defined as a top-level property in the Terraform/Magic Modules schema.")
                    print("        Check if a custom flattener or decoder is needed to extract it from the nested map.")
                print("    --> Action: Check/implement the flattener (e.g. flatten<FieldName>) in the cai2hcl.go or templates.")
    else:
        print("  --> TEST FAILED IN STEP 2 (tfplan2cai / Round-trip validation phase).")
        print("  --> Analyzing data loss across HCL -> CAI -> HCL round-trip...")
        
        roundtrip_attrs = parse_tf_file(roundtrip_tf_path)
        with open(roundtrip_cai_path, 'r') as f:
            roundtrip_cai = json.load(f)
        roundtrip_cai_data = roundtrip_cai[0]["resource"]["data"] if len(roundtrip_cai) > 0 else {}
        roundtrip_cai_attrs = flatten_json(roundtrip_cai_data)
        
        for key, val in export_attrs.items():
            if key in ("project", "location", "zone", "id"):
                continue
                
            rt_val = roundtrip_attrs.get(key)
            
            if not rt_val or rt_val == "":
                failures_found += 1
                print(f"\n[DATA LOSS DETECTED] Field: {key}")
                print(f"  - Export Value (cai2hcl output): '{val}'")
                print(f"  - Roundtrip Value (cai2hcl input): '{rt_val or '(empty)'}'")
                
                matched_cai_key = fuzzy_match_keys(key, roundtrip_cai_attrs.keys())
                
                if matched_cai_key:
                    cai_val = roundtrip_cai_attrs[matched_cai_key]
                    print(f"  - Roundtrip CAI Value: '{cai_val}' (Path: {matched_cai_key})")
                    print("  - Diagnosis:")
                    print("    --> [cai2hcl FLATTENER BUG]: The field successfully exists in the round-trip CAI JSON,")
                    print("        but the flattener failed to map it back into the Terraform HCL config.")
                    if "." in matched_cai_key:
                        print(f"    --> [STRUCTURAL MISMATCH WARNING]: The field is nested inside the CAI asset as '{matched_cai_key}',")
                        print("        but the flattener is likely attempting to extract it from the top-level structure.")
                        print("        Check if a custom flattener/decoder is needed to resolve the hierarchy.")
                    print("    --> Action: Update the flattener (e.g. flatten<Name>) in the corresponding cai2hcl.go or templates.")
                else:
                    print("  - Roundtrip CAI Value: (NOT FOUND)")
                    print("  - Diagnosis:")
                    print("    --> [tfplan2cai EXPANDER BUG]: The field is missing completely from the round-trip CAI JSON.")
                    print("        The expander failed to map the Terraform HCL config property into the CAI asset.")
                    print("    --> Action: Update the expander (e.g. expand<Name>) in the corresponding tfplan2cai.go or templates.")
                    
                    input_matched = fuzzy_match_keys(key, input_cai_attrs.keys())
                    if input_matched:
                        print(f"    --> Note: Field is present in original CAI asset: '{input_cai_attrs[input_matched]}' (Path: {input_matched})")
                        print("        Using 'is_missing_in_cai: true' is FORBIDDEN here because the CAI asset natively supports the field.")
                    else:
                        print("    --> Note: Field is also missing from the original CAI asset.")
                        print("        This is a valid candidate for 'is_missing_in_cai: true' mapping if the CAI asset never has it.")
                        
    if failures_found == 0:
        print("\nAll fields round-trip successfully! Parity check PASSED.")
    else:
        print(f"\nDeep comparative diagnostics complete. Found {failures_found} data loss issues.")
    print("="*80 + "\n")

def main():
    if len(sys.argv) < 3:
        print("Usage: diagnose_test_failure.py <log_file> <test_path>")
        sys.exit(1)
        
    log_file = sys.argv[1]
    test_path = sys.argv[2]
    
    if not os.path.exists(log_file):
        print(f"Log file not found: {log_file}")
        sys.exit(1)
        
    # 1. Parse log file for "missing fields" and failed subtest names
    missing_fields = set()
    missing_fields_pattern = re.compile(r'missing fields:\s*\[(.*?)\]')
    
    failed_subtests = []
    fail_regex = re.compile(r'--- FAIL:\s+(\S+)')
    
    with open(log_file, 'r') as f:
        for line in f:
            m = missing_fields_pattern.search(line)
            if m:
                fields = [f.strip() for f in m.group(1).split(',')]
                for field in fields:
                    if field:
                        missing_fields.add(field)
                        
            m2 = fail_regex.search(line)
            if m2:
                subtest = m2.group(1)
                if '/' in subtest and 'step' in subtest:
                    parts = subtest.split('/')
                    if len(parts) >= 3:
                        leaf_name = f"{parts[1]}_{parts[2]}"
                        if leaf_name not in failed_subtests:
                            failed_subtests.append(leaf_name)
                            
    # 2. Run deep comparative diagnostics if intermediate files exist
    for subtest_name in failed_subtests:
        run_deep_comparative_diagnostics(subtest_name, test_path)
                            
    if not missing_fields:
        return
        
    print("\n" + "="*60)
    print("               TGC INTEGRATION TEST DIAGNOSTICS")
    print("="*60)
    print(f"Detected failed test with missing fields: {list(missing_fields)}")
    
    # 3. Find recently modified JSON files in test_path
    json_pattern = os.path.join(test_path, "*.json")
    json_files = glob.glob(json_pattern)
    
    now = time.time()
    recent_json_files = [f for f in json_files if now - os.path.getmtime(f) < 600]
    
    if not recent_json_files:
        print("Warning: No recently generated CAI JSON asset files found in test path.")
        print("Make sure WRITE_FILES=true is enabled to capture asset dumps.")
        print("="*60 + "\n")
        return
        
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
        
        # Target exactly the specific failed step's raw input CAI asset file first
        target_json_filename = f"{failed_subtests[0]}.json" if failed_subtests else None
        target_jfile = os.path.join(test_path, target_json_filename) if target_json_filename else None
        
        if target_jfile and os.path.exists(target_jfile):
            try:
                with open(target_jfile, 'r') as f:
                    content = json.load(f)
                checked_local_files_count += 1
                found_in_cai, actual_key = search_key_in_dict(content, camel_field.lower())
                if found_in_cai:
                    missing_in_all_local_cai = False
                    print(f"  [CURRENT RUN: FOUND] Field exists in the targeted raw CAI asset '{target_json_filename}' as '{actual_key}'")
                else:
                    print(f"  [CURRENT RUN: ABSENT] Field is completely absent from the targeted raw CAI asset '{target_json_filename}'")
            except Exception:
                pass
        else:
            # Fallback to recent JSON files if target step JSON is not found
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
                        break
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
            print(f"  --> The field is systematically ABSENT from BOTH current and historical CAI assets.")
            print(f"  --> This proves a permanent omission on the CAI side.")
            print("\nRECOMMENDED ACTION:")
            print("  >> Follow the strict logical Debugging Order of Operations before applying any fix:")
            print("     1. Check if the resource is handwritten or generated.")
            print("     2. If handwritten, check if Go-level conversion support (schema/flattener/expander)")
            print(f"        for '{field}' exists in tgc_next. If missing, you MUST implement Go support first!")
            print("        (Do not skip Go-level support to avoid masking missing features with false-positive ignores).")
            print(f"     3. If the field is still missing after Go support is implemented, set 'is_missing_in_cai: true' on '{camel_field}' in the Magic Modules YAML configuration.")
        else:
            if missing_in_all_local_cai:
                print(f"DIAGNOSTIC RESULT for '{field}':")
                print(f"  --> The field is specified in the HCL configuration but is completely absent from the local pre-recorded/mocked CAI asset file.")
                print(f"  --> (Note: It was found in some historical/live runs, suggesting it is supported by the live API but missing from this specific mock asset).")
                print("\nRECOMMENDED ACTION:")
                print("  >> Option A (Recommended): Follow the strict logical Debugging Order of Operations before applying any fix:")
                print("     1. Check if the resource is handwritten or generated.")
                print("     2. If handwritten, check if Go-level conversion support (schema/flattener/expander)")
                print(f"        for '{field}' exists in tgc_next. If missing, you MUST implement Go support first!")
                print(f"     3. If the field is still missing after Go support is implemented, set 'is_missing_in_cai: true' on '{camel_field}' in the Magic Modules YAML configuration.")
                print("  >> Option B: If you have active live GCP credentials, run the test with live connection to regenerate and update the local CAI asset files.")
            else:
                print(f"DIAGNOSTIC RESULT for '{field}':")
                print(f"  --> The field exists in the current CAI asset, but was lost during HCL conversion.")
                print("\nRECOMMENDED ACTION:")
                print(f"  >> Inspect 'cai2hcl' flatteners, encoders, or decoders in mmv1/templates/tgc_next/.")
                print(f"  >> A custom flattener/decoder may be required to map the field to HCL.")
            
    print("="*60 + "\n")

if __name__ == "__main__":
    main()
