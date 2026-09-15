#!/usr/bin/env python3
"""
Helper script to fetch and filter nightly test results from GCS test-metadata.
"""

import argparse
import datetime
import json
import os
import subprocess
import sys

def get_date_str(days_ago):
    today = datetime.date.today()
    target_date = today - datetime.timedelta(days=days_ago)
    return target_date.strftime("%Y-%m-%d")

def is_permission_error(stderr_text, stdout_text=""):
    combined = (stderr_text + " " + stdout_text).lower()
    auth_indicators = [
        "permission denied",
        "permissiondenied",
        "403",
        "401",
        "access_token_scope_insufficient",
        "invalid_grant",
        "unauthenticated",
        "invalid authentication credentials",
        "does not have storage.objects.get access",
        "could not refresh access token",
        "insufficient permissions",
        "refresherror",
    ]
    return any(indicator in combined for indicator in auth_indicators)

def fetch_metadata(provider, date_str):
    uri = f"gs://nightly-test-data/test-metadata/{provider}/{date_str}-{provider}.json"
    res = subprocess.run(["gcloud", "storage", "cat", uri], capture_output=True, text=True)
    if res.returncode != 0:
        if is_permission_error(res.stderr, res.stdout):
            error_details = (res.stderr or res.stdout or "").strip()
            print(f"\n[ERROR] Permission or authentication failure while fetching {uri}:", file=sys.stderr)
            if error_details:
                print(f"{error_details}\n", file=sys.stderr)
            print("Please check your gcloud authentication ('gcloud auth login') and ensure you have read access ('roles/storage.objectViewer') to gs://nightly-test-data.", file=sys.stderr)
            sys.exit(1)
        return None
    try:
        return json.loads(res.stdout)
    except json.JSONDecodeError:
        return None

def main():
    parser = argparse.ArgumentParser(description="Fetch nightly test metadata from GCS.")
    parser.add_argument("--provider", choices=["ga", "beta", "both"], default="both", help="Provider version (ga, beta, or both)")
    parser.add_argument("--date", help="Specific date YYYY-MM-DD (defaults to latest available)")
    parser.add_argument("--days", type=int, default=1, help="Number of past days to check (default: 1)")
    parser.add_argument("--status", choices=["FAILURE", "SUCCESS", "ALL"], default="FAILURE", help="Filter by test status (default: FAILURE)")
    parser.add_argument("--test", help="Filter by test name substring")
    parser.add_argument("--json", action="store_true", help="Output raw JSON format")

    args = parser.parse_args()

    providers = ["beta", "ga"] if args.provider == "both" else [args.provider]
    results = []

    for prov in providers:
        if args.date:
            data = fetch_metadata(prov, args.date)
            date_used = args.date
            if data:
                for item in data:
                    item["provider"] = prov
                    item["date"] = date_used
                    results.append(item)
        else:
            # Check backwards for latest available days
            found_days = 0
            for day_offset in range(14):
                if found_days >= args.days:
                    break
                d_str = get_date_str(day_offset)
                data = fetch_metadata(prov, d_str)
                if data:
                    found_days += 1
                    for item in data:
                        item["provider"] = prov
                        item["date"] = d_str
                        results.append(item)

    # Filter results
    filtered = []
    for item in results:
        if args.status != "ALL" and item.get("status") != args.status:
            continue
        if args.test and args.test.lower() not in item.get("name", "").lower():
            continue
        filtered.append(item)

    if args.json:
        print(json.dumps(filtered, indent=2))
    else:
        print(f"Total matching tests: {len(filtered)}\n")
        print(f"| Test Name | Provider | Date | Status | Duration (s) | Log Link | Error Snippet |")
        print(f"| --- | --- | --- | --- | --- | --- | --- |")
        for item in filtered:
            name = item.get("name")
            prov = item.get("provider").upper()
            d_str = item.get("date")
            st = item.get("status")
            dur = f"{item.get('duration', 0) / 1000.0:.1f}s"
            log = f"[Log]({item.get('log_link')})" if item.get("log_link") else "N/A"
            err = item.get("error_message", "").replace("\n", " ")[:100]
            print(f"| `{name}` | {prov} | {d_str} | **{st}** | {dur} | {log} | `{err}` |")

if __name__ == "__main__":
    main()
