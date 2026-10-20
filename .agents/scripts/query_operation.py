#!/usr/bin/env python3
"""Query a Google Cloud Long-Running Operation (LRO) URL and report its status and duration.

Validates the target URL against a strict googleapis.com allowlist and uses
subprocess.run(..., shell=False) to obtain an access token without shell expansion.
"""

import argparse
import datetime
import json
import re
import subprocess
import sys
import urllib.error
import urllib.request

OPERATION_URL_REGEX = re.compile(
    r"^https://[a-z0-9.-]*googleapis\.com/[a-zA-Z0-9_.\-/]+$"
)


def validate_operation_url(url: str) -> str:
    candidate = (url or "").strip()
    if not OPERATION_URL_REGEX.match(candidate) or ".." in candidate.split("/"):
        raise ValueError(
            f"Invalid operation URL '{candidate}'. Must match {OPERATION_URL_REGEX.pattern}."
        )
    return candidate


def parse_iso_timestamp(ts: str) -> datetime.datetime:
    # Normalize trailing Z for Python 3.9/3.10 compatibility
    if ts.endswith("Z"):
        ts = ts[:-1] + "+00:00"
    # Truncate nanoseconds to microseconds if present
    ts = re.sub(r"(\.\d{6})\d+", r"\1", ts)
    return datetime.datetime.fromisoformat(ts)


def extract_timestamps(payload: dict) -> tuple:
    metadata = payload.get("metadata") if isinstance(payload.get("metadata"), dict) else {}
    start_ts = (
        metadata.get("createTime")
        or payload.get("startTime")
        or payload.get("insertTime")
        or payload.get("createTime")
    )
    end_ts = (
        metadata.get("endTime")
        or payload.get("endTime")
        or (metadata.get("updateTime") if payload.get("done") is True else None)
    )
    return start_ts, end_ts


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Query a GCP operation URL and report its status and total duration."
    )
    parser.add_argument("operation_url", help="Full https://*.googleapis.com/... operation URL")
    args = parser.parse_args()

    try:
        op_url = validate_operation_url(args.operation_url)
    except ValueError as err:
        print(f"[ERROR] {err}", file=sys.stderr)
        sys.exit(1)

    try:
        token_proc = subprocess.run(
            ["gcloud", "auth", "print-access-token"],
            check=True,
            capture_output=True,
            text=True,
            shell=False,
        )
        token = token_proc.stdout.strip()
    except subprocess.CalledProcessError as err:
        print(
            f"[ERROR] Failed to obtain gcloud access token: {(err.stderr or '').strip()}",
            file=sys.stderr,
        )
        sys.exit(1)

    req = urllib.request.Request(
        op_url,
        headers={
            "Authorization": f"Bearer {token}",
            "Accept": "application/json",
        },
    )

    try:
        with urllib.request.urlopen(req) as resp:
            payload = json.load(resp)
    except urllib.error.HTTPError as err:
        body = err.read().decode("utf-8", errors="replace")
        print(f"[ERROR] HTTP {err.code} querying {op_url}:\n{body}", file=sys.stderr)
        sys.exit(1)
    except Exception as err:
        print(f"[ERROR] Failed to query {op_url}: {err}", file=sys.stderr)
        sys.exit(1)

    print(json.dumps(payload, indent=2))

    start_ts, end_ts = extract_timestamps(payload)
    if start_ts and end_ts:
        try:
            start_dt = parse_iso_timestamp(start_ts)
            end_dt = parse_iso_timestamp(end_ts)
            elapsed_sec = (end_dt - start_dt).total_seconds()
            print(
                f"\n# Duration: {elapsed_sec:.1f}s ({elapsed_sec / 60.0:.2f}m) "
                f"[{start_ts} -> {end_ts}]"
            )
        except Exception:
            pass


if __name__ == "__main__":
    main()
