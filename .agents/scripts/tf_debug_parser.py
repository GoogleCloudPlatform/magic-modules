#!/usr/bin/env python3
import argparse
import json
import os
import re
import sys
import time
import urllib.parse
from typing import Any, Dict, List, Optional

NOISE_QUERY_PARAMS = frozenset({"alt", "prettyPrint"})


def parse_request_target(raw_target: str) -> Dict[str, Any]:
    """Parses an HTTP request target losslessly into path, full URL, and query metadata."""
    raw_target = raw_target.strip()
    split = urllib.parse.urlsplit(raw_target)
    if split.scheme:
        url_path = f"{split.scheme}://{split.netloc}{split.path}"
    else:
        url_path = split.path

    query_params = urllib.parse.parse_qs(split.query, keep_blank_values=True)
    result: Dict[str, Any] = {
        "url": url_path,
        "url_path": url_path,
        "url_full": raw_target,
        "query_params": query_params,
    }

    mask_values: Optional[List[str]] = None
    if "updateMask" in query_params:
        mask_values = query_params["updateMask"]
    elif "update_mask" in query_params:
        mask_values = query_params["update_mask"]

    if mask_values is not None:
        result["update_mask"] = [
            field.strip()
            for raw_mask in mask_values
            for field in raw_mask.split(",")
            if field.strip()
        ]

    return result


def format_concise_url(event: Dict[str, Any]) -> str:
    """Formats a clean URL path while preserving meaningful query params like updateMask."""
    url_path = event.get("url_path") or event.get("url", "UNKNOWN")
    url_full = event.get("url_full", "")
    if not url_full or "?" not in url_full:
        return url_path

    split = urllib.parse.urlsplit(url_full)
    pairs = urllib.parse.parse_qsl(split.query, keep_blank_values=True)
    kept_parts: List[str] = []
    for key, val in pairs:
        if key in NOISE_QUERY_PARAMS:
            continue
        if val == "" and "=" not in split.query:
            kept_parts.append(key)
        else:
            kept_parts.append(f"{key}={val}")

    if not kept_parts:
        return url_path
    return f"{url_path}?{'&'.join(kept_parts)}"


def format_phase_prefix(event: Dict[str, Any]) -> str:
    """Returns a concise human-readable step or phase prefix for summary output."""
    step = event.get("test_step")
    if step is not None:
        return f"[Step {step}] "
    if event.get("phase") == "destroy":
        return "[Destroy] "
    return ""


class TFDebugParser:
    def __init__(self, filepath: str):
        self.filepath = filepath
        self.events: List[Dict[str, Any]] = []
        self.failures: List[Dict[str, Any]] = []
        self.diffs: List[str] = []
        self.steps: List[Dict[str, Any]] = []
        self.outline: List[str] = []
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
        with open(self.filepath, 'r', encoding='utf-8', errors='replace') as f:
            lines = f.readlines()

        current_state: Optional[str] = None
        current_event: Dict[str, Any] = {}
        current_buffer: List[str] = []

        current_step: Optional[int] = None
        current_phase = "setup"
        destroy_banner_emitted = False

        diff_state = False
        diff_buffer: List[str] = []

        # When we are in a payload block, we omit it from the outline,
        # and instead inject a marker when the block ends.
        in_payload_block = False

        googleapi_error_state = False
        warn_error_state = False

        req_start_re = re.compile(r'---\[ REQUEST \]---')
        res_start_re = re.compile(r'---\[ RESPONSE \]---')
        end_re = re.compile(r'-----------------------------------------------------')

        method_url_re = re.compile(r'^(GET|POST|PUT|DELETE|PATCH) (.*) HTTP/\d\.\d')
        status_re = re.compile(r'^HTTP/\d\.\d (\d{3} .*)')
        fail_re = re.compile(r'(FAIL:|Error:|Step \d+/\d+ error:)')
        diff_start_re = re.compile(
            r'Terraform used the selected providers to generate the following execution|'
            r'Terraform will perform the following actions|'
            r'plan\. Resource actions are indicated with the following symbols:'
        )
        diff_end_re = re.compile(r'Plan: \d+ to add, \d+ to change, \d+ to destroy\.')

        step_start_re = re.compile(
            r'\[DEBUG\] sdk\.helper_resource: Starting TestStep:.*?test_step_number=(\d+)'
        )
        step_fallback_re = re.compile(r'\btest_step_number=(\d+)\b')
        step_finish_re = re.compile(r'\[DEBUG\] sdk\.helper_resource: Finished TestStep')
        step_abort_re = re.compile(
            r'\[ERROR\] sdk\.helper_resource: (?:Unexpected error|Expected an error)'
        )
        destroy_marker_re = re.compile(
            r'\[(?:DEBUG|TRACE)\] sdk\.helper_resource: Calling (?:TestCase CheckDestroy|Terraform CLI destroy command)|'
            r'Error running post-test destroy'
        )
        timestamp_re = re.compile(
            r'^\d{4}[/-]\d{2}[/-]\d{2}[T\s]\d{2}:\d{2}:\d{2}(?:\.\d{3}[+-]\d{4})?\s+'
        )
        log_ts_re = re.compile(r'^\d{4}[/-]\d{2}[/-]\d{2}[T\s]')
        map_dump_re = re.compile(r':\s*map\[string\]interface\s*\{\}.*$')

        for i, line in enumerate(lines):
            clean_line = line.strip()

            # 0. Check for Test Name Execution and TestStep / Destroy Phase Transitions
            if "=== RUN   TestAcc" in clean_line:
                self.test_name = clean_line.split()[-1]

            step_match = step_start_re.search(clean_line)
            if step_match:
                current_step = int(step_match.group(1))
                current_phase = f"step_{current_step}"
                destroy_banner_emitted = False
                self.steps.append({"step_number": current_step, "line": i + 1})
            elif current_step is None and current_phase == "setup":
                fallback_match = step_fallback_re.search(clean_line)
                if fallback_match:
                    current_step = int(fallback_match.group(1))
                    current_phase = f"step_{current_step}"
                    self.steps.append({"step_number": current_step, "line": i + 1})

            if step_finish_re.search(clean_line) or destroy_marker_re.search(clean_line):
                current_step = None
                current_phase = "destroy"

            # 1. Check for Failures/Errors
            if fail_re.search(clean_line) and not diff_state and current_state is None:
                self.failures.append({
                    "line": i + 1,
                    "message": clean_line,
                    "test_step": current_step,
                    "phase": current_phase,
                })

            if step_abort_re.search(clean_line):
                current_step = None
                current_phase = "destroy"

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

                    if not is_duplicate:
                        self.diffs.append(diff_text)
                        self.events.append({
                            "type": "diff",
                            "content": diff_text,
                            "line": i + 1,
                            "test_step": current_step,
                            "phase": current_phase,
                        })

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
                current_event = {
                    "type": "request",
                    "line": i + 1,
                    "test_step": current_step,
                    "phase": current_phase,
                }
                current_buffer = []
                continue

            if res_start_re.search(clean_line):
                current_state = "RESPONSE"
                in_payload_block = True
                current_event = {
                    "type": "response",
                    "line": i + 1,
                    "test_step": current_step,
                    "phase": current_phase,
                }
                current_buffer = []
                continue

            if end_re.search(clean_line) and current_state is not None:
                # Process the buffer we just collected
                if current_state == "REQUEST":
                    for bline in current_buffer:
                        m = method_url_re.search(bline)
                        if m:
                            current_event["method"] = m.group(1)
                            current_event.update(parse_request_target(m.group(2)))
                            break
                elif current_state == "RESPONSE":
                    for bline in current_buffer:
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
                        current_event["body"] = raw_body  # fallback to string

                if current_phase == "destroy" and not destroy_banner_emitted:
                    if self.outline and self.outline[-1] != "\n":
                        self.outline.append("\n")
                    self.outline.append(">>> POST-TEST DESTROY <<<\n")
                    destroy_banner_emitted = True

                self.events.append(current_event)

                # Add marker to outline
                event_index = len(self.events)
                prefix = f"{event_index:02d}_{current_event['type'].upper()}"
                if current_event["type"] == "request":
                    method = current_event.get("method", "UNKNOWN")
                    filename = f"{prefix}_{method}.json"
                    concise_url = format_concise_url(current_event)
                    self.outline.append(f"    -> [{filename}] {method} {concise_url}\n")
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
                if log_ts_re.match(clean_line):
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

                    formatted_line = timestamp_re.sub('', clean_line)

                    # Truncate struct data dumps
                    if map_dump_re.search(formatted_line):
                        formatted_line = map_dump_re.sub('...', formatted_line)

                    # Format Test Steps Visually
                    m = step_start_re.match(formatted_line)
                    if m:
                        step_num = m.group(1)
                        formatted_line = f">>> STARTING TEST STEP {step_num} <<<"

                    # Strip redundant error noise
                    if formatted_line.startswith("test_step_number="):
                        continue
                    if "error: After applying this test step" in formatted_line:
                        continue
                    if formatted_line.strip() == "stdout:":
                        continue

                    self.outline.append(f"{formatted_line}\n")
                    self.last_was_identity = is_identity

    def to_dict(self) -> Dict[str, Any]:
        return {
            "test_name": self.test_name,
            "steps": self.steps,
            "timeline": self.events,
            "failures": self.failures,
            "diffs": self.diffs,
        }

    def export_json(self) -> str:
        return json.dumps(self.to_dict(), indent=2)

    def print_summary(self):
        print("====== API TIMELINE ======")
        for i, ev in enumerate(self.events):
            phase_prefix = format_phase_prefix(ev)
            if ev["type"] == "request":
                concise_url = format_concise_url(ev)
                print(f"[{i+1}] {phase_prefix}REQ: {ev.get('method', 'UNKNOWN')} {concise_url}")
            elif ev["type"] == "diff":
                print(f"[{i+1}] {phase_prefix}DIFF: Terraform Plan Diff")
            else:
                status = ev.get('status', 'UNKNOWN')
                body_keys = "No Body"
                if "body" in ev and isinstance(ev["body"], dict):
                    body_keys = f"Keys: {', '.join(ev['body'].keys())}"
                print(f"[{i+1}] {phase_prefix}RES: Status {status} | {body_keys}")

        print("\n====== FAILURES ======")
        for f in self.failures:
            phase_prefix = format_phase_prefix(f)
            print(f"Line {f['line']}: {phase_prefix}{f['message']}")

        print("\n====== LATEST DIFF ======")
        if self.diffs:
            print(self.diffs[-1])
        else:
            print("No plan diff found.")

    def export_dir(self, base_dir: str) -> str:
        timestamp = int(time.time())
        out_dir = os.path.join(base_dir, f"{self.test_name}_{timestamp}")

        if not os.path.exists(out_dir):
            os.makedirs(out_dir)

        for i, ev in enumerate(self.events):
            prefix = f"{i+1:02d}_{ev['type'].upper()}"
            if ev["type"] == "request":
                filename = f"{prefix}_{ev.get('method', 'UNKNOWN')}.json"
                filepath = os.path.join(out_dir, filename)
                with open(filepath, 'w', encoding='utf-8') as f:
                    json.dump(ev, f, indent=2)
            elif ev["type"] == "diff":
                filename = f"{prefix}.txt"
                filepath = os.path.join(out_dir, filename)
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.write(ev.get("content", ""))
            else:
                filename = f"{prefix}_{ev.get('status', 'UNKNOWN').replace(' ', '_')}.json"
                filepath = os.path.join(out_dir, filename)
                with open(filepath, 'w', encoding='utf-8') as f:
                    json.dump(ev, f, indent=2)

        # Write outline
        with open(os.path.join(out_dir, "outline.txt"), 'w', encoding='utf-8') as f:
            for line in self.outline:
                f.write(line)

        # Write failures
        with open(os.path.join(out_dir, "failures.json"), 'w', encoding='utf-8') as f:
            json.dump(self.failures, f, indent=2)

        print(f"Extracted API timeline and errors to {out_dir}/")
        return self.export_json()


def main(argv: Optional[List[str]] = None) -> int:
    parser = argparse.ArgumentParser(description="Parse Terraform Debug Logs for API interaction and errors.")
    parser.add_argument("logfile", help="Path to the TF_LOG=DEBUG output file")
    parser.add_argument("--summary", action="store_true", help="Print a human-readable summary to stdout")
    parser.add_argument("--json", action="store_true", help="Output the full parsed structured timeline as JSON")
    parser.add_argument("--extract-dir", help="Directory to write extracted API JSON files and errors to")

    args = parser.parse_args(argv)

    if not os.path.exists(args.logfile):
        print(f"Error: File {args.logfile} not found.")
        return 1

    tf_parser = TFDebugParser(args.logfile)
    tf_parser.parse()

    if args.extract_dir:
        tf_parser.export_dir(args.extract_dir)
    elif args.json:
        print(tf_parser.export_json())
    else:
        tf_parser.print_summary()
    return 0


if __name__ == "__main__":
    sys.exit(main())
