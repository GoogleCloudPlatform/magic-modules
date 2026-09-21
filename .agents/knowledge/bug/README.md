# Bug Knowledge Entries

This directory contains guidelines, recurring patterns, and troubleshooting strategies for fixing bugs in the Google Terraform providers.

## Entry Template

When adding a new bug entry, create a Markdown file in this directory with the following structure:

```yaml
---
name: kebab-case-name
description: Succinct description of the bug pattern or triage rule (<=140 chars).
topics: [bug]
task_types: [bug-fix]
source: doc, PR link, or "authored"
status: draft
last_verified: YYYY-MM-DD
---
```

# Heading (Same as Name)

### The Bug Pattern
* Explain the root cause pattern or recurring scenario.
* Provide a concrete code example of the bug.

### Triage & Resolution Rules
* Actionable, rule-based steps for the agent to identify and fix the issue.
* Include "the why" for each rule.
* Provide a code example showing the correct implementation.
