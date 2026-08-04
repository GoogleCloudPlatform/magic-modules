#!/usr/bin/env python3
import sys
import re
import json
import os
import argparse
import time
from typing import List, Dict, Any, Optional

class TFDebugParser:
    def __init__(self, filepath: str):
        self.filepath = filepath
        self.events = []
        self.failures = []
        self.diffs = []
        self.outline = []
        self.test_name = "unknown_test"
        
        # We want an opt-out approach for the outline to preserve unexpected logs by default.
        self.outline_exclusions = [
            re.compile(r'\[DEBUG\] sdk\.helper_resource: (?!Starting TestStep:)'),
            re.compile(r'\[DEBUG\] sdk\.helper_schema:'),
            re.compile(r'\[DEBUG\] Retry Transport:'),
            re.compile(r'\[DEBUG\] Waiting for state to become:'),
            re.compile(r'Authenticating using configured Google JSON'),
            re.compile(r'-- Scopes: \['),
            re.compile(r'\[DEBUG\] VCR_PATH'),
            re.compile(r'\[DEBUG\] Google API Request Details:'),
            re.compile(r'\[DEBUG\] Google API Response Details:'),
            re.compile(r'\[DEBUG\]TGC Terraform metadata:'),
            re.compile(r'\[DEBUG\]TGC Terraform error:'),
            re.compile(r'\[DEBUG\] CaiAssetNames'),
            re.compile(r'\[DEBUG\] tgcPayload'),
            re.compile(r'\[WARN\]  sdk\.helper_schema:'),
            re.compile(r'\[DEBUG\] setting computed for'),
            re.compile(r'==> Checking that code complies'),
            re.compile(r'go vet'),
        ]
        
        self.chunk_separators = [
            re.compile(r'^=== RUN'), 
            re.compile(r'\[DEBUG\] sdk\.helper_resource: Starting TestStep:'),
            re.compile(r'\[DEBUG\] Creating new '),
            re.compile(r'\[DEBUG\] Finished creating '),
            re.compile(r'\[DEBUG\] Deleting '),
            re.compile(r'\[DEBUG\] Finished deleting '),
            re.compile(r'\[DEBUG\]test_step_number=\d+ TGC Terraform metadata:'),
            re.compile(r'\[ERROR\]'),
            re.compile(r'^--- (FAIL|PASS)'),
        ]
        self.identity_re = re.compile(r'\[INFO\] Terraform is using this identity:')
        self.last_was_identity = False
        self.identity_seen = False
        
    def should_include_in_outline(self, line: str) -> bool:
        for ex in self.outline_exclusions:
            if ex.search(line):
                return False
                
        # Deduplicate identity logs
        if self.identity_re.search(line):
            if self.identity_seen:
                return False
            self.identity_seen = True
            
        return True
        
    def parse(self):
        with open(self.filepath, 'r') as f:
            lines = f.readlines()
            
        current_state = None
        current_event: Dict[str, Any] = {}
        current_buffer = []
        
        diff_state = False
        diff_buffer = []

        # basic regexes for extraction
        req_start_re = re.compile(r'---\[ REQUEST \]---')
        res_start_re = re.compile(r'---\[ RESPONSE \]---')
        end_re = re.compile(r'-----------------------------------------------------')
        
        # When we are in a payload block, we omit it from the outline, 
        # and instead inject a marker when the block ends.
        in_payload_block = False
        
        googleapi_error_state = False
        warn_error_state = False

        # basic regexes for extraction
        req_start_re = re.compile(r'---\[ REQUEST ]---')
        res_start_re = re.compile(r'---\[ RESPONSE ]---')
        end_re = re.compile(r'-----------------------------------------------------')
        
        method_url_re = re.compile(r'^(GET|POST|PUT|DELETE|PATCH) (.*) HTTP/\d\.\d')
        status_re = re.compile(r'^HTTP/\d\.\d (\d{3} .*)')
        fail_re = re.compile(r'(FAIL:|Error:|Step \d+/\d+ error:)')
        diff_start_re = re.compile(r'Terraform used the selected providers to generate the following execution|Terraform will perform the following actions|plan. Resource actions are indicated with the following symbols:')
        diff_end_re = re.compile(r'Plan: \d+ to add, \d+ to change, \d+ to destroy.')

        for i, line in enumerate(lines):
            clean_line = line.strip()
            
            # 0. Check for Test Name Execution
            if "=== RUN   TestAcc" in clean_line:
                self.test_name = clean_line.split(" ")[-1]
            
            # 1. Check for Failures/Errors
            if fail_re.search(clean_line) and not diff_state and current_state is None:
                self.failures.append({"line": i+1, "message": clean_line})

            # 2. Check for Diff Block
            if diff_start_re.search(clean_line) and not diff_state:
                diff_state = True
                diff_buffer = [clean_line]
                continue
            
            if diff_state:
                diff_buffer.append(line.rstrip('\r\n'))
                if diff_end_re.search(clean_line):
                    diff_state = False
                    diff_text = "\n".join(diff_buffer)
                    
                    # Deduplicate consecutive identical diffs (ignoring pipes/whitespace)
                    raw_current = re.sub(r'[\s\|]', '', diff_text)
                    is_duplicate = False
                    if self.diffs and re.sub(r'[\s\|]', '', self.diffs[-1]) == raw_current:
                        is_duplicate = True
                        if '|' in self.diffs[-1] and '|' not in diff_text:
                            self.diffs[-1] = diff_text
                            if self.events and self.events[-1]["type"] == "diff":
                                self.events[-1]["content"] = diff_text
                                
                            # Find the last diff marker in the outline and replace the subsequent text
                            # This is a bit complex, so instead of trying to perfectly replace it, 
                            # we'll just let the deduplication prefer the better text in the payloads directory.
                            # The outline will just contain whatever the first diff was (usually with pipes).
                    
                    if not is_duplicate:
                        self.diffs.append(diff_text)
                        self.events.append({"type": "diff", "content": diff_text, "line": i+1})
                        
                        # Add diff marker sequence to outline
                        event_index = len(self.events)
                        prefix = f"{event_index:02d}_DIFF.txt"
                        if self.outline and self.outline[-1] != "\n":
                            self.outline.append("\n")
                        self.outline.append(f"-> [{prefix}]\n")
                        self.outline.append(f"{diff_text}\n")
                        self.last_was_identity = False
                continue


            # 3. Check for Request/Response blocks
            if req_start_re.search(clean_line):
                current_state = "REQUEST"
                in_payload_block = True
                current_event = {"type": "request", "line": i+1}
                current_buffer = []
                continue
                
            if res_start_re.search(clean_line):
                current_state = "RESPONSE"
                in_payload_block = True
                current_event = {"type": "response", "line": i+1}
                current_buffer = []
                continue
                
            if end_re.search(clean_line) and current_state is not None:
                # Process the buffer we just collected
                if current_state == "REQUEST":
                    # Parse Method and URL
                    for j, bline in enumerate(current_buffer):
                        m = method_url_re.search(bline)
                        if m:
                            current_event["method"] = m.group(1)
                            current_event["url"] = m.group(2).split('?')[0] # strip query for cleaner overview
                            break
                elif current_state == "RESPONSE":
                    # Parse Status Code
                    for j, bline in enumerate(current_buffer):
                        m = status_re.search(bline)
                        if m:
                            current_event["status"] = m.group(1)
                            break
                
                # Try to extract JSON body if it exists
                body_start = -1
                for j, bline in enumerate(current_buffer):
                    if bline.startswith('{') or bline.startswith('['):
                        body_start = j
                        break
                
                if body_start != -1:
                    raw_body = "\n".join(current_buffer[body_start:])
                    try:
                        current_event["body"] = json.loads(raw_body)
                    except json.JSONDecodeError:
                        current_event["body"] = raw_body # fallback to string

                self.events.append(current_event)
                
                # Add marker to outline
                event_index = len(self.events)
                prefix = f"{event_index:02d}_{current_event['type'].upper()}"
                if current_event["type"] == "request":
                    filename = f"{prefix}_{current_event.get('method', 'UNKNOWN')}.json"
                    self.outline.append(f"    -> [{filename}]\n")
                else:
                    filename = f"{prefix}_{current_event.get('status', 'UNKNOWN').replace(' ', '_')}.json"
                    self.outline.append(f"    <- [{filename}]\n")

                current_state = None
                in_payload_block = False
                current_event = {}
                self.last_was_identity = False
                continue
                
            if current_state is not None:
                current_buffer.append(clean_line)
                
            # 4. Check for unformatted googleapi HTTP dump noise
            if "googleapi: got HTTP response code" in clean_line:
                googleapi_error_state = True
                continue
            if googleapi_error_state:
                # The payload is over once we hit a standard log line with a timestamp
                if re.match(r'^\d{4}[/-]\d{2}[/-]\d{2}[T\s]', clean_line):
                    googleapi_error_state = False
                else:
                    continue
                    
            # 5. Check for redundant [WARN] error block
            if "[WARN]  sdk.helper_resource: Error running Terraform CLI command:" in clean_line:
                warn_error_state = True
                continue
            if warn_error_state:
                # The redundant block ends when the actual [ERROR] block starts
                if "[ERROR]" in clean_line:
                    warn_error_state = False
                else:
                    continue
                
            # 6. Add line to outline if it isn't an exclusion and we aren't inside a raw payload block
            if not in_payload_block and not diff_state and self.should_include_in_outline(clean_line):
                # Don't add the start/end lines of diffs since we inject the marker
                if not diff_start_re.search(clean_line):
                    is_identity = bool(self.identity_re.search(clean_line))
                    
                    needs_newline = False
                    for sep in self.chunk_separators:
                        if sep.search(clean_line):
                            needs_newline = True
                            break
                            
                    if is_identity and not self.last_was_identity:
                        needs_newline = True
                        
                    if needs_newline and self.outline and self.outline[-1] != "\n":
                        self.outline.append("\n")
                        
                    formatted_line = clean_line
                    
                    # 1. Strip timestamps
                    timestamp_re = re.compile(r'^\d{4}[/-]\d{2}[/-]\d{2}[T\s]\d{2}:\d{2}:\d{2}(?:\.\d{3}[+-]\d{4})?\s+')
                    formatted_line = timestamp_re.sub('', formatted_line)
                    
                    # 2. Truncate struct data dumps
                    map_dump_re = re.compile(r':\s*map\[string\]interface\s*\{\}.*$')
                    if map_dump_re.search(formatted_line):
                        formatted_line = map_dump_re.sub('...', formatted_line)
                        
                    # 3. Format Test Steps Visually
                    step_re = re.compile(r'\[DEBUG\] sdk\.helper_resource: Starting TestStep:.*?test_step_number=(\d+).*')
                    m = step_re.match(formatted_line)
                    if m:
                        step_num = m.group(1)
                        formatted_line = f">>> STARTING TEST STEP {step_num} <<<"
                        
                    # 4. Strip redundant error noise
                    # The lines immediately following a failure contain duplicate info
                    if formatted_line.startswith("test_step_number="):
                        continue
                    if "error: After applying this test step" in formatted_line:
                        continue
                    if formatted_line.strip() == "stdout:":
                        continue
                        
                    self.outline.append(f"{formatted_line}\n")
                    self.last_was_identity = is_identity
                
    def print_summary(self):
        print("====== API TIMELINE ======")
        for i, ev in enumerate(self.events):
            if ev["type"] == "request":
                print(f"[{i+1}] REQ: {ev.get('method', 'UNKNOWN')} {ev.get('url', 'UNKNOWN')}")
            elif ev["type"] == "diff":
                print(f"[{i+1}] DIFF: Terraform Plan Diff")
            else:
                status = ev.get('status', 'UNKNOWN')
                body_keys = "No Body"
                if "body" in ev and isinstance(ev["body"], dict):
                    body_keys = f"Keys: {', '.join(ev['body'].keys())}"
                print(f"[{i+1}] RES: Status {status} | {body_keys}")
                
        print("\n====== FAILURES ======")
        for f in self.failures:
            print(f"Line {f['line']}: {f['message']}")
            
        print("\n====== LATEST DIFF ======")
        if self.diffs:
            print(self.diffs[-1])
        else:
            print("No plan diff found.")

    def export_dir(self, base_dir: str):
        timestamp = int(time.time())
        out_dir = os.path.join(base_dir, f"{self.test_name}_{timestamp}")
        
        if not os.path.exists(out_dir):
            os.makedirs(out_dir)
            
        for i, ev in enumerate(self.events):
            prefix = f"{i+1:02d}_{ev['type'].upper()}"
            if ev["type"] == "request":
                filename = f"{prefix}_{ev.get('method', 'UNKNOWN')}.json"
                filepath = os.path.join(out_dir, filename)
                with open(filepath, 'w') as f:
                    json.dump(ev, f, indent=2)
            elif ev["type"] == "diff":
                filename = f"{prefix}.txt"
                filepath = os.path.join(out_dir, filename)
                with open(filepath, 'w') as f:
                    f.write(ev.get("content", ""))
            else:
                filename = f"{prefix}_{ev.get('status', 'UNKNOWN').replace(' ', '_')}.json"
                filepath = os.path.join(out_dir, filename)
                with open(filepath, 'w') as f:
                    json.dump(ev, f, indent=2)
                
        # Write outline
        with open(os.path.join(out_dir, "outline.txt"), 'w') as f:
            for line in self.outline:
                f.write(line)
                
        # Write failures
        with open(os.path.join(out_dir, "failures.json"), 'w') as f:
            json.dump(self.failures, f, indent=2)
            
        print(f"Extracted API timeline and errors to {out_dir}/")
        return json.dumps({
            "timeline": self.events,
            "failures": self.failures,
            "diffs": self.diffs
        }, indent=2)

def main():
    parser = argparse.ArgumentParser(description="Parse Terraform Debug Logs for API interaction and errors.")
    parser.add_argument("logfile", help="Path to the TF_LOG=DEBUG output file")
    parser.add_argument("--summary", action="store_true", help="Print a human-readable summary to stdout")
    parser.add_argument("--json", action="store_true", help="Output the full parsed structured timeline as JSON")
    
    parser.add_argument("--extract-dir", help="Directory to write extracted API JSON files and errors to")
    
    args = parser.parse_args()
    
    if not os.path.exists(args.logfile):
        print(f"Error: File {args.logfile} not found.")
        sys.exit(1)
        
    tf_parser = TFDebugParser(args.logfile)
    tf_parser.parse()
    
    if args.extract_dir:
        tf_parser.export_dir(args.extract_dir)
    elif args.json:
        print(tf_parser.export_json())
    else:
        tf_parser.print_summary()

if __name__ == "__main__":
    main()
