# TGC Resource Main Loop

This document serves as the entry point for TGC (Terraform Google Conversion) development workflows. 

Depending on your specific task, please refer to and strictly follow the appropriate workflow file:

- **For Adding a New Resource**: Follow the workflow defined in [GEMINI_ADD.md](file:///Users/zhenhuali/Documents/workspace/feature-a/GEMINI_ADD.md).
- **For Fixing a Resource or Test Failure**: Follow the workflow defined in [GEMINI_FIX.md](file:///Users/zhenhuali/Documents/workspace/feature-a/GEMINI_FIX.md).

Please ensure you follow the sequence of phases and mandatory skill checks specified in those files.

## Error Reporting Template

If a failure is detected or provided, you **MUST** report it using the following template before proceeding.

```markdown
# Error Report

## Failed Command
`[Paste the exact command that failed here]`

## Detailed Logs
```
[Paste relevant log snippets or error messages here]
```

## Analysis
[Provide a detailed analysis of the failure. Explain what happened, which stage of the workflow failed, and why.]

## Tracing Evidence (Required for Integration Test failures only)
[Cite specific lines or file contents from `Test_export.tf`, `Test_roundtrip.json`, or `Test_roundtrip.tf` to prove at which stage the data was lost.]

## Proposed Solution
[Describe the proposed solution to fix the failure. Include file paths and code changes if known.]

> [!IMPORTANT]
> I have stopped execution and am waiting for your approval to proceed with this solution.
```
