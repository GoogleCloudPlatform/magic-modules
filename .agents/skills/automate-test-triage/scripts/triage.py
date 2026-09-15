import argparse
import datetime
import json
import os
import re
import subprocess
import sys
from collections import defaultdict

def get_date_str(days_ago, base_date=None):
    if base_date is None:
        today = datetime.date.today()
    else:
        today = base_date
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

def fetch_results(date_str, provider_type):
    uri = f"gs://nightly-test-data/test-metadata/{provider_type}/{date_str}-{provider_type}.json"
    print(f"Fetching {uri}...")
    result = subprocess.run(["gcloud", "storage", "cat", uri], capture_output=True, text=True)
    if result.returncode != 0:
        if is_permission_error(result.stderr, result.stdout):
            error_details = (result.stderr or result.stdout or "").strip()
            print(f"\n[ERROR] Permission or authentication failure while fetching {uri}:", file=sys.stderr)
            if error_details:
                print(f"{error_details}\n", file=sys.stderr)
            print("Please check your gcloud authentication ('gcloud auth login') and ensure you have read access ('roles/storage.objectViewer') to gs://nightly-test-data.", file=sys.stderr)
            sys.exit(1)
        print(f"Failed to fetch {uri}")
        return None
    try:
        return json.loads(result.stdout)
    except json.JSONDecodeError as e:
        print(f"Failed to parse JSON for {uri}: {e}")
        return None

def has_actionable_error_context(msg_lower, raw_msg):
    """
    Returns True if an error message contains specific, actionable Terraform/provider/test error details
    that should prevent it from being dismissed as a generic non-actionable Internal Error.
    """
    actionable_patterns = [
        r'\b(?:panic:|runtime\s+error:|sigsegv|nil\s+pointer\s+dereference)\b',
        r'\b(?:plan\s+was\s+not\s+empty|inconsistent\s+(?:result|final\s+plan)|root\s+object\s+was\s+present,\s+but\s+now\s+absent)\b',
        r'\b(?:importstateverify|cannot\s+import\s+non-existent|check\s+failed|expected\s+to\s+be\s+set|expected\s+state)\b',
        r'\b(?:conflicting\s+configuration\s+arguments|invalid\s+resource\s+type|blocks\s+of\s+type|inconsistent\s+dependency\s+lock\s+file)\b',
        r'\b(?:error\s+40[049]|invalid_argument|failed_precondition|permission_denied|missing\s+required\s+argument|unsupported\s+argument)\b',
    ]
    return any(re.search(pat, msg_lower) for pat in actionable_patterns)

def is_quota_or_stockout_error(msg_lower, raw_msg):
    # 1. Structured GCP error codes and protobuf types
    structured_markers = [
        "google.rpc.quotafailure",
        "rate_limit_exceeded",
        "zone_resource_pool_exhausted",
        "resource_exhausted",
        "stockout",
        "error 429",
        "429 too many requests",
        "the folder operation violates fanout constraints",
    ]
    if any(k in msg_lower for k in structured_markers):
        return True

    # 2. Semantic regex / concept matching (grammar-agnostic)
    # Quota & Rate-Limit concept (e.g. "quota exceeded", "quotas are exceeded", "rate limit exceeded", "you do not have quota")
    if re.search(r'\b(quota|rate\s*limit)s?\b.*\b(exceed|exhaust|limit|violate)e?d?\b', msg_lower):
        return True
    if "you do not have quota" in msg_lower or "quota limit" in msg_lower or "has been exceeded" in msg_lower:
        return True

    # Stockout & Capacity concept (e.g. "no resource available", "not enough resources available", "insufficient capacity")
    if re.search(r'\b(insufficient|not\s+enough|no|lack\s+of)\b.*\b(capacity|resource|stock)s?\b', msg_lower):
        return True

    # Location / Zone / Region retry advice concept (e.g. "try a different zone, or try again later", "try again in a different zone")
    if re.search(r'\btry\s+(again\s+in\s+|a\s+different\s+)?(later|zone|region|location)\b', msg_lower) and ("resource" in msg_lower or "capacity" in msg_lower or "available" in msg_lower):
        return True

    return False

def is_internal_error_13(msg_lower, raw_msg):
    if has_actionable_error_context(msg_lower, raw_msg):
        return False

    structured_markers = [
        "grpc.status\": 13",
        "error code 13",
        "error 13",
        "code: 'internal'",
        "error 500",
        "error 503",
        "error 502",
        "backenderror",
    ]
    if any(k in msg_lower for k in structured_markers):
        return True

    if re.search(r'\b(?:internal\s+(?:server\s+)?error|an?\s+internal\s+error\s+has\s+occurred|internal\s+error\s+during\s+operation)\b', msg_lower):
        return True

    return False

def is_tenant_project_creation_error(msg_lower, raw_msg):
    return bool(re.search(r'\b(?:fail(?:ed|ure)?|error|unable|could\s+not)\b.*\btenant\s+project\b|\btenant\s+project\b.*\b(?:fail(?:ed|ure)?|error|unable|creation)\b', msg_lower))

TEST_ENV_PROJECTS = [
    "ci-test-project",
    "ci-test-project-188019", "1067888929963",
    "ci-test-project-nightly-ga", "594424405950",
    "ci-test-project-nightly-beta", "653407317329",
    "tf-vcr-private", "808590572184"
]

def is_project_allowlist_or_permission_error(msg_lower, raw_msg):
    # Do NOT classify as a non-actionable allowlist error if it is simply an API enablement error in a test environment project!
    # (An agent can automatically run 'gcloud services enable <api>' to fix those!)
    is_api_disabled = bool(re.search(r'\b(?:has\s+not\s+been\s+used\s+in\s+project|before\s+or\s+it\s+is\s+disabled|reason:\s*"?service_disabled"?|api_not_enabled|enable\s+it\s+by\s+visiting)\b', msg_lower))
    if is_api_disabled and any(p in msg_lower for p in TEST_ENV_PROJECTS):
        return False

    structured_markers = [
        "reason: \"project_not_allowlisted\"",
        "reason: \"consumer_invalid\"",
    ]
    if any(k in msg_lower for k in structured_markers):
        return True

    if re.search(r'\b(?:not\s+allowlisted|not\s+in\s+allowlist|require(?:s|d)?\s+allowlist(?:ing)?|unallowlisted|allowlisted\s+for)\b', msg_lower):
        return True

    if re.search(r'\b(?:not\s+allowed|prohibited|unauthorized|access\s+denied)\b.*\bproject\b|\bproject\b.*\b(?:not\s+allowed|not\s+allowlisted|allowlist|unauthorized)\b', msg_lower):
        if any(w in msg_lower for w in ["api", "engine", "service", "allowlist", "terraform"]):
            return True

    return False

HUMAN_ACTION_RULES = [
    ("Quota / Resource Availability", is_quota_or_stockout_error),
    ("Internal Error (Error Code 13)", is_internal_error_13),
    ("Tenant Project Creation Failure", is_tenant_project_creation_error),
    ("Project Allowlist / API Permission Required", is_project_allowlist_or_permission_error),
]

def classify_human_action(error_msg):
    """
    Returns category name if the error requires human action (non-actionable by agents), else None.
    """
    if not error_msg:
        return "Empty / Generic Error"
    raw_msg = error_msg.strip()
    msg_lower = raw_msg.lower()
    for category_name, matcher in HUMAN_ACTION_RULES:
        if matcher(msg_lower, raw_msg):
            return category_name
    return None

def is_generic_error(error_msg):
    return classify_human_action(error_msg) is not None

SEVERITY_RULES = [
    # (Priority, Category ID, Display Badge, Matcher Function)
    (100, "PANIC", "🚨 **CRITICAL (Panic/Crash)**", lambda msg_lower, raw_msg: bool(re.search(r'\b(?:panic:|runtime\s+error:|sigsegv|nil\s+pointer\s+dereference)\b', raw_msg, re.IGNORECASE))),
    (90,  "API_ENV", "🚨 **CRITICAL (API Not Enabled in Test Env)**", lambda msg_lower, raw_msg: (
        bool(re.search(r'\b(?:has\s+not\s+been\s+used\s+in\s+project|before\s+or\s+it\s+is\s+disabled|reason:\s*"?service_disabled"?|api_not_enabled|enable\s+it\s+by\s+visiting)\b', msg_lower))
        and any(p in msg_lower for p in TEST_ENV_PROJECTS)
    )),
]

def classify_severity(error_msg):
    """
    Classifies the severity of an actionable error message based on SEVERITY_RULES.
    Returns (priority_score, badge_string).
    If no critical severity rule matches, returns (0, None).
    """
    if not error_msg:
        return (0, None)
    raw_msg = error_msg.strip()
    msg_lower = raw_msg.lower()
    for priority, cat_id, badge, matcher in SEVERITY_RULES:
        if matcher(msg_lower, raw_msg):
            return (priority, badge)
    return (0, None)

def is_panic_or_crash(error_msg):
    priority, _ = classify_severity(error_msg)
    return priority == 100

def get_failures(provider_type, base_date=None):
    failure_counts = defaultdict(int)
    latest_failures = {} # Store name -> {error: msg, log: link}
    latest_available_date = None
    latest_sha = ""
    service_stats = defaultdict(lambda: {"total": 0, "failed": 0})

    # Check past 7 days (including today)
    for i in range(7):
        date_str = get_date_str(i, base_date=base_date)
        data = fetch_results(date_str, provider_type)
        if data:
            if latest_available_date is None:
                latest_available_date = date_str
                for item in data:
                    if not latest_sha and item.get("commit_sha"):
                        latest_sha = item.get("commit_sha")
                    service = item.get("service") or item.get("Service") or "unknown"
                    service_stats[service]["total"] += 1
                    if item.get("status") == "FAILURE":
                        service_stats[service]["failed"] += 1
                        error_msg = item.get("error_message", "")
                        log_link = item.get("log_link") or item.get("LogLink", "")
                        latest_failures[item.get("name")] = {
                            "error": error_msg,
                            "log": log_link,
                            "service": service
                        }

            # Count failures across all days
            for item in data:
                if item.get("status") == "FAILURE":
                    failure_counts[item.get("name")] += 1
                    
    return latest_failures, failure_counts, latest_available_date, service_stats, latest_sha

def get_actual_error(error_str):
    lines = [line.strip() for line in error_str.split('\n') if not (line.startswith("=== RUN") or line.startswith("=== PAUSE") or line.startswith("=== CONT") or line.startswith("--- FAIL:") or line.strip() == "FAIL")]
    clean_error = "\n".join(lines).strip()
    
    # Check for panic first
    for line in lines:
        if line.startswith("panic: ") or "panic: " in line:
            return line
            
    error_keywords = ["runtime error:", "Error:", "googleapi: Error", "Check failed:"]
    for kw in error_keywords:
        idx = clean_error.find(kw)
        if idx != -1:
            err_block = clean_error[idx:]
            return re.split(r'\n+(?:Error|Check failed|panic):', err_block)[0].strip()
            
    return re.split(r'\n+(?:Error|Check failed|panic):', clean_error)[0].strip()

def sanitize_for_comparison(error_str):
    # Replace project IDs and resource names like tf-test... or tf_test...
    s = re.sub(r'tf[-_]test[a-z0-9-_]+', 'tf-test-ID', error_str)
    # Replace ci-test-project-nightly-beta/ga
    s = re.sub(r'ci-test-project-nightly-[a-z]+', 'ci-test-project', s)
    # Replace project numbers like projects/12345678
    s = re.sub(r'projects/\d+', 'projects/NUMBER', s)
    # Replace project numbers in messages like project number: 123456
    s = re.sub(r'project number: \d+', 'project number: NUMBER', s)
    # Replace folder and organization numbers like folders/123456
    s = re.sub(r'(folders|organizations|reasoningEngines)/\d+', r'\g<1>/NUMBER', s)
    # Replace random numbers in subject or violations
    s = re.sub(r'project:\d+', 'project:NUMBER', s)
    # Replace service account project numbers like service-123456@
    s = re.sub(r'service-\d+@', 'service-NUMBER@', s)
    # Replace custom service account names before @ci-test-project...
    s = re.sub(r'[a-z0-9-]+@ci-test-project', 'sa@ci-test-project', s)
    # Replace service account key IDs
    s = re.sub(r'/keys/[a-f0-9]+', '/keys/KEY_ID', s)
    # Replace random hex strings ending in -tp (common in Apigee tests)
    s = re.sub(r'[a-z0-9]{10,20}-tp', 'RANDOM-tp', s)
    # Replace UUIDs
    s = re.sub(r'[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}', 'UUID', s)
    # Normalize v1beta vs v1 in Google API URLs
    s = re.sub(r'googleapis\.com/v1beta/', 'googleapis.com/v1/', s)
    # Replace resource array indices like [0]
    s = re.sub(r'\[\d+\]', '[X]', s)
    # Replace Help Tokens
    s = re.sub(r'Help Token: [a-zA-Z0-9_-]+', 'Help Token: TOKEN', s)
    # Replace Request IDs
    s = re.sub(r'"requestId":\s*"[a-f0-9]+"', '"requestId": "ID"', s)
    # Replace tag values
    s = re.sub(r'tagValues/\d+', 'tagValues/NUMBER', s)
    # Replace timestamps in JSON metadata like "time":"2026-05-13T04:20:17.485043Z"
    s = re.sub(r'"time":\s*"\d{4}-\d{2}-\d{2}T[0-9:.]+Z"', '"time": "TIMESTAMP"', s)
    # Replace Go test file line numbers
    s = re.sub(r'[a-zA-Z0-9_]+_test\.go:\d+:', 'test.go:LINE:', s)
    # Replace terraform plan unchanged hidden counts
    s = re.sub(r'\(\d+ unchanged (attributes|blocks|elements) hidden\)', r'(X unchanged \g<1> hidden)', s)
    # Replace google_<resource_type>.<resource_name> with google_<resource_type>.RESOURCE_NAME
    s = re.sub(r'\b(google_[a-zA-Z0-9_-]+)\.([a-zA-Z0-9_-]+)\b', r'\1.RESOURCE_NAME', s)
    # Replace resource "google_<resource_type>" "<resource_name>" with resource "google_<resource_type>" "RESOURCE_NAME"
    s = re.sub(r'\bresource\s+"(google_[a-zA-Z0-9_-]+)"\s+"([a-zA-Z0-9_-]+)"', r'resource "\1" "RESOURCE_NAME"', s)
    
    # 1. Truncate post-test destroy failures/warnings
    s = re.sub(r'(?s)(?:\b\w+\.go:\d+:\s*)?Error running post-test destroy.*$', '', s)
    # 2. Sanitize Terraform configuration line numbers (e.g. on terraform_plugin_test.tf line 12)
    s = re.sub(r'on\s+([a-zA-Z0-9_-]+\.tf)\s+line\s+\d+', r'on \1 line LINE', s)
    # 3. Sanitize line prefix numbers (e.g. 12: resource "google_...")
    s = re.sub(r'(?m)^\s*\d+:\s*', 'LINE: ', s)
    
    # 4. Sanitize random suffixes (length 8-10 or 16)
    s = re.sub(r'([/_-])[a-z0-9]{8,10}\b', r'\1ID', s)
    s = re.sub(r'([/_-])[a-z0-9]{16}\b', r'\1ID', s)
    
    # 5. Sanitize Terraform attributes containing resource names/paths
    s = re.sub(r'\b(id|name|project|bucket|namespace|table|email|member|role|unique_id)\s*=\s*".*?"', r'\1 = "VALUE"', s)
    
    # 6. Sanitize Step X/Y numbers
    s = re.sub(r'\bStep\s+\d+/\d+\b', 'Step X/Y', s)
    
    return s.strip()


def find_issue_link(test_name, issues):
    exact_target = f"Failing test(s): {test_name}"
    for issue in issues:
        title = issue.get("title", "").strip()
        if title == exact_target:
            return f"[#{issue['number']}]({issue['url']})"
            
    # Try wildcard/substring matching e.g. "Failing test(s): TestAccAccessContextManager*"
    for issue in issues:
        title = issue.get("title", "").strip()
        if title.startswith("Failing test(s): "):
            pattern = title.replace("Failing test(s): ", "").strip()
            if pattern.endswith("*"):
                prefix = pattern[:-1]
                if test_name.startswith(prefix):
                    return f"[#{issue['number']}]({issue['url']})"
            elif test_name in title or pattern in test_name:
                return f"[#{issue['number']}]({issue['url']})"
                
    return "N/A"

def format_service_summary(srv, beta_stats, ga_stats, fail_count):
    b_tot = beta_stats.get(srv, {}).get("total", 0)
    b_fail = beta_stats.get(srv, {}).get("failed", 0)
    g_tot = ga_stats.get(srv, {}).get("total", 0)
    g_fail = ga_stats.get(srv, {}).get("failed", 0)
    
    parts = []
    if g_tot > 0:
        g_pct = (g_fail / g_tot) * 100.0
        parts.append(f"GA: <b>{g_fail} / {g_tot}</b> failed ({g_pct:.1f}%)")
    if b_tot > 0:
        b_pct = (b_fail / b_tot) * 100.0
        parts.append(f"Beta: <b>{b_fail} / {b_tot}</b> failed ({b_pct:.1f}%)")
        
    stats_str = " | ".join(parts) if parts else "N/A"
    return f"📦 <code>{srv}</code> — {stats_str}"

def format_human_action_cell(error_str):
    cat = classify_human_action(error_str)
    if cat is not None:
        return f"👤 **Yes ({cat})**"
    return "🤖 **No (Agent-Actionable)**"

def main():
    parser = argparse.ArgumentParser(description="Automate test triage and reporting.")
    parser.add_argument("--date", help="Target end date (YYYY-MM-DD) for 7-day window. Defaults to today.", default=None)
    args = parser.parse_args()

    base_date = None
    if args.date:
        base_date = datetime.datetime.strptime(args.date, "%Y-%m-%d").date()

    beta_failures, beta_counts, beta_date, beta_stats, beta_sha = get_failures("beta", base_date=base_date)
    ga_failures, ga_counts, ga_date, ga_stats, ga_sha = get_failures("ga", base_date=base_date)

    # Combine failures for 7-day persistent window
    all_failures = defaultdict(dict)
    count_threshold = 4
    
    for name, details in beta_failures.items():
        count = beta_counts[name]
        if count >= count_threshold:
            if name not in all_failures:
                all_failures[name] = {}
            all_failures[name]["Beta"] = {
                "count": count,
                "error": details["error"],
                "log": details["log"],
                "service": details["service"]
            }
            
    for name, details in ga_failures.items():
        count = ga_counts[name]
        if count >= count_threshold:
            if name not in all_failures:
                all_failures[name] = {}
            all_failures[name]["GA"] = {
                "count": count,
                "error": details["error"],
                "log": details["log"],
                "service": details["service"]
            }

    # Fetch GitHub issues
    issues = []
    print("Fetching open test-failure issues from GitHub...")
    gh_res = subprocess.run(["gh", "issue", "list", "--repo", "hashicorp/terraform-provider-google", "--label", "test-failure", "--state", "open", "--limit", "1000", "--json", "number,title,url"], capture_output=True, text=True)
    if gh_res.returncode == 0:
        try:
            issues = json.loads(gh_res.stdout)
        except Exception as e:
            print(f"Failed to parse gh issue list JSON: {e}")
    else:
        print(f"Failed to fetch GitHub issues: {gh_res.stderr}")

    report_date = beta_date or ga_date or get_date_str(0, base_date=base_date)
    output_file = f"tmp/test-status/test-report-{report_date}.md"
    os.makedirs(os.path.dirname(output_file), exist_ok=True)
    
    # 1. Group latest run failures across ALL tests in latest run
    latest_error_groups = defaultdict(lambda: {"tests": set(), "issues": set(), "sample_error": "", "sample_log": "", "providers": set(), "services": set()})
    for name, details in beta_failures.items():
        actual_err = get_actual_error(details["error"])
        sanitized_err = sanitize_for_comparison(actual_err)
        grp = latest_error_groups[sanitized_err]
        grp["tests"].add(name)
        grp["providers"].add("Beta")
        grp["services"].add(details.get("service", "unknown"))
        issue = find_issue_link(name, issues)
        if issue != "N/A":
            grp["issues"].add(issue)
        if not grp["sample_error"]:
            grp["sample_error"] = sanitized_err
        if not grp["sample_log"] and details.get("log"):
            grp["sample_log"] = details.get("log")

    for name, details in ga_failures.items():
        actual_err = get_actual_error(details["error"])
        sanitized_err = sanitize_for_comparison(actual_err)
        grp = latest_error_groups[sanitized_err]
        grp["tests"].add(name)
        grp["providers"].add("GA")
        grp["services"].add(details.get("service", "unknown"))
        issue = find_issue_link(name, issues)
        if issue != "N/A":
            grp["issues"].add(issue)
        if not grp["sample_error"]:
            grp["sample_error"] = sanitized_err
        if not grp["sample_log"] and details.get("log"):
            grp["sample_log"] = details.get("log")

    # 2. Group persistent failures (past 7 days)
    error_groups = defaultdict(lambda: {"tests": set(), "issues": set(), "sample_error": "", "sample_log": "", "providers": set(), "services": set()})
    for name, providers in all_failures.items():
        issue_link = find_issue_link(name, issues)
        for prov, details in providers.items():
            actual_err = get_actual_error(details["error"])
            sanitized_err = sanitize_for_comparison(actual_err)
            group = error_groups[sanitized_err]
            group["tests"].add(name)
            group["providers"].add(prov)
            group["services"].add(details.get("service", "unknown"))
            if issue_link != "N/A":
                group["issues"].add(issue_link)
            if not group["sample_error"]:
                group["sample_error"] = sanitized_err
            if not group["sample_log"] and details.get("log"):
                group["sample_log"] = details.get("log")

    with open(output_file, "w") as f:
        f.write("# Nightly Test Failures Triage & Monitoring Report\n\n")
        
        beta_total_tests = sum(s.get("total", 0) for s in beta_stats.values())
        ga_total_tests = sum(s.get("total", 0) for s in ga_stats.values())
        
        if beta_date:
            b_pct_str = f" — {(len(beta_failures)/beta_total_tests)*100.0:.1f}% failure rate" if beta_total_tests > 0 else ""
            b_sha_str = f" ([{beta_sha}](https://github.com/hashicorp/terraform-provider-google-beta/commit/{beta_sha}))" if beta_sha else ""
            f.write(f"**Latest Beta run**: {beta_date}{b_sha_str} ({len(beta_failures)} / {beta_total_tests} failing tests{b_pct_str})\n\n")
        if ga_date:
            g_pct_str = f" — {(len(ga_failures)/ga_total_tests)*100.0:.1f}% failure rate" if ga_total_tests > 0 else ""
            g_sha_str = f" ([{ga_sha}](https://github.com/hashicorp/terraform-provider-google/commit/{ga_sha}))" if ga_sha else ""
            f.write(f"**Latest GA run**: {ga_date}{g_sha_str} ({len(ga_failures)} / {ga_total_tests} failing tests{g_pct_str})\n")
        latest_run_date = beta_date or ga_date or ""
        date_str = f" ({latest_run_date})" if latest_run_date else ""
        f.write("\n### 📑 Table of Contents\n")
        f.write("* [1. High-Impact Agent-Actionable Errors in Latest Run](#section-1)\n")
        f.write("* [2. Test Failures Requiring Human Action (Non-Actionable by Agents)](#section-2)\n")
        f.write("* [3. Persistent Agent-Actionable Failures Grouped by Error Signature (Past 7 Days)](#section-3)\n")
        f.write(f"* [4. Detailed Test Failures Grouped by Service Package{date_str}](#section-4)\n")
        f.write("\n---\n\n")
        
        # Section 1: High-Impact Agent-Actionable Errors in Latest Run (Critical Severity OR High Volume >= 3 tests)
        f.write("<details open>\n<summary><h2 id=\"section-1\">1. High-Impact Agent-Actionable Errors in Latest Run</h2></summary>\n<br>\n\n")
        f.write("High-impact agent-actionable errors are flagged based on **Critical Severity** (e.g., provider panic/crash, API enablement error in test environment) or **High Volume** (affecting $\\ge 3$ tests).\n")
        f.write("> 🤖 **What does \"Agent-Actionable\" mean?** An error is **agent-actionable** if it represents a provider code defect, schema regression, or permadiff that an autonomous AI agent can automatically detect, triage, and repair through a sequence of automated code remediation actions—enabling headless automated remediation without requiring human intervention.\n")
        f.write("> ℹ️ **Note:** Error messages in this table are truncated to 500 characters for readability. For complete, untruncated error messages and per-test details, see **[Section 4: Detailed Test Failures Grouped by Service Package](#section-4)**.\n\n")
        f.write("| # | Impact Category | Affected Tests | Provider | GCP Service(s) | GitHub Issue(s) | Log Link | Error Signature / Sample Message | Sample Affected Tests |\n")
        f.write("| --- | --- | --- | --- | --- | --- | --- | --- | --- |\n")
        
        def high_impact_sort_key(item):
            sanitized_err, grp = item
            priority, _ = classify_severity(grp["sample_error"])
            num_tests = len(grp["tests"])
            return (priority, num_tests)

        sorted_latest_groups = sorted(latest_error_groups.items(), key=high_impact_sort_key, reverse=True)
        high_impact_idx = 1
        for sanitized_err, grp_data in sorted_latest_groups:
            if classify_human_action(grp_data["sample_error"]) is not None:
                continue
            num_tests = len(grp_data["tests"])
            priority, severity_badge = classify_severity(grp_data["sample_error"])
            
            # Highlight if it is a Critical Severity error (priority > 0) OR affects >= 3 tests
            if priority == 0 and num_tests < 3:
                continue
                
            badges = []
            if severity_badge:
                badges.append(severity_badge)
            if num_tests >= 3:
                badges.append("⚠️ **High Volume**")
                
            impact_badge = "<br>".join(badges) if badges else "⚠️ **High Volume**"
            
            test_list = ", ".join(sorted(list(grp_data["tests"])))
            if len(test_list) > 150:
                test_list = test_list[:147] + "..."
            
            prov_str = "Both" if len(grp_data["providers"]) > 1 else list(grp_data["providers"])[0]
            issues_str = ", ".join(sorted(list(grp_data["issues"]))) if grp_data["issues"] else "N/A"
            services_str = ", ".join(sorted(list(grp_data["services"]))) if grp_data["services"] else "N/A"
            
            log_display = f"[Log]({grp_data['sample_log']})" if grp_data.get("sample_log") else "N/A"
            err_summary = grp_data["sample_error"][:500].replace("|", "\\|").replace("\n", "<br>")
            err_cell = f"<pre>{err_summary}</pre>"
            
            has_s4 = any(t in all_failures for t in grp_data["tests"])
            s4_link = "<br>[See full details in Section 4](#section-4)" if has_s4 else ""
            test_col = f"`{test_list}`{s4_link}"
            
            f.write(f"| {high_impact_idx} | {impact_badge} | **{num_tests}** | {prov_str} | {services_str} | {issues_str} | {log_display} | {err_cell} | {test_col} |\n")
            high_impact_idx += 1

        f.write("\n</details>\n<br>\n\n---\n\n")
        
        # Section 2: Test Failures Requiring Human Action (Non-Actionable by Agents)
        f.write("<details open>\n<summary><h2 id=\"section-2\">2. Test Failures Requiring Human Action (Non-Actionable by Agents)</h2></summary>\n<br>\n\n")
        f.write("These failures are caused by infrastructure, quota, or test account/environment issues that cannot be fixed by code changes in Terraform providers.\n")
        f.write("> ℹ️ **Note:** Error messages in this table are truncated to 500 characters for readability. For complete, untruncated error messages and per-test details, see **[Section 4: Detailed Test Failures Grouped by Service Package](#section-4)**.\n\n")
        
        human_action_groups = {}
        for grp_dict in [latest_error_groups, error_groups]:
            for sanitized_err, grp_data in grp_dict.items():
                cat = classify_human_action(grp_data["sample_error"])
                if cat is not None:
                    if sanitized_err not in human_action_groups:
                        human_action_groups[sanitized_err] = {
                            "category": cat,
                            "tests": set(grp_data["tests"]),
                            "providers": set(grp_data["providers"]),
                            "services": set(grp_data["services"]),
                            "issues": set(grp_data["issues"]),
                            "sample_error": grp_data["sample_error"],
                            "sample_log": grp_data.get("sample_log", "")
                        }
                    else:
                        human_action_groups[sanitized_err]["tests"].update(grp_data["tests"])
                        human_action_groups[sanitized_err]["providers"].update(grp_data["providers"])
                        human_action_groups[sanitized_err]["services"].update(grp_data["services"])
                        human_action_groups[sanitized_err]["issues"].update(grp_data["issues"])
                        if not human_action_groups[sanitized_err]["sample_log"] and grp_data.get("sample_log"):
                            human_action_groups[sanitized_err]["sample_log"] = grp_data.get("sample_log")
                        
        if not human_action_groups:
            f.write("No test failures requiring human action were detected in the current monitoring window.\n\n")
        else:
            f.write("| # | Human Action Category | Affected Tests | Provider | GCP Service(s) | GitHub Issue(s) | Log Link | Error Signature / Sample Message | Sample Affected Tests |\n")
            f.write("| --- | --- | --- | --- | --- | --- | --- | --- | --- |\n")
            
            sorted_human_groups = sorted(human_action_groups.items(), key=lambda x: len(x[1]["tests"]), reverse=True)
            human_idx = 1
            for sanitized_err, grp_data in sorted_human_groups:
                num_tests = len(grp_data["tests"])
                cat_badge = f"👤 **{grp_data['category']}**"
                test_list = ", ".join(sorted(list(grp_data["tests"])))
                if len(test_list) > 150:
                    test_list = test_list[:147] + "..."
                
                prov_str = "Both" if len(grp_data["providers"]) > 1 else list(grp_data["providers"])[0]
                issues_str = ", ".join(sorted(list(grp_data["issues"]))) if grp_data["issues"] else "N/A"
                services_str = ", ".join(sorted(list(grp_data["services"]))) if grp_data["services"] else "N/A"
                
                log_display = f"[Log]({grp_data['sample_log']})" if grp_data.get("sample_log") else "N/A"
                err_summary = grp_data["sample_error"][:500].replace("|", "\\|").replace("\n", "<br>")
                err_cell = f"<pre>{err_summary}</pre>"
                
                f.write(f"| {human_idx} | {cat_badge} | **{num_tests}** | {prov_str} | {services_str} | {issues_str} | {log_display} | {err_cell} | `{test_list}` |\n")
                human_idx += 1
                
        f.write("\n</details>\n<br>\n\n---\n\n")
        
        # Section 3: Persistent Agent-Actionable Failures Grouped by Error Signature (Past 7 Days)
        f.write("<details open>\n<summary><h2 id=\"section-3\">3. Persistent Agent-Actionable Failures Grouped by Error Signature (Past 7 Days)</h2></summary>\n<br>\n\n")
        f.write("Criteria: Persistent **agent-actionable** failures (failing in latest run and at least 4 days in past 7 days) that can be automatically detected and remediated by autonomous agents without human intervention.\n")
        f.write("> ℹ️ **Note:** Error messages in this table are truncated to 500 characters for readability. For complete, untruncated error messages and per-test details, see **[Section 4: Detailed Test Failures Grouped by Service Package](#section-4)**.\n\n")
        f.write("| # | Affected Tests Count | GCP Service(s) | Failure Category / Error Signature | GitHub Issue(s) | Log Link | Affected Test Names |\n")
        f.write("| --- | --- | --- | --- | --- | --- | --- |\n")
        
        sorted_groups = sorted(error_groups.items(), key=lambda x: len(x[1]["tests"]), reverse=True)
        grp_idx = 1
        for sanitized_err, grp_data in sorted_groups:
            if classify_human_action(grp_data["sample_error"]) is not None:
                continue
            num_tests = len(grp_data["tests"])
            test_list = ", ".join(sorted(list(grp_data["tests"])))
            if len(test_list) > 150:
                test_list = test_list[:147] + "..."
            
            issues_str = ", ".join(sorted(list(grp_data["issues"]))) if grp_data["issues"] else "N/A"
            services_str = ", ".join(sorted(list(grp_data["services"]))) if grp_data["services"] else "N/A"
            
            log_display = f"[Log]({grp_data['sample_log']})" if grp_data.get("sample_log") else "N/A"
            err_summary = grp_data["sample_error"][:500].replace("|", "\\|").replace("\n", "<br>")
            err_cell = f"<pre>{err_summary}</pre>"
            
            has_s4 = any(t in all_failures for t in grp_data["tests"])
            s4_link = "<br>[See full details in Section 4](#section-4)" if has_s4 else ""
            test_col = f"`{test_list}`{s4_link}"
            
            f.write(f"| {grp_idx} | **{num_tests}** | {services_str} | {err_cell} | {issues_str} | {log_display} | {test_col} |\n")
            grp_idx += 1
            
        f.write("\n</details>\n<br>\n\n---\n\n")
        
        # Section 4: Detailed Test Failures Grouped by Service Package
        latest_run_date = beta_date or ga_date or ""
        date_str = f" ({latest_run_date})" if latest_run_date else ""
        f.write(f"<details open>\n<summary><h2 id=\"section-4\">4. Detailed Test Failures Grouped by Service Package{date_str}</h2></summary>\n<br>\n\n")
        f.write("All test failures from the latest run are grouped by GCP Service package. Each collapsible section displays how many tests failed out of total tests run in that package, with a column indicating whether human intervention is required (`👤 **Yes**` for infrastructure/quota issues vs. `🤖 **No (Agent-Actionable)**` for provider code defects fixable by agents).\n")
        f.write("> ℹ️ **Note:** Unlike Sections 1–3, error messages in this section are displayed in **full without truncation** for comprehensive debugging.\n\n")
        
        combined_latest_tests = defaultdict(dict)
        for name, details in beta_failures.items():
            combined_latest_tests[name]["Beta"] = {
                "count": beta_counts[name],
                "error": details["error"],
                "log": details["log"],
                "service": details["service"]
            }
        for name, details in ga_failures.items():
            combined_latest_tests[name]["GA"] = {
                "count": ga_counts[name],
                "error": details["error"],
                "log": details["log"],
                "service": details["service"]
            }

        service_to_tests = defaultdict(list)
        for name, providers in combined_latest_tests.items():
            srv = "unknown"
            if "Beta" in providers:
                srv = providers["Beta"].get("service", "unknown")
            elif "GA" in providers:
                srv = providers["GA"].get("service", "unknown")
            service_to_tests[srv].append(name)
            
        sorted_services = sorted(service_to_tests.keys(), key=lambda s: (s == "unknown", s.lower()))
        
        if not sorted_services:
            f.write("No test failures to report.\n\n")
        else:
            for idx, srv in enumerate(sorted_services):
                srv_summary = format_service_summary(srv, beta_stats, ga_stats, len(service_to_tests[srv]))
                open_attr = " open" if idx == 0 else ""
                f.write(f"<details{open_attr}>\n<summary><b>{srv_summary}</b></summary>\n\n")
                f.write("| # | Test Name | Provider | Failures (7d) | Human Action Required? | GitHub Issue | Log Link | Error Message |\n")
                f.write("| --- | --- | --- | --- | --- | --- | --- | --- |\n")
                
                row_idx = 1
                for name in sorted(service_to_tests[srv]):
                    providers = combined_latest_tests[name]
                    issue_link = find_issue_link(name, issues)
                    
                    if "Beta" in providers and "GA" in providers:
                        beta_details = providers["Beta"]
                        ga_details = providers["GA"]
                        
                        beta_error = sanitize_for_comparison(get_actual_error(beta_details["error"]))
                        ga_error = sanitize_for_comparison(get_actual_error(ga_details["error"]))
                        
                        if beta_error == ga_error:
                            table_error = ga_error.replace("|", "\\|").replace("\n", "<br>")
                            table_error = f"<pre>{table_error}</pre>"
                            log_display = f"[Log]({ga_details['log']})" if ga_details['log'] else "N/A"
                            ha_cell = format_human_action_cell(ga_details["error"])
                            
                            f.write(f"| {row_idx} | {name} | Both (GA shown) | {ga_details['count']} | {ha_cell} | {issue_link} | {log_display} | {table_error} |\n")
                            row_idx += 1
                            continue
                    
                    for prov in ["Beta", "GA"]:
                        if prov in providers:
                            details = providers[prov]
                            actual_error = sanitize_for_comparison(get_actual_error(details["error"]))
                            
                            table_error = actual_error.replace("|", "\\|").replace("\n", "<br>")
                            table_error = f"<pre>{table_error}</pre>"
                            
                            log_display = f"[Log]({details['log']})" if details["log"] else "N/A"
                            ha_cell = format_human_action_cell(details["error"])
                            
                            f.write(f"| {row_idx} | {name} | {prov} | {details['count']} | {ha_cell} | {issue_link} | {log_display} | {table_error} |\n")
                            row_idx += 1
                            
                f.write("\n</details>\n<br>\n\n")

        f.write("\n</details>\n<br>\n")

    print(f"\nResults written to {output_file}")

if __name__ == "__main__":
    main()
