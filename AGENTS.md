# Agent Instructions

Magic Modules (MMv1) defines the Google Cloud Terraform providers: resource schemas in `mmv1/products/`,
handwritten code and tests under `mmv1/third_party/terraform/`, generated into the downstream
`terraform-provider-google` and `terraform-provider-google-beta` repositories. Changes are always made
here, upstream — never in the downstream repos.

## Entry points

- **[.agents/WORKFLOWS.md](.agents/WORKFLOWS.md)** — the workflows for provider tasks (adding
  resources/fields, fixing issues, syncing) and the rules that govern them. Read it before starting a
  provider task.
- **`.agents/skills/`** — reusable skills the workflows compose (generation, testing, log parsing,
  troubleshooting). Each skill's `SKILL.md` frontmatter states when to use it.
- **`.agents/knowledge/`** — the curated knowledge base (initial seeding in progress). When present,
  consult its `index.md` at decision points and open only the entries the task needs.
- **`.agents/archive/`** — parked tracks (currently TGC). Not maintained; do not use as reference.

## Ground rules

- **Never weaken a test or check to make it pass.** No disabling or skipping tests, and no test-dodging
  behavior flags (`ignore_read`, `default_from_api`, `ImportStateVerifyIgnore`) without an adjacent
  comment justifying the API behavior that requires them.
- **Verify before opening a PR**: generate, build, and run the tests relevant to what you changed.
