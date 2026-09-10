# Agent knowledge base

The **funnel to the knowledge an agent needs**: [`index.md`](index.md) routes to the right source for a
given decision. Agents read the index at decision points and open only what the task needs — never the
whole base (loading everything degrades reasoning and wastes context).

Most knowledge already lives in the **contributor documentation** (`docs/content/`), and that stays the
single source of truth for everything it covers — the index points there. **Agent-only entries** exist in
this directory only for knowledge that has no home in the docs: judgment rules that were never written
down, pitfall catalogs, and (later) lessons proposed by agents from completed tasks. If something is
covered by the contributor docs but covered badly, the fix is improving the docs — not writing a copy here.

Planned entries: [`BACKLOG.md`](BACKLOG.md).

## Entry format

One Markdown file, one topic, ≤120 lines, YAML frontmatter:

```yaml
---
name: enums-vs-strings              # kebab-case, unique, matches filename
description: Model an API enum as Enum (strict) or String (forward-compatible).  # the index line; <=140 chars
topics: [field]                     # index section(s) that list it
task_types: [field-add, new-resource]
source: docs/content/best-practices/validation.md   # provenance: doc, PR, or "authored"
status: certified                   # certified (human-reviewed) | draft (proposed, unreviewed)
last_verified: 2026-07-09
---
```

Body rules: **rule + the why**, never bare imperatives; a real example per rule; a "do NOT use for" line
where misuse is likely; reference the contributor docs rather than duplicating them.

## Curation rules

- **Nothing enters unreviewed.** Entries land via PR like code, **one entry at a time** — each one codifies
  team judgment and deserves a real review conversation. Agent-proposed entries ship as `status: draft` in
  the task PR that motivated them; a human review flips them to `certified`.
- **Never bulk-rewrite.** Edit entry by entry. Do not have a model re-summarize or reorganize the base
  wholesale — that is the documented failure mode (drift and self-contradiction).
- **Contradictions are quarantined.** A proposed entry that contradicts an existing one is not merged; a
  human adjudicates: update or retire the old entry, or discard the proposal.
- **Keep the index true.** Index lines come from entry frontmatter (name + description). Adding or changing
  an entry updates the index in the same commit.
