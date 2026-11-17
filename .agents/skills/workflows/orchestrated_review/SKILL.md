---
name: orchestrated-review-workflow
description: "Opt-in two-wave multi-agent code review workflow for large Magic Modules changes. Invoke only when the user explicitly requests an orchestrated or multi-agent review; use code-review for standard reviews."
---

# Orchestrated Code Review Workflow (Opt-In)

Use [`.agents/skills/utils/code-review/SKILL.md`](../../utils/code-review/SKILL.md) directly for standard reviews. When explicitly asked for an orchestrated review on a large multi-resource, Beta-to-GA, or custom-Go change:

1. **Shared Pre-Check:** Run Step 1 of [`code-review`](../../utils/code-review/SKILL.md) once in `/tmp/` (PR diff/metadata, [`verify-api-schema`](../../utils/verify-api-schema/SKILL.md) Discovery schemas, and any needed pre-gen, provider-validation, or parsed debug-log outputs) so subagents share one verified baseline instead of repeating checks.
2. **Wave 1 — Parallel Domain Slices + Holistic Pass:** Perform your own end-to-end pass while running three read-only [`code-reviewer`](../../../agents/code-reviewer.md) subagents (or sequential passes if subagents are unavailable), each grounded in `review-pr.md` and `.agents/knowledge/index.md`:
   - **Schema & API contract:** Live Discovery alignment (`api_name`, types, GA/beta endpoints), `min_version` gating, enum conventions, TGC flags, and sibling schemas.
   - **State lifecycle & custom Go:** Immutability/`ForceNew`, permadiffs, zero-value clearing, import-safe flatteners, custom hooks, and backward compatibility.
   - **Test adequacy & runtime proof:** Configurable field coverage, meaningful 2-step updates, `ImportStateVerify`, test skips, `CheckDestroy`, and unit-test scope.
3. **Wave 2 — Verify & Synthesize:** Read every repository file cited in subagent findings or cross-file leads directly before accepting it, deduplicate findings sharing one root cause, clean up any temporary worktree, and output the standard `code-review` report (**Review Summary**, **Blocking Findings**, **Non-Blocking Observations**, **Verdict**).
