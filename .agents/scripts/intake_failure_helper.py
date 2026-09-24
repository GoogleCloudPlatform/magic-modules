#!/usr/bin/env python3
"""Deterministic helper for ingesting and validating test failure information.

This script safely ingests test failure inputs from GitHub issues, GCS error/debug
logs, or local log files. To prevent command injection and prompt injection:
1. All external fields (test_name, target_provider, GCS URIs, GitHub issue URLs,
   local file paths) are validated against strict regex allowlists.
2. GCS URIs are restricted to the trusted CI bucket (gs://nightly-test-data/...)
   and untrusted GitHub issue prose is never written into error log files.
3. All external CLI invocations (gh, gcloud, tf_debug_parser.py) execute via
   subprocess.run(..., shell=False) using direct argument vectors.
4. Error messages and stack traces are isolated into local files under
   debug_output/<test_name>/raw_error.log rather than printed raw into the
   agent's prompt context.
"""

import argparse
import json
import os
import re
import subprocess
import sys
from typing import Any, Dict, List, Optional, Tuple

ALLOWED_PROVIDERS = frozenset({"ga", "beta", "both"})
TEST_NAME_REGEX = re.compile(r"^TestAcc[A-Za-z0-9_]+$")
GCS_URI_REGEX = re.compile(r"^gs://nightly-test-data/[a-zA-Z0-9_.\-/]+$")
HTTPS_GCS_REGEX = re.compile(
    r"^https://storage\.(?:googleapis|cloud\.google)\.com/(nightly-test-data|teamcity-logs)/([a-zA-Z0-9_.\-/]+)$"
)
GITHUB_ISSUE_URL_REGEX = re.compile(
    r"^https://github\.com/hashicorp/terraform-provider-google/issues/(\d+)$"
)
AUTH_ERROR_REGEX = re.compile(
    r"\b(?:"
    r"permission\s+denied|"
    r"permissiondenied|"
    r"403|"
    r"401|"
    r"access_token_scope_insufficient|"
    r"invalid_grant|"
    r"unauthenticated|"
    r"invalid\s+authentication\s+credentials|"
    r"does\s+not\s+have\s+storage\.objects\.get\s+access|"
    r"could\s+not\s+refresh\s+access\s+token|"
    r"insufficient\s+permissions|"
    r"refresherror"
    r")\b",
    re.IGNORECASE,
)


def is_permission_error(stderr_text: str, stdout_text: str = "") -> bool:
    combined = f"{stderr_text} {stdout_text}"
    return bool(AUTH_ERROR_REGEX.search(combined))


def validate_test_name(test_name: str) -> str:
    candidate = (test_name or "").strip()
    if not TEST_NAME_REGEX.match(candidate):
        raise ValueError(
            f"Invalid test name format '{candidate}'. Must match {TEST_NAME_REGEX.pattern}."
        )
    return candidate


def validate_provider(provider: str) -> str:
    candidate = (provider or "").strip().lower()
    if candidate not in ALLOWED_PROVIDERS:
        raise ValueError(
            f"Invalid target provider '{candidate}'. Must be one of {sorted(ALLOWED_PROVIDERS)}."
        )
    return candidate


def normalize_and_validate_gcs_uri(uri: str) -> str:
    candidate = (uri or "").strip()
    https_match = HTTPS_GCS_REGEX.match(candidate)
    if https_match:
        bucket, object_path = https_match.group(1), https_match.group(2)
        candidate = f"gs://{bucket}/{object_path}"

    if not GCS_URI_REGEX.match(candidate):
        raise ValueError(
            f"Invalid GCS URI '{uri}'. Must match {GCS_URI_REGEX.pattern}."
        )

    # Reject any path traversal or empty segments
    path_parts = candidate[len("gs://") :].split("/")
    if any(part in ("", ".", "..") for part in path_parts):
        raise ValueError(f"Invalid GCS URI path segments in '{uri}'.")

    return candidate


def validate_github_issue_url(issue_url: str) -> str:
    candidate = (issue_url or "").strip()
    if not GITHUB_ISSUE_URL_REGEX.match(candidate):
        raise ValueError(
            f"Invalid GitHub issue URL '{candidate}'. Expected a hashicorp/terraform-provider-google issue URL."
        )
    return candidate


def validate_local_path(filepath: str, base_dir: Optional[str] = None) -> str:
    candidate = (filepath or "").strip()
    if not candidate or "\x00" in candidate:
        raise ValueError("Invalid empty or null-byte file path.")
    resolved = os.path.realpath(candidate)
    resolved_base = os.path.realpath(base_dir if base_dir is not None else os.getcwd())
    if os.path.commonpath([resolved, resolved_base]) != resolved_base:
        raise ValueError(
            f"Path '{filepath}' resolves outside allowed directory '{resolved_base}'."
        )
    if not os.path.isfile(resolved):
        raise ValueError(f"Local log file does not exist: '{filepath}'.")
    return resolved


def fetch_gcs_content(gcs_uri: str) -> str:
    validated_uri = normalize_and_validate_gcs_uri(gcs_uri)
    res = subprocess.run(
        ["gcloud", "storage", "cat", validated_uri],
        shell=False,
        check=False,
        capture_output=True,
        text=True,
    )
    if res.returncode != 0:
        if is_permission_error(res.stderr, res.stdout):
            raise RuntimeError(
                f"Permission or authentication failure while fetching {validated_uri}. "
                "Please verify 'gcloud auth login' and ensure 'roles/storage.objectViewer' access."
            )
        raise RuntimeError(
            f"Failed to fetch GCS URI {validated_uri}: {(res.stderr or res.stdout).strip()}"
        )
    return res.stdout


def download_and_parse_debug_log(
    test_name: str,
    gcs_debug_uri: Optional[str] = None,
    local_debug_log: Optional[str] = None,
    output_dir: str = "debug_output",
    base_dir: Optional[str] = None,
) -> str:
    validated_test = validate_test_name(test_name)
    os.makedirs(output_dir, exist_ok=True)
    extract_dir = os.path.join(output_dir, validated_test)
    os.makedirs(extract_dir, exist_ok=True)

    if gcs_debug_uri:
        validated_uri = normalize_and_validate_gcs_uri(gcs_debug_uri)
        log_path = os.path.join(extract_dir, "raw_test.log")
        res = subprocess.run(
            ["gcloud", "storage", "cp", validated_uri, log_path],
            shell=False,
            check=False,
            capture_output=True,
            text=True,
        )
        if res.returncode != 0:
            if is_permission_error(res.stderr, res.stdout):
                raise RuntimeError(
                    f"Permission or authentication failure while copying {validated_uri}. "
                    "Please verify 'gcloud auth login' and ensure 'roles/storage.objectViewer' access."
                )
            raise RuntimeError(
                f"Failed to copy debug log {validated_uri}: {(res.stderr or res.stdout).strip()}"
            )
    elif local_debug_log:
        log_path = validate_local_path(local_debug_log, base_dir=base_dir)
    else:
        raise ValueError("Either gcs_debug_uri or local_debug_log must be provided.")

    script_dir = os.path.dirname(os.path.abspath(__file__))
    parser_script = os.path.join(script_dir, "tf_debug_parser.py")
    parse_res = subprocess.run(
        [sys.executable, parser_script, log_path, "--extract-dir", extract_dir],
        shell=False,
        check=False,
        capture_output=True,
        text=True,
    )
    if parse_res.returncode != 0:
        raise RuntimeError(
            f"Failed to parse debug log {log_path}: {(parse_res.stderr or parse_res.stdout).strip()}"
        )

    # tf_debug_parser.py writes into a nested <extract_dir>/<test_name>_<timestamp>/ directory
    # and prints: "Extracted API timeline and errors to <out_dir>/"
    match = re.search(
        r"Extracted API timeline and errors to\s+(.+?)/?\s*$",
        parse_res.stdout.strip(),
        re.MULTILINE,
    )
    if match:
        return match.group(1).rstrip("/") + "/"

    return extract_dir.rstrip("/") + "/"


def fetch_issue_payload(issue_url: str) -> Dict[str, Any]:
    validated_url = validate_github_issue_url(issue_url)
    cmd = ["gh", "issue", "view", validated_url, "--json", "title,body,labels"]
    res = subprocess.run(
        cmd,
        capture_output=True,
        text=True,
        check=False,
        shell=False,
    )
    if res.returncode != 0:
        raise RuntimeError(
            f"Failed to fetch GitHub issue {validated_url}: {(res.stderr or res.stdout).strip()}"
        )
    return json.loads(res.stdout)


def parse_failure_rates(body: str) -> Tuple[Optional[float], Optional[float]]:
    ga_rate: Optional[float] = None
    beta_rate: Optional[float] = None

    ga_match = re.search(r"\bGA:\s*(\d+(?:\.\d+)?)%", body, re.IGNORECASE)
    if ga_match:
        ga_rate = float(ga_match.group(1))

    beta_match = re.search(r"\bBeta:\s*(\d+(?:\.\d+)?)%", body, re.IGNORECASE)
    if beta_match:
        beta_rate = float(beta_match.group(1))

    return ga_rate, beta_rate


def extract_gcs_links_from_issue(body: str) -> Tuple[Dict[str, str], Dict[str, str]]:
    """Extracts validated error and debug log GCS URIs keyed by provider ('ga', 'beta')."""
    error_links: Dict[str, str] = {}
    debug_links: Dict[str, str] = {}

    url_pattern = re.compile(
        r"(?:gs://nightly-test-data/[a-zA-Z0-9_.\-/]+|https://storage\.(?:googleapis|cloud\.google)\.com/nightly-test-data/[a-zA-Z0-9_.\-/]+)"
    )

    for raw_url in url_pattern.findall(body):
        try:
            normalized = normalize_and_validate_gcs_uri(raw_url)
        except ValueError:
            continue

        lower_url = normalized.lower()
        if "/test-errors/ga/" in lower_url or (
            "error" in lower_url and "/ga/" in lower_url
        ):
            error_links.setdefault("ga", normalized)
        elif "/test-errors/beta/" in lower_url or (
            "error" in lower_url and "/beta/" in lower_url
        ):
            error_links.setdefault("beta", normalized)
        elif "debug" in lower_url or lower_url.endswith(".log"):
            if "/beta/" in lower_url:
                debug_links.setdefault("beta", normalized)
            else:
                debug_links.setdefault("ga", normalized)

    return error_links, debug_links


def validate_and_extract_issue(issue_data: Dict[str, Any]) -> Dict[str, Any]:
    title = issue_data.get("title") or ""
    body = issue_data.get("body") or ""
    labels_raw = issue_data.get("labels") or []
    labels = [
        lbl.get("name", "") if isinstance(lbl, dict) else str(lbl)
        for lbl in labels_raw
    ]

    has_test_failure_label = any(
        lbl.startswith("test-failure") for lbl in labels
    )
    if not has_test_failure_label:
        raise ValueError(
            "GitHub issue does not have a 'test-failure*' label; refusing to ingest non-test-failure issue."
        )

    # Extract test name from title first, then body
    test_match = re.search(r"\b(TestAcc[A-Za-z0-9_]+)\b", title) or re.search(
        r"\b(TestAcc[A-Za-z0-9_]+)\b", body
    )
    if not test_match:
        raise ValueError(
            "Could not find a valid TestAcc function name in issue title or body."
        )
    test_name = validate_test_name(test_match.group(1))

    ga_rate, beta_rate = parse_failure_rates(body)
    error_links, debug_links = extract_gcs_links_from_issue(body)

    if ga_rate is not None or beta_rate is not None:
        ga_failing = (ga_rate or 0.0) > 0.0
        beta_failing = (beta_rate or 0.0) > 0.0
        if ga_failing and beta_failing:
            target_provider = "both"
        elif beta_failing:
            target_provider = "beta"
        elif ga_failing:
            target_provider = "ga"
        else:
            target_provider = (
                "both"
                if ("ga" in error_links and "beta" in error_links)
                else ("beta" if "beta" in error_links else "ga")
            )
    else:
        if "ga" in error_links and "beta" in error_links:
            target_provider = "both"
        elif "beta" in error_links:
            target_provider = "beta"
        else:
            target_provider = "ga"

    target_provider = validate_provider(target_provider)

    return {
        "test_name": test_name,
        "target_provider": target_provider,
        "has_test_failure_label": has_test_failure_label,
        "labels": labels,
        "error_links": error_links,
        "debug_links": debug_links,
    }


def write_isolated_error_log(
    test_name: str, error_text: str, output_dir: str = "debug_output"
) -> str:
    validated_test = validate_test_name(test_name)
    test_dir = os.path.join(output_dir, validated_test)
    os.makedirs(test_dir, exist_ok=True)
    error_file_path = os.path.join(test_dir, "raw_error.log")
    with open(error_file_path, "w", encoding="utf-8") as f:
        f.write(error_text)
    return error_file_path


def run_cli(argv: Optional[List[str]] = None, base_dir: Optional[str] = None) -> Dict[str, Any]:
    parser = argparse.ArgumentParser(
        description="Deterministic, validated test-failure ingestion helper."
    )
    parser.add_argument(
        "--issue-url",
        help="GitHub issue URL (e.g., https://github.com/hashicorp/terraform-provider-google/issues/28244)",
    )
    parser.add_argument(
        "--test-name",
        help="Acceptance test function name (must match ^TestAcc[A-Za-z0-9_]+$)",
    )
    parser.add_argument(
        "--target-provider",
        choices=sorted(ALLOWED_PROVIDERS),
        default=None,
        help="Target provider version (ga, beta, or both)",
    )
    parser.add_argument(
        "--gcs-error-uri",
        action="append",
        default=[],
        help="GCS URI (gs://nightly-test-data/... or https://storage.googleapis.com/nightly-test-data/...) for error log file(s)",
    )
    parser.add_argument(
        "--gcs-debug-uri",
        help="GCS URI (gs://nightly-test-data/... or https://storage.googleapis.com/nightly-test-data/...) for TF_LOG=DEBUG log file",
    )
    parser.add_argument(
        "--local-log",
        help="Path to a local debug or error log file within the workspace",
    )
    parser.add_argument(
        "--parse-debug-log",
        action="store_true",
        help="Download and parse debug log via tf_debug_parser.py",
    )
    parser.add_argument(
        "--output-dir",
        default="debug_output",
        help="Base directory for isolated error and debug log outputs (default: debug_output)",
    )

    args = parser.parse_args(argv)

    if args.issue_url and (args.gcs_error_uri or args.local_log):
        raise ValueError(
            "--issue-url cannot be combined with --gcs-error-uri or --local-log."
        )

    error_chunks: List[str] = []
    parsed_logs_dir: Optional[str] = None

    if args.issue_url:
        issue_data = fetch_issue_payload(args.issue_url)
        extracted = validate_and_extract_issue(issue_data)
        test_name = validate_test_name(args.test_name or extracted["test_name"])
        target_provider = validate_provider(
            args.target_provider or extracted["target_provider"]
        )

        error_links: Dict[str, str] = extracted["error_links"]
        providers_to_fetch = (
            ["ga", "beta"]
            if target_provider == "both"
            else [target_provider]
        )
        for prov in providers_to_fetch:
            if prov in error_links:
                content = fetch_gcs_content(error_links[prov])
                error_chunks.append(f"=== [{prov.upper()} ERROR LOG] ===\n{content}")

        if not error_chunks:
            for prov, uri in error_links.items():
                content = fetch_gcs_content(uri)
                error_chunks.append(f"=== [{prov.upper()} ERROR LOG] ===\n{content}")

        if args.parse_debug_log:
            debug_links: Dict[str, str] = extracted["debug_links"]
            chosen_debug_uri = (
                args.gcs_debug_uri
                or debug_links.get(target_provider)
                or debug_links.get("ga")
                or debug_links.get("beta")
            )
            if chosen_debug_uri:
                parsed_logs_dir = download_and_parse_debug_log(
                    test_name=test_name,
                    gcs_debug_uri=chosen_debug_uri,
                    output_dir=args.output_dir,
                    base_dir=base_dir,
                )

        if not error_chunks and not parsed_logs_dir:
            raise ValueError(
                "No valid gs://nightly-test-data/ error or debug log links found in GitHub issue."
            )
    else:
        if not args.test_name:
            raise ValueError("--test-name is required when --issue-url is not set.")
        test_name = validate_test_name(args.test_name)
        target_provider = validate_provider(args.target_provider or "ga")

        for uri in args.gcs_error_uri:
            error_chunks.append(fetch_gcs_content(uri))

        if args.local_log:
            validated_local = validate_local_path(args.local_log, base_dir=base_dir)
            with open(validated_local, "r", encoding="utf-8", errors="replace") as f:
                error_chunks.append(f.read())

        if args.parse_debug_log and (args.gcs_debug_uri or args.local_log):
            parsed_logs_dir = download_and_parse_debug_log(
                test_name=test_name,
                gcs_debug_uri=args.gcs_debug_uri,
                local_debug_log=args.local_log,
                output_dir=args.output_dir,
                base_dir=base_dir,
            )

    combined_error = "\n\n".join(error_chunks).strip()
    error_log_file = write_isolated_error_log(
        test_name=test_name,
        error_text=combined_error,
        output_dir=args.output_dir,
    )

    return {
        "normalized_failure_payload": {
            "test_name": test_name,
            "target_provider": target_provider,
            "error_log_file": error_log_file,
            "parsed_logs_dir": parsed_logs_dir or f"{args.output_dir}/{test_name}/",
        }
    }


def main() -> None:
    try:
        payload = run_cli()
        print(json.dumps(payload, indent=2))
    except Exception as exc:
        print(f"[ERROR] {exc}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
