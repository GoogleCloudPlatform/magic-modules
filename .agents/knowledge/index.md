# Knowledge index

Read this index at decision points; open only the source the task needs. Format and curation rules:
[README.md](README.md).

## Contributor documentation (`docs/content/`)

**Reference (what the schema properties mean):**

- Field properties (`output`, `immutable`, defaults, `sensitive`, validation, enums, arrays, `api_name`, conflicts) — `docs/content/reference/field.md`
- Resource-level properties (URLs, timeouts, async, custom code hooks) — `docs/content/reference/resource.md`
- Resource metadata — `docs/content/reference/metadata.md`
- Samples format — `docs/content/reference/sample.md`
- Update test expectations — `docs/content/reference/update-test-changes.md`
- Make commands — `docs/content/reference/make-commands.md`

**Procedures (how to do a task):**

- Add or change a field — `docs/content/develop/add-fields.md`
- Add a resource — `docs/content/develop/add-resource.md`
- Add a handwritten data source — `docs/content/develop/add-handwritten-datasource.md`
- Add IAM support — `docs/content/develop/add-iam-support.md`
- Custom code (expanders, flatteners, hooks) — `docs/content/develop/custom-code.md`
- Fix diffs and permadiffs — `docs/content/develop/diffs.md`
- Client-side fields — `docs/content/develop/client-side-fields.md`
- Promote beta → GA — `docs/content/develop/promote-to-ga.md`
- Generate the providers / set up the environment — `docs/content/develop/generate-providers.md`, `docs/content/develop/set-up-dev-environment.md`

**Testing:**

- Write acceptance tests (create, update, import) — `docs/content/test/test.md`
- Run tests — `docs/content/test/run-tests.md`

**Contributing (PRs, review, release notes):**

- Contribution process end to end — `docs/content/contribution-process.md`
- Create a PR — `docs/content/code-review/create-pr.md`
- Write release notes — `docs/content/code-review/release-notes.md`
- Review a PR — `docs/content/code-review/review-pr.md`

**Documentation:**

- Add resource documentation — `docs/content/document/add-documentation.md`
- Handwritten docs style guide — `docs/content/document/handwritten-docs-style-guide.md`

**Best practices (judgment calls with team positions):**

- Immutable fields / ForceNew — `docs/content/best-practices/immutable-fields.md`
- Deletion policy — `docs/content/best-practices/deletion-policy.md`
- Labels and annotations — `docs/content/best-practices/labels-and-annotations.md`
- Client-side validation — `docs/content/best-practices/validation.md`
- Common resource patterns (singletons) — `docs/content/best-practices/common-resource-patterns.md`

## Agent-only entries

- **enums-vs-strings** — Model an API enum as Enum (strict, plan-time) or String (forward-compatible): the deliberate trade-off. — [field/enums-vs-strings.md](field/enums-vs-strings.md)
