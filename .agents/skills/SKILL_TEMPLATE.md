---
name: template-skill
description: "A single-sentence description of what this skill does and when an agent should invoke it."
---

# `template-skill`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Prerequisites
List any environment variables, workspace parity checks, or directory contexts required before executing this skill.
* Example: "You must be in the `magic-modules` root directory."

## Execution Steps

### 1. Verification
Provide the exact bash commands the agent should run to verify the prerequisites. Do NOT let the agent proceed if this fails.

#### Example: Verify Directory
```bash
pwd # Verify we are in /Users/camthornton/magic-modules
```

### 2. The Core Commands
Provide the precise, bulletproof bash commands, python script, or instruction set the agent needs to execute. This replaces the agent having to guess CLI flags or complex arguments.

#### Example: Compile Downstream Provider
```bash
# Make terraform provider for Beta
make terraform VERSION=beta OUTPUT_PATH=$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta PRODUCT=compute
```

### 3. Verification & Handoff
Instructions on how the agent should verify the command succeeded, and what workflow or rule it should return to next.

* Example: "If the `make testacc` command fails, you MUST immediately invoke the `parse-debug-logs` skill. Do not attempt a blind fix."
