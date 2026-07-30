import datetime
import json
import os
import re
import subprocess
from collections import defaultdict

def get_date_str(days_ago):
    today = datetime.date.today()
    target_date = today - datetime.timedelta(days=days_ago)
    return target_date.strftime("%Y-%m-%d")

def fetch_results(date_str, provider_type):
    uri = f"gs://nightly-test-data/test-metadata/{provider_type}/{date_str}-{provider_type}.json"
    print(f"Fetching {uri}...")
    result = subprocess.run(["gcloud", "storage", "cat", uri], capture_output=True, text=True)
    if result.returncode != 0:
        print(f"Failed to fetch {uri}")
        return None
    try:
        return json.loads(result.stdout)
    except json.JSONDecodeError as e:
        print(f"Failed to parse JSON for {uri}: {e}")
        return None

def is_generic_error(error_msg):
    if not error_msg:
        return True
    error_msg_lower = error_msg.lower()
    if "error 13" in error_msg_lower:
        return True
    if "internal error" in error_msg_lower:
        return True
    if "failed to perform tenant project creation" in error_msg_lower:
        return True
    return False

def is_panic_or_crash(error_msg):
    if not error_msg:
        return False
    msg = error_msg.strip()
    if "panic: " in msg:
        return True
    if msg.startswith("panic: "):
        return True
    msg_lower = msg.lower()
    if "runtime error:" in msg_lower or "sigsegv" in msg_lower:
        return True
    return False

def get_failures(provider_type):
    failure_counts = defaultdict(int)
    latest_failures = {} # Store name -> {error: msg, log: link}
    latest_available_date = None

    # Check past 7 days (including today)
    for i in range(7):
        date_str = get_date_str(i)
        data = fetch_results(date_str, provider_type)
        if data:
            if latest_available_date is None:
                latest_available_date = date_str
                for item in data:
                    if item.get("status") == "FAILURE":
                        error_msg = item.get("error_message", "")
                        log_link = item.get("log_link") or item.get("LogLink", "")
                        service = item.get("service") or item.get("Service") or "unknown"
                        if not is_generic_error(error_msg):
                            latest_failures[item.get("name")] = {
                                "error": error_msg,
                                "log": log_link,
                                "service": service
                            }

            # Count failures across all days
            for item in data:
                if item.get("status") == "FAILURE":
                    failure_counts[item.get("name")] += 1
                    
    return latest_failures, failure_counts, latest_available_date

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

def main():
    beta_failures, beta_counts, beta_date = get_failures("beta")
    ga_failures, ga_counts, ga_date = get_failures("ga")

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

    report_date = beta_date or ga_date or get_date_str(0)
    output_file = f"tmp/test-status/test-report-{report_date}.md"
    os.makedirs(os.path.dirname(output_file), exist_ok=True)
    
    # 1. Group latest run failures across ALL tests in latest run
    latest_error_groups = defaultdict(lambda: {"tests": set(), "issues": set(), "sample_error": "", "providers": set(), "services": set()})
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

    # 2. Group persistent failures (past 7 days)
    error_groups = defaultdict(lambda: {"tests": set(), "issues": set(), "sample_error": "", "providers": set(), "services": set()})
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

    with open(output_file, "w") as f:
        f.write("# Nightly Test Failures Triage & Monitoring Report\n\n")
        
        if beta_date:
            f.write(f"**Latest Beta run**: {beta_date} ({len(beta_failures)} failing tests)\n")
        if ga_date:
            f.write(f"**Latest GA run**: {ga_date} ({len(ga_failures)} failing tests)\n")
        f.write("\n---\n\n")
        
        # Section 1: High-Impact Errors in Latest Run (High Volume >= 3 tests OR Critical Panic/Crash)
        f.write("## 1. High-Impact Errors in Latest Run\n\n")
        f.write("High-impact errors are flagged based on **Critical Severity** (provider panic/crash) or **High Volume** (affecting $\\ge 3$ tests).\n\n")
        f.write("| # | Impact Category | Affected Tests | Provider | GCP Service(s) | GitHub Issue(s) | Error Signature / Sample Message | Sample Affected Tests |\n")
        f.write("| --- | --- | --- | --- | --- | --- | --- | --- |\n")
        
        def high_impact_sort_key(item):
            sanitized_err, grp = item
            is_panic = is_panic_or_crash(grp["sample_error"])
            num_tests = len(grp["tests"])
            return (1 if is_panic else 0, num_tests)

        sorted_latest_groups = sorted(latest_error_groups.items(), key=high_impact_sort_key, reverse=True)
        high_impact_idx = 1
        for sanitized_err, grp_data in sorted_latest_groups:
            num_tests = len(grp_data["tests"])
            is_panic = is_panic_or_crash(grp_data["sample_error"])
            
            # Highlight if it is a panic/crash OR affects >= 3 tests
            if not is_panic and num_tests < 3:
                continue
                
            impact_badge = "🚨 **CRITICAL (Panic/Crash)**" if is_panic else "⚠️ **High Volume**"
            
            test_list = ", ".join(sorted(list(grp_data["tests"])))
            if len(test_list) > 150:
                test_list = test_list[:147] + "..."
            
            prov_str = "Both" if len(grp_data["providers"]) > 1 else list(grp_data["providers"])[0]
            issues_str = ", ".join(sorted(list(grp_data["issues"]))) if grp_data["issues"] else "N/A"
            services_str = ", ".join(sorted(list(grp_data["services"]))) if grp_data["services"] else "N/A"
            
            err_summary = grp_data["sample_error"][:250].replace("|", "\\|").replace("\n", "<br>")
            err_cell = f"<pre>{err_summary}</pre>"
            
            f.write(f"| {high_impact_idx} | {impact_badge} | **{num_tests}** | {prov_str} | {services_str} | {issues_str} | {err_cell} | `{test_list}` |\n")
            high_impact_idx += 1

        f.write("\n---\n\n")
        
        # Section 2: Persistent Failures Grouped by Error Signature (Past 7 Days)
        f.write("## 2. Persistent Failures Grouped by Error Signature (Past 7 Days)\n\n")
        f.write("Criteria: Failed in latest run and at least 4 days in past 7 days, excluding generic errors.\n\n")
        f.write("| # | Affected Tests Count | GCP Service(s) | Failure Category / Error Signature | GitHub Issue(s) | Affected Test Names |\n")
        f.write("| --- | --- | --- | --- | --- | --- |\n")
        
        sorted_groups = sorted(error_groups.items(), key=lambda x: len(x[1]["tests"]), reverse=True)
        grp_idx = 1
        for sanitized_err, grp_data in sorted_groups:
            num_tests = len(grp_data["tests"])
            test_list = ", ".join(sorted(list(grp_data["tests"])))
            if len(test_list) > 150:
                test_list = test_list[:147] + "..."
            
            issues_str = ", ".join(sorted(list(grp_data["issues"]))) if grp_data["issues"] else "N/A"
            services_str = ", ".join(sorted(list(grp_data["services"]))) if grp_data["services"] else "N/A"
            
            err_summary = grp_data["sample_error"][:250].replace("|", "\\|").replace("\n", "<br>")
            err_cell = f"<pre>{err_summary}</pre>"
            
            f.write(f"| {grp_idx} | **{num_tests}** | {services_str} | {err_cell} | {issues_str} | `{test_list}` |\n")
            grp_idx += 1
            
        f.write("\n---\n\n")
        
        # Section 3: Detailed Per-Test Failure Table
        f.write("## 3. Detailed Per-Test Persistent Failures Table\n\n")
        f.write("| # | Test Name | GCP Service | Provider | Failures (Days) | GitHub Issue | Log Link | Error Message |\n")
        f.write("| --- | --- | --- | --- | --- | --- | --- | --- |\n")
        
        row_idx = 1
        for name in sorted(all_failures.keys()):
            providers = all_failures[name]
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
                    service_str = ga_details.get("service", "N/A")
                    
                    f.write(f"| {row_idx} | {name} | {service_str} | Both (GA shown) | {ga_details['count']} | {issue_link} | {log_display} | {table_error} |\n")
                    row_idx += 1
                    continue
            
            for prov in ["Beta", "GA"]:
                if prov in providers:
                    details = providers[prov]
                    actual_error = sanitize_for_comparison(get_actual_error(details["error"]))
                    
                    table_error = actual_error.replace("|", "\\|").replace("\n", "<br>")
                    table_error = f"<pre>{table_error}</pre>"
                    
                    log_display = f"[Log]({details['log']})" if details["log"] else "N/A"
                    service_str = details.get("service", "N/A")
                    
                    f.write(f"| {row_idx} | {name} | {service_str} | {prov} | {details['count']} | {issue_link} | {log_display} | {table_error} |\n")
                    row_idx += 1

    print(f"\nResults written to {output_file}")

if __name__ == "__main__":
    main()
