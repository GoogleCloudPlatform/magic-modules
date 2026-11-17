---
name: code-review
description: "Review a Magic Modules pull request, branch, or diff for API contract alignment, state lifecycle safety, backward compatibility, and test adequacy."
---

# Reviewing a Magic Modules Change

Authoritative review criteria live in [`docs/content/code-review/review-pr.md`](../../../../docs/content/code-review/review-pr.md) and [`.agents/knowledge/index.md`](../../../knowledge/index.md). Read `review-pr.md` and `index.md` from your primary `magic-modules` workspace first, then open every knowledge entry and reference doc matching the touched surface (plus [`schema-review`](../schema-review/SKILL.md) when schemas change). This skill defines only how to gather evidence and structure the report.

## 1. Gather Context & Evidence Without Mutating Your Checkout

- Read the **full surrounding files** touched by the change (`product.yaml`, resource YAMLs, custom `.go.tmpl` / Go hooks, and tests/samples), not diff hunks alone. Do not post external GitHub comments unless explicitly asked.
- When schemas, fields, or endpoints are added, modified, or promoted, run [`verify-api-schema`](../verify-api-schema/SKILL.md) against the live GCP REST Discovery document (`v1` for GA, `beta` for beta) to verify `api_name`, types, and GA endpoint presence.
- When CI checks have not already covered the change locally, check out the PR commit in a temporary `git worktree` outside your working tree before running [`run-pre-gen-checks`](../run-pre-gen-checks/SKILL.md) or [`validate-provider-changes`](../validate-provider-changes/SKILL.md) (both inspect `HEAD` rather than a `.diff` file), and remove the temporary worktree and ref when finished. Use [`parse-debug-logs`](../parse-debug-logs/SKILL.md) (writing output outside the repo) when `TF_LOG=DEBUG` logs are available.

## 2. Evidence Gate & Report Format

- Review the entire change across all touched resources, fields, hooks, and test steps, reporting each distinct root-cause defect once.
- **Blocking Findings:** Require exact `file:line` citations, concrete proof (code trace, Discovery schema mismatch, compiler/test failure, or wire log), and what breaks for users. Unverified suspicions belong as questions for the author, never blockers.
- **Non-Blocking Observations:** Style nits, missing pasted manual GA test output when test configs exist and CI passes, and out-of-scope follow-ups (such as mirroring fields into a sibling `<Resource>Version.yaml` when the linked issue is scoped strictly to the parent resource). Never demote failing/skipped tests or missing required test coverage to non-blocking.
- End with **Review Summary**, **Blocking Findings** (omit if none), **Non-Blocking Observations**, and **Verdict** (`APPROVE` if zero blocking findings, else `REJECT`).
