# Knowledge backlog

Tracks known work for growing the agent knowledge base. This document is not intended for agents themselves.

## Not in the docs

- **casing-and-pluralization** — how API names map to Terraform names; when `api_name` is required.
- **raw-config-access** — when to use `d.GetRawConfig()` / `GetRawPlan()` / `GetRawState()` instead of
  `Get`/`GetOk`, which conflate "unset in config" with "set to the zero value" (false/0/""); the raw cty
  values distinguish null from zero for detecting whether a user actually set a field.

## Overlaps the docs (entry would add the judgment layer)

- **permadiff-decision-path** — choosing between `output`, `default_from_api`, diff suppression, or a real
  fix; mechanics are in `docs/content/develop/diffs.md`.
- **data-source-idioms** — pitfalls beyond the procedure in `docs/content/develop/add-handwritten-datasource.md`.
- **test-adequacy-traps** — cases `docs/content/test/test.md` doesn't cover (identical-config update steps,
  missing import-and-recheck).

## Mostly covered by the docs (revisit only if agents misread them)

- **immutability-nuances** — `docs/content/reference/field.md`, `docs/content/best-practices/immutable-fields.md`.
- **sensitive-and-write-only** — `docs/content/reference/field.md`.

## Lives elsewhere

- **failure-troubleshooting-catalog** — `.agents/skills/operations/troubleshooting_reference.md`; migrates
  here if it outgrows the skill.

## Other future work

- **PR mining** — mine historical PRs to discover new entries. These could be recurring review catches and per-service quirks that were only ever captured in the PR reviews themselves.
