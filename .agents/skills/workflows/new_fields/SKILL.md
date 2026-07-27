---
name: new-fields-workflow
description: "Workflow specifically for adding new fields to an existing resource."
---

# `new-fields-workflow`

This document guides the agent through the task of adding new fields to an existing Terraform resource in Magic Modules.

## Execution Steps

### 1. Implement change

- Consult `.agents/knowledge/index.md` for the topics the resource touches and open the relevant sources.
- Follow the standard process to add the fields to the resource's YAML definition in `mmv1/products/...` based on API documentation and repository patterns.
- Follow the standard process to add or update tests for the resource to cover the new fields.

### 2. Run pre-gen tests

- Use [run-pre-gen-checks](.agents/skills/utils/run-pre-gen-checks/SKILL.md). If issues are found, analyze and fix them.
- Repeat this step until all issues are resolved. Don't proceed until the run-pre-gen-checks skill reports complete success. If you're stuck, ask the user for assistance.

### 3. Build and test downstreams

- Use [build-and-test-downstreams](.agents/skills/utils/build-and-test-downstreams/SKILL.md). If issues are found, analyze and fix them.
- Repeat this step until all issues are resolved. Don't proceed until the build-and-test-downstream skill reports complete success. If you're stuck, ask the user for assistance.

### 4. Run acceptance tests

- Use [`repo-sync`](`.agents/skills/operations/repo-sync/`) to ensure the downstream repositories are in sync with magic-modules.
- Use [`generate-provider`](.agents/skills/operations/generate-provider/`) to generate the provider code into the downstream repositories.
- Invoke the specialized `qa-test-runner` subagent using the `invoke_subagent` tool to run all acceptance tests for the modified resource. The subagent evaluates if the test fails/passes and returns a human-readable interpretation of the results.
- If failures fall within known best practices / troubleshooting procedures, and the tests ran in less than 5 minutes, report the results to the user but go ahead and try to fix the failures. Otherwise, explain the situation and ask for user input.
- Repeat this step until all issues are resolved. Don't proceed until the `qa-test-runner` reports complete success. If you're stuck, ask for user assistance.
